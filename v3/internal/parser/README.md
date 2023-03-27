# Parser

This package contains the static analyser used for parsing Wails projects so that we may:

- Generate the bindings for the frontend
- Generate Typescript definitions for the structs used by the bindings

## Implemented

- [ ] Bound types
  - [x] Struct Literal Pointer
  - [ ] Variable
    - [ ] Assignment
      - [x] Struct Literal Pointer
      - [ ] Function
        - [ ] Same package
        - [ ] Different package
  - [ ] Function

- [x] Parsing of bound methods
  - [x] Method names
  - [x] Method parameters
    - [x] Zero parameters
    - [x] Single input parameter
    - [x] Single output parameter
    - [x] Multiple input parameters
    - [x] Multiple output parameters
    - [x] Named output parameters
    - [x] int/8/16/32/64
      - [x] Pointer
    - [x] uint/8/16/32/64
      - [x] Pointer
    - [x] float
      - [x] Pointer
    - [x] string
      - [x] Pointer
    - [x] bool
      - [x] Pointer
    - [x] Struct
      - [x] Pointer
    - [x] Slices 
      - [x] Pointer
    - [x] Maps
      - [x] Pointer
- [x] Model Parsing
  - [x] In same package
  - [x] In different package
  - [x] Multiple named fields, e.g. one,two,three string
  - [x] Scalars
  - [x] Arrays
  - [x] Maps
  - [x] Structs
    - [x] Recursive
    - [x] Anonymous
- [ ] Generation of models
  - [x] Scalars
  - [ ] Arrays
  - [ ] Maps
  - [x] Structs
- [ ] Generation of bindings

## Limitations

There are many ways to write a Go program so there are many program structures that we would need to support. This is a work in progress and will be improved over time. The current limitations are:

- The call to `application.New()` must be in the `main` package
- Bound structs must be declared as struct literals

