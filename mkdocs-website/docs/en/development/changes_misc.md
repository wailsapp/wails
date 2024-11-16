### Misc

## Windows Application Options

### WndProcInterceptor

If this is set, the WndProc will be intercepted and the function will be called.
This allows you to handle Windows messages directly. The function should have
the following signature:

```go
func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnValue uintptr, shouldReturn)
```

The `shouldReturn` value should be set to `true` if the returnValue should be
returned by the main wndProc method. If it is set to `false`, the return value
will be ignored and the message will continue to be processed by the main
wndProc method.

## Hide Window on Close + OnBeforeClose

In v2, there was the `HideWindowOnClose` flag to hide the window when it closed.
There was a logical overlap between this flag and the `OnBeforeClose` callback.
In v3, the `HideWindowOnClose` flag has been removed and the `OnBeforeClose`
callback has been renamed to `ShouldClose`. The `ShouldClose` callback is called
when the user attempts to close a window. If the callback returns `true`, the
window will close. If it returns `false`, the window will not close. This can be
used to hide the window instead of closing it.

## Window Drag

In v2, the `--wails-drag` attribute was used to indicate that an element could
be used to drag the window. In v3, this has been replaced with
`--webkit-app-region` to be more in line with the way other frameworks handle
this. The `--webkit-app-region` attribute can be set to any of the following
values:

- `drag` - The element can be used to drag the window
- `no-drag` - The element cannot be used to drag the window

We would have ideally liked to use `app-region`, however this is not supported
by the `getComputedStyle` call on webkit on macOS.
