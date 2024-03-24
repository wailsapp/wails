#### Enawau

Yng Ngo, mae enawau yn aml yn cael eu diffinio fel math a set o gysonau. Er enghraifft:

```go
type MyEnum int

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

Oherwydd anghydnawsedd rhwng Go a JavaScript, ni ellir defnyddio mathau custom mewn
ffordd hon. Y strategaeth orau yw defnyddio alias math ar gyfer float64:

```go
type MyEnum = float64

const (
    MyEnumOne MyEnum = iota
    MyEnumTwo
    MyEnumThree
)
```

Yn JavaScript, gallwch chi wedyn ddefnyddio'r canlynol:

```js
const MyEnum = {
  MyEnumOne: 0,
  MyEnumTwo: 1,
  MyEnumThree: 2,
};
```

- Pam defnyddio `float64`? Oni allwn ni ddefnyddio `int`?
    - Oherwydd nad oes gan JavaScript gysyniad o `int`. Mae popeth yn `number`, sy'n cyfieithu i `float64` yn Go. Mae hefyd cyfyngiadau
      ar daflu mathau yn pecyn adlewyrchu Go, sy'n golygu nad yw defnyddio `int` yn
      gweithio.