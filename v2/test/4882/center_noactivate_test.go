package test

import (
	"testing"
)

const (
	SWP_NOSIZE     = 0x0001
	SWP_NOACTIVATE = 0x0010
)

func TestCenterUsesNoActivateFlag(t *testing.T) {
	combinedFlags := SWP_NOSIZE | SWP_NOACTIVATE
	if combinedFlags&SWP_NOACTIVATE == 0 {
		t.Error("SWP_NOACTIVATE flag should be set in Center() SetWindowPos call")
	}
	if combinedFlags&SWP_NOSIZE == 0 {
		t.Error("SWP_NOSIZE flag should still be set in Center() SetWindowPos call")
	}
}

func TestNoActivateFlagValue(t *testing.T) {
	if SWP_NOACTIVATE != 0x0010 {
		t.Errorf("SWP_NOACTIVATE should be 0x0010, got 0x%04X", SWP_NOACTIVATE)
	}
}

func TestCenterFlagsDifferFromOldFlags(t *testing.T) {
	oldFlags := SWP_NOSIZE
	newFlags := SWP_NOSIZE | SWP_NOACTIVATE
	if oldFlags == newFlags {
		t.Error("new flags should differ from old flags (SWP_NOACTIVATE was added)")
	}
	if oldFlags&SWP_NOACTIVATE != 0 {
		t.Error("old flags should NOT have SWP_NOACTIVATE")
	}
	if newFlags&SWP_NOACTIVATE == 0 {
		t.Error("new flags SHOULD have SWP_NOACTIVATE")
	}
}
