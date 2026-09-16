# macOS native feature audit

**Branch:** `feat/macos-window-chrome` at `64adeb2f5` (rebased onto master 2026-09-16)
**Date:** 2026-09-16
**Question answered:** which macOS-native capabilities could Wails v3 expose to Go developers that this branch does not expose today?

## 1. Summary

The branch already gives Wails a real AppKit application structure: `NSToolbar`, an `NSSplitViewController` with semantic sidebar, primary and inspector items, an `NSOutlineView` source list, a typed inspector, a native sharing provider, `NSPanel` and notch windows, window tabbing, Liquid Glass, and an experimental WebView-free `NativeWindow` with an `NSTextView` editor.

This audit found **85 candidate features** across nine areas that AppKit or WebKit provides on public API and Wails does not yet expose. They fall into three bands:

| Band | Count | What it means |
|---|---|---|
| **Tier 1: finish the chrome model** | 14 | Completes APIs that exist on this branch but stop short (toolbar, sidebar, inspector, accessories, tabs). Mostly additive Go surface over Objective-C that is already written. |
| **Tier 2: platform integration** | 48 | Things every polished Mac app does that need no new window architecture: dock, menus, dialogs, clipboard, drag-out, sounds, permissions, sleep control, Services, Quick Look. Several are open GitHub requests. |
| **Tier 3: native content and deeper WebKit** | 23 | New content surfaces (PDFKit, Quick Look, AVKit, tables) and public WebKit APIs that replace private hooks or unlock downloads, snapshots, profiles and extensions. |

Two findings deserve attention before any of this is built:

- **Titlebar accessories lost their constructor.** Commit `8531d4501` added `AddTitlebarAccessory`, but at the time of the audit only the `WrapMacAccessoryViewController(unsafe.Pointer)` wrapper survived. Resolved the same day by `MacAccessory` (Tier 1, item 1).
- **The `server` cross-build was broken.** `CGO_ENABLED=0 GOOS=linux go build -tags server ./pkg/application` failed because `webview.NewRequest` is unavailable without cgo. Resolved the same day by moving the constructors to a darwin-only file.

