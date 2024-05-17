package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct{}

// Person represents a person
type Person struct {
	Name string
}

// Greet does XYZ
// It has a multiline doc comment
// The comment has even some */ traps!!
func (*GreetService) Greet(str string, people []Person, _ struct {
	AnotherCount int
	AnotherOne   *Person
}, assoc map[int]*bool, _ []*float32, other ...string) (person Person, _ any, err1 error, _ []int, err error) {
	return
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
