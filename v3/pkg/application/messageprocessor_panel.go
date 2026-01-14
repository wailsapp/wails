package application

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

// Panel method constants for frontend-to-backend communication
const (
	PanelSetBounds    = 0
	PanelGetBounds    = 1
	PanelSetZIndex    = 2
	PanelSetURL       = 3
	PanelSetHTML      = 4
	PanelExecJS       = 5
	PanelReload       = 6
	PanelForceReload  = 7
	PanelShow         = 8
	PanelHide         = 9
	PanelIsVisible    = 10
	PanelSetZoom      = 11
	PanelGetZoom      = 12
	PanelFocus        = 13
	PanelIsFocused    = 14
	PanelOpenDevTools = 15
	PanelDestroy      = 16
	PanelName         = 17
)

var panelMethodNames = map[int]string{
	PanelSetBounds:    "SetBounds",
	PanelGetBounds:    "GetBounds",
	PanelSetZIndex:    "SetZIndex",
	PanelSetURL:       "SetURL",
	PanelSetHTML:      "SetHTML",
	PanelExecJS:       "ExecJS",
	PanelReload:       "Reload",
	PanelForceReload:  "ForceReload",
	PanelShow:         "Show",
	PanelHide:         "Hide",
	PanelIsVisible:    "IsVisible",
	PanelSetZoom:      "SetZoom",
	PanelGetZoom:      "GetZoom",
	PanelFocus:        "Focus",
	PanelIsFocused:    "IsFocused",
	PanelOpenDevTools: "OpenDevTools",
	PanelDestroy:      "Destroy",
	PanelName:         "Name",
}

func (m *MessageProcessor) processPanelMethod(
	req *RuntimeRequest,
	window Window,
) (any, error) {
	args := req.Args.AsMap()

	// Get the WebviewWindow to access panels
	ww, ok := window.(*WebviewWindow)
	if !ok {
		return nil, errs.NewInvalidRuntimeCallErrorf("window is not a WebviewWindow")
	}

	// Get panel name from args
	panelName := args.String("panel")
	if panelName == nil || *panelName == "" {
		return nil, errs.NewInvalidRuntimeCallErrorf("panel name is required")
	}

	// Get the panel
	panel := ww.GetPanel(*panelName)
	if panel == nil {
		// Try by ID
		panelID := args.UInt("panelId")
		if panelID != nil && *panelID > 0 {
			panel = ww.GetPanelByID(uint(*panelID))
		}
	}
	if panel == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("panel not found: %s", *panelName)
	}

	switch req.Method {
	case PanelSetBounds:
		x := args.Int("x")
		y := args.Int("y")
		width := args.Int("width")
		height := args.Int("height")
		if x == nil || y == nil || width == nil || height == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("x, y, width, and height are required")
		}
		panel.SetBounds(Rect{X: *x, Y: *y, Width: *width, Height: *height})
		return unit, nil

	case PanelGetBounds:
		bounds := panel.Bounds()
		return map[string]interface{}{
			"x":      bounds.X,
			"y":      bounds.Y,
			"width":  bounds.Width,
			"height": bounds.Height,
		}, nil

	case PanelSetZIndex:
		zIndex := args.Int("zIndex")
		if zIndex == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("zIndex is required")
		}
		panel.SetZIndex(*zIndex)
		return unit, nil

	case PanelSetURL:
		url := args.String("url")
		if url == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("url is required")
		}
		panel.SetURL(*url)
		return unit, nil

	case PanelSetHTML:
		html := args.String("html")
		if html == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("html is required")
		}
		panel.SetHTML(*html)
		return unit, nil

	case PanelExecJS:
		js := args.String("js")
		if js == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("js is required")
		}
		panel.ExecJS(*js)
		return unit, nil

	case PanelReload:
		panel.Reload()
		return unit, nil

	case PanelForceReload:
		panel.ForceReload()
		return unit, nil

	case PanelShow:
		panel.Show()
		return unit, nil

	case PanelHide:
		panel.Hide()
		return unit, nil

	case PanelIsVisible:
		return panel.IsVisible(), nil

	case PanelSetZoom:
		zoom := args.Float64("zoom")
		if zoom == nil {
			return nil, errs.NewInvalidRuntimeCallErrorf("zoom is required")
		}
		panel.SetZoom(*zoom)
		return unit, nil

	case PanelGetZoom:
		return panel.GetZoom(), nil

	case PanelFocus:
		panel.Focus()
		return unit, nil

	case PanelIsFocused:
		return panel.IsFocused(), nil

	case PanelOpenDevTools:
		panel.OpenDevTools()
		return unit, nil

	case PanelDestroy:
		panel.Destroy()
		return unit, nil

	case PanelName:
		return panel.Name(), nil

	default:
		return nil, fmt.Errorf("unknown panel method: %d", req.Method)
	}
}
