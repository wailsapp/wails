//go:build linux && !android

/*
Portions of this code are derived from the project:
- https://github.com/fyne-io/systray
*/
package application

import "C"
import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/wailsapp/wails/v3/internal/dbus/menu"
	"github.com/wailsapp/wails/v3/internal/dbus/notifier"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

const (
	itemPath = "/StatusNotifierItem"
	menuPath = "/StatusNotifierMenu"
)

type linuxSystemTray struct {
	parent *SystemTray

	id    uint
	label string
	icon  []byte
	menu  *Menu

	iconPosition   IconPosition
	isTemplateIcon bool

	quitChan  chan struct{}
	conn      *dbus.Conn
	props     *prop.Properties
	menuProps *prop.Properties

	menuVersion uint32
	itemMap     map[int32]*systrayMenuItem
	tooltip     string

	lastClickX int
	lastClickY int
}

func (s *linuxSystemTray) getScreen() (*Screen, error) {
	_, _, result := getMousePosition()
	return result, nil
}

// dbusMenu is a named struct to map into generated bindings.
// It represents the layout of a menu item
type dbusMenu = struct {
	V0 int32                   // items' unique id
	V1 map[string]dbus.Variant // layout properties
	V2 []dbus.Variant          // child menu(s)
}

// systrayMenuItem is an implementation of the menuItemImpl interface
type systrayMenuItem struct {
	sysTray  *linuxSystemTray
	menuItem *MenuItem
	dbusItem *dbusMenu
}

func (s *systrayMenuItem) setBitmap(data []byte) {
	s.dbusItem.V1["icon-data"] = dbus.MakeVariant(data)
	s.sysTray.update(s)
}

func (s *systrayMenuItem) setTooltip(v string) {
	s.dbusItem.V1["tooltip"] = dbus.MakeVariant(v)
	s.sysTray.update(s)
}

func (s *systrayMenuItem) setLabel(v string) {
	s.dbusItem.V1["label"] = dbus.MakeVariant(v)
	s.sysTray.update(s)
}

func (s *systrayMenuItem) setDisabled(disabled bool) {
	v := dbus.MakeVariant(!disabled)
	if s.dbusItem.V1["toggle-state"] != v {
		s.dbusItem.V1["enabled"] = v
		s.sysTray.update(s)
	}
}

func (s *systrayMenuItem) destroy() {}

func (s *systrayMenuItem) setChecked(checked bool) {
	v := dbus.MakeVariant(0)
	if checked {
		v = dbus.MakeVariant(1)
	}
	if s.dbusItem.V1["toggle-state"] != v {
		s.dbusItem.V1["toggle-state"] = v
		s.sysTray.update(s)
	}
}

func (s *systrayMenuItem) setAccelerator(accelerator *accelerator) {}
func (s *systrayMenuItem) setHidden(hidden bool) {
	s.dbusItem.V1["visible"] = dbus.MakeVariant(!hidden)
	s.sysTray.update(s)
}

func (s *systrayMenuItem) dbus() *dbusMenu {
	item := &dbusMenu{
		V0: int32(s.menuItem.id),
		V1: map[string]dbus.Variant{},
		V2: []dbus.Variant{},
	}
	return item
}

func (s *linuxSystemTray) setIconPosition(position IconPosition) {
	s.iconPosition = position
}

