//go:build linux && !android && !server

package application

import (
	"fmt"
	"runtime"
	"syscall"
)

func (g *x11GlobalShortcuts) drainOps() {
	for {
		select {
		case op := <-g.opCh:
			op()
		default:
			return
		}
	}
}

func (g *x11GlobalShortcuts) modMask(accel *accelerator) uint {
	var mask uint
	for _, m := range accel.Modifiers {
		switch m {
		case CmdOrCtrlKey, ControlKey:
			mask |= x11ControlMask
		case OptionOrAltKey:
			mask |= x11Mod1Mask
		case ShiftKey:
			mask |= x11ShiftMask
		case SuperKey:
			mask |= x11Mod4Mask
		}
	}
	return mask
}

// run executes fn on the event-loop goroutine and waits for it to complete.
func (g *x11GlobalShortcuts) run(fn func()) {
	done := make(chan struct{})
	g.opCh <- func() {
		fn()
		close(done)
	}
	g.wake()
	<-done
}

func (g *x11GlobalShortcuts) wake() {
	var b [1]byte
	_, _ = syscall.Write(g.wakeW, b[:])
}

func (g *x11GlobalShortcuts) eventLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for {
		r, keycode, state := g.waitForEvent()
		switch r {
		case 0:
			g.drainOps()
		case 1:
			b := x11Binding{keycode: uint(keycode), modMask: uint(state) & x11SignificantMask}
			g.mu.RLock()
			id, ok := g.match[b]
			g.mu.RUnlock()
			if ok {
				g.manager.dispatch(id)
			}
		default:
			return
		}
	}
}

func (g *x11GlobalShortcuts) register(id int, accel *accelerator) error {
	if g.startErr != nil {
		return g.startErr
	}
	keysymName, ok := x11KeysymNames[accel.Key]
	if !ok {
		return fmt.Errorf("key %q is not supported as a global shortcut", accel.Key)
	}
	modMask := g.modMask(accel)

	var regErr error
	var binding x11Binding
	g.run(func() {
		keycode := g.keycodeForName(keysymName)
		if keycode == 0 {
			regErr = fmt.Errorf("key %q has no keycode on this keyboard", accel.Key)
			return
		}
		if g.grabKey(uint32(keycode), uint32(modMask)) != 0 {
			regErr = fmt.Errorf("the shortcut is already registered (possibly by another application)")
			return
		}
		binding = x11Binding{keycode: keycode, modMask: modMask}
	})
	if regErr != nil {
		return regErr
	}

	g.mu.Lock()
	g.bindings[id] = binding
	g.match[binding] = id
	g.mu.Unlock()
	return nil
}

func (g *x11GlobalShortcuts) unregister(id int) error {
	g.mu.Lock()
	binding, ok := g.bindings[id]
	if ok {
		delete(g.bindings, id)
		delete(g.match, binding)
	}
	g.mu.Unlock()
	if !ok {
		return nil
	}
	g.run(func() {
		g.ungrabKey(uint32(binding.keycode), uint32(binding.modMask))
	})
	return nil
}

func (g *x11GlobalShortcuts) unregisterAll() error {
	g.mu.Lock()
	bindings := g.bindings
	g.bindings = make(map[int]x11Binding)
	g.match = make(map[x11Binding]int)
	g.mu.Unlock()
	if len(bindings) == 0 {
		return nil
	}
	g.run(func() {
		for _, b := range bindings {
			g.ungrabKey(uint32(b.keycode), uint32(b.modMask))
		}
	})
	return nil
}
