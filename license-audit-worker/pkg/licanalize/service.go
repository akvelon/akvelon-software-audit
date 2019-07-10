package licanalize

import (
	"fmt"
	"log"

	"akvelon/akvelon-software-audit/license-audit-service/pkg/download/vcs"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize/boyterlc"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/storage/mongo"
)

// Service provides license analize operations.
type Service interface {
	Scan(repo AnalizedRepo) error
}

type service struct {
	r Repository
}

// Repository provides access to db.
type Repository interface {
	InitStorage() error
	UpdateRecentlyViewed(repo string) error
	SaveRepoToDB(repo string, data []mongo.RepoScanResult) error
}

// NewService creates new service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// Scan scans given repository for license
func (s *service) Scan(repo AnalizedRepo) error {
	// Download repo before we can scan
	log.Printf("Start downloading repo %s \n", repo.URL)
	repoRoot, err := vcs.Download(repo.URL, "_repos/src")
	if err != nil {
		return fmt.Errorf("Failed do download repository: %v", err)
	}

	results, err := boyterlc.Scan(repoRoot)
	if err != nil {
		return err
	}

	// convert to proper type TODO: move to converter package?
	var recentRepos = make([]mongo.RepoScanResult, len(results))
	var j = len(results) - 1
	for _, r := range results {
		recentRepos[j] = mongo.RepoScanResult{
			File:       r.File,
			License:    r.License,
			Confidence: r.Confidence,
			Size:       r.Size,
		}
		j--
	}

	err = s.r.SaveRepoToDB(repo.URL, recentRepos)
	if err != nil {
		return fmt.Errorf("failed to save results to db: %v", err)
	}

	err = s.r.UpdateRecentlyViewed(repo.URL)
	if err != nil {
		return fmt.Errorf("failed to update recently viewed to db: %v", err)
	}

	return nil
}
