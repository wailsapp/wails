package sqlite

import "testing"

func TestStatementHandleFormat(t *testing.T) {
	for _, tc := range []struct {
		id   uint64
		want string
	}{
		{0, "0000000000000000"},
		{1, "0000000000000001"},
		{15, "000000000000000f"},
		{16, "0000000000000010"},
		{0x0123456789abcdef, "0123456789abcdef"},
		{^uint64(0), "ffffffffffffffff"},
	} {
		stmt := &Stmt{id: tc.id}
		got, err := stmt.MarshalText()
		if err != nil || string(got) != tc.want {
			t.Errorf("MarshalText(%x)=%q, %v; want %q", tc.id, got, err, tc.want)
		}
	}
}

var statementHandleSink []byte

func BenchmarkStatementHandleMarshalText(b *testing.B) {
	stmt := &Stmt{id: 0x123456789abcdef0}
	b.ReportAllocs()
	for b.Loop() {
		encoded, err := stmt.MarshalText()
		if err != nil {
			b.Fatal(err)
		}
		statementHandleSink = encoded
	}
}
