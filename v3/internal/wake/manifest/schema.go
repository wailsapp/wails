package manifest

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ValidationError struct {
	Field  string
	Detail string
	Range  SourceRange
	Cause  error
}

func (e *ValidationError) Unwrap() error { return e.Cause }

type SchemaField struct {
	Path        string
	Type        string
	Required    bool
	Default     string
	Example     string
	Description string
	Platforms   string
	Formats     string
}

func SchemaReference() []SchemaField {
	var result []SchemaField
	collectSchemaFields(manifestSchema, "", &result)
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	return result
}

func SchemaReferenceMarkdown() []byte {
	var output strings.Builder
	output.WriteString("| Field | Type | Required | Default | Example | Applies to | Description |\n")
	output.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	for _, field := range SchemaReference() {
		defaultValue := ""
		if field.Default != "" {
			defaultValue = "`" + field.Default + "`"
		}
		appliesTo := strings.Join(nonempty(field.Platforms, field.Formats), ", ")
		fmt.Fprintf(&output, "| `%s` | %s | %s | %s | `%s` | %s | %s |\n", field.Path, field.Type, map[bool]string{true: "yes", false: "no"}[field.Required], defaultValue, field.Example, appliesTo, field.Description)
	}
	return []byte(output.String())
}

func collectSchemaFields(node *schemaNode, parent string, result *[]SchemaField) {
	for _, name := range node.attributeOrder {
		descriptor := node.attributes[name]
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		field := node.typeInfo.Field(descriptor.fieldIndex)
		path := schemaPath(parent, name)
		*result = append(*result, SchemaField{
			Path: path, Type: schemaTypeName(field.Type), Required: descriptor.required, Default: descriptor.defaultText,
			Example: schemaExample(path, name, field.Type, descriptor.defaultText), Description: schemaDescription(path, name),
			Platforms: descriptor.platforms, Formats: descriptor.formats,
		})
	}
	for _, name := range node.blockOrder {
		descriptor := node.blocks[name]
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		path := schemaPath(parent, name)
		for _, label := range descriptor.node.labelNames {
			path += "[" + strconv.Quote(label) + "]"
		}
		collectSchemaFields(descriptor.node, path, result)
	}
}

