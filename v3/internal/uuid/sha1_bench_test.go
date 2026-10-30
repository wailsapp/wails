package uuid

import "testing"

var namedUUID UUID

func BenchmarkNewSHA1(b *testing.B) {
	namespace := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	data := []byte("www.example.com")
	if got := NewSHA1(namespace, data).String(); got != "2ed6657d-e927-568b-95e1-2665a8aea6a2" {
		b.Fatalf("unexpected name-based UUID: %s", got)
	}
	b.ReportAllocs()
	for b.Loop() {
		namedUUID = NewSHA1(namespace, data)
	}
}
