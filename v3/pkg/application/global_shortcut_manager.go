package application

import (
	"fmt"
	"sort"
	"sync"
)

// globalShortcutImpl is the platform-specific implementation of global
// (system-wide) keyboard shortcuts. Each platform registers a shortcut with the
// operating system against an integer id that the native event handler reports
// back when the shortcut fires.
//
// All methods are called on the main thread (see GlobalShortcutManager which
// wraps every call in InvokeSync*). Implementations must not assume otherwise.
type globalShortcutImpl interface {
	// register asks the OS to bind the given accelerator to id. It returns an
	// error if the OS rejects the registration (for example, because another
	// application already owns the shortcut).
	register(id int, accel *accelerator) error
	// unregister releases the OS binding for id.
	unregister(id int) error
	// unregisterAll releases every binding owned by this application.
	unregisterAll() error
}

// globalShortcut is a single registered shortcut.
type globalShortcut struct {
	id          int
	accelerator string // canonical, normalized accelerator string
	parsed      *accelerator
	callback    func()
}

// GlobalShortcutManager manages application-wide (global) keyboard shortcuts.
//
// Unlike menu accelerators or [KeyBindingManager] - which only fire while a
// Wails window has focus - a global shortcut fires regardless of which
// application is currently focused, as long as the Wails application is
// running.
//
// Global shortcuts are owned by the application, not by an individual window.
// Registering the same accelerator twice within the same application is
// reported as an error and the original binding is preserved; see [Register].
//
// Shortcuts may be registered before [App.Run] is called: the binding with the
// operating system is then deferred until the application starts.
type GlobalShortcutManager struct {
	app  *App
	impl globalShortcutImpl

	mu      sync.Mutex
	started bool                       // set once the app is running and pending shortcuts are flushed
	byName  map[string]*globalShortcut // keyed by canonical accelerator string
	byID    map[int]*globalShortcut    // keyed by native id
	pending []*globalShortcut          // registered before the app started; bound on start
	nextID  int
}

// newGlobalShortcutManager creates a new GlobalShortcutManager instance.
func newGlobalShortcutManager(app *App) *GlobalShortcutManager {
	return &GlobalShortcutManager{
		app:    app,
		byName: make(map[string]*globalShortcut),
		byID:   make(map[int]*globalShortcut),
	}
}

// getImpl returns the platform implementation, creating it lazily on first use
// so that platforms which do not support global shortcuts do not pay any cost
// unless the feature is actually used. Callers must hold m.mu.
func (m *GlobalShortcutManager) getImpl() globalShortcutImpl {
	if m.impl == nil {
		m.impl = newGlobalShortcutImpl(m)
	}
	return m.impl
}

// Register binds the given accelerator (for example "Ctrl+Shift+P" or
// "Cmd+Option+K") to callback. The callback is invoked - on its own goroutine -
// whenever the shortcut is pressed, even when the application does not have
// focus.
//
// Accelerators use the same syntax as menu accelerators (see SetAccelerator).
// "CmdOrCtrl" resolves to Command on macOS and Control elsewhere.
//
// Register may be called before [App.Run]; the OS binding is then performed
// when the application starts and any OS-level rejection is reported via the
// application's error handler rather than returned here.
//
// Register returns an error when:
//   - the accelerator string cannot be parsed;
//   - the accelerator is already registered by this application (the existing
//     binding is left untouched - this is "error and preserve" semantics, see
//     the package documentation on conflicting shortcuts);
//   - the application is already running and the operating system rejects the
//     registration, typically because another application has already claimed
//     the shortcut. Behaviour in this case is platform dependent; see the
//     documentation.
func (m *GlobalShortcutManager) Register(accelerator string, callback func()) error {
	if callback == nil {
		return fmt.Errorf("global shortcut callback must not be nil")
	}

	parsed, err := parseAccelerator(accelerator)
	if err != nil {
		return fmt.Errorf("invalid global shortcut %q: %w", accelerator, err)
	}
	name := parsed.String()

	m.mu.Lock()
	if _, exists := m.byName[name]; exists {
		m.mu.Unlock()
		return fmt.Errorf("global shortcut %q is already registered", name)
	}
	shortcut := &globalShortcut{
		id:          m.nextID,
		accelerator: name,
		parsed:      parsed,
		callback:    callback,
	}
	// Reserve the id and slots before the (blocking) native call so that a
	// concurrent Register of the same accelerator loses the race cleanly.
	m.nextID++
	m.byName[name] = shortcut
	m.byID[shortcut.id] = shortcut

	if !m.started {
		// The application is not running yet: defer the OS binding until start.
		m.pending = append(m.pending, shortcut)
		m.mu.Unlock()
		return nil
	}
	impl := m.getImpl()
	m.mu.Unlock()

	if regErr := InvokeSyncWithError(func() error {
		return impl.register(shortcut.id, parsed)
	}); regErr != nil {
		// Roll back the reservation so the accelerator can be retried later.
		m.mu.Lock()
		delete(m.byName, name)
		delete(m.byID, shortcut.id)
		m.mu.Unlock()
		return fmt.Errorf("failed to register global shortcut %q: %w", name, regErr)
	}
	return nil
}

