package main

import (
	_ "embed"
	"log"

	nobindingshere "github.com/wailsapp/wails/v3/internal/parser/testcases/no_bindings_here"
	"github.com/wailsapp/wails/v3/internal/parser/testcases/no_bindings_here/other"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

// EmbedService is tricky.
type EmbedService struct {
	nobindingshere.SomeMethods
}

// EmbedOther is even trickier.
type EmbedOther struct {
	other.OtherMethods
}

// Greet someone
func (*GreetService) Greet(string) {}

func main() {
	app := application.New(application.Options{
		Bind: []interface{}{
			new(GreetService),
			&EmbedService{},
			&EmbedOther{},
			&nobindingshere.SomeMethods{},
			&other.OtherMethods{},
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
