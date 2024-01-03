package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseStructLiteralMultiple(t *testing.T) {
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
			name: "should find multiple bound services",
			dir:  "testdata/struct_literal_multiple",
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Package:    "main",
							Name:       "Greet",
							DocComment: "",
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
					"OtherService": {
						{
							Package: "main",
							Name:    "Hello",
							ID:      4249972365,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should find multiple bound services over multiple files",
			dir:  "testdata/struct_literal_multiple_files",
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Package:    "main",
							Name:       "Greet",
							DocComment: "",
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
					"OtherService": {
						{
							Package: "main",
							Name:    "Hello",
							ID:      4249972365,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "should find multiple bound services over multiple packages",
			dir:     "testdata/struct_literal_multiple_other",
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
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services": {
					"OtherService": {
						{
							Package:    "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
							Name:       "Yay",
							DocComment: "Yay does this and that",
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Address",
										IsStruct:  true,
										IsPointer: true,
										Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
									},
								},
							},
							ID: 469445984,
						},
					},
				},
			},
			wantModels: map[string]map[string]*StructDef{
				"main": {
					"Person": {
						Name: "Person",
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
									Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services": {
					"Address": {
						Name: "Address",
						Fields: []*Field{
							{
								Name: "Street",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
									Name:    "string",
								},
							},
							{
								Name: "State",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
									Name:    "string",
								},
							},
							{
								Name: "Country",
								Type: &ParameterType{
									Package: "github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services",
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

			if diff := cmp.Diff(tt.wantBoundMethods, got.BoundMethods); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
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
