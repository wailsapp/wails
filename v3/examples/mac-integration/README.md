# Platform integration (macOS)

This example exercises the system integration managers on `application.App`:

- `app.ServicesProvider`: a "Summarise with Wails Example" entry in the
  Services menu of every application, taking the selected text and writing a
  summary back.
- `app.Activity`: a Handoff activity that mirrors a text field
  (`Publish`, `PublishedActivity.Update`, `Invalidate`), plus `OnWillContinue`,
  `OnContinue` and `OnFailed` for activities and universal links handed to
  the app.
- `app.Browser`: `ApplicationsForFile`, `OpenWith` and `ActivateApplication`.
- `app.Env`: `DefaultHandler` for a content type or URL scheme.

## Running

```bash
go run .
```

Then:

- Type into the **Handoff activity** field. The first change publishes an
  `NSUserActivity`; later changes update its `userInfo`. With Handoff enabled
  on a nearby device signed into the same iCloud account, the activity appears
  in its Dock or App Switcher once the app is bundled (see below).
- Enter a file path and click **Applications for file**; pick **Open with** or
  **Activate** on a row, or use the bundle identifier field (defaults to
  TextEdit).
- Click **Default app for content type** or **Default app for scheme** to see
  which application currently owns `public.plain-text` or `mailto`.

Everything is logged at the bottom of the page.

## Bundling

The Services menu, Handoff and universal links only work for a bundled
application with the right Info.plist keys (and, for universal links, an
entitlement). `go run` binaries have no Info.plist, so those three features
stay silent until you build a bundle.

### NSServices (Services menu)

The service registered at runtime must be declared in Info.plist. The CLI
renders the `services` list from `build/config.yml` into the generated
`build/darwin/Info.plist` and `Info.dev.plist` when you run
`wails3 task common:update:build-assets`:

```yaml
services:
  - name: SummariseWithWailsExample   # ServiceDefinition.Name
    menuTitle: Summarise with Wails Example
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S                    # optional, Cmd+Shift+S
    # portName: My App                  # optional, defaults to the product name
```

`portName` must equal the bundle's `CFBundleName`, which is the `Name` passed
to `application.New`. For a hand-written Info.plist,
`app.ServicesProvider.InfoPlistXML()` returns the exact XML to paste (the page
shows it), and `InfoPlistEntries()` returns the same data as maps for a plist
serialiser. Run `/System/Library/CoreServices/pbs -update` after installing a
new bundle if the entry does not appear in the Services menu straight away.

### NSUserActivityTypes (Handoff)

Every activity type the app publishes or continues must be listed:

```xml
<key>NSUserActivityTypes</key>
<array>
    <string>com.wails.example.mac-integration.editing</string>
</array>
```

Add this to `build/darwin/Info.plist` and `Info.dev.plist`; the CLI keeps
keys it does not manage when it regenerates those files.

### Universal links

Universal links reach the app as `NSUserActivityTypeBrowsingWeb` activities.
They are delivered twice on purpose: `app.Activity.OnContinue` receives the
activity with `WebpageURL` set, and `events.Common.ApplicationLaunchedWithUrl`
fires with the same URL, so an app that already handles a custom URL scheme
needs no second code path.

They require the associated-domains entitlement:

```xml
<key>com.apple.developer.associated-domains</key>
<array>
    <string>applinks:example.com</string>
</array>
```

`wails3 setup entitlements` writes `build/darwin/entitlements.plist` and
`entitlements.dev.plist` for boolean entitlements; associated domains take an
array, so add the key above to those files by hand. The app must be signed
with a provisioning profile that includes the capability, and
`https://example.com/.well-known/apple-app-site-association` must list the
app's team and bundle identifiers.

## Notes

- Service handlers run on the main thread while the requesting application
  waits, so keep them quick and hand long work to a goroutine.
- `OpenWith` accepts a bundle identifier or a `.app` path and returns once the
  target application has accepted the request.
- `app.Env.SetDefaultHandler` (not wired to a button, since it changes system
  settings) needs a bundle identifier too; macOS asks the user to confirm
  when a browser or mail client is being replaced.
