package ipcauth

import "testing"

func TestTokenIsStableAndNonEmpty(t *testing.T) {
	a, b := Token(), Token()
	if a == "" {
		t.Fatal("token is empty")
	}
	if a != b {
		t.Fatalf("token is not stable across calls: %q vs %q", a, b)
	}
}

func TestValidAcceptsOnlyTheToken(t *testing.T) {
	if !Valid(Token()) {
		t.Error("the real token was rejected")
	}
	for _, bad := range []string{"", "nope", Token() + "x", Token()[:len(Token())-1]} {
		if Valid(bad) {
			t.Errorf("an invalid capability %q was accepted", bad)
		}
	}
}
