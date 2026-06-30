package setupwizard

import "testing"

func TestParseGPGSecretKeys(t *testing.T) {
	// Representative `gpg --list-secret-keys --with-colons --keyid-format long`
	// output: two keys, the second UID contains an escaped colon (\x3a).
	sample := `sec:u:4096:1:ABCDEF0123456789:1700000000:::u:::scESC:::+:::23::0:
fpr:::::::::1111111111111111ABCDEF0123456789:
grp:::::::::AAAA:
uid:u::::1700000000::HASH1::Lea Anthony <lea@wails.io>::::::::::0:
ssb:u:4096:1:1111:1700000000::::::e:::+:::23:
sec:u:255:22:FEDCBA9876543210:1700000100:::u:::scESC:::+:::ed25519:::0:
fpr:::::::::2222222222222222FEDCBA9876543210:
uid:u::::1700000100::HASH2::Test User (work\x3a ci) <ci@example.com>::::::::::0:
`
	keys := parseGPGSecretKeys([]byte(sample))
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %+v", len(keys), keys)
	}
	if keys[0].KeyID != "ABCDEF0123456789" {
		t.Errorf("key0 ID = %q", keys[0].KeyID)
	}
	if keys[0].UID != "Lea Anthony <lea@wails.io>" {
		t.Errorf("key0 UID = %q", keys[0].UID)
	}
	if keys[1].KeyID != "FEDCBA9876543210" {
		t.Errorf("key1 ID = %q", keys[1].KeyID)
	}
	if keys[1].UID != "Test User (work: ci) <ci@example.com>" {
		t.Errorf("key1 UID = %q (colon unescape failed?)", keys[1].UID)
	}
}

func TestParseGPGSecretKeys_Empty(t *testing.T) {
	if got := parseGPGSecretKeys(nil); got != nil {
		t.Errorf("expected nil for empty input, got %+v", got)
	}
}

func TestGPGFieldClean(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"Lea Anthony", "Lea Anthony", true},
		{"  spaced  ", "spaced", true},
		{"", "", false},
		{"line\nbreak", "", false},
		{"%commit", "", false},        // batch-directive injection
		{"carriage\rreturn", "", false},
	}
	for _, c := range cases {
		got, ok := gpgFieldClean(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("gpgFieldClean(%q) = (%q,%v), want (%q,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}
