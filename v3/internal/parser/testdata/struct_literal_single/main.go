package main

import (
	_ "embed"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name    string
	Parent  *Person
	Details struct {
		Age     int
		Address struct {
			Street string
		}
	}
}

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
}

// Greet someone
func (*GreetService) Greet(name string) string {
	return "Hello " + name
}

func (*GreetService) NoInputsStringOut() string {
	return "Hello"
}

func (*GreetService) StringArrayInputStringOut(in []string) string {
	return strings.Join(in, ",")
}

func (*GreetService) StringArrayInputStringArrayOut(in []string) []string {
	return in
}

func (*GreetService) StringArrayInputNamedOutput(in []string) (output []string) {
	return in
}

func (*GreetService) StringArrayInputNamedOutputs(in []string) (output []string, err error) {
	return in, nil
}

func (*GreetService) IntPointerInputNamedOutputs(in *int) (output *int, err error) {
	return in, nil
}

func (*GreetService) UIntPointerInAndOutput(in *uint) *uint {
	return in
}

func (*GreetService) UInt8PointerInAndOutput(in *uint8) *uint8 {
	return in
}

func (*GreetService) UInt16PointerInAndOutput(in *uint16) *uint16 {
	return in
}

func (*GreetService) UInt32PointerInAndOutput(in *uint32) *uint32 {
	return in
}

func (*GreetService) UInt64PointerInAndOutput(in *uint64) *uint64 {
	return in
}

func (*GreetService) IntPointerInAndOutput(in *int) *int {
	return in
}

func (*GreetService) Int8PointerInAndOutput(in *int8) *int8 {
	return in
}

func (*GreetService) Int16PointerInAndOutput(in *int16) *int16 {
	return in
}

func (*GreetService) Int32PointerInAndOutput(in *int32) *int32 {
	return in
}

func (*GreetService) Int64PointerInAndOutput(in *int64) *int64 {
	return in
}

func (*GreetService) IntInIntOut(in int) int {
	return in
}

func (*GreetService) Int8InIntOut(in int8) int8 {
	return in
}
func (*GreetService) Int16InIntOut(in int16) int16 {
	return in
}
func (*GreetService) Int32InIntOut(in int32) int32 {
	return in
}
func (*GreetService) Int64InIntOut(in int64) int64 {
	return in
}

func (*GreetService) UIntInUIntOut(in uint) uint {
	return in
}

func (*GreetService) UInt8InUIntOut(in uint8) uint8 {
	return in
}
func (*GreetService) UInt16InUIntOut(in uint16) uint16 {
	return in
}
func (*GreetService) UInt32InUIntOut(in uint32) uint32 {
	return in
}
func (*GreetService) UInt64InUIntOut(in uint64) uint64 {
	return in
}

func (*GreetService) Float32InFloat32Out(in float32) float32 {
	return in
}

func (*GreetService) Float64InFloat64Out(in float64) float64 {
	return in
}

func (*GreetService) PointerFloat32InFloat32Out(in *float32) *float32 {
	return in
}

func (*GreetService) PointerFloat64InFloat64Out(in *float64) *float64 {
	return in
}

func (*GreetService) BoolInBoolOut(in bool) bool {
	return in
}

func (*GreetService) PointerBoolInBoolOut(in *bool) *bool {
	return in
}

func (*GreetService) PointerStringInStringOut(in *string) *string {
	return in
}

func (*GreetService) StructPointerInputErrorOutput(in *Person) error {
	return nil
}

func (*GreetService) StructInputStructOutput(in Person) Person {
	return in
}

func (*GreetService) StructPointerInputStructPointerOutput(in *Person) *Person {
	return in
}

func (*GreetService) MapIntInt(in map[int]int) {
}

func (*GreetService) PointerMapIntInt(in *map[int]int) {
}

func (*GreetService) MapIntPointerInt(in map[*int]int) {
}

func (*GreetService) MapIntSliceInt(in map[int][]int) {
}

func (*GreetService) MapIntSliceIntInMapIntSliceIntOut(in map[int][]int) (out map[int][]int) {
	return nil
}

func (*GreetService) ArrayInt(in [4]int) {
}

func main() {
	app := application.New(application.Options{
		Bind: []interface{}{
			&GreetService{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
