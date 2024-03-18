Dyma'r testun wedi'i gyfieithu i'r Gymraeg:

### NewWebviewWindow

API: `NewWebviewWindow() *WebviewWindow`

Mae `NewWebviewWindow()` yn creu ffenestr Webview newydd gyda'r opsiynau rhagosodedig, ac yn ei dychwelyd.

```go
    // Creu ffenestr webview newydd
    window := app.NewWebviewWindow()
```

### NewWebviewWindowWithOptions

API:
`NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow`

Mae `NewWebviewWindowWithOptions()` yn creu ffenestr webview newydd gydag opsiynau custom. Caiff y ffenestr newydd ei ychwanegu at fap o ffenestri a reolir gan y cymhwysiad.

```go
    // Creu ffenestr webview newydd gydag opsiynau custom
    window := app.NewWebviewWindowWithOptions(WebviewWindowOptions{
		Name: "Main",
        Title: "Fy Ffenestr",
        Width: 800,
        Height: 600,
    })
```

### OnWindowCreation

API: `OnWindowCreation(callback func(window *WebviewWindow))`

Mae `OnWindowCreation()` yn cofrestru ffwythiant alw-nôl i'w alw pan grëir ffenestr.

```go
    // Cofrestru ffwythiant alw-nôl i'w alw pan grëir ffenestr
    app.OnWindowCreation(func(window *WebviewWindow) {
        // Gwneud rhywbeth
    })
```

### GetWindowByName

API: `GetWindowByName(name string) *WebviewWindow`

Mae `GetWindowByName()` yn nôl ac yn dychwelyd ffenestr gyda enw penodol.

```go
    // Cael ffenestr drwy ei henw
    window := app.GetWindowByName("Main")
```

### CurrentWindow

API: `CurrentWindow() *WebviewWindow`

Mae `CurrentWindow()` yn nôl ac yn dychwelyd cyfeiriad at y ffenestr weithredol yn y cymhwysiad. Os nad oes ffenestr, mae'n dychwelyd nil.

```go
    // Cael y ffenestr gyfredol
    window := app.CurrentWindow()
```