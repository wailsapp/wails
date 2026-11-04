package signingconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"gopkg.in/yaml.v3"
)

func TestApplyDefaultsOnlyReadsSignedRequests(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, _ := defaults.GetDefaultsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("signing: [invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	config := manifest.Config{Root: t.TempDir()}
	if _, err := ApplyDefaults(config); err != nil {
		t.Fatal("unsigned build read malformed defaults", err)
	}
	if _, err := ApplyDefaults(config, "windows"); err == nil {
		t.Fatal("signed build ignored malformed defaults")
	}
}

func TestApplyDefaultsPreservesExplicitChoicesAndInput(t *testing.T) {
	global := defaults.SigningDefaults{}
	global.Windows.CertificatePath = filepath.Join(t.TempDir(), "global.pfx")
	global.Windows.Thumbprint = "global-thumbprint"
	global.Windows.TimestampServer = "https://timestamp.example.test"
	global.Darwin.Identity = "global-identity"
	global.Darwin.P12Path = filepath.Join(t.TempDir(), "global.p12")
	global.Darwin.KeychainProfile = "global-profile"
	global.Linux.GPGKeyPath = filepath.Join(t.TempDir(), "global.asc")
	original := manifest.Config{Root: t.TempDir()}
	original.Signing.Windows.Certificate = "project.pfx"
	original.Signing.Darwin.Certificate = "project.p12"
	original.Signing.Darwin.NotarizationKey = "project.p8"
	original.Signing.Darwin.NotarizationKeyID = "project-id"
	original.Signing.Darwin.NotarizationIssuerID = "project-issuer"
	before := original
	resolved, err := apply(original, map[string]bool{"windows": true, "darwin": true, "linux": true}, global)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, before) {
		t.Fatal("input config mutated")
	}
	if resolved.Signing.Windows.Certificate != "project.pfx" || resolved.Signing.Windows.Thumbprint != "" || resolved.Signing.Darwin.Certificate != "project.p12" || resolved.Signing.Darwin.Identity != "" || resolved.Signing.Darwin.NotarizationCredential != "" {
		t.Fatalf("explicit signing choices overwritten: %+v", resolved.Signing)
	}
	if resolved.Signing.Linux.Certificate != global.Linux.GPGKeyPath || !IsGlobalPath(resolved.Signing.Linux, "certificate", global.Linux.GPGKeyPath) {
		t.Fatal("missing trusted global key")
	}
	// Runtime-only origin data must neither be emitted nor accepted as manifest input.
	trusted := manifest.SigningPlatform{GlobalPaths: manifest.GlobalSigningPaths{Certificate: "secret-origin"}}
	for _, marshal := range []func(any) ([]byte, error){json.Marshal, yaml.Marshal} {
		data, err := marshal(trusted)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "secret-origin") || strings.Contains(string(data), "GlobalPaths") {
			t.Fatalf("serialized provenance: %s", data)
		}
	}
}

func TestSelectedProfileResolvesOnlyRequestedPlatform(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	global := defaults.Default()
	global.Signing.Windows.CertificatePath = filepath.Join(t.TempDir(), "outside.pfx")
	global.Signing.Darwin.P12Path = "unused.p12"
	if err := defaults.Save(global); err != nil {
		t.Fatal(err)
	}
	config := manifest.Config{Root: t.TempDir()}
	config.Selected = manifest.Profile{Name: "setup", Targets: []manifest.ProfileTarget{{Target: "windows/amd64", Sign: true}, {Target: "darwin/arm64"}}}
	got, err := ApplyDefaults(config)
	if err != nil {
		t.Fatal(err)
	}
	if got.Signing.Windows.Certificate != global.Signing.Windows.CertificatePath || got.Signing.Darwin.Certificate != "" || config.Signing.Windows.Certificate != "" {
		t.Fatal("incorrect platform selection or input mutation")
	}
}

func TestResolvePathRequiresExactGlobalOrigin(t *testing.T) {
	root, outside := t.TempDir(), filepath.Join(t.TempDir(), "key.pfx")
	if err := os.WriteFile(outside, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	config := manifest.SigningPlatform{Certificate: outside, GlobalPaths: manifest.GlobalSigningPaths{Certificate: outside}}
	if _, err := ResolvePath(root, config, "certificate", outside); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ field, path string }{{"entitlements", outside}, {"certificate", outside + ".other"}} {
		if _, err := ResolvePath(root, config, tc.field, tc.path); err == nil {
			t.Fatal("trusted different field or path")
		}
	}
	if _, err := ResolvePath(root, manifest.SigningPlatform{}, "certificate", outside); err == nil {
		t.Fatal("untrusted external file accepted")
	}
	inside := filepath.Join(root, "key.pfx")
	if err := os.WriteFile(inside, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolvePath(root, manifest.SigningPlatform{}, "certificate", inside); err != nil {
		t.Fatal("derived contained path rejected", err)
	}
	config.GlobalPaths.Certificate = filepath.Dir(outside)
	if _, err := ResolvePath(root, config, "certificate", filepath.Dir(outside)); err == nil {
		t.Fatal("directory accepted")
	}
}
