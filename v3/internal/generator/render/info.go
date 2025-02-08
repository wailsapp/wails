package render

import (
	"strings"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// modelInfo gathers useful information about a model.
type modelInfo struct {
	HasValues bool
	IsEnum    bool

	IsAlias      bool
	IsClassAlias bool
	IsTypeAlias  bool

	IsClassOrInterface bool
	IsInterface        bool
	IsClass            bool

	Template struct {
		Params     string
		ParamList  string
		CreateList string
	}
}

// modelinfo gathers and returns useful information about the given model.
func modelinfo(model *collect.ModelInfo, useInterfaces bool) (info modelInfo) {
	info.HasValues = len(model.Values) > 0
	info.IsEnum = info.HasValues && !model.Alias

	info.IsAlias = !info.IsEnum && model.Type != nil
	info.IsClassAlias = info.IsAlias && model.Predicates.IsClass && !useInterfaces
	info.IsTypeAlias = info.IsAlias && !info.IsClassAlias

	info.IsClassOrInterface = !info.IsEnum && !info.IsAlias
	info.IsInterface = info.IsClassOrInterface && (model.Alias || useInterfaces)
	info.IsClass = info.IsClassOrInterface && !info.IsInterface

	if len(model.TypeParams) > 0 {
		var params, paramList, createList strings.Builder

		paramList.WriteRune('<')
		createList.WriteRune('(')

		for i, param := range model.TypeParams {
			if i > 0 {
				params.WriteRune(',')
				paramList.WriteString(", ")
				createList.WriteString(", ")
			}
			params.WriteString(param)
			paramList.WriteString(param)

			createList.WriteString("$$createParam")
			createList.WriteString(param)
		}

		paramList.WriteRune('>')
		createList.WriteRune(')')

		info.Template.Params = params.String()
		info.Template.ParamList = paramList.String()
		info.Template.CreateList = createList.String()
	}

	return
}
