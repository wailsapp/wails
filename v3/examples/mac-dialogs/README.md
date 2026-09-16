# Native macOS dialogs and panels

This example exercises the macOS extras on the dialog API. Every button on
the page (and every entry in the Dialogs menu) calls a bound Go method that
shows a native AppKit alert, panel or picker and reports what the user did.

It demonstrates:

- an alert with a suppression checkbox (`SetSuppression`, `OnSuppression`,
  `Suppressed`) and a help button (`SetHelp`);
- a text prompt as a sheet (`app.Dialog.Prompt` with `Window`) and a secure
  modal prompt (`Secure: true`);
- an open panel that accepts any file conforming to `public.image` next to an
  extension filter for PDF (`AddContentType` alongside `AddFilter`);
- a save panel with a Format pop-up accessory that swaps the allowed type and
  the name field's extension (`SetFormats`, `SelectedFormat`), a custom name
  field label (`SetNameFieldLabel`) and Finder tags (`SetTags`);
- the system colour panel (`app.Dialog.PickColor`) and font panel
  (`app.Dialog.PickFont`), both with live `OnChange` updates pushed to the
  page as events.

## Running

```shell
go run .
```

The blocking helpers (`Prompt`, `PickColor`, `PickFont`) wait on a Go channel
until the panel is dismissed, so call them from a goroutine (bound service
methods and menu callbacks already are). They return
`application.ErrDialogNotSupported` on platforms without a native
implementation; the example still builds and runs there, and the buttons
report the error.
