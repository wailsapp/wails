package frontend

type Frontend interface {

	// Main methods
	Run() error
	Quit()

	//// Events
	//NotifyEvent(message string)
	//CallResult(message string)
	//
	//// Dialog
	//OpenFileDialog(dialogOptions dialog.OpenDialogOptions, callbackID string)
	//OpenMultipleFilesDialog(dialogOptions dialog.OpenDialogOptions, callbackID string)
	//OpenDirectoryDialog(dialogOptions dialog.OpenDialogOptions, callbackID string)
	//SaveDialog(dialogOptions dialog.SaveDialogOptions, callbackID string)
	//MessageDialog(dialogOptions dialog.MessageDialogOptions, callbackID string)

	// Window
	//WindowSetTitle(title string)
	WindowShow()
	WindowHide()
	//WindowCenter()
	//WindowMaximise()
	//WindowUnmaximise()
	//WindowMinimise()
	//WindowUnminimise()
	//WindowPosition(x int, y int)
	//WindowSize(width int, height int)
	//WindowSetMinSize(width int, height int)
	//WindowSetMaxSize(width int, height int)
	WindowFullscreen()
	WindowUnFullscreen()
	//WindowSetColour(colour int)
	//
	//// Menus
	//SetApplicationMenu(menu *menu.Menu)
	//SetTrayMenu(menu *menu.TrayMenu)
	//UpdateTrayMenuLabel(menu *menu.TrayMenu)
	//UpdateContextMenu(contextMenu *menu.ContextMenu)
	//DeleteTrayMenuByID(id string)
}
