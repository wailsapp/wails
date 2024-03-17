//go:build production && !devtools

package application

func newOpenDevToolsMenuItem() *MenuItem {
	return nil
}
