package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type schemaNode struct {
	typeInfo       reflect.Type
	attributes     map[string]schemaAttribute
	attributeOrder []string
	blocks         map[string]schemaBlock
	blockOrder     []string
	labels         []int
	labelNames     []string
}

type schemaAttribute struct {
	fieldIndex   int
	required     bool
	nonempty     bool
	defaultText  string
	defaultValue reflect.Value
	platforms    string
	formats      string
	platformMask applicabilityMask
	formatMask   applicabilityMask
	path         bool
}

type schemaBlock struct {
	fieldIndex   int
	required     bool
	repeated     bool
	node         *schemaNode
	platforms    string
	formats      string
	platformMask applicabilityMask
	formatMask   applicabilityMask
}

type applicabilityMask uint16

const (
	applicabilityWindows applicabilityMask = 1 << iota
	applicabilityDarwin
	applicabilityLinux
	applicabilityIOS
	applicabilityAndroid
	applicabilityNSIS
	applicabilityMSIX
	applicabilityDMG
	applicabilityAppImage
	applicabilityDeb
	applicabilityRPM
	applicabilityArchLinux
)

var schemaNodesByType map[reflect.Type]*schemaNode
var manifestSchema = buildManifestSchema()

func buildManifestSchema() *schemaNode {
	schemaNodesByType = make(map[reflect.Type]*schemaNode)
	return buildSchemaNode(reflect.TypeOf(manifestDocument{}))
}

func buildSchemaNode(typeInfo reflect.Type) *schemaNode {
	node := &schemaNode{typeInfo: typeInfo, attributes: make(map[string]schemaAttribute), blocks: make(map[string]schemaBlock)}
	schemaNodesByType[typeInfo] = node
	for index := 0; index < typeInfo.NumField(); index++ {
		field := typeInfo.Field(index)
		tag, tagged := field.Tag.Lookup("manifest")
		if !tagged {
			continue
		}
		name, mode, _ := strings.Cut(tag, ",")
		switch mode {
		case "label":
			node.labels = append(node.labels, index)
			node.labelNames = append(node.labelNames, field.Tag.Get("schema_label"))
		case "block":
			child := field.Type
			repeated := child.Kind() == reflect.Slice
			if repeated {
				child = child.Elem()
			}
			if child.Kind() == reflect.Pointer {
				child = child.Elem()
			}
			platforms, formats := field.Tag.Get("platforms"), field.Tag.Get("formats")
			node.blocks[name] = schemaBlock{
				fieldIndex: index, required: field.Tag.Get("required") == "true", repeated: repeated,
				node: buildSchemaNode(child), platforms: platforms, formats: formats,
				platformMask: parseApplicabilityMask(platforms), formatMask: parseApplicabilityMask(formats),
			}
		default:
			platforms, formats := field.Tag.Get("platforms"), field.Tag.Get("formats")
			descriptor := schemaAttribute{
				fieldIndex: index, required: mode != "optional" || field.Tag.Get("required") == "true",
				nonempty: field.Tag.Get("nonempty") == "true", defaultText: field.Tag.Get("default"),
				platforms: platforms, formats: formats, platformMask: parseApplicabilityMask(platforms), formatMask: parseApplicabilityMask(formats), path: field.Tag.Get("path") == "true",
			}
			if descriptor.defaultText != "" && !strings.HasPrefix(descriptor.defaultText, "$") {
				descriptor.defaultValue = decodeSchemaDefault(field.Type, descriptor.defaultText)
			}
			node.attributes[name] = descriptor
			node.attributeOrder = append(node.attributeOrder, name)
		}
		if mode == "block" {
			node.blockOrder = append(node.blockOrder, name)
		}
	}
	return node
}

func decodeSchemaDefault(destination reflect.Type, encoded string) reflect.Value {
	decoded := reflect.New(destination.Elem())
	switch decoded.Elem().Kind() {
	case reflect.String:
		decoded.Elem().SetString(encoded)
	case reflect.Int:
		parsed, err := strconv.ParseInt(encoded, 10, 64)
		if err != nil {
			panic(err)
		}
		decoded.Elem().SetInt(parsed)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(encoded)
		if err != nil {
			panic(err)
		}
		decoded.Elem().SetBool(parsed)
	case reflect.Slice, reflect.Map:
		if err := json.Unmarshal([]byte(encoded), decoded.Interface()); err != nil {
			panic(err)
		}
	default:
		panic("unsupported manifest schema default kind " + decoded.Elem().Kind().String())
	}
	return decoded
}

