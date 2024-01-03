# Roadmap

The roadmap is a living document and is subject to change. If you have any
suggestions, please open an issue. Each milestone will have a set of goals that
we are aiming to achieve. These are subject to change.

## Alpha milestones

### Current: Alpha 3

#### Goals

The Alpha 3 cycle aims to provide bindings support. Wails 3 uses a new static analysis approach which allows us to provide 
a better bindings experience than in Wails 2. 
We also want to get all examples working on Linux.

#### How Can I Help?

You can generate bindings using the `wails3 generate bindings` command. This will generate bindings for all exported struct methods bound to your project.
Run `wails3 generate bindings -help` to view options that govern how bindings are generated.
 
The tests for the bindings generator can be found [here](https://github.com/wailsapp/wails/tree/v3-alpha/v3/internal/parser) with the test data located in the `testdata` directory. 

Review the table below and look for untested scenarios. The parser code and tests are located in `v3/internal/parser`. All tests can be run using `go test ./...` from the `v3` directory.
Basically, try to break it and let us know if you find any issues! :smile:

#### Status

Bindings for struct (CallByID):

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha3-bindings-callbyid.csv") }}

Bindings for struct (CallByName):

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha3-bindings-callbyname.csv") }}

Models:

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha3-models.csv") }}


Examples:

- [ ] All examples working on Linux

## Upcoming milestones

### Alpha 4

#### Goals

- [ ] `wails3 dev` command
- [ ] `wails3 package` command

### Alpha 5

#### Goals

- [ ] Keyboard shortcuts
    - Window Level shortcuts
    - Application Level shortcuts (applies to all windows)
    - Ensure Keydown/Keyup events are sent to JS if not handled by Go

## Previous milestones

### Alpha 2

#### Goals

Alpha 2 aims to introduce [Taskfile](https://taskfile.dev) support. This will
allow us to have a single, extensible build system that works on all platforms.
We also want to get all examples working on Linux.

#### Status

- [ ] All examples working on Linux
- [x] Init & Build commands


- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha2.csv") }}

### Alpha 1

#### Goals

Alpha 1 is the initial release. It is intended to get feedback on the new API
and to get people experimenting with it. The main goal is to get most of the
examples working on all platforms.

#### Status

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha1.csv") }}
