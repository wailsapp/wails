package binding_test

type SingleField struct {
	Name string `json:"name"`
}

func (s SingleField) Get() SingleField {
	return s
}

var SingleFieldTest = BindingTest{
	name: "SingleField",
	structs: []interface{}{
		&SingleField{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class SingleField {
		name: string;
		static createFrom(source: any = {}) {
			return new SingleField(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
}
`,
}
