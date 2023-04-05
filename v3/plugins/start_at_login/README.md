# start_at_login Plugin

This example plugin provides a way to generate hashes of strings.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
    Plugins: map[string]application.Plugin{
        "start_at_login": start_at_login.NewPlugin(),
    },
```

## Usage

You can then call the methods from the frontend:

```js
    wails.Plugin("start_at_login","StartAtLogin", true).then((result) => console.log(result))
    wails.Plugin("start_at_login","IsStartAtLogin").then((result) => console.log(result))
```

To use this from Go, create a new instance of the plugin, then call the methods on that:

```go
    start_at_login := start_at_login.NewPlugin()
	start_at_login.StartAtLogin(true)
```

## Support

If you find a bug in this plugin, please raise a ticket [here](https://github.com/plugin/repository). 
Please do not contact the Wails team for support.