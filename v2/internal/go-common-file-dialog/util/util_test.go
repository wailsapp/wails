package util

import (
	"github.com/go-ole/go-ole"
	"testing"
)

func TestStringToUUID(t *testing.T) {
	generated := *StringToUUID("TestTestTest")
	expected := *ole.NewGUID("7933985F-2C87-5A5B-A26E-5D0326829AC2")
	if generated != expected {
		t.Errorf("not equal. expected %s, found %s", expected.String(), generated.String())
	}
}
