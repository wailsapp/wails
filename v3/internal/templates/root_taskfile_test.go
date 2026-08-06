package templates

import (
	"strings"
	"testing"
)

func TestRootTaskfile(t *testing.T) {
	content, err := RootTaskfile("my-app")
	if err != nil {
		t.Fatalf("RootTaskfile failed: %v", err)
	}
	rendered := string(content)
	if !strings.Contains(rendered, `APP_NAME: "my-app"`) {
		t.Errorf("expected APP_NAME to be set, got:\n%s", rendered)
	}
	// The task template placeholders must survive rendering.
	if !strings.Contains(rendered, `{{.PACKAGE_MANAGER | default "npm"}}`) {
		t.Errorf("expected task template placeholders to be preserved")
	}
	if strings.Contains(rendered, "{{.Opn}}") || strings.Contains(rendered, "{{.Cls}}") {
		t.Errorf("unrendered Opn/Cls placeholders remain")
	}
}
