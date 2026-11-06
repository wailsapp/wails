package uuid

import "testing"

var formattedUUID string

func BenchmarkUUIDString(b *testing.B) {
	id := UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	if got := id.String(); got != "00010203-0405-0607-0809-0a0b0c0d0e0f" {
		b.Fatalf("unexpected UUID: %q", got)
	}
	b.ReportAllocs()
	for b.Loop() {
		formattedUUID = id.String()
	}
}
