package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is great
type GreetService int

// A nice type Alias.
type Alias = int

// A class alias.
type AliasedPerson = Person

// An empty struct alias.
type EmptyAliasStruct = struct{}

// A struct alias.
// This should be rendered as a typedef or interface in every mode.
type AliasStruct = struct {
	// A field with a comment.
	Foo      []int
	Bar, Baz string `json:",omitempty"` // Definitely not Foo.

	Other OtherAliasStruct // A nested alias struct.
}

// Another struct alias.
type OtherAliasStruct = struct {
	NoMoreIdeas []rune
}

// An empty struct.
type EmptyStruct struct{}

// A non-generic struct containing an alias.
type Person struct {
	Name         string // The Person's name.
	AliasedField Alias  // A random alias field.
}

// A generic struct containing an alias.
type GenericPerson[T any] struct {
	Name         T
	AliasedField Alias
}

// Another class alias, but ordered after its aliased class.
type StrangelyAliasedPerson = Person

// Get someone.
func (GreetService) Get(aliasValue Alias) Person {
	return Person{"hello", aliasValue}
}

// Get someone quite different.
func (GreetService) GetButDifferent() GenericPerson[bool] {
	return GenericPerson[bool]{true, 13}
}

// Apparently, aliases are all the rage right now.
func (GreetService) GetButAliased(p AliasedPerson) StrangelyAliasedPerson {
	return p
}

// Greet a lot of unusual things.
func (GreetService) Greet(EmptyAliasStruct, EmptyStruct) AliasStruct {
	return AliasStruct{}
}

func NewGreetService() application.Service {
	return application.NewService(new(GreetService))
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			NewGreetService(),
		},
	})

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
