//go:build wails_native && !wails_single_instance

package application

import (
	"strings"
	"testing"
)

func TestNativeSingleInstanceStubDisabled(t *testing.T) {
	manager, err := newSingleInstanceManager(nil, nil)
	if err != nil || manager != nil {
		t.Fatalf("nil options returned manager %v and error %v", manager, err)
	}

	manager, err = newSingleInstanceManager(nil, &SingleInstanceOptions{})
	if manager != nil {
		t.Fatalf("disabled single-instance support returned manager %v", manager)
	}
	if err == nil || !strings.Contains(err.Error(), "wails_single_instance") {
		t.Fatalf("disabled single-instance error = %v, want build-tag guidance", err)
	}
}
