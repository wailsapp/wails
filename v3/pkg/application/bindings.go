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
	MethodID    uint32            `json:"methodID"`
	PackageName string            `json:"packageName"`
	StructName  string            `json:"structName"`
	MethodName  string            `json:"methodName"`
	Args        []json.RawMessage `json:"args"`
}

func (c *CallOptions) Name() string {
	return fmt.Sprintf("%s.%s.%s", c.PackageName, c.StructName, c.MethodName)
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
	PackageName string
	StructName  string
	PackagePath string

	needsContext bool
}

type Bindings struct {
	boundMethods  map[string]map[string]map[string]*BoundMethod
	boundByID     map[uint32]*BoundMethod
	methodAliases map[uint32]uint32
}

func NewBindings(structs []any, aliases map[uint32]uint32) (*Bindings, error) {
	b := &Bindings{
		boundMethods:  make(map[string]map[string]map[string]*BoundMethod),
		boundByID:     make(map[uint32]*BoundMethod),
		methodAliases: aliases,
	}
	for _, binding := range structs {
		err := b.Add(binding)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Add the given struct methods to the Bindings
func (b *Bindings) Add(structPtr interface{}) error {

	methods, err := b.getMethods(structPtr, false)
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
		b.boundByID[method.ID] = method
	}
	return nil
}

func (b *Bindings) AddPlugins(plugins map[string]Plugin) error {
	for pluginID, plugin := range plugins {
		methods, err := b.getMethods(plugin, true)
		if err != nil {
			return fmt.Errorf("cannot add plugin '%s' to app: %s", pluginID, err.Error())
		}

		exportedMethods := plugin.CallableByJS()

		for _, method := range methods {
			// Do not expose reserved methods
			if lo.Contains(reservedPluginMethods, method.Name) {
				continue
			}
			// Do not expose methods that are not in the exported list
			if !lo.Contains(exportedMethods, method.Name) {
				continue
			}
			packageName := "wails-plugins"
			structName := pluginID
			methodName := method.Name

			// Add it as a regular method
			if _, ok := b.boundMethods[packageName]; !ok {
				b.boundMethods[packageName] = make(map[string]map[string]*BoundMethod)
			}
			if _, ok := b.boundMethods[packageName][structName]; !ok {
				b.boundMethods[packageName][structName] = make(map[string]*BoundMethod)
			}
			b.boundMethods[packageName][structName][methodName] = method
			b.boundByID[method.ID] = method
			globalApplication.debug("Added plugin method: "+structName+"."+methodName, "id", method.ID)
		}
	}
	return nil
}

// Get returns the bound method with the given name
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
	return fmt.Sprintf("%s.%s.%s", b.PackageName, b.StructName, b.Name)
}

func (b *Bindings) getMethods(value interface{}, isPlugin bool) ([]*BoundMethod, error) {

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
			return nil, fmt.Errorf("%s is a function, not a pointer to a struct. Wails v2 has deprecated the binding of functions. Please wrap your functions up in a struct and bind a pointer to that struct", name)
		}

		return nil, fmt.Errorf("not a pointer to a struct")
	}

	// Process Struct
	structType := reflect.TypeOf(value)
	structValue := reflect.ValueOf(value)
	structTypeString := structType.String()
	baseName := structTypeString[1:]

	ctxType := reflect.TypeFor[context.Context]()

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
		var err error
		boundMethod.ID, err = hash.Fnv(boundMethod.String())
		if err != nil {
			return nil, err
		}

		if !isPlugin {
			args := []any{"name", boundMethod, "id", boundMethod.ID}
			if b.methodAliases != nil {
				alias, found := lo.FindKey(b.methodAliases, boundMethod.ID)
				if found {
					args = append(args, "alias", alias)
				}
			}
			globalApplication.debug("Adding method:", args...)
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

var errorType = reflect.TypeFor[error]()

// Call will attempt to call this bound method with the given args
func (b *BoundMethod) Call(ctx context.Context, args []json.RawMessage) (returnValue interface{}, err error) {
	// Use a defer statement to capture panics
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				if strings.HasPrefix(str, "reflect: Call using") {
					// Remove prefix
					str = strings.Replace(str, "reflect: Call using ", "", 1)
					// Split on "as"
					parts := strings.Split(str, " as type ")
					if len(parts) == 2 {
						err = fmt.Errorf("invalid argument type: got '%s', expected '%s'", parts[0], parts[1])
						return
					}
				}
			}
			err = fmt.Errorf("%v", r)
		}
	}()

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
