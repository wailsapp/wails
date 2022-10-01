package binding_test

type NonStringMapKey struct {
	NumberMap map[uint]any `json:"numberMap"`
}

func (s NonStringMapKey) Get() NonStringMapKey {
	return s
}

var NonStringMapKeyTest = BindingTest{
	name: "NonStringMapKey",
	structs: []interface{}{
		&NonStringMapKey{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class NonStringMapKey {
		numberMap: {[key: number]: any};
		static createFrom(source: any = {}) {
			return new NonStringMapKey(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.numberMap = source["numberMap"];
		}
	}
}
`,
}
