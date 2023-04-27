//go:build linux

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <stdio.h>

static GtkWidget* new_about_dialog(GtkWindow *parent, const gchar *msg) {
   // gtk_message_dialog_new is variadic!  Can't call from cgo
   GtkWidget *dialog;
   dialog = gtk_message_dialog_new(
       parent,
       GTK_DIALOG_MODAL | GTK_DIALOG_DESTROY_WITH_PARENT,
	   GTK_MESSAGE_INFO,
	   GTK_BUTTONS_CLOSE,
       msg);

   g_signal_connect_swapped (dialog,
                             "response",
                             G_CALLBACK (gtk_widget_destroy),
                             dialog);
   return dialog;
};

*/
import "C"
import (
	"fmt"
	"unsafe"
)

const AlertStyleWarning = C.int(0)
const AlertStyleInformational = C.int(1)
const AlertStyleCritical = C.int(2)

var alertTypeMap = map[DialogType]C.int{
	WarningDialog:  AlertStyleWarning,
	InfoDialog:     AlertStyleInformational,
	ErrorDialog:    AlertStyleCritical,
	QuestionDialog: AlertStyleInformational,
}

func setWindowIcon(window *C.GtkWindow, icon []byte) {
	fmt.Println("setWindowIcon", len(icon))
	loader := C.gdk_pixbuf_loader_new()
	if loader == nil {
		return
	}
	written := C.gdk_pixbuf_loader_write(
		loader,
		(*C.uchar)(&icon[0]),
		C.ulong(len(icon)),
		nil)
	if written == 0 {
		fmt.Println("failed to write icon")
		return
	}
	C.gdk_pixbuf_loader_close(loader, nil)
	pixbuf := C.gdk_pixbuf_loader_get_pixbuf(loader)
	if pixbuf != nil {
		fmt.Println("gtk_window_set_icon", window)
		C.gtk_window_set_icon((*C.GtkWindow)(window), pixbuf)
	}
	C.g_object_unref(C.gpointer(loader))
}

func (m *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	globalApplication.dispatchOnMainThread(func() {
		parent := C.gtk_application_get_active_window((*C.GtkApplication)(m.application))
		cMsg := C.CString(message)
		cTitle := C.CString(title)
		defer C.free(unsafe.Pointer(cMsg))
		defer C.free(unsafe.Pointer(cTitle))
		dialog := C.new_about_dialog(parent, cMsg)
		C.gtk_window_set_title(
			(*C.GtkWindow)(unsafe.Pointer(dialog)),
			cTitle)
		//		setWindowIcon((*C.GtkWindow)(dialog), icon)
		C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(dialog)))

	})
}

type linuxDialog struct {
	dialog *MessageDialog

	//nsDialog unsafe.Pointer
}

func (m *linuxDialog) show() {
	globalApplication.dispatchOnMainThread(func() {

		// Mac can only have 4 Buttons on a dialog
		if len(m.dialog.Buttons) > 4 {
			m.dialog.Buttons = m.dialog.Buttons[:4]
		}

		// if m.nsDialog != nil {
		// 	//C.releaseDialog(m.nsDialog)
		// }
		// var title *C.char
		// if m.dialog.Title != "" {
		// 	title = C.CString(m.dialog.Title)
		// }
		// var message *C.char
		// if m.dialog.Message != "" {
		// 	message = C.CString(m.dialog.Message)
		// }
		// var iconData unsafe.Pointer
		// var iconLength C.int
		// if m.dialog.Icon != nil {
		// 	iconData = unsafe.Pointer(&m.dialog.Icon[0])
		// 	iconLength = C.int(len(m.dialog.Icon))
		// } else {
		// 	// if it's an error, use the application Icon
		// 	if m.dialog.DialogType == ErrorDialog {
		// 		iconData = unsafe.Pointer(&globalApplication.options.Icon[0])
		// 		iconLength = C.int(len(globalApplication.options.Icon))
		// 	}
		// }

		// alertType, ok := alertTypeMap[m.dialog.DialogType]
		// if !ok {
		// 	alertType = AlertStyleInformational
		// }

		//		m.nsDialog = C.createAlert(alertType, title, message, iconData, iconLength)

		// Reverse the Buttons so that the default is on the right
		reversedButtons := make([]*Button, len(m.dialog.Buttons))
		var count = 0
		for i := len(m.dialog.Buttons) - 1; i >= 0; i-- {
			//button := m.dialog.Buttons[i]
			//C.alertAddButton(m.nsDialog, C.CString(button.Label), C.bool(button.IsDefault), C.bool(button.IsCancel))
			reversedButtons[count] = m.dialog.Buttons[i]
			count++
		}

		buttonPressed := int(0) //C.dialogRunModal(m.nsDialog))
		if len(m.dialog.Buttons) > buttonPressed {
			button := reversedButtons[buttonPressed]
			if button.callback != nil {
				button.callback()
			}
		}
	})

}

