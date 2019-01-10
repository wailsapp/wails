package frameworks

import (
	"log"

	"github.com/gobuffalo/packr"
)

// Framework has details about a specific framework
type Framework struct {
	Name    string
	JS      string
	CSS     string
	Options string
}

// FrameworkToUse is the framework we will use when building
// Set by `wails init`, used by `wails build`
var FrameworkToUse *Framework

// BoxString extracts a string from a packr box
func BoxString(box *packr.Box, filename string) string {
	result, err := box.FindString(filename)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
