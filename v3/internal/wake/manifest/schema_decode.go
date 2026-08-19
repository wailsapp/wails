package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
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
	return buildSchemaNode(reflect.TypeOf(hclDocument{}))
}

func buildSchemaNode(typeInfo reflect.Type) *schemaNode {
	node := &schemaNode{
		typeInfo: typeInfo, attributes: make(map[string]schemaAttribute), blocks: make(map[string]schemaBlock),
	}
	schemaNodesByType[typeInfo] = node
	for index := 0; index < typeInfo.NumField(); index++ {
		field := typeInfo.Field(index)
		tag, tagged := field.Tag.Lookup("hcl")
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

func decodeManifestSchema(body *hclsyntax.Body) (hclDocument, error) {
	var result hclDocument
	errorsFound := validateVersionPosition(body)
	errorsFound = append(errorsFound, decodeSchemaBody(body, manifestSchema, reflect.ValueOf(&result).Elem(), "", body.SrcRange)...)
	if version := body.Attributes["version"]; version != nil && result.Version != 0 && result.Version != 3 {
		errorsFound = append(errorsFound, &ValidationError{Field: "version", Detail: "must be 3", Range: sourceRange(version.Range())})
	}
	if len(errorsFound) == 0 {
		return result, nil
	}
	return hclDocument{}, errors.Join(errorsFound...)
}

func validateVersionPosition(body *hclsyntax.Body) []error {
	version := body.Attributes["version"]
	if version == nil {
		return nil
	}
	first := version.Range().Start.Byte
	var offender *hclsyntax.Attribute
	for _, attribute := range body.Attributes {
		if attribute.Range().Start.Byte < first && (offender == nil || attribute.Range().Start.Byte < offender.Range().Start.Byte) {
			offender = attribute
		}
	}
	if offender != nil {
		return []error{&ValidationError{Field: "version", Detail: "version must be the first attribute", Range: sourceRange(offender.Range())}}
	}
	for _, block := range body.Blocks {
		if block.TypeRange.Start.Byte < first {
			return []error{&ValidationError{Field: "version", Detail: "version must be the first attribute", Range: sourceRange(block.TypeRange)}}
		}
	}
	return nil
}

func decodeSchemaBody(body *hclsyntax.Body, node *schemaNode, target reflect.Value, parent string, owner hcl.Range) []error {
	var errorsFound []error
	attributeNames := make([]string, 0, len(body.Attributes))
	for name := range body.Attributes {
		attributeNames = append(attributeNames, name)
	}
	sort.Slice(attributeNames, func(left, right int) bool {
		return body.Attributes[attributeNames[left]].SrcRange.Start.Byte < body.Attributes[attributeNames[right]].SrcRange.Start.Byte
	})
	seenAttributes := make(map[string]bool, len(attributeNames))
	for _, name := range attributeNames {
		attribute := body.Attributes[name]
		field := schemaPath(parent, name)
		descriptor, known := node.attributes[name]
		if !known || !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			errorsFound = append(errorsFound, &ValidationError{
				Field: field, Detail: "Unsupported argument: an argument with this name is not expected here", Range: sourceRange(attribute.NameRange),
			})
			continue
		}
		seenAttributes[name] = true
		if err := validateLiteralExpression(attribute.Expr); err != nil {
			validation := err.(*ValidationError)
			copy := *validation
			copy.Field = field
			errorsFound = append(errorsFound, &copy)
			continue
		}
		value, err := decodeSchemaAttribute(attribute.Expr, target.Field(descriptor.fieldIndex).Type())
		if err != nil {
			errorsFound = append(errorsFound, &ValidationError{
				Field: field, Detail: "Unsuitable value type: " + err.Error(), Range: sourceRange(attribute.Expr.Range()),
			})
			continue
		}
		valueElement := value
		if valueElement.Kind() == reflect.Pointer {
			valueElement = valueElement.Elem()
		}
		if descriptor.nonempty && valueElement.Kind() == reflect.String && valueElement.String() == "" {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "must not be empty", Range: sourceRange(attribute.Expr.Range())})
			continue
		}
		target.Field(descriptor.fieldIndex).Set(value)
	}
	for _, name := range node.attributeOrder {
		descriptor := node.attributes[name]
		if descriptor.required && schemaFieldAllowed(parent, descriptor.platformMask) && schemaFormatAllowed(parent, descriptor.formatMask) && !seenAttributes[name] {
			errorsFound = append(errorsFound, &ValidationError{
				Field: schemaPath(parent, name), Detail: "required field is missing", Range: sourceRange(owner),
			})
		}
	}

	seenSingletons := make(map[string]bool, len(node.blocks))
	seenLabels := make(map[string]bool)
	seenBlockTypes := make(map[string]bool, len(node.blocks))
	for _, block := range body.Blocks {
		field := schemaPath(parent, block.Type)
		descriptor, known := node.blocks[block.Type]
		if !known || !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
			errorsFound = append(errorsFound, &ValidationError{
				Field: field, Detail: "Unsupported block type: blocks of this type are not expected here", Range: sourceRange(block.TypeRange),
			})
			continue
		}
		seenBlockTypes[block.Type] = true
		for _, label := range block.Labels {
			field += "[" + strconv.Quote(label) + "]"
		}
		if len(block.Labels) != len(descriptor.node.labels) {
			errorsFound = append(errorsFound, &ValidationError{
				Field: field, Detail: fmt.Sprintf("block requires %d label(s), got %d", len(descriptor.node.labels), len(block.Labels)), Range: sourceRange(block.DefRange()),
			})
			continue
		}
		if descriptor.repeated {
			labelKey := block.Type + "\x00" + strings.Join(block.Labels, "\x00")
			if seenLabels[labelKey] {
				diagnosticRange := block.TypeRange
				if len(block.LabelRanges) > 0 {
					diagnosticRange = block.LabelRanges[0]
				}
				errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "duplicate labeled block", Range: sourceRange(diagnosticRange)})
				continue
			}
			seenLabels[labelKey] = true
		} else if seenSingletons[block.Type] {
			errorsFound = append(errorsFound, &ValidationError{Field: field, Detail: "duplicate block", Range: sourceRange(block.TypeRange)})
			continue
		} else {
			seenSingletons[block.Type] = true
		}

		child := reflect.New(descriptor.node.typeInfo).Elem()
		for index, fieldIndex := range descriptor.node.labels {
			child.Field(fieldIndex).SetString(block.Labels[index])
		}
		errorsFound = append(errorsFound, decodeSchemaBody(block.Body, descriptor.node, child, field, block.TypeRange)...)
		destination := target.Field(descriptor.fieldIndex)
		if descriptor.repeated {
			destination.Set(reflect.Append(destination, child))
		} else {
			pointer := reflect.New(descriptor.node.typeInfo)
			pointer.Elem().Set(child)
			destination.Set(pointer)
		}
	}
	for _, name := range node.blockOrder {
		descriptor := node.blocks[name]
		if descriptor.required && schemaFieldAllowed(parent, descriptor.platformMask) && schemaFormatAllowed(parent, descriptor.formatMask) && !seenBlockTypes[name] {
			errorsFound = append(errorsFound, &ValidationError{
				Field: schemaPath(parent, name), Detail: "required block is missing", Range: sourceRange(owner),
			})
		}
	}
	return errorsFound
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

