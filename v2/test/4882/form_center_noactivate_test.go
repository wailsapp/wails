package main

import (
	"os"
	"regexp"
	"testing"
)

func TestCenterUsesNoActivateFlag(t *testing.T) {
	paths := []string{
		"../../internal/frontend/desktop/windows/winc/form.go",
		"internal/frontend/desktop/windows/winc/form.go",
	}

	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Skip("form.go not found (only runs in full repo)")
	}

	content := string(data)

	re := regexp.MustCompile(`SWP_NOSIZE.*SWP_NOACTIVATE`)
	if !re.MatchString(content) {
		t.Error("Center() in form.go does not include SWP_NOACTIVATE flag alongside SWP_NOSIZE in SetWindowPos call")
	}
}
