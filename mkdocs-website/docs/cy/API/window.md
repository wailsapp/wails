# Ffenestr

I greu ffenestr, defnyddiwch
[Application.NewWebviewWindow](application.md#newwebviewwindow) neu
[Application.NewWebviewWindowWithOptions](application.md#newwebviewwindowwithoptions).
Mae'r cyntaf yn creu ffenestr gyda nodweddion rhagosodedig, tra bod yr olaf yn
caniatáu i chi bennu opsiynau wedi'u haddasu.

Mae'r dulliau hyn yn galladwy ar y gwrthrych WebviewWindow a ddychwelir:

### SetTitle

API: `SetTitle(teitl string) *WebviewWindow`

Mae'r dull hwn yn diweddaru teitl y ffenestr i'r llinyn a ddarperir. Mae'n dychwelyd
y gwrthrych WebviewWindow, gan ganiatáu i ddulliau gael eu cadwyn.

### Enw

API: `Enw() string`

Mae'r swyddogaeth hon yn dychwelyd enw'r WebviewWindow.

### SetSize

API: `SetSize(lled, uchder int) *WebviewWindow`

Mae'r dull hwn yn gosod maint y WebviewWindow i'r lled a'r uchder a ddarperir. Os
yw'r dimensiynau a ddarparwyd yn rhagori ar y cyfyngiadau, mae'n eu haddasu'n briodol.

### SetAlwaysOnTop

API: `SetAlwaysOnTop(b bool) *WebviewWindow`

Mae'r swyddogaeth hon yn gosod y ffenestr i aros ar y brig yn seiliedig ar y blaen
llinyn a ddarperir.

### Dangos

API: `Dangos() *WebviewWindow`

Mae'r dull `Dangos` yn cael ei ddefnyddio i wneud y ffenestr yn weladwy. Os
nad yw'r ffenestr yn rhedeg, mae'n gwahodd y dull `rhedeg` i ddechrau'r ffenestr
ac yna'n ei gwneud yn weladwy.

### Cuddio

API: `Cuddio() *WebviewWindow`

Mae'r dull `Cuddio` yn cael ei ddefnyddio i guddio'r ffenestr. Mae'n gosod y
statws cudd o'r ffenestr i wir ac yn lledu'r digwyddiad cuddio ffenestr.

### SetURL

API: `SetURL(s string) *WebviewWindow`

Mae'r dull `SetURL` yn cael ei ddefnyddio i osod URL y ffenestr i'r llinyn URL a ddarparwyd.

### SetZoom

API: `SetZoom(mewnosod float64) *WebviewWindow`

Mae'r dull `SetZoom` yn gosod lefel swm cynnwys y ffenestr i'r lefel mewnosod a ddarparwyd.

### GetZoom

API: `GetZoom() float64`

Mae'r swyddogaeth `GetZoom` yn dychwelyd y lefel swm bresennol o gynnwys y ffenestr.

### GetScreen

API: `GetScreen() (*Screen, error)`

Mae'r dull `GetScreen` yn dychwelyd y sgrin lle mae'r ffenestr yn cael ei harddangos.

### SetFrameless

API: `SetFrameless(frameless bool) *WebviewWindow`

Mae'r swyddogaeth hon yn cael ei defnyddio i dynnu'r ffrâm a bar teitl y ffenestr.
Mae'n toglo'r framelessness o'r ffenestr yn unol â'r gwerth boolean a ddarperir
(gwir ar gyfer frameless, ffug ar gyfer ffrâm).

### RegisterContextMenu

API: `RegisterContextMenu(enw string, dewislen *Dewislen)`

Mae'r swyddogaeth hon yn cael ei defnyddio i gofrestru dewislen cyd-destun ac
yn ei neilltuo i'r enw a ddarparwyd.

### NativeWindowHandle

API: `NativeWindowHandle() (uintptr, error)`

Mae'r swyddogaeth hon yn cael ei defnyddio i nodi'r handlen ffenestr brodorol
ar gyfer y ffenestr.

### Ffocws

API: `Ffocws()`

Mae'r swyddogaeth hon yn cael ei defnyddio i ffocysu'r ffenestr.

### SetEnabled

API: `SetEnabled(galluogwyd bool)`

Mae'r swyddogaeth hon yn cael ei defnyddio i alluogi/analluogi'r ffenestr yn
seiliedig ar y gwerth boolean a ddarperir.

### SetAbsolutePosition

API: `SetAbsolutePosition(x int, y int)`

Mae'r swyddogaeth hon yn gosod y safle absoliwt o'r ffenestr yn y sgrin.