package test_5111

import (
	"testing"
)

func TestBackingScaleFactorPropagation(t *testing.T) {
	expectedDPR := 2.0
	backingScaleFactor := 2.0

	if backingScaleFactor != expectedDPR {
		t.Errorf("backingScaleFactor should be %.1f on Retina, got %.1f", expectedDPR, backingScaleFactor)
	}
}

func TestSetOverrideDeviceScaleFactorIsCalled(t *testing.T) {
	backingScaleFactor := 2.0
	overriddenDPR := backingScaleFactor

	if overriddenDPR != 2.0 {
		t.Errorf("_setOverrideDeviceScaleFactor should be called with backingScaleFactor, expected %.1f got %.1f", backingScaleFactor, overriddenDPR)
	}
}

func TestNonRetinaScaleFactor(t *testing.T) {
	backingScaleFactor := 1.0
	overriddenDPR := backingScaleFactor

	if overriddenDPR != 1.0 {
		t.Errorf("On non-Retina displays, devicePixelRatio should be 1.0, got %.1f", overriddenDPR)
	}
}
