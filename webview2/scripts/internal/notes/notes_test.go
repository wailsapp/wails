package notes

import "testing"

const sampleNotes = `# WebView2 SDK Release Notes

## Stable Release Notes

[NuGet package for WebView2 1.0.2903.40](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.2903.40)

This release requires WebView2 Runtime version 121.0.2277.83 or higher.

Release Date: October 2024

* Added the ` + "`ICoreWebView2_27`" + ` interface for new feature foo.
* Added the ` + "`ICoreWebView2Profile`" + ` interface.

## Some other section

[NuGet package for WebView2 1.0.2739.15](https://www.nuget.org/packages/Microsoft.Web.WebView2/1.0.2739.15)

This release requires WebView2 Runtime version 119.0.2151.97 or higher.

Release Date: August 2024

* Added the ` + "`ICoreWebView2_26`" + ` interface.
* Added the ` + "`ICoreWebView2_27`" + ` interface in preview (should not win — earlier release wins).
`

func TestParse(t *testing.T) {
	releases, err := Parse([]byte(sampleNotes))
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 2 {
		t.Fatalf("expected 2 releases, got %d", len(releases))
	}
	if releases[0].SDKVersion != "1.0.2903.40" {
		t.Errorf("first release SDK = %q", releases[0].SDKVersion)
	}
	if releases[0].RuntimeVersion != "121.0.2277.83" {
		t.Errorf("first release runtime = %q", releases[0].RuntimeVersion)
	}
	if len(releases[0].Notes) == 0 {
		t.Error("first release should have notes")
	}
	if releases[1].SDKVersion != "1.0.2739.15" {
		t.Errorf("second release SDK = %q", releases[1].SDKVersion)
	}
}

func TestInterfaceMinimumVersions(t *testing.T) {
	releases, err := Parse([]byte(sampleNotes))
	if err != nil {
		t.Fatal(err)
	}
	got := InterfaceMinimumVersions(releases)

	// _26 only appears in the older release.
	if got["ICoreWebView2_26"] != "1.0.2739.15" {
		t.Errorf("ICoreWebView2_26 = %q, want 1.0.2739.15", got["ICoreWebView2_26"])
	}
	// _27 appears in both — the older one (2739.15) should win.
	if got["ICoreWebView2_27"] != "1.0.2739.15" {
		t.Errorf("ICoreWebView2_27 = %q, want 1.0.2739.15 (oldest mention)", got["ICoreWebView2_27"])
	}
	// Profile only in the newer release.
	if got["ICoreWebView2Profile"] != "1.0.2903.40" {
		t.Errorf("ICoreWebView2Profile = %q, want 1.0.2903.40", got["ICoreWebView2Profile"])
	}
}
