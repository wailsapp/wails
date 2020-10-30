package backendjs

import (
	"reflect"
)

// Parameter defines a parameter used by a struct method
type Parameter struct {
	Name string
	Type reflect.Type
}

// Method defines a struct method
type Method struct {
	Name     string
	Inputs   []*Parameter
	Outputs  []*Parameter
	Comments []string
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
