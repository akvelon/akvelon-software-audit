package boyterlc

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const outFormat = "json"

type LicenseMatch struct {
	LicenseId string
	Score     float64
}

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

// LCScanResult shows meaningful results of license scan.
type LCScanResult struct {
	File string
	License string
	Confidence string
	Size string
}

// Scan license with https://github.com/boyter/lc tool.
func Scan(path string) ([]LCScanResult, error) {
	fmt.Printf("Start lc at path: %s\n", path)
	cmd := exec.Command("lc", "-f", outFormat, path)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return nil, err
	}

	var res []FileResult
	jsonErr := json.Unmarshal(out, &res)
	if jsonErr != nil {
		log.Printf("Failed to parse output json %s\n", jsonErr)
		return nil, jsonErr
	}

	var output []LCScanResult
	for _, item := range res {

		licenseConcluded, confidence := determineLicense(item)
		output = append(output, LCScanResult {
			File: item.Filename,
			License: licenseConcluded,
			Confidence: confidence,
			Size: item.BytesHuman,
		})
	}
	fmt.Printf("Finish running the command at path: %s\n", path)
	return output, nil
}

func determineLicense(result FileResult) (string, string) {
	license := ""
	confidence := 100.00
	var licenseMatches []LicenseMatch

	if len(result.LicenseIdentified) != 0 {
		license = joinLicenseList(result.LicenseIdentified, result.LicenseRoots, " AND ")
		confidence = 100.00
	} else if len(result.LicenseGuesses) != 0 {
		license = result.LicenseGuesses[0].LicenseId
		confidence = result.LicenseGuesses[0].Score
		licenseMatches = append(licenseMatches, result.LicenseGuesses[0])
	}

	rootLicenses := joinLicenseList(result.LicenseRoots, licenseMatches, " OR ")
	if rootLicenses != "" {
		if license == "" {
			license = rootLicenses
		} else {
			license = rootLicenses + " AND " + license
		}
	}

	if license == "" {
		license = "NOASSERTION"
	}

	return license, fmt.Sprintf("%.2f%%", confidence)
}

func joinLicenseList(licenseList []LicenseMatch, ignore []LicenseMatch, operator string) string {
	licenseDeclared := ""

	if len(licenseList) == 1 {
		if licenceListHasLicense(licenseList[0], ignore) == false {
			licenseDeclared = licenseList[0].LicenseId
		}
	} else if len(licenseList) >= 2 {
		var licenseNames []string
		for _, v := range licenseList {
			if licenceListHasLicense(v, ignore) == false {
				licenseNames = append(licenseNames, v.LicenseId)
			}
		}

		if len(licenseNames) == 1 {
			licenseDeclared = licenseNames[0]
		} else if len(licenseNames) != 0 {

			licenseDeclared = strings.Join(licenseNames, operator)

			if operator == " OR " {
				licenseDeclared = "(" + licenseDeclared + ")"
			}
		}
	}

	return licenseDeclared
}

// Returns true if a license list contains the license
func licenceListHasLicense(license LicenseMatch, licenseList []LicenseMatch) bool {
	for _, v := range licenseList {
		if v.LicenseId == license.LicenseId {
			return true
		}
	}

	return false
}