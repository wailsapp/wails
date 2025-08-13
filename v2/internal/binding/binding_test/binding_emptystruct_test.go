package binding_test

type EmptyStruct struct {
	Empty struct{} `json:"empty"`
}

func (s EmptyStruct) Get() EmptyStruct {
	return s
}

var EmptyStructTest = BindingTest{
	name: "EmptyStruct",
	structs: []interface{}{
		&EmptyStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class EmptyStruct {
		// Go type: struct {}

		empty: any;

		static createFrom(source: any = {}) {
			return new EmptyStruct(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.empty = this.convertValues(source["empty"], Object);
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
