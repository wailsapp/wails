package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

type PluginCallOptions struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

var reservedPluginMethods = []string{
	"Name",
	"Init",
	"Shutdown",
	"Exported",
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
	ID          uint32        `json:"id"`
	Name        string        `json:"name"`
	Inputs      []*Parameter  `json:"inputs,omitempty"`
	Outputs     []*Parameter  `json:"outputs,omitempty"`
	Comments    string        `json:"comments,omitempty"`
	Method      reflect.Value `json:"-"`
	TypeName    string
	PackagePath string

	needsContext bool
}

type Bindings struct {
	boundMethods  map[string]*BoundMethod
	boundByID     map[uint32]*BoundMethod
	methodAliases map[uint32]uint32
}

func NewBindings(instances []Service, aliases map[uint32]uint32) (*Bindings, error) {
	app := Get()
	b := &Bindings{
		boundMethods:  make(map[string]*BoundMethod),
		boundByID:     make(map[uint32]*BoundMethod),
		methodAliases: aliases,
	}
	for _, binding := range instances {
		handler, ok := binding.Instance().(http.Handler)
		if ok && binding.options.Route != "" {
			app.assets.AttachServiceHandler(binding.options.Route, handler)
		}
		err := b.Add(binding.Instance())
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Add the given named type pointer methods to the Bindings
func (b *Bindings) Add(namedPtr interface{}) error {
	methods, err := b.getMethods(namedPtr)
	if err != nil {
		return fmt.Errorf("cannot bind value to app: %s", err.Error())
	}

	for _, method := range methods {
		// Add it as a regular method
		b.boundMethods[method.String()] = method
		b.boundByID[method.ID] = method
	}
	return nil
}

// Get returns the bound method with the given name
func (b *Bindings) Get(options *CallOptions) *BoundMethod {
	method, ok := b.boundMethods[options.MethodName]
	if !ok {
		return nil
	}
	return method
}

// GetByID returns the bound method with the given ID
func (b *Bindings) GetByID(id uint32) *BoundMethod {
	// Check method aliases
	if b.methodAliases != nil {
		if alias, ok := b.methodAliases[id]; ok {
			id = alias
		}
	}
	result := b.boundByID[id]
	return result
}

// GenerateID generates a unique ID for a binding
func (b *Bindings) GenerateID(name string) (uint32, error) {
	id, err := hash.Fnv(name)
	if err != nil {
		return 0, err
	}
	// Check if we already have it
	boundMethod, ok := b.boundByID[id]
	if ok {
		return 0, fmt.Errorf("oh wow, we're sorry about this! Amazingly, a hash collision was detected for method '%s' (it generates the same hash as '%s'). To continue, please rename it. Sorry :(", name, boundMethod.String())
	}
	return id, nil
}

func (b *BoundMethod) String() string {
	return fmt.Sprintf("%s.%s.%s", b.PackagePath, b.TypeName, b.Name)
}

func (b *Bindings) getMethods(value interface{}) ([]*BoundMethod, error) {
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

	ctxType := reflect.TypeFor[context.Context]()

	// Process Methods
	for i := 0; i < ptrType.NumMethod(); i++ {
		methodDef := ptrType.Method(i)
		methodName := methodDef.Name
		method := namedValue.MethodByName(methodName)

		if b.internalMethod(methodDef) {
			continue
		}

		// Create new method
		boundMethod := &BoundMethod{
			Name:        methodName,
			PackagePath: packagePath,
			TypeName:    typeName,
			Inputs:      nil,
			Outputs:     nil,
			Comments:    "",
			Method:      method,
		}
		var err error
		boundMethod.ID, err = b.GenerateID(boundMethod.String())
		if err != nil {
			return nil, err
		}

		args := []any{"name", boundMethod, "id", boundMethod.ID}
		if b.methodAliases != nil {
			alias, found := lo.FindKey(b.methodAliases, boundMethod.ID)
			if found {
				args = append(args, "alias", alias)
			}
		}
		globalApplication.debug("Adding method:", args...)

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

func (b *Bindings) internalMethod(def reflect.Method) bool {
	// Get the receiver type
	receiverType := def.Type.In(0)

	// Create a new instance of the receiver type
	instance := reflect.New(receiverType.Elem()).Interface()

	// Check if the instance implements any of our service interfaces
	// and if the method matches the interface method
	switch def.Name {
	case "ServiceName":
		if _, ok := instance.(ServiceName); ok {
			return true
		}
	case "ServiceStartup":
		if _, ok := instance.(ServiceStartup); ok {
			return true
		}
	case "ServiceShutdown":
		if _, ok := instance.(ServiceShutdown); ok {
			return true
		}
	}

	return false
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
