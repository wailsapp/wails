# WebKitGTK 6.0 / GTK4 Implementation Tracker

## Overview

This document tracks the implementation of WebKitGTK 6.0 (GTK4) support for Wails v3 on Linux.

**Goal**: Provide GTK4/WebKitGTK 6.0 support as an EXPERIMENTAL opt-in via `-tags gtk4`, while maintaining GTK3/WebKit2GTK 4.1 as the stable default.

## Architecture Decisions

### Decision 1: GTK3 as Default, GTK4 Opt-In (2026-02-04)
**Context**: Need to support modern Linux distributions with GTK4 while maintaining stability for existing apps.

**Decision**: GTK3 remains the stable default (no build tag required). GTK4 is available as experimental via `-tags gtk4`.

**Rationale**:
- GTK3/WebKit2GTK 4.1 is battle-tested and widely deployed
- GTK4 support needs more community testing before becoming default
- Allows gradual migration and feedback collection
- Protects existing apps from unexpected breakage

**Build Tags**:
- Default (no tag): `//go:build linux && cgo && !gtk4 && !android`
- Experimental GTK4: `//go:build linux && cgo && gtk4 && !android`

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

Package key naming convention: `gtk3`, `webkit2gtk-4.1` (primary/default), `gtk4`, `webkitgtk-6.0` (experimental, optional)

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
- GTK3 + WebKit2GTK 4.1 (default build target)
- GTK4 + WebKitGTK 6.0 (for experimental `-tags gtk4` builds)

Build scripts now support `BUILD_TAGS` environment variable:
- Default: Builds with GTK3/WebKit2GTK 4.1
- `BUILD_TAGS=gtk4`: Builds with GTK4/WebKitGTK 6.0 (experimental)

#### 6.2 Taskfile Targets
New targets added to `v3/Taskfile.yaml`:

| Target | Description |
|--------|-------------|
| `test:example:linux` | Build single example with GTK3 (native, default) |
| `test:example:linux:gtk4` | Build single example with GTK4 (native, experimental) |
| `test:examples:linux:docker:x86_64` | Build all examples with GTK3 in Docker |
| `test:examples:linux:docker:x86_64:gtk4` | Build all examples with GTK4 in Docker (experimental) |
| `test:examples:linux:docker:arm64` | Build all examples with GTK3 in Docker (ARM64) |
| `test:examples:linux:docker:arm64:gtk4` | Build all examples with GTK4 in Docker (ARM64, experimental) |

TODO (deferred):
- [ ] Update CI/CD workflows to test both GTK versions

### Phase 8: Dialog System ✅ COMPLETE

GTK4 completely replaced the dialog APIs. GTK3's `GtkFileChooserDialog` and
`gtk_message_dialog_new` are deprecated/removed.

#### 8.1 File Dialogs
GTK4 uses `GtkFileDialog` with async API:
- `gtk_file_dialog_open()` - Open single file
- `gtk_file_dialog_open_multiple()` - Open multiple files
- `gtk_file_dialog_select_folder()` - Select folder
- `gtk_file_dialog_select_multiple_folders()` - Select multiple folders
- `gtk_file_dialog_save()` - Save file

Key differences:
- No more `gtk_dialog_run()` - everything is async with callbacks
- Filters use `GListStore` of `GtkFileFilter` objects
- Results delivered via `GAsyncResult` callbacks
- Custom button text via `gtk_file_dialog_set_accept_label()`

#### 8.1.1 GTK4 File Dialog Limitations (Portal-based)

GTK4's `GtkFileDialog` uses **xdg-desktop-portal** for native file dialogs. This provides
better desktop integration but removes some application control:

| Feature | GTK3 | GTK4 | Notes |
|---------|------|------|-------|
| `ShowHiddenFiles()` | ✅ Works | ❌ No effect | User controls via portal UI toggle |
| `CanCreateDirectories()` | ✅ Works | ❌ No effect | Always enabled in portal |
| `ResolvesAliases()` | ✅ Works | ❌ No effect | Portal handles symlinks |
| `SetButtonText()` | ✅ Works | ✅ Works | `gtk_file_dialog_set_accept_label()` |
| Multiple folders | ✅ Works | ✅ Works | `gtk_file_dialog_select_multiple_folders()` |

