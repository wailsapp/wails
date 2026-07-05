//go:build windows

package webview2

type COREWEBVIEW2_PHYSICAL_KEY_STATUS struct {
	RepeatCount uint32
	ScanCode uint32
	IsExtendedKey int32
	IsMenuKeyDown int32
	WasKeyDown int32
	IsKeyReleased int32
}
