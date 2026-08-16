package flags

type Build struct {
	Common
	Tags       string `name:"tags" description:"Additional build tags to pass to the Go compiler (comma-separated)"`
	Obfuscated bool   `name:"obfuscated" description:"Build with garble and stable obfuscated binding IDs"`
	GarbleArgs string `name:"garbleargs" description:"Additional arguments to pass to garble before the build command"`
	Profile    string `name:"profile" description:"Manifest profile to apply"`
	Target     string `name:"target" description:"Comma-separated platform/architecture targets (for example darwin/arm64,linux/amd64)"`
	Force      bool   `name:"force" description:"Ignore Wake cache entries"`
}

type Dev struct {
	Common
}

type Package struct {
	Common
	Profile string `name:"profile" description:"Manifest profile to apply"`
	Target  string `name:"target" description:"Comma-separated platform/architecture targets"`
	Formats string `name:"formats" description:"Comma-separated package formats"`
	Force   bool   `name:"force" description:"Ignore Wake cache entries"`
}

type SignWrapper struct {
	Common
	Profile string `name:"profile" description:"Manifest profile to apply"`
	Target  string `name:"target" description:"Comma-separated platform/architecture targets"`
	Formats string `name:"formats" description:"Comma-separated package formats to build and sign"`
}
