package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// JSDefault renders the Javascript representation
// of the zero value of the given type,
// using the receiver's import map to resolve dependencies.
//
// JSDefault's output may be incorrect
// if imports.AddType has not been called for the given type.
func (m *module) JSDefault(typ types.Type, quoted bool) (result string) {
	switch t := typ.(type) {
	case *types.Alias, *types.Named:
		result, ok := m.renderNamedDefault(t.(aliasOrNamed), quoted)
		if ok {
			return result
		}

	case *types.Array:
		if t.Len() == 0 {
			return "[]"
		} else {
			// Initialise array with expected number of elements
			return fmt.Sprintf("Array.from({ length: %d }, () => %s)", t.Len(), m.JSDefault(t.Elem(), false))
		}

	case *types.Slice:
		if types.Identical(typ, typeByteSlice) {
			return `""`
		} else {
			return "[]"
		}

	case *types.Basic:
		return m.renderBasicDefault(t, quoted)

	case *types.Map:
		return "{}"

	case *types.Struct:
		return m.renderStructDefault(t)

	case *types.TypeParam:
		// Should be unreachable
		panic("type parameters have no default value")
	}

	// Fall back to null.
	// encoding/json ignores null values so this is safe.
	return "null"
}

// renderBasicDefault outputs the Javascript representation
// of the zero value for the given basic type.
func (*module) renderBasicDefault(typ *types.Basic, quoted bool) string {
	switch {
	case typ.Info()&types.IsBoolean != 0:
		if quoted {
			return `"false"`
		} else {
			return "false"
		}

	case typ.Info()&types.IsNumeric != 0 && typ.Info()&types.IsComplex == 0:
		if quoted {
			return `"0"`
		} else {
			return "0"
		}

	case typ.Info()&types.IsString != 0:
		if quoted {
			return `'""'`
		} else {
			return `""`
		}
	}

	// Fall back to untyped mode.
	if quoted {
		return `""`
	} else {
		// encoding/json ignores null values so this is safe.
		return "null"
	}
}

// renderNamedDefault outputs the Javascript representation
// of the zero value for the given alias or named type.
// The result field named 'ok' is true when the resulting code is valid.
// If false, it must be discarded.
func (m *module) renderNamedDefault(typ aliasOrNamed, quoted bool) (result string, ok bool) {
	if typ.Obj().Pkg() == nil {
		// Builtin alias or named type: render underlying type.
		return m.JSDefault(typ.Underlying(), quoted), true
	}

	if quoted {
		// WARN: Do not test with IsAny/IsStringAlias here!! We only want to catch marshalers.
		if collect.MaybeJSONMarshaler(typ) == collect.NonMarshaler && collect.MaybeTextMarshaler(typ) == collect.NonMarshaler {
			if basic, ok := typ.Underlying().(*types.Basic); ok {
				// Quoted mode for basic alias/named type that is not a marshaler: delegate.
				return m.renderBasicDefault(basic, quoted), true
			}
			// No need to handle typeparams: they are initialised to null anyways.
		}
	}

	prefix := ""
	if m.Imports.ImportModels {
		prefix = "$models."
	}

	if collect.IsAny(typ) {
		return "", false
	} else if collect.MaybeTextMarshaler(typ) != collect.NonMarshaler {
		return `""`, true
	} else if collect.IsClass(typ) && !istpalias(typ) {
		if typ.Obj().Pkg().Path() == m.Imports.Self {
			return fmt.Sprintf("(new %s%s())", prefix, jsid(typ.Obj().Name())), true
		} else {
			return fmt.Sprintf("(new %s.%s())", jsimport(m.Imports.External[typ.Obj().Pkg().Path()]), jsid(typ.Obj().Name())), true
		}
	} else if _, isAlias := typ.(*types.Alias); isAlias {
		return m.JSDefault(types.Unalias(typ), quoted), true
	} else if len(m.collector.Model(typ.Obj()).Collect().Values) > 0 {
		if typ.Obj().Pkg().Path() == m.Imports.Self {
			return fmt.Sprintf("%s%s.$zero", prefix, jsid(typ.Obj().Name())), true
		} else {
			return fmt.Sprintf("%s.%s.$zero", jsimport(m.Imports.External[typ.Obj().Pkg().Path()]), jsid(typ.Obj().Name())), true
		}
	} else {
		return m.JSDefault(typ.Underlying(), quoted), true
	}
}

// renderStructDefault outputs the Javascript representation
// of the zero value for the given struct type.
func (m *module) renderStructDefault(typ *types.Struct) string {
	if collect.MaybeJSONMarshaler(typ) != collect.NonMarshaler {
		return "null"
	} else if collect.MaybeTextMarshaler(typ) != collect.NonMarshaler {
		return `""`
	}

	info := m.collector.Struct(typ)
	info.Collect()

	var builder strings.Builder

	builder.WriteRune('{')
	for i, field := range info.Fields {
		if field.Optional {
			continue
		}

		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteRune('"')
		template.JSEscape(&builder, []byte(field.JsonName))
		builder.WriteRune('"')

		builder.WriteString(": ")

		builder.WriteString(m.JSDefault(field.Type, field.Quoted))
	}
	builder.WriteRune('}')

	return builder.String()
}
