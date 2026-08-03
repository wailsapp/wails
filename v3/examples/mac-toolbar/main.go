package main

import (
	"embed"
	"log"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

// infoClicks backs the Info toolbar item's BadgeCount. There is no
// per-item update method: bumping a badge means rebuilding the whole
// toolbar and calling SetToolbar again, which is exactly what this demo
// is here to prove.
var infoClicks atomic.Int32

func main() {
	app := application.New(application.Options{
		Name:        "mac-toolbar",
		Description: "A demo of MacToolbar: buttons, a segmented group, a search field and a badged button",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Notes",
		Width:  900,
		Height: 620,
		URL:    "/",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	// SetToolbar is called before app.Run() starts the native event loop
	// (window.impl doesn't exist yet at this point), which exercises the
	// stash-and-apply-during-run() path rather than the immediate one.
	window.SetToolbar(buildToolbar(window, 0))

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func buildToolbar(window *application.WebviewWindow, badgeCount int) *application.MacToolbar {
	emit := func(name string, data any) {
		application.Get().Event.Emit(name, data)
	}

	return &application.MacToolbar{
		Items: []application.MacToolbarItem{
			{
				ID:    "view-mode",
				Kind:  application.ToolbarGroup,
				Label: "View",
				Items: []application.MacToolbarItem{
					{
						ID: "edit", Label: "Edit", SymbolName: "pencil", Bordered: true,
						OnClick: func(*application.Context) { emit("view:mode", "edit") },
					},
					{
						ID: "preview", Label: "Preview", SymbolName: "eye", Bordered: true,
						OnClick: func(*application.Context) { emit("view:mode", "preview") },
					},
				},
			},
			{ID: "flex", Kind: application.ToolbarFlexibleSpace},
			{
				ID: "search", Kind: application.ToolbarSearchField, Label: "Search",
				OnSearch: func(_ *application.Context, query string) {
					emit("toolbar:search", query)
				},
			},
			{
				ID: "share", Kind: application.ToolbarButton, Label: "Share",
				SymbolName: "square.and.arrow.up", Bordered: true,
				OnClick: func(*application.Context) {
					emit("toolbar:action", map[string]any{"action": "share"})
				},
			},
			{
				ID: "info", Kind: application.ToolbarButton, Label: "Info",
				SymbolName: "info.circle", Bordered: true, Prominent: true,
				TintColor:  &application.RGBA{Red: 10, Green: 132, Blue: 255, Alpha: 255},
				BadgeCount: badgeCount,
				OnClick: func(*application.Context) {
					count := int(infoClicks.Add(1))
					window.SetToolbar(buildToolbar(window, count))
					emit("toolbar:action", map[string]any{"action": "info", "count": count})
				},
			},
		},
	}
}
