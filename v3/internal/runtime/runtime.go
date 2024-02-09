package runtime

var runtimeInit = `window._wails=window._wails||{};window.wails=window.wails||{};`

func Core() string {
	return runtimeInit + flags + invoke + environment
}
