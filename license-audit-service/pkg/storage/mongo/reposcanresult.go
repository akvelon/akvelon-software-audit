package mongo

// RepoScanResult represents single result of License Scan
type RepoScanResult struct {
	File       string `json:"file" bson:"file"`
	License    string `json:"license" bson:"license"`
	Confidence string `json:"confidence" bson:"confidence"`
	Size       string `json:"size" bson:"size"`
}

// RepoScanItem represents all results for given repository
type RepoScanItem struct {
	Repo string `json:"repo" bson:"repo"`
	Results []RepoScanResult `json:"reposcanresults" bson:"reposcanresults"`
}