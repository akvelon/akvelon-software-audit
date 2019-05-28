package main

import (
	_ "akvelon/akvelon-software-audit/routers"
	"log"
	"os"

	"github.com/astaxie/beego"
)

func main() {
	if err := os.MkdirAll("_repos", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}
	beego.Run()
}
