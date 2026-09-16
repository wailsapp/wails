package application

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"
)

// These tests cover constructing split WebView windows and native windows
// after App.Run without a GUI: the platform layer is the controllable
// threadProbeApp from event_ordering_test.go and the window implementations
// are fakes, so what is exercised is the Go dispatch, validation, ownership,
// and lifecycle logic that the native code sits behind.

// runtimeProbe installs a fake running application whose reported errors are
// collected for assertions.
type runtimeProbe struct {
	app    *threadProbeApp
	errors []string
	mu     sync.Mutex
}

func newRuntimeProbe(t *testing.T, running bool) *runtimeProbe {
	t.Helper()
	probe := &runtimeProbe{app: &threadProbeApp{}}
	probe.app.onMain.Store(true)
	previous := globalApplication
	globalApplication = &App{
		impl:          probe.app,
		running:       running,
		nativeWindows: map[uint]*NativeWindow{},
		options: Options{ErrorHandler: func(err error) {
			probe.mu.Lock()
			probe.errors = append(probe.errors, err.Error())
			probe.mu.Unlock()
		}},
	}
	previousFactory := nativeWindowImplFactory
	t.Cleanup(func() {
		globalApplication = previous
		nativeWindowImplFactory = previousFactory
	})
	return probe
}

func (p *runtimeProbe) reported() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return strings.Join(p.errors, "\n")
}

// fakeSplitInstaller stands in for the macOS window implementation. It
// records whether the late installation reached the application thread and
// can simulate native success or failure.
type fakeSplitInstaller struct {
	webviewWindowImpl
	mu       sync.Mutex
	calls    int
	onMain   bool
	err      error
	markDone bool
	parent   *WebviewWindow
}

func (f *fakeSplitInstaller) installSplitViewLate() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.onMain = globalApplication.impl.isOnMainThread()
	if f.err != nil {
		return f.err
	}
	if f.markDone {
		f.parent.splitViewLock.RLock()
		split := f.parent.splitView
		f.parent.splitViewLock.RUnlock()
		split.lock.Lock()
		split.installed = true
		split.lock.Unlock()
	}
	return nil
}

func newLateSplitWindow(installer *fakeSplitInstaller) *WebviewWindow {
	window := NewWindow(WebviewWindowOptions{Name: "late"})
	installer.parent = window
	window.impl = installer
	return window
}

func TestSetSplitViewAfterCreationDispatchesToApplicationThread(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	installer := &fakeSplitInstaller{markDone: true}
	window := newLateSplitWindow(installer)
	split, _, _, _ := newTestSidebarSplit()

	// Call from a goroutine, as a tray or service callback would, so the
	// installation has to hop to the application thread.
	probe.app.onMain.Store(false)
	done := make(chan struct{})
	go func() {
		defer close(done)
		window.SetSplitView(split)
	}()
	deadline := time.Now().Add(5 * time.Second)
	for probe.app.pendingCount() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	probe.app.runPending()
	<-done

	if installer.calls != 1 || !installer.onMain {
		t.Fatalf("late install calls = %d on main thread = %v; want one call on the application thread", installer.calls, installer.onMain)
	}
	if split.ownerWindow() != window || !split.isInstalled() {
		t.Fatal("late SetSplitView should claim and install the layout")
	}
	window.splitViewLock.RLock()
	stored := window.splitView
	window.splitViewLock.RUnlock()
	if stored != split {
		t.Fatal("late SetSplitView should record the layout on the window")
	}
	if got := probe.reported(); got != "" {
		t.Fatalf("unexpected errors: %s", got)
	}
}

func TestSetSplitViewAfterCreationRejectsReplacingInstalledLayout(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	installer := &fakeSplitInstaller{markDone: true}
	window := newLateSplitWindow(installer)
	first, _, _, _ := newTestSidebarSplit()
	second, _, _, _ := newTestSidebarSplit()

	window.SetSplitView(first)
	window.SetSplitView(second)

	if installer.calls != 1 {
		t.Fatalf("install calls = %d, want the second layout to be rejected before dispatch", installer.calls)
	}
	if !strings.Contains(probe.reported(), ErrMacSplitViewAlreadyInstalled.Error()) {
		t.Fatalf("errors = %q, want ErrMacSplitViewAlreadyInstalled", probe.reported())
	}
	if second.ownerWindow() != nil {
		t.Fatal("a rejected layout must stay unclaimed")
	}
	// Removing an installed layout is equally unsupported.
	window.SetSplitView(nil)
	if !first.isInstalled() || first.ownerWindow() != window {
		t.Fatal("SetSplitView(nil) must not disturb an installed layout")
	}
}

