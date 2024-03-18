
### On

API:
`On(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`On()` registers an event listener for specific application events. The callback
function provided will be triggered when the corresponding event occurs. The
function returns a function that can be called to remove the listener.

### RegisterHook

API:
`RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`RegisterHook()` registers a callback to be run as a hook during specific
events. These hooks are run before listeners attached with `On()`. The function
returns a function that can be called to remove the hook.
