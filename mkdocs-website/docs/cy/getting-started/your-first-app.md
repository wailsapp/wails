# Dy Cymhwysiad Cyntaf

Mae creu eich cymhwysiad cyntaf gyda Wails v3 Alpha yn daith gyffrous i mewn i fyd datblygu apiau bwrdd gwaith modern. Bydd y canllaw hwn yn mynd â chi drwy'r broses o greu cymhwysiad sylfaenol, gan ddangos pŵer a symlrwydd Wails.

## Gofynion Rhewydd

Cyn dechrau, sicrhewch eich bod wedi gosod y canlynol:

- Go (fersiwn 1.21 neu ddiweddarach)
- Node.js (fersiwn LTS)
- Wails v3 Alpha (gweler y [canllaw gosod](installation.md) am gyfarwyddiadau)

## Cam 1: Creu Prosiect Newydd

Agorwch eich terfynell a rhedeg y gorchymyn canlynol i greu prosiect Wails newydd:

`wails3 init -n myfirstapp`

Mae'r gorchymyn hwn yn creu cyfeiriadur newydd o'r enw `myfirstapp` gyda'r holl ffeiliau angenrheidiol.

## Cam 2: Archwilio Strwythur y Prosiect

Ewch i'r cyfeiriadur `myfirstapp`. Byddwch yn canfod nifer o ffeiliau a ffolderi:

- `build`: Yn cynnwys ffeiliau a ddefnyddir gan y broses adeiladu.
- `frontend`: Yn cynnwys cod rhagflaen eich rhyngrwyd.
- `go.mod` a `go.sum`: Ffeiliau modiwl Go.
- `main.go`: Pwynt mynediad eich cymhwysiad Wails.
- `Taskfile.yml`: Yn diffinio'r holl dasgau a ddefnyddir gan y system adeiladu. Dysgu rhagor ar wefan [Task](https://taskfile.dev/).

Cymerwch ennyd i archwilio'r ffeiliau hyn a'ch cyfarwyddo â'r strwythur.

!!! note
    Er bod Wails v3 yn defnyddio [Task](https://taskfile.dev/) fel ei system adeiladu ddiofyn, does dim byd yn atal chi rhag defnyddio `make` neu unrhyw system adeiladu amgen.  

## Cam 3: Adeiladu Eich Cymhwysiad

I adeiladu eich cymhwysiad, rhedwch:

`wails3 build`

Mae'r gorchymyn hwn yn cyfansoddi fersiwn dadfygio o'ch cymhwysiad ac yn ei gadw mewn cyfeiriadur `bin` newydd. 
Gallwch ei redeg fel unrhyw gymhwysiad arferol:

=== "Mac"

    `./bin/myfirstapp`

=== "Windows"

    `bin\myfirstapp.exe`

=== "Linux"

    `./bin/myfirstapp`

Byddwch yn gweld rhyngwyneb defnyddiwr syml, pwynt cychwyn eich cymhwysiad. Gan ei fod yn fersiwn dadfygio, byddwch hefyd yn gweld logiau yn y ffenestr gonsol. Mae hyn yn ddefnyddiol at ddibenion dadfygio.

## Cam 4: Modd Datblygu

Gallwn hefyd redeg y cymhwysiad yn y modd datblygu. Mae'r modd hwn yn caniatáu i chi wneud newidiadau i'ch cod rhagflaen a gweld y newidiadau'n cael eu hadlewyrchu yn y cymhwysiad sy'n rhedeg heb orfod ailadeiladu'r cyfan.

1. Agorwch ffenestr derfynell newydd.
2. Rhedwch `wails3 dev`.
3. Agorwch `frontend/main.js`.
4. Newid y llinell sydd â `<h1>Hello Wails!</h1>` i `<h1>Helo Byd!</h1>`.
5. Cadwch y ffeil.

Bydd y cymhwysiad yn diweddaru'n awtomatig, a byddwch yn gweld y newidiadau'n cael eu hadlewyrchu yn y cymhwysiad sy'n rhedeg. 

## Cam 5: Ailadeiladu'r Cymhwysiad

Pan fyddwch yn hapus gyda'ch newidiadau, ailadeiladu'r cymhwysiad eto:

`wails3 build`

Byddwch yn sylwi bod yr amser adeiladu wedi bod yn gyflymach y tro hwn. Mae hynny oherwydd bod y system adeiladu newydd yn unig yn adeiladu'r rhannau o'ch cymhwysiad sydd wedi newid.

Dylech weld gweithrediannol newydd yn y cyfeiriadur `build`.

## Casgliad

Llongyfarchiadau! Rydych newydd greu ac adeiladu eich cymhwysiad Wails cyntaf. Dyma ddechrau'r hyn y gallwch ei gyflawni gyda Wails v3 Alpha. Archwiliwch y ddogfennaeth, profwch y gwahanol nodweddion, a dechrau adeiladu apiau rhyfeddol!