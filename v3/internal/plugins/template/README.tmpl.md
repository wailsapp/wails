# {{.Name}} Plugin

This example plugin provides a way to generate hashes of strings.

## Installation

Add the plugin to the `Plugins` option in the Applications options:

```go
    Plugins: map[string]application.Plugin{
        "{{.Name}}": {{.Name}}.NewPlugin(),
    },
```

## Usage

You can then call the methods from the frontend:

```js
    wails.Plugin("{{.Name}}","All","hello world").then((result) => console.log(result))
```

This method returns a struct with the following fields:

```typescript
    interface Hashes {
        MD5: string;
        SHA1: string;
        SHA256: string;
    }
```

A TypeScript definition file is provided for this interface.

## Support

If you find a bug in this plugin, please raise a ticket [here](https://github.com/plugin/repository). 
Please do not contact the Wails team for support.