//go:build windows || darwin

package hostinfo

import (
	"os"
	"runtime"
	"testing"
)

func TestNativeHardware(t *testing.T) {
	if os.Getenv("WAILS_TEST_HOSTINFO") != "1" {
		t.Skip("set WAILS_TEST_HOSTINFO=1 to query the native hardware providers")
	}
	processors, err := Processors()
	if err != nil || len(processors) == 0 {
		t.Fatalf("processors: %+v, %v", processors, err)
	}
	for _, p := range processors {
		if p.Model == "" || p.Cores < 1 {
			t.Fatalf("invalid processor: %+v", p)
		}
	}
	memory, err := Memory()
	if err != nil || memory == 0 {
		t.Fatalf("memory: %d, %v", memory, err)
	}
	if runtime.GOOS == "windows" {
		cards, err := Graphics()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("graphics: %+v", cards)
	}
	t.Logf("processors: %+v; memory: %d", processors, memory)
}
