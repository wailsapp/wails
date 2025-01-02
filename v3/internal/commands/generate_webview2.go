package commands

import (
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/term"
	"os"
	"path/filepath"
)

//go:embed webview2/MicrosoftEdgeWebview2Setup.exe
var webview2Bootstrapper []byte

type GenerateWebView2Options struct {
	Directory string `json:"directory"`
}

func GenerateWebView2Bootstrapper(options *GenerateWebView2Options) error {
	// If the file already exists, exit early
	if _, err := os.Stat(filepath.Join(options.Directory, "MicrosoftEdgeWebview2Setup.exe")); err == nil {
		return nil
	}

	// Create target directory if it doesn't exist
	err := os.MkdirAll(options.Directory, 0755)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write to target directory
	targetPath := filepath.Join(options.Directory, "MicrosoftEdgeWebview2Setup.exe")
	err = os.WriteFile(targetPath, webview2Bootstrapper, 0644)
	if err != nil {
		return fmt.Errorf("failed to write WebView2 bootstrapper: %w", err)
	}

	term.Success("Generated WebView2 bootstrapper at: " + targetPath + "\n")
	return nil
}
