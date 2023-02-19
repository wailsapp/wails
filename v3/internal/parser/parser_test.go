package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseDirectory(t *testing.T) {
	tests := []struct {
		name             string
		dir              string
		wantBoundMethods map[string]map[string][]*BoundMethod
		wantErr          bool
	}{
		{
			name: "should find single bound service",
			dir:  "testdata/struct_literal_single",
			//wantModels: []string{"main.GreetService"},
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "Greet someone\n",
							Inputs: []*Parameter{
								{
									Name:     "name",
									Type:     "string",
									IsStruct: false,
									IsSlice:  false,
								},
							},
							Outputs: []*Parameter{
								{
									Name:     "",
									Type:     "string",
									IsStruct: false,
									IsSlice:  false,
								},
							},
						},
						{
							Name:       "NoInputsStringOut",
							DocComment: "",
							Inputs:     nil,
							Outputs: []*Parameter{
								{
									Name:     "",
									Type:     "string",
									IsStruct: false,
									IsSlice:  false,
								},
							},
						},
						{
							Name: "StringArrayInputStringOut",
							Inputs: []*Parameter{
								{
									Name:    "in",
									Type:    "string",
									IsSlice: true,
								},
							},
							Outputs: []*Parameter{
								{
									Type: "string",
								},
							},
						},
						{
							Name: "StringArrayInputStringArrayOut",
							Inputs: []*Parameter{
								{
									Name:    "in",
									Type:    "string",
									IsSlice: true,
								},
							},
							Outputs: []*Parameter{
								{
									Type:    "string",
									IsSlice: true,
								},
							},
						},
						{
							Name: "StringArrayInputNamedOutput",
							Inputs: []*Parameter{
								{
									Name:    "in",
									Type:    "string",
									IsSlice: true,
								},
							},
							Outputs: []*Parameter{
								{
									Name:    "output",
									Type:    "string",
									IsSlice: true,
								},
							},
						},
						{
							Name: "StringArrayInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name:    "in",
									Type:    "string",
									IsSlice: true,
								},
							},
							Outputs: []*Parameter{
								{
									Name:    "output",
									Type:    "string",
									IsSlice: true,
								},
								{
									Name: "err",
									Type: "error",
								},
							},
						},
						{
							Name: "IntPointerInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name:      "in",
									Type:      "int",
									IsPointer: true,
								},
							},
							Outputs: []*Parameter{
								{
									Name:      "output",
									Type:      "int",
									IsPointer: true,
								},
								{
									Name: "err",
									Type: "error",
								},
							},
						},
						{
							Name: "StructPointerInputErrorOutput",
							Inputs: []*Parameter{
								{
									Name:      "in",
									Type:      "Person",
									IsStruct:  true,
									IsPointer: true,
								},
							},
							Outputs: []*Parameter{
								{
									Type: "error",
								},
							},
						},
						{
							Name: "StructPointerInputStructPointerOutput",
							Inputs: []*Parameter{
								{
									Name:      "in",
									Type:      "Person",
									IsStruct:  true,
									IsPointer: true,
								},
							},
							Outputs: []*Parameter{
								{
									Type:      "Person",
									IsPointer: true,
									IsStruct:  true,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		//{
		//	name: "should find multiple bound services",
		//	dir:  "testdata/struct_literal_multiple",
		//	//wantModels: []string{"main.GreetService", "main.OtherService"},
		//	wantErr: false,
		//},
		//{
		//	name: "should find multiple bound services over multiple files",
		//	dir:  "testdata/struct_literal_multiple_files",
		//	//wantModels: []string{"main.GreetService", "main.OtherService"},
		//	wantErr: false,
		//},
		//{
		//	name: "should find multiple bound services over multiple packages",
		//	dir:  "testdata/struct_literal_multiple_other",
		//	//wantModels: []string{"main.GreetService", "services.OtherService", "main.Person"},
		//	wantErr: false,
		//},
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
		})
	}

}

//func TestGenerateTypeScript(t *testing.T) {
//	tests := []struct {
//		name    string
//		dir     string
//		wantModels    string
//		wantErr bool
//	}{
//		{
//			name: "should find single bound service",
//			dir:  "testdata/struct_literal_single",
//			wantModels: `namespace main {
//  class GreetService {
//    SomeVariable: number;
//  }
//}
//`,
//			wantErr: false,
//		},
//		{
//			name: "should find multiple bound services",
//			dir:  "testdata/struct_literal_multiple",
//			wantModels: `namespace main {
//  class GreetService {
//    SomeVariable: number;
//  }
//  class OtherService {
//  }
//}
//`,
//			wantErr: false,
//		},
//		{
//			name: "should find multiple bound services over multiple files",
//			dir:  "testdata/struct_literal_multiple_files",
//			wantModels: `namespace main {
//  class GreetService {
//    SomeVariable: number;
//  }
//  class OtherService {
//  }
//}
//`,
//			wantErr: false,
//		},
//		{
//			name: "should find bound services from other packages",
//			dir:  "../../examples/binding",
//			wantModels: `namespace main {
//  class localStruct {
//  }
//}
//namespace models {
//  class Person {
//    Name: string;
//  }
//}
//namespace services {
//  class GreetService {
//    SomeVariable: number;
//    Parent: models.Person;
//  }
//}
//`,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			Debug = true
//			context, err := ParseDirectory(tt.dir)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			ts, err := GenerateModels(context)
//			require.NoError(t, err)
//			require.Equal(t, tt.wantModels, string(ts))
//
//		})
//	}
//}
