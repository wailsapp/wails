# README

The `main.js` file in this directory is the entrypoint for the `runtime.js` file that may be
loaded at runtime. This will add `window.wails` and `window._wails` to the global scope.

NOTE: It is preferable to use the `@wailsio/runtime` package to use the runtime.

After updating any files in this directory, you must run `wails3 task build:runtime` to regenerate the compiled JS. 