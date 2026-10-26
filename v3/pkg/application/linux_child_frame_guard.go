//go:build linux && cgo && !android

package application

import _ "embed"

// WebKitGTK delivers "external" script messages without frame information, so
// the childFrame gating used on macOS/iOS cannot be done natively. Instead this
// script is injected at document start into every frame and, in child frames,
// replaces window.webkit.messageHandlers.external with a wrapper that drops
// reserved "wails:" messages. Non-reserved messages still reach
// RawMessageHandler, matching macOS (which cannot be told apart from main-frame
// ones here; OriginInfo.IsMainFrame is not populated on Linux).
//
// Limitation: a same-origin child frame (about:blank, srcdoc, same-origin URL)
// can still reach window.top.webkit directly. Cross-origin and sandboxed
// frames cannot.
//
//go:embed linux_child_frame_guard.js
var linuxChildFrameGuardJS string
