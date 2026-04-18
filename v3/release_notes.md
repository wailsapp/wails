## Added
- Add native Liquid Glass effect support for macOS with NSGlassEffectView (macOS 15.0+) and NSVisualEffectView fallback, including comprehensive material customization options by @leaanthony in [#4534](https://github.com/wailsapp/wails/pull/4534)

## Fixed
- Fix [#5089](https://github.com/wailsapp/wails/issues/5089): on macOS, URL-scheme launches no longer lose the URL when `SingleInstance` is active. The second instance now briefly runs a minimal NSApplication event loop to capture the `kAEGetURL` Apple Event and relay it to the first instance via `SecondInstanceData.Args`, matching Windows/Linux behaviour.