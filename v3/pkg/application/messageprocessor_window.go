package application

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	WindowPosition                   = 0
	WindowCenter                     = 1
	WindowClose                      = 2
	WindowDisableSizeConstraints     = 3
	WindowEnableSizeConstraints      = 4
	WindowFocus                      = 5
	WindowForceReload                = 6
	WindowFullscreen                 = 7
	WindowGetScreen                  = 8
	WindowGetZoom                    = 9
	WindowHeight                     = 10
	WindowHide                       = 11
	WindowIsFocused                  = 12
	WindowIsFullscreen               = 13
	WindowIsMaximised                = 14
	WindowIsMinimised                = 15
	WindowMaximise                   = 16
	WindowMinimise                   = 17
	WindowName                       = 18
	WindowOpenDevTools               = 19
	WindowRelativePosition           = 20
	WindowReload                     = 21
	WindowResizable                  = 22
	WindowRestore                    = 23
	WindowSetPosition                = 24
	WindowSetAlwaysOnTop             = 25
	WindowSetBackgroundColour        = 26
	WindowSetFrameless               = 27
	WindowSetFullscreenButtonEnabled = 28
	WindowSetMaxSize                 = 29
	WindowSetMinSize                 = 30
	WindowSetRelativePosition        = 31
	WindowSetResizable               = 32
	WindowSetSize                    = 33
	WindowSetTitle                   = 34
	WindowSetZoom                    = 35
	WindowShow                       = 36
	WindowSize                       = 37
	WindowToggleFullscreen           = 38
	WindowToggleMaximise             = 39
	WindowUnFullscreen               = 40
	WindowUnMaximise                 = 41
	WindowUnMinimise                 = 42
	WindowWidth                      = 43
	WindowZoom                       = 44
	WindowZoomIn                     = 45
	WindowZoomOut                    = 46
	WindowZoomReset                  = 47
)

var windowMethodNames = map[int]string{
	WindowPosition:                   "Position",
	WindowCenter:                     "Center",
	WindowClose:                      "Close",
	WindowDisableSizeConstraints:     "DisableSizeConstraints",
	WindowEnableSizeConstraints:      "EnableSizeConstraints",
	WindowFocus:                      "Focus",
	WindowForceReload:                "ForceReload",
	WindowFullscreen:                 "Fullscreen",
	WindowGetScreen:                  "GetScreen",
	WindowGetZoom:                    "GetZoom",
	WindowHeight:                     "Height",
	WindowHide:                       "Hide",
	WindowIsFocused:                  "IsFocused",
	WindowIsFullscreen:               "IsFullscreen",
	WindowIsMaximised:                "IsMaximised",
	WindowIsMinimised:                "IsMinimised",
	WindowMaximise:                   "Maximise",
	WindowMinimise:                   "Minimise",
	WindowName:                       "Name",
	WindowOpenDevTools:               "OpenDevTools",
	WindowRelativePosition:           "RelativePosition",
	WindowReload:                     "Reload",
	WindowResizable:                  "Resizable",
	WindowRestore:                    "Restore",
	WindowSetPosition:                "SetPosition",
	WindowSetAlwaysOnTop:             "SetAlwaysOnTop",
	WindowSetBackgroundColour:        "SetBackgroundColour",
	WindowSetFrameless:               "SetFrameless",
	WindowSetFullscreenButtonEnabled: "SetFullscreenButtonEnabled",
	WindowSetMaxSize:                 "SetMaxSize",
	WindowSetMinSize:                 "SetMinSize",
	WindowSetRelativePosition:        "SetRelativePosition",
	WindowSetResizable:               "SetResizable",
	WindowSetSize:                    "SetSize",
	WindowSetTitle:                   "SetTitle",
	WindowSetZoom:                    "SetZoom",
	WindowShow:                       "Show",
	WindowSize:                       "Size",
	WindowToggleFullscreen:           "ToggleFullscreen",
	WindowToggleMaximise:             "ToggleMaximise",
	WindowUnFullscreen:               "UnFullscreen",
	WindowUnMaximise:                 "UnMaximise",
	WindowUnMinimise:                 "UnMinimise",
	WindowWidth:                      "Width",
	WindowZoom:                       "Zoom",
	WindowZoomIn:                     "ZoomIn",
	WindowZoomOut:                    "ZoomOut",
	WindowZoomReset:                  "ZoomReset",
}

