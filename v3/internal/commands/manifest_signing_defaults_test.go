package commands

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestGlobalSigningDefaultsSignOwnedFixtureOutsideProject(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("PATH", t.TempDir())
	t.Setenv("WAILS_PGP_PASSWORD", "package-test-password")
	root, external := t.TempDir(), t.TempDir()
	key, err := os.ReadFile("../packagesign/testdata/private.asc")
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile("../packagesign/testdata/unsigned.deb")
	if err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(external, "key.asc")
	for path, data := range map[string][]byte{keyPath: key, filepath.Join(root, "input.deb"): original} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	global := defaults.Default()
	global.Signing.Linux.GPGKeyPath = keyPath
	if err := defaults.Save(global); err != nil {
		t.Fatal(err)
	}
	resolved, err := signingconfig.ApplyDefaults(manifest.Config{Root: root}, "linux")
	if err != nil {
		t.Fatal(err)
	}
	handler := &manifestHandler{root: root}
	spec := pipeline.SignSpec{TargetOS: "linux", Format: "deb", Input: "input.deb", Config: resolved.Signing.Linux}
	if _, err := handler.sign(t.Context(), spec); err != nil {
		t.Fatal(err)
	}
	signed, err := os.ReadFile(filepath.Join(root, "input.deb.signed"))
	if err != nil || bytes.Equal(signed, original) {
		t.Fatal("missing signed output", err)
	}
	for path, want := range map[string][]byte{keyPath: key, filepath.Join(root, "input.deb"): original} {
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatal("signer changed source", path, err)
		}
	}
	spec.Config.Credential = "EXPLICIT_PASSWORD"
	handler.environment = []string{"EXPLICIT_PASSWORD=wrong"}
	if _, err := handler.sign(t.Context(), spec); err == nil {
		t.Fatal("global password overrode explicit incorrect password")
	}
	got, err := os.ReadFile(filepath.Join(root, "input.deb.signed"))
	if err != nil || !bytes.Equal(got, signed) {
		t.Fatal("failed signing replaced prior output", err)
	}
}

type globalSigningInspectionHandler struct{}

func (globalSigningInspectionHandler) Identity(context.Context, pipeline.Node) (string, error) {
	return "fixture", nil
}
func (globalSigningInspectionHandler) Run(context.Context, pipeline.Node) (pipeline.RunResult, error) {
	return pipeline.RunResult{}, nil
}

