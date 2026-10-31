package cfd

import (
	"errors"
	"testing"

	"github.com/go-ole/go-ole"
)

func TestHRESULTClassification(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  uint64
	}{
		{"S_OK", 0},
		{"S_FALSE", 1},
		{"positive status", 0x7fffffff},
		{"failure boundary", 0x80000000},
		{"E_FAIL", 0x80004005},
		{"E_ACCESSDENIED", 0x80070005},
		{"zero extended failure", 0x0000000080004005},
		{"sign extended failure", 0xffffffff80004005},
		{"upper bits with success", 0xffffffff00000001},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := hresultToError(uintptr(tc.raw))
			code := uint32(tc.raw)
			if code < 0x80000000 {
				if err != nil {
					t.Fatalf("success returned %T", err)
				}
				return
			}
			var failure *ole.OleError
			if !errors.As(err, &failure) {
				t.Fatalf("HRESULT %#x returned %T, want OleError", tc.raw, err)
			}
			if failure.Code() != uintptr(code) {
				t.Errorf("code = %#x, want %#x", failure.Code(), code)
			}
		})
	}
}

func TestHRESULTCancellation(t *testing.T) {
	for _, raw := range []uint64{0x800704c7, 0xffffffff800704c7} {
		if err := hresultToError(uintptr(raw)); !errors.Is(err, ErrorCancelled) {
			t.Errorf("cancellation %#x returned %T, want ErrorCancelled", raw, err)
		}
	}
}
