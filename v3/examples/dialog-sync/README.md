# dialog-sync

This example demonstrates a synchronous (blocking) message dialog flow using `MessageDialog.Result()`.

## What it shows

- Creating a `Question` dialog via `app.Dialog.Question()`
- Configuring buttons via the fluent API (`AddDefaultButton`, `AddCancelButton`)
- Blocking until the user selects a button using `Result()`
- Handling the dialog result and showing a follow-up dialog

## Run

From the `wails/v3` directory:

```bash
go run ./examples/dialog-sync
```

## Expected behaviour

- A window titled "Dialog Sync" will open
- A native question dialog will appear, attached to that window
- Clicking a button will log the result to the terminal and show a follow-up info dialog
