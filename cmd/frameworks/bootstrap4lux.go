// +build frameworkbootstrap4lux

package frameworks

import (
	"github.com/gobuffalo/packr"
)

func init() {
	assets := packr.NewBox("./bootstrap4lux/assets")
	FrameworkToUse = &Framework{
		Name: "Bootstrap 4 (Lux)",
		JS:   BoxString(&assets, "bootstrap.bundle.min.js"),
		CSS:  BoxString(&assets, "bootstrap.min.css"),
	}
}
