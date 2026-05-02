//go:build android

package application

// dialogsImpl implements dialogs for Android
type dialogsImpl struct {
	// Android-specific fields if needed
}

func newDialogsImpl() *dialogsImpl {
	return &dialogsImpl{}
}

// Android dialog implementations would use AlertDialog
// These are placeholder implementations for now

func (d *dialogsImpl) info(id uint, param MessageDialogOptions) {
	// TODO: Implement using AlertDialog
}

func (d *dialogsImpl) warning(id uint, param MessageDialogOptions) {
	// TODO: Implement using AlertDialog
}

func (d *dialogsImpl) error(id uint, param MessageDialogOptions) {
	// TODO: Implement using AlertDialog
}

func (d *dialogsImpl) question(id uint, param MessageDialogOptions) chan bool {
	// TODO: Implement using AlertDialog
	ch := make(chan bool, 1)
	ch <- false
	return ch
}

func (d *dialogsImpl) openFile(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) openMultipleFiles(id uint, param OpenFileDialogOptions) chan []string {
	// TODO: Implement using Android file picker intent
	ch := make(chan []string, 1)
	ch <- []string{}
	return ch
}

func (d *dialogsImpl) openDirectory(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) saveFile(id uint, param SaveFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

type androidDialog struct {
	dialog *MessageDialog
}

func (d *androidDialog) show() {
	// TODO: Implement using AlertDialog
}

func newDialogImpl(d *MessageDialog) *androidDialog {
	return &androidDialog{
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
