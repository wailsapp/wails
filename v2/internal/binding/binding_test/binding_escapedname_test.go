package binding_test

type EscapedName struct {
	Name string `json:"n.a.m.e"`
}

func (s EscapedName) Get() EscapedName {
	return s
}

var EscapedNameTest = BindingTest{
	name: "EscapedName",
	structs: []interface{}{
		&EscapedName{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class EscapedName {
		"n.a.m.e": string;
		static createFrom(source: any = {}) {
			return new EscapedName(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this["n.a.m.e"] = source["n.a.m.e"];
		}
	}
}
`,
}
