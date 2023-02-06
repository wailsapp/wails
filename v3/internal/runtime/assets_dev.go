//go:build !production

package runtime

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
