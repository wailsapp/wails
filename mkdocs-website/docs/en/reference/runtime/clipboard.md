# Clipboard

The Clipboard API provides access to the system clipboard, allowing you to read and write text data.

## SetText(text string) bool

The SetText() method sets the text content of the clipboard.

```go
success := clipboard.SetText("Hello, World!")
if !success {
    // Handle error
}
```

# Text() (string, bool)

The Text() method retrieves the current text content of the clipboard.

```go
text, success := clipboard.Text()
if !success {
    // Handle error
} else {
    fmt.Println(text)
}
```

