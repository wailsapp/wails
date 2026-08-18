package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ErrEjectionSuggestionsUnavailable means the selected scope is already
// frozen but this CLI does not retain the historical default snapshot needed
// to produce safe three-way upgrade comments. The manifest is unchanged.
var ErrEjectionSuggestionsUnavailable = errors.New("ejection upgrade suggestions unavailable")

func Minimal(project Project) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.SetAttributeValue("version", cty.NumberIntVal(3))
	body.AppendNewline()
	projectBlock := hclwrite.NewBlock("project", nil)
	projectBody := projectBlock.Body()
	projectBody.SetAttributeValue("name", cty.StringVal(project.Name))
	projectBody.SetAttributeValue("product_name", cty.StringVal(project.ProductName))
	projectBody.SetAttributeValue("identifier", cty.StringVal(project.Identifier))
	projectBody.SetAttributeValue("version", cty.StringVal(project.Version))
	projectBody.SetAttributeValue("binary_name", cty.StringVal(deriveBinaryName(project.Name)))
	body.AppendBlock(projectBlock)
	body.AppendNewline()
	frontendBlock := hclwrite.NewBlock("frontend", nil)
	frontendBody := frontendBlock.Body()
	frontendBody.SetAttributeValue("directory", cty.StringVal("frontend"))
	frontendBody.SetAttributeValue("install", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("install")}))
	frontendBody.SetAttributeValue("build", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("build")}))
	frontendBody.SetAttributeValue("dev", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("dev")}))
	frontendBody.SetAttributeValue("output", cty.StringVal("frontend/dist"))
	body.AppendBlock(frontendBlock)
	body.AppendNewline()
	buildBlock := hclwrite.NewBlock("build", nil)
	buildBlock.Body().SetAttributeValue("output", cty.StringVal("bin"))
	body.AppendBlock(buildBlock)
	body.AppendNewline()
	return file.Bytes()
}

func WriteMinimal(root string, project Project) error {
	if err := validateProject(project); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), Minimal(project), 0o644)
}

func EncodeConfig(config Config) ([]byte, error) {
	return encodeConfigHCL(config, "")
}

func EncodeDocument(doc Document) ([]byte, error) {
	if err := validateProject(doc.Project); err != nil {
		return nil, err
	}
	return EncodeConfig(configFromDocument(".", "", doc))
}

// sparseValue turns a programmatically assembled Document into the same sparse
// shape a user would write. Struct zero values inherit compiled defaults;
// explicit values inside profile and extension maps are preserved verbatim.
func sparseValue(value, defaultValue reflect.Value, preserveZero bool) (any, bool) {
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return nil, false
		}
		value = value.Elem()
	}
	for defaultValue.IsValid() && (defaultValue.Kind() == reflect.Interface || defaultValue.Kind() == reflect.Pointer) {
		if defaultValue.IsNil() {
			defaultValue = reflect.Value{}
			break
		}
		defaultValue = defaultValue.Elem()
	}
	if !value.IsValid() {
		return nil, false
	}
	switch value.Kind() {
	case reflect.Struct:
		result := map[string]any{}
		typeInfo := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := typeInfo.Field(i)
			name := strings.Split(field.Tag.Get("toml"), ",")[0]
			if name == "" {
				name = field.Name
			}
			if name == "-" {
				continue
			}
			childPreservesZero := preserveZero || field.Name == "Profiles" || field.Name == "Extensions"
			var childDefault reflect.Value
			if defaultValue.IsValid() && defaultValue.Kind() == reflect.Struct && i < defaultValue.NumField() {
				childDefault = defaultValue.Field(i)
			}
			child, include := sparseValue(value.Field(i), childDefault, childPreservesZero)
			if include {
				result[name] = child
			}
		}
		return result, preserveZero || len(result) > 0
	case reflect.Map:
		if value.IsNil() || value.Len() == 0 {
			return nil, false
		}
		result := map[string]any{}
		iterator := value.MapRange()
		for iterator.Next() {
			child, include := sparseValue(iterator.Value(), reflect.Value{}, true)
			if include {
				result[fmt.Sprint(iterator.Key().Interface())] = child
			}
		}
		return result, true
	case reflect.Slice, reflect.Array:
		if value.Len() == 0 {
			return []any{}, preserveZero || !valuesEqual(value, defaultValue)
		}
		result := make([]any, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			child, include := sparseValue(value.Index(i), reflect.Value{}, true)
			if include {
				result = append(result, child)
			}
		}
		return result, true
	default:
		if !preserveZero && valuesEqual(value, defaultValue) {
			return nil, false
		}
		return value.Interface(), true
	}
}

func valuesEqual(value, defaultValue reflect.Value) bool {
	return defaultValue.IsValid() && value.Type() == defaultValue.Type() && reflect.DeepEqual(value.Interface(), defaultValue.Interface())
}

func WriteDocument(root string, doc Document) error {
	if err := validateProject(doc.Project); err != nil {
		return err
	}
	data, err := EncodeDocument(doc)
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), data, 0o644)
}

// WriteMigrationDraft writes the inactive result of legacy analysis. Only
// `wails3 migrate --activate` may rename this file into the opt-in manifest.
func WriteMigrationDraft(root string, doc Document) error {
	if err := validateProject(doc.Project); err != nil {
		return err
	}
	data, err := EncodeDocument(doc)
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, MigratedFilename), data, 0o644)
}

func Eject(root, profile, cliVersion string, backup bool) error {
	if profile != "" {
		return fmt.Errorf("wails3 eject does not accept profiles; it writes the complete resolved manifest")
	}
	if backup {
		return fmt.Errorf("wails3 eject never overwrites a file, so --backup is not supported")
	}
	loaded, err := Load(root, "")
	if err != nil {
		return err
	}
	output := filepath.Join(loaded.Config.Root, EjectedFilename)
	data, err := EncodeEjectedHCL(loaded.Config, cliVersion)
	if err != nil {
		return err
	}
	if err := exclusiveWrite(output, data, 0o644); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%s already exists; refusing to overwrite it", EjectedFilename)
		}
		return err
	}
	return nil
}

func removeProfileTargetIdentity(profile map[string]any) {
	targets, ok := profile["targets"].(map[string]any)
	if !ok {
		return
	}
	var visit func(map[string]any)
	visit = func(table map[string]any) {
		for _, key := range []string{"identifier", "product_name", "version", "build_number"} {
			delete(table, key)
		}
		for _, value := range table {
			if child, ok := value.(map[string]any); ok {
				visit(child)
			}
		}
	}
	visit(targets)
}

func clearProfileTargetIdentity(targets *Targets) {
	for _, platform := range []*Platform{&targets.Windows, &targets.Darwin, &targets.Linux, &targets.IOS, &targets.Android} {
		platform.ProductName = ""
		platform.Identifier = ""
		platform.BuildNumber = 0
		platform.AMD64.BuildNumber = 0
		platform.ARM64.BuildNumber = 0
		platform.ARM.BuildNumber = 0
		platform.X86.BuildNumber = 0
		platform.Universal.BuildNumber = 0
	}
}

func mapTable(parent map[string]any, key string) map[string]any {
	if table, ok := parent[key].(map[string]any); ok {
		return table
	}
	table := map[string]any{}
	parent[key] = table
	return table
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".wails-toml-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func exclusiveWrite(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	complete := false
	defer func() {
		if !complete {
			_ = os.Remove(path)
		}
	}()
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	complete = true
	return nil
}
