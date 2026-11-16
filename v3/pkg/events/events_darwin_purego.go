//go:build darwin && !ios && purego

package events

// This is the CGO-free (purego) counterpart of events_darwin.go.
//
// The cgo build of this package uses events_darwin.go purely to compile two
// C helper functions (registerListener / hasListeners) that are consumed by the
// Objective-C event-handling layer. That file exposes NO Go-visible symbols
// (no //export directives, no Go functions), and no Go code in this package
// references those C helpers. Under the purego build there is no cgo and no
// Objective-C compilation unit, so there is nothing to reimplement here — this
// file exists only to provide the package declaration under the purego tag.
//
// All event constants and maps live in the platform-independent files
// (known_events.go, defaults.go, events.go) and are shared across builds.
