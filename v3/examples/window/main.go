package main

import (
	"embed"
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

// This is a stub for non-windows platforms
var getExStyle = func() int {
	return 0
}

//go:embed assets/*
var assets embed.FS

type WindowService struct{}

// ==============================================
func (s *WindowService) SetPos(relative bool, x, y float64) {
	win := application.Get().CurrentWindow()
	initX, initY := win.Position()
	if relative {
		x += float64(initX)
		y += float64(initY)
	}
	win.SetPosition(int(x), int(y))
	currentX, currentY := win.Position()
	fmt.Printf("SetPos: %d, %d => %d, %d\n", initX, initY, currentX, currentY)
}
func (s *WindowService) SetSize(relative bool, wdt, hgt float64) {
	win := application.Get().CurrentWindow()
	initW, initH := win.Size()
	if relative {
		wdt += float64(initW)
		hgt += float64(initH)
	}
	win.SetSize(int(wdt), int(hgt))
	currentW, currentH := win.Size()
	fmt.Printf("SetSize: %d, %d => %d, %d\n", initW, initH, currentW, currentH)
}
func (s *WindowService) SetBounds(x, y, w, h float64) {
	win := application.Get().CurrentWindow()
	initR := win.Bounds()
	win.SetBounds(application.Rect{
		X:      int(x),
		Y:      int(y),
		Width:  int(w),
		Height: int(h),
	})
	currentR := win.Bounds()
	fmt.Printf("SetBounds: %+v => %+v\n", initR, currentR)
}
func (s *WindowService) GetBounds() application.Rect {
	win := application.Get().CurrentWindow()
	r := win.Bounds()
	mid := r.X + (r.Width-1)/2
	fmt.Printf("GetBounds: %+v: mid: %d\n", r, mid)
	return r
}

// ==============================================

func main() {
	app := application.New(application.Options{
		Name:        "WebviewWindow Demo",
		Description: "A demo of the WebviewWindow API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
	})
	app.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
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
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	} else {
		menu.AddRole(application.FileMenu)
	}
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
	if runtime.GOOS != "linux" {
		myMenu.Add("New WebviewWindow (Disable Minimise)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					MinimiseButtonState: application.ButtonDisabled,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Disable Maximise)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					MaximiseButtonState: application.ButtonDisabled,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Hide Minimise)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					MinimiseButtonState: application.ButtonHidden,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Always on top)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					AlwaysOnTop: true,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Hide Maximise)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					MaximiseButtonState: application.ButtonHidden,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		myMenu.Add("New WebviewWindow (Disable Close)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					CloseButtonState: application.ButtonDisabled,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
		myMenu.Add("New WebviewWindow (Hide Close)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					CloseButtonState: application.ButtonHidden,
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})

	}

	if runtime.GOOS == "windows" {
		myMenu.Add("New WebviewWindow (Custom ExStyle)").
			OnClick(func(ctx *application.Context) {
				app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
					Windows: application.WindowsWindow{
						DisableMenu: true,
						ExStyle:     getExStyle(),
					},
				}).
					SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
					SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
					SetURL("https://wails.io").
					Show()
				windowCounter++
			})
	}
	myMenu.Add("New WebviewWindow (Listen to Move)").
		OnClick(func(ctx *application.Context) {
			w := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				Windows: application.WindowsWindow{
					DisableMenu: true,
				},
			}).
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			w.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
				x, y := w.Position()
				fmt.Printf("WindowDidMove event triggered. New position: (%d, %d)\n", x, y)
			})
			windowCounter++
		})
	myMenu.Add("New WebviewWindow (Listen to Resize)").
		OnClick(func(ctx *application.Context) {
			w := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				Windows: application.WindowsWindow{
					DisableMenu: true,
				},
			}).
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			w.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
				width, height := w.Size()

				fmt.Printf("WindowDidResize event triggered. New size: (%d, %d)\n", width, height)
			})
			windowCounter++
		})
	myMenu.Add("New WebviewWindow (Hides on Close one time)").
		SetAccelerator("CmdOrCtrl+H").
		OnClick(func(ctx *application.Context) {
			app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
				// This will be called when the user clicks the close button
				// on the window. It will hide the window for 5 seconds.
				// If the user clicks the close button again, the window will
				// close.
				ShouldClose: func(window *application.WebviewWindow) bool {
					if !lo.Contains(hiddenWindows, window) {
						hiddenWindows = append(hiddenWindows, window)
						go func() {
							time.Sleep(5 * time.Second)
							window.Show()
						}()
						window.Hide()
						return false
					}
					// Remove the window from the hiddenWindows list
					hiddenWindows = lo.Without(hiddenWindows, window)
					return true
				},
			}).
				SetTitle("WebviewWindow "+strconv.Itoa(windowCounter)).
				SetRelativePosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Show()
			windowCounter++
		})
	myMenu.Add("New WebviewWindow (Frameless)").
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
	myMenu.Add("New WebviewWindow (Ignores mouse events)").
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
			w.SetMaxSize(600, 600)
			w.SetMaximiseButtonState(application.ButtonDisabled)
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

	positionMenu.Add("Set Relative Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetRelativePosition(0, 0)
		})
	})
	positionMenu.Add("Set Relative Position (Corner)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			screen, _ := w.GetScreen()
			w.SetRelativePosition(screen.WorkArea.Width-w.Width(), screen.WorkArea.Height-w.Height())
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

	positionMenu.Add("Center").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.Center()
		})
	})
	titleBarMenu := menu.AddSubmenu("Controls")
	titleBarMenu.Add("Disable Minimise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMinimiseButtonState(application.ButtonDisabled)
		})
	})
	titleBarMenu.Add("Enable Minimise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMinimiseButtonState(application.ButtonEnabled)
		})
	})
	titleBarMenu.Add("Hide Minimise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMinimiseButtonState(application.ButtonHidden)
		})
	})
	titleBarMenu.Add("Disable Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMaximiseButtonState(application.ButtonDisabled)
		})
	})
	titleBarMenu.Add("Enable Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMaximiseButtonState(application.ButtonEnabled)
		})
	})
	titleBarMenu.Add("Hide Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetMaximiseButtonState(application.ButtonHidden)
		})
	})
	titleBarMenu.Add("Disable Close").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetCloseButtonState(application.ButtonDisabled)
		})
	})
	titleBarMenu.Add("Enable Close").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetCloseButtonState(application.ButtonEnabled)
		})
	})
	titleBarMenu.Add("Hide Close").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.SetCloseButtonState(application.ButtonHidden)
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
	stateMenu.Add("Open Dev Tools").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.WebviewWindow) {
			w.OpenDevTools()
		})
	})

	if runtime.GOOS != "darwin" {
		stateMenu.Add("Flash for 5s").OnClick(func(ctx *application.Context) {
			currentWindow(func(w *application.WebviewWindow) {
				time.Sleep(2 * time.Second)
				w.Flash(true)
				time.Sleep(5 * time.Second)
				w.Flash(false)
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
		Title:            "Window Demo",
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
