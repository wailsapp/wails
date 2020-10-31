package backendjs

import "reflect"

type Struct struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name string
	Type reflect.Value
}
