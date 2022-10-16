package binding_test

type Multistruct1 struct {
	Name string `json:"name"`
}

func (s Multistruct1) Get() Multistruct1 {
	return s
}

type Multistruct2 struct {
	Name string `json:"name"`
}

func (s Multistruct2) Get() Multistruct2 {
	return s
}

type Multistruct3 struct {
	Name string `json:"name"`
}

func (s Multistruct3) Get() Multistruct3 {
	return s
}

type Multistruct4 struct {
	Name string `json:"name"`
}

func (s Multistruct4) Get() Multistruct4 {
	return s
}

var MultistructTest = BindingTest{
	name: "Multistruct",
	structs: []interface{}{
		&Multistruct1{},
		&Multistruct2{},
		&Multistruct3{},
		&Multistruct4{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `export namespace binding_test {
	export class Multistruct1 {
		name: string;
		static createFrom(source: any = {}) {
			return new Multistruct1(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
	export class Multistruct2 {
		name: string;
		static createFrom(source: any = {}) {
			return new Multistruct2(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
	export class Multistruct3 {
		name: string;
		static createFrom(source: any = {}) {
			return new Multistruct3(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
	export class Multistruct4 {
		name: string;
		static createFrom(source: any = {}) {
			return new Multistruct4(source);
		}
		constructor(source: any = {}) {
			if ('string' === typeof source) source = JSON.parse(source);
			this.name = source["name"];
		}
	}
}
`,
}
