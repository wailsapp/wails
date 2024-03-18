# Cynllun Gweithredu

Mae'r cynllun gweithredu yn ddogfen fyw ac yn agored i newid. Os oes gennych unrhyw awgrymiadau, agorwch fater. Bydd gan bob cam bwysig gyfres o nodau yr ydym yn anelu at eu cyflawni. Mae'r rhain yn agored i newid.

## Materion Hysbys

- Nid yw cynhyrchu rhwymau ar gyfer dull sy'n mewnforio pecyn sydd â'r un enw â phecyn arall a fewnforiwyd yn cael ei gefnogi ar hyn o bryd.

## Camau Milltir Alpha  

### Presennol: Alpha 5

#### Nodau

Mae cylch Alpha 5 yn anelu at ddod â Linux i gydraddoldeb (Alpha 4) â'r platfformau eraill.

#### Sut Gallaf Helpu?

!!! note
Adroddwch unrhyw faterion a ganfyddwch gan ddefnyddio [y canllaw hwn](./getting-started/feedback.md).

- Profi'r cwbl ar Linux!

#### Statws

Enghreifftiau Linux:

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha5.csv") }}

## Camau Milltir Nesaf

## Alpha 6

## Camau Milltir Blaenorol

### Alpha 4 - Wedi'i Chwblhau 2024-02-01

#### Nodau

Mae cylch Alpha 4 yn anelu at ddarparu'r gorchmynion `dev` a `package`. 
Dylai'r gorchymyn `wails dev` wneud y canlynol:
- Adeiladu'r cais
- Cychwyn y cais
- Cychwyn y gweinydd datblygu blaen 
- Gwylio am newidiadau i'r cod cais ac ail-adeiladu/ailgychwyn yn ôl yr angen

Dylai'r gorchymyn `wails package` wneud y canlynol:
- Adeiladu'r cais
- Pecynnu'r cais mewn fformat penodol i'r platfform
  - Windows: Rhaglen weithredol safonol, Gosodwr NSIS
  - Linux: AppImage
  - MacOS: Rhaglen weithredol safonol, Pecyn App
- Cefnogi gwrthodiad y cod cais

- Hefyd, rydym am gael pob enghraifft yn gweithio ar Linux.

#### Sut Gallaf Helpu?

!!! note
    Adroddwch unrhyw faterion a ganfyddwch gan ddefnyddio [y canllaw hwn](./getting-started/feedback.md).


- Gosod y CLI gan ddefnyddio'r cyfarwyddiadau [yma](./getting-started/installation.md).
- Rhedeg `wails3 doctor` a sicrhau bod yr holl ddibyniaeth wedi'u gosod. 
- Cynhyrchu project newydd gan ddefnyddio `wails3 init`.

Profi'r gorchymyn `wails3 dev`:

- Rhedeg `wails3 dev` yn y cyfeiriadur project. Dylai redeg y cais mewn modd datblygu.
- Ceisiwch newid ffeiliau a sicrhau bod y cais yn cael ei ail-adeiladu a'i ailgychwyn.
- Rhedeg `wails3 dev -help` i weld yr opsiynau.
- Ceisiwch wahanol opsiynau a sicrhau eu bod yn gweithio fel y disgwylir.

Profi'r gorchymyn `wails3 package`:

- Rhedeg `wails3 package` yn y cyfeiriadur project.
- Gwiriwch fod y cais wedi'i becynnu'n gywir ar gyfer y platfform presennol.
- Rhedeg `wails3 package -help` i weld yr opsiynau.
- Ceisiwch wahanol opsiynau a sicrhau eu bod yn gweithio fel y disgwylir.

Adolygwch y tabl isod a chwilio am senarios heb eu profi. 
Yn y bôn, ceisiwch ei dorri a rhowch wybod i ni os ydych yn dod o hyd i unrhyw faterion! :smile:

#### Statws

Gorchymyn `wails3 dev`:

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha4-wails3-dev.csv") }}

- Mae Windows yn gweithio'n rhannol:
  - Mae newidiadau i'r blaen-ddalen yn gweithio fel y disgwylir
  - Mae newidiadau i Go yn achosi i'r cais gael ei adeiladu ddwywaith

- Mae Mac yn gweithio'n rhannol:
  - Mae newidiadau i'r blaen-ddalen yn gweithio fel y disgwylir
  - Mae newidiadau i Go yn gorffen y broses vite

Gorchymyn `wails3 package`:

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio
- :material-cancel: - Heb ei Gefnogi

{{ read_csv("alpha4-wails3-package.csv") }}


### Alpha 3 - Wedi'i Chwblhau 2024-01-14

#### Nodau

Mae cylch Alpha 3 yn anelu at ddarparu cefnogaeth rhwymau. Mae Wails 3 yn defnyddio dull dadansoddi státig newydd sy'n ein galluogi i ddarparu profiad rhwymau gwell nag yn Wails 2.
Hefyd, rydym am gael pob enghraifft yn gweithio ar Linux.

#### Sut Gallaf Helpu?

Gallwch gynhyrchu rhwymau gan ddefnyddio'r gorchymyn `wails3 generate bindings`. Bydd hyn yn cynhyrchu rhwymau ar gyfer holl fethodoedd strwythur a ranbir gyda'ch project.
Rhedeg `wails3 generate bindings -help` i weld yr opsiynau sy'n llywodraethu sut caiff rhwymau eu cynhyrchu.
 
Mae'r profion ar gyfer y cynhyrchwr rhwymau i'w canfod [yma](https://github.com/wailsapp/wails/tree/v3-alpha/v3/internal/parser) gyda'r data profion wedi'i leoli yn y cyfeiriadur `testdata`. 

Adolygwch y tabl isod a chwilio am senarios heb eu profi. Mae'r cod parser a'r profion wedi'u lleoli yn `v3/internal/parser`. Gellir rhedeg yr holl brofion gan ddefnyddio `go test ./...` o'r cyfeiriadur `v3`.
Yn y bôn, ceisiwch ei dorri a rhowch wybod i ni os ydych yn dod o hyd i unrhyw faterion! :smile:

#### Statws

Rhwymau ar gyfer strwythur (CallByID):

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol 
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha3-bindings-callbyid.csv") }}

Rhwymau ar gyfer strwythur (CallByName):

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha3-bindings-callbyname.csv") }}

Modelau:

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha3-models.csv") }}


Enghreifftiau:

- [ ] Pob enghraifft yn gweithio ar Linux

### Alpha 2

#### Nodau

Mae Alpha 2 yn anelu at gyflwyno cefnogaeth [Taskfile](https://taskfile.dev). Bydd hyn yn
caniatáu i ni gael system adeiladu unigol, estynnadwy sy'n gweithio ar bob platfform.
Hefyd, rydym am gael pob enghraifft yn gweithio ar Linux.

#### Statws

- [ ] Pob enghraifft yn gweithio ar Linux
- [x] Gorchmynion Cychwyn a Adeiladu


- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha2.csv") }}

### Alpha 1

#### Nodau

Mae Alpha 1 y rhyddhad cychwynnol. Mae'n fwriadedig i gael adborth ar yr API newydd
ac i bobl ddechreu arbrofi ag ef. Y nod pennaf yw cael y rhan fwyaf o'r
enghreifftiau yn gweithio ar bob platfform.

#### Statws

- :material-check-bold: - Yn gweithio
- :material-minus: - Yn gweithio'n rhannol
- :material-close: - Ddim yn gweithio

{{ read_csv("alpha1.csv") }}