func newDialogImpl(d *MessageDialog) *linuxDialog {
	return &linuxDialog{
		dialog: d,
	}
}

type linuxOpenFileDialog struct {
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *linuxOpenFileDialog {
	return &linuxOpenFileDialog{
		dialog: d,
	}
}

func toCString(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func (m *linuxOpenFileDialog) show() ([]string, error) {
	openFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//nsWindow = m.dialog.window.impl.(*macosWebviewWindow).nsWindow
	}

	// Massage filter patterns into macOS format
	// We iterate all filter patterns, tidy them up and then join them with a semicolon
	// This should produce a single string of extensions like "png;jpg;gif"
	// 	var filterPatterns string
	// if len(m.dialog.filters) > 0 {
	// 	var allPatterns []string
	// 	for _, filter := range m.dialog.filters {
	// 		patternComponents := strings.Split(filter.Pattern, ";")
	// 		for i, component := range patternComponents {
	// 			filterPattern := strings.TrimSpace(component)
	// 			filterPattern = strings.TrimPrefix(filterPattern, "*.")
	// 			patternComponents[i] = filterPattern
	// 		}
	// 		allPatterns = append(allPatterns, strings.Join(patternComponents, ";"))
	// 	}
	// 	filterPatterns = strings.Join(allPatterns, ";")
	// }

	// C.showOpenFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canChooseFiles),
	// 	C.bool(m.dialog.canChooseDirectories),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.allowsMultipleSelection),
	// 	C.bool(m.dialog.resolvesAliases),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowsOtherFileTypes),
	// 	toCString(filterPatterns),
	// 	C.uint(len(filterPatterns)),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	nsWindow)
	var result []string
	for filename := range openFileResponses[m.dialog.id] {
		result = append(result, filename)
	}
	return result, nil
}

//export openFileDialogCallback
func openFileDialogCallback(cid C.uint, cpath *C.char) {
	path := C.GoString(cpath)
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		channel <- path
	} else {
		panic("No channel found for open file dialog")
	}
}

//export openFileDialogCallbackEnd
func openFileDialogCallbackEnd(cid C.uint) {
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		close(channel)
		delete(openFileResponses, id)
		freeDialogID(id)
	} else {
		panic("No channel found for open file dialog")
	}
}

type linuxSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *linuxSaveFileDialog {
	return &linuxSaveFileDialog{
		dialog: d,
	}
}

func (m *linuxSaveFileDialog) show() (string, error) {
	saveFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//		nsWindow = m.dialog.window.impl.(*linuxWebviewWindow).nsWindow
	}

	// C.showSaveFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.canSelectHiddenExtension),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowOtherFileTypes),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	toCString(m.dialog.filename),
	// 	nsWindow)
	return <-saveFileResponses[m.dialog.id], nil
}

//export saveFileDialogCallback
func saveFileDialogCallback(cid C.uint, cpath *C.char) {
	// Covert the path to a string
	path := C.GoString(cpath)
	id := uint(cid)
	// put response on channel
	channel, ok := saveFileResponses[id]
	if ok {
		channel <- path
		close(channel)
		delete(saveFileResponses, id)
		freeDialogID(id)

	} else {
		panic("No channel found for save file dialog")
	}
}
