package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/internal/parser/testdata/enum_from_imported_package/services"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
}

// Greet does XYZ
func (*GreetService) Greet(name string, title services.Title) string {
	return "Hello " + title.String() + " " + name
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