func decodeManifestSchema(src []byte, filename string) (manifestDocument, map[string]Origin, error) {
	decoder := yaml.NewDecoder(strings.NewReader(string(src)))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return manifestDocument{}, nil, yamlParseError(filename, err)
	}
	var trailing yaml.Node
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err != nil {
			return manifestDocument{}, nil, yamlParseError(filename, err)
		}
		return manifestDocument{}, nil, &ValidationError{Field: "manifest", Detail: "only one YAML document is allowed", Range: sourceRangeYAML(filename, &trailing)}
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return manifestDocument{}, nil, &ValidationError{Field: "manifest", Detail: "a YAML mapping is required", Range: sourceRangeYAML(filename, &document)}
	}
	root := document.Content[0]
	if err := validateYAMLTree(root, filename); err != nil {
		return manifestDocument{}, nil, err
	}
	var result manifestDocument
	origins := map[string]Origin{}
	errorsFound := validateYAMLVersion(root, filename)
	errorsFound = append(errorsFound, decodeSchemaYAML(root, manifestSchema, reflect.ValueOf(&result).Elem(), "", filename, origins)...)
	if version := mappingValue(root, "version"); version != nil && result.Version != 0 && result.Version != 3 {
		errorsFound = append(errorsFound, &ValidationError{Field: "version", Detail: "must be 3", Range: sourceRangeYAML(filename, version)})
	}
	if len(errorsFound) != 0 {
		return manifestDocument{}, nil, errors.Join(errorsFound...)
	}
	return result, origins, nil
}

func yamlParseError(filename string, err error) error {
	return &ValidationError{Field: "manifest", Detail: "parse " + filepath.Base(filename) + ": " + err.Error(), Range: SourceRange{Filename: filename, StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 2}, Cause: err}
}

func validateYAMLTree(node *yaml.Node, filename string) error {
	var errorsFound []error
	var visit func(*yaml.Node, string)
	visit = func(current *yaml.Node, field string) {
		if current == nil {
			return
		}
		if current.Kind == yaml.AliasNode || current.Alias != nil {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "YAML aliases are not supported", Range: sourceRangeYAML(filename, current)})
			return
		}
		if current.Anchor != "" {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "YAML anchors are not supported", Range: sourceRangeYAML(filename, current)})
		}
		if current.Tag == "!!null" {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "null values are not supported", Range: sourceRangeYAML(filename, current)})
		}
		if strings.HasPrefix(current.Tag, "!") && !strings.HasPrefix(current.Tag, "!!") {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "custom YAML tags are not supported", Range: sourceRangeYAML(filename, current)})
		}
		if current.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for index := 0; index+1 < len(current.Content); index += 2 {
				key, value := current.Content[index], current.Content[index+1]
				childField := schemaPath(field, key.Value)
				if key.Value == "<<" || key.Tag == "!!merge" {
					errorsFound = append(errorsFound, &ValidationError{Field: childField, Detail: "YAML merge keys are not supported", Range: sourceRangeYAML(filename, key)})
					continue
				}
				if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
					errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "mapping keys must be strings", Range: sourceRangeYAML(filename, key)})
					continue
				}
				if seen[key.Value] {
					errorsFound = append(errorsFound, &ValidationError{Field: childField, Detail: "duplicate mapping key", Range: sourceRangeYAML(filename, key)})
					continue
				}
				seen[key.Value] = true
				visit(value, childField)
			}
			return
		}
		for _, child := range current.Content {
			visit(child, field)
		}
	}
	visit(node, "")
	return errors.Join(errorsFound...)
}

func validateYAMLVersion(root *yaml.Node, filename string) []error {
	if len(root.Content) == 0 || root.Content[0].Value != "version" {
		return []error{&ValidationError{Field: "version", Detail: "version must be the first field", Range: sourceRangeYAML(filename, root)}}
	}
	return nil
}

