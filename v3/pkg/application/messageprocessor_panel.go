package application

import (
	"math"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

// Panel method constants for frontend-to-backend communication
const (
	PanelSetBounds = 0
	PanelGetBounds = 1
	PanelSetZIndex = 2
	// IDs 3-5 are reserved for removed content-control methods. Never reuse them:
	// a stale frontend must not dispatch navigation or script execution.
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
	if req == nil || req.Args == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("panel arguments are required")
	}
	handler, ok := panelMethodHandlers[req.Method]
	if !ok {
		return nil, errs.NewInvalidRuntimeCallErrorf("unknown panel method: %d", req.Method)
	}
	args := req.Args.AsMap()

	// Get the WebviewWindow to access panels
	ww, ok := window.(*WebviewWindow)
	if !ok || ww == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("window is not a WebviewWindow")
	}

	// Get panel name from args
	panelName := args.String("panel")
	var panel *WebviewPanel
	if panelName != nil && *panelName != "" {
		panel = ww.GetPanel(*panelName)
	}
	if panel == nil {
		panelID := panelInteger(args, "panelId")
		if panelID != nil && *panelID > 0 {
			panel = ww.GetPanelByID(uint(*panelID))
		}
	}
	if panel == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("panel not found: a valid panel name or panelId is required")
	}

	return handler(panel, args)
}

func handlePanelSetBounds(panel *WebviewPanel, args *MapArgs) (any, error) {
	x := panelInteger(args, "x")
	y := panelInteger(args, "y")
	width := panelInteger(args, "width")
	height := panelInteger(args, "height")
	if x == nil || y == nil || width == nil || height == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("x, y, width, and height are required")
	}
	if *width <= 0 || *height <= 0 {
		return nil, errs.NewInvalidRuntimeCallErrorf("width and height must be positive")
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
	zIndex := panelInteger(args, "zIndex")
	if zIndex == nil {
		return nil, errs.NewInvalidRuntimeCallErrorf("zIndex is required")
	}
	panel.SetZIndex(*zIndex)
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
	if zoom == nil || *zoom <= 0 || math.IsNaN(*zoom) || math.IsInf(*zoom, 0) {
		return nil, errs.NewInvalidRuntimeCallErrorf("zoom must be a positive finite number")
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

// panelInteger rejects truncation and overflow before passing geometry or IDs
// to native APIs, whose integer coordinates are 32-bit on desktop platforms.
func panelInteger(args *MapArgs, key string) *int {
	value := args.Int(key)
	if value == nil || *value < math.MinInt32 || *value > math.MaxInt32 {
		return nil
	}
	if raw := args.Float64(key); raw != nil && *raw != float64(*value) {
		return nil
	}
	return value
}
