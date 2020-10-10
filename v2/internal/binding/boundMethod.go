package binding

import (
	"fmt"
	"reflect"
	"strings"
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

// IsWailsInit returns true if the method name is "WailsInit"
func (b *BoundMethod) IsWailsInit() bool {
	return strings.HasSuffix(b.Name, "WailsInit")
}

// IsWailsShutdown returns true if the method name is "WailsShutdown"
func (b *BoundMethod) IsWailsShutdown() bool {
	return strings.HasSuffix(b.Name, "WailsShutdown")
}

// VerifyWailsInit checks if the WailsInit signature is correct
func (b *BoundMethod) VerifyWailsInit() error {
	// Must only have 1 input
	if b.InputCount() != 1 {
		return fmt.Errorf("invalid method signature for %s: expected `WailsInit(*wails.Runtime) error`", b.Name)
	}

	// Check input type
	if !b.Inputs[0].IsType("*runtime.Runtime") {
		return fmt.Errorf("invalid method signature for %s: expected `WailsInit(*wails.Runtime) error`", b.Name)
	}

	// Must only have 1 output
	if b.OutputCount() != 1 {
		return fmt.Errorf("invalid method signature for %s: expected `WailsInit(*wails.Runtime) error`", b.Name)
	}

	// Check output type
	if !b.Outputs[0].IsError() {
		return fmt.Errorf("invalid method signature for %s: expected `WailsInit(*wails.Runtime) error`", b.Name)
	}

	// Input must be of type Runtime
	return nil
}

// VerifyWailsShutdown checks if the WailsShutdown signature is correct
func (b *BoundMethod) VerifyWailsShutdown() error {
	// Must have no inputs
	if b.InputCount() != 0 {
		return fmt.Errorf("invalid method signature for WailsShutdown: expected `WailsShutdown()`")
	}

	// Must have no outputs
	if b.OutputCount() != 0 {
		return fmt.Errorf("invalid method signature for WailsShutdown: expected `WailsShutdown()`")
	}

	// Input must be of type Runtime
	return nil
}

// InputCount returns the number of inputs this bound method has
func (b *BoundMethod) InputCount() int {
	return len(b.Inputs)
}

// OutputCount returns the number of outputs this bound method has
func (b *BoundMethod) OutputCount() int {
	return len(b.Outputs)
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

		// Attempt to convert the argument to the type expected by the method
		value, err := convertArgToValue(arg, b.Inputs[index])

		// If it fails, return a suitable error
		if err != nil {
			return nil, fmt.Errorf("%s (parameter %d): %s", b.Name, index+1, err.Error())
		}

		// Save the converted argument
		callArgs[index] = value
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
