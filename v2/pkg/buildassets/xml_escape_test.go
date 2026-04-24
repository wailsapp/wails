package buildassets

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/internal/project"
)

func TestXmlEscapeAmpersand(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Joe & Bill", "Joe &amp; Bill"},
		{"<tag>", "&lt;tag&gt;"},
		{`"quoted"`, "&#34;quoted&#34;"},
		{"normal", "normal"},
	}
	for _, tt := range tests {
		got := xmlEscape(tt.input)
		if got != tt.expected {
			t.Errorf("xmlEscape(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestResolveProjectDataEscapesAmpersand(t *testing.T) {
	copyright := "Joe & Bill, Inc."
	comments := "A <test> & more"
	pd := &project.Project{
		Name:           "My App & Co",
		OutputFilename: "my-app",
		Info: project.Info{
			CompanyName:    "Test & Co",
			ProductName:    "Test\"Product",
			ProductVersion: "1.0.0",
			Copyright:      &copyright,
			Comments:       &comments,
		},
	}

	tmpl := `<key>Copyright</key>
<string>{{xml .Info.Copyright}}</string>
<key>Comments</key>
<string>{{xml .Info.Comments}}</string>
<key>Name</key>
<string>{{xml .Name}}</string>`

	content, err := resolveProjectData([]byte(tmpl), pd)
	if err != nil {
		t.Fatalf("resolveProjectData() error = %v", err)
	}

	result := string(content)
	if strings.Contains(result, "Joe & Bill") && !strings.Contains(result, "Joe &amp; Bill") {
		t.Error("Ampersand in copyright was not escaped")
	}
	if !strings.Contains(result, "Joe &amp; Bill, Inc.") {
		t.Errorf("Expected escaped copyright, got: %s", result)
	}
	if !strings.Contains(result, "A &lt;test&gt; &amp; more") {
		t.Errorf("Expected escaped comments, got: %s", result)
	}
	if !strings.Contains(result, "My App &amp; Co") {
		t.Errorf("Expected escaped name, got: %s", result)
	}
}
