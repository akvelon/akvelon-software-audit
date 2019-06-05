package analyzer

import (
	"akvelon/akvelon-software-audit/internals/analyzer/boyterlc"
)

// Service provides analyze operations.
type Service interface {
	Run() ([]ScanResult, error)
}

type service struct {
	sources string
}

// ScanResult is a combined result of repo analysis.
type ScanResult struct {
	File string `json:"File"`
	License string `json:"License"`
	Confidence string `json:"Confidence"`
	Size string `json:"Size"`
}

// NewService creates an analize service with the necessary dependencies.
func NewService(path string) Service {
	return &service{sources: path}
}

func (s *service) Run() ([]ScanResult, error) {
	// Let's omit DI pattern for various analyzers here for simplicity
	results, err := boyterlc.Scan(s.sources) 
	if err != nil {
		return nil, err
	}
	var res []ScanResult
	for _, result := range results {
		res = append(res, ScanResult {
			File: result.File,
			License: result.License,
			Confidence: result.Confidence,
			Size: result.Size,
		})
	}
	return res, nil	
}
