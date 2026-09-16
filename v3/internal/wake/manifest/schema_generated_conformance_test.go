package manifest

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
			if schemaFieldAllowed(parent, descriptor.platformMask) && schemaFormatAllowed(parent, descriptor.formatMask) {
				cases = append(cases, attributeCase{path: schemaPath(parent, name), destination: node.typeInfo.Field(descriptor.fieldIndex).Type})
			}
		}
		for _, name := range node.blockOrder {
			descriptor := node.blocks[name]
			if !schemaFieldAllowed(parent, descriptor.platformMask) || !schemaFormatAllowed(parent, descriptor.formatMask) {
				continue
			}
			pathName := name
			if descriptor.repeated {
				pathName = yamlRepeatedName(name)
			}
			path := schemaPath(parent, pathName)
			for _, label := range descriptor.node.labelNames {
				path += `[` + strconv.Quote(label) + `]`
			}
			collect(descriptor.node, path)
		}
	}
	collect(manifestSchema, "")
	sort.Slice(cases, func(i, j int) bool { return cases[i].path < cases[j].path })

	var reference []SchemaField
	for _, field := range SchemaReference() {
		if field.Type != "object" && field.Type != "map(object)" {
			reference = append(reference, field)
		}
	}
	require.Len(t, reference, len(cases))
	for index, test := range cases {
		field := reference[index]
		assert.Equal(t, field.Path, test.path)
		t.Run(test.path, func(t *testing.T) {
			require.NotEmpty(t, field.Description)
			require.NotEmpty(t, field.Example)
			valid, invalid := yamlConformanceNodes(test.destination)
			decoded, err := decodeYAMLAttribute(valid, test.destination)
			require.NoError(t, err)
			assert.Equal(t, test.destination, decoded.Type())
			_, err = decodeYAMLAttribute(invalid, test.destination)
			require.Error(t, err)
		})
	}
}

func yamlConformanceNodes(destination reflect.Type) (*yaml.Node, *yaml.Node) {
	element := destination
	for element.Kind() == reflect.Pointer {
		element = element.Elem()
	}
	switch element.Kind() {
	case reflect.String:
		return yamlString("value"), yamlBool(true)
	case reflect.Int:
		return yamlInt(7), yamlString("seven")
	case reflect.Bool:
		return yamlBool(true), yamlString("true")
	case reflect.Slice:
		return yamlStrings([]string{"value"}), &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: []*yaml.Node{yamlBool(true)}}
	case reflect.Map:
		return yamlStringMap(map[string]string{"KEY": "value"}), &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{yamlString("KEY"), yamlBool(true)}}
	default:
		panic("unsupported schema conformance type " + element.String())
	}
}

func TestYAMLSchemaDecoderRejectsContainerAndIntegerEdgeCases(t *testing.T) {
	for _, test := range []struct {
		name        string
		node        *yaml.Node
		destination reflect.Type
	}{
		{"fractional integer", &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: "1.5"}, reflect.TypeOf((*int)(nil))},
		{"overflowing integer", &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "999999999999999999999999999999999999"}, reflect.TypeOf((*int)(nil))},
		{"list scalar", yamlString("value"), reflect.TypeOf((*[]string)(nil))},
		{"map tuple", yamlStrings(nil), reflect.TypeOf((*map[string]string)(nil))},
		{"null", &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, reflect.TypeOf((*string)(nil))},
		{"unsupported destination", yamlString("1.0"), reflect.TypeOf((*float64)(nil))},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := decodeYAMLAttribute(test.node, test.destination)
			require.Error(t, err)
		})
	}
}

func TestSchemaMetadataDefensiveContracts(t *testing.T) {
	for _, name := range []string{"windows", "darwin", "linux", "ios", "android", "nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"} {
		assert.NotZero(t, applicabilityBit(name), name)
	}
	assert.Zero(t, applicabilityBit("unknown"))
	assert.Panics(t, func() { schemaTypeName(reflect.TypeOf(float64(0))) })
	assert.Panics(t, func() { schemaExample("synthetic", "synthetic", reflect.TypeOf(struct{}{}), "") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*int)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*bool)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*[]string)(nil)), "nope") })
	assert.Panics(t, func() { decodeSchemaDefault(reflect.TypeOf((*float64)(nil)), "1") })
	assert.Panics(t, func() { applySchemaDefaults(nil) })
	assert.Panics(t, func() { applySchemaDefaults(&struct{}{}) })
}
