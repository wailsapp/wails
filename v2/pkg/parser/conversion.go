package parser

import (
	"fmt"
	"strings"

	"github.com/leaanthony/slicer"
)

// JSType represents a javascript type
type JSType string

const (
	// JsString is a JS string
	JsString JSType = "string"
	// JsBoolean is a JS bool
	JsBoolean = "boolean"
	// JsInt is a JS number
	JsInt = "number"
	// JsFloat is a JS number
	JsFloat = "number"
	// JsArray is a JS array
	JsArray = "Array"
	// JsObject is a JS object
	JsObject = "Object"
	// JsUnsupported represents a type that cannot be converted
	JsUnsupported = "*"
)

func goTypeToJS(input *Field) string {
	switch input.Type {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	// case reflect.Array, reflect.Slice:
	// 	return JsArray
	// case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Interface:
	// 	return JsObject
	case "struct":
		return input.Struct.Name
	default:
		fmt.Printf("Unsupported input to goTypeToJS: %+v", input)
		return "*"
	}
}

// goTypeToTS converts the given field into a Typescript type
// The pkgName is the package that the field is being output in.
// This is used to ensure we don't qualify local structs.
func goTypeToTS(input *Field, pkgName string) string {
	var result string
	switch input.Type {
	case "string":
		result = "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		result = "number"
	case "float32", "float64":
		result = "number"
	case "bool":
		result = "boolean"
	case "struct":
		if input.Struct.Package.Name != "" {
			if input.Struct.Package.Name != pkgName {
				result = input.Struct.Package.Name + "."
			}
		}
		result += input.Struct.Name
	// case reflect.Array, reflect.Slice:
	// 	return string(JsArray)
	// case reflect.Ptr, reflect.Struct:
	// 	fqt := input.Type().String()
	// 	return strings.Split(fqt, ".")[1]
	// case reflect.Map, reflect.Interface:
	// 	return string(JsObject)
	default:
		fmt.Printf("Unsupported input to goTypeToTS: %+v", input)
		return JsUnsupported
	}

	if input.IsArray {
		result = result + "[]"
	}

	return result
}

func goTypeToTSDeclaration(input *Field, pkgName string) string {
	var result string
	switch input.Type {
	case "string":
		result = "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		result = "number"
	case "float32", "float64":
		result = "number"
	case "bool":
		result = "boolean"
	case "struct":
		if input.Struct.Package.Name != "" {
			if input.Struct.Package.Name != pkgName {
				result = `import("./_` + input.Struct.Package.Name + `").`
			}
		}
		result += input.Struct.Name
	// case reflect.Array, reflect.Slice:
	// 	return string(JsArray)
	// case reflect.Ptr, reflect.Struct:
	// 	fqt := input.Type().String()
	// 	return strings.Split(fqt, ".")[1]
	// case reflect.Map, reflect.Interface:
	// 	return string(JsObject)
	default:
		fmt.Printf("Unsupported input to goTypeToTS: %+v", input)
		return JsUnsupported
	}

	if input.IsArray {
		result = result + "[]"
	}

	return result
}

func isUnresolvedType(typeName string) bool {
	switch typeName {
	case "string":
		return false
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return false
	case "float32", "float64":
		return false
	case "bool":
		return false
	case "struct":
		return false
	default:
		return true
	}
}

var reservedJSWords []string = []string{"abstract", "arguments", "await", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "debugger", "default", "delete", "do", "double", "else", "enum", "eval", "export", "extends", "false", "final", "finally", "float", "for", "function", "goto", "if", "implements", "import", "in", "instanceof", "int", "interface", "let", "long", "native", "new", "null", "package", "private", "protected", "public", "return", "short", "static", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "true", "try", "typeof", "var", "void", "volatile", "while", "with", "yield", "Array", "Date", "eval", "function", "hasOwnProperty", "Infinity", "isFinite", "isNaN", "isPrototypeOf", "length", "Math", "NaN", "Number", "Object", "prototype", "String", "toString", "undefined", "valueOf"}
var jsReservedWords *slicer.StringSlicer = slicer.String(reservedJSWords)

func isJSReservedWord(input string) bool {
	return jsReservedWords.Contains(input)
}

func startsWithLowerCaseLetter(input string) bool {
	firstLetter := string(input[0])
	return strings.ToLower(firstLetter) == firstLetter
}