func TestGlobalEntitlementsParticipateInPrebuiltPlanAndCache(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	key := filepath.Join(t.TempDir(), "identity.p12")
	entitlements := filepath.Join(t.TempDir(), "entitlements.plist")
	source := []byte(`version: 3
project:
  name: app
  product_name: App
  identifier: com.example.app
  version: 1.0.0
profiles:
  setup:
    targets:
      darwin/arm64:
        sign: true
`)
	for path, data := range map[string][]byte{filepath.Join(root, "wails.yaml"): source, key: []byte("fixture-key"), entitlements: []byte("first-entitlements")} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	global := defaults.Default()
	global.Signing.Darwin.P12Path = key
	global.Signing.Darwin.Entitlements = entitlements
	if err := defaults.Save(global); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(root, "setup")
	if err != nil {
		t.Fatal(err)
	}
	_, resolved, plan, err := resolveManifestPlanWithOperations(manifestRunOptions{Verb: "build", Loaded: loaded, SkipBuild: true}, manifestPlanOperations{
		getwd: func() (string, error) { return root, nil },
		plan: func(config manifest.Config, request pipeline.Request) (pipeline.Plan, error) {
			return pipeline.PlanBuildForHost(config, request, pipeline.NewHostCapabilities("linux", "amd64", nil, nil))
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Signing.Darwin.Certificate != "" {
		t.Fatal("caller config changed")
	}
	var assetKey pipeline.NodeKey
	for k, node := range plan.Nodes {
		for _, input := range node.Inputs {
			for _, path := range input.Files {
				if path == key {
					t.Fatal("private key included in cached/staged assets")
				}
				if path == entitlements {
					assetKey = k
				}
			}
		}
		if node.Kind == pipeline.CompileApplication {
			t.Fatal("prebuilt plan compiles")
		}
	}
	if assetKey == "" {
		t.Fatal("global entitlements missing from plan inputs")
	}
	executor := pipeline.Executor{Handler: globalSigningInspectionHandler{}}
	first, err := executor.Inspect(t.Context(), plan, root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entitlements, []byte("other-entitlements"), 0600); err != nil {
		t.Fatal(err)
	}
	second, err := executor.Inspect(t.Context(), plan, root)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(first.Operations[assetKey].Inputs, second.Operations[assetKey].Inputs) {
		t.Fatal("external entitlements content omitted from cache key")
	}
	stage := filepath.Join(root, ".wails", "test-assets")
	handler := &manifestHandler{root: root, config: resolved.Config}
	if err := handler.applyUserSigningInputs(stage, "darwin"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(stage, "signing", "entitlements.plist"))
	if err != nil || string(got) != "other-entitlements" {
		t.Fatal("global entitlements not staged", err)
	}
	after, err := os.ReadFile(filepath.Join(root, "wails.yaml"))
	if err != nil || !bytes.Equal(after, source) {
		t.Fatal("project YAML changed", err)
	}
}

func TestGlobalSigningPasswordPrecedence(t *testing.T) {
	t.Setenv("WAILS_WINDOWS_CERT_PASSWORD", "global-password")
	path := filepath.Join(t.TempDir(), "fixture.pfx")
	config := manifest.SigningPlatform{Certificate: path, GlobalPaths: manifest.GlobalSigningPaths{Certificate: path}}
	got, err := manifestSigningPassword(config, "windows", nil)
	if err != nil || got != "global-password" {
		t.Fatal(got, err)
	}
	config.Credential = "PROJECT_PASSWORD"
	got, err = manifestSigningPassword(config, "windows", []string{"PROJECT_PASSWORD=explicit"})
	if err != nil || got != "explicit" {
		t.Fatal(got, err)
	}
	config.Credential = ""
	config.GlobalPaths = manifest.GlobalSigningPaths{}
	got, err = manifestSigningPassword(config, "windows", nil)
	if err != nil || got != "" {
		t.Fatal("project key inherited unrelated global password", got, err)
	}
}

func TestGlobalSigningDefaultsFollowCompleteProfileIntent(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	if err := manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	loaded.Config.Selected = manifest.Profile{Name: "unsigned", Targets: []manifest.ProfileTarget{{Target: "darwin/arm64"}}}
	configPath, _ := defaults.GetDefaultsPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte("signing: [invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	operations := manifestPlanOperations{getwd: func() (string, error) { return root, nil }, plan: pipeline.PlanBuild}
	_, _, plan, err := resolveManifestPlanWithOperations(manifestRunOptions{Verb: "sign", Loaded: loaded, SkipBuild: true}, operations)
	if err != nil {
		t.Fatal("complete unsigned profile read signing defaults", err)
	}
	for _, node := range plan.Nodes {
		if node.Kind == pipeline.SignArtifact {
			t.Fatal("verb overrode unsigned profile")
		}
	}
	global := defaults.Default()
	global.Signing.Darwin.P12Path = filepath.Join(t.TempDir(), "global.p12")
	if err := defaults.Save(global); err != nil {
		t.Fatal(err)
	}
	loaded.Config.Selected = manifest.Profile{}
	_, resolved, plan, err := resolveManifestPlanWithOperations(manifestRunOptions{Verb: "sign", Loaded: loaded, SkipBuild: true, TargetOS: "darwin", TargetArch: "arm64"}, operations)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Config.Signing.Darwin.Certificate != global.Signing.Darwin.P12Path {
		t.Fatal("anonymous sign missed defaults")
	}
	found := false
	for _, node := range plan.Nodes {
		found = found || node.Kind == pipeline.SignArtifact
	}
	if !found {
		t.Fatal("anonymous sign omitted signer")
	}
}

func TestGlobalSigningDoesNotReadUnusedNotarizationKey(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("PATH", t.TempDir())
	root := t.TempDir()
	keyPath := filepath.Join(t.TempDir(), "fixture.p12")
	for path, fixture := range map[string]string{keyPath: "../codesignature/testdata/modern.pfx", filepath.Join(root, "program"): "../applesign/testdata/arm64.macho"} {
		data, err := os.ReadFile(fixture)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0700); err != nil {
			t.Fatal(err)
		}
	}
	global := defaults.Default()
	global.Signing.Darwin.P12Path = keyPath
	global.Signing.Darwin.APIKeyPath = filepath.Join(t.TempDir(), "not-needed.p8")
	global.Signing.Darwin.APIKeyID = "unused-key-id"
	global.Signing.Darwin.APIIssuerID = "unused-issuer"
	if err := defaults.Save(global); err != nil {
		t.Fatal(err)
	}
	config := manifest.Config{Root: root}
	config.Signing.Darwin.Credential = "APPLE_TEST_PASSWORD"
	resolved, err := signingconfig.ApplyDefaults(config, "darwin")
	if err != nil {
		t.Fatal(err)
	}
	handler := &manifestHandler{root: root, environment: []string{"APPLE_TEST_PASSWORD=fixture-only"}}
	if _, err := handler.sign(t.Context(), pipeline.SignSpec{TargetOS: "darwin", Format: "binary", Input: "program", Config: resolved.Signing.Darwin}); err != nil {
		t.Fatal("sign-only operation used unrelated notarization defaults", err)
	}
	assertAppleSignedArtifact(t, filepath.Join(root, "program.signed"))
	config.Selected = manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "darwin/arm64", Sign: true, Notarize: true}}}
	resolved, err = signingconfig.ApplyDefaults(config)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Signing.Darwin.NotarizationKey != global.Signing.Darwin.APIKeyPath {
		t.Fatal("notarized profile lost global API key")
	}
}
