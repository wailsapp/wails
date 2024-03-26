package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseStructLiteralSingle(t *testing.T) {
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
			name: "should find single bound service",
			dir:  "testdata/struct_literal_single",
			//wantModels: []string{"main.GreetService"},
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
						{
							Name:       "NoInputsStringOut",
							DocComment: "",
							Inputs:     nil,
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							ID: 1075577233,
						},
						{
							Name: "StringArrayInputStringOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "string",
									},
								},
							},
							ID: 1091960237,
						},
						{
							Name: "StringArrayInputStringArrayOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							ID: 383995060,
						},
						{
							Name: "StringArrayInputNamedOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							ID: 3678582682,
						},
						{
							Name: "StringArrayInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Name:    "string",
										IsSlice: true,
									},
								},
								{
									Name: "err",
									Type: &ParameterType{
										Name: "error",
									},
								},
							},
							ID: 319259595,
						},
						{
							Name: "IntPointerInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Name:      "int",
										IsPointer: true,
									},
								},
								{
									Name: "err",
									Type: &ParameterType{
										Name: "error",
									}},
							},
							ID: 2718999663,
						},
						{
							Name: "UIntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "uint",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "uint",
										IsPointer: true,
									},
								},
							},
							ID: 1367187362,
						},
						{
							Name: "UInt8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "uint8",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "uint8",
										IsPointer: true,
									},
								},
							},
							ID: 518250834,
						},
						{
							Name: "UInt16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "uint16",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "uint16",
										IsPointer: true,
									},
								},
							},
							ID: 1236957573,
						},
						{
							Name: "UInt32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "uint32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "uint32",
										IsPointer: true,
									},
								},
							},
							ID: 1739300671,
						},
						{
							Name: "UInt64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "uint64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "uint64",
										IsPointer: true,
									},
								},
							},
							ID: 1403757716,
						},
						{
							Name: "IntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							ID: 1066151743,
						},
						{
							Name: "Int8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int8",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "int8",
										IsPointer: true,
									},
								},
							},
							ID: 2189402897,
						},
						{
							Name: "Int16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int16",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "int16",
										IsPointer: true,
									},
								},
							},
							ID: 1754277916,
						},
						{
							Name: "Int32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "int32",
										IsPointer: true,
									},
								},
							},
							ID: 4251088558,
						},
						{
							Name: "Int64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "int64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "int64",
										IsPointer: true,
									},
								},
							},
							ID: 2205561041,
						},
						{
							Name: "IntInIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "int",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "int",
									},
								},
							},
							ID: 642881729,
						},
						{
							Name: "Int8InIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "int8",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "int8",
									},
								},
							},
							ID: 572240879,
						},
						{
							Name: "Int16InIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "int16",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "int16",
									},
								},
							},
							ID: 3306292566,
						},
						{
							Name: "Int32InIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "int32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "int32",
									},
								},
							},
							ID: 1909469092,
						},
						{
							Name: "Int64InIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "int64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "int64",
									},
								},
							},
							ID: 1343888303,
						},
						{
							Name: "UIntInUIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "uint",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "uint",
									},
								},
							},
							ID: 2836661285,
						},
						{
							Name: "UInt8InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "uint8",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "uint8",
									},
								},
							},
							ID: 2988345717,
						},
						{
							Name: "UInt16InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "uint16",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "uint16",
									},
								},
							},
							ID: 3401034892,
						},
						{
							Name: "UInt32InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "uint32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "uint32",
									},
								},
							},
							ID: 1160383782,
						},
						{
							Name: "UInt64InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "uint64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "uint64",
									},
								},
							},
							ID: 793803239,
						},
						{
							Name: "Float32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "float32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "float32",
									},
								},
							},
							ID: 3132595881,
						},
						{
							Name: "Float64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "float64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "float64",
									},
								},
							},
							ID: 2182412247,
						},
						{
							Name: "PointerFloat32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "float32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "float32",
										IsPointer: true,
									},
								},
							},
							ID: 224675106,
						},
						{
							Name: "PointerFloat64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "float64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "float64",
										IsPointer: true,
									},
								},
							},
							ID: 2124953624,
						},
						{
							Name: "BoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "bool",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "bool",
									},
								},
							},
							ID: 2424639793,
						},
						{
							Name: "PointerBoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "bool",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "bool",
										IsPointer: true,
									},
								},
							},
							ID: 3589606958,
						},
						{
							Name: "PointerStringInStringOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "string",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "string",
										IsPointer: true,
									},
								},
							},
							ID: 229603958,
						},
						{
							Name: "StructPointerInputErrorOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name: "error",
									},
								},
							},
							ID: 2447692557,
						},
						{
							Name: "StructInputStructOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:     "Person",
										IsStruct: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:     "Person",
										IsStruct: true,
									},
								},
							},
							ID: 3835643147,
						},
						{
							Name: "StructPointerInputStructPointerOutput",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							ID: 2943477349,
						},
						{
							Name: "MapIntInt",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "map",
										MapKey: &ParameterType{
											Name: "int",
										},
										MapValue: &ParameterType{
											Name: "int",
										},
									},
								},
							},
							ID: 2386486356,
						},
						{
							Name: "PointerMapIntInt",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:      "map",
										IsPointer: true,
										MapKey: &ParameterType{
											Name: "int",
										},
										MapValue: &ParameterType{
											Name: "int",
										},
									},
								},
							},
							ID: 3516977899,
						},
						{
							Name: "MapIntPointerInt",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "map",
										MapKey: &ParameterType{
											Name:      "int",
											IsPointer: true,
										},
										MapValue: &ParameterType{
											Name: "int",
										},
									},
								},
							},
							ID: 550413585,
						},
						{
							Name: "MapIntSliceInt",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "map",
										MapKey: &ParameterType{
											Name: "int",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
										},
									},
								},
							},
							ID: 2900172572,
						},
						{
							Name: "MapIntSliceIntInMapIntSliceIntOut",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name: "map",
										MapKey: &ParameterType{
											Name: "int",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
										},
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "out",
									Type: &ParameterType{
										Name: "map",
										MapKey: &ParameterType{
											Name: "int",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
										},
									},
								},
							},
							ID: 881980169,
						},
						{
							Name: "ArrayInt",
							Inputs: []*Parameter{
								{
									Name: "$in",
									Type: &ParameterType{
										Name:    "int",
										IsSlice: true,
									},
								},
							},
							ID: 3862002418,
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
									Name: "string",
								},
							},
							{
								Name: "Parent",
								Type: &ParameterType{
									Name:      "Person",
									IsStruct:  true,
									IsPointer: true,
								},
							},
							{
								Name: "Details",
								Type: &ParameterType{
									Name:     "$anon1",
									IsStruct: true,
								},
							},
						},
					},
					"$anon1": {
						Name: "$anon1",
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
									Name:     "$anon2",
									IsStruct: true,
								},
							},
						},
					},
					"$anon2": {
						Name: "$anon2",
						Fields: []*Field{
							{
								Name: "Street",
								Type: &ParameterType{
									Name: "string",
								},
							},
						},
					},
				},
			},
			wantErr: false,
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
