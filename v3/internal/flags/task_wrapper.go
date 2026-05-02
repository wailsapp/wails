package flags

type Build struct {
	Common
	Tags string `name:"tags" description:"Additional build tags to pass to the Go compiler (comma-separated)"`
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
