# Splash screen lifecycle example

This example demonstrates the `ApplicationStarting` and
`ApplicationInitialized` lifecycle events. It shows a responsive splash window
while an unchanged `ServiceStartup` implementation deliberately waits for five
seconds, then closes the splash and reveals the main window.

Run it from the `v3` directory:

```sh
go run ./examples/splashscreen
```
