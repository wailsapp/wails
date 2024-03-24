Dyma'r cyfieithiad Cymraeg (cym) o'r testun Saesneg:

# Cyflwyniad

!!! note
    Mae'r canllaw hwn yn gweithio ymlaen.

Diolch am ddymuno helpu gyda datblygu Wails! Bydd y canllaw hwn yn eich helpu i
gychwyn.

## Cychwyn

- Cloniwch y storfa Git. Symudwch i'r gangen `v3-alpha`.
- Gosodwch y CLI: `cd v3/cmd/wails3 && go install`

- Dewisol: Os ydych am ddefnyddio'r system adeiladu i adeiladu cod blaen, bydd
  angen i chi osod [npm](https://nodejs.org/en/download).

## Adeiladu

Ar gyfer rhaglenni syml, gallwch ddefnyddio'r gorchymyn `go build` safonol. Mae
modd defnyddio `go run` hefyd.

Mae Wails hefyd yn cynnwys system adeiladu y gellir ei defnyddio i adeiladu
prosiectau mwy cymhleth. Mae'n defnyddio'r system adeiladu wych
[Task](https://taskfile.dev). Am fwy o wybodaeth, gwiriwch y dudalen gartref Task
neu runnwch `wails task --help`.

## Cynllun prosiect

Mae'r prosiect yn cael y strwythur canlynol:

    ```
    v3
    ├── cmd/wails3                  // CLI
    ├── examples                   // Enghreifftiau o apiau Wails
    ├── internal                   // Pecynnau mewnol
    |   ├── runtime                // Y runtime JS Wails
    |   └── templates              // Y templed prosiect a gynhelir
    ├── pkg
    |   ├── application            // Y llyfrgell Wails craidd
    |   └── events                 // Diffiniadau digwyddiadau
    |   └── mac                    // Cod penodol i macOS a ddefnyddir gan addasiadau
    |   └── w32                    // Cod penodol i Windows
    ├── plugins                    // Addasiadau a gynhelir
    ├── tasks                      // Tasgau cyffredinol
    └── Taskfile.yaml              // Ffurfweddiad tasgau datblygu
    ```

## Datblygu

### Rhestr Tasgau Alpha

Rydym yn monitro materion hysbys a thasgau ar hyn o bryd yn y
[Rhestr Tasgau Alpha](https://github.com/orgs/wailsapp/projects/6). Os ydych am
helpu, edrychwch ar y rhestr hon a dilyn y cyfarwyddiadau yn y
[Adborth](../getting-started/feedback.md) tudalen.

### Ychwanegu swyddogaeth ffenestr

Y ffordd well o ychwanegu swyddogaeth ffenestr yw ychwanegu swyddogaeth newydd i'r
ffeil `pkg/application/webview_window.go`. Dylai hon weithredu'r holl
swyddogaeth sydd ei hangen ar gyfer pob platfform. Dylid galw unrhyw god platfform
penodol drwy fethôd `webviewWindowImpl`. Gweithredir y rhyngwyneb hwn gan bob un
o'r platfformau targed i ddarparu'r swyddogaeth benodol i'r platfform. Mewn rhai
achosion, efallai na fydd yn gwneud dim. Ar ôl ychwanegu'r dull rhyngwyneb,
sicrhewch fod pob platfform yn ei weithredu. Mae'r dull `SetMinSize` yn enghraifft
dda o hyn.

- Mac: `webview_window_darwin.go`
- Windows: `webview_window_windows.go`
- Linux: `webview_window_linux.go`

Dylai'r rhan fwyaf, os nad y cyfan, o'r cod platfform penodol gael ei redeg ar y
prif drywydd. Er mwyn symleiddio hyn, mae nifer o ddulliau `invokeSync` wedi'u
diffinio yn `application.go`.

### Diweddaru'r runtime

Mae'r runtime wedi'i leoli yn `v3/internal/runtime`. Pan diweddarir y runtime,
rhaid cymryd y camau canlynol:

```shell
wails3 task runtime:build
```

### Digwyddiadau

Diffinnir digwyddiadau yn `v3/pkg/events`. Wrth ychwanegu digwyddiad newydd, rhaid
cymryd y camau canlynol:

- Ychwanegu'r digwyddiad i'r ffeil `events.txt`
- Rhedeg `wails3 task events:generate`

Mae nifer o fathau o ddigwyddiadau: digwyddiadau platfform penodol i'r ap a'r
ffenestr + digwyddiadau cyffredin. Mae'r digwyddiadau cyffredin yn ddefnyddiol ar
gyfer trin digwyddiadau ar draws platfformau, ond nid ydych wedi'ch cyfyngu i'r "isaf
cyffredin". Gallwch ddefnyddio'r digwyddiadau platfform penodol os oes angen i chi.

Wrth ychwanegu digwyddiad cyffredin, sicrhewch fod y digwyddiadau platfform penodol
wedi'u mapio. Mae enghraifft o hyn yn `window_webview_darwin.go`:

```go
		// Translate ShouldClose to common WindowClosing event
		w.parent.On(events.Mac.WindowShouldClose, func(_ *WindowEventContext) {
			w.parent.emit(events.Common.WindowClosing)
		})
```

NODYN: Efallai y byddwn yn ceisio awtomeiddio hyn yn y dyfodol drwy ychwanegu'r
mapio at y diffiniad digwyddiad.

### Addasiadau

Mae addasiadau yn ffordd o estyn swyddogaeth eich ap Wails.

#### Creu addasiad

Mae addasiadau yn strwythur Go safonol sy'n cydymffurfio â'r rhyngwyneb canlynol:

```go
type Plugin interface {
    Name() string
    Init(*application.App) error
    Shutdown()
    CallableByJS() []string
    InjectJS() string
}
```

Mae'r dull `Name()` yn dychwelyd enw'r addasiad. Defnyddir hwn at ddibenion
cofnodi.

Mae'r dull `Init(*application.App) error` yn cael ei alw pan gaiff yr addasiad ei
lwytho. Mae'r paramedr `*application.App` yn yr ap y caiff yr addasiad ei lwytho
iddo. Bydd unrhyw wallau yn atal yr ap rhag cychwyn.

Gelwir y dull `Shutdown()` pan gaiff yr ap ei ddidoli.

Mae'r dull `CallableByJS()` yn dychwelyd rhestr o swyddogaethau alladwy y gellir
eu galw o'r blaen. Rhaid i enwau'r dulliau hyn gwatsh yn union ag enwau'r dulliau
a alluogeir gan yr addasiad.

Mae'r dull `InjectJS()` yn dychwelyd JavaScript y dylid ei fewnosod i bob ffenestr
wrth iddynt gael eu creu. Mae hyn yn ddefnyddiol ar gyfer ychwanegu swyddogaethau
JavaScript cyfaddas i'r addasiad.

Cewch hyd i'r addasiadau mewnol yn y cyfeiriadur `v3/plugins`. Edrychwch arnynt am
ysbrydoliaeth.

## Tasgau

Mae'r CLI Wails yn defnyddio'r system adeiladu [Task](https://taskfile.dev). Fe'i
mewnforiwyd fel llyfrgell a'i ddefnyddio i redeg y tasgau a ddiffinnir yn
`Taskfile.yaml`. Y prif gyswllt â Task ddigwydd yn `v3/internal/commands/task.go`.

### Uwchraddio Taskfile

I wirio a oes diweddariad ar gyfer Taskfile, rhedwch `wails3 task -version` a
gwiriwch yn erbyn gwefan Task.

I uwchraddio'r fersiwn o Taskfile a ddefnyddir, rhedwch:

```shell
wails3 task taskfile:upgrade
```

Os oes anghydnawsedd, dylai'r rhain ymddangos yn y ffeil
`v3/internal/commands/task.go`.

Fel arfer, y ffordd orau o drwsio anghydnawsedd yw clonio'r storfa dasg yn
`https://github.com/go-task/task` a gwirio hanes y git i benderfynu beth sydd
wedi newid a pham.

I wirio bod yr holl newidiadau wedi gweithio'n gywir, ail-osodwch y CLI a gwirio'r
fersiwn eto:

```shell
wails3 task cli:install
wails3 task -version
```

## Agor PR

Gwnewch yn siŵr bod gan bob PR docyn cysylltiedig â nhw sy'n darparu cyd-destun y
newid. Os nad oes tocyn, crëwch un yn gyntaf. Sicrhewch fod pob PR wedi
diweddaru'r ffeil CHANGELOG.md gyda'r newidiadau a wnaed. Mae'r ffeil CHANGELOG.md
wedi'i lleoli yn y cyfeiriadur `mkdocs-website/docs`.

## Tasgau Amrywiol

### Uwchraddio Taskfile

Mae'r CLI Wails yn defnyddio'r system adeiladu [Task](https://taskfile.dev). Fe'i
mewnforiwyd fel llyfrgell a'i ddefnyddio i redeg y tasgau a ddiffinnir yn
`Taskfile.yaml`. Y prif gyswllt â Task ddigwydd yn `v3/internal/commands/task.go`.

I wirio a oes diweddariad ar gyfer Taskfile, rhedwch `wails3 task -version` a
gwiriwch yn erbyn gwefan Task.

I uwchraddio'r fersiwn o Taskfile a ddefnyddir, rhedwch:

```shell
wails3 task taskfile:upgrade
```

Os oes anghydnawsedd, dylai'r rhain ymddangos yn y ffeil
`v3/internal/commands/task.go`.

Fel arfer, y ffordd orau o drwsio anghydnawsedd yw clonio'r storfa dasg yn
`https://github.com/go-task/task` a gwirio hanes y git i benderfynu beth sydd
wedi newid a pham.

I wirio bod yr holl newidiadau wedi gweithio'n gywir, ail-osodwch y CLI a gwirio'r
fersiwn eto:

```shell
wails3 task cli:install
wails3 task -version
```