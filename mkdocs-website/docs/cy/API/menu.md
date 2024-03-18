# Dewislen

Gellir creu a chynnwys dewislenni yn y rhaglen. Gellir eu defnyddio i greu
dewislenni cyd-destun, dwylo system a dewislenni rhaglen.

I greu dewislen newydd, galwch:

```go
    // Creu dewislen newydd
    dewislen := app.NewMenu()
```

Mae'r gweithrediadau canlynol ar gael ar y `Dewislen` math:

### Ychwanegu

API: `Ychwanegu(label string) *EitemDewislen`

Mae'r dull hwn yn cymryd `label` o fath `string` fel mewnbwn ac yn ychwanegu
`EitemDewislen` newydd gyda'r label a roddir at y ddewislen. Mae'n dychwelyd yr
`EitemDewislen` a ychwanegwyd.

### YchwaneguSeparwr

API: `YchwaneguSeparwr()`

Mae'r dull hwn yn ychwanegu `EitemDewislen` gwahanol newydd at y ddewislen.

### YchwaneguBlwch

API: `YchwaneguBlwch(label string, galluogedig bool) *EitemDewislen`

Mae'r dull hwn yn cymryd `label` o fath `string` a `galluogedig` o fath `bool`
fel mewnbwn ac yn ychwanegu `EitemDewislen` blwch ticio newydd gyda'r label a'r
cyflwr galluogedig a roddir at y ddewislen. Mae'n dychwelyd yr `EitemDewislen`
a ychwanegwyd.

### YchwaneguRadio

API: `YchwaneguRadio(label string, galluogedig bool) *EitemDewislen`

Mae'r dull hwn yn cymryd `label` o fath `string` a `galluogedig` o fath `bool`
fel mewnbwn ac yn ychwanegu `EitemDewislen` radio newydd gyda'r label a'r
cyflwr galluogedig a roddir at y ddewislen. Mae'n dychwelyd yr `EitemDewislen`
a ychwanegwyd.

### Diweddaru

API: `Diweddaru()`

Mae'r dull hwn yn prosesu unrhyw grwpiau radio ac yn diweddaru'r ddewislen os
na chaiff y rhyngwyneb dewislen ei gychwyn.

### YchwaneguIsddewislen

API: `YchwaneguIsddewislen(s string) *Dewislen`

Mae'r dull hwn yn cymryd `s` o fath `string` fel mewnbwn ac yn ychwanegu
`EitemDewislen` isddewislen newydd gyda'r label a roddir at y ddewislen. Mae'n
dychwelyd yr isddewislen a ychwanegwyd.

### YchwaneguRôl

API: `YchwaneguRôl(rôl Rôl) *Dewislen`

Mae'r dull hwn yn cymryd `rôl` o fath `Rôl` fel mewnbwn, yn ei ychwanegu at y
ddewislen os nad yw'n `nil` ac yn dychwelyd y `Dewislen`.

### SetLabel

API: `SetLabel(label string)`

Mae'r dull hwn yn gosod `label` y `Dewislen`.