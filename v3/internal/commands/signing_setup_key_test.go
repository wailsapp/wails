package commands

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGeneratedSigningKeysRoundTrip(t *testing.T) {
	for _, test := range []struct{ name, password string }{{"plaintext", ""}, {"encrypted", "test-only-setup-passphrase"}} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			privateKey, publicKey := filepath.Join(root, "private.asc"), filepath.Join(root, "public.asc")
			if err := generatePGPKeyForSetup("Signing setup fixture", "fixture@example.invalid", test.password, privateKey, publicKey); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(privateKey)
			if err != nil {
				t.Fatal(err)
			}
			input := filepath.Join(root, "fixture.pkg.tar.zst")
			if err := os.WriteFile(input, []byte("native Go signing setup round-trip payload"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := signArchPackage(t.Context(), input, privateKey, test.password); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(privateKey)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("signing mutated exported private key: %v", err)
			}
			signature, err := os.ReadFile(input + ".sig")
			if err != nil || len(signature) == 0 {
				t.Fatalf("missing native Go signature: %v", err)
			}
			if test.password != "" {
				if err := signArchPackage(t.Context(), input, privateKey, "incorrect"); err == nil {
					t.Fatal("encrypted exported key accepted the wrong passphrase")
				}
			}
			gpg, gpgErr := exec.LookPath("gpg")
			gpgv, gpgvErr := exec.LookPath("gpgv")
			if gpgErr == nil && gpgvErr == nil {
				keyring := filepath.Join(root, "public.gpg")
				dearmor := exec.Command(gpg, "--batch", "--dearmor", "--output", filepath.Base(keyring), filepath.Base(publicKey))
				dearmor.Dir = root
				if output, err := dearmor.CombinedOutput(); err != nil {
					t.Fatalf("decode exported public key: %s, %v", output, err)
				}
				verify := exec.Command(gpgv, "--homedir", ".", "--keyring", "./"+filepath.Base(keyring), filepath.Base(input)+".sig", filepath.Base(input))
				verify.Dir = root
				output, err := verify.CombinedOutput()
				if err != nil {
					t.Fatalf("independent generated-key signature verification: %s, %v", output, err)
				}
				t.Logf("Independent gpgv verification: %s", output)
			}
		})
	}
}
