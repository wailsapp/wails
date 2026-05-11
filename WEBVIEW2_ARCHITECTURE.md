# Webview2 Architecture Decision Record

This document captures key architectural decisions and trade-offs made in the Webview2 v2 architecture implementation for Wails.

## Overview

The Webview2 v2 architecture represents a complete rewrite of the Windows webview binding layer, moving from the deprecated `pkg/edge/` implementation to a modern, maintainable `pkg/webview2/` based on direct COM access to the Microsoft Edge WebView2 SDK.

## Design Decisions

### 1. Programmatic Emitter vs Text Templates

**Decision:** Use programmatic Go code emission (`fmt.Fprintf`, `bytes.Buffer`) instead of text templates.

**Rationale:**
- Better debugging experience with full IDE support (step through generator code)
- Type safety: Go compiler catches errors at build time
- Easier to maintain: Go developers are more comfortable than template maintainers
- Better error messages: Go panics show exact line numbers vs template runtime errors

**Trade-offs:**
- Slightly more verbose than templates
- Requires Go code changes to modify output format
- Less separation between logic and presentation

**Impact:** All emitter logic is in Go, not in template files. This means the generator codebase is larger but more maintainable for Go developers.

### 2. File-per-Interface Organization

**Decision:** One file per interface (e.g., `ICoreWebView2.go`, `ICoreWebView2_2.go`, ..., `ICoreWebView2_27.go`).

**Rationale:**
- Readable diffs when interfaces change
- Enables parallel development by multiple contributors
- Easier code review: reviewers see only what changed in one interface
- Follows the existing pattern from `pkg/edge/`

**Trade-offs:**
- More files (30+ generated files vs one large file)
- Slightly more complex build system
- Requires directory navigation to find specific methods

**Impact:** Developers can easily see what changed in each interface version, and git history shows clear evolution of each interface independently.

### 3. Single Unified IDL Preprocessing

**Decision:** Inline all imports (objidl.idl, oaidl.idl, etc.) into one IDL before parsing.

**Rationale:**
- Parser sees a single, flat IDL tree without complex import resolution
- Simpler parser grammar with no import statements to handle
- Eliminates circular dependency issues in import resolution
- Microsoft IDL files are stable, so inlining is safe

**Trade-offs:**
- Preprocessing step adds complexity to the generator
- Larger memory usage during parsing
- Must be updated if Microsoft reorganizes their IDL includes

**Impact:** Parser implementation is significantly simpler and more maintainable. The preprocessing step is a one-time cost per generation.

### 4. Capabilities from Multiple Sources

**Decision:** Primary source: SDK release notes scraping. Secondary validation: IDL comments.

**Rationale:**
- Release notes are authoritative and published with each SDK version
- IDL comments provide validation cross-check
- Enables automatic capability detection for new features
- Reduces manual maintenance burden

**Trade-offs:**
- Scraping is fragile if Microsoft changes docs format
- May miss features not documented in release notes
- Requires occasional manual updates to `interface_versions.json`

**Impact:** Most capability detection is automatic, but manual intervention may be needed when Microsoft changes documentation format or adds undocumented features.

### 5. pkg/edge/ Deprecation (Not Deletion)

**Decision:** Archive `pkg/edge/` with deprecation notice, keep in git history.

**Rationale:**
- Safe fallback if `pkg/webview2/` has issues during rollout
- Allows gradual migration for users with custom `pkg/edge/` forks
- Preserves git history for reference and debugging
- Enables comparison between old and new implementations

**Trade-offs:**
- Two implementations in repo temporarily
- Potential confusion about which to use
- Slightly larger repository size

**Impact:** Users can pin to old behavior if needed, git history is preserved for reference, and there's a rollback path if issues arise.

### 6. VARIANT as Proper Struct

**Decision:** Define VARIANT as 16-byte struct matching Windows ABI.

**Rationale:**
- Enables correct `AddHostObjectToScript` implementation
- Enables correct `TryGetWebMessageAsString` implementation
- Matches Microsoft's COM ABI specification
- Required for proper COM interop

**Trade-offs:**
- Developers must understand VARIANT struct layout
- Requires careful handling of union fields
- More complex than simple type aliases

**Impact:** Enables critical COM features and future enhancements requiring VARIANT parameters. Correct memory layout is essential for COM interop.

## Rejected Approaches

### 1. Hand-Maintaining Generated Code

**Rejected because:** 306 files × N interfaces is unsustainable.

**Confirmation:** Current state has bugs that would take weeks to fix by hand. Auto-generation reduces the surface area for human error.

**Alternative:** Use code generator that parses Microsoft IDL files and emits Go bindings.

### 2. Using go-ole or Other COM Libraries

**Rejected because:** Tight coupling to third-party library; Wails wants control over COM interop.

**Confirmation:** Design of `webviewloader` and `combridge` shows preference for thin, custom wrappers.

**Alternative:** Implement thin COM bridge layer with direct vtable calls, giving Wails full control over memory and lifetime management.

### 3. Monolithic webview2.go File

**Rejected because:** 3000+ line files are hard to review and diff.

**Confirmation:** File-per-interface pattern in `pkg/edge/` already demonstrates the benefits of separation.

