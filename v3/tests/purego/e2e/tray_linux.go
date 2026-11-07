package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// The tray scenario plays the desktop's panel. A Linux tray icon is a
// StatusNotifierItem: the app registers it with the session bus's
// StatusNotifierWatcher, and the panel reads its title, tooltip, icon and
// menu over D-Bus and calls back when the user clicks. Xvfb has no panel, so
// the probe is one: it owns the watcher name, takes the registration, and
// acts on the item as a panel would.

const (
	watcherName  = "org.kde.StatusNotifierWatcher"
	watcherPath  = "/StatusNotifierWatcher"
	itemIface    = "org.kde.StatusNotifierItem"
	itemPath     = "/StatusNotifierItem"
	menuIface    = "com.canonical.dbusmenu"
	propsGet     = "org.freedesktop.DBus.Properties.Get"
	trayWaitTime = 5 * time.Second
)

// fakeWatcher takes the items' registrations.
type fakeWatcher struct{ registered chan string }

func (w *fakeWatcher) RegisterStatusNotifierItem(sender dbus.Sender, service string) *dbus.Error {
	select {
	case w.registered <- string(sender):
	default:
	}
	return nil
}

func (w *fakeWatcher) RegisterStatusNotifierHost(sender dbus.Sender, service string) *dbus.Error {
	return nil
}

// trayItem is one registered item as the panel sees it.
type trayItem struct {
	conn  *dbus.Conn
	owner string
}

func (t trayItem) item() dbus.BusObject { return t.conn.Object(t.owner, itemPath) }
func (t trayItem) menu() dbus.BusObject { return t.conn.Object(t.owner, "/StatusNotifierMenu") }

func (t trayItem) prop(name string) (dbus.Variant, error) {
	var v dbus.Variant
	err := t.item().Call(propsGet, 0, itemIface, name).Store(&v)
	return v, err
}

// title reads the item's Title property.
func (t trayItem) title() string {
	v, err := t.prop("Title")
	if err != nil {
		return "error: " + err.Error()
	}
	s, _ := v.Value().(string)
	return s
}

// tooltip reads the description in the item's ToolTip property.
func (t trayItem) tooltip() string {
	v, err := t.prop("ToolTip")
	if err != nil {
		return "error: " + err.Error()
	}
	// (s a(iiay) s s): icon name, icons, title, description.
	if fields, ok := v.Value().([]any); ok && len(fields) == 4 {
		title, _ := fields[2].(string)
		text, _ := fields[3].(string)
		return title + "|" + text
	}
	return fmt.Sprintf("unexpected %v", v)
}

// layout reads the whole menu and returns its revision and each item's
// label by id.
func (t trayItem) layout() (uint32, map[string]int32, error) {
	var revision uint32
	var root []any
	call := t.menu().Call(menuIface+".GetLayout", 0, int32(0), int32(-1), []string{})
	if call.Err != nil {
		return 0, nil, call.Err
	}
	if err := call.Store(&revision, &root); err != nil {
		return 0, nil, err
	}
	labels := map[string]int32{}
	var walk func(node []any)
	walk = func(node []any) {
		if len(node) != 3 {
			return
		}
		id, _ := node[0].(int32)
		if props, ok := node[1].(map[string]dbus.Variant); ok {
			if l, ok := props["label"].Value().(string); ok {
				labels[l] = id
			}
		}
		children, _ := node[2].([]dbus.Variant)
		for _, c := range children {
			if n, ok := c.Value().([]any); ok {
				walk(n)
			}
		}
	}
	walk(root)
	return revision, labels, nil
}

func (t trayItem) event(id int32, kind string) error {
	return t.menu().Call(menuIface+".Event", 0, id, kind, dbus.MakeVariant(""), uint32(0)).Err
}

