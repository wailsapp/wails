## Dewislenni Cyd-destun

Mae dewislenni cyd-destun yn ddewislenni cyd-destunol a ddangosir pan fydd y
defnyddiwr yn clicio'n dde ar elfen. Mae creu dewislen gyd-destun yr un peth â
chreu dewislen safonol, gan ddefnyddio `app.NewMenu()`. I wneud y ddewislen
gyd-destun ar gael i ffenestr, galwch `window.RegisterContextMenu(name, menu)`.
Bydd y enw yn yr id o'r ddewislen gyd-destun ac a ddefnyddir gan y rhaglen
wynebu.

I nodi bod gan elfen ddewislen gyd-destun, ychwanegwch y priodoledd
`data-contextmenu` at yr elfen. Dylai gwerth y priodoledd hwn fod yn enw o
ddewislen gyd-destun a gofrestrwyd yn flaenorol gyda'r ffenestr.

Mae'n bosibl cofrestru dewislen gyd-destun ar lefel y cymhwysiad, gan ei
gwneud ar gael i bob ffenestr. Gellir gwneud hyn gan ddefnyddio
`app.RegisterContextMenu(name, menu)`. Os na ellir dod o hyd i ddewislen
gyd-destun ar lefel y ffenestr, bydd y dewislenni cyd-destun cymhwyso yn cael
eu gwirio. Ceir demo o hyn yn `v3/examples/contextmenus`.