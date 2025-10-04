# Wails v3 Developer Guide (Top-Level Overview)

This guide provides a technical, developer-focused overview of the Wails v3 codebase located under the `v3` directory. It is the start of a comprehensive series that will later dive into each package and file in detail. For now, it establishes orientation, architecture, responsibilities, and workflows so contributors have a reliable high-level map before we document every component thoroughly.

Scope: Only the `v3` directory is in scope for this guide.

Last updated: 2025-09-24


## Goals of Wails v3

Wails lets you build desktop applications using Go for the backend and web technologies (HTML/CSS/JS) for the frontend, packaging them into a native desktop experience across platforms. Version 3 focuses on modular internal packages, improved runtime and tooling, and clearer boundaries between build-time and run-time responsibilities.


## High-Level Architecture

- Command-line tooling (CLI): developer-facing entry points for building, generating, packaging, and inspecting Wails apps.
- Internal libraries: build orchestration, runtime integration, OS abstractions, packaging, capabilities, templates, and utilities. These are mostly under `v3/internal` and not intended as public API.
- Public Go packages: simple APIs for applications to use (e.g., `v3/pkg/application`, `v3/pkg/ui`, `v3/pkg/services`, etc.).
- Runtime (desktop): the bridge between the Go app and the web frontend, with bindings, window control, events, and platform integrations.
- Examples and tests: reference implementations and validation.


## Directory Map (top-level under v3)

- `cmd/` — CLI entrypoints.
  - `wails3/` — The main Wails v3 CLI. Includes `main.go`, usage docs, and a prebuilt binary (when present in repo artifacts).

- `examples/` — Working sample apps demonstrating various capabilities.
  - Examples include, among others: drag-and-drop, raw message handling, etc. These show best practices and can be used to validate changes.

- `internal/` — Internal (non-public) packages used by the CLI and by build/runtime components.
  - `assetserver/` — Static file serving logic for development or build steps.
  - `buildinfo/` — Build metadata and version stamping.
  - `capabilities/` — Capabilities gating or feature flags across subsystems.
  - `changelog/` — Utilities for changelog parsing/validation.
  - `commands/` — Implementation of CLI subcommands and orchestration logic.
  - `dbus/` — D-Bus integration for Linux desktop features.
  - `debug/` — Debug utilities or dev-mode helpers.
  - `doctor/` — Environment diagnostics (e.g., checking toolchains).
  - `fileexplorer/` — File chooser/dialog logic beyond platform primitives.
  - `flags/` — Centralized command-line flag definitions and parsing helpers.
  - `generator/` — Code and project generators (scaffolding, bindings, etc.).
  - `github/` — GitHub integration (release automation or metadata helpers).
  - `go-common-file-dialog/` — Common file dialog wrappers (likely cross-platform helpers).
  - `hash/` — Hashing utilities for caching, integrity, or content addressing.
  - `operatingsystem/` — OS detection and OS-specific helpers.
  - `packager/` — Packaging/build pipelines to produce platform-specific bundles.
  - `runtime/` — Runtime glue for the desktop environment.
    - `desktop/` — Desktop runtime resources including the browser window bridge and TS/JS runtime.
      - `@wailsio/runtime/` — TypeScript source for the frontend runtime bridge (window and app bindings, event bus, etc.).
  - `s/` — Small internal helpers or shared primitives.
  - `service/` — Service abstractions and service lifecycle helpers used internally.
  - `signal/` — Signal handling (OS process signals).
  - `templates/` — Project and code templates (used by the generator and packager).
  - `term/`, `term2/` — Terminal rendering, progress, and rich text helpers.
  - `version/` — Versioning helpers shared across commands.

- `pkg/` — Public packages intended for application developers.
  - `application/` — Main application abstraction (app lifecycle, window management hooks, startup/shutdown).
  - `events/` — Event bus/interfaces for app<->frontend communication.
  - `icons/` — Icon helpers/types for application and window icons.
  - `mac/` — macOS-specific helpers.
  - `mcp/` — Platform or protocol utilities (exact scope to be detailed in the deep dive).
  - `services/` — Public service interfaces and service registration.
  - `ui/` — UI helpers and types exposed to apps.
  - `w32/` — Windows-specific helpers and bindings.

- `scripts/` — Maintenance or developer scripts used in CI or local workflows.
  - `validate-changelog.go` — Validates changelog format and rules.

- `tasks/` — Task-based workflow helpers (used by maintainers and CI).
  - `cleanup/`, `contribs/`, `events/`, `fix-bindings/`, `release/`, `sed/` — Per-task utilities and scripts.

