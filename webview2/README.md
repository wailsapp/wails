# WebView2 Bindings for Go

This directory contains automatically generated Go bindings for Microsoft's WebView2 control, used by Wails v3 on Windows.

## Important: Do Not Edit Generated Files

All files in `pkg/webview2/` are **automatically generated** from the official Microsoft WebView2 IDL (Interface Definition Language) files. Do not manually edit these files, as your changes will be overwritten when the bindings are regenerated.

## Regenerating Bindings

To regenerate the WebView2 bindings:

```bash
cd webview2
go run scripts/generator/cmd/webview2gen/main.go full
```

Or using the built binary:

```bash
cd webview2
go build -o bin/webview2gen scripts/generator/cmd/webview2gen
./bin/webview2gen full
```

You can also use `go generate`:

```bash
cd webview2
go generate ./pkg/webview2/...
```

## Current Version

The bindings are generated from WebView2 SDK version: `1.0.2903.40`

See `scripts/latest_version.txt` for the current version.

## v1 → v2 Migration

The WebView2 bindings have been completely rewritten (v2) with the following improvements:

### Fixed Bugs

1. **Memory Leak in String Handling**: The v1 generator had a critical bug where string pointers returned by COM methods were not properly freed using `CoTaskMemFree`, causing memory leaks.
2. **Incorrect COM Reference Counting**: Interface references were not properly managed, leading to premature releases or dangling pointers.
3. **Missing Error Propagation**: Many COM error codes were silently ignored, making debugging extremely difficult.
4. **Incorrect VARIANT Handling**: The `VARIANT` type was incorrectly implemented as a simple `uintptr`, losing type information and causing crashes.
5. **Missing BSTR Conversion**: BSTR (Basic String) types were not properly converted to/from Go strings, causing encoding issues.
6. **Incomplete Interface Support**: Many newer interfaces were missing, limiting access to WebView2 features.

### Impact of Fixes

- **Memory Usage**: Applications using WebView2 no longer leak memory over time.
- **Stability**: COM reference counting is now correct, preventing crashes.
- **Debugging**: Errors are now properly propagated, making issues easier to diagnose.
- **Feature Support**: All WebView2 interfaces up to version 1.0.2903.40 are available.

## Capability System

The WebView2 SDK uses a "capability" system where newer APIs are only available on specific runtime versions. The generated bindings include all APIs, and it's the application's responsibility to:

1. Check the installed WebView2 runtime version
2. Query interface availability using COM's `QueryInterface`
3. Gracefully handle unsupported capabilities

See the [Wails documentation](https://wails.io/docs/reference/runtime/windows#webview2-version-requirements) for version-specific requirements.

## Testing

To verify that generated bindings are up-to-date:

```bash
cd webview2
go run scripts/generator/cmd/webview2gen/main.go verify
```

This will regenerate all bindings from IDL and compare them with the committed files. Any differences indicate that the bindings need to be regenerated.

## Generator Source Code

The generator source code is in `scripts/generator/` and consists of:

- **Parser**: Parses IDL files using the `participle` library
- **Emitter**: Generates Go code from the parsed IDL AST
- **Templates**: Go templates used for code generation

The generator is written in Go and can be modified if needed to handle special cases or add new features.

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed documentation on the generator's design, including:

- Why we use a programmatic emitter instead of templates
- File-per-interface organization
- Parser grammar overview
- Type system design
- Capability system implementation

## Contributing

When contributing changes to the WebView2 bindings:

1. **Do NOT** modify files in `pkg/webview2/` directly
2. If you need to change bindings, modify the generator in `scripts/generator/`
3. Regenerate bindings using `webview2gen full`
4. Run `webview2gen verify` to ensure everything matches
5. Include both the generator changes and regenerated bindings in your PR

## Version Compatibility

| Wails Version | Min WebView2 Runtime |
|---------------|---------------------|
| v3.0.0       | 1.0.1418.22        |

Always ensure users have the minimum required runtime version installed. Applications can check the installed version using Wails' runtime APIs.

## Troubleshooting

### "Unsupported capability" Error

This error occurs when your application tries to use a WebView2 API that isn't available in the installed runtime version. Solutions:

1. Update the WebView2 runtime: https://go.microsoft.com/fwlink/p/?LinkId=2124703
2. Check the installed version using Wails' runtime APIs
3. Use version-specific APIs and provide fallbacks for older runtimes

### "Generated output does not match committed files"

Run `webview2gen full` to regenerate the bindings. This error typically occurs when:
- A new IDL file was added to `scripts/`
- The generator code was modified
- `scripts/latest_version.txt` was updated
