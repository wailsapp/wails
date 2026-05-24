package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("npx", "all-contributors-cli", "check")
	//cmd.Stdin = strings.NewReader("some input")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	missingSplit := strings.Split(out.String(), "\n")
	if len(missingSplit) < 2 {
		log.Fatal(out.String())
	}
	missing := missingSplit[1]
	missing = strings.TrimSpace(missing)
	// Split on comma
	for _, contrib := range strings.Split(missing, ",") {
		// Trim whitespace
		contrib = strings.TrimSpace(contrib)
		if contrib == "dependabot[bot]" || contrib == "" {
			continue
		}
		// Add contributor
		cmd := exec.Command("npx", "all-contributors-cli", "add", contrib, "code")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}
