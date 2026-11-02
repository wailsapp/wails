package manifest

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSigningReferenceUsesPlatformSpecificExamples(t *testing.T) {
	fields := map[string]SchemaField{}
	for _, field := range SchemaReference() {
		fields[field.Path] = field
	}
	for path, expected := range map[string]string{
		"linux.signing.certificate":      `"signing/private-key.asc"`,
		"linux.signing.identity":         `"builder"`,
		"linux.signing.credential":       `"LINUX_SIGNING_PASSPHRASE"`,
		"android.signing.certificate":    `"signing/upload.jks"`,
		"android.signing.identity":       `"upload"`,
		"windows.signing.certificate":    `"signing/windows.pfx"`,
		"darwin.notarization.credential": `"release-notary"`,
	} {
		if field := fields[path]; field.Example != expected {
			t.Errorf("%s example %q; want %q", path, field.Example, expected)
		}
	}
	if !strings.Contains(fields["linux.signing.certificate"].Description, "armored") || !strings.Contains(fields["linux.signing.credential"].Description, "environment variable") {
		t.Fatal("Linux signing key/passphrase contract is missing")
	}
	if !strings.Contains(fields["android.signing.entitlements"].Description, "Not used") {
		t.Fatal("Apple entitlements were documented as an Android signing input")
	}
}

func TestEditorSchemaExposesExistingVoidPackageConfiguration(t *testing.T) {
	data, err := JSONSchema()
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties map[string]struct {
			Properties map[string]struct {
				Properties map[string]json.RawMessage
			}
		}
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	properties := schema.Properties["packages"].Properties["void"].Properties
	if properties["dependencies"] == nil || properties["maintainer"] == nil {
		t.Fatal("editor schema omits supported Void package metadata")
	}
	if properties["defines"] != nil {
		t.Fatal("NSIS-only definitions leaked into the Void schema")
	}
}