func nonempty(values ...string) []string {
	result := values[:0]
	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

var schemaAttributeDescriptions = map[string]string{
	"application_id":       "Android application identifier.",
	"assets_car":           "User-owned compiled Apple asset catalogue.",
	"background":           "User-owned DMG background image.",
	"background_modes":     "iOS background execution modes.",
	"binary_name":          "Base name of the native executable.",
	"build":                "Frontend production build command and arguments.",
	"cache":                "Allow this hook to use the artifact cache; complete inputs and outputs are required.",
	"build_number":         "Platform build number.",
	"bundle_id":            "iOS bundle identifier.",
	"capabilities":         "Native platform capabilities requested by the application.",
	"categories":           "Desktop menu categories.",
	"certificate":          "Project-relative signing certificate file.",
	"cf_bundle_icon_name":  "Apple asset-catalogue icon set name.",
	"comments":             "Additional package metadata comments.",
	"company":              "Company or organisation name.",
	"compiler_flags":       "Additional Go compiler flags.",
	"copyright":            "Copyright notice embedded in artifacts.",
	"credential":           "Named external credential reference; never secret material.",
	"debounce_ms":          "Delay used to coalesce development file changes.",
	"dependencies":         "Native package dependencies.",
	"description":          "Human-readable description.",
	"desktop_entry":        "User-owned Linux desktop-entry file.",
	"destination":          "Mobile build destination.",
	"dev":                  "Frontend development-server command and arguments.",
	"directory":            "Project-relative frontend source directory.",
	"display_name":         "Name displayed by the mobile platform.",
	"entitlements":         "Project-relative Apple entitlements file.",
	"environment":          "Environment variables added to this operation.",
	"exclude":              "Development watcher exclusion patterns.",
	"extensions":           "File extensions handled by the application.",
	"file_icon":            "User-owned DMG file icon.",
	"files":                "User-owned files copied into the DMG, keyed by destination name.",
	"formats":              "Artifact formats produced for this profile target.",
	"garble_args":          "Additional arguments passed to garble.",
	"grace_period_ms":      "Time allowed for a development process to stop cleanly.",
	"icon":                 "Project-relative user-owned icon file.",
	"identifier":           "Reverse-domain application identifier.",
	"identity":             "Signing identity selected through the platform toolchain.",
	"index_filename":       "Base filename for generated binding exports.",
	"info_plist":           "User-owned Apple Info.plist file.",
	"install":              "Frontend dependency installation command and arguments.",
	"install_scope":        "NSIS installation scope.",
	"inputs":               "Project-relative files or directories that determine a cached hook result.",
	"interfaces":           "Generate TypeScript interfaces instead of classes where supported.",
	"key_alias":            "Alias of the Android signing key.",
	"ldflags":              "Additional Go linker flags.",
	"log_level":            "Development log verbosity.",
	"maintainer":           "Native Linux package maintainer.",
	"manifest":             "Project-relative user-owned platform or package manifest.",
	"mime_type":            "MIME type registered for a file association.",
	"minimum_sdk":          "Minimum supported Android SDK level.",
	"minimum_version":      "Minimum supported operating-system version.",
	"models_filename":      "Base filename for generated binding models.",
	"name":                 "Stable project or registration name.",
	"notarize":             "Notarize the produced artifact.",
	"obfuscated":           "Obfuscate Go code with garble.",
	"output":               "Project-relative output directory.",
	"outputs":              "Project-relative files or directories produced completely by a cached hook.",
	"platforms":            "Platforms receiving this registration.",
	"post_install":         "Project-relative Linux post-install script.",
	"post_remove":          "Project-relative Linux post-remove script.",
	"pre_install":          "Project-relative Linux pre-install script.",
	"pre_remove":           "Project-relative Linux pre-remove script.",
	"product_name":         "Human-readable product name.",
	"provisioning_profile": "Project-relative Apple provisioning profile.",
	"publisher":            "Windows package publisher identity.",
	"role":                 "Application role for a file association.",
	"section":              "Linux package repository section.",
	"script":               "Project-relative executable script invoked directly without a shell command string.",
	"sign":                 "Sign the produced artifact.",
	"strip":                "Strip debug information from production binaries.",
	"tags":                 "Additional Go build tags.",
	"target_sdk":           "Android SDK level targeted by the application.",
	"template":             "Complete user-owned package template replacement.",
	"thumbprint":           "Windows certificate thumbprint.",
	"time_type":            "TypeScript representation of Go time values.",
	"timestamp_server":     "URL of the signing timestamp service.",
	"toolchain":            "Compiler toolchain policy for this target.",
	"trim_path":            "Remove local filesystem paths from production binaries.",
	"typescript":           "Generate TypeScript bindings.",
	"use_git_ignore":       "Apply .gitignore rules to the development watcher.",
	"vcs_info":             "Embed version-control metadata in the Go binary.",
	"version":              "Manifest or application version.",
	"version_code":         "Monotonic Android release number.",
	"version_name":         "Human-readable Android release version.",
	"volume_icon":          "User-owned DMG volume icon.",
	"watch":                "Development watcher inclusion patterns.",
	"window_height":        "DMG Finder window height in pixels.",
	"window_width":         "DMG Finder window width in pixels.",
}

var schemaAttributeExamples = map[string]string{
	"application_id": `"com.example.app"`, "assets_car": `"build/Assets.car"`, "background": `"assets/dmg-background.png"`,
	"background_modes": `["fetch"]`, "bundle_id": `"com.example.app"`, "capabilities": `["internetClient"]`,
	"certificate": `"signing/certificate.p12"`, "categories": `["Development"]`, "cf_bundle_icon_name": `"AppIcon"`,
	"compiler_flags": `["all=-l"]`, "credential": `"release-signing"`, "dependencies": `["libgtk-4-1"]`,
	"desktop_entry": `"assets/app.desktop"`, "destination": `"device"`, "display_name": `"Example App"`,
	"entitlements": `"signing/entitlements.plist"`, "environment": `{ RELEASE = "true" }`, "exclude": `["node_modules"]`,
	"extensions": `["example"]`, "file_icon": `"assets/file.icns"`, "files": `{ "License.pdf" = "LICENSE.pdf" }`,
	"formats": `["nsis"]`, "garble_args": `["-literals"]`, "icon": `"assets/appicon.png"`,
	"inputs":     `["version.txt"]`,
	"identifier": `"com.example.app"`, "identity": `"Developer ID Application: Example"`, "info_plist": `"assets/Info.plist"`,
	"install_scope": `"user"`, "key_alias": `"upload"`, "ldflags": `["-X example/build.version=1.0.0"]`,
	"log_level": `"info"`, "maintainer": `"Example <release@example.com>"`, "manifest": `"assets/manifest.xml"`,
	"mime_type": `"application/x-example"`, "minimum_sdk": `24`, "minimum_version": `"12.0"`,
	"platforms": `["windows", "darwin", "linux"]`, "post_install": `"packaging/postinstall.sh"`,
	"post_remove": `"packaging/postremove.sh"`, "pre_install": `"packaging/preinstall.sh"`, "pre_remove": `"packaging/preremove.sh"`,
	"provisioning_profile": `"signing/app.mobileprovision"`, "publisher": `"CN=Example"`, "role": `"editor"`,
	"section": `"utils"`, "template": `"packaging/package.tmpl"`, "thumbprint": `"0123456789ABCDEF"`,
	"outputs": `["generated/version.go"]`, "script": `"scripts/generate-version.sh"`,
	"time_type": `"string"`, "timestamp_server": `"https://timestamp.example.com"`, "toolchain": `"auto"`,
	"version_name": `"1.0.0"`, "volume_icon": `"assets/volume.icns"`,
	"watch": `["**/*.go", "wails.hcl"]`,
}

func schemaDescription(path, name string) string {
	if path == "version" {
		return "Wails manifest schema version; this must be the first attribute."
	}
	if strings.HasPrefix(path, `hook[`) && name == "directory" {
		return "Project-relative working directory used to invoke the hook."
	}
	return schemaAttributeDescriptions[name]
}

func schemaExample(path, name string, destination reflect.Type, defaultText string) string {
	if path == "version" {
		return "3"
	}
	if example := schemaAttributeExamples[name]; example != "" {
		return example
	}
	element := destination
	for element.Kind() == reflect.Pointer {
		element = element.Elem()
	}
	if defaultText != "" {
		if strings.HasPrefix(defaultText, "$") {
			return `"my-app"`
		}
		if element.Kind() == reflect.String {
			return strconv.Quote(defaultText)
		}
		return defaultText
	}
	switch element.Kind() {
	case reflect.String:
		return `"value"`
	case reflect.Int:
		return "1"
	case reflect.Bool:
		return "true"
	case reflect.Slice:
		return `["value"]`
	case reflect.Map:
		return `{ KEY = "value" }`
	default:
		panic("unsupported manifest schema example type " + element.String())
	}
}

func schemaTypeName(value reflect.Type) string {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int:
		return "integer"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		return "list(" + schemaTypeName(value.Elem()) + ")"
	case reflect.Map:
		return "map(" + schemaTypeName(value.Elem()) + ")"
	}
	panic("unsupported manifest schema type " + value.String())
}

