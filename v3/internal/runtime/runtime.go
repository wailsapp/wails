package runtime

import (
	"encoding/json"
	"fmt"
)

var runtimeInit = `window._wails=window._wails||{};window.wails=window.wails||{};`

func Core(flags map[string]any) string {
	flagsStr := ""
	if len(flags) > 0 {
		f, err := json.Marshal(flags)
		if err == nil {
			flagsStr += fmt.Sprintf("window._wails.flags=%s;", f)
		}
	}

	return runtimeInit + flagsStr + invoke + environment
}
