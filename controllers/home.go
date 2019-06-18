package controllers

import (
	"akvelon/akvelon-software-audit/internals/analyzer"
	"akvelon/akvelon-software-audit/internals/storage/bolt"
	"akvelon/akvelon-software-audit/internals/vcs"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	beego.ReadFromRequest(&this.Controller)

	recent, errDb := bolt.GetRecentlyViewed()
	if errDb != nil {
		fmt.Printf("failed to get results from db: %v", errDb)
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
	repoResult, _ := bolt.GetRepoFromDB(repoURL)

	this.Data["analyzeResult"] = repoResult

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

	err := doAnalyze(repo)
	if err != nil {
		flash.Error("Couldn't analyze the repository: " + err.Error())
		flash.Store(&this.Controller)
		this.Redirect("/", 302)
		return
	}

	this.Ctx.Redirect(302, fmt.Sprintf("/report/%v", repoLink))
}

func doAnalyze(repo *vcs.Repository) error {
	var reposDest = filepath.Join(".", "_repos")
	_, err := repo.Download(reposDest)
	if err != nil {
		return fmt.Errorf("Failed do download repository: %v", err)
	}
	analyzer := analyzer.NewService(reposDest)
	res, analyzerErr := analyzer.Run()

	if analyzerErr != nil {
		return fmt.Errorf("Fatal error analizing repo %s: %s", reposDest, analyzerErr.Error())
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not marshal json: %v", err)
	}

	errDb := bolt.SaveRepoToDB(repo.URL, resBytes)
	if errDb != nil {
		return fmt.Errorf("failed to save results to db: %v", errDb)
	}

	errDb = bolt.UpdateRecentlyViewed(repo.URL)
	if errDb != nil {
		return fmt.Errorf("failed to update recently viewed to db: %v", errDb)
	}
	return nil
}
