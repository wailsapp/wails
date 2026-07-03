package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// FetchAvatars downloads every sponsor's avatar at pixel size px (already
// scaled for hi-dpi), re-encodes it as JPEG and returns data URIs keyed by
// login. Failures are reported but non-fatal: missing entries fall back to
// a placeholder at render time.
func FetchAvatars(client *http.Client, sponsors []Sponsor, sizes map[string]int, quality int) map[string]string {
	var mu sync.Mutex
	uris := make(map[string]string, len(sponsors))
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup

	for _, s := range sponsors {
		wg.Add(1)
		go func(s Sponsor) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			uri, err := fetchAvatar(client, s.AvatarURL, sizes[s.Login], quality)
			if err != nil {
				fmt.Printf("  warning: avatar for %s: %v\n", s.Login, err)
				return
			}
			mu.Lock()
			uris[s.Login] = uri
			mu.Unlock()
		}(s)
	}
	wg.Wait()
	return uris
}

func fetchAvatar(client *http.Client, avatarURL string, px, quality int) (string, error) {
	u, err := url.Parse(avatarURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("s", strconv.Itoa(px))
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %s", resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		// Unknown format (e.g. webp): embed the original bytes untouched.
		mime := resp.Header.Get("Content-Type")
		if mime == "" {
			return "", fmt.Errorf("undecodable avatar with no content type")
		}
		return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// Flatten any transparency onto the card background, then JPEG-encode.
	bounds := img.Bounds()
	flat := image.NewRGBA(bounds)
	draw.Draw(flat, bounds, image.NewUniform(cardBackground), image.Point{}, draw.Src)
	draw.Draw(flat, bounds, img, bounds.Min, draw.Over)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, flat, &jpeg.Options{Quality: quality}); err != nil {
		return "", err
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
