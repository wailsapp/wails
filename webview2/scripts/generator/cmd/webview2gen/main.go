package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"updater/generator"
)

var (
	version string
)

func init() {
	flag.StringVar(&version, "version", "", "WebView2 version to use (e.g., 1.0.2903.40)")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: webview2gen <command> [options]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  full      Generate all bindings from IDL")
		fmt.Fprintln(os.Stderr, "  verify    Verify committed output matches generated output")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	command := flag.Args()[0]

	switch command {
	case "full":
		runFull()
	case "verify":
		runVerify()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runFull() {
	idlVersion := version
	if idlVersion == "" {
		data, err := os.ReadFile("scripts/latest_version.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read latest version: %v\n", err)
			os.Exit(1)
		}
		idlVersion = string(data)
	}

	idlFile := filepath.Join("scripts", fmt.Sprintf("WebView2.%s.idl", idlVersion))
	idlData, err := os.ReadFile(idlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read IDL file %s: %v\n", idlFile, err)
		os.Exit(1)
	}

	files, err := generator.ParseIDL(idlData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse IDL: %v\n", err)
		os.Exit(1)
	}

	outputDir := "pkg/webview2"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		filePath := filepath.Join(outputDir, file.FileName)
		err := os.WriteFile(filePath, file.Content.Bytes(), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s\n", filePath)
	}

	fmt.Printf("\nSuccessfully generated %d files\n", len(files))
}

func runVerify() {
	idlVersion := version
	if idlVersion == "" {
		data, err := os.ReadFile("scripts/latest_version.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read latest version: %v\n", err)
			os.Exit(1)
		}
		idlVersion = string(data)
	}

	idlFile := filepath.Join("scripts", fmt.Sprintf("WebView2.%s.idl", idlVersion))
	idlData, err := os.ReadFile(idlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read IDL file %s: %v\n", idlFile, err)
		os.Exit(1)
	}

	files, err := generator.ParseIDL(idlData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse IDL: %v\n", err)
		os.Exit(1)
	}

	hasDifferences := false
	outputDir := "pkg/webview2"

	for _, file := range files {
		filePath := filepath.Join(outputDir, file.FileName)
		existingData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read existing file %s: %v\n", filePath, err)
			os.Exit(1)
		}

		generatedData := file.Content.Bytes()
		if string(existingData) != string(generatedData) {
			fmt.Printf("DIFFERENCE FOUND in %s\n", filePath)
			hasDifferences = true
		}
	}

	if hasDifferences {
		fmt.Fprintln(os.Stderr, "\nERROR: Generated output does not match committed files.")
		fmt.Fprintln(os.Stderr, "Run 'webview2gen full' to regenerate the bindings.")
		os.Exit(1)
	}

	fmt.Printf("All %d files match generated output\n", len(files))
}
