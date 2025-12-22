# Task 24: Improve pkg/application Test Coverage

## Summary

This PR adds comprehensive unit tests to the `pkg/application` package, improving test coverage from **13.6%** to **17.7%**.

## Changes

### New Test Files Created

1. **`context_test.go`** - Tests for Context struct operations including:
   - `newContext()` initialization
   - `ClickedMenuItem()` getter with exists/not-exists cases
   - `IsChecked()` getter with true/false/not-set cases
   - `ContextMenuData()` with various edge cases
   - Method chaining tests

2. **`services_test.go`** - Tests for Service management including:
   - `NewService()` and `NewServiceWithOptions()`
   - `getServiceName()` priority (options > interface > type name)
   - `Service.Instance()` accessor
   - `DefaultServiceOptions` defaults
   - Service lifecycle interfaces (`ServiceStartup`, `ServiceShutdown`)

3. **`parameter_test.go`** - Tests for Parameter and CallError types:
   - `Parameter.IsType()` and `Parameter.IsError()`
   - `newParameter()` factory function
   - `CallError.Error()` method
   - Error kinds (ReferenceError, TypeError, RuntimeError)
   - `CallOptions` struct fields

4. **`dialogs_test.go`** - Tests for dialog utilities:
   - `getDialogID()` and `freeDialogID()` ID management
   - `Button` methods (`OnClick`, `SetAsDefault`, `SetAsCancel`)
   - Method chaining on buttons
   - DialogType constants
   - `MessageDialogOptions`, `FileFilter`, `OpenFileDialogOptions`, `SaveFileDialogOptions` fields

5. **`webview_window_options_test.go`** - Tests for window options:
   - `NewRGBA()`, `NewRGB()`, `NewRGBPtr()` helper functions
   - All constant types (BackgroundType, BackdropType, DragEffect, Theme, etc.)
   - MacTitleBar preset configurations
   - Default values for platform-specific window options

6. **`application_options_test.go`** - Tests for application configuration:
   - `ActivationPolicy` and `NativeTabIcon` constants
   - `ChainMiddleware()` function with empty/single/multiple middleware
   - Middleware short-circuit behavior
   - Default values for Options, MacOptions, WindowsOptions, LinuxOptions, IOSOptions, AndroidOptions

7. **`keys_test.go`** - Tests for keyboard accelerator parsing:
   - `parseKey()` for named keys, single chars, and special cases
   - `parseAccelerator()` for full keyboard shortcuts
   - Modifier constant uniqueness
   - `accelerator.String()` formatting
   - `accelerator.clone()` behavior
   - `modifierMap` and `namedKeys` contents

8. **`single_instance_test.go`** - Tests for single instance management:
   - `encrypt()` and `decrypt()` AES-256-GCM operations
   - Edge cases (empty data, wrong key, invalid data)
   - `getLockPath()` path construction
   - `getCurrentWorkingDir()` utility
   - `SecondInstanceData` and `SingleInstanceOptions` field defaults
   - Nil-safety in `singleInstanceManager.cleanup()`

9. **`menuitem_internal_test.go`** - Internal tests for menu items:
   - Menu item type constants
   - All `NewMenuItem*` factory functions
   - Menu item map operations (`addToMenuItemMap`, `getMenuItemByID`, `removeMenuItemByID`)
   - All setter methods (`SetLabel`, `SetEnabled`, `SetChecked`, `SetHidden`, etc.)
   - Accelerator methods
   - `Clone()` behavior
   - Method chaining

10. **`menu_internal_test.go`** - Internal tests for menus:
    - `NewMenu()` and all `Add*` methods
    - `FindByLabel()` including nested submenus
    - `RemoveMenuItem()`, `Clear()`, `Append()`, `Prepend()`
    - `Clone()` for menus
    - `processRadioGroups()` for radio button grouping
    - `setContextData()` propagation

11. **`screenmanager_internal_test.go`** - Internal tests for screen management:
    - Alignment and OffsetReference constants
    - `Rect` methods: `Origin`, `Corner`, `InsideCorner`, `right`, `bottom`, `Size`, `IsEmpty`, `Contains`, `Intersect`, `distanceFromRectSquared`
    - `Screen` methods: `Origin`, `scale`, `right`, `bottom`, `intersects`
    - `Point`, `Size`, `ScreenPlacement` field tests
    - `newScreenManager()` initialization

## Coverage Analysis

### Why 40% Target Wasn't Fully Achieved

The `pkg/application` package presents unique testing challenges:

1. **Platform-Specific Code (~50% of codebase)**
   - `*_darwin.go`, `*_windows.go`, `*_linux.go`, `*_ios.go`, `*_android.go` files
   - CGO code (`linux_cgo.go`, `linux_purego.go`)
   - These files can only be tested on their respective platforms

2. **Runtime Dependencies**
   - Many functions require `globalApplication` to be initialized
   - Window management functions require active webview instances
   - Platform-specific system calls that can't be mocked easily

3. **GUI Dependencies**
   - Functions that interact with window systems (X11, Cocoa, Win32)
   - Event loop integration that requires running application

### What Can Be Tested

The tests focus on:
- Pure Go logic (struct methods, helper functions)
- Type constants and default values
- Data structures and their operations
- Utility functions (encryption, path handling, parsing)
- Internal state management (menu items, dialogs, contexts)

### Recommendations for Further Coverage

1. **Integration Tests**: Run platform-specific tests in CI with appropriate runners
2. **Mock Interfaces**: Create mock implementations for `menuItemImpl`, `platformLock`, etc.
3. **Test Fixtures**: Set up test applications for window-related tests
4. **Platform-Specific Test Files**: Create `*_test.go` files with build tags for each platform

## Files Modified

- `v3/UNRELEASED_CHANGELOG.md` - Added changelog entry

## Testing

All tests pass:
```bash
cd v3/pkg/application && go test ./...
# ok  github.com/wailsapp/wails/v3/pkg/application  coverage: 17.7% of statements
```
