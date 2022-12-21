//go:build darwin

package application

/*
#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#import <Cocoa/Cocoa.h>

extern void openFileDialogCallback(uint id, char* path);
extern void openFileDialogCallbackEnd(uint id);
extern void saveFileDialogCallback(uint id, char* path);

static void showAboutBox(char* title, char *message, void *icon, int length) {

	// run on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSAlert *alert = [[NSAlert alloc] init];
		if (title != NULL) {
			[alert setMessageText:[NSString stringWithUTF8String:title]];
			free(title);
		}
		if (message != NULL) {
			[alert setInformativeText:[NSString stringWithUTF8String:message]];
			free(message);
		}
		if (icon != NULL) {
			NSImage *image = [[NSImage alloc] initWithData:[NSData dataWithBytes:icon length:length]];
			[alert setIcon:image];
		}
		[alert setAlertStyle:NSAlertStyleInformational];
		[alert runModal];
	});
}


// Create an NSAlert
static void* createAlert(int alertType, char* title, char *message, void *icon, int length) {
	NSAlert *alert = [[NSAlert alloc] init];
	[alert setAlertStyle:alertType];
	if (title != NULL) {
		[alert setMessageText:[NSString stringWithUTF8String:title]];
		free(title);
	}
	if (message != NULL) {
		[alert setInformativeText:[NSString stringWithUTF8String:message]];
		free(message);
	}
	if (icon != NULL) {
		NSImage *image = [[NSImage alloc] initWithData:[NSData dataWithBytes:icon length:length]];
		[alert setIcon:image];
	} else {
		if(alertType == NSAlertStyleCritical || alertType == NSAlertStyleWarning) {
			NSImage *image = [NSImage imageNamed:NSImageNameCaution];
			[alert setIcon:image];
		} else {
			NSImage *image = [NSImage imageNamed:NSImageNameInfo];
			[alert setIcon:image];
		}
}
	return alert;

}

// Run the dialog
static int dialogRunModal(void *dialog) {
	NSAlert *alert = (__bridge NSAlert *)dialog;
    long response = [alert runModal];
    int result;

    if( response == NSAlertFirstButtonReturn ) {
        result = 0;
    }
    else if( response == NSAlertSecondButtonReturn ) {
        result = 1;
    }
    else if( response == NSAlertThirdButtonReturn ) {
        result = 2;
    } else {
        result = 3;
    }
	return result;
}

// Release the dialog
static void releaseDialog(void *dialog) {
	NSAlert *alert = (__bridge NSAlert *)dialog;
	[alert release];
}

// Add a button to the dialog
static void alertAddButton(void *dialog, char *label, bool isDefault, bool isCancel) {
	NSAlert *alert = (__bridge NSAlert *)dialog;
	NSButton *button = [alert addButtonWithTitle:[NSString stringWithUTF8String:label]];
	free(label);
    if( isDefault ) {
        [button setKeyEquivalent:@"\r"];
    } else if( isCancel ) {
        [button setKeyEquivalent:@"\033"];
    } else {
        [button setKeyEquivalent:@""];
    }
}

static void processOpenFileDialogResults(NSOpenPanel *panel, NSInteger result, bool allowsMultipleSelection, uint dialogID) {
	const char *path = NULL;
	if (result == NSModalResponseOK) {
		if (allowsMultipleSelection) {
			NSArray *urls = [panel URLs];
			for (NSURL *url in urls) {
				path = [[url path] UTF8String];
				openFileDialogCallback(dialogID, (char *)path);
			}

		} else {
			NSURL *url = [panel URL];
			path = [[url path] UTF8String];
			openFileDialogCallback(dialogID, (char *)path);
		}
	}
	openFileDialogCallbackEnd(dialogID);
}


static void showOpenFileDialog(unsigned int dialogID,
	bool canChooseFiles,
	bool canChooseDirectories,
	bool canCreateDirectories,
	bool showHiddenFiles,
	bool allowsMultipleSelection,
	bool resolvesAliases,
	bool hideExtension,
	bool treatsFilePackagesAsDirectories,
	bool allowsOtherFileTypes,
	char* message,
	char* directory,
	char* buttonText,
	void *window) {

	// run on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSOpenPanel *panel = [NSOpenPanel openPanel];

		if (message != NULL) {
			[panel setMessage:[NSString stringWithUTF8String:message]];
			free(message);
		}

		if (directory != NULL) {
			[panel setDirectoryURL:[NSURL fileURLWithPath:[NSString stringWithUTF8String:directory]]];
			free(directory);
		}

		if (buttonText != NULL) {
			[panel setPrompt:[NSString stringWithUTF8String:buttonText]];
			free(buttonText);
		}

		[panel setCanChooseFiles:canChooseFiles];
		[panel setCanChooseDirectories:canChooseDirectories];
		[panel setCanCreateDirectories:canCreateDirectories];
		[panel setShowsHiddenFiles:showHiddenFiles];
		[panel setAllowsMultipleSelection:allowsMultipleSelection];
		[panel setResolvesAliases:resolvesAliases];
		[panel setExtensionHidden:hideExtension];
		[panel setTreatsFilePackagesAsDirectories:treatsFilePackagesAsDirectories];
		[panel setAllowsOtherFileTypes:allowsOtherFileTypes];

		if (window != NULL) {
			[panel beginSheetModalForWindow:(__bridge NSWindow *)window completionHandler:^(NSInteger result) {
				processOpenFileDialogResults(panel, result, allowsMultipleSelection, dialogID);
			}];
		} else {
			[panel beginWithCompletionHandler:^(NSInteger result) {
				processOpenFileDialogResults(panel, result, allowsMultipleSelection, dialogID);
			}];
		}
	});
}

static void showSaveFileDialog(unsigned int dialogID,
	bool canCreateDirectories,
	bool showHiddenFiles,
	bool canSelectHiddenExtension,
	bool hideExtension,
	bool treatsFilePackagesAsDirectories,
	bool allowOtherFileTypes,
	char* message,
	char* directory,
	char* buttonText,
	char* filename,
	void *window) {

// run on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSSavePanel *panel = [NSSavePanel savePanel];

		if (message != NULL) {
			[panel setMessage:[NSString stringWithUTF8String:message]];
			free(message);
		}

		if (directory != NULL) {
			[panel setDirectoryURL:[NSURL fileURLWithPath:[NSString stringWithUTF8String:directory]]];
			free(directory);
		}

		if (filename != NULL) {
			[panel setNameFieldStringValue:[NSString stringWithUTF8String:filename]];
			free(filename);
		}

		if (buttonText != NULL) {
			[panel setPrompt:[NSString stringWithUTF8String:buttonText]];
			free(buttonText);
		}

		[panel setCanCreateDirectories:canCreateDirectories];
		[panel setShowsHiddenFiles:showHiddenFiles];
		[panel setCanSelectHiddenExtension:canSelectHiddenExtension];
		[panel setExtensionHidden:hideExtension];
		[panel setTreatsFilePackagesAsDirectories:treatsFilePackagesAsDirectories];
		[panel setAllowsOtherFileTypes:allowOtherFileTypes];

		if (window != NULL) {
			[panel beginSheetModalForWindow:(__bridge NSWindow *)window completionHandler:^(NSInteger result) {
				const char *path = NULL;
				if (result == NSModalResponseOK) {
					NSURL *url = [panel URL];
					const char *path = [[url path] UTF8String];
				}
				saveFileDialogCallback(dialogID, (char *)path);
			}];
		} else {
			[panel beginWithCompletionHandler:^(NSInteger result) {
				const char *path = NULL;
				if (result == NSModalResponseOK) {
					NSURL *url = [panel URL];
					const char *path = [[url path] UTF8String];
				}
				saveFileDialogCallback(dialogID, (char *)path);
			}];
		}
	});
}


*/
import "C"
import (
	"unsafe"
)

