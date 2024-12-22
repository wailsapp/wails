package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "WebviewWindow Demo",
		Description: "A demo of the WebviewWindow API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	app.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
		log.Println("ApplicationDidFinishLaunching")
	})

	var hiddenWindows []*application.WebviewWindow

	currentWindow := func(fn func(window *application.WebviewWindow)) {
		if app.CurrentWindow() != nil {
			fn(app.CurrentWindow())
		} else {
			println("Current WebviewWindow is nil")
		}
	}

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	windowCounter := 1

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("New")

	myMenu.Add("New WebviewWindow").
		SetAccelerator("CmdOrCtrl+N").
		OnClick(func(ctx *application.Context) {
			app.NewWebviewWindow().
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			windowCounter++
		})
	myMenu.Add("New WebviewWindow (Hides on Close one time)").
		SetAccelerator("CmdOrCtrl+H").
		OnClick(func(ctx *application.Context) {
			w := app.NewWebviewWindow()
			w.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
				if !lo.Contains(hiddenWindows, w) {
					hiddenWindows = append(hiddenWindows, w)
					go func() {
						time.Sleep(5 * time.Second)
						w.Show()
					}()
					w.Hide()
					e.Cancel()
				}
				// Remove the window from the hiddenWindows list
				hiddenWindows = lo.Without(hiddenWindows, w)
			})
			w.SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			windowCounter++
		})
	myMenu.Add("New Frameless WebviewWindow").
		SetAccelerator("CmdOrCtrl+F").
		OnClick(func(ctx *application.Context) {
			app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				X:                rand.Intn(1000),
				Y:                rand.Intn(800),
				BackgroundColour: application.NewRGB(33, 37, 41),
				Frameless:        true,
				Mac: application.MacWindow{
					InvisibleTitleBarHeight: 50,
				},
			}).Show()
			windowCounter++
		})
	myMenu.Add("New WebviewWindow (ignores mouse events").
		SetAccelerator("CmdOrCtrl+F").
		OnClick(func(ctx *application.Context) {
			app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				HTML:              "<div style='width: 100%; height: 95%; border: 3px solid red; background-color: \"0000\";'></div>",
				X:                 rand.Intn(1000),
				Y:                 rand.Intn(800),
				IgnoreMouseEvents: true,
				BackgroundType:    application.BackgroundTypeTransparent,
				Mac: application.MacWindow{
					InvisibleTitleBarHeight: 50,
				},
			}).Show()
			windowCounter++
		})
	if runtime.GOOS == "darwin" {
		myMenu.Add("New WebviewWindow (MacTitleBarHiddenInset)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					Mac: application.MacWindow{
						TitleBar:                application.MacTitleBarHiddenInset,
						InvisibleTitleBarHeight: 25,
					},
				}).
					SetBackgroundColour(application.NewRGB(33, 37, 41)).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetHTML("<br/><br/><p>A MacTitleBarHiddenInset WebviewWindow example</p>").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (MacTitleBarHiddenInsetUnified)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
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
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
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
	if runtime.GOOS == "windows" {
		myMenu.Add("New WebviewWindow (Mica)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					Title:          "WebviewWindow " + strconv.Itoa(windowCounter),
					X:              rand.Intn(1000),
					Y:              rand.Intn(800),
					BackgroundType: application.BackgroundTypeTranslucent,
					HTML:           "<html style='background-color: rgba(0,0,0,0);'><body></body></html>",
					Windows: application.WindowsWindow{
						BackdropType: application.Mica,
					},
				}).Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Acrylic)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					Title:          "WebviewWindow " + strconv.Itoa(windowCounter),
					X:              rand.Intn(1000),
					Y:              rand.Intn(800),
					BackgroundType: application.BackgroundTypeTranslucent,
					HTML:           "<html style='background-color: rgba(0,0,0,0);'><body></body></html>",
					Windows: application.WindowsWindow{
						BackdropType: application.Acrylic,
					},
				}).Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Tabbed)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					Title:          "WebviewWindow " + strconv.Itoa(windowCounter),
					X:              rand.Intn(1000),
					Y:              rand.Intn(800),
					BackgroundType: application.BackgroundTypeTranslucent,
					HTML:           "<html style='background-color: rgba(0,0,0,0);'><body></body></html>",
					Windows: application.WindowsWindow{
						BackdropType: application.Tabbed,
					},
				}).Show()
				windowCounter++
			})
	}

	sizeMenu := menu.AddSubmenu("Size")
	sizeMenu.Add("Set Size (800,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetSize(800, 600)
		})
	})

	sizeMenu.Add("Set Size (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetSize(rand.Intn(800)+200, rand.Intn(600)+200)
		})
	})
	sizeMenu.Add("Set Min Size (200,200)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMinSize(200, 200)
		})
	})
	sizeMenu.Add("Set Max Size (600,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMaximiseButtonState(application.ButtonDisabled)
			w.SetMaxSize(600, 600)
		})
	})
	sizeMenu.Add("Get Current WebviewWindow Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			width, height := w.Size()
			application.InfoDialog().SetTitle("Current WebviewWindow Size").SetMessage("Width: " + strconv.Itoa(width) + " Height: " + strconv.Itoa(height)).Show()
		})
	})

	sizeMenu.Add("Reset Min Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMinSize(0, 0)
		})
	})

	sizeMenu.Add("Reset Max Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMaxSize(0, 0)
			w.SetMaximiseButtonState(application.ButtonEnabled)
		})
	})
	positionMenu := menu.AddSubmenu("Position")
	positionMenu.Add("Set Relative Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetRelativePosition(0, 0)
		})
	})
	positionMenu.Add("Set Relative Position (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetRelativePosition(rand.Intn(1000), rand.Intn(800))
		})
	})

	positionMenu.Add("Get Relative Position").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			x, y := w.RelativePosition()
			application.InfoDialog().SetTitle("Current WebviewWindow Position").SetMessage("X: " + strconv.Itoa(x) + " Y: " + strconv.Itoa(y)).Show()
		})
	})

	positionMenu.Add("Set Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetPosition(0, 0)
		})
	})

	positionMenu.Add("Set Position (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetPosition(rand.Intn(1000), rand.Intn(800))
		})
	})

	positionMenu.Add("Get Position").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			x, y := w.Position()
			application.InfoDialog().SetTitle("Current WebviewWindow Position").SetMessage("X: " + strconv.Itoa(x) + " Y: " + strconv.Itoa(y)).Show()
		})
	})

	positionMenu.Add("Center").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Center()
		})
	})
	stateMenu := menu.AddSubmenu("State")
	stateMenu.Add("Minimise (for 2 secs)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Minimise()
			time.Sleep(2 * time.Second)
			w.Restore()
		})
	})
	stateMenu.Add("Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Maximise()
		})
	})
	stateMenu.Add("Fullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Fullscreen()
		})
	})
	stateMenu.Add("UnFullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.UnFullscreen()
		})
	})
	stateMenu.Add("Restore").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Restore()
		})
	})
	stateMenu.Add("Hide (for 2 seconds)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Hide()
			time.Sleep(2 * time.Second)
			w.Show()
		})
	})
	stateMenu.Add("Always on Top").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetAlwaysOnTop(true)
		})
	})
	stateMenu.Add("Not always on Top").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetAlwaysOnTop(false)
		})
	})
	stateMenu.Add("Google.com").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetURL("https://google.com")
		})
	})
	stateMenu.Add("wails.io").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetURL("https://wails.io")
		})
	})
	stateMenu.Add("Get Primary Screen").OnClick(func(ctx *application.Context) {
		screen, err := app.GetPrimaryScreen()
		if err != nil {
			application.ErrorDialog().SetTitle("Error").SetMessage(err.Error()).Show()
			return
		}
		msg := fmt.Sprintf("Screen: %+v", screen)
		application.InfoDialog().SetTitle("Primary Screen").SetMessage(msg).Show()
	})
	stateMenu.Add("Get Screens").OnClick(func(ctx *application.Context) {
		screens, err := app.GetScreens()
		if err != nil {
			application.ErrorDialog().SetTitle("Error").SetMessage(err.Error()).Show()
			return
		}
		for _, screen := range screens {
			msg := fmt.Sprintf("Screen: %+v", screen)
			application.InfoDialog().SetTitle(fmt.Sprintf("Screen %s", screen.ID)).SetMessage(msg).Show()
		}
	})
	stateMenu.Add("Get Screen for WebviewWindow").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			screen, err := w.GetScreen()
			if err != nil {
				application.ErrorDialog().SetTitle("Error").SetMessage(err.Error()).Show()
				return
			}
			msg := fmt.Sprintf("Screen: %+v", screen)
			application.InfoDialog().SetTitle(fmt.Sprintf("Screen %s", screen.ID)).SetMessage(msg).Show()
		})
	})
	stateMenu.Add("Disable for 5s").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetEnabled(false)
			time.Sleep(5 * time.Second)
			w.SetEnabled(true)
		})
	})

	if runtime.GOOS == "windows" {
		stateMenu.Add("Flash Start").OnClick(func(ctx *application.Context) {
			currentWindow(func(w *application.WebviewWindow) {
				time.Sleep(2 * time.Second)
				w.Flash(true)
			})
		})
	}

	printMenu := menu.AddSubmenu("Print")
	printMenu.Add("Print").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			_ = w.Print()
		})
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		BackgroundColour: application.NewRGB(33, 37, 41),
		Mac: application.MacWindow{
			DisableShadow: true,
		},
	})

	app.SetMenu(menu)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
