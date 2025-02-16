package main

import (
	_ "embed"
	"log"

	nobindingshere "github.com/wailsapp/wails/v3/internal/generator/testcases/no_bindings_here"
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

// A generic alias that forwards to a type parameter.
// type GenericAlias[T any] = T

// A generic alias that wraps a pointer type.
// type GenericPtrAlias[T any] = *GenericAlias[T]

// A generic alias that wraps a map.
// type GenericMapAlias[T interface {
// 	comparable
// 	encoding.TextMarshaler
// }, U any] = map[T]U

// A generic alias that wraps a generic struct.
// type GenericPersonAlias[T any] = GenericPerson[[]GenericPtrAlias[T]]

// An alias that wraps a class through a non-typeparam alias.
// type IndirectPersonAlias = GenericPersonAlias[bool]

// An alias that wraps a class through a typeparam alias.
// type TPIndirectPersonAlias = GenericAlias[GenericPerson[bool]]

// A class whose fields have various aliased types.
// type AliasGroup struct {
// 	GAi   GenericAlias[int]
// 	GAP   GenericAlias[GenericPerson[bool]]
// 	GPAs  GenericPtrAlias[[]string]
// 	GPAP  GenericPtrAlias[GenericPerson[[]int]]
// 	GMA   GenericMapAlias[struct{ encoding.TextMarshaler }, float32]
// 	GPA   GenericPersonAlias[bool]
// 	IPA   IndirectPersonAlias
// 	TPIPA TPIndirectPersonAlias
// }

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

func (GreetService) GetButForeignPrivateAlias() (_ nobindingshere.PrivatePerson) {
	return
}

// func (GreetService) GetButGenericAliases() (_ AliasGroup) {
//	return
// }

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
