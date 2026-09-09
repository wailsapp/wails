package application

type titlebarDoubleClickAction uint8

const (
	titlebarDoubleClickNone titlebarDoubleClickAction = iota
	titlebarDoubleClickToggleMaximise
	titlebarDoubleClickMinimise
)

func parseMacTitlebarDoubleClickAction(preference string) titlebarDoubleClickAction {
	switch preference {
	case "", "Maximize":
		return titlebarDoubleClickToggleMaximise
	case "Minimize":
		return titlebarDoubleClickMinimise
	default:
		return titlebarDoubleClickNone
	}
}
