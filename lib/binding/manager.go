package binding

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
)

var typescriptDefinitionFilename = ""

// Manager handles method binding
type Manager struct {
	methods          map[string]*boundMethod
	functions        map[string]*boundFunction
	internalMethods  *internalMethods
	initMethods      []*boundMethod
	shutdownMethods  []*boundMethod
	log              *logger.CustomLogger
	renderer         interfaces.Renderer
	runtime          interfaces.Runtime // The runtime object to pass to bound structs
	objectsToBind    []interface{}
	bindPackageNames bool                // Package name should be considered when binding
	structList       map[string][]string // structList["mystruct"] = []string{"Method1", "Method2"}
}

// NewManager creates a new Manager struct
func NewManager() interfaces.BindingManager {

	result := &Manager{
		methods:         make(map[string]*boundMethod),
		functions:       make(map[string]*boundFunction),
		log:             logger.NewCustomLogger("Bind"),
		internalMethods: newInternalMethods(),
		structList:      make(map[string][]string),
	}
	return result
}

// BindPackageNames sets a flag to indicate package names should be considered when binding
func (b *Manager) BindPackageNames() {
	b.bindPackageNames = true
}

// Start the binding manager
func (b *Manager) Start(renderer interfaces.Renderer, runtime interfaces.Runtime) error {
	b.log.Info("Starting")
	b.renderer = renderer
	b.runtime = runtime
	err := b.initialise()
	if err != nil {
		b.log.Errorf("Binding error: %s", err.Error())
		return err
	}
	err = b.callWailsInitMethods()
	return err
}

