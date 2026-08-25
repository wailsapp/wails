# The Wails CLI

The Wails CLI is a command line tool that allows you to create, build and run Wails applications.
There are a number of commands related to tooling, such as icon generation and asset bundling.

## MCP project server

`wails3 mcp` starts a local Model Context Protocol server for agent-assisted project
management. It is separate from the experimental MCP server compiled into a running
Wails application: this server manages the project lifecycle, while the application
server controls the running WebView.

Transport is selected automatically. When an MCP host launches the command with
piped stdin/stdout, Wails uses stdio. When a person runs it in a terminal, Wails
uses Streamable HTTP on loopback and an OS-selected free port. Use `-stdio` or
`-http` to override this behavior; `-port 0` selects a free loopback port.

The server is confined to the current working directory by default. An explicit root
may be supplied with `-root`; paths and symlinks that escape that root are rejected.
Mutating and process-control tools require a per-process session token. If `-token`
and `WAILS_MCP_TOKEN` are both absent, a cryptographically random token is generated
and disclosed in the MCP initialize instructions (and on stderr for human operators).

```bash
cd /path/to/project-or-workspace
wails3 mcp
```

In terminal mode, the command prints the endpoint and bearer token to stderr. In
stdio mode, the token is included in the MCP initialize instructions.

Example agent configuration:

```json
{
  "mcpServers": {
    "wails3": {
      "command": "wails3",
      "args": ["mcp"]
    }
  }
}
```

Available tools include project inspection and initialization, diagnostics, task
listing and execution, builds, development-mode startup, binding generation, and
the `wails_job_status` and `wails_job_stop` tools; status returns bounded output.
The server does not expose arbitrary shell execution. Remote templates and Git
remotes require `allowExternal=true`; signing,
publishing, and deployment are not exposed and should be approved explicitly by the
user when added in a future capability.

## Commands

### task

The `task` command is for running tasks defined in `Taskfile.yml`. It is a wrapper around [Task](https://taskfile.dev).

### generate

The `generate` command is used to generate resources and assets for your Wails project.
It can be used to generate many things including: 
  - application icons, 
  - resource files for Windows applications
  - Info.plist files for macOS deployments

#### icon

The `icon` command generates icons for your project. 

| Flag               | Type   | Description                                          | Default              |
|--------------------|--------|------------------------------------------------------|----------------------|
| `-example`         | bool   | Generates example icon file (appicon.png)            |                      |
| `-input`           | string | The input image file                                 |                      |
| `-sizes`           | string | The sizes to generate in .ico file (comma separated) | "256,128,64,48,32,16" |
| `-windowsFilename` | string | The output filename for the Windows icon             | icon.ico             |
| `-macFilename`     | string | The output filename for the Mac icon bundle          | icons.icns           |

```bash
wails3 generate icon -input myicon.png -sizes "32,64,128" -windowsFilename myicon.ico -macFilename myicon.icns       
```

This will generate icons for mac and windows and save them in the current directory as `myicon.ico`
and `myicons.icns`.

#### syso

The `syso` command generates a Windows resource file (aka `.syso`).

```bash
wails3 generate syso <options>
```

| Flag        | Type   | Description                                | Default          |
|-------------|--------|--------------------------------------------|------------------|
| `-example`  | bool   | Generates example manifest & info files    |                  |
| `-manifest` | string | The manifest file                          |                  |
| `-info`     | string | The info.json file                         |                  |
| `-icon`     | string | The icon file                              |                  |
| `-out`      | string | The output filename for the syso file      | `wails.exe.syso` |
| `-arch`     | string | The target architecture  (amd64,arm64,386) | `runtime.GOOS`   |

If `-example` is provided, the command will generate example manifest and info files
in the current directory and exit.

If `-manifest` is provided, the command will use the provided manifest file to generate
the syso file.

If `-info` is provided, the command will use the provided info.json file to set the version
information in the syso file.

NOTE: We use [winres](https://github.com/tc-hib/winres) to generate the syso file. Please
refer to the winres documentation for more information.

NOTE: Whilst the tool will work for 32-bit Windows, it is not supported. Please use 64-bit.

#### defaults

```bash
wails3 generate defaults      
```
This will generate all the default assets and resources in the current directory.

#### bindings

```bash
wails3 generate bindings
```

Generates bindings and models for your bound Go methods and structs.
