package wails

import (
	"strings"
)

func escapeJS(js string) (string, error) {
	result := strings.Replace(js, "\\", "\\\\", -1)
	result = strings.Replace(result, "'", "\\'", -1)
	result = strings.Replace(result, "\n", "\\n", -1)
	return result, nil
}
