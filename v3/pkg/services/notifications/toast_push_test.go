package notifications

import (
	"os"
	"strings"
	"testing"
)

// TestWindowsToastDoesNotUsePowershell guards against re-enabling go-toast's
// PowerShell delivery path, which renders notification content into a
// generated script. The Windows notifier must push through WinRT only.
func TestWindowsToastDoesNotUsePowershell(t *testing.T) {
	src, err := os.ReadFile("notifications_windows.go")
	if err != nil {
		t.Fatalf("read notifications_windows.go: %v", err)
	}
	for _, opt := range []string{"wintoast.PowershellFallback", "wintoast.PreferPowershell"} {
		if strings.Contains(string(src), opt) {
			t.Errorf("notifications_windows.go must not use %s", opt)
		}
	}
}
