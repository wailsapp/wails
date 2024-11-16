package config

// WailsAppPkgPath is the official import path of Wails v3's application package.
const WailsAppPkgPath = "github.com/wailsapp/wails/v3/pkg/application"

// SystemPaths holds resolved paths of required system packages.
type SystemPaths struct {
	ContextPackage     string
	ApplicationPackage string
}
