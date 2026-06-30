//go:build linux && cgo && !android && !server

package application

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	portalService    = "org.freedesktop.portal.Desktop"
	portalPath       = "/org/freedesktop/portal/desktop"
	portalShortcutIf = "org.freedesktop.portal.GlobalShortcuts"
	portalRequestIf  = "org.freedesktop.portal.Request"
)

// portalShortcut is one shortcut as known to the portal backend.
type portalShortcut struct {
	id      int
	trigger string // preferred trigger, e.g. "CTRL+SHIFT+a"
	desc    string
}

// portalGlobalShortcuts implements globalShortcutImpl on Wayland using the XDG
// Desktop Portal's org.freedesktop.portal.GlobalShortcuts interface.
//
// IMPORTANT semantic difference from the X11/macOS/Windows backends: the portal
// only takes a *preferred* trigger. The compositor - and ultimately the user -
// decides the final key combination, and may change or reject it. The callback
// still fires when the (possibly remapped) shortcut is activated, but the exact
// keys are not guaranteed to match what was requested. IsRegistered/GetAll
// therefore report what the application asked for, not what the compositor
// bound.
//
// All D-Bus work runs on a dedicated goroutine so that register/unregister
// never block the UI thread on a (potentially interactive) portal call.
type portalGlobalShortcuts struct {
	manager *GlobalShortcutManager

	mu        sync.Mutex
	desired   map[int]portalShortcut // current desired set, keyed by numeric id
	cmds      chan func()            // serialized onto the worker goroutine
	tokenSeq  int
	sessionWG sync.Once
}

func newPortalGlobalShortcuts(manager *GlobalShortcutManager) globalShortcutImpl {
	return &portalGlobalShortcuts{
		manager: manager,
		desired: make(map[int]portalShortcut),
		cmds:    make(chan func(), 32),
	}
}

func (p *portalGlobalShortcuts) ensureWorker() {
	p.sessionWG.Do(func() {
		go p.worker()
	})
}

func (p *portalGlobalShortcuts) register(id int, accel *accelerator) error {
	p.mu.Lock()
	p.desired[id] = portalShortcut{
		id:      id,
		trigger: portalTrigger(accel),
		desc:    accel.String(),
	}
	p.mu.Unlock()
	p.ensureWorker()
	p.requestRebind()
	return nil
}

func (p *portalGlobalShortcuts) unregister(id int) error {
	p.mu.Lock()
	_, ok := p.desired[id]
	delete(p.desired, id)
	p.mu.Unlock()
	if ok {
		p.requestRebind()
	}
	return nil
}

func (p *portalGlobalShortcuts) unregisterAll() error {
	p.mu.Lock()
	p.desired = make(map[int]portalShortcut)
	p.mu.Unlock()
	p.requestRebind()
	return nil
}

// requestRebind asks the worker to (re)bind the full desired set. Binds are
// debounced so that the burst of register() calls applications make at startup
// collapses into a single portal interaction.
func (p *portalGlobalShortcuts) requestRebind() {
	select {
	case p.cmds <- p.rebind:
	default:
		// A rebind is already queued; the worker reads the latest desired set.
	}
}

func (p *portalGlobalShortcuts) snapshot() []portalShortcut {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]portalShortcut, 0, len(p.desired))
	for _, s := range p.desired {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].id < out[j].id })
	return out
}

