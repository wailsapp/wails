package application

import (
	"errors"
	"sync"
	"time"
)

type clipboardImpl interface {
	setText(text string) bool
	text() (string, bool)
}

// clipboardExtendedImpl is the optional rich clipboard surface. A platform
// implementation that also satisfies it unlocks the image, file, HTML, RTF
// and raw-data methods on Clipboard and ClipboardManager; the manager
// type-asserts for it so platforms that only implement clipboardImpl keep
// working unchanged. Every method is called on the main thread.
type clipboardExtendedImpl interface {
	setImage(png []byte) error
	image() ([]byte, error)
	setFiles(paths []string) error
	files() ([]string, error)
	setHTML(html string, plain string) error
	html() (string, error)
	setRTF(rtf []byte) error
	rtf() ([]byte, error)
	setData(uti string, data []byte) error
	data(uti string) ([]byte, error)
	types() []string
	clear()
	changeCount() int
}

// ErrClipboardNotSupported is returned by the rich clipboard methods on
// platforms that only support plain text (currently everything except
// macOS).
var ErrClipboardNotSupported = errors.New("this clipboard operation is not supported on this platform")

// clipboardPollInterval is how often OnChange listeners poll the pasteboard
// change count. macOS has no clipboard change notification, so polling is
// the only option; the ticker only runs while at least one listener exists.
const clipboardPollInterval = 500 * time.Millisecond

type Clipboard struct {
	impl clipboardImpl

	// dispatch runs fn on the main thread. It is nil in production, which
	// means InvokeSync; tests replace it so no application is needed.
	dispatch func(fn func())

	changeLock      sync.Mutex
	changeListeners map[uint]func()
	nextListenerID  uint
	pollStop        chan struct{}
	pollInterval    time.Duration
}

func newClipboard() *Clipboard {
	return &Clipboard{
		impl: newClipboardImpl(),
	}
}

func (c *Clipboard) invoke(fn func()) {
	if c.dispatch != nil {
		c.dispatch(fn)
		return
	}
	InvokeSync(fn)
}

// extended returns the rich clipboard implementation when the platform
// provides one.
func (c *Clipboard) extended() (clipboardExtendedImpl, bool) {
	impl, ok := c.impl.(clipboardExtendedImpl)
	return impl, ok
}

func (c *Clipboard) SetText(text string) bool {
	return InvokeSyncWithResult(func() bool {
		return c.impl.setText(text)
	})
}

func (c *Clipboard) Text() (string, bool) {
	return InvokeSyncWithResultAndOther(c.impl.text)
}

// SetImage replaces the clipboard contents with a PNG image. On macOS the
// image is written as both PNG and TIFF so every paste target can read it.
func (c *Clipboard) SetImage(png []byte) error {
	impl, ok := c.extended()
	if !ok {
		return ErrClipboardNotSupported
	}
	if len(png) == 0 {
		return errors.New("clipboard: image data is empty")
	}
	var err error
	c.invoke(func() { err = impl.setImage(png) })
	return err
}

// Image returns the clipboard image as PNG. TIFF clipboard images (the
// format most macOS apps copy) are converted; a clipboard without an image
// returns an error.
func (c *Clipboard) Image() ([]byte, error) {
	impl, ok := c.extended()
	if !ok {
		return nil, ErrClipboardNotSupported
	}
	var (
		result []byte
		err    error
	)
	c.invoke(func() { result, err = impl.image() })
	return result, err
}

// SetFiles replaces the clipboard contents with file references, as Finder
// does for Edit > Copy. Paths should be absolute.
func (c *Clipboard) SetFiles(paths []string) error {
	impl, ok := c.extended()
	if !ok {
		return ErrClipboardNotSupported
	}
	if len(paths) == 0 {
		return errors.New("clipboard: no file paths given")
	}
	var err error
	c.invoke(func() { err = impl.setFiles(paths) })
	return err
}

// Files returns the file paths on the clipboard. An empty slice and a nil
// error means the clipboard holds no file references.
func (c *Clipboard) Files() ([]string, error) {
	impl, ok := c.extended()
	if !ok {
		return nil, ErrClipboardNotSupported
	}
	var (
		result []string
		err    error
	)
	c.invoke(func() { result, err = impl.files() })
	return result, err
}

