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
	PanelReload       = 4
	PanelForceReload  = 5
	PanelShow         = 6
	PanelHide         = 7
	PanelIsVisible    = 8
	PanelSetZoom      = 9
	PanelGetZoom      = 10
	PanelFocus        = 11
	PanelIsFocused    = 12
	PanelOpenDevTools = 13
	PanelDestroy      = 14
	PanelName         = 15
)

var panelMethodNames = map[int]string{
	PanelSetBounds:    "SetBounds",
	PanelGetBounds:    "GetBounds",
	PanelSetZIndex:    "SetZIndex",
	PanelSetURL:       "SetURL",
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

// panelMethodHandler handles a specific panel method
type panelMethodHandler func(panel *WebviewPanel, args *MapArgs) (any, error)

// panelMethodHandlers maps method IDs to their handlers
var panelMethodHandlers = map[int]panelMethodHandler{
	PanelSetBounds:    handlePanelSetBounds,
	PanelGetBounds:    handlePanelGetBounds,
	PanelSetZIndex:    handlePanelSetZIndex,
	PanelSetURL:       handlePanelSetURL,
	PanelReload:       handlePanelReload,
	PanelForceReload:  handlePanelForceReload,
	PanelShow:         handlePanelShow,
	PanelHide:         handlePanelHide,
	PanelIsVisible:    handlePanelIsVisible,
	PanelSetZoom:      handlePanelSetZoom,
	PanelGetZoom:      handlePanelGetZoom,
	PanelFocus:        handlePanelFocus,
	PanelIsFocused:    handlePanelIsFocused,
	PanelOpenDevTools: handlePanelOpenDevTools,
	PanelDestroy:      handlePanelDestroy,
	PanelName:         handlePanelName,
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

	// Get the panel by name or ID
	panel := ww.GetPanel(*panelName)
	if panel == nil {
		panelID := args.UInt("panelId")
		if panelID != nil && *panelID > 0 {
			panel = ww.GetPanelByID(uint(*panelID))
		}
	}
	if panel == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("panel not found: %s", *panelName)
	}

	// Look up and execute the handler
	handler, ok := panelMethodHandlers[req.Method]
	if !ok {
		return nil, fmt.Errorf("unknown panel method: %d", req.Method)
	}
	return handler(panel, args)
}

func handlePanelSetBounds(panel *WebviewPanel, args *MapArgs) (any, error) {
	x := args.Int("x")
	y := args.Int("y")
	width := args.Int("width")
	height := args.Int("height")
	if x == nil || y == nil || width == nil || height == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("x, y, width, and height are required")
	}
	panel.SetBounds(Rect{X: *x, Y: *y, Width: *width, Height: *height})
	return unit, nil
}

func handlePanelGetBounds(panel *WebviewPanel, _ *MapArgs) (any, error) {
	bounds := panel.Bounds()
	return map[string]interface{}{
		"x":      bounds.X,
		"y":      bounds.Y,
		"width":  bounds.Width,
		"height": bounds.Height,
	}, nil
}

func handlePanelSetZIndex(panel *WebviewPanel, args *MapArgs) (any, error) {
	zIndex := args.Int("zIndex")
	if zIndex == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("zIndex is required")
	}
	panel.SetZIndex(*zIndex)
	return unit, nil
}

func handlePanelSetURL(panel *WebviewPanel, args *MapArgs) (any, error) {
	url := args.String("url")
	if url == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("url is required")
	}
	panel.SetURL(*url)
	return unit, nil
}

func handlePanelReload(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.Reload()
	return unit, nil
}

func handlePanelForceReload(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.ForceReload()
	return unit, nil
}

func handlePanelShow(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.Show()
	return unit, nil
}

func handlePanelHide(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.Hide()
	return unit, nil
}

func handlePanelIsVisible(panel *WebviewPanel, _ *MapArgs) (any, error) {
	return panel.IsVisible(), nil
}

func handlePanelSetZoom(panel *WebviewPanel, args *MapArgs) (any, error) {
	zoom := args.Float64("zoom")
	if zoom == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("zoom is required")
	}
	panel.SetZoom(*zoom)
	return unit, nil
}

func handlePanelGetZoom(panel *WebviewPanel, _ *MapArgs) (any, error) {
	return panel.GetZoom(), nil
}

func handlePanelFocus(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.Focus()
	return unit, nil
}

func handlePanelIsFocused(panel *WebviewPanel, _ *MapArgs) (any, error) {
	return panel.IsFocused(), nil
}

func handlePanelOpenDevTools(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.OpenDevTools()
	return unit, nil
}

func handlePanelDestroy(panel *WebviewPanel, _ *MapArgs) (any, error) {
	panel.Destroy()
	return unit, nil
}

func handlePanelName(panel *WebviewPanel, _ *MapArgs) (any, error) {
	return panel.Name(), nil
}
