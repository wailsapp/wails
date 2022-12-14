package project_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/internal/project"
)

func TestProject_GetFrontendDir(t *testing.T) {
	cwd := lo.Must(os.Getwd())
	tests := []struct {
		name      string
		inputJSON string
		want      string
		wantError bool
	}{
		{
			name:      "Should use 'frontend' by default",
			inputJSON: "{}",
			want:      filepath.ToSlash(filepath.Join(cwd, "frontend")),
			wantError: false,
		},
		{
			name:      "Should resolve a relative path with no project path",
			inputJSON: `{"frontend:dir": "./frontend"}`,
			want:      filepath.ToSlash(filepath.Join(cwd, "frontend")),
			wantError: false,
		},
		{
			name: "Should resolve a relative path with project path set",
			inputJSON: func() string {
				if runtime.GOOS == "windows" {
					return `{"frontend:dir": "./frontend", "projectdir": "C:\\project"}`
				} else {
					return `{"frontend:dir": "./frontend", "projectdir": "/home/user/project"}`
				}
			}(),
			want: func() string {
				if runtime.GOOS == "windows" {
					return `C:/project/frontend`
				} else {
					return `/home/user/project/frontend`
				}
			}(),
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
func TestProject_GetBuildDir(t *testing.T) {
	cwd := lo.Must(os.Getwd())
	tests := []struct {
		name      string
		inputJSON string
		want      string
		wantError bool
	}{
		{
			name:      "Should use 'build' by default",
			inputJSON: "{}",
			want:      filepath.ToSlash(filepath.Join(cwd, "build")),
			wantError: false,
		},
		{
			name:      "Should resolve a relative path with no project path",
			inputJSON: `{"build:dir": "./build"}`,
			want:      filepath.ToSlash(filepath.Join(cwd, "build")),
			wantError: false,
		},
		{
			name:      "Should resolve a relative path with project path set",
			inputJSON: `{"build:dir": "./build", "projectdir": "/home/user/project"}`,
			want:      "/home/user/project/build",
			wantError: false,
		},
		{
			name: "Should honour an absolute path",
			inputJSON: func() string {
				if runtime.GOOS == "windows" {
					return `{"build:dir": "C:\\build", "projectdir": "C:\\project"}`
				} else {
					return `{"build:dir": "/home/myproject/build", "projectdir": "/home/user/project"}`
				}
			}(),
			want: func() string {
				if runtime.GOOS == "windows" {
					return `C:/build`
				} else {
					return `/home/myproject/build`
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
			got := proj.GetBuildDir()
			got = filepath.ToSlash(got)
			if got != tt.want {
				t.Errorf("GetFrontendDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
