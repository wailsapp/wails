# Amser Rhedeg

Mae amser rhedeg Wails yn llyfrgell safonol ar gyfer ceisiadau Wails. Mae'n darparu nifer o nodweddion y gellir eu defnyddio yn eich ceisiadau, gan gynnwys:

- Rheolaeth ffenestr
- Deialogau
- Integreiddio porwr
- Clipfwrdd
- Llusgian diframwaith
- Eiconau ardal gwaith
- Rheoli dewislen
- Gwybodaeth system
- Digwyddiadau
- Galw cod Go
- Dewislenni Cyd-destun
- Sgrîn
- WML (Iaith Marcio Wails)

Mae'r amser rhedeg yn ofynnol ar gyfer integreiddio rhwng Go a'r rhaglen blaen. Mae 2 ffordd o integreiddio'r amser rhedeg:

- Gan ddefnyddio'r pecyn `@wailsio/runtime`
- Gan ddefnyddio fersiwn wedi'i chyn-adeiladu o'r amser rhedeg

## Defnyddio'r pecyn `@wailsio/runtime`

Mae'r pecyn `@wailsio/runtime` yn becyn JavaScript sy'n darparu mynediad at amser rhedeg Wails. Fe'i defnyddir gan yr holl dempled safonol ac mae'n y ffordd a argymhellir i integreiddio'r amser rhedeg i'ch cais. Drwy ddefnyddio'r pecyn, dim ond y rhannau o'r amser rhedeg yr ydych yn eu defnyddio a gaiff eu cynnwys.

Mae'r pecyn ar gael ar npm a gellir ei osod gan ddefnyddio:

```shell
npm install --save @wailsio/runtime
```

## Defnyddio fersiwn wedi'i chyn-adeiladu o'r amser rhedeg

Bydd rhai prosiectau heb ddefnyddio pecynwr JavaScript ac efallai y byddant yn well ganddynt ddefnyddio fersiwn wedi'i chyn-adeiladu o'r amser rhedeg. Dyma'r rhagosodiad ar gyfer yr enghreifftiau yn `v3/examples`. Gellir cynhyrchu'r fersiwn wedi'i chyn-adeiladu o'r amser rhedeg gan ddefnyddio'r gorchymyn canlynol:

```shell
wails3 generate runtime
```

Bydd hyn yn cynhyrchu ffeil `runtime.js` (a `runtime.debug.js`) yn y cyfeiriadur presennol.
Gellir defnyddio'r ffeil hon gan eich cais drwy ei hychwanegu at eich cyfeiriadur asedau (fel arfer `frontend/dist`) ac yna ei chynnwys yn eich HTML:

```html
<html>
    <head>
        <script src="/runtime.js"></script>
    </head>
    <!--- ... -->
</>
```