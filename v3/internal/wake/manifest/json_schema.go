package manifest

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"

	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
)

const JSONSchemaURL = "https://v3.wails.io/schemas/wails.v3.json"

// JSONSchema returns the editor-facing JSON Schema generated from the same
// closed descriptor used by the runtime decoder.
func JSONSchema() ([]byte, error) {
	root := jsonSchemaNode(manifestSchema, "")
	root["$schema"] = "https://json-schema.org/draft/2020-12/schema"
	root["$id"] = JSONSchemaURL
	root["title"] = "Wails v3 build configuration"
	root["description"] = "Declarative build configuration for Wails v3 projects."
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func jsonSchemaNode(node *schemaNode, parent string) map[string]any {
	properties := map[string]any{}
	required := []string{}
	for _, name := range node.attributeOrder {
		descriptor := node.attributes[name]
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		path := schemaPath(parent, name)
		property := jsonSchemaAttribute(path, name, node.typeInfo.Field(descriptor.fieldIndex).Type, descriptor)
		if path == "version" {
			property["const"] = 3
		}
		properties[name] = property
		if descriptor.required {
			required = append(required, name)
		}
	}
	for _, name := range node.blockOrder {
		descriptor := node.blocks[name]
		if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			continue
		}
		path := schemaPath(parent, name)
		propertyName := name
		var property map[string]any
		if descriptor.repeated {
			propertyName = yamlRepeatedName(name)
			property = jsonSchemaNamedMap(name, descriptor.node, parent)
		} else {
			property = jsonSchemaNode(descriptor.node, path)
		}
		properties[propertyName] = property
		if descriptor.required {
			required = append(required, propertyName)
		}
	}
	result := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"properties":           properties,
	}
	if len(required) > 0 {
		result["required"] = required
	}
	return result
}

func jsonSchemaAttribute(path, name string, destination reflect.Type, descriptor schemaAttribute) map[string]any {
	element := destination
	for element.Kind() == reflect.Pointer {
		element = element.Elem()
	}
	result := map[string]any{"description": schemaDescription(path, name)}
	switch element.Kind() {
	case reflect.String:
		result["type"] = "string"
		if descriptor.nonempty {
			result["minLength"] = 1
		}
	case reflect.Int:
		result["type"] = "integer"
	case reflect.Bool:
		result["type"] = "boolean"
	case reflect.Slice:
		result["type"] = "array"
		result["items"] = map[string]any{"type": "string"}
	case reflect.Map:
		result["type"] = "object"
		result["additionalProperties"] = map[string]any{"type": "string"}
	}
	if descriptor.defaultValue.IsValid() {
		result["default"] = descriptor.defaultValue.Elem().Interface()
	}
	return result
}

func jsonSchemaNamedMap(name string, child *schemaNode, parent string) map[string]any {
	result := map[string]any{"type": "object"}
	switch name {
	case "package":
		properties := map[string]any{}
		for _, format := range []string{"nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"} {
			properties[format] = jsonSchemaNode(child, `packages[`+strconv.Quote(format)+`]`)
		}
		result["properties"] = properties
		result["additionalProperties"] = false
	case "target":
		result["additionalProperties"] = jsonSchemaNode(child, schemaPath(parent, `targets["target"]`))
		names := buildinfo.SupportedTargetNames()
		if parent == "" {
			names = append(names, "windows", "darwin", "linux", "ios", "android")
		}
		sort.Strings(names)
		result["propertyNames"] = map[string]any{"enum": names}
	case "profile":
		profile := jsonSchemaNode(child, `profiles["profile"]`)
		if properties, ok := profile["properties"].(map[string]any); ok {
			if targets, ok := properties["targets"].(map[string]any); ok {
				targets["minProperties"] = 1
			}
		}
		result["additionalProperties"] = profile
		result["propertyNames"] = map[string]any{"pattern": `^(?!default$)[a-z0-9]+(?:-[a-z0-9]+)*$`}
	case "hook":
		result["additionalProperties"] = jsonSchemaNode(child, `hooks["phase"]`)
		phases := make([]string, 0, len(HookPhases))
		for _, phase := range HookPhases {
			phases = append(phases, string(phase))
		}
		result["propertyNames"] = map[string]any{"enum": phases}
	case "file_association":
		association := jsonSchemaNode(child, `file_associations["association"]`)
		if properties, ok := association["properties"].(map[string]any); ok {
			if extensions, ok := properties["extensions"].(map[string]any); ok {
				extensions["minItems"] = 1
			}
		}
		result["additionalProperties"] = association
		result["propertyNames"] = map[string]any{"minLength": 1}
	case "protocol":
		result["additionalProperties"] = jsonSchemaNode(child, `protocols["scheme"]`)
		result["propertyNames"] = map[string]any{"minLength": 1}
	default:
		result["additionalProperties"] = jsonSchemaNode(child, schemaPath(parent, name+`["name"]`))
	}
	return result
}
