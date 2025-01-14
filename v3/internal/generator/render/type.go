package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// aliasOrNamed is a common interface for *types.Alias and *types.Named.
type aliasOrNamed interface {
	types.Type
	Obj() *types.TypeName
}

// typeByteSlice caches the type-checker type for a slice of bytes.
var typeByteSlice = types.NewSlice(types.Universe.Lookup("byte").Type())

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
func (m *module) JSFieldType(field *collect.StructField) string {
	result, _ := m.renderType(field.Type, field.Quoted)
	return result
}

// renderType provides the actual implementation of [module.Type].
// It returns the rendered type and a boolean indicating whether
// the resulting expression describes a nullable type.
func (m *module) renderType(typ types.Type, quoted bool) (result string, nullable bool) {
	switch t := typ.(type) {
	case *types.Alias, *types.Named:
		return m.renderNamedType(typ.(aliasOrNamed), quoted)

	case *types.Array, *types.Slice:
		null := ""
		if _, isSlice := typ.(*types.Slice); isSlice && m.UseInterfaces {
			// In interface mode, record the fact that encoding/json marshals nil slices as null.
			null = " | null"
		}

		if types.Identical(typ, typeByteSlice) {
			// encoding/json marshals byte slices as base64 strings
			return "string" + null, null != ""
		}

		elem, ptr := m.renderType(typ.(interface{ Elem() types.Type }).Elem(), false)
		if ptr {
			return fmt.Sprintf("(%s)[]%s", elem, null), null != ""
		} else {
			return fmt.Sprintf("%s[]%s", elem, null), null != ""
		}

	case *types.Basic:
		return m.renderBasicType(t, quoted), false

	case *types.Map:
		return m.renderMapType(t)

	case *types.Pointer:
		elem, ptr := m.renderType(t.Elem(), false)
		if ptr {
			return elem, true
		} else {
			return fmt.Sprintf("%s | null", elem), true
		}

	case *types.Struct:
		return m.renderStructType(t), false

	case *types.TypeParam:
		str := ""
		if quoted {
			str = "| string "
		}
		return fmt.Sprintf("%s %s| null", typeparam(t.Index(), t.Obj().Name()), str), true
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
func (m *module) renderMapType(typ *types.Map) (result string, nullable bool) {
	null := ""
	if m.UseInterfaces {
		// In interface mode, record the fact that encoding/json marshals nil slices as null.
		null = " | null"
	}

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
		if collect.IsMapKey(typ) {
			if collect.IsString(typ) {
				// Alias or named type is a string and therefore
				// safe to use as a JS object key.
				if ptr, ok := k.(*types.Pointer); ok {
					// Unwrap pointer to named string type, but not pointer aliases.
					key, _ = m.renderType(ptr.Elem(), false)
				} else {
					key, _ = m.renderType(k, false)
				}
			} else if basic, ok := typ.Underlying().(*types.Basic); ok && basic.Info()&types.IsString == 0 {
				// Render non-string basic type in quoted mode.
				key = m.renderBasicType(basic, true)
			}
		}

	case *types.Pointer:
		if collect.IsMapKey(typ) && collect.IsString(typ.Elem()) {
			// Base type is a string alias and therefore
			// safe to use as a JS object key.
			key, _ = m.renderType(k.Elem(), false)
		}
	}

	return fmt.Sprintf("{ [_: %s]: %s }%s", key, elem, null), m.UseInterfaces
}

// renderNamedType outputs the TS representation
// of the given named or alias type.
func (m *module) renderNamedType(typ aliasOrNamed, quoted bool) (result string, nullable bool) {
	if typ.Obj().Pkg() == nil {
		// Builtin alias or named type: render underlying type.
		return m.renderType(typ.Underlying(), quoted)
	}

	if quoted {
		switch a := types.Unalias(typ).(type) {
		case *types.Basic:
			// Quoted mode for (alias of?) basic type: delegate.
			return m.renderBasicType(a, quoted), false
		case *types.TypeParam:
			// Quoted mode for (alias of?) typeparam: delegate.
			return m.renderType(a, quoted)
		case *types.Named:
			// Quoted mode for (alias of?) named type.
			// WARN: Do not test with IsString here!! We only want to catch marshalers.
			if !collect.IsAny(typ) && !collect.MaybeTextMarshaler(typ) {
				// No custom marshaling for this type.
				switch u := a.Underlying().(type) {
				case *types.Basic:
					// Quoted mode for basic named type that is not a marshaler: delegate.
					return m.renderBasicType(u, quoted), false
				case *types.TypeParam:
					// Quoted mode for generic type that maps to typeparam: delegate.
					return m.renderType(u, quoted)
				}
			}
		}
	}

	var builder strings.Builder

	if typ.Obj().Pkg().Path() == m.Imports.Self {
		if m.Imports.ImportModels {
			builder.WriteString("$models.")
		}
	} else {
		builder.WriteString(jsimport(m.Imports.External[typ.Obj().Pkg().Path()]))
		builder.WriteRune('.')
	}
	builder.WriteString(jsid(typ.Obj().Name()))

	instance, _ := typ.(interface{ TypeArgs() *types.TypeList })
	if instance != nil {
		// Render type arguments.
		if targs := instance.TypeArgs(); targs != nil && targs.Len() > 0 {
			builder.WriteRune('<')
			for i := range targs.Len() {
				if i > 0 {
					builder.WriteString(", ")
				}
				arg, _ := m.renderType(targs.At(i), false)
				builder.WriteString(arg)
			}
			builder.WriteRune('>')
		}
	}

	return builder.String(), false
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
		template.JSEscape(&builder, []byte(field.JsonName))
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
