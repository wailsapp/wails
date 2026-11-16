package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNotaryAPICredentialsRoundTrip(t *testing.T) {
	root := t.TempDir()
	source := `version: 3
project:
  name: notarized
  product_name: Notarized
  identifier: com.example.notarized
  version: 1.0.0
darwin:
  signing:
    identity: Developer ID Application
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
`
	path := filepath.Join(root, Filename)
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeConfig(loaded.Config)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, err := decodeYAML(root, "round-trip.yaml", encoded, "")
	if err != nil || roundTrip.Config.Signing.Darwin != loaded.Config.Signing.Darwin {
		t.Fatal("API credentials did not round trip", err)
	}
	for _, invalid := range []string{
		strings.Replace(source, "    key_id: ABCDEFGHIJ\n", "", 1),
		source + "    credential: profile\n",
		strings.Replace(source, "signing/AuthKey.p8", "../outside.p8", 1),
	} {
		if _, err := decodeYAML(root, "invalid.yaml", []byte(invalid), ""); err == nil {
			t.Fatal("invalid API configuration accepted")
		}
	}
}
