package render

import (
	"fmt"
	"go/types"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// NoCreate returns true if the given array of types needs no creation code.
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
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return m.NeedsCreate(t.Underlying())
		}

		if collect.IsClass(t) {
			return true
		} else {
			return m.NeedsCreate(types.Unalias(t))
		}

	case *types.Array, *types.Map, *types.Slice:
		return true

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return m.NeedsCreate(t.Underlying())
		}

		if collect.IsAny(t) || collect.IsString(t) {
			break
		} else if collect.IsClass(t) {
			return true
		} else {
			return m.NeedsCreate(t.Underlying())
		}

	case *types.Pointer:
		return m.NeedsCreate(t.Elem())

	case *types.Struct:
		info := m.collector.Struct(t)
		info.Collect()

		for _, field := range info.Fields {
			if m.NeedsCreate(field.Type) {
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
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return m.JSCreateWithParams(t.Underlying(), params)
		}

		return m.JSCreateWithParams(types.Unalias(t), params)

	case *types.Array:
		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if !ok {
			m.JSCreateWithParams(t.Elem(), params)
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
		}

		return fmt.Sprintf("$$createType%d%s", pp.index, params)

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

		if collect.IsAny(t) || collect.IsString(t) {
			break
		} else if !collect.IsClass(t) {
			return m.JSCreateWithParams(t.Underlying(), params)
		}

		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if !ok {
			if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
				// Postpone type args.
				for i, length := 0, t.TypeArgs().Len(); i < length; i++ {
					m.JSCreateWithParams(t.TypeArgs().At(i), params)
				}
			}

			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
		}

		return fmt.Sprintf("$$createType%d%s", pp.index, params)

	case *types.Pointer:
		pp, ok := m.postponedCreates.At(typ).(*postponed)
		if ok {
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

		createElement := m.JSCreateWithParams(t.Elem(), params)
		if createElement != "$Create.Any" {
			pp = &postponed{m.postponedCreates.Len(), params}
			m.postponedCreates.Set(typ, pp)
			return fmt.Sprintf("$$createType%d%s", pp.index, params)
		}

	case *types.Slice:
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
		case *types.Array:
			result[pp.index] = fmt.Sprintf("%s$Create.Array(%s)", pre, m.JSCreateWithParams(t.Elem(), pp.params))

		case *types.Map:
			result[pp.index] = fmt.Sprintf(
				"%s$Create.Map(%s, %s)",
				pre,
				m.JSCreateWithParams(t.Key(), pp.params),
				m.JSCreateWithParams(t.Elem(), pp.params),
			)

		case *types.Named:
			var builder strings.Builder

			builder.WriteString(pre)

			if t.Obj().Pkg().Path() == m.Imports.Self {
				if t.Obj().Exported() && m.Imports.ImportModels {
					builder.WriteString("$models.")
				} else if !t.Obj().Exported() && m.Imports.ImportInternal {
					builder.WriteString("$internal.")
				}
			} else {
				builder.WriteString(jsimport(m.Imports.External[t.Obj().Pkg().Path()]))
				builder.WriteRune('.')
			}
			builder.WriteString(jsid(t.Obj().Name()))
			builder.WriteString(".createFrom")

			if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
				builder.WriteString("(")
				for i, length := 0, t.TypeArgs().Len(); i < length; i++ {
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

		case *types.Slice:
			result[pp.index] = fmt.Sprintf("%s$Create.Array(%s)", pre, m.JSCreateWithParams(t.Elem(), pp.params))

		case *types.Struct:
			info := m.collector.Struct(t)
			info.Collect()

			var builder strings.Builder
			builder.WriteString(pre)
			builder.WriteString("$Create.Struct({")

			for _, field := range info.Fields {
				createField := m.JSCreateWithParams(field.Type, pp.params)
				if createField == "" {
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

	return result
}

type postponed struct {
	index  int
	params string
}

// hasTypeParams returns true if the given type depends upon type parameters.
func (m *module) hasTypeParams(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Array:
		return m.hasTypeParams(t.Elem())

	case *types.Map:
		return m.hasTypeParams(t.Key()) || m.hasTypeParams(t.Elem())

	case *types.Named:
		if t.TypeArgs() == nil || t.TypeArgs().Len() <= 0 {
			return false
		}

		for i, length := 0, t.TypeArgs().Len(); i < length; i++ {
			if m.hasTypeParams(t.TypeArgs().At(i)) {
				return true
			}
		}

	case *types.Pointer:
		return m.hasTypeParams(t.Elem())

	case *types.Slice:
		return m.hasTypeParams(t.Elem())

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
