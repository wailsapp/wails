# Ardal Hysbysu

Mae'r ardal hysbysu yn cynnwys ardal hysbysu ar amgylchedd bwrdd gwaith, a all
gynnwys eiconau o'r rhaglenni sy'n rhedeg ar hyn o bryd a hysbysiadau system
penodol.

Rydych yn creu ardal hysbysu trwy alw `app.NewSystemTray()`:

```go
    // Creu ardal hysbysu newydd
tray := app.NewSystemTray()
```

Mae'r dulliau canlynol ar gael ar y `SystemTray` math:

### SetLabel

API: `SetLabel(label string)`

Mae'r dull `SetLabel` yn gosod label yr ardal hysbysu.

### Label

API: `Label() string`

Mae'r dull `Label` yn adfer label yr ardal hysbysu.

### PositionWindow

API: `PositionWindow(*WebviewWindow, offset int) error`

Mae'r dull `PositionWindow` yn galw'r dulliau `AttachWindow` a `WindowOffset`.

### SetIcon

API: `SetIcon(icon []byte) *SystemTray`

Mae'r dull `SetIcon` yn gosod eicon yr ardal hysbysu system.

### SetDarkModeIcon

API: `SetDarkModeIcon(icon []byte) *SystemTray`

Mae'r dull `SetDarkModeIcon` yn gosod eicon yr ardal hysbysu system pan mewn modd tywyll.

### SetMenu

API: `SetMenu(menu *Menu) *SystemTray`

Mae'r dull `SetMenu` yn gosod dewislen yr ardal hysbysu.

### Destroy

API: `Destroy()`

Mae'r dull `Destroy` yn dinistrio'r enghraifft ardal hysbysu.

### OnClick

API: `OnClick(handler func()) *SystemTray`

Mae'r dull `OnClick` yn gosod y swyddogaeth i'w gweithredu pan fo'r eicon ardal hysbysu wedi'i glicio.

### OnRightClick

API: `OnRightClick(handler func()) *SystemTray`

Mae'r dull `OnRightClick` yn gosod y swyddogaeth i'w gweithredu pan fo'r eicon ardal hysbysu wedi'i glicio â'r dde.

### OnDoubleClick

API: `OnDoubleClick(handler func()) *SystemTray`

Mae'r dull `OnDoubleClick` yn gosod y swyddogaeth i'w gweithredu pan fo'r eicon ardal hysbysu wedi'i glicio ddwywaith.

### OnRightDoubleClick

API: `OnRightDoubleClick(handler func()) *SystemTray`

Mae'r dull `OnRightDoubleClick` yn gosod y swyddogaeth i'w gweithredu pan fo'r eicon ardal hysbysu wedi'i glicio ddwywaith â'r dde.

### AttachWindow

API: `AttachWindow(window *WebviewWindow) *SystemTray`

Mae'r dull `AttachWindow` yn atodi ffenestr i'r ardal hysbysu system. Bydd y ffenestr yn cael ei dangos pan fo'r eicon ardal hysbysu wedi'i glicio.

### WindowOffset

API: `WindowOffset(offset int) *SystemTray`

Mae'r dull `WindowOffset` yn gosod y bwlch mewn picselau rhwng yr ardal hysbysu system a'r ffenestr.

### WindowDebounce

API: `WindowDebounce(debounce time.Duration) *SystemTray`

Mae'r dull `WindowDebounce` yn gosod amser diddymu. Yng nghyd-destun Windows, defnyddir hyn i bennu faint o amser i aros cyn ymateb i ddigwyddiad clic llygoden i fyny ar yr eicon hysbysu.

### OpenMenu

API: `OpenMenu()`

Mae'r dull `OpenMenu` yn agor y ddewislen sy'n gysylltiedig â'r ardal hysbysu system.