**Alternative:** One file per interface with clear naming convention (`ICoreWebView2.go`, `ICoreWebView2_2.go`, etc.).

## Future Considerations

### 1. Generics for Type-Safe COM Calls

**Status:** ✅ **IMPLEMENTED** (2026-05-11)

**Previous:** All vtable calls used `unsafe.Pointer` and casts.

**Now:** Type-safe generic helpers provide compile-time type safety and better IDE support.

**Files:**
- `webview2/pkg/webview2/com_generic.go` - Generic COM call helpers
- `webview2/pkg/webview2/com_generic_examples.go` - Usage examples and migration guide
- `webview2/pkg/webview2/com_generic_test.go` - Type safety tests

**Impact:** Requires Go 1.18+ (already required for combridge). Provides better developer experience with:
- IDE autocompletion
- Compile-time type checking
- Automatic memory management for COM strings
- Consistent error handling
- Reduced boilerplate

**Example:**
```go
// Old pattern
func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
	var value POINT
	hr, _, _ := i.Vtbl.GetDefaultDownloadDialogMargin.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return value, nil
}

// New pattern with generics
func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
	var value POINT
	return InvokeValue(i.Vtbl.GetDefaultDownloadDialogMargin, &value, uintptr(unsafe.Pointer(i)))
}
```

**Generic helpers available:**
- `InvokeVoid` - Methods returning only HRESULT
- `InvokeValue[T]` - Getters returning struct/primitive values
- `InvokeBool` - Getters returning bool
- `InvokeString` - Getters returning strings (with automatic memory cleanup)
- `InvokeInterface[T]` - Getters returning interface pointers
- `InvokeToken` - Event registration methods

### 2. Automatic Version Detection

**Current:** Capability checking requires explicit version query and manual configuration.

**Future:** Could auto-detect installed WebView2 version at environment creation time.

**Impact:** Simplifies consumer code (no manual version checks), but adds startup cost. Could provide caching to mitigate performance impact.

**Example:**
```go
// Current
env := webview2.NewEnvironmentWithOptions(version)

// Future
env := webview2.NewEnvironment() // Auto-detects latest installed version
```

### 3. Cross-Platform COM Simulation

**Current:** Windows only (COM is Windows-specific).

**Future:** Could add fake COM backend for testing on non-Windows platforms.

**Impact:** Better test coverage on CI, faster development iteration, but requires maintaining simulation layer that diverges from real COM behavior.

**Example:**
```go
// Add build tag: +build !windows

// Fake COM for testing on Linux/Mac
type fakeICoreWebView2 struct {}
```

### 4. Event Handler Optimization

**Current:** Each event registration allocates a new handler object.

**Future:** Could use object pooling or reuse handler objects.

**Impact:** Reduced GC pressure and memory usage, but adds complexity to lifecycle management.

### 5. Async/Await Pattern Support

**Current:** COM async methods use callback-based pattern.

**Future:** Could wrap callbacks in Go channels for async/await style code.

**Impact:** More idiomatic Go code, but requires careful error handling and goroutine management.

**Example:**
```go
// Current
webview.ExecuteScript(script, func(result *string, err error) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*result)
})

// Future
result, err := await(webview.ExecuteScriptAsync(script))
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

## Architecture Principles

1. **Maintainability over Convenience:** Code is written once, read many times. Prioritize readability and ease of modification.

2. **Type Safety over Flexibility:** Use Go's type system to catch errors at compile time rather than runtime.

3. **Explicit over Implicit:** Make memory management and lifetimes clear in the API surface.

4. **Standard Go Patterns:** Use idiomatic Go (errors, interfaces, channels) rather than porting COM patterns directly.

5. **Testing over Documentation:** Complex COM interop should be exercised by automated tests rather than just documented.

## Version Compatibility Matrix

| Webview2 Version | SDK Minimum | Wails Support | Notes |
|-----------------|-------------|---------------|-------|
| 1.0.0 - 1.0.1150 | 1.0.0 | v2.0+ | Initial release |
| 1.0.1150 - 1.0.1180 | 1.0.1150 | v2.1+ | Added Composition API |
| 1.0.1180+ | 1.0.1180 | v2.2+ | Latest features |

## Migration Guide for pkg/edge/ Users

1. Update imports: `github.com/wailsapp/wails/v2/pkg/edge` → `github.com/wailsapp/wails/webview2/pkg/webview2`

2. Environment creation changes slightly (new options API)

3. Event handler registration now uses explicit registration methods

4. Interface methods are now version-specific (ICoreWebView2_7.AddHostObjectToScript)

See pkg/webview2/examples/ for complete migration examples.

## Acceptance Criteria

- ✅ All design decisions documented with rationale
- ✅ Trade-offs clearly explained
- ✅ Architectural decisions recorded in this document
- ✅ Future maintainers can understand why choices were made
- ✅ Rejected approaches documented with reasoning
- ✅ Future considerations identified for evolution

---

**Document Version:** 1.0  
**Last Updated:** 2026-05-11  
**Maintainer:** Wails Team  
**Related Issues:** [WAI-299](mention://issue/5391dfb2-8651-4e39-a4f2-49db1428f98a)
