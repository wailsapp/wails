package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"log"
)

func main() {

	// Create application with options
	app, err := wails.CreateAppWithOptions(&options.App{
		Title:     "Kitchen Sink",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		//Tray:      menu.NewMenuFromItems(menu.AppMenu()),
		//Menu:      menu.NewMenuFromItems(menu.AppMenu()),
		ContextMenus: map[string]*menu.Menu{
			"test": menu.NewMenuFromItems(menu.Text("Test Menu", "Test Context Menu")),
		},
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			// Comment out line below to see Window.SetTitle() work
			TitleBar: mac.TitleBarHiddenInset(),
			Menu:     createApplicationMenu(),
			Tray:     createApplicationTray(),
		},
		LogLevel: logger.TRACE,
	})

	if err != nil {
		log.Fatal(err)
	}

	app.Bind(&Events{})
	app.Bind(&Logger{})
	app.Bind(&Browser{})
	app.Bind(&System{})
	app.Bind(&Dialog{})
	app.Bind(&Window{})
	app.Bind(&Menu{})
	app.Bind(&Tray{})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
