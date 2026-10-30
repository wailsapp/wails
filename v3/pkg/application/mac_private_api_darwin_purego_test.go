//go:build darwin && purego && !ios && !server && private_mac_apis

package application

import "testing"

// wailsSetKVCValue guards each write on the selectors KVC itself looks for, so
// the selector names must match KVC's capitalisation exactly. A wrong name
// silently skips a write the cgo bridge performs.
func TestKVCSetterSelectors(t *testing.T) {
	tests := map[string][2]string{
		"drawsBackground":        {"setDrawsBackground:", "_setDrawsBackground:"},
		"backgroundColor":        {"setBackgroundColor:", "_setBackgroundColor:"},
		"developerExtrasEnabled": {"setDeveloperExtrasEnabled:", "_setDeveloperExtrasEnabled:"},
		"groupSpacing":           {"setGroupSpacing:", "_setGroupSpacing:"},
		"style":                  {"setStyle:", "_setStyle:"},
		"a":                      {"setA:", "_setA:"},
	}
	for key, want := range tests {
		setter, privateSetter := kvcSetterSelectors(key)
		if setter != want[0] || privateSetter != want[1] {
			t.Errorf("kvcSetterSelectors(%q) = (%q, %q), want (%q, %q)", key, setter, privateSetter, want[0], want[1])
		}
	}
}
