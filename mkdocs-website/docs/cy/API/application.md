# Cais

Mae'r API cais yn cynorthwyo i greu cais gan ddefnyddio fframwaith Wails.

### Newydd

API: `New(appOptions Options) *App`

`New(appOptions Options)` yn creu cais newydd gan ddefnyddio'r opsiynau cais a ddarperir. Mae'n cymhwyso gwerthoedd rhagosodedig ar gyfer opsiynau heb eu pennu, yn eu cyfuno â'r rhai a ddarparwyd, yn eu cychwyn a'n dychwelyd enghraifft o'r cais.

Os bydd gwall yn ystod y cychwyn, caiff y cais ei atal gyda'r neges gwall a ddarperir.

Dylid nodi, os oes enghraifft gyffredinol o gais yn bodoli eisoes, y bydd yr enghraifft honno'n cael ei dychwelyd yn hytrach na chreu un newydd.

```go title="main.go" hl_lines="6-9"
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name:        "Demo Ffenestr Gweddarlunydd",
		// Opsiynau eraill
    })

	// Gweddill y cais
}
```

### Cael

`Get()` yn dychwelyd yr enghraifft gyffredinol o'r cais. Mae'n ddefnyddiol pan fydd angen mynediad i'r cais o wahanol rannau o'ch cod.

```go
    // Cael enghraifft o'r cais
    app := application.Get()
```

### Galluoedd

API: `Capabilities() capabilities.Capabilities`

`Capabilities()` yn adfer map o'r galluoedd sydd gan y cais ar hyn o bryd. Gall y galluoedd fod ynghylch y nodweddion gwahanol y system weithredu sy'n darparu, fel nodweddion gweddarlunydd.

```go
    // Cael galluoedd y cais
    capabilities := app.Capabilities()
	if capabilities.HasNativeDrag {
		// Gwneud rhywbeth
    }
```

### GetPID

API: `GetPID() int`

`GetPID()` yn dychwelyd ID y Broses y cais.

```go
    pid := app.GetPID()
```

### Rhedeg

API: `Run() error`

`Run()` yn dechrau gweithredu'r cais a'i gydrannau.

```go
    app := application.New(application.Options{
	    //options
	})
    // Rhedeg y cais
    err := application.Run()
    if err != nil {
        // Ymdrin â'r gwall
    }
```

### Gadael

API: `Quit()`

`Quit()` yn gadael y cais trwy ddinistrio ffenestri a rhai cydrannau eraill o bosibl.

```go
    // Gadael y cais
    app.Quit()
```

### AydunDdyryslyd

API: `IsDarkMode() bool`

`IsDarkMode()` yn gwirio a yw'r cais yn rhedeg mewn modd tywyll. Mae'n dychwelyd gwerth boolean yn nodi a yw'r modd tywyll wedi'i alluogi.

```go
    // Gwiriwch a yw'r modd tywyll wedi'i alluogi
    if app.IsDarkMode() {
        // Gwneud rhywbeth
    }
```

### Cuddio

API: `Hide()`

`Hide()` yn cuddio ffenestr y cais.

```go
    // Cuddio ffenestr y cais
    app.Hide()
```

### Dangos

API: `Show()`

`Show()` yn dangos ffenestr y cais.

```go
    // Dangos ffenestr y cais
    app.Show()
```

--8<--
./docs/cy/API/application_window.md
./docs/cy/API/application_menu.md
./docs/cy/API/application_dialogs.md
./docs/cy/API/application_events.md
./docs/cy/API/application_screens.md
--8<--


## Opsiynau

```go title="pkg/application/application_options.go"
--8<--
../v3/pkg/application/application_options.go
--8<--
```