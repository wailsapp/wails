package commands

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestInjectOpenHarmonySigningConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "build-profile.json5")
	original := `{
  "app": {
    "signingConfigs": [],
    "products": [{"name": "default", "signingConfig": "default"}]
  },
  "modules": [{"name": "entry", "srcPath": "./app"}]
}`
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := injectOpenHarmonySigningConfig(path, "key.p12", "cert.cer", "profile.p7b", "myAlias", "secret"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		App struct {
			SigningConfigs []struct {
				Name      string `json:"name"`
				Keystore  string `json:"keystore"`
				Certpath  string `json:"certpath"`
				Profile   string `json:"profile"`
				Keyalias  string `json:"keyalias"`
				Storepass string `json:"storepass"`
			} `json:"signingConfigs"`
			Products []struct {
				SigningConfig string `json:"signingConfig"`
			} `json:"products"`
		} `json:"app"`
		Modules []json.RawMessage `json:"modules"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.App.SigningConfigs) != 1 {
		t.Fatalf("want one signing config, got %d", len(document.App.SigningConfigs))
	}
	config := document.App.SigningConfigs[0]
	if config.Keystore != "key.p12" || config.Certpath != "cert.cer" || config.Profile != "profile.p7b" || config.Keyalias != "myAlias" || config.Storepass != "secret" {
		t.Fatalf("unexpected signing config: %+v", config)
	}
	if len(document.App.Products) != 1 || document.App.Products[0].SigningConfig != "default" {
		t.Fatalf("product signingConfig not set: %+v", document.App.Products)
	}
	if len(document.Modules) != 1 {
		t.Fatalf("modules lost during injection: %+v", document.Modules)
	}
}

func TestApplyHarmonyOSFormatProfileOnlyPatchesStagedHosHap(t *testing.T) {
	embedded, err := buildAssets.ReadFile("build_assets/openharmony/build-profile.json5")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		format, runtimeOS string
	}{
		{format: "hap", runtimeOS: "OpenHarmony"},
		{format: "hos_hap", runtimeOS: "HarmonyOS"},
	} {
		t.Run(test.format, func(t *testing.T) {
			staged := t.TempDir()
			path := filepath.Join(staged, "build-profile.json5")
			if err := os.WriteFile(path, embedded, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := applyHarmonyOSFormatProfile(staged, test.format); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var document struct {
				App struct {
					Products []struct {
						RuntimeOS string `json:"runtimeOS"`
					} `json:"products"`
				} `json:"app"`
				Modules []json.RawMessage `json:"modules"`
			}
			if err := json.Unmarshal(data, &document); err != nil {
				t.Fatal(err)
			}
			if len(document.App.Products) != 1 || document.App.Products[0].RuntimeOS != test.runtimeOS {
				t.Fatalf("%s runtimeOS: %+v", test.format, document.App.Products)
			}
			if len(document.Modules) != 1 {
				t.Fatalf("%s lost modules: %+v", test.format, document.Modules)
			}
		})
	}
	after, err := buildAssets.ReadFile("build_assets/openharmony/build-profile.json5")
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(embedded) {
		t.Fatal("embedded OpenHarmony build profile was modified by the staged patch")
	}
}

func TestPackageArtifactDispatchesBothHapFormatsToOpenHarmony(t *testing.T) {
	t.Setenv("DEVECO_SDK_HOME", "")
	h := &manifestHandler{root: t.TempDir()}
	for _, format := range []string{"hap", "hos_hap"} {
		_, err := h.packageArtifact(context.Background(), pipeline.PackageSpec{TargetOS: "harmonyos", TargetArch: "amd64", Format: format, Assets: "assets", Output: "bin/app.hap"})
		if err == nil || !strings.Contains(err.Error(), "DEVECO_SDK_HOME") {
			t.Fatalf("%s did not reach the OpenHarmony packager: %v", format, err)
		}
	}
}

func TestSignDispatchesHarmonyOSToOpenHarmonySigner(t *testing.T) {
	t.Setenv("DEVECO_SDK_HOME", "")
	h := &manifestHandler{root: t.TempDir()}
	_, err := h.sign(context.Background(), pipeline.SignSpec{TargetOS: "harmonyos", TargetArch: "amd64", Format: "hos_hap", Input: "bin/app.hap", Config: manifest.SigningPlatform{Enabled: true}})
	if err == nil || !strings.Contains(err.Error(), "DEVECO_SDK_HOME") {
		t.Fatalf("harmonyos signing did not reach the OpenHarmony signer: %v", err)
	}
}

func TestHarmonyOSDeployRequiresAnExplicitHAP(t *testing.T) {
	err := HarmonyOSDeploy(&HarmonyOSDeployOptions{})
	if err == nil || !strings.Contains(err.Error(), "-targets harmonyos/<arch> -formats hos_hap") {
		t.Fatalf("unexpected HarmonyOS deploy guidance: %v", err)
	}
}
