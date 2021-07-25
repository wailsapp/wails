// +build !experimental\

package menu

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

func UpdateApplicationMenu(ctx context.Context) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("menu:updateappmenu", nil)
}

func UpdateContextMenu(ctx context.Context, contextMenu *menu.ContextMenu) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("menu:updatecontextmenu", contextMenu)
}

func SetTrayMenu(ctx context.Context, trayMenu *menu.TrayMenu) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("menu:settraymenu", trayMenu)
}

func UpdateTrayMenuLabel(ctx context.Context, trayMenu *menu.TrayMenu) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("menu:updatetraymenulabel", trayMenu)
}

func DeleteTrayMenu(ctx context.Context, trayMenu *menu.TrayMenu) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("menu:deletetraymenu", trayMenu)
}
