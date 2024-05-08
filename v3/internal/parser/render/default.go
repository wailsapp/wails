package render

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// RenderDefault renders the Javascript representation
// of the zero value of the given type,
// using the given import map to resolve dependencies.
//
// RenderDefault's output may be incorrect
// if imports.AddType has not been called for the given type.
func RenderDefault(typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool, typeScript bool) string {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return RenderDefault(t.Underlying(), imports, collector, quoted, typeScript)
		}

		if collect.IsClass(typ) {
			if t.Obj().Pkg().Path() == imports.Self {
				return fmt.Sprintf("(new %s())", jsid(t.Obj().Name()))
			} else {
				return fmt.Sprintf("(new %s.%s())", jsimport(imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name()))
			}
		} else {
			return RenderDefault(types.Unalias(t), imports, collector, quoted, typeScript)
		}

	case *types.Array:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte arrays as base64 strings
			return `""`
		} else {
			return "[]"
		}

	case *types.Basic:
		return renderBasicDefault(t, quoted)

	case *types.Map:
		return "{}"

	case *types.Named:
		result, ok := renderNamedDefault(t, imports, collector, quoted, typeScript)
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
		return renderStructDefault(t, imports, collector, typeScript)
	}

	// Fall back to null.
	// encoding/json ignores null values so this is safe.
	return "null"
}

// renderBasicDefault outputs the Javascript representation
// of the zero value for the given basic type.
func renderBasicDefault(typ *types.Basic, quoted bool) string {
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
func renderNamedDefault(named *types.Named, imports *collect.ImportMap, collector *collect.Collector, quoted bool, typeScript bool) (result string, ok bool) {
	if named.Obj().Pkg() == nil {
		// Builtin named type: render underlying type.
		return RenderDefault(named.Underlying(), imports, collector, quoted, typeScript), true
	}

	if quoted {
		if basic, ok := named.Underlying().(*types.Basic); ok && !collect.IsAny(named) && !collect.MaybeTextMarshaler(named) {
			// Quoted mode for basic named type that is not a marshaler: render underlying type.
			return renderBasicDefault(basic, quoted), true
		}
	}

	if collect.IsAny(named) {
		return "", false
	} else if collect.IsString(named) {
		return `""`, true
	} else if collect.IsClass(named) {
		if named.Obj().Pkg().Path() == imports.Self {
			return fmt.Sprintf("(new %s())", jsid(named.Obj().Name())), true
		} else {
			return fmt.Sprintf("(new %s.%s())", jsimport(imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name())), true
		}
	} else {
		// Inject a type assertion in case we are breaking an enum.
		// Using the true Go zero value is preferrable to selecting an arbitrary enum value.
		value := RenderDefault(named.Underlying(), imports, collector, quoted, typeScript)
		if named.Obj().Pkg().Path() == imports.Self {
			if typeScript {
				return fmt.Sprintf("(%s as %s)", value, jsid(named.Obj().Name())), true
			} else {
				return fmt.Sprintf("(/** @type {%s} */(%s))", jsid(named.Obj().Name()), value), true
			}
		} else {
			if typeScript {
				return fmt.Sprintf("(%s as %s.%s)", value, jsimport(imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name())), true
			} else {
				return fmt.Sprintf("(/** @type {%s.%s} */(%s))", jsimport(imports.External[named.Obj().Pkg().Path()]), jsid(named.Obj().Name()), value), true
			}
		}
	}
}

// renderStructDefault outputs the Javascript representation
// of the zero value for the given struct type.
func renderStructDefault(typ *types.Struct, imports *collect.ImportMap, collector *collect.Collector, typeScript bool) string {
	info := collector.Struct(typ)
	info.Collect()

	var builder strings.Builder
	tmplStructDefault.Execute(&builder, &struct {
		*collect.StructInfo
		Imports   *collect.ImportMap
		Collector *collect.Collector
		TS        bool
	}{
		info, imports, collector, typeScript,
	})

	return builder.String()
}
