package application

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// plainClipboard only implements clipboardImpl, like Windows and Linux.
type plainClipboard struct {
	value string
}

func (p *plainClipboard) setText(text string) bool { p.value = text; return true }
func (p *plainClipboard) text() (string, bool)     { return p.value, true }

// richClipboard implements clipboardExtendedImpl in memory so the manager
// logic can be tested without a pasteboard.
type richClipboard struct {
	plainClipboard
	mu      sync.Mutex
	store   map[string][]byte
	count   int64
	polls   int64
	cleared int
}

func newRichClipboard() *richClipboard {
	return &richClipboard{store: map[string][]byte{}}
}

func (r *richClipboard) bump() { r.count++ }

func (r *richClipboard) setImage(png []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = map[string][]byte{"public.png": png}
	r.bump()
	return nil
}

func (r *richClipboard) image() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if png, ok := r.store["public.png"]; ok {
		return png, nil
	}
	return nil, errors.New("no image")
}

func (r *richClipboard) setFiles(paths []string) error { r.bump(); return nil }
func (r *richClipboard) files() ([]string, error)      { return nil, nil }
func (r *richClipboard) setHTML(html, plain string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = map[string][]byte{"public.html": []byte(html)}
	if plain != "" {
		r.store["public.utf8-plain-text"] = []byte(plain)
	}
	r.bump()
	return nil
}
func (r *richClipboard) html() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if html, ok := r.store["public.html"]; ok {
		return string(html), nil
	}
	return "", errors.New("no html")
}
func (r *richClipboard) setRTF(rtf []byte) error { return r.setData("public.rtf", rtf) }
func (r *richClipboard) rtf() ([]byte, error)    { return r.data("public.rtf") }
func (r *richClipboard) setData(uti string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = map[string][]byte{uti: data}
	r.bump()
	return nil
}
func (r *richClipboard) data(uti string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if data, ok := r.store[uti]; ok {
		return data, nil
	}
	return nil, errors.New("no data for " + uti)
}
func (r *richClipboard) types() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]string, 0, len(r.store))
	for uti := range r.store {
		result = append(result, uti)
	}
	return result
}
func (r *richClipboard) clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = map[string][]byte{}
	r.cleared++
	r.bump()
}
func (r *richClipboard) changeCount() int {
	atomic.AddInt64(&r.polls, 1)
	r.mu.Lock()
	defer r.mu.Unlock()
	return int(r.count)
}

// newTestClipboard wires a clipboard that runs "main thread" work inline.
func newTestClipboard(impl clipboardImpl) *Clipboard {
	return &Clipboard{
		impl:         impl,
		dispatch:     func(fn func()) { fn() },
		pollInterval: 5 * time.Millisecond,
	}
}

func TestClipboardRichMethodsReportNotSupportedWithoutExtendedImpl(t *testing.T) {
	clipboard := newTestClipboard(&plainClipboard{})
	if err := clipboard.SetImage([]byte{1}); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("SetImage error = %v, want ErrClipboardNotSupported", err)
	}
	if _, err := clipboard.Image(); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("Image error = %v, want ErrClipboardNotSupported", err)
	}
	if err := clipboard.SetFiles([]string{"/tmp/a"}); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("SetFiles error = %v", err)
	}
	if _, err := clipboard.Files(); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("Files error = %v", err)
	}
	if err := clipboard.SetHTML("<b>x</b>", "x"); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("SetHTML error = %v", err)
	}
	if _, err := clipboard.HTML(); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("HTML error = %v", err)
	}
	if err := clipboard.SetRTF([]byte("{\\rtf1}")); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("SetRTF error = %v", err)
	}
	if _, err := clipboard.RTF(); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("RTF error = %v", err)
	}
	if err := clipboard.SetData("com.example.x", []byte{1}); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("SetData error = %v", err)
	}
	if _, err := clipboard.Data("com.example.x"); !errors.Is(err, ErrClipboardNotSupported) {
		t.Fatalf("Data error = %v", err)
	}
	if types := clipboard.Types(); types != nil {
		t.Fatalf("Types() = %v, want nil", types)
	}
	if count := clipboard.ChangeCount(); count != 0 {
		t.Fatalf("ChangeCount() = %d, want 0", count)
	}
	// Clear falls back to setting empty text so it still empties the clipboard.
	plain := clipboard.impl.(*plainClipboard)
	plain.value = "something"
	clipboard.Clear()
	if plain.value != "" {
		t.Fatalf("Clear() left %q on the plain clipboard", plain.value)
	}
	cancel := clipboard.OnChange(func() {})
	if clipboard.isPolling() {
		t.Fatal("OnChange started polling on a plain clipboard")
	}
	cancel()
}

