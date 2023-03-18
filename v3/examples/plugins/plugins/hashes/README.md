# Hashes Plugin

This example plugin provides a way to generate hashes of strings.

## Usage

Add the plugin to the `Plugins` option in the Applications options:

```go
    Plugins: map[string]application.Plugin{
        "hashes": hashes.NewPlugin(),
    },
```

You can then call the Generate method from the frontend:

```js
    wails.Plugin("hashes","Generate","hello world").then((result) => console.log(result))
```

This method returns a struct with the following fields:

```typescript
    interface Hashes {
        md5: string;
        sha1: string;
        sha256: string;
    }
```

A TypeScript definition file is provided for this interface.
