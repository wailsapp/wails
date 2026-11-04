//go:build windows

package cfd

import (
	"errors"
	"syscall"
	"testing"

	"github.com/go-ole/go-ole"
)

func TestShowPropagatesHRESULT(t *testing.T) {
	for _, code := range []uint32{0, 0x80004005, 0x800704c7} {
		vtbl := iModalWindowVtbl{Show: syscall.NewCallback(func(this, owner uintptr) uintptr {
			return uintptr(code)
		})}
		err := vtbl.show(nil, 0)
		switch code {
		case 0:
			if err != nil {
				t.Fatalf("success returned %T", err)
			}
		case 0x800704c7:
			if !errors.Is(err, ErrorCancelled) {
				t.Fatalf("cancel returned %T", err)
			}
		default:
			var failure *ole.OleError
			if !errors.As(err, &failure) || failure.Code() != uintptr(code) {
				t.Fatalf("failure %#x returned %T", code, err)
			}
		}
	}
}

func TestReleaseUsesReferenceCount(t *testing.T) {
	calls := 0
	vtbl := iUnknownVtbl{Release: syscall.NewCallback(func(this uintptr) uintptr {
		calls++
		return 0x80000000
	})}
	if err := vtbl.release(nil); err != nil {
		t.Fatalf("reference count was interpreted as %T", err)
	}
	if calls != 1 {
		t.Fatalf("Release called %d times, want once", calls)
	}
}
