package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// SkipCreate returns true if the given array of types needs no creation code.
func (m *module) SkipCreate(ts []types.Type) bool {
	for _, typ := range ts {
		if m.NeedsCreate(typ) {
			return false
		}
	}
	return true
}

// NeedsCreate returns true if the given type needs some creation code.
func (m *module) NeedsCreate(typ types.Type) bool {
	return m.needsCreateImpl(typ, make(map[*types.TypeName]bool))
}

// needsCreateImpl provides the actual implementation of NeedsCreate.
// The visited parameter is used to break cycles.
func (m *module) needsCreateImpl(typ types.Type, visited map[*types.TypeName]bool) bool {
	switch t := typ.(type) {
	case *types.Alias, *types.Named:
		obj := typ.(interface{ Obj() *types.TypeName }).Obj()
		if visited[obj] {
			return false
		}
		visited[obj] = true

		if obj.Pkg() == nil {
			// Builtin alias or named type: render underlying type.
			return m.needsCreateImpl(t.Underlying(), visited)
		}

		if collect.IsAny(t) || collect.IsString(t) {
			break
		} else if collect.IsClass(t) {
			return true
		} else if _, isAlias := typ.(*types.Alias); isAlias {
			return m.needsCreateImpl(types.Unalias(t), visited)
		} else {
			return m.needsCreateImpl(t.Underlying(), visited)
		}

	case *types.Array, *types.Pointer:
		return m.needsCreateImpl(typ.(interface{ Elem() types.Type }).Elem(), visited)

	case *types.Map, *types.Slice:
		return true

	case *types.Struct:
		info := m.collector.Struct(t)
		info.Collect()

		for _, field := range info.Fields {
			if m.needsCreateImpl(field.Type, visited) {
				return true
			}
		}

	case *types.TypeParam:
		return true
	}

	return false
}

// JSCreate renders JS/TS code that creates an instance
// of the given type from JSON data.
//
// JSCreate's output may be incorrect
// if m.Imports.AddType has not been called for the given type.
func (m *module) JSCreate(typ types.Type) string {
	return m.JSCreateWithParams(typ, "")
}