- `test/`, `tests/`, `tooltest/` — Test suites, fixtures, and tools used for verifying functionality.

- `.task/`, `wep/` — Project tooling/configuration directories (to be detailed later as needed).


## Build and Development Overview

- CLI usage: The `v3/cmd/wails3` command is the primary interface. Typical workflows include:
  - Project creation via templates (generator).
  - Development server and live reload for the frontend (assetserver, runtime dev helpers).
  - Building a production bundle and packaging for the target OS (packager).
- Internal packages are not part of the public API contract and may change. Public Go APIs under `v3/pkg/` are maintained for application developers.
- The desktop runtime bridges Go and the Web frontend via a JS/TS runtime (`internal/runtime/desktop/@wailsio/runtime`). This provides window controls, eventing, and invocation bridges to Go.


## Key Data Flows (Conceptual)

1. Developer runs a CLI command (wails3) → CLI delegates to `internal/commands` and uses helpers under `internal/*` for build, generate, package, etc.
2. The application code (using `v3/pkg/*`) defines the Go-side application, services, and event handlers.
3. At runtime, the desktop bridge injects a JS runtime into the web app so frontend can call Go methods and subscribe to events. The Go-side emits events to the frontend via the runtime channel.
4. Packaging uses `internal/packager` to produce OS-native bundles, using OS-specific helpers from `pkg/mac`, `pkg/w32`, etc., as needed.


## Contributing Workflow (High-Level)

- Make small, focused changes. For cross-platform features, ensure Linux, macOS, and Windows implications are considered.
- Run example apps under `v3/examples` to validate changes.
- Use scripts under `v3/scripts` and tasks under `v3/tasks` where applicable (e.g., changelog validation or release workflows).
- When modifying the runtime bridge, ensure TypeScript builds for `@wailsio/runtime` are up-to-date and integration-tested against example apps.


## Next Steps in This Guide

This is the top-level overview. In subsequent iterations, we will deep-dive into every package and file in the `v3` directory with:
- Detailed per-package responsibilities and structure.
- File-by-file documentation including key types, functions, and data flow.
- Build, test, and troubleshooting instructions for each component.

If you need a specific package documented next, please indicate which one, and we’ll start the in-depth coverage there.


# Deep Dive: pkg/application

Scope: This section documents every file under v3/pkg/application and explains the end-to-end data flows for how an application is initialised and then run. It complements the top-level overview above.

Last reviewed: 2025-09-24


Overview of the application package
- Purpose: Provides the public API and core runtime glue for building and running a Wails desktop application. It owns the App type and orchestrates windows, menus, dialogs, system trays, eventing, and platform-specific main loops.
- Key types:
  - App: The core application object exposed to developers.
  - Window/WebviewWindow: Represents a browser-hosted UI Window (WebView2 on Windows, WKWebView on macOS, WebKitGTK on Linux).
  - Managers: Subsystems that manage windows, menus, dialogs, events, clipboard, screens, environment, system trays, etc.
  - Options: Configuration used to create the App and its initial windows and services.
  - EventProcessor and event contexts: Abstractions for custom events and application/window event flows.


Lifecycle and data flow (initialisation -> run -> shutdown)
1) Construction via New(options)
   - File: application.go (New)
   - Actions:
     - Merge defaults (mergeApplicationDefaults) and normalize Options (application_options.go).
     - Construct the App: allocate fields, set logger, prepare asset serving (middleware, asset handlers), store any initial services from options.
     - Prepare developer tooling (dev vs production variants in application_dev.go/application_production.go). In dev, static server may proxy to dev server; in production, use bundled assets.
     - Defer platform-specific initialization; platform-specific App (platformApp) is created later in Run() via newPlatformApp(...) implemented per-OS.

2) App.init()
   - File: application.go (init method on App)
   - Actions:
     - Create root context (a.ctx) and cancellation function.
     - Instantiate and wire all managers:
       - newWindowManager, newContextMenuManager, newKeyBindingManager, newBrowserManager, newEnvironmentManager, newDialogManager, newEventManager, newMenuManager, newScreenManager, newClipboardManager, newSystemTrayManager.
     - Initialize internal maps: windows, systemTrays, contextMenus, keyBindings, listener registries.

3) Pre-run hooks and service startup
   - File: application.go (Run -> preRun -> service startup)
   - Actions:
     - preRun() executes any post-creation hooks (e.g., build/runtime settling, menu prep, assetserver finalization). If this fails, Run aborts.
     - Services provided in Options.Services are started in order (startupService). On failure, an error is returned and started services are tracked for orderly shutdown.

