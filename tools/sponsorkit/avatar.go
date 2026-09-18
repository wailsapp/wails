package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// Avatar mask exponents: the superellipse |x/r|^n + |y/r|^n = 1 is a circle
// at n=2 and the contributor squircle at n=4 (matching squirclePathD).
const (
	maskCircle   = 2.0
	maskSquircle = 4.0
)

// FetchAvatars downloads every sponsor's avatar at pixel size px (already
// scaled for hi-dpi), bakes the shape mask into the bitmap and returns data
// URIs keyed by login. The mask is applied to the pixels rather than left to
// a runtime clip-path because browsers do not clip <image> elements reliably
// (WebKit, and Chromium's software rasterizer, intermittently paint the raw
// square bitmap). Failures are reported but non-fatal: missing entries fall
// back to a placeholder at render time.
func FetchAvatars(client *http.Client, sponsors []Sponsor, sizes map[string]int, quality int, maskExp float64) map[string]string {
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

			uri, err := fetchAvatar(client, s.AvatarURL, sizes[s.Login], quality, maskExp)
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

func fetchAvatar(client *http.Client, avatarURL string, px, quality int, maskExp float64) (string, error) {
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
		// Unknown format (e.g. webp): embed the original bytes untouched and
		// leave the shaping to the renderer's clip-path.
		mime := resp.Header.Get("Content-Type")
		if mime == "" {
			return "", fmt.Errorf("undecodable avatar with no content type")
		}
		return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// Flatten any source transparency onto the card background so the
	// interior looks identical in every embed context, then cut the shape
	// mask into the alpha channel.
	bounds := img.Bounds()
	flat := image.NewRGBA(bounds)
	draw.Draw(flat, bounds, image.NewUniform(cardBackground), image.Point{}, draw.Src)
	draw.Draw(flat, bounds, img, bounds.Min, draw.Over)
	masked := maskShape(flat, maskExp)

	// JPEG cannot hold alpha, so its corners are flattened back onto the
	// card background colour; photos compress far better this way. Flat
	// identicon-style avatars compress better as PNG, which also keeps the
	// corners genuinely transparent. Embed whichever is smaller.
	flatCorners := image.NewRGBA(bounds)
	draw.Draw(flatCorners, bounds, image.NewUniform(cardBackground), image.Point{}, draw.Src)
	draw.Draw(flatCorners, bounds, masked, bounds.Min, draw.Over)
	var jbuf bytes.Buffer
	if err := jpeg.Encode(&jbuf, flatCorners, &jpeg.Options{Quality: quality}); err != nil {
		return "", err
	}
	var pbuf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if pal := toPaletted(masked); pal != nil {
		// Flat-colour avatars (identicons, logos) fit a lossless palette,
		// which encodes far smaller than full RGBA.
		if err := enc.Encode(&pbuf, pal); err != nil {
			return "", err
		}
	} else if err := enc.Encode(&pbuf, masked); err != nil {
		return "", err
	}
	if pbuf.Len() < jbuf.Len() {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(pbuf.Bytes()), nil
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(jbuf.Bytes()), nil
}

// toPaletted losslessly converts img to an indexed-colour image, or returns
// nil if it uses more than 256 distinct colours.
func toPaletted(img *image.NRGBA) *image.Paletted {
	b := img.Bounds()
	index := map[color.NRGBA]uint8{}
	var palette color.Palette
	out := image.NewPaletted(b, nil)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			idx, ok := index[c]
			if !ok {
				if len(palette) == 256 {
					return nil
				}
				idx = uint8(len(palette))
				index[c] = idx
				palette = append(palette, c)
			}
			out.SetColorIndex(x, y, idx)
		}
	}
	out.Palette = palette
	return out
}

// maskShape zeroes the alpha of every pixel outside the centred superellipse
// |x/rx|^exp + |y/ry|^exp = 1, with an anti-aliased edge roughly one pixel
// wide, and returns the result.
func maskShape(src *image.RGBA, exp float64) *image.NRGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(out, out.Bounds(), src, b.Min, draw.Src)
	rx, ry := float64(w)/2, float64(h)/2
	for y := 0; y < h; y++ {
		fy := math.Abs((float64(y) + 0.5 - ry) / ry)
		for x := 0; x < w; x++ {
			fx := math.Abs((float64(x) + 0.5 - rx) / rx)
			// Normalised superellipse "radius" of this pixel; the shape edge
			// sits at d=1 and d scales roughly linearly with distance near it.
			d := math.Pow(math.Pow(fx, exp)+math.Pow(fy, exp), 1/exp)
			cov := (1-d)*rx + 0.5
			if cov >= 1 {
				continue
			}
			i := out.PixOffset(x, y)
			if cov <= 0 {
				out.Pix[i+3] = 0
				continue
			}
			out.Pix[i+3] = uint8(float64(out.Pix[i+3])*cov + 0.5)
		}
	}
	return out
}
