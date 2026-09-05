package application

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type panelProbe struct {
	webviewPanelImpl
	destroyed atomic.Int32
	focused   atomic.Int32
	zoom      float64
}

func (p *panelProbe) destroy()             { p.destroyed.Add(1) }
func (p *panelProbe) isFocused() bool      { p.focused.Add(1); return true }
func (p *panelProbe) setZoom(zoom float64) { p.zoom = zoom }

func panelTestApp(t *testing.T, main bool) *threadProbeApp {
	t.Helper()
	old := globalApplication
	probe := &threadProbeApp{}
	probe.onMain.Store(main)
	globalApplication = &App{impl: probe}
	t.Cleanup(func() { globalApplication = old })
	return probe
}

func TestPanelRuntimeFocusDoesNotDestroy(t *testing.T) {
	panelTestApp(t, true)
	window := &WebviewWindow{}
	panel := window.NewPanel(WebviewPanelOptions{Name: "browser"})
	native := &panelProbe{}
	panel.impl = native
	processor := &MessageProcessor{}
	// Literal IDs are the wire contract, independent of the Go constants.
	result, err := processor.processPanelMethod(windowMethodRequest(t, 14, `{"panel":"browser"}`), window)
	if err != nil || result != true {
		t.Fatalf("IsFocused: %v, %v", result, err)
	}
	if native.destroyed.Load() != 0 || window.GetPanel("browser") != panel {
		t.Fatal("focus query destroyed panel")
	}
	if _, err := processor.processPanelMethod(windowMethodRequest(t, 11, `{"panel":"browser","zoom":1.5}`), window); err != nil {
		t.Fatal(err)
	}
	if native.zoom != 1.5 {
		t.Fatalf("SetZoom dispatched incorrectly: %v", native.zoom)
	}
	if _, err := processor.processPanelMethod(windowMethodRequest(t, 16, `{"panel":"browser"}`), window); err != nil {
		t.Fatal(err)
	}
	if native.destroyed.Load() != 1 || window.GetPanel("browser") != nil {
		t.Fatal("Destroy did not remove native panel")
	}
}

func TestPanelRuntimeRejectsContentControl(t *testing.T) {
	window := &WebviewWindow{}
	panel := window.NewPanel(WebviewPanelOptions{Name: "browser", URL: "https://example.com"})
	for _, id := range []int{3, 4, 5, -1, 999} {
		_, err := (&MessageProcessor{}).processPanelMethod(windowMethodRequest(t, id, `{"panel":"browser","url":"file:///etc/passwd","html":"unsafe","js":"unsafe"}`), window)
		if err == nil {
			t.Errorf("method %d accepted", id)
		}
	}
	if panel.URL() != "https://example.com" || panel.isDestroyed() {
		t.Fatal("rejected method mutated panel")
	}
}

func TestPanelRuntimeLookupAndValidation(t *testing.T) {
	window := &WebviewWindow{}
	panel := window.NewPanel(WebviewPanelOptions{Name: "browser"})
	processor := &MessageProcessor{}
	for _, args := range []string{`{"panel":"browser"}`, fmt.Sprintf(`{"panelId":%d}`, panel.ID()), fmt.Sprintf(`{"panel":"missing","panelId":%d}`, panel.ID())} {
		result, err := processor.processPanelMethod(windowMethodRequest(t, PanelName, args), window)
		if err != nil || result != "browser" {
			t.Errorf("lookup %s: %v, %v", args, result, err)
		}
	}
	for _, args := range []string{`{}`, `{"panelId":0}`, `{"panelId":-1}`, `{"panel":"missing"}`} {
		if _, err := processor.processPanelMethod(windowMethodRequest(t, PanelName, args), window); err == nil {
			t.Errorf("accepted %s", args)
		}
	}
	for _, args := range []string{`{}`, `{"panel":"browser","x":0,"y":0,"width":0,"height":10}`, `{"panel":"browser","x":0,"y":0,"width":-1,"height":10}`} {
		if _, err := processor.processPanelMethod(windowMethodRequest(t, PanelSetBounds, args), window); err == nil {
			t.Errorf("accepted bounds %s", args)
		}
	}
	for _, zoom := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if _, err := handlePanelSetZoom(panel, &MapArgs{data: map[string]any{"zoom": zoom}}); err == nil {
			t.Errorf("accepted zoom %v", zoom)
		}
	}
	if _, err := processor.processPanelMethod(&RuntimeRequest{}, window); err == nil {
		t.Fatal("nil args accepted")
	}
}

