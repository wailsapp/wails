# HTML Drag and Drop API Example

This example should demonstrate whether the [HTML Drag and Drop API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API") works correctly.

## Expected Behaviour

When dragging the "draggable", in the console should be printed:
1. "dragstart" once
2. "drag" many times
3. "dragend" once

When dragging the "draggable" on the drop target, the inner text of the latter shoud change and in the console should be printed:
1. "dragstart" once
2. "drag" many times
3. "dragenter" once
4. "dragover" many times (alternating with "drag")
5.  - "drop" once (in case of a drop inside the drop target)
    - "dragleave" once (in case the draggable div leaves the drop target)
6. "dragend" once

## Running the example

To run the example, simply run the following command:

```bash
go run main.go
```

## Building the example

To build the example in debug mode, simply run the following command:

```bash
wails3 task build
```

# Status

| Platform | Status      |
|----------|-------------|
| Mac      | Working     |
| Windows  | Not Working |
| Linux    |             |
