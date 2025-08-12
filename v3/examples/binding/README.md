# Bindings Example

This example demonstrates how to generate bindings for your application.

To generate bindings, run `wails3 generate bindings -clean -b -d assets/bindings` in this directory.

See more options by running `wails3 generate bindings --help`.

## Notes
  - The bindings generator is still a work in progress and is subject to change.
  - The generated code uses `wails.CallByID` by default. This is the most robust way to call a function.
