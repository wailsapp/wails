package runtime

import (
	"github.com/wailsapp/wails/v2/pkg/runtime/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime/events"
	"github.com/wailsapp/wails/v2/pkg/runtime/log"
	"github.com/wailsapp/wails/v2/pkg/runtime/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime/window"
)

var (
	Window = window.Window{}
	Menu   = menu.Menu{}
	Log    = log.Log{}
	Events = events.Events{}
	Dialog = dialog.Dialog{}
)
