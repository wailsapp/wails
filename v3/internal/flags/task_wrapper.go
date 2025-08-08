package flags

type Build struct {
	Common
	Static   bool   `name:"static" description:"Enable static linking using musl-gcc (Linux only)"`
	Compiler string `name:"cc" description:"C compiler to use for compilation"`
}

type Dev struct {
	Common
}

type Package struct {
	Common
	Static   bool   `name:"static" description:"Enable static linking using musl-gcc (Linux only)"`
	Compiler string `name:"cc" description:"C compiler to use for compilation"`
}
