# v3

*NOTE*: The examples in this directory may or may not compile / run at any given time during alpha development.


## Running the examples

    cd v3/examples/<example>
    go mod tidy
    go run .

### macOS private APIs

Default macOS builds include Wails' private WebKit and AppKit integrations.
Examples that use transparent or translucent web content, Liquid Glass, or
automatic Web Inspector opening exercise those APIs.

To run any example using public macOS APIs only:

    go run -tags noprivateapis .

The example source and window options remain unchanged. Effects without public
equivalents degrade as documented in the
[macOS build guide](https://v3alpha.wails.io/guides/build/macos/#private-macos-apis).

## Compiling the examples

    cd v3/examples/<example>
    go mod tidy
    go build
    ./<example>

Use `go build -tags noprivateapis` to compile an example without Wails' private
macOS API integrations.
