## Ffenestr

Mae'r API Ffenestr wedi aros yn yr un fath i raddau helaeth, fodd bynnag mae'r dulliau yn awr ar enghraifft o ffenestr yn hytrach na'r amser gweithredu. Rhai gwahaniaeth nodedig yw:

- Mae gan Ffenestri nawr Enw sy'n eu hadnabod. Defnyddir hyn i adnabod y ffenestr wrth yrru digwyddiadau.
- Mae gan Ffenestri lawer mwy o ddulliau ar y rhai nad oeddent ar gael o'r blaen, fel `AbsolutePosition` a `ToggleDevTools`.
- Gall Ffenestri nawr dderbyn ffeiliau drwy lusgo a gollwng brodorol. Gweler yr adran Lusgo a Gollwng am fwy o fanylion.

### ColourCefndir

Yn v2, roedd hwn yn bwynt i strwythur `RGBA`. Yn v3, mae hwn yn werthhRGBA` strwythur.

### FfenestrnynTranslucent

Mae'r fflach hon wedi'i thynnu. Erbyn hyn mae gan `BackgroundType` fflach y gellir ei defnyddio i osod y math o gefndir y dylai'r ffenestr ei chael. Gellir gosod y fflach hon i unrhyw un o'r gwerthoedd canlynol:

- `BackgroundTypeSolid` - Bydd gan y ffenestr gefndir solet
- `BackgroundTypeTransparent` - Bydd gan y ffenestr gefndir tryloyw
- `BackgroundTypeTranslucent` - Bydd gan y ffenestr gefndir trawslucent

Ar Windows, os yw'r `BackgroundType` wedi'i osod i `BackgroundTypeTranslucent`, gellir gosod y math o drawslucedd gan ddefnyddio'r fflach `BackdropType` yn opsiynau `WindowsWindow`. Gellir gosod hon i unrhyw un o'r gwerthoedd canlynol:

- `Auto` - Bydd y ffenestr yn defnyddio effaith a benderfynir gan y system
- `None` - Ni fydd gan y ffenestr gefndir
- `Mica` - Bydd y ffenestr yn defnyddio'r effaith Mica
- `Acrylic` - Bydd y ffenestr yn defnyddio'r effaith acrylig
- `Tabbed` - Bydd y ffenestr yn defnyddio'r effaith tabbed