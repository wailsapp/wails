package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// JSDefault renders the Javascript representation
// of the zero value of the given type,
// using the receiver's import map to resolve dependencies.
//
// JSDefault's output may be incorrect
// if imports.AddType has not been called for the given type.
func (m *module) JSDefault(typ types.Type, quoted bool) (result string) {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return m.JSDefault(t.Underlying(), quoted)
		}

		if collect.IsClass(typ) {
			if t.Obj().Pkg().Path() == m.Imports.Self {
				prefix := ""
				if t.Obj().Exported() && m.Imports.ImportModels {
					prefix = "$models."
				} else if !t.Obj().Exported() && m.Imports.ImportInternal {
					prefix = "$internal."
				}

				return fmt.Sprintf("(new %s%s())", prefix, jsid(t.Obj().Name()))
			} else {
				return fmt.Sprintf("(new %s.%s())", jsimport(m.Imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name()))
			}
		} else {
			return m.JSDefault(types.Unalias(t), quoted)
		}

	case *types.Array:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte arrays as base64 strings
			return `""`
		} else {
			return "[]"
		}

	case *types.Basic:
		return m.renderBasicDefault(t, quoted)

	case *types.Map:
		return "{}"

	case *types.Named:
		result, ok := m.renderNamedDefault(t, quoted)
		if ok {
			return result
		}

	case *types.Pointer:
		return "null"

	case *types.Slice:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte slices as base64 strings
			return `""`
		} else {
			return "[]"
		}

	case *types.Struct:
		return m.renderStructDefault(t)
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
// of the zero value for the given named type.
// The result field named 'ok' is true when the resulting code is valid.
// If false, it must be discarded.
func (m *module) renderNamedDefault(named *types.Named, quoted bool) (result string, ok bool) {
	if named.Obj().Pkg() == nil {
		// Builtin named type: render underlying type.
		return m.JSDefault(named.Underlying(), quoted), true
	}

	if quoted {
		// WARN: Do not test with IsString here!! We only want to catch marshalers.
		if !collect.IsAny(named) && !collect.MaybeTextMarshaler(named) {
			if basic, ok := named.Underlying().(*types.Basic); ok {
				// Quoted mode for basic named type that is not a marshaler: render underlying type.
				return m.renderBasicDefault(basic, quoted), true
			}
			// No need to handle typeparams: they are initialised to null anyways.
		}
	}

	prefix := ""
	if named.Obj().Exported() && m.Imports.ImportModels {
		prefix = "$models."
	} else if !named.Obj().Exported() && m.Imports.ImportInternal {
		prefix = "$internal."
	}

	if collect.IsAny(named) {
		return "", false
	} else if collect.IsString(named) {
		return `""`, true
	} else if collect.IsClass(named) {
		if named.Obj().Pkg().Path() == m.Imports.Self {
			return fmt.Sprintf("(new %s%s())", prefix, jsid(named.Obj().Name())), true
		} else {
			return fmt.Sprintf("(new %s.%s())", jsimport(m.Imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name())), true
		}
	} else {
		// Inject a type assertion in case we are breaking an enum.
		// Using the true Go zero value is preferrable to selecting an arbitrary enum value.
		value := m.JSDefault(named.Underlying(), quoted)
		if named.Obj().Pkg().Path() == m.Imports.Self {
			if m.TS {
				return fmt.Sprintf("(%s as %s%s)", value, prefix, jsid(named.Obj().Name())), true
			} else {
				return fmt.Sprintf("(/** @type {%s%s} */(%s))", prefix, jsid(named.Obj().Name()), value), true
			}
		} else {
			if m.TS {
				return fmt.Sprintf("(%s as %s.%s)", value, jsimport(m.Imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name())), true
			} else {
				return fmt.Sprintf("(/** @type {%s.%s} */(%s))", jsimport(m.Imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name()), value), true
			}
		}
	}
}

// renderStructDefault outputs the Javascript representation
// of the zero value for the given struct type.
func (m *module) renderStructDefault(typ *types.Struct) string {
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