4) Platform App creation
   - File: application.go (Run)
   - Calls newPlatformApp(a) which selects implementation based on GOOS:
     - application_windows.go -> windowsApp
     - application_darwin.go -> darwinApp (plus associated Objective-C glue .m/.h for Cocoa/WKWebView)
     - application_linux.go -> gtkApp (with CGO or purego variants)
   - The platformApp implements methods like run(), init(), setApplicationMenu(), setIcon(), show/hide, main-thread dispatch, window/tray registration, etc.

5) Event pump setup (concurrency and channels)
   - File: application.go (Run)
   - Starts goroutines to drain and dispatch internal channels:
     - applicationEvents -> a.Event.handleApplicationEvent
     - windowEvents -> a.handleWindowEvent
     - webviewRequests -> a.handleWebViewRequest
     - windowMessageBuffer -> a.handleWindowMessage
     - windowKeyEvents -> a.handleWindowKeyEvent
     - windowDragAndDropBuffer -> a.handleDragAndDropMessage
     - menuItemClicked -> a.Menu.handleMenuItemClicked
   - These goroutines run for the lifetime of the app and route messages from the platform layer and WebView bridges into App managers and windows.

6) Final run
   - File: application.go (Run)
   - Marks a.running = true, drains any a.pendingRun tasks scheduled before the run loop started, sets the application menu (darwin) and icon (if provided), and then calls a.impl.run() which enters the platform event loop (message pump on Windows, CFRunLoop on macOS, GTK main loop on Linux).

7) Shutdown
   - File: application.go (cleanup, shutdownServices, Quit)
   - Cancel root context, tear down windows/system trays/menus, stop services in reverse order, and exit the platform loop. Fatal errors use handleFatalError to log and os.Exit(1).


File-by-file guide (v3/pkg/application)
Note: Brief purpose for each file and how it participates in init/run and data flows.

Core App and options
- application.go: Defines App and its lifecycle. Key functions/methods: New, init (method), Run, cleanup, Quit, event and message handlers, Hide/Show, SetIcon, runOrDeferToAppRun. Creates managers and coordinates inter-manager communication via channels and callbacks.
- application_options.go: Defines Options for configuring the App: window defaults, asset options, logger, middlewares, platform-specific options (MacOptions, WindowsOptions, LinuxOptions), and asset server helpers. Contains middleware chaining and helpers for serving assets from fs.FS.
- services.go: Declares Service interface and service lifecycle integration with App (startupService, shutdownServices) and registration via App.RegisterService.
- environment.go / environment_manager.go: Detects and exposes environment info to the app (dev/prod, capabilities, platform variables), often referenced when preparing runtime flags and window defaults.
- errors.go: Error types, including FatalError used by handleFatalError in App.
- panic_handler.go: Centralized panic recovery where goroutines defer handlePanic() to avoid crashing the app silently.
- mainthread*.go: Utilities to ensure code runs on the GUI main thread when required by the OS toolkit; used by platformApp implementations and App.dispatchOnMainThread.
- messageprocessor*.go: Message router between WebView frontend and Go backend; different files focus on different domains (application, window, dialogs, clipboard, context menus, screens, system). They parse inbound messages (e.g., from JS runtime) and call into managers/App methods.
- bindings.go: Handles Go<->JS binding registration and call routing, including returning values/errors to the frontend via WebView.

Platform-specific app implementations
- application_windows.go: windowsApp implements platformApp for Windows. Responsibilities:
  - init(): OS integration (DPI awareness, COM init as needed via WebView2 loader), building application menu, creating hidden “application” window for message routing if needed.
  - run(): Enters the Windows message loop, dispatches messages to wndProc, and integrates with WebView2 (go-webview2).
  - wndProc(): Translates native window messages into Wails events: focus, resize, menu select, accelerator keys, system tray clicks, drag/drop, etc., forwarding into channels consumed by App.Run goroutines.
  - setApplicationMenu(), setIcon(), show/hide, registration of windows and system trays.
  - logPlatformInfo() and platformEnvironment() provide environment metadata used during startup.
- application_darwin.go and related Objective-C files (.h/.m): darwinApp uses Cocoa and WKWebView. Objective-C delegate files integrate with NSApplication, NSWindow, menus, and app activation policy; bridge callbacks enqueue events into Go channels.
- application_linux.go plus linux_cgo.go/linux_purego.go: GTK-based implementation; CGO variant integrates with GTK main loop and dialogs; purego variant offers a fallback depending on build tags.
- application_dev.go / application_production.go: Build-tagged files that adjust behavior for development vs production (e.g., logger defaults, asset server wiring, runtime flags). logger_dev*.go/logger_prod.go customize logging sinks by target OS and mode.

