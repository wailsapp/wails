package binding_test

import (
	"github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"
)

type GeneratedJsEntity struct {
	Name string `json:"name"`
}

func (s GeneratedJsEntity) Get() GeneratedJsEntity {
	return s
}

var GeneratedJsEntityTest = BindingTest{
	name: "GeneratedJsEntityTest",
	structs: []interface{}{
		&GeneratedJsEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
	
	export class MY_PREFIX_GeneratedJsEntity_MY_SUFFIX {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new MY_PREFIX_GeneratedJsEntity_MY_SUFFIX(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}

}

`,
}

type ParentEntity struct {
	Name       string      `json:"name"`
	Ref        ChildEntity `json:"ref"`
	ParentProp string      `json:"parentProp"`
}

func (p ParentEntity) Get() ParentEntity {
	return p
}

type ChildEntity struct {
	Name      string `json:"name"`
	ChildProp int    `json:"childProp"`
}

var GeneratedJsEntityWithNestedStructTest = BindingTest{
	name: "GeneratedJsEntityWithNestedStructTest",
	structs: []interface{}{
		&ParentEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
				
				export class MY_PREFIX_ChildEntity_MY_SUFFIX {
					name: string;
					childProp: number;
				
					static createFrom(source: any = {}) {
						return new MY_PREFIX_ChildEntity_MY_SUFFIX(source);
					}
				
					constructor(source: any = {}) {
						if ('string' === typeof source) source = JSON.parse(source);
						this.name = source["name"];
						this.childProp = source["childProp"];
					}
				}
				export class MY_PREFIX_ParentEntity_MY_SUFFIX {
					name: string;
					ref: MY_PREFIX_ChildEntity_MY_SUFFIX;
					parentProp: string;
				
					static createFrom(source: any = {}) {
						return new MY_PREFIX_ParentEntity_MY_SUFFIX(source);
					}
				
					constructor(source: any = {}) {
						if ('string' === typeof source) source = JSON.parse(source);
						this.name = source["name"];
						this.ref = this.convertValues(source["ref"], MY_PREFIX_ChildEntity_MY_SUFFIX);
						this.parentProp = source["parentProp"];
					}
				
					convertValues(a: any, classs: any, asMap: boolean = false): any {
						if (!a) {
							return a;
						}
						if (a.slice && a.map) {
							return (a as any[]).map(elem => this.convertValues(elem, classs));
						} else if ("object" === typeof a) {
							if (asMap) {
								for (const key of Object.keys(a)) {
									a[key] = new classs(a[key]);
								}
								return a;
							}
							return new classs(a);
						}
						return a;
					}
				}

			}
`,
}

type ParentPackageEntity struct {
	Name string             `json:"name"`
	Ref  ChildPackageEntity `json:"ref"`
}

func (p ParentPackageEntity) Get() ParentPackageEntity {
	return p
}

type ChildPackageEntity struct {
	Name            string                       `json:"name"`
	ImportedPackage binding_test_import.AWrapper `json:"importedPackage"`
}

var EntityWithDiffNamespacesTest = BindingTest{
	name: "EntityWithDiffNamespaces ",
	structs: []interface{}{
		&ParentPackageEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
            	
            	export class MY_PREFIX_ChildPackageEntity_MY_SUFFIX {
            	    name: string;
            	    importedPackage: binding_test_import.MY_PREFIX_AWrapper_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_ChildPackageEntity_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.name = source["name"];
            	        this.importedPackage = this.convertValues(source["importedPackage"], binding_test_import.MY_PREFIX_AWrapper_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice && a.map) {
            		        return (a as any[]).map(elem => this.convertValues(elem, classs));
            		    } else if ("object" === typeof a) {
            		        if (asMap) {
            		            for (const key of Object.keys(a)) {
            		                a[key] = new classs(a[key]);
            		            }
            		            return a;
            		        }
            		        return new classs(a);
            		    }
            		    return a;
            		}
            	}
            	export class MY_PREFIX_ParentPackageEntity_MY_SUFFIX {
            	    name: string;
            	    ref: MY_PREFIX_ChildPackageEntity_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_ParentPackageEntity_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.name = source["name"];
            	        this.ref = this.convertValues(source["ref"], MY_PREFIX_ChildPackageEntity_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice && a.map) {
            		        return (a as any[]).map(elem => this.convertValues(elem, classs));
            		    } else if ("object" === typeof a) {
            		        if (asMap) {
            		            for (const key of Object.keys(a)) {
            		                a[key] = new classs(a[key]);
            		            }
            		            return a;
            		        }
            		        return new classs(a);
            		    }
            		    return a;
            		}
            	}

            }

            export namespace binding_test_import {
            	
            	export class MY_PREFIX_AWrapper_MY_SUFFIX {
            	    AWrapper: binding_test_nestedimport.MY_PREFIX_A_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_AWrapper_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.AWrapper = this.convertValues(source["AWrapper"], binding_test_nestedimport.MY_PREFIX_A_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice && a.map) {
            		        return (a as any[]).map(elem => this.convertValues(elem, classs));
            		    } else if ("object" === typeof a) {
            		        if (asMap) {
            		            for (const key of Object.keys(a)) {
            		                a[key] = new classs(a[key]);
            		            }
            		            return a;
            		        }
            		        return new classs(a);
            		    }
            		    return a;
            		}
            	}

            }

            export namespace binding_test_nestedimport {
            	
            	export class MY_PREFIX_A_MY_SUFFIX {
            	    A: string;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_A_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.A = source["A"];
            	    }
            	}

            }

