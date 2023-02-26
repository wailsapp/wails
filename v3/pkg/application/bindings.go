package application

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

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
		packageName, structName, _ := strings.Cut(baseName, ".")
		fullMethodName := baseName + "." + methodName
		method := structValue.MethodByName(methodName)

		// Create new method
		boundMethod := &BoundMethod{
			Name:        fullMethodName,
			PackageName: packageName,
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
			//
			//thisInput := input
			//
			//if thisInput.Kind() == reflect.Slice {
			//	thisInput = thisInput.Elem()
			//}
			//
			//// Process struct pointer params
			//if thisInput.Kind() == reflect.Ptr {
			//	if thisInput.Elem().Kind() == reflect.Struct {
			//		typ := thisInput.Elem()
			//		a := reflect.New(typ)
			//		s := reflect.Indirect(a).Interface()
			//		name := typ.Name()
			//		packageName := getPackageName(thisInput.String())
			//		b.AddStructToGenerateTS(packageName, name, s)
			//	}
			//}
			//
			//// Process struct params
			//if thisInput.Kind() == reflect.Struct {
			//	a := reflect.New(thisInput)
			//	s := reflect.Indirect(a).Interface()
			//	name := thisInput.Name()
			//	packageName := getPackageName(thisInput.String())
			//	b.AddStructToGenerateTS(packageName, name, s)
			//}

			inputs = append(inputs, thisParam)
		}

		boundMethod.Inputs = inputs

		outputParamCount := methodType.NumOut()
		var outputs []*Parameter
		for outputIndex := 0; outputIndex < outputParamCount; outputIndex++ {
			output := methodType.Out(outputIndex)
			thisParam := newParameter("", output)
			//
			//thisOutput := output
			//
			//if thisOutput.Kind() == reflect.Slice {
			//	thisOutput = thisOutput.Elem()
			//}
			//
			//// Process struct pointer params
			//if thisOutput.Kind() == reflect.Ptr {
			//	if thisOutput.Elem().Kind() == reflect.Struct {
			//		typ := thisOutput.Elem()
			//		a := reflect.New(typ)
			//		s := reflect.Indirect(a).Interface()
			//		name := typ.Name()
			//		packageName := getPackageName(thisOutput.String())
			//	}
			//}
			//
			//// Process struct params
			//if thisOutput.Kind() == reflect.Struct {
			//	a := reflect.New(thisOutput)
			//	s := reflect.Indirect(a).Interface()
			//	name := thisOutput.Name()
			//	packageName := getPackageName(thisOutput.String())
			//	b.AddStructToGenerateTS(packageName, name, s)
			//}

			outputs = append(outputs, thisParam)
		}
		boundMethod.Outputs = outputs

		// Save method in result
		result = append(result, boundMethod)

	}
	return result, nil
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

func getPackageName(in string) string {
	result := strings.Split(in, ".")[0]
	result = strings.ReplaceAll(result, "[]", "")
	result = strings.ReplaceAll(result, "*", "")
	return result
}
