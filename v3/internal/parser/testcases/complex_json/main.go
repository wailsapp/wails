package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Title is a title
type Title string

const (
	// Mister is a title
	Mister Title = "Mr"
	Miss   Title = "Miss"
	Ms     Title = "Ms"
	Mrs    Title = "Mrs"
	Dr     Title = "Dr"
)

// GreetService is great
type GreetService struct{}

type Embedded1 struct {
	// Friends should be shadowed in Person by a field of lesser depth
	Friends int

	// Vanish should be omitted from Person because there is another field with same depth and no tag
	Vanish float32

	// StillThere should be shadowed in Person by other field with same depth and a json tag
	StillThere string

	// embedded4 should effectively appear as an embedded field
	embedded4

	// unexported should be invisible
	unexported bool
}

type Embedded2 struct {
	// Vanish should be omitted from Person because there is another field with same depth and no tag
	Vanish bool

	// StillThereButRenamed should shadow in Person the other field with same depth and no json tag
	StillThereButRenamed *Embedded3 `json:"StillThere"`
}

type Embedded3 string

// Person represents a person
type Person struct {
	// Titles is optional in JSON
	Titles []Title `json:",omitempty"`

	// Names has a
	// multiline comment
	Names []string

	// Partner has a custom and complex JSON key
	Partner *Person `json:"the person's partner ❤️"`
	Friends []*Person

	Embedded1
	Embedded2

	// UselessMap is invisible to JSON
	UselessMap map[int]bool `json:"-"`

	// StrangeNumber maps to "-"
	StrangeNumber float32 `json:"-,"`

	// Embedded3 should appear with key "Embedded3"
	Embedded3

	// StrangerNumber is serialized as a string
	StrangerNumber int `json:",string"`
	// StrangestString is optional and serialized as a JSON string
	StrangestString string `json:",omitempty,string"`
	// StringStrangest is serialized as a JSON string and optional
	StringStrangest string `json:",string,omitempty"`

	// unexportedToo should be invisible even with a json tag
	unexportedToo bool `json:"Unexported"`

	// embedded4 should be optional and appear with key "emb4"
	embedded4 `json:"emb4,omitempty"`
}

type embedded4 struct {
	// NamingThingsIsHard is a law of programming
	NamingThingsIsHard bool `json:",string"`

	// Friends should not be shadowed in Person as embedded4 is not embedded
	// from encoding/json's point of view;
	// however, it should be shadowed in Embedded1
	Friends bool

	// Embedded string should be invisible because it's unexported
	string
}

// Greet does XYZ
func (*GreetService) Greet(person Person, emb Embedded1) string {
	return ""
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
