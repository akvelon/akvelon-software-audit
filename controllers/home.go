package controllers

import (
	"akvelon/akvelon-software-audit/internals/analyzer"
	"akvelon/akvelon-software-audit/internals/vcs"
	"fmt"
	"path/filepath"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	beego.ReadFromRequest(&this.Controller)

	this.Layout = "layout_main.tpl"
	this.LayoutSections = make(map[string]string)

	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}

func (this *MainController) Report() {
	provider := this.Ctx.Input.Param(":provider")
    if provider == "" {
        this.Ctx.WriteString("provider is empty")
        return
	}
	
	orgname := this.Ctx.Input.Param(":orgname")
    if orgname == "" {
        this.Ctx.WriteString("orgname is empty")
        return
	}
	
	reponame := this.Ctx.Input.Param(":reponame")
    if reponame == "" {
        this.Ctx.WriteString("reponame is empty")
        return
    }

	this.Data["provider"] = provider
	this.Data["orgname"] = orgname
	this.Data["reponame"] = reponame

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

	repo := vcs.NewRepository(repoLink)

	_, err := doAnalyze(repo)
	if err != nil {
		flash.Error("Couldn't analyze the repository: " + err.Error())
		flash.Store(&this.Controller)
		this.Redirect("/", 302)
		return
	}

	// flash.Success("Thanks, repository submitted for analyze.")
	// flash.Store(&this.Controller)
	this.Ctx.Redirect(302, fmt.Sprintf("/report/%v", repoLink))
}

func doAnalyze(repo *vcs.Repository) (analyzer.ScanResult, error) {
	// fetch repo for further analyzis
	var reposDest = filepath.Join(".", "_repos")
	_, err := repo.Download(reposDest)
	if err != nil {
		return analyzer.ScanResult{}, fmt.Errorf("Failed do download repository: %v", err)
	}
	return analyzer.ScanResult{}, nil
}
