package frameworks

// Framework has details about a specific framework
type Framework struct {
	Name    string
	JS      string
	CSS     string
	Options string
}

// FrameworkToUse is the framework we will use when building
// Set by `wails init`, used by `wails build`
var FrameworkToUse *Framework
