package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// JSType renders a Go type to its TypeScript representation,
// using the receiver's import map to resolve dependencies.
//
// JSType's output may be incorrect if m.Imports.AddType
// has not been called for the given type.
func (m *module) JSType(typ types.Type) string {
	result, _ := m.renderType(typ, false)
	return result
}

// JSFieldType renders a struct field type to its TypeScript representation,
// using the receiver's import map to resolve dependencies.
//
// JSFieldType's output may be incorrect if m.Imports.AddType
// has not been called for the given type.
func (m *module) JSFieldType(field *collect.FieldInfo) string {
	result, _ := m.renderType(field.Type, field.Quoted)
	return result
}

// renderType provides the actual implementation of [module.Type].
// It returns the rendered type and a boolean indicating whether
// the resulting expression describes a nullable type.
func (m *module) renderType(typ types.Type, quoted bool) (result string, nullable bool) {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return m.renderType(t.Underlying(), quoted)
		}

		if quoted {
			if _, isBasic := t.Underlying().(*types.Basic); isBasic {
				switch u := types.Unalias(t).(type) {
				case *types.Basic:
					// Quoted mode for alias of basic type: render underlying type.
					return m.renderBasicType(u, quoted), false
				case *types.Named:
					// Quoted mode for alias of named type: delegate.
					return m.renderType(u, quoted)
				}
			}
		}

		if t.Obj().Pkg().Path() == m.Imports.Self {
			return jsid(t.Obj().Name()), false
		} else {
			return fmt.Sprintf("%s.%s", jsimport(m.Imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name())), false
		}

	case *types.Array:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte arrays as base64 strings
			return "string", false
		}

		elem, ptr := m.renderType(t.Elem(), false)
		if ptr {
			return fmt.Sprintf("(%s)[]", elem), false
		} else {
			return fmt.Sprintf("%s[]", elem), false
		}

	case *types.Basic:
		return m.renderBasicType(t, quoted), false

	case *types.Map:
		return m.renderMapType(t), false

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return m.renderType(t.Underlying(), quoted)
		}

		if quoted {
			if basic, ok := t.Underlying().(*types.Basic); ok && !collect.IsAny(typ) && !collect.MaybeTextMarshaler(typ) {
				// Quoted mode for basic named type that is not a marshaler: render underlying type.
				return m.renderBasicType(basic, quoted), false
			}
		}

		var builder strings.Builder

		if t.Obj().Pkg().Path() != m.Imports.Self {
			builder.WriteString(jsimport(m.Imports.External[t.Obj().Pkg().Path()]))
			builder.WriteRune('.')
		}
		builder.WriteString(jsid(t.Obj().Name()))

		if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
			builder.WriteRune('<')
			for i, length := 0, t.TypeArgs().Len(); i < length; i++ {
				if i > 0 {
					builder.WriteString(", ")
				}
				arg, _ := m.renderType(t.TypeArgs().At(i), false)
				builder.WriteString(arg)
			}
			builder.WriteRune('>')
		}

		return builder.String(), false

	case *types.Pointer:
		elem, ptr := m.renderType(t.Elem(), false)
		if ptr {
			return elem, true
		} else {
			return fmt.Sprintf("%s | null", elem), true
		}

	case *types.Slice:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte slices as base64 strings
			return "string", false
		}

		elem, ptr := m.renderType(t.Elem(), false)
		if ptr {
			return fmt.Sprintf("(%s)[]", elem), false
		} else {
			return fmt.Sprintf("%s[]", elem), false
		}

	case *types.Struct:
		return m.renderStructType(t), false

	case *types.TypeParam:
		str := ""
		if quoted {
			str = "| string "
		}

		if t.Obj().Name() == "" || t.Obj().Name() == "_" {
			return fmt.Sprintf("T$$%d %s| null", t.Index(), str), true
		} else {
			return fmt.Sprintf("%s %s| null", jsid(t.Obj().Name()), str), true
		}
	}

	// Fall back to untyped mode.
	return "any", false
}

// renderBasicType outputs the TypeScript representation
// of the given basic type.
func (*module) renderBasicType(typ *types.Basic, quoted bool) string {
	switch {
	case typ.Info()&types.IsBoolean != 0:
		if quoted {
			return "`${boolean}`"
		} else {
			return "boolean"
		}

	case typ.Info()&types.IsNumeric != 0 && typ.Info()&types.IsComplex == 0:
		if quoted {
			return "`${number}`"
		} else {
			return "number"
		}

	case typ.Info()&types.IsString != 0:
		if quoted {
			return "`\"${string}\"`"
		} else {
			return "string"
		}
	}

	// Fall back to untyped mode.
	if quoted {
		return "string"
	} else {
		return "any"
	}
}

// renderMapType outputs the TypeScript representation of the given map type.
func (m *module) renderMapType(typ *types.Map) string {
	key := "string"
	elem, _ := m.renderType(typ.Elem(), false)

	// Test whether we can upgrade key rendering.
	switch k := typ.Key().(type) {
	case *types.Basic:
		if k.Info()&types.IsString == 0 && collect.IsMapKey(k) {
			// Render non-string basic type in quoted mode.
			key = m.renderBasicType(k, true)
		}

	case *types.Alias, *types.Named:
		if collect.IsString(typ) {
			// Named type is a string alias and therefore
			// safe to use as a JS object key.
			key, _ = m.renderType(k, false)
		}

	case *types.Pointer:
		if collect.IsMapKey(typ) && collect.IsString(typ.Elem()) {
			// Base type is a string alias and therefore
			// safe to use as a JS object key.
			key, _ = m.renderType(k.Elem(), false)
		}
	}

	return fmt.Sprintf("{ [_: %s]: %s }", key, elem)
}

// renderStructType outputs the TS representation
// of the given anonymous struct type.
func (m *module) renderStructType(typ *types.Struct) string {
	info := m.collector.Struct(typ)
	info.Collect()

	var builder strings.Builder

	builder.WriteRune('{')
	for i, field := range info.Fields {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteRune('"')
		template.JSEscape(&builder, []byte(field.Name))
		builder.WriteRune('"')

		if field.Optional {
			builder.WriteRune('?')
		}

		builder.WriteString(": ")
		builder.WriteString(m.JSFieldType(field))
	}
	builder.WriteRune('}')

	return builder.String()
}
