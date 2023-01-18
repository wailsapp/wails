# Parser

This package contains the static analyser used for parsing Wails projects so that we may:

- Generate the bindings for the frontend
- Generate Typescript definitions for the structs used by the bindings

## Implemented

- [x] Parsing of structs
- [x] Generation of models
  - [x] Scalars
  - [x] Arrays
  - [ ] Maps
  - [x] Structs
- [ ] Generation of bindings

## Limitations

There are many ways to write a Go program so there are many different program structures that we would need to support. This is a work in progress and will be improved over time. The current limitations are:

- The call to `application.New()` must be in the `main` package
- Bound structs must be declared as struct literals

