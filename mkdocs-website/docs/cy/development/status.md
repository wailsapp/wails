Dyma'r cyfieithiad Cymraeg (cym):

# Statws

Statws nodweddion yn v3.

!!! note

        Mae'r rhestr hon yn gymysgedd o gymorth API cyhoeddus a mewnol.<br/>
        Nid yw'n gyflawn ac efallai nad yw'n gyfoes.

## Problemau Hysbys

- Nid yw Linux eto ar barity nodwedd gyda Windows/Mac

## Cymhwyster

Dulliau rhyngwyneb cymhwyster

| Dull                                                          | Windows | Linux | Mac | Nodiadau |
| ------------------------------------------------------------- | ------- | ----- | --- | -------- |
| run() gwall                                                   | I       | I     | I   |          |
| destroy()                                                     |         | I     | I   |          |
| setApplicationMenu(menu \*Menu)                               | I       | I     | I   |          |
| name() llinell                                                |         | I     | I   |          |
| getCurrentWindowID() uint                                     | I       | I     | I   |          |
| showAboutDialog(name llinell, description llinell, icon []byte) |         | I     | I   |          |
| setIcon(icon []byte)                                          | -       | I     | I   |          |
| on(id uint)                                                   |         |       | I   |          |
| dispatchOnMainThread(fn func())                               | I       | I     | I   |          |
| cuddio()                                                        | I       | I     | I   |          |
| dangos()                                                        | I       | I     | I   |          |
| getPrimaryScreen() (\*Screen, gwall)                          |         | I     | I   |          |
| getScreens() ([]\*Screen, gwall)                              |         | I     | I   |          |

## Ffenestr Gwe-weld

Dulliau Rhyngwyneb Ffenestr Gwe-weld

| Dull                                             | Windows | Linux | Mac | Nodiadau                                    |
| -------------------------------------------------- | ------- | ----- | --- | ---------------------------------------- |
| canolbwyntio()                                           | I       | I     | I   |                                          |
| cau()                                            | i       | I     | I   |                                          |
| destroy()                                          |         | I     | I   |                                          |
| execJS(js llinell)                                  | i       | I     | I   |                                          |
| ffocws()                                            | I       | I     |     |                                          |
| forceReload()                                      |         | I     | I   |                                          |
| fullscreen()                                       | I       | I     | I   |                                          |
| getScreen() (\*Screen, gwall)                      | i       | I     | I   |                                          |
| getZoom() float64                                  |         | I     | I   |                                          |
| uchder() int                                       | I       | I     | I   |                                          |
| cuddio()                                             | I       | I     | I   |                                          |
| isFullscreen() bool                                | I       | I     | I   |                                          |
| isMaximised() bool                                 | I       | I     | I   |                                          |
| isMinimised() bool                                 | I       | I     | I   |                                          |
| mwyhau()                                         | I       | I     | I   |                                          |
| lleihau()                                         | I       | I     | I   |                                          |
| nativeWindowHandle() (uintptr, gwall)              | I       | I     | I   |                                          |
| on(eventID uint)                                   | i       |       | I   |                                          |
| openContextMenu(menu *Menu, data *ContextMenuData) | i       | I     | I   |                                          |
| positionsberthol() (int, int)                      | I       | I     | I   |                                          |
| ail-lwytho()                                           | i       | I     | I   |                                          |
| rhedeg()                                              | I       | I     | I   |                                          |
| setAlwaysOnTop(alwaysOnTop bool)                   | I       | I     | I   |                                          |
| setBackgroundColour(color RGBA)                    | I       | I     | I   |                                          |
| setEnabled(bool)                                   |         | I     | I   |                                          |
| setFrameless(bool)                                 |         | I     | I   |                                          |
| setFullscreenButtonEnabled(enabled bool)           | -       | I     | I   | Nid oes botwm sgrin lawn yn Windows |
| setHTML(html llinell)                               | I       | I     | I   |                                          |
| setMaxSize(width, uchder int)                      | I       | I     | I   |                                          |
| setMinSize(width, uchder int)                      | I       | I     | I   |                                          |
| setRelativePosition(x int, y int)                  | I       | I     | I   |                                          |
| setResizable(resizable bool)                       | I       | I     | I   |                                          |
| setSize(width, uchder int)                         | I       | I     | I   |                                          |
| setTitle(title llinell)                             | I       | I     | I   |                                          |
| setURL(url llinell)                                 | I       | I     | I   |                                          |
| setZoom(zoom float64)                              | I       | I     | I   |                                          |
| dangos()                                             | I       | I     | I   |                                          |
| maint() (int, int)                                  | I       | I     | I   |                                          |
| toggleDevTools()                                   | I       | I     | I   |                                          |
| un-fullscreen()                                     | I       | I     | I   |                                          |
| un-mwyhau()                                       | I       | I     | I   |                                          |
| un-lleihau()                                       | I       | I     | I   |                                          |
| lled() int                                        | I       | I     | I   |                                          |
| chwyddo()                                             |         | I     | I   |                                          |
| chwyddo()                                           | I       | I     | I   |                                          |
| chwyddo()                                          | I       | I     | I   |                                          |
| chwyddo()                                        | I       | I     | I   |                                          |

