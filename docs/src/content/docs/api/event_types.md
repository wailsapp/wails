---
title: Event Types
sidebar:
  order: 140
---

### ApplicationEvent

Returned when an application hook event is triggered. The event can be cancelled
by calling the `Cancel()` method on the event.

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

Returned when a window hook event is triggered. The event can be cancelled by
calling the `Cancel()` method on the event.

```go
type WindowEvent struct {
    ctx       *WindowEventContext
    Cancelled bool
}

// Cancel the event
func (w *WindowEvent) Cancel() {}
```

### CustomEvent

CustomEvent is returned when an event is being received it includes the name of
the event, the data that was sent with the event, the sender of the event,
application or a specific window. The event can be cancelled by calling the
`Cancel()` method on the event.

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
