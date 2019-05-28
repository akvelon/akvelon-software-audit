package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const testDownloadDir = "test_downloads"

func TestRepoRootForImportPath(t *testing.T) {
	cases := []struct {
		giveURL  string
		wantPath string
	}{
		{"github.com/akvelon/PowerBI-Stacked-Column-Chart", "github.com/akvelon/PowerBI-Stacked-Column-Chart"},
	}

	for _, tt := range cases {
		repo := Repository{tt.giveURL}

		fullPath, err := repo.Download(testDownloadDir)
		if err != nil {
			t.Fatalf("Error calling Download(%q): %v", tt.giveURL, err)
		}

		fmt.Printf("fullPath at: %s\n", fullPath)
		fmt.Printf("wantPath at: %s\n", filepath.Join(testDownloadDir, tt.wantPath))

		if fullPath != filepath.Join(testDownloadDir, tt.wantPath) {
			t.Errorf("Download(%q): root.Repo = %q, want %q", tt.giveURL, fullPath, tt.wantPath)
		}

		wantPath := filepath.Join(testDownloadDir, tt.wantPath)
		ex, _ := exists(wantPath)
		if !ex {
			t.Errorf("Download(%q): %q was not created", tt.giveURL, wantPath)
		}
	}

	// clean up the test
	os.RemoveAll(testDownloadDir)
}
