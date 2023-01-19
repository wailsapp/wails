//go:build dev

package runtime

var RuntimeAssetsBundle = &RuntimeAssets{
	desktopIPC:       DesktopIPC,
	websocketIPC:     WebsocketIPC,
	runtimeDesktopJS: RuntimeDesktopJS,
}

type RuntimeAssets struct {
	desktopIPC       []byte
	websocketIPC     []byte
	runtimeDesktopJS []byte
}

func (r *RuntimeAssets) DesktopIPC() []byte {
	return r.desktopIPC
}

func (r *RuntimeAssets) WebsocketIPC() []byte {
	return r.websocketIPC
}

func (r *RuntimeAssets) RuntimeDesktopJS() []byte {
	return r.runtimeDesktopJS
}
