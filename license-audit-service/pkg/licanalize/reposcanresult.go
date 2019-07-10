package licanalize

// RepoScanResult represents structure for scan results.
type RepoScanResult struct {
	File       string `json:"file" bson:"file"`
	License    string `json:"license" bson:"license"`
	Confidence string `json:"confidence" bson:"confidence"`
	Size       string `json:"size" bson:"size"`
}
