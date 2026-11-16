package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestNativeLinuxSigningRunsWithoutExternalTools(t *testing.T) {
	for _, format := range []string{"deb", "rpm"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			input, err := os.ReadFile(filepath.Join("..", "packagesign", "testdata", "unsigned."+format))
			if err != nil {
				t.Fatal(err)
			}
			key, err := os.ReadFile(filepath.Join("..", "packagesign", "testdata", "private.asc"))
			if err != nil {
				t.Fatal(err)
			}
			for name, data := range map[string][]byte{"input." + format: input, "key.asc": key} {
				if err := os.WriteFile(filepath.Join(root, name), data, 0600); err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("PATH", t.TempDir())
			handler := &manifestHandler{root: root, environment: []string{"TEST_SIGN_PASSWORD=package-test-password"}}
			spec := pipeline.SignSpec{TargetOS: "linux", Format: format, Input: "input." + format, Config: manifest.SigningPlatform{Enabled: true, Certificate: "key.asc", Credential: "TEST_SIGN_PASSWORD"}}
			if _, err := handler.sign(t.Context(), spec); err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(root, spec.Input+".signed")
			signed, err := os.ReadFile(output)
			if err != nil || bytes.Equal(input, signed) {
				t.Fatal("native signer did not produce a signature", err)
			}
			handler.environment = []string{"TEST_SIGN_PASSWORD=incorrect"}
			if _, err := handler.sign(t.Context(), spec); err == nil {
				t.Fatal("wrong key password accepted")
			}
			unchanged, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(signed, unchanged) {
				t.Fatal("failed signing replaced previous output", err)
			}
			original, err := os.ReadFile(filepath.Join(root, spec.Input))
			if err != nil || !bytes.Equal(input, original) {
				t.Fatal("signing changed input", err)
			}
			if err := Sign(&flags.Sign{Input: filepath.Join(root, spec.Input), PGPKey: filepath.Join(root, "key.asc"), PGPPassword: "package-test-password", Output: filepath.Join(root, "cli."+format)}); err != nil {
				t.Fatal("standalone native signer", err)
			}
		})
	}
}
