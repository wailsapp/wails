# v3

*NOTE*: The examples in this directory may or may not compile / run at any given time during alpha development.


## Running the examples

    cd v3/examples/<example>
    export WAILS_EXP_USE_WAKE=1
    wails3 run

## Compiling the examples

    cd v3/examples/<example>
    export WAILS_EXP_USE_WAKE=1
    wails3 build

The build output lists the resulting executable or application bundle.

Each runnable example has a local `wails.hcl` recording its build requirements.
`wails3 run --plan` shows the selected tags and launch settings without building.
Use `wails3 run -- --application-flag` to pass arguments to the example.
Without a local or discoverable parent manifest, `wails3 run` uses `go run .`.

Mobile examples default to the host desktop target. Select a mobile target with
`--target android/arm64 --device <serial>` or, on macOS,
`--target ios/arm64 --device <simulator-udid>`. Physical iOS devices also require
`--destination device` and configured signing credentials.

A Linux build and launch smoke sweep is available from the repository root:

```sh
python3 v3/scripts/test-examples-run.py --cli /absolute/path/to/wails3 --output /tmp/wails-example-tests
```

The sweep opens applications briefly and requests their termination. Its report
separates unsupported hosts, build failures, launch failures and launch smoke
checks; it does not replace interactive checks of each example's demonstrated
features. Use `--example badge-custom` to rerun a single example.
