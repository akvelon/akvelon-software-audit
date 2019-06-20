package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

var (
	auditSrv             = beego.AppConfig.String("auditservice")
	getHealthURL         = fmt.Sprintf("%s/health", auditSrv)
	getRecentlyViewedURL = fmt.Sprintf("%s/recent", auditSrv)
	getAnalizeByRepoURL  = fmt.Sprintf("%s/analize", auditSrv)
	postAnalizeByRepoURL = fmt.Sprintf("%s/analize", auditSrv)
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

	body := strings.NewReader(fmt.Sprintf(`url=%s`, repoLink))
	req, err := http.NewRequest("POST", postAnalizeByRepoURL, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		flash.Error("Failed to analyze the repository.")
		flash.Store(&this.Controller)
		this.Redirect("/", 302)
		return
	}

	defer resp.Body.Close()

	flash.Success("Thanks, results will be ready soon...")
	flash.Store(&this.Controller)
	this.Redirect("/", 302)
}
