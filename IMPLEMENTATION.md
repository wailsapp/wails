# WebKitGTK 6.0 / GTK4 Implementation Tracker

## Overview

This document tracks the implementation of WebKitGTK 6.0 (GTK4) support for Wails v3 on Linux.

**Goal**: Make GTK4/WebKitGTK 6.0 the DEFAULT build target, with GTK3/WebKit2GTK 4.1 available via `-tags gtk3` for legacy systems.

## Architecture Decisions

### Decision 1: GTK4 as Default (2026-01-04)
**Context**: Need to support modern Linux distributions with GTK4 while maintaining backward compatibility.

**Decision**: GTK4 is the new default (no build tag required). GTK3 requires explicit `-tags gtk3`.

**Rationale**:
- Ubuntu 22.04+ now has WebKitGTK 6.0 in official repos
- GTK4 is the future direction for Linux desktop
- Matches the pattern used for other platform-specific features

**Build Tags**:
- Default (no tag): `//go:build linux && cgo && !gtk3 && !android`
- Legacy: `//go:build linux && cgo && gtk3 && !android`

### Decision 2: pkg-config Libraries (2026-01-04)
**GTK4/WebKitGTK 6.0**:
```
#cgo linux pkg-config: gtk4 webkitgtk-6.0 libsoup-3.0
```

**GTK3/WebKit2GTK 4.1** (legacy):
```
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1 libsoup-3.0
```

### Decision 3: Wayland Window Positioning (2026-01-04)
**Context**: GTK4/Wayland doesn't support arbitrary window positioning - this is a Wayland protocol limitation.

**Decision**: Window positioning functions (`move()`, `setPosition()`, `center()`) are documented NO-OPs on GTK4/Wayland.

**Rationale**: This is a fundamental Wayland design decision, not a limitation we can work around. Users need to be aware of this behavioral difference.

### Decision 4: Menu System Architecture (2026-01-04)
**Context**: GTK4 removes GtkMenu/GtkMenuItem in favor of GMenu/GAction.

**Decision**: Complete rewrite of menu system for GTK4 using GMenu/GAction/GtkPopoverMenuBar.

**Status**: Stub implementations only. Full implementation pending.

### Decision 5: System Tray Compatibility (2026-01-04)
**Context**: v3's system tray uses D-Bus StatusNotifierItem protocol.

**Decision**: No changes needed - system tray is already GTK-agnostic.

## Implementation Progress

### Phase 1: Build Infrastructure ✅ COMPLETE

**Commit**: `a0ca13fdc` (2026-01-04)

#### 1.1 Add gtk3 constraint to existing files
Files modified:
- `v3/pkg/application/application_linux.go` - Added `gtk3` constraint
- `v3/pkg/application/linux_cgo.go` - Added `gtk3` constraint  
- `v3/internal/assetserver/webview/request_linux.go` - Added `gtk3` constraint
- `v3/internal/assetserver/webview/responsewriter_linux.go` - Added `gtk3` constraint
- `v3/internal/assetserver/webview/webkit2.go` - Added `gtk3` constraint

#### 1.2 Create GTK4 stub files
Files created:
- `v3/pkg/application/linux_cgo_gtk4.go` (~1000 lines)
  - Main CGO file with GTK4 bindings
  - Implements: window management, clipboard, basic menu stubs
  - Uses `gtk4 webkitgtk-6.0` pkg-config
  
- `v3/pkg/application/application_linux_gtk4.go` (~250 lines)
  - Application lifecycle management
  - System theme detection via D-Bus
  - NVIDIA DMA-BUF workaround for Wayland

#### 1.3 Create WebKitGTK 6.0 asset server stubs
Files created:
- `v3/internal/assetserver/webview/webkit6.go`
- `v3/internal/assetserver/webview/request_linux_gtk4.go`
- `v3/internal/assetserver/webview/responsewriter_linux_gtk4.go`

### Phase 2: Doctor & Capabilities 🔄 IN PROGRESS

**Goal**: Update `wails doctor` to check for GTK4 as primary, GTK3 as secondary.

