package application

import (
	"sync"
	"sync/atomic"
)

// MacTextEditor is a native plain-text editor backed by NSTextView inside an
// NSScrollView. It can be used as the primary content of a MacSplitView hosted
// by NativeWindow. The API is experimental in v3.
type MacTextEditor struct {
	lock sync.RWMutex

	internalID  uint64
	text        string
	textVersion uint64
	editable    bool
	pane        *MacSplitPane
	dead        bool
	onChange    func(*Context)
}

var macTextEditorID uint64

// NewMacTextEditor creates an editable, plain-text native editor.
func NewMacTextEditor() *MacTextEditor {
	return &MacTextEditor{
		internalID: atomic.AddUint64(&macTextEditorID, 1),
		editable:   true,
	}
}

// SetText replaces the editor contents. Programmatic changes do not invoke
// OnChange, which makes loading another file safe without marking it dirty.
func (e *MacTextEditor) SetText(text string) *MacTextEditor {
	if e == nil {
		return e
	}
	e.lock.Lock()
	if e.dead {
		e.lock.Unlock()
		return e
	}
	e.text = text
	e.textVersion++
	version := e.textVersion
	e.lock.Unlock()
	if macTextEditorApplyText(e, text) {
		e.clearCachedText(version)
	}
	return e
}

// Text returns the current complete plain-text contents. When attached, the
// read is dispatched synchronously to AppKit's main thread.
func (e *MacTextEditor) Text() string {
	if e == nil {
		return ""
	}
	if text, ok := macTextEditorReadText(e); ok {
		return text
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.text
}

// SetEditable controls whether the user can modify the text.
func (e *MacTextEditor) SetEditable(editable bool) *MacTextEditor {
	if e == nil {
		return e
	}
	e.lock.Lock()
	if e.dead {
		e.lock.Unlock()
		return e
	}
	e.editable = editable
	e.lock.Unlock()
	macTextEditorApplyEditable(e, editable)
	return e
}

// OnChange sets the callback fired for user edits. The callback is a signal;
// call Text when the full document is actually needed (for example, Save) so
// large documents are not copied across the bridge on every keystroke.
func (e *MacTextEditor) OnChange(callback func(*Context)) *MacTextEditor {
	if e == nil {
		return e
	}
	e.lock.Lock()
	if !e.dead {
		e.onChange = callback
	}
	e.lock.Unlock()
	return e
}

// Focus makes the NSTextView first responder.
func (e *MacTextEditor) Focus() *MacTextEditor {
	if e != nil {
		macTextEditorFocus(e)
	}
	return e
}

func (e *MacTextEditor) snapshot() (id uint64, text string, editable bool, textVersion uint64) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.internalID, e.text, e.editable, e.textVersion
}

// clearCachedText drops the staging copy once AppKit owns the same version.
// Before installation, text remains cached so it can seed the NSTextView.
func (e *MacTextEditor) clearCachedText(version uint64) {
	e.lock.Lock()
	if e.textVersion == version {
		e.text = ""
	}
	e.lock.Unlock()
}

func (e *MacTextEditor) markDead() {
	if e == nil {
		return
	}
	unregisterMacTextEditor(e.internalID)
	e.lock.Lock()
	e.dead = true
	e.pane = nil
	e.onChange = nil
	e.lock.Unlock()
}

var macTextEditorRegistry = make(map[uint64]*MacTextEditor)
var macTextEditorRegistryLock sync.RWMutex
var macTextEditorChanged = make(chan uint64, 64)

func registerMacTextEditor(editor *MacTextEditor) {
	if editor == nil {
		return
	}
	macTextEditorRegistryLock.Lock()
	macTextEditorRegistry[editor.internalID] = editor
	macTextEditorRegistryLock.Unlock()
}

func unregisterMacTextEditor(id uint64) {
	macTextEditorRegistryLock.Lock()
	delete(macTextEditorRegistry, id)
	macTextEditorRegistryLock.Unlock()
}

func handleMacTextEditorChanged(id uint64) {
	defer handlePanic()
	macTextEditorRegistryLock.RLock()
	editor := macTextEditorRegistry[id]
	macTextEditorRegistryLock.RUnlock()
	if editor == nil {
		return
	}
	editor.lock.RLock()
	callback := editor.onChange
	dead := editor.dead
	editor.lock.RUnlock()
	if !dead && callback != nil {
		callback(newContext())
	}
}
