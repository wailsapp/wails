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

### Phase 2: Doctor & Capabilities ✅ COMPLETE

**Goal**: Update `wails doctor` to check for GTK4 as primary, GTK3 as secondary.

#### 2.1 Package Manager Updates
All 7 package managers updated to check GTK4/WebKitGTK 6.0 as primary, GTK3 as optional/legacy:
- `v3/internal/doctor/packagemanager/apt.go` ✅
- `v3/internal/doctor/packagemanager/dnf.go` ✅
- `v3/internal/doctor/packagemanager/pacman.go` ✅
- `v3/internal/doctor/packagemanager/zypper.go` ✅
- `v3/internal/doctor/packagemanager/emerge.go` ✅
- `v3/internal/doctor/packagemanager/eopkg.go` ✅
- `v3/internal/doctor/packagemanager/nixpkgs.go` ✅

Package key naming convention: `gtk4`, `webkitgtk-6.0` (primary), `gtk3 (legacy)`, `webkit2gtk (legacy)` (optional)

#### 2.2 Capabilities Detection
Files created/updated:
- `v3/internal/capabilities/capabilities.go` - Added `GTKVersion` (int) and `WebKitVersion` (string) fields
- `v3/internal/capabilities/capabilities_linux.go` - GTK4 default: `GTKVersion: 4, WebKitVersion: "6.0"`
- `v3/internal/capabilities/capabilities_linux_gtk3.go` - GTK3 legacy: `GTKVersion: 3, WebKitVersion: "4.1"`

TODO (deferred to Phase 3):
- [ ] Update `v3/internal/doctor/doctor_linux.go` - Improve output to show GTK4 vs GTK3 status

### Phase 3: Window Management ✅ COMPLETE

#### 3.1 GTK4 Event Controllers
GTK4 replaces direct signal handlers with `GtkEventController` objects:
- `GtkEventControllerFocus` for focus in/out events
- `GtkGestureClick` for button press/release events
- `GtkEventControllerKey` for keyboard events
- Window signals: `close-request`, `notify::maximized`, `notify::fullscreened`

New C function `setupWindowEventControllers()` sets up all event controllers.

#### 3.2 Window Drag and Resize
GTK4 uses `GdkToplevel` API instead of GTK3's `gtk_window_begin_move_drag`:
- `gdk_toplevel_begin_move()` for window drag
- `gdk_toplevel_begin_resize()` for window resize
- Requires `gtk_native_get_surface()` to get the GdkSurface

#### 3.3 Drag-and-Drop with GtkDropTarget
Complete implementation using GTK4's `GtkDropTarget`:
- `on_drop_enter` / `on_drop_leave` for drag enter/exit events
- `on_drop_motion` for drag position updates
- `on_drop` handles file drops via `GDK_TYPE_FILE_LIST`
- Go callbacks: `onDropEnter`, `onDropLeave`, `onDropMotion`, `onDropFiles`

#### 3.4 Window State Detection
- `isMinimised()` uses `gdk_toplevel_get_state()` with `GDK_TOPLEVEL_STATE_MINIMIZED`
- `isMaximised()` uses `gtk_window_is_maximized()`
- `isFullscreen()` uses `gtk_window_is_fullscreen()`

#### 3.5 Size Constraints
GTK4 removed `gtk_window_set_geometry_hints()`. Now using `gtk_widget_set_size_request()` for minimum size.

TODO (deferred):
- [ ] Test window lifecycle on GTK4 with actual GTK4 libraries

### Phase 4: Menu System ✅ COMPLETE

GTK4 completely replaced the menu system. GTK3's GtkMenu/GtkMenuItem are gone.

#### 4.1 GMenu/GAction Architecture
- `GMenu` - Menu model (data structure, not a widget)
- `GMenuItem` - Individual menu item in the model
- `GSimpleAction` - Action that gets triggered when menu item is activated
- `GSimpleActionGroup` - Container for actions, attached to widgets

#### 4.2 Menu Bar Implementation
- `GtkPopoverMenuBar` created from `GMenu` model via `create_menu_bar_from_model()`
- Action group attached to window with `attach_action_group_to_widget()`
- Actions use "app.action_name" namespace

#### 4.3 New Files Created
- `v3/pkg/application/menu_linux_gtk4.go` - GTK4 menu processing
- `v3/pkg/application/menuitem_linux_gtk4.go` - GTK4 menu item handling

#### 4.4 Build Tag Changes
- `menu_linux.go` - Added `gtk3` tag
- `menuitem_linux.go` - Added `gtk3` tag

#### 4.5 Key Functions
- `menuActionActivated()` - Callback when GAction is triggered
- `menuItemNewWithId()` - Creates GMenuItem + associated GSimpleAction
- `menuCheckItemNewWithId()` - Creates stateful toggle action
- `menuRadioItemNewWithId()` - Creates radio action
- `set_action_enabled()` / `set_action_state()` - Manage action state

TODO (deferred):
- [ ] Context menus with GtkPopoverMenu
- [ ] Keyboard accelerators with GtkShortcut

### Phase 5: Asset Server ✅ COMPLETE

