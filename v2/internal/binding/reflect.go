package binding

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// isStructPtr returns true if the value given is a
// pointer to a struct
func isStructPtr(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr &&
		reflect.ValueOf(value).Elem().Kind() == reflect.Struct
}

// isStructPtr returns true if the value given is a struct
func isStruct(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Struct
}

func getMethods(value interface{}) ([]*BoundMethod, error) {

	// Create result placeholder
	var result []*BoundMethod

	// Check type
	if !isStructPtr(value) {

		if isStruct(value) {
			name := reflect.ValueOf(value).Type().Name()
			return nil, fmt.Errorf("%s is a struct, not a pointer to a struct", name)
		}

		return nil, fmt.Errorf("not a pointer to a struct")
	}

	// Process Struct
	structType := reflect.TypeOf(value)
	structValue := reflect.ValueOf(value)
	baseName := structType.String()[1:]

	// Process Methods
	for i := 0; i < structType.NumMethod(); i++ {
		methodDef := structType.Method(i)
		methodName := methodDef.Name
		fullMethodName := baseName + "." + methodName
		method := structValue.MethodByName(methodName)

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
			outputs = append(outputs, thisParam)
		}
		boundMethod.Outputs = outputs

		// Save method in result
		result = append(result, boundMethod)

	}
	return result, nil
}

// convertArgToValue
func convertArgToValue(input json.RawMessage, target *Parameter) (result reflect.Value, err error) {

	// Catch type conversion panics thrown by convert
	defer func() {
		if r := recover(); r != nil {
			// Modify error
			err = fmt.Errorf("%s", r.(string)[23:])
		}
	}()

	// Do the conversion

	// Handle nil values
	if input == nil {
		switch target.reflectType.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			result = reflect.ValueOf(input).Convert(target.reflectType)
		default:
			return reflect.Zero(target.reflectType), fmt.Errorf("Unable to use null value")
		}
	} else {
		result = reflect.ValueOf(input).Convert(target.reflectType)
	}

	// We don't like doing this but it's the only way to
	// handle recover() correctly
	return

}
