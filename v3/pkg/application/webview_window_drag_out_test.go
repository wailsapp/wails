package application

import (
	"errors"
	"strings"
	"testing"
)

func TestDragItemsValidateRequiresAnItem(t *testing.T) {
	if err := (DragItems{}).validate(); !errors.Is(err, ErrDragItemsEmpty) {
		t.Fatalf("empty items error = %v, want ErrDragItemsEmpty", err)
	}
	if err := (DragItems{Image: []byte{1}, ImageOffset: Point{X: 4}}).validate(); !errors.Is(err, ErrDragItemsEmpty) {
		t.Fatalf("image-only items error = %v, want ErrDragItemsEmpty", err)
	}
	if err := (DragItems{Text: "hello"}).validate(); err != nil {
		t.Fatalf("text items rejected: %v", err)
	}
	if err := (DragItems{Files: []string{"/tmp/a.txt"}}).validate(); err != nil {
		t.Fatalf("file items rejected: %v", err)
	}
	if err := (DragItems{Files: []string{" "}}).validate(); err == nil {
		t.Fatal("blank file path accepted")
	}
}

func TestDragItemsValidatePromises(t *testing.T) {
	data := func() ([]byte, error) { return []byte("x"), nil }
	cases := map[string]DragPromise{
		"missing filename": {Data: data},
		"path separator":   {Filename: "dir/file.txt", Data: data},
		"backslash":        {Filename: `dir\file.txt`, Data: data},
		"missing data":     {Filename: "file.txt"},
	}
	for name, promise := range cases {
		err := (DragItems{Promises: []DragPromise{promise}}).validate()
		if err == nil {
			t.Fatalf("%s: promise accepted", name)
		}
		if !strings.Contains(err.Error(), "promise 0") {
			t.Fatalf("%s: error %q does not name the promise", name, err)
		}
	}
	if err := (DragItems{Promises: []DragPromise{{Filename: "file.txt", Data: data}}}).validate(); err != nil {
		t.Fatalf("valid promise rejected: %v", err)
	}
}

func TestStartDragValidatesBeforeTouchingNative(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{}}
	if err := window.StartDrag(DragItems{}); !errors.Is(err, ErrDragItemsEmpty) {
		t.Fatalf("StartDrag(empty) = %v, want ErrDragItemsEmpty", err)
	}
	var nilWindow *WebviewWindow
	if err := nilWindow.StartDrag(DragItems{Text: "x"}); !errors.Is(err, ErrDragOutWindowNotCreated) {
		t.Fatalf("nil window StartDrag = %v, want ErrDragOutWindowNotCreated", err)
	}
}

func TestDragItemsOperationsDefaultToCopy(t *testing.T) {
	if got := (DragItems{}).operations(); got != DragOperationCopy {
		t.Fatalf("default operations = %v, want copy", got)
	}
	if got := (DragItems{Operations: DragOperationMove | DragOperationLink}).operations(); got != DragOperationMove|DragOperationLink {
		t.Fatalf("explicit operations = %v", got)
	}
}

func TestDragOperationString(t *testing.T) {
	cases := map[DragOperation]string{
		DragOperationNone:                     "none",
		DragOperationCopy:                     "copy",
		DragOperationMove | DragOperationCopy: "copy|move",
		DragOperationDelete:                   "delete",
		DragOperationEvery:                    "every",
		DragOperation(64):                     "DragOperation(64)",
	}
	for op, want := range cases {
		if got := op.String(); got != want {
			t.Fatalf("%d.String() = %q, want %q", uint(op), got, want)
		}
	}
}

func TestUTIForFilename(t *testing.T) {
	cases := map[string]string{
		"photo.PNG":     "public.png",
		"report.pdf":    "com.adobe.pdf",
		"notes.txt":     "public.plain-text",
		"archive.zip":   "public.zip-archive",
		"data.json":     "public.json",
		"mystery.xyz":   "public.data",
		"noextension":   "public.data",
		"page.html":     "public.html",
		"sheet.xlsx":    "org.openxmlformats.spreadsheetml.sheet",
		"clip.mov":      "com.apple.quicktime-movie",
		"song.mp3":      "public.mp3",
		"README.md":     "net.daringfireball.markdown",
		"picture.jpeg":  "public.jpeg",
		"animation.gif": "com.compuserve.gif",
	}
	for name, want := range cases {
		if got := utiForFilename(name); got != want {
			t.Fatalf("utiForFilename(%q) = %q, want %q", name, got, want)
		}
	}
	if got := (DragPromise{Filename: "a.png", UTI: "com.example.custom"}).uti(); got != "com.example.custom" {
		t.Fatalf("explicit UTI ignored: %q", got)
	}
	if got := (DragPromise{Filename: "a.png"}).uti(); got != "public.png" {
		t.Fatalf("derived UTI = %q", got)
	}
}

