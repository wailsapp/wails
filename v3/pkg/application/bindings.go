package application

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/wailsapp/wails/v3/internal/hash"
	"github.com/wailsapp/wails/v3/internal/sliceutil"
)

// init forces goccy/go-json to initialize its type address cache at program startup.
// This prevents a Windows-specific index out-of-bounds panic that can occur when the decoder is first invoked later (see https://github.com/goccy/go-json/issues/474).
func init() {
	// Force goccy/go-json to initialize its type address cache early.
	// On Windows, if the decoder is first invoked later (e.g., during tests),
	// the type address calculation can fail with an index out of bounds panic.
	// See: https://github.com/goccy/go-json/issues/474
	var si []int
	_ = json.Unmarshal([]byte(`[]`), &si)
}

// CallOptions defines the options for a method call.
// Field order is optimized to minimize struct padding.
type CallOptions struct {
	MethodName string            `json:"methodName"`
	Args       []json.RawMessage `json:"args"`
	MethodID   uint32            `json:"methodID"`
}

type ErrorKind string

const (
	ReferenceError ErrorKind = "ReferenceError"
	TypeError      ErrorKind = "TypeError"
	RuntimeError   ErrorKind = "RuntimeError"
)

// CallError represents an error that occurred during a method call.
// Field order is optimized to minimize struct padding.
type CallError struct {
	Message string    `json:"message"`
	Cause   any       `json:"cause,omitempty"`
	Kind    ErrorKind `json:"kind"`
}

