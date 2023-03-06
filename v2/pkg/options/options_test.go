package options

import (
	"testing"
)

func TestMergeDefaultsWH(t *testing.T) {
	tests := []struct {
		name       string
		appoptions *App
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "No width and height",
			appoptions: &App{},
			wantWidth:  1024,
			wantHeight: 768,
		},
		{
			name: "Basic width and height",
			appoptions: &App{
				Width:  800,
				Height: 600,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
		{
			name: "With MinWidth and MinHeight",
			appoptions: &App{
				Width:     200,
				MinWidth:  800,
				Height:    100,
				MinHeight: 600,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
		{
			name: "With MaxWidth and MaxHeight",
			appoptions: &App{
				Width:     900,
				MaxWidth:  800,
				Height:    700,
				MaxHeight: 600,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
		{
			name: "With MinWidth more than MaxWidth",
			appoptions: &App{
				Width:    900,
				MinWidth: 900,
				MaxWidth: 800,
				Height:   600,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
		{
			name: "With MinHeight more than MaxHeight",
			appoptions: &App{
				Width:     800,
				Height:    700,
				MinHeight: 900,
				MaxHeight: 600,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeDefaults(tt.appoptions)
			if tt.appoptions.Width != tt.wantWidth {
				t.Errorf("MergeDefaults().Width =%v, want %v", tt.appoptions.Width, tt.wantWidth)
			}
			if tt.appoptions.Height != tt.wantHeight {
				t.Errorf("MergeDefaults().Height =%v, want %v", tt.appoptions.Height, tt.wantHeight)
			}
		})
	}
}
