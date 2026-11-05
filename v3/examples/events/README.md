# Events Example

This example is a demonstration of using the new events API.
It has 2 windows that can emit events from the frontend and the backend emits an event every 10 seconds.
All events emitted are logged either to the console or the window.

It also demonstrates the use of `RegisterHook` to register a function to be called when an event is emitted. 
For one window, it captures the `WindowClosing` event and prevents the window from closing twice.
The other window uses both hooks and events to show the window is gaining focus.

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

## Running the example

To run the example, simply run the following command:

```bash
go run main.go
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      |         |
| Windows  | Working |
| Linux    |         |