### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

Mae `RegisterContextMenu()` yn cofrestru dewislen cyd-destun gyda enw penodol. Gellir defnyddio'r dewislen hon yn ddiweddarach yn yr ap.

```go
    // Creu dewislen newydd
    ctxmenu := app.NewMenu()

    // Cofrestru'r dewislen fel dewislen cyd-destun
    app.RegisterContextMenu("MyContextMenu", ctxmenu)
```

### SetMenu

API: `SetMenu(menu *Menu)`

Mae `SetMenu()` yn gosod y ddewislen ar gyfer yr ap. Ar Mac, bydd hyn yn fod y ddewislen fyd-eang. Ar gyfer Windows a Linux, bydd hyn yn fod y ddewislen ddiofyn ar gyfer unrhyw ffenestr newydd a grëir.

```go
    // Creu dewislen newydd
    menu := app.NewMenu()

    // Gosod y ddewislen ar gyfer yr ap
    app.SetMenu(menu)
```