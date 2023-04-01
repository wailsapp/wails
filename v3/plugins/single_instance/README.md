# single-instance Plugin

This example plugin provides a way to generate hashes of strings.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
    Plugins: map[string]application.Plugin{
        "single_instance": single_instance.NewPlugin(&single_instance.Config{
            // When true, the original app will be activated when a second instance is launched
            ActivateAppOnSubsequentLaunch: true,
        }
    },
```

## Usage

This plugin prevents the launch of multiple copies of your application. 
If you set `ActivateAppOnSubsequentLaunch` to true the original app will be activated when a second instance is launched.

## Support

If you find a bug in this plugin, please raise a ticket [here](https://github.com/plugin/repository). 
Please do not contact the Wails team for support.

## Credit

This plugin contains modified code from the awesome [go-singleinstance](https://github.com/allan-simon/go-singleinstance) module (c) 2015 Allan Simon.
Original license file has been renamed `go-singleinstance.LICENSE` and is available [here](./singleinstance_LICENSE).