func TestSetSplitViewAfterCreationReleasesLayoutWhenInstallFails(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	installer := &fakeSplitInstaller{err: errors.New("native split view refused")}
	window := newLateSplitWindow(installer)
	split, _, _, _ := newTestSidebarSplit()

	window.SetSplitView(split)

	if !strings.Contains(probe.reported(), "native split view refused") {
		t.Fatalf("errors = %q, want the native failure reported", probe.reported())
	}
	window.splitViewLock.RLock()
	stored := window.splitView
	window.splitViewLock.RUnlock()
	if stored != nil {
		t.Fatal("a failed installation must not leave the layout on the window")
	}
	if split.ownerWindow() != nil {
		t.Fatal("a failed installation must release the layout for another window")
	}
	// The released layout is usable by a second window.
	other := newLateSplitWindow(&fakeSplitInstaller{markDone: true})
	other.SetSplitView(split)
	if split.ownerWindow() != other || !split.isInstalled() {
		t.Fatal("the released layout should install into another window")
	}
}

func TestSetSplitViewAfterCreationValidatesBeforeDispatch(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	installer := &fakeSplitInstaller{markDone: true}
	window := newLateSplitWindow(installer)
	invalid := NewMacSplitView()
	invalid.AddPrimaryContent()

	window.SetSplitView(invalid)

	if installer.calls != 0 {
		t.Fatal("an invalid layout must be rejected before reaching the application thread")
	}
	if !strings.Contains(probe.reported(), "at least two panes") {
		t.Fatalf("errors = %q, want validation failure", probe.reported())
	}
}

func TestSetSplitViewAfterCreationWithoutAppKitStoresLayout(t *testing.T) {
	newRuntimeProbe(t, true)
	window := NewWindow(WebviewWindowOptions{Name: "plain"})
	window.impl = &stubWindowImpl{}
	split, _, _, _ := newTestSidebarSplit()

	window.SetSplitView(split)

	if split.ownerWindow() != window || split.isInstalled() {
		t.Fatal("without a late installer the layout is claimed and stored but never installed")
	}
	window.SetSplitView(nil)
	if split.ownerWindow() != nil {
		t.Fatal("a stored, uninstalled layout can still be released")
	}
}

func TestWebviewWindowOptionsClaimSplitViewAndToolbar(t *testing.T) {
	newRuntimeProbe(t, false)
	split, _, _, _ := newTestSidebarSplit()
	toolbar := NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()

	window := NewWindow(WebviewWindowOptions{
		Name: "options",
		Mac:  MacWindow{SplitView: split, Toolbar: toolbar},
	})

	window.splitViewLock.RLock()
	storedSplit := window.splitView
	window.splitViewLock.RUnlock()
	window.toolbarLock.RLock()
	storedToolbar := window.toolbar
	window.toolbarLock.RUnlock()
	if storedSplit != split || split.ownerWindow() != window {
		t.Fatal("Mac.SplitView should be queued and claimed before the window runs")
	}
	if storedToolbar != toolbar {
		t.Fatal("Mac.Toolbar should be queued before the window runs")
	}
	if _, err := claimMacToolbar(toolbar, NewWindow(WebviewWindowOptions{Name: "other"})); err == nil {
		t.Fatal("the options toolbar should already be owned by the new window")
	}
}

func TestWebviewWindowOptionsRejectInvalidSplitView(t *testing.T) {
	probe := newRuntimeProbe(t, false)
	invalid := NewMacSplitView()
	invalid.AddPrimaryContent()

	window := NewWindow(WebviewWindowOptions{Name: "invalid", Mac: MacWindow{SplitView: invalid}})

	window.splitViewLock.RLock()
	stored := window.splitView
	window.splitViewLock.RUnlock()
	if stored != nil || invalid.ownerWindow() != nil {
		t.Fatal("an invalid options layout must not be stored or claimed")
	}
	if !strings.Contains(probe.reported(), "at least two panes") {
		t.Fatalf("errors = %q, want validation failure", probe.reported())
	}
}

