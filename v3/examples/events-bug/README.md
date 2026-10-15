# Events Bug Example

Demonstrates window events and an inspector keyboard shortcut.

## Running on macOS with private APIs

This example calls `OpenDevTools()`. On macOS, opening the inspector from code requires the `private_mac_apis` build tag. Without it, that call is a no-op; the rest of the example still runs.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.
