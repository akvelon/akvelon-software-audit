package licanalize

import (
	"fmt"

	"akvelon/akvelon-software-audit/license-audit-service/pkg/storage/mongo"
)

// Service provides license analize operations.
type Service interface {
	CheckHealth() bool
	GetRecent() ([]string, error)
	GetRepoResultFromDB(repo string) ([]mongo.RepoScanResult, error)
}

type service struct {
	r Repository
}

// Repository provides access to db.
type Repository interface {
	InitStorage() error
	GetRecentlyViewed() ([]string, error)
	GetRepoFromDB(repo string) ([]mongo.RepoScanResult, error)
}

// NewService creates new service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// GetRepoResultFromDB returns scan result from DB
func (s *service) GetRepoResultFromDB(repo string) ([]mongo.RepoScanResult, error) {
	results, err := s.r.GetRepoFromDB(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan result from DB")
	}

	return results, nil
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
