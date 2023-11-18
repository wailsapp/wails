//go:build full_test

package commands_test

import (
	"github.com/wailsapp/wails/v3/internal/commands"
	"github.com/wailsapp/wails/v3/internal/s"
	"testing"
)

func Test_generateAppImage(t *testing.T) {

	tests := []struct {
		name     string
		options  *commands.GenerateAppImageOptions
		wantErr  bool
		setup    func()
		teardown func()
	}{
		{
			name:    "Should fail if binary path is not provided",
			options: &commands.GenerateAppImageOptions{},
			wantErr: true,
		},
		{
			name: "Should fail if Icon is not provided",
			options: &commands.GenerateAppImageOptions{
				Binary: "testapp",
			},
			wantErr: true,
		},

		{
			name: "Should fail if desktop file is not provided",
			options: &commands.GenerateAppImageOptions{
				Binary: "testapp",
				Icon:   "testicon",
			},
			wantErr: true,
		},
		{
			name: "Should work if inputs are valid",
			options: &commands.GenerateAppImageOptions{
				Binary:      "testapp",
				Icon:        "appicon.png",
				DesktopFile: "testapp.desktop",
			},
			setup: func() {
				// Compile the test application
				s.CD("appimage_testfiles")
				testDir := s.CWD()
				_, err := s.EXEC(`go build -ldflags="-s -w" -o testapp`)
				if err != nil {
					t.Fatal(err)
				}
				s.DEFER(func() {
					s.CD(testDir)
					s.RM("testapp")
					s.RM("testapp-x86_64.AppImage")
				})
			},
			teardown: func() {
				s.CALLDEFER()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if err := commands.GenerateAppImage(tt.options); (err != nil) != tt.wantErr {
				t.Errorf("generateAppImage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.teardown != nil {
				tt.teardown()
			}
		})
	}
}
