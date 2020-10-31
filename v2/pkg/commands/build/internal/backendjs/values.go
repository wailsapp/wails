package backendjs

import "reflect"

var BoolValue bool = true
var IntValue int = 0
var Int8Value int8 = 0
var Int16Value int16 = 0
var Int32Value int32 = 0
var Int64Value int64 = 0
var UintValue uint = 0
var Uint8Value uint8 = 0
var Uint16Value uint16 = 0
var Uint32Value uint32 = 0
var Uint64Value uint64 = 0
var UintptrValue uintptr = 0
var Float32Value float32 = 0
var Float64Value float64 = 0
var Complex64Value complex64 = 0
var Complex128Value complex128 = 0
var StringValue string = ""

type Person struct {
	Name string
	Age  uint8
}

type GuestList struct {
	People []*Person
}

func NewPerson(name string, age uint8) reflect.Value {
	return reflect.New(reflect.TypeOf(&Person{
		Name: name,
		Age:  age,
	}))
}

var Bool = reflect.New(reflect.TypeOf(BoolValue))
var Int = reflect.New(reflect.TypeOf(IntValue))
var Int8 = reflect.New(reflect.TypeOf(Int8Value))
var Int16 = reflect.New(reflect.TypeOf(Int16Value))
var Int32 = reflect.New(reflect.TypeOf(Int32Value))
var Int64 = reflect.New(reflect.TypeOf(Int64Value))
var Uint = reflect.New(reflect.TypeOf(UintValue))
var Uint8 = reflect.New(reflect.TypeOf(Uint8Value))
var Uint16 = reflect.New(reflect.TypeOf(Uint16Value))
var Uint32 = reflect.New(reflect.TypeOf(Uint32Value))
var Uint64 = reflect.New(reflect.TypeOf(Uint64Value))
var Uintptr = reflect.New(reflect.TypeOf(UintptrValue))
var Float32 = reflect.New(reflect.TypeOf(Float32Value))
var Float64 = reflect.New(reflect.TypeOf(Float64Value))
var Complex64 = reflect.New(reflect.TypeOf(Complex64Value))
var Complex128 = reflect.New(reflect.TypeOf(Complex128Value))
var String = reflect.New(reflect.TypeOf(StringValue))
