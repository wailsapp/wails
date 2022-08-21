package binding

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/internal/fs"

	"github.com/leaanthony/slicer"
)

func (b *Bindings) GenerateGoBindings(baseDir string) error {
	store := b.db.store
	var obfuscatedBindings map[string]int
	if b.obfuscate {
		obfuscatedBindings = b.db.UpdateObfuscatedCallMap()
	}
	for packageName, structs := range store {
		packageDir := filepath.Join(baseDir, packageName)
		err := fs.Mkdir(packageDir)
		if err != nil {
			return err
		}
		for structName, methods := range structs {
			var jsoutput bytes.Buffer
			jsoutput.WriteString(`// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
`)
			var tsBody bytes.Buffer
			var tsContent bytes.Buffer
			tsContent.WriteString(`// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
`)
			// Sort the method names alphabetically
			methodNames := make([]string, 0, len(methods))
			for methodName := range methods {
				methodNames = append(methodNames, methodName)
			}
			sort.Strings(methodNames)

			var importNamespaces slicer.StringSlicer
			for _, methodName := range methodNames {
				// Get the method details
				methodDetails := methods[methodName]

				// Generate JS
				var args slicer.StringSlicer
				for count := range methodDetails.Inputs {
					arg := fmt.Sprintf("arg%d", count+1)
					args.Add(arg)
				}
				argsString := args.Join(", ")
				jsoutput.WriteString(fmt.Sprintf("\nexport function %s(%s) {", methodName, argsString))
				jsoutput.WriteString("\n")
				if b.obfuscate {
					id := obfuscatedBindings[strings.Join([]string{packageName, structName, methodName}, ".")]
					jsoutput.WriteString(fmt.Sprintf("  return ObfuscatedCall(%d, [%s]);", id, argsString))
				} else {
					jsoutput.WriteString(fmt.Sprintf("  return window['go']['%s']['%s']['%s'](%s);", packageName, structName, methodName, argsString))
				}
				jsoutput.WriteString("\n")
				jsoutput.WriteString(fmt.Sprintf("}"))
				jsoutput.WriteString("\n")

				// Generate TS
				tsBody.WriteString(fmt.Sprintf("\nexport function %s(", methodName))

				args.Clear()
				for count, input := range methodDetails.Inputs {
					arg := fmt.Sprintf("arg%d", count+1)
					args.Add(arg + ":" + goTypeToTypescriptType(input.TypeName, &importNamespaces))
				}
				tsBody.WriteString(args.Join(",") + "):")
				returnType := "Promise"
				if methodDetails.OutputCount() > 0 {
					firstType := goTypeToTypescriptType(methodDetails.Outputs[0].TypeName, &importNamespaces)
					returnType += "<" + firstType
					if methodDetails.OutputCount() == 2 {
						if methodDetails.Outputs[1].TypeName != "error" {
							secondType := goTypeToTypescriptType(methodDetails.Outputs[1].TypeName, &importNamespaces)
							returnType += "|" + secondType
						}
					}
					returnType += ">"
				} else {
					returnType = "void"
				}
				tsBody.WriteString(returnType + ";\n")
			}

			importNamespaces.Deduplicate()
			importNamespaces.Each(func(namespace string) {
				tsContent.WriteString("import {" + namespace + "} from '../models';\n")
			})
			tsContent.WriteString(tsBody.String())

			jsfilename := filepath.Join(packageDir, structName+".js")
			err = os.WriteFile(jsfilename, jsoutput.Bytes(), 0755)
			if err != nil {
				return err
			}
			tsfilename := filepath.Join(packageDir, structName+".d.ts")
			err = os.WriteFile(tsfilename, tsContent.Bytes(), 0755)
			if err != nil {
				return err
			}
		}
	}
	err := b.WriteModels(baseDir)
	if err != nil {
		return err
	}
	return nil
}

func goTypeToJSDocType(input string, importNamespaces *slicer.StringSlicer) string {
	switch true {
	case input == "interface {}" || input == "interface{}":
		return "any"
	case input == "string":
		return "string"
	case input == "error":
		return "Error"
	case
		strings.HasPrefix(input, "int"),
		strings.HasPrefix(input, "uint"),
		strings.HasPrefix(input, "float"):
		return "number"
	case input == "bool":
		return "boolean"
	case input == "[]byte":
		return "string"
	case strings.HasPrefix(input, "map"):
		temp := strings.TrimPrefix(input, "map[")
		// Split the string into the key and value types
		tempSplit := strings.SplitN(temp, "]", 2)
		if len(tempSplit) < 2 {
			panic("Invalid map type provided: " + input)
		}
		keyType := tempSplit[0]
		valueType := tempSplit[1]
		return fmt.Sprintf("{[key: %s]: %s}", goTypeToJSDocType(keyType, importNamespaces), goTypeToJSDocType(valueType, importNamespaces))
	case strings.HasPrefix(input, "[]"):
		arrayType := goTypeToJSDocType(input[2:], importNamespaces)
		return "Array<" + arrayType + ">"
	default:
		if strings.ContainsRune(input, '.') {
			namespace := getPackageName(input)
			importNamespaces.Add(namespace)
			return namespace + "." + strings.Split(input, ".")[1]
		}
		return "any"
	}
}

func goTypeToTypescriptType(input string, importNamespaces *slicer.StringSlicer) string {
	if strings.HasPrefix(input, "[]") {
		arrayType := goTypeToJSDocType(input[2:], importNamespaces)
		return "Array<" + arrayType + ">"
	}
	return goTypeToJSDocType(input, importNamespaces)
}
