package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/hash"

	"github.com/samber/lo"
)

type CallOptions struct {
	MethodID   uint32            `json:"methodID"`
	MethodName string            `json:"methodName"`
	Args       []json.RawMessage `json:"args"`
}

// Parameter defines a Go method parameter
type Parameter struct {
	Name        string `json:"name,omitempty"`
	TypeName    string `json:"type"`
	ReflectType reflect.Type
}

func newParameter(Name string, Type reflect.Type) *Parameter {
	return &Parameter{
		Name:        Name,
		TypeName:    Type.String(),
		ReflectType: Type,
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
	ID       uint32        `json:"id"`
	Name     string        `json:"name"`
	Inputs   []*Parameter  `json:"inputs,omitempty"`
	Outputs  []*Parameter  `json:"outputs,omitempty"`
	Comments string        `json:"comments,omitempty"`
	Method   reflect.Value `json:"-"`
	FQN      string

	needsContext bool
}

type Bindings struct {
	boundMethods  map[string]*BoundMethod
	boundByID     map[uint32]*BoundMethod
	methodAliases map[uint32]uint32
}

func NewBindings(aliases map[uint32]uint32) *Bindings {
	return &Bindings{
		boundMethods:  make(map[string]*BoundMethod),
		boundByID:     make(map[uint32]*BoundMethod),
		methodAliases: aliases,
	}
}

// Add adds the given service to the bindings.
func (b *Bindings) Add(service Service) error {
	methods, err := getMethods(service.Instance())
	if err != nil {
		return err
	}

	// Validate and log methods.
	for _, method := range methods {
		if _, ok := b.boundMethods[method.FQN]; ok {
			return fmt.Errorf("bound method '%s' is already registered. Please note that you can register at most one service of each type; additional instances must be wrapped in dedicated structs", method.FQN)
		}
		if boundMethod, ok := b.boundByID[method.ID]; ok {
			return fmt.Errorf("oh wow, we're sorry about this! Amazingly, a hash collision was detected for method '%s' (it generates the same hash as '%s'). To use this method, please rename it. Sorry :(", method.FQN, boundMethod.FQN)
		}

		// Log
		attrs := []any{"name", method.Name, "id", method.ID}
		if alias, ok := lo.FindKey(b.methodAliases, method.ID); ok {
			attrs = append(attrs, "alias", alias)
		}
		globalApplication.debug("Registering bound method:", attrs...)
	}

	for i, method := range methods {
		// Register method
		b.boundMethods[method.FQN] = method
		b.boundByID[method.ID] = method
	}

	return nil
}

// Get returns the bound method with the given name
func (b *Bindings) Get(options *CallOptions) *BoundMethod {
	return b.boundMethods[options.MethodName]
}

// GetByID returns the bound method with the given ID
func (b *Bindings) GetByID(id uint32) *BoundMethod {
	// Check method aliases
	if b.methodAliases != nil {
		if alias, ok := b.methodAliases[id]; ok {
			id = alias
		}
	}

	return b.boundByID[id]
}

// internalServiceMethod is a set of methods
// that are handled specially by the binding engine
// and must not be exposed to the frontend.
//
// For simplicity we exclude these by name
// without checking their signatures,
// and so does the binding generator.
var internalServiceMethods = map[string]bool{
	"ServiceName":     true,
	"ServiceStartup":  true,
	"ServiceShutdown": true,
	"ServeHTTP":       true,
}

var ctxType = reflect.TypeFor[context.Context]()

func getMethods(value any) ([]*BoundMethod, error) {
	// Create result placeholder
	var result []*BoundMethod

	// Check type
	if !isNamed(value) {
		if isFunction(value) {
			name := runtime.FuncForPC(reflect.ValueOf(value).Pointer()).Name()
			return nil, fmt.Errorf("%s is a function, not a pointer to named type. Wails v2 has deprecated the binding of functions. Please define your functions as methods on a struct and bind a pointer to that struct", name)
		}

		return nil, fmt.Errorf("%s is not a pointer to named type", reflect.ValueOf(value).Type().String())
	} else if !isPtr(value) {
		return nil, fmt.Errorf("%s is a named type, not a pointer to named type", reflect.ValueOf(value).Type().String())
	}

	// Process Named Type
	namedValue := reflect.ValueOf(value)
	ptrType := namedValue.Type()
	namedType := ptrType.Elem()
	typeName := namedType.Name()
	packagePath := namedType.PkgPath()

	if strings.Contains(namedType.String(), "[") {
		return nil, fmt.Errorf("%s.%s is a generic type. Generic bound types are not supported", packagePath, namedType.String())
	}

	// Process Methods
	for i := range ptrType.NumMethod() {
		methodName := ptrType.Method(i).Name
		method := namedValue.Method(i)

		if internalServiceMethods[methodName] {
			continue
		}

		fqn := fmt.Sprintf("%s.%s.%s", packagePath, typeName, methodName)

		// Create new method
		boundMethod := &BoundMethod{
			ID:       hash.Fnv(fqn),
			FQN:      fqn,
			Name:     methodName,
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
			if inputIndex == 0 && input.AssignableTo(ctxType) {
				boundMethod.needsContext = true
			}
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

func (b *BoundMethod) String() string {
	return b.FQN
}

var errorType = reflect.TypeFor[error]()

// Call will attempt to call this bound method with the given args
func (b *BoundMethod) Call(ctx context.Context, args []json.RawMessage) (returnValue interface{}, err error) {
	// Use a defer statement to capture panics
	defer handlePanic(handlePanicOptions{skipEnd: 5})
	argCount := len(args)
	if b.needsContext {
		argCount++
	}

	if argCount != len(b.Inputs) {
		err = fmt.Errorf("%s expects %d arguments, received %d", b.Name, len(b.Inputs), argCount)
		return
	}

	// Convert inputs to values of appropriate type
	callArgs := make([]reflect.Value, argCount)
	base := 0

	if b.needsContext {
		callArgs[0] = reflect.ValueOf(ctx)
		base++
	}

	// Iterate over given arguments
	for index, arg := range args {
		value := reflect.New(b.Inputs[base+index].ReflectType)
		err = json.Unmarshal(arg, value.Interface())
		if err != nil {
			err = fmt.Errorf("could not parse argument #%d: %w", index, err)
			return
		}
		callArgs[base+index] = value.Elem()
	}

	// Do the call
	var callResults []reflect.Value
	if b.Method.Type().IsVariadic() {
		callResults = b.Method.CallSlice(callArgs)
	} else {
		callResults = b.Method.Call(callArgs)
	}

	var nonErrorOutputs = make([]any, 0, len(callResults))
	var errorOutputs []error

	for _, result := range callResults {
		if result.Type() == errorType {
			if result.IsNil() {
				continue
			}
			if errorOutputs == nil {
				errorOutputs = make([]error, 0, len(callResults)-len(nonErrorOutputs))
				nonErrorOutputs = nil
			}
			errorOutputs = append(errorOutputs, result.Interface().(error))
		} else if nonErrorOutputs != nil {
			nonErrorOutputs = append(nonErrorOutputs, result.Interface())
		}
	}

	if errorOutputs != nil {
		err = errors.Join(errorOutputs...)
	} else if len(nonErrorOutputs) == 1 {
		returnValue = nonErrorOutputs[0]
	} else if len(nonErrorOutputs) > 1 {
		returnValue = nonErrorOutputs
	}

	return
}

// isPtr returns true if the value given is a pointer.
func isPtr(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}

// isFunction returns true if the given value is a function
func isFunction(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Func
}

// isNamed returns true if the given value is of named type
// or pointer to named type.
func isNamed(value interface{}) bool {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv.Type().Name() != ""
}
