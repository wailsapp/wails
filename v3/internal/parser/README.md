# Parser

This package contains the static analyser used for parsing Wails projects so that we may:

- Generate the bindings for the frontend
- Generate Typescript definitions for the structs used by the bindings

## Implemented

- [ ] Parsing of bound methods
  - [x] Method names
  - [ ] Method parameters
    - [x] Zero parameters
    - [x] Single input parameter
    - [x] Single output parameter
    - [x] Multiple input parameters
    - [x] Multiple output parameters
    - [x] Named output parameters
    - [x] int
      - [x] Pointer
    - [x] uint
      - [ ] Pointer
    - [x] float
      - [ ] Pointer
    - [x] string
      - [ ] Pointer
    - [x] bool
      - [ ] Pointer
    - [ ] Struct
      - [x] Pointer
    - [x] Slices 
      - [ ] Pointer
      - [x] Recursive
    - [x] Maps
      - [ ] Pointer
- [ ] Model Parsing
  - [x] In same package
  - [ ] In different package
  - [ ] Multiple named fields, e.g. one,two,three string
  - [ ] Scalars
  - [ ] Arrays
  - [ ] Maps
  - [ ] Structs
    - [ ] Anonymous structs
- [ ] Generation of models
  - [ ] Scalars
  - [ ] Arrays
  - [ ] Maps
  - [ ] Structs
- [ ] Generation of bindings

## Limitations

There are many ways to write a Go program so there are many different program structures that we would need to support. This is a work in progress and will be improved over time. The current limitations are:

- The call to `application.New()` must be in the `main` package
- Bound structs must be declared as struct literals

