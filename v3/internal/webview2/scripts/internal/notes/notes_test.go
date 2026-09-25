package notes

import (
	"strings"
	"testing"
)

const sampleNotes = `# WebView2 SDK Release Notes

## Stable Release Notes

[NuGet package for WebView2 1.0.2903.40](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.2903.40)

This release requires WebView2 Runtime version 121.0.2277.83 or higher.

Release Date: October 2024

* [ICoreWebView2_27 interface](/microsoft-edge/webview2/reference/win32/icorewebview2_27?view=webview2-1.0.2903.40)
* [ICoreWebView2Profile interface](/microsoft-edge/webview2/reference/win32/icorewebview2profile?view=webview2-1.0.2903.40)
* [ICoreWebView2_27::add_NewFeature](url) — a method link, must not be counted as a new interface.

## Some other section

[NuGet package for WebView2 1.0.2739.15](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.2739.15)

This release requires WebView2 Runtime version 119.0.2151.97 or higher.

Release Date: August 2024

* [ICoreWebView2_26 interface](/microsoft-edge/webview2/reference/win32/icorewebview2_26?view=webview2-1.0.2739.15)

[NuGet package for WebView2 1.0.2739.15-prerelease](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.2739.15-prerelease)

This release requires WebView2 Runtime version 119.0.2151.97 or higher.

* [ICoreWebView2_26 interface](/microsoft-edge/webview2/reference/win32/icorewebview2_26?view=webview2-1.0.2739.15-prerelease) — prerelease, must be skipped.
* Mentioned ICoreWebView2Host inline as a rename — no link, no match.
`

func TestParse(t *testing.T) {
	releases, err := Parse([]byte(sampleNotes))
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 3 {
		t.Fatalf("expected 3 releases (2 stable + 1 prerelease), got %d", len(releases))
	}
	if releases[0].SDKVersion != "1.0.2903.40" {
		t.Errorf("first release SDK = %q", releases[0].SDKVersion)
	}
	if releases[0].RuntimeVersion != "121.0.2277.83" {
		t.Errorf("first release runtime = %q", releases[0].RuntimeVersion)
	}
}

func TestInterfaceMinimumVersions(t *testing.T) {
	releases, err := Parse([]byte(sampleNotes))
	if err != nil {
		t.Fatal(err)
	}
	got := InterfaceMinimumVersions(releases)

	if got["ICoreWebView2_26"] != "1.0.2739.15" {
		t.Errorf("ICoreWebView2_26 = %q, want 1.0.2739.15", got["ICoreWebView2_26"])
	}
	if got["ICoreWebView2_27"] != "1.0.2903.40" {
		t.Errorf("ICoreWebView2_27 = %q, want 1.0.2903.40 (only stable mention)", got["ICoreWebView2_27"])
	}
	if got["ICoreWebView2Profile"] != "1.0.2903.40" {
		t.Errorf("ICoreWebView2Profile = %q, want 1.0.2903.40", got["ICoreWebView2Profile"])
	}
	// Renames mentioned without backticks must NOT appear.
	if _, ok := got["ICoreWebView2Host"]; ok {
		t.Error("ICoreWebView2Host should not be in mapping (mentioned only as rename)")
	}
	if _, ok := got["ICoreWebView2Controller"]; ok {
		t.Error("ICoreWebView2Controller should not be in mapping (mentioned only as rename target, no backticks)")
	}
}

// Microsoft restructured the release notes around the "Phase 1/2/3" promotion
// vocabulary and dropped the " interface" suffix from top-level interface
// links. The parser must keep accepting both forms.
const sampleNotesNewFormat = ` ## Release SDK 1.0.3405.78, for Runtime 134 (Aug. 13, 2025)

[NuGet package for WebView2 SDK 1.0.3405.78](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.3405.78)

For full API compatibility, this Release version of the WebView2 SDK requires WebView2 Runtime version 134.0.3405.78 or higher.

#### Promotions to Phase 3 (Stable in Release)

##### [Win32/C++](#tab/win32cpp)

* [ICoreWebView2_28](/microsoft-edge/webview2/reference/win32/icorewebview2_28?view=webview2-1.0.3405.78&preserve-view=true)
   * [ICoreWebView2_28::get_Find](/microsoft-edge/webview2/reference/win32/icorewebview2_28?view=webview2-1.0.3405.78&preserve-view=true#get_find)

* [ICoreWebView2Find](/microsoft-edge/webview2/reference/win32/icorewebview2find?view=webview2-1.0.3405.78&preserve-view=true)
   * [ICoreWebView2Find::Start](/microsoft-edge/webview2/reference/win32/icorewebview2find?view=webview2-1.0.3405.78&preserve-view=true#start)
`

func TestInterfaceMinimumVersionsNewFormat(t *testing.T) {
	releases, err := Parse([]byte(sampleNotesNewFormat))
	if err != nil {
		t.Fatal(err)
	}
	got := InterfaceMinimumVersions(releases)
	if got["ICoreWebView2_28"] != "1.0.3405.78" {
		t.Errorf("ICoreWebView2_28 = %q, want 1.0.3405.78 (new bullet format without ' interface')", got["ICoreWebView2_28"])
	}
	if got["ICoreWebView2Find"] != "1.0.3405.78" {
		t.Errorf("ICoreWebView2Find = %q, want 1.0.3405.78", got["ICoreWebView2Find"])
	}
	// Method links must not be counted as interfaces.
	for k := range got {
		if strings.Contains(k, "::") {
			t.Errorf("method link leaked into interface map: %q", k)
		}
	}
}

func TestIsPrerelease(t *testing.T) {
	cases := map[string]bool{
		"1.0.2903.40":            false,
		"1.0.2739.15-prerelease": true,
		"0.9.515 prerelease":     true,
		"1.0.500-preview":        true,
	}
	for v, want := range cases {
		if got := IsPrerelease(v); got != want {
			t.Errorf("IsPrerelease(%q) = %v, want %v", v, got, want)
		}
	}
}
