package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestGeneratedSchemaConformanceCoversEveryAttribute(t *testing.T) {
	type attributeCase struct {
		path        string
		destination reflect.Type
	}
	var cases []attributeCase
	var collect func(*schemaNode, string)
	collect = func(node *schemaNode, parent string) {
		for _, name := range node.attributeOrder {
			descriptor := node.attributes[name]
			if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
				continue
			}
			cases = append(cases, attributeCase{path: schemaPath(parent, name), destination: node.typeInfo.Field(descriptor.fieldIndex).Type})
		}
		for _, name := range node.blockOrder {
			descriptor := node.blocks[name]
			if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
				continue
			}
			path := schemaPath(parent, name)
			for _, label := range descriptor.node.labelNames {
				path += `[` + strconv.Quote(label) + `]`
			}
			collect(descriptor.node, path)
		}
	}
	collect(manifestSchema, "")
	sort.Slice(cases, func(left, right int) bool { return cases[left].path < cases[right].path })

	reference := SchemaReference()
	require.Len(t, cases, len(reference))
	for index, test := range cases {
		field := reference[index]
		assert.Equal(t, field.Path, test.path)
		t.Run(test.path, func(t *testing.T) {
			require.NotEmpty(t, field.Description)
			require.NotEmpty(t, field.Example)
			_, err := decodeSchemaAttribute(parseSchemaExpression(t, field.Example), test.destination)
			require.NoError(t, err, "generated example %s", field.Example)
			valid, invalid := schemaConformanceExpressions(test.destination)
			decoded, err := decodeSchemaAttribute(parseSchemaExpression(t, valid), test.destination)
			require.NoError(t, err)
			assert.Equal(t, test.destination, decoded.Type())
			_, err = decodeSchemaAttribute(parseSchemaExpression(t, invalid), test.destination)
			require.Error(t, err)
		})
	}
}

func TestSchemaDecoderRejectsContainerAndIntegerEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		destination reflect.Type
	}{
		{"fractional integer", "1.5", reflect.TypeOf((*int)(nil))},
		{"overflowing integer", "999999999999999999999999999999999999", reflect.TypeOf((*int)(nil))},
		{"list scalar", `"value"`, reflect.TypeOf((*[]string)(nil))},
		{"map tuple", `[]`, reflect.TypeOf((*map[string]string)(nil))},
		{"evaluation diagnostic", `missing`, reflect.TypeOf((*string)(nil))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := decodeSchemaAttribute(parseSchemaExpression(t, test.expression), test.destination)
			require.Error(t, err)
		})
	}
	_, err := decodeSchemaAttribute(&hclsyntax.LiteralValueExpr{Val: cty.NullVal(cty.String)}, reflect.TypeOf((*string)(nil)))
	require.Error(t, err)
	_, err = decodeSchemaAttribute(parseSchemaExpression(t, `1.0`), reflect.TypeOf((*float64)(nil)))
	require.Error(t, err)

	synthetic := buildSchemaNode(reflect.TypeOf(struct {
		Ignored string
		Value   *string `hcl:"value,optional"`
	}{}))
	assert.Contains(t, synthetic.attributes, "value")
	assert.NotContains(t, synthetic.attributes, "")
}

func TestSchemaDecoderRejectsWrongLabelsAndDuplicateSingletonBlocks(t *testing.T) {
	source := []byte(`version = 3

project {
  name = "labels"
  product_name = "Labels"
  identifier = "com.example.labels"
  version = "1.0.0"
}

build {}
build {}
profile {}
`)
	_, err := decodeHCL(t.TempDir(), "labels.hcl", source, "")
	require.Error(t, err)
	var fields []string
	for _, item := range joinedErrors(err) {
		var validation *ValidationError
		if errors.As(item, &validation) {
			fields = append(fields, validation.Field)
		}
	}
	assert.Equal(t, []string{"build", "profile"}, fields)
}

func schemaConformanceExpressions(destination reflect.Type) (valid, invalid string) {
	element := destination
	if element.Kind() == reflect.Pointer {
		element = element.Elem()
	}
	switch element.Kind() {
	case reflect.String:
		return `"value"`, `true`
	case reflect.Int:
		return `7`, `"seven"`
	case reflect.Bool:
		return `true`, `"true"`
	case reflect.Slice:
		return `["value"]`, `[true]`
	case reflect.Map:
		return `{ KEY = "value" }`, `{ KEY = true }`
	default:
		panic("unsupported schema conformance type " + element.String())
	}
}

