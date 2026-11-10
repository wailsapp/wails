# Wails Enhancement Proposal (WEP) 

## Do Windows GUI setup when an App is built, not at package init

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Srinivas Iyer  
**Created**: 2026-11-08  
**Discussion**: Split out of #6583; the Linux part is #6586 (bug fix, #6584)  
**Implementor**: Srinivas Iyer
**Target**: Wails v3

## Summary

On Windows, `pkg/application` and the packages it imports do their GUI setup in package `init()` and package-level variables. Any program that links Wails therefore pays for a GUI before `main`, even if it never builds an `App`:

- 18 to 19 DLLs load beyond what a plain Go program loads (user32, gdi32, ole32, combase, shell32, uxtheme, advapi32, …);
- the main thread becomes a single-threaded COM apartment;
- startup is 10 to 15 ms slower.

This proposal moves that work to `application.New` and first use. It adds one function to the public `pkg/w32` package, and it changes when `pkg/w32` sets up dark mode and when `Path`/`Paths` take their snapshot.

## Motivation

A common shape for a Wails app is one binary that is both the GUI and a CLI, or a GUI plus a background daemon, so the app and its command line can't drift apart. Every CLI call and the daemon link `pkg/application` but never open a window. Today they still load the GUI stack on Windows. For a program that only links the package and prints a line:

| | plain Go | v3 master | with this change |
|---|---|---|---|
| Windows Server 2025 x64 (`windows-latest`) | 7.6 ms, 6.1 MB | 19.5 ms, 13.2 MB | 9.3 ms, 7.7 MB |
| Windows 11 Arm64 (`windows-11-arm`) | 11.2 ms, 7.4 MB | 27.5 ms, 14.8 MB | 13.0 ms, 9.0 MB |

These are medians of interleaved runs, with peak working set.

The COM side effect matters beyond speed. A program that links Wails and wants a multithreaded apartment on its main thread, or a library that calls `CoInitializeEx` itself, finds the thread already initialised as an STA before its own code runs.

Where each load comes from was found with a `LdrRegisterDllNotification` hook that prints the Go stack behind each DLL load before `main`:

| Source | What it loads |
|---|---|
| `application`: `RegisterWindowMessage` for the invoke and `TaskbarCreated` messages | user32, gdi32, win32u, imm32 |
| `application/path.go`: `userpaths.Resolve()` (`SHGetKnownFolderPath`) | shell32, windows.storage, shcore, shlwapi, profapi |
| `w32/theme.go`: dark-mode uxtheme ordinals, `SetPreferredAppMode` | uxtheme |
| `w32/consts.go`: the Windows version, from the registry | advapi32, sechost |
| `webview2/pkg/edge`: `CoInitializeEx` | ole32, combase |
| `webview2/pkg/combridge`: `windows.GUIDFromString` and `GUID.String` while registering vtables | ole32 (and user32, gdi32 through it) |
| `go-common-file-dialog/cfd`: `ole.IIDFromString` | ole32 |

## Detailed Design

1. **`platformInit()` in `application.New`.** It's a no-op on most platforms; the Linux purego fix (#6586) introduces the hook. On Windows it:
   - registers the two window messages;
   - calls `w32.InitDarkMode()`;
   - calls `edge.InitCOM()`, but only when `New` runs on the startup thread, which keeps the main thread an STA from `New` on, as package init does today.
2. **`newPlatformApp`**, which `Run` calls on the startup thread before it starts services, also calls `edge.InitCOM()` and `w32.InitDarkMode()`. An App built on another goroutine therefore still has COM and dark mode on the window thread before any `ServiceStartup` can open a window.
3. **`w32.InitDarkMode()`** (new, exported). It does what `theme.go`'s `init()` does today, once: it resolves the uxtheme ordinals into the exported function variables (`AllowDarkModeForWindow` and the others) and opts the process into dark menus. `SetTheme` and `SetMenuTheme` call it.
4. **The Windows version** is read from the registry on first use (`sync.OnceValue`). `IsWindowsVersionAtLeast` returns false if it can't be read; it currently dereferences nil.
5. **`Path`/`Paths`** resolve the user's folders on their first call (`sync.OnceValue`). This is the one cross-platform change.
6. **combridge** parses and formats GUIDs in Go instead of through ole32, with a test against ole32's results. The file dialog package drops an IID that it parses at init and that nothing reads.
7. **`edge.InitCOM()`** replaces the `init()` that calls `CoInitializeEx`, and checks the HRESULT as an `int32`. The current check, `int(r) < 0`, never fires on 64-bit.

After this, a program that only links `pkg/application` loads the same DLLs as a plain Go program.

## Non-Goals