WebKitGTK 6.0 uses the same URI scheme handler API as WebKitGTK 4.1.
The asset server implementation is identical between GTK3 and GTK4.

#### 5.1 Asset Server Files (already created in Phase 1)
- `v3/internal/assetserver/webview/webkit6.go` - WebKitGTK 6.0 helpers
- `v3/internal/assetserver/webview/request_linux_gtk4.go` - Request handling
- `v3/internal/assetserver/webview/responsewriter_linux_gtk4.go` - Response writing

#### 5.2 Missing Exports Added
The GTK4 CGO file was missing two critical exports that were in the GTK3 file:
- `onProcessRequest` - Handles URI scheme requests from WebKit
- `sendMessageToBackend` - Handles JavaScript to Go communication

Both exports were added to `linux_cgo_gtk4.go`.

#### 5.3 Key Differences from GTK3
| Aspect | GTK3 | GTK4 |
|--------|------|------|
| pkg-config | `webkit2gtk-4.1` | `webkitgtk-6.0` |
| Headers | `webkit2/webkit2.h` | `webkit/webkit.h` |
| Min version | 2.40 | 6.0 |
| URI scheme API | Same | Same |

TODO (deferred to testing phase):
- [ ] Test asset loading on actual GTK4 system
- [ ] Verify JavaScript execution works correctly

### Phase 6: Docker & Build System ✅ COMPLETE

#### 6.1 Docker Container Updates
Updated both Dockerfile.linux-x86_64 and Dockerfile.linux-arm64 to install:
- GTK4 + WebKitGTK 6.0 (default build target)
- GTK3 + WebKit2GTK 4.1 (for legacy `-tags gtk3` builds)

Build scripts now support `BUILD_TAGS` environment variable:
- Default: Builds with GTK4/WebKitGTK 6.0
- `BUILD_TAGS=gtk3`: Builds with GTK3/WebKit2GTK 4.1

#### 6.2 Taskfile Targets
New targets added to `v3/Taskfile.yaml`:

| Target | Description |
|--------|-------------|
| `test:example:linux` | Build single example with GTK4 (native) |
| `test:example:linux:gtk3` | Build single example with GTK3 (native, legacy) |
| `test:examples:linux:docker:x86_64` | Build all examples with GTK4 in Docker |
| `test:examples:linux:docker:x86_64:gtk3` | Build all examples with GTK3 in Docker |
| `test:examples:linux:docker:arm64` | Build all examples with GTK4 in Docker (ARM64) |
| `test:examples:linux:docker:arm64:gtk3` | Build all examples with GTK3 in Docker (ARM64) |

TODO (deferred):
- [ ] Update CI/CD workflows to test both GTK versions

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

v3/internal/capabilities/
  capabilities_linux_gtk3.go # GTK3 capabilities (gtk3 tag)
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

v3/internal/capabilities/
  capabilities_linux.go      # GTK4 capabilities (!gtk3 tag)
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

### 2026-01-04 (Session 5 continued)
- Completed Phase 6: Docker & Build System
- Updated Dockerfile.linux-x86_64 and Dockerfile.linux-arm64 for GTK4 + GTK3
- Added BUILD_TAGS environment variable support in build scripts
- Added Taskfile targets for GTK4 (default) and GTK3 (legacy) builds

### 2026-01-04 (Session 5)
- Completed Phase 5: Asset Server
- Verified WebKitGTK 6.0 uses same URI scheme handler API as WebKitGTK 4.1
- Added missing `onProcessRequest` export to linux_cgo_gtk4.go
- Added missing `sendMessageToBackend` export to linux_cgo_gtk4.go
- Confirmed asset server files (webkit6.go, request/responsewriter) are complete

### 2026-01-04 (Session 4)
- Completed Phase 4: Menu System
- Implemented GMenu/GAction architecture for GTK4 menus
- Created GtkPopoverMenuBar integration
- Added menu_linux_gtk4.go and menuitem_linux_gtk4.go
- Added gtk3 build tags to original menu files
- Implemented stateful actions for checkboxes and radio items

### 2026-01-04 (Session 3)
- Completed Phase 3: Window Management
- Implemented GTK4 event controllers (GtkEventControllerFocus, GtkGestureClick, GtkEventControllerKey)
- Implemented window drag using GdkToplevel API (gdk_toplevel_begin_move/resize)
- Implemented complete drag-and-drop with GtkDropTarget
- Fixed window state detection (isMinimised, isMaximised, isFullscreen)
- Fixed size() function to properly return window dimensions
- Updated windowSetGeometryHints for GTK4 (uses gtk_widget_set_size_request)

### 2026-01-04 (Session 2)
- Completed Phase 2: Doctor & Capabilities
- Updated all 7 package managers for GTK4/WebKitGTK 6.0 as primary
- Added GTKVersion and WebKitVersion fields to Capabilities struct
- Created capabilities_linux_gtk3.go for legacy build path

### 2026-01-04 (Session 1)
- Initial implementation of GTK4 build infrastructure
- Added `gtk3` constraint to 5 existing files
- Created 5 new GTK4 stub files
- Updated UNRELEASED_CHANGELOG.md
