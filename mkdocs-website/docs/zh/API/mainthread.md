# 主线程函数

这些方法是在主线程上运行代码的实用函数。当您想要在UI线程上运行自定义代码时，这是必需的。

### InvokeSync

API: `InvokeSync(fn func())`

此函数以同步方式运行传入的函数（`fn`）。它使用一个`WaitGroup`（`wg`）确保主线程在`fn`函数完成之前等待，然后才继续执行。如果在`fn`内部发生恐慌，它将传递给应用程序选项中定义的处理程序函数`PanicHandler`。

### InvokeSyncWithResult

API: `InvokeSyncWithResult[T any](fn func() T) (res T)`

此函数与`InvokeSync(fn func())`类似，但它返回一个结果。可用于调用具有单个返回值的任何函数。

### InvokeSyncWithError

API: `InvokeSyncWithError(fn func() error) (err error)`

此函数同步运行`fn`并返回`fn`产生的任何错误。请注意，如果在`fn`执行期间发生恐慌，此函数将从恢复。

### InvokeSyncWithResultAndError

API: `InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error)`

此函数同步运行`fn`并返回类型为`T`的结果和一个错误。

### InvokeAsync

API: `InvokeAsync(fn func())`

此函数以异步方式运行`fn`。它在主线程上运行给定的函数。如果在`fn`内部发生恐慌，它将传递给应用程序选项中定义的处理程序函数`PanicHandler`。

---

注意：这些函数将阻塞执行，直到`fn`完成。确保`fn`不会阻塞至关重要。如果需要运行阻塞函数，请改用`InvokeAsync`。