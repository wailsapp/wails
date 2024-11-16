# Bindings Example

This example demonstrates how to generate bindings for your application.

To generate bindings, run `wails3 generate bindings -b -d assets/bindings` in this directory.

See more options by running `wails3 generate bindings --help`.

## Notes
  - The bindings generator is still a work in progress and is subject to change.
  - The generated code uses `wails.CallByID` by default. This is the most secure and quickest way to call a function. In a future release, we will look at generating `wails.CallByName` too.
