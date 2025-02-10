package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"runtime"
)

//go:embed assets
var assets embed.FS

// Define our custom themes
var (
	// Midnight Ocean - A deep, calming blue theme
	midnightOcean = application.ThemeSettings{
		DarkModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(0, 96, 122),    // Bright blue accent
			TitleBarColour:  application.NewRGBPtr(0, 32, 41),     // Rich deep blue
			TitleTextColour: application.NewRGBPtr(220, 230, 242), // Soft white-blue
		},
		DarkModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(0, 48, 61),     // Muted blue accent
			TitleBarColour:  application.NewRGBPtr(0, 24, 31),     // Darker blue
			TitleTextColour: application.NewRGBPtr(180, 190, 202), // Muted white-blue
		},
		LightModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(0, 96, 122),    // Bright blue accent
			TitleBarColour:  application.NewRGBPtr(240, 245, 250), // Very light blue
			TitleTextColour: application.NewRGBPtr(0, 32, 41),     // Rich deep blue
		},
		LightModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(0, 48, 61),     // Muted blue accent
			TitleBarColour:  application.NewRGBPtr(230, 235, 240), // Slightly darker light blue
			TitleTextColour: application.NewRGBPtr(0, 24, 31),     // Darker blue
		},
		DarkModeMenuBar: &application.MenuBarTheme{
			Default: &application.TextTheme{
				Text:       application.NewRGBPtr(220, 230, 242), // Soft white-blue
				Background: application.NewRGBPtr(0, 32, 41),     // Rich deep blue
			},
			Hover: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(0, 96, 122),    // Bright blue accent
			},
			Selected: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(0, 72, 92),     // Medium bright blue
			},
		},
	}

	// Forest Haven - An elegant green theme
	forestHaven = application.ThemeSettings{
		DarkModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(76, 145, 65),   // Vibrant forest green
			TitleBarColour:  application.NewRGBPtr(22, 42, 33),    // Deep forest green
			TitleTextColour: application.NewRGBPtr(230, 240, 230), // Soft sage
		},
		DarkModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(56, 105, 45),   // Muted forest green
			TitleBarColour:  application.NewRGBPtr(18, 32, 25),    // Darker forest green
			TitleTextColour: application.NewRGBPtr(190, 200, 190), // Muted sage
		},
		LightModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(76, 145, 65),   // Vibrant forest green
			TitleBarColour:  application.NewRGBPtr(240, 245, 240), // Very light sage
			TitleTextColour: application.NewRGBPtr(22, 42, 33),    // Deep forest green
		},
		LightModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(56, 105, 45),   // Muted forest green
			TitleBarColour:  application.NewRGBPtr(230, 235, 230), // Slightly darker light sage
			TitleTextColour: application.NewRGBPtr(18, 32, 25),    // Darker forest green
		},
		DarkModeMenuBar: &application.MenuBarTheme{
			Default: &application.TextTheme{
				Text:       application.NewRGBPtr(230, 240, 230), // Soft sage
				Background: application.NewRGBPtr(22, 42, 33),    // Deep forest green
			},
			Hover: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(76, 145, 65),   // Vibrant forest green
			},
			Selected: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(66, 125, 55),   // Medium forest green
			},
		},
	}

	// Royal Purple - A sophisticated purple theme
	royalPurple = application.ThemeSettings{
		DarkModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(116, 80, 170),  // Rich purple
			TitleBarColour:  application.NewRGBPtr(50, 33, 58),    // Deep purple
			TitleTextColour: application.NewRGBPtr(235, 230, 240), // Soft lavender
		},
		DarkModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(86, 60, 130),   // Muted purple
			TitleBarColour:  application.NewRGBPtr(40, 26, 46),    // Darker purple
			TitleTextColour: application.NewRGBPtr(195, 190, 200), // Muted lavender
		},
		LightModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(116, 80, 170),  // Rich purple
			TitleBarColour:  application.NewRGBPtr(245, 240, 250), // Very light lavender
			TitleTextColour: application.NewRGBPtr(50, 33, 58),    // Deep purple
		},
		LightModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(86, 60, 130),   // Muted purple
			TitleBarColour:  application.NewRGBPtr(235, 230, 240), // Slightly darker light lavender
			TitleTextColour: application.NewRGBPtr(40, 26, 46),    // Darker purple
		},
		DarkModeMenuBar: &application.MenuBarTheme{
			Default: &application.TextTheme{
				Text:       application.NewRGBPtr(235, 230, 240), // Soft lavender
				Background: application.NewRGBPtr(50, 33, 58),    // Deep purple
			},
			Hover: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(116, 80, 170),  // Rich purple
			},
			Selected: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(96, 65, 140),   // Medium purple
			},
		},
	}

	// Desert Sand - A warm, neutral theme
	desertSand = application.ThemeSettings{
		DarkModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(190, 149, 109), // Warm brown
			TitleBarColour:  application.NewRGBPtr(70, 63, 58),    // Dark taupe
			TitleTextColour: application.NewRGBPtr(242, 235, 228), // Warm sand
		},
		DarkModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(150, 119, 89),  // Muted brown
			TitleBarColour:  application.NewRGBPtr(60, 53, 48),    // Darker taupe
			TitleTextColour: application.NewRGBPtr(202, 195, 188), // Muted sand
		},
		LightModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(190, 149, 109), // Warm brown
			TitleBarColour:  application.NewRGBPtr(242, 235, 228), // Warm sand
			TitleTextColour: application.NewRGBPtr(70, 63, 58),    // Dark taupe
		},
		LightModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(150, 119, 89),  // Muted brown
			TitleBarColour:  application.NewRGBPtr(232, 225, 218), // Slightly darker sand
			TitleTextColour: application.NewRGBPtr(60, 53, 48),    // Darker taupe
		},
		DarkModeMenuBar: &application.MenuBarTheme{
			Default: &application.TextTheme{
				Text:       application.NewRGBPtr(242, 235, 228), // Warm sand
				Background: application.NewRGBPtr(70, 63, 58),    // Dark taupe
			},
			Hover: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(190, 149, 109), // Warm brown
			},
			Selected: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(170, 129, 89),  // Medium brown
			},
		},
	}

	// Arctic Frost - A clean, crisp light theme
	arcticFrost = application.ThemeSettings{
		DarkModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(127, 179, 213), // Glacier blue
			TitleBarColour:  application.NewRGBPtr(44, 62, 80),    // Deep slate
			TitleTextColour: application.NewRGBPtr(240, 245, 250), // Ice blue-white
		},
		DarkModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(97, 149, 183),  // Muted glacier blue
			TitleBarColour:  application.NewRGBPtr(34, 52, 70),    // Darker slate
			TitleTextColour: application.NewRGBPtr(200, 205, 210), // Muted ice blue
		},
		LightModeActive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(127, 179, 213), // Glacier blue
			TitleBarColour:  application.NewRGBPtr(240, 245, 250), // Ice blue-white
			TitleTextColour: application.NewRGBPtr(44, 62, 80),    // Deep slate
		},
		LightModeInactive: &application.WindowTheme{
			BorderColour:    application.NewRGBPtr(97, 149, 183),  // Muted glacier blue
			TitleBarColour:  application.NewRGBPtr(230, 235, 240), // Slightly darker ice blue
			TitleTextColour: application.NewRGBPtr(34, 52, 70),    // Darker slate
		},
		DarkModeMenuBar: &application.MenuBarTheme{
			Default: &application.TextTheme{
				Text:       application.NewRGBPtr(240, 245, 250), // Ice blue-white
				Background: application.NewRGBPtr(44, 62, 80),    // Deep slate
			},
			Hover: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(127, 179, 213), // Glacier blue
			},
			Selected: &application.TextTheme{
				Text:       application.NewRGBPtr(255, 255, 255), // Pure white
				Background: application.NewRGBPtr(107, 159, 193), // Medium glacier blue
			},
		},
	}
)

