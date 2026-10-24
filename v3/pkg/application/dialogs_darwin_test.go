//go:build darwin && !ios && !server

package application

import "testing"

func TestSaveDialogExtensions(t *testing.T) {
	tests := []struct {
		name    string
		filters []FileFilter
		want    string
	}{
		{"no filters", nil, ""},
		{"single pattern", []FileFilter{{DisplayName: "PDF Document", Pattern: "*.pdf"}}, "pdf"},
		{"bare extension", []FileFilter{{DisplayName: "PDF Document", Pattern: "pdf"}}, "pdf"},
		{"patterns of one filter", []FileFilter{{DisplayName: "Images", Pattern: "*.png; *.jpg"}}, "png;jpg"},
		{
			"filters keep their order",
			[]FileFilter{{DisplayName: "Text", Pattern: "*.txt"}, {DisplayName: "Markdown", Pattern: "*.md"}},
			"txt;md",
		},
		{"wildcards are skipped", []FileFilter{{DisplayName: "All Files", Pattern: "*.*;*"}}, ""},
		{
			"wildcard next to an extension",
			[]FileFilter{{DisplayName: "All Files", Pattern: "*.*"}, {DisplayName: "PDF Document", Pattern: "*.pdf"}},
			"pdf",
		},
		{"dotted extension is skipped", []FileFilter{{DisplayName: "Archive", Pattern: "*.tar.gz"}}, ""},
		{"empty pattern", []FileFilter{{DisplayName: "Empty", Pattern: ""}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveDialogExtensions(tt.filters); got != tt.want {
				t.Errorf("saveDialogExtensions(%v) = %q, want %q", tt.filters, got, tt.want)
			}
		})
	}
}
