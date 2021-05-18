# Build

The build command processes the Wails project and generates an application binary. 

## Usage

`wails build <flags>`

### Flags

| Flag           | Details      | Default |
| :------------- | :----------- | :------ |
| -clean | Clean the bin directory before building   | |
| -compiler path/to/compiler  | Use a different go compiler, eg go1.15beta1 | go |
| -ldflags "custom ld flags" | Use given ldflags | | 
| -o path/to/binary | Compile to given path/filename | |
| -k | Keep generated assets | |
| -package | Create a platform specific package | |
| -production | Compile in production mode: -ldflags="-w -s" + "-h windows" on Windows | |
| -tags | Build tags to pass to Go compiler (quoted and space separated) | |
| -upx | Compress final binary with UPX (if installed) | |
| -upxflags "custom flags" | Flags to pass to upx | |
| -v int | Verbosity level (0 - silent, 1 - default, 2 - verbose) | 1 |

## The Build Process

The build process is as follows:

  - The flags are processed, and an Options struct built containing the build context.
  - The type of target is determined, and a custom build process is followed for target.

### Desktop Target 

  - The frontend dependencies are installed. The command is read from the project file `wails.json` under the key `frontend:install` and executed in the `frontend` directory. If this is not defined, it is ignored.
  - The frontend is then built. This command is read from the project file `wails.json` under the key `frontend:install` and executed in the `frontend` directory. If this is not defined, it is ignored.
  - The project directory is checked to see if the `build` directory exists. If not, it is created and default project assets are copied to it.
  - An asset bundle is then created by reading the `html` key from `wails.json` and loading the referenced file. This is then parsed, looking for local Javascript and CSS references. Those files are in turn loaded into memory, converted to C data and saved into the asset bundle located at `build/assets.h`, which also includes the original HTML.
  - The application icon is then processed: if there is no `build/appicon.png`, a default icon is copied. On Windows, an `app.ico` file is generated from this png. On Mac, `icons.icns` is generated.
  - If there are icons in the `build/tray` directory, these are processed, converted to C data and saved as `build/trayicons.h`, ready for the compilation step.
  - If there are icons in the `build/dialog` directory, these are processed, converted to C data and saved as `build/userdialogicons.h`, ready for the compilation step.
  - If the `-package` flag is given for a Windows target, the Windows assets in the `build/windows` directory are processed: manifest + icons compiled to a `.syso` file (deleted after compilation).
  - If we are building a universal binary for Mac, the application is compiled for both `arm64` and `amd64`. The `lipo` tool is then executed to create the universal binary.
  - If we are not building a universal binary for Mac, the application is built using `go build`, using build tags to indicate type of application and build mode (debug/production).
  - If the `-upx` flag was provided, `upx` is invoked to compress the binary. Custom flags may be provided using the `-upxflags` flag.
  - If the `package` flag is given for a non Windows target, the application is bundled for the platform. On Mac, this creates a `.app` with the processed icons, the `Info.plist` in `build/darwin` and the compiled binary.

### Server Target

TBD

### Hybrid Target

TBD

