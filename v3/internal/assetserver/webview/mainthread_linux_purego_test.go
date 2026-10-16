//go:build linux && !wailscgo && !android

package webview

import (
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
)

var testLoop struct {
	sync.Mutex
	loop, thread     uintptr
	ready            chan struct{}
	new              func(uintptr, int32) uintptr
	run, quit, unref func(uintptr)
	self             func() uintptr
	idle             func(uintptr, uintptr) uint32
}

var testLoopReady = purego.NewCallback(func(_ uintptr) int32 {
	close(testLoop.ready)
	return 0
})

func testRunMainLoop() {
	ensureGLib()
	testLoop.Lock()
	for name, target := range map[string]any{
		"g_main_loop_new": &testLoop.new, "g_main_loop_run": &testLoop.run,
		"g_main_loop_quit": &testLoop.quit, "g_main_loop_unref": &testLoop.unref,
		"g_thread_self": &testLoop.self, "g_idle_add": &testLoop.idle,
	} {
		mustRegisterWebviewFunc(target, webviewLibGLib, name)
	}
	testLoop.loop, testLoop.thread = testLoop.new(0, 0), testLoop.self()
	testLoop.ready = make(chan struct{})
	testLoop.idle(testLoopReady, 0)
	testLoop.Unlock()
	testLoop.run(testLoop.loop)
	testLoop.unref(testLoop.loop)
}

func testWaitLoopRunning() {
	for {
		testLoop.Lock()
		ready := testLoop.ready
		testLoop.Unlock()
		if ready != nil {
			<-ready
			return
		}
		runtime.Gosched()
	}
}

func testQuitMainLoop()      { testLoop.quit(testLoop.loop) }
func testOnLoopThread() bool { return testLoop.self() == testLoop.thread }

func TestRequestLifetimeNativeCancellation(t *testing.T) {
	// A real GObject exercises the native toggle callback and handle lifetime
	// without requiring a display or manufacturing a WebKit request.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ensureWebviewLibs()
	lib := dlopenWebviewLib("libgobject-2.0.so.0")
	var getType func() uintptr
	var create func(uintptr, uint32, uintptr, uintptr) uintptr
	mustRegisterWebviewFunc(&getType, lib, "g_object_get_type")
	mustRegisterWebviewFunc(&create, lib, "g_object_new_with_properties")
	object := create(getType(), 0, 0, 0)
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&object))
	lifetime := retainRequest(ptr)
	if lifetime.Err() != nil {
		t.Fatal("request cancelled while native owner holds a reference")
	}
	g_object_unref(object)
	if lifetime.Err() == nil {
		t.Fatal("dropping native owner's reference did not cancel request")
	}
	lifetime.releaseNative()
	lifetime.releaseNative()
	requestCancelHandles.Lock()
	_, retained := requestCancelHandles.values[lifetime.handle]
	requestCancelHandles.Unlock()
	if retained {
		t.Fatal("released request retained its callback handle")
	}
}

func TestRequestLifetimeAfterShutdown(t *testing.T) {
	webviewDispatchMu.Lock()
	previous := webviewMainDispatchEnabled
	webviewMainDispatchEnabled = false
	webviewDispatchMu.Unlock()
	defer func() {
		webviewDispatchMu.Lock()
		webviewMainDispatchEnabled = previous
		webviewDispatchMu.Unlock()
	}()
	if invokeOnMainSyncIfEnabled(func() { t.Error("native release ran after shutdown") }) {
		t.Fatal("dispatch accepted a native release after shutdown")
	}
	mainSyncMu.Lock()
	n := len(mainSyncCalls)
	mainSyncMu.Unlock()
	if n != 0 {
		t.Fatalf("shutdown dispatch leaked %d callbacks", n)
	}
}

func TestShutdownReleasesQueuedWorkers(t *testing.T) {
	ran := false
	normal := &mainSyncCall{fn: func() { ran = true }, done: make(chan struct{}), runIfDisabled: true}
	native := &mainSyncCall{fn: func() { t.Error("native release executed after shutdown") }, done: make(chan struct{})}
	webviewDispatchMu.Lock()
	previous := webviewMainDispatchEnabled
	webviewMainDispatchEnabled = true
	mainSyncMu.Lock()
	mainSyncNextID++
	mainSyncCalls[mainSyncNextID] = normal
	mainSyncNextID++
	mainSyncCalls[mainSyncNextID] = native
	mainSyncMu.Unlock()
	webviewDispatchMu.Unlock()
	DisableMainThreadDispatch()
	defer func() { webviewDispatchMu.Lock(); webviewMainDispatchEnabled = previous; webviewDispatchMu.Unlock() }()
	for _, call := range []*mainSyncCall{normal, native} {
		select {
		case <-call.done:
		default:
			t.Fatal("shutdown left a queued worker blocked")
		}
	}
	if !ran || !normal.executed || native.executed {
		t.Fatal("incorrect queued callback completion state")
	}
}

func TestDispatchPanicReleasesWaiters(t *testing.T) {
	for _, shutdown := range []bool{false, true} {
		call := &mainSyncCall{fn: func() { panic("callback failed") }, done: make(chan struct{}), runIfDisabled: true}
		webviewDispatchMu.Lock()
		previous := webviewMainDispatchEnabled
		mainSyncMu.Lock()
		mainSyncNextID++
		id := mainSyncNextID
		mainSyncCalls[id] = call
		mainSyncMu.Unlock()
		webviewDispatchMu.Unlock()
		func() {
			defer func() {
				if recover() != "callback failed" {
					t.Error("callback panic was not propagated")
				}
			}()
			if shutdown {
				DisableMainThreadDispatch()
			} else {
				runMainSyncCallback(id)
			}
		}()
		webviewDispatchMu.Lock()
		webviewMainDispatchEnabled = previous
		webviewDispatchMu.Unlock()
		select {
		case <-call.done:
		default:
			t.Fatalf("callback panic left its waiter blocked (shutdown=%v)", shutdown)
		}
	}
}
