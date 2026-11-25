//go:build linux && gtk3 && !android && !server

package application

// GTK3 widgets have one parent. Each window therefore owns a snapshot of the
// menu tree, while the original Menu remains the source of truth.
var gtk3MenuViews = make(map[*Menu]map[*Menu]struct{})
var gtk3ViewSources = make(map[*Menu]*Menu)
var gtk3ItemViews = make(map[*MenuItem]map[*linuxMenuItem]struct{})

func registerGTK3View(source, view *Menu) {
	if gtk3MenuViews[source] == nil {
		gtk3MenuViews[source] = make(map[*Menu]struct{})
	}
	gtk3MenuViews[source][view] = struct{}{}
	gtk3ViewSources[view] = source
	for i, item := range source.items {
		if item.submenu != nil && view.items[i].submenu != nil {
			registerGTK3View(item.submenu, view.items[i].submenu)
		}
	}
}

func unregisterGTK3View(view *Menu) {
	for _, item := range view.items {
		if item.submenu != nil {
			unregisterGTK3View(item.submenu)
		}
	}
	if source := gtk3ViewSources[view]; source != nil {
		delete(gtk3MenuViews[source], view)
		if len(gtk3MenuViews[source]) == 0 {
			delete(gtk3MenuViews, source)
		}
		delete(gtk3ViewSources, view)
	}
}

func windowMenuNative(menu *Menu) pointer {
	menu.Update()
	view := menu.Clone()
	registerGTK3View(menu, view)
	view.Update()
	impl := view.impl.(*linuxMenu)
	impl.onDestroy = func() { unregisterGTK3View(view) }
	return impl.native
}

func refreshGTK3MenuViews(source *Menu) {
	// A view is rebuilt directly below, not recursively through this function.
	if gtk3ViewSources[source] != nil {
		return
	}
	for view := range gtk3MenuViews[source] {
		for _, item := range view.items {
			if item.submenu != nil {
				unregisterGTK3View(item.submenu)
			}
		}
		snapshot := source.Clone()
		view.items = snapshot.items
		for i, item := range source.items {
			if item.submenu != nil {
				registerGTK3View(item.submenu, view.items[i].submenu)
			}
		}
		view.processRadioGroups()
		view.impl.(*linuxMenu).processMenu(view)
	}
}

func registerGTK3ItemView(view *linuxMenuItem) {
	original := getMenuItemByID(view.menuItem.id)
	if original == nil {
		return
	}
	view.owner = original
	if gtk3ItemViews[original] == nil {
		gtk3ItemViews[original] = make(map[*linuxMenuItem]struct{})
	}
	gtk3ItemViews[original][view] = struct{}{}
}

func unregisterGTK3ItemView(view *linuxMenuItem) {
	original := view.owner
	if original == nil {
		return
	}
	delete(gtk3ItemViews[original], view)
	if len(gtk3ItemViews[original]) == 0 {
		delete(gtk3ItemViews, original)
		delete(gtkSignalToMenuItem, original.id)
	}
}

func gtk3ItemPeers(original *MenuItem) []*linuxMenuItem {
	peers := make([]*linuxMenuItem, 0, len(gtk3ItemViews[original]))
	for peer := range gtk3ItemViews[original] {
		if peer.native != nil {
			peers = append(peers, peer)
		}
	}
	return peers
}
