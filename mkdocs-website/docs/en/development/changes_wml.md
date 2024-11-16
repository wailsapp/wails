## Wails Markup Language (WML)

The Wails Markup Language is a simple markup language that allows you to add
functionality to standard HTML elements without the use of Javascript.

The following tags are currently supported:

### `data-wml-event`

This specifies that a Wails event will be emitted when the element is clicked.
The value of the attribute should be the name of the event to emit.

Example:

```html
<button data-wml-event="myevent">Click Me</button>
```

Sometimes you need the user to confirm an action. This can be done by adding the
`data-wml-confirm` attribute to the element. The value of this attribute will be
the message to display to the user.

Example:

```html
<button data-wml-event="delete-all-items" data-wml-confirm="Are you sure?">
  Delete All Items
</button>
```

### `data-wml-window`

Any `wails.window` method can be called by adding the `data-wml-window`
attribute to an element. The value of the attribute should be the name of the
method to call. The method name should be in the same case as the method.

```html
<button data-wml-window="Close">Close Window</button>
```

### `data-wml-trigger`

This attribute specifies which javascript event should trigger the action. The
default is `click`.

```html
<button data-wml-event="hover-box" data-wml-trigger="mouseover">
  Hover over me!
</button>
```
