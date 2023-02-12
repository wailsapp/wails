package application

import (
	"net/http"
)

func (m *MessageProcessor) processWindowMethod(method string, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err)
		return
	}

	switch method {
	case "SetTitle":
		title := args.String("title")
		if title == nil {
			m.Error("SetTitle: title is required")
			return
		}
		window.SetTitle(*title)
		m.ok(rw)
	case "SetSize":
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetSize Message")
			return
		}
		window.SetSize(*width, *height)
		m.ok(rw)
	case "SetPosition":
		x := args.Int("x")
		y := args.Int("y")
		if x == nil || y == nil {
			m.Error("Invalid SetPosition Message")
			return
		}
		window.SetPosition(*x, *y)
		m.ok(rw)
	case "Fullscreen":
		window.Fullscreen()
		m.ok(rw)
	case "UnFullscreen":
		window.UnFullscreen()
		m.ok(rw)
	case "Minimise":
		window.Minimize()
		m.ok(rw)
	case "UnMinimise":
		window.UnMinimise()
		m.ok(rw)
	case "Maximise":
		window.Maximise()
		m.ok(rw)
	case "UnMaximise":
		window.UnMaximise()
		m.ok(rw)
	case "Show":
		window.Show()
		m.ok(rw)
	case "Hide":
		window.Hide()
		m.ok(rw)
	case "Close":
		window.Close()
		m.ok(rw)
	case "Center":
		window.Center()
		m.ok(rw)
	case "Size":
		width, height := window.Size()
		m.json(rw, map[string]interface{}{
			"width":  width,
			"height": height,
		})
	case "Position":
		x, y := window.Position()
		m.json(rw, map[string]interface{}{
			"x": x,
			"y": y,
		})
	case "SetBackgroundColour":
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
		window.SetBackgroundColour(&RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		m.ok(rw)
	case "SetAlwaysOnTop":
		alwaysOnTop := args.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			m.Error("Invalid SetAlwaysOnTop Message: 'alwaysOnTop' value required")
			return
		}
		window.SetAlwaysOnTop(*alwaysOnTop)
		m.ok(rw)
	case "SetResizable":
		resizable := args.Bool("resizable")
		if resizable == nil {
			m.Error("Invalid SetResizable Message: 'resizable' value required")
			return
		}
		window.SetResizable(*resizable)
		m.ok(rw)
	case "SetMinSize":
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMinSize Message")
			return
		}
		window.SetMinSize(*width, *height)
		m.ok(rw)
	case "SetMaxSize":
		width := args.Int("width")
		height := args.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMaxSize Message")
			return
		}
		window.SetMaxSize(*width, *height)
		m.ok(rw)
	case "Width":
		width := window.Width()
		m.json(rw, map[string]interface{}{
			"width": width,
		})
	case "Height":
		height := window.Height()
		m.json(rw, map[string]interface{}{
			"height": height,
		})
	case "ZoomIn":
		window.ZoomIn()
		m.ok(rw)
	case "ZoomOut":
		window.ZoomOut()
		m.ok(rw)
	case "ZoomReset":
		window.ZoomReset()
		m.ok(rw)
	case "GetZoom":
		zoomLevel := window.GetZoom()
		m.json(rw, map[string]interface{}{
			"zoomLevel": zoomLevel,
		})
	case "Screen":
		screen, err := window.GetScreen()
		if err != nil {
			m.httpError(rw, err.Error())
			return
		}
		m.json(rw, screen)
	case "SetZoom":
		zoomLevel := args.Float64("zoomLevel")
		if zoomLevel == nil {
			m.Error("Invalid SetZoom Message: invalid 'zoomLevel' value")
			return
		}
		window.SetZoom(*zoomLevel)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown window method: %s", method)
	}

}
