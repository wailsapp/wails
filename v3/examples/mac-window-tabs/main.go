package main

import (
	"embed"
	_ "embed"
	"log"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "mac-window-tabs",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// The Tabs menu drives the tab group API for the key window. AppKit adds
	// its own tab items to the Window menu (Show Next Tab, Merge All Windows
	// and so on); this menu performs the same operations from Go.
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.FileMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)

	tabs := menu.AddSubmenu("Tabs")
	tabs.Add("Add Tab").SetAccelerator("CmdOrCtrl+T").OnClick(func(*application.Context) {
		if err := addTabToCurrentWindow(); err != nil {
			app.Logger.Error("Add Tab failed", "error", err)
		}
	})
	tabs.Add("Add Native Window Tab").OnClick(func(*application.Context) {
		if err := addNativeTabToCurrentWindow(); err != nil {
			app.Logger.Error("Add Native Window Tab failed", "error", err)
		}
	})
	tabs.AddSeparator()
	tabs.Add("Select Next Tab").OnClick(func(*application.Context) {
		currentTabGroup().SelectNext()
	})
	tabs.Add("Select Previous Tab").OnClick(func(*application.Context) {
		currentTabGroup().SelectPrevious()
	})
	tabs.Add("Select First Tab").OnClick(func(*application.Context) {
		group := currentTabGroup()
		if windows := group.Windows(); len(windows) > 0 {
			group.Select(windows[0])
		}
	})
	tabs.AddSeparator()
	tabs.Add("Toggle Tab Bar").OnClick(func(*application.Context) {
		currentTabGroup().ToggleTabBar()
	})
	tabs.Add("Toggle Tab Overview").OnClick(func(*application.Context) {
		currentTabGroup().ToggleTabOverview()
	})
	tabs.AddSeparator()
	tabs.Add("Rename Tab").OnClick(func(*application.Context) {
		if window := currentWebviewWindow(); window != nil {
			stamp := time.Now().Format("15:04:05")
			window.SetTabTitle("Renamed at " + stamp)
			window.SetTabTooltip("Tab title set from Go at " + stamp)
		}
	})
	tabs.Add("Move Tab to New Window").OnClick(func(*application.Context) {
		if window := currentWebviewWindow(); window != nil {
			window.MoveTabToNewWindow()
		}
	})
	tabs.Add("Merge All Windows").OnClick(func(*application.Context) {
		if window := currentWebviewWindow(); window != nil {
			window.MergeAllWindows()
		}
	})
	tabs.AddSeparator()
	tabs.Add("Log Tab Groups").OnClick(func(*application.Context) {
		groups := app.Window.TabGroups()
		app.Logger.Info("tab groups", "count", len(groups))
		for index, group := range groups {
			names := make([]string, 0, group.Count())
			for _, window := range group.Windows() {
				names = append(names, window.Name())
			}
			for _, window := range group.NativeWindows() {
				names = append(names, window.Name()+" (native)")
			}
			selected := "none"
			if window := group.SelectedWindow(); window != nil {
				selected = window.Name()
			}
			app.Logger.Info("tab group",
				"index", index,
				"identifier", group.Identifier(),
				"count", group.Count(),
				"windows", strings.Join(names, ", "),
				"selected", selected,
				"tabBarVisible", group.IsTabBarVisible(),
				"overviewVisible", group.IsOverviewVisible(),
			)
		}
	})
	app.Menu.Set(menu)

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	windowBackground := application.NewRGB(27, 38, 54)

	// Open a single window at startup. It uses TabbingModePreferred so that any
	// "tabbed" window opened from its buttons will merge into it as a new tab.
	// The buttons in the frontend drive the demo:
	//   - "Open tabbed window"     -> TabbingModePreferred     (joins this window)
	//   - "Open non-tabbed window" -> TabbingModeDisallowed    (opens standalone)
	//   - "Add tab via AddTab"     -> TabbingModeAutomatic     (joined explicitly)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "macOS Window Tabs",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			TabbingMode:             application.MacWindowTabbingModePreferred,
		},
		BackgroundColour: windowBackground,
		URL:              "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
