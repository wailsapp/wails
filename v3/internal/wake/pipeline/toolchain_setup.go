package pipeline

// CompileToolchain resolves the same compiler selection used by host planning.
// Setup uses it to explain requirements before tools have been installed. This
// does not validate installed compilers or container images, or run any tools.
func CompileToolchain(spec CompileSpec, host HostCapabilities) (string, error) {
	return resolveToolchain(spec, host)
}