**Why these limitations exist**: GTK4's portal-based dialogs delegate UI control to the
desktop environment (GNOME, KDE, etc.). This is intentional - the portal provides
consistent UX across applications and respects user preferences.

#### 8.2 Message Dialogs
GTK4 uses `GtkAlertDialog`:
- `gtk_alert_dialog_choose()` - Show dialog with buttons
- Buttons specified as NULL-terminated string array
- Default and cancel button indices configurable

#### 8.3 Implementation Details
- Request ID tracking for async callback matching
- `fileDialogCallback` / `alertDialogCallback` C exports for results
- `runChooserDialog()` and `runQuestionDialog()` Go wrappers
- `runOpenFileDialog()` and `runSaveFileDialog()` convenience functions

| GTK3 | GTK4 |
|------|------|
| `GtkFileChooserDialog` | `GtkFileDialog` |
| `gtk_dialog_run()` | Async callbacks |
| `gtk_message_dialog_new()` | `GtkAlertDialog` |
| `gtk_widget_destroy()` | `g_object_unref()` |

### Phase 9: Keyboard Accelerators ✅ COMPLETE

GTK4 uses `gtk_application_set_accels_for_action()` to bind keyboard shortcuts to GActions.

#### 9.1 Key Components

**C Helper Functions** (in `linux_cgo_gtk4.go`):
- `set_action_accelerator(app, action_name, accel)` - Sets accelerator for a GAction
- `build_accelerator_string(key, mods)` - Converts key+modifiers to GTK accelerator string

**Go Functions** (in `linux_cgo_gtk4.go`):
- `namedKeysToGTK` - Map of key names to GDK keysym values (e.g., "backspace" → 0xff08)
- `parseKeyGTK(key)` - Converts Wails key string to GDK keysym
- `parseModifiersGTK(modifiers)` - Converts Wails modifiers to GdkModifierType
- `acceleratorToGTK(accel)` - Converts full accelerator to GTK format
- `setMenuItemAccelerator(itemId, accel)` - Sets accelerator for a menu item

**Integration** (in `menuitem_linux_gtk4.go`):
- `setAccelerator()` method on `linuxMenuItem` calls `setMenuItemAccelerator()`
- `newMenuItemImpl()`, `newCheckMenuItemImpl()`, `newRadioMenuItemImpl()` all set accelerators during creation

#### 9.2 Accelerator String Format

GTK accelerator strings use format like:
- `<Control>q` - Ctrl+Q
- `<Control><Shift>s` - Ctrl+Shift+S
- `<Alt>F4` - Alt+F4
- `<Super>e` - Super+E (Windows/Command key)

#### 9.3 Modifier Mapping

| Wails Modifier | GDK Modifier |
|----------------|--------------|
| `CmdOrCtrlKey` | `GDK_CONTROL_MASK` |
| `ControlKey` | `GDK_CONTROL_MASK` |
| `OptionOrAltKey` | `GDK_ALT_MASK` |
| `ShiftKey` | `GDK_SHIFT_MASK` |
| `SuperKey` | `GDK_SUPER_MASK` |

### Phase 10: Testing 📋 PENDING

TODO:
- [ ] Test on Ubuntu 24.04 (native GTK4)
- [ ] Test on Ubuntu 22.04 (backported WebKitGTK 6.0)
- [ ] Test legacy build on older systems
- [ ] Performance benchmarks
- [ ] Verify file dialogs work correctly
- [ ] Verify message dialogs work correctly

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