// worker owns the D-Bus connection and serializes all portal interaction.
func (p *portalGlobalShortcuts) worker() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		p.manager.app.handleError(fmt.Errorf("global shortcuts: cannot connect to session bus: %w", err))
		return
	}
	defer conn.Close()

	w := &portalWorker{
		p:        p,
		conn:     conn,
		pending:  make(map[dbus.ObjectPath]chan portalResponse),
		shortcut: make(map[string]int),
	}

	sender := conn.Names()[0]
	w.senderToken = strings.ReplaceAll(strings.TrimPrefix(sender, ":"), ".", "_")

	// One signal channel for everything: Request.Response (method results) and
	// GlobalShortcuts.Activated (shortcut presses). A dedicated goroutine reads
	// it so that synchronous portal calls (which block waiting for their
	// Response) do not starve signal delivery.
	sigs := make(chan *dbus.Signal, 32)
	conn.Signal(sigs)
	_ = conn.AddMatchSignal(dbus.WithMatchInterface(portalRequestIf), dbus.WithMatchMember("Response"))
	_ = conn.AddMatchSignal(dbus.WithMatchInterface(portalShortcutIf), dbus.WithMatchMember("Activated"))
	go func() {
		for sig := range sigs {
			w.handleSignal(sig)
		}
	}()

	if err := w.createSession(); err != nil {
		p.manager.app.handleError(fmt.Errorf("global shortcuts: portal session failed: %w", err))
		return
	}

	// Debounce timer for coalescing the startup burst of register() calls into
	// a single BindShortcuts (and thus a single portal interaction).
	var debounce *time.Timer
	var debounceC <-chan time.Time
	rebindPending := false

	for {
		select {
		case _, ok := <-p.cmds:
			if !ok {
				return
			}
			rebindPending = true
			if debounce == nil {
				debounce = time.NewTimer(150 * time.Millisecond)
				debounceC = debounce.C
			} else {
				debounce.Reset(150 * time.Millisecond)
			}
		case <-debounceC:
			if rebindPending {
				rebindPending = false
				w.bindShortcuts()
			}
		}
	}
}

// rebind is the command pushed onto cmds; the actual work happens in the worker
// loop via the debounce timer.
func (p *portalGlobalShortcuts) rebind() {}

type portalResponse struct {
	code    uint32
	results map[string]dbus.Variant
}

type portalWorker struct {
	p           *portalGlobalShortcuts
	conn        *dbus.Conn
	senderToken string
	sessionPath dbus.ObjectPath

	mu       sync.Mutex
	pending  map[dbus.ObjectPath]chan portalResponse
	shortcut map[string]int // portal shortcut id string -> numeric id
}

func (w *portalWorker) nextToken(prefix string) string {
	w.p.mu.Lock()
	w.p.tokenSeq++
	seq := w.p.tokenSeq
	w.p.mu.Unlock()
	return fmt.Sprintf("wails_%s_%d", prefix, seq)
}

func (w *portalWorker) requestPath(token string) dbus.ObjectPath {
	return dbus.ObjectPath("/org/freedesktop/portal/desktop/request/" + w.senderToken + "/" + token)
}

// call invokes a portal method that follows the Request/Response pattern and
// blocks (on the worker goroutine) until the Response signal arrives.
func (w *portalWorker) call(method string, options map[string]dbus.Variant, args ...interface{}) (portalResponse, error) {
	token := w.nextToken("req")
	options["handle_token"] = dbus.MakeVariant(token)
	expected := w.requestPath(token)

	respCh := make(chan portalResponse, 1)
	w.mu.Lock()
	w.pending[expected] = respCh
	w.mu.Unlock()
	defer func() {
		w.mu.Lock()
		delete(w.pending, expected)
		w.mu.Unlock()
	}()

	// Build the argument list: portal request methods take their own args
	// followed by the options dictionary last.
	callArgs := append(append([]interface{}{}, args...), options)
	obj := w.conn.Object(portalService, portalPath)
	var handle dbus.ObjectPath
	if err := obj.Call(portalShortcutIf+"."+method, 0, callArgs...).Store(&handle); err != nil {
		return portalResponse{}, err
	}
	// The real request path is the one the portal returns; results arrive there.
	if handle != expected {
		w.mu.Lock()
		w.pending[handle] = respCh
		w.mu.Unlock()
		defer func() {
			w.mu.Lock()
			delete(w.pending, handle)
			w.mu.Unlock()
		}()
	}

	select {
	case resp := <-respCh:
		return resp, nil
	case <-time.After(30 * time.Second):
		return portalResponse{}, fmt.Errorf("timed out waiting for portal response to %s", method)
	}
}

