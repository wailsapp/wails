package installerbuild

import (
	"os"
	"strings"
	"testing"
)

func TestInstallerWebView2ControlLeavesProjectScriptUntouched(t *testing.T) {
	options := testOptions(t)
	original, err := os.ReadFile(options.Script)
	if err != nil {
		t.Fatal(err)
	}
	for _, skip := range []bool{false, true} {
		options.SkipWebView2 = skip
		wrapper, err := wrapperScript(options.Script, options.Output, options)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(wrapper, "!define WAILS_SKIP_WEBVIEW2 ") != skip {
			t.Fatalf("skip=%v wrapper=%s", skip, wrapper)
		}
		after, err := os.ReadFile(options.Script)
		if err != nil || string(after) != string(original) {
			t.Fatalf("project source modified: %v", err)
		}
	}
}