func (s *linuxSystemTray) processMenu(menu *Menu, parentId int32) {
	parentItem, ok := s.itemMap[int32(parentId)]
	if !ok {
		return
	}
	parent := parentItem.dbusItem

	for _, item := range menu.items {
		menuItem := &dbusMenu{
			V0: int32(item.id),
			V1: map[string]dbus.Variant{},
			V2: []dbus.Variant{},
		}
		item.impl = &systrayMenuItem{
			sysTray:  s,
			menuItem: item,
			dbusItem: menuItem,
		}
		s.itemMap[int32(item.id)] = item.impl.(*systrayMenuItem)

		menuItem.V1["enabled"] = dbus.MakeVariant(!item.disabled)
		menuItem.V1["visible"] = dbus.MakeVariant(!item.hidden)
		if item.label != "" {
			menuItem.V1["label"] = dbus.MakeVariant(item.label)
		}
		if item.bitmap != nil {
			menuItem.V1["icon-data"] = dbus.MakeVariant(item.bitmap)
		}
		switch item.itemType {
		case checkbox:
			menuItem.V1["toggle-type"] = dbus.MakeVariant("checkmark")
			v := dbus.MakeVariant(0)
			if item.checked {
				v = dbus.MakeVariant(1)
			}
			menuItem.V1["toggle-state"] = v
		case submenu:
			menuItem.V1["children-display"] = dbus.MakeVariant("submenu")
			s.processMenu(item.submenu, int32(item.id))
		case text:
		case radio:
			menuItem.V1["toggle-type"] = dbus.MakeVariant("radio")
			v := dbus.MakeVariant(0)
			if item.checked {
				v = dbus.MakeVariant(1)
			}
			menuItem.V1["toggle-state"] = v
		case separator:
			menuItem.V1["type"] = dbus.MakeVariant("separator")
		}

		parent.V2 = append(parent.V2, dbus.MakeVariant(menuItem))
	}
}

func (s *linuxSystemTray) refresh() {
	s.menuVersion++
	if err := s.menuProps.Set("com.canonical.dbusmenu", "Version",
		dbus.MakeVariant(s.menuVersion)); err != nil {
		globalApplication.error("systray error: failed to update menu version: %w", err)
		return
	}
	if err := menu.Emit(s.conn, &menu.Dbusmenu_LayoutUpdatedSignal{
		Path: menuPath,
		Body: &menu.Dbusmenu_LayoutUpdatedSignalBody{
			Revision: s.menuVersion,
		},
	}); err != nil {
		globalApplication.error("systray error: failed to emit layout updated signal: %w", err)
	}
}

func (s *linuxSystemTray) setMenu(menu *Menu) {
	s.itemMap = map[int32]*systrayMenuItem{}
	s.itemMap[0] = &systrayMenuItem{
		menuItem: nil,
		dbusItem: &dbusMenu{
			V0: int32(0),
			V1: map[string]dbus.Variant{},
			V2: []dbus.Variant{},
		},
	}
	if menu != nil {
		menu.processRadioGroups()
		s.processMenu(menu, 0)
	}
	s.menu = menu
	s.refresh()
}

func (s *linuxSystemTray) positionWindow(window Window, offset int) error {
	_, _, currentScreen := getMousePosition()
	if currentScreen == nil {
		return fmt.Errorf("unable to get screen information")
	}

	screenX := currentScreen.X
	screenY := currentScreen.Y
	screenWidth := currentScreen.Size.Width
	screenHeight := currentScreen.Size.Height
	windowWidth := window.Width()
	windowHeight := window.Height()

	if isTilingWM() {
		newX := screenX + (screenWidth-windowWidth)/2
		newY := screenY + (screenHeight-windowHeight)/2
		window.SetPosition(newX, newY)
		return nil
	}

	clickX, clickY := s.lastClickX, s.lastClickY
	if clickX == 0 && clickY == 0 {
		if cx, cy, ok := getCursorPositionFromCompositor(); ok {
			clickX, clickY = cx, cy
		} else {
			clickX = screenX + screenWidth/2
			clickY = screenY + screenHeight/2
		}
	}

	newX := clickX - (windowWidth / 2)
	if newX < screenX {
		newX = screenX
	}
	if newX+windowWidth > screenX+screenWidth {
		newX = screenX + screenWidth - windowWidth
	}

	relativeY := clickY - screenY
	topThreshold := screenHeight / 5
	bottomThreshold := screenHeight * 4 / 5

	var newY int
	if relativeY < topThreshold {
		newY = clickY + offset
	} else if relativeY > bottomThreshold {
		newY = clickY - windowHeight - offset
	} else {
		newY = clickY - (windowHeight / 2)
	}

	if newY < screenY {
		newY = screenY
	}
	if newY+windowHeight > screenY+screenHeight {
		newY = screenY + screenHeight - windowHeight
	}

	window.SetPosition(newX, newY)
	return nil
}

