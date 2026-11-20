//go:build windows

package w32

import (
	"strings"
	"syscall"
	"testing"
)

// TestGetAllScreensReusesCallback is the regression test for the callback slot
// leak: GetAllScreens used to pass a new closure to syscall.NewCallback on every
// call. The Go runtime never frees those slots and allows only about 2000 per
// process, so the process died with "fatal error: too many callback functions"
// after fewer than 2000 display changes. That is a throw, not a panic, so a
// regression takes the whole test binary down.
func TestGetAllScreensReusesCallback(t *testing.T) {
	for i := 0; i < 2500; i++ {
		// The result is display dependent (a headless runner may have no
		// monitor); only the number of callbacks created matters here.
		_, _ = GetAllScreens()
	}
}

func TestMonitorEnumerationIncompleteError(t *testing.T) {
	// Three monitors side by side: primary, middle, far.
	primary := &Screen{IsPrimary: true, MONITORINFOEX: monitorAt(0, 1920)}
	middle := &Screen{MONITORINFOEX: monitorAt(1920, 3840)}
	far := &Screen{MONITORINFOEX: monitorAt(3840, 5760)}
	secondary := middle

	tests := []struct {
		name     string
		e        monitorEnumeration
		wantErr  bool
		contains []string
	}{
		{
			name: "nothing failed",
			e:    monitorEnumeration{screens: []*Screen{primary, secondary}},
		},
		{
			name: "a secondary monitor was skipped",
			e: monitorEnumeration{
				screens:  []*Screen{primary},
				failures: []string{"GetMonitorInfo failed for monitor 0x2"},
			},
		},
		{
			name: "the skipped monitor was last in the chain",
			e: monitorEnumeration{
				screens:  []*Screen{primary, middle},
				failures: []string{"GetMonitorInfo failed for monitor 0x3"},
			},
		},
		{
			name: "the skipped monitor disconnects a retained one",
			e: monitorEnumeration{
				screens:  []*Screen{primary, far},
				failures: []string{"GetMonitorInfo failed for monitor 0x2"},
			},
			wantErr:  true,
			contains: []string{"not all connected", "monitor 0x2"},
		},
		{
			name: "the primary monitor was skipped",
			e: monitorEnumeration{
				screens:  []*Screen{secondary},
				failures: []string{"GetMonitorInfo failed for monitor 0x1", "GetMonitorInfo failed for monitor 0x3"},
			},
			wantErr:  true,
			contains: []string{"primary monitor", "monitor 0x1", "monitor 0x3"},
		},
		{
			name: "every monitor was skipped",
			e: monitorEnumeration{
				failures: []string{"GetMonitorInfo failed for monitor 0x1"},
			},
			wantErr:  true,
			contains: []string{"monitor 0x1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.e.incompleteError()
			if (err != nil) != tt.wantErr {
				t.Fatalf("incompleteError() = %v, want error: %v", err, tt.wantErr)
			}
			for _, want := range tt.contains {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err, want)
				}
			}
		})
	}
}

// monitorAt returns monitor info for a 1080 pixel high monitor spanning the
// given horizontal range.
func monitorAt(left, right int32) MONITORINFOEX {
	return MONITORINFOEX{MONITORINFO: MONITORINFO{RcMonitor: RECT{Left: left, Right: right, Bottom: 1080}}}
}

func TestLastErrorSuffix(t *testing.T) {
	if got := lastErrorSuffix(syscall.Errno(0)); got != "" {
		t.Errorf("zero errno: got %q, want empty", got)
	}
	if got := lastErrorSuffix(nil); got != "" {
		t.Errorf("nil: got %q, want empty", got)
	}
	if got := lastErrorSuffix(syscall.Errno(5)); got == "" {
		t.Error("non-zero errno: got empty, want the error text")
	}
}
