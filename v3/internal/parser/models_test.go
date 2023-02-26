package parser

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

const expected = `
  export class Person {
    name: string;
    parent: main.Person;
    details: main.anon1;
    
    static createFrom(source: any = {}) {
      return new Person(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) {
        source = JSON.parse(source);
      }

      this.name = source["name"]
      this.parent = source["parent"]
      this.details = source["details"]
      
    }
  }
`

func TestGenerateClass(t *testing.T) {
	person := StructDef{
		Name: "Person",
		Fields: []*Field{
			{
				Name: "Name",
				Type: &ParameterType{
					Name: "string",
				},
			},
			{
				Name: "Parent",
				Type: &ParameterType{
					Name:      "Person",
					IsStruct:  true,
					IsPointer: true,
					Package:   "main",
				},
			},
			{
				Name: "Details",
				Type: &ParameterType{
					Name:     "anon1",
					IsStruct: true,
					Package:  "main",
				},
			},
		},
	}

	var builder strings.Builder
	err := GenerateClass(&builder, &person)
	if err != nil {
		t.Fatal(err)
	}

	text := builder.String()
	println("Built string")
	println(text)
	if diff := cmp.Diff(strings.TrimPrefix(expected, "\n"), text); diff != "" {
		t.Errorf("GenerateClass() failed:\n" + diff)
	}
}
