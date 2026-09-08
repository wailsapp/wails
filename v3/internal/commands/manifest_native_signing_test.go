package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestHCLNativeSigningSettingsReachWindowsAndDarwinArtifacts(t *testing.T) {
	t.Setenv("WINDOWS_PASSWORD", "test-pfx-password")
	tests := []struct {
		name, platform, format, hcl string
	}{
		{
			name:     "windows",
			platform: "windows",
			format:   "msix",
			hcl: `version = 3
project {
  name = "windows"
  product_name = "Windows"
  company = "Example"
  identifier = "com.example.windows"
  version = "1.0.0"
}
windows {
  signing {
    credential = "WINDOWS_PASSWORD"
    identity = "Example Publisher"
    certificate = "release.pfx"
    thumbprint = "ABC123"
    timestamp_server = "https://timestamp.example.test"
  }
}
profile "release" {
  target "windows/amd64" {
    formats = ["msix"]
    sign = true
  }
}
`,
		},
		{
			name:     "darwin",
			platform: "darwin",
			format:   "app",
			hcl: `version = 3
project {
  name = "darwin"
  product_name = "Darwin"
  identifier = "com.example.darwin"
  version = "1.0.0"
}
darwin {
  signing {
    credential = "MACOS_KEYCHAIN"
    identity = "Developer ID Application: Example"
    entitlements = "release.entitlements"
  }
  notarization {
    credential = "NOTARY_PROFILE"
  }
}
profile "release" {
  target "darwin/arm64" {
    sign = true
    notarize = true
  }
}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(test.hcl), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(root, "release.pfx"), []byte("certificate"), 0o644))
			require.NoError(t, os.WriteFile(filepath.Join(root, "release.entitlements"), []byte("entitlements"), 0o644))
			loaded, err := manifest.Load(root, "release")
			require.NoError(t, err)
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
			require.NoError(t, err)
			signKey := pipeline.NodeKey("package:windows/amd64:msix:sign")
			if test.platform == "darwin" {
				signKey = "assemble:darwin/arm64:sign"
			}
			signNode := plan.Nodes[signKey]
			spec := signNode.Spec.(pipeline.SignSpec)
			input := filepath.Join(root, spec.Input)
			if test.platform == "darwin" {
				require.NoError(t, os.MkdirAll(filepath.Join(input, "Contents"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(input, "Contents", "Info.plist"), []byte("plist"), 0o644))
			} else {
				require.NoError(t, os.MkdirAll(filepath.Dir(input), 0o755))
				require.NoError(t, os.WriteFile(input, []byte("package"), 0o644))
			}

			var received flags.Sign
			previousSign := manifestSign
			manifestSign = func(options *flags.Sign) error {
				received = *options
				return nil
			}
			t.Cleanup(func() { manifestSign = previousSign })
			handler := &manifestHandler{root: root, config: loaded.Config}
			if spec.Config.Entitlements != "" {
				assetsRoot := filepath.Dir(filepath.Dir(filepath.Join(root, spec.Config.Entitlements)))
				relativeAssets, relativeErr := filepath.Rel(root, assetsRoot)
				require.NoError(t, relativeErr)
				require.NoError(t, handler.applyUserSigningInputs(filepath.Join(root, relativeAssets), test.platform))
				require.NoError(t, os.WriteFile(filepath.Join(root, "release.entitlements"), []byte("mutated"), 0o644))
			}
			_, err = handler.Run(t.Context(), signNode)
			require.NoError(t, err)
			assert.True(t, spec.Config.Enabled)
			assert.NotEqual(t, input, received.Input, "signing must operate on a private staged copy")
			assert.Contains(t, received.Input, string(filepath.Separator)+".sign-output-")
			assert.Equal(t, spec.Config.Identity, received.Identity)
			expectedCertificate := spec.Config.Certificate
			if expectedCertificate != "" {
				expectedCertificate, err = filepath.EvalSymlinks(filepath.Join(root, expectedCertificate))
				require.NoError(t, err)
			}
			assert.Equal(t, expectedCertificate, received.Certificate)
			assert.Equal(t, spec.Config.Thumbprint, received.Thumbprint)
			assert.Equal(t, spec.Config.TimestampServer, received.Timestamp)
			expectedEntitlements := spec.Config.Entitlements
			if expectedEntitlements != "" {
				expectedEntitlements, err = filepath.EvalSymlinks(filepath.Join(root, expectedEntitlements))
				require.NoError(t, err)
			}
			assert.Equal(t, expectedEntitlements, received.Entitlements)
			if test.platform == "darwin" {
				assert.True(t, received.Notarize)
				assert.Equal(t, "NOTARY_PROFILE", received.KeychainProfile)
				assert.DirExists(t, input+".signed")
			} else {
				assert.FileExists(t, input+".signed")
				assert.Equal(t, "test-pfx-password", received.Password)
				t.Setenv("WINDOWS_PASSWORD", "")
				manifestSign = func(*flags.Sign) error {
					t.Fatal("missing credentials must fail before invoking the signer")
					return nil
				}
				_, err = handler.Run(t.Context(), signNode)
				require.ErrorContains(t, err, "WINDOWS_PASSWORD is empty")
			}
		})
	}
}