func (e *CallError) Error() string {
	return e.Message
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
// bound to the Wails application.
// Field order is optimized to minimize struct padding (136 bytes vs 144 bytes).
type BoundMethod struct {
	Method       reflect.Value `json:"-"`
	Name         string        `json:"name"`
	FQN          string        `json:"-"`
	Comments     string        `json:"comments,omitempty"`
	Inputs       []*Parameter  `json:"inputs,omitempty"`
	Outputs      []*Parameter  `json:"outputs,omitempty"`
	marshalError func(error) []byte
	ID           uint32 `json:"id"`
	needsContext bool
	isVariadic   bool // cached at registration to avoid reflect call per invocation
}

type Bindings struct {
	marshalError  func(error) []byte
	boundMethods  map[string]*BoundMethod
	boundByID     map[uint32]*BoundMethod
	methodAliases map[uint32]uint32
}

func NewBindings(marshalError func(error) []byte, aliases map[uint32]uint32) *Bindings {
	return &Bindings{
		marshalError:  wrapErrorMarshaler(marshalError, defaultMarshalError),
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

	marshalError := wrapErrorMarshaler(service.options.MarshalError, defaultMarshalError)

	// Validate and log methods.
	for _, method := range methods {
		if _, ok := b.boundMethods[method.FQN]; ok {
			return fmt.Errorf("bound method '%s' is already registered. Please note that you can register at most one service of each type; additional instances must be wrapped in dedicated structs", method.FQN)
		}
		if boundMethod, ok := b.boundByID[method.ID]; ok {
			return fmt.Errorf("oh wow, we're sorry about this! Amazingly, a hash collision was detected for method '%s' (it generates the same hash as '%s'). To use this method, please rename it. Sorry :(", method.FQN, boundMethod.FQN)
		}

		// Log
		attrs := []any{"fqn", method.FQN, "id", method.ID}
		if alias, ok := sliceutil.FindMapKey(b.methodAliases, method.ID); ok {
			attrs = append(attrs, "alias", alias)
		}
		globalApplication.debug("Registering bound method:", attrs...)
	}

	for _, method := range methods {
		// Store composite error marshaler
		method.marshalError = marshalError

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

// getMethods returns the list of BoundMethod descriptors for the methods of the named pointer type provided by value.
// 
// It returns an error if value is not a pointer to a named type, if a function value is supplied (binding functions is deprecated), or if a generic type is supplied.
// The returned BoundMethod slice includes only exported methods that are not listed in internalServiceMethods. Each BoundMethod has its FQN, ID (computed from the FQN), Method reflect.Value, Inputs and Outputs populated, isVariadic cached from the method signature, and needsContext set when the first parameter is context.Context.
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

		// Iterate inputs
		methodType := method.Type()

		// Create new method with cached flags
		boundMethod := &BoundMethod{
			ID:         hash.Fnv(fqn),
			FQN:        fqn,
			Name:       methodName,
			Inputs:     nil,
			Outputs:    nil,
			Comments:   "",
			Method:     method,
			isVariadic: methodType.IsVariadic(), // cache to avoid reflect call per invocation
		}
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

// Call will attempt to call this bound method with the given args.
// If the call succeeds, result will be either a non-error return value (if there is only one)
// or a slice of non-error return values (if there are more than one).
//
// If the arguments are mistyped or the call returns one or more non-nil error values,
// result is nil and err is an instance of *[CallError].
func (b *BoundMethod) Call(ctx context.Context, args []json.RawMessage) (result any, err error) {
	// Use a defer statement to capture panics
	defer handlePanic(handlePanicOptions{skipEnd: 5})
	argCount := len(args)
	if b.needsContext {
		argCount++
	}

	if argCount != len(b.Inputs) {
		err = &CallError{
			Message: fmt.Sprintf("%s expects %d arguments, got %d", b.FQN, len(b.Inputs), argCount),
			Kind:    TypeError,
		}
		return
	}

	// Use stack-allocated buffer for common case (<=8 args), heap for larger
	var argBuffer [8]reflect.Value
	var callArgs []reflect.Value
	if argCount <= len(argBuffer) {
		callArgs = argBuffer[:argCount]
	} else {
		callArgs = make([]reflect.Value, argCount)
	}
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
			err = &CallError{
				Message: fmt.Sprintf("could not parse argument #%d: %s", index, err),
				Cause:   json.RawMessage(b.marshalError(err)),
				Kind:    TypeError,
			}
			return
		}
		callArgs[base+index] = value.Elem()
	}

	// Do the call using cached isVariadic flag
	var callResults []reflect.Value
	if b.isVariadic {
		callResults = b.Method.CallSlice(callArgs)
	} else {
		callResults = b.Method.Call(callArgs)
	}

	// Process results - optimized for common case of 0-2 return values
	// to avoid slice allocation
	var firstResult any
	var hasFirstResult bool
	var nonErrorOutputs []any // only allocated if >1 non-error results
	var errorOutputs []error

	for _, field := range callResults {
		if field.Type() == errorType {
			if field.IsNil() {
				continue
			}
			if errorOutputs == nil {
				errorOutputs = make([]error, 0, len(callResults))
			}
			errorOutputs = append(errorOutputs, field.Interface().(error))
		} else if errorOutputs == nil {
			// Only collect non-error outputs if no errors yet
			val := field.Interface()
			if !hasFirstResult {
				firstResult = val
				hasFirstResult = true
			} else if nonErrorOutputs == nil {
				// Second result - need to allocate slice
				nonErrorOutputs = make([]any, 0, len(callResults))
				nonErrorOutputs = append(nonErrorOutputs, firstResult, val)
			} else {
				nonErrorOutputs = append(nonErrorOutputs, val)
			}
		}
	}

	if len(errorOutputs) > 0 {
		info := make([]json.RawMessage, len(errorOutputs))
		for i, err := range errorOutputs {
			info[i] = b.marshalError(err)
		}

		cerr := &CallError{
			Message: errors.Join(errorOutputs...).Error(),
			Cause:   info,
			Kind:    RuntimeError,
		}
		if len(info) == 1 {
			cerr.Cause = info[0]
		}

		err = cerr
	} else if nonErrorOutputs != nil {
		result = nonErrorOutputs
	} else if hasFirstResult {
		result = firstResult
	}

	return
}

// wrapErrorMarshaler returns an error marshaling functions
// that calls the primary marshaler first,
// then falls back to the secondary one.
func wrapErrorMarshaler(primary func(error) []byte, secondary func(error) []byte) func(error) []byte {
	if primary == nil {
		return secondary
	}

	return func(err error) []byte {
		result := primary(err)
		if result == nil {
			result = secondary(err)
		}

		return result
	}
}

// defaultMarshalError implements the default error marshaling mechanism.
func defaultMarshalError(err error) []byte {
	result, jsonErr := json.Marshal(&err)
	if jsonErr != nil {
		return nil
	}
	return result
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