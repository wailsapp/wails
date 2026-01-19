# Spotlight Example

This example demonstrates how to create a Spotlight-like launcher window using the `CollectionBehavior` option on macOS.

## Features

- **Appears on all Spaces**: Using `MacWindowCollectionBehaviorCanJoinAllSpaces`, the window is visible across all virtual desktops
- **Overlays fullscreen apps**: Using `MacWindowCollectionBehaviorFullScreenAuxiliary`, the window can appear over fullscreen applications
- **Combined behaviors**: Demonstrates combining multiple behaviors with bitwise OR
- **Floating window**: `MacWindowLevelFloating` keeps the window above other windows
- **Accessory app**: Doesn't appear in the Dock (uses `ActivationPolicyAccessory`)
- **Frameless design**: Clean, borderless appearance with translucent backdrop

## Running the example

```bash
go run .
```

**Note**: This example is macOS-specific due to the use of `CollectionBehavior`.

## Combining CollectionBehaviors

Behaviors can be combined using bitwise OR (`|`):

```go
CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                    application.MacWindowCollectionBehaviorFullScreenAuxiliary,
```

## CollectionBehavior Options

These are bitmask values that can be combined:

**Space behavior:**
| Option | Description |
|--------|-------------|
| `MacWindowCollectionBehaviorDefault` | Uses FullScreenPrimary (default) |
| `MacWindowCollectionBehaviorCanJoinAllSpaces` | Window appears on all Spaces |
| `MacWindowCollectionBehaviorMoveToActiveSpace` | Moves to active Space when shown |
| `MacWindowCollectionBehaviorManaged` | Default managed window behavior |
| `MacWindowCollectionBehaviorTransient` | Temporary/transient window |
| `MacWindowCollectionBehaviorStationary` | Stays stationary during Space switches |

**Fullscreen behavior:**
| Option | Description |
|--------|-------------|
| `MacWindowCollectionBehaviorFullScreenPrimary` | Can enter fullscreen mode |
| `MacWindowCollectionBehaviorFullScreenAuxiliary` | Can overlay fullscreen apps |
| `MacWindowCollectionBehaviorFullScreenNone` | Disables fullscreen |
| `MacWindowCollectionBehaviorFullScreenAllowsTiling` | Allows side-by-side tiling |

## Use Cases

- **Launcher apps** (like Spotlight, Alfred, Raycast)
- **Quick capture tools** (notes, screenshots)
- **System utilities** that need to be accessible anywhere
- **Overlay widgets** that should appear over fullscreen apps

## Status

| Platform | Status |
|----------|--------|
| Mac      | Working |
| Windows  | N/A (macOS-specific feature) |
| Linux    | N/A (macOS-specific feature) |
