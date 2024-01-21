# Roadmap

The roadmap is a living document and is subject to change. If you have any
suggestions, please open an issue. Each milestone will have a set of goals that
we are aiming to achieve. These are subject to change.

## Known Issues

- Generating bindings for a method that imports a package that has the same name as another imported package is currently not supported.

## Alpha milestones

### Current: Alpha 4

#### Goals

The Alpha 4 cycle aims to provide the `dev` and `package` commands. 
The `wails dev` command should do the following:
- Build the application
- Start the application
- Start the frontend dev server
- Watch for changes to the application code and rebuild/restart as necessary

The `wails package` command should do the following:
- Build the application
- Package the application in a platform specific format
  - Windows: Standard executable, NSIS Installer
  - Linux: AppImage
  - MacOS: Standard executable, App Bundle
- Support obfuscation of the application code

- We also want to get all examples working on Linux.

#### How Can I Help?

!!! note
    Report any issues you find using [this guide](./getting-started/feedback.md).


- Install the CLI using the instructions [here](./getting-started/installation).
- Run `wails3 doctor` and ensure that all dependencies are installed. 
- Generate a new project using `wails3 init`.

Test the `wails3 dev` command:

- Run `wails3 dev` in the project directory. It should run the application in development mode.
- Try changing files and ensure that the application is rebuilt and restarted.
- Run `wails3 dev -help` to view options.
- Try different options and ensure that they work as expected.

Test the `wails3 package` command:

- Run `wails3 package` in the project directory.
- Check that the application is packaged correctly for the current platform.
- Run `wails3 package -help` to view options.
- Try different options and ensure that they work as expected.

Review the table below and look for untested scenarios. 
Basically, try to break it and let us know if you find any issues! :smile:

#### Status

`wails3 dev` command:

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("alpha4-wails3-dev.csv") }}

- Windows is partially working:
  - Frontend changes work as expected
  - Go changes cause the application to be built twice

- Mac is partially working:
  - Frontend changes work as expectedFrontend changes work as expected
  - Go changes terminates the vite process

`wails3 package` command:

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working
- :material-cancel: - Not Supported

{{ read_csv("alpha4-wails3-package.csv") }}


### Alpha 3 - Completed 2024-01-14

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
