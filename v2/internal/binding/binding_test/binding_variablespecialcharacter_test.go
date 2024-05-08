package binding_test

type SpecialCharacterField struct {
	ID string `json:"@ID,omitempty"`
}

func (s SpecialCharacterField) Get() SpecialCharacterField {
	return s
}

var SpecialCharacterFieldTest = BindingTest{
	name: "SpecialCharacterField",
	structs: []interface{}{
		&SpecialCharacterField{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
		export class SpecialCharacterField {
			"@ID"?: string;
        	static createFrom(source: any = {}) {
			return new SpecialCharacterField(source);
		}
		constructor(source: any = {}) {
        	if ('string' === typeof source) source = JSON.parse(source);
			this["@ID"] = source["@ID"];
		}
	} 
}
`,
}
