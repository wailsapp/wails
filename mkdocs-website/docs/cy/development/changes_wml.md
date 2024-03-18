## Iaith Marcio Wails (WML)

Mae'r Iaith Marcio Wails yn iaith farcio syml sy'n caniatáu i chi ychwanegu
swyddogaeth at elfennau HTML safonol heb ddefnyddio JavaScript.

Mae'r tagiau canlynol yn cael eu cefnogi ar hyn o bryd:

### `data-wml-event`

Mae hyn yn pennu y bydd digwyddiad Wails yn cael ei allyrru pan gliciwyd ar yr
elfen. Dylai gwerth yr priodoledd fod yn enw'r digwyddiad i'w allyrru.

Enghraifft:

```html
<button data-wml-event="myevent">Cliciwch Fi</button>
```

Weithiau mae angen i'r defnyddiwr gadarnhau gweithred. Gellir gwneud hyn drwy
ychwanegu'r briodoledd `data-wml-confirm` at yr elfen. Bydd gwerth y briodoledd
hwn yn fesur i'w ddangos i'r defnyddiwr.

Enghraifft:

```html
<button data-wml-event="delete-all-items" data-wml-confirm="Ydych chi'n siŵr?">
  Dileu Pob Eitem
</button>
```

### `data-wml-window`

Gellir galw unrhyw fethododd `wails.window` drwy ychwanegu'r briodoledd
`data-wml-window` at elfen. Dylai gwerth y briodoledd fod yn enw'r
dull i'w alw. Dylai enw'r dull fod yn yr un acen â'r dull.

```html
<button data-wml-window="Close">Cau'r Ffenestr</button>
```

### `data-wml-trigger`

Mae'r briodoledd hwn yn pennu pa ddigwyddiad JavaScript ddylai ysgogi'r
weithred. Y rhagosodiad yw `click`.

```html
<button data-wml-event="hover-box" data-wml-trigger="mouseover">
  Gallwch hofran drosodd fi!
</button>
```