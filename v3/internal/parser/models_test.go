package parser

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

const expected = `
export namespace main {
  
  export class Person {
    name: string;
    parent: Person;
    details: anon1;
    address: package.Address;
    
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
      this.address = source["address"]
      
    }
  }
  
  export class anon1 {
    age: int;
    address: string;
    
    static createFrom(source: any = {}) {
      return new anon1(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) {
        source = JSON.parse(source);
      }

      this.age = source["age"]
      this.address = source["address"]
      
    }
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
			{
				Name: "Address",
				Type: &ParameterType{
					Name:      "Address",
					IsStruct:  true,
					IsPointer: true,
					Package:   "github.com/some/other/package",
				},
			},
		},
	}
	anon1 := StructDef{
		Name: "anon1",
		Fields: []*Field{
			{
				Name: "Age",
				Type: &ParameterType{
					Name: "int",
				},
			},
			{
				Name: "Address",
				Type: &ParameterType{
					Name: "string",
				},
			},
		},
	}

	var builder strings.Builder
	models := make(map[string]*StructDef)
	models["Person"] = &person
	models["anon1"] = &anon1
	def := ModelDefinitions{
		Package: "main",
		Models:  models,
	}

	err := GenerateModel(&builder, &def)
	if err != nil {
		t.Fatal(err)
	}

	text := builder.String()
	println("Built string")
	println(text)
	if diff := cmp.Diff(expected, text); diff != "" {
		t.Errorf("GenerateClass() failed:\n" + diff)
	}
}
