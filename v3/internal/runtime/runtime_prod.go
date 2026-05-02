//go:build production

package runtime

import "fmt"
import goruntime "runtime"

var environment = fmt.Sprintf(`window._wails.environment={"OS":"%s","Arch":"%s","Debug":false};`, goruntime.GOOS, goruntime.GOARCH)
