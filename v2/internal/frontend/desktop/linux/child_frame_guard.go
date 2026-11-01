//go:build linux

package linux

import _ "embed"

//go:embed child_frame_guard.js
var childFrameGuardJS string
