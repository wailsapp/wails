package parser

import (
	"strings"

	"github.com/samber/lo"
)

const helperTemplate = `function {{structName}}(method) {
    return {
        packageName: "{{packageName}}",
        serviceName: "{{structName}}",
        methodName: method,
        args: Array.prototype.slice.call(arguments, 1),
    };
}
`

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
 **/
function {{methodName}}({{inputs}}) {
    return wails.Call({{structName}}("{{methodName}}"{{args}}));
}
`

func sanitiseJSVarName(name string) string {
	// if the name is a reserved word, prefix with an
	// underscore
	if strings.Contains("break,case,catch,class,const,continue,debugger,default,delete,do,else,enum,export,extends,false,finally,for,function,if,implements,import,in,instanceof,interface,let,new,null,package,private,protected,public,return,static,super,switch,this,throw,true,try,typeof,var,void,while,with,yield", name) {
		return "_" + name
	}
	return name
}

func GenerateBinding(structName string, method *BoundMethod) (string, []string) {
	var models []string
	result := strings.ReplaceAll(bindingTemplate, "{{structName}}", structName)
	result = strings.ReplaceAll(result, "{{methodName}}", method.Name)
	comments := strings.TrimSpace(method.DocComment)
	result = strings.ReplaceAll(result, "Comments", comments)
	var params string
	for _, input := range method.Inputs {
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		params += " * @param " + sanitiseJSVarName(input.Name) + " {" + input.JSType() + "}\n"
	}
	params = strings.TrimSuffix(params, "\n")
	if len(params) == 0 {
		params = " *"
	}
	////params += "\n"
	result = strings.ReplaceAll(result, " * @param name {string}", params)
	var inputs string
	for _, input := range method.Inputs {
		pkgName := getPackageName(input)
		if pkgName != "" {
			models = append(models, pkgName)
		}
		inputs += sanitiseJSVarName(input.Name) + ", "
	}
	inputs = strings.TrimSuffix(inputs, ", ")
	args := inputs
	if len(args) > 0 {
		args = ", " + args
	}
	result = strings.ReplaceAll(result, "{{inputs}}", inputs)
	result = strings.ReplaceAll(result, "{{args}}", args)

	// outputs
	var returns string
	if len(method.Outputs) == 0 {
		returns = " * @returns {Promise<void>}"
	} else {
		returns = " * @returns {Promise<"
		for _, output := range method.Outputs {
			pkgName := getPackageName(output)
			if pkgName != "" {
				models = append(models, pkgName)
			}
			jsType := output.JSType()
			if jsType == "error" {
				jsType = "void"
			}
			returns += jsType + ", "
		}
		returns = strings.TrimSuffix(returns, ", ")
		returns += ">}"
	}
	result = strings.ReplaceAll(result, " * @returns {Promise<string>}", returns)

	return result, lo.Uniq(models)
}

func getPackageName(input *Parameter) string {
	if !input.Type.IsStruct {
		return ""
	}
	result := input.Type.Package
	if result == "" {
		result = "main"
	}
	return result
}

func GenerateBindings(bindings map[string]map[string][]*BoundMethod) string {

	var result string
	var allModels []string
	for packageName, packageBindings := range bindings {
		for structName, bindings := range packageBindings {
			result += GenerateHelper(packageName, structName)
			for _, binding := range bindings {
				thisBinding, models := GenerateBinding(structName, binding)
				result += thisBinding
				allModels = append(allModels, models...)
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

	// add imports
	imports := "import {" + strings.Join(lo.Uniq(allModels), ", ") + "} from './models';\n"
	result = imports + "\n" + result

	return result
}