func (e *ValidationError) Error() string {
	if e.Range.Filename == "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Detail)
	}
	return fmt.Sprintf("%s:%d:%d: %s: %s", e.Range.Filename, e.Range.StartLine, e.Range.StartColumn, e.Field, e.Detail)
}

func fieldValidationError(field, format string, args ...any) error {
	return &ValidationError{Field: field, Detail: fmt.Sprintf(format, args...)}
}

func fieldValidationCause(field string, cause error, format string, args ...any) error {
	return &ValidationError{Field: field, Detail: fmt.Sprintf(format, args...), Cause: cause}
}

func attachValidationRange(err error, origins map[string]Origin) error {
	var validation *ValidationError
	if !errors.As(err, &validation) || validation.Range.Filename != "" {
		return err
	}
	origin, ok := origins[validation.Field]
	if !ok || origin.Kind != OriginManifest {
		return err
	}
	copy := *validation
	copy.Range = origin.Range
	return &copy
}

func validationFromDiagnostics(diagnostics hcl.Diagnostics) error {
	var errorsFound []error
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity != hcl.DiagError {
			continue
		}
		field := "manifest"
		var diagnosticRange hcl.Range
		if diagnostic.Subject != nil {
			diagnosticRange = *diagnostic.Subject
		}
		detail := "parse " + Filename + ": " + diagnostic.Summary
		if diagnostic.Detail != "" {
			detail += ": " + diagnostic.Detail
		}
		errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: detail, Range: sourceRange(diagnosticRange)})
	}
	return errors.Join(errorsFound...)
}

func schemaPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "." + name
}

func defaultOrigins() map[string]Origin {
	result := make(map[string]Origin, len(cachedSchemaDefaultPaths))
	for _, field := range cachedSchemaDefaultPaths {
		result[field] = Origin{Kind: OriginDefault}
	}
	return result
}

func applySchemaDefaults(target any) {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		panic("manifest schema defaults require a non-nil pointer to a struct")
	}
	value = value.Elem()
	node := schemaNodesByType[value.Type()]
	if node == nil {
		panic("manifest schema defaults require a schema struct")
	}
	for _, name := range node.attributeOrder {
		descriptor := node.attributes[name]
		if !descriptor.defaultValue.IsValid() {
			continue
		}
		value.Field(descriptor.fieldIndex).Set(cloneSchemaDefault(descriptor.defaultValue))
	}
}

func cloneSchemaDefault(value reflect.Value) reflect.Value {
	result := reflect.New(value.Elem().Type())
	switch value.Elem().Kind() {
	case reflect.Slice:
		result.Elem().Set(reflect.AppendSlice(reflect.MakeSlice(value.Elem().Type(), 0, value.Elem().Len()), value.Elem()))
	case reflect.Map:
		copy := reflect.MakeMapWithSize(value.Elem().Type(), value.Elem().Len())
		iterator := value.Elem().MapRange()
		for iterator.Next() {
			copy.SetMapIndex(iterator.Key(), iterator.Value())
		}
		result.Elem().Set(copy)
	default:
		result.Elem().Set(value.Elem())
	}
	return result
}

var cachedSchemaDefaultPaths = buildSchemaDefaultPaths()

func buildSchemaDefaultPaths() []string {
	var result []string
	collectSchemaDefaultPaths(manifestSchema, "", &result)
	sort.Strings(result)
	return result
}

func collectSchemaDefaultPaths(node *schemaNode, parent string, result *[]string) {
	for _, name := range node.attributeOrder {
		if node.attributes[name].defaultText != "" {
			*result = append(*result, schemaPath(parent, name))
		}
	}
	for _, name := range node.blockOrder {
		descriptor := node.blocks[name]
		path := schemaPath(parent, name)
		for _, label := range descriptor.node.labelNames {
			path += "[" + strconv.Quote(label) + "]"
		}
		collectSchemaDefaultPaths(descriptor.node, path, result)
	}
}

func manifestOrigins(body *hclsyntax.Body) map[string]Origin {
	result := map[string]Origin{}
	collectManifestOrigins(result, "", body)
	return result
}

func collectManifestOrigins(result map[string]Origin, parent string, body *hclsyntax.Body) {
	for name, attribute := range body.Attributes {
		field := name
		if parent != "" {
			field = parent + "." + name
		}
		result[field] = Origin{Kind: OriginManifest, Range: sourceRange(attribute.Range())}
	}
	for _, block := range body.Blocks {
		field := block.Type
		if parent != "" {
			field = parent + "." + field
		}
		for _, label := range block.Labels {
			field += "[" + strconv.Quote(label) + "]"
		}
		collectManifestOrigins(result, field, block.Body)
	}
}

func sourceRange(value hcl.Range) SourceRange {
	return SourceRange{
		Filename: value.Filename, StartLine: value.Start.Line, StartColumn: value.Start.Column,
		EndLine: value.End.Line, EndColumn: value.End.Column,
	}
}
