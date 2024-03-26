package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseEnum(t *testing.T) {
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
			name:    "should find a bound service with an enum",
			dir:     "testdata/enum",
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
								{
									Name: "title",
									Type: &ParameterType{
										Name:   "Title",
										IsEnum: true,
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
									Type: &ParameterType{
										Name:      "Person",
										IsStruct:  true,
										IsPointer: true,
									},
								},
							},
							ID: 1661412647,
						},
					},
				},
			},
			wantTypes: map[string]map[string]*TypeDef{
				"main": {
					"Title": {
						Name:       "Title",
						DocComment: "Title is a title",
						Type:       "string",
						Consts: []*ConstDef{
							{
								Name:       "Mister",
								DocComment: "Mister is a title",
								Value:      `"Mr"`,
							},
							{
								Name:  "Miss",
								Value: `"Miss"`,
							},
							{
								Name:  "Ms",
								Value: `"Ms"`,
							},
							{
								Name:  "Mrs",
								Value: `"Mrs"`,
							},
							{
								Name:  "Dr",
								Value: `"Dr"`,
							},
						},
						ShouldGenerate: true,
					},
				},
			},
			wantModels: map[string]map[string]*StructDef{
				"main": {
					"Person": {
						Name:       "Person",
						DocComment: "Person represents a person",
						Fields: []*Field{
							{
								Name: "Title",
								Type: &ParameterType{
									Name:   "Title",
									IsEnum: true,
								},
							},
							{
								Name: "Name",
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

			patchParserOutput(got)

			// Loop over the things we want
			for packageName, packageData := range tt.wantBoundMethods {
				for structName, wantBoundMethods := range packageData {
					gotBoundMethods := got.BoundMethods[packageName][structName]
					if diff := cmp.Diff(wantBoundMethods, gotBoundMethods, cmp.AllowUnexported(Parameter{})); diff != "" {
						t.Errorf("ParseDirectory() failed:\n" + diff)
					}
				}
			}

			if diff := cmp.Diff(tt.wantBoundMethods, got.BoundMethods, cmp.AllowUnexported(Parameter{})); diff != "" {
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
