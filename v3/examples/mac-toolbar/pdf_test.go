package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestRenderDaymarkPDF(t *testing.T) {
	document, err := renderDaymarkPDF(sharePayload{
		Title:    "Saturday, slowly.",
		Subtitle: "A slow day is still a day well spent.",
		Body:     "A good day has room around it.\n\nLeave the phone at home.",
	})
	if err != nil {
		t.Fatalf("render PDF: %v", err)
	}
	if !bytes.HasPrefix(document, []byte("%PDF-1.4")) || !bytes.HasSuffix(document, []byte("%%EOF\n")) {
		t.Fatal("renderer did not produce a complete PDF document")
	}
	if !bytes.Contains(document, []byte("/Type /Page")) ||
		!bytes.Contains(document, []byte("Saturday, slowly.")) {
		t.Fatal("PDF is missing its page or note content")
	}

	marker := []byte("startxref\n")
	markerIndex := bytes.LastIndex(document, marker)
	if markerIndex < 0 {
		t.Fatal("PDF is missing startxref")
	}
	value := strings.TrimSpace(strings.SplitN(string(document[markerIndex+len(marker):]), "\n", 2)[0])
	offset, err := strconv.Atoi(value)
	if err != nil || offset < 0 || offset >= len(document) || !bytes.HasPrefix(document[offset:], []byte("xref\n")) {
		t.Fatalf("invalid xref offset %q", value)
	}
}

func TestDaymarkShareProviderOffersPDF(t *testing.T) {
	provider := &daymarkShareProvider{note: sharePayload{Title: "A note", Body: "Share me"}}
	representations := provider.ShareRepresentations()
	if len(representations) == 0 || representations[0].ContentType != "com.adobe.pdf" {
		t.Fatalf("representations = %#v", representations)
	}
	document, err := provider.ShareData(application.MacShareRequest{
		ContentType:   application.MacShareTypePDF,
		SuggestedName: "A note",
	})
	if err != nil || !bytes.HasPrefix(document, []byte("%PDF-1.4")) {
		t.Fatalf("PDF representation failed: %v", err)
	}
}
