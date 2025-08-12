//go:build windows

package runtime

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

var invoke = `window._wails.invoke=window.chrome.webview.postMessage;`
var flags = fmt.Sprintf(
	`window._wails.flags={"system":{"resizeHandleWidth":%d,"resizeHandleHeight":%d}};`,
	w32.GetSystemMetrics(w32.SM_CXSIZEFRAME),
	w32.GetSystemMetrics(w32.SM_CYSIZEFRAME))