func decodeSchemaYAML(body *yaml.Node, schema *schemaNode, target reflect.Value, parent, filename string, origins map[string]Origin) []error {
	var errorsFound []error
	seenAttributes := map[string]bool{}
	seenBlocks := map[string]bool{}
	for index := 0; index+1 < len(body.Content); index += 2 {
		key, value := body.Content[index], body.Content[index+1]
		name := key.Value
		field := schemaPath(parent, name)
		if descriptor, known := schema.attributes[name]; known {
			if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
				errorsFound = append(errorsFound, unsupportedYAMLField(field, filename, key))
				continue
			}
			seenAttributes[name] = true
			origins[field] = Origin{Kind: OriginManifest, Range: sourceRangeYAML(filename, value)}
			decoded, err := decodeYAMLAttribute(value, target.Field(descriptor.fieldIndex).Type())
			if err != nil {
				errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "unsuitable value type: " + err.Error(), Range: sourceRangeYAML(filename, value)})
				continue
			}
			element := decoded
			if element.Kind() == reflect.Pointer {
				element = element.Elem()
			}
			if descriptor.nonempty && element.Kind() == reflect.String && element.String() == "" {
				errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "must not be empty", Range: sourceRangeYAML(filename, value)})
				continue
			}
			target.Field(descriptor.fieldIndex).Set(decoded)
			continue
		}

		blockName, descriptor, known := yamlSchemaBlock(schema, name)
		if !known || !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			errorsFound = append(errorsFound, unsupportedYAMLField(field, filename, key))
			continue
		}
		seenBlocks[blockName] = true
		if value.Kind != yaml.MappingNode {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "a mapping is required", Range: sourceRangeYAML(filename, value)})
			continue
		}
		if descriptor.repeated {
			for childIndex := 0; childIndex+1 < len(value.Content); childIndex += 2 {
				label, childBody := value.Content[childIndex], value.Content[childIndex+1]
				semanticField := schemaPath(parent, yamlRepeatedName(blockName)) + "[" + strconv.Quote(label.Value) + "]"
				if label.Kind != yaml.ScalarNode || label.Tag != "!!str" {
					errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "entry names must be strings", Range: sourceRangeYAML(filename, label)})
					continue
				}
				if childBody.Kind != yaml.MappingNode {
					errorsFound = append(errorsFound, &ValidationError{Field: semanticField, Detail: "a mapping is required", Range: sourceRangeYAML(filename, childBody)})
					continue
				}
				child := reflect.New(descriptor.node.typeInfo).Elem()
				if len(descriptor.node.labels) != 1 {
					panic("YAML manifest supports exactly one label per repeated section")
				}
				child.Field(descriptor.node.labels[0]).SetString(label.Value)
				origins[semanticField] = Origin{Kind: OriginManifest, Range: sourceRangeYAML(filename, label)}
				errorsFound = append(errorsFound, decodeSchemaYAML(childBody, descriptor.node, child, semanticField, filename, origins)...)
				destination := target.Field(descriptor.fieldIndex)
				destination.Set(reflect.Append(destination, child))
			}
			continue
		}
		semanticField := schemaPath(parent, blockName)
		origins[semanticField] = Origin{Kind: OriginManifest, Range: sourceRangeYAML(filename, key)}
		child := reflect.New(descriptor.node.typeInfo).Elem()
		errorsFound = append(errorsFound, decodeSchemaYAML(value, descriptor.node, child, semanticField, filename, origins)...)
		pointer := reflect.New(descriptor.node.typeInfo)
		pointer.Elem().Set(child)
		target.Field(descriptor.fieldIndex).Set(pointer)
	}
	for _, name := range schema.attributeOrder {
		descriptor := schema.attributes[name]
		if descriptor.required && schemaFieldAllowed(parent, descriptor.platformMask) && schemaFormatAllowed(parent, descriptor.formatMask) && !seenAttributes[name] {
			errorsFound = append(errorsFound, &ValidationError{Field: schemaPath(parent, name), Detail: "required field is missing", Range: sourceRangeYAML(filename, body)})
		}
	}
	for _, name := range schema.blockOrder {
		descriptor := schema.blocks[name]
		if descriptor.required && schemaFieldAllowed(parent, descriptor.platformMask) && schemaFormatAllowed(parent, descriptor.formatMask) && !seenBlocks[name] {
			errorsFound = append(errorsFound, &ValidationError{Field: schemaPath(parent, name), Detail: "required section is missing", Range: sourceRangeYAML(filename, body)})
		}
	}
	return errorsFound
}

func unsupportedYAMLField(field, filename string, node *yaml.Node) error {
	return &ValidationError{Field: field, Detail: "unsupported field: a field with this name is not expected here", Range: sourceRangeYAML(filename, node)}
}

func yamlSchemaBlock(schema *schemaNode, yamlName string) (string, schemaBlock, bool) {
	for name, block := range schema.blocks {
		candidate := name
		if block.repeated {
			candidate = yamlRepeatedName(name)
		}
		if candidate == yamlName {
			return name, block, true
		}
	}
	return "", schemaBlock{}, false
}

func yamlRepeatedName(name string) string {
	switch name {
	case "target":
		return "targets"
	case "package":
		return "packages"
	case "profile":
		return "profiles"
	case "file_association":
		return "file_associations"
	case "protocol":
		return "protocols"
	case "hook":
		return "hooks"
	default:
		return name + "s"
	}
}

