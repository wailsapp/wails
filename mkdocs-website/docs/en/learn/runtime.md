# Runtime

The Wails runtime is the standard library for Wails applications. It provides a number of features that may
be used in your applications, including:

- Window management
- Dialogs
- Browser integration
- Clipboard
- Frameless dragging
- Tray icons
- Menu management
- System information
- Events
- Calling Go code
- Context Menus
- Screens
- WML (Wails Markup Language)

The runtime is required for integration between Go and the frontend. There are 2 ways to integrate the runtime:

- Using the `@wailsio/runtime` package
- Using a pre-built version of the runtime

## Using the `@wailsio/runtime` package

The `@wailsio/runtime` package is a JavaScript package that provides access to the Wails runtime. It is used in by all
the standard templates and is the recommended way to integrate the runtime into your application. By using the package,
you will only include the parts of the runtime that you use.

The package is available on npm and can be installed using:

```shell
npm install --save @wailsio/runtime
```

## Using a pre-built version of the runtime

Some projects will not use a Javascript bundler and may prefer to use a pre-built version of the runtime. This is
the default for the examples in `v3/examples`. The pre-built version of the runtime can be generated using the
following command:

```shell
wails3 generate runtime
```

This will generate a `runtime.js` (and `runtime.debug.js`) file in the current directory.
This file can be used by your application by adding it to your assets directory (normally `frontend/dist`) and then including it in your HTML:

```html
<html>
    <head>
        <script src="/runtime.js"></script>
    </head>
    <!--- ... -->
</>
```
