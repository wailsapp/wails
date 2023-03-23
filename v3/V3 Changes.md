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

The Window API has largely remained the same, however the methods are now on an instance of a window rather than the runtime. 
Some notable differences are:
- Windows now have a Name that identifies them. This is used to identify the window when emitting events.

## ClipBoard

The clipboard API has been simplified. There is now a single `Clipboard` object that can be used to read and write to the clipboard. The `Clipboard` object is available in both Go and JS. `SetText()` to set the text and `Text()` to get the text.

## Bindings

TBD

## Dialogs

Dialogs are now available in JavaScript! 

## Drag and Drop

Native drag and drop can be enabled per-window. Simply set the `EnableDragAndDrop` window config option to `true` and the window will allow files to be dragged onto it. When this happens, the `events.FilesDropped` event will be emitted. The filenames can then be retrieved from the WindowEventContext using the `DroppedFiles()` method. This returns a slice of strings containing the filenames.

## Context Menus

Context menus are contextual menus that are shown when the user right-clicks on an element. Creating a context menu is the same as creating a standard menu , by using `app.NewMenu()`. To make the context menu available to a window, call `window.RegisterContextMenu(name, menu)`. The name will be the id of the context menu and used by the frontend.

To indicate that an element has a context menu, add the `data-contextmenu` attribute to the element. The value of this attribute should be the name of a context menu previously registered with the window.

It is possible to register a context menu at the application level, making it available to all windows. This can be done using `app.RegisterContextMenu(name, menu)`. If a context menu cannot be found at the window level, the application context menus will be checked. A demo of this can be found in `v3/examples/contextmenus`.

## Wails Markup Language (WML)

The Wails Markup Language is a simple markup language that allows you to add functionality to standard HTML elements without the use of Javascript. 

The following tags are currently supported:

### `data-wml-event` 

This specifies that a Wails event will be emitted when the element is clicked. The value of the attribute should be the name of the event to emit. 

Example:
```html
<button data-wml-event="myevent">Click Me</button>
```
Sometimes you need the user to confirm an action. This can be done by adding the `data-wml-confirm` attribute to the element. The value of this attribute will be the message to display to the user.

Example:
```html
<button data-wml-event="delete-all-items" data-wml-confirm="Are you sure?">Delete All Items</button>
```

### `data-wml-window`

Any `wails.window` method can be called by adding the `data-wml-window` attribute to an element. The value of the attribute should be the name of the method to call. The method name should be in the same case as the method.

```html
<button data-wml-window="Close">Close Window</button>
```

### `data-wml-trigger`

This attribute specifies which javascript event should trigger the action. The default is `click`. 

```html
<button data-wml-event="hover-box" data-wml-trigger="mouseover">Hover over me!</button>
```

## Plugins

Plugins are a way to extend the functionality of your Wails application.

### Creating a plugin

Plugins are standard Go structure that adhere to the following interface:

```go
type Plugin interface {
    Name() string
    Init(*application.App) error
    Shutdown() 
    CallableByJS() []string
    InjectJS() string
}
```

The `Name()` method returns the name of the plugin. This is used for logging purposes.

The `Init(*application.App) error` method is called when the plugin is loaded. The `*application.App`
parameter is the application that the plugin is being loaded into. Any errors will prevent
the application from starting.

The `Shutdown()` method is called when the application is shutting down.

The `CallableByJS()` method returns a list of exported functions that can be called from
the frontend. These method names must exactly match the names of the methods exported
by the plugin.

The `InjectJS()` method returns JavaScript that should be injected into all windows as they are created. This is useful for adding custom JavaScript functions that complement the plugin.

### Tips

#### Enums

In Go, enums are often defined as a type and a set of constants. For example:

```go
type MyEnum int

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

Due to incompatibility between Go and JavaScript, custom types cannot be used in this way. The best strategy is to use a type alias for float64:

```go
type MyEnum = float64

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

In Javascript, you can then use the following:

```js
const MyEnum = {
    MyEnumOne: 0,
    MyEnumTwo: 1,
    MyEnumThree: 2
}
```

- Why use `float64`? Can't we use `int`? 
  - Because JavaScript doesn't have a concept of `int`. Everything is a `number`, which translates to `float64` in Go. There are also restrictions on casting types in Go's reflection package, which means using `int` doesn't work. 
