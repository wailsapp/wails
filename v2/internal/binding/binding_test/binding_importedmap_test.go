package binding_test

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"

type ImportedMap struct {
	AMapWrapperContainer binding_test_import.AMapWrapper `json:"AMapWrapperContainer"`
}

func (s ImportedMap) Get() ImportedMap {
	return s
}

var ImportedMapTest = BindingTest{
	name: "ImportedMap",
	structs: []interface{}{
		&ImportedMap{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class ImportedMap {
		AMapWrapperContainer: binding_test_import.AMapWrapper;
		static createFrom(source: any = {}) {
			return new ImportedMap(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.AMapWrapperContainer = this.convertValues(source["AMapWrapperContainer"], binding_test_import.AMapWrapper);
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

export namespace binding_test_import {
	export class AMapWrapper {
		AMap: {[key: string]: binding_test_nestedimport.A};
		static createFrom(source: any = {}) {
			return new AMapWrapper(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.AMap = this.convertValues(source["AMap"], binding_test_nestedimport.A, true);
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

export namespace binding_test_nestedimport {
	export class A {
		A: string;
		static createFrom(source: any = {}) {
			return new A(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.A = source["A"];
		}
	}
}
`,
}