// SetHTML replaces the clipboard contents with an HTML fragment plus a plain
// text fallback for targets that cannot paste HTML. When plain is empty only
// the HTML representation is written.
func (c *Clipboard) SetHTML(html string, plain string) error {
	impl, ok := c.extended()
	if !ok {
		return ErrClipboardNotSupported
	}
	var err error
	c.invoke(func() { err = impl.setHTML(html, plain) })
	return err
}

// HTML returns the HTML representation on the clipboard, or an error when
// there is none.
func (c *Clipboard) HTML() (string, error) {
	impl, ok := c.extended()
	if !ok {
		return "", ErrClipboardNotSupported
	}
	var (
		result string
		err    error
	)
	c.invoke(func() { result, err = impl.html() })
	return result, err
}

// SetRTF replaces the clipboard contents with rich text (RTF bytes).
func (c *Clipboard) SetRTF(rtf []byte) error {
	impl, ok := c.extended()
	if !ok {
		return ErrClipboardNotSupported
	}
	var err error
	c.invoke(func() { err = impl.setRTF(rtf) })
	return err
}

// RTF returns the RTF representation on the clipboard, or an error when
// there is none.
func (c *Clipboard) RTF() ([]byte, error) {
	impl, ok := c.extended()
	if !ok {
		return nil, ErrClipboardNotSupported
	}
	var (
		result []byte
		err    error
	)
	c.invoke(func() { result, err = impl.rtf() })
	return result, err
}

// SetData replaces the clipboard contents with raw bytes under the given
// type identifier (a UTI such as "public.png" or a custom
// "com.example.myapp.record"). Use it for application-private formats.
func (c *Clipboard) SetData(uti string, data []byte) error {
	impl, ok := c.extended()
	if !ok {
		return ErrClipboardNotSupported
	}
	if uti == "" {
		return errors.New("clipboard: a type identifier is required")
	}
	var err error
	c.invoke(func() { err = impl.setData(uti, data) })
	return err
}

// Data returns the raw bytes stored under the given type identifier, or an
// error when the clipboard has no such representation.
func (c *Clipboard) Data(uti string) ([]byte, error) {
	impl, ok := c.extended()
	if !ok {
		return nil, ErrClipboardNotSupported
	}
	if uti == "" {
		return nil, errors.New("clipboard: a type identifier is required")
	}
	var (
		result []byte
		err    error
	)
	c.invoke(func() { result, err = impl.data(uti) })
	return result, err
}

// Types returns the type identifiers currently on the clipboard, most
// preferred first. It is empty on platforms without rich clipboard support.
func (c *Clipboard) Types() []string {
	impl, ok := c.extended()
	if !ok {
		return nil
	}
	var result []string
	c.invoke(func() { result = impl.types() })
	return result
}

// Clear empties the clipboard. On platforms without rich clipboard support
// it sets the text to "".
func (c *Clipboard) Clear() {
	impl, ok := c.extended()
	if !ok {
		c.invoke(func() { c.impl.setText("") })
		return
	}
	c.invoke(impl.clear)
}

// ChangeCount returns a counter that increases every time the clipboard
// contents change, from any application. It returns 0 on platforms without
// rich clipboard support.
func (c *Clipboard) ChangeCount() int {
	impl, ok := c.extended()
	if !ok {
		return 0
	}
	var result int
	c.invoke(func() { result = impl.changeCount() })
	return result
}

// OnChange calls fn whenever the clipboard contents change, including
// changes made by other applications. The returned function removes the
// listener.
//
// There is no system notification for clipboard changes, so the change
// count is polled every 500 ms; the poll runs only while at least one
// listener is registered and stops when the last one is removed. Listeners
// run on their own goroutine. On platforms without rich clipboard support
// the listener is never called and the returned function is a no-op.
func (c *Clipboard) OnChange(fn func()) func() {
	if fn == nil {
		return func() {}
	}
	if _, ok := c.extended(); !ok {
		return func() {}
	}
	c.changeLock.Lock()
	defer c.changeLock.Unlock()
	if c.changeListeners == nil {
		c.changeListeners = make(map[uint]func())
	}
	id := c.nextListenerID
	c.nextListenerID++
	c.changeListeners[id] = fn
	if c.pollStop == nil {
		c.startPollingLocked()
	}
	return func() {
		c.changeLock.Lock()
		defer c.changeLock.Unlock()
		delete(c.changeListeners, id)
		if len(c.changeListeners) == 0 && c.pollStop != nil {
			close(c.pollStop)
			c.pollStop = nil
		}
	}
}

