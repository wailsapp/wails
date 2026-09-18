# Dev Example

**NOTE**: This example is currently a work in progress. It is not yet ready for use.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Build the frontend first as required by this work-in-progress example; the tag does not change its existing setup requirements.

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.