// fakeNativeWindowImpl stands in for the AppKit NativeWindow implementation.
type fakeNativeWindowImpl struct {
	runErr  error
	runs    int
	visible bool
}

func (f *fakeNativeWindowImpl) run() error {
	f.runs++
	if f.runErr != nil {
		return f.runErr
	}
	f.visible = true
	return nil
}
func (f *fakeNativeWindowImpl) show()                      { f.visible = true }
func (f *fakeNativeWindowImpl) hide()                      { f.visible = false }
func (*fakeNativeWindowImpl) focus()                       {}
func (f *fakeNativeWindowImpl) close()                     { f.visible = false }
func (f *fakeNativeWindowImpl) isVisible() bool            { return f.visible }
func (*fakeNativeWindowImpl) setTitle(string)              {}
func (*fakeNativeWindowImpl) nativeWindow() unsafe.Pointer { return nil }
func (*fakeNativeWindowImpl) setToolbar(*MacToolbar) error { return nil }
func (*fakeNativeWindowImpl) installSplitView() error      { return nil }

func useFakeNativeWindowImpl(t *testing.T, fake *fakeNativeWindowImpl) {
	t.Helper()
	nativeWindowImplFactory = func(*NativeWindow) nativeWindowImpl { return fake }
}

func TestNativeWindowRunRequiresEditorContent(t *testing.T) {
	newRuntimeProbe(t, false)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	window := newNativeWindow(NativeWindowOptions{Name: "empty"})

	if err := window.Run(); !errors.Is(err, ErrNativeWindowContentRequired) {
		t.Fatalf("Run without content = %v, want ErrNativeWindowContentRequired", err)
	}
	if fake.runs != 0 || window.implementation() != nil {
		t.Fatal("a window without content must not create its implementation")
	}

	webviewSplit, _, _, _ := newTestSidebarSplit()
	if err := window.SetSplitView(webviewSplit); !errors.Is(err, ErrNativeWindowEditorRequired) {
		t.Fatalf("SetSplitView with a WebView primary pane = %v, want ErrNativeWindowEditorRequired", err)
	}
	if webviewSplit.ownerWindow() != nil {
		t.Fatal("a rejected layout must stay unclaimed")
	}
}

func TestNativeWindowRunPropagatesImplementationError(t *testing.T) {
	probe := newRuntimeProbe(t, false)
	failing := &fakeNativeWindowImpl{runErr: errors.New("NSWindow refused")}
	useFakeNativeWindowImpl(t, failing)
	window := newNativeWindow(NativeWindowOptions{Name: "failing"})
	split, _ := nativeEditorSplitForTest()
	if err := window.SetSplitView(split); err != nil {
		t.Fatalf("SetSplitView: %v", err)
	}

	err := window.Run()
	if err == nil || !strings.Contains(err.Error(), "NSWindow refused") {
		t.Fatalf("Run = %v, want the implementation error returned", err)
	}
	if !strings.Contains(probe.reported(), "NSWindow refused") {
		t.Fatalf("errors = %q, want the failure logged as well", probe.reported())
	}
	if window.implementation() != nil {
		t.Fatal("a failed Run must not leave a half-built implementation behind")
	}

	// The window stays usable: a working implementation runs it.
	working := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, working)
	if err := window.Run(); err != nil {
		t.Fatalf("Run after fixing the failure: %v", err)
	}
	if err := window.Run(); err != nil || working.runs != 1 {
		t.Fatalf("second Run = %v with %d runs; want a nil no-op", err, working.runs)
	}
	if !window.IsVisible() {
		t.Fatal("a created window should report visible through its implementation")
	}
	if err := window.SetSplitView(nativeSplitForTest()); !errors.Is(err, ErrMacSplitViewAlreadyInstalled) {
		t.Fatalf("SetSplitView after creation = %v, want ErrMacSplitViewAlreadyInstalled", err)
	}
	window.Close()
	if err := window.Run(); !errors.Is(err, ErrNativeWindowClosed) {
		t.Fatalf("Run after Close = %v, want ErrNativeWindowClosed", err)
	}
}

func nativeSplitForTest() *MacSplitView {
	split, _ := nativeEditorSplitForTest()
	return split
}

