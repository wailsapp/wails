package wails

import (
	"log"
	"strings"

	"github.com/wailsapp/wails/cmd"

	"github.com/gobuffalo/packr"
)

func escapeJS(js string) (string, error) {
	result := strings.Replace(js, "\\", "\\\\", -1)
	result = strings.Replace(result, "'", "\\'", -1)
	result = strings.Replace(result, "\n", "\\n", -1)
	return result, nil
}

// BoxString wraps packr.FindString
func BoxString(box *packr.Box, filename string) string {
	result, err := box.FindString(filename)
	if err != nil {
		if BuildMode == cmd.BuildModeBridge {
			return ""
		}
		log.Fatal(err)
	}
	return result
}