`,
}

type IntEnum int

const (
	IntEnumValue1 IntEnum = iota
	IntEnumValue2
	IntEnumValue3
)

var AllIntEnumValues = []struct {
	Value  IntEnum
	TSName string
}{
	{IntEnumValue1, "Value1"},
	{IntEnumValue2, "Value2"},
	{IntEnumValue3, "Value3"},
}

type EntityWithIntEnum struct {
	Name string  `json:"name"`
	Enum IntEnum `json:"enum"`
}

func (e EntityWithIntEnum) Get() EntityWithIntEnum {
	return e
}

var GeneratedJsEntityWithIntEnumTest = BindingTest{
	name: "GeneratedJsEntityWithIntEnumTest",
	structs: []interface{}{
		&EntityWithIntEnum{},
	},
	enums: []interface{}{
		AllIntEnumValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `export namespace binding_test {
        	
        	export enum MY_PREFIX_IntEnum_MY_SUFFIX {
        	    Value1 = 0,
        	    Value2 = 1,
        	    Value3 = 2,
        	}
        	export class MY_PREFIX_EntityWithIntEnum_MY_SUFFIX {
        	    name: string;
        	    enum: MY_PREFIX_IntEnum_MY_SUFFIX;
        	
        	    static createFrom(source: any = {}) {
        	        return new MY_PREFIX_EntityWithIntEnum_MY_SUFFIX(source);
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

type StringEnum string

const (
	StringEnumValue1 StringEnum = "value1"
	StringEnumValue2 StringEnum = "value2"
	StringEnumValue3 StringEnum = "value3"
)

var AllStringEnumValues = []struct {
	Value  StringEnum
	TSName string
}{
	{StringEnumValue1, "Value1"},
	{StringEnumValue2, "Value2"},
	{StringEnumValue3, "Value3"},
}

type EntityWithStringEnum struct {
	Name string     `json:"name"`
	Enum StringEnum `json:"enum"`
}

func (e EntityWithStringEnum) Get() EntityWithStringEnum {
	return e
}

var GeneratedJsEntityWithStringEnumTest = BindingTest{
	name: "GeneratedJsEntityWithStringEnumTest",
	structs: []interface{}{
		&EntityWithStringEnum{},
	},
	enums: []interface{}{
		AllStringEnumValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `export namespace binding_test {
        	
        	export enum MY_PREFIX_StringEnum_MY_SUFFIX {
        	    Value1 = "value1",
        	    Value2 = "value2",
        	    Value3 = "value3",
        	}
        	export class MY_PREFIX_EntityWithStringEnum_MY_SUFFIX {
        	    name: string;
        	    enum: MY_PREFIX_StringEnum_MY_SUFFIX;
        	
        	    static createFrom(source: any = {}) {
        	        return new MY_PREFIX_EntityWithStringEnum_MY_SUFFIX(source);
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

type EnumWithTsName string

const (
	EnumWithTsName1 EnumWithTsName = "value1"
	EnumWithTsName2 EnumWithTsName = "value2"
	EnumWithTsName3 EnumWithTsName = "value3"
)

var AllEnumWithTsNameValues = []EnumWithTsName{EnumWithTsName1, EnumWithTsName2, EnumWithTsName3}

func (v EnumWithTsName) TSName() string {
	switch v {
	case EnumWithTsName1:
		return "TsName1"
	case EnumWithTsName2:
		return "TsName2"
	case EnumWithTsName3:
		return "TsName3"
	default:
		return "???"
	}
}

type EntityWithEnumTsName struct {
	Name string         `json:"name"`
	Enum EnumWithTsName `json:"enum"`
}

func (e EntityWithEnumTsName) Get() EntityWithEnumTsName {
	return e
}

var GeneratedJsEntityWithEnumTsName = BindingTest{
	name: "GeneratedJsEntityWithEnumTsName",
	structs: []interface{}{
		&EntityWithEnumTsName{},
	},
	enums: []interface{}{
		AllEnumWithTsNameValues,
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `export namespace binding_test {
        	
        	export enum MY_PREFIX_EnumWithTsName_MY_SUFFIX {
        	    TsName1 = "value1",
        	    TsName2 = "value2",
        	    TsName3 = "value3",
        	}
        	export class MY_PREFIX_EntityWithEnumTsName_MY_SUFFIX {
        	    name: string;
        	    enum: MY_PREFIX_EnumWithTsName_MY_SUFFIX;
        	
        	    static createFrom(source: any = {}) {
        	        return new MY_PREFIX_EntityWithEnumTsName_MY_SUFFIX(source);
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

var GeneratedJsEntityWithNestedStructInterfacesTest = BindingTest{
	name: "GeneratedJsEntityWithNestedStructInterfacesTest",
	structs: []interface{}{
		&ParentEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix:     "MY_PREFIX_",
		TsSuffix:     "_MY_SUFFIX",
		TsOutputType: "interfaces",
	},
	want: `export namespace binding_test {
        	
        	export interface MY_PREFIX_ChildEntity_MY_SUFFIX {
        	    name: string;
        	    childProp: number;
        	}
        	export interface MY_PREFIX_ParentEntity_MY_SUFFIX {
        	    name: string;
        	    ref: MY_PREFIX_ChildEntity_MY_SUFFIX;
        	    parentProp: string;
        	}
        
        }
`,
}
