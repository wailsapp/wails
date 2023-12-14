# Wails API

This package provides a typed Javascript API for Wails applications.

It provides methods for the following components:

- [Dialog](#dialog)
- [Events](#events)
- [Window](#window)
- [Plugin](#plugin)
- [Screens](#screens)
- [Application](#application)

## Installation

In your Wails application, run the following command in the frontend project directory:

```bash
npm install -D @wailsapp/api
```

## Usage

Import the API into your application:

```javascript
import * as Wails from "@wailsapp/api";
```

Then use the API components:

```javascript
function showDialog() {
    Wails.Dialog.Info({
        Title: "Hello",
    }).then((result) => {
        console.log("Result: " + result);
    });
}
```

Individual components of the API can also be imported directly.



## API

### Dialog

The Dialog API provides access to the native system dialogs.

```javascript
import { Dialog } from "@wailsapp/api";

function example() {
    Dialog.Info({
        Title: "Hello",
    }).then((result) => {
        console.log("Result: " + result);
    });
}
```

#### Message Dialogs

Message dialogs are used to display a message to the user.
They can be used to display information, errors, warnings and questions.
Each method returns the button that was pressed by the user.

- `Info(options: MessageDialogOptions): Promise<string>`
- `Error(options: MessageDialogOptions): Promise<string>`
- `Warning(options: MessageDialogOptions): Promise<string>`
- `Question(options: MessageDialogOptions): Promise<string>`

#### Open Dialog

The Open Dialog is used to open a file or directory. It returns the path of the selected file or directory.
If the `AllowsMultipleFiles` option is set, multiple files or directories can be selected and are returned
as an array of file paths.

- `Open(options: OpenDialogOptions): Promise<string[]|string>`

#### Save Dialog

The Save Dialog is used to save a file. It returns the path of the selected file.

- `Save(options: SaveDialogOptions): Promise<string>`

### Events

The Events API provides access to the Wails event system. This is a global event system
that can be used to send events between the Go and Javascript.
Events are available to every window in a multi-window application.
These API methods are specific to the window in which they are called in.

```javascript
import { Events } from "@wailsapp/api";

function example() {
    // Emit an event
    Events.Emit("myevent", { message: "Hello" });

    // Subscribe to an event
    let unsub = Events.On("otherEvent", (data) => {
        console.log("Received event: " + data);
    });

    // Unsubscribe from the event
    unsub();
}
```

#### Emit

Emit an event with optional data.

- `Emit(eventName: string, data?: any): void`

#### Subscribe

Three methods are provided to subscribe to an event:
  - `On(eventName: string, callback: (data: any) => void): () => void` - Subscribe to all events of the given name
  - `Once(eventName: string, callback: (data: any) => void): () => void` - Subscribe to one event of the given name
  - `OnMultiple(eventName: string, callback: (data: any) => void, count: number): () => void` - Subscribe to multiple events of the given name

The callback will be called when the event is emitted.
The returned function can be called to unsubscribe from the event.

#### Unsubscribe

As well as unsubscribing from a single event, you can unsubscribe from events of a given name or all events.
  - `Off(eventName: string, additionalEventNames: ...string): void` - Unsubscribe from all events of the given name(s)
  - `OffAll(): void` - Unsubscribe all events

### Window

The Window API provides a number of methods that interact with the window in which the API is called.

- `Center: (void) => void` - Center the window
- `SetTitle: (title) => void` - Set the window title
- `Fullscreen: () => void` - Set the window to fullscreen
- `UnFullscreen: () => void` - Restore a fullscreen window
- `SetSize: (width: number, height: number) => void` - Set the window size
- `Size: () => Size` - Get the window size
- `SetMaxSize: (width, height) => void` - Set the window maximum size
- `SetMinSize: (width, height) => void` - Set the window minimum size
- `SetAlwaysOnTop: (onTop) => void` - Set window to be always on top
- `SetPosition: (x, y) => void` - Set the window position
- `Position: () => Position` - Get the window position
- `SetResizable: (resizable) => void` - Set whether the window is resizable
- `Screen: () => Screen` - Get information about the screen the window is on
- `Hide: () => void` - Hide the window
- `Show: () => void` - Show the window
- `Maximise: () => void` - Maximise the window
- `Close: () => void` - Close the window
- `ToggleMaximise: () => void` - Toggle the window maximise state
- `UnMaximise: () => void` - UnMaximise the window
- `Minimise: () => void` - Minimise the window
- `UnMinimise: () => void` - UnMinimise the window
- `SetBackgroundColour: (r, g, b, a) => void` - Set the background colour of the window

### Plugin

The Plugin API provides access to the Wails plugin system.
This method provides the ability to call a plugin method from the frontend.

```javascript

import { Plugin } from "@wailsapp/api";

function example() {
    // Call a plugin method
    Plugin.Call("myplugin", "MyMethod", { message: "Hello" }).then((result) => {
        console.log("Result: " + result);
    });
}
```

### Screens

The Screens API provides access to the Wails screen system.

```javascript
import { Screens } from "@wailsapp/api";

function example() {
    // Get all attatched screens
    Screens.GetAll().then((screens) => {
        console.log("Screens: " + screens);
    });

    // Get the primary screen
    Screens.GetPrimary().then((screen) => {
        console.log("Primary screen: " + screen);
    });

    // Get the screen the window is on
    Screens.GetCurrent().then((screen) => {
        console.log("Window screen: " + screen);
    });
}
```

- `GetAll: () => Promise<Screen[]>` - Get all screens
- `GetPrimary: () => Promise<Screen>` - Get the primary screen
- `GetCurrent: () => Promise<Screen>` - Get the screen the window is on

### Application

The Application API provides access to the Wails application system.

```javascript
import { Application } from "@wailsapp/api";

function example() {

    // Hide the application
    Application.Hide();

    // Shopw the application
    Application.Show();

    // Quit the application
    Application.Quit();

}
```

- `Hide: () => void` - Hide the application
- `Show: () => void` - Show the application
- `Quit: () => void` - Quit the application
