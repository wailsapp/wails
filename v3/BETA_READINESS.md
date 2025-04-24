# Beta Readiness

This document is for tracking the status of the v3-alpha branch in readiness for the beta release.

## Examples

| Example            | Linux | Windows                                          | macOS |
|--------------------|-------|--------------------------------------------------|-------|
| badge              |       | ✅                                                |       |
| binding            |       | ✅                                                |       |
| cancel-async       |       | ✅                                                |       |
| cancel-chaining    |       | ✅                                                |       |
| clipboard          |       | ✅                                                |       |
| context-menus      |       | ✅                                                |       |
| dialogs            |       | ⚠️ custom icon not working                       |       |
| dialogs-basic      |       | ⚠️ cancel crashes the app                        |       |
| drag-n-drop        |       | ✅                                                |       |
| environment        |       | ✅                                                |       |
| events             |       | ✅                                                |       |
| file-association   |       | ✅                                                |       |
| frameless          |       | ✅                                                |       |
| gin-example        |       | ✅                                                |       |
| gin-routing        |       | ✅                                                |       |
| gin-service        |       | ✅                                                |       |
| html-dnd-api       |       | 🚫                                               |       |
| ignore-mouse       |       | ✅                                                |       |
| keybindings        |       | ✅                                                |       |
| menu               |       | ⚠️ Hide/Unhide issue                             |       |
| notifications      |       | ✅                                                |       |
| panic-handling     |       | ✅                                                |       |
| plain              |       | ✅                                                |       |
| raw-message        |       | ✅                                                |       |
| screen             |       | ✅                                                |       |
| services           |       | ⚠️ KV needs refreshing after save                |       |
| show-macos-toolbar |       | ➖                                                |       |
| single-instance    |       | ✅                                                |       |
| systray-basic      |       | ✅                                                |       |
| systray-custom     |       | ✅                                                |       |
| systray-menu       |       | ⚠️ Check systray.Hide/Show                       |       |
| video              |       | ✅                                                |       |
| window             |       | ⚠️ SetPos 0,0 is going to 5,0. GetPos is correct |       |
| window-api         |       | ✅                                                |       |
| window-call        |       | ✅                                                |       |
| window-menubar     |       | ✅                                                |       |
| wml                |       | ✅                                                |       |

## Open Bugs

- https://github.com/wailsapp/wails/issues/4151
- https://github.com/wailsapp/wails/issues/4131
- https://github.com/wailsapp/wails/issues/3743
- https://github.com/wailsapp/wails/issues/3683 - needs checking
- https://github.com/wailsapp/wails/issues/1503
- https://github.com/wailsapp/wails/issues/4235
- https://github.com/wailsapp/wails/issues/4236

## Todo

- [ ] [Custom Protocol Support](https://github.com/wailsapp/wails/issues/4026)
- [ ] [Implement Window.SetScreen](https://github.com/wailsapp/wails/issues/4000)
- [ ] [Port DLL Directory Initialisation](https://github.com/wailsapp/wails/pull/4207)
- [ ] Check if [this](https://github.com/wailsapp/wails/pull/4047#issuecomment-2814676117) needs porting.
- [ ] Update docs
  - [ ] Add tutorials