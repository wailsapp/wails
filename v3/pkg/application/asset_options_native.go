//go:build wails_native

package application

// AssetOptions remains in Options so shared option construction can compile
// in both modes. Native builds have no HTTP asset server.
type AssetOptions struct {
	DisableLogging bool
}
