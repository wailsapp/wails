//go:build linux && purego && !android && !server

package application

import (
	"sync/atomic"
	"testing"
)

func TestX11GrabRestoresErrorHandler(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[fail], func(t *testing.T) {
			current := uintptr(42)
			grabs := 0
			g := &x11GlobalShortcuts{lib: &gsX11Funcs{
				defaultRootWindow: func(uintptr) uintptr { return 1 },
				setErrorHandler:   func(handler uintptr) uintptr { old := current; current = handler; return old },
				grabKey:           func(uintptr, int32, uint32, uintptr, int32, int32, int32) int32 { grabs++; return 0 },
				sync: func(uintptr, int32) int32 {
					if fail {
						atomic.StoreInt32(&gsX11GrabError, 1)
					}
					return 0
				},
			}}
			got := g.grabKey(10, 0)
			if (got != 0) != fail {
				t.Fatalf("grab result = %d, failure = %v", got, fail)
			}
			if current != 42 || grabs != len(gsX11LockMasks) {
				t.Fatalf("handler = %d, grabs = %d", current, grabs)
			}
		})
	}
}
