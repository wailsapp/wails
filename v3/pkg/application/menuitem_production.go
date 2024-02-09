//go:build production && !devtools

package application

func newShowDevToolsMenuItem() *MenuItem {
	return nil
}
