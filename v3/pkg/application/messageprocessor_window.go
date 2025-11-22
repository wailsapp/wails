package application

import (
	"fmt"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/errs"
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
	WindowToggleFrameless            = 40
	WindowUnFullscreen               = 41
	WindowUnMaximise                 = 42
	WindowUnMinimise                 = 43
	WindowWidth                      = 44
	WindowZoom                       = 45
	WindowZoomIn                     = 46
	WindowZoomOut                    = 47
	WindowZoomReset                  = 48
	WindowSnapAssist                 = 49
	WindowDropZoneDropped            = 50
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
	WindowToggleFrameless:            "ToggleFrameless",
	WindowUnFullscreen:               "UnFullscreen",
	WindowUnMaximise:                 "UnMaximise",
	WindowUnMinimise:                 "UnMinimise",
	WindowWidth:                      "Width",
	WindowZoom:                       "Zoom",
	WindowZoomIn:                     "ZoomIn",
	WindowZoomOut:                    "ZoomOut",
	WindowZoomReset:                  "ZoomReset",
	WindowDropZoneDropped:            "DropZoneDropped",
	WindowSnapAssist:                 "SnapAssist",
}

var unit = struct{}{}

func (m *MessageProcessor) processWindowMethod(
	req *RuntimeRequest,
	window Window,
) (any, error) {
	args := req.Args.AsMap()

	switch req.Method {
	case WindowPosition:
		x, y := window.Position()
		return map[string]interface{}{
			"x": x,
			"y": y,
		}, nil
	case WindowCenter:
		window.Center()
		return unit, nil
	case WindowClose:
		window.Close()
		return unit, nil
	case WindowDisableSizeConstraints:
		window.DisableSizeConstraints()
		return unit, nil
	case WindowEnableSizeConstraints:
		window.EnableSizeConstraints()
		return unit, nil
	case WindowFocus:
		window.Focus()
		return unit, nil
	case WindowForceReload:
		window.ForceReload()
		return unit, nil
	case WindowFullscreen:
		window.Fullscreen()
		return unit, nil
	case WindowGetScreen:
		screen, err := window.GetScreen()
		if err != nil {
			return nil, fmt.Errorf("Window.GetScreen failed: %w", err)
		}
		return screen, nil
	case WindowGetZoom:
		return window.GetZoom(), nil
	case WindowHeight:
		return window.Height(), nil
	case WindowHide:
		window.Hide()
		return unit, nil
	case WindowIsFocused:
		return window.IsFocused(), nil
	case WindowIsFullscreen:
		return window.IsFullscreen(), nil
	case WindowIsMaximised:
		return window.IsMaximised(), nil
	case WindowIsMinimised:
		return window.IsMinimised(), nil
	case WindowMaximise:
		window.Maximise()
		return unit, nil
	case WindowMinimise:
		window.Minimise()
		return unit, nil
	case WindowName:
		return window.Name(), nil
	case WindowRelativePosition:
		x, y := window.RelativePosition()
		return map[string]interface{}{
			"x": x,
			"y": y,
		}, nil
	case WindowReload:
		window.Reload()
		return unit, nil
	case WindowResizable:
		return window.Resizable(), nil
	case WindowRestore:
		window.Restore()
		return unit, nil
	case WindowSetPosition:
		x := args.Int("x")
		if x == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'x'")
		}
		y := args.Int("y")
		if y == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'y'")
		}
		window.SetPosition(*x, *y)
		return unit, nil
	case WindowSetAlwaysOnTop:
		alwaysOnTop := args.Bool("alwaysOnTop")
		if alwaysOnTop == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'alwaysOnTop'")
		}
		window.SetAlwaysOnTop(*alwaysOnTop)
		return unit, nil
	case WindowSetBackgroundColour:
		r := args.UInt8("r")
		if r == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'r'")
		}
		g := args.UInt8("g")
		if g == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'g'")
		}
		b := args.UInt8("b")
		if b == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'b'")
		}
		a := args.UInt8("a")
		if a == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'a'")
		}
		window.SetBackgroundColour(RGBA{
			Red:   *r,
			Green: *g,
			Blue:  *b,
			Alpha: *a,
		})
		return unit, nil
	case WindowSetFrameless:
		frameless := args.Bool("frameless")
		if frameless == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'frameless'")
		}
		window.SetFrameless(*frameless)
		return unit, nil
	case WindowSetMaxSize:
		width := args.Int("width")
		if width == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'width'")
		}
		height := args.Int("height")
		if height == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'height'")
		}
		window.SetMaxSize(*width, *height)
		return unit, nil
	case WindowSetMinSize:
		width := args.Int("width")
		if width == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'width'")
		}
		height := args.Int("height")
		if height == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'height'")
		}
		window.SetMinSize(*width, *height)
		return unit, nil
	case WindowSetRelativePosition:
		x := args.Int("x")
		if x == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'x'")
		}
		y := args.Int("y")
		if y == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'y'")
		}
		window.SetRelativePosition(*x, *y)
		return unit, nil
	case WindowSetResizable:
		resizable := args.Bool("resizable")
		if resizable == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'resizable'")
		}
		window.SetResizable(*resizable)
		return unit, nil
	case WindowSetSize:
		width := args.Int("width")
		if width == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'width'")
		}
		height := args.Int("height")
		if height == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'height'")
		}
		window.SetSize(*width, *height)
		return unit, nil
	case WindowSetTitle:
		title := args.String("title")
		if title == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'title'")
		}
		window.SetTitle(*title)
		return unit, nil
	case WindowSetZoom:
		zoom := args.Float64("zoom")
		if zoom == nil {
			return nil, errs.NewInvalidWindowCallErrorf("missing or invalid argument 'zoom'")
		}
		window.SetZoom(*zoom)
		return unit, nil
	case WindowShow:
		window.Show()
		return unit, nil
	case WindowSize:
		width, height := window.Size()
		return map[string]interface{}{
			"width":  width,
			"height": height,
		}, nil
	case WindowOpenDevTools:
		window.OpenDevTools()
		return unit, nil
	case WindowToggleFullscreen:
		window.ToggleFullscreen()
		return unit, nil
	case WindowToggleMaximise:
		window.ToggleMaximise()
		return unit, nil
	case WindowToggleFrameless:
		window.ToggleFrameless()
		return unit, nil
	case WindowUnFullscreen:
		window.UnFullscreen()
		return unit, nil
	case WindowUnMaximise:
		window.UnMaximise()
		return unit, nil
	case WindowUnMinimise:
		window.UnMinimise()
		return unit, nil
	case WindowWidth:
		return window.Width(), nil
	case WindowZoom:
		window.Zoom()
		return unit, nil
	case WindowZoomIn:
		window.ZoomIn()
		return unit, nil
	case WindowZoomOut:
		window.ZoomOut()
		return unit, nil
	case WindowZoomReset:
		window.ZoomReset()
		return unit, nil
	case WindowDropZoneDropped:
		m.Info(
			"[DragDropDebug] processWindowMethod: Entered WindowDropZoneDropped case",
		)

		slog.Info("[DragDropDebug] Raw 'args' payload string:", "data", string(req.Args.rawData))

		var payload fileDropPayload

		err := req.Args.ToStruct(&payload)
		if err != nil {
			return nil, errs.WrapInvalidWindowCallErrorf(err, "Error decoding file drop payload from 'args' parameter")
		}
		m.Info(
			"[DragDropDebug] processWindowMethod: Decoded payload from 'args'",
			"payload",
			fmt.Sprintf("%+v", payload),
		)

		dropDetails := &DropZoneDetails{
			X:          payload.X,
			Y:          payload.Y,
			ElementID:  payload.ElementDetails.ID,
			ClassList:  payload.ElementDetails.ClassList,
			Attributes: payload.ElementDetails.Attributes, // Assumes DropZoneDetails struct is updated to include this field
		}

		wvWindow, ok := window.(*WebviewWindow)
		if !ok {
			return nil, errs.NewInvalidWindowCallErrorf("Error: Target window is not a WebviewWindow for FilesDroppedWithContext")
		}

		msg := &dragAndDropMessage{
			windowId:  wvWindow.id,
			filenames: payload.Filenames,
			DropZone:  dropDetails,
		}

		m.Info(
			"[DragDropDebug] processApplicationMethod: Sending message to windowDragAndDropBuffer",
			"message",
			fmt.Sprintf("%+v", msg),
		)
		windowDragAndDropBuffer <- msg
		return unit, nil
	case WindowSnapAssist:
		window.SnapAssist()
		return unit, nil
	default:
		return nil, errs.NewInvalidWindowCallErrorf("Unknown method %d", req.Method)
	}
}

// ElementDetailsPayload holds detailed information about the drop target element.
type ElementDetailsPayload struct {
	ID         string            `json:"id"`
	ClassList  []string          `json:"classList"`
	Attributes map[string]string `json:"attributes"`
}

// Define a struct for the JSON payload from HandlePlatformFileDrop
type fileDropPayload struct {
	Filenames      []string              `json:"filenames"`
	X              int                   `json:"x"`
	Y              int                   `json:"y"`
	ElementDetails ElementDetailsPayload `json:"elementDetails"`
}
