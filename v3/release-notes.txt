## Added
- Add Content Protection on Windows/Mac by [@leaanthony](https://github.com/leaanthony) based on the original work of [@Taiterbase](https://github.com/Taiterbase) in this [PR](https://github.com/wailsapp/wails/pull/4241)
- Add support for passing CLI variables to Task commands through `wails3 build` and `wails3 package` aliases (#4422) by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4488)

## Changed
- `window.NativeWindowHandle()` -> `window.NativeWindow()` by @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refactor internal window handling by @leaanthony in [#4471](https://github.com/wailsapp/wails/pull/4471)
+ Fix extra-broad Linux package dependencies, fix outdated RPM dependencies.