//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func main() {
	output := flag.String("output", "../docs/public/schemas/wails.v3.json", "generated JSON Schema output path")
	flag.Parse()
	data, err := manifest.JSONSchema()
	if err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(*output, data, 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
