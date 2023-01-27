package application

import (
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/options"
)

func (m *MessageProcessor) processWindowMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "SetTitle":
		title := params.String("title")
		if title == nil {
			m.Error("SetTitle: title is required")
			return
		}
		window.SetTitle(*title)
		m.ok(rw)
	case "SetSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetSize message")
			return
		}
		window.SetSize(*width, *height)
		m.ok(rw)
	case "SetPosition":
		x := params.Int("x")
		y := params.Int("y")
		if x == nil || y == nil {
			m.Error("Invalid SetPosition message")
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
		r := params.UInt8("r")
		if r == nil {
			m.Error("Invalid SetBackgroundColour message: 'r' value required")
			return
		}
		g := params.UInt8("g")
		if g == nil {
			m.Error("Invalid SetBackgroundColour message: 'g' value required")
			return
		}
		b := params.UInt8("b")
		if b == nil {
			m.Error("Invalid SetBackgroundColour message: 'b' value required")
			return
		}
		a := params.UInt8("a")
		if a == nil {
			m.Error("Invalid SetBackgroundColour message: 'a' value required")
			return
		}
		window.SetBackgroundColour(&options.RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		m.ok(rw)
	case "SetAlwaysOnTop":
		alwaysOnTop := params.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			m.Error("Invalid SetAlwaysOnTop message: 'alwaysOnTop' value required")
			return
		}
		window.SetAlwaysOnTop(*alwaysOnTop)
		m.ok(rw)
	case "SetResizable":
		resizable := params.Bool("resizable")
		if resizable == nil {
			m.Error("Invalid SetResizable message: 'resizable' value required")
			return
		}
		window.SetResizable(*resizable)
		m.ok(rw)
	case "SetMinSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMinSize message")
			return
		}
		window.SetMinSize(*width, *height)
		m.ok(rw)
	case "SetMaxSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMaxSize message")
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
	case "SetZoom":
		zoomLevel := params.Float64("zoomLevel")
		if zoomLevel == nil {
			m.Error("Invalid SetZoom message: invalid 'zoomLevel' value")
			return
		}
		window.SetZoom(*zoomLevel)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown window method: %s", method)
	}

}
