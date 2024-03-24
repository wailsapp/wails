# Prif Swyddogaethau Trywydd

Mae'r dulliau hyn yn swyddogaethau cymorth i redeg cod ar y prif drywydd. Mae hyn yn
ofynnol pan fyddwch am redeg cod cyfaddas ar y llwyfan UI.

### InvokeSync

API: `InvokeSync(fn func())`

Mae'r swyddogaeth hon yn rhedeg y swyddogaeth a drosglwyddwyd (`fn`) yn ddilynebol. Mae'n defnyddio WaitGroup
(`wg`) i sicrhau bod y prif drywydd yn aros i `fn` swyddogaeth orffen
cyn iddo barhau. Os bydd panig yn digwydd o fewn `fn`, bydd yn cael ei drosglwyddo i'r
swyddogaeth trin panig `PanicHandler`, a ddiffinnir yn opsiynau'r cymhwysiad.

### InvokeSyncWithResult

API: `InvokeSyncWithResult[T any](fn func() T) (res T)`

Mae'r swyddogaeth hon yn gweithio'n debyg i `InvokeSync(fn func())`, fodd bynnag, mae'n rhoi
canlyniad. Defnyddiwch hyn ar gyfer galw unrhyw swyddogaeth gyda un canlyniad yn unig.

### InvokeSyncWithError

API: `InvokeSyncWithError(fn func() error) (err error)`

Mae'r swyddogaeth hon yn rhedeg `fn` yn ddilynebol ac yn dychwelyd unrhyw wall a gynhyrchir gan `fn`.
Sylwch y bydd y swyddogaeth hon yn adfer o banig os bydd un yn digwydd yn ystod
gweithrediad `fn`.

### InvokeSyncWithResultAndError

API:
`InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error)`

Mae'r swyddogaeth hon yn rhedeg `fn` yn ddilynebol ac yn dychwelyd canlyniad o fath `T` a
gwall.

### InvokeAsync

API: `InvokeAsync(fn func())`

Mae'r swyddogaeth hon yn rhedeg `fn` yn asyng. Mae'n rhedeg y swyddogaeth a roddir ar y
prif drywydd. Os bydd panig yn digwydd o fewn `fn`, bydd yn cael ei drosglwyddo i'r
swyddogaeth trin panig `PanicHandler`, a ddiffinnir yn opsiynau'r cymhwysiad.

---

_Sylw_: Bydd y swyddogaethau hyn yn rhwystro gweithrediad nes bod `fn` wedi gorffen. Mae'n
hanfodol sicrhau nad yw `fn` yn rhwystro. Os bydd angen i chi redeg swyddogaeth sy'n
rhwystro, defnyddiwch `InvokeAsync` yn lle.