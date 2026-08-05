package main

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestDaymarkShareProviderRendersTextAndHTML(t *testing.T) {
	provider := &daymarkShareProvider{note: sharePayload{
		Title:    "  Saturday <slowly>  ",
		Subtitle: "  A good day  ",
		Body:     "  First line\nSecond & final line.  ",
	}}

	plainText, err := provider.ShareData(application.MacShareRequest{
		ContentType: application.MacShareTypePlainText,
	})
	if err != nil {
		t.Fatalf("render plain text: %v", err)
	}
	if got, want := string(plainText), "Saturday <slowly>\n\nA good day\n\nFirst line\nSecond & final line."; got != want {
		t.Fatalf("plain text = %q, want %q", got, want)
	}

	html, err := provider.ShareData(application.MacShareRequest{
		ContentType: application.MacShareTypeHTML,
	})
	if err != nil {
		t.Fatalf("render HTML: %v", err)
	}
	htmlText := string(html)
	for _, expected := range []string{
		"<title>Saturday &lt;slowly&gt;</title>",
		"<em>A good day</em>",
		"First line<br>\nSecond &amp; final line.",
	} {
		if !strings.Contains(htmlText, expected) {
			t.Errorf("HTML does not contain %q:\n%s", expected, htmlText)
		}
	}
}

func TestDaymarkShareProviderUsesUntitledFallback(t *testing.T) {
	provider := &daymarkShareProvider{}
	data, err := provider.ShareData(application.MacShareRequest{
		ContentType: application.MacShareTypePlainText,
	})
	if err != nil {
		t.Fatalf("render plain text: %v", err)
	}
	if got, want := string(data), "Untitled note"; got != want {
		t.Fatalf("plain text = %q, want %q", got, want)
	}
}
