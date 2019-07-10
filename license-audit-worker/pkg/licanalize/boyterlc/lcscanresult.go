package boyterlc

// LCScanResult shows meaningful results of license scan.
type LCScanResult struct {
	File       string
	License    string
	Confidence string
	Size       string
}