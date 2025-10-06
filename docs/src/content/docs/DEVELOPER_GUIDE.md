---
title: Wails v3 Developer Guide
description: Comprehensive onboarding for engineering and debugging the Wails v3 codebase
slug: developer-guide
sidebar:
  label: Developer Guide
  order: 50
---

## How to Use This Guide
- **Scope**: Everything documented here applies to the `v3/` branch of Wails, tracking the current workspace state (commit you have checked out) and intentionally ignores the legacy `v2/` runtime.
- **Audience**: Senior Go + desktop developers who need to become productive at fixing bugs or extending the v3 runtime, CLI, or build tooling.
- **Format**: Top-down description → component drill-down → pseudocode for critical flows → teardown guidance → reference tables.
- **Navigation**: Skim the diagrams first, then expand sections relevant to the subsystem you are editing.

```d2
direction: down
Developer: "Developer"
CLI: "CLI Layer\n(cmd/wails3)"
Commands: "Command Handlers\n(internal/commands)"
TaskRunner: "Task Runner\n(internal/commands/task.go)"
Taskfile: "Taskfile\n(v3/Taskfile.yaml)"
BuildPipeline: "Build/Packaging\n(internal/*)"
Templates: "Templates\n(internal/templates)"
RuntimeGen: "Runtime Assets\n(internal/runtime)"
Application: "Application Core\n(pkg/application)"
AssetServer: "Asset Server\n(internal/assetserver)"
MessageBridge: "Message Processor\n(pkg/application/messageprocessor*.go)"
Webview: "WebView Implementations\n(pkg/application/webview_window_*.go)"
Services: "Services & Bindings\n(pkg/services/*)"
Platform: "Platform Layers\n(pkg/mac, pkg/w32, pkg/events, pkg/ui)"

Developer -> CLI: "runs wails3"
CLI -> Commands: "registers subcommands"
Commands -> TaskRunner: "wraps build/package/dev"
TaskRunner -> Taskfile: "executes tasks"
Commands -> BuildPipeline: "invokes packagers & generators"
BuildPipeline -> Templates: "renders scaffolds"
Commands -> RuntimeGen: "builds runtime JS"
RuntimeGen -> Application: "embedded assets"
Application -> AssetServer: "serves HTTP"
Application -> MessageBridge: "routes runtime calls"
MessageBridge -> Webview: "bridge over WebView"
Application -> Services: "binds Go services"
Application -> Platform: "dispatches to OS"
AssetServer -> Webview: "feeds frontend"
Services -> MessageBridge: "expose methods"
```

---

## Repository Layout (v3 only)
| Path | Purpose | Highlights |
| --- | --- | --- |
| `v3/cmd/wails3` | CLI entrypoint | Command registration, build info capture (`main.go`). |
| `v3/internal/commands` | Implementation for each CLI subcommand | Includes build/package/dev tooling, template generators, diagnostics. |
| `v3/internal/flags` | Typed flag definitions for CLI commands | Shared flag structs (e.g. `Build`, `Package`, `ServiceInit`). |
| `v3/internal/assetserver` | Embedded HTTP server powering the runtime | Request routing, middleware, webview handoff. |
| `v3/internal/runtime` | Generates runtime JS payload | Platform-specific flag injection, dev/prod toggles. |
| `v3/internal/templates` | Project/frontend scaffolding | Multiple frontend stacks with shared `_common` assets. |
| `v3/internal/service` | Service template generator | Creates Go service skeletons. |
| `v3/internal/packager` | Wrapper around nfpm for Linux packaging | Generates `.deb`, `.rpm`, `.apk`, etc. |
| `v3/internal/term`, `term2` | Console styling helpers | Shared across commands. |
| `v3/pkg/application` | Core runtime used by user apps | Windowing, events, bindings, asset hosting, shutdown. |
| `v3/pkg/events` | Enumerates runtime event IDs | Shared constants used on Go and JS sides. |
| `v3/pkg/services` | Built-in services (badge, notifications, sqlite, etc.) | Reference implementations for custom services. |
| `v3/pkg/icons`, `v3/pkg/ui` | Icon tooling & UI examples | Support packages used by CLI utilities. |
| `v3/pkg/mac`, `v3/pkg/w32`, `v3/pkg/mcp` | Platform bindings | cgo / pure Go glue for OS integrations. |
| `v3/scripts`, `v3/tasks`, `v3/test`, `v3/tests` | Automation helpers | Task runners, Docker harnesses, regression suites. |
| `v3/Taskfile.yaml` | Root Taskfile consumed by CLI wrappers | Defines build/test workflows referenced by commands. |

