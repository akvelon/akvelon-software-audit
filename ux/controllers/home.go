package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/streadway/amqp"
)

var (
	auditSrv  = beego.AppConfig.String("auditservice")
	rabbitSrv = beego.AppConfig.String("rabbit")

	getRecentlyViewedURL = fmt.Sprintf("%s/recent", auditSrv)
	getAnalizeByRepoURL  = fmt.Sprintf("%s/analize", auditSrv)

	uXAuditQueueName = "audit-queue"
)

type MainController struct {
	beego.Controller
}

type RepoScanResult struct {
	File       string
	License    string
	Confidence string
	Size       string
}

func (this *MainController) Get() {
	beego.ReadFromRequest(&this.Controller)
	req := httplib.Get(getRecentlyViewedURL)
	var recent []string
	rec, err := req.String()
	if err != nil {
		fmt.Printf("failed to get results from audit service: %v", err)
	}

	dec := json.NewDecoder(strings.NewReader(rec))
	err = dec.Decode(&recent)
	if err != nil {
		fmt.Printf("failed to parse results from audit service: %v", err)
	}

	this.Data["Recent"] = recent

	this.Layout = "layout_main.tpl"
	this.LayoutSections = make(map[string]string)

	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}

func (this *MainController) Report() {
	provider := this.Ctx.Input.Param(":provider")
	orgname := this.Ctx.Input.Param(":orgname")
	reponame := this.Ctx.Input.Param(":reponame")

	if provider == "" || orgname == "" || reponame == "" {
		this.Ctx.WriteString("Sorry, invalid query string parameter.")
		return
	}

	repoURL := fmt.Sprintf("%s/%s/%s", provider, orgname, reponame)
	this.Data["repoURL"] = repoURL

	req := httplib.Get(fmt.Sprintf("%s?url=%s", getAnalizeByRepoURL, repoURL))
	var result []RepoScanResult
	r, err := req.String()
	if err != nil {
		fmt.Printf("failed to get results from audit service: %v", err)
	}

	dec := json.NewDecoder(strings.NewReader(r))
	err = dec.Decode(&result)
	if err != nil {
		fmt.Printf("failed to parse results from audit service: %v", err)
	}

	this.Data["analyzeResult"] = result

	this.Layout = "layout_main.tpl"
	this.LayoutSections = make(map[string]string)

	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}

func (this *MainController) Analyze() {
	repoLink := this.GetString("repo")
	flash := beego.NewFlash()
	if repoLink == "" {
		flash.Error("Couldn't analyze the repository, empty string provided.")
		flash.Store(&this.Controller)
		this.Redirect("/", 302)
		return
	}

	conn, err := amqp.Dial(rabbitSrv)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a rmq channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		uXAuditQueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a rabbit queue")

	body := repoLink
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	flash.Success("Thanks, results are submitted and will be ready soon...")
	flash.Store(&this.Controller)
	this.Redirect("/", 302)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
