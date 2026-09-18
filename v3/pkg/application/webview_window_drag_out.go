package application

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// DragOperation is the operation a drop target performed on a drag started
// with WebviewWindow.StartDrag, and the set of operations a drag source
// allows. The values match NSDragOperation so they can be combined with |.
type DragOperation uint

const (
	// DragOperationNone means the drag was cancelled or rejected.
	DragOperationNone DragOperation = 0
	// DragOperationCopy copies the data to the target.
	DragOperationCopy DragOperation = 1
	// DragOperationLink creates a link or alias at the target.
	DragOperationLink DragOperation = 2
	// DragOperationGeneric lets the target choose the operation.
	DragOperationGeneric DragOperation = 4
	// DragOperationPrivate is an operation the source and target agree on
	// privately.
	DragOperationPrivate DragOperation = 8
	// DragOperationMove moves the data to the target.
	DragOperationMove DragOperation = 16
	// DragOperationDelete means the data was dropped on the Trash.
	DragOperationDelete DragOperation = 32
	// DragOperationEvery allows any operation.
	DragOperationEvery DragOperation = 0xFFFFFFFF
)

// String returns a readable name for a single operation, or the combined
// names for a mask.
func (o DragOperation) String() string {
	if o == DragOperationNone {
		return "none"
	}
	if o == DragOperationEvery {
		return "every"
	}
	names := []string{}
	add := func(flag DragOperation, name string) {
		if o&flag != 0 {
			names = append(names, name)
		}
	}
	add(DragOperationCopy, "copy")
	add(DragOperationLink, "link")
	add(DragOperationGeneric, "generic")
	add(DragOperationPrivate, "private")
	add(DragOperationMove, "move")
	add(DragOperationDelete, "delete")
	if len(names) == 0 {
		return fmt.Sprintf("DragOperation(%d)", uint(o))
	}
	return strings.Join(names, "|")
}

// DragPromise describes a file that does not exist yet. The destination
// (Finder, Mail, another app) receives a file promise; when it accepts the
// drop, Data is called on a background goroutine and the bytes are written
// to the destination path. This is how apps drag out downloads or exports
// without writing temp files first.
type DragPromise struct {
	// Filename is the name the destination writes, including extension.
	Filename string
	// UTI is the uniform type identifier of the file, such as "public.png".
	// Leave it empty to derive it from the Filename extension.
	UTI string
	// Data produces the file contents once the drop has been accepted.
	Data func() ([]byte, error)
}

// DragItems is what WebviewWindow.StartDrag hands to the system. At least
// one of Files, Promises or Text must be set. When several are set they are
// offered together and the destination takes what it understands.
type DragItems struct {
	// Files are absolute paths of existing files to drag out.
	Files []string
	// Promises are files produced on demand; see DragPromise.
	Promises []DragPromise
	// Text is plain text to drag out.
	Text string
	// Image is a PNG shown under the cursor while dragging. When empty a
	// file icon (files and promises) or the rendered text is used.
	Image []byte
	// ImageOffset is the point inside the drag image, measured from its top
	// left corner, that sits under the cursor. Zero puts the top left corner
	// under the cursor.
	ImageOffset Point
	// Operations is the set of operations the drop target may perform. Zero
	// means DragOperationCopy.
	Operations DragOperation
}

var (
	// ErrDragItemsEmpty means StartDrag was given no files, promises or
	// text.
	ErrDragItemsEmpty = errors.New("drag out: at least one file, promise or text item is required")
	// ErrDragOutNoGesture means StartDrag was called while no mouse button
	// was held down over the window. Call it from a Go method the page
	// invokes in its mousedown or dragstart handler.
	ErrDragOutNoGesture = errors.New("drag out: no mouse drag gesture is in progress; call StartDrag from a mousedown or dragstart handler")
	// ErrDragOutWindowNotCreated means the native window does not exist
	// yet.
	ErrDragOutWindowNotCreated = errors.New("drag out: the native window has not been created yet")
	// ErrDragOutUnsupported means drag out is unavailable on this
	// platform.
	ErrDragOutUnsupported = errors.New("drag out is not supported on this platform")
)

