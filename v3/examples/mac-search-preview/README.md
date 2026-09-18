# Search and preview

This example exercises three macOS integration managers on `application.App`:

- `app.Spotlight`: indexes three sample notes as `SearchableItem`s in the
  domain `com.wails.mac-search-preview.notes`, removes them again with
  `DeleteDomain`, and receives Spotlight selections through `OnOpen`.
- `app.QuickLook`: opens the Quick Look panel on one note or on every sample
  file (`Preview`), closes it (`ClosePreview`), and renders thumbnails of the
  sample files with `Thumbnail` (shown in the page as data URLs; toggle icon
  mode to see the Finder-style rendering).
- `app.AppleEvents`: handles the custom `WAIL`/`note` event, which opens a
  note by ID or title and replies with its text, and shows the generated
  scripting definition (`ScriptingDefinition`). A form sends the same event
  with `Send` to any bundle identifier.

The sample files are written to a temporary directory at startup and removed
on shutdown, together with the Spotlight domain.

## Running

```bash
go run .
```

Then:

- Click **Index in Spotlight**, open Spotlight (Cmd+Space) and type
  "Beta feedback". Spotlight only lists items from bundled applications, so
  with `go run` the call succeeds but nothing appears; see below for the
  bundled workflow.
- Click **Preview** on a note, or **Preview all** to page through every file
  with the arrows in the panel. **Close** dismisses it.
- Thumbnails are generated as the page loads. Tick **Icon mode** to
  regenerate them with the document border.

## Driving the app with AppleScript

The handler is registered for event class `WAIL` and ID `note`. Scripts can
send it using the raw event syntax; the direct parameter is the note ID or
title and the reply is the note body:

```bash
osascript -e 'tell application id "com.wails.mac-search-preview" to «event WAILnote» "alpha"'
```

Targeting by name or bundle identifier needs the running instance to be a
bundled application, so package it first (for example with
`wails3 task package` in a generated project) with these keys in
`Info.plist`:

```xml
<key>CFBundleIdentifier</key>
<string>com.wails.mac-search-preview</string>
<key>NSAppleScriptEnabled</key>
<true/>
<key>OSAScriptingDefinition</key>
<string>mac-search-preview.sdef</string>
```

`OSAScriptingDefinition` names the sdef file inside `Contents/Resources`.
Write the output of `app.AppleEvents.ScriptingDefinition()` (the page shows
it) to that file and Script Editor's File > Open Dictionary lists the
`note` command, so the script becomes:

```applescript
tell application "mac-search-preview" to note "alpha"
```

Unknown notes reply with error -1728 (no such object), which `osascript`
prints as a script error.

The **Send** form calls `app.AppleEvents.Send(bundleID, "WAIL", "note", text)`
so the same event can be sent from Go to another running application. A
bundled build can target itself with its own bundle identifier; sending to
other applications additionally needs the `NSAppleEventsUsageDescription`
Info.plist key and the user's consent under System Settings > Privacy &
Security > Automation.

## Spotlight results

`Index` reports success from an unbundled binary but macOS only surfaces
items indexed by a bundled application, so run the packaged app to see the
notes in Spotlight. Selecting one delivers a
`CSSearchableItemActionType` user activity to
`application:continueUserActivity:`, which the application delegate routes
to `app.Spotlight.OnOpen`; the example then emits `note:opened` to the page
and focuses the window.
