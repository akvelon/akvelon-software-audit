package boyterlc

type FileResult struct {
	Directory         string         `json:"Directory"`
	Filename          string         `json:"Filename"`
	LicenseGuesses    []LicenseMatch `json:"LicenseGuesses"`
	LicenseRoots      []LicenseMatch `json:"LicenseRoots"`
	LicenseIdentified []LicenseMatch `json:"LicenseIdentified"`
	Md5Hash           string         `json:"Md5Hash"`
	Sha1Hash          string         `json:"Sha1Hash"`
	Sha256Hash        string         `json:"Sha256Hash"`
	BytesHuman        string         `json:"BytesHuman"`
	Bytes             int            `json:"Bytes"`
}