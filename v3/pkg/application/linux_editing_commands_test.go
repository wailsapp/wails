//go:build linux && !android && !server

package application

import "testing"

// TestEditingRoles checks what each Edit-menu role sends, on every Linux
// backend: WebKit's editing command for the clipboard and selection roles,
// document.execCommand for undo and redo (webview_window_linux_editing.go).
func TestEditingRoles(t *testing.T) {
	cases := []struct {
		name string
		run  func(*linuxWebviewWindow)
		want string
	}{
		{"cut", (*linuxWebviewWindow).cut, "Cut"},
		{"copy", (*linuxWebviewWindow).copy, "Copy"},
		{"paste", (*linuxWebviewWindow).paste, "Paste"},
		{"delete", (*linuxWebviewWindow).delete, "Delete"},
		{"selectAll", (*linuxWebviewWindow).selectAll, "SelectAll"},
		{"undo", (*linuxWebviewWindow).undo, "document.execCommand('undo')"},
		{"redo", (*linuxWebviewWindow).redo, "document.execCommand('redo')"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sent []string
			prev := editingHook
			editingHook = func(s string) { sent = append(sent, s) }
			t.Cleanup(func() { editingHook = prev })
			tc.run(&linuxWebviewWindow{})
			if len(sent) != 1 || sent[0] != tc.want {
				t.Fatalf("%s sent %q, want [%q]", tc.name, sent, tc.want)
			}
		})
	}
}
