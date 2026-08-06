package main

import (
	"fmt"
	"html"
	"log"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type sharePayload struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

type daymarkShareProvider struct {
	lock            sync.RWMutex
	note            sharePayload
	representations []application.MacShareRepresentation
}

func (p *daymarkShareProvider) ShareRepresentations() []application.MacShareRepresentation {
	if len(p.representations) != 0 {
		return append([]application.MacShareRepresentation(nil), p.representations...)
	}
	return []application.MacShareRepresentation{
		{ContentType: application.MacShareTypePDF},
		{ContentType: application.MacShareTypeHTML},
		{ContentType: application.MacShareTypePlainText},
	}
}

func newDaymarkPDFShareProvider(note sharePayload) *daymarkShareProvider {
	return &daymarkShareProvider{
		note: note,
		representations: []application.MacShareRepresentation{
			{ContentType: application.MacShareTypePDF},
		},
	}
}

func (p *daymarkShareProvider) ShareData(request application.MacShareRequest) ([]byte, error) {
	note := p.snapshot()
	log.Printf("native share requested %s", request.ContentType)

	switch request.ContentType {
	case application.MacShareTypePDF:
		return renderDaymarkPDF(note)
	case application.MacShareTypeHTML:
		return renderDaymarkHTML(note), nil
	case application.MacShareTypePlainText:
		return renderDaymarkPlainText(note), nil
	default:
		return nil, fmt.Errorf("Daymark cannot provide %q", request.ContentType)
	}
}

func (p *daymarkShareProvider) snapshot() sharePayload {
	p.lock.RLock()
	note := p.note
	p.lock.RUnlock()
	return normaliseSharePayload(note)
}

func (p *daymarkShareProvider) update(note sharePayload) {
	p.lock.Lock()
	p.note = note
	p.lock.Unlock()
}

func normaliseSharePayload(note sharePayload) sharePayload {
	note.Title = strings.TrimSpace(note.Title)
	if note.Title == "" {
		note.Title = "Untitled note"
	}
	note.Subtitle = strings.TrimSpace(note.Subtitle)
	note.Body = strings.TrimSpace(note.Body)
	return note
}

func renderDaymarkPlainText(note sharePayload) []byte {
	parts := []string{note.Title}
	if note.Subtitle != "" {
		parts = append(parts, note.Subtitle)
	}
	if note.Body != "" {
		parts = append(parts, note.Body)
	}
	return []byte(strings.Join(parts, "\n\n"))
}

func renderDaymarkHTML(note sharePayload) []byte {
	body := strings.ReplaceAll(html.EscapeString(note.Body), "\n", "<br>\n")
	return []byte(fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>%s</title></head>
<body><article><h1>%s</h1><p><em>%s</em></p><p>%s</p></article></body></html>`,
		html.EscapeString(note.Title),
		html.EscapeString(note.Title),
		html.EscapeString(note.Subtitle),
		body,
	))
}
