package project_test

import (
	"github.com/wailsapp/wails/v2/internal/project"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProject_GetFrontendDir(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON string
		want      string
		wantError bool
	}{
		{
			name:      "Should use 'frontend' by default",
			inputJSON: "{}",
			want:      "frontend",
			wantError: false,
		},
		{
			name:      "Should resolve a relative path with no project path",
			inputJSON: `{"frontend:dir": "./frontend"}`,
			want:      "frontend",
			wantError: false,
		},
		{
			name:      "Should resolve a relative path with project path set",
			inputJSON: `{"frontend:dir": "./frontend", "projectdir": "/home/user/project"}`,
			want:      "/home/user/project/frontend",
			wantError: false,
		},
		{
			name: "Should honour an absolute path",
			inputJSON: func() string {
				if runtime.GOOS == "windows" {
					return `{"frontend:dir": "C:\\frontend", "projectdir": "C:\\project"}`
				} else {
					return `{"frontend:dir": "/home/myproject/frontend", "projectdir": "/home/user/project"}`
				}
			}(),
			want: func() string {
				if runtime.GOOS == "windows" {
					return `C:/frontend`
				} else {
					return `/home/myproject/frontend`
				}
			}(),
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj, err := project.Parse([]byte(tt.inputJSON))
			if err != nil && !tt.wantError {
				t.Errorf("Error parsing project: %s", err)
			}
			got := proj.GetFrontendDir()
			got = filepath.ToSlash(got)
			if got != tt.want {
				t.Errorf("GetFrontendDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
