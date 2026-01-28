package binding_test

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"

type ImportedSlice struct {
	ASliceWrapperContainer binding_test_import.ASliceWrapper `json:"ASliceWrapperContainer"`
}

func (s ImportedSlice) Get() ImportedSlice {
	return s
}

var ImportedSliceTest = BindingTest{
	name: "ImportedSlice",
	structs: []interface{}{
		&ImportedSlice{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class ImportedSlice {
		ASliceWrapperContainer: binding_test_import.ASliceWrapper;
		static createFrom(source: any = {}) {
			return new ImportedSlice(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.ASliceWrapperContainer = this.convertValues(source["ASliceWrapperContainer"], binding_test_import.ASliceWrapper);
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
	export class ASliceWrapper {
		ASlice: binding_test_nestedimport.A[];
		static createFrom(source: any = {}) {
			return new ASliceWrapper(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.ASlice = this.convertValues(source["ASlice"], binding_test_nestedimport.A);
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
