package analizer

import (
	"fmt"
	"log"
	"os/exec"
)

const outJSON = "licInfo.json"
const outFormat = "json"

// Scan license scan againt https://github.com/boyter/lc tool.
func Scan(path string) {
	fmt.Printf("Start executing command for repo analize at path: %s\n", path)
	cmd := exec.Command("lc", "-f", outFormat, "-o", outJSON, path)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("Finish running the command at path: %s\n", path)
}
