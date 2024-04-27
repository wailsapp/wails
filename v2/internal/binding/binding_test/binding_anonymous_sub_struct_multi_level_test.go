package binding_test

type StructWithAnonymousSubMultiLevelStruct struct {
	Name string `json:"name"`
	Meta struct {
		Age  int `json:"age"`
		More struct {
			Info       string `json:"info"`
			MoreInMore struct {
				Demo string `json:"demo"`
			} `json:"more_in_more"`
		} `json:"more"`
	} `json:"meta"`
}

func (s StructWithAnonymousSubMultiLevelStruct) Get() StructWithAnonymousSubMultiLevelStruct {
	return s
}

var AnonymousSubStructMultiLevelTest = BindingTest{
	name: "StructWithAnonymousSubMultiLevelStruct",
	structs: []interface{}{
		&StructWithAnonymousSubMultiLevelStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class StructWithAnonymousSubMultiLevelStruct {
		name: string;
		// Go type: struct { Age int "json:\"age\""; More struct { Info string "json:\"info\""; MoreInMore struct { Demo string "json:\"demo\"" } "json:\"more_in_more\"" } "json:\"more\"" }
		meta: any;
	
		static createFrom(source: any = {}) {
			return new StructWithAnonymousSubMultiLevelStruct(source);
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
