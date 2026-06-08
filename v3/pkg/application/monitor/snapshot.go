package monitor

import "sync/atomic"

// Snapshot is a point-in-time description of the running app: its windows, the
// bound methods the frontend can call, the events the backend listens for, and
// app metadata. It is requested on demand by a consumer (the TUI) and is
// independent of the live trace stream.
type Snapshot struct {
	App       AppInfo             `json:"app"`
	Windows   []WindowInfo        `json:"windows"`
	Bindings  []BindingInfo       `json:"bindings"`
	Listeners []EventListenerInfo `json:"listeners"`
}

// AppInfo is static application metadata.
type AppInfo struct {
	Name      string `json:"name"`
	PID       int    `json:"pid"`
	Platform  string `json:"platform"`
	Transport string `json:"transport,omitempty"`
	DevServer string `json:"devServer,omitempty"`
	DebugMode bool   `json:"debugMode"`
}

// WindowInfo describes one live window.
type WindowInfo struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	X          int     `json:"x"` // absolute desktop coordinates
	Y          int     `json:"y"`
	RelX       int     `json:"relX"` // position relative to current screen
	RelY       int     `json:"relY"`
	Border     Insets  `json:"border"` // window chrome thickness
	Focused    bool    `json:"focused"`
	Fullscreen bool    `json:"fullscreen"`
	Maximised  bool    `json:"maximised"`
	Minimised  bool    `json:"minimised"`
	Resizable  bool    `json:"resizable"`
	IgnoreMouse bool   `json:"ignoreMouse"`
	Zoom       float64 `json:"zoom"`
	Screen     *ScreenInfo `json:"screen,omitempty"`
}

// Insets are per-edge sizes (window borders).
type Insets struct {
	Left   int `json:"left"`
	Right  int `json:"right"`
	Top    int `json:"top"`
	Bottom int `json:"bottom"`
}

// ScreenInfo describes the display a window is currently on.
type ScreenInfo struct {
	Name        string  `json:"name"`
	ID          string  `json:"id"`
	IsPrimary   bool    `json:"isPrimary"`
	ScaleFactor float32 `json:"scaleFactor"` // DPI / 96
	Rotation    float32 `json:"rotation"`
	BoundsW     int     `json:"boundsW"`
	BoundsH     int     `json:"boundsH"`
	BoundsX     int     `json:"boundsX"`
	BoundsY     int     `json:"boundsY"`
	WorkW       int     `json:"workW"`
	WorkH       int     `json:"workH"`
}

// ParamInfo is one input/output parameter of a bound method.
type ParamInfo struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

// BindingInfo describes one bound method the frontend can call.
type BindingInfo struct {
	FQN     string      `json:"fqn"`
	Name    string      `json:"name"`
	ID      uint32      `json:"id"`
	Service string      `json:"service,omitempty"`
	Comment string      `json:"comment,omitempty"`
	Inputs  []ParamInfo `json:"inputs,omitempty"`
	Outputs []ParamInfo `json:"outputs,omitempty"`
}

// EventListenerInfo describes a backend-registered listener for an event name.
type EventListenerInfo struct {
	Name  string `json:"name"`
	Count int    `json:"count"` // number of registered listeners
}

// describeFn is the app-provided snapshot collector. It is registered by
// pkg/application (which can see globalApplication) so the monitor package
// stays free of an import cycle. nil when the app has not registered one.
var describeFn atomic.Pointer[func() *Snapshot]

// SetDescribeFunc registers the snapshot collector. Pass nil to clear.
func SetDescribeFunc(fn func() *Snapshot) {
	if fn == nil {
		describeFn.Store(nil)
		return
	}
	describeFn.Store(&fn)
}

// describe invokes the registered collector, or returns nil if none is set.
func describe() *Snapshot {
	if p := describeFn.Load(); p != nil {
		return (*p)()
	}
	return nil
}
