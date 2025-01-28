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
		Windows.WindowDPIChanged:   Common.WindowDPIChanged,
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
		Mac.WindowDidZoom:            Common.WindowMaximise,
		Mac.WindowShow:               Common.WindowShow,
		Mac.WindowHide:               Common.WindowHide,
		Mac.WindowZoomIn:             Common.WindowZoomIn,
		Mac.WindowZoomOut:            Common.WindowZoomOut,
		Mac.WindowZoomReset:          Common.WindowZoomReset,
		Mac.WindowShouldClose:        Common.WindowClosing,
		Mac.WindowDidResignKey:       Common.WindowLostFocus,
		Mac.WindowDidResignMain:      Common.WindowLostFocus,
		Mac.WindowDidResize:          Common.WindowDidResize,
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