func TestPanelNamesAndWindowTeardown(t *testing.T) {
	window := &WebviewWindow{}
	var wg sync.WaitGroup
	panels := make(chan *WebviewPanel, 32)
	for range 32 {
		wg.Add(1)
		go func() { defer wg.Done(); panels <- window.NewPanel(WebviewPanelOptions{Name: "same"}) }()
	}
	wg.Wait()
	close(panels)
	first := window.GetPanel("same")
	for panel := range panels {
		if panel != first {
			t.Fatal("duplicate named panel")
		}
	}
	window.destroyAllPanels()
	if len(window.GetPanels()) != 0 || !first.isDestroyed() {
		t.Fatal("window did not destroy its panels")
	}
	if window.NewPanel(WebviewPanelOptions{}) != nil {
		t.Fatal("created panel during teardown")
	}
}

func TestPanelDestroyExactlyOnce(t *testing.T) {
	panelTestApp(t, true)
	window := &WebviewWindow{}
	panel := window.NewPanel(WebviewPanelOptions{Name: "browser"})
	native := &panelProbe{}
	panel.impl = native
	var wg sync.WaitGroup
	for range 64 {
		wg.Add(1)
		go func() { defer wg.Done(); panel.Destroy() }()
	}
	wg.Wait()
	if native.destroyed.Load() != 1 {
		t.Fatalf("destroy called %d times", native.destroyed.Load())
	}
	if panel.IsVisible() {
		t.Fatal("destroyed panel is visible")
	}
}

func TestPanelDestroyDuringCreation(t *testing.T) {
	for _, dispatchDuringCreate := range []bool{false, true} {
		t.Run(fmt.Sprintf("dispatchDuringCreate=%t", dispatchDuringCreate), func(t *testing.T) {
			dispatcher := panelTestApp(t, false)
			panel := NewPanel(WebviewPanelOptions{})
			native := &panelProbe{}
			panel.impl = native
			panel.creating = true
			done := make(chan struct{})
			go func() { panel.Destroy(); close(done) }()
			deadline := time.After(time.Second)
			for dispatcher.pendingCount() == 0 {
				select {
				case <-deadline:
					t.Fatal("Destroy was not queued")
				default:
					time.Sleep(time.Millisecond)
				}
			}
			if dispatchDuringCreate {
				// A native create call may dispatch the queued callback itself.
				dispatcher.runPending()
			}
			if native.destroyed.Load() != 0 {
				t.Fatal("destroy called before native creation finished")
			}
			// Creation returns on the UI thread and sees the pending destruction.
			panel.creating = false
			panel.destroyNative(native)
			dispatcher.runPending()
			<-done
			panel.Destroy()
			if native.destroyed.Load() != 1 {
				t.Fatalf("destroy called %d times", native.destroyed.Load())
			}
		})
	}
}

func TestPanelQueuedOperationSkipsDestroyedNativeView(t *testing.T) {
	dispatcher := panelTestApp(t, false)
	panel := NewPanel(WebviewPanelOptions{})
	native := &panelProbe{}
	panel.impl = native
	done := make(chan struct{})
	go func() { panel.IsFocused(); close(done) }()
	deadline := time.After(time.Second)
	for {
		dispatcher.mu.Lock()
		count := len(dispatcher.pending)
		dispatcher.mu.Unlock()
		if count > 0 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("operation was not queued")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	panel.destroyedLock.Lock()
	panel.destroyed = true
	panel.destroyedLock.Unlock()
	dispatcher.runPending()
	<-done
	if native.focused.Load() != 0 {
		t.Fatal("queued query accessed destroyed view")
	}
}

func TestPanelAnchorBaselineAndClamping(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Width: 800, Height: 600}}
	panel := window.NewPanel(WebviewPanelOptions{X: 10, Y: 20, Width: 300, Height: 200, Anchor: AnchorFill})
	panel.initializeAnchor()
	panel.handleWindowResize(900, 650)
	want := Rect{X: 10, Y: 20, Width: 400, Height: 250}
	if got := panel.Bounds(); got != want {
		t.Fatalf("resize: %+v, want %+v", got, want)
	}
	panel.handleWindowResize(800, 600)
	if got := panel.Bounds(); got != (Rect{X: 10, Y: 20, Width: 300, Height: 200}) {
		t.Fatalf("resize drift: %+v", got)
	}
	panel.SetBounds(Rect{X: 50, Y: 60, Width: 200, Height: 100})
	panel.handleWindowResize(900, 650)
	if got := panel.Bounds(); got != (Rect{X: 50, Y: 60, Width: 300, Height: 150}) {
		t.Fatalf("manual baseline ignored: %+v", got)
	}
	panel.handleWindowResize(1, 1)
	if got := panel.Bounds(); got.Width != 1 || got.Height != 1 {
		t.Fatalf("negative native dimensions: %+v", got)
	}
	before := panel.Bounds()
	panel.FillBeside((&WebviewWindow{}).NewPanel(WebviewPanelOptions{}), "right")
	if panel.Bounds() != before {
		t.Fatal("used reference from another window")
	}
}

