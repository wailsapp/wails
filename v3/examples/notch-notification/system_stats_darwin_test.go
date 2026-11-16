//go:build darwin

package main

import "testing"

func TestReadSystemSample(t *testing.T) {
	first, ok := readSystemSample()
	if !ok {
		t.Fatal("readSystemSample failed")
	}
	if first.totalTicks == 0 || first.activeTicks > first.totalTicks {
		t.Fatalf("cpu ticks: active %d of %d", first.activeTicks, first.totalTicks)
	}
	if first.memoryTotal == 0 || first.memoryUsed == 0 || first.memoryUsed > first.memoryTotal {
		t.Fatalf("memory: used %d of %d", first.memoryUsed, first.memoryTotal)
	}
	if first.diskTotal == 0 || first.diskUsed == 0 || first.diskUsed > first.diskTotal {
		t.Fatalf("disk: used %d of %d", first.diskUsed, first.diskTotal)
	}
	second, ok := readSystemSample()
	if !ok || second.totalTicks < first.totalTicks {
		t.Fatalf("second sample: ok=%v ticks %d after %d", ok, second.totalTicks, first.totalTicks)
	}
}
