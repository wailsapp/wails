# README

The `index.js` file in the `compiled` directory is the entrypoint for the `runtime.js` file that may be
loaded at runtime. This will add `window.wails` and `window._wails` to the global scope.

NOTE: It is preferable to use the `@wailsio/runtime` package to use the runtime.

⚠️ Do not rebuild the runtime manually after updating TS code:
the CI pipeline will take care of this.
PRs that touch build artifacts will be blocked from merging.