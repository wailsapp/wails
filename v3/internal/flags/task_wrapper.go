package flags

type Build struct {
	Common
	Tags       string `name:"tags" description:"Additional build tags to pass to the Go compiler (comma-separated)"`
	Obfuscated bool   `name:"obfuscated" description:"Build with garble and stable obfuscated binding IDs"`
	GarbleArgs string `name:"garbleargs" description:"Additional arguments to pass to garble before the build command"`
}

type Dev struct {
	Common
}

type Package struct {
	Common
}

type SignWrapper struct {
	Common
}
