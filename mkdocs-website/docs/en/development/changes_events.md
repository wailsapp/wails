## Events

In v3, there are 3 types of events:

- Application Events
- Window Events
- Custom Events

### Application Events

Application events are events that are emitted by the application. These events
include native events such as `ApplicationDidFinishLaunching` on macOS.

### Window Events

Window events are events that are emitted by a window. These events include
native events such as `WindowDidBecomeMain` on macOS. Common events are also
defined, so they work cross-platform, e.g. `WindowClosing`.

### Custom Events

Events that the user defines are called `WailsEvents`. This is to differentiate
them from the `Event` object that is used to communicate with the browser.
WailsEvents are now objects that encapsulate all the details of an event. This
includes the event name, the data, and the source of the event.

The data associated with a WailsEvent is now a single value. If multiple values
are required, then a struct can be used.

### Event callbacks and `Emit` function signature

The signatures events callbacks (as used by `On`, `Once` & `OnMultiple`) have
changed. In v2, the callback function received optional data. In v3, the
callback function receives a `WailsEvent` object that contains all data related
to the event.

Similarly, the `Emit` function has changed. Instead of taking a name and
optional data, it now takes a single `WailsEvent` object that it will emit.

### `Off` and `OffAll`

In v2, `Off` and `OffAll` calls would remove events in both JS and Go. Due to
the multi-window nature of v3, this has been changed so that these methods only
apply to the context they are called in. For example, if you call `Off` in a
window, it will only remove events for that window. If you use `Off` in Go, it
will only remove events for Go.

### Hooks

Event Hooks are a new feature in v3. They allow you to hook into the event
system and perform actions when certain events are emitted. For example, you can
hook into the `WindowClosing` event and perform some cleanup before the window
closes. Hooks can be registered at the application level or at the window level
using `RegisterHook`. Application level are for application events. Window level
hooks will only be called for the window they are registered with.

### Developer notes

When emitting an event in Go, it will dispatch the event to local Go listeners
and also each window in the application. When emitting an event in JS, it now
sends the event to the application. This will be processed as if it was emitted
in Go, however the sender ID will be that of the window.
