package main

import (
	_ "embed"
	services2 "github.com/wailsapp/wails/v3/internal/parser/testdata/enum_from_imported_package_same_name/other/services"
	"github.com/wailsapp/wails/v3/internal/parser/testdata/enum_from_imported_package_same_name/services"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService struct {
	SomeVariable int
	lowerCase    string
}

// Greet does XYZ
func (*GreetService) Greet(name string, title services.Title, title2 services2.Title2) string {
	return "Hello " + title.String() + " " + name + title2.String()
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
