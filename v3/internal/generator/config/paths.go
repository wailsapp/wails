package config

// WailsAppPkgPath is the official import path of Wails v3's application package.
const WailsAppPkgPath = "github.com/wailsapp/wails/v3/pkg/application"

// WailsInternalPkgPath is the official import path of Wails v3's internal package.
const WailsInternalPkgPath = "github.com/wailsapp/wails/v3/internal"

// SystemPaths holds resolved paths of required system packages.
type SystemPaths struct {
	ContextPackage     string
	ApplicationPackage string
	InternalPackage    string
}
