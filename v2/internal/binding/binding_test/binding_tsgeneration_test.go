package binding_test

import (
	"github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import"
)

type GeneratedJsEntity struct {
	Name string `json:"name"`
}

func (s GeneratedJsEntity) Get() GeneratedJsEntity {
	return s
}

var GeneratedJsEntityTest = BindingTest{
	name: "GeneratedJsEntityTest",
	structs: []interface{}{
		&GeneratedJsEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
	
	export class MY_PREFIX_GeneratedJsEntity_MY_SUFFIX {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new MY_PREFIX_GeneratedJsEntity_MY_SUFFIX(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}

}

`,
}

type ParentEntity struct {
	Name       string      `json:"name"`
	Ref        ChildEntity `json:"ref"`
	ParentProp string      `json:"parentProp"`
}

func (p ParentEntity) Get() ParentEntity {
	return p
}

type ChildEntity struct {
	Name      string `json:"name"`
	ChildProp int    `json:"childProp"`
}

var GeneratedJsEntityWithNestedStructTest = BindingTest{
	name: "GeneratedJsEntityWithNestedStructTest",
	structs: []interface{}{
		&ParentEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
				
				export class MY_PREFIX_ChildEntity_MY_SUFFIX {
					name: string;
					childProp: number;
				
					static createFrom(source: any = {}) {
						return new MY_PREFIX_ChildEntity_MY_SUFFIX(source);
					}
				
					constructor(source: any = {}) {
						if ('string' === typeof source) source = JSON.parse(source);
						this.name = source["name"];
						this.childProp = source["childProp"];
					}
				}
				export class MY_PREFIX_ParentEntity_MY_SUFFIX {
					name: string;
					ref: MY_PREFIX_ChildEntity_MY_SUFFIX;
					parentProp: string;
				
					static createFrom(source: any = {}) {
						return new MY_PREFIX_ParentEntity_MY_SUFFIX(source);
					}
				
					constructor(source: any = {}) {
						if ('string' === typeof source) source = JSON.parse(source);
						this.name = source["name"];
						this.ref = this.convertValues(source["ref"], MY_PREFIX_ChildEntity_MY_SUFFIX);
						this.parentProp = source["parentProp"];
					}
				
					convertValues(a: any, classs: any, asMap: boolean = false): any {
						if (!a) {
							return a;
						}
						if (a.slice) {
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

type ParentPackageEntity struct {
	Name string             `json:"name"`
	Ref  ChildPackageEntity `json:"ref"`
}

func (p ParentPackageEntity) Get() ParentPackageEntity {
	return p
}

type ChildPackageEntity struct {
	Name            string                       `json:"name"`
	ImportedPackage binding_test_import.AWrapper `json:"importedPackage"`
}

var EntityWithDiffNamespaces = BindingTest{
	name: "EntityWithDiffNamespaces ",
	structs: []interface{}{
		&ParentPackageEntity{},
	},
	exemptions:  nil,
	shouldError: false,
	TsGenerationOptionsTest: TsGenerationOptionsTest{
		TsPrefix: "MY_PREFIX_",
		TsSuffix: "_MY_SUFFIX",
	},
	want: `
export namespace binding_test {
            	
            	export class MY_PREFIX_ChildPackageEntity_MY_SUFFIX {
            	    name: string;
            	    importedPackage: binding_test_import.MY_PREFIX_AWrapper_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_ChildPackageEntity_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.name = source["name"];
            	        this.importedPackage = this.convertValues(source["importedPackage"], binding_test_import.MY_PREFIX_AWrapper_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice) {
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
            	export class MY_PREFIX_ParentPackageEntity_MY_SUFFIX {
            	    name: string;
            	    ref: MY_PREFIX_ChildPackageEntity_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_ParentPackageEntity_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.name = source["name"];
            	        this.ref = this.convertValues(source["ref"], MY_PREFIX_ChildPackageEntity_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice) {
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
            	
            	export class MY_PREFIX_AWrapper_MY_SUFFIX {
            	    AWrapper: binding_test_nestedimport.MY_PREFIX_A_MY_SUFFIX;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_AWrapper_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.AWrapper = this.convertValues(source["AWrapper"], binding_test_nestedimport.MY_PREFIX_A_MY_SUFFIX);
            	    }
            	
            		convertValues(a: any, classs: any, asMap: boolean = false): any {
            		    if (!a) {
            		        return a;
            		    }
            		    if (a.slice) {
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
            	
            	export class MY_PREFIX_A_MY_SUFFIX {
            	    A: string;
            	
            	    static createFrom(source: any = {}) {
            	        return new MY_PREFIX_A_MY_SUFFIX(source);
            	    }
            	
            	    constructor(source: any = {}) {
            	        if ('string' === typeof source) source = JSON.parse(source);
            	        this.A = source["A"];
            	    }
            	}

            }

`,
}
