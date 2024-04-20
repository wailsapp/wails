// This test produces different bindings depending on whether GODEBUG=gotypesalias=1 is set
// https://pkg.go.dev/go/types#Alias
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Person struct {
	Name string `json:"name"`
}

// Human is a copy of Person
// it is only visible without GODEBUG=gotypesalias=1
type Human Person

// APerson is an alias of Human
// it is only visible with GODEBUG=gotypesalias=1
type APerson = Human

type AString = string
type AIntSlice = []int
type AStringMap = map[string]string

type Aliases struct {
	Name     AString
	Numbers  AIntSlice
	Settings AStringMap
	Person   APerson
}

type GreetService struct{}

func (*GreetService) Greet(person APerson) string {
	return "Hello " + person.Name
}

func (*GreetService) UnwrapAliases(aliases Aliases) (AString, AIntSlice, AStringMap, APerson) {
	return aliases.Name, aliases.Numbers, aliases.Settings, aliases.Person
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
