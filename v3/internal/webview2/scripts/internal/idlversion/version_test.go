package idlversion

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in      string
		want    Version
		wantErr bool
	}{
		{"1.0.2903.40", Version{1, 0, 2903, 40, ""}, false},
		{"1.0.515-prerelease", Version{1, 0, 515, 0, "prerelease"}, false},
		{"121.0.2277.83 stable", Version{121, 0, 2277, 83, "stable"}, false},
		{"1.0", Version{1, 0, 0, 0, ""}, false},
		{"1.2.3.4.5", Version{}, true},
		{"x.y.z", Version{}, true},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("Parse(%q) — expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q) — unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Parse(%q) — got %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		// Spec calls out 1.0.0 < 10.0.0 (lexical vs numeric).
		{"1.0.0", "10.0.0", -1},
		{"10.0.0", "1.0.0", 1},
		{"1.0.2903.40", "1.0.2903.40", 0},
		{"1.0.2903.41", "1.0.2903.40", 1},
		{"1.0.2739.15", "1.0.2903.40", -1},
		{"1.0.515-prerelease", "1.0.515-stable", 0}, // channel ignored
		{"1.0", "1.0.0.0", 0},
	}
	for _, c := range cases {
		got, err := Compare(c.a, c.b)
		if err != nil {
			t.Errorf("Compare(%q, %q) — error: %v", c.a, c.b, err)
			continue
		}
		if got != c.want {
			t.Errorf("Compare(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	v, _ := Parse("1.0.2903.40")
	if v.String() != "1.0.2903.40" {
		t.Errorf("String() = %q, want %q", v.String(), "1.0.2903.40")
	}
	v, _ = Parse("1.0.515 prerelease")
	if v.String() != "1.0.515.0 prerelease" {
		t.Errorf("String() = %q, want %q", v.String(), "1.0.515.0 prerelease")
	}
}
