package render

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// RenderCreate renders JS/TS code that creates an instance
// of the given type from JSON data.
//
// If condition is not empty, the output will be wrapped
// by an if statement whose guard is the specified condition.
//
// The resulting string starts always with a newline character
// unless it is empty.
//
// RenderCreate's output may be incorrect
// if imports.AddType has not been called for the given type.
func RenderCreate(condition string, target string, source string, typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool, typeScript bool, indent string) string {
	return renderConditionalCreate(0, condition, target, source, typ, imports, collector, quoted, typeScript, indent)
}

// renderConditionalCreate renders creation code
// and optionally wraps it with an if statement.
func renderConditionalCreate(depth int, condition string, target string, source string, typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool, typeScript bool, indent string) string {
	originalIndent := indent
	if condition != "" {
		indent = indent + "    "
	}

	result := renderCreate(depth, target, source, typ, imports, collector, quoted, typeScript, indent)

	if result != "" && condition != "" {
		// Wrap if statement around result.
		return fmt.Sprintf("\n%sif (%s) {%s\n%s}", originalIndent, condition, result, originalIndent)
	}

	return result
}

// renderCreate renders unconditional creation code.
// The depth parameter is used to avoid naming collisions between temporary variables.
func renderCreate(depth int, target string, source string, typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool, typeScript bool, indent string) string {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return renderCreate(depth, target, source, t.Underlying(), imports, collector, quoted, typeScript, indent)
		}

		if collect.IsClass(t) {
			if t.Obj().Pkg().Path() == imports.Self {
				return fmt.Sprintf("\n%s%s = %s.createFrom(%s);", indent, target, jsid(t.Obj().Name()), source)
			} else {
				return fmt.Sprintf("\n%s%s = %s.%s.createFrom(%s);", indent, target, jsimport(imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name()), source)
			}
		} else {
			return renderCreate(depth, target, source, types.Unalias(t), imports, collector, quoted, typeScript, indent)
		}

	case *types.Array:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte arrays as base64 strings
			return fmt.Sprintf("\n%s%s = (%s === null) ? \"\" : %s;", indent, target, source, source)
		}

		tmp := fmt.Sprintf("index$$%d", depth)
		elTgt, elSrc := fmt.Sprintf("%s[%s]", target, tmp), fmt.Sprintf("%s[%s]", source, tmp)
		createElement := renderCreate(depth+1, elTgt, elSrc, t.Elem(), imports, collector, false, typeScript, indent+"    ")
		// Avoid unnecessary work.
		if createElement == "" {
			return fmt.Sprintf("\n%s%s = (%s === null) ? [] : %s;", indent, target, source, source)
		} else {
			return fmt.Sprintf(
				"\n%sif (%s === null) {\n%s    %s = [];\n%s} else for (const %s in %s) {%s\n%s}",
				indent, source, indent, target, indent, tmp, source, createElement, indent,
			)
		}

	case *types.Basic:
		// No need to convert basic types.
		return ""

	case *types.Map:
		tmp := fmt.Sprintf("key$$%d", depth)
		elTgt, elSrc := fmt.Sprintf("%s[%s]", target, tmp), fmt.Sprintf("%s[%s]", source, tmp)
		createElement := renderCreate(depth+1, elTgt, elSrc, t.Elem(), imports, collector, false, typeScript, indent+"    ")
		// Avoid unnecessary work.
		if createElement != "" {
			return fmt.Sprintf("\n%sfor (const %s in %s) {%s\n%s}", indent, tmp, source, createElement, indent)
		}

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return renderCreate(depth, target, source, t.Underlying(), imports, collector, quoted, typeScript, indent)
		}

		if collect.IsAny(t) || collect.IsString(t) {
			break
		} else if collect.IsClass(t) {
			if t.Obj().Pkg().Path() == imports.Self {
				return fmt.Sprintf("\n%s%s = %s.createFrom(%s);", indent, target, jsid(t.Obj().Name()), source)
			} else {
				return fmt.Sprintf("\n%s%s = %s.%s.createFrom(%s);", indent, target, jsimport(imports.External[t.Obj().Pkg().Path()]), jsid(t.Obj().Name()), source)
			}
		} else {
			return renderCreate(depth, target, source, t.Underlying(), imports, collector, quoted, typeScript, indent)
		}

	case *types.Pointer:
		createElement := renderCreate(depth, target, source, t.Elem(), imports, collector, false, typeScript, indent+"    ")
		if createElement != "" {
			prepare := ""
			if target != source {
				prepare = fmt.Sprintf("\n%s%s = null;", indent, target)
			}
			return fmt.Sprintf("%s\n%sif (%s !== null) {%s\n%s}", prepare, indent, source, createElement, indent)
		}

	case *types.Slice:
		if types.Identical(t.Elem(), types.Universe.Lookup("byte").Type()) {
			// encoding/json marshals byte arrays as base64 strings
			return fmt.Sprintf("\n%s%s = (%s === null) ? \"\" : %s;", indent, target, source, source)
		}

		tmp := fmt.Sprintf("index$$%d", depth)
		elTgt, elSrc := fmt.Sprintf("%s[%s]", target, tmp), fmt.Sprintf("%s[%s]", source, tmp)
		createElement := renderCreate(depth+1, elTgt, elSrc, t.Elem(), imports, collector, false, typeScript, indent+"    ")
		// Avoid unnecessary work.
		if createElement == "" {
			return fmt.Sprintf("\n%s%s = (%s === null) ? [] : %s;", indent, target, source, source)
		} else {
			return fmt.Sprintf(
				"\n%sif (%s === null) {\n%s    %s = [];\n%s} else for (const %s in %s) {%s\n%s}",
				indent, source, indent, target, indent, tmp, source, createElement, indent,
			)
		}

	case *types.Struct:
		return renderStructCreate(depth, target, source, t, imports, collector, typeScript, indent)
	}

	if target == source {
		// Do not render type assertion for self assignment.
		return ""
	}

	// Fall back to type assertion.
	jstype, _ := renderType(typ, imports, collector, quoted)
	if typeScript {
		return fmt.Sprintf("\n%s%s = %s as %s;", indent, target, source, jstype)
	} else {
		return fmt.Sprintf("\n%s%s = /** @type {%s} */(%s);", indent, target, jstype, source)
	}
}

// renderStructDefault Javascript creation code for the given struct type.
func renderStructCreate(depth int, target string, source string, typ *types.Struct, imports *collect.ImportMap, collector *collect.Collector, typeScript bool, indent string) string {
	info := collector.Struct(typ)
	info.Collect()

	var builder strings.Builder
	tmplStructCreate.Execute(&builder, &struct {
		*collect.StructInfo
		Imports   *collect.ImportMap
		Collector *collect.Collector
		Depth     int
		Target    string
		Source    string
		TS        bool
		Indent    string
	}{
		info, imports, collector,
		depth, target, source, typeScript, indent,
	})

	return builder.String()
}
