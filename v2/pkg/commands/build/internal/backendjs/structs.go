package backendjs

import (
	"reflect"
	"strings"
)

// Parameter defines a parameter used by a struct method
type Parameter struct {
	Name string
	Type reflect.Kind
}

// JSType returns the Javascript equivalent of the
// parameter type
func (p *Parameter) JSType() string {
	return string(goTypeToJS(p.Type))
}

// Method defines a struct method
type Method struct {
	Name     string
	Inputs   []*Parameter
	Outputs  []*Parameter
	Comments []string
}

// InputsAsJSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) InputsAsJSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputs = append(inputs, input.Name)
	}

	return strings.Join(inputs, ", ")
}

// func generateStructFile() {
// 	// Create string buffer
// 	var result bytes.Buffer

// 	// Add some standard comments
// 	_, err := result.WriteString(structJSHeader + )
// 	if err != nil {
// 		return errors.Wrap(err, "Error writing string")
// 	}

// 	// Loop over the methods
// 	for _, method := range methods {
// 		generatedCode := generateMethodWrapper(method) {

// 		}
// 	}
// 	return nil
// }