func main() {
	app := application.New(application.Options{
		Name:        "Theme Custom",
		Description: "An example of custom window themes",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create our main menu
	mainMenu := app.NewMenu()

	// Add standard menus for macOS
	if runtime.GOOS == "darwin" {
		mainMenu.AddRole(application.AppMenu)
	}
	mainMenu.AddRole(application.FileMenu)
	mainMenu.AddRole(application.EditMenu)
	mainMenu.AddRole(application.WindowMenu)

	// Create our themes menu
	themesMenu := mainMenu.AddSubmenu("Themes")
	themesMenu.Add("Midnight Ocean").OnClick(func(ctx *application.Context) {
		createWindow(app, "Midnight Ocean Theme", midnightOcean, mainMenu)
	})
	themesMenu.Add("Forest Haven").OnClick(func(ctx *application.Context) {
		createWindow(app, "Forest Haven Theme", forestHaven, mainMenu)
	})
	themesMenu.Add("Royal Purple").OnClick(func(ctx *application.Context) {
		createWindow(app, "Royal Purple Theme", royalPurple, mainMenu)
	})
	themesMenu.Add("Desert Sand").OnClick(func(ctx *application.Context) {
		createWindow(app, "Desert Sand Theme", desertSand, mainMenu)
	})
	themesMenu.Add("Arctic Frost").OnClick(func(ctx *application.Context) {
		createWindow(app, "Arctic Frost Theme", arcticFrost, mainMenu)
	})

	app.SetMenu(mainMenu)

	// Create our first window with the default theme
	createWindow(app, "Welcome to Theme Custom", midnightOcean, mainMenu)

	app.Run()
}

func createWindow(app *application.App, title string, theme application.ThemeSettings, mainMenu *application.Menu) *application.WebviewWindow {
	// Create a new window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:   title,
		Width:   800,
		Height:  600,
		Windows: application.WindowsWindow{CustomTheme: theme, Menu: mainMenu},
	})

	// Set up the window content
	window.SetURL("index.html")

	window.Show()
	return window
}
