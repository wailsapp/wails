### OnEvent

API:
`OnEvent(name string, callback func(event *CustomEvent)) func()`

`OnEvent()` registers an event listener for specific application events. The callback
function provided will be triggered when the corresponding event occurs.

### OffEvent
API:
`OffEvent(name string)`

`OffEvent()` removes an event listener for a specific named event specified.

### OnMultipleEvent
API:
`OnMultipleEvent(name string, callback func(event *CustomEvent), counter int) func()`

`OnMultipleEvent()` registers an event listener for X number of Events. The callback 
function provided will be triggered `counter` times when the corresponding event occurs.

### ResetEvents
API:
`ResetEvents()`

`ResetEvents()` removes all event listeners for all application events.

### OnApplicationEvent
API:
`OnApplicationEvent(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func()`

`OnApplicationEvent()` registers an event listener for specific application events.
The `eventType` is based on events.ApplicationEventType. See [ApplicationEventType](/API/events/#applicationevent)

### RegisterApplicationHook
API:
`RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func()`

`RegisterApplicationEventHook()` registers a callback to be triggered based on specific application events.

### RegisterHook

API:
`RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`RegisterHook()` registers a callback to be run as a hook during specific
events. These hooks are run before listeners attached with `On()`. The function
returns a function that can be called to remove the hook.
