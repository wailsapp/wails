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
	}

	return false
}

// JSCreate renders JS/TS code that creates an instance
// of the given type from JSON data.
//
// JSCreate's output may be incorrect
// if m.Imports.AddType has not been called for the given type.
func (m *module) JSCreate(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			// Builtin alias: render underlying type.
			return m.JSCreate(t.Underlying())
		}

		return m.JSCreate(types.Unalias(t))

	case *types.Array:
		id, ok := m.postponedCreates.At(typ).(int)
		if !ok {
			m.JSCreate(t.Elem())
			id = m.postponedCreates.Len()
			m.postponedCreates.Set(typ, id)
		}

		return fmt.Sprintf("$$createType%d", id)

	case *types.Named:
		if t.Obj().Pkg() == nil {
			// Builtin named type: render underlying type.
			return m.JSCreate(t.Underlying())
		}

		if collect.IsAny(t) || collect.IsString(t) {
			break
		} else if !collect.IsClass(t) {
			return m.JSCreate(t.Underlying())
		}

		var builder strings.Builder

		if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
			builder.WriteString("(($$source) => ")
		}

		if t.Obj().Pkg().Path() != m.Imports.Self {
			builder.WriteString(jsimport(m.Imports.External[t.Obj().Pkg().Path()]))
			builder.WriteRune('.')
		}
		builder.WriteString(jsid(t.Obj().Name()))
		builder.WriteString(".createFrom")

		if t.TypeArgs() != nil && t.TypeArgs().Len() > 0 {
			builder.WriteString("(")
			for i, length := 0, t.TypeArgs().Len(); i < length; i++ {
				builder.WriteString(m.JSCreate(t.TypeArgs().At(i)))
				builder.WriteString(", ")
			}
			builder.WriteString("$$source))")
		}

		return builder.String()

	case *types.Map:
		id, ok := m.postponedCreates.At(typ).(int)
		if !ok {
			m.JSCreate(t.Key())
			m.JSCreate(t.Elem())
			id = m.postponedCreates.Len()
			m.postponedCreates.Set(typ, id)
		}

		return fmt.Sprintf("$$createType%d", id)

	case *types.Pointer:
		id, ok := m.postponedCreates.At(typ).(int)
		if ok {
			return fmt.Sprintf("$$createType%d", id)
		}

		createElement := m.JSCreate(t.Elem())
		if createElement != "$Create.Any" {
			id = m.postponedCreates.Len()
			m.postponedCreates.Set(typ, id)
			return fmt.Sprintf("$$createType%d", id)
		}

	case *types.Slice:
		id, ok := m.postponedCreates.At(typ).(int)
		if !ok {
			m.JSCreate(t.Elem())
			id = m.postponedCreates.Len()
			m.postponedCreates.Set(typ, id)
		}

		return fmt.Sprintf("$$createType%d", id)

	case *types.Struct:
		id, ok := m.postponedCreates.At(typ).(int)
		if ok {
			return fmt.Sprintf("$$createType%d", id)
		}

		info := m.collector.Struct(t)
		info.Collect()

		postpone := false
		for _, field := range info.Fields {
			if m.JSCreate(field.Type) != "$Create.Any" {
				postpone = true
			}
		}

		if postpone {
			id = m.postponedCreates.Len()
			m.postponedCreates.Set(typ, id)
			return fmt.Sprintf("$$createType%d", id)
		}

	case *types.TypeParam:
		if t.Obj().Name() == "" || t.Obj().Name() == "_" {
			return fmt.Sprintf("$$createT$$%d", t.Index())
		} else {
			return fmt.Sprintf("$$create%s", jsid(t.Obj().Name()))
		}
	}

	return "$Create.Any"
}

// PostponedCreates returns the list of postponed create functions
// for the given module.
func (m *module) PostponedCreates() []string {
	result := make([]string, m.postponedCreates.Len())

	m.postponedCreates.Iterate(func(key types.Type, value any) {
		id := value.(int)
		switch t := key.(type) {
		case *types.Array:
			result[id] = fmt.Sprintf("$Create.Array(%s)", m.JSCreate(t.Elem()))

		case *types.Map:
			result[id] = fmt.Sprintf("$Create.Map(%s, %s)", m.JSCreate(t.Key()), m.JSCreate(t.Elem()))

		case *types.Pointer:
			result[id] = fmt.Sprintf("$Create.Nullable(%s)", m.JSCreate(t.Elem()))

		case *types.Slice:
			result[id] = fmt.Sprintf("$Create.Array(%s)", m.JSCreate(t.Elem()))

		case *types.Struct:
			info := m.collector.Struct(t)
			info.Collect()

			var builder strings.Builder
			builder.WriteString("$Create.Struct({")

			for _, field := range info.Fields {
				createField := m.JSCreate(field.Type)
				if createField == "" {
					continue
				}

				builder.WriteString("\n    \"")
				template.JSEscape(&builder, []byte(field.Name))
				builder.WriteString("\": ")
				builder.WriteString(createField)
				builder.WriteRune(',')
			}

			if len(info.Fields) > 0 {
				builder.WriteRune('\n')
			}
			builder.WriteString("})")

			result[id] = builder.String()

		default:
			result[id] = "$Create.Any"
		}
	})

	return result
}
