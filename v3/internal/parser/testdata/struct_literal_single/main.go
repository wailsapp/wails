package main

import (
	_ "embed"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name   string
	Parent *Person
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

func (*GreetService) IntInIntOut(in int) int {
	return in
}

func (*GreetService) UIntInUIntOut(in uint) uint {
	return in
}
func (*GreetService) Float32InFloat32Out(in float32) float32 {
	return in
}
func (*GreetService) Float64InFloat64Out(in float64) float64 {
	return in
}

func (*GreetService) BoolInBoolOut(in bool) bool {
	return in
}

func (*GreetService) StructPointerInputErrorOutput(in *Person) error {
	return nil
}

func (*GreetService) StructPointerInputStructPointerOutput(in *Person) *Person {
	return in
}

func (*GreetService) MapIntInt(in map[int]int) {
}

func (*GreetService) MapIntSliceInt(in map[int][]int) {
}

func (*GreetService) MapIntSliceIntInMapIntSliceIntOut(in map[int][]int) (out map[int][]int) {
	return nil
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
