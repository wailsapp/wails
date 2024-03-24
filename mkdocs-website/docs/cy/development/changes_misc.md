### Misc

## Opsiynau Cymhwyso Windows

### WndProcInterceptor

Os caiff hwn ei osod, bydd WndProc yn cael ei ychwanegu ac fe gaiff y swyddogaeth ei galw.
Mae hyn yn caniatáu i chi ddelio â negeseuon Windows yn uniongyrchol. Dylai'r swyddogaeth
fod â'r llofnod canlynol:

```go
func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnValue uintptr, shouldReturn)
```

Dylid gosod y gwerth `shouldReturn` i `true` os dylai'r `returnValue` gael ei
ddychwelyd gan y prif ddull wndProc. Os caiff ei osod i `false`, bydd y gwerth
dychwelyd yn cael ei anwybyddu a bydd y neges yn parhau i gael ei phrosesu gan y prif
ddull wndProc.

## Cuddio'r Ffenestr wrth Gau + OnBeforeClose

Yn v2, roedd y fflag `HideWindowOnClose` i guddio'r ffenestr pan gaiff ei chau.
Roedd gorgyffwrdd rhesymegol rhwng y fflag hon a'r galwad `OnBeforeClose`.
Yn v3, mae'r fflag `HideWindowOnClose` wedi'i thynnu ac mae'r galwad `OnBeforeClose`
wedi'i ailenwi i `ShouldClose`. Caiff y galwad `ShouldClose` ei galw pan fydd y
defnyddiwr yn ceisio cau ffenestr. Os bydd y galwad yn dychwelyd `true`, caiff y
ffenestr ei chau. Os yw'n dychwelyd `false`, ni chaiff y ffenestr ei chau. Gellir
ei ddefnyddio i guddio'r ffenestr yn hytrach na'i chau.

## Llusgo Ffenestr

Yn v2, defnyddiwyd yr ymddangosiad `--wails-drag` i nodi y gallai elfen gael ei
defnyddio i lusgo'r ffenestr. Yn v3, mae hwn wedi'i ddisodli gan `--webkit-app-region`
i fod yn fwy yn unol â'r ffordd y mae fframweithiau eraill yn ymdrin â hyn. Gellir
gosod yr ymddangosiad `--webkit-app-region` i unrhyw un o'r gwerthoedd canlynol:

- `drag` - Gellir defnyddio'r elfen i lusgo'r ffenestr
- `no-drag` - Ni ellir defnyddio'r elfen i lusgo'r ffenestr

Byddem wedi hoffi defnyddio `app-region`, fodd bynnag, nid yw hwn yn cael ei
gefnogi gan yr alwad `getComputedStyle` ar webkit ar macOS.