// Command starhistory generates Wails' self-contained GitHub star-history SVG.
//
// The generator intentionally uses only the Go standard library so it can run
// in the repository's existing Go workflows without Python, Docker, or Node.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// main fetches current stargazer data and atomically replaces the chart SVG.
func main() {
	repo := flag.String("repo", "wailsapp/wails", "GitHub repository in owner/name form")
	background := flag.String("background", "../../docs/public/digital_wales_master.webp", "background image to embed in the SVG")
	logo := flag.String("logo", "../../website/static/img/wails-logo-horizontal-dark.svg", "horizontal Wails logo to embed in the SVG")
	out := flag.String("out", "../../website/static/img/star-history.svg", "output SVG path")
	flag.Parse()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("STAR_HISTORY_TOKEN")
	}
	if token == "" {
		fatal("GITHUB_TOKEN or STAR_HISTORY_TOKEN must be set")
	}

	backgroundData, err := os.ReadFile(*background)
	if err != nil {
		fatal("read background image: %v", err)
	}
	logoData, err := os.ReadFile(*logo)
	if err != nil {
		fatal("read Wails logo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client := &GitHubClient{}
	stars, err := client.FetchStars(ctx, token, *repo)
	if err != nil {
		fatal("fetch star history: %v", err)
	}
	history := BuildWeeklyHistory(stars)
	starCount, err := client.FetchStarCount(ctx, token, *repo)
	if err != nil {
		fatal("fetch current star count: %v", err)
	}
	history = ReconcileCurrentCount(history, starCount, time.Now().UTC())

	output := RenderSVG(RenderOptions{
		Repo:           *repo,
		History:        history,
		BackgroundData: backgroundData,
		LogoData:       logoData,
		RefreshedAt:    time.Now().UTC(),
	})
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fatal("create output directory: %v", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(*out), ".star-history-*.svg")
	if err != nil {
		fatal("create temporary output: %v", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.WriteString(output); err != nil {
		_ = tmp.Close()
		fatal("write output: %v", err)
	}
	if err := tmp.Close(); err != nil {
		fatal("close output: %v", err)
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		fatal("set output permissions: %v", err)
	}
	if err := os.Rename(tmpName, *out); err != nil {
		fatal("replace output: %v", err)
	}

	fmt.Printf("Generated %s from %d visible stargazers (%d weekly points)\n", *out, len(stars), len(history))
}

// fatal reports a command error and terminates with a non-zero status.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "starhistory: "+format+"\n", args...)
	os.Exit(1)
}
