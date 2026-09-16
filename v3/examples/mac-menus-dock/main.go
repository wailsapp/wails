package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/dock"
)

//go:embed assets
var assets embed.FS

func assetsFS() fs.FS {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatal(err)
	}
	return sub
}

func main() {
	dockService := dock.New()

	app := application.New(application.Options{
		Name:        "macOS Menus and Dock",
		Description: "SF Symbols, badges, palettes, dock menu, Open Recent and dock progress",
		Services: []application.Service{
			application.NewService(dockService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assetsFS()),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	logLine := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		log.Println(line)
		app.Event.Emit("demo:log", line)
	}

	// Files opened through Open Recent (or Finder) arrive here.
	app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
		logLine("opened file: %s", event.Context().Filename())
	})

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	// File menu with the standard Open Recent submenu.
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("Add a recent document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
		path := filepath.Join(os.TempDir(), fmt.Sprintf("wails-recent-%d.txt", time.Now().Unix()))
		if err := os.WriteFile(path, []byte("hello from wails\n"), 0o644); err != nil {
			logLine("could not create %s: %v", path, err)
			return
		}
		app.Menu.AddRecentDocument(path)
		logLine("added recent document %s (%d in list)", path, len(app.Menu.RecentDocuments()))
	})
	fileMenu.AddRole(application.OpenRecent)
	fileMenu.AddSeparator()
	if runtime.GOOS == "darwin" {
		fileMenu.AddRole(application.CloseWindow)
	} else {
		fileMenu.AddRole(application.Quit)
	}
	menu.AddRole(application.EditMenu)

	// Demo menu: symbols, badges, section headers, mixed state, alternates.
	demo := menu.AddSubmenu("Demo")
	demo.AddSectionHeader("Symbols and badges")
	inbox := demo.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
	inbox.OnClick(func(ctx *application.Context) {
		next := inbox.BadgeCount() + 1
		inbox.SetBadge(next)
		logLine("inbox badge is now %d", next)
	})
	demo.Add("Updates").SetSymbol("arrow.down.circle").SetBadgeText("New").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().ClearBadge()
		logLine("cleared the Updates badge")
	})

	demo.AddSeparator()
	demo.AddSectionHeader("States")
	wrap := demo.AddCheckbox("Wrap lines (mixed)", false).SetMixed()
	wrap.OnClick(func(ctx *application.Context) {
		logLine("wrap lines: checked=%v mixed=%v", wrap.Checked(), wrap.Mixed())
	})
	demo.Add("Reset to mixed").SetIndentationLevel(1).OnClick(func(*application.Context) {
		wrap.SetMixed()
		logLine("wrap lines back to mixed")
	})
	demo.Add("Indented level 2").SetIndentationLevel(2).SetEnabled(false)

	demo.AddSeparator()
	demo.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").SetSymbol("xmark").OnClick(func(*application.Context) {
		logLine("close tab")
	})
	// Shown in place of "Close Tab" while Option is held.
	demo.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).OnClick(func(*application.Context) {
		logLine("close all tabs")
	})

	demo.AddSeparator()
	// Palette (macOS 14+). Give it a label to present it as a titled submenu.
	colours := []application.RGBA{
		application.NewRGB(255, 59, 48),
		application.NewRGB(255, 149, 0),
		application.NewRGB(52, 199, 89),
		application.NewRGB(0, 122, 255),
		application.NewRGB(175, 82, 222),
	}
	demo.AddPalette([]string{"tag.fill"}, colours, 3, func(ctx *application.Context, index int) {
		logLine("tag colour %d selected", index)
	}).SetLabel("Tag colour")

	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)

	// Dock menu, rebuilt on every right-click so it reflects the current
	// badge and progress.
	progress := 0.0
	app.Menu.OnDockMenu(func() *application.Menu {
		dockMenu := application.NewMenu()
		dockMenu.AddSectionHeader("Dock demo")
		dockMenu.Add("Increment badge").SetSymbol("plus.circle").OnClick(func(*application.Context) {
			current := 0
			if badge := dockService.GetBadge(); badge != nil {
				fmt.Sscanf(*badge, "%d", &current)
			}
			_ = dockService.SetBadge(fmt.Sprintf("%d", current+1))
			logLine("dock badge is now %d", current+1)
		})
		dockMenu.Add("Clear badge").OnClick(func(*application.Context) {
			_ = dockService.RemoveBadge()
			logLine("dock badge cleared")
		})
		dockMenu.AddSeparator()
		dockMenu.Add(fmt.Sprintf("Advance progress (%.0f%%)", progress*100)).SetSymbol("gauge").OnClick(func(*application.Context) {
			progress += 0.25
			if progress > 1 {
				progress = 0.25
			}
			_ = dockService.SetProgress(progress)
			logLine("dock progress %.0f%%", progress*100)
		})
		dockMenu.Add("Clear progress").OnClick(func(*application.Context) {
			progress = 0
			_ = dockService.ClearProgress()
			logLine("dock progress cleared")
		})
		return dockMenu
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "macOS menus and Dock",
		Width:  640,
		Height: 480,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
