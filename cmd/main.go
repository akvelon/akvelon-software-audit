package main

import (
	"akvelon/akvelon-software-audit/internals/analyzer"
	"akvelon/akvelon-software-audit/internals/vcs"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	repo      = flag.String("repo", "", "URL of remote repository")
	reposDest = filepath.Join("..", "_repos")
)

func main() {
	flag.Parse()
	if *repo == "" {
		log.Println("Usage: ./cmd -repo `github-repo-url`")
		return
	}

	repository := vcs.NewRepository(*repo)
	fullPath, err := repository.Download(reposDest)
	fmt.Printf("Downloaded at: %s\n", fullPath)

	if err != nil {
		log.Fatalf("Fatal error downloading %s: %s", *repo, err.Error())
		os.Exit(1)
	}

	// TODO: decide where to store results, e.g. MongoDB? BoltDB?
	//storage := mongoDB.NewStorage()
	analyzer := analyzer.NewService(fullPath)

	fmt.Printf("Starting license analize at %s\n for %s\n", time.Now().Format(time.RFC850), fullPath)
	res, analyzerErr := analyzer.Run()

	if analyzerErr != nil {
		log.Fatalf("Fatal error analizing repo %s: %s", fullPath, analyzerErr.Error())
	}

	for _, item := range res {
		fmt.Printf("FileName: %s\n", item.File)
		fmt.Printf("License: %s\n", item.License)
		fmt.Printf("Confidence: %s\n", item.Confidence)
		fmt.Printf("Size: %s\n\n\n", item.Size)
	}
}
