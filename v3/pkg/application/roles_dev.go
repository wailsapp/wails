//go:build !production || devtools

package application

func addDevToolMenuItem(viewMenu *Menu) {
	viewMenu.AddRole(ShowDevTools)
}
