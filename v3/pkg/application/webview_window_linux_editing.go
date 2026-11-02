//go:build linux && !android && !server

package application

// The Edit menu's roles must edit the page: an Edit menu binds Ctrl+X/C/V/A/Z,
// and both GTK stacks hand those shortcuts to the menu before the webview
// sees the key (#4723).
//
// The clipboard and selection roles run WebKit's own editing commands (the
// WEBKIT_EDITING_COMMAND_* names): document.execCommand from an evaluated
// script cannot paste, WebKit refuses it. Undo and redo stay in the page,
// because WebKitGTK's native undo is unreliable for <input> elements.
//
// This is the logic every Linux backend shares; each one only supplies
// executeEditingCommand, the call to webkit_web_view_execute_editing_command.

// editingHook, when set, receives each role's command name or script instead
// of the webview. Tests use it to check the roles without a WebKit view.
var editingHook func(commandOrScript string)

func (w *linuxWebviewWindow) editingCommand(command string) {
	if editingHook != nil {
		editingHook(command)
		return
	}
	InvokeAsync(func() { w.executeEditingCommand(command) })
}

func (w *linuxWebviewWindow) editingScript(js string) {
	if editingHook != nil {
		editingHook(js)
		return
	}
	w.execJS(js)
}

func (w *linuxWebviewWindow) cut()       { w.editingCommand("Cut") }
func (w *linuxWebviewWindow) copy()      { w.editingCommand("Copy") }
func (w *linuxWebviewWindow) paste()     { w.editingCommand("Paste") }
func (w *linuxWebviewWindow) delete()    { w.editingCommand("Delete") }
func (w *linuxWebviewWindow) selectAll() { w.editingCommand("SelectAll") }
func (w *linuxWebviewWindow) undo()      { w.editingScript("document.execCommand('undo')") }
func (w *linuxWebviewWindow) redo()      { w.editingScript("document.execCommand('redo')") }
