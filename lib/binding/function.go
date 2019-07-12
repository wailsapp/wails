package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/wailsapp/wails/lib/logger"
)

type boundFunction struct {
	fullName           string
	function           reflect.Value
	functionType       reflect.Type
	inputs             []reflect.Type
	returnTypes        []reflect.Type
	log                *logger.CustomLogger
	hasErrorReturnType bool
}

// Creates a new bound function based on the given method + type
func newBoundFunction(object interface{}) (*boundFunction, error) {

	objectValue := reflect.ValueOf(object)
	objectType := reflect.TypeOf(object)

	name := runtime.FuncForPC(objectValue.Pointer()).Name()

	result := &boundFunction{
		fullName:     name,
		function:     objectValue,
		functionType: objectType,
		log:          logger.NewCustomLogger(name),
	}

	err := result.processParameters()

	return result, err
}

func (b *boundFunction) processParameters() error {

	// Param processing
	functionType := b.functionType

	// Input parameters
	inputParamCount := functionType.NumIn()
	if inputParamCount > 0 {
		b.inputs = make([]reflect.Type, inputParamCount)
		// We start at 1 as the first param is the struct
		for index := 0; index < inputParamCount; index++ {
			param := functionType.In(index)
			name := param.Name()
			kind := param.Kind()
			b.inputs[index] = param
			typ := param
			index := index
			b.log.DebugFields("Input param", logger.Fields{
				"index": index,
				"name":  name,
				"kind":  kind,
				"typ":   typ,
			})
		}
	}

	// Process return/output declarations
	returnParamsCount := functionType.NumOut()
	// Guard against bad number of return types
	switch returnParamsCount {
	case 0:
	case 1:
		// Check if it's an error type
		param := functionType.Out(0)
		paramName := param.Name()
		if paramName == "error" {
			b.hasErrorReturnType = true
		}
		// Save return type
		b.returnTypes = append(b.returnTypes, param)
	case 2:
		// Check the second return type is an error
		secondParam := functionType.Out(1)
		secondParamName := secondParam.Name()
		if secondParamName != "error" {
			return fmt.Errorf("last return type of method '%s' must be an error (got %s)", b.fullName, secondParamName)
		}

		// Check the second return type is an error
		firstParam := functionType.Out(0)
		firstParamName := firstParam.Name()
		if firstParamName == "error" {
			return fmt.Errorf("first return type of method '%s' must not be an error", b.fullName)
		}
		b.hasErrorReturnType = true

		// Save return types
		b.returnTypes = append(b.returnTypes, firstParam)
		b.returnTypes = append(b.returnTypes, secondParam)

	default:
		return fmt.Errorf("cannot register method '%s' with %d return parameters. Please use up to 2", b.fullName, returnParamsCount)
	}

	return nil
}

// call the method with the given data
func (b *boundFunction) call(data string) ([]reflect.Value, error) {

	// The data will be an array of values so we will decode the
	// input data into
	var jsArgs []interface{}
	d := json.NewDecoder(bytes.NewBufferString(data))
	// d.UseNumber()
	err := d.Decode(&jsArgs)
	if err != nil {
		return nil, fmt.Errorf("Invalid data passed to method call: %s", err.Error())
	}

	// Check correct number of inputs
	if len(jsArgs) != len(b.inputs) {
		return nil, fmt.Errorf("Invalid number of parameters given to %s. Expected %d but got %d", b.fullName, len(b.inputs), len(jsArgs))
	}

	// Set up call
	args := make([]reflect.Value, len(b.inputs))
	for index := 0; index < len(b.inputs); index++ {

		// Set the input values
		value, err := b.setInputValue(index, b.inputs[index], jsArgs[index])
		if err != nil {
			return nil, err
		}
		args[index] = value
	}
	b.log.Debugf("Unmarshalled Args: %+v\n", jsArgs)
	b.log.Debugf("Converted Args: %+v\n", args)
	results := b.function.Call(args)

	b.log.Debugf("results = %+v", results)
	return results, nil
}

// Attempts to set the method input <typ> for parameter <index> with the given value <val>
func (b *boundFunction) setInputValue(index int, typ reflect.Type, val interface{}) (result reflect.Value, err error) {

	// Catch type conversion panics thrown by convert
	defer func() {
		if r := recover(); r != nil {
			// Modify error
			err = fmt.Errorf("%s for parameter %d of function %s", r.(string)[23:], index+1, b.fullName)
		}
	}()

	// Translate javascript null values
	if val == nil {
		result = reflect.Zero(typ)
	} else {
		result = reflect.ValueOf(val).Convert(typ)
	}
	return result, err
}
