package binding

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// isStructPtr returns true if the value given is a
// pointer to a struct
func isStructPtr(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr &&
		reflect.ValueOf(value).Elem().Kind() == reflect.Struct
}

// isFunction returns true if the given value is a function
func isFunction(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Func
}

// isStructPtr returns true if the value given is a struct
func isStruct(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Struct
}

func (b *Bindings) getMethods(value interface{}) ([]*BoundMethod, error) {
	// Create result placeholder
	var result []*BoundMethod

	// Check type
	if !isStructPtr(value) {

		if isStruct(value) {
			name := reflect.ValueOf(value).Type().Name()
			return nil, fmt.Errorf("%s is a struct, not a pointer to a struct", name)
		}

		if isFunction(value) {
			name := runtime.FuncForPC(reflect.ValueOf(value).Pointer()).Name()
			return nil, fmt.Errorf("%s is a function, not a pointer to a struct. Wails v2 has deprecated the binding of functions. Please wrap your functions up in a struct and bind a pointer to that struct.", name)
		}

		return nil, fmt.Errorf("not a pointer to a struct.")
	}

	// Process Struct
	structType := reflect.TypeOf(value)
	structValue := reflect.ValueOf(value)
	structTypeString := structType.String()
	baseName := structTypeString[1:]

	// Process Methods
	for i := 0; i < structType.NumMethod(); i++ {
		methodDef := structType.Method(i)
		methodName := methodDef.Name
		fullMethodName := baseName + "." + methodName
		method := structValue.MethodByName(methodName)

		methodReflectName := runtime.FuncForPC(methodDef.Func.Pointer()).Name()
		if b.exemptions.Contains(methodReflectName) {
			continue
		}

		// Create new method
		boundMethod := &BoundMethod{
			Name:     fullMethodName,
			Inputs:   nil,
			Outputs:  nil,
			Comments: "",
			Method:   method,
		}

		// Iterate inputs
		methodType := method.Type()
		inputParamCount := methodType.NumIn()
		var inputs []*Parameter
		for inputIndex := 0; inputIndex < inputParamCount; inputIndex++ {
			input := methodType.In(inputIndex)
			thisParam := newParameter("", input)

			thisInput := input

			if thisInput.Kind() == reflect.Slice {
				thisInput = thisInput.Elem()
			}

			// Process struct pointer params
			if thisInput.Kind() == reflect.Ptr {
				if thisInput.Elem().Kind() == reflect.Struct {
					typ := thisInput.Elem()
					a := reflect.New(typ)
					s := reflect.Indirect(a).Interface()
					name := typ.Name()
					packageName := getPackageName(thisInput.String())
					b.AddStructToGenerateTS(packageName, name, s)
				}
			}

			// Process struct params
			if thisInput.Kind() == reflect.Struct {
				a := reflect.New(thisInput)
				s := reflect.Indirect(a).Interface()
				name := thisInput.Name()
				packageName := getPackageName(thisInput.String())
				b.AddStructToGenerateTS(packageName, name, s)
			}

			inputs = append(inputs, thisParam)
		}

		boundMethod.Inputs = inputs

		// Iterate outputs
		// TODO: Determine what to do about limiting return types
		//       especially around errors.
		outputParamCount := methodType.NumOut()
		var outputs []*Parameter
		for outputIndex := 0; outputIndex < outputParamCount; outputIndex++ {
			output := methodType.Out(outputIndex)
			thisParam := newParameter("", output)

			thisOutput := output

			if thisOutput.Kind() == reflect.Slice {
				thisOutput = thisOutput.Elem()
			}

			// Process struct pointer params
			if thisOutput.Kind() == reflect.Ptr {
				if thisOutput.Elem().Kind() == reflect.Struct {
					typ := thisOutput.Elem()
					a := reflect.New(typ)
					s := reflect.Indirect(a).Interface()
					name := typ.Name()
					packageName := getPackageName(thisOutput.String())
					b.AddStructToGenerateTS(packageName, name, s)
				}
			}

			// Process struct params
			if thisOutput.Kind() == reflect.Struct {
				a := reflect.New(thisOutput)
				s := reflect.Indirect(a).Interface()
				name := thisOutput.Name()
				packageName := getPackageName(thisOutput.String())
				b.AddStructToGenerateTS(packageName, name, s)
			}

			outputs = append(outputs, thisParam)
		}
		boundMethod.Outputs = outputs

		// Save method in result
		result = append(result, boundMethod)

	}
	return result, nil
}

func getPackageName(in string) string {
	result := strings.Split(in, ".")[0]
	result = strings.ReplaceAll(result, "[]", "")
	result = strings.ReplaceAll(result, "*", "")
	return result
}

func getSplitReturn(in string) (string, string) {
	result := strings.Split(in, ".")
	return result[0], result[1]
}

func hasElements(typ reflect.Type) bool {
	kind := typ.Kind()
	return kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}
