// +build frameworkbootstrap4

package frameworks

import (
	"github.com/gobuffalo/packr"
)

func init() {
	assets := packr.NewBox("./bootstrap4default/assets")
	FrameworkToUse = &Framework{
		Name: "Bootstrap 4",
		JS:   BoxString(&assets, "bootstrap.bundle.min.js"),
		CSS:  BoxString(&assets, "bootstrap.min.css"),
	}
}