func decodeYAMLAttribute(node *yaml.Node, destination reflect.Type) (reflect.Value, error) {
	pointer := destination.Kind() == reflect.Pointer
	element := destination
	if pointer {
		element = destination.Elem()
	}
	decoded := reflect.New(element)
	switch element.Kind() {
	case reflect.String:
		if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
			return reflect.Value{}, fmt.Errorf("a string is required")
		}
		decoded.Elem().SetString(node.Value)
	case reflect.Bool:
		if node.Kind != yaml.ScalarNode || node.Tag != "!!bool" {
			return reflect.Value{}, fmt.Errorf("a boolean is required")
		}
		parsed, err := strconv.ParseBool(node.Value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("a boolean is required")
		}
		decoded.Elem().SetBool(parsed)
	case reflect.Int:
		if node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
			return reflect.Value{}, fmt.Errorf("an integer is required")
		}
		integer, err := strconv.ParseInt(node.Value, 10, 64)
		if err != nil || decoded.Elem().OverflowInt(integer) {
			return reflect.Value{}, fmt.Errorf("an integer in range is required")
		}
		decoded.Elem().SetInt(integer)
	case reflect.Slice:
		if element.Elem().Kind() != reflect.String || node.Kind != yaml.SequenceNode {
			return reflect.Value{}, fmt.Errorf("a list of strings is required")
		}
		items := reflect.MakeSlice(element, 0, len(node.Content))
		for _, item := range node.Content {
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
				return reflect.Value{}, fmt.Errorf("a list of strings is required")
			}
			items = reflect.Append(items, reflect.ValueOf(item.Value))
		}
		decoded.Elem().Set(items)
	case reflect.Map:
		if element.Key().Kind() != reflect.String || element.Elem().Kind() != reflect.String || node.Kind != yaml.MappingNode {
			return reflect.Value{}, fmt.Errorf("a map of strings is required")
		}
		items := reflect.MakeMapWithSize(element, len(node.Content)/2)
		for index := 0; index+1 < len(node.Content); index += 2 {
			key, value := node.Content[index], node.Content[index+1]
			if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
				return reflect.Value{}, fmt.Errorf("a map of strings is required")
			}
			items.SetMapIndex(reflect.ValueOf(key.Value), reflect.ValueOf(value.Value))
		}
		decoded.Elem().Set(items)
	default:
		return reflect.Value{}, fmt.Errorf("unsupported schema type %s", element)
	}
	if pointer {
		return decoded, nil
	}
	return decoded.Elem(), nil
}

func mappingValue(mapping *yaml.Node, name string) *yaml.Node {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == name {
			return mapping.Content[index+1]
		}
	}
	return nil
}

func sourceRangeYAML(filename string, node *yaml.Node) SourceRange {
	if node == nil {
		return SourceRange{Filename: filename}
	}
	endColumn := node.Column + len(node.Value)
	if endColumn <= node.Column {
		endColumn = node.Column + 1
	}
	return SourceRange{Filename: filename, StartLine: node.Line, StartColumn: node.Column, EndLine: node.Line, EndColumn: endColumn}
}

func schemaFieldAllowed(parent string, mask applicabilityMask) bool {
	if mask == 0 {
		return true
	}
	platform, _, _ := strings.Cut(parent, ".")
	return mask&applicabilityBit(platform) != 0
}

func schemaFormatAllowed(parent string, mask applicabilityMask) bool {
	if mask == 0 {
		return true
	}
	format := strings.TrimPrefix(parent, `package["`)
	if format == parent {
		format = strings.TrimPrefix(parent, `packages["`)
	}
	format, _, _ = strings.Cut(format, `"]`)
	if format == "format" {
		return true
	}
	return schemaFormatNameAllowed(format, mask)
}

func schemaFormatNameAllowed(format string, mask applicabilityMask) bool {
	return mask == 0 || mask&applicabilityBit(format) != 0
}

func parseApplicabilityMask(value string) applicabilityMask {
	var result applicabilityMask
	for _, name := range strings.Split(value, ",") {
		result |= applicabilityBit(name)
	}
	return result
}

func applicabilityBit(name string) applicabilityMask {
	switch name {
	case "windows":
		return applicabilityWindows
	case "darwin":
		return applicabilityDarwin
	case "linux":
		return applicabilityLinux
	case "ios":
		return applicabilityIOS
	case "android":
		return applicabilityAndroid
	case "nsis":
		return applicabilityNSIS
	case "msix":
		return applicabilityMSIX
	case "dmg":
		return applicabilityDMG
	case "appimage":
		return applicabilityAppImage
	case "deb":
		return applicabilityDeb
	case "rpm":
		return applicabilityRPM
	case "archlinux":
		return applicabilityArchLinux
	default:
		return 0
	}
}
