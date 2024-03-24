Dyma'r cyfieithiad i'r gymraeg:

## Rhwymedigaethau

Mae rhwymedigaethau yn gweithio mewn modd tebyg i v2, drwy ddarparu ffordd i rwymo
dulliau strwythur i'r rhyngwyneb blaen. Gellir eu galw yn y rhyngwyneb blaen gan
ddefnyddio'r wraperi rhwymedigaeth a gynhyrchwyd gan y gorchymyn `wails3 generate bindings`:

```javascript
// @ts-check
// Mae'r ffeil hon wedi'i chynhyrchu'n awtomatig. PEIDIWCH Â'I GOLYGU

import { main } from "./models";

window.go = window.go || {};
window.go.main = {
  GreetService: {
    /**
     * GreetService.Greet
     * Mae Greet yn cyfarch rhywun
     * @param name {string}
     * @returns {Promise<string>}
     **/
    Greet: function (name) {
      wails.CallByID(1411160069, ...Array.prototype.slice.call(arguments, 0));
    },

    /**
     * GreetService.GreetPerson
     * Mae GreetPerson yn cyfarch rhywun
     * @param person {main.Person}
     * @returns {Promise<string>}
     **/
    GreetPerson: function (person) {
      wails.CallByID(4021313248, ...Array.prototype.slice.call(arguments, 0));
    },
  },
};
```

Mae dulliau rhwymo wedi'u cuddio'n ddiofyn, ac maent yn cael eu hadnabod gan IDs uint32,
a gyfrifir gan ddefnyddio'r [algorithm hasio FNV](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function).
Mae hyn er mwyn atal enw'r dull rhag cael ei ddatgelu mewn adeiladau cynhyrchiol. Mewn
modd dadfygio, mae'r IDs dull yn cael eu logio ynghyd â'r ID a gyfrifwyd o'r dull
i helpu i ddadfygio. Os ydych chi am ychwanegu haen arall o guddio, gallwch
ddefnyddio'r opsiwn `BindAliases`. Mae hyn yn caniatáu ichi bennu map o IDs alias i
IDs dull. Pan fydd y rhyngwyneb blaen yn galw dull gan ddefnyddio ID, bydd yr ID dull
yn cael ei chwilio yn y map alias yn gyntaf am gywiro. Os nad yw'n ei ganfod,
mae'n tybio mai ID dull safonol yw ac yn ceisio canfod y dull yn y ffordd arferol.

Enghraifft:

```go
	app := application.New(application.Options{
		Bind: []any{
			&GreetService{},
		},
		BindAliases: map[uint32]uint32{
			1: 1411160069,
			2: 4021313248,
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
```

Nawr gallwn alw gan ddefnyddio'r alias hwn yn y rhyngwyneb blaen: `wails.Call(1, "byd!")`.

### Galwadau anniogel

Os nad ydych chi'n poeni am eich galwadau yn cael eu cyhoeddi mewn testun plaen yn eich binari
ac nid oes gennych fwriad o ddefnyddio [garble](https://github.com/burrowers/garble), yna
gallwch ddefnyddio'r dull `wails.CallByName()` anniogel. Mae'r dull hwn yn cymryd enw
llawn cymhwysol y dull i'w alw a'r arguments i'w pasio iddo.
Enghraifft:

    ```go
    wails.CallByName("main.GreetService.Greet", "byd!")
    ```

!!! angen gofal

    Darperir hwn dim ond fel dull cyfleustra ar gyfer datblygu. Ni chyngherir i'w ddefnyddio mewn cynhyrchiad.
