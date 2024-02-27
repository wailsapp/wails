package application

import (
	"net/http"
)

const (
	WindowCenter              = 0
	WindowSetTitle            = 1
	WindowFullscreen          = 2
	WindowUnFullscreen        = 3
	WindowSetSize             = 4
	WindowSize                = 5
	WindowSetMaxSize          = 6
	WindowSetMinSize          = 7
	WindowSetAlwaysOnTop      = 8
	WindowSetRelativePosition = 9
	WindowRelativePosition    = 10
	WindowScreen              = 11
	WindowHide                = 12
	WindowMaximise            = 13
	WindowUnMaximise          = 14
	WindowToggleMaximise      = 15
	WindowMinimise            = 16
	WindowUnMinimise          = 17
	WindowRestore             = 18
	WindowShow                = 19
	WindowClose               = 20
	WindowSetBackgroundColour = 21
	WindowSetResizable        = 22
	WindowWidth               = 23
	WindowHeight              = 24
	WindowZoomIn              = 25
	WindowZoomOut             = 26
	WindowZoomReset           = 27
	WindowGetZoomLevel        = 28
	WindowSetZoomLevel        = 29
)

var windowMethodNames = map[int]string{
	WindowCenter:              "Center",
	WindowSetTitle:            "SetTitle",
	WindowFullscreen:          "Fullscreen",
	WindowUnFullscreen:        "UnFullscreen",
	WindowSetSize:             "SetSize",
	WindowSize:                "Size",
	WindowSetMaxSize:          "SetMaxSize",
	WindowSetMinSize:          "SetMinSize",
	WindowSetAlwaysOnTop:      "SetAlwaysOnTop",
	WindowSetRelativePosition: "SetRelativePosition",
	WindowRelativePosition:    "RelativePosition",
	WindowScreen:              "Screen",
	WindowHide:                "Hide",
	WindowMaximise:            "Maximise",
	WindowUnMaximise:          "UnMaximise",
	WindowToggleMaximise:      "ToggleMaximise",
	WindowMinimise:            "Minimise",
	WindowUnMinimise:          "UnMinimise",
	WindowRestore:             "Restore",
	WindowShow:                "Show",
	WindowClose:               "Close",
	WindowSetBackgroundColour: "SetBackgroundColour",
	WindowSetResizable:        "SetResizable",
	WindowWidth:               "Width",
	WindowHeight:              "Height",
	WindowZoomIn:              "ZoomIn",
	WindowZoomOut:             "ZoomOut",
	WindowZoomReset:           "ZoomReset",
	WindowGetZoomLevel:        "GetZoomLevel",
	WindowSetZoomLevel:        "SetZoomLevel",
}

func (m *MessageProcessor) processWindowMethod(method int, rw http.ResponseWriter, _ *http.Request, window Window, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}

	switch method {
	case WindowSetTitle:
		title := args.String("title")
		if title == nil {
			m.Error("SetTitle: title is required")
			return
		}
		window.SetTitle(*title)
		m.ok(rw)
	case WindowSetSize:
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetSize Message")
			return
		}
		window.SetSize(*width, *height)
		m.ok(rw)
	case WindowSetRelativePosition:
		x := args.Int("x")
		y := args.Int("y")
		if x == nil || y == nil {
			m.Error("Invalid SetRelativePosition Message")
			return
		}
		window.SetRelativePosition(*x, *y)
		m.ok(rw)
	case WindowFullscreen:
		window.Fullscreen()
		m.ok(rw)
	case WindowUnFullscreen:
		window.UnFullscreen()
		m.ok(rw)
	case WindowMinimise:
		window.Minimise()
		m.ok(rw)
	case WindowUnMinimise:
		window.UnMinimise()
		m.ok(rw)
	case WindowMaximise:
		window.Maximise()
		m.ok(rw)
	case WindowUnMaximise:
		window.UnMaximise()
		m.ok(rw)
	case WindowToggleMaximise:
		window.ToggleMaximise()
		m.ok(rw)
	case WindowRestore:
		window.Restore()
		m.ok(rw)
	case WindowShow:
		window.Show()
		m.ok(rw)
	case WindowHide:
		window.Hide()
		m.ok(rw)
	case WindowClose:
		window.Close()
		m.ok(rw)
	case WindowCenter:
		window.Center()
		m.ok(rw)
	case WindowSize:
		width, height := window.Size()
		m.json(rw, map[string]interface{}{
			"width":  width,
			"height": height,
		})
	case WindowRelativePosition:
		x, y := window.RelativePosition()
		m.json(rw, map[string]interface{}{
			"x": x,
			"y": y,
		})
	case WindowSetBackgroundColour:
		r := args.UInt8("r")
		if r == nil {
			m.Error("Invalid SetBackgroundColour Message: 'r' value required")
			return
		}
		g := args.UInt8("g")
		if g == nil {
			m.Error("Invalid SetBackgroundColour Message: 'g' value required")
			return
		}
		b := args.UInt8("b")
		if b == nil {
			m.Error("Invalid SetBackgroundColour Message: 'b' value required")
			return
		}
		a := args.UInt8("a")
		if a == nil {
			m.Error("Invalid SetBackgroundColour Message: 'a' value required")
			return
		}
		window.SetBackgroundColour(RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		m.ok(rw)
	case WindowSetAlwaysOnTop:
		alwaysOnTop := args.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			m.Error("Invalid SetAlwaysOnTop Message: 'alwaysOnTop' value required")
			return
		}
		window.SetAlwaysOnTop(*alwaysOnTop)
		m.ok(rw)
	case WindowSetResizable:
		resizable := args.Bool("resizable")
		if resizable == nil {
			m.Error("Invalid SetResizable Message: 'resizable' value required")
			return
		}
		window.SetResizable(*resizable)
		m.ok(rw)
	case WindowSetMinSize:
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMinSize Message")
			return
		}
		window.SetMinSize(*width, *height)
		m.ok(rw)
	case WindowSetMaxSize:
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMaxSize Message")
			return
		}
		window.SetMaxSize(*width, *height)
		m.ok(rw)
	case WindowWidth:
		width := window.Width()
		m.json(rw, map[string]interface{}{
			"width": width,
		})
	case WindowHeight:
		height := window.Height()
		m.json(rw, map[string]interface{}{
			"height": height,
		})
	case WindowZoomIn:
		window.ZoomIn()
		m.ok(rw)
	case WindowZoomOut:
		window.ZoomOut()
		m.ok(rw)
	case WindowZoomReset:
		window.ZoomReset()
		m.ok(rw)
	case WindowGetZoomLevel:
		zoomLevel := window.GetZoom()
		m.json(rw, map[string]interface{}{
			"zoomLevel": zoomLevel,
		})
	case WindowScreen:
		screen, err := window.GetScreen()
		if err != nil {
			m.httpError(rw, err.Error())
			return
		}
		m.json(rw, screen)
	case WindowSetZoomLevel:
		zoomLevel := args.Float64("zoomLevel")
		if zoomLevel == nil {
			m.Error("Invalid SetZoom Message: invalid 'zoomLevel' value")
			return
		}
		window.SetZoom(*zoomLevel)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown window method id: %d", method)
	}

	m.Info("Runtime Call:", "method", "Window."+windowMethodNames[method])
}
