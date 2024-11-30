# Events

## Event Hooks

--8<--
./docs/en/API/event_hooks.md
--8<--

## Custom Events
You can create your own custom events that can be emitted and received on both the frontend and backend. 
Events are able to emitted at both the application and the window level. The receiver of the event gets data of where the
event was emitted from along with the data that was sent with the event. Events can be cancelled by the receiver.

=== "Go"
    ```go
    app.OnEvent("event1", func(e *application.CustomEvent) {
        app.Logger.Info("[Go] CustomEvent received", "name", e.Name, "data", e.Data, "sender", e.Sender, "cancelled", e.Cancelled)
        app.Logger.Info("[Go]", e.Data[0].(string)) // Logs "Hello from JS" to the terminal
    })

    window.EmitEvent("event2", "Hello from Go")
    ```
=== "JS"
    ```javascript
    wails.Events.Emit("event1", "Hello from JS")

    wails.Events.On("event2", function(event) {
        console.log("[JS] CustomEvent received", event)
        console.log(event.data) // prints "Hello from Go" to the webview console
    })

    ```

### Emitting Events

`application.EmitEvent(name string, data ...any)`
:   Emits an event from the application instance

`window.EmitEvent(name string, data ...any)`
:   Emits an event from the window instance

`wails.Events.Emit(event:wails.Events.EventData)`
:   Emits an event from the frontend sending an object with `name` and `data` properties or the typescript type WailsEvent 

### Receiving Events
Events can be received on the application instance and the frontend with a couple options of how 
you chose to receive them. You can register a single event listener that will trigger every time the event is emitted 
or you can register an event listener that will only trigger a specific number of times. 

Golang

`application.OnEvent(name string, handler func(data ...any))`
:   Registers an event on the application instance this will trigger every time the event is emitted

`application.OnMultipleEvent(name string, handler func(data ...any), count int)`
:   Registers an event on the application instance this will trigger every time the event is emitted up to the count specified 

Frontend

`wails.Events.On(name: string, callback: ()=>void)`, 
:   Registers an event on the frontend, this function returns a function that can be called to remove the event listener

`wails.Events.Once(name: string, callback: ()=>void)`, 
:   Registers an event on the frontend that will only be called once, this function returns a function that can be called to remove the event listener

`wails.Events.OnMultiple(name: string, callback: ()=>void, count: number)`, 
:   Registers an event on the frontend that will only be called `count` times, this function returns a function that can be called to remove the event listener

### Removing Events
There are a few ways to remove events that are registered. All of the registration functions return a function that can be called to remove the event listeneer
in the frontend. There are additional functions provided to help remove events as well.

Golang

`application.OffEvent(name string, ...additionalNames string)`
:   Removes an event listener with the specificed name

`application.ResetEvents()`
:   Removes all registered events and hooks

Frontend

`wails.Events.OffAll()`
:   Removes all registered events

`wails.Events.Off(name: string)`
:   Removes an event listener with the specified name


## Event Types

### ApplicationEvent
Returned when an application hook event is triggered. The event can be cancelled by calling the `Cancel()` method on the event.
```go
type ApplicationEvent struct {
  Id        uint
  ctx       *ApplicationEventContext
  Cancelled bool
}

// Cancel the event
func (a *ApplicationEvent) Cancel() {}
```

### WindowEvent
Returned when a window hook event is triggered. The event can be cancelled by calling the `Cancel()` method on the event.
```go
type WindowEvent struct {
	ctx       *WindowEventContext
	Cancelled bool
}

// Cancel the event
func (w *WindowEvent) Cancel() {}
```

### CustomEvent

CustomEvent is returned when an event is being recieved it includes the name of the event, the data that was sent with the event,
the sender of the event, application or a specific window. The event can be cancelled by calling the `Cancel()` method on the event.
```go
type CustomEvent struct {
	Name      string `json:"name"`
	Data      any    `json:"data"`
	Sender    string `json:"sender"`
	Cancelled bool
}

// Cancel the event
func (c *CustomEvent) Cancel() {}
```
