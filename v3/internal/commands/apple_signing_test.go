package commands

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/applesign"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/notaryservice"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestPortableAppleSigningCommandAndWakePreserveInputs(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, artifact := range []struct{ name, fixture string }{{"program", "arm64.macho"}, {"image.dmg", "unsigned.dmg"}} {
		for _, certificate := range []bool{false, true} {
			t.Run(artifact.name+map[bool]string{false: "-adhoc", true: "-certificate"}[certificate], func(t *testing.T) {
				root := t.TempDir()
				original, err := os.ReadFile(filepath.Join("..", "applesign", "testdata", artifact.fixture))
				if err != nil {
					t.Fatal(err)
				}
				input := filepath.Join(root, artifact.name)
				if err := os.WriteFile(input, original, 0755); err != nil {
					t.Fatal(err)
				}
				config := manifest.SigningPlatform{Enabled: true, Identity: "-"}
				options := &flags.Sign{Input: input, Output: input + ".output", Identity: "-"}
				if certificate {
					key, err := os.ReadFile(filepath.Join("..", "codesignature", "testdata", "modern.pfx"))
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(root, "identity.pfx"), key, 0600); err != nil {
						t.Fatal(err)
					}
					options.Certificate, options.Password = filepath.Join(root, "identity.pfx"), "fixture-only"
					config.Certificate, config.Credential = "identity.pfx", "APPLE_TEST_PASSWORD"
				}
				if err := Sign(options); err != nil {
					t.Fatal(err)
				}
				assertAppleSignedArtifact(t, options.Output)
				handler := &manifestHandler{root: root, environment: []string{"APPLE_TEST_PASSWORD=fixture-only"}}
				format := "binary"
				if artifact.fixture == "unsigned.dmg" {
					format = "dmg"
				}
				if _, err := handler.sign(t.Context(), pipeline.SignSpec{TargetOS: "darwin", Format: format, BundleSigned: true, Input: artifact.name, Config: config}); err != nil {
					t.Fatal(err)
				}
				assertAppleSignedArtifact(t, input+".signed")
				after, err := os.ReadFile(input)
				if err != nil || !bytes.Equal(original, after) {
					t.Fatal("signing changed input", err)
				}
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				before, err := os.ReadFile(options.Output)
				if err != nil {
					t.Fatal(err)
				}
				if err := signPortableApple(ctx, options); !errors.Is(err, context.Canceled) {
					t.Fatal("cancellation lost", err)
				}
				after, err = os.ReadFile(options.Output)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatal("failure replaced signed output", err)
				}
			})
		}
	}
}

func assertAppleSignedArtifact(t *testing.T, path string) {
	t.Helper()
	identities, err := applesign.Inspect(path)
	if err != nil || len(identities) == 0 {
		t.Fatal("signed artifact cannot be inspected", err)
	}
	if runtime.GOOS == "darwin" {
		if output, err := exec.Command("/usr/bin/codesign", "--verify", "--strict", "--all-architectures", path).CombinedOutput(); err != nil {
			t.Fatalf("native verifier: %v: %s", err, output)
		}
	}
}

