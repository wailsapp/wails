package application

import (
	"os"
	"strings"
	"testing"
)

func TestPopupmenuBuildMenuUsesErrorNotFatal(t *testing.T) {
	data, err := os.ReadFile("popupmenu_windows.go")
	if err != nil {
		t.Skip("popupmenu_windows.go not available")
	}
	content := string(data)

	if strings.Contains(content, `globalApplication.fatal("error adding menu item`) {
		t.Error("buildMenu should not call fatal() for menu item errors - use error() instead")
	}
	if strings.Contains(content, `globalApplication.fatal("error setting menu icons`) {
		t.Error("buildMenu should not call fatal() for menu icon errors - use error() instead")
	}

	if !strings.Contains(content, `globalApplication.error("error adding menu item`) {
		t.Error("buildMenu should use error() for menu item failures")
	}
}
