# Changes for v3

## Events

In v3, there are 3 types of events:

- Application Events
- Window Events
- Custom Events

### Application Events

Application events are events that are emitted by the application. These events include native events such as `ApplicationDidFinishLaunching` on macOS.

### Window Events

Window events are events that are emitted by a window. These events include native events such as `WindowDidBecomeMain` on macOS.

### Custom Events

Events that the user defines are called `CustomEvents`. This is to differentiate them from the `Event` object that is used to communicate with the browser. CustomEvents are now objects that encapsulate all the details of an event. This includes the event name, the data, and the source of the event.

The data associated with a CustomEvent is now a single value. If multiple values are required, then a struct can be used.

### Event callbacks and `Emit` function signature

The signatures events callbacks (as used by `On`, `Once` & `OnMultiple`) have changed. In v2, the callback function received optional data. In v3, the callback function receives a `CustomEvent` object that contains all data related to the event.

Similarly, the `Emit` function has changed. Instead of taking a name and optional data, it now takes a single `CustomEvent` object that it will emit.

### `Off` and `OffAll`

In v2, `Off` and `OffAll` calls would remove events in both JS and Go. Due to the multi-window nature of v3, this has been changed so that these methods only apply to the context they are called in. For example, if you call `Off` in a window, it will only remove events for that window. If you use `Off` in Go, it will only remove events for Go.

### Logging

There was a lot of requests for different types of logging in v2 so for v3 we have simplified things to make it as customisable as you want. There is now a single call `Log` that takes a LogMessage object. This object contains the message, the level, and the source, the log message and any data to be printed out. The default logger is the Console logger, however any number of outputs to log to can be added. Simply add custom loggers to the `options.Application.Logger.CustomLoggers` slice. The default logger does not have log level filtering, however custom loggers can be added that do.

### Developer notes

When emitting an event in Go, it will dispatch the event to local Go listeners and also each window in the application.
When emitting an event in JS, it now sends the event to the application. This will be processed as if it was emitted in Go, however the sender ID will be that of the window.

## Window

The Window API has largely remained the same, however there are a few changes to note. 

## ClipBoard

The clipboard API has been simplified. There is now a single `Clipboard` object that can be used to read and write to the clipboard. The `Clipboard` object is available in both Go and JS. `SetText()` to set the text and `Text()` to get the text.

## Bindings

TBD

## Dialogs

Dialogs are now available in JavaScript! 

