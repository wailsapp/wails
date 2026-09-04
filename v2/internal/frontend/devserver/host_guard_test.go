//go:build dev

package devserver

import (
	"testing"
)

func TestHostAllowed(t *testing.T) {
	cases := []struct {
		name               string
		requestHost        string
		bindHost, bindPort string
		extraHosts         []string
		want               bool
	}{
		// Default localhost:34115 bind.
		{"localhost matches", "localhost:34115", "localhost", "34115", nil, true},
		{"localhost is case-insensitive", "LOCALHOST:34115", "localhost", "34115", nil, true},
		{"wrong port is rejected", "localhost:9999", "localhost", "34115", nil, false},
		{"elided port defaults to 80 and misses", "localhost", "localhost", "34115", nil, false},
		{"trailing dot is not localhost", "localhost.:34115", "localhost", "34115", nil, false},
		{"subdomain of localhost is rejected", "app.localhost:34115", "localhost", "34115", nil, false},
		{"loopback IPv4 literal", "127.0.0.1:34115", "localhost", "34115", nil, true},
		{"loopback IPv6 literal", "[::1]:34115", "localhost", "34115", nil, true},
		{"unspecified IPv4 literal", "0.0.0.0:34115", "localhost", "34115", nil, true},
		{"unspecified IPv6 literal", "[::]:34115", "localhost", "34115", nil, true},
		{"IPv6 literal with zone", "[fe80::1%en0]:34115", "localhost", "34115", nil, true},
		{"unbracketed IPv6 is rejected", "::1:34115", "localhost", "34115", nil, false},
		{"shorthand IPv4 is rejected", "127.1:34115", "localhost", "34115", nil, false},
		{"octal IPv4 is rejected", "0177.0.0.1:34115", "localhost", "34115", nil, false},
		{"foreign name is rejected (rebinding)", "attacker.example:34115", "localhost", "34115", nil, false},
		{"localhost-prefixed name is rejected", "localhost.attacker.example:34115", "localhost", "34115", nil, false},
		{"empty host is rejected", "", "localhost", "34115", nil, false},
		{"malformed host is rejected", "a:b:c", "localhost", "34115", nil, false},

		// Wildcard bind ":34115" — bind host is empty; only the CLI's own
		// empty-host poll matches, never a browser.
		{"empty bind host matches empty request host", ":34115", "", "34115", nil, true},
		{"empty bind host still rejects a name", "evil.example:34115", "", "34115", nil, false},

		// 0.0.0.0 bind — LAN device testing reaches the machine by IP.
		{"LAN IP under wildcard bind", "192.168.1.5:34115", "0.0.0.0", "34115", nil, true},
		{"LAN hostname under wildcard bind is rejected", "mylaptop.local:34115", "0.0.0.0", "34115", nil, false},

		// Explicit hostname bind.
		{"configured bind host is case-insensitive", "MYBOX.LAN:34115", "mybox.lan", "34115", nil, true},

		// Env-var escape hatch (entries are already lower-cased by parseAllowedHosts).
		{"extra host without port", "dev.example.com", "localhost", "34115", []string{"dev.example.com"}, true},
		{"extra host is port-exempt", "dev.example.com:443", "localhost", "34115", []string{"dev.example.com"}, true},
		{"extra host is case-insensitive", "DEV.EXAMPLE.COM:443", "localhost", "34115", []string{"dev.example.com"}, true},
		{"subdomain of an extra host is rejected", "sub.dev.example.com:34115", "localhost", "34115", []string{"dev.example.com"}, false},
	}

	for _, c := range cases {
		if got := hostAllowed(c.requestHost, c.bindHost, c.bindPort, c.extraHosts); got != c.want {
			t.Errorf("%s: hostAllowed(%q, %q, %q, %v) = %v, want %v",
				c.name, c.requestHost, c.bindHost, c.bindPort, c.extraHosts, got, c.want)
		}
	}
}

func TestParseAllowedHosts(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  []string
	}{
		{"empty", "", nil},
		{"whitespace only", "   ", nil},
		{"simple list", "a,b", []string{"a", "b"}},
		{"trims, lower-cases, drops empties", " A , ,b ", []string{"a", "b"}},
		{"lower-cases a hostname", "Dev.Example.COM", []string{"dev.example.com"}},
	}

	for _, c := range cases {
		got := parseAllowedHosts(c.value)
		if len(got) != len(c.want) {
			t.Errorf("%s: parseAllowedHosts(%q) = %v, want %v", c.name, c.value, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: parseAllowedHosts(%q) = %v, want %v", c.name, c.value, got, c.want)
				break
			}
		}
	}
}
