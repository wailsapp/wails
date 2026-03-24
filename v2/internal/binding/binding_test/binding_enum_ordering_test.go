package binding_test

// Test for PR #4664: Fix generated enums ordering
// This test ensures that enum output is deterministic regardless of map iteration order

// ZFirstEnum - named with Z prefix to test alphabetical sorting
type ZFirstEnum int

const (
	ZFirstEnumValue1 ZFirstEnum = iota
	ZFirstEnumValue2
)

var AllZFirstEnumValues = []struct {
	Value  ZFirstEnum
	TSName string
}{
	{ZFirstEnumValue1, "ZValue1"},
	{ZFirstEnumValue2, "ZValue2"},
}

// ASecondEnum - named with A prefix to test alphabetical sorting
type ASecondEnum int

const (
	ASecondEnumValue1 ASecondEnum = iota
	ASecondEnumValue2
)

var AllASecondEnumValues = []struct {
	Value  ASecondEnum
	TSName string
}{
	{ASecondEnumValue1, "AValue1"},
	{ASecondEnumValue2, "AValue2"},
}

// MMiddleEnum - named with M prefix to test alphabetical sorting
type MMiddleEnum int

const (
	MMiddleEnumValue1 MMiddleEnum = iota
	MMiddleEnumValue2
)

var AllMMiddleEnumValues = []struct {
	Value  MMiddleEnum
	TSName string
}{
	{MMiddleEnumValue1, "MValue1"},
	{MMiddleEnumValue2, "MValue2"},
}

type EntityWithMultipleEnums struct {
	Name   string      `json:"name"`
	EnumZ  ZFirstEnum  `json:"enumZ"`
	EnumA  ASecondEnum `json:"enumA"`
	EnumM  MMiddleEnum `json:"enumM"`
}

func (e EntityWithMultipleEnums) Get() EntityWithMultipleEnums {
	return e
}

// EnumOrderingTest tests that multiple enums in the same package are output
// in alphabetical order by enum name. Before PR #4664, the order was
// non-deterministic due to Go map iteration order.
var EnumOrderingTest = BindingTest{
	name: "EnumOrderingTest",
	structs: []interface{}{
		&EntityWithMultipleEnums{},
	},
	enums: []interface{}{
		// Intentionally add enums in non-alphabetical order
		AllZFirstEnumValues,
		AllASecondEnumValues,
		AllMMiddleEnumValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "",
		TsSuffix: "",
	},
	// Expected output should have enums in alphabetical order: ASecondEnum, MMiddleEnum, ZFirstEnum
	want: `export namespace binding_test {

	export enum ASecondEnum {
	    AValue1 = 0,
	    AValue2 = 1,
	}
	export enum MMiddleEnum {
	    MValue1 = 0,
	    MValue2 = 1,
	}
	export enum ZFirstEnum {
	    ZValue1 = 0,
	    ZValue2 = 1,
	}
	export class EntityWithMultipleEnums {
	    name: string;
	    enumZ: ZFirstEnum;
	    enumA: ASecondEnum;
	    enumM: MMiddleEnum;

	    static createFrom(source: any = {}) {
	        return new EntityWithMultipleEnums(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.enumZ = source["enumZ"];
	        this.enumA = source["enumA"];
	        this.enumM = source["enumM"];
	    }
	}

}
`,
}

// EnumElementOrderingEnum tests sorting of enum elements by TSName
type EnumElementOrderingEnum string

const (
	EnumElementZ EnumElementOrderingEnum = "z_value"
	EnumElementA EnumElementOrderingEnum = "a_value"
	EnumElementM EnumElementOrderingEnum = "m_value"
)

// AllEnumElementOrderingValues intentionally lists values out of alphabetical order
// to test that AddEnum sorts them
var AllEnumElementOrderingValues = []struct {
	Value  EnumElementOrderingEnum
	TSName string
}{
	{EnumElementZ, "Zebra"},
	{EnumElementA, "Apple"},
	{EnumElementM, "Mango"},
}

type EntityWithUnorderedEnumElements struct {
	Name string                  `json:"name"`
	Enum EnumElementOrderingEnum `json:"enum"`
}

func (e EntityWithUnorderedEnumElements) Get() EntityWithUnorderedEnumElements {
	return e
}

// EnumElementOrderingTest tests that enum elements are sorted alphabetically
// by their TSName within an enum. Before PR #4664, elements appeared in the
// order they were added, which could be arbitrary.
var EnumElementOrderingTest = BindingTest{
	name: "EnumElementOrderingTest",
	structs: []interface{}{
		&EntityWithUnorderedEnumElements{},
	},
	enums: []interface{}{
		AllEnumElementOrderingValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "",
		TsSuffix: "",
	},
	// Expected output should have enum elements sorted: Apple, Mango, Zebra
	want: `export namespace binding_test {

	export enum EnumElementOrderingEnum {
	    Apple = "a_value",
	    Mango = "m_value",
	    Zebra = "z_value",
	}
	export class EntityWithUnorderedEnumElements {
	    name: string;
	    enum: EnumElementOrderingEnum;

	    static createFrom(source: any = {}) {
	        return new EntityWithUnorderedEnumElements(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.enum = source["enum"];
	    }
	}

}
`,
}

// TSNameEnumElementOrdering tests sorting with TSName() method enum
type TSNameEnumElementOrdering string

const (
	TSNameEnumZ TSNameEnumElementOrdering = "z_value"
	TSNameEnumA TSNameEnumElementOrdering = "a_value"
	TSNameEnumM TSNameEnumElementOrdering = "m_value"
)

func (v TSNameEnumElementOrdering) TSName() string {
	switch v {
	case TSNameEnumZ:
		return "Zebra"
	case TSNameEnumA:
		return "Apple"
	case TSNameEnumM:
		return "Mango"
	default:
		return "Unknown"
	}
}

// AllTSNameEnumValues intentionally out of order
var AllTSNameEnumValues = []TSNameEnumElementOrdering{TSNameEnumZ, TSNameEnumA, TSNameEnumM}

type EntityWithTSNameEnumOrdering struct {
	Name string                    `json:"name"`
	Enum TSNameEnumElementOrdering `json:"enum"`
}

func (e EntityWithTSNameEnumOrdering) Get() EntityWithTSNameEnumOrdering {
	return e
}

// TSNameEnumElementOrderingTest tests that enums using TSName() method
// also have their elements sorted alphabetically by the TSName.
var TSNameEnumElementOrderingTest = BindingTest{
	name: "TSNameEnumElementOrderingTest",
	structs: []interface{}{
		&EntityWithTSNameEnumOrdering{},
	},
	enums: []interface{}{
		AllTSNameEnumValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "",
		TsSuffix: "",
	},
	// Expected output should have enum elements sorted: Apple, Mango, Zebra
	want: `export namespace binding_test {

	export enum TSNameEnumElementOrdering {
	    Apple = "a_value",
	    Mango = "m_value",
	    Zebra = "z_value",
	}
	export class EntityWithTSNameEnumOrdering {
	    name: string;
	    enum: TSNameEnumElementOrdering;

	    static createFrom(source: any = {}) {
	        return new EntityWithTSNameEnumOrdering(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.enum = source["enum"];
	    }
	}

}
`,
}
