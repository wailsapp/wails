package binding_test

type InnerStruct struct {
	IsInfinite bool  `json:"isInfinite"`
	StartTime  int64 `json:"startTime"`
	EndTime    int64 `json:"endTime"`
}

type AppInfo struct {
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
}

type StructWithNamedEmbeddedStruct struct {
	Name        string `json:"name"`
	Typ         string `json:"typ"`
	Desc        string `json:"desc"`
	InnerStruct `json:"timeLimitDef"`
	AppInfo     `json:"application"`
}

func (s StructWithNamedEmbeddedStruct) Get() StructWithNamedEmbeddedStruct {
	return s
}

var NamedEmbeddedStructTest = BindingTest{
	name: "StructWithNamedEmbeddedStruct",
	structs: []interface{}{
		&StructWithNamedEmbeddedStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class AppInfo {
		appName: string;
		appVersion: string;
	
		static createFrom(source: any = {}) {
			return new AppInfo(source);
		}
	
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.appName = source["appName"];
			this.appVersion = source["appVersion"];
		}
	}
	export class InnerStruct {
		isInfinite: boolean;
		startTime: number;
		endTime: number;
	
		static createFrom(source: any = {}) {
			return new InnerStruct(source);
		}
	
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.isInfinite = source["isInfinite"];
			this.startTime = source["startTime"];
			this.endTime = source["endTime"];
		}
	}
	export class StructWithNamedEmbeddedStruct {
		name: string;
		typ: string;
		desc: string;
		timeLimitDef: InnerStruct;
		application: AppInfo;
	
		static createFrom(source: any = {}) {
			return new StructWithNamedEmbeddedStruct(source);
		}
	
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
			this.typ = source["typ"];
			this.desc = source["desc"];
			this.timeLimitDef = this.convertValues(source["timeLimitDef"], InnerStruct);
			this.application = this.convertValues(source["application"], AppInfo);
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
