package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// TestResolveSigningDefaults verifies that signing options unset via flags are
// backfilled from the global config for every platform, and that explicit flags
// take precedence.
func TestResolveSigningDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, "wails")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := `signing:
  darwin:
    identity: "Developer ID Application: Test (T123)"
    keychainProfile: notary
    entitlements: ent.plist
  windows:
    certificatePath: /tmp/cert.pfx
    thumbprint: ABC123
    timestampServer: http://ts.example
  linux:
    gpgKeyPath: /tmp/key.asc
    signRole: maint
`
	if err := os.WriteFile(filepath.Join(cfgDir, "defaults.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	// Empty options → everything backfilled.
	opts := &flags.Sign{}
	resolveSigningDefaults(opts)
	checks := map[string]struct{ got, want string }{
		"Identity":        {opts.Identity, "Developer ID Application: Test (T123)"},
		"KeychainProfile": {opts.KeychainProfile, "notary"},
		"Entitlements":    {opts.Entitlements, "ent.plist"},
		"Certificate":     {opts.Certificate, "/tmp/cert.pfx"},
		"Thumbprint":      {opts.Thumbprint, "ABC123"},
		"Timestamp":       {opts.Timestamp, "http://ts.example"},
		"PGPKey":          {opts.PGPKey, "/tmp/key.asc"},
		"Role":            {opts.Role, "maint"},
	}
	for name, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", name, c.got, c.want)
		}
	}

	// Explicit flags must win over the config.
	explicit := &flags.Sign{PGPKey: "/explicit/key.asc", Certificate: "/explicit/cert.pfx"}
	resolveSigningDefaults(explicit)
	if explicit.PGPKey != "/explicit/key.asc" {
		t.Errorf("flag PGPKey overridden by config: %q", explicit.PGPKey)
	}
	if explicit.Certificate != "/explicit/cert.pfx" {
		t.Errorf("flag Certificate overridden by config: %q", explicit.Certificate)
	}
}

// TestResolveSigningDefaults_GPGKeyIDOnly documents a known limitation: a config
// with only a key *ID* (no exported key path) does not yield a usable PGPKey,
// because the signer consumes a key file. Keys generated before the wizard's
// export feature are in this state and must be re-exported.
func TestResolveSigningDefaults_GPGKeyIDOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, "wails")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "signing:\n  linux:\n    gpgKeyID: 70016B305B5DB108\n"
	if err := os.WriteFile(filepath.Join(cfgDir, "defaults.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &flags.Sign{}
	resolveSigningDefaults(opts)
	if opts.PGPKey != "" {
		t.Errorf("expected PGPKey empty when only gpgKeyID is set, got %q", opts.PGPKey)
	}
}
