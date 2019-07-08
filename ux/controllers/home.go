package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"strconv"
	"net/http"
	"akvelon/akvelon-software-audit/ux/monitor"
	"akvelon/akvelon-software-audit/ux/lib/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/astaxie/beego"
)

var (
	auditSrv  = beego.AppConfig.String("auditservice")

	recentlyViewedURL = fmt.Sprintf("%s/recent", auditSrv)
	analizeByRepoURL  = fmt.Sprintf("%s/analize", auditSrv)
)

type MainController struct {
	beego.Controller
	Tracer opentracing.Tracer
}

type RepoScanResult struct {
	File       string
	License    string
	Confidence string
	Size       string
}

func (this *MainController) Get() {
	monitor.HttpRequestsTotal.Inc()
	span := this.Tracer.StartSpan("Get-MainController")
	defer span.Finish()

	beego.ReadFromRequest(&this.Controller)

	req, err := http.NewRequest("GET", recentlyViewedURL, nil)

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, recentlyViewedURL)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	resp, err := xhttp.Do(req)

	if err != nil {
		span.LogKV("getting-recent-scans", "failed to get recent results from audit service")
		fmt.Printf("failed to get results from audit service: %v", err)
	}

	var recent []string
	rec := string(resp)

	dec := json.NewDecoder(strings.NewReader(rec))
	err = dec.Decode(&recent)
	if err != nil {
		span.LogKV("getting-recent-scans", "failed to parse results from audit service")
		fmt.Printf("failed to parse results from audit service: %v", err)
	}
	
	this.Data["Recent"] = recent
	
	this.Layout = "layout_main.tpl"
	this.LayoutSections = make(map[string]string)

	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}

func (this *MainController) Report() {
	monitor.HttpRequestsTotal.Inc()
	span := this.Tracer.StartSpan("Report-MainController")
	defer span.Finish()

	provider := this.Ctx.Input.Param(":provider")
	orgname := this.Ctx.Input.Param(":orgname")
	reponame := this.Ctx.Input.Param(":reponame")

	if provider == "" || orgname == "" || reponame == "" {
		this.Ctx.WriteString("Sorry, invalid query string parameter.")
		span.LogKV("getting-report-results", "invalid query string parameter")
		return
	}

	repoURL := fmt.Sprintf("%s/%s/%s", provider, orgname, reponame)
	span.SetTag("analized-repo", repoURL)
	this.Data["repoURL"] = repoURL

	url := fmt.Sprintf("%s?url=%s", analizeByRepoURL, repoURL)
	req, err := http.NewRequest("GET", url, nil)

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	resp, err := xhttp.Do(req)

	if err != nil {
		span.LogKV("getting-report-results", "failed to get results from audit service")
		fmt.Printf("failed to get results from audit service: %v", err)
	}

	r := string(resp)

	var result []RepoScanResult

	dec := json.NewDecoder(strings.NewReader(r))
	err = dec.Decode(&result)
	if err != nil {
		span.LogKV("getting-report-results", "failed to parse results from audit service")
		fmt.Printf("failed to parse results from audit service: %v", err)
	}

	this.Data["analyzeResult"] = result

	this.Layout = "layout_main.tpl"
	this.LayoutSections = make(map[string]string)

	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}

func (this *MainController) Analyze() {
	monitor.HttpRequestsTotal.Inc()
	span := this.Tracer.StartSpan("Analyze-MainController")
	defer span.Finish()

	repoLink := this.GetString("repo")
	flash := beego.NewFlash()
	if repoLink == "" {
		flash.Error("Couldn't analyze the repository, empty string provided.")
		flash.Store(&this.Controller)
		span.LogKV("analize-repo", "couldn't analyze the repository, empty string provided")
		this.Redirect("/", 302)
		return
	}

	span.SetTag("analized-repo", repoLink)

	data := url.Values{}
	data.Set("url", repoLink)

	req, err := http.NewRequest("POST", analizeByRepoURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, analizeByRepoURL)
	ext.HTTPMethod.Set(span, "POST")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	_, err = xhttp.Do(req)

	if err != nil {
		span.LogKV("analized-repo", "failed to execute analize repo from audit service")
		fmt.Printf("failed to execute analize repo from audit service: %v", err)
	}

	flash.Success("Thanks, results are submitted and will be ready soon...")
	flash.Store(&this.Controller)
	this.Redirect("/", 302)
}