Windows/macOS/Linux specific helpers
- keys*.go: Keyboard key codes and modifier mappings per OS; used to parse accelerators and handle key bindings.
- dialogs*.go: Platform-specific dialog implementations (message box, file open/save) backing the cross-platform dialogs.go facade.
- clipboard*.go and clipboard_manager.go: Clipboard APIs per OS.
- screen_*.go and screenmanager.go: Screen enumeration, DPI/scale factors, active screen selection; used when positioning windows.
- systemtray_*.go and system_tray_manager.go: System tray icon/menu management; events funneled through channels to App and Menu managers.

Windowing: webview window and options
- webview_window.go: The core window abstraction backed by a webviewWindowImpl (implemented per-OS). Responsibilities:
  - Window creation via NewWindow(options) after App is created.
  - Eventing: WindowEvent and listener registration (OnWindowEvent, RegisterHook), mapping of events to callbacks.
  - Lifecycle methods: Show/Hide, Run (window-level), Close/Destroy, Focus, Resize/Move, Fullscreen/Maximise/Minimise, Zoom, ContentProtection, Frameless toggling, DevTools (per OS), ExecJS, SetHTML/URL, printing.
  - Input handling: processKeyBinding and HandleKeyEvent for accelerators; integrates with menu accelerator mappings.
  - Messaging bridge: HandleMessage (routes call responses and RPC from JS), DialogResponse/CallResponse, drag-and-drop (HandleDragAndDropMessage), context menus.
  - Internal queues: dispatchWindowEvent, isDestroyed checks, destroyed state management.
- webview_window_options.go: Defines WebviewWindowOptions, styling and behavior (title, size, min/max constraints, position, background color, frameless, URL/HTML/Assets, debug flags, dev tools, zoom policy, theme, user agent, drag regions, etc.). Also includes GPU policy and per-OS defaults. Options feed directly into platform window creation.
- webview_window_* per-OS files: Implement webviewWindowImpl for each platform. They translate cross-platform window API to native toolkit calls and send back events via channels (e.g., windowEvents, windowMessageBuffer, windowKeyEvents).
- window.go and window_manager.go: The common Window interface and its manager. The manager tracks windows by ID, exposes Find/All/Focused, and routes broadcast events.

Menus, roles, and accelerators
- menu.go/menu_darwin.go/menu_linux.go/menu_windows.go: Cross-platform Menu representation with platform renderers. Supports app menu (macOS), window menus, and menu bar visibility.
- menuitem.go and family: MenuItem definition, roles, selectors (macOS), per-OS renderers. menuitem_roles.go defines standard roles (Copy, Paste, Quit, etc.).
- roles*.go: Dev vs production role sets.
- popupmenu_windows.go and context_menu_manager.go: Context menu helpers and manager lifecycle.
- key_binding_manager.go: Manages registration of accelerators (e.g., "CmdOrCtrl+Shift+P") and dispatches to window/menu handlers.

Dialogs
- dialogs.go: Cross-platform API for dialogs. Provides MessageDialog, OpenFileDialog, SaveFileDialog with builder-style options; delegates to per-OS implementations via messageDialogImpl/openFileDialogImpl/saveFileDialogImpl.
- dialogs_* per-OS: Bridge to native dialog toolkits (NSAlert/NSSavePanel on macOS, IFileDialog/MessageBox on Windows, GTK dialogs on Linux).

Eventing
- events.go: Custom event system for the application package. Provides:
  - CustomEvent (name, data, cancellation) and ApplicationEvent/WindowEvent wrappers with Context() access to typed event contexts.
  - EventProcessor to register listeners (On/Once/OnMultiple), hooks, and to Emit events. Safe for concurrent use via sync primitives.
  - Bridges to pkg/events for enum types (ApplicationEventType, WindowEventType) used throughout.
- context_application_event.go, context_window_event.go: Define typed context payloads for app/window events passed to listeners (e.g., theme change info, file open URL, window IDs).
- events_common_*.go and event_manager.go: Glue code that maps platform events to high-level events and exposes App.Event API used by apps and internal subsystems.

