# Bindings Example

This example demonstrates how to manually generate bindings for your application.

To generate bindings, run `wails3 generate bindings -clean -b -d assets/bindings` in this directory.
The bindings will be available in `assets/bindings/github.com/wailsapp/wails/v3/examples/binding`.
Examine `assets/index.html` to see how these bindings are used.
Run the main application to see it in action.

See more options by running `wails3 generate bindings --help`.

## Notes
  - The generated code uses `wails.CallByID` by default. This is the most robust way to call a function.