- **Linux** is a separate bug fix, #6586: it has no public API or behaviour change beyond startup.
- **macOS** needs nothing. Linking the package costs about 1 ms over plain Go there, nearly all of it from `net/http`/`crypto/x509` linking CoreFoundation and Security.
- **Startup work that loads no DLLs** is out of scope: building COM callback tables costs 1 to 2 ms on Windows.
- **The `WEBVIEW2_*` environment variables** that the webview2 loader sets at init are left as they are.

## Platform Considerations

- **Windows:** everything above.
- **Linux:** only the `Path`/`Paths` timing (5).
- **macOS:** only the `Path`/`Paths` timing (5).
- **Mobile and `server` builds:** they keep the no-op `platformInit`; the `Path` change applies where the package builds.

## Pros/Cons

**Pros**

- A CLI or daemon sharing a binary with a Wails app starts as fast as plain Go, and doesn't have COM set up on its main thread.
- Nothing changes for an app that creates an `App`: all the same setup happens before its first window and its services.

**Cons**

- Code using `pkg/w32` directly, without an `App`, must call `w32.InitDarkMode()` (or `SetTheme`/`SetMenuTheme`) before it relies on the dark-mode function variables or on dark menu text.
- The setup is now called from two places (`New` and `newPlatformApp`), each guarded by a `sync.Once`.

## Alternatives Considered

- **Leave `pkg/w32` as it is** and change only `pkg/application`. uxtheme, the registry read, and with it advapi32 still load at startup, about half the gain is lost, and `pkg/w32`'s own `init()` still opts every importer into dark menus.
- **Keep COM in `init()`** and move the rest. This leaves the STA side effect, and ole32/combase, in every program.
- **Split the GUI into its own binary** is what apps do today to avoid the cost. It gives up the single binary that keeps the app and its CLI in step.

## Backwards Compatibility

For apps, nothing changes: the setup still happens before the first window, the first service and the first menu.

For direct users of `pkg/w32` or `application.Path`:

- Importing `pkg/w32` no longer fills in the dark-mode function variables or opts the process into dark menus; call `w32.InitDarkMode()`. The variables can already be nil on Windows before 1809, so callers already check them.
- `Path`/`Paths` take their snapshot on the first call, not at program start. A program that changes `HOME`, `XDG_*` or `APPDATA` before its first call now sees the change.
- `IsWindowsVersionAtLeast` returns false instead of panicking when the version can't be read.

## Security and Privacy

None. No new permissions, data or network access. COM is initialised with the same apartment model, on the same thread, before any WebView2 is created.

## Test Plan

- **Wails' tests:** `go test ./pkg/w32/... ./internal/go-common-file-dialog/... ./internal/webview2/... ./pkg/application/` on `windows-latest` and `windows-11-arm`.
- **GUID helpers:** a unit test compares them with `windows.GUIDFromString`/`GUID.String`.
- **Real windows:** a smoke test opens a window, waits for the page to load, and runs an `InvokeSync` from a goroutine, in three cases:
  - the App built on the main goroutine;
  - the App built on another goroutine;
  - the window opened from a `ServiceStartup`.

  It runs against master and the branch.
- **Startup cost:** measured, and the DLLs loaded before `main` recorded with the load hook, against plain Go and master.
- **Builds:** windows amd64/arm64/386 and `server`; `go vet` unchanged from master.

The probe program, the measuring harness, the smoke test and the DLL hook are in [doze-dev/doze `tools/wailswin`](https://github.com/doze-dev/doze/tree/main/tools/wailswin). The workflow that runs them is [`wails-windows.yml`](https://github.com/doze-dev/doze/blob/main/.github/workflows/wails-windows.yml).

## Reference Implementation

This has been implemented and tested on [doze-dev/wails `fix/linux-purego-lazy-gtk`](https://github.com/wailsapp/wails/compare/master...doze-dev:wails:fix/linux-purego-lazy-gtk), which was closed PR #6583. Its first commit is the Linux fix (#6586). The other three are this proposal:

- "resolve the user's folders on first use";
- "stop calling ole32 at package init on Windows";
- "do the Windows GUI setup when an App is built".

All the numbers above come from that branch.

## Maintenance Plan

The pattern is small: a `sync.Once` per piece of setup, called from `platformInit` and `newPlatformApp`. New Windows setup that needs a GUI goes there rather than into `init()`. A check that a program linking `pkg/application` loads no DLLs beyond plain Go before `main` could run in CI so this doesn't regress; the probe and hook above are enough for that.

## Conclusion

Linking Wails on Windows would cost what linking any Go package costs, and nothing more. Apps keep their startup order. `pkg/w32` gains one explicit call, `InitDarkMode`, in place of an import side effect. I'll implement it, starting from the reference branch, once the proposal is accepted.
