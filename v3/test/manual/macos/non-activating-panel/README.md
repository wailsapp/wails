# Non-activating NSPanel manual test

Run on macOS from `v3`:

```sh
go run ./test/manual/macos/non-activating-panel
```

The left window is a normal Wails `NSWindow`. The floating window is an opt-in
`NSPanel` using `NSWindowStyleMaskNonactivatingPanel`.

## Checks

1. Activate another app, such as TextEdit, then click the panel and its action
   button. The other app must remain active and keep its menu bar.
2. Type in the panel input. Text entry and selection must work. The button must
   update the status exactly once per click, and **Actions** must increase by
   exactly one.
3. Press Escape while the input is focused. The Wails key binding must hide the
   panel without inserting text or firing twice. Show the panel again and
   verify **Escape callbacks** and **Hide events** each increased by one.
4. Switch back to the normal control window. It must activate Wails and behave
   as an ordinary main window. **Main events** must remain zero; the panel must
   never emit a will-become-main or did-become-main event.
5. Activate and deactivate Wails repeatedly. The panel must stay visible; an
   `NSPanel` normally hides on deactivation, so this verifies the Wails default.
6. Use **Close panel**, then **Recreate panel**, several times. Quit with either
   window visible and hidden. There must be no crash, stale pointer, or duplicate
   event handling. The four counters must continue updating on every recreated
   panel.
7. Move between Spaces. The panel should appear on every Space and stay fixed.
   Put another app in fullscreen and confirm the panel can appear over it.
8. Verify the normal window remains the control: it participates in ordinary
   activation, window cycling, and main-window behaviour.

The panel deliberately uses `MacWindowLevelFloating`. To test a higher window
level, change it to `MacWindowLevelPopUpMenu`; window level is independent of
whether an `NSPanel` activates the application.