---

## CLI and Tooling Layer
### Entry Point: `v3/cmd/wails3/main.go`
The CLI uses `github.com/leaanthony/clir` to declaratively register subcommands. Build metadata is captured in `init()` by reading `debug.ReadBuildInfo` and storing it into `internal/commands.BuildSettings` for later introspection (`wails3 tool buildinfo`).

```pseudo
func main():
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    register simple actions (docs, sponsor) -> open browser
    register functional subcommands that delegate to internal/commands
        init -> commands.Init
        build -> commands.Build (wraps task)
        dev -> commands.Dev (starts watcher)
        package -> commands.Package
        doctor -> commands.Doctor
        releasenotes -> commands.ReleaseNotes
        task -> commands.RunTask (Taskfile wrapper)
        generate subtree -> {build-assets, icons, syso, runtime, template, ...}
        update subtree -> {build-assets, cli}
        service -> {init}
        tool subtree -> {checkport, watcher, cp, buildinfo, package, version}
        version -> commands.Version
    defer printFooter() (prints docs/sponsor hints unless disabled)
    on error -> log with pterm and exit 1
```
Key takeaways:
- All heavy lifting happens inside `internal/commands`. The CLI just constructs flag structs and hands them off.
- Build/package commands are aliases for `wails3 task <name>`; real work is defined in `v3/Taskfile.yaml`.
- `commands.DisableFooter` prevents duplicate footer output when commands open browsers or produce their own footers.

### Flag Definitions: `v3/internal/flags`
- Each subcommand gets a struct tagged with descriptions/defaults used by `clir`.
- Example: `flags.Build` holds options for `wails3 build`, while `flags.GenerateBindingsOptions` describes binding generator inputs.
- Flag structs double as configuration objects passed into command implementations, so the same struct layout must be respected by tests.

### Command Implementations: Highlights
| File | Responsibility |
| --- | --- |
| `internal/commands/init.go` | Scaffolds new projects from templates, resolves template metadata, writes config (`wails.json`). |
| `internal/commands/dev.go` | Sets up the dev server port, populates `FRONTEND_DEVSERVER_URL`, and forwards to the file watcher. |
| `internal/commands/watcher.go` | Loads `taskfile`-compatible YAML (`dev_mode`), instantiates `refresh/engine`, registers signal handlers, and keeps the engine alive until interrupted. |
| `internal/commands/build-assets.go` | Materializes packaging assets from embedded templates (`gosod` extractor). Handles defaults, path normalization, and YAML config ingestion. |
| `internal/commands/generate_template.go` | Creates template stubs by reading `internal/templates`. |
| `internal/commands/tool_*` | Misc utilities: port checks, file copy, semantic version bump (`tool_version.go`), etc. |
| `internal/commands/task_wrapper.go` | Implements `wails3 build/package` alias behavior by re-invoking the `task` subcommand with rewired `os.Args`. |

```pseudo
func wrapTask(command string, otherArgs []string):
    warn("alias for task")
    newArgs := ["wails3", "task", command] + otherArgs
    os.Args = newArgs
    return RunTask(&RunTaskOptions{Name: command}, otherArgs)
```

### Task Runner Integration: `internal/commands/task.go`
- Wraps `github.com/wailsapp/task/v3` for cross-platform task execution.
- Accepts both positional task names (`wails3 task dev -- <extra>`) and `--name` flags.
- Validates mutually exclusive options (e.g., `--dir` vs `--taskfile`).
- Supports list/status/JSON output inherited from upstream Task library.
- Uses `BuildSettings` to print the bundled Task version when `--version` is passed.

### Dev Mode Flow (`wails3 dev`)
1. Resolve Vite port preference (`--port`, `WAILS_VITE_PORT`, default `9245`).
2. Check that the port is free by attempting to `net.Listen` and immediately closing.
3. Export `WAILS_VITE_PORT` and set `FRONTEND_DEVSERVER_URL` (`http` vs `https` based on `--secure`).
4. Invoke `Watcher` with the configured Taskfile path (defaults to `./build/config.yml`).

