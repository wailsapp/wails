package binding_test

import (
	"time"

	"github.com/google/uuid"
)

type CustomTypeStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s CustomTypeStruct) Get() CustomTypeStruct {
	return s
}

var CustomTypeGenerationTest = BindingTest{
	name: "CustomTypeGenerationTest",
	structs: []interface{}{
		&CustomTypeStruct{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
	export class CustomTypeStruct {
		id: string;
		createdAt: string;
		static createFrom(source: any = {}) {
			return new CustomTypeStruct(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.id = source["id"];
			this.createdAt = source["createdAt"];
		}
	}
}
`,
}
