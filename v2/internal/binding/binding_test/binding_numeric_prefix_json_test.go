package binding_test

type NumericPrefixJSON struct {
	Name string `json:"200name"`
}

func (n NumericPrefixJSON) Get() NumericPrefixJSON {
	return n
}

var NumericPrefixJSONTest = BindingTest{
	name: "NumericPrefixJSON",
	structs: []interface{}{
		&NumericPrefixJSON{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class NumericPrefixJSON {
		"200name"?: string;
		static createFrom(source: any = {}) {
			return new NumericPrefixJSON(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this["200name] = source["200name"];
		}
	}
}
`,
}
