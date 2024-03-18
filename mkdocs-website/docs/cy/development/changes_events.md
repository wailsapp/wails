## Digwyddiadau

Yn v3, mae 3 math o ddigwyddiadau:

- Digwyddiadau Cymhwysiad
- Digwyddiadau Ffenestr
- Digwyddiadau Cyfaddas

### Digwyddiadau Cymhwysiad

Mae digwyddiadau cymhwysiad yn ddigwyddiadau a allbynir gan y cymhwysiad. Mae'r digwyddiadau hyn yn cynnwys digwyddiadau brodorol fel `ApplicationDidFinishLaunching` ar macOS.

### Digwyddiadau Ffenestr

Mae digwyddiadau ffenestr yn ddigwyddiadau a allbynir gan ffenestr. Mae'r digwyddiadau hyn yn cynnwys digwyddiadau brodorol fel `WindowDidBecomeMain` ar macOS. Diffinnir digwyddiadau cyffredin hefyd, fel y maent yn gweithio ar draws platfformau, e.e. `WindowClosing`.

### Digwyddiadau Cyfaddas

Mae'r digwyddiadau y mae'r defnyddiwr yn eu diffinio yn cael eu galw `WailsEvents`. Mae hyn er mwyn eu gwahaniaethu o'r gwrthrych `Event` a ddefnyddir i gyfathrebu gyda'r porwr. Mae WailsEvents bellach yn wrthrychau sy'n erynu holl fanylion digwyddiad. Mae hyn yn cynnwys enw'r digwyddiad, y data, a ffynhonnell y digwyddiad.

Mae'r data sy'n gysylltiedig â WailsEvent bellach yn un gwerth. Os oes angen mwy nag un gwerth, gellir defnyddio strwythur.

### Galwadau digwyddiad a llofnod swyddogaeth `Emit`

Mae llofnodion y galwadau digwyddiad (fel y defnyddir gan `On`, `Once` & `OnMultiple`) wedi newid. Yn v2, naeth y swyddogaeth alwad dderbyn data dewisol. Yn v3, mae'r swyddogaeth alwad yn derbyn gwrthrych `WailsEvent` sy'n cynnwys yr holl ddata sy'n berthnasol i'r digwyddiad.

Yn yr un modd, mae'r swyddogaeth `Emit` wedi newid. Yn lle cymryd enw a data dewisol, mae'n cymryd un gwrthrych `WailsEvent` y bydd yn ei allbynnu.

### `Off` a `OffAll`

Yn v2, byddai galwadau `Off` a `OffAll` yn tynnu digwyddiadau i ffwrdd yn JS ac yn Go. Oherwydd natur aml-ffenestr v3, mae hyn wedi newid fel bod y dulliau hyn ond yn berthnasol i'r cyd-destun y'u galwyd. Er enghraifft, os ydych yn galw `Off` mewn ffenestr, dim ond digwyddiadau ar gyfer y ffenestr honno y bydd yn eu tynnu. Os ydych yn defnyddio `Off` yn Go, dim ond digwyddiadau ar gyfer Go y bydd yn eu tynnu.

### Bachau

Mae Bachau Digwyddiad yn nodwedd newydd yn v3. Maent yn caniatáu i chi fachlu i mewn i'r system ddigwyddiadau a chyflawni gweithredoedd pan fydd digwyddiadau penodol yn cael eu hallbynnu. Er enghraifft, gallwch fachlu i mewn i'r digwyddiad `WindowClosing` a chyflawni rhywfaint o lanhau cyn i'r ffenestr gau. Gellir cofrestru bachau ar lefel y cymhwysiad neu ar lefel y ffenestr gan ddefnyddio `RegisterHook`. Bydd bachau lefel cymhwysiad ar gyfer digwyddiadau cymhwysiad. Bydd bachau lefel ffenestr ond yn cael eu galw ar gyfer y ffenestr y'u cofrestrir.

### Nodiadau datblygwr

Pan allbynwch ddigwyddiad yn Go, bydd yn dosbarthu'r digwyddiad i wrrandawyr Go lleol a hefyd i bob ffenestr yn y cymhwysiad. Pan allbynwch ddigwyddiad yn JS, mae'n nawr yn anfon y digwyddiad at y cymhwysiad. Caiff hwn ei brosesu fel petai wedi ei allbynnu yn Go, fodd bynnag bydd ID y anfonwr yn bod hwnnw o'r ffenestr.