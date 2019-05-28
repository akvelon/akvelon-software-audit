package main

import (
	"akvelon/akvelon-software-audit/internals/vcs"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	repo = flag.String("repo", "", "URL of remote repository")
)

var reposDest = filepath.Join("..", "_repos")

func main() {
	flag.Parse()
	if *repo == "" {
		log.Println("Usage: ./cmd -repo `github-repo-url`")
		return
	}

	repository := vcs.Repository{URL: *repo}
	fullLocalPath, err := repository.Download(reposDest)

	if err != nil {
		log.Fatalf("Fatal error downloading %s: %s", *repo, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Downloaded at: %s\n", fullLocalPath)
}