// flushPending binds every shortcut that was registered before the application
// started. It is called once, on the main thread, during application startup.
// OS-level rejections are reported through the application error handler since
// the original Register caller has already returned.
func (m *GlobalShortcutManager) flushPending() {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return
	}
	m.started = true
	pending := m.pending
	m.pending = nil
	var impl globalShortcutImpl
	if len(pending) > 0 {
		impl = m.getImpl()
	}
	m.mu.Unlock()

	for _, shortcut := range pending {
		if err := impl.register(shortcut.id, shortcut.parsed); err != nil {
			m.mu.Lock()
			// Only roll back if it is still the shortcut we registered (it may
			// have been Unregistered in the meantime).
			if current, ok := m.byID[shortcut.id]; ok && current == shortcut {
				delete(m.byName, shortcut.accelerator)
				delete(m.byID, shortcut.id)
			}
			m.mu.Unlock()
			m.app.handleError(fmt.Errorf("failed to register global shortcut %q: %w", shortcut.accelerator, err))
		}
	}
}

// Unregister releases the given accelerator. It returns an error if the
// accelerator is not currently registered or if the OS rejects the request.
func (m *GlobalShortcutManager) Unregister(accelerator string) error {
	parsed, err := parseAccelerator(accelerator)
	if err != nil {
		return fmt.Errorf("invalid global shortcut %q: %w", accelerator, err)
	}
	name := parsed.String()

	m.mu.Lock()
	shortcut, exists := m.byName[name]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("global shortcut %q is not registered", name)
	}
	delete(m.byName, name)
	delete(m.byID, shortcut.id)
	started := m.started
	impl := m.getImpl()
	m.mu.Unlock()

	if !started {
		// Never bound with the OS yet; just drop it from the pending list.
		m.removePending(shortcut)
		return nil
	}
	return InvokeSyncWithError(func() error {
		return impl.unregister(shortcut.id)
	})
}

// removePending drops a shortcut from the pending queue (used when a shortcut is
// unregistered before the application has started).
func (m *GlobalShortcutManager) removePending(shortcut *globalShortcut) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, p := range m.pending {
		if p == shortcut {
			m.pending = append(m.pending[:i], m.pending[i+1:]...)
			return
		}
	}
}

// UnregisterAll releases every global shortcut registered by this application.
// It is called automatically during application shutdown.
func (m *GlobalShortcutManager) UnregisterAll() error {
	m.mu.Lock()
	hadShortcuts := len(m.byName) > 0
	m.byName = make(map[string]*globalShortcut)
	m.byID = make(map[int]*globalShortcut)
	m.pending = nil
	impl := m.impl
	started := m.started
	m.mu.Unlock()

	// Nothing was ever bound with the OS: avoid forcing the platform impl into
	// existence just to tear nothing down.
	if impl == nil || !started || !hadShortcuts {
		return nil
	}
	return InvokeSyncWithError(impl.unregisterAll)
}

// IsRegistered reports whether the given accelerator is currently registered by
// this application. It returns false for accelerators that cannot be parsed.
//
// On Wayland the returned value reflects what the application requested, not
// necessarily what the compositor ultimately bound; see the package
// documentation on the global shortcuts portal.
func (m *GlobalShortcutManager) IsRegistered(accelerator string) bool {
	parsed, err := parseAccelerator(accelerator)
	if err != nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.byName[parsed.String()]
	return exists
}

// GetAll returns the canonical accelerator strings of all shortcuts currently
// registered by this application, sorted for stable output.
func (m *GlobalShortcutManager) GetAll() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, 0, len(m.byName))
	for name := range m.byName {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// dispatch is called by the platform implementation (from the native event
// handler) when a shortcut with the given id fires. The user callback is run on
// its own goroutine so that it cannot block the platform's main event loop -
// callbacks that need to touch the UI should marshal onto the main thread
// themselves (for example via InvokeSync).
func (m *GlobalShortcutManager) dispatch(id int) {
	m.mu.Lock()
	shortcut, ok := m.byID[id]
	m.mu.Unlock()
	if !ok || shortcut.callback == nil {
		return
	}
	go func() {
		defer handlePanic()
		shortcut.callback()
	}()
}