func parseSchemaExpression(t *testing.T, source string) hcl.Expression {
	t.Helper()
	expression, diagnostics := hclsyntax.ParseExpression([]byte(source), "expression.hcl", hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	return expression
}

func TestSchemaMetadataDefensiveContracts(t *testing.T) {
	for _, name := range []string{"windows", "darwin", "linux", "ios", "android", "nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"} {
		assert.NotZero(t, applicabilityBit(name), name)
	}
	assert.Zero(t, applicabilityBit("unknown"))

	assert.Panics(t, func() { schemaTypeName(reflect.TypeOf(float64(0))) })
	assert.Panics(t, func() { schemaExample("synthetic", "synthetic", reflect.TypeOf(struct{}{}), "") })
	assert.Equal(t, `{ KEY = "value" }`, schemaExample("synthetic", "synthetic", reflect.TypeOf(map[string]string{}), ""))
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*int)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*bool)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*[]string)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*float64)(nil)), "1") })

	assert.Panics(t, func() { applySchemaDefaults(nil) })
	assert.Panics(t, func() { applySchemaDefaults(&struct{}{}) })
	project := &hclProject{}
	applySchemaDefaults(project)
	assert.Nil(t, project.Name)

	nonemptyMap := decodeSchemaDefault(reflect.TypeOf((*map[string]string)(nil)), `{"A":"B"}`)
	cloned := cloneSchemaDefault(nonemptyMap)
	assert.Equal(t, map[string]string{"A": "B"}, cloned.Elem().Interface())

	original := fmt.Errorf("plain")
	assert.Same(t, original, attachValidationRange(original, nil))
	ranged := &ValidationError{Field: "field", Range: SourceRange{Filename: "already.hcl"}}
	assert.Same(t, ranged, attachValidationRange(ranged, map[string]Origin{"field": {Kind: OriginManifest}}))
	unranged := &ValidationError{Field: "field"}
	assert.Same(t, unranged, attachValidationRange(unranged, map[string]Origin{"field": {Kind: OriginDefault}}))

	diagnostics := hcl.Diagnostics{
		&hcl.Diagnostic{Severity: hcl.DiagWarning, Summary: "ignored"},
		&hcl.Diagnostic{Severity: hcl.DiagError, Summary: "broken"},
	}
	err := validationFromDiagnostics(diagnostics)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "broken")
	assert.NotContains(t, err.Error(), "ignored")
}

func TestSchemaDefaultsReturnIndependentCollections(t *testing.T) {
	first := defaults(Project{Name: "first"})
	first.Frontend.Environment["MUTATED"] = "yes"
	first.Frontend.Install[0] = "mutated"

	second := defaults(Project{Name: "second"})
	assert.Empty(t, second.Frontend.Environment)
	assert.Equal(t, []string{"npm", "install"}, second.Frontend.Install)
}

func TestGeneratedSchemaDefaultsApplyAndExposeOrigins(t *testing.T) {
	origins := defaultOrigins()
	var walk func(*schemaNode, string)
	walk = func(node *schemaNode, parent string) {
		value := reflect.New(node.typeInfo)
		applySchemaDefaults(value.Interface())
		for _, name := range node.attributeOrder {
			descriptor := node.attributes[name]
			if descriptor.defaultText == "" {
				continue
			}
			path := schemaPath(parent, name)
			assert.Equal(t, OriginDefault, origins[path].Kind, path)
			actual := value.Elem().Field(descriptor.fieldIndex)
			if descriptor.defaultValue.IsValid() {
				assert.Equal(t, descriptor.defaultValue.Interface(), actual.Interface(), path)
			} else {
				assert.True(t, actual.IsZero(), path)
			}
		}
		for _, name := range node.blockOrder {
			descriptor := node.blocks[name]
			path := schemaPath(parent, name)
			for _, label := range descriptor.node.labelNames {
				path += "[" + strconv.Quote(label) + "]"
			}
			walk(descriptor.node, path)
		}
	}
	walk(manifestSchema, "")
}

