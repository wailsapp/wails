package flags

// Lipo represents the options for creating macOS universal binaries
type Lipo struct {
	Common

	// Output is the path for the universal binary
	Output string `name:"output" short:"o" description:"Output path for the universal binary" default:""`

	// Inputs are the architecture-specific binaries to combine
	Inputs []string `name:"input" short:"i" description:"Input binaries to combine (specify multiple times)" default:""`
}
