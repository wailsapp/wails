# Linux implementation tracker

This branch's tracker records changes to the default GTK4 / WebKitGTK 6.0
implementation and the legacy GTK3 path. Earlier implementation history is not
reconstructed here.

## 2026-09-10: Production builds with devtools

Status: ✅ COMPLETE — corrected build constraint verified by native GTK4 example builds and launches.

The example build sweep reproduced missing `openDevTools` and `enableDevTools`
methods when both `production` and `devtools` are enabled. The production stubs
exclude `devtools`, but the development implementation previously excluded all
production builds.

Decision: `v3/pkg/application/webview_window_linux_dev.go` uses
`(!production || devtools)`, matching the Windows and macOS selection rule.
Both GTK4 and legacy GTK3 use these Go methods and retain their existing native
implementations; no native API differences are introduced.

Verification: the final example sweep compiled and launched 100 Linux examples using `production,devtools` (with additional example-specific tags where needed). `TestLinuxDevtoolsImplementationSelectedForEveryRunMode` verifies exactly one implementation for development, devtools, production and production+devtools on both GTK stacks. Native GTK3 compilation remains in the platform handoff because GTK3 / WebKit2GTK 4.1 development libraries are not installed on this host.