func (m *MessageProcessor) processWindowMethod(method int, rw http.ResponseWriter, _ *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid window call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	switch method {
	case WindowPosition:
		x, y := window.Position()
		m.json(rw, map[string]interface{}{
			"x": x,
			"y": y,
		})
	case WindowCenter:
		window.Center()
		m.ok(rw)
	case WindowClose:
		window.Close()
		m.ok(rw)
	case WindowDisableSizeConstraints:
		window.DisableSizeConstraints()
		m.ok(rw)
	case WindowEnableSizeConstraints:
		window.EnableSizeConstraints()
		m.ok(rw)
	case WindowFocus:
		window.Focus()
		m.ok(rw)
	case WindowForceReload:
		window.ForceReload()
		m.ok(rw)
	case WindowFullscreen:
		window.Fullscreen()
		m.ok(rw)
	case WindowGetScreen:
		screen, err := window.GetScreen()
		if err != nil {
			m.httpError(rw, "Window.GetScreen failed:", err)
			return
		}
		m.json(rw, screen)
	case WindowGetZoom:
		zoom := window.GetZoom()
		m.json(rw, zoom)
	case WindowHeight:
		height := window.Height()
		m.json(rw, height)
	case WindowHide:
		window.Hide()
		m.ok(rw)
	case WindowIsFocused:
		isFocused := window.IsFocused()
		m.json(rw, isFocused)
	case WindowIsFullscreen:
		isFullscreen := window.IsFullscreen()
		m.json(rw, isFullscreen)
	case WindowIsMaximised:
		isMaximised := window.IsMaximised()
		m.json(rw, isMaximised)
	case WindowIsMinimised:
		isMinimised := window.IsMinimised()
		m.json(rw, isMinimised)
	case WindowMaximise:
		window.Maximise()
		m.ok(rw)
	case WindowMinimise:
		window.Minimise()
		m.ok(rw)
	case WindowName:
		name := window.Name()
		m.json(rw, name)
	case WindowRelativePosition:
		x, y := window.RelativePosition()
		m.json(rw, map[string]interface{}{
			"x": x,
			"y": y,
		})
	case WindowReload:
		window.Reload()
		m.ok(rw)
	case WindowResizable:
		resizable := window.Resizable()
		m.json(rw, resizable)
	case WindowRestore:
		window.Restore()
		m.ok(rw)
	case WindowSetPosition:
		x := args.Int("x")
		if x == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'x'"))
			return
		}
		y := args.Int("y")
		if y == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'y'"))
			return
		}
		window.SetPosition(*x, *y)
		m.ok(rw)
	case WindowSetAlwaysOnTop:
		alwaysOnTop := args.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'alwaysOnTop'"))
			return
		}
		window.SetAlwaysOnTop(*alwaysOnTop)
		m.ok(rw)
	case WindowSetBackgroundColour:
		r := args.UInt8("r")
		if r == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'r'"))
			return
		}
		g := args.UInt8("g")
		if g == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'g'"))
			return
		}
		b := args.UInt8("b")
		if b == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'b'"))
			return
		}
		a := args.UInt8("a")
		if a == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'a'"))
			return
		}
		window.SetBackgroundColour(RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		m.ok(rw)
	case WindowSetFrameless:
		frameless := args.Bool("frameless")
		if frameless == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'frameless'"))
			return
		}
		window.SetFrameless(*frameless)
		m.ok(rw)
	case WindowSetMaxSize:
		width := args.Int("width")
		if width == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'width'"))
			return
		}
		height := args.Int("height")
		if height == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'height'"))
			return
		}
		window.SetMaxSize(*width, *height)
		m.ok(rw)
	case WindowSetMinSize:
		width := args.Int("width")
		if width == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'width'"))
			return
		}
		height := args.Int("height")
		if height == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'height'"))
			return
		}
		window.SetMinSize(*width, *height)
		m.ok(rw)
	case WindowSetRelativePosition:
		x := args.Int("x")
		if x == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'x'"))
			return
		}
		y := args.Int("y")
		if y == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'y'"))
			return
		}
		window.SetRelativePosition(*x, *y)
		m.ok(rw)
	case WindowSetResizable:
		resizable := args.Bool("resizable")
		if resizable == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'resizable'"))
			return
		}
		window.SetResizable(*resizable)
		m.ok(rw)
	case WindowSetSize:
		width := args.Int("width")
		if width == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'width'"))
			return
		}
		height := args.Int("height")
		if height == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'height'"))
			return
		}
		window.SetSize(*width, *height)
		m.ok(rw)
	case WindowSetTitle:
		title := args.String("title")
		if title == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing argument 'title'"))
			return
		}
		window.SetTitle(*title)
		m.ok(rw)
	case WindowSetZoom:
		zoom := args.Float64("zoom")
		if zoom == nil {
			m.httpError(rw, "Invalid window call:", errors.New("missing or invalid argument 'zoom'"))
			return
		}
		window.SetZoom(*zoom)
		m.ok(rw)
	case WindowShow:
		window.Show()
		m.ok(rw)
	case WindowSize:
		width, height := window.Size()
		m.json(rw, map[string]interface{}{
			"width":  width,
			"height": height,
		})
	case WindowOpenDevTools:
		window.OpenDevTools()
		m.ok(rw)
	case WindowToggleFullscreen:
		window.ToggleFullscreen()
		m.ok(rw)
	case WindowToggleMaximise:
		window.ToggleMaximise()
		m.ok(rw)
	case WindowUnFullscreen:
		window.UnFullscreen()
		m.ok(rw)
	case WindowUnMaximise:
		window.UnMaximise()
		m.ok(rw)
	case WindowUnMinimise:
		window.UnMinimise()
		m.ok(rw)
	case WindowWidth:
		width := window.Width()
		m.json(rw, width)
	case WindowZoom:
		window.Zoom()
		m.ok(rw)
	case WindowZoomIn:
		window.ZoomIn()
		m.ok(rw)
	case WindowZoomOut:
		window.ZoomOut()
		m.ok(rw)
	case WindowZoomReset:
		window.ZoomReset()
		m.ok(rw)
	default:
		m.httpError(rw, "Invalid window call:", fmt.Errorf("unknown method %d", method))
		return
	}

	m.Info("Runtime call:", "method", "Window."+windowMethodNames[method])
}
