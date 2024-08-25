package binding_test

type NoFieldTags struct {
	Name    string
	Address string
	Zip     *string
	Spouse  *NoFieldTags
}

func (n NoFieldTags) Get() NoFieldTags {
	return n
}

var NoFieldTagsTest = BindingTest{
	name: "NoFieldTags",
	structs: []interface{}{
		&NoFieldTags{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class NoFieldTags {
		Name: string;
		Address: string;
		Zip?: string;
		Spouse?: NoFieldTags;
		static createFrom(source: any = {}) {
			return new NoFieldTags(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.Name = source["Name"];
			this.Address = source["Address"];
			this.Zip = source["Zip"];
			this.Spouse = this.convertValues(source["Spouse"], NoFieldTags);
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
