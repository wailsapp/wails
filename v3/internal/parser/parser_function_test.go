package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
							Name:       "Greet",
							DocComment: "Greet someone",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Name: "string",
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
							Name:       "Greet",
							DocComment: "Greet does XYZ",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							ID: 1411160069,
						},
						{
							Name:       "NewPerson",
							DocComment: "NewPerson creates a new person",
							Inputs: []*Parameter{
								{
									Name: "name",
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
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
							Name:       "Yay",
							DocComment: "Yay does this and that",
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Address",
										IsStruct:  true,
										IsPointer: true,
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
						Name:       "Person",
						DocComment: "Person is a person",
						Fields: []*Field{
							{
								Name: "Name",
								Type: &ParameterType{
									Name: "string",
								},
							},
							{
								Name: "Address",
								Type: &ParameterType{
									Name:      "Address",
									IsStruct:  true,
									IsPointer: true,
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
									Name: "string",
								},
							},
							{
								Name: "State",
								Type: &ParameterType{
									Name: "string",
								},
							},
							{
								Name: "Country",
								Type: &ParameterType{
									Name: "string",
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

			cmpOptions := []cmp.Option{
				cmpopts.IgnoreTypes(Project{}, &Project{}, ParsedPackage{}, &ParsedPackage{}),
				cmpopts.IgnoreUnexported(Field{}),
			}

			// Loop over the things we want
			for packageName, packageData := range tt.wantBoundMethods {
				for structName, wantBoundMethods := range packageData {
					gotBoundMethods := got.BoundMethods[packageName][structName]
					if diff := cmp.Diff(wantBoundMethods, gotBoundMethods, cmpOptions...); diff != "" {
						t.Errorf("ParseDirectory() failed:\n" + diff)
					}
				}
			}

			if diff := cmp.Diff(tt.wantBoundMethods, got.BoundMethods, cmpOptions...); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
			}
			if diff := cmp.Diff(tt.wantModels, got.Models, cmpOptions...); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
			}
			if diff := cmp.Diff(tt.wantTypes, got.Types, cmpOptions...); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
			}
		})
	}

}
