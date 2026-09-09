# Window Example

This example is a demonstration of the Windows API.

## Running on macOS with private APIs

This example calls `OpenDevTools()`. On macOS, opening the inspector from code requires the `private_mac_apis` build tag. Without it, that call is a no-op; the rest of the example still runs.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

## Running the example

To run the example, simply run the following command:

```bash
# Public macOS APIs only; private effects and inspector opening are disabled.
go run .
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      |         |
| Windows  | Working |
| Linux    |         |
