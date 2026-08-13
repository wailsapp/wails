# linux-flatpak-appid (fix side)

Fix verification for the startup crash reproduced by the example of the same
name in the unpatched checkout. Same manifest, same app id, same launch path.
`main.go` differs by one option:

```go
Linux: application.LinuxOptions{
    ApplicationID: "com.example.WailsFlatpakAppId",
},
```

Unlike the GTK4 uniqueness fix, this one is opt-in. The patched tree alone
changes nothing — an app that does not set `ApplicationID` still gets
`org.wails.<name>` and still crashes, which is what keeps the patch backward
compatible.

## Test

```sh
cd v3/examples/linux-flatpak-appid
wails3 task bundle
flatpak install --user ./flatpak/com.example.WailsFlatpakAppId.flatpak
```

Then launch **Wails Flatpak App ID** from the desktop with `journalctl --user -f`
open. The window should open and the journal should be quiet.

On the unpatched tree the same steps produce:

```
Portal call failed: Invalid sandbox a11y own name:
  'org.wails.linux-flatpak-appid.Sandboxed.WebProcess-<uuid>' doesn't match app id
SIGTRAP: trace trap
```

Both bundles carry the same app id, so install one at a time — installing either
replaces the other, which is what makes them directly comparable.

Launching from the menu matters. The web process only makes the portal call when
the accessibility bus is reachable inside the sandbox, so a terminal launch can
pass on the unpatched build and prove nothing.

## The fix

`Options.Linux.ApplicationID` overrides the id the GtkApplication is built with,
defaulting to the derived `org.wails.<name>` when unset. Matching it to the
flatpak app id makes the accessibility name WebKit requests fall under the app
id prefix the sandbox permits.

It also lets the desktop match windows to the installed `.desktop` entry, which
the derived id never does.
