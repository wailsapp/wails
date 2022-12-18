package application

type DialogType int

const (
	InfoDialog DialogType = iota
	QuestionDialog
	WarningDialog
	ErrorDialog
)

type Button struct {
	label     string
	isCancel  bool
	isDefault bool
	callback  func()
}

func (b *Button) OnClick(callback func()) {
	b.callback = callback
}

type dialogImpl interface {
	show()
}

type Dialog struct {
	dialogType DialogType
	title      string
	message    string
	buttons    []*Button

	// platform independent
	impl dialogImpl
	icon []byte
}

var defaultTitles = map[DialogType]string{
	InfoDialog:     "Information",
	QuestionDialog: "Question",
	WarningDialog:  "Warning",
	ErrorDialog:    "Error",
}

func newDialog(dialogType DialogType) *Dialog {
	return &Dialog{
		dialogType: dialogType,
		title:      defaultTitles[dialogType],
	}
}

func (d *Dialog) SetTitle(title string) *Dialog {
	d.title = title
	return d
}

func (d *Dialog) SetMessage(message string) *Dialog {
	d.message = message
	return d
}

func (d *Dialog) Show() {
	if d.impl == nil {
		d.impl = newDialogImpl(d)
	}
	d.impl.show()
}

func (d *Dialog) SetIcon(icon []byte) *Dialog {
	d.icon = icon
	return d
}

func (d *Dialog) AddButton(s string) *Button {
	result := &Button{
		label: s,
	}
	d.buttons = append(d.buttons, result)
	return result
}

func (d *Dialog) SetDefaultButton(button *Button) *Dialog {
	for _, b := range d.buttons {
		b.isDefault = false
	}
	button.isDefault = true
	return d
}

func (d *Dialog) SetCancelButton(button *Button) *Dialog {
	for _, b := range d.buttons {
		b.isCancel = false
	}
	button.isCancel = true
	return d
}
