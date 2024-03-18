# Gosod

I osod yr Wails CLI, sicrhewch eich bod wedi gosod [Go 1.21+](https://go.dev/dl/)
a rhedwch:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3-alpha
cd v3/cmd/wails3
go install
```

## Platfformau a Gefnogir

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Ubuntu 22.04 AMD64/ARM64 (gall dosbarthiadau Linux eraill weithio hefyd!)

## Dibyniaeth

Mae gan Wails nifer o ddibyniaeth cyffredin sydd eu hangen cyn gosod:

=== "Go 1.21+"

    Lawrlwythwch Go o'r [Dudalen Lawrlwytho Go](https://go.dev/dl/).

    Sicrhewch eich bod yn dilyn y [cyfarwyddiadau gosod Go swyddogol](https://go.dev/doc/install). Bydd angen i chi hefyd sicrhau bod eich amrywiol amgylchedd `PATH` yn cynnwys llwybr eich cyfeiriadur `~/go/bin`. Ailgychwynnwch eich terfynell a gwiriwch y canlynol:

    - Gwiriwch fod Go wedi'i osod yn gywir: `go version`
    - Gwiriwch fod `~/go/bin` yn eich `PATH`: `echo $PATH | grep go/bin`

=== "npm (Dewisol)"

    Er nad oes angen npm i'w osod ar Wails, mae ei angen os ydych am ddefnyddio'r templed sydd wedi'i gynnwys.

    Lawrlwythwch y gosodwr node diweddaraf o'r [Dudalen Lawrlwytho Node](https://nodejs.org/en/download/). Mae'n well defnyddio'r rhyddhad diweddaraf gan mai hwnnw yr ydym fel arfer yn ei brofi yn erbyn.

    Rhedwch `npm --version` i wirio.

=== "Task (Dewisol)"

    Mae gan yr Wails CLI rhedwr tasg wedi'i ymgorffori o'r enw [Task](https://taskfile.dev/#/installation). Mae'n ddewisol, ond fe'i argymhellir. Os nad ydych am osod Task, gallwch ddefnyddio'r gorchymyn `wails3 task` yn lle `task`.
    Bydd gosod Task yn rhoi'r hyblygrwydd mwyaf i chi.

## Dibyniaeth Penodol i'r Platfform

Bydd angen i chi hefyd osod dibyniaeth penodol i'r platfform:

=== "Mac"

    Mae angen i Wails fod â'r offer llinell orchymyn xcode wedi'u gosod. Gellir gwneud hyn trwy redeg:

    ```
    xcode-select --install
    ```

=== "Windows"

    Mae angen i Wails fod â [Rhedwr WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) wedi'i osod. Bydd rhai gosodiadau Windows eisoes wedi'i gael hwn wedi'i osod. Gallwch wirio hyn gan ddefnyddio'r gorchymyn `wails doctor`.

=== "Linux"

    Mae angen offer adeiladu `gcc` safonol a `libgtk3` a `libwebkit` ar Linux. Yn hytrach nag rhestru llawer o orchmynion ar gyfer gwahanol ddosbarthiadau, gall Wails geisio penderfynu beth yw'r gorchmynion gosod ar gyfer eich dosbarthiad penodol. Rhedwch <code>wails doctor</code> ar ôl gosod i weld y cyfarwyddiadau ar sut i osod y dibyniaeth. Os na chefnogir eich dosbarthiad/rheolwr pecyn, rhowch wybod i ni ar discord.

## Gwirio'r System

Bydd rhedeg `wails3 doctor` yn gwirio a oes gennych y dibyniaeth cywir
wedi'i gosod. Os nad oes, bydd yn rhoi cyngor ar yr hyn sydd ar goll a sut i
unioni
unrhyw broblemau.

## Mae'r gorchymyn `wails3` yn ymddangos fel ei fod ar goll?

Os yw eich system yn adrodd bod y gorchymyn `wails3` ar goll, gwiriwch y
canlynol:

- Sicrhewch eich bod wedi dilyn y cyfarwyddiadau gosod Go yn gywir.
- Gwiriwch fod cyfeiriadur `go/bin` yn y newidyn amgylchedd `PATH`.
- Cau/Agor terminal presennol i ddewis y `PATH` newydd.