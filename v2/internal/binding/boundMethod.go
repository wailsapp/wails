package binding

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// BoundMethod defines all the data related to a Go method that is
// bound to the Wails application
type BoundMethod struct {
	Name     string        `json:"name"`
	Inputs   []*Parameter  `json:"inputs,omitempty"`
	Outputs  []*Parameter  `json:"outputs,omitempty"`
	Comments string        `json:"comments,omitempty"`
	Method   reflect.Value `json:"-"`
}

// InputCount returns the number of inputs this bound method has
func (b *BoundMethod) InputCount() int {
	return len(b.Inputs)
}

// OutputCount returns the number of outputs this bound method has
func (b *BoundMethod) OutputCount() int {
	return len(b.Outputs)
}

// ParseArgs method converts the input json into the types expected by the method
func (b *BoundMethod) ParseArgs(args []json.RawMessage) ([]interface{}, error) {
	result := make([]interface{}, b.InputCount())
	if len(args) != b.InputCount() {
		return nil, fmt.Errorf("received %d arguments to method '%s', expected %d", len(args), b.Name, b.InputCount())
	}
	for index, arg := range args {
		typ := b.Inputs[index].reflectType
		inputValue := reflect.New(typ).Interface()
		err := json.Unmarshal(arg, inputValue)
		if err != nil {
			return nil, err
		}
		if inputValue == nil {
			result[index] = reflect.Zero(typ).Interface()
		} else {
			result[index] = reflect.ValueOf(inputValue).Elem().Interface()
		}
	}
	return result, nil
}

// Call will attempt to call this bound method with the given args
func (b *BoundMethod) Call(args []interface{}) (interface{}, error) {
	// Check inputs
	expectedInputLength := len(b.Inputs)
	actualInputLength := len(args)
	if expectedInputLength != actualInputLength {
		return nil, fmt.Errorf("%s takes %d inputs. Received %d", b.Name, expectedInputLength, actualInputLength)
	}

	/** Convert inputs to reflect values **/

	// Create slice for the input arguments to the method call
	callArgs := make([]reflect.Value, expectedInputLength)

	// Iterate over given arguments
	for index, arg := range args {
		// Save the converted argument
		callArgs[index] = reflect.ValueOf(arg)
	}

	// Do the call
	callResults := b.Method.Call(callArgs)

	//** Check results **//
	var returnValue interface{}
	var err error

	switch b.OutputCount() {
	case 1:
		// Loop over results and determine if the result
		// is an error or not
		for _, result := range callResults {
			interfac := result.Interface()
			temp, ok := interfac.(error)
			if ok {
				err = temp
			} else {
				returnValue = interfac
			}
		}
	case 2:
		returnValue = callResults[0].Interface()
		if temp, ok := callResults[1].Interface().(error); ok {
			err = temp
		}
	}

	return returnValue, err
}
