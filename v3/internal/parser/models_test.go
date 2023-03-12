package parser

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
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

//func TestGenerateClass(t *testing.T) {
//	person := StructDef{
//		Name: "Person",
//		Fields: []*Field{
//			{
//				Name: "Name",
//				Type: &ParameterType{
//					Name: "string",
//				},
//			},
//			{
//				Name: "Parent",
//				Type: &ParameterType{
//					Name:      "Person",
//					IsStruct:  true,
//					IsPointer: true,
//					Package:   "main",
//				},
//			},
//			{
//				Name: "Details",
//				Type: &ParameterType{
//					Name:     "anon1",
//					IsStruct: true,
//					Package:  "main",
//				},
//			},
//			{
//				Name: "Address",
//				Type: &ParameterType{
//					Name:      "Address",
//					IsStruct:  true,
//					IsPointer: true,
//					Package:   "github.com/some/other/package",
//				},
//			},
//		},
//	}
//	anon1 := StructDef{
//		Name: "anon1",
//		Fields: []*Field{
//			{
//				Name: "Age",
//				Type: &ParameterType{
//					Name: "int",
//				},
//			},
//			{
//				Name: "Address",
//				Type: &ParameterType{
//					Name: "string",
//				},
//			},
//		},
//	}
//
//	var builder strings.Builder
//	models := make(map[string]*StructDef)
//	models["Person"] = &person
//	models["anon1"] = &anon1
//	def := ModelDefinitions{
//		Package: "main",
//		Models:  models,
//	}
//
//	err := GenerateModel(&builder, &def)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	text := builder.String()
//	println("Built string")
//	println(text)
//	if diff := cmp.Diff(expected, text); diff != "" {
//		t.Errorf("GenerateClass() failed:\n" + diff)
//	}
//}

func TestGenerateModels(t *testing.T) {

	tests := []struct {
		dir  string
		want string
	}{
		{
			"testdata/function_single",
			getFile("testdata/function_single/models.ts"),
		},
		{
			"testdata/function_from_imported_package",
			getFile("testdata/function_from_imported_package/models.ts"),
		},
		{
			"testdata/variable_single",
			getFile("testdata/variable_single/models.ts"),
		},
		{
			"testdata/variable_single_from_function",
			getFile("testdata/variable_single_from_function/models.ts"),
		},
		{
			"testdata/variable_single_from_other_function",
			getFile("testdata/variable_single_from_other_function/models.ts"),
		},
		{
			"testdata/struct_literal_single",
			getFile("testdata/struct_literal_single/models.ts"),
		},
		{
			"testdata/struct_literal_multiple",
			getFile("testdata/struct_literal_multiple/models.ts"),
		},
		{
			"testdata/struct_literal_multiple_other",
			getFile("testdata/struct_literal_multiple_other/models.ts"),
		},
		{
			"testdata/struct_literal_multiple_files",
			getFile("testdata/struct_literal_multiple_files/models.ts"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			// Run parser on directory
			project, err := ParseProject(tt.dir)
			if err != nil {
				t.Fatalf("ParseProject() error = %v", err)
			}

			// Generate Models
			got, err := GenerateModels(project.Models)
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				err = os.WriteFile(filepath.Join(tt.dir, "models.got.ts"), []byte(got), 0644)
				if err != nil {
					t.Errorf("os.WriteFile() error = %v", err)
					return
				}
				t.Fatalf("GenerateModels() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
