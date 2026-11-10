// Command sponsorkit generates the Wails sponsors and contributors SVGs
// from live GitHub data. It is a dependency-free Go replacement for the
// Node sponsorkit package.
//
// Usage:
//
//	SPONSORKIT_GITHUB_TOKEN=... go run . -login leaanthony -out sponsors.svg
//	SPONSORKIT_GITHUB_TOKEN=... go run . -mode contributors -repo wailsapp/wails -out contributors.svg
//
// For sponsors, the token needs to belong to (or have sponsor visibility
// for) the sponsored account so that tier amounts are included; without
// them every sponsor lands in the catch-all tier. For contributors any
// token works.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	mode := flag.String("mode", "sponsors", `what to render: "sponsors" or "contributors"`)
	login := flag.String("login", "leaanthony", "GitHub login whose sponsors to fetch")
	repo := flag.String("repo", "wailsapp/wails", "owner/name repository whose contributors to fetch")
	metric := flag.String("metric", "prs", `how contributors are ranked: "prs" (merged pull requests) or "commits"`)
	changelogs := flag.String("changelogs", "", "comma-separated changelog paths scanned for @login credits (contributors mode)")
	out := flag.String("out", "sponsors.svg", "output SVG path")
	width := flag.Float64("width", 800, "SVG width in CSS pixels")
	scale := flag.Int("scale", 2, "avatar oversampling factor for hi-dpi displays")
	quality := flag.Int("quality", 80, "JPEG quality for embedded avatars")
	flag.Parse()

	token := os.Getenv("SPONSORKIT_GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: set SPONSORKIT_GITHUB_TOKEN or GITHUB_TOKEN")
		os.Exit(1)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	var svg string
	var err error
	switch *mode {
	case "sponsors":
		svg, err = generateSponsors(client, token, *login, *width, *scale, *quality)
	case "contributors":
		switch *metric {
		case "commits":
		case "prs":
			creditMetric = "prs"
			for i := range bands {
				bands[i].MinCredit = bands[i].PRMinCredit
			}
		default:
			fmt.Fprintf(os.Stderr, "error: unknown -metric %q\n", *metric)
			os.Exit(1)
		}
		svg, err = generateContributors(client, token, *repo, splitPaths(*changelogs), *width, *scale, *quality)
	default:
		err = fmt.Errorf("unknown -mode %q", *mode)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*out, []byte(svg), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s (%d KiB)\n", *out, len(svg)/1024)
}

func generateSponsors(client *http.Client, token, login string, width float64, scale, quality int) (string, error) {
	fmt.Printf("Fetching sponsors of %s...\n", login)
	sponsors, err := FetchSponsors(client, token, login)
	if err != nil {
		return "", err
	}
	fmt.Printf("Found %d sponsors\n", len(sponsors))

	groups := GroupByTier(sponsors)
	sizes := map[string]int{}
	for i, group := range groups {
		if len(group) > 0 {
			fmt.Printf("  %-18s %d\n", tiers[i].Title, len(group))
		}
		for _, s := range group {
			sizes[s.Login] = int(tiers[i].Avatar) * scale
		}
	}

	fmt.Println("Fetching avatars...")
	uris := FetchAvatars(client, sponsors, sizes, quality, maskCircle)

	return Render(sponsors, RenderOptions{
		Width:      width,
		SidePad:    40,
		SponsorURL: "https://github.com/sponsors/" + login,
		AvatarURIs: uris,
	}), nil
}

func generateContributors(client *http.Client, token, repo string, changelogs []string, width float64, scale, quality int) (string, error) {
	fmt.Printf("Fetching contributors of %s...\n", repo)
	contributors, err := FetchContributors(client, token, repo, changelogs)
	if err != nil {
		return "", err
	}
	fmt.Printf("Found %d contributors\n", len(contributors))

	groups := GroupByBand(contributors)
	sizes := map[string]int{}
	// With hundreds of contributors the embedded avatars dominate the file
	// size, so the small squircles trade JPEG quality for weight.
	byQuality := map[int][]Sponsor{}
	for i, group := range groups {
		if len(group) > 0 {
			fmt.Printf("  %4d+ credits: %d\n", bands[i].MinCredit, len(group))
		}
		q := quality
		switch {
		case bands[i].Avatar < 40:
			q = 55
		case !bands[i].ShowName:
			q = 68
		}
		for _, c := range group {
			sizes[c.Login] = int(bands[i].Avatar) * scale
			byQuality[q] = append(byQuality[q], Sponsor{Login: c.Login, AvatarURL: c.AvatarURL})
		}
	}

	fmt.Println("Fetching avatars...")
	uris := map[string]string{}
	for q, list := range byQuality {
		for login, uri := range FetchAvatars(client, list, sizes, q, maskSquircle) {
			uris[login] = uri
		}
	}

	return RenderContributors(contributors, ContributorRenderOptions{
		Width:         width,
		SidePad:       40,
		ContributeURL: "https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md",
		AvatarURIs:    uris,
	}), nil
}

func splitPaths(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