func TestAppleInputDoesNotInheritWindowsCertificateDefaults(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	if err := os.Mkdir(filepath.Join(root, "wails"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wails", "defaults.yaml"), []byte("signing:\n  windows:\n    certificatePath: windows.pfx\n    timestampServer: https://windows.example\n    thumbprint: WINDOWS\n"), 0600); err != nil {
		t.Fatal(err)
	}
	options := &flags.Sign{Input: "Application.app", Identity: "-"}
	resolveSigningDefaults(options)
	if options.Certificate != "" || options.Timestamp != "" || options.Thumbprint != "" {
		t.Fatal("Windows defaults leaked into Apple signing", options)
	}
}

func TestAppleSigningDefaultsReuseConfiguredIdentityAndAPIKey(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	if err := os.Mkdir(filepath.Join(root, "wails"), 0700); err != nil {
		t.Fatal(err)
	}
	config := "signing:\n  darwin:\n    p12Path: apple.p12\n    apiKeyPath: api.p8\n    apiKeyID: KEYID\n    apiIssuerID: ISSUER\n"
	if err := os.WriteFile(filepath.Join(root, "wails", "defaults.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	options := &flags.Sign{Input: "Application.app"}
	resolveSigningDefaults(options)
	if options.Certificate != "apple.p12" || options.NotaryKey != "api.p8" || options.NotaryKeyID != "KEYID" || options.NotaryIssuerID != "ISSUER" {
		t.Fatal("configured Apple identity was lost", options)
	}
	explicit := &flags.Sign{Input: "Application.app", Identity: "-", KeychainProfile: "explicit-profile"}
	resolveSigningDefaults(explicit)
	if explicit.Certificate != "" || usesNotaryAPI(explicit) {
		t.Fatal("defaults overrode explicit signer or notarization profile", explicit)
	}
}

func TestSignRoutesStandaloneAppleBundlesWithoutPlatformTools(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, extension := range []string{".framework", ".FRAMEWORK", ".xpc", ".bundle", ".appex"} {
		t.Run(extension, func(t *testing.T) {
			root := t.TempDir()
			input := filepath.Join(root, "Probe"+extension)
			output := filepath.Join(root, "Signed"+extension)
			binaryPath := filepath.Join("Contents", "MacOS", "Probe")
			infoPath := filepath.Join("Contents", "Info.plist")
			resourcePath := filepath.Join("Contents", "Resources", "owned.txt")
			fixture := "arm64.macho"
			if strings.EqualFold(extension, ".framework") {
				binaryPath = "Probe"
				infoPath = filepath.Join("Resources", "Info.plist")
				resourcePath = filepath.Join("Resources", "owned.txt")
				fixture = "framework.macho"
			}
			original, err := os.ReadFile(filepath.Join("..", "applesign", "testdata", fixture))
			if err != nil {
				t.Fatal(err)
			}
			info := []byte(`<?xml version="1.0"?><plist version="1.0"><dict><key>CFBundleExecutable</key><string>Probe</string><key>CFBundleIdentifier</key><string>org.wails.routingprobe</string></dict></plist>`)
			for _, entry := range []struct {
				name string
				data []byte
				mode os.FileMode
			}{{binaryPath, original, 0755}, {infoPath, info, 0644}, {resourcePath, []byte("owned resource"), 0644}} {
				path := filepath.Join(input, entry.name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, entry.data, entry.mode); err != nil {
					t.Fatal(err)
				}
			}
			if err := Sign(&flags.Sign{Input: input, Output: output, Identity: "-"}); err != nil {
				t.Fatal(err)
			}
			if identities, err := applesign.Inspect(filepath.Join(output, binaryPath)); err != nil || len(identities) != 1 {
				t.Fatalf("missing signed bundle executable: %v %v", identities, err)
			}
			seal := filepath.Join(output, "Contents", "_CodeSignature", "CodeResources")
			if strings.EqualFold(extension, ".framework") {
				seal = filepath.Join(output, "_CodeSignature", "CodeResources")
			}
			if data, err := os.ReadFile(seal); err != nil || len(data) == 0 {
				t.Fatalf("missing bundle resource seal: %v", err)
			}
			if after, err := os.ReadFile(filepath.Join(input, binaryPath)); err != nil || !bytes.Equal(original, after) {
				t.Fatal("input executable changed", err)
			}
			if after, err := os.ReadFile(filepath.Join(input, infoPath)); err != nil || !bytes.Equal(info, after) {
				t.Fatal("input metadata changed", err)
			}
			if runtime.GOOS == "darwin" {
				if log, err := exec.Command("/usr/bin/codesign", "--verify", "--strict", output).CombinedOutput(); err != nil {
					t.Fatalf("native bundle verification: %v: %s", err, log)
				}
			}
		})
	}
}

func init() {
	if os.Getenv("WAILS_TEST_NOTARY_PROFILE_HELPER") != "1" {
		return
	}
	if err := os.WriteFile(os.Getenv("WAILS_TEST_NOTARY_PROFILE_READY"), []byte(strings.Join(os.Args[1:], "\n")), 0600); err != nil {
		os.Exit(23)
	}
	time.Sleep(time.Minute)
	os.Exit(24)
}

func TestSubmitNotarizationProfileCancelsRunningTool(t *testing.T) {
	root := t.TempDir()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	name := "xcrun"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.WriteFile(filepath.Join(root, name), data, 0755); err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(root, "ready")
	t.Setenv("PATH", root)
	t.Setenv("WAILS_TEST_NOTARY_PROFILE_HELPER", "1")
	t.Setenv("WAILS_TEST_NOTARY_PROFILE_READY", ready)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	input := filepath.Join(root, "Probe.dmg")
	go func() { result <- submitNotarizationProfile(ctx, input, "test-profile") }()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	poll := time.NewTicker(10 * time.Millisecond)
	defer poll.Stop()
	for {
		select {
		case err := <-result:
			t.Fatalf("tool exited before cancellation: %v", err)
		case <-deadline.C:
			cancel()
			t.Fatal("notarytool did not start")
		case <-poll.C:
			arguments, err := os.ReadFile(ready)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				t.Fatal(err)
			}
			want := strings.Join([]string{"notarytool", "submit", input, "--keychain-profile", "test-profile", "--wait"}, "\n")
			if string(arguments) != want {
				t.Fatalf("wrong notarytool arguments: %s", arguments)
			}
			cancel()
			select {
			case err := <-result:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation lost: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("notarytool survived cancellation")
			}
			return
		}
	}
}

type uppercaseAppNotaryClient struct{ fixtureNotaryClient }

func (c *uppercaseAppNotaryClient) Submit(ctx context.Context, path string) (notaryservice.Submission, error) {
	c.archive = path
	archive, err := zip.OpenReader(path)
	if err != nil {
		return notaryservice.Submission{}, err
	}
	defer archive.Close()
	if len(archive.File) != 2 || archive.File[1].Name != "Example.APP/binary" {
		return notaryservice.Submission{}, errors.New("uppercase application was not archived with its bundle root")
	}
	return notaryservice.Submission{ID: "uppercase-app"}, nil
}

func TestNotarizationArchivesUppercaseAppExtension(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "Example.APP")
	key := filepath.Join(root, "AuthKey.p8")
	if err := os.Mkdir(app, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "binary"), []byte("owned test executable"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(key, []byte("test-only-key"), 0600); err != nil {
		t.Fatal(err)
	}
	client := &uppercaseAppNotaryClient{fixtureNotaryClient: fixtureNotaryClient{t: t, status: "Accepted"}}
	previousClient, previousStaple := newNotaryClient, stapleNotarization
	t.Cleanup(func() { newNotaryClient, stapleNotarization = previousClient, previousStaple })
	newNotaryClient = func(notaryservice.Credentials) (notaryClient, error) { return client, nil }
	stapled := false
	stapleNotarization = func(ctx context.Context, path string) error {
		if path != app {
			return errors.New("stapled upload archive instead of source application")
		}
		stapled = true
		return nil
	}
	if err := notarizeMacOSWithAPI(t.Context(), &flags.Sign{Input: app, NotaryKey: key, NotaryKeyID: "KEYID", NotaryIssuerID: "ISSUER"}); err != nil {
		t.Fatal(err)
	}
	if !stapled {
		t.Fatal("accepted app was not stapled")
	}
	if _, err := os.Stat(client.archive); !os.IsNotExist(err) {
		t.Fatal("temporary upload archive leaked")
	}
}
