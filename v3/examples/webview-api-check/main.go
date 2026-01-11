package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

// PlatformInfo contains information about the platform and webview
type PlatformInfo struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	GoVersion   string `json:"goVersion"`
	WailsInfo   string `json:"wailsInfo"`
	WebViewInfo string `json:"webViewInfo"`
	GTKVersion  string `json:"gtkVersion,omitempty"`
	Timestamp   string `json:"timestamp"`
}

// APIReport represents the full API compatibility report
type APIReport struct {
	Platform PlatformInfo           `json:"platform"`
	APIs     map[string]interface{} `json:"apis"`
}

// APICheckService provides methods for the frontend
type APICheckService struct{}

// GetPlatformInfo returns information about the current platform
func (s *APICheckService) GetPlatformInfo() PlatformInfo {
	info := PlatformInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
		WailsInfo: "v3.0.0-dev",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Platform-specific webview info
	switch runtime.GOOS {
	case "linux":
		info.WebViewInfo = getLinuxWebViewInfo()
		info.GTKVersion = getGTKVersionInfo()
	case "darwin":
		info.WebViewInfo = "WKWebView (WebKit)"
	case "windows":
		info.WebViewInfo = "WebView2 (Chromium-based)"
	default:
		info.WebViewInfo = "Unknown"
	}

	return info
}

// SaveReport saves the API report to a file
func (s *APICheckService) SaveReport(report APIReport) error {
	filename := fmt.Sprintf("webview-api-report-%s-%s.json",
		report.Platform.OS,
		time.Now().Format("20060102-150405"))

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	fmt.Printf("Report saved to: %s\n", filename)
	return nil
}

func main() {
	app := application.New(application.Options{
		Name:        "WebView API Check",
		Description: "Check which Web APIs are available in the webview",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(&APICheckService{}),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "WebView API Compatibility Check",
		Width:  1200,
		Height: 800,
		URL:    "/",
	})

	err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
