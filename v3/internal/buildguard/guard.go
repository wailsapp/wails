//go:build ((linux && !android) || (darwin && !ios)) && !cgo && !purego && !server

package buildguard

// Building with CGO_ENABLED=0 on macOS or Linux needs the experimental
// CGO-free backend: add -tags purego (or set WAILS_PUREGO=1 for wails3), or
// build with CGO_ENABLED=1.
var _ = requires_CGO_ENABLED_1_or_build_tag_purego
