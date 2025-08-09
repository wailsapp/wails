package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

	app := application.New(application.Options{
		Name:        "WebviewWindow Demo (debug)",
		Description: "A demo of the WebviewWindow API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(*application.ApplicationEvent) {
		log.Println("ApplicationDidFinishLaunching")
	})

	currentWindow := func(fn func(window application.Window)) {
		current := app.Window.Current()
		if current != nil {
			fn(current)
		} else {
			println("Current WebviewWindow is nil")
		}
	}

	// Create a custom menu
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	windowCounter := 1

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("New")

	myMenu.Add("New WebviewWindow").
		SetAccelerator("CmdOrCtrl+N").
		OnClick(func(ctx *application.Context) {
			app.Window.New().
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			windowCounter++
		})
	myMenu.Add("New Frameless WebviewWindow").
		SetAccelerator("CmdOrCtrl+F").
		OnClick(func(ctx *application.Context) {
			app.Window.NewWithOptions(application.WebviewWindowOptions{
				X:         rand.Intn(1000),
				Y:         rand.Intn(800),
				Frameless: true,
				Mac: application.MacWindow{
					InvisibleTitleBarHeight: 50,
				},
			}).Show()
			windowCounter++
		})
	if runtime.GOOS == "darwin" {
		myMenu.Add("New WebviewWindow (MacTitleBarHiddenInset)").
			OnClick(func(ctx *application.Context) {
				app.Window.NewWithOptions(application.WebviewWindowOptions{
					Mac: application.MacWindow{
						TitleBar:                application.MacTitleBarHiddenInset,
						InvisibleTitleBarHeight: 25,
					},
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetHTML("<br/><br/><p>A MacTitleBarHiddenInset WebviewWindow example</p>").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (MacTitleBarHiddenInsetUnified)").
			OnClick(func(ctx *application.Context) {
				app.Window.NewWithOptions(application.WebviewWindowOptions{
					Mac: application.MacWindow{
						TitleBar:                application.MacTitleBarHiddenInsetUnified,
						InvisibleTitleBarHeight: 50,
					},
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetHTML("<br/><br/><p>A MacTitleBarHiddenInsetUnified WebviewWindow example</p>").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (MacTitleBarHidden)").
			OnClick(func(ctx *application.Context) {
				app.Window.NewWithOptions(application.WebviewWindowOptions{
					Mac: application.MacWindow{
						TitleBar:                application.MacTitleBarHidden,
						InvisibleTitleBarHeight: 25,
					},
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetHTML("<br/><br/><p>A MacTitleBarHidden WebviewWindow example</p>").
					Show()
				windowCounter++
			})
	}

	sizeMenu := menu.AddSubmenu("Size")
	sizeMenu.Add("Set Size (800,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetSize(800, 600)
		})
	})

	sizeMenu.Add("Set Size (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetSize(rand.Intn(800)+200, rand.Intn(600)+200)
		})
	})
	sizeMenu.Add("Set Min Size (200,200)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetMinSize(200, 200)
		})
	})
	sizeMenu.Add("Set Max Size (600,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetMaximiseButtonState(application.ButtonDisabled)
			w.SetMaxSize(600, 600)
		})
	})
	sizeMenu.Add("Get Current WebviewWindow Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			width, height := w.Size()
			application.InfoDialog().SetTitle("Current WebviewWindow Size").SetMessage("Width: " + strconv.Itoa(width) + " Height: " + strconv.Itoa(height)).Show()
		})
	})

	sizeMenu.Add("Reset Min Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetMinSize(0, 0)
		})
	})

	sizeMenu.Add("Reset Max Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetMaxSize(0, 0)
			w.SetMaximiseButtonState(application.ButtonEnabled)
		})
	})
	positionMenu := menu.AddSubmenu("Position")
	positionMenu.Add("Set Relative Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetRelativePosition(0, 0)
		})
	})
	positionMenu.Add("Set Relative Position (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetRelativePosition(rand.Intn(1000), rand.Intn(800))
		})
	})

	positionMenu.Add("Get Relative Position").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			x, y := w.RelativePosition()
			application.InfoDialog().SetTitle("Current WebviewWindow Relative Position").SetMessage("X: " + strconv.Itoa(x) + " Y: " + strconv.Itoa(y)).Show()
		})
	})

	positionMenu.Add("Center").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Center()
		})
	})
	stateMenu := menu.AddSubmenu("State")
	stateMenu.Add("Minimise (for 2 secs)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Minimise()
			time.Sleep(2 * time.Second)
			w.Restore()
		})
	})
	stateMenu.Add("Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Maximise()
		})
	})
	stateMenu.Add("Fullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Fullscreen()
		})
	})
	stateMenu.Add("UnFullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.UnFullscreen()
		})
	})
	stateMenu.Add("Restore").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Restore()
		})
	})
	stateMenu.Add("Hide (for 2 seconds)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.Hide()
			time.Sleep(2 * time.Second)
			w.Show()
		})
	})
	stateMenu.Add("Always on Top").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetAlwaysOnTop(true)
		})
	})
	stateMenu.Add("Not always on Top").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetAlwaysOnTop(false)
		})
	})
	stateMenu.Add("Google.com").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetURL("https://google.com")
		})
	})
	stateMenu.Add("wails.io").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			w.SetURL("https://wails.io")
		})
	})
	stateMenu.Add("Get Primary Screen").OnClick(func(ctx *application.Context) {
		screen := app.Screen.GetPrimary()
		msg := fmt.Sprintf("Screen: %+v", screen)
		application.InfoDialog().SetTitle("Primary Screen").SetMessage(msg).Show()
	})
	stateMenu.Add("Get Screens").OnClick(func(ctx *application.Context) {
		screens := app.Screen.GetAll()
		for _, screen := range screens {
			msg := fmt.Sprintf("Screen: %+v", screen)
			application.InfoDialog().SetTitle(fmt.Sprintf("Screen %s", screen.ID)).SetMessage(msg).Show()
		}
	})
	stateMenu.Add("Get Screen for WebviewWindow").OnClick(func(ctx *application.Context) {
		currentWindow(func(w application.Window) {
			screen, err := w.GetScreen()
			if err != nil {
				application.ErrorDialog().SetTitle("Error").SetMessage(err.Error()).Show()
				return
			}
			msg := fmt.Sprintf("Screen: %+v", screen)
			application.InfoDialog().SetTitle(fmt.Sprintf("Screen %s", screen.ID)).SetMessage(msg).Show()
		})
	})
	app.Window.New()

	app.Menu.Set(menu)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
