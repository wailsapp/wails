package events

import "runtime"

var defaultWindowEventMapping = map[string]map[WindowEventType]WindowEventType{
	"windows": {
		Windows.WindowClosing:      Common.WindowClosing,
		Windows.WindowInactive:     Common.WindowLostFocus,
		Windows.WindowClickActive:  Common.WindowFocus,
		Windows.WindowActive:       Common.WindowFocus,
		Windows.WindowMaximise:     Common.WindowMaximise,
		Windows.WindowMinimise:     Common.WindowMinimise,
		Windows.WindowRestore:      Common.WindowRestore,
		Windows.WindowUnMaximise:   Common.WindowUnMaximise,
		Windows.WindowUnMinimise:   Common.WindowUnMinimise,
		Windows.WindowFullscreen:   Common.WindowFullscreen,
		Windows.WindowUnFullscreen: Common.WindowUnFullscreen,
		Windows.WindowShow:         Common.WindowShow,
		Windows.WindowHide:         Common.WindowHide,
		Windows.WindowDidMove:      Common.WindowDidMove,
		Windows.WindowDidResize:    Common.WindowDidResize,
		Windows.WindowSetFocus:     Common.WindowFocus,
		Windows.WindowKillFocus:    Common.WindowLostFocus,
	},
	"darwin": {
		Mac.WindowDidResignKey:       Common.WindowLostFocus,
		Mac.WindowDidBecomeKey:       Common.WindowFocus,
		Mac.WindowDidMiniaturize:     Common.WindowMinimise,
		Mac.WindowDidDeminiaturize:   Common.WindowUnMinimise,
		Mac.WindowDidEnterFullScreen: Common.WindowFullscreen,
		Mac.WindowDidExitFullScreen:  Common.WindowUnFullscreen,
		Mac.WindowMaximise:           Common.WindowMaximise,
		Mac.WindowUnMaximise:         Common.WindowUnMaximise,
		Mac.WindowDidMove:            Common.WindowDidMove,
		Mac.WindowDidResize:          Common.WindowDidResize,
		Mac.WindowDidUpdate:          Common.WindowShow,
		Mac.WindowDidZoom:            Common.WindowMaximise,
		Mac.WindowZoomIn:             Common.WindowZoomIn,
		Mac.WindowZoomOut:            Common.WindowZoomOut,
		Mac.WindowZoomReset:          Common.WindowZoomReset,
	},
	"linux": {
		Linux.WindowDeleteEvent: Common.WindowClosing,
		Linux.WindowFocusIn:     Common.WindowFocus,
		Linux.WindowFocusOut:    Common.WindowLostFocus,
		Linux.WindowDidMove:     Common.WindowDidMove,
		Linux.WindowDidResize:   Common.WindowDidResize,
		Linux.WindowLoadChanged: Common.WindowShow,
	},
}

func DefaultWindowEventMapping() map[WindowEventType]WindowEventType {
	platform := runtime.GOOS
	return defaultWindowEventMapping[platform]
}