### Packaging & Distribution Commands
- `GenerateSyso` (Windows resource) writes `.syso` files for icon/resource embedding.
- `GenerateIcons` converts SVG/PNGs into platform icon bundles.
- `ToolPackage` uses `internal/packager` to drive `nfpm` for Linux packages.
- `GenerateAppImage`, `generate_webview2`, and `package/msix.go` provide OS-specific installers.

---

## Template & Asset Generation
### Project Templates: `v3/internal/templates`
- Each frontend stack lives under its own directory (React, Vue, Svelte, Solid, Qwik, Vanilla, Lit, Preact) with `-ts` variants.
- `_common` contains shared scaffolding (Go module layout, `wails.json`, default assets).
- `generate template` extracts template files via `gosod` into the target project directory, honoring parameters collected during `wails3 init`.

### Build Assets: `internal/commands/build-assets.go`
Embedded assets define packaging metadata such as installers, file associations, protocol handlers, and Windows installer scripts.

```pseudo
func GenerateBuildAssets(opts):
    opts.Dir = abs(opts.Dir); mkdir if missing
    fill empty fields (ProductComments, Identifier, BinaryName, etc.) with sensible defaults
    load `build_assets` FS subtree and render via gosod using opts (includes file associations + protocols)
    render `updatable_build_assets` on top so user edits stay intact
```

### Runtime Bundles: `internal/runtime`
- `Core()` concatenates `runtimeInit + flags + invoke + environment` to produce the JS bootstrap injected into every window.
- Platform-specific files (`runtime_windows.go`, `runtime_linux.go`, etc.) define how Go exposes `window._wails.invoke` and system flags (resize handle sizes on Windows via `pkg/w32`).
- `GenerateRuntime` command serializes pre-built runtime assets for embedding.

---

## Runtime Boot Sequence (pkg/application)
### Application Construction: `pkg/application/application.go`
A singleton `App` is created via `application.New(options)` and stored in `globalApplication`.

```pseudo
func New(options):
    if globalApplication exists -> return it (enforces singleton)
    mergeApplicationDefaults(options)
    app := newApplication(options) // debug vs production build tags
    globalApplication = app
    fatalHandler(app.handleFatalError)
    configure Logger (debug -> DefaultLogger, prod -> discard)
    install default signal handler unless disabled (ctrl+c, SIGTERM -> App.Quit)
    log startup + platform info
    customEventProcessor := NewWailsEventProcessor(app.Event.dispatch)
    messageProc := NewMessageProcessor(app.Logger)
    assetOpts := assetserver.Options {
        Handler: options.Assets.Handler (default BundledAssetFileServer)
        Middleware: Chain(user middleware, internal middleware for /wails endpoints)
        Logger: app.Logger (or discard when DisableLogging)
    }
    asset server intercepts:
        /wails/runtime.js -> bundled runtime asset
        /wails/runtime    -> messageProc.ServeHTTP
        /wails/capabilities -> emits capabilities JSON
        /wails/flags -> marshals platform flags via impl.GetFlags
    app.assets = AssetServer(assetOpts)
    app.bindings = NewBindings(options.MarshalError, options.BindAliases)
    app.options.Services = clone(options.Services)
    process key bindings if provided
    register OnShutdown hook if provided
    if SingleInstance configured -> newSingleInstanceManager
    return app
```

Key supporting structures:
- `Options` (`application_options.go`) configures assets, services, platform-specific knobs, signal handling, keybindings, custom marshaling, and single-instance behavior.
- `signal.NewSignalHandler` (`internal/signal`) watches OS signals, invoking `App.Quit` while printing customizable exit messages.

### Run Loop & Lifecycle
`App.Run()` orchestrates startup, service lifecycle, channel fan-out, and the platform event loop.

```pseudo
func (a *App) Run():
    lock runLock; guard against double runs
    defer cancel app.Context (ensures goroutines exit on failure)
    execute a.preRun() (noop in production, debug logging otherwise)
    a.impl = newPlatformApp(a) // windowsApp, darwinApp, linuxApp depending on GOOS
    defer a.shutdownServices()
    services := clone(options.Services); options.Services = nil
    for service in services:
        startupService(service) // binds methods, registers HTTP routes, calls ServiceStartup
        append to options.Services so shutdown order is reversed
    spawn goroutines reading buffered channels:
        applicationEvents -> EventManager.handleApplicationEvent
        windowEvents -> handleWindowEvent
        webviewRequests -> assets.ServeWebViewRequest
        windowMessageBuffer -> handleWindowMessage (custom vs wails: prefix)
        windowKeyEvents -> handleWindowKeyEvent
        windowDragAndDropBuffer -> handleDragAndDropMessage (debug logs)
        menuItemClicked -> Menu.handleMenuItemClicked
    mark running=true; flush pendingRun queue by invoking runnable.Run() asynchronously
    if GOOS == darwin -> set application menu immediately
    if Icon provided -> impl.setIcon
    return a.impl.run() (enters platform main loop)
```