func TestClipboardRichMethodsUseExtendedImpl(t *testing.T) {
	rich := newRichClipboard()
	clipboard := newTestClipboard(rich)

	if err := clipboard.SetImage(nil); err == nil {
		t.Fatal("SetImage(nil) should fail")
	}
	if err := clipboard.SetImage([]byte("png")); err != nil {
		t.Fatalf("SetImage: %v", err)
	}
	png, err := clipboard.Image()
	if err != nil || string(png) != "png" {
		t.Fatalf("Image() = %q, %v", png, err)
	}
	if err := clipboard.SetHTML("<b>hi</b>", "hi"); err != nil {
		t.Fatalf("SetHTML: %v", err)
	}
	if html, err := clipboard.HTML(); err != nil || html != "<b>hi</b>" {
		t.Fatalf("HTML() = %q, %v", html, err)
	}
	if err := clipboard.SetData("", []byte{1}); err == nil {
		t.Fatal("SetData with an empty type should fail")
	}
	if _, err := clipboard.Data(""); err == nil {
		t.Fatal("Data with an empty type should fail")
	}
	if err := clipboard.SetData("com.example.record", []byte{1, 2, 3}); err != nil {
		t.Fatalf("SetData: %v", err)
	}
	data, err := clipboard.Data("com.example.record")
	if err != nil || len(data) != 3 {
		t.Fatalf("Data() = %v, %v", data, err)
	}
	if types := clipboard.Types(); len(types) != 1 || types[0] != "com.example.record" {
		t.Fatalf("Types() = %v", types)
	}
	if err := clipboard.SetFiles(nil); err == nil {
		t.Fatal("SetFiles(nil) should fail")
	}
	before := clipboard.ChangeCount()
	clipboard.Clear()
	if rich.cleared != 1 {
		t.Fatalf("Clear() called the impl %d times", rich.cleared)
	}
	if after := clipboard.ChangeCount(); after != before+1 {
		t.Fatalf("ChangeCount() = %d after clear, want %d", after, before+1)
	}
}

func TestClipboardOnChangeStartsAndStopsPollingWithListeners(t *testing.T) {
	rich := newRichClipboard()
	clipboard := newTestClipboard(rich)

	if clipboard.isPolling() {
		t.Fatal("polling before any listener")
	}
	if cancel := clipboard.OnChange(nil); cancel == nil {
		t.Fatal("OnChange(nil) must return a cancel function")
	} else {
		cancel()
	}
	if clipboard.isPolling() {
		t.Fatal("OnChange(nil) started polling")
	}

	var fired int32
	first := clipboard.OnChange(func() { atomic.AddInt32(&fired, 1) })
	if !clipboard.isPolling() {
		t.Fatal("first listener did not start polling")
	}
	second := clipboard.OnChange(func() { atomic.AddInt32(&fired, 1) })
	if !clipboard.isPolling() {
		t.Fatal("polling stopped while listeners remain")
	}

	// A change to the pasteboard is seen by the poll and fans out to both
	// listeners.
	waitFor(t, func() bool { return atomic.LoadInt64(&rich.polls) > 0 }, "initial poll")
	rich.mu.Lock()
	rich.bump()
	rich.mu.Unlock()
	waitFor(t, func() bool { return atomic.LoadInt32(&fired) == 2 }, "listeners notified")

	// Removing one listener keeps the ticker; removing the last stops it.
	first()
	if !clipboard.isPolling() {
		t.Fatal("polling stopped with one listener left")
	}
	second()
	if clipboard.isPolling() {
		t.Fatal("polling did not stop after the last listener was removed")
	}
	// Cancelling twice is harmless.
	second()
	if clipboard.isPolling() {
		t.Fatal("double cancel restarted polling")
	}

	// The goroutine has actually stopped: the poll counter settles.
	polls := atomic.LoadInt64(&rich.polls)
	time.Sleep(30 * time.Millisecond)
	if again := atomic.LoadInt64(&rich.polls); again != polls {
		t.Fatalf("poll goroutine still running: %d -> %d", polls, again)
	}

	// A new listener starts a fresh ticker.
	third := clipboard.OnChange(func() {})
	if !clipboard.isPolling() {
		t.Fatal("re-adding a listener did not restart polling")
	}
	third()
}

func TestClipboardManagerDelegatesRichMethods(t *testing.T) {
	rich := newRichClipboard()
	manager := &ClipboardManager{clipboard: newTestClipboard(rich)}
	if err := manager.SetImage([]byte("x")); err != nil {
		t.Fatalf("manager SetImage: %v", err)
	}
	if png, err := manager.Image(); err != nil || string(png) != "x" {
		t.Fatalf("manager Image() = %q, %v", png, err)
	}
	if manager.ChangeCount() != 1 {
		t.Fatalf("manager ChangeCount() = %d", manager.ChangeCount())
	}
	cancel := manager.OnChange(func() {})
	if !manager.clipboard.isPolling() {
		t.Fatal("manager OnChange did not start polling")
	}
	cancel()
	if manager.clipboard.isPolling() {
		t.Fatal("manager cancel did not stop polling")
	}
}

func waitFor(t *testing.T, condition func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}
