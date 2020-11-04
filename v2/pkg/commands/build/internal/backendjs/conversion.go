package backendjs

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

func goTypeToJS(input string) JSType {
	switch input {
	case "string":
		return JsString
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return JsInt
	case "float32", "float64":
		return JsFloat
	case "bool":
		return JsBoolean
	// case reflect.Array, reflect.Slice:
	// 	return JsArray
	// case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Interface:
	// 	return JsObject
	default:
		println("UNSUPPORTED: ", input)
		return JsUnsupported
	}
}

func goTypeToTS(input *Field) string {
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
		if input.Struct.Package != "" {
			result = input.Struct.Package + "."
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
		println("UNSUPPORTED: ", input)
		return JsUnsupported
	}

	// if input.IsArray {
	// 	result = "Array<" + result + ">"
	// }

	return result
}
