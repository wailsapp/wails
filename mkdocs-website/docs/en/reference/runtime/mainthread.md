# Main Thread Functions

These methods are utility functions to run code on the main thread. This is
required when you want to run custom code on the UI thread.

### InvokeSync

API: `InvokeSync(fn func())`

This function runs the passed function (`fn`) synchronously. It uses a WaitGroup
(`wg`) to ensure that the main thread waits for the `fn` function to finish
before it continues. If a panic occurs inside `fn`, it will be passed to the
handler function `PanicHandler`, defined in the application options.

### InvokeSyncWithResult

API: `InvokeSyncWithResult[T any](fn func() T) (res T)`

This function works similarly to `InvokeSync(fn func())`, however, it yields a
result. Use this for calling any function with a single return.

### InvokeSyncWithError

API: `InvokeSyncWithError(fn func() error) (err error)`

This function runs `fn` synchronously and returns any error that `fn` produces.
Note that this function will recover from a panic if one occurs during `fn`'s
execution.

### InvokeSyncWithResultAndError

API:
`InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error)`

This function runs `fn` synchronously and returns both a result of type `T` and
an error.

### InvokeAsync

API: `InvokeAsync(fn func())`

This function runs `fn` asynchronously. It runs the given function on the main
thread. If a panic occurs inside `fn`, it will be passed to the handler function
`PanicHandler`, defined in the application options.

---

_Note_: These functions will block execution until `fn` has finished. It's
critical to ensure that `fn` doesn't block. If you need to run a function that
blocks, use `InvokeAsync` instead.
