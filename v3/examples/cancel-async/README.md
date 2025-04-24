# Binding Call Cancelling Example

Sometimes, we need to execute functions that become long-running processes. In Go, we normally use context to cancel
such processes but how can we achieve this in Javascript? 

Automatically, that's how! 

If your Service method's first parameter is a `context.Context`, then the bindings generator will generate a 
"cancellable promise". This allows you to specify the signal to cancel the request.

This method: 
```go
    func (*Service) LongRunning(ctx context.Context, milliseconds int) error {
        select {
        case <-time.After(time.Duration(milliseconds) * time.Millisecond):
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
```
can be called, and cancelled, from the frontend using the following semantics:

```js
    controller = new AbortController();
    
    try {
        await Service.LongRunning(duration).cancelOn(controller.signal);
        document.getElementById("result").innerText = "Long running call complete!";
    } catch (err) {
        document.getElementById("result").innerText = (err instanceof CancelError)
            ? "Call cancelled: " + err.cause
            : "Call failed: " + err;
    } finally {
        lastCall = null;
    }
```

This may now be cancelled at any time using: 
```js
controller.abort("cancel button pressed!");
```

## Notes

To regenerate bindings, run `wails3 generate bindings -clean -b -d assets/bindings` in this directory.
