// Command sponsorkit generates the Wails sponsors SVG from live GitHub
// Sponsors data. It is a dependency-free Go replacement for the Node
// sponsorkit package.
//
// Usage:
//
//	SPONSORKIT_GITHUB_TOKEN=... go run . -login leaanthony -out sponsors.svg
//
// The token needs to belong to (or have sponsor visibility for) the
// sponsored account so that tier amounts are included; without them every
// sponsor lands in the catch-all tier.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	login := flag.String("login", "leaanthony", "GitHub login whose sponsors to fetch")
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

	fmt.Printf("Fetching sponsors of %s...\n", *login)
	sponsors, err := FetchSponsors(client, token, *login)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d sponsors\n", len(sponsors))

	groups := GroupByTier(sponsors)
	sizes := map[string]int{}
	for i, group := range groups {
		if len(group) > 0 {
			fmt.Printf("  %-18s %d\n", tiers[i].Title, len(group))
		}
		for _, s := range group {
			sizes[s.Login] = int(tiers[i].Avatar) * *scale
		}
	}

	fmt.Println("Fetching avatars...")
	uris := FetchAvatars(client, sponsors, sizes, *quality)

	svg := Render(sponsors, RenderOptions{
		Width:      *width,
		SidePad:    40,
		SponsorURL: "https://github.com/sponsors/" + *login,
		AvatarURIs: uris,
	})
	if err := os.WriteFile(*out, []byte(svg), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	info, _ := os.Stat(*out)
	fmt.Printf("Wrote %s (%d KiB)\n", *out, info.Size()/1024)
}
