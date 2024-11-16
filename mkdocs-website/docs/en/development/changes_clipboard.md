## ClipBoard

The clipboard API has been simplified. There is now a single `Clipboard` object
that can be used to read and write to the clipboard. The `Clipboard` object is
available in both Go and JS. `SetText()` to set the text and `Text()` to get the
text.