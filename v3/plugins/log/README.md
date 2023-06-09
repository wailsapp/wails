# log Plugin

This example plugin provides a way to generate hashes of strings.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
    Plugins: map[string]application.Plugin{
        "log": log.NewPlugin(),
    },
```

## Usage

You can then call the methods from the frontend:

```js
    wails.Plugin("log","Debug","hello world")
```

### Methods

- Trace
- Debug
- Info
- Warning
- Error
- Fatal
- SetLevel

SetLevel takes an integer value from JS:

```js
    wails.Plugin("log","SetLevel",1)
```

Levels are:

 - Trace: 1 
 - Debug: 2
 - Info: 3
 - Warning: 4
 - Error: 5
 - Fatal: 6

## Support

If you find a bug in this plugin, please raise a ticket [here](https://github.com/plugin/repository). 
Please do not contact the Wails team for support.