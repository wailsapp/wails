package binding_test

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"

type ImportedEnumStruct struct {
	EnumValue binding_test_import.ImportedEnum `json:"EnumValue"`
}

func (s ImportedEnumStruct) Get() ImportedEnumStruct {
	return s
}

var ImportedEnumTest = BindingTest{
	name: "ImportedEnum",
	structs: []interface{}{
		&ImportedEnumStruct{},
	},
	enums: []interface{}{
		binding_test_import.AllImportedEnumValues,
	},
	exemptions:  nil,
	shouldError: false,
	want: `export namespace binding_test {
        	
        	export class ImportedEnumStruct {
        	    EnumValue: binding_test_import.ImportedEnum;
        	
        	    static createFrom(source: any = {}) {
        	        return new ImportedEnumStruct(source);
        	    }
        	
        	    constructor(source: any = {}) {
        	        if ('string' === typeof source) source = JSON.parse(source);
        	        this.EnumValue = source["EnumValue"];
        	    }
        	}
        
        }
        
        export namespace binding_test_import {
        	
        	export enum ImportedEnum {
        	    Value1 = "value1",
        	    Value2 = "value2",
        	    Value3 = "value3",
        	}
        
        }
`,
}