func (w *portalWorker) createSession() error {
	sessionToken := w.nextToken("session")
	opts := map[string]dbus.Variant{
		"session_handle_token": dbus.MakeVariant(sessionToken),
	}
	resp, err := w.call("CreateSession", opts)
	if err != nil {
		return err
	}
	if resp.code != 0 {
		return fmt.Errorf("CreateSession refused (response %d)", resp.code)
	}
	if v, ok := resp.results["session_handle"]; ok {
		if s, ok := v.Value().(string); ok {
			w.sessionPath = dbus.ObjectPath(s)
		}
	}
	if w.sessionPath == "" {
		return fmt.Errorf("portal did not return a session handle")
	}
	return nil
}

func (w *portalWorker) bindShortcuts() {
	if w.sessionPath == "" {
		return
	}
	shortcuts := w.p.snapshot()

	// Build the a(sa{sv}) shortcuts argument and refresh the id mapping.
	newShortcut := make(map[string]int, len(shortcuts))
	type entry struct {
		ID    string
		Props map[string]dbus.Variant
	}
	list := make([]entry, 0, len(shortcuts))
	for _, s := range shortcuts {
		idStr := strconv.Itoa(s.id)
		newShortcut[idStr] = s.id
		props := map[string]dbus.Variant{
			"description": dbus.MakeVariant(s.desc),
		}
		if s.trigger != "" {
			props["preferred_trigger"] = dbus.MakeVariant(s.trigger)
		}
		list = append(list, entry{ID: idStr, Props: props})
	}
	w.mu.Lock()
	w.shortcut = newShortcut
	w.mu.Unlock()

	resp, err := w.call("BindShortcuts", map[string]dbus.Variant{}, w.sessionPath, list, "")
	if err != nil {
		w.p.manager.app.handleError(fmt.Errorf("global shortcuts: BindShortcuts failed: %w", err))
		return
	}
	if resp.code != 0 {
		// response 1 = user cancelled the portal dialog, 2 = ended.
		w.p.manager.app.handleError(fmt.Errorf("global shortcuts: portal did not grant shortcuts (response %d); they will not fire", resp.code))
	}
}

func (w *portalWorker) handleSignal(sig *dbus.Signal) {
	switch {
	case strings.HasSuffix(sig.Name, ".Response"):
		w.mu.Lock()
		ch, ok := w.pending[sig.Path]
		w.mu.Unlock()
		if ok {
			var resp portalResponse
			if len(sig.Body) >= 1 {
				if code, ok := sig.Body[0].(uint32); ok {
					resp.code = code
				}
			}
			if len(sig.Body) >= 2 {
				if results, ok := sig.Body[1].(map[string]dbus.Variant); ok {
					resp.results = results
				}
			}
			select {
			case ch <- resp:
			default:
			}
		}
	case strings.HasSuffix(sig.Name, ".Activated"):
		// Activated(o session_handle, s shortcut_id, t timestamp, a{sv} options)
		if len(sig.Body) < 2 {
			return
		}
		idStr, ok := sig.Body[1].(string)
		if !ok {
			return
		}
		w.mu.Lock()
		id, ok := w.shortcut[idStr]
		w.mu.Unlock()
		if ok {
			w.p.manager.dispatch(id)
		}
	}
}

// portalTrigger renders an accelerator as a portal "preferred_trigger" string.
// The portal trigger syntax uses the modifier names CTRL, ALT, SHIFT and LOGO
// joined to the key with "+". This is only a preference; the compositor may
// bind a different combination.
func portalTrigger(accel *accelerator) string {
	var parts []string
	hasCtrl, hasAlt, hasShift, hasSuper := false, false, false, false
	for _, m := range accel.Modifiers {
		switch m {
		case CmdOrCtrlKey, ControlKey:
			hasCtrl = true
		case OptionOrAltKey:
			hasAlt = true
		case ShiftKey:
			hasShift = true
		case SuperKey:
			hasSuper = true
		}
	}
	if hasCtrl {
		parts = append(parts, "CTRL")
	}
	if hasAlt {
		parts = append(parts, "ALT")
	}
	if hasShift {
		parts = append(parts, "SHIFT")
	}
	if hasSuper {
		parts = append(parts, "LOGO")
	}
	key, ok := x11KeysymNames[accel.Key]
	if !ok {
		key = accel.Key
	}
	parts = append(parts, key)
	return strings.Join(parts, "+")
}
