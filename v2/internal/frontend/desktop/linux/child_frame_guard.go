//go:build linux

package linux

import _ "embed"

// The browser enforces cross-origin access to the trusted parent bridge.
// Same-origin child frames remain within the parent document trust boundary.
//
//go:embed child_frame_guard.js
var childFrameGuardJS string
