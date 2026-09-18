# mac-windows-extra

Four macOS platform features on one window: a window-modal sheet, popovers
anchored to a toolbar item, to a rectangle in the page and to the system
tray, application presentation options (kiosk mode), and window state
restoration with a counter that survives a relaunch.

Everything here is macOS only. On other platforms the calls are documented
no-ops or return `ErrMacOnly`, `ErrMacSheetUnsupported` or
`ErrMacPopoverUnsupported`.

## Running

From the `v3` directory:

```bash
GOWORK=off go run ./examples/mac-windows-extra
```

The example is a plain package inside the Wails module, so it needs no
binding generation or frontend build.

## What to expect

**Sheet.** The Sheet toolbar button presents a second WebView window as a
sheet (`WebviewWindow.PresentSheet`). The sheet window is created with
`Frameless: false` and `Hidden: true`: AppKit dictates a titled window and
hides its title bar buttons while it is attached, and creating it hidden
stops it flashing on screen before `beginSheet:` runs. Its OK and Cancel
buttons call `EndSheet` with a response code, which arrives in the sheet
window's `OnSheetEnd` callback and is shown on the page.
`PresentCriticalSheet` is the same call for sheets that must not queue
behind an existing one; `SheetParent`, `AttachedSheet` and `IsSheet` answer
the AppKit queries.

**Popovers.** `NewMacPopover` builds an `NSPopover` whose content is a
native `MacAccessory` strip: a label, a symbol button that increments the
counter and a Close button. The same popover is shown from the Popover
toolbar item (`MacToolbarItem.ShowPopover`) and from the button in the page
(`MacPopover.ShowRelativeTo` with a rectangle in window points, attached to
its top edge). The tray icon shows a second popover
(`SystemTray.ShowPopover`) with Show Window and Quit buttons, the native
replacement for a window positioned under the status item. Transient
popovers close on an outside click; `OnClose` reports it.

A popover hosts native controls. Hosting a WebView page in a popover is a
follow-up: the window's WebView cannot be reparented safely and a second
WebView needs its own runtime bridge.

**Presentation options.** The Kiosk toolbar button and View > Toggle Kiosk
Presentation (Cmd-K) switch `App.SetPresentationOptions` between the default
and a kiosk set that auto-hides the Dock and menu bar and disables
Command-Tab and Hide. Force Quit stays enabled. View > Invalid Combination
asks to hide the menu bar without the Dock, which `Validate` rejects in Go
before AppKit would raise. The options are reset when the application
shuts down; `MacOptions.PresentationOptions` applies a set at launch.

**State restoration.** The main window is created with
`Mac.RestorationID: "main"` and stores the counter with
`SetRestorationData` on every increment. `app.Window.OnRestore` recreates it
from that data, and the default window is only created when
`app.Window.WillRestoreWindows()` is false, so a relaunch shows one window.

macOS restores windows after a crash, a forced quit or a reboot with
"Reopen windows when logging back in", and after a normal quit only when
System Settings > Desktop & Dock > "Close windows when quitting an
application" is off. To see it after a normal quit on a box with the
default setting, allow it for this example only:

```bash
defaults write mac-windows-extra NSQuitAlwaysKeepsWindows -bool true
```

(`mac-windows-extra` is the defaults domain of the unbundled binary; use
the bundle identifier for a bundled build.) Increment the counter, quit
with Cmd-Q, relaunch: the window comes back with the counter and logs
"restoring the main window". Remove the override afterwards with
`defaults delete mac-windows-extra NSQuitAlwaysKeepsWindows`.

`MacOptions.SupportsSecureRestorableState` is set, opting in to secure
coding for the saved state; the delegate method it feeds is always
implemented, so the macOS 14 "Secure coding is automatically enabled for
restorable state" warning never appears either way.