func (s *linuxSystemTray) bounds() (*Rect, error) {

	// Best effort guess at the screen bounds

	return &Rect{}, nil

}

func (s *linuxSystemTray) run() {
	conn, err := dbus.SessionBus()
	if err != nil {
		globalApplication.error("systray error: failed to connect to DBus: %w\n", err)
		return
	}
	err = notifier.ExportStatusNotifierItem(conn, itemPath, s)
	if err != nil {
		globalApplication.error("systray error: failed to export status notifier item: %w\n", err)
	}

	err = menu.ExportDbusmenu(conn, menuPath, s)
	if err != nil {
		globalApplication.error("systray error: failed to export status notifier menu: %w", err)
		return
	}

	name := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid()) // register id 1 for this process
	_, err = conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		globalApplication.error("systray error: failed to request name: %w", err)
		// it's not critical error: continue
	}
	props, err := prop.Export(conn, itemPath, s.createPropSpec())
	if err != nil {
		globalApplication.error("systray error: failed to export notifier item properties to bus: %w", err)
		return
	}
	menuProps, err := prop.Export(conn, menuPath, s.createMenuPropSpec())
	if err != nil {
		globalApplication.error("systray error: failed to export notifier menu properties to bus: %w", err)
		return
	}

	s.conn = conn
	s.props = props
	s.menuProps = menuProps

	node := introspect.Node{
		Name: itemPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			notifier.IntrospectDataStatusNotifierItem,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&node), itemPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		globalApplication.error("systray error: failed to export node introspection: %w", err)
		return
	}
	menuNode := introspect.Node{
		Name: menuPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			menu.IntrospectDataDbusmenu,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&menuNode), menuPath,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		globalApplication.error("systray error: failed to export menu node introspection: %w", err)
		return
	}
	s.setLabel(s.label)
	go func() {
		defer handlePanic()
		s.register()

		if err := conn.AddMatchSignal(
			dbus.WithMatchObjectPath("/org/freedesktop/DBus"),
			dbus.WithMatchInterface("org.freedesktop.DBus"),
			dbus.WithMatchSender("org.freedesktop.DBus"),
			dbus.WithMatchMember("NameOwnerChanged"),
			dbus.WithMatchArg(0, "org.kde.StatusNotifierWatcher"),
		); err != nil {
			globalApplication.error("systray error: failed to register signal matching: %w", err)
			return
		}

		sc := make(chan *dbus.Signal, 10)
		conn.Signal(sc)

		for {
			select {
			case sig := <-sc:
				if sig == nil {
					return // We get a nil signal when closing the window.
				}
				// sig.Body has the args, which are [name old_owner new_owner]
				if sig.Body[2] != "" {
					s.register()
				}

			case <-s.quitChan:
				return
			}
		}
	}()

	if s.parent.label != "" {
		s.setLabel(s.parent.label)
	}

	if s.parent.tooltip != "" {
		s.setTooltip(s.parent.tooltip)
	}
	s.setMenu(s.menu)
}

func (s *linuxSystemTray) setTooltip(_ string) {
	// TBD
}

func (s *linuxSystemTray) setIcon(icon []byte) {

	s.icon = icon

	iconPx, err := iconToPX(icon)
	if err != nil {
		globalApplication.error("systray error: failed to convert icon to PX: %w", err)
		return
	}
	s.props.SetMust("org.kde.StatusNotifierItem", "IconPixmap", []PX{iconPx})

	if s.conn == nil {
		return
	}

	err = notifier.Emit(s.conn, &notifier.StatusNotifierItem_NewIconSignal{
		Path: itemPath,
		Body: &notifier.StatusNotifierItem_NewIconSignalBody{},
	})
	if err != nil {
		globalApplication.error("systray error: failed to emit new icon signal: %w", err)
		return
	}
}

func (s *linuxSystemTray) setDarkModeIcon(icon []byte) {
	s.setIcon(icon)
}

func (s *linuxSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	s.setIcon(icon)
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	label := s.label
	if label == "" {
		label = "Wails"
	}

	return &linuxSystemTray{
		parent:         s,
		id:             s.id,
		label:          label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
		quitChan:       make(chan struct{}),
		menuVersion:    1,
	}
}

