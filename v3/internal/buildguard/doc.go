// Package buildguard turns a desktop build with no native backend into one
// clear error. On macOS and Linux the default backend uses cgo and the
// experimental CGO-free one needs the purego build tag, so CGO_ENABLED=0
// without -tags purego selects neither; without this package that fails
// with a long list of undefined symbols.
package buildguard
