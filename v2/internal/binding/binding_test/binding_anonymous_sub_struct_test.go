package binding_test

type StructWithAnonymousSubStruct struct {
	Name string `json:"name"`
	Meta struct {
		Age int `json:"age"`
	} `json:"meta"`
}

func (s StructWithAnonymousSubStruct) Get() StructWithAnonymousSubStruct {
	return s
}

var AnonymousSubStructTest = BindingTest{
	name: "StructWithAnonymousSubStruct",
	structs: []interface{}{
		&StructWithAnonymousSubStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class StructWithAnonymousSubStruct {
		name: string;
		// Go type: struct { Age int "json:\"age\"" }
		meta: any;
	
		static createFrom(source: any = {}) {
			return new StructWithAnonymousSubStruct(source);
		}
	
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
			this.meta = this.convertValues(source["meta"], Object);
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
