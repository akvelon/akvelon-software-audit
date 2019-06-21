package main

import (
	"akvelon/akvelon-software-audit/license-audit-service/pkg/http/rest"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/storage/bolt"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

const (
	downloadRepoPath = "_repos/src/github.com"
)

var (
	addr   = flag.String("http", ":8000", "HTTP listen address")
	router *httprouter.Router
)

func main() {
	flag.Parse()
	if err := os.MkdirAll(downloadRepoPath, 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	s := new(bolt.Storage)
	s.InitStorage()
	licAnalizer := licanalize.NewService(s)

	router := rest.Handler(licAnalizer)

	log.Printf("The license-audit-service is running on: http://localhost:%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}
