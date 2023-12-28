//go:build production

package runtime

import _ "embed"

//go:embed runtime.js
var DesktopRuntime []byte

var RuntimeAssetsBundle = &RuntimeAssets{
	runtimeDesktopJS: DesktopRuntime,
}

type RuntimeAssets struct {
	runtimeDesktopJS []byte
}

func (r *RuntimeAssets) DesktopIPC() []byte {
	return []byte("")
}

func (r *RuntimeAssets) WebsocketIPC() []byte {
	return []byte("")
}

func (r *RuntimeAssets) RuntimeDesktopJS() []byte {
	return r.runtimeDesktopJS
}