// JSCreateWithParams renders JS/TS code that creates an instance
// of the given type from JSON data. For generic types,
// it renders parameterised code.
//
// JSCreateWithParams's output may be incorrect
// if m.Imports.AddType has not been called for the given type.
func (m *module) JSCreateWithParams(typ types.Type, params string) string {
	if len(params) > 0 && !m.hasTypeParams(typ) {
		// Forget params for non-generic types.
		params = ""
	}

	switch t := typ.(type) {
	case *types.Alias:
		return m.JSCreateWithParams(types.Unalias(typ), params)

	case *types.Array, *types.Pointer:
		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if ok {
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

		createElement := m.JSCreateWithParams(typ.(interface{ Elem() types.Type }).Elem(), params)
		if createElement != "$Create.Any" {
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

	case *types.Map:
		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if !ok {
			m.JSCreateWithParams(t.Key(), params)
			m.JSCreateWithParams(t.Elem(), params)
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
		}

		return fmt.Sprintf("$$createType%d%s", pp.index, params)

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return m.JSCreateWithParams(t.Underlying(), params)
		}

		if collect.IsAny(typ) || collect.IsString(typ) || !m.NeedsCreate(typ) {
			break
		}

		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if !ok {
			if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
				// Postpone type args.
				for i := range t.TypeArgs().Len() {
					m.JSCreateWithParams(t.TypeArgs().At(i), params)
				}
			}

			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)

			if !collect.IsClass(typ) {
				m.JSCreateWithParams(t.Underlying(), params)
			}
		}

		return fmt.Sprintf("$$createType%d%s", pp.index, params)

	case *types.Slice:
		if types.Identical(typ, typeByteSlice) {
			return "$Create.ByteSlice"
		}

		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if !ok {
			m.JSCreateWithParams(t.Elem(), params)
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
		}

		return fmt.Sprintf("$$createType%d%s", pp.index, params)

	case *types.Struct:
		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if ok {
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

		info := m.collector.Struct(t)
		info.Collect()

		postpone := false
		for _, field := range info.Fields {
			if m.JSCreateWithParams(field.Type, params) != "$Create.Any" {
				postpone = true
			}
		}

		if postpone {
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

	case *types.TypeParam:
		return fmt.Sprintf("$$createParam%s", typeparam(t.Index(), t.Obj().Name()))
	}

	return "$Create.Any"
}

// PostponedCreates returns the list of postponed create functions
// for the given module.
func (m *module) PostponedCreates() []string {
	result := make([]string, m.postponedCreates.Len())

	m.postponedCreates.Iterate(func(key types.Type, value any) {
		pp := value.(*postponed)

		pre := ""
		if pp.params != "" {
			pre = pp.params + " => "
		}

		switch t := key.(type) {
		case *types.Array, *types.Slice:
			result[pp.index] = fmt.Sprintf("%s$Create.Array(%s)", pre, m.JSCreateWithParams(t.(interface{ Elem() types.Type }).Elem(), pp.params))

		case *types.Map:
			result[pp.index] = fmt.Sprintf(
				"%s$Create.Map(%s, %s)",
				pre,
				m.JSCreateWithParams(t.Key(), pp.params),
				m.JSCreateWithParams(t.Elem(), pp.params),
			)

		case *types.Named:
			if !collect.IsClass(key) {
				// Creation functions for non-struct named types
				// require an indirect assignment to break cycles.

				// Typescript cannot infer the return type on its own: add hints.
				cast, returnType := "", ""
				if m.TS {
					returnType = ": any"
				} else {
					cast = "/** @type {(...args: any[]) => any} */"
				}

				result[pp.index] = fmt.Sprintf(`
%s(function $$initCreateType%d(...args)%s {
    if ($$createType%d === $$initCreateType%d) {
        $$createType%d = %s%s;
    }
    return $$createType%d(...args);
})`,
					cast, pp.index, returnType,
					pp.index, pp.index,
					pp.index, pre, m.JSCreateWithParams(t.Underlying(), pp.params),
					pp.index,
				)[1:] // Remove initial newline.

				// We're done.
				break
			}

			var builder strings.Builder

			builder.WriteString(pre)

			if t.Obj().Pkg().Path() == m.Imports.Self {
				if m.Imports.ImportModels {
					builder.WriteString("$models.")
				}
			} else {
				builder.WriteString(jsimport(m.Imports.External[t.Obj().Pkg().Path()]))
				builder.WriteRune('.')
			}
			builder.WriteString(jsid(t.Obj().Name()))
			builder.WriteString(".createFrom")

			if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
				builder.WriteString("(")
				for i := range t.TypeArgs().Len() {
					if i > 0 {
						builder.WriteString(", ")
					}
					builder.WriteString(m.JSCreateWithParams(t.TypeArgs().At(i), pp.params))
				}
				builder.WriteString(")")
			}

			result[pp.index] = builder.String()

		case *types.Pointer:
			result[pp.index] = fmt.Sprintf("%s$Create.Nullable(%s)", pre, m.JSCreateWithParams(t.Elem(), pp.params))

		case *types.Struct:
			info := m.collector.Struct(t)
			info.Collect()

			var builder strings.Builder
			builder.WriteString(pre)
			builder.WriteString("$Create.Struct({")

			for _, field := range info.Fields {
				createField := m.JSCreateWithParams(field.Type, pp.params)
				if createField == "$Create.Any" {
					continue
				}

				builder.WriteString("\n    \"")
				template.JSEscape(&builder, []byte(field.JsonName))
				builder.WriteString("\": ")
				builder.WriteString(createField)
				builder.WriteRune(',')
			}

			if len(info.Fields) > 0 {
				builder.WriteRune('\n')
			}
			builder.WriteString("})")

			result[pp.index] = builder.String()

		default:
			result[pp.index] = pre + "$Create.Any"
		}
	})

	if newline != "\n" {
		// Replace newlines according to local git config.
		for i := range result {
			result[i] = strings.ReplaceAll(result[i], "\n", newline)
		}
	}

	return result
}

type postponed struct {
	index  int
	params string
}

// hasTypeParams returns true if the given type depends upon type parameters.
func (m *module) hasTypeParams(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: these are never rendered as templates.
			return false
		}

		return m.hasTypeParams(types.Unalias(typ))

	case *types.Array, *types.Pointer, *types.Slice:
		return m.hasTypeParams(typ.(interface{ Elem() types.Type }).Elem())

	case *types.Map:
		return m.hasTypeParams(t.Key()) || m.hasTypeParams(t.Elem())

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: these are never rendered as templates.
			return false
		}

		if targs := t.TypeArgs(); targs != nil {
			for i := range targs.Len() {
				if m.hasTypeParams(targs.At(i)) {
					return true
				}
			}
		}

	case *types.Struct:
		info := m.collector.Struct(t)
		info.Collect()

		for _, field := range info.Fields {
			if m.hasTypeParams(field.Type) {
				return true
			}
		}

	case *types.TypeParam:
		return true
	}

	return false
}
