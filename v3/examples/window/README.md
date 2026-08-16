# Window Example

This example is a demonstration of the Windows API.

## Running the example

To run the example, simply run the following command:

```bash
go run .
```

On macOS, the transparent and translucent window actions use Wails' default
private WebKit integration. Run `go run -tags noprivateapis .` to exercise the
public-only fallback; the rest of the example is unchanged, but those webviews
remain opaque.

# Status

| Platform | Status  |
|----------|---------|
| Mac      |         |
| Windows  | Working |
| Linux    |         |
