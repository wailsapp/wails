package binding_test

type As struct {
	B Bs `json:"b"`
}

type Bs struct {
	Name string `json:"name"`
}

func (a As) Get() As {
	return a
}

var NestedFieldTest = BindingTest{
	name: "NestedField",
	structs: []interface{}{
		&As{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class Bs {
		name: string;
		static createFrom(source: any = {}) {
			return new Bs(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
	export class As {
		b: Bs;
		static createFrom(source: any = {}) {
			return new As(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.b = this.convertValues(source["b"], Bs);
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
