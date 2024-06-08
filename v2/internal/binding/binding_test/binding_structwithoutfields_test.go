package binding_test

type WithoutFields struct {
}

func (s WithoutFields) Get() WithoutFields {
	return s
}

var WithoutFieldsTest = BindingTest{
	name: "StructWithoutFields",
	structs: []interface{}{
		&WithoutFields{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
        	
        	export class WithoutFields {
        	
        	
        	    static createFrom(source: any = {}) {
        	        return new WithoutFields(source);
        	    }
        	
        	    constructor(source: any = {}) {
        	        if ('string' === typeof source) source = JSON.parse(source);
        	
        	    }
        	}
        
}`,
}