// trayHost runs the panel's side of a tray icon's life.
func trayHost(p *probe) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("session bus: %w", err)
	}
	defer conn.Close()
	watcher := &fakeWatcher{registered: make(chan string, 8)}
	if err := conn.Export(watcher, watcherPath, watcherName); err != nil {
		return err
	}
	if reply, err := conn.RequestName(watcherName, dbus.NameFlagDoNotQueue); err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("the watcher name is taken (%v, %v): a real panel is running", reply, err)
	}
	if err := conn.AddMatchSignal(dbus.WithMatchInterface(menuIface), dbus.WithMatchMember("LayoutUpdated")); err != nil {
		return err
	}
	signals := make(chan *dbus.Signal, 64)
	conn.Signal(signals)

	awaitRegistration := func(what string) (trayItem, error) {
		select {
		case owner := <-watcher.registered:
			return trayItem{conn, owner}, nil
		case <-time.After(trayWaitTime):
			return trayItem{}, fmt.Errorf("%s: the tray did not register with the watcher", what)
		}
	}

	menu := application.NewMenu()
	action := menu.Add("Tray action").OnClick(func(*application.Context) { emit("tray-action") })
	menu.AddCheckbox("Tray check", false)
	sub := menu.AddSubmenu("Tray more")
	sub.Add("Tray nested").OnClick(func(*application.Context) { emit("tray-nested") })
	menu.AddSeparator()
	menu.Add("Tray disabled").SetEnabled(false)
	tray := p.app.SystemTray.New().SetIcon(bitmap()).SetMenu(menu)
	tray.SetLabel("E2E tray")
	tray.SetTooltip("E2E tip")
	tray.OnClick(func() { emit("tray-click") })
	// The app is running, so New ran the tray already: no Run.
	destroyed := false
	defer func() {
		if !destroyed {
			tray.Destroy()
		}
	}()

	t, err := awaitRegistration("Run")
	if err != nil {
		return err
	}
	// The panel reads what to show.
	if got := t.title(); got != "E2E tray" {
		return fmt.Errorf("Title %q, want %q", got, "E2E tray")
	}
	if !within(trayWaitTime, func() bool { return strings.HasSuffix(t.tooltip(), "|E2E tip") }) {
		return fmt.Errorf("ToolTip %q, want description %q", t.tooltip(), "E2E tip")
	}
	if v, err := t.prop("IconPixmap"); err != nil || fmt.Sprint(v) == "@a(iiay) []" {
		return fmt.Errorf("IconPixmap %v, %v: want an icon", v, err)
	}
	if v, err := t.prop("Menu"); err != nil || v.Value() != dbus.ObjectPath("/StatusNotifierMenu") {
		return fmt.Errorf("Menu %v, %v", v, err)
	}
	revision, labels, err := t.layout()
	if err != nil {
		return fmt.Errorf("GetLayout: %w", err)
	}
	for _, want := range []string{"Tray action", "Tray check", "Tray more", "Tray nested", "Tray disabled"} {
		if _, ok := labels[want]; !ok {
			return fmt.Errorf("GetLayout has no %q: %v", want, labels)
		}
	}

	// The user clicks: the icon, then menu items.
	drain()
	if err := t.item().Call(itemIface+".Activate", 0, int32(10), int32(10)).Err; err != nil {
		return err
	}
	if err := await("tray-click", trayWaitTime); err != nil {
		return fmt.Errorf("Activate: %w", err)
	}
	for _, m := range []string{"ContextMenu", "SecondaryActivate"} {
		if err := t.item().Call(itemIface+"."+m, 0, int32(10), int32(10)).Err; err != nil {
			return fmt.Errorf("%s: %w", m, err)
		}
	}
	if err := t.item().Call(itemIface+".Scroll", 0, int32(1), "vertical").Err; err != nil {
		return fmt.Errorf("Scroll: %w", err)
	}
	if err := t.menu().Call(menuIface+".AboutToShow", 0, int32(0)).Err; err != nil {
		return fmt.Errorf("AboutToShow: %w", err)
	}
	if err := t.event(0, "opened"); err != nil {
		return err
	}
	if err := t.event(labels["Tray action"], "clicked"); err != nil {
		return err
	}
	if err := await("tray-action", trayWaitTime); err != nil {
		return fmt.Errorf(`Event("clicked") on Tray action: %w`, err)
	}
	if err := t.event(labels["Tray nested"], "clicked"); err != nil {
		return err
	}
	if err := await("tray-nested", trayWaitTime); err != nil {
		return fmt.Errorf(`Event("clicked") on a submenu's item: %w`, err)
	}
	group := []struct {
		ID        int32
		EventID   string
		Data      dbus.Variant
		Timestamp uint32
	}{{labels["Tray action"], "clicked", dbus.MakeVariant(""), 0}}
	if err := t.menu().Call(menuIface+".EventGroup", 0, group).Err; err != nil {
		return fmt.Errorf("EventGroup: %w", err)
	}
	if err := await("tray-action", trayWaitTime); err != nil {
		return fmt.Errorf(`EventGroup "clicked" on Tray action: %w`, err)
	}
	if err := t.event(0, "closed"); err != nil {
		return err
	}
	var label dbus.Variant
	if err := t.menu().Call(menuIface+".GetProperty", 0, labels["Tray action"], "label").Store(&label); err != nil || label.Value() != "Tray action" {
		return fmt.Errorf("GetProperty label: %v, %v", label, err)
	}
	var groupProps []any
	if err := t.menu().Call(menuIface+".GetGroupProperties", 0, []int32{labels["Tray action"]}, []string{}).Store(&groupProps); err != nil || len(groupProps) != 1 {
		return fmt.Errorf("GetGroupProperties: %v, %v", groupProps, err)
	}

	// The app changes the tray: the panel sees each change.
	tray.SetLabel("E2E tray renamed")
	if !within(trayWaitTime, func() bool { return t.title() == "E2E tray renamed" }) {
		return fmt.Errorf("Title %q after SetLabel", t.title())
	}
	tray.SetTooltip("E2E tip 2")
	if !within(trayWaitTime, func() bool { return strings.HasSuffix(t.tooltip(), "|E2E tip 2") }) {
		return fmt.Errorf("ToolTip %q after SetTooltip", t.tooltip())
	}
	tray.SetIcon(bitmap())
	for drained := false; !drained; {
		select {
		case <-signals:
		default:
			drained = true
		}
	}
	action.SetLabel("Tray action renamed")
	menu.Update()
	if err := awaitLayoutUpdated(signals, t.owner); err != nil {
		return err
	}
	newRevision, labels, err := t.layout()
	if err != nil {
		return err
	}
	if _, ok := labels["Tray action renamed"]; !ok || newRevision <= revision {
		return fmt.Errorf("after the label change: revision %d (was %d), labels %v", newRevision, revision, labels)
	}

	// The panel restarts: the tray registers again with the new watcher.
	if _, err := conn.ReleaseName(watcherName); err != nil {
		return err
	}
	if _, err := conn.RequestName(watcherName, dbus.NameFlagDoNotQueue); err != nil {
		return err
	}
	if _, err := awaitRegistration("a restarted panel"); err != nil {
		return err
	}

	// Destroy removes the icon: its connection leaves the bus, which is how
	// a watcher learns an item has gone.
	tray.Destroy()
	destroyed = true
	gone := func() bool {
		var has bool
		err := conn.BusObject().Call("org.freedesktop.DBus.NameHasOwner", 0, t.owner).Store(&has)
		return err == nil && !has
	}
	if !within(trayWaitTime, gone) {
		return fmt.Errorf("the destroyed tray's connection is still on the bus: the panel keeps showing it")
	}
	return nil
}

// awaitLayoutUpdated waits for the menu's LayoutUpdated signal from owner.
func awaitLayoutUpdated(signals chan *dbus.Signal, owner string) error {
	deadline := time.After(trayWaitTime)
	for {
		select {
		case s := <-signals:
			if s != nil && s.Sender == owner && strings.HasSuffix(s.Name, ".LayoutUpdated") {
				return nil
			}
		case <-deadline:
			return fmt.Errorf("no LayoutUpdated signal after the menu changed")
		}
	}
}
