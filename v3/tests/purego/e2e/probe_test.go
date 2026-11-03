package main

import (
	"math"
	"testing"
)

// The probe's own logic, apart from any OS: these run on every platform.

func TestSolveMapping(t *testing.T) {
	for _, test := range []struct {
		name       string
		a, b       pointerReport
		ox, oy, sc float64
	}{
		// GTK3 menu bar 27 px tall, scale 1: page (120, 113) is raw (120, 140).
		{"menu bar", pointerReport{x: 120, y: 113}, pointerReport{x: 520, y: 313}, 0, 27, 1},
		// GDK_SCALE=2 with a 54 px (27 logical) bar.
		{"scale 2", pointerReport{x: 60, y: 43}, pointerReport{x: 260, y: 143}, 0, 54, 2},
		// Windows at 150%: client origin, raw = 1.5 * page.
		{"scale 1.5", pointerReport{x: 80, y: 93.333}, pointerReport{x: 346.667, y: 226.667}, 0, 0, 1.5},
	} {
		t.Run(test.name, func(t *testing.T) {
			m, err := solveMapping(120, 140, 520, 340, test.a, test.b)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(m.ox-test.ox) > 0.5 || math.Abs(m.oy-test.oy) > 0.5 || math.Abs(m.scale-test.sc) > 0.01 {
				t.Errorf("mapping = %+v, want origin (%v, %v) scale %v", m, test.ox, test.oy, test.sc)
			}
			// The mapping sends the page points back to the raw points.
			if x, y := m.raw(int(math.Round(test.a.x)), int(math.Round(test.a.y))); math.Abs(float64(x-120)) > 1 || math.Abs(float64(y-140)) > 1 {
				t.Errorf("raw(a) = (%d, %d), want (120, 140)", x, y)
			}
		})
	}
}

func TestSolveMappingRejectsNonsense(t *testing.T) {
	for name, pair := range map[string][2]pointerReport{
		"no motion":        {{x: 10, y: 10}, {x: 10, y: 10}},
		"backwards":        {{x: 400, y: 300}, {x: 10, y: 10}},
		"axes disagree":    {{x: 0, y: 0}, {x: 400, y: 100}},
		"stale first read": {{x: 520, y: 313}, {x: 520, y: 313}},
	} {
		if _, err := solveMapping(120, 140, 520, 340, pair[0], pair[1]); err == nil {
			t.Errorf("%s: solveMapping accepted %v", name, pair)
		}
	}
}

func TestScenarioSelection(t *testing.T) {
	names := func(ss []scenario) map[string]bool {
		out := map[string]bool{}
		for _, s := range ss {
			out[s.name] = true
		}
		return out
	}
	all := names(selected("linux"))
	for _, want := range []string{"calibrate", "menu-navigate", "type-text", "save-file"} {
		if !all[want] {
			t.Errorf("linux runs no %q", want)
		}
	}
	if names(selected("darwin"))["menu-navigate"] {
		t.Error("darwin runs menu-navigate, which it has no keyboard path for")
	}
	if !names(selected("windows"))["calibrate"] {
		t.Error("windows does not calibrate")
	}

	t.Setenv("E2E_SCENARIOS", "keybinding, type-text")
	t.Setenv("E2E_SKIP", "type-text")
	if got := names(selected("linux")); len(got) != 1 || !got["keybinding"] {
		t.Errorf("E2E_SCENARIOS and E2E_SKIP select %v, want only keybinding", got)
	}
}

func TestScenarioNamesAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range scenarios {
		if seen[s.name] {
			t.Errorf("two scenarios are called %q", s.name)
		}
		seen[s.name] = true
	}
	if scenarios[0].name != "calibrate" {
		t.Error("calibrate must run first: the pointer scenarios click through its mapping")
	}
}

func TestWindowsVirtualKey(t *testing.T) {
	for key, want := range map[string]uint16{
		"return": 0x0D, "escape": 0x1B, "down": 0x28, "delete": 0x2E,
		"f1": 0x70, "f10": 0x79, "f12": 0x7B, "a": 'A', "z": 'Z', "0": '0', "9": '9',
		"numlock": 0x90, "capslock": 0x14,
	} {
		if got, ok := windowsVirtualKey(key); !ok || got != want {
			t.Errorf("windowsVirtualKey(%q) = %#x, %v; want %#x", key, got, ok, want)
		}
	}
	// Punctuation depends on the layout: the driver asks Windows at run time.
	for _, key := range []string{"/", ";", "f13", "", "ctrl"} {
		if _, ok := windowsVirtualKey(key); ok {
			t.Errorf("windowsVirtualKey(%q) has a fixed code", key)
		}
	}
}