// isPolling reports whether the change-count ticker is running.
func (c *Clipboard) isPolling() bool {
	c.changeLock.Lock()
	defer c.changeLock.Unlock()
	return c.pollStop != nil
}

// startPollingLocked starts the change-count ticker. The caller holds
// changeLock.
func (c *Clipboard) startPollingLocked() {
	stop := make(chan struct{})
	c.pollStop = stop
	interval := c.pollInterval
	if interval <= 0 {
		interval = clipboardPollInterval
	}
	go c.poll(stop, interval)
}

// ready reports whether main-thread dispatch is available: either a test
// dispatcher is installed or the application is running. OnChange may be
// called before App.Run, so the poll waits for this.
func (c *Clipboard) ready() bool {
	if c.dispatch != nil {
		return true
	}
	if globalApplication == nil {
		return false
	}
	globalApplication.runLock.Lock()
	defer globalApplication.runLock.Unlock()
	return globalApplication.running
}

func (c *Clipboard) poll(stop chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for !c.ready() {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
	last := c.ChangeCount()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			count := c.ChangeCount()
			if count == last {
				continue
			}
			last = count
			c.notifyChange()
		}
	}
}

func (c *Clipboard) notifyChange() {
	c.changeLock.Lock()
	listeners := make([]func(), 0, len(c.changeListeners))
	for _, listener := range c.changeListeners {
		listeners = append(listeners, listener)
	}
	c.changeLock.Unlock()
	for _, listener := range listeners {
		go listener()
	}
}

// Rich clipboard access on the manager. Off macOS these return
// ErrClipboardNotSupported (or zero values for Types/ChangeCount) so
// cross-platform code can feature-detect at runtime.

// SetImage replaces the clipboard contents with a PNG image.
func (cm *ClipboardManager) SetImage(png []byte) error { return cm.getClipboard().SetImage(png) }

// Image returns the clipboard image as PNG.
func (cm *ClipboardManager) Image() ([]byte, error) { return cm.getClipboard().Image() }

// SetFiles replaces the clipboard contents with file references.
func (cm *ClipboardManager) SetFiles(paths []string) error {
	return cm.getClipboard().SetFiles(paths)
}

// Files returns the file paths on the clipboard.
func (cm *ClipboardManager) Files() ([]string, error) { return cm.getClipboard().Files() }

// SetHTML replaces the clipboard contents with HTML plus a plain text fallback.
func (cm *ClipboardManager) SetHTML(html string, plain string) error {
	return cm.getClipboard().SetHTML(html, plain)
}

// HTML returns the HTML representation on the clipboard.
func (cm *ClipboardManager) HTML() (string, error) { return cm.getClipboard().HTML() }

// SetRTF replaces the clipboard contents with rich text.
func (cm *ClipboardManager) SetRTF(rtf []byte) error { return cm.getClipboard().SetRTF(rtf) }

// RTF returns the RTF representation on the clipboard.
func (cm *ClipboardManager) RTF() ([]byte, error) { return cm.getClipboard().RTF() }

// SetData replaces the clipboard contents with raw bytes under a type identifier.
func (cm *ClipboardManager) SetData(uti string, data []byte) error {
	return cm.getClipboard().SetData(uti, data)
}

// Data returns the raw bytes stored under a type identifier.
func (cm *ClipboardManager) Data(uti string) ([]byte, error) { return cm.getClipboard().Data(uti) }

// Types returns the type identifiers currently on the clipboard.
func (cm *ClipboardManager) Types() []string { return cm.getClipboard().Types() }

// Clear empties the clipboard.
func (cm *ClipboardManager) Clear() { cm.getClipboard().Clear() }

// ChangeCount returns the clipboard change counter.
func (cm *ClipboardManager) ChangeCount() int { return cm.getClipboard().ChangeCount() }

// OnChange calls fn whenever the clipboard changes; see Clipboard.OnChange.
func (cm *ClipboardManager) OnChange(fn func()) func() { return cm.getClipboard().OnChange(fn) }
