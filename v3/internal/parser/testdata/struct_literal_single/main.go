package main

import (
	_ "embed"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name string
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

func (*GreetService) StructPointerInputErrorOutput(in *Person) error {
	return nil
}

func (*GreetService) StructPointerInputStructPointerOutput(in *Person) *Person {
	return in
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
