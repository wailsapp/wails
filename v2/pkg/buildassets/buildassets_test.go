package buildassets

import (
	"testing"

	"github.com/wailsapp/wails/v2/internal/project"
)

func TestResolveProjectData_SanitizeIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		template string
		project  *project.Project
		want     string
	}{
		{
			name:     "spaces replaced with hyphens",
			template: `com.wails.{{.Name | sanitizeIdentifier}}`,
			project: &project.Project{
				Name: "My Fab Application",
			},
			want: `com.wails.My-Fab-Application`,
		},
		{
			name:     "underscores replaced with hyphens",
			template: `com.wails.{{.Name | sanitizeIdentifier}}`,
			project: &project.Project{
				Name: "my_app_name",
			},
			want: `com.wails.my-app-name`,
		},
		{
			name:     "alphanumeric and dots preserved",
			template: `com.wails.{{.Name | sanitizeIdentifier}}`,
			project: &project.Project{
				Name: "My.App.v2",
			},
			want: `com.wails.My.App.v2`,
		},
		{
			name:     "already clean identifier",
			template: `com.wails.{{.Name | sanitizeIdentifier}}`,
			project: &project.Project{
				Name: "MyApp",
			},
			want: `com.wails.MyApp`,
		},
		{
			name:     "special characters replaced",
			template: `com.wails.{{.Name | sanitizeIdentifier}}`,
			project: &project.Project{
				Name: "App@Beta#3!",
			},
			want: `com.wails.App-Beta-3-`,
		},
		{
			name:     "outputfilename used for executable",
			template: `<string>{{.OutputFilename}}</string>`,
			project: &project.Project{
				Name:           "My Fab Application",
				OutputFilename: "my-fab-application",
			},
			want: `<string>my-fab-application</string>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveProjectData([]byte(tt.template), tt.project)
			if err != nil {
				t.Fatalf("resolveProjectData() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("resolveProjectData() = %q, want %q", string(got), tt.want)
			}
		})
	}
}
