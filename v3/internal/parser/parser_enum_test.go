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
			name:    "should find a bound services with an enum",
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
										Package: "main",
										Name:    "string",
									},
								},
								{
									Name: "title",
									Type: &ParameterType{
										Package: "main",
										Name:    "Title",
										IsEnum:  true,
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
									Type: &ParameterType{
										Package:   "main",
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
						Name: "Title",
						Type: "string",
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
						Name: "Person",
						Fields: []*Field{
							{
								Name: "Title",
								Type: &ParameterType{
									Package: "main",
									Name:    "Title",
									IsEnum:  true,
								},
							},
							{
								Name: "Name",
								Type: &ParameterType{
									Package: "main",
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
