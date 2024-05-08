package render

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// RenderType renders a Go type to its TypeScript representation,
// using the given import map to resolve dependencies.
//
// RenderType's output may be incorrect
// if imports.AddType has not been called for the given type.
func RenderType(typ types.Type, imports *collect.ImportMap, collector *collect.Collector) string {
	result, _ := renderType(typ, imports, collector, false)
	return result
}

// renderType provides the actual implementation of [RenderType].
// It returns the rendered type and a boolean indicating whether
// the resulting expression describes a pointer type.
func renderType(typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool) (result string, ptr bool) {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return renderType(t.Underlying(), imports, collector, quoted)
		}

		if quoted {
			if _, isBasic := t.Underlying().(*types.Basic); isBasic {
				switch u := types.Unalias(t).(type) {
				case *types.Basic:
					// Quoted mode for alias of basic type: render underlying type.
					return renderBasicType(u, quoted), false
				case *types.Named:
					// Quoted mode for alias of named type: delegate.
					return renderType(u, imports, collector, quoted)
				}
			}
		}

		if t.Obj().Pkg().Path() == imports.Self {
			return jsid(t.Obj().Name()), false
		} else {
			return fmt.Sprintf("%s.%s", jsimport(imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name())), false
		}

	case *types.Array:
		elem, ptr := renderType(t.Elem(), imports, collector, false)
		if ptr {
			return fmt.Sprintf("(%s)[]", elem), false
		} else {
			return fmt.Sprintf("%s[]", elem), false
		}

	case *types.Basic:
		return renderBasicType(t, quoted), false

	case *types.Map:
		return renderMapType(t, imports, collector), false

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return renderType(t.Underlying(), imports, collector, quoted)
		}

		if quoted {
			if basic, ok := t.Underlying().(*types.Basic); ok && !collect.IsAny(typ) && !collect.MaybeTextMarshaler(typ) {
				// Quoted mode for basic named type that is not a marshaler: render underlying type.
				return renderBasicType(basic, quoted), false
			}
		}

		if t.Obj().Pkg().Path() == imports.Self {
			return jsid(t.Obj().Name()), false
		} else {
			return fmt.Sprintf("%s.%s", jsimport(imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name())), false
		}

	case *types.Pointer:
		elem, ptr := renderType(t.Elem(), imports, collector, false)
		if ptr {
			return elem, true
		} else {
			return fmt.Sprintf("%s | null", elem), true
		}

	case *types.Slice:
		elem, ptr := renderType(t.Elem(), imports, collector, false)
		if ptr {
			return fmt.Sprintf("(%s)[]", elem), false
		} else {
			return fmt.Sprintf("%s[]", elem), false
		}

	case *types.Struct:
		return renderStructType(t, imports, collector), false
	}

	// Fall back to untyped mode.
	return "any", false
}

// renderBasicType outputs the TypeScript representation
// of the given basic type.
func renderBasicType(typ *types.Basic, quoted bool) string {
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
func renderMapType(typ *types.Map, imports *collect.ImportMap, collector *collect.Collector) string {
	key := "string"
	elem, _ := renderType(typ.Elem(), imports, collector, false)

	// Test whether we can upgrade key rendering.
	switch k := typ.Key().(type) {
	case *types.Basic:
		if k.Info()&types.IsString == 0 && collect.IsMapKey(k) {
			// Render non-string basic type in quoted mode.
			key = renderBasicType(k, true)
		}

	case *types.Alias, *types.Named:
		if collect.IsString(typ) {
			// Named type is a string alias and therefore
			// safe to use as a JS object key.
			key, _ = renderType(k, imports, collector, false)
		}

	case *types.Pointer:
		if collect.IsMapKey(typ) && collect.IsString(typ.Elem()) {
			// Base type is a string alias and therefore
			// safe to use as a JS object key.
			key, _ = renderType(k.Elem(), imports, collector, false)
		}
	}

	return fmt.Sprintf("{ [_: %s]: %s }", key, elem)
}

// renderStructType outputs the TS representation
// of the given anonymous struct type.
func renderStructType(typ *types.Struct, imports *collect.ImportMap, collector *collect.Collector) string {
	info := collector.Struct(typ)
	info.Collect()

	var builder strings.Builder
	tmplStructType.Execute(&builder, &struct {
		*collect.StructInfo
		Imports   *collect.ImportMap
		Collector *collect.Collector
	}{
		info, imports, collector,
	})

	return builder.String()
}
