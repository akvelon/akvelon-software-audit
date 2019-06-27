package main

import (
	"akvelon/akvelon-software-audit/license-audit-service/pkg/http/rest"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/storage/bolt"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/monitor"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/tracing"

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

	t, _ := tracing.InitTracer(os.Getenv("JAEGER_SERVICE_NAME"))

	s := new(bolt.Storage)
	s.InitStorage()
	la := licanalize.NewService(s)
	m := &monitor.Monitor{}

	router := rest.Handler(la, m, t)

	log.Printf("The license-audit-service is running on: http://localhost:%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}
