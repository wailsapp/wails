package binding_test

type GeneratedJsEntity struct {
	Name string `json:"name"`
}

func (s GeneratedJsEntity) Get() GeneratedJsEntity {
	return s
}

var GeneratedJsEntityTest = BindingTest{
	name: "GeneratedJsEntityTest ",
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
