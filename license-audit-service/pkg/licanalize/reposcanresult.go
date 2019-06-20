package licanalize

type RepoScanResult struct {
	File       string `json:"File"`
	License    string `json:"License"`
	Confidence string `json:"Confidence"`
	Size       string `json:"Size"`
}