const NSAlertStyleWarning = C.int(0)
const NSAlertStyleInformational = C.int(1)
const NSAlertStyleCritical = C.int(2)

var alertTypeMap = map[DialogType]C.int{
	WarningDialog:  NSAlertStyleWarning,
	InfoDialog:     NSAlertStyleInformational,
	ErrorDialog:    NSAlertStyleCritical,
	QuestionDialog: NSAlertStyleInformational,
}

func (m *macosApp) showAboutDialog(title string, message string, icon []byte) {
	var iconData unsafe.Pointer
	if icon != nil {
		iconData = unsafe.Pointer(&icon[0])
	}
	C.showAboutBox(C.CString(title), C.CString(message), iconData, C.int(len(icon)))
}

type macosDialog struct {
	dialog *MessageDialog

	nsDialog unsafe.Pointer
}

func (m *macosDialog) show() {
	DispatchOnMainThread(func() {

		// Mac can only have 4 buttons on a dialog
		if len(m.dialog.buttons) > 4 {
			m.dialog.buttons = m.dialog.buttons[:4]
		}

		if m.nsDialog != nil {
			C.releaseDialog(m.nsDialog)
		}
		var title *C.char
		if m.dialog.title != "" {
			title = C.CString(m.dialog.title)
		}
		var message *C.char
		if m.dialog.message != "" {
			message = C.CString(m.dialog.message)
		}
		var iconData unsafe.Pointer
		var iconLength C.int
		if m.dialog.icon != nil {
			iconData = unsafe.Pointer(&m.dialog.icon[0])
			iconLength = C.int(len(m.dialog.icon))
		} else {
			// if it's an error, use the application icon
			if m.dialog.dialogType == ErrorDialog {
				iconData = unsafe.Pointer(&globalApplication.icon[0])
				iconLength = C.int(len(globalApplication.icon))
			}
		}

		alertType, ok := alertTypeMap[m.dialog.dialogType]
		if !ok {
			alertType = C.NSAlertStyleInformational
		}

		m.nsDialog = C.createAlert(alertType, title, message, iconData, iconLength)

		// Reverse the buttons so that the default is on the right
		reversedButtons := make([]*Button, len(m.dialog.buttons))
		var count = 0
		for i := len(m.dialog.buttons) - 1; i >= 0; i-- {
			button := m.dialog.buttons[i]
			C.alertAddButton(m.nsDialog, C.CString(button.label), C.bool(button.isDefault), C.bool(button.isCancel))
			reversedButtons[count] = m.dialog.buttons[i]
			count++
		}

		buttonPressed := int(C.dialogRunModal(m.nsDialog))
		if len(m.dialog.buttons) > buttonPressed {
			button := reversedButtons[buttonPressed]
			if button.callback != nil {
				button.callback()
			}
		}
	})

}

