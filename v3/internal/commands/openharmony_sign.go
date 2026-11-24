package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"github.com/wailsapp/wails/v3/internal/signingconfig"
)

// signOpenHarmony signs a HAP by re-running hvigor with a signingConfigs
// section injected into a staging copy of the generated hvigor project.
//
// Field mapping in signing.openharmony:
//   certificate          .p12 keystore
//   entitlements         .cer signing certificate
//   provisioning_profile .p7b profile
//   identity/key_alias   key alias inside the keystore
//   credential           environment variable holding the keystore password
func (h *manifestHandler) signOpenHarmony(ctx context.Context, s pipeline.SignSpec) (string, error) {
	sdk := os.Getenv("DEVECO_SDK_HOME")
	if sdk == "" {
		return "", fmt.Errorf("OpenHarmony signing requires DEVECO_SDK_HOME")
	}
	studio := filepath.Dir(sdk)
	hvigorw := filepath.Join(studio, "tools", "hvigor", "bin", "hvigorw")
	if manifestHostOS == "windows" {
		hvigorw += ".bat"
	}
	if _, err := os.Stat(hvigorw); err != nil {
		return "", fmt.Errorf("hvigorw not found at %s: %w", hvigorw, err)
	}

	keystore, err := signingconfig.ResolvePath(h.root, s.Config, "certificate", s.Config.Certificate)
	if err != nil {
		return "", fmt.Errorf("signing.openharmony.certificate: %w", err)
	}
	certificate, err := signingconfig.ResolvePath(h.root, s.Config, "entitlements", s.Config.Entitlements)
	if err != nil {
		return "", fmt.Errorf("signing.openharmony.entitlements: %w", err)
	}
	profile, err := signingconfig.ResolvePath(h.root, s.Config, "provisioning_profile", s.Config.ProvisioningProfile)
	if err != nil {
		return "", fmt.Errorf("signing.openharmony.provisioning_profile: %w", err)
	}
	alias := chooseString(s.Config.KeyAlias != "", s.Config.KeyAlias, s.Config.Identity)
	password := os.Getenv(s.Config.Credential)
	switch {
	case keystore == "":
		return "", fmt.Errorf("openharmony signing requires signing.openharmony.certificate (.p12 keystore)")
	case certificate == "":
		return "", fmt.Errorf("openharmony signing requires signing.openharmony.entitlements (.cer certificate)")
	case profile == "":
		return "", fmt.Errorf("openharmony signing requires signing.openharmony.provisioning_profile (.p7b profile)")
	case alias == "":
		return "", fmt.Errorf("openharmony signing requires signing.openharmony.identity or key_alias")
	case s.Config.Credential == "" || password == "":
		return "", fmt.Errorf("openharmony signing requires signing.openharmony.credential naming a non-empty password environment variable")
	}

	// The package stage writes <build>/artifacts/<name>.hap and stages the
	// hvigor project at <build>/assets/openharmony.
	buildRoot := filepath.Dir(filepath.Dir(filepath.Join(h.root, s.Input)))
	assets := filepath.Join(buildRoot, "assets", "openharmony")
	if _, err := os.Stat(filepath.Join(assets, "build-profile.json5")); err != nil {
		return "", fmt.Errorf("generated OpenHarmony project not found at %s: %w", assets, err)
	}

	staged, err := os.MkdirTemp(filepath.Dir(buildRoot), ".ohos-sign-stage-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(staged)
	if err := copyManifestPath(assets, staged); err != nil {
		return "", err
	}
	if err := injectOpenHarmonySigningConfig(filepath.Join(staged, "build-profile.json5"), keystore, certificate, profile, alias, password); err != nil {
		return "", err
	}
	abi, err := openHarmonyABI(s.TargetArch)
	if err != nil {
		return "", err
	}
	if err := copyManifestPath(filepath.Join(buildRoot, "libwails_go.so"), filepath.Join(staged, "app", "src", "main", "cpp", "libs", abi, "libwails_go.so")); err != nil {
		return "", err
	}
	if err := copyManifestPath(filepath.Join(buildRoot, "libwails_go.h"), filepath.Join(staged, "app", "src", "main", "cpp", "include", "libwails_go.h")); err != nil {
		return "", err
	}
	output, err := runManifestCommand(ctx, staged, []string{"OHOS_BASE_SDK_HOME=" + filepath.Join(sdk, "default")}, hvigorw, "assembleHap", "--no-daemon")
	if err != nil {
		return output, err
	}
	produced := filepath.Join(staged, "app", "build", "default", "outputs", "default", "entry-default-unsigned.hap")
	return output, copyPathTransactional(produced, filepath.Join(h.root, s.Input+".signed"), ".ohos-sign-output-*")
}

func injectOpenHarmonySigningConfig(path, keystore, certificate, profile, alias, password string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var document struct {
		App struct {
			SigningConfigs []map[string]any `json:"signingConfigs"`
			Products       []map[string]any `json:"products"`
		} `json:"app"`
		Modules []any `json:"modules"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("parse build-profile.json5: %w", err)
	}
	if len(document.App.Products) == 0 {
		return fmt.Errorf("build-profile.json5 has no products")
	}
	document.App.SigningConfigs = []map[string]any{{
		"name":       "default",
		"type":       "pkcs12",
		"certpath":   certificate,
		"keystore":   keystore,
		"storepass":  password,
		"keyalias":   alias,
		"profile":    profile,
		"signAlg":    "SHA256withECDSA",
		"storepassType": "pkcs12",
	}}
	for _, product := range document.App.Products {
		product["signingConfig"] = "default"
	}
	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0o644)
}
