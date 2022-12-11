package events

var Mac = newMacEvents()

type macEvents struct {
	ApplicationWillFinishLaunching           string
	ApplicationDidFinishLaunching            string
	ApplicationWillBecomeActive              string
	ApplicationDidBecomeActive               string
	ApplicationWillUpdate                    string
	ApplicationDidUpdate                     string
	ApplicationWillHide                      string
	ApplicationDidHide                       string
	ApplicationWillUnhide                    string
	ApplicationDidUnhide                     string
	ApplicationWillResignActive              string
	ApplicationDidResignActive               string
	ApplicationWillTerminate                 string
	ApplicationDidChangeOcclusionState       string
	ApplicationDidChangeScreenParameters     string
	ApplicationDidChangeBackingProperties    string
	ApplicationDidChangeIcon                 string
	ApplicationDidChangeStatusBarOrientation string
	ApplicationDidChangeStatusBarFrame       string
	ApplicationDidChangeEffectiveAppearance  string
}

func newMacEvents() macEvents {
	return macEvents{
		ApplicationWillFinishLaunching:           "mac:ApplicationWillFinishLaunching",
		ApplicationDidFinishLaunching:            "mac:ApplicationDidFinishLaunching",
		ApplicationWillBecomeActive:              "mac:ApplicationWillBecomeActive",
		ApplicationDidBecomeActive:               "mac:ApplicationDidBecomeActive",
		ApplicationWillUpdate:                    "mac:ApplicationWillUpdate",
		ApplicationDidUpdate:                     "mac:ApplicationDidUpdate",
		ApplicationWillHide:                      "mac:ApplicationWillHide",
		ApplicationDidHide:                       "mac:ApplicationDidHide",
		ApplicationWillUnhide:                    "mac:ApplicationWillUnhide",
		ApplicationDidUnhide:                     "mac:ApplicationDidUnhide",
		ApplicationWillResignActive:              "mac:ApplicationWillResignActive",
		ApplicationDidResignActive:               "mac:ApplicationDidResignActive",
		ApplicationWillTerminate:                 "mac:ApplicationWillTerminate",
		ApplicationDidChangeOcclusionState:       "mac:ApplicationDidChangeOcclusionState",
		ApplicationDidChangeScreenParameters:     "mac:ApplicationDidChangeScreenParameters",
		ApplicationDidChangeBackingProperties:    "mac:ApplicationDidChangeBackingProperties",
		ApplicationDidChangeIcon:                 "mac:ApplicationDidChangeIcon",
		ApplicationDidChangeStatusBarOrientation: "mac:ApplicationDidChangeStatusBarOrientation",
		ApplicationDidChangeStatusBarFrame:       "mac:ApplicationDidChangeStatusBarFrame",
		ApplicationDidChangeEffectiveAppearance:  "mac:ApplicationDidChangeEffectiveAppearance",
	}
}
