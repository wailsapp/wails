# linux-gtk4-uniqueness (fix side)

Fix verification for the hang reproduced by the example of the same name in the
unpatched checkout. `main.go` is byte-identical — this bug needs no source
change, only the patched `APPLICATION_DEFAULT_FLAGS`.

## Test

```sh
cd v3/examples/linux-gtk4-uniqueness
go run .        # window appears, leave it running
go run .        # second process: window also appears
```

Both processes must print `ApplicationStarted fired` and open a window. On the
unpatched tree the second prints `still waiting` forever and draws nothing.

Confirm you are actually testing the patched tree — the two checkouts differ
only in the module they resolve, so it is easy to build the wrong one:

```sh
go list -m -f '{{.Dir}}' github.com/wailsapp/wails/v3
strings <binary> | grep -c 'wails/v3/pkg/application'
```

## The fix

`pkg/application/linux_cgo.h`:

```c
#define APPLICATION_DEFAULT_FLAGS G_APPLICATION_NON_UNIQUE
```

GTK no longer makes the application unique, so a second launch never registers
as a remote instance, `activate` always fires locally, and
`linuxApp.waitForActivation` never blocks. Single instance behaviour is left to
`SingleInstanceOptions`, which is what asks for it explicitly.

## Note

Running this needs a reachable accessibility bus — WebKit's sandbox launches a
dbus-proxy against `$XDG_RUNTIME_DIR/at-spi` and aborts the process if it is
missing, which will stop the example from starting at all inside a container
that does not forward that socket.
