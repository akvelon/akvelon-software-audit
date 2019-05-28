package main

import (
	"akvelon/akvelon-software-audit/internals/vcs"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	dir = flag.String("repo", "", "URL of remote repository")
)

// CLIDest is a predefined folder to download repos
const CLIDest = "CLIRepos"

func main() {
	flag.Parse()
	repo := vcs.Repository{URL: *dir}

	fullLocalPath, err := repo.Download(CLIDest)
	if err != nil {
		log.Fatalf("Fatal error downloading %s: %s", *dir, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Downloaded at: %s\n", fullLocalPath)
}
