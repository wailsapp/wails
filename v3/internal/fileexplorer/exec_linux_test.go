//go:build linux

package fileexplorer

import (
	"slices"
	"testing"
)

func TestSplitDesktopExec(t *testing.T) {
	for _, tc := range []struct {
		name, command string
		want          []string
	}{
		{"plain", "manager --new-window %U", []string{"manager", "--new-window", "%U"}},
		{"executable", `"/opt/My Files/files" %f`, []string{"/opt/My Files/files", "%f"}},
		{"argument", `manager "two words"`, []string{"manager", "two words"}},
		{"empty argument", `manager "" %f`, []string{"manager", "", "%f"}},
		{"quote", `manager "say \"hello\""`, []string{"manager", `say "hello"`}},
		{"dollar", `manager "price \\$5"`, []string{"manager", "price $5"}},
		{"backslash", `manager "C:\\\\files"`, []string{"manager", `C:\files`}},
		{"space escape", `manager "two\swords"`, []string{"manager", "two words"}},
		{"tab escape", `manager "two\twords"`, []string{"manager", "two\twords"}},
		{"unbalanced", `manager "open`, nil},
		{"empty executable", `"" %f`, nil},
		{"invalid escape", `manager "bad\q"`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := splitDesktopExec(tc.command); !slices.Equal(got, tc.want) {
				t.Errorf("arguments = %q, want %q", got, tc.want)
			}
		})
	}
}