```d2
direction: down
root: "App.Run dispatch hub"
root -> applicationEvents: "goroutine"
applicationEvents -> EventManager: "Event.handleApplicationEvent"
root -> windowEvents: "goroutine"
windowEvents -> Windows: "handleWindowEvent"
root -> webviewRequests
webviewRequests -> AssetServer: "ServeWebViewRequest"
root -> windowMessages
windowMessages -> MessageHandlers: "HandleMessage / RawMessageHandler"
root -> windowKeyEvents
windowKeyEvents -> KeyBinding: "HandleKeyEvent"
root -> dragDropBuffer
dragDropBuffer -> Windows: "HandleDragAndDrop"
root -> menuItemClicked
menuItemClicked -> MenuManager
```

### Shutdown Path
- `shutdownServices()` iterates bound services in reverse start order, invoking `ServiceShutdown` when implemented and cancelling the app context.
- `OnShutdown` hooks run synchronously on the main thread; `PostShutdown` runs last (useful on macOS where `Run` may block indefinitely).
- `cleanup()` (triggered via `impl.destroy()` and `App.Quit`) sets `performingShutdown`, cancels the context, runs queued shutdown tasks, releases the single-instance manager, and closes windows/system trays.

---

## Message Bridge & Asset Server
### Message Processor (`messageprocessor.go` et al.)
- Handles `/wails/runtime` POST requests from the frontend.
- Dispatch keyed by `object` query parameter:
  - `callRequest` (0) → service bindings.
  - `clipboardRequest`, `applicationRequest`, `eventsRequest`, `contextMenuRequest`, `dialogRequest`, `windowRequest`, `screensRequest`, `systemRequest`, `browserRequest`, `cancelCallRequest`.
- Uses HTTP headers `x-wails-window-id` / `name` to resolve the target `Window` (`getTargetWindow`).
- Maintains `runningCalls` map to support cancellation by call ID.

```pseudo
func processCallMethod(method, rw, req, window, params):
    args := params.Args()
    callID := args.String("call-id")
    if method == CallBinding:
        options := params.ToStruct(CallOptions)
        ctx, cancel := context.WithCancel(request.Context)
        register cancel in runningCalls[callID]
        respond 200 OK immediately
        go func():
            boundMethod := lookup by name or ID (honour aliases)
            if not found -> CallError(kind=ReferenceError)
            if window != nil -> ctx = context.WithValue(ctx, WindowKey, window)
            result, err := boundMethod.Call(ctx, options.Args)
            on CallError -> window.CallError(callID, json, knownError=true)
            marshal result -> jsonResult
            window.CallResponse(callID, jsonResult)
            cleanup runningCalls entry and cancel context
        ```