## Amser Gweithredol

### Cymhwyster

| Nodwedd | Windows | Linux | Mac | Nodiadau |
| ------- | ------- | ----- | --- | -------- |
| Gadael  | I       | I     | I   |          |
| Cuddio  | I       | I     | I   |          |
| Dangos  | I       |       | I   |          |

### Deialogau

| Nodwedd  | Windows | Linux | Mac | Nodiadau |
| -------- | ------- | ----- | --- | -------- |
| Gwybodaeth | I       | I     | I   |          |
| Rhybudd  | I       | I     | I   |          |
| Gwall    | I       | I     | I   |          |
| Cwestiwn | I       | I     | I   |          |
| OpenFile | I       | I     | I   |          |
| SaveFile | I       | I     | I   |          |

### Clipfwrdd

| Nodwedd | Windows | Linux | Mac | Nodiadau |
|---------|---------|-------|-----|----------|
| SetText | I       | I     | I   |          |
| Text    | I       | I     | I   |          |

### ContextMenu

| Nodwedd           | Windows | Linux | Mac | Nodiadau |
|------------------|---------|-------|-----|----------|
| OpenContextMenu  | I       | I     | I   |          |
| Ar Ddiofyn        |         |       |     |          |
| Rheoli drwy HTML  | I       |       |     |          |

Mae'r ddewislen cyd-destun rhagosodedig wedi'i galluogi'n rhagosodedig ar gyfer yr holl elfennau sydd â `contentEditable: true`, `<input>` neu `<textarea>` tagiau neu â'r arddull `--default-contextmenu: true` osodwyd. Bydd arddull `--default-contextmenu: show` bob amser yn dangos y ddewislen cyd-destun Mae arddull `--default-contextmenu: hide` bob amser yn cuddio'r ddewislen cyd-destun

Ni fydd unrhywbeth wedi'i nythu o dan tag â'r arddull `--default-contextmenu: hide` yn dangos y ddewislen cyd-destun oni bai ei fod yn cael ei osod yn benodol gyda `--default-contextmenu: show`.

### Sgriniau

| Nodwedd    | Windows | Linux | Mac | Nodiadau |
| ---------- | ------- | ----- | --- | -------- |
| GetAll     | I       | I     | I   |          |
| GetPrimary | I       | I     | I   |          |
| GetCurrent | I       | I     | I   |          |

### System

| Nodwedd    | Windows | Linux | Mac | Nodiadau |
| ---------- | ------- | ----- | --- | -------- |
| IsDarkMode |         |       | I   |          |

### Ffenestr

I = Cefnogir U = Heb ei brofi

- = Ddim ar gael