func decodeSchemaAttribute(expression hcl.Expression, destination reflect.Type) (reflect.Value, error) {
	value, diagnostics := expression.Value(nil)
	if diagnostics.HasErrors() {
		return reflect.Value{}, fmt.Errorf("%s", diagnostics.Error())
	}
	if value.IsNull() || !value.IsKnown() {
		return reflect.Value{}, fmt.Errorf("null and unknown values are not supported")
	}
	pointer := destination.Kind() == reflect.Pointer
	element := destination
	if pointer {
		element = destination.Elem()
	}
	decoded := reflect.New(element)
	switch element.Kind() {
	case reflect.String:
		if value.Type() != cty.String {
			return reflect.Value{}, fmt.Errorf("a string is required")
		}
		decoded.Elem().SetString(value.AsString())
	case reflect.Bool:
		if value.Type() != cty.Bool {
			return reflect.Value{}, fmt.Errorf("a boolean is required")
		}
		decoded.Elem().SetBool(value.True())
	case reflect.Int:
		if value.Type() != cty.Number {
			return reflect.Value{}, fmt.Errorf("an integer is required")
		}
		integer, accuracy := value.AsBigFloat().Int64()
		if accuracy != big.Exact || decoded.Elem().OverflowInt(integer) {
			return reflect.Value{}, fmt.Errorf("an integer in range is required")
		}
		decoded.Elem().SetInt(integer)
	case reflect.Slice:
		if element.Elem().Kind() != reflect.String || (!value.Type().IsTupleType() && !value.Type().IsListType()) {
			return reflect.Value{}, fmt.Errorf("a list of strings is required")
		}
		items := reflect.MakeSlice(element, 0, value.LengthInt())
		iterator := value.ElementIterator()
		for iterator.Next() {
			_, item := iterator.Element()
			if item.Type() != cty.String {
				return reflect.Value{}, fmt.Errorf("a list of strings is required")
			}
			items = reflect.Append(items, reflect.ValueOf(item.AsString()))
		}
		decoded.Elem().Set(items)
	case reflect.Map:
		if element.Key().Kind() != reflect.String || element.Elem().Kind() != reflect.String || (!value.Type().IsObjectType() && !value.Type().IsMapType()) {
			return reflect.Value{}, fmt.Errorf("a map of strings is required")
		}
		items := reflect.MakeMapWithSize(element, value.LengthInt())
		iterator := value.ElementIterator()
		for iterator.Next() {
			key, item := iterator.Element()
			if key.Type() != cty.String || item.Type() != cty.String {
				return reflect.Value{}, fmt.Errorf("a map of strings is required")
			}
			items.SetMapIndex(reflect.ValueOf(key.AsString()), reflect.ValueOf(item.AsString()))
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
