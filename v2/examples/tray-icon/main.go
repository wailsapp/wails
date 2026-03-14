package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func resolveTrayIconPath() string {
	return resolveTrayAssetPath("tray-icon.png")
}

func resolveTrayAssetPath(filename string) string {
	candidates := []string{
		filepath.Join("examples", "tray-icon", "trayicons", filename),
		filepath.Join("trayicons", filename),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if abs, err := filepath.Abs(p); err == nil {
				return abs
			}
			return p
		}
	}

	return candidates[0]
}

func main() {
	app := NewApp()
	primaryIcon := resolveTrayIconPath()
	altIcon := resolveTrayAssetPath("tray-icon-alt.png")

	state := struct {
		mu      sync.Mutex
		mode    int
		clicks  int
		updates int
	}{}

	var buildTray func() *menu.TrayMenu
	buildTray = func() *menu.TrayMenu {
		state.mu.Lock()
		mode := state.mode
		clicks := state.clicks
		updates := state.updates
		state.mu.Unlock()

		icon := primaryIcon
		label := "TrayIcon-A"
		actionText := "Primary Action: switch to B"
		actionLog := "action A"
		if mode == 1 {
			icon = altIcon
			label = "TrayIcon-B"
			actionText = "Primary Action: switch to A"
			actionLog = "action B"
		}

		trayMenu := menu.NewMenu()

		trayMenu.AddText("Show", nil, func(_ *menu.CallbackData) {
			runtime.Show(app.ctx)
		})
		trayMenu.AddText(actionText, nil, func(_ *menu.CallbackData) {
			state.mu.Lock()
			state.mode = 1 - state.mode
			state.clicks++
			state.updates++
			currentClicks := state.clicks
			state.mu.Unlock()

			fmt.Printf("[tray-test] %s clicked (clicks=%d)\n", actionLog, currentClicks)
			runtime.TraySetSystemTray(app.ctx, buildTray())
		})

		trayMenu.AddText("Ping", nil, func(_ *menu.CallbackData) {
			fmt.Println("[tray-test] Ping clicked")
		})

		trayMenu.AddText("Toggle icon + text now", nil, func(_ *menu.CallbackData) {
			state.mu.Lock()
			state.mode = 1 - state.mode
			state.clicks++
			state.updates++
			currentClicks := state.clicks
			state.mu.Unlock()

			fmt.Printf("[tray-test] manual toggle clicked (clicks=%d)\n", currentClicks)
			runtime.TraySetSystemTray(app.ctx, buildTray())
		})

		trayMenu.AddSeparator()
		trayMenu.AddText("Quit", nil, func(_ *menu.CallbackData) {
			runtime.Quit(app.ctx)
		})

		return &menu.TrayMenu{
			Label:   fmt.Sprintf("%s (%d)", label, updates),
			Tooltip: fmt.Sprintf("Mode=%s, Clicks=%d", label, clicks),
			Image:   icon,
			Menu:    trayMenu,
		}
	}

	onStartup := func(ctx context.Context) {
		app.startup(ctx)

		go func() {
			time.Sleep(4 * time.Second)
			state.mu.Lock()
			state.mode = 1 - state.mode
			state.updates++
			state.mu.Unlock()
			fmt.Println("[tray-test] timed update: 4s")
			runtime.TraySetSystemTray(app.ctx, buildTray())

			time.Sleep(4 * time.Second)
			state.mu.Lock()
			state.mode = 1 - state.mode
			state.updates++
			state.mu.Unlock()
			fmt.Println("[tray-test] timed update: 8s")
			runtime.TraySetSystemTray(app.ctx, buildTray())
		}()
	}

	err := wails.Run(&options.App{
		Title:  "Wails Tray Icon Test",
		Width:  720,
		Height: 420,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:         onStartup,
		Tray:              buildTray(),
		HideWindowOnClose: true,
	})

	if err != nil {
		fmt.Println("[tray-test] Error:", err)
	}
}
