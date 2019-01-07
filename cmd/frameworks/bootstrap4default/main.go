// +build frameworkbootstrap4

package bootstrap4

import (
	"github.com/gobuffalo/packr"
	"github.com/wailsapp/wails/frameworks"
)

func init() {
	assets := packr.NewBox("./assets")
	frameworks.FrameworkToUse = &frameworks.Framework{
		Name: "Bootstrap 4",
		JS:   assets.String("bootstrap.bundle.min.js"),
		CSS:  assets.String("bootstrap.min.css"),
	}
}
