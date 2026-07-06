//go:build windows

package edge

type COREWEBVIEW2_PHYSICAL_KEY_STATUS struct {
	RepeatCount   uint32
	ScanCode      uint32
	IsExtendedKey bool
	IsMenuKeyDown bool
	WasKeyDown    bool
	IsKeyReleased bool
}

// Bools need to be int32 in the native struct otherwise we end up in memory corruption. Using the internal
// struct is a hacky way so we don't break the public interface.
type internal_COREWEBVIEW2_PHYSICAL_KEY_STATUS struct {
	RepeatCount   uint32
	ScanCode      uint32
	IsExtendedKey int32
	IsMenuKeyDown int32
	WasKeyDown    int32
	IsKeyReleased int32
}
