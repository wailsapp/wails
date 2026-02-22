//go:build windows

package application

import (
	"runtime"
	"sort"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

var (
	wmInvokeCallback uint32
)

func init() {
	wmInvokeCallback = w32.RegisterWindowMessage(w32.MustStringToUTF16Ptr("WailsV0.InvokeCallback"))
}

// initMainLoop must be called with the same OSThread that is used to call runMainLoop() later.
func (m *windowsApp) initMainLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if m.mainThreadWindowHWND != 0 {
		panic("initMainLoop was already called")
	}

	// We need a hidden window so we can PostMessage to it, if we don't use PostMessage for dispatching to a HWND
	// messages might get lost if a modal inner loop is being run.
	// We had this once in V2: https://github.com/wailsapp/wails/issues/969
	// See: https://devblogs.microsoft.com/oldnewthing/20050426-18/?p=35783
	// See also: https://learn.microsoft.com/en-us/windows/win32/winmsg/using-messages-and-message-queues#creating-a-message-loop
	// > Because the system directs messages to individual windows in an application, a thread must create at least one window before starting its message loop.
	m.mainThreadWindowHWND = w32.CreateWindowEx(
		0,
		w32.MustStringToUTF16Ptr(m.parent.options.Windows.WndClass),
		w32.MustStringToUTF16Ptr("__wails_hidden_mainthread"),
		w32.WS_DISABLED,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		0,
		0,
		0,
		0,
		w32.GetModuleHandle(""),
		nil)

	m.mainThreadID, _ = w32.GetWindowThreadProcessId(m.mainThreadWindowHWND)
}

func (m *windowsApp) runMainLoop() int {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if m.invokeRequired() {
		panic("invokeRequired for runMainLoop, the mainloop must be running on the same OSThread as the mainThreadWindow has been created on")
	}

	msg := (*w32.MSG)(unsafe.Pointer(w32.GlobalAlloc(0, uint32(unsafe.Sizeof(w32.MSG{})))))
	defer w32.GlobalFree(w32.HGLOBAL(unsafe.Pointer(msg)))

	for w32.GetMessage(msg, 0, 0, 0) != 0 {
		w32.TranslateMessage(msg)
		w32.DispatchMessage(msg)
	}

	return int(msg.WParam)
}

func (m *windowsApp) dispatchOnMainThread(id uint) {
	mainThreadHWND := m.mainThreadWindowHWND
	if mainThreadHWND == 0 {
		panic("initMainLoop was not called")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if m.invokeRequired() {
		w32.PostMessage(mainThreadHWND, wmInvokeCallback, uintptr(id), 0)
	} else {
		mainThreadFunctionStoreLock.Lock()
		fn := mainThreadFunctionStore[id]
		delete(mainThreadFunctionStore, id)
		mainThreadFunctionStoreLock.Unlock()

		if fn == nil {
			Fatal("dispatchOnMainThread called with invalid id: %v", id)
		}
		fn()
	}
}

func (m *windowsApp) invokeRequired() bool {
	mainThreadID := m.mainThreadID
	if mainThreadID == 0 {
		panic("initMainLoop was not called")
	}

	return mainThreadID != w32.GetCurrentThreadId()
}

func (m *windowsApp) invokeCallback(wParam, lParam uintptr) {
	// TODO: Should we invoke just one or all queued? In v2 we always invoked all pendings...
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if m.invokeRequired() {
		panic("invokeCallback must always be called on the MainOSThread")
	}

	mainThreadFunctionStoreLock.Lock()
	fnIDs := make([]uint, 0, len(mainThreadFunctionStore))
	for id := range mainThreadFunctionStore {
		fnIDs = append(fnIDs, id)
	}
	sort.Slice(fnIDs, func(i, j int) bool { return fnIDs[i] < fnIDs[j] })

	fns := make([]func(), len(fnIDs))
	for i, id := range fnIDs {
		fns[i] = mainThreadFunctionStore[id]
		delete(mainThreadFunctionStore, id)
	}
	mainThreadFunctionStoreLock.Unlock()

	for _, fn := range fns {
		fn()
	}
}
