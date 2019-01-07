package wails

import (
	"fmt"
	"reflect"
	"unicode"
)

/**

binding:
 Name() // Full name (package+name)
 Call(params)

**/

type bindingManager struct {
	methods          map[string]*boundMethod
	functions        map[string]*boundFunction
	initMethods      []*boundMethod
	log              *CustomLogger
	renderer         Renderer
	runtime          *Runtime // The runtime object to pass to bound structs
	objectsToBind    []interface{}
	bindPackageNames bool // Package name should be considered when binding
}

func newBindingManager() *bindingManager {
	result := &bindingManager{
		methods:   make(map[string]*boundMethod),
		functions: make(map[string]*boundFunction),
		log:       newCustomLogger("Bind"),
	}
	return result
}

// Sets flag to indicate package names should be considered when binding
func (b *bindingManager) BindPackageNames() {
	b.bindPackageNames = true
}

func (b *bindingManager) start(renderer Renderer, runtime *Runtime) error {
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

func (b *bindingManager) initialise() error {

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
	return nil
}

// bind the given struct method
func (b *bindingManager) bindMethod(object interface{}) error {

	objectType := reflect.TypeOf(object)
	baseName := objectType.String()

	// Strip pointer if there
	if baseName[0] == '*' {
		baseName = baseName[1:]
	}

	b.log.Debugf("Processing struct: %s", baseName)

	// Iterate over method definitions
	for i := 0; i < objectType.NumMethod(); i++ {

		// Get method definition
		methodDef := objectType.Method(i)
		methodName := methodDef.Name
		fullMethodName := baseName + "." + methodName
		method := reflect.ValueOf(object).MethodByName(methodName)

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
func (b *bindingManager) bindFunction(object interface{}) error {

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

// Save the given object to be bound at start time
func (b *bindingManager) bind(object interface{}) {
	// Store binding
	b.objectsToBind = append(b.objectsToBind, object)
}

// process an incoming call request
func (b *bindingManager) processCall(callData *callData) (interface{}, error) {
	b.log.Debugf("Wanting to call %s", callData.BindingName)

	// Determine if this is function call or method call by the number of
	// dots in the binding name
	dotCount := 0
	for _, character := range callData.BindingName {
		if character == '.' {
			dotCount++
		}
	}

	// Return values
	var result []reflect.Value
	var err error

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
		return result[0].Interface(), nil
	case 2:
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

	default:
		return nil, fmt.Errorf("Invalid binding name '%s'", callData.BindingName)
	}
}

// callWailsInitMethods calls all of the WailsInit methods that were
// registered with the runtime object
func (b *bindingManager) callWailsInitMethods() error {
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
