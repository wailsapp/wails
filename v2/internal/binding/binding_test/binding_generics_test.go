package binding_test

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import/float_package"

// Issues 3900, 3371, 2323 (no TS generics though)

type ListData[T interface{}] struct {
	Total     int64 `json:"Total"`
	TotalPage int64 `json:"TotalPage"`
	PageNum   int   `json:"PageNum"`
	List      []T   `json:"List,omitempty"`
}

func (x ListData[T]) Get() ListData[T] {
	return x
}

var Generics1Test = BindingTest{
	name: "Generics1",
	structs: []interface{}{
		&ListData[string]{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
        
                export class ListData_string_ {
                    Total: number;
                    TotalPage: number;
                    PageNum: number;
                    List?: string[];
        
                    static createFrom(source: any = {}) {
                        return new ListData_string_(source);
                    }
        
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Total = source["Total"];
                        this.TotalPage = source["TotalPage"];
                        this.PageNum = source["PageNum"];
                        this.List = source["List"];
                    }
                }
        
        }
`,
}

var Generics2Test = BindingTest{
	name: "Generics2",
	structs: []interface{}{
		&ListData[float_package.SomeStruct]{},
		&ListData[*float_package.SomeStruct]{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
        
                export class ListData__github_com_wailsapp_wails_v2_internal_binding_binding_test_binding_test_import_float_package_SomeStruct_ {
                    Total: number;
                    TotalPage: number;
                    PageNum: number;
                    List?: float_package.SomeStruct[];
        
                    static createFrom(source: any = {}) {
                        return new ListData__github_com_wailsapp_wails_v2_internal_binding_binding_test_binding_test_import_float_package_SomeStruct_(source);
                    }
        
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Total = source["Total"];
                        this.TotalPage = source["TotalPage"];
                        this.PageNum = source["PageNum"];
                        this.List = this.convertValues(source["List"], float_package.SomeStruct);
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
                export class ListData_github_com_wailsapp_wails_v2_internal_binding_binding_test_binding_test_import_float_package_SomeStruct_ {
                    Total: number;
                    TotalPage: number;
                    PageNum: number;
                    List?: float_package.SomeStruct[];
        
                    static createFrom(source: any = {}) {
                        return new ListData_github_com_wailsapp_wails_v2_internal_binding_binding_test_binding_test_import_float_package_SomeStruct_(source);
                    }
        
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Total = source["Total"];
                        this.TotalPage = source["TotalPage"];
                        this.PageNum = source["PageNum"];
                        this.List = this.convertValues(source["List"], float_package.SomeStruct);
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
        
        export namespace float_package {
        
                export class SomeStruct {
                    string: string;
        
                    static createFrom(source: any = {}) {
                        return new SomeStruct(source);
                    }
        
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.string = source["string"];
                    }
                }
        
        }
`,
}
