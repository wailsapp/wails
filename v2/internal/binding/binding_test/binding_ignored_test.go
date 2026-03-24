package binding_test

import (
	"unsafe"
)

// Issues 3755, 3809

type Ignored struct {
	Valid      bool
	Total      func() int `json:"Total"`
	UnsafeP    unsafe.Pointer
	Complex64  complex64 `json:"Complex"`
	Complex128 complex128
	StringChan chan string
}

func (x Ignored) Get() Ignored {
	return x
}

var IgnoredTest = BindingTest{
	name: "Ignored",
	structs: []interface{}{
		&Ignored{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
            
                export class Ignored {
                    Valid: boolean;
            
                    static createFrom(source: any = {}) {
                        return new Ignored(source);
                    }
            
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Valid = source["Valid"];
                    }
                }
            
            }
`,
}