func (b *Manager) initialise() error {

	var err error
	// var binding *boundMethod

	b.log.Info("Binding Go Functions/Methods")

	// Create bindings for objects
	for _, object := range b.objectsToBind {

		// Safeguard against nils
		if object == nil {
			return fmt.Errorf("attempted to bind nil object")
		}

		// Determine kind of object
		objectType := reflect.TypeOf(object)
		objectKind := objectType.Kind()

		switch objectKind {
		case reflect.Ptr:
			err = b.bindMethod(object)
		case reflect.Func:
			// spew.Dump(result.objectType.String())
			err = b.bindFunction(object)
		default:
			err = fmt.Errorf("cannot bind object of type '%s'", objectKind.String())
		}

		// Return error if set
		if err != nil {
			return err
		}
	}

	// If we wish to generate a typescript definition file...
	if typescriptDefinitionFilename != "" {
		err := b.generateTypescriptDefinitions()
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate typescript
func (b *Manager) generateTypescriptDefinitions() error {

	var output strings.Builder

	for structname, methodList := range b.structList {
		structname = strings.SplitN(structname, ".", 2)[1]
		output.WriteString(fmt.Sprintf("interface %s {\n", structname))
		for _, method := range methodList {
			output.WriteString(fmt.Sprintf("\t%s(...args : any[]):Promise<any>\n", method))
		}
		output.WriteString("}\n")
	}

	output.WriteString("\n")
	output.WriteString("interface Backend {\n")

	for structname := range b.structList {
		structname = strings.SplitN(structname, ".", 2)[1]
		output.WriteString(fmt.Sprintf("\t%[1]s: %[1]s\n", structname))
	}
	output.WriteString("}\n")

	globals := `
declare global {
	interface Window {
		backend: Backend;
	}
}
export {};`
	output.WriteString(globals)

	b.log.Info("Written Typescript file: " + typescriptDefinitionFilename)

	dir := filepath.Dir(typescriptDefinitionFilename)
	os.MkdirAll(dir, 0755)
	return ioutil.WriteFile(typescriptDefinitionFilename, []byte(output.String()), 0755)
}

// bind the given struct method
func (b *Manager) bindMethod(object interface{}) error {

	objectType := reflect.TypeOf(object)
	baseName := objectType.String()

	// Strip pointer if there
	if baseName[0] == '*' {
		baseName = baseName[1:]
	}

	b.log.Debugf("Processing struct: %s", baseName)

	// Calc actual name
	actualName := strings.TrimPrefix(baseName, "main.")
	if b.structList[actualName] == nil {
		b.structList[actualName] = []string{}
	}

	// Iterate over method definitions
	for i := 0; i < objectType.NumMethod(); i++ {

		// Get method definition
		methodDef := objectType.Method(i)
		methodName := methodDef.Name
		fullMethodName := baseName + "." + methodName
		method := reflect.ValueOf(object).MethodByName(methodName)

		b.structList[actualName] = append(b.structList[actualName], methodName)

		// Skip unexported methods
		if !unicode.IsUpper([]rune(methodName)[0]) {
			continue
		}

		// Create a new boundMethod
		newMethod, err := newBoundMethod(methodName, fullMethodName, method, objectType)
		if err != nil {
			return err
		}

		// Check if it's a wails init function
		if newMethod.isWailsInit {
			b.log.Debugf("Detected WailsInit function: %s", fullMethodName)
			b.initMethods = append(b.initMethods, newMethod)
		} else if newMethod.isWailsShutdown {
			b.log.Debugf("Detected WailsShutdown function: %s", fullMethodName)
			b.shutdownMethods = append(b.shutdownMethods, newMethod)
		} else {
			// Save boundMethod
			b.log.Infof("Bound Method: %s()", fullMethodName)
			b.methods[fullMethodName] = newMethod

			// Inform renderer of new binding
			b.renderer.NewBinding(fullMethodName)
		}
	}

	return nil
}

// bind the given function object
func (b *Manager) bindFunction(object interface{}) error {

	newFunction, err := newBoundFunction(object)
	if err != nil {
		return err
	}

	// Save method
	b.log.Infof("Bound Function: %s()", newFunction.fullName)
	b.functions[newFunction.fullName] = newFunction

	// Register with Renderer
	b.renderer.NewBinding(newFunction.fullName)

	return nil
}

// Bind saves the given object to be bound at start time
func (b *Manager) Bind(object interface{}) {
	// Store binding
	b.objectsToBind = append(b.objectsToBind, object)
}

func (b *Manager) processInternalCall(callData *messages.CallData) (interface{}, error) {
	// Strip prefix
	return b.internalMethods.processCall(callData)
}

func (b *Manager) processFunctionCall(callData *messages.CallData) (interface{}, error) {
	// Return values
	var result []reflect.Value
	var err error

	function := b.functions[callData.BindingName]
	if function == nil {
		return nil, fmt.Errorf("Invalid function name '%s'", callData.BindingName)
	}
	result, err = function.call(callData.Data)
	if err != nil {
		return nil, err
	}

	// Do we have an error return type?
	if function.hasErrorReturnType {
		// We do - last result is an error type
		// Check if the last result was nil
		b.log.Debugf("# of return types: %d", len(function.returnTypes))
		b.log.Debugf("# of results: %d", len(result))
		errorResult := result[len(function.returnTypes)-1]
		if !errorResult.IsNil() {
			// It wasn't - we have an error
			return nil, errorResult.Interface().(error)
		}
	}
	// fmt.Printf("result = '%+v'\n", result)
	if len(result) > 0 {
		return result[0].Interface(), nil
	}
	return nil, nil
}

func (b *Manager) processMethodCall(callData *messages.CallData) (interface{}, error) {
	// Return values
	var result []reflect.Value
	var err error

	// do we have this method?
	method := b.methods[callData.BindingName]
	if method == nil {
		return nil, fmt.Errorf("Invalid method name '%s'", callData.BindingName)
	}

	result, err = method.call(callData.Data)
	if err != nil {
		return nil, err
	}

	// Do we have an error return type?
	if method.hasErrorReturnType {
		// We do - last result is an error type
		// Check if the last result was nil
		b.log.Debugf("# of return types: %d", len(method.returnTypes))
		b.log.Debugf("# of results: %d", len(result))
		errorResult := result[len(method.returnTypes)-1]
		if !errorResult.IsNil() {
			// It wasn't - we have an error
			return nil, errorResult.Interface().(error)
		}
	}
	if result != nil {
		return result[0].Interface(), nil
	}
	return nil, nil
}

// ProcessCall processes the given call request
func (b *Manager) ProcessCall(callData *messages.CallData) (result interface{}, err error) {
	b.log.Debugf("Wanting to call %s", callData.BindingName)

	// Determine if this is function call or method call by the number of
	// dots in the binding name
	dotCount := 0
	for _, character := range callData.BindingName {
		if character == '.' {
			dotCount++
		}
	}

	// We need to catch reflect related panics and return
	// a decent error message
	// TODO: DEBUG THIS!

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r.(string))
		}
	}()

	switch dotCount {
	case 1:
		result, err = b.processFunctionCall(callData)
	case 2:
		result, err = b.processMethodCall(callData)
	case 3:
		result, err = b.processInternalCall(callData)
	default:
		result = nil
		err = fmt.Errorf("Invalid binding name '%s'", callData.BindingName)
	}
	return
}

// callWailsInitMethods calls all of the WailsInit methods that were
// registered with the runtime object
func (b *Manager) callWailsInitMethods() error {
	// Create reflect value for runtime object
	runtimeValue := reflect.ValueOf(b.runtime)
	params := []reflect.Value{runtimeValue}

	// Iterate initMethods
	for _, initMethod := range b.initMethods {
		// Call
		result := initMethod.method.Call(params)
		// Check errors
		err := result[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	return nil
}

// Shutdown the binding manager
func (b *Manager) Shutdown() {
	b.log.Debug("Shutdown called")
	for _, method := range b.shutdownMethods {
		b.log.Debugf("Calling Shutdown for method: %s", method.fullName)
		method.call("[]")
	}
	b.log.Debug("Shutdown complete")
}