func TestOnDragEndListenerBookkeeping(t *testing.T) {
	window := &WebviewWindow{id: 9001}
	if dragOutListenerCount(window.id) != 0 {
		t.Fatal("listeners registered before OnDragEnd")
	}
	noop := window.OnDragEnd(nil)
	noop()
	if dragOutListenerCount(window.id) != 0 {
		t.Fatal("OnDragEnd(nil) registered a listener")
	}
	got := make(chan DragOperation, 2)
	first := window.OnDragEnd(func(op DragOperation) { got <- op })
	second := window.OnDragEnd(func(op DragOperation) { got <- op })
	if dragOutListenerCount(window.id) != 2 {
		t.Fatalf("listener count = %d, want 2", dragOutListenerCount(window.id))
	}
	dispatchDragEnd(window.id, DragOperationMove)
	for i := 0; i < 2; i++ {
		if op := <-got; op != DragOperationMove {
			t.Fatalf("listener received %v", op)
		}
	}
	first()
	first()
	if dragOutListenerCount(window.id) != 1 {
		t.Fatalf("listener count after one cancel = %d, want 1", dragOutListenerCount(window.id))
	}
	second()
	if dragOutListenerCount(window.id) != 0 {
		t.Fatal("listeners remain after both cancelled")
	}
	dispatchDragEnd(window.id, DragOperationCopy)
	select {
	case op := <-got:
		t.Fatalf("cancelled listener received %v", op)
	default:
	}
	dragOutListenersLock.Lock()
	delete(dragOutListeners, window.id)
	dragOutListenersLock.Unlock()
}

func TestEffectiveDropTypesDefaultToFiles(t *testing.T) {
	if got := effectiveDropTypes(nil); len(got) != 1 || got[0] != DropFiles {
		t.Fatalf("effectiveDropTypes(nil) = %v, want [DropFiles]", got)
	}
	if got := effectiveDropTypes([]DropType{DropType(0), DropType(99)}); len(got) != 1 || got[0] != DropFiles {
		t.Fatalf("unknown drop types = %v, want [DropFiles]", got)
	}
	got := effectiveDropTypes([]DropType{DropText, DropText, DropURLs})
	if len(got) != 2 || got[0] != DropText || got[1] != DropURLs {
		t.Fatalf("effectiveDropTypes deduplicated = %v", got)
	}
	if mask := dropTypeMask(nil); mask != dropMaskFiles {
		t.Fatalf("default mask = %d, want %d", mask, dropMaskFiles)
	}
	if mask := dropTypeMask([]DropType{DropFiles, DropText, DropURLs, DropImages}); mask != dropMaskFiles|dropMaskText|dropMaskURLs|dropMaskImages {
		t.Fatalf("full mask = %d", mask)
	}
	if mask := dropTypeMask([]DropType{DropImages}); mask != dropMaskImages {
		t.Fatalf("images-only mask = %d, want %d", mask, dropMaskImages)
	}
	if dropTypesNeedOverlay(nil) || dropTypesNeedOverlay([]DropType{DropFiles}) {
		t.Fatal("file-only drop types reported as needing the overlay")
	}
	if !dropTypesNeedOverlay([]DropType{DropText}) {
		t.Fatal("text drops do not report needing the overlay")
	}
	for i, want := range []string{"DropFiles", "DropText", "DropURLs", "DropImages"} {
		if got := DropType(i + 1).String(); got != want {
			t.Fatalf("DropType(%d).String() = %q, want %q", i+1, got, want)
		}
	}
	if got := DropType(7).String(); got != "DropType(7)" {
		t.Fatalf("unknown DropType string = %q", got)
	}
}

func TestNewWindowEnablesOverlayForNonFileDropTypes(t *testing.T) {
	window := NewWindow(WebviewWindowOptions{DropTypes: []DropType{DropURLs}})
	if !window.options.EnableFileDrop {
		t.Fatal("DropURLs did not enable the drop overlay")
	}
	window = NewWindow(WebviewWindowOptions{DropTypes: []DropType{DropFiles}})
	if window.options.EnableFileDrop {
		t.Fatal("DropFiles alone enabled the drop overlay")
	}
}

func TestOnDropListenerBookkeeping(t *testing.T) {
	window := &WebviewWindow{id: 9002}
	if dispatchDrop(window.id, DropData{Text: "x"}) {
		t.Fatal("dispatchDrop reported listeners before any were added")
	}
	got := make(chan DropData, 1)
	cancel := window.OnDrop(func(ctx *Context, data DropData) {
		if ctx == nil {
			t.Error("nil context")
		}
		got <- data
	})
	if dropListenerCount(window.id) != 1 {
		t.Fatalf("listener count = %d", dropListenerCount(window.id))
	}
	if !dispatchDrop(window.id, DropData{Text: "hello", URLs: []string{"https://wails.io"}, X: 3, Y: 4}) {
		t.Fatal("dispatchDrop reported no listeners")
	}
	data := <-got
	if data.Text != "hello" || len(data.URLs) != 1 || data.X != 3 || data.Y != 4 {
		t.Fatalf("listener received %+v", data)
	}
	if data.IsEmpty() {
		t.Fatal("populated DropData reported empty")
	}
	if !(DropData{}).IsEmpty() {
		t.Fatal("zero DropData not reported empty")
	}
	cancel()
	if dropListenerCount(window.id) != 0 {
		t.Fatal("listener remains after cancel")
	}
	dropListenersLock.Lock()
	delete(dropListeners, window.id)
	dropListenersLock.Unlock()
}