func (s *linuxSystemTray) openMenu() {
	// FIXME: Emit com.canonical to open?
	globalApplication.info("systray error: openMenu not implemented on Linux")
}

func (s *linuxSystemTray) setLabel(label string) {
	s.label = label

	if err := s.props.Set("org.kde.StatusNotifierItem", "Title", dbus.MakeVariant(label)); err != nil {
		globalApplication.error("systray error: failed to set Title prop: %w", err)
		return
	}

	if s.conn == nil {
		return
	}

	if err := notifier.Emit(s.conn, &notifier.StatusNotifierItem_NewTitleSignal{
		Path: itemPath,
		Body: &notifier.StatusNotifierItem_NewTitleSignalBody{},
	}); err != nil {
		globalApplication.error("systray error: failed to emit new title signal: %w", err)
		return
	}

}

func (s *linuxSystemTray) destroy() {
	close(s.quitChan)
}

func (s *linuxSystemTray) createMenuPropSpec() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		"com.canonical.dbusmenu": {
			// update version each time we change something
			"Version": {
				Value:    s.menuVersion,
				Writable: true,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"TextDirection": {
				Value:    "ltr",
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"Status": {
				Value:    "normal",
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"IconThemePath": {
				Value:    []string{},
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
		},
	}
}

func (s *linuxSystemTray) createPropSpec() map[string]map[string]*prop.Prop {
	props := map[string]*prop.Prop{
		"Status": {
			Value:    "Active", // Passive, Active or NeedsAttention
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"Title": {
			Value:    s.label,
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"Id": {
			Value:    s.label,
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"Category": {
			Value:    "ApplicationStatus",
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"IconData": {
			Value:    "",
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},

		"IconName": {
			Value:    "",
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"IconThemePath": {
			Value:    "",
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"ItemIsMenu": {
			Value:    false,
			Writable: false,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
		"ToolTip": {
			Value:    tooltip{V2: s.label},
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: nil,
		},
	}

	if s.icon == nil {
		s.icon = icons.WailsLogoWhiteTransparent
	}
	if iconPx, err := iconToPX(s.icon); err == nil {
		props["IconPixmap"] = &prop.Prop{
			Value:    []PX{iconPx},
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: nil,
		}
	}

	if s.menu != nil {
		props["Menu"] = &prop.Prop{
			Value:    dbus.ObjectPath(menuPath),
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: nil,
		}
	}

	return map[string]map[string]*prop.Prop{
		"org.kde.StatusNotifierItem": props,
	}
}

func (s *linuxSystemTray) update(i *systrayMenuItem) {
	s.itemMap[int32(i.menuItem.id)] = i
	s.refresh()
}

func (s *linuxSystemTray) register() bool {
	obj := s.conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	call := obj.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, itemPath)
	if call.Err != nil {
		globalApplication.error("systray error: failed to register: %w", call.Err)
		return false
	}

	return true
}

type PX struct {
	W, H int
	Pix  []byte
}

func iconToPX(icon []byte) (PX, error) {
	img, err := pngToImage(icon)
	if err != nil {
		return PX{}, err
	}
	w, h, bytes := ToARGB(img)
	return PX{
		W:   w,
		H:   h,
		Pix: bytes,
	}, nil
}

// AboutToShow is an implementation of the com.canonical.dbusmenu.AboutToShow method.
func (s *linuxSystemTray) AboutToShow(id int32) (needUpdate bool, err *dbus.Error) {
	return
}

// AboutToShowGroup is an implementation of the com.canonical.dbusmenu.AboutToShowGroup method.
func (s *linuxSystemTray) AboutToShowGroup(ids []int32) (updatesNeeded []int32, idErrors []int32, err *dbus.Error) {
	return
}

// GetProperty is an implementation of the com.canonical.dbusmenu.GetProperty method.
func (s *linuxSystemTray) GetProperty(id int32, name string) (value dbus.Variant, err *dbus.Error) {
	if item, ok := s.itemMap[id]; ok {
		if p, ok := item.dbusItem.V1[name]; ok {
			return p, nil
		}
	}
	return
}

func (s *linuxSystemTray) Event(id int32, eventID string, data dbus.Variant, timestamp uint32) (err *dbus.Error) {
	globalApplication.debug("systray Event called", "id", id, "eventID", eventID, "lastClick", fmt.Sprintf("(%d,%d)", s.lastClickX, s.lastClickY))
	switch eventID {
	case "clicked":
		if item, ok := s.itemMap[id]; ok {
			InvokeAsync(item.menuItem.handleClick)
		}
	case "opened":
		if s.parent.clickHandler != nil {
			s.parent.clickHandler()
		}
		if s.parent.onMenuOpen != nil {
			s.parent.onMenuOpen()
		}
	case "closed":
		if s.parent.onMenuClose != nil {
			s.parent.onMenuClose()
		}
	}
	return
}

// EventGroup is an implementation of the com.canonical.dbusmenu.EventGroup method.
func (s *linuxSystemTray) EventGroup(events []struct {
	V0 int32
	V1 string
	V2 dbus.Variant
	V3 uint32
}) (idErrors []int32, err *dbus.Error) {
	for _, event := range events {
		fmt.Printf("EventGroup: %v, %v, %v, %v\n", event.V0, event.V1, event.V2, event.V3)
		if event.V1 == "clicked" {
			item, ok := s.itemMap[event.V0]
			if ok {
				InvokeAsync(item.menuItem.handleClick)
			}
		}
	}
	return
}

// GetGroupProperties is an implementation of the com.canonical.dbusmenu.GetGroupProperties method.
func (s *linuxSystemTray) GetGroupProperties(ids []int32, propertyNames []string) (properties []struct {
	V0 int32
	V1 map[string]dbus.Variant
}, err *dbus.Error) {
	// FIXME: RLock?
	/*	instance.menuLock.Lock()
		defer instance.menuLock.Unlock()
	*/
	for _, id := range ids {
		if m, ok := s.itemMap[id]; ok {
			p := struct {
				V0 int32
				V1 map[string]dbus.Variant
			}{
				V0: m.dbusItem.V0,
				V1: make(map[string]dbus.Variant, len(m.dbusItem.V1)),
			}
			for k, v := range m.dbusItem.V1 {
				p.V1[k] = v
			}
			properties = append(properties, p)
		}
	}
	return properties, nil
}

// GetLayout is an implementation of the com.canonical.dbusmenu.GetLayout method.
func (s *linuxSystemTray) GetLayout(parentID int32, recursionDepth int32, propertyNames []string) (revision uint32, layout dbusMenu, err *dbus.Error) {
	// FIXME: RLock?
	if m, ok := s.itemMap[parentID]; ok {
		return s.menuVersion, *m.dbusItem, nil
	}

	return
}

func (s *linuxSystemTray) Activate(x int32, y int32) (err *dbus.Error) {
	s.lastClickX = int(x)
	s.lastClickY = int(y)
	globalApplication.debug("systray Activate called", "x", x, "y", y)
	if s.parent.clickHandler != nil {
		s.parent.clickHandler()
	}
	return
}

func (s *linuxSystemTray) ContextMenu(x int32, y int32) (err *dbus.Error) {
	s.lastClickX = int(x)
	s.lastClickY = int(y)
	return nil
}

func (s *linuxSystemTray) Scroll(delta int32, orientation string) (err *dbus.Error) {
	return
}

func (s *linuxSystemTray) SecondaryActivate(x int32, y int32) (err *dbus.Error) {
	s.lastClickX = int(x)
	s.lastClickY = int(y)
	if s.parent.rightClickHandler != nil {
		s.parent.rightClickHandler()
	} else if s.menu != nil {
		s.parent.OpenMenu()
	}
	return
}

// Show is a no-op for Linux
func (s *linuxSystemTray) Show() {
	// No-op
}

// Hide is a no-op for Linux
func (s *linuxSystemTray) Hide() {
	// No-op
}

// tooltip is our data for a tooltip property.
// Param names need to match the generated code...
type tooltip = struct {
	V0 string // name
	V1 []PX   // icons
	V2 string // title
	V3 string // description
}
