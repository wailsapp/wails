//go:build production && !devtools

package application

func NewOpenDevToolsMenuItem() *MenuItem {
	return nil
}
