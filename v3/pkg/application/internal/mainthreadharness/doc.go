// Package mainthreadharness drives an AppKit run loop so the main thread
// dispatch tests can schedule work at it and nest a loop inside a dispatched
// callback, as a modal dialog or an open menu does.
//
// It is a package of its own, rather than a file in package application,
// because the go tool does not allow cgo in _test.go files and these helpers
// must not be linked into the shipped package. Only the tests import it.
//
// Everything here is macOS only; on other platforms the package is empty.
package mainthreadharness
