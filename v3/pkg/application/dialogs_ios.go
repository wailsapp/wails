//go:build ios

package application

// dialogsImpl implements dialogs for iOS
type dialogsImpl struct {
	// iOS-specific fields if needed
}

func newDialogsImpl() *dialogsImpl {
	return &dialogsImpl{}
}

// iOS dialog implementations would use UIAlertController
// These are placeholder implementations for now

func (d *dialogsImpl) info(id uint, param MessageDialogOptions) {
	// TODO: Implement using UIAlertController
}

func (d *dialogsImpl) warning(id uint, param MessageDialogOptions) {
	// TODO: Implement using UIAlertController
}

func (d *dialogsImpl) error(id uint, param MessageDialogOptions) {
	// TODO: Implement using UIAlertController
}

func (d *dialogsImpl) question(id uint, param MessageDialogOptions) chan bool {
	// TODO: Implement using UIAlertController
	ch := make(chan bool, 1)
	ch <- false
	return ch
}

func (d *dialogsImpl) openFile(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using UIDocumentPickerViewController
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) openMultipleFiles(id uint, param OpenFileDialogOptions) chan []string {
	// TODO: Implement using UIDocumentPickerViewController
	ch := make(chan []string, 1)
	ch <- []string{}
	return ch
}

func (d *dialogsImpl) openDirectory(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using UIDocumentPickerViewController
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) saveFile(id uint, param SaveFileDialogOptions) chan string {
	// TODO: Implement using UIDocumentPickerViewController
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

type iosDialog struct {
	dialog *MessageDialog
}

func (d *iosDialog) show() (string, error) {
	// TODO: Implement using UIAlertController
	return "", nil
}

func newDialogImpl(d *MessageDialog) *iosDialog {
	return &iosDialog{
		dialog: d,
	}
}

func (d *dialogsImpl) show() (chan string, error) {
	ch := make(chan string, 1)
	ch <- ""
	return ch, nil
}

func newOpenFileDialogImpl(_ *OpenFileDialogStruct) openFileDialogImpl {
	return &dialogsImpl{}
}

func newSaveFileDialogImpl(_ *SaveFileDialogStruct) saveFileDialogImpl {
	return &dialogsImpl{}
}
