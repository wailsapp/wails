package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFunction(t *testing.T) {
	tests := []struct {
		name             string
		dir              string
		wantBoundMethods map[string]map[string][]*BoundMethod
		wantEnums        map[string]map[string]*EnumDef
		wantModels       map[string]map[string]*StructDef
		wantTypes        map[string]map[string]*TypeDef
		wantErr          bool
	}{
		{
			name:    "should find a bound service returned from a function call",
			dir:     "testdata/function_single",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Package:    "main",
							Name:       "Greet",
							DocComment: "Greet someone",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							ID: 1411160069,
						},
					},
				},
			},
		},
		{
			name:    "should find a bound service returned from a function call in another package",
			dir:     "testdata/function_from_imported_package",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Package:    "main",
							Name:       "Greet",
							DocComment: "Greet does XYZ",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							ID: 1411160069,
						},
						{
							Package:    "main",
							Name:       "NewPerson",
							DocComment: "NewPerson creates a new person",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Package:   "main",
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							ID: 1661412647,
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services": {
					"OtherService": {
						{
							Package:    "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
							Name:       "Yay",
							DocComment: "Yay does this and that",
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Address",
										IsStruct:  true,
										IsPointer: true,
										Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
									},
								},
							},
							ID: 1592414782,
						},
					},
				},
			},
			wantModels: map[string]map[string]*StructDef{
				"main": {
					"Person": {
						Name:        "Person",
						DocComments: []string{"// Person is a person"},
						Fields: []*Field{
							{
								Name: "Name",
								Type: &ParameterType{
									Package: "main",
									Name:    "string",
								},
							},
							{
								Name: "Address",
								Type: &ParameterType{
									Name:      "Address",
									IsStruct:  true,
									IsPointer: true,
									Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services": {
					"Address": {
						Name: "Address",
						Fields: []*Field{
							{
								Name: "Street",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
									Name:    "string",
								},
							},
							{
								Name: "State",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
									Name:    "string",
								},
							},
							{
								Name: "Country",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services",
									Name:    "string",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProject(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Patch the PackageDir in the wantBoundMethods
			for _, packageData := range got.BoundMethods {
				for _, boundMethods := range packageData {
					for _, boundMethod := range boundMethods {
						boundMethod.PackageDir = ""
					}
				}
			}

			// Loop over the things we want
			for packageName, packageData := range tt.wantBoundMethods {
				for structName, wantBoundMethods := range packageData {
					gotBoundMethods := got.BoundMethods[packageName][structName]
					if diff := cmp.Diff(wantBoundMethods, gotBoundMethods); diff != "" {
						t.Errorf("ParseDirectory() failed:\n" + diff)
					}
				}
			}

			// Loop over the models
			for _, packageData := range got.Models {
				for _, wantModel := range packageData {
					// Loop over the Fields
					for _, field := range wantModel.Fields {
						field.Project = nil
					}
				}
			}

			if !reflect.DeepEqual(tt.wantBoundMethods, got.BoundMethods) {
				t.Errorf("ParseDirectory() failed:\n" + cmp.Diff(tt.wantBoundMethods, got.BoundMethods))
				//spew.Dump(tt.wantBoundMethods)
				//spew.Dump(got.BoundMethods)
			}
			if !reflect.DeepEqual(tt.wantModels, got.Models) {
				t.Errorf("ParseDirectory() failed:\n" + cmp.Diff(tt.wantModels, got.Models))
			}
			if diff := cmp.Diff(tt.wantTypes, got.Types); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
			}
		})
	}

}
