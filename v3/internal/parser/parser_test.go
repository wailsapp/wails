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
		wantModels       map[string]map[string]*StructDef
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
						},
						{
							Name: "StringArrayInputStringOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "StringArrayInputStringArrayOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "StringArrayInputNamedOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "StringArrayInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "IntPointerInputNamedOutputs",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UIntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "IntPointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int8PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int16PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int32PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int64PointerInAndOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "IntInIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int8InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int16InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int32InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Int64InIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UIntInUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt8InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt16InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt32InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "UInt64InUIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Float32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "Float64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "PointerFloat32InFloat32Out",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "PointerFloat64InFloat64Out",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "BoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "PointerBoolInBoolOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "PointerStringInStringOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "StructPointerInputErrorOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
										Package:   "main",
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
						},
						{
							Name: "StructInputStructOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Name:     "Person",
										IsStruct: true,
										Package:  "main",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:     "Person",
										IsStruct: true,
										Package:  "main",
									},
								},
							},
						},
						{
							Name: "StructPointerInputStructPointerOutput",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
										Package:   "main",
									},
								},
							},
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Person",
										IsPointer: true,
										IsStruct:  true,
										Package:   "main",
									},
								},
							},
						},
						{
							Name: "MapIntInt",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "PointerMapIntInt",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "MapIntPointerInt",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "MapIntSliceInt",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "MapIntSliceIntInMapIntSliceIntOut",
							Inputs: []*Parameter{
								{
									Name: "in",
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
						},
						{
							Name: "ArrayInt",
							Inputs: []*Parameter{
								{
									Name: "in",
									Type: &ParameterType{
										Name:    "int",
										IsSlice: true,
									},
								},
							},
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
						},
					},
					"anon1": {
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
									Name:     "anon2",
									IsStruct: true,
									Package:  "main",
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
									Name: "string",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should find multiple bound services",
			dir:  "testdata/struct_literal_multiple",
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "",
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
						},
					},
					"OtherService": {
						{
							Name: "Hello",
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
							Name:       "Greet",
							DocComment: "",
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
						},
					},
					"OtherService": {
						{
							Name: "Hello",
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
							Name:       "Greet",
							DocComment: "Greet does XYZ\n",
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
						},
						{
							Name:       "NewPerson",
							DocComment: "NewPerson creates a new person\n",
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
										Package:   "main",
									},
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services": {
					"OtherService": {
						{
							Name: "Yay",
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
		{
			name:    "should find a bound services using a variable",
			dir:     "testdata/variable_single",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "Greet someone\n",
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
						},
					},
				},
			},
		},
		{
			name:    "should find a bound services using a variable from function call",
			dir:     "testdata/variable_single_from_function",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "Greet someone\n",
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
						},
					},
				},
			},
		},
		{
			name:    "should find a bound services using a variable from function call in another package",
			dir:     "testdata/variable_single_from_other_function",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "Greet does XYZ\n",
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
						},
						{
							Name:       "NewPerson",
							DocComment: "NewPerson creates a new person\n",
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
										Package:   "main",
									},
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services": {
					"OtherService": {
						{
							Name: "Yay",
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Address",
										IsStruct:  true,
										IsPointer: true,
										Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services",
									},
								},
							},
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
								Name: "Address",
								Type: &ParameterType{
									Name:      "Address",
									IsStruct:  true,
									IsPointer: true,
									Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services",
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services": {
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
		{
			name:    "should find a bound service returned from a function call",
			dir:     "testdata/function_single",
			wantErr: false,
			wantBoundMethods: map[string]map[string][]*BoundMethod{
				"main": {
					"GreetService": {
						{
							Name:       "Greet",
							DocComment: "Greet someone\n",
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
							DocComment: "Greet does XYZ\n",
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
						},
						{
							Name:       "NewPerson",
							DocComment: "NewPerson creates a new person\n",
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
										Package:   "main",
									},
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services": {
					"OtherService": {
						{
							Name: "Yay",
							Outputs: []*Parameter{
								{
									Type: &ParameterType{
										Name:      "Address",
										IsStruct:  true,
										IsPointer: true,
										Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services",
									},
								},
							},
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
								Name: "Address",
								Type: &ParameterType{
									Name:      "Address",
									IsStruct:  true,
									IsPointer: true,
									Package:   "github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services",
								},
							},
						},
					},
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services": {
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
			if diff := cmp.Diff(tt.wantBoundMethods, got.BoundMethods); diff != "" {
				t.Errorf("ParseDirectory() failed:\n" + diff)
			}
			if diff := cmp.Diff(tt.wantModels, got.Models); diff != "" {
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
