// Package runtime is a temporary compatibility bridge for projects migrated
// from Wails v2 via the `wails3 migrate` command.
//
// It mirrors the Wails v2 runtime API (github.com/wailsapp/wails/v2/pkg/runtime),
// a set of context-first free functions, implemented on top of the Wails v3
// application API (github.com/wailsapp/wails/v3/pkg/application). Migrated v2
// code keeps calling functions such as runtime.WindowSetTitle(a.ctx, ...)
// unchanged; only the import path changes. The context parameter is accepted
// for source compatibility and ignored.
//
// New code should use github.com/wailsapp/wails/v3/pkg/application directly.
// Each function in this package documents its v3 equivalent so that call
// sites can be migrated incrementally and this package eventually removed.
package runtime
