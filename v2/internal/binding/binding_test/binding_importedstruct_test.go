package binding_test

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"

type ImportedStruct struct {
	AWrapperContainer binding_test_import.AWrapper `json:"AWrapperContainer"`
}

func (s ImportedStruct) Get() ImportedStruct {
	return s
}

var ImportedStructTest = BindingTest{
	name: "ImportedStruct",
	structs: []interface{}{
		&ImportedStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class ImportedStruct {
		AWrapperContainer: binding_test_import.AWrapper;
		static createFrom(source: any = {}) {
			return new ImportedStruct(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.AWrapperContainer = this.convertValues(source["AWrapperContainer"], binding_test_import.AWrapper);
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
	export class AWrapper {
		AWrapper: binding_test_nestedimport.A;
		static createFrom(source: any = {}) {
			return new AWrapper(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.AWrapper = this.convertValues(source["AWrapper"], binding_test_nestedimport.A);
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
