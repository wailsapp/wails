package application

// chromeEventLoops holds the goroutines that drain native chrome callback
// channels (toolbar clicks, sidebar selection, inspector edits, ...) once
// App.Run starts. Each chrome feature registers its own loop from an init
// function in its own file, so adding a feature never touches App.Run.
var chromeEventLoops []func(*App)

func registerChromeEventLoop(loop func(*App)) {
	if loop != nil {
		chromeEventLoops = append(chromeEventLoops, loop)
	}
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			itemID := <-toolbarItemClicked
			go handleToolbarItemClicked(itemID)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-toolbarSearchTriggered
			go handleToolbarSearch(event.itemID, event.query)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-toolbarShareCompleted
			go handleToolbarShareResult(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-splitPaneCollapseEvents
			handleMacSplitPaneCollapsed(event.paneID, event.collapsed)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			itemID := <-macSidebarItemSelected
			go handleMacSidebarItemSelected(itemID)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macInspectorControlEvents
			go handleMacInspectorControlEvent(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			editorID := <-macTextEditorChanged
			go handleMacTextEditorChanged(editorID)
		}
	})
	registerChromeEventLoop(func(a *App) {
		for {
			windowID := <-nativeWindowClosed
			if window, ok := a.NativeWindow.GetByID(windowID); ok {
				go window.Close()
			}
		}
	})
}
