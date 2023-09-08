package templates

import (
	"os"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestInstall(t *testing.T) {

	tests := []struct {
		name    string
		options *flags.Init
		wantErr bool
	}{
		{
			name: "should install template",
			options: &flags.Init{
				ProjectName:  "test",
				TemplateName: "svelte",
				Quiet:        false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// Remove test directory if it exists
		if _, err := os.Stat(tt.options.ProjectName); err == nil {
			_ = os.RemoveAll(tt.options.ProjectName)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := Install(tt.options); (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
