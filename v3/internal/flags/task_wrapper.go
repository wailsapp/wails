package flags

type Build struct {
	Common
	Tags            string `name:"tags" description:"Additional build tags to pass to Go compiler (comma-separated)"`
	RuntimeDevtools bool   `name:"runtimedevtools" description:"Enable runtime devtools API support (allows programmatic opening of devtools)"`
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
