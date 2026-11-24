package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
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
