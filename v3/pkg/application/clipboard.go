package application

type clipboardImpl interface {
	setText(text string) bool
	text() (string, bool)
}

type Clipboard struct {
	impl clipboardImpl
}

func newClipboard() *Clipboard {
	return &Clipboard{
		impl: newClipboardImpl(),
	}
}

func (c *Clipboard) SetText(text string) bool {
	return c.impl.setText(text)
}

func (c *Clipboard) Text() (string, bool) {
	return c.impl.text()
}
