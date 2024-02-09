package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
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
						{
							Package:    "main",
							Name:       "NoInputsStringOut",
							DocComment: "",
							Inputs:     nil,
							Outputs: []*Parameter{
								{
									Name: "",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							ID: 1075577233,
						},
						{
							Package: "main",
							Name:    "StringArrayInputStringOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
									},
								},
							},
							ID: 1091960237,
						},
						{
							Package: "main",
							Name:    "StringArrayInputStringArrayOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							ID: 383995060,
						},
						{
							Package: "main",
							Name:    "StringArrayInputNamedOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							ID: 3678582682,
						},
						{
							Package: "main",
							Name:    "StringArrayInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Package: "main",
										Name:    "string",
										IsSlice: true,
									},
								},
								{
									Name: "err",
									Type: &ParameterType{
										Package: "main",
										Name:    "error",
									},
								},
							},
							ID: 319259595,
						},
						{
							Package: "main",
							Name:    "IntPointerInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "output",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int",
										IsPointer: true,
									},
								},
								{
									Name: "err",
									Type: &ParameterType{
										Package: "main",
										Name:    "error",
									}},
							},
							ID: 2718999663,
						},
						{
							Package: "main",
							Name:    "UIntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint",
										IsPointer: true,
									},
								},
							},
							ID: 1367187362,
						},
						{
							Package: "main",
							Name:    "UInt8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint8",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint8",
										IsPointer: true,
									},
								},
							},
							ID: 518250834,
						},
						{
							Package: "main",
							Name:    "UInt16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint16",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint16",
										IsPointer: true,
									},
								},
							},
							ID: 1236957573,
						},
						{
							Package: "main",
							Name:    "UInt32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint32",
										IsPointer: true,
									},
								},
							},
							ID: 1739300671,
						},
						{
							Package: "main",
							Name:    "UInt64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "uint64",
										IsPointer: true,
									},
								},
							},
							ID: 1403757716,
						},
						{
							Package: "main",
							Name:    "IntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "int",
										IsPointer: true,
									},
								},
							},
							ID: 1066151743,
						},
						{
							Package: "main",
							Name:    "Int8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int8",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "int8",
										IsPointer: true,
									},
								},
							},
							ID: 2189402897,
						},
						{
							Package: "main",
							Name:    "Int16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int16",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "int16",
										IsPointer: true,
									},
								},
							},
							ID: 1754277916,
						},
						{
							Package: "main",
							Name:    "Int32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "int32",
										IsPointer: true,
									},
								},
							},
							ID: 4251088558,
						},
						{
							Package: "main",
							Name:    "Int64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "int64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "int64",
										IsPointer: true,
									},
								},
							},
							ID: 2205561041,
						},
						{
							Package: "main",
							Name:    "IntInIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "int",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "int",
									},
								},
							},
							ID: 642881729,
						},
						{
							Package: "main",
							Name:    "Int8InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "int8",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "int8",
									},
								},
							},
							ID: 572240879,
						},
						{
							Package: "main",
							Name:    "Int16InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "int16",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "int16",
									},
								},
							},
							ID: 3306292566,
						},
						{
							Package: "main",
							Name:    "Int32InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "int32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "int32",
									},
								},
							},
							ID: 1909469092,
						},
						{
							Package: "main",
							Name:    "Int64InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "int64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "int64",
									},
								},
							},
							ID: 1343888303,
						},
						{
							Package: "main",
							Name:    "UIntInUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "uint",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "uint",
									},
								},
							},
							ID: 2836661285,
						},
						{
							Package: "main",
							Name:    "UInt8InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "uint8",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "uint8",
									},
								},
							},
							ID: 2988345717,
						},
						{
							Package: "main",
							Name:    "UInt16InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "uint16",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "uint16",
									},
								},
							},
							ID: 3401034892,
						},
						{
							Package: "main",
							Name:    "UInt32InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "uint32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "uint32",
									},
								},
							},
							ID: 1160383782,
						},
						{
							Package: "main",
							Name:    "UInt64InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "uint64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "uint64",
									},
								},
							},
							ID: 793803239,
						},
						{
							Package: "main",
							Name:    "Float32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "float32",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "float32",
									},
								},
							},
							ID: 3132595881,
						},
						{
							Package: "main",
							Name:    "Float64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "float64",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "float64",
									},
								},
							},
							ID: 2182412247,
						},
						{
							Package: "main",
							Name:    "PointerFloat32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "float32",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "float32",
										IsPointer: true,
									},
								},
							},
							ID: 224675106,
						},
						{
							Package: "main",
							Name:    "PointerFloat64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "float64",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "float64",
										IsPointer: true,
									},
								},
							},
							ID: 2124953624,
						},
						{
							Package: "main",
							Name:    "BoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "bool",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "bool",
									},
								},
							},
							ID: 2424639793,
						},
						{
							Package: "main",
							Name:    "PointerBoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "bool",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "bool",
										IsPointer: true,
									},
								},
							},
							ID: 3589606958,
						},
						{
							Package: "main",
							Name:    "PointerStringInStringOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "string",
										IsPointer: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "string",
										IsPointer: true,
									},
								},
							},
							ID: 229603958,
						},
						{
							Package: "main",
							Name:    "StructPointerInputErrorOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package: "main",
										Name:    "error",
									},
								},
							},
							ID: 2447692557,
						},
						{
							Package: "main",
							Name:    "StructInputStructOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:  "main",
										Name:     "Person",
										IsStruct: true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:  "main",
										Name:     "Person",
										IsStruct: true,
									},
								},
							},
							ID: 3835643147,
						},
						{
							Package: "main",
							Name:    "StructPointerInputStructPointerOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Package:   "main",
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
									},
								},
							},
							ID: 2943477349,
						},
						{
							Package: "main",
							Name:    "MapIntInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "map",
										MapKey: &ParameterType{
											Name:    "int",
											Package: "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											Package: "main",
										},
									},
								},
							},
							ID: 2386486356,
						},
						{
							Package: "main",
							Name:    "PointerMapIntInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package:   "main",
										Name:      "map",
										IsPointer: true,
										MapKey: &ParameterType{
											Name:    "int",
											Package: "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											Package: "main",
										},
									},
								},
							},
							ID: 3516977899,
						},
						{
							Package: "main",
							Name:    "MapIntPointerInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "map",
										MapKey: &ParameterType{
											Name:      "int",
											IsPointer: true,
											Package:   "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											Package: "main",
										},
									},
								},
							},
							ID: 550413585,
						},
						{
							Package: "main",
							Name:    "MapIntSliceInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "map",
										MapKey: &ParameterType{
											Name:    "int",
											Package: "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
											Package: "main",
										},
									},
								},
							},
							ID: 2900172572,
						},
						{
							Package: "main",
							Name:    "MapIntSliceIntInMapIntSliceIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
										Name:    "map",
										MapKey: &ParameterType{
											Name:    "int",
											Package: "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
											Package: "main",
										},
									},
								},
							},
							Outputs: []*Parameter{
								{
									Name: "out",
									Type: &ParameterType{
										Package: "main",
										Name:    "map",
										MapKey: &ParameterType{
											Name:    "int",
											Package: "main",
										},
										MapValue: &ParameterType{
											Name:    "int",
											IsSlice: true,
											Package: "main",
										},
									},
								},
							},
							ID: 881980169,
						},
						{
							Package: "main",
							Name:    "ArrayInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Package: "main",
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
									Package: "main",
									Name:    "string",
								},
							},
							{
								Name: "Parent",
								Type: &ParameterType{
									Package:   "main",
									Name:      "Person",
									IsStruct:  true,
									IsPointer: true,
								},
							},
							{
								Name: "Details",
								Type: &ParameterType{
									Package:  "main",
									Name:     "anon1",
									IsStruct: true,
								},
							},
						},
					},
					"anon1": {
						Name: "anon1",
						Fields: []*Field{
							{
								Name: "Age",
								Type: &ParameterType{
									Package: "main",
									Name:    "int",
								},
							},
							{
								Name: "Address",
								Type: &ParameterType{
									Package:  "main",
									Name:     "anon2",
									IsStruct: true,
								},
							},
						},
					},
					"anon2": {
						Name: "anon2",
						Fields: []*Field{
							{
								Name: "Street",
								Type: &ParameterType{
									Package: "main",
									Name:    "string",
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
					if diff := cmp.Diff(wantBoundMethods, gotBoundMethods, cmp.AllowUnexported(Parameter{})); diff != "" {
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
