//go:build dev

package devserver

import (
	"net/http"
	"testing"
)

func req(origin, host string) *http.Request {
	r := &http.Request{Host: host, Header: http.Header{}}
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	return r
}

func TestSameOrigin(t *testing.T) {
	cases := []struct {
		name         string
		origin, host string
		want         bool
	}{
		{"same origin page", "http://localhost:34115", "localhost:34115", true},
		{"same origin over https", "https://localhost:34115", "localhost:34115", true},
		{"host comparison is case-insensitive", "http://LocalHost:34115", "localhost:34115", true},
		{"same origin IPv6", "http://[::1]:34115", "[::1]:34115", true},
		{"headerless client is now refused", "", "localhost:34115", false},
		{"foreign https page", "https://evil.example", "localhost:34115", false},
		{"foreign loopback port", "http://localhost:9999", "localhost:34115", false},
		{"unparseable origin", "://nope", "localhost:34115", false},
		{"null origin", "null", "localhost:34115", false},
		{"file origin", "file:///etc/passwd", "localhost:34115", false},
		{"wails scheme origin", "wails://wails", "localhost:34115", false},
		{"extension origin", "chrome-extension://abcdef", "localhost:34115", false},
	}
	for _, c := range cases {
		if got := sameOrigin(req(c.origin, c.host)); got != c.want {
			t.Errorf("%s: origin=%q host=%q got %v want %v", c.name, c.origin, c.host, got, c.want)
		}
	}
}
