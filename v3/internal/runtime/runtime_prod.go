//go:build production

package runtime

import "fmt"
import goruntime "runtime"

var environment = fmt.Sprintf(`window._wails.environment={"OS":"%s","Arch":"%s","Debug":true};`, goruntime.GOOS, goruntime.GOARCH)
