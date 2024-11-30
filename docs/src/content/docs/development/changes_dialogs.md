---
title: Changes (Dialogs)
sidebar:
  order: 100
---

Dialogs are now available in JavaScript!

### Windows

Dialog buttons in Windows are not configurable and are constant depending on the
type of dialog. To trigger a callback when a button is pressed, create a button
with the same name as the button you wish to have the callback attached to.
Example: Create a button with the label `Ok` and use `OnClick()` to set the
callback method:

```go
// Create a question dialog
dialog := app.QuestionDialog().

// Configure dialog title and message
SetTitle("Update").
SetMessage("The cancel button is selected when pressing escape")

// Add "Ok" button with callback
ok := dialog.AddButton("Ok")
ok.OnClick(func() {
  // Handle successful confirmation (Or do something else)
  if err := handleConfirmation(); err != nil {
    log.Printf("Error handling confirmation: %v", err)
  }
})

// Add "Cancel" button and configure dialog behavior
no := dialog.AddButton("Cancel")
dialog.SetDefaultButton(ok)
dialog.SetCancelButton(no)

// Show dialog and handle potential errors
if err := dialog.Show(); err != nil {
  log.Printf("Error showing dialog: %v", err)
}
```
