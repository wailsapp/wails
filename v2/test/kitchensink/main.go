package main

import (
	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

func main() {

	// Create menu
	myMenu := menu.DefaultMacMenu()

	windowMenu := menu.SubMenu("Test", []*menu.MenuItem{
		menu.Togglefullscreen(),
		menu.Minimize(),
		menu.Zoom(),

		menu.Separator(),

		menu.Copy(),
		menu.Cut(),
		menu.Delete(),

		menu.Separator(),

		menu.Front(),

		menu.SubMenu("Test Submenu", []*menu.MenuItem{
			menu.Text("Hi!", "hello"), // Label = "Hi!", ID= "hello"
			&menu.MenuItem{
				Label:    "Disabled Menu",
				Type:     menu.TextType,
				Disabled: true,
			},
			&menu.MenuItem{
				Label:  "Hidden Menu",
				Type:   menu.TextType,
				Hidden: true,
			},
			&menu.MenuItem{
				ID:      "checkbox-menu",
				Label:   "Checkbox Menu",
				Type:    menu.CheckboxType,
				Checked: true,
			},
			menu.Separator(),
			menu.Radio("😀 Option 1", "😀option-1", true),
			menu.Radio("😺 Option 2", "option-2", false),
			menu.Radio("❤️ Option 3", "option-3", false),
		}),
	})

	myMenu.Append(windowMenu)

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:     "Kitchen Sink",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			// Comment out line below to see Window.SetTitle() work
			TitleBar: mac.TitleBarHiddenInset(),
			Menu:     myMenu,
		},
		LogLevel: logger.TRACE,
	})

	app.Bind(&Events{})
	app.Bind(&Logger{})
	app.Bind(&Browser{})
	app.Bind(&System{})
	app.Bind(&Dialog{})
	app.Bind(&Window{})

	app.Run()
}