Three items on this list are already tracked as issues: window cascading (#4762), drag-out to Finder and other apps (#4648) and the cross-platform `Permissions` option being ignored on macOS (#6067). Fn/Globe key handling (#6033) and ProMotion frame rate (#6056) are bugs rather than features and are not scored here.

## 2. How the audit was done

1. Inventoried every macOS-facing Go type, constructor and method in `v3/pkg/application`, `v3/pkg/services`, `v3/pkg/events` and the darwin Objective-C files on this branch, including the four WEP proposal directories and `review.md`.
2. Grepped the tree for the AppKit, WebKit, Foundation and framework symbols each candidate needs, to confirm a zero-hit before listing it as a gap.
3. Pulled the open GitHub issues that mention macOS to weigh demand.
4. Scored each candidate on value, effort and fit with the branch's existing model.

Effort scale: **S** under a day of native plus Go work, **M** a few days, **L** a week or more or needs a design decision first.

## 3. What exists today, in one screen

| Area | Exposed on this branch |
|---|---|
| App options | `ActivationPolicy`, `ApplicationShouldTerminateAfterLastWindowClosed`, `NativeOnly`, `FileAssociations` |
| Window options | Backdrop (normal, transparent, translucent, Liquid Glass), corners, titlebar presets and toolbar style, content layout, appearance, window level, collection behaviour, tabbing mode, panel class and preferences, WebView preferences |
| Toolbar | Buttons, search, share, item groups, flexible space, sidebar and inspector toggles and tracking separators; SF Symbols, badges, tint, prominence, display mode |
| Split view | Sidebar, primary WebView, text editor and inspector panes; thickness, priority, collapse, autosave name |
| Sidebar | One section level, SF Symbol rows, single selection, enabled, hidden, tooltip, click |
| Inspector | Label, text field, checkbox, popup |
| Native window | `NativeWindow` with toolbar and split view; plain-text `MacTextEditor` |
| Menus | 48 roles, submenus, checkbox, radio, key equivalents, tooltips, bitmap icons, Services and Window menus, context menus |
| Dialogs | `NSAlert` modal or sheet with custom buttons and icon; open and save panels with extension filters, sheets, directory, hidden files |
| System tray | Icon, template icon, dark mode icon, label, menu, click handlers, attached window positioning |
| Services | Dock badge and icon visibility, UNUserNotificationCenter notifications with actions and categories |
| Managers | Autostart (SMAppService or LaunchAgent), single instance with URL scheme capture, global shortcuts (Carbon hot keys), clipboard (text), browser, environment (dark mode, accent colour, reveal in Finder), screens |
| Events | 137 `events.Mac.*` values covering app, menu, WebView navigation and NSWindow lifecycle |
| Other | Drop-in file drag, `Print()` with fixed settings, `Flash()` informational bounce, reduce-motion read |

## 4. Gap catalogue

Each row gives the AppKit or framework API, why it matters, effort, and a sketch of the Go surface that fits the branch's style.

### 4.1 Tier 1: finish the chrome model

**Status, 2026-09-16: all 14 items implemented** on this branch, together with the three review blockers (runtime split construction, the server cross-build, and the native window option contract). One deviation from the sketches: split-item accessories require macOS 26, not 11, because `NSSplitViewItemAccessoryViewController` is a macOS 26 API. The docs guide `guides/macos-native-chrome.mdx` covers the resulting surface.

These are the highest-leverage items because the Objective-C host code, callback plumbing and snapshot model already exist. Most rows here also appear in `review.md`'s "recommended next work".

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 1 | **Titlebar and split-item accessories with Go constructors** | `NSTitlebarAccessoryViewController`, `NSSplitViewItem.addTopAlignedAccessoryViewController` | Restores the feature that was removed before HEAD. Finder-style filter strips above a sidebar and status labels in the titlebar are common. The existing `MacAccessoryViewController` wrapper becomes usable without unsafe pointers. | M | `toolbar.AddTitlebarAccessory(kind, layout)` and `pane.AddTopAccessory(...)` returning a typed handle with search, segmented, button, label and native-escape-hatch variants |
| 2 | **Menu toolbar item** | `NSMenuToolbarItem` | Dropdown actions in a toolbar. Reuses the existing `Menu` model, so the API is one constructor. | S | `toolbar.AddMenu(label, *Menu)` with `SetShowsIndicator` |
| 3 | **Standard toolbar identifiers** | `NSToolbarSpaceItemIdentifier`, `NSToolbarToggleSidebarItemIdentifier` variants, `centeredItemIdentifiers` | Fixed spaces and centred items are needed for Safari-style layouts. | S | `AddSpace()`, `toolbar.SetCentered(item)` |
| 4 | **Toolbar visibility priority and navigational items** | `NSToolbarItem.visibilityPriority`, `isNavigational` | Controls which items collapse into the overflow menu first; navigational items pin to the leading edge on macOS 11+. | S | `item.SetVisibilityPriority(p)`, `item.SetNavigational(bool)` |
| 5 | **Live toolbar mutation** | `insertItemWithItemIdentifier:atIndex:`, `removeItemAtIndex:` | Today `Add*` after attachment only changes the Go model and the live toolbar goes stale silently (review P2). | M | Make `Add*`, `Remove(item)` and `Move(item, index)` apply to the live `NSToolbar` |
| 6 | **User customisation and autosave** | `allowsUserCustomization`, `autosavesConfiguration`, `identifier` | The standard "Customize Toolbar..." sheet. Needs an optional stable persistence key per item. | M | `toolbar.SetCustomizable(persistKey)`, `item.SetPersistKey(key)`; keep generated IDs when not enabled |
| 7 | **Search recents and search menu** | `NSSearchField.recentsAutosaveName`, `searchMenuTemplate`, `sendsSearchStringImmediately` | The search item exists but has no recents, no scope menu and no incremental-versus-submit control. | S | `item.SetSearchRecentsKey(key)`, `SetSearchMenu(*Menu)`, `SetSearchIncremental(bool)` |
| 8 | **Content-list pane** | `NSSplitViewItem contentListWithViewController:` plus `NSTableView` | The missing middle column of Finder, Mail and every document browser. Needs titles, subtitles, badges, multi-selection and sorting. Designed separately from `MacSidebar`. | L | `split.AddContentList(*MacContentList)` with `AddRow`, `SetColumns`, `OnSelectionChange`, `SetSortDescriptors` |
| 9 | **Sidebar hierarchy, badges and accessories** | `NSOutlineView` nested items, `NSTableCellView` badge and trailing views | Mail-style unread counts, nested folders, disclosure state. | M | `item.AddItem(label)`, `item.SetBadge(count)`, `item.SetExpanded(bool)`, `SetAccessorySymbol` |
| 10 | **Sidebar interaction** | `NSOutlineView` drag types, `menuForEvent:`, inline editing via `NSTextField` delegate | Reordering, context menus and rename are expected in any source list. | M | `sidebar.SetReorderable(bool)`, `item.SetContextMenu(*Menu)`, `item.SetEditable(bool)`, `OnRename`, `OnMove` |
| 11 | **Sidebar multi-selection and removal** | `allowsMultipleSelection`, `removeItemsAtIndexes:` | Batch actions and dynamic lists. Today items can only be hidden. | S | `sidebar.SetMultipleSelection(bool)`, `section.Remove(item)`, `sidebar.SelectedItems()` |
| 12 | **Inspector controls** | `NSSlider`, `NSStepper`, `NSSegmentedControl`, `NSColorWell`, `NSDatePicker`, `NSButton`, disclosure via `NSStackView` sections | Four control kinds is not enough for a real inspector. Each follows the existing snapshot and callback pattern. | M | `AddSlider(label, min, max, value)`, `AddStepper`, `AddSegmented(labels)`, `AddColorWell(rgba)`, `AddDatePicker`, `AddButton`, `section.SetCollapsible(bool)` |
| 13 | **Inspector removal and reordering** | Snapshot diff | Sections and controls can be hidden but not removed. | S | `section.Remove(control)`, `inspector.RemoveSection(section)` |
| 14 | **Window tab groups** | `NSWindowTabGroup`, `addTabbedWindow:ordered:`, `selectNextTab`, `toggleTabBar:`, `toggleTabOverview:`, `moveTabToNewWindow:` | Tabbing mode exists but the group is invisible to Go. Requested in #4964, which the option only half answers. | M | `window.AddTab(other)`, `window.TabGroup()` with `Windows()`, `SelectNext()`, `SelectPrevious()`, `ToggleTabBar()`, `ToggleOverview()`, `Detach()` |

### 4.2 Tier 2: platform integration

**Status, 2026-09-17: 47 of 48 items implemented** on this branch across ten commits, each with a runtime probe on macOS 26.4.1 and an example under `v3/examples/mac-*`. The one exception is trackpad gestures (item 26), left for the open macOS 27 preparation PR #5760, which already adds gesture recognisers and would conflict. Two deviations from the sketches: popovers host native accessory controls rather than a WebView, because the WebView configuration is bound to a window delegate (a WebView popover is a follow-up), and the permission kinds carry a `PermissionKind` prefix because `PermissionCamera` and friends already existed as window option values. The guide `guides/macos-platform-integration.mdx` covers the resulting surface.

Features that make an app feel native without touching the window architecture. Several are one-line AppKit calls behind a small Go API.

#### Window behaviour

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 15 | **Represented file and proxy icon** | `representedURL`, `setTitleWithRepresentedFilename:` | Document apps show the file icon in the titlebar; Command-click reveals the path. | S | `window.SetRepresentedFile(path)` |
| 16 | **Document-edited indicator** | `setDocumentEdited:` | The dot in the close button that tells users a save is pending. | S | `window.SetDocumentEdited(bool)` |
| 17 | **Window subtitle** | `NSWindow.subtitle` (macOS 11) | Second line under the title in unified toolbars. | S | `window.SetSubtitle(text)` |
| 18 | **Window cascading** | `cascadeTopLeftFromPoint:` | New document windows should step diagonally. Open issue #4762. | S | `WebviewWindowOptions.InitialPosition = WindowCascade` |
| 19 | **Frame autosave** | `setFrameAutosaveName:`, `setFrameUsingName:` | Windows reopen where the user left them. Complements the split view's `SetAutosaveName`. | S | `MacWindow.FrameAutosaveName` |
| 20 | **State restoration** | `NSWindowRestoration`, `restorableStateKeyPaths`, `NSApplicationDelegate applicationSupportsSecureRestorableState:` | Relaunch after a crash or reboot restores the window set. Also silences the secure-restoration warning that new SDKs log. | M | `MacOptions.SupportsSecureRestorableState`, `window.SetRestorationID(id)`, `OnRestore` |
| 21 | **Sheets** | `beginSheet:completionHandler:`, `beginCriticalSheet:` | Attach a second WebView window or native window as a window-modal sheet. Today only `NSAlert` and file panels can be sheets. | M | `window.PresentSheet(other)`, `other.EndSheet(result)` |
| 22 | **Popovers** | `NSPopover` anchored to a view or rect | Anchored transient content from a toolbar item, sidebar row, status item or a rectangle in the WebView. | M | `toolbarItem.ShowPopover(window)`, `tray.ShowPopover(window)`, `window.ShowPopover(rect, other)` |
| 23 | **Presentation options** | `NSApplication.presentationOptions` | Kiosk mode, hiding the Dock and menu bar, disabling process switching. | S | `MacOptions.PresentationOptions` bitmask |
| 24 | **Critical attention request** | `requestUserAttention:NSCriticalRequest`, `cancelUserAttentionRequest:` | `Flash()` only bounces once. Critical bounces until the app is activated, and the request cannot be cancelled today. | S | `window.Flash(FlashCritical)`, `app.CancelAttentionRequest()` |
| 25 | **Traffic light position** | `standardWindowButton:` frame adjustment on layout | Open issue #4227. The branch's toolbar layout makes the timing reliable. | M | `MacTitleBar.WindowButtonsOffset Point` |
| 26 | **Trackpad gestures** | `NSMagnificationGestureRecognizer`, `NSRotationGestureRecognizer`, `swipeWithEvent:`, `pressureChangeWithEvent:` | Pinch, rotate, swipe and force-click events delivered to Go. The macOS 27 preparation PR #5760 adds gesture recognisers but is not merged. | M | `events.Mac.WindowMagnify`, `WindowRotate`, `WindowSwipe`, `WindowPressure` with delta payloads |

#### Menus and dock

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 27 | **SF Symbols on menu items** | `imageWithSystemSymbolName:` | Toolbar and sidebar items already take symbols; menus still need bitmap bytes. | S | `menuItem.SetSymbol(name)` |
| 28 | **Dock menu** | `applicationDockMenu:` | Right-click actions on the Dock icon. Zero hits in the delegate. | S | `app.Menu.SetDockMenu(*Menu)` |
| 29 | **Open Recent** | `NSDocumentController.noteNewRecentDocumentURL:`, `clearRecentDocuments:` | The system-managed recent list in the File menu, also surfaced in the Dock menu. | S | `app.Menu.AddRecentDocument(path)`, `ClearRecentDocuments()`, `NewRole(OpenRecent)` |
| 30 | **Menu item badges, section headers and palettes** | `NSMenuItemBadge`, `NSMenuItem.sectionHeaderWithTitle:`, `NSMenu.paletteMenuWithColors:` (all macOS 14) | Modern menu affordances: unread counts, labelled groups, colour or symbol pickers. | S | `menuItem.SetBadge(count)`, `menu.AddSectionHeader(title)`, `menu.AddPalette(...)` |
| 31 | **Mixed state, alternates and indentation** | `NSControlStateValueMixed`, `isAlternate`, `indentationLevel` | Indeterminate checkboxes, Option-key alternate items, nested visual grouping. | S | `menuItem.SetMixed()`, `SetAlternate(bool)`, `SetIndentation(level)` |
| 32 | **Dock tile progress and custom view** | `NSDockTile.contentView`, `NSProgressIndicator` inside it | Download and export progress on the Dock icon. The dock service already owns badges. | S | `dock.SetProgress(fraction)`, `dock.ClearProgress()` |

#### Dialogs and panels

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 33 | **Alert suppression checkbox** | `showsSuppressionButton`, `suppressionButton.state` | "Do not show again" is the standard way to stop nagging. | S | `dialog.SetSuppression(label)`, result carries `Suppressed bool` |
| 34 | **Alert accessory view and text prompt** | `NSAlert.accessoryView` with `NSTextField` or `NSSecureTextField` | A native "enter a name" or password prompt without a WebView window. | S | `app.Dialog.Prompt(title, message, placeholder, secure bool)` |
| 35 | **Alert help button** | `showsHelp`, `helpAnchor`, delegate `alertShowHelp:` | Links an alert to documentation. | S | `dialog.SetHelp(func())` |
| 36 | **Content-type filters** | `allowedContentTypes` with `UTType` | Filters today are extension strings converted at runtime. Accepting UTIs directly supports `public.image`, `com.adobe.pdf` and app-defined types. | S | `dialog.AddContentType(uti)` |
| 37 | **Panel accessory view and format popup** | `NSSavePanel.accessoryView`, `nameFieldLabel`, `tagNames`, `prompt` | Export dialogs with a format chooser and Finder tags. | M | `saveDialog.SetFormats(labels, OnChange)`, `SetNameFieldLabel`, `SetTags` |
| 38 | **Colour and font panels** | `NSColorPanel`, `NSFontPanel`, `NSFontManager` | System pickers with the user's swatches and collections. | S | `app.Dialog.PickColor(initial) (RGBA, error)`, `PickFont()` |
| 39 | **Print configuration and silent printing** | `NSPrintInfo` orientation, margins, `showsPrintPanel`, `jobDisposition` | `Print()` hardcodes landscape and modal panels with a TODO to expose config. | S | `window.PrintWithOptions(PrintOptions{Orientation, Margins, Silent, Printer})` |
| 40 | **PDF and image export** | `WKWebView createPDFWithConfiguration:` (macOS 11), `takeSnapshotWithConfiguration:` | The mac-toolbar example hand-writes a PDF in Go to avoid AppKit. WebKit can render the page to PDF or PNG directly. | S | `window.ExportPDF(rect) ([]byte, error)`, `window.Snapshot(rect) ([]byte, error)` |

#### System tray

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 41 | **Tray tooltip** | `NSStatusItem.button.toolTip` | `SetTooltip` is a documented no-op, but the API exists. One line. | S | Implement existing `SetTooltip` |
| 42 | **SF Symbol tray icon** | `imageWithSystemSymbolName:` as template | Avoids shipping PNGs for the menu bar. | S | `tray.SetSymbol(name)` |
| 43 | **Status item behaviour and visibility** | `NSStatusItem.behavior`, `isVisible`, `autosaveName` | Lets users Command-drag the item out of the menu bar and have that choice remembered. | S | `tray.SetRemovable(bool, autosaveKey)` |
| 44 | **Native tray popover** | `NSPopover` from the status button | Replaces the positioned-window approximation with the real thing, including the arrow and auto-dismiss. Shares #22. | M | `tray.ShowPopover(window)` |

#### Input, clipboard, feedback

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 45 | **Rich clipboard** | `NSPasteboard` types for PNG, TIFF, file URLs, RTF, HTML, custom UTIs; `changeCount` | Clipboard is text-only today. Image paste and file copy are common asks. | M | `clipboard.SetImage(png)`, `Image()`, `SetFiles(paths)`, `Files()`, `SetHTML`, `SetData(uti, bytes)`, `OnChange` |
| 46 | **Drag out** | `beginDraggingSessionWithItems:`, `NSFilePromiseProvider` | Drag a file or generated content out of the app to Finder or another app. Open issue #4648. | M | `window.StartDrag(DragItems{Files, Promise func() []byte, Image})` from a JS `dragstart` bridge |
| 47 | **Non-file drop types** | `registerForDraggedTypes:` with text, URL, image types | Drop-in only accepts file paths. | S | `WebviewWindowOptions.DropTypes` |
| 48 | **Haptics** | `NSHapticFeedbackManager` | iOS and Android already have haptics managers; macOS is the gap. | S | `app.Haptics.Perform(kind)` |
| 49 | **Sounds** | `NSSound soundNamed:`, `NSBeep()` | System alert sounds and named sounds without shipping audio. | S | `app.Sound.Beep()`, `app.Sound.Play(name)` |
| 50 | **Speech** | `AVSpeechSynthesizer`, `SFSpeechRecognizer` | Read text aloud and dictate. The speech menu roles exist but only for WebView selections. | M | `app.Speech.Speak(text, voice)`, `Recognize(...)` |

#### System integration

| # | Feature | Native API | Why it matters | Effort | Go surface sketch |
|---|---|---|---|---|---|
| 51 | **Permission requests** | `AVCaptureDevice requestAccessForMediaType:`, `CGRequestScreenCaptureAccess`, `AXIsProcessTrustedWithOptions`, `CLLocationManager`, `UNUserNotificationCenter` | The cross-platform `Permissions` window option is ignored on macOS (#6067) and dev-mode capture prompts fail (#4270). A single manager that requests and reports TCC state closes both. | M | `app.Permissions.Request(kind) (Status, error)`, `Status(kind)`, plus honouring `Permissions` in `WKUIDelegate requestMediaCapturePermission` |
| 52 | **Sleep and idle control** | `NSProcessInfo beginActivityWithOptions:`, `IOPMAssertionCreateWithName` | Keep the display or system awake during a long export or playback. | S | `app.Power.PreventSleep(reason) (release func())` |
| 53 | **Power and thermal state** | `isLowPowerModeEnabled` (read exists internally), `thermalState`, IOKit power source | Adapt work to battery and thermal pressure. | S | `app.Env.PowerState()` with events |
| 54 | **Sudden and automatic termination** | `NSProcessInfo disableSuddenTermination`, `disableAutomaticTermination:` | Lets macOS kill an idle app instantly at logout unless work is pending. | S | `app.Lifecycle.HoldTermination(reason) (release func())` |
| 55 | **Services provider** | `NSApplication.servicesProvider`, `NSServices` Info.plist entries | Receive selected text or files from other apps via the Services menu. Wails shows the Services menu but cannot provide services. | M | `app.Services.Register(name, types, handler)` plus packaging support |
| 56 | **AppleScript and Apple Events** | `NSAppleEventManager`, `sdef` | Scripting support for automation-heavy users. `kAEGetURL` is already handled for URL schemes. | L | `app.AppleEvents.Handle(class, id, func)`, optional generated `sdef` |
| 57 | **Handoff and universal links** | `NSUserActivity`, `userActivityWillSave`, `continueUserActivity:` | Continue a document on another device; open `applinks` URLs. Single-instance already routes custom schemes. | M | `app.Activity.Publish(type, userInfo)`, `OnContinue` |
| 58 | **Quick Look** | `QLPreviewPanel`, `QLPreviewView`, `QLThumbnailGenerator` | Preview any file the way Finder does, and generate thumbnails for lists. | M | `app.QuickLook.Preview(paths)`, `Thumbnail(path, size) ([]byte, error)` |
| 59 | **Spotlight indexing** | `CoreSpotlight CSSearchableIndex` | Make in-app content searchable from the system. | M | `app.Spotlight.Index(items)`, `OnOpen(id)` |
| 60 | **Workspace helpers** | `NSWorkspace openURLs:withApplicationAtURL:`, `setDefaultApplicationAtURL:toOpenContentType:` (macOS 12), `activateWithOptions:`, `LSSetDefaultHandlerForURLScheme` | Open with a specific app, register as default for a type or scheme at runtime. | S | `app.Browser.OpenWith(path, app)`, `app.Env.SetDefaultHandler(uti or scheme)` |
| 61 | **Accessibility state** | `accessibilityDisplayShouldIncreaseContrast`, `ShouldReduceTransparency`, `isVoiceOverEnabled`, `NSWorkspaceAccessibilityDisplayOptionsDidChange` | Reduce-motion is read internally for the zoom animation but nothing reaches Go. Apps adapt UI to these settings. | S | `app.Env.Accessibility()` struct plus a change event |
| 62 | **Keyboard layout and locale changes** | `kTISNotifySelectedKeyboardInputSourceChanged`, `NSCurrentLocaleDidChangeNotification`, `CFBundleLocalizations` | Global shortcuts are layout-aware on the follow-up branch; apps also want locale change events. Issue #4582 asks how to localise a Mac app at all. | S | `events.Mac.ApplicationDidChangeKeyboardLayout`, `DidChangeLocale`; packaging writes `CFBundleLocalizations` |

### 4.3 Tier 3: native content and deeper WebKit

These change what a window can contain or replace private WebKit hooks with public API. They should follow the window and content model decisions in the two macOS WEP drafts rather than be added ad hoc.

#### Native content views

| # | Feature | Native API | Why it matters | Effort |
|---|---|---|---|---|
| 63 | **Rich text editor** | `NSTextView` attributed strings, `NSFontPanel`, find bar (`usesFindBar`), rulers, undo manager, spell and grammar checking, Writing Tools | `MacTextEditor` is plain text. Rich text makes the WebView-free window useful for notes, mail and documents. | M |
| 64 | **PDF view** | `PDFKit PDFView`, `PDFThumbnailView` | Native PDF reading with search, selection and annotation. | M |
| 65 | **Quick Look view** | `QLPreviewView` as a split pane | Preview any file type in a pane. Shares the #58 bridge. | S |
| 66 | **Media player** | `AVKit AVPlayerView` | Native video with Picture in Picture and AirPlay. | M |
| 67 | **Map view** | `MapKit MKMapView` | Native maps without a web tile provider. | M |
| 68 | **Table and collection views** | `NSTableView`, `NSCollectionView` | Beyond the content list in #8, general data grids and galleries as panes. | L |
| 69 | **Image view** | `NSImageView` with `imageScaling`, `animates` | Cheap native image display for previews. | S |
| 70 | **SwiftUI and custom view hosting** | `NSHostingView` via a Swift dynamic library, or a generic `NSView` escape hatch | Lets teams add views Wails does not model. Needs a toolchain decision. | L |
| 71 | **Metal or CALayer surface** | `MTKView`, `CAMetalLayer` | GPU-rendered panes for visualisation tools. | L |

#### Public WebKit APIs not yet used

| # | Feature | Native API | Why it matters | Effort |
|---|---|---|---|---|
| 72 | **Public inspector toggle** | `WKWebView.isInspectable` (macOS 13.3) | Replaces the `_inspector` private-API path that today needs `private_mac_apis`. Related to #6121. | S |
| 73 | **Under-page background colour** | `underPageBackgroundColor` (macOS 12) | Public alternative for part of what `drawsBackground` does; the split view transparency is a no-op without the private tag. | S |
| 74 | **Downloads** | `WKDownload`, `WKDownloadDelegate` (macOS 11.3) | Native download handling with progress and destination prompts. Zero hits today. | M |
| 75 | **`window.open` and popups** | `WKUIDelegate createWebViewWithConfiguration:` | Open issue #5043. Decide between a new `WebviewWindow` and delegating to the browser. | M |
| 76 | **Navigation interception and cookies** | `decidePolicyForNavigationAction:` exposed to Go, `WKHTTPCookieStore` | Issue #5799. One hit exists for policy but nothing reaches Go. | M |
| 77 | **Website data stores and profiles** | `WKWebsiteDataStore nonPersistentDataStore`, `dataStoreForIdentifier:` (macOS 14) | Private windows, per-account profiles, clearing site data. | M |
| 78 | **Session state** | `interactionState` (macOS 12) | Save and restore scroll, form and history state across launches. Pairs with #20. | S |
| 79 | **Isolated script worlds and reply handlers** | `WKContentWorld`, `WKScriptMessageHandlerWithReply` | Injected runtime scripts that cannot collide with page scripts, and request-response messaging without the HTTP transport. | M |
| 80 | **Capture state** | `cameraCaptureState`, `microphoneCaptureState`, `setCameraCaptureState:` (macOS 12) | Show and control live capture indicators; part of #51. | S |
| 81 | **Page zoom and find** | `pageZoom` (macOS 11), `findInteraction` | Public zoom that persists layout, and a native find bar. `ZoomIn` roles use CSS zoom today. | S |
| 82 | **Writing Tools** | `writingToolsBehavior` (macOS 15.1) | Apple Intelligence proofreading and rewriting inside the WebView. Off by default in apps. | S |
| 83 | **Web extensions** | `WKWebExtension`, `WKWebExtensionController` (macOS 15.4) | Load browser-style extensions into the app WebView. | L |
| 84 | **App-bound domains and HTTPS upgrade** | `limitsNavigationsToAppBoundDomains`, `upgradeKnownHostsToHTTPS` | Hardening options for apps that load remote content. | S |
| 85 | **Theme colour and link previews** | `themeColor`, `allowsLinkPreview` (two hits, not exposed) | Read the page's theme colour for chrome tinting; toggle Force-click previews. | S |

## 5. Recommended sequencing

1. **Stabilise first.** Fix the `server` cross-build, make split and toolbar construction work for windows created after `App.Run`, and make `NativeWindowOptions.Mac` truthful. `review.md` covers all three.
2. **Restore accessories (#1)**, since the branch history shows the API existed and the wrapper is stranded without it.
3. **Cheap Tier 2 wins in one batch**: tray tooltip (#41), tray symbol (#42), menu symbols (#27), dock menu (#28), Open Recent (#29), represented file and edited dot (#15, #16), subtitle (#17), cascading (#18), frame autosave (#19), critical attention (#24), alert suppression and prompt (#33, #34), content-type filters (#36), print options (#39), PDF and snapshot export (#40), haptics and sounds (#48, #49), sleep control (#52), accessibility state (#61). Each is a small commit and most close a known gap or issue.
4. **Toolbar completion (#2 to #7)** and **inspector controls (#12)**, which reuse existing plumbing.
5. **Permissions manager (#51)** with the WebKit capture-state APIs (#80), since it fixes two open bugs.
6. **Rich clipboard and drag-out (#45 to #47)**, which are cross-platform features with the macOS side as the first implementation.
7. **Content list (#8)**, **window tab groups (#14)**, **sheets and popovers (#21, #22)** once the runtime construction model is settled.
8. **Public WebKit replacements (#72, #73)** whenever the private-API tag becomes a maintenance burden.
9. Tier 3 content views after the v4 window model decision that the native-window WEP defers to.

## 6. Out of scope for this audit

- Bugs that need fixing rather than features: Fn/Globe key swallowed by WKWebView (#6033), ProMotion capped at 60 FPS (#6056), title bar resize animation desync (#6070), mixed-DPI screen bounds (#5409).
- Swift-only frameworks with no Objective-C surface: App Intents and Shortcuts, WidgetKit, Live Activities. These need a Swift bridge decision first, which #70 would provide.
- Store and account services: StoreKit, Sign in with Apple, passkeys via `ASAuthorizationController`. Worth a separate proposal because they are entitlement and packaging heavy.
- Packaging and tooling gaps such as pkg installers (#2413) and entitlement generation. They are CLI work, not runtime API.
