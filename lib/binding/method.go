package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/wailsapp/wails/lib/logger"
)

type boundMethod struct {
	Name               string
	fullName           string
	method             reflect.Value
	inputs             []reflect.Type
	returnTypes        []reflect.Type
	log                *logger.CustomLogger
	hasErrorReturnType bool // Indicates if there is an error return type
	isWailsInit        bool
	isWailsShutdown    bool
}

// Creates a new bound method based on the given method + type
func newBoundMethod(name string, fullName string, method reflect.Value, objectType reflect.Type) (*boundMethod, error) {
	result := &boundMethod{
		Name:     name,
		method:   method,
		fullName: fullName,
	}

	// Setup logger
	result.log = logger.NewCustomLogger(result.fullName)

	// Check if Parameters are valid
	err := result.processParameters()

	// Are we a WailsInit method?
	if result.Name == "WailsInit" {
		err = result.processWailsInit()
	}

	// Are we a WailsShutdown method?
	if result.Name == "WailsShutdown" {
		err = result.processWailsShutdown()
	}

	return result, err
}

func (b *boundMethod) processParameters() error {

	// Param processing
	methodType := b.method.Type()

	// Input parameters
	inputParamCount := methodType.NumIn()
	if inputParamCount > 0 {
		b.inputs = make([]reflect.Type, inputParamCount)
		// We start at 1 as the first param is the struct
		for index := 0; index < inputParamCount; index++ {
			param := methodType.In(index)
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
	returnParamsCount := methodType.NumOut()
	// Guard against bad number of return types
	switch returnParamsCount {
	case 0:
	case 1:
		// Check if it's an error type
		param := methodType.Out(0)
		paramName := param.Name()
		if paramName == "error" {
			b.hasErrorReturnType = true
		}
		// Save return type
		b.returnTypes = append(b.returnTypes, param)
	case 2:
		// Check the second return type is an error
		secondParam := methodType.Out(1)
		secondParamName := secondParam.Name()
		if secondParamName != "error" {
			return fmt.Errorf("last return type of method '%s' must be an error (got %s)", b.Name, secondParamName)
		}

		// Check the second return type is an error
		firstParam := methodType.Out(0)
		firstParamName := firstParam.Name()
		if firstParamName == "error" {
			return fmt.Errorf("first return type of method '%s' must not be an error", b.Name)
		}
		b.hasErrorReturnType = true

		// Save return types
		b.returnTypes = append(b.returnTypes, firstParam)
		b.returnTypes = append(b.returnTypes, secondParam)

	default:
		return fmt.Errorf("cannot register method '%s' with %d return parameters. Please use up to 2", b.Name, returnParamsCount)
	}

	return nil
}

// call the method with the given data
func (b *boundMethod) call(data string) ([]reflect.Value, error) {

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
	results := b.method.Call(args)

	b.log.Debugf("results = %+v", results)
	return results, nil
}

// Attempts to set the method input <typ> for parameter <index> with the given value <val>
func (b *boundMethod) setInputValue(index int, typ reflect.Type, val interface{}) (result reflect.Value, err error) {

	// Catch type conversion panics thrown by convert
	defer func() {
		if r := recover(); r != nil {
			// Modify error
			fmt.Printf("Recovery message: %+v\n", r)
			err = fmt.Errorf("%s for parameter %d of method %s", r.(string)[23:], index+1, b.fullName)
		}
	}()

	// Do the conversion
	// Handle nil values
	if val == nil {
		switch typ.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			b.log.Debug("Converting nil to type")
			result = reflect.ValueOf(val).Convert(typ)
		default:
			b.log.Debug("Cannot convert nil to type, returning error")
			return reflect.Zero(typ), fmt.Errorf("Unable to use null value for parameter %d of method %s", index+1, b.fullName)
		}
	} else {
		result = reflect.ValueOf(val).Convert(typ)
	}

	return result, err
}

func (b *boundMethod) processWailsInit() error {
	// We must have only 1 input, it must be *wails.Runtime
	if len(b.inputs) != 1 {
		return fmt.Errorf("Invalid WailsInit() definition. Expected 1 input, but got %d", len(b.inputs))
	}

	// It must be *wails.Runtime
	inputName := b.inputs[0].String()
	b.log.Debugf("WailsInit input type: %s", inputName)
	if inputName != "*runtime.Runtime" {
		return fmt.Errorf("Invalid WailsInit() definition. Expected input to be wails.Runtime, but got %s", inputName)
	}

	// We must have only 1 output, it must be error
	if len(b.returnTypes) != 1 {
		return fmt.Errorf("Invalid WailsInit() definition. Expected 1 return type, but got %d", len(b.returnTypes))
	}

	// It must be *wails.Runtime
	outputName := b.returnTypes[0].String()
	b.log.Debugf("WailsInit output type: %s", outputName)
	if outputName != "error" {
		return fmt.Errorf("Invalid WailsInit() definition. Expected input to be error, but got %s", outputName)
	}

	// We are indeed a wails Init method
	b.isWailsInit = true

	return nil
}

func (b *boundMethod) processWailsShutdown() error {
	// We must not have any inputs
	if len(b.inputs) != 0 {
		return fmt.Errorf("Invalid WailsShutdown() definition. Expected 0 inputs, but got %d", len(b.inputs))
	}

	// We must have only 1 output, it must be error
	if len(b.returnTypes) != 0 {
		return fmt.Errorf("Invalid WailsShutdown() definition. Expected 0 return types, but got %d", len(b.returnTypes))
	}

	// We are indeed a wails Shutdown method
	b.isWailsShutdown = true

	return nil
}
