package manifest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func Minimal(project Project) []byte {
	return []byte(fmt.Sprintf(`[project]
name = %q
product_name = %q
identifier = %q
version = %q
# binary_name = %q
`, project.Name, project.ProductName, project.Identifier, project.Version, deriveBinaryName(project.Name)))
}

func WriteMinimal(root string, project Project) error {
	if err := validateProject(project); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), Minimal(project), 0o644)
}

func EncodeConfig(config Config) ([]byte, error) {
	doc := Document{Project: config.Project, Frontend: config.Frontend, Build: config.Build, Dev: config.Dev, Targets: config.Targets, Package: config.Package, Signing: config.Signing, Associations: config.Associations, Protocols: config.Protocols, Hooks: config.Hooks, Wake: config.Wake, Extensions: config.Extensions}
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(doc); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func EncodeDocument(doc Document) ([]byte, error) {
	defaultDoc := defaults(Project{})
	defaultDoc.Project.BinaryName = deriveBinaryName(doc.Project.Name)
	sparse, ok := sparseValue(reflect.ValueOf(doc), reflect.ValueOf(defaultDoc), false)
	if !ok {
		return nil, fmt.Errorf("manifest document is empty")
	}
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(sparse); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
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
			child, _ := sparseValue(value.Index(i), reflect.Value{}, false)
			result = append(result, child)
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

func Eject(root, profile, cliVersion string, backup bool) error {
	loaded, err := Load(root, profile)
	if err != nil {
		return err
	}
	if profile == "" && loaded.Document.Wake.EjectedBy != "" {
		return fmt.Errorf("default configuration was already ejected by %s; edit %s directly", loaded.Document.Wake.EjectedBy, Filename)
	}
	if profile != "" && loaded.Document.Wake.EjectedProfiles[profile] != "" {
		return fmt.Errorf("profile %q was already ejected by %s; edit %s directly", profile, loaded.Document.Wake.EjectedProfiles[profile], Filename)
	}
	if backup {
		stamp := time.Now().Format("20060102-150405")
		if err := os.WriteFile(loaded.Path+"."+stamp+".bak", loaded.Raw, 0o644); err != nil {
			return fmt.Errorf("write backup: %w", err)
		}
	}
	if profile == "" {
		config := loaded.Config
		config.Wake.EjectedBy = cliVersion
		if config.Wake.EjectedProfiles == nil {
			config.Wake.EjectedProfiles = loaded.Document.Wake.EjectedProfiles
		}
		data, err := EncodeConfig(config)
		if err != nil {
			return err
		}
		// Preserve sparse named Profiles independently of the frozen base.
		var doc Document
		if _, err := toml.Decode(string(data), &doc); err != nil {
			return err
		}
		doc.Profiles = loaded.Document.Profiles
		var output bytes.Buffer
		if err := toml.NewEncoder(&output).Encode(doc); err != nil {
			return err
		}
		return atomicWrite(loaded.Path, output.Bytes(), 0o644)
	}

	base, err := Load(root, "")
	if err != nil {
		return err
	}
	frozen := profileLayerFromConfig(loaded.Config)
	var encoded bytes.Buffer
	if err := toml.NewEncoder(&encoded).Encode(frozen); err != nil {
		return err
	}
	var profileMap map[string]any
	if _, err := toml.Decode(encoded.String(), &profileMap); err != nil {
		return err
	}
	// Modify the raw sparse document so ejecting one profile does not also
	// freeze every compiled base default.
	var sparseBase map[string]any
	if _, err := toml.Decode(string(base.Raw), &sparseBase); err != nil {
		return err
	}
	profiles := mapTable(sparseBase, "profiles")
	profiles[profile] = profileMap
	wake := mapTable(sparseBase, "wake")
	ejectedProfiles := mapTable(wake, "ejected_profiles")
	ejectedProfiles[profile] = cliVersion
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(sparseBase); err != nil {
		return err
	}
	return atomicWrite(base.Path, output.Bytes(), 0o644)
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
