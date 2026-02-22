package flags

type BuildCommon struct {
	LdFlags      string `description:"Additional ldflags to pass to the compiler"`
	Compiler     string `description:"Use a different go compiler to build, eg go1.15beta1"`
	SkipBindings bool   `description:"Skips generation of bindings"`
	RaceDetector bool   `name:"race" description:"Build with Go's race detector"`
	SkipFrontend bool   `name:"s" description:"Skips building the frontend"`
	Verbosity    int    `name:"v" description:"Verbosity level (0 = quiet, 1 = normal, 2 = verbose)"`
	Tags         string `description:"Build tags to pass to Go compiler. Must be quoted. Space or comma (but not both) separated"`
	NoSyncGoMod  bool   `description:"Don't sync go.mod"`
	SkipModTidy  bool   `name:"m" description:"Skip mod tidy before compile"`
}

func (c BuildCommon) Default() BuildCommon {
	return BuildCommon{
		Compiler:  "go",
		Verbosity: 1,
	}
}
