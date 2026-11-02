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

Each runnable example has a local `wails.yaml` recording its build requirements.
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
## Private macOS APIs

Some examples use Wails features that require private macOS APIs. To see those effects on macOS, opt in with the Go build tag `private_mac_apis`. The public Wails API stays the same without the tag; private-only operations become no-ops, and documented alternatives remain available.

For examples that run directly with Go:

    cd v3/examples/spotlight
    go run -tags private_mac_apis .

For frontend examples with a `run` task, use their build pipeline so bindings and frontend assets are ready:

    cd v3/examples/badge
    wails3 build -tags private_mac_apis
    wails3 task run

For live reload in those examples, use `EXTRA_TAGS=private_mac_apis wails3 dev`. The `ios`, `mobile`, and `mac-window-tabs` examples have their own modules; their READMEs include `GOWORK=off` where needed. Follow each example's existing setup instructions first. The `dev` example remains a work in progress.

| Examples | Why the tag is needed on macOS | Without the tag |
| --- | --- | --- |
| [badge](badge/), [badge-custom](badge-custom/), [contextmenus](contextmenus/), [custom-protocol-example](custom-protocol-example/), [dev](dev/), [dock](dock/), [drag-n-drop](drag-n-drop/), [environment](environment/), [events](events/), [file-association](file-association/), [ios](ios/), [mac-window-tabs](mac-window-tabs/), [mobile](mobile/), [notifications](notifications/), [plain](plain/), [raw-message](raw-message/), [screen](screen/), [spotlight](spotlight/), [window-api](window-api/), [wml](wml/) | Transparent webview above a translucent native backdrop | Webview remains opaque |
| [liquid-glass](liquid-glass/) | Webview transparency and existing native glass styles | Opaque webview; public style alternatives |
| [notch-notification](notch-notification/) | Webview transparency inside the native notch window | Webview remains opaque |
| [events-bug](events-bug/), [keybindings](keybindings/), [window](window/) | Programmatic `OpenDevTools()` | Inspector-opening call is a no-op |

No example sets Liquid Glass `GroupID` or `GroupSpacing`. Those options also require the tag and are otherwise no-ops. The iOS and mobile examples need this tag only when running their macOS desktop variant; iOS and Android builds are unaffected.

To make a production binary directly with Go, use:

    go build -tags production,private_mac_apis .

Production builds still disable inspector support by default. To keep an inspector shortcut working in a production example, use `go build -tags production,devtools,private_mac_apis .`. On macOS 13.3+, Safari inspection remains available without private APIs in development builds or builds with `devtools`.

For task-based examples, `wails3 build -tags private_mac_apis` passes the tag to the build task. Omit the private API tag from any of these commands to use public APIs only. See the [private macOS API reference](https://v3.wails.io/guides/build/private-macos-apis/) for the complete fallback behavior.
