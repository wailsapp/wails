//go:build wails_obfuscated

package application

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Enumerating with reflect.Type.Method here would make the linker keep every exported method of
// every reachable type (cmd/link reflectSeen); only the registry keeps dead-method elimination on.
func boundMethodsOf(namedValue reflect.Value, ptrType reflect.Type, typeName, packagePath string) ([]*BoundMethod, error) {
	var result []*BoundMethod

	registeredBindingMethodIDs.Range(func(_, value any) bool {
		entry := value.(registeredBindingMethod)
		fnType := entry.fn.Type()
		if fnType.NumIn() == 0 || fnType.In(0) != ptrType {
			return true
		}

		result = append(result, newRegisteredBoundMethod(namedValue, entry, typeName, packagePath))
		return true
	})

	if len(result) == 0 {
		return nil, fmt.Errorf("%s.%s has no registered binding IDs: regenerate the bindings with `wails3 generate bindings -obfuscated`", packagePath, typeName)
	}

	return result, nil
}

func newRegisteredBoundMethod(receiver reflect.Value, entry registeredBindingMethod, typeName, packagePath string) *BoundMethod {
	fn := entry.fn
	fnType := fn.Type()
	variadic := fnType.IsVariadic()

	inTypes := make([]reflect.Type, 0, fnType.NumIn()-1)
	for i := 1; i < fnType.NumIn(); i++ {
		inTypes = append(inTypes, fnType.In(i))
	}
	outTypes := make([]reflect.Type, 0, fnType.NumOut())
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}

	bound := reflect.MakeFunc(reflect.FuncOf(inTypes, outTypes, variadic), func(args []reflect.Value) []reflect.Value {
		callArgs := append([]reflect.Value{receiver}, args...)
		if variadic {
			return fn.CallSlice(callArgs)
		}
		return fn.Call(callArgs)
	})

	methodName := registeredMethodName(fn, entry.id)
	boundMethod := &BoundMethod{
		ID:         entry.id,
		FQN:        fmt.Sprintf("%s.%s.%s", packagePath, typeName, methodName),
		Name:       methodName,
		Method:     bound,
		isVariadic: variadic,
	}

	for i, in := range inTypes {
		if i == 0 && in.AssignableTo(ctxType) {
			boundMethod.needsContext = true
		}
		boundMethod.Inputs = append(boundMethod.Inputs, newParameter("", in))
	}
	for _, out := range outTypes {
		boundMethod.Outputs = append(boundMethod.Outputs, newParameter("", out))
	}

	return boundMethod
}

func registeredMethodName(fn reflect.Value, id uint32) string {
	if f := runtime.FuncForPC(fn.Pointer()); f != nil {
		name := f.Name()
		if i := strings.LastIndex(name, "."); i >= 0 && i < len(name)-1 {
			return name[i+1:]
		}
	}
	return fmt.Sprintf("method%d", id)
}
