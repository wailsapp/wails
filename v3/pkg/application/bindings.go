package application

import (
	"fmt"
	"github.com/samber/lo"
	"reflect"
	"runtime"
	"strings"
)

type CallOptions struct {
	PackageName string `json:"packageName"`
	StructName  string `json:"structName"`
	MethodName  string `json:"methodName"`
	Args        []any  `json:"args"`
}

// Parameter defines a Go method parameter
type Parameter struct {
	Name        string `json:"name,omitempty"`
	TypeName    string `json:"type"`
	reflectType reflect.Type
}

func newParameter(Name string, Type reflect.Type) *Parameter {
	return &Parameter{
		Name:        Name,
		TypeName:    Type.String(),
		reflectType: Type,
	}
}

// IsType returns true if the given
func (p *Parameter) IsType(typename string) bool {
	return p.TypeName == typename
}

// IsError returns true if the parameter type is an error
func (p *Parameter) IsError() bool {
	return p.IsType("error")
}

// BoundMethod defines all the data related to a Go method that is
// bound to the Wails application
type BoundMethod struct {
	Name        string        `json:"name"`
	Inputs      []*Parameter  `json:"inputs,omitempty"`
	Outputs     []*Parameter  `json:"outputs,omitempty"`
	Comments    string        `json:"comments,omitempty"`
	Method      reflect.Value `json:"-"`
	PackageName string
	StructName  string
	PackagePath string
}

type Bindings struct {
	boundMethods map[string]map[string]map[string]*BoundMethod
}

func NewBindings(bindings []any) (*Bindings, error) {
	b := &Bindings{
		boundMethods: make(map[string]map[string]map[string]*BoundMethod),
	}
	for _, binding := range bindings {
		err := b.Add(binding)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Add the given struct methods to the Bindings
func (b *Bindings) Add(structPtr interface{}) error {

	methods, err := b.getMethods(structPtr)
	if err != nil {
		return fmt.Errorf("cannot bind value to app: %s", err.Error())
	}

	for _, method := range methods {
		packageName := method.PackageName
		structName := method.StructName
		methodName := method.Name

		// Add it as a regular method
		if _, ok := b.boundMethods[packageName]; !ok {
			b.boundMethods[packageName] = make(map[string]map[string]*BoundMethod)
		}
		if _, ok := b.boundMethods[packageName][structName]; !ok {
			b.boundMethods[packageName][structName] = make(map[string]*BoundMethod)
		}
		b.boundMethods[packageName][structName][methodName] = method
		//b.db.AddMethod(packageName, structName, methodName, method)
	}
	return nil
}

func (b *Bindings) Get(options *CallOptions) *BoundMethod {
	_, ok := b.boundMethods[options.PackageName]
	if !ok {
		return nil
	}
	_, ok = b.boundMethods[options.PackageName][options.StructName]
	if !ok {
		return nil
	}
	method, ok := b.boundMethods[options.PackageName][options.StructName][options.MethodName]
	if !ok {
		return nil
	}
	return method
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

		return nil, fmt.Errorf("not a pointer to a struct")
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
		packageName, structName, _ := strings.Cut(baseName, ".")
		method := structValue.MethodByName(methodName)
		packagePath, _ := lo.Coalesce(structType.PkgPath(), "main")

		// Create new method
		boundMethod := &BoundMethod{
			Name:        methodName,
			PackageName: packageName,
			PackagePath: packagePath,
			StructName:  structName,
			Inputs:      nil,
			Outputs:     nil,
			Comments:    "",
			Method:      method,
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

	switch len(b.Outputs) {
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
