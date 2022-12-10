package events

var Mac = newMacEvents()

type macEvents struct {
	ApplicationDidFinishLaunching  string
	ApplicationWillTerminate       string
	ApplicationDidBecomeActive     string
	ApplicationWillUpdate          string
	ApplicationDidUpdate           string
	ApplicationWillFinishLaunching string
	ApplicationWillHide            string
	ApplicationWillUnhide          string
	ApplicationDidHide             string
	ApplicationDidUnhide           string
}

func newMacEvents() macEvents {
	return macEvents{
		ApplicationDidFinishLaunching:  "mac:ApplicationDidFinishLaunching",
		ApplicationWillTerminate:       "mac:ApplicationWillTerminate",
		ApplicationDidBecomeActive:     "mac:ApplicationDidBecomeActive",
		ApplicationWillUpdate:          "mac:ApplicationWillUpdate",
		ApplicationDidUpdate:           "mac:ApplicationDidUpdate",
		ApplicationWillFinishLaunching: "mac:ApplicationWillFinishLaunching",
		ApplicationWillHide:            "mac:ApplicationWillHide",
		ApplicationWillUnhide:          "mac:ApplicationWillUnhide",
		ApplicationDidHide:             "mac:ApplicationDidHide",
		ApplicationDidUnhide:           "mac:ApplicationDidUnhide",
	}
}
