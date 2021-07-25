//+build experimental

package runtime

import (
	"context"
)

func UpdateApplicationMenu(ctx context.Context) {
	frontend := getFrontend(ctx)
	frontend.UpdateApplicationMenu()
}

/*
func UpdateContextMenu(ctx context.Context, contextMenu *menu.ContextMenu) {
	frontend := getFrontend(ctx)
	bus.Publish("menu:updatecontextmenu", contextMenu)
}

func SetTrayMenu(ctx context.Context, trayMenu *menu.TrayMenu) {
	frontend := getFrontend(ctx)
	bus.Publish("menu:settraymenu", trayMenu)
}

func UpdateTrayMenuLabel(ctx context.Context, trayMenu *menu.TrayMenu) {
	frontend := getFrontend(ctx)
	bus.Publish("menu:updatetraymenulabel", trayMenu)
}

func DeleteTrayMenu(ctx context.Context, trayMenu *menu.TrayMenu) {
	frontend := getFrontend(ctx)
	bus.Publish("menu:deletetraymenu", trayMenu)
}
*/
