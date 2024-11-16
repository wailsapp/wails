## Dialogs

Dialogs are now available in JavaScript!

### Windows

Dialog buttons in Windows are not configurable and are constant depending on the
type of dialog. To trigger a callback when a button is pressed, create a button
with the same name as the button you wish to have the callback attached to.
Example: Create a button with the label `Ok` and use `OnClick()` to set the
callback method:

```go
        dialog := app.QuestionDialog().
			SetTitle("Update").
			SetMessage("The cancel button is selected when pressing escape")
		ok := dialog.AddButton("Ok")
		ok.OnClick(func() {
			// Do something
		})
		no := dialog.AddButton("Cancel")
		dialog.SetDefaultButton(ok)
		dialog.SetCancelButton(no)
		dialog.Show()
```