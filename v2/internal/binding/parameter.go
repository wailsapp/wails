package binding

import "reflect"

// Parameter defines a Go method parameter
type Parameter struct {
	Name        string `json:"name,omitempty"`
	TypeName    string `json:"type"`
	reflectType reflect.Type
}

func newParameter(Name string, Type reflect.Type) *Parameter {
	return &Parameter{
		Name:        Name,
		TypeName:    Type.String(),
		reflectType: Type,
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
