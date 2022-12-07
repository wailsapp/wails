package events

var Mac = newMacEvents()

type macEvents struct {
	ApplicationDidFinishLaunching string
	ApplicationWillTerminate      string
}

func newMacEvents() macEvents {
	return macEvents{
		ApplicationDidFinishLaunching: "mac:ApplicationDidFinishLaunching",
		ApplicationWillTerminate:      "mac:ApplicationWillTerminate",
	}
}