func TestGeneratedApplicabilityRejectsEveryInvalidPlatformAndFormatField(t *testing.T) {
	platforms := []string{"windows", "darwin", "linux", "ios", "android"}
	platformNode := manifestSchema.blocks["windows"].node
	for _, name := range platformNode.attributeOrder {
		descriptor := platformNode.attributes[name]
		if descriptor.platformMask == 0 {
			continue
		}
		field := platformNode.typeInfo.Field(descriptor.fieldIndex)
		for _, platform := range platforms {
			t.Run(platform+"."+name, func(t *testing.T) {
				source := generatedApplicabilitySource(platform, "", name, schemaRoundTripExample(schemaPath(platform, name), name, field.Type, descriptor.defaultText))
				err := decodeGeneratedSchemaSource(source)
				if descriptor.platformMask&applicabilityBit(platform) != 0 {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, schemaPath(platform, name))
				}
			})
		}
	}
	for _, platform := range platforms {
		t.Run(platform+".notarization", func(t *testing.T) {
			source := generatedApplicabilitySource(platform, "", "notarization", "{\n    credential = \"notary\"\n  }")
			err := decodeGeneratedSchemaSource(source)
			if platform == "darwin" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, schemaPath(platform, "notarization"))
			}
		})
	}

	formats := []string{"nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"}
	packageNode := schemaNodesByType[reflect.TypeOf(PackageFormat{})]
	for _, name := range packageNode.attributeOrder {
		descriptor := packageNode.attributes[name]
		field := packageNode.typeInfo.Field(descriptor.fieldIndex)
		for _, format := range formats {
			t.Run("package/"+format+"."+name, func(t *testing.T) {
				parent := `package["` + format + `"]`
				source := generatedApplicabilitySource("package", format, name, schemaRoundTripExample(schemaPath(parent, name), name, field.Type, descriptor.defaultText))
				err := decodeGeneratedSchemaSource(source)
				if descriptor.formatMask&applicabilityBit(format) != 0 {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, schemaPath(parent, name))
				}
			})
		}
	}
}

func generatedApplicabilitySource(block, label, attribute, example string) []byte {
	var output bytes.Buffer
	output.WriteString("version = 3\nproject {\n  name = \"app\"\n  product_name = \"App\"\n  identifier = \"com.example.app\"\n  version = \"1.0.0\"\n}\n")
	output.WriteString(block)
	if label != "" {
		output.WriteByte(' ')
		output.WriteString(strconv.Quote(label))
	}
	output.WriteString(" {\n  ")
	output.WriteString(attribute)
	if attribute == "notarization" {
		output.WriteByte(' ')
	} else {
		output.WriteString(" = ")
	}
	output.WriteString(example)
	output.WriteString("\n}\n")
	return output.Bytes()
}

func decodeGeneratedSchemaSource(source []byte) error {
	file, diagnostics := hclsyntax.ParseConfig(source, "generated.hcl", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return errors.New(diagnostics.Error())
	}
	_, err := decodeManifestSchema(file.Body.(*hclsyntax.Body))
	return err
}

func TestPublicManifestRoundTripCoversEverySchemaField(t *testing.T) {
	root := t.TempDir()
	template := filepath.Join(root, "packaging", "package.tmpl")
	require.NoError(t, os.MkdirAll(filepath.Dir(template), 0o755))
	require.NoError(t, os.WriteFile(template, []byte("template"), 0o644))
	hook := filepath.Join(root, filepath.FromSlash(schemaRoundTripHookScript()))
	require.NoError(t, os.MkdirAll(filepath.Dir(hook), 0o755))
	require.NoError(t, os.WriteFile(hook, schemaRoundTripHookContents(), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("1.0.0\n"), 0o644))
	source := generatedSchemaRoundTripManifest()
	path := filepath.Join(root, Filename)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err, string(source))
	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	output := parseEjectedBody(t, encoded)
	outputOrigins := manifestOrigins(output)
	input := parseEjectedBody(t, source)
	inputOrigins := manifestOrigins(input)
	for field := range inputOrigins {
		assert.Contains(t, outputOrigins, field, field)
	}

	covered := make(map[string]bool, len(inputOrigins))
	for field := range inputOrigins {
		covered[normalizeGeneratedSchemaPath(field)] = true
	}
	for _, field := range SchemaReference() {
		assert.True(t, covered[field.Path], "generated public round trip does not cover %s", field.Path)
	}
}

func generatedSchemaRoundTripManifest() []byte {
	var output bytes.Buffer
	writeGeneratedSchemaNode(&output, manifestSchema, "", "")
	return output.Bytes()
}

func writeGeneratedSchemaNode(output *bytes.Buffer, node *schemaNode, parent, indent string) {
	for _, name := range node.attributeOrder {
		descriptor := node.attributes[name]
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		path := schemaPath(parent, name)
		field := node.typeInfo.Field(descriptor.fieldIndex)
		example := schemaRoundTripExample(path, name, field.Type, descriptor.defaultText)
		fmt.Fprintf(output, "%s%s = %s\n", indent, name, example)
	}
	for _, name := range node.blockOrder {
		descriptor := node.blocks[name]
		if name == "package" && parent == "" {
			writeGeneratedPackageBlocks(output, descriptor.node, indent)
			continue
		}
		if name == "profile" && parent == "" {
			writeGeneratedProfileBlock(output, indent)
			continue
		}
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		label := generatedSchemaLabel(name, parent)
		path := schemaPath(parent, name)
		output.WriteString(indent)
		output.WriteString(name)
		if label != "" {
			output.WriteByte(' ')
			output.WriteString(strconv.Quote(label))
			path += "[" + strconv.Quote(label) + "]"
		}
		output.WriteString(" {\n")
		writeGeneratedSchemaNode(output, descriptor.node, path, indent+"  ")
		output.WriteString(indent)
		output.WriteString("}\n\n")
	}
}

