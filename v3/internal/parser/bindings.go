package parser

import (
	"strings"
)

const helperTemplate = `function {{structName}}(method) {
    return {
        packageName: "{{packageName}}",
        serviceName: "{{structName}}",
        methodName: method,
        args: Array.prototype.slice.call(arguments, 1),
    };
}`

func GenerateHelper(packageName, structName string) string {
	result := strings.ReplaceAll(helperTemplate, "{{packageName}}", packageName)
	result = strings.ReplaceAll(result, "{{structName}}", structName)
	return result
}

const bindingTemplate = `

/**
 * {{structName}}.{{methodName}}
 * Comments
 * @param name {string}
 * @returns {Promise<string>}
 */
function {{methodName}}({{args}}) {
    return wails.Call({{structName}}("{{methodName}}", {{args}}));
}
`

func GenerateBinding(structName string, method *BoundMethod) string {
	result := strings.ReplaceAll(bindingTemplate, "{{structName}}", structName)
	result = strings.ReplaceAll(result, "{{methodName}}", method.Name)
	result = strings.ReplaceAll(result, "Comments", strings.TrimSpace(method.DocComment))
	var params string
	for _, input := range method.Inputs {
		params += " * @param " + input.Name + " {" + input.JSType() + "}\n"
	}
	params = strings.TrimSuffix(params, "\n")
	result = strings.ReplaceAll(result, " * @param name {string}", params)
	var args string
	for _, input := range method.Inputs {
		args += input.Name + ", "
	}
	args = strings.TrimSuffix(args, ", ")
	result = strings.ReplaceAll(result, "{{args}}", args)
	return result
}

func GenerateBindings(bindings map[string]map[string][]*BoundMethod) string {

	var result string
	for packageName, packageBindings := range bindings {
		for structName, bindings := range packageBindings {
			result += GenerateHelper(packageName, structName)
			for _, binding := range bindings {
				result += GenerateBinding(structName, binding)
			}
		}
	}
	result += `
window.go = window.go || {};
`
	for packageName, packageBindings := range bindings {
		result += "Object.window.go." + packageName + " = {\n"
		for structName, methods := range packageBindings {
			result += "    " + structName + ": {\n"
			for _, method := range methods {
				result += "        " + method.Name + ",\n"
			}
			result += "    }\n"
		}
		result += "};\n"
	}
	return result
}