// validate checks the items before anything native is touched.
func (d DragItems) validate() error {
	if len(d.Files) == 0 && len(d.Promises) == 0 && d.Text == "" {
		return ErrDragItemsEmpty
	}
	for i, file := range d.Files {
		if strings.TrimSpace(file) == "" {
			return fmt.Errorf("drag out: file %d has an empty path", i)
		}
	}
	for i, promise := range d.Promises {
		if strings.TrimSpace(promise.Filename) == "" {
			return fmt.Errorf("drag out: promise %d needs a Filename", i)
		}
		if strings.ContainsAny(promise.Filename, `/\`) {
			return fmt.Errorf("drag out: promise %d Filename %q must not contain path separators", i, promise.Filename)
		}
		if promise.Data == nil {
			return fmt.Errorf("drag out: promise %d needs a Data function", i)
		}
	}
	return nil
}

// operations returns the effective operation mask.
func (d DragItems) operations() DragOperation {
	if d.Operations == DragOperationNone {
		return DragOperationCopy
	}
	return d.Operations
}

// utiByExtension maps common file extensions to their uniform type
// identifiers. Unknown extensions fall back to public.data, which every
// destination accepts as an opaque file.
var utiByExtension = map[string]string{
	".png":    "public.png",
	".jpg":    "public.jpeg",
	".jpeg":   "public.jpeg",
	".gif":    "com.compuserve.gif",
	".tif":    "public.tiff",
	".tiff":   "public.tiff",
	".bmp":    "com.microsoft.bmp",
	".webp":   "org.webmproject.webp",
	".svg":    "public.svg-image",
	".heic":   "public.heic",
	".pdf":    "com.adobe.pdf",
	".txt":    "public.plain-text",
	".md":     "net.daringfireball.markdown",
	".html":   "public.html",
	".htm":    "public.html",
	".xml":    "public.xml",
	".json":   "public.json",
	".csv":    "public.comma-separated-values-text",
	".rtf":    "public.rtf",
	".zip":    "public.zip-archive",
	".gz":     "org.gnu.gnu-zip-archive",
	".tar":    "public.tar-archive",
	".mp3":    "public.mp3",
	".m4a":    "com.apple.m4a-audio",
	".wav":    "com.microsoft.waveform-audio",
	".mp4":    "public.mpeg-4",
	".mov":    "com.apple.quicktime-movie",
	".js":     "com.netscape.javascript-source",
	".go":     "public.source-code",
	".py":     "public.python-script",
	".sh":     "public.shell-script",
	".yaml":   "public.yaml",
	".yml":    "public.yaml",
	".toml":   "public.data",
	".log":    "public.plain-text",
	".dmg":    "com.apple.disk-image-udif",
	".pkg":    "com.apple.installer-package-archive",
	".app":    "com.apple.application-bundle",
	".docx":   "org.openxmlformats.wordprocessingml.document",
	".xlsx":   "org.openxmlformats.spreadsheetml.sheet",
	".pptx":   "org.openxmlformats.presentationml.presentation",
	".epub":   "org.idpf.epub-container",
	".ics":    "com.apple.ical.ics",
	".vcf":    "public.vcard",
	".sqlite": "public.data",
}

// utiForFilename returns the uniform type identifier for a filename based on
// its extension, or public.data when the extension is unknown.
func utiForFilename(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if uti, ok := utiByExtension[ext]; ok {
		return uti
	}
	return "public.data"
}

// uti returns the promise's type identifier, deriving it from the filename
// when it was not given.
func (p DragPromise) uti() string {
	if strings.TrimSpace(p.UTI) != "" {
		return p.UTI
	}
	return utiForFilename(p.Filename)
}

// dragOutListeners holds OnDragEnd callbacks per window. The WebviewWindow
// struct stays untouched; entries are removed when the window closes.
var (
	dragOutListenersLock sync.Mutex
	dragOutListeners     = map[uint]map[uint]func(DragOperation){}
	dragOutListenerSeq   uint
)

// StartDrag begins a system drag from this window with the given items, as
// if the user had picked them up in Finder. Files land in the destination
// as copies, promises are written on demand, and text is offered as a plain
// string. It returns as soon as the drag session has been started; use
// OnDragEnd to learn what the destination did.
//
// The drag must be started during a mouse gesture: the OS attaches the
// session to the mouse-down or mouse-dragged event that is being processed.
// Bind a Go method and call it from the page's mousedown, pointerdown or
// dragstart handler on the draggable element (set the element's HTML5
// draggable attribute to false so WebKit does not start its own drag). When
// no mouse button is down over the window it returns ErrDragOutNoGesture.
//
// Drag out is implemented on macOS. Other platforms return
// ErrDragOutUnsupported.
func (w *WebviewWindow) StartDrag(items DragItems) error {
	if w == nil {
		return ErrDragOutWindowNotCreated
	}
	if err := items.validate(); err != nil {
		return err
	}
	return w.startDragOut(items)
}

// OnDragEnd calls fn when a drag started with StartDrag finishes, with the
// operation the destination performed (DragOperationNone when it was
// cancelled). The returned function removes the listener. Listeners run on
// their own goroutine.
func (w *WebviewWindow) OnDragEnd(fn func(operation DragOperation)) func() {
	if w == nil || fn == nil {
		return func() {}
	}
	dragOutListenersLock.Lock()
	listeners, ok := dragOutListeners[w.id]
	if !ok {
		listeners = map[uint]func(DragOperation){}
		dragOutListeners[w.id] = listeners
		if w.eventListeners != nil {
			w.OnWindowEvent(events.Common.WindowClosing, func(*WindowEvent) {
				dragOutListenersLock.Lock()
				delete(dragOutListeners, w.id)
				dragOutListenersLock.Unlock()
			})
		}
	}
	dragOutListenerSeq++
	id := dragOutListenerSeq
	listeners[id] = fn
	dragOutListenersLock.Unlock()
	return func() {
		dragOutListenersLock.Lock()
		defer dragOutListenersLock.Unlock()
		if listeners, ok := dragOutListeners[w.id]; ok {
			delete(listeners, id)
		}
	}
}

// dragOutListenerCount reports how many OnDragEnd listeners a window has.
func dragOutListenerCount(windowID uint) int {
	dragOutListenersLock.Lock()
	defer dragOutListenersLock.Unlock()
	return len(dragOutListeners[windowID])
}

// dispatchDragEnd delivers a finished drag to the window's listeners.
func dispatchDragEnd(windowID uint, operation DragOperation) {
	dragOutListenersLock.Lock()
	listeners := make([]func(DragOperation), 0, len(dragOutListeners[windowID]))
	for _, listener := range dragOutListeners[windowID] {
		listeners = append(listeners, listener)
	}
	dragOutListenersLock.Unlock()
	for _, listener := range listeners {
		go listener(operation)
	}
}

// DropType selects which kinds of content dragged from other applications
// a window accepts. Files use the existing file drop pipeline (the
// WindowFilesDropped event and the runtime's data-file-drop-target
// handling); the other types are delivered to Go through OnDrop.
type DropType int

const (
	// DropFiles accepts files and folders. This is the default.
	DropFiles DropType = iota + 1
	// DropText accepts plain text, for example a selection dragged from
	// another app.
	DropText
	// DropURLs accepts URLs, for example a link dragged from a browser.
	DropURLs
	// DropImages accepts bitmap images dragged from another app (as PNG).
	DropImages
)

// String returns the option name.
func (t DropType) String() string {
	switch t {
	case DropFiles:
		return "DropFiles"
	case DropText:
		return "DropText"
	case DropURLs:
		return "DropURLs"
	case DropImages:
		return "DropImages"
	default:
		return fmt.Sprintf("DropType(%d)", int(t))
	}
}

// Drop type bits shared with the native drag view.
const (
	dropMaskFiles  = 1
	dropMaskText   = 2
	dropMaskURLs   = 4
	dropMaskImages = 8
)

// effectiveDropTypes returns the drop types a window accepts: the ones
// configured, or DropFiles when none were.
func effectiveDropTypes(types []DropType) []DropType {
	if len(types) == 0 {
		return []DropType{DropFiles}
	}
	result := make([]DropType, 0, len(types))
	seen := map[DropType]bool{}
	for _, t := range types {
		if t < DropFiles || t > DropImages || seen[t] {
			continue
		}
		seen[t] = true
		result = append(result, t)
	}
	if len(result) == 0 {
		return []DropType{DropFiles}
	}
	return result
}

// dropTypeMask converts drop types to the bitmask the native drag view
// registers pasteboard types from.
func dropTypeMask(types []DropType) int {
	mask := 0
	for _, t := range effectiveDropTypes(types) {
		switch t {
		case DropFiles:
			mask |= dropMaskFiles
		case DropText:
			mask |= dropMaskText
		case DropURLs:
			mask |= dropMaskURLs
		case DropImages:
			mask |= dropMaskImages
		}
	}
	return mask
}

// dropTypesNeedOverlay reports whether the configured drop types include
// anything other than files, which needs the drop overlay even when
// EnableFileDrop was left off.
func dropTypesNeedOverlay(types []DropType) bool {
	return dropTypeMask(types)&^dropMaskFiles != 0
}

// DropData is what OnDrop receives for a non-file drop. Only the fields
// matching the window's DropTypes are filled; a drag can carry several
// representations at once (a browser link is both text and a URL).
type DropData struct {
	// Files are the dropped file paths. They are only set when the drop
	// carried both files and a requested non-file type; plain file drops go
	// through the WindowFilesDropped event instead.
	Files []string
	// Text is the plain text representation.
	Text string
	// URLs are the dropped URLs as strings.
	URLs []string
	// Images are the dropped bitmaps as PNG.
	Images [][]byte
	// X and Y are the drop point in window content coordinates, origin top
	// left, in CSS pixels.
	X int
	Y int
}

// IsEmpty reports whether nothing usable was delivered.
func (d DropData) IsEmpty() bool {
	return len(d.Files) == 0 && d.Text == "" && len(d.URLs) == 0 && len(d.Images) == 0
}

// dropListeners holds OnDrop callbacks per window.
var (
	dropListenersLock sync.Mutex
	dropListeners     = map[uint]map[uint]func(*Context, DropData){}
	dropListenerSeq   uint
)

// OnDrop calls fn when text, URLs or images are dropped on the window from
// another application. The window must be created with EnableFileDrop and
// the wanted types listed in DropTypes; without DropText, DropURLs or
// DropImages nothing is delivered here. File drops keep using the
// WindowFilesDropped event. The returned function removes the listener.
// Listeners run on their own goroutine.
//
// Non-file drops are implemented on macOS.
func (w *WebviewWindow) OnDrop(fn func(ctx *Context, data DropData)) func() {
	if w == nil || fn == nil {
		return func() {}
	}
	dropListenersLock.Lock()
	listeners, ok := dropListeners[w.id]
	if !ok {
		listeners = map[uint]func(*Context, DropData){}
		dropListeners[w.id] = listeners
		if w.eventListeners != nil {
			w.OnWindowEvent(events.Common.WindowClosing, func(*WindowEvent) {
				dropListenersLock.Lock()
				delete(dropListeners, w.id)
				dropListenersLock.Unlock()
			})
		}
	}
	dropListenerSeq++
	id := dropListenerSeq
	listeners[id] = fn
	dropListenersLock.Unlock()
	return func() {
		dropListenersLock.Lock()
		defer dropListenersLock.Unlock()
		if listeners, ok := dropListeners[w.id]; ok {
			delete(listeners, id)
		}
	}
}

// dropListenerCount reports how many OnDrop listeners a window has.
func dropListenerCount(windowID uint) int {
	dropListenersLock.Lock()
	defer dropListenersLock.Unlock()
	return len(dropListeners[windowID])
}

// dispatchDrop delivers a non-file drop to the window's listeners and
// reports whether anyone was listening.
func dispatchDrop(windowID uint, data DropData) bool {
	dropListenersLock.Lock()
	listeners := make([]func(*Context, DropData), 0, len(dropListeners[windowID]))
	for _, listener := range dropListeners[windowID] {
		listeners = append(listeners, listener)
	}
	dropListenersLock.Unlock()
	for _, listener := range listeners {
		go listener(newContext(), data)
	}
	return len(listeners) > 0
}