func newDialogImpl(d *MessageDialog) *macosDialog {
	return &macosDialog{
		dialog: d,
	}
}

type macosOpenFileDialog struct {
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *macosOpenFileDialog {
	return &macosOpenFileDialog{
		dialog: d,
	}
}

func toCString(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func (m *macosOpenFileDialog) show() ([]string, error) {
	openFileResponses[dialogID] = make(chan string)
	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		nsWindow = m.dialog.window.impl.(*macosWindow).nsWindow
	}
	C.showOpenFileDialog(C.uint(m.dialog.id),
		C.bool(m.dialog.canChooseFiles),
		C.bool(m.dialog.canChooseDirectories),
		C.bool(m.dialog.canCreateDirectories),
		C.bool(m.dialog.showHiddenFiles),
		C.bool(m.dialog.allowsMultipleSelection),
		C.bool(m.dialog.resolvesAliases),
		C.bool(m.dialog.hideExtension),
		C.bool(m.dialog.treatsFilePackagesAsDirectories),
		C.bool(m.dialog.allowsOtherFileTypes),
		toCString(m.dialog.message),
		toCString(m.dialog.directory),
		toCString(m.dialog.buttonText),
		nsWindow)
	var result []string
	for filename := range openFileResponses[m.dialog.id] {
		result = append(result, filename)
	}
	return result, nil
}

//export openFileDialogCallback
func openFileDialogCallback(id C.uint, path *C.char) {
	// Covert the path to a string
	filePath := C.GoString(path)
	// put response on channel
	channel, ok := openFileResponses[uint(id)]
	if ok {
		channel <- filePath
	} else {
		panic("No channel found for open file dialog")
	}
}

//export openFileDialogCallbackEnd
func openFileDialogCallbackEnd(id C.uint) {
	channel, ok := openFileResponses[uint(id)]
	if ok {
		close(channel)
	} else {
		panic("No channel found for open file dialog")
	}
}

type macosSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *macosSaveFileDialog {
	return &macosSaveFileDialog{
		dialog: d,
	}
}

func (m *macosSaveFileDialog) show() (string, error) {
	saveFileResponses[dialogID] = make(chan string)
	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		nsWindow = m.dialog.window.impl.(*macosWindow).nsWindow
	}
	C.showSaveFileDialog(C.uint(m.dialog.id),
		C.bool(m.dialog.canCreateDirectories),
		C.bool(m.dialog.showHiddenFiles),
		C.bool(m.dialog.canSelectHiddenExtension),
		C.bool(m.dialog.hideExtension),
		C.bool(m.dialog.treatsFilePackagesAsDirectories),
		C.bool(m.dialog.allowOtherFileTypes),
		toCString(m.dialog.message),
		toCString(m.dialog.directory),
		toCString(m.dialog.buttonText),
		toCString(m.dialog.filename),
		nsWindow)
	return <-saveFileResponses[m.dialog.id], nil
}

//export saveFileDialogCallback
func saveFileDialogCallback(id C.uint, path *C.char) {
	// Covert the path to a string
	filePath := C.GoString(path)
	// put response on channel
	channel, ok := saveFileResponses[uint(id)]
	if ok {
		channel <- filePath
	} else {
		panic("No channel found for save file dialog")
	}
}