func TestPanelCopiesOptionsAndRejectsInvalidURL(t *testing.T) {
	visible, devtools := false, false
	headers := map[string]string{"X-Test": "original"}
	panel := NewPanel(WebviewPanelOptions{Visible: &visible, DevToolsEnabled: &devtools, Headers: headers, URL: "/panel.html"})
	visible, devtools = true, true
	headers["X-Test"] = "changed"
	if panel.IsVisible() || *panel.snapshotOptions().DevToolsEnabled || panel.snapshotOptions().Headers["X-Test"] != "original" {
		t.Fatal("options alias caller memory")
	}
	original := panel.URL()
	for _, url := range []string{"", "http://%zz", "\x00"} {
		panel.SetURL(url)
		if panel.URL() != original {
			t.Errorf("invalid URL changed navigation: %q", url)
		}
	}
	if invalid := NewPanel(WebviewPanelOptions{URL: "http://%zz"}); invalid.URL() != "" {
		t.Fatal("invalid initial URL retained")
	}
}

func TestPanelConcurrentStateBeforeRun(t *testing.T) {
	panel := NewPanel(WebviewPanelOptions{})
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				panel.SetBounds(Rect{Width: 100, Height: 80})
				panel.Bounds()
				panel.Hide()
				panel.Show()
				panel.IsVisible()
				panel.SetZoom(1.25)
				panel.GetZoom()
				panel.SetZIndex(2)
				panel.ZIndex()
				panel.SetURL("https://example.com")
				panel.URL()
			}
		}()
	}
	wg.Wait()
}

func TestPanelIntegerRejectsTruncationAndOverflow(t *testing.T) {
	for _, value := range []float64{1.5, -1.5, math.NaN(), math.Inf(1), math.MaxInt32 + 1.0} {
		if panelInteger(&MapArgs{data: map[string]any{"value": value}}, "value") != nil {
			t.Errorf("accepted invalid native integer %v", value)
		}
	}
}

func TestPanelStackingOrder(t *testing.T) {
	window := &WebviewWindow{}
	a := window.NewPanel(WebviewPanelOptions{Name: "a", ZIndex: 3})
	b := window.NewPanel(WebviewPanelOptions{Name: "b", ZIndex: 1})
	c := window.NewPanel(WebviewPanelOptions{Name: "c", ZIndex: 3})
	ordered := a.sortedSiblings()
	if ordered[0] != b || ordered[1] != a || ordered[2] != c {
		t.Fatal("incorrect stacking order")
	}
	b.SetZIndex(4)
	ordered = a.sortedSiblings()
	if ordered[0] != a || ordered[1] != c || ordered[2] != b {
		t.Fatal("SetZIndex did not reorder siblings")
	}
}

func TestPanelSingleEdgeAnchoring(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Width: 800, Height: 600}}
	panel := window.NewPanel(WebviewPanelOptions{X: 600, Y: 400, Width: 100, Height: 100, Anchor: AnchorRight | AnchorBottom})
	panel.initializeAnchor()
	panel.handleWindowResize(900, 650)
	if got := panel.Bounds(); got != (Rect{X: 700, Y: 450, Width: 100, Height: 100}) {
		t.Fatalf("single edge changed dimensions: %+v", got)
	}
}
