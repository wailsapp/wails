package events_only

import (
	"github.com/wailsapp/wails/v3/internal/generator/testcases/no_bindings_here/more"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const eventPrefix = "events_only" + `:`

var registerStringEvent = application.RegisterEvent[string]

func registerIntEvent(name string) {
	application.RegisterEvent[int](name)
}

func registerSliceEvent[T any]() {
	application.RegisterEvent[[]T]("parametric")
}

func init() {
	application.RegisterEvent[[]more.StringPtr](eventPrefix + "other")
	application.RegisterEvent[string]("common:ApplicationStarted")
	registerStringEvent("indirect_var")
	registerIntEvent("indirect_fn")
	registerSliceEvent[uintptr]()
}
