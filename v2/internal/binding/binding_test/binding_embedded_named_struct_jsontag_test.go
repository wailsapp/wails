package binding_test

type EmbeddedWithJSONTagTimeLimitDef struct {
	IsInfinite bool  `json:"isInfinite"`
	StartTime  int64 `json:"startTime"`
	EndTime    int64 `json:"endTime"`
}

type EmbeddedWithJSONTagApplication struct {
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
	AuthUser   string `json:"authUser"`
	AuthModule string `json:"authModule"`
}

type StructWithEmbeddedNamedStructJSONTag struct {
	Name                            string `json:"name"`
	Typ                             string `json:"typ"`
	Desc                            string `json:"desc"`
	EmbeddedWithJSONTagTimeLimitDef `json:"timeLimitDef"`
	EmbeddedWithJSONTagApplication  `json:"application"`
}

func (s StructWithEmbeddedNamedStructJSONTag) Get() StructWithEmbeddedNamedStructJSONTag {
	return s
}

var EmbeddedNamedStructJSONTagTest = BindingTest{
	name: "EmbeddedNamedStructJSONTag",
	structs: []interface{}{
		&StructWithEmbeddedNamedStructJSONTag{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class EmbeddedWithJSONTagApplication {
		appName: string;
		appVersion: string;
		authUser: string;
		authModule: string;
		static createFrom(source: any = {}) {
			return new EmbeddedWithJSONTagApplication(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.appName = source["appName"];
			this.appVersion = source["appVersion"];
			this.authUser = source["authUser"];
			this.authModule = source["authModule"];
		}
	}
	export class EmbeddedWithJSONTagTimeLimitDef {
		isInfinite: boolean;
		startTime: number;
		endTime: number;
		static createFrom(source: any = {}) {
			return new EmbeddedWithJSONTagTimeLimitDef(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.isInfinite = source["isInfinite"];
			this.startTime = source["startTime"];
			this.endTime = source["endTime"];
		}
	}
	export class StructWithEmbeddedNamedStructJSONTag {
		name: string;
		typ: string;
		desc: string;
		timeLimitDef: EmbeddedWithJSONTagTimeLimitDef;
		application: EmbeddedWithJSONTagApplication;
		static createFrom(source: any = {}) {
			return new StructWithEmbeddedNamedStructJSONTag(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
			this.typ = source["typ"];
			this.desc = source["desc"];
			this.timeLimitDef = this.convertValues(source["timeLimitDef"], EmbeddedWithJSONTagTimeLimitDef);
			this.application = this.convertValues(source["application"], EmbeddedWithJSONTagApplication);
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
