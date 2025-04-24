# Notifications Example

This is an example of the Notifications Service.

## Running

### Windows

```sh
cd examples/notifications
go run main.go
```

### macOS

macOS requires a bundle to be built for notifications to work correctly:

```sh
wails3 package
```
Then run the application built in the `bin` directory.