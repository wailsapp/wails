# Croniclau

<!--
Bydd pob newid sylweddol i'r prosiect hwn yn cael ei ddogfennu yn y ffeil hon.

Mae'r fformat yn seiliedig ar [Cadw Croniclau](https://keepachangelog.com/en/1.0.0/),
ac mae'r prosiect hwn yn cydymffurfio â [Fersiwneiddio Semantig](https://semver.org/spec/v2.0.0.html).

- `Ychwanegwyd` ar gyfer nodweddion newydd.
- `Newidiwyd` ar gyfer newidiadau mewn swyddogaeth bresennol.
- `Wedi Dibrisio` ar gyfer nodweddion a fydd yn cael eu dileu yn fuan.
- `Tynnwyd` ar gyfer nodweddion a gafodd eu tynnu yn awr.
- `Wedi Trwsio` ar gyfer unrhyw ddatrysiadau grwydro.
- `Diogelwch` os oes agoreiddiadau diamddiffyn.

-->

## [Heb ei ryddhau]

### Ychwanegwyd
- [darwin] ychwanegu Digwyddiad ApplicationShouldHandleReopen i fedru ymdrin â chlicio ar yr eicon doc gan @5aaee9 yn [#2991](https://github.com/wailsapp/wails/pull/2991)
- [darwin] ychwanegu getPrimaryScreen/getScreens i impl gan @tmclane yn [#2618](https://github.com/wailsapp/wails/pull/2618)
- [darwin] ychwanegu opsiwn ar gyfer dangos y bar offer mewn modd sgrin llawn ar macOS gan [@fbbdev](https://github.com/fbbdev) yn [#3282](https://github.com/wailsapp/wails/pull/3282)
- [linux] ychwanegu rhesymeg onKeyPress i drosi allwedd linux i gyflymydd gan @[Atterpac](https://github.com/Atterpac) yn [#3022](https://github.com/wailsapp/wails/pull/3022])
- [linux] ychwanegu tasg `rhedeg:linux` gan [@marcus-crane](https://github.com/marcus-crane) yn [#3146](https://github.com/wailsapp/wails/pull/3146)
- allforio dull `SetIcon` gan @almas1992 yn [PR](https://github.com/wailsapp/wails/pull/3147)
- Gwella `OnShutdown` gan @almas1992 yn [PR](https://github.com/wailsapp/wails/pull/3189)
- adfer dull `ToggleMaximise` yn y rhyngwyneb `Window` gan [@fbbdev](https://github.com/fbbdev) yn [#3281](https://github.com/wailsapp/wails/pull/3281)

### Wedi Trwsio

- Wedi trwsio prosesau zombie wrth weithio mewn modd datblygu drwy ddiweddaru i'r diweddaraf gan [Atterpac](https://github.com/atterpac) yn [#3320](https://github.com/wailsapp/wails/pull/3320).
- Wedi trwsio ffynhonnell ffeil webkit appimage gan [Atterpac](https://github.com/atterpac) yn [#3306](https://github.com/wailsapp/wails/pull/3306).
- Wedi trwsio Doctor fygythiad pecyn apt gan [Atterpac](https://github.com/Atterpac) yn [#2972](https://github.com/wailsapp/wails/pull/2972).
- Wedi trwsio'r cais wedi rhewi wrth ddod allan (Darwin) gan @5aaee9 yn [#2982](https://github.com/wailsapp/wails/pull/2982)
- Wedi trwsio lliwiau cefndir yr enghreifftiau ar Windows gan [mmgvh](https://github.com/mmghv) yn [#2750](https://github.com/wailsapp/wails/pull/2750).
- Wedi trwsio dewislenni cyd-destun rhagosodedig gan [mmgvh](https://github.com/mmghv) yn [#2753](https://github.com/wailsapp/wails/pull/2753).
- Wedi trwsio gwerth hecsadegol ar gyfer bysellau saeth ar Darwin gan [jaybeecave](https://github.com/jaybeecave) yn [#3052](https://github.com/wailsapp/wails/pull/3052).
- Gosod llusgo-a-gollwng ar gyfer Windows i weithio. Ychwanegwyd gan [@pylotlight](https://github.com/pylotlight) yn [PR](https://github.com/wailsapp/wails/pull/3039)
- Wedi trwsio bygiau ar gyfer linux yn y meddyg os nad oes gan y defnyddiwr y gyrwyr priodol wedi'u gosod. Ychwanegwyd gan [@pylotlight](https://github.com/pylotlight) yn [PR](https://github.com/wailsapp/wails/pull/3032)
- Trwsio graddio dpi wrth gychwyn (windows). Newidiwyd gan @almas1992 yn [PR](https://github.com/wailsapp/wails/pull/3145)
- Trwsio'r llinell amnewid yn `go.mod` i ddefnyddio llwybrau cymharol - Trwsio llwybrau Windows gyda gofodau gan @leaanthony.
- Trwsio gweithredu clicio Maclanwad system wrth ddim cysylltiedig â ffenestr gan [thomas-senechal](https://github.com/thomas-senechal) yn PR [#3207](https://github.com/wailsapp/wails/pull/3207)
- Trwsio adeiladu Windows yn methu oherwydd opsiwn anhysbys gan [thomas-senechal](https://github.com/thomas-senechal) yn PR [#3208](https://github.com/wailsapp/wails/pull/3208)
- Trwsio URL sylfaenol anghywir wrth agor ffenestr ddwywaith gan @5aaee9 yn PR [#3273](https://github.com/wailsapp/wails/pull/3273)
- Trwsio trefn brigiau os yn y dull `WebviewWindow.Restore` gan [@fbbdev](https://github.com/fbbdev) yn [#3279](https://github.com/wailsapp/wails/pull/3279)
- Cyfrifo `startURL` yn gywir ar draws galwadau lluosog `GetStartURL` pan fo `FRONTEND_DEVSERVER_URL` yn bresennol. [#3299](https://github.com/wailsapp/wails/pull/3299)

### Newidiwyd

### Tynnwyd

### Wedi Dibrisio

### Diogelwch