| Nodwedd              | Windows | Linux | Mac | Nodiadau                                                                                |
| ------------------- | ------- | ----- | --- | ------------------------------------------------------------------------------------ |
| Canolbwyntio        | I       | I     |     |                                                                                      |
| Ffocws               | I       | I     |     |                                                                                      |
| FullScreen          | I       | I     | I   |                                                                                      |
| GetZoom             | I       | I     | I   | Cael graddfa golwg gyfredol                                                               |
| Uchder              | I       | I     | I   |                                                                                      |
| Cuddio              | I       | I     | I   |                                                                                      |
| Mwyhau            | I       | I     | I   |                                                                                      |
| Lleihau            | I       | I     | I   |                                                                                      |
| PositionsbertholRhyngwladol    | I       | I     | I   |                                                                                      |
| Sgrin              | I       | I     | I   | Cael sgrin ar gyfer ffenestr                                                                |
| SetAlwaysOnTop      | I       | I     | I   |                                                                                      |
| SetBackgroundColour | I       | I     | I   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| SetEnabled          | I       | U     | -   | Gosod y ffenestr i fod wedi'i galluogi/analluogi                                                |
| SetMaxSize          | I       | I     | I   |                                                                                      |
| SetMinSize          | I       | I     | I   |                                                                                      |
| SetRelativePosition | I       | I     | I   |                                                                                      |
| SetResizable        | I       | I     | I   |                                                                                      |
| SetSize             | I       | I     | I   |                                                                                      |
| SetTitle            | I       | I     | I   |                                                                                      |
| SetZoom             | I       | I     | I   | Gosod graddfa golwg                                                                       |
| Dangos              | I       | I     | I   |                                                                                      |
| Maint                | I       | I     | I   |                                                                                      |
| UnFullscreen        | I       | I     | I   |                                                                                      |
| UnMaximise          | I       | I     | I   |                                                                                      |
| UnMinimise          | I       | I     | I   |                                                                                      |
| Lled               | I       | I     | I   |                                                                                      |
| ZoomIn              | I       | I     | I   | Cynyddu graddfa golwg                                                                  |
| ZoomOut             | I       | I     | I   | Gostwng graddfa golwg                                                                  |
| ZoomReset           | I       | I     | I   | Ailosod graddfa golwg                                                                     |

### Opsiynau Ffenestr

Mae 'I' yn y tabl isod yn dangos bod yr opsiwn wedi'i brofi ac yn cael ei gymhwyso wrth greu'r ffenestr. Mae 'X' yn dangos nad yw'r opsiwn yn cael ei gefnogi gan y platfform.

| Nodwedd                         | Windows | Linux | Mac | Nodiadau                                      |
|---------------------------------|---------|-------|-----|--------------------------------------------|
| AlwaysOnTop                     | I       | I     |     |                                            |
| BackgroundColour                | I       | I     |     |                                            |
| BackgroundType                  |         |       |     | Mae'n ymddangos bod Acrylic yn gweithio ond nid y lleill |
| CSS                             | I       | I     |     |                                            |
| DevToolsEnabled                 | I       | I     | I   |                                            |
| DisableResize                   | I       | I     |     |                                            |
| EnableDragAndDrop               |         | I     |     |                                            |
| EnableFraudulentWebsiteWarnings |         |       |     |                                            |
| Ffocws                         | I       | I     |     |                                            |
| Frameless                       | I       | I     |     |                                            |
| FullscreenButtonEnabled         | I       |       |     |                                            |
| Uchder                          | I       | I     |     |                                            |
| Hidden                          | I       | I     |     |                                            |
| HTML                            | I       | I     |     |                                            |
| JS                              | I       | I     |     |                                            |
| Mac                             | -       | -     |     |                                            |
| MaxHeight                       | I       | I     |     |                                            |
| MaxWidth                        | I       | I     |     |                                            |
| MinHeight                       | I       | I     |     |                                            |
| MinWidth                        | I       | I     |     |                                            |
| Enw                            | I       | I     |     |                                            |
| OpenInspectorOnStartup          |         |       |     |                                            |
| StartState                      | I       |       |     |                                            |
| Teitl                           | I       | I     |     |                                            |
| URL                             | I       | I     |     |                                            |
| Lled                           | I       | I     |     |                                            |
| Windows                         | I       | -     | -   |                                            |
| X                               | I