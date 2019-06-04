package controllers

import (
	"akvelon/akvelon-software-audit/internals/analizer"
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

	flash.Success("Thanks, repository submitted for analyze.")
	flash.Store(&this.Controller)
	this.Ctx.Redirect(302, "/")
}

func doAnalyze(repo *vcs.Repository) (analizer.Result, error) {
	// fetch repo for further analyzis
	var reposDest = filepath.Join(".", "_repos")
	_, err := repo.Download(reposDest)
	if err != nil {
		return analizer.Result{}, fmt.Errorf("Failed do download repository: %v", err)
	}
	return analizer.Result{}, nil
}
