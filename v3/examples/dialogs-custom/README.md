# Custom Dialogs Example

This example demonstrates the new custom dialog functionality in Wails v3, which allows you to:

- Create dialogs with custom button text (including emojis)
- Use individual button callbacks
- Set default and cancel buttons
- Maintain the same API as standard dialogs

## Features Demonstrated

### 1. Custom Button Text
Shows how to create dialogs with descriptive, custom button labels instead of generic "OK/Cancel" buttons.

### 2. Action Buttons
Demonstrates using emojis and descriptive text for better user experience.

### 3. Custom Confirmation
Example of a delete confirmation dialog with clear, specific button labels.

### 4. Multiple Choice
Shows how to present multiple options in a single dialog.

### 5. Game Style
Example of how custom dialogs can enhance gaming or entertainment applications.

## Technical Implementation

The custom dialog functionality automatically detects when you're using custom button text and switches from the standard Windows MessageBox API to the more flexible TaskDialog API on Windows. This provides:

- Support for any number of custom buttons
- Custom button text with Unicode support (including emojis)
- Individual button callbacks
- Proper default and cancel button handling
- Consistent API across all dialog types

## Running the Example

```bash
cd v3/examples/dialogs-custom
go run main.go
```

Navigate through the menus to see different examples of custom dialogs in action. Compare them with the traditional dialogs to see the difference.

## API Usage

```go
dialog := application.QuestionDialog()
dialog.SetTitle("Custom Dialog")
dialog.SetMessage("Choose your action:")

btn1 := dialog.AddButton("🎯 Action 1")
btn1.OnClick(func() {
    // Handle button 1 click
})

btn2 := dialog.AddButton("🔧 Action 2") 
btn2.OnClick(func() {
    // Handle button 2 click
})

dialog.SetDefaultButton(btn1) // Set default button
dialog.Show()
```

The same API works for all dialog types: `InfoDialog()`, `WarningDialog()`, `ErrorDialog()`, and `QuestionDialog()`.