### GTK3 (Default) Files
```
v3/pkg/application/
  linux_cgo.go              # Main CGO (!gtk4 tag - default)
  application_linux.go       # App lifecycle (!gtk4 tag - default)

v3/internal/assetserver/webview/
  webkit2.go                 # WebKit2GTK helpers (!gtk4 tag - default)
  request_linux.go           # Request handling (!gtk4 tag - default)
  responsewriter_linux.go    # Response writing (!gtk4 tag - default)

v3/internal/capabilities/
  capabilities_linux_gtk3.go # GTK3 capabilities (!gtk4 tag - default)

v3/internal/operatingsystem/
  webkit_linux.go           # WebKit version info (!gtk4 tag - default)
```

### GTK4 (Experimental) Files
```
v3/pkg/application/
  linux_cgo_gtk4.go          # Main CGO (gtk4 tag - experimental)
  application_linux_gtk4.go   # App lifecycle (gtk4 tag - experimental)

v3/internal/assetserver/webview/
  webkit6.go                 # WebKitGTK 6.0 helpers (gtk4 tag - experimental)
  request_linux_gtk4.go      # Request handling (gtk4 tag - experimental)
  responsewriter_linux_gtk4.go # Response writing (gtk4 tag - experimental)

v3/internal/capabilities/
  capabilities_linux.go      # GTK4 capabilities (gtk4 tag - experimental)

v3/internal/operatingsystem/
  webkit_linux_gtk4.go       # WebKit version info (gtk4 tag - experimental)
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

### 2026-01-07 (Session 11)
- Fixed GTK4 dialog system bugs
- **File Dialog Fix**: Removed premature `g_object_unref()` that freed dialog before async callback
  - GTK4 async dialogs manage their own lifecycle
  - Commit: `6f9c5beb5`
- **Alert Dialog Fixes**:
  - Removed premature `g_object_unref(dialog)` from `show_alert_dialog()` (same issue as file dialogs)
  - Fixed deadlock in `dialogs_linux.go` - `InvokeAsync` → `go func()` since `runQuestionDialog` blocks internally
  - Fixed `runQuestionDialog` to use `options.Title` as message (was using `options.Message`)
  - Added default "OK" button when no buttons specified
  - Commit: `1a77e6091`
- **Other Fixes**:
  - Fixed checkptr errors with `-race` flag by changing C signal functions to accept `uintptr_t` (`3999f1f24`)
  - Fixed ExecJS race condition by adding mutex for `runtimeLoaded`/`pendingJS` (`8e386034e`)
- Added DEBUG_LOG macro for compile-time debug output: `CGO_CFLAGS="-DWAILS_GTK_DEBUG" go build ...`
- Added manual dialog test suite in `v3/test/manual/dialog/`
- **Additional Dialog Fixes** (Session 11 continued):
  - Added `gtk_file_dialog_set_accept_label()` for custom button text
  - Added `gtk_file_dialog_select_multiple_folders()` for multiple directory selection
  - Fixed data race in `application.go` cleanup - was using RLock() when writing `a.windows = nil`
  - Documented GTK4 portal limitations (ShowHiddenFiles, CanCreateDirectories have no effect)
- Files modified:
  - `v3/pkg/application/linux_cgo_gtk4.go` - dialog fixes, race fixes, accept label, multiple folders
  - `v3/pkg/application/linux_cgo_gtk4.c` - DEBUG_LOG macro, alert dialog lifecycle fix, select_multiple_folders callback
  - `v3/pkg/application/linux_cgo_gtk4.h` - uintptr_t for signal functions
  - `v3/pkg/application/dialogs_linux.go` - deadlock fix
  - `v3/pkg/application/webview_window.go` - pendingJS mutex
  - `v3/pkg/application/application.go` - RLock → Lock for cleanup writes
  - `docs/src/content/docs/reference/dialogs.mdx` - documented GTK4 limitations

### 2026-01-04 (Session 10)
- Fixed Window → Zoom menu behavior to toggle maximize/restore (was incorrectly calling webview zoomIn)
- Fixed radio button styling in GTK4 GMenu (now shows dots instead of checkmarks)
  - Implemented proper GMenu radio groups with string-valued stateful actions
  - All items in group share same action name with unique target values
  - Added `create_radio_menu_item()` C helper and `menuRadioItemNewWithGroup()` Go wrapper
- Researched Wayland minimize behavior:
  - `gtk_window_minimize()` works on GNOME/KDE (sends xdg_toplevel_set_minimized)
  - May be no-op on tiling WMs (Sway, etc.) per Wayland protocol design
- Fixed app not terminating when last window closed
  - Added quit logic to `unregisterWindow()` in `application_linux_gtk4.go`
  - Respects `DisableQuitOnLastWindowClosed` option
- Fixed menu separators not showing
  - GMenu uses sections for visual separators (not separate separator items)
  - Rewrote menu processing to group items into sections, separators create new sections
  - Added `menuNewSection()`, `menuAppendSection()`, `menuAppendItemToSection()` helpers
- Added CSS provider to reduce popover menu padding
- Removed all debug println statements
- Files modified:
  - `v3/pkg/application/linux_cgo_gtk4.go` - added radio group support, section helpers
  - `v3/pkg/application/linux_cgo_gtk4.c` - added create_radio_menu_item(), init_menu_css()
  - `v3/pkg/application/linux_cgo_gtk4.h` - added function declaration
  - `v3/pkg/application/application_linux_gtk4.go` - added quit-on-last-window logic
  - `v3/pkg/application/menu_linux_gtk4.go` - section-based menu processing, radio groups
  - `v3/pkg/application/menuitem_linux_gtk4.go` - updated radio item creation
  - `v3/pkg/application/webview_window_linux.go` - fixed zoom() to toggle maximize
  - `v3/pkg/application/window_manager.go` - removed debug output

### 2026-01-04 (Session 9)
- Fixed GTK4 window creation crash (SIGSEGV in gtk_application_window_new)
- **Root Cause**: GTK4 requires app to be "activated" before creating windows
- **Solution**: Added activation synchronization mechanism:
  - Added `activated` channel and `sync.Once` to `linuxApp` struct
  - Added `markActivated()` method called from `activateLinux()` callback
  - Added `waitForActivation()` method for callers to block until ready
  - Modified `WebviewWindow.Run()` to wait for activation before `InvokeSync`
- Files modified:
  - `v3/pkg/application/application_linux_gtk4.go` - activation gate
  - `v3/pkg/application/linux_cgo_gtk4.go` - call markActivated() in activateLinux
  - `v3/pkg/application/webview_window.go` - wait for activation on GTK4
- GTK4 apps now create windows successfully without crashes

### 2026-01-04 (Session 8)
- Fixed GTK3/GTK4 symbol conflict in operatingsystem package
- Added `gtk3` build tag to `v3/internal/operatingsystem/webkit_linux.go`
- Created `v3/internal/operatingsystem/webkit_linux_gtk4.go` with GTK4/WebKitGTK 6.0
- Moved app initialization from `init()` to `newPlatformApp()` for cleaner setup
- Resolved runtime crash: "GTK 2/3 symbols detected in GTK 4 process"
- Verified menu example runs successfully with GTK 4.20.3 and WebKitGTK 2.50.3

### 2026-01-04 (Session 7)
- Completed Phase 9: Keyboard Accelerators
- Added namedKeysToGTK map with GDK keysym values for all special keys
- Added parseKeyGTK() and parseModifiersGTK() conversion functions
- Added acceleratorToGTK() to convert Wails accelerator format to GTK
- Added setMenuItemAccelerator() Go wrapper that calls C helpers
- Integrated accelerator setting in all menu item creation functions
- Uses gtk_application_set_accels_for_action() for GTK4 shortcut binding

### 2026-01-04 (Session 6)
- Completed Phase 8: Dialog System
- Implemented GtkFileDialog for file open/save/folder dialogs
- Implemented GtkAlertDialog for message dialogs
- Added async callback system for GTK4 dialogs (no more gtk_dialog_run)
- Added C helper functions and Go wrapper functions

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
