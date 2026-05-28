package buildassets

import (
	"testing"

	"github.com/wailsapp/wails/v2/internal/project"
)

func strPtr(s string) *string { return &s }

func TestResolveProjectData_XMLEscaping(t *testing.T) {
	tests := []struct {
		name     string
		template string
		project  *project.Project
		want     string
	}{
		{
			name:     "ampersand in name",
			template: `<string>{{.Name}}</string>`,
			project: &project.Project{
				Name: "Tom & Jerry",
			},
			want: `<string>Tom &amp; Jerry</string>`,
		},
		{
			name:     "ampersand in copyright pointer",
			template: `<string>{{.Info.Copyright}}</string>`,
			project: &project.Project{
				Info: project.Info{
					Copyright: strPtr("Joe & Bill, Inc."),
				},
			},
			want: `<string>Joe &amp; Bill, Inc.</string>`,
		},
		{
			name:     "angle brackets in name",
			template: `<string>{{.Name}}</string>`,
			project: &project.Project{
				Name: "<script>alert(1)</script>",
			},
			want: `<string>&lt;script&gt;alert(1)&lt;/script&gt;</string>`,
		},
		{
			name:     "plain text no escaping needed",
			template: `<string>{{.Name}}</string>`,
			project: &project.Project{
				Name: "MyApp",
			},
			want: `<string>MyApp</string>`,
		},
		{
			name:     "multiple ampersands",
			template: `<string>{{.Name}}</string>`,
			project: &project.Project{
				Name: "A&B&C & Co",
			},
			want: `<string>A&amp;B&amp;C &amp; Co</string>`,
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
