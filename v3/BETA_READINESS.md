# Beta Readiness

This document is for tracking the status of the v3-alpha branch in readiness for the beta release.

## Examples

| Example            | Linux | Windows                    | macOS                                                 |
|--------------------|-------|----------------------------|-------------------------------------------------------|
| badge              | no folder for it?      | ✅                          |                                                       |
| binding            |✅     | ✅                          | ✅                                                     |
| cancel-async       |✅     | ✅                          | ✅                                                     |
| cancel-chaining    |✅     | ✅                          | ✅                                                     |
| clipboard          |✅     | ✅                          | ✅                                                     |
| context-menus      |🚫 GTK error| ✅                          | 🚫 panic                                              |
| dialogs            |⚠️ services menu empty<br>open and save -> no dialogs open| ⚠️ custom icon not working <br>⚠️question (with cancel) (no cancel btn, esc does nothing)<br> ⚠️save (full example) dialogue not opening|⚠️                              |
| dialogs-basic      |⚠️ multiple filters, multiple selection and full featured dialog don't show files in dialog (filter issue?)| ⚠️ basic open is behind app window<br>⚠️ save dialogue is behind window | ✅                            |
| drag-n-drop        |✅     | ✅                          | ✅                                                     |
| environment        |✅     | ✅                          | ✅                                                     |
| events             |✅     | ✅ but window 1 hidden behind window2| ✅                                            |
| file-association   |⚠️ message pops up but no filename| ✅                  | ✅                                  |
| frameless          |✅     | ✅                          | ⚠️ minimise for 3 not working                         |
| gin-example        |✅     | ✅                          | ✅                                                     |
| gin-routing        |✅     | ✅                          | ⚠️ cant see difference from gin-example (copy/paste?) |
| gin-service        |✅     | ✅                          | ⚠️ half buttons does nothing? ( getuserbyid )         |
| html-dnd-api       |✅     | 🚫                         | ✅                                                     |
| ignore-mouse       |⚠️can interact normally| ✅                          | ✅                                    |
| keybindings        |✅     | ✅                          | ✅                                                     |
| menu               |⚠️ click me and submenus do nothing| ⚠️ Hide/Unhide issue       | ✅                         |
| notifications      |✅     | ✅                          | ⚠️ nothing happens on button click                    |
| panic-handling     |✅     | ✅                          | ✅                                                     |
| plain              |✅     | ✅                          | ✅                                                     |
| raw-message        |✅     | ✅                          | ✅                                                     |
| screen             |✅     | ✅                          | ⚠️ slider bubble drags window                          |
| services           |✅     | ✅ ?windows threat protection blocks it| ✅                                          |
| show-macos-toolbar |➖     | ➖                          | ✅                                                     |
| single-instance    |✅     | ✅                          | ✅                                                     |
| systray-basic      |✅     | ⚠️white window in centre of screen on launch, no right click menu|✅                                      |
| systray-custom     |⚠️nothing on click/doubleclick.<br>window opens on right click (centre of screen), quits on close.<br>context menu pops in for a sec only.| ✅ | ✅|
| systray-menu       |⚠️nothing on click/doubleclick.<br>window and menu toggle on right click.<br>systray does not hide. | ⚠️white window in centre of screen on launch| ✅ |
| video              |✅     | ✅  | ✅ |
| window             |✅| ⚠️hide minimise and hide maximise don't work<br> ⚠️hide close hides all three buttons| ✅                         |
| window-api         |✅     | ✅                          | ✅                                                     |
| window-call        |✅     | ✅                          | ✅                                                     |
| window-menubar     |⚠️menubar doesn't show| ✅  | ⚠️ not sure what should happen in osx                 |
| wml                |⚠️restore does nothing after window moved/resized| ✅| ✅|

## Open Bugs
- 

- https://github.com/wailsapp/wails/issues/3743
- https://github.com/wailsapp/wails/issues/3683 - needs checking
- https://github.com/wailsapp/wails/issues/4235
- https://github.com/wailsapp/wails/issues/4236

## Todo

- [ ] [Custom Protocol Support](https://github.com/wailsapp/wails/issues/4026)
- [ ] [Implement Window.SetScreen](https://github.com/wailsapp/wails/issues/4000)
- [ ] [Port DLL Directory Initialisation](https://github.com/wailsapp/wails/pull/4207)
- [ ] Check if [this](https://github.com/wailsapp/wails/pull/4047#issuecomment-2814676117) needs porting.
- [ ] Update docs
    - [ ] Add tutorials
