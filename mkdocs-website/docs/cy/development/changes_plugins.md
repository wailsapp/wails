## Ategion

Mae ategion yn ffordd o ymestyn swyddogaeth eich cais Wails.

### Creu ategyn

Mae ategion yn strwythur Go safonol sy'n cydymffurfio â'r rhyngwyneb canlynol:

```go
type Plugin interface {
    Name() string
    Init(*application.App) error
    Shutdown()
    CallableByJS() []string
    InjectJS() string
}
```

Mae'r dull `Name()` yn dychwelyd enw'r ategyn. Defnyddir hwn at ddibenion cofnodi.

Mae'r dull `Init(*application.App) error` yn cael ei alw pan gaiff yr ategyn ei lwytho. 
Mae'r paramedr `*application.App` yn gymhwysiad y caiff yr ategyn ei lwytho iddo. Bydd unrhyw
wallau yn atal y cais rhag dechrau.

Mae'r dull `Shutdown()` yn cael ei alw pan fydd y cais yn cau.

Mae'r dull `CallableByJS()` yn dychwelyd rhestr o swyddogaethau alladwy y gellir eu galw o'r
blaen-wyneb. Rhaid i enwau'r dulliau hyn gyfateb yn union i enwau'r dulliau a allodir
gan yr ategyn.

Mae'r dull `InjectJS()` yn dychwelyd JavaScript y dylid ei fewnosod i bob ffenestr wrth iddynt
gael eu creu. Mae hyn yn ddefnyddiol ar gyfer ychwanegu swyddogaethau JavaScript 
cyfatebol i'r ategyn.