Miscellaneous utilities
- browser_manager.go: Tracks browser/webview state across windows; used by certain operations that need to coordinate with the web runtime.
- image.go: Image helpers (e.g., RGBA), icons.
- path.go: Small helpers for path normalization when interacting with OS file dialogs or asset paths.
- urlvalidator.go: Validates URLs provided for window content (security and correctness). Unit tests in urlvalidator_test.go.
- single_instance*.go: Optional single-instance enforcement per OS; integrates with platform IPC or file locks.
- systemtray.go: Cross-platform SystemTray API given to app developers; platform-specific implementations in systemtray_*.go.
- clipboard.go: Cross-platform Clipboard API facade; platform-specific glue in clipboard_*.go.
- logger_dev*.go/logger_prod.go: Logging helpers per build mode and OS.
- TODO.md: Notes for maintainers.


Detailed init/run sequence with call graph highlights
- App creation:
  1. application.New(options)
     - Merges defaults; stores Options; sets up assets/middleware stack; chooses dev/prod behaviors; prepares logger.
     - Does NOT start platform loop yet.
  2. App.init() (method) is called by New internally to initialise managers and state.

- Run startup (application.go: Run):
  3. Guard against double Run; mark starting; arrange context cancellation on failure.
  4. preRun(): perform any post-construction checks and assetserver finalisation.
  5. Create platform app: a.impl = newPlatformApp(a) -> application_{os}.go
     - Platform impl may do platform-specific pre-main-loop init in impl.init().
  6. Start services: for each options.Service, a.startupService(service) and record for shutdown order.
  7. Start event pumps: spin goroutines to drain internal channels and dispatch to managers and windows (see channels list above).
  8. Flip a.running = true and execute any pending runnables queued before Run (runOrDeferToAppRun).
  9. Apply app-level UI state: set application menu (macOS), set icon if provided.
 10. Enter platform loop: return a.impl.run(). From here, the native main loop owns the thread and dispatches native events that are bridged back into Go via callbacks and channels.

- Runtime interactions:
  - Window creation: NewWindow(options) constructs a WebviewWindow, invokes platform creation, and registers it with WindowManager and platformApp. Events from native layer (resize, focus, key) become messages in windowEvents/windowMessageBuffer/windowKeyEvents and are dispatched by App handlers to per-window callbacks.
  - Dialogs: The public dialog builders call through to per-OS implementations; responses are sent back via DialogResponse/DialogError into the originating Window.
  - Menus and accelerators: Menu items with accelerators are registered with KeyBindingManager. Native key events are turned into accelerator strings per-OS (keys_*.go) and routed back to WebviewWindow.processKeyBinding.
  - Custom events: App.Event.Emit(CustomEvent) dispatches to registered listeners and also to windows (DispatchWailsEvent) so frontends subscribed via the JS runtime receive them.

- Shutdown:
  - Triggered by App.Quit(), window close that results in shouldQuit() returning true, or fatal errors.
  - App.cleanup() tears down windows, trays, menus, and services in order; cancels root context; releases native resources via platformApp.destroy().


Data channels and queues (core ones referenced in application.go)
- applicationEvents: carries ApplicationEvent to App.Event handler.
- windowEvents: carries internal windowEvent to App.handleWindowEvent which locates the Window and emits a typed WindowEvent.
- webviewRequests: transports webview asset HTTP requests to the asset server pipeline (headers, middleware, handlers).
- windowMessageBuffer: messages from webview to window (RPC responses, console/log, menu clicks, etc.).
- windowKeyEvents: key accelerator strings per window; processed by WebviewWindow to run actions or menu accelerators.
- windowDragAndDropBuffer: drag-and-drop messages containing filenames and drop-zone details; dispatched to the destination window.
- menuItemClicked: menu item selection IDs enqueued by platform menu handlers and dispatched to Menu manager.


How to trace the flow yourself
- Set breakpoints in application.go Run() and platform newPlatformApp(...) to observe startup.
- On Windows, observe application_windows.go run() and wndProc() to see message dispatch into channels.
- Create a simple Example (see v3/examples) and instrument NewWindow() and WebviewWindow.HandleMessage() to see frontend<->backend calls.
- Enable dev logging (logger_dev*.go) to see debug lines emitted when events and drag-and-drop messages are processed.


Practical usage pattern (developer API surface)
- Create and configure the app:
  - app := application.New(application.Options{ /* set windows, assets, menu, services, etc. */ })
  - mainWindow := application.NewWindow(webview_window_options)
  - mainWindow.Show()
  - return app.Run()
- Register services via app.RegisterService(...) before Run().
- Use app.Event.On/Once to register for application events (ApplicationStarted, ThemeChanged, etc.).
- Use window.OnWindowEvent to subscribe to per-window events.
- Use dialogs via application.InfoDialog()/OpenFileDialog()/SaveFileDialog(), etc.

This deep dive should equip you to understand and trace how an application is initialised and run across platforms within the v3/pkg/application package. 