TODO:
- [ ] Update `v3/internal/doctor/doctor_linux.go`
- [ ] Update `v3/internal/capabilities/capabilities_linux.go`
- [ ] Add `wails3 capabilities` command output for GTK version detection

### Phase 3: Window Management 📋 PENDING

TODO:
- [ ] Implement GTK4 signal handlers (different event model)
- [ ] Implement window state management (fullscreen, maximize, minimize)
- [ ] Implement GTK4 drag-and-drop with GtkDropTarget
- [ ] Test window lifecycle on GTK4

### Phase 4: Menu System 📋 PENDING

TODO:
- [ ] Implement GMenu/GAction menu system
- [ ] Implement GtkPopoverMenuBar for application menus
- [ ] Implement context menus with GtkPopoverMenu
- [ ] Handle menu item states (checked, disabled)
- [ ] Implement accelerators with GtkShortcut

### Phase 5: Asset Server 📋 PENDING

TODO:
- [ ] Verify WebKitGTK 6.0 URI scheme handler API
- [ ] Test asset loading
- [ ] Verify JavaScript execution API changes

### Phase 6: Docker & Build System 📋 PENDING

TODO:
- [ ] Update Docker container with both GTK3 and GTK4 libraries
- [ ] Add Taskfile targets: `build:linux` (GTK4), `build:linux:gtk3` (legacy)
- [ ] Update CI/CD workflows

### Phase 7: Testing 📋 PENDING

TODO:
- [ ] Test on Ubuntu 24.04 (native GTK4)
- [ ] Test on Ubuntu 22.04 (backported WebKitGTK 6.0)
- [ ] Test legacy build on older systems
- [ ] Performance benchmarks

## API Differences: GTK3 vs GTK4

| Feature | GTK3 | GTK4 |
|---------|------|------|
| Init | `gtk_init(&argc, &argv)` | `gtk_init_check()` |
| Container | `gtk_container_add()` | `gtk_window_set_child()` |
| Show | `gtk_widget_show_all()` | Widgets visible by default |
| Hide | `gtk_widget_hide()` | `gtk_widget_set_visible(w, FALSE)` |
| Clipboard | `GtkClipboard` | `GdkClipboard` |
| Menu | `GtkMenu/GtkMenuItem` | `GMenu/GAction` |
| Menu Bar | `GtkMenuBar` | `GtkPopoverMenuBar` |
| Window Move | `gtk_window_move()` | NO-OP on Wayland |
| Window Position | `gtk_window_get_position()` | Not available on Wayland |
| Destroy | `gtk_widget_destroy()` | `gtk_window_destroy()` |
| Drag Start | `gtk_window_begin_move_drag()` | `gtk_native_get_surface()` + surface drag |

## Files Reference

### GTK3 (Legacy) Files
```
v3/pkg/application/
  linux_cgo.go              # Main CGO (gtk3 tag)
  application_linux.go       # App lifecycle (gtk3 tag)

v3/internal/assetserver/webview/
  webkit2.go                 # WebKit2GTK helpers (gtk3 tag)
  request_linux.go           # Request handling (gtk3 tag)
  responsewriter_linux.go    # Response writing (gtk3 tag)
```

### GTK4 (Default) Files
```
v3/pkg/application/
  linux_cgo_gtk4.go          # Main CGO (!gtk3 tag)
  application_linux_gtk4.go   # App lifecycle (!gtk3 tag)

v3/internal/assetserver/webview/
  webkit6.go                 # WebKitGTK 6.0 helpers (!gtk3 tag)
  request_linux_gtk4.go      # Request handling (!gtk3 tag)
  responsewriter_linux_gtk4.go # Response writing (!gtk3 tag)
```

### Shared Files (no GTK-specific code)
```
v3/pkg/application/
  webview_window_linux.go    # Window wrapper (uses methods from CGO files)
  systemtray_linux.go        # D-Bus based, GTK-agnostic
  
v3/internal/assetserver/webview/
  request.go                 # Interface definitions
  responsewriter.go          # Interface definitions
```

## Changelog

### 2026-01-04
- Initial implementation of GTK4 build infrastructure
- Added `gtk3` constraint to 5 existing files
- Created 5 new GTK4 stub files
- Updated UNRELEASED_CHANGELOG.md
