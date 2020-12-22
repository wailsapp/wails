package menu

type TrayType string

const (
	TrayIcon  TrayType = "icon"
	TrayLabel TrayType = "label"
)

type TrayOptions struct {
	Type  TrayType
	Label string
	Menu  *Menu
}
