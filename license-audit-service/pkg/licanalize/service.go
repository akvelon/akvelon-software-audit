package licanalize

import (
	"encoding/json"
	"fmt"
	"log"

	"akvelon/akvelon-software-audit/license-audit-service/pkg/download/vcs"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize/boyterlc"
)

// Service provides license analize operations.
type Service interface {
	CheckHealth() bool
	Scan(repo AnalizedRepo) error
	GetRecent() ([]string, error)
	GetRepoResultFromDB(repo string) ([]RepoScanResult, error)
}

type service struct {
	r Repository
}

// Repository provides access to db.
type Repository interface {
	InitStorage() error
	GetRecentlyViewed() ([]string, error)
	UpdateRecentlyViewed(repo string) error
	SaveRepoToDB(key string, data []byte) error
	GetRepoFromDB(repo string) ([]byte, error)
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

	var res []boyterlc.LCScanResult
	for _, result := range results {
		res = append(res, boyterlc.LCScanResult{
			File:       result.File,
			License:    result.License,
			Confidence: result.Confidence,
			Size:       result.Size,
		})
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not marshal json: %v", err)
	}

	err = s.r.SaveRepoToDB(repo.URL, resBytes)
	if err != nil {
		return fmt.Errorf("failed to save results to db: %v", err)
	}

	err = s.r.UpdateRecentlyViewed(repo.URL)
	if err != nil {
		return fmt.Errorf("failed to update recently viewed to db: %v", err)
	}

	return nil
}

// GetRepoResultFromDB returns scan result from DB
func (s *service) GetRepoResultFromDB(repo string) ([]RepoScanResult, error) {
	b, err := s.r.GetRepoFromDB(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan result from DB")
	}

	resp := []RepoScanResult{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON for %q in result", repo)
	}
	return resp, nil
}

// GetRecent returns top recent repos scanned
func (s *service) GetRecent() ([]string, error) {
	recent, err := s.r.GetRecentlyViewed()
	if err != nil {
		return nil, fmt.Errorf("Failed to get recent data from DB")
	}
	return recent, nil
}

// Scan scans given repository for license
func (s *service) CheckHealth() bool {
	// TODO: check DB is avaliable
	return true
}