```

- Errors bubble to the frontend via `CallError` JSON with `ReferenceError`, `TypeError`, or `RuntimeError` kinds.
- Cancellation (`cancelCallRequest`) removes the call ID from `runningCalls` and cancels the context.

### Asset Server (`internal/assetserver`)
- `AssetServer.ServeHTTP` wraps responses with MIME sniffing, logs duration, window metadata, and status codes.
- `serveHTTP` intercepts root/index requests to optionally serve a localized index fallback (`defaultIndexHTML` uses `accept-language`).
- `/wails/*` special routes and user-defined services share the same middleware chain, allowing injection of auth, logging, or routing.
- Dev vs production logic uses environment variable `FRONTEND_DEVSERVER_URL`; when present, requests proxy to the external dev server instead of the embedded assets (`asset_fileserver.go`).
- `AttachServiceHandler` mounts service-provided HTTP handlers under custom routes (`ServiceOptions.Route`).
- WebView-specific request/response types are implemented in `internal/assetserver/webview`, providing platform-native bridges to feed asset bytes directly into the webview without round-tripping through TCP when possible.

### Runtime JS Exposure
- `/wails/runtime.js` serves the concatenated runtime string produced by `internal/runtime`, ensuring the frontend has access to `window.wails` helpers.
- `/wails/flags` serializes `Options.Flags` extended by the platform implementation (`windowsApp.GetFlags` injects resize handles, etc.), allowing frontend startup logic to adjust to platform constraints.
- `/wails/capabilities` exposes a JSON describing features like native drag support (populated by `internal/capabilities`).

---

## Windowing & UI Layer (`pkg/application`)
### Window Manager (`window_manager.go`)
- Maintains the authoritative map of `Window` instances keyed by numeric ID.
- Defers actual `Run` execution if the app has not started (`runOrDeferToAppRun`).
- Provides lookup by name/ID, iteration (`GetAll`), and lifecycle hooks (`OnCreate`).
- Works in tandem with `App.pendingRun` to ensure windows created before `App.Run` are executed after the platform loop is ready.

### Webview Windows (`webview_window.go` + `webview_window_<platform>.go`)
- `WebviewWindow` wraps a platform-specific `webviewWindowImpl` with synchronized maps for event listeners, key bindings, menus, and asynchronous cancellers.
- Supports full window control API: sizing, positioning, zoom, devtools, menu bar toggles, context menu injection, border size queries, drag/resize operations.
- `HandleMessage`, `HandleKeyEvent`, and drag-and-drop handlers are invoked by the central channels in `App.Run`.
- On runtime readiness, windows emit `events.Common.WindowRuntimeReady` to allow frontends to hydrate state once the JS bridge is loaded.

### Event System (`events.go`, `context_*`)
- `ApplicationEvent` and `WindowEvent` objects carry strongly-typed contexts (`ApplicationEventContext`, `WindowEventContext`) to provide structured data (files dropped, URLs, screen info, etc.).
- Channels `applicationEvents`, `windowEvents`, and `menuItemClicked` are buffered to avoid blocking the OS event loop.
- `EventProcessor` manages custom user events, offering `On`, `Once`, `Emit`, and hook registration for pre-dispatch inspection.
- Platform-specific files (`events_common_windows.go`, `events_common_darwin.go`, etc.) translate native callbacks into the common channel structure.

### Managers & Subsystems
| Manager | File | Responsibility |
| --- | --- | --- |
| `ContextMenuManager` | `context_menu_manager.go` | Creates native context menus, wires handlers. |
| `DialogManager` | `dialog_manager.go` + OS-specific files | Wraps native file/message dialogs. |
| `ClipboardManager` | `clipboard_manager.go` | Provides cross-platform clipboard access, proxies to OS-specific implementations. |
| `ScreenManager` | `screenmanager.go` | Exposes multi-monitor info, resolution, scaling. |
| `SystemTrayManager` | `system_tray_manager.go` | Manages tray icons, menu interactions. |
| `BrowserManager` | `browser_manager.go` | Handles window navigation and devtools. |
| `KeyBindingManager` | `key_binding_manager.go` | Registers accelerators and callbacks. |
| `EnvironmentManager` | `environment_manager.go` | Tracks env state (dark mode, accent color) exposed to the frontend. |

---

## Services & Binding Engine
### Service Definition (`service.go`, `bindings.go`)
- Services are user-provided structs implementing optional interfaces:
  - `ServiceStartup(ctx context.Context, opts ServiceOptions) error`
  - `ServiceShutdown() error`
  - `ServeHTTP` (when exposing HTTP routes)
  - Methods to be bound must be exported, live on pointer receivers of named types, and cannot be generic.

```pseudo
func startupService(service):
    bindings.Add(service) // reflect, hash method signatures, register alias map
    if service.Route != "":
        if instance implements http.Handler -> attach to asset server
        else -> attach fallback handler returning 503
    if instance implements ServiceStartup -> call with app context
```

- `Bindings.Add` enumerates exported methods via reflection, hashes fully qualified names using `internal/hash`, and stores them in `boundMethods` (by name) and `boundByID` (by numeric ID). Hash collisions are guarded with explicit error messages instructing developers to rename methods.
- `CallOptions` submitted from the frontend include `MethodID`, `MethodName`, and JSON-encoded args; the binding engine supports both to allow smaller payloads in production (IDs) while keeping dev ergonomics (names).

### Error Marshalling & Aliases
- Custom error marshaling can be provided per service (`ServiceOptions.MarshalError`) or globally (`Options.MarshalError`), allowing user-defined JSON payloads.
- `BindAliases` maps alternative IDs to primary method IDs, helpful when generated bindings are versioned and need to stay stable across refactors.

### Context Propagation
- During calls, the bridge injects `context.Context` with the window (`context.WithValue(ctx, WindowKey, window)`) so services can inspect which window invoked them (e.g., to push events back via `window.Emit`).
- Cancellation is propagated when the frontend aborts a promise or the app shuts down (`App.cancel()`).

---

## Configuration Options & Their Effects
### `application.Options` Highlights (`application_options.go`)
| Field | Effect on Runtime |
| --- | --- |
| `Assets.Handler` | Base HTTP handler serving static assets. Overrides default embedded bundle, but `/wails/*` middleware still executes. |
| `Assets.Middleware` | Injects custom middleware before internal routes; can short-circuit requests to implement routing or authentication. |
| `Assets.DisableLogging` | Swaps the asset server logger to a discard handler to avoid noisy logs. |
| `Flags` | Merged with platform flags and exposed to the frontend at `/wails/flags`. Changing affects frontend boot configuration. |
| `Services` | List of services auto-registered during `Run`. Order matters for startup/shutdown. |
| `BindAliases` | Remaps method IDs used by the runtime. Critical when regenerating bindings without breaking existing frontend code. |
| `KeyBindings` | Global accelerator map executed per window. Processed during `New`, stored per window at runtime. |
| `OnShutdown` / `PostShutdown` | Lifecycle hooks executed during teardown. `OnShutdown` runs before services shut down; `PostShutdown` runs after platform loop returns (if ever). |
| `ShouldQuit` | Gatekeeper invoked when the user attempts to quit (e.g., Cmd+Q). Returning `false` keeps the app alive. |
| `RawMessageHandler` | Receives messages that do not start with the `wails:` prefix, enabling custom bridge protocols aside from service calls. |
| `WarningHandler` / `ErrorHandler` | Overrides default slog warnings/errors for system-level diagnostics. |
| `FileAssociations` | Used during packaging and when launch arguments are parsed on Windows/macOS to emit `ApplicationOpenedWithFile`. |
| `SingleInstance` | Triggers single-instance manager setup; options include encryption key, exit code, and `OnSecondInstanceLaunch` callback. |

### Single Instance Workflow (`single_instance.go`)
```pseudo
func newSingleInstanceManager(app, opts):
    if opts == nil -> return nil
    start goroutine that reads secondInstanceBuffer and dispatches OnSecondInstanceLaunch
    lock := newPlatformLock(opts.UniqueID) // OS-specific
    if lock.acquire fails -> already running
    return manager

func manager.notifyFirstInstance():
    data := {Args, WorkingDir, AdditionalData}
    payload := json.Marshal(data) or encrypt(AES-256-GCM)
    lock.notify(payload)
```
- On startup, if another instance is detected, the CLI prints a warning, invokes `notifyFirstInstance`, and exits with the configured code.
- Platform locks live in `single_instance_<os>.go`, using named pipes, mutexes, or DBus depending on OS.

---

## Platform Implementations
| Platform | Entry Files | Notes |
| --- | --- | --- |
| Windows | `application_windows.go`, `webview_window_windows*.go`, `pkg/w32` | Integrates with Win32 APIs, WebView2 via `go-webview2`. Handles taskbar recreation, dark mode, accent color, custom title bars, drag/resize, `WndProcInterceptor`. |
| macOS | `application_darwin.go`, `webview_window_darwin*.go`, `pkg/mac` | Objective-C bridges (via cgo) for NSApplication, implements activation policy, app/URL events, and main thread run loops. |
| Linux | `application_linux.go`, `webview_window_linux*.go`, `pkg/application/linux_*` | GTK/WebKit2 integration, theme detection, DBus single-instance hooks. Pure Go fallback provided (`linux_purego`). |
| Runtime JS | `runtime_<os>.go` | Defines `window._wails.invoke` binding for each platform (e.g., `webkit.messageHandlers.external.postMessage` on Linux/macOS, `chrome.webview.postMessage` on Windows). |

Key responsibilities of each platform struct (e.g., `windowsApp`):
- Manage native window/class registration, system tray IDs, focus management.
- Emit platform-specific events into `applicationEvents` when OS notifications occur (power events, taskbar resets, theme changes).
- Provide platform-specific implementations of `setApplicationMenu`, `dispatchOnMainThread`, `GetFlags`, etc.

---

## Packaging & Distribution Pipeline
- **Linux NFPM**: `internal/packager/packager.go` wraps `nfpm` to parse YAML recipes and build `.deb`, `.rpm`, `.apk`, `.ipk`, or Arch packages. CLI `wails3 tool package` selects the format and output path.
- **Windows**: `msix.go`, `generate_syso.go`, and build assets templates generate installers, embedding certificate details and executable paths from `BuildAssetsOptions`.
- **macOS**: `dmg` command scaffolds DMG packaging scripts.
- **AppImage**: `GenerateAppImage` compiles a portable AppImage by invoking the bundled `linuxdeploy` scripts (`appimage_testfiles` contains test fixtures).
- **Build Info**: `tool_buildinfo.go` prints deduced module versions using the captured `BuildSettings` map.

During release automation (see `v3/tasks/release`), these commands are orchestrated by Taskfile recipes to build cross-platform artifacts.

---

## Diagnostics, Logging & Error Handling
- Logging defaults differ: debug builds use structured logging (`DefaultLogger`), production builds discard framework logs unless a logger is provided.
- `fatalHandler` sets a package-level panic handler that exits the process on unrecoverable errors (panic + formatted message).
- `App.handleError`/`handleWarning` route errors to user-provided handlers or `Logger.Error/Warn`.
- Signal handling (SIGINT/SIGTERM) ensures the watcher and the app cleanly shut down, releasing OS resources and closing windows.
- `internal/term` and `term2` provide colored console output; used heavily in CLI warnings/errors to standardize UX (`term.Warningf`, `term.Hyperlink`).

---

## Developer Workflow Checklist
1. **Familiarize with Taskfile**: Run `task -l` or `wails3 task --list` to see available workflows (`install`, `precommit`, example builds).
2. **Local CLI Install**: `task v3:install` compiles the CLI (`go install` inside `cmd/wails3`).
3. **Run Dev Mode**: `wails3 dev` (or `task dev`) spawns the watcher, sets `FRONTEND_DEVSERVER_URL`, and tails the Go backend.
4. **Build Release**: `wails3 build` → alias for `task build` (see `Taskfile.yaml` for pipeline details, including `go build`, asset bundling, templating, packaging).
5. **Packaging Tests**: Use `wails3 tool package --type deb --config ...` or `GenerateBuildAssets` to refresh installer scaffolding.
6. **Pre-commit routine**: `task v3:precommit` runs `go test ./...` and repository formatters before opening PRs.

---

## Lifecycles at a Glance
### Application Startup Timeline
1. CLI entry (user application) calls `wails.Run()` (inside user code) → constructs `application.Options`.
2. `application.New` merges defaults, sets up asset server, bindings, single-instance manager, event processor.
3. User code configures windows/services (often via `app.NewWebviewWindow` or `app.RegisterService`).
4. `App.Run`:
   - Validates no concurrent runs.
   - Instantiates platform app (calls Objective-C/Win32/GTK setup).
   - Starts services, mounts HTTP routes.
   - Spins up channel goroutines and flushes deferred runnables.
   - Enters platform main loop (`impl.run()`), bridging native events into Go.
5. Frontend loads `runtime.js` → initializes `window.wails`, fetches `/wails/flags` & `/wails/capabilities`, and begins invoking bound methods via JSON payloads.

### Shutdown Timeline
1. `App.Quit` (called from user code, signal handler, or platform event) sets `performingShutdown`.
2. Cancels `context.Context` to abort long-running service calls.
3. Executes `OnShutdown` tasks in registration order.
4. Invokes `shutdownServices` in reverse startup order; each service can release resources.
5. Releases single-instance lock, destroys windows/system trays, and calls platform `destroy()`.
6. Executes `PostShutdown` hook (if provided) then exits the process.

---

## Key Channels & Buffers
| Channel | Producer | Consumer | Purpose |
| --- | --- | --- | --- |
| `applicationEvents` (`events.go`) | Platform apps (`application_<os>.go`) and system signals | `App.Event.handleApplicationEvent` | Broadcast application-level events (startup, theme changes, power events). |
| `windowEvents` | Platform window callbacks | `App.handleWindowEvent` | Window lifecycle events (move, resize, focus). |
| `webviewRequests` | Asset server (`webview/request.go`) | `App.handleWebViewRequest` | Direct webview resource streaming. |
| `windowMessageBuffer` | Platform message pump | `App.handleWindowMessage` | Messages from WebView (prefixed `wails:`) -> service invocations or raw handlers. |
| `windowDragAndDropBuffer` | Drag/drop handlers | `App.handleDragAndDropMessage` | File drop notifications with DOM metadata. |
| `windowKeyEvents` | Accelerator handlers | `App.handleWindowKeyEvent` | Dispatch registered key bindings to the appropriate window. |
| `menuItemClicked` | Native menu callbacks | `Menu.handleMenuItemClicked` | Execute Go handlers for menu events. |

All buffers default to size 5 to absorb bursty events without blocking UI threads.

---

## Environment Variables & Configuration Hooks
| Variable | Used By | Effect |
| --- | --- | --- |
| `FRONTEND_DEVSERVER_URL` | Asset server (`GetDevServerURL`) | Redirects asset requests to external dev server instead of embedded bundle. |
| `WAILS_VITE_PORT` | `commands.Dev` | Controls Vite dev server port advertised to frontend tasks. |
| `NO_COLOR` | Task runner | Disables colored CLI output. |
| `WAILSENV` | User applications (common pattern) | Choose between dev/prod logic when bootstrapping `Options`. |

---

## Debugging Tips
- **Inspect Call Flow**: Enable verbose logging (`Options.LogLevel = slog.LevelDebug`) to trace binding calls and HTTP requests.
- **Capture Runtime JS**: Hit `http://127.0.0.1:<port>/wails/runtime.js` while running in dev to verify injected flags and capabilities.
- **Watch Service Registration**: Look for `Registering bound method` debug lines (`bindings.go`) to confirm service methods were detected.
- **Single Instance Issues**: Check platform lock files (e.g., `%AppData%\Wails\Locks` on Windows, `$XDG_RUNTIME_DIR` on Linux) and ensure the encryption key is consistent.
- **Asset Server Paths**: Use `assetserver.ServeFile` error output (logged via slog) to diagnose missing assets or MIME type issues.

---

## Extending the Codebase
1. **Adding a CLI Command**: Implement function in `internal/commands`, define flags under `internal/flags`, register in `cmd/wails3/main.go`. Document the Taskfile hook if it wraps tasks.
2. **New Service Template**: Update `internal/service/template`, ensure `gosod` placeholders map to new options, and adjust `flags.ServiceInit` accordingly.
3. **Runtime Feature Flags**: Modify `internal/runtime/runtime_<os>.go` to expose additional data via `/wails/flags`; update frontend expectations.
4. **Custom Middleware**: Provide `Options.Assets.Middleware` to inject auth/logging; remember middleware runs before internal routes, so call `next` for default behavior.
5. **Platform-Specific Fixes**: Locate corresponding `application_<os>.go` and `webview_window_<os>*.go` files. Keep cross-platform interfaces (`platformApp`, `webviewWindowImpl`) stable.

---

## Reference: Critical Files by Responsibility
| Responsibility | File(s) | Notes |
| --- | --- | --- |
| CLI bootstrap | `v3/cmd/wails3/main.go` | Command registration, browser helpers. |
| Build info capture | `v3/cmd/wails3/main.go:init` | Populates `commands.BuildSettings`. |
| Task aliasing | `v3/internal/commands/task_wrapper.go` | Rewrites `os.Args` and invokes `RunTask`. |
| Asset server core | `v3/internal/assetserver/assetserver.go` | Middleware, logging, fallback handling. |
| Application singleton | `v3/pkg/application/application.go:56` | Global `App` creation. |
| Service binding | `v3/pkg/application/bindings.go` | Reflection, aliasing, error marshaling. |
| Message bridge | `v3/pkg/application/messageprocessor*.go` | Runtime call routing. |
| Event channels | `v3/pkg/application/events.go` | Buffered channels + processors. |
| Window API | `v3/pkg/application/webview_window.go` | Platform interface and user-facing API. |
| Single-instance control | `v3/pkg/application/single_instance*.go` | Locking and IPC. |
| Platform adapters | `v3/pkg/application/application_<os>.go` | Native message loops, start/stop. |
| Packaging | `v3/internal/packager/packager.go`, `internal/commands/tool_package.go` | Linux packaging wrappers. |

---

Armed with these maps, pseudocode, and lifecycle notes, you can confidently trace any bug from CLI invocation through runtime dispatch, into platform-specific glue, and back. Use the pseudocode as a mental model, verify concrete behavior by jumping into the referenced files, and lean on the diagrams to understand how data flows between layers.