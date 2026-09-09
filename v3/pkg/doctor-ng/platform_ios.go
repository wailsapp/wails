//go:build ios

package doctorng

// The doctor runs on development machines, not on iOS devices. These stubs
// exist so the package (and anything importing it) still compiles under
// GOOS=ios.

func collectPlatformExtras() map[string]string {
	return nil
}

func (d *Doctor) collectDependencies() error {
	return nil
}

func (d *Doctor) runDiagnostics() {
}
