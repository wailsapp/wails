package uuid

import (
	"regexp"
	"testing"
)

// uuidPattern matches the standard 8-4-4-4-12 hex format.
var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestNil(t *testing.T) {
	var zero UUID
	if Nil != zero {
		t.Errorf("Nil should be the zero UUID, got %v", Nil)
	}
}

func TestNilString(t *testing.T) {
	got := Nil.String()
	want := "00000000-0000-0000-0000-000000000000"
	if got != want {
		t.Errorf("Nil.String() = %q, want %q", got, want)
	}
}

func TestNew_NotNil(t *testing.T) {
	id := New()
	if id == Nil {
		t.Error("New() returned nil UUID")
	}
}

func TestNew_VersionBits(t *testing.T) {
	id := New()
	if v := id[6] >> 4; v != 4 {
		t.Errorf("New() version bits: got %d, want 4", v)
	}
}

func TestNew_VariantBits(t *testing.T) {
	id := New()
	// RFC 4122 variant: top two bits of byte 8 must be 10xxxxxx
	if id[8]&0xc0 != 0x80 {
		t.Errorf("New() variant bits: got 0x%02x, want 0x80 mask", id[8]&0xc0)
	}
}

func TestNew_Unique(t *testing.T) {
	a := New()
	b := New()
	if a == b {
		t.Error("Two calls to New() returned the same UUID")
	}
}

func TestNewSHA1_Deterministic(t *testing.T) {
	ns := New()
	data := []byte("hello world")
	a := NewSHA1(ns, data)
	b := NewSHA1(ns, data)
	if a != b {
		t.Errorf("NewSHA1 is not deterministic: %v != %v", a, b)
	}
}

func TestNewSHA1_DifferentNamespaces(t *testing.T) {
	ns1 := New()
	ns2 := New()
	data := []byte("same data")
	a := NewSHA1(ns1, data)
	b := NewSHA1(ns2, data)
	if a == b {
		t.Error("NewSHA1 with different namespaces returned the same UUID")
	}
}

func TestNewSHA1_DifferentData(t *testing.T) {
	ns := New()
	a := NewSHA1(ns, []byte("foo"))
	b := NewSHA1(ns, []byte("bar"))
	if a == b {
		t.Error("NewSHA1 with different data returned the same UUID")
	}
}

func TestNewSHA1_VersionBits(t *testing.T) {
	id := NewSHA1(Nil, []byte("test"))
	if v := id[6] >> 4; v != 5 {
		t.Errorf("NewSHA1() version bits: got %d, want 5", v)
	}
}

func TestNewSHA1_VariantBits(t *testing.T) {
	id := NewSHA1(Nil, []byte("test"))
	if id[8]&0xc0 != 0x80 {
		t.Errorf("NewSHA1() variant bits: got 0x%02x, want 0x80 mask", id[8]&0xc0)
	}
}

func TestString_Format(t *testing.T) {
	id := New()
	s := id.String()
	if !uuidPattern.MatchString(s) {
		t.Errorf("String() %q does not match UUID pattern", s)
	}
}

func TestString_SegmentLengths(t *testing.T) {
	id := New()
	s := id.String()
	parts := regexp.MustCompile(`-`).Split(s, -1)
	if len(parts) != 5 {
		t.Fatalf("String() has %d segments, want 5", len(parts))
	}
	want := []int{8, 4, 4, 4, 12}
	for i, p := range parts {
		if len(p) != want[i] {
			t.Errorf("segment %d: length %d, want %d", i, len(p), want[i])
		}
	}
}
