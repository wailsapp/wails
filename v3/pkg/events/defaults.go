package events

import "runtime"

var defaultWindowEventMapping = map[string]map[WindowEventType]WindowEventType{
	"windows": {
		Windows.WindowClose:        Common.WindowClosing,
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
	},
	"darwin": {
		Mac.WindowDidResignKey:       Common.WindowLostFocus,
		Mac.WindowDidResignKey:       Common.WindowLostFocus,
		Mac.WindowDidBecomeKey:       Common.WindowFocus,
		Mac.WindowDidMiniaturize:     Common.WindowMinimise,
		Mac.WindowDidDeminiaturize:   Common.WindowUnMinimise,
		Mac.WindowDidEnterFullScreen: Common.WindowFullscreen,
		Mac.WindowDidExitFullScreen:  Common.WindowUnFullscreen,
	},
	"linux": {},
}

func DefaultWindowEventMapping() map[WindowEventType]WindowEventType {
	platform := runtime.GOOS
	return defaultWindowEventMapping[platform]
}
