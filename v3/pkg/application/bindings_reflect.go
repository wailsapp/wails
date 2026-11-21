//go:build !wails_obfuscated

package application

import (
	"fmt"
	"reflect"

	"github.com/wailsapp/wails/v3/internal/hash"
)

func getRegisteredBindingMethodID(method reflect.Method) (uint32, bool) {
	entry, ok := registeredBindingMethodIDs.Load(method.Func.Pointer())
	if !ok {
		return 0, false
	}
	return entry.(registeredBindingMethod).id, true
}

func boundMethodsOf(namedValue reflect.Value, ptrType reflect.Type, typeName, packagePath string) ([]*BoundMethod, error) {
	var result []*BoundMethod

	for i := range ptrType.NumMethod() {
		methodDef := ptrType.Method(i)
		methodName := methodDef.Name
		method := namedValue.Method(i)

		if internalServiceMethods[methodName] {
			continue
		}

		fqn := fmt.Sprintf("%s.%s.%s", packagePath, typeName, methodName)

		// Iterate inputs
		methodType := method.Type()

		// Create new method with cached flags
		methodID := hash.Fnv(fqn)
		if registeredID, ok := getRegisteredBindingMethodID(methodDef); ok {
			methodID = registeredID
		}

		boundMethod := &BoundMethod{
			ID:         methodID,
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