func TestNativeWindowManagerDefersCreationUntilContentAfterAppRun(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	manager := newNativeWindowManager(globalApplication)

	window := manager.New()
	if fake.runs != 0 || window.implementation() != nil {
		t.Fatal("New after App.Run must defer creation until the window has content")
	}
	if got := probe.reported(); got != "" {
		t.Fatalf("deferral must not be reported as an error: %s", got)
	}
	if _, ok := manager.GetByID(window.ID()); !ok {
		t.Fatal("a deferred window is still registered with the manager")
	}

	if err := window.SetSplitView(nativeSplitForTest()); err != nil {
		t.Fatalf("SetSplitView on a running application: %v", err)
	}
	if fake.runs != 1 || window.implementation() == nil || !window.IsVisible() {
		t.Fatalf("SetSplitView should create the deferred window: runs=%d visible=%v", fake.runs, window.IsVisible())
	}
}

func TestNativeWindowManagerCreatesFromOptionsAfterAppRun(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	manager := newNativeWindowManager(globalApplication)
	split, _ := nativeEditorSplitForTest()
	toolbar := NewMacToolbar()
	toolbar.AddSidebarToggle()

	window := manager.NewWithOptions(NativeWindowOptions{Name: "atomic", SplitView: split, Toolbar: toolbar})

	if fake.runs != 1 || window.implementation() == nil {
		t.Fatalf("NewWithOptions after App.Run should create the window once, runs=%d", fake.runs)
	}
	window.lock.RLock()
	storedSplit, storedToolbar := window.split, window.toolbar
	window.lock.RUnlock()
	if storedSplit != split || storedToolbar != toolbar || split.ownerWindow() != window {
		t.Fatal("options-time split view and toolbar should be claimed before creation")
	}
	if got := probe.reported(); got != "" {
		t.Fatalf("unexpected errors: %s", got)
	}
}

func TestNativeWindowManagerQueuesBeforeAppRun(t *testing.T) {
	newRuntimeProbe(t, false)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	manager := newNativeWindowManager(globalApplication)
	split, _ := nativeEditorSplitForTest()

	window := manager.NewWithOptions(NativeWindowOptions{Name: "queued", SplitView: split})

	if fake.runs != 0 || window.implementation() != nil {
		t.Fatal("before App.Run the window must wait for the start-up queue")
	}
	if len(globalApplication.pendingRun) != 1 {
		t.Fatalf("pendingRun = %d entries, want the window queued once", len(globalApplication.pendingRun))
	}
	globalApplication.pendingRun[0].Run()
	if fake.runs != 1 || window.implementation() == nil {
		t.Fatal("the start-up queue should create the queued window")
	}
}

func TestNativeWindowOptionsReportInvalidSplitViewAndStayDeferred(t *testing.T) {
	probe := newRuntimeProbe(t, true)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	manager := newNativeWindowManager(globalApplication)
	webviewSplit, _, _, _ := newTestSidebarSplit()

	window := manager.NewWithOptions(NativeWindowOptions{Name: "bad", SplitView: webviewSplit})

	if fake.runs != 0 || window.implementation() != nil {
		t.Fatal("an invalid options layout must leave the window deferred")
	}
	if !strings.Contains(probe.reported(), ErrNativeWindowEditorRequired.Error()) {
		t.Fatalf("errors = %q, want ErrNativeWindowEditorRequired", probe.reported())
	}
	if err := window.SetSplitView(nativeSplitForTest()); err != nil {
		t.Fatalf("a valid layout should still create the window: %v", err)
	}
	if fake.runs != 1 {
		t.Fatal("SetSplitView should create the window after the option was rejected")
	}
}

func TestNativeWindowShowCreatesDeferredWindow(t *testing.T) {
	newRuntimeProbe(t, true)
	fake := &fakeNativeWindowImpl{}
	useFakeNativeWindowImpl(t, fake)
	window := newNativeWindow(NativeWindowOptions{Name: "show", Hidden: true})

	window.Show()
	if fake.runs != 0 {
		t.Fatal("Show without content cannot create the window")
	}
	window.lock.Lock()
	window.split = nativeSplitForTest()
	window.lock.Unlock()
	window.Show()
	if fake.runs != 1 || !window.IsVisible() {
		t.Fatalf("Show should create and show a deferred window with content: runs=%d", fake.runs)
	}
}
