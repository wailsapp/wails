// +build frameworkbootstrap4lux

package frameworks

import (
	"github.com/gobuffalo/packr"
)

func init() {
	assets := packr.NewBox("./bootstrap4lux/assets")
	FrameworkToUse = &Framework{
		Name: "Bootstrap 4 (Lux)",
		JS:   assets.String("bootstrap.bundle.min.js"),
		CSS:  assets.String("bootstrap.min.css"),
	}
}