func writeGeneratedProfileBlock(output *bytes.Buffer, indent string) {
	fmt.Fprintf(output, `%sprofile "release" {
%s  target "windows/amd64" {
%s    formats = ["nsis"]
%s    sign = true
%s  }

%s  target "darwin/arm64" {
%s    sign = true
%s    notarize = true
%s  }

%s  target "ios/arm64" {
%s    destination = "device"
%s  }
%s}

`, indent, indent, indent, indent, indent, indent, indent, indent, indent, indent, indent, indent, indent)
}

func writeGeneratedPackageBlocks(output *bytes.Buffer, node *schemaNode, indent string) {
	for _, format := range []string{"nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"} {
		parent := `package["` + format + `"]`
		var attributes []string
		for _, name := range node.attributeOrder {
			descriptor := node.attributes[name]
			if schemaFormatAllowed(parent, descriptor.formatMask) && (name != "template" || format == "archlinux") && (format != "archlinux" || name == "template") {
				attributes = append(attributes, name)
			}
		}
		if len(attributes) == 0 {
			continue
		}
		fmt.Fprintf(output, "%spackage %s {\n", indent, strconv.Quote(format))
		for _, name := range attributes {
			descriptor := node.attributes[name]
			field := node.typeInfo.Field(descriptor.fieldIndex)
			fmt.Fprintf(output, "%s  %s = %s\n", indent, name, schemaRoundTripExample(schemaPath(parent, name), name, field.Type, descriptor.defaultText))
		}
		output.WriteString(indent)
		output.WriteString("}\n\n")
	}
}

func generatedSchemaLabel(block, parent string) string {
	switch block {
	case "target":
		return "windows/amd64"
	case "profile":
		return "release"
	case "file_association":
		return "association"
	case "protocol":
		return "example"
	case "hook":
		return "before_build"
	default:
		return ""
	}
}

func schemaRoundTripExample(path, name string, destination reflect.Type, defaultText string) string {
	if strings.HasPrefix(path, `file_association[`) && name == "name" {
		return `"association"`
	}
	if strings.HasPrefix(path, `hook[`) {
		switch name {
		case "script":
			return strconv.Quote(schemaRoundTripHookScript())
		case "directory":
			return `"scripts"`
		case "inputs":
			return `["version.txt"]`
		case "outputs":
			return `["generated/version.go"]`
		}
	}
	element := destination
	for element.Kind() == reflect.Pointer {
		element = element.Elem()
	}
	if element.Kind() == reflect.Bool {
		if defaultText == "true" {
			return "false"
		}
		return "true"
	}
	if (element.Kind() == reflect.Slice || element.Kind() == reflect.Map) && (defaultText == "[]" || defaultText == "{}") {
		valid, _ := schemaConformanceExpressions(destination)
		return valid
	}
	return schemaExample(path, name, destination, defaultText)
}

func schemaRoundTripHookScript() string {
	if runtime.GOOS == "windows" {
		return "scripts/generate-version.cmd"
	}
	return "scripts/generate-version.sh"
}

func schemaRoundTripHookContents() []byte {
	if runtime.GOOS == "windows" {
		return []byte("@echo off\r\n")
	}
	return []byte("#!/bin/sh\n")
}

func normalizeGeneratedSchemaPath(path string) string {
	for _, replacement := range []struct{ actual, schema string }{
		{`package["nsis"]`, `package["format"]`}, {`package["msix"]`, `package["format"]`},
		{`package["dmg"]`, `package["format"]`}, {`package["appimage"]`, `package["format"]`},
		{`package["deb"]`, `package["format"]`}, {`package["rpm"]`, `package["format"]`},
		{`package["archlinux"]`, `package["format"]`}, {`profile["release"]`, `profile["profile"]`},
		{`target["windows/amd64"]`, `target["target"]`}, {`file_association["association"]`, `file_association["association"]`},
		{`target["darwin/arm64"]`, `target["target"]`}, {`target["ios/arm64"]`, `target["target"]`},
		{`protocol["example"]`, `protocol["scheme"]`},
		{`hook["before_build"]`, `hook["phase"]`},
	} {
		path = strings.ReplaceAll(path, replacement.actual, replacement.schema)
	}
	return path
}
