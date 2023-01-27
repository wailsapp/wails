package application

import (
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/options"
)

func (m *MessageProcessor) processWindowMethod(method string, rw http.ResponseWriter, r *http.Request) {

	params := QueryParams(r.URL.Query())

	var targetWindow = m.window
	windowID := params.UInt("windowID")
	if windowID != nil {
		// Get window for ID
		targetWindow = globalApplication.getWindowForID(*windowID)
		if targetWindow == nil {
			m.Error("Window ID %s not found", *windowID)
			return
		}
	}

	switch method {
	case "SetTitle":
		title := params.String("title")
		if title == nil {
			m.Error("SetTitle: title is required")
			return
		}
		targetWindow.SetTitle(*title)
		m.json(rw, nil)
	case "SetSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetSize message")
			return
		}
		targetWindow.SetSize(*width, *height)
		m.json(rw, nil)
	case "SetPosition":
		x := params.Int("x")
		y := params.Int("y")
		if x == nil || y == nil {
			m.Error("Invalid SetPosition message")
			return
		}
		targetWindow.SetPosition(*x, *y)
		m.json(rw, nil)
	case "Fullscreen":
		targetWindow.Fullscreen()
		m.json(rw, nil)
	case "UnFullscreen":
		targetWindow.UnFullscreen()
		m.json(rw, nil)
	case "Minimise":
		targetWindow.Minimize()
		m.json(rw, nil)
	case "UnMinimise":
		targetWindow.UnMinimise()
		m.json(rw, nil)
	case "Maximise":
		targetWindow.Maximise()
		m.json(rw, nil)
	case "UnMaximise":
		targetWindow.UnMaximise()
		m.json(rw, nil)
	case "Show":
		targetWindow.Show()
		m.json(rw, nil)
	case "Hide":
		targetWindow.Hide()
		m.json(rw, nil)
	case "Close":
		targetWindow.Close()
		m.json(rw, nil)
	case "Center":
		targetWindow.Center()
		m.json(rw, nil)
	case "Size":
		width, height := targetWindow.Size()
		m.json(rw, map[string]interface{}{
			"width":  width,
			"height": height,
		})
	case "Position":
		x, y := targetWindow.Position()
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
		targetWindow.SetBackgroundColour(&options.RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		m.json(rw, nil)
	case "SetAlwaysOnTop":
		alwaysOnTop := params.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			m.Error("Invalid SetAlwaysOnTop message: 'alwaysOnTop' value required")
			return
		}
		targetWindow.SetAlwaysOnTop(*alwaysOnTop)
		m.json(rw, nil)
	case "SetResizable":
		resizable := params.Bool("resizable")
		if resizable == nil {
			m.Error("Invalid SetResizable message: 'resizable' value required")
			return
		}
		targetWindow.SetResizable(*resizable)
		m.json(rw, nil)
	case "SetMinSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMinSize message")
			return
		}
		targetWindow.SetMinSize(*width, *height)
		m.json(rw, nil)
	case "SetMaxSize":
		width := params.Int("width")
		height := params.Int("height")
		if width == nil || height == nil {
			m.Error("Invalid SetMaxSize message")
			return
		}
		targetWindow.SetMaxSize(*width, *height)
		m.json(rw, nil)
	case "Width":
		width := targetWindow.Width()
		m.json(rw, map[string]interface{}{
			"width": width,
		})
	case "Height":
		height := targetWindow.Height()
		m.json(rw, map[string]interface{}{
			"height": height,
		})
	case "ZoomIn":
		targetWindow.ZoomIn()
		m.json(rw, nil)
	case "ZoomOut":
		targetWindow.ZoomOut()
		m.json(rw, nil)
	case "ZoomReset":
		targetWindow.ZoomReset()
		m.json(rw, nil)
	case "GetZoom":
		zoomLevel := targetWindow.GetZoom()
		m.json(rw, map[string]interface{}{
			"zoomLevel": zoomLevel,
		})
	case "SetZoom":
		zoomLevel := params.Float64("zoomLevel")
		if zoomLevel == nil {
			m.Error("Invalid SetZoom message: invalid 'zoomLevel' value")
			return
		}
		targetWindow.SetZoom(*zoomLevel)
		m.json(rw, nil)
	default:
		m.httpError(rw, "Unknown window method: %s", method)
	}

}
