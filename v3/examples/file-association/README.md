# File Association Sample Project

This sample project demonstrates how to associate a file type with an application.
More info at: https://v3.wails.io/learn/guides/file-associations/

To run the sample, follow these steps:

1. Run `wails3 package` to generate the package.
2. On Windows, run the installer that was built in the `bin` directory.
3. Double-click on the `test.wails` file to open it with the application.
4. On macOS, double-click on the `test.wails` file and select the built application.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.
