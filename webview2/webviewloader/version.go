//go:build windows && !native_webview2loader

package webviewloader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// UsingGoWebview2Loader is set to true when the go webview2loader is used.
var UsingGoWebview2Loader bool

// CompareBrowserVersions will compare the 2 given versions and return:
//
//	-1 = v1 < v2
//	 0 = v1 == v2
//	 1 = v1 > v2
func CompareBrowserVersions(v1 string, v2 string) (int, error) {
	v, err := parseVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("v1 invalid: %w", err)
	}

	w, err := parseVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("v2 invalid: %w", err)
	}

	return v.compare(w), nil
}

// GetAvailableCoreWebView2BrowserVersionString get the browser version info including channel name
// if it is the WebView2 Runtime.
// Channel names are Beta, Dev, and Canary.
func GetAvailableCoreWebView2BrowserVersionString(browserExecutableFolder string) (string, error) {
	if browserExecutableFolder != "" {
		clientPath, err := findEmbeddedClientDll(browserExecutableFolder)
		if errors.Is(err, errNoClientDLLFound) {
			// WebView2 is not found
			return "", nil
		} else if err != nil {
			return "", err
		}

		return findEmbeddedBrowserVersion(clientPath)
	}

	_, version, err := findInstalledClientDll(false)
	if errors.Is(err, errNoClientDLLFound) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return version.String(), nil
}

type version struct {
	major int
	minor int
	patch int
	build int

	channel string
}

func (v version) String() string {
	vv := fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.build)
	if v.channel != "" {
		vv += " " + v.channel
	}

	return vv
}

func (v version) compare(o version) int {
	if c := compareInt(v.major, o.major); c != 0 {
		return c
	}
	if c := compareInt(v.minor, o.minor); c != 0 {
		return c
	}
	if c := compareInt(v.patch, o.patch); c != 0 {
		return c
	}
	return compareInt(v.build, o.build)
}

func parseVersion(v string) (version, error) {
	var p version

	// Split away channel information...
	if i := strings.Index(v, " "); i > 0 {
		p.channel = v[i+1:]
		v = v[:i]
	}

	vv := strings.Split(v, ".")
	if len(vv) > 4 {
		return p, fmt.Errorf("too many version parts")
	}

	var err error
	vv, p.major, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad major version: %w", err)
	}

	vv, p.minor, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad minor version: %w", err)
	}

	vv, p.patch, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad patch version: %w", err)
	}

	_, p.build, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad build version: %w", err)
	}

	return p, nil
}

func parseInt(v []string) ([]string, int, error) {
	if len(v) == 0 {
		return nil, 0, nil
	}

	p, err := strconv.ParseInt(v[0], 10, 32)
	if err != nil {
		return nil, 0, err
	}
	return v[1:], int(p), nil
}

func compareInt(v1, v2 int) int {
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	} else {
		return +1
	}
}
