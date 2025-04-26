# Beta Readiness

This document is for tracking the status of the v3-alpha branch in readiness for the beta release.

## Examples

| Example            | Linux | Windows                                          | macOS                                                 |
|--------------------|-------|--------------------------------------------------|-------------------------------------------------------|
| badge              |       | ✅                                                |                                                       |
| binding            |       | ✅                                                | ✅                                                     |
| cancel-async       |       | ✅                                                | ✅                                                     |
| cancel-chaining    |       | ✅                                                | ✅                                                     |
| clipboard          |       | ✅                                                | ✅                                                     |
| context-menus      |       | ✅                                                | 🚫 panic                                              |
| dialogs            |       | ⚠️ custom icon not working                       | ⚠️                                                    |
| dialogs-basic      |       | ✅                                                | ✅                                                     |
| drag-n-drop        |       | ✅                                                | ✅                                                     |
| environment        |       | ✅                                                | ✅                                                     |
| events             |       | ✅                                                | ✅                                                     |
| file-association   |       | ✅                                                | ✅                                                     |
| frameless          |       | ✅                                                | ⚠️ minimise for 3 not working                         |
| gin-example        |       | ✅                                                | ✅                                                     |
| gin-routing        |       | ✅                                                | ⚠️ cant see difference from gin-example (copy/paste?) |
| gin-service        |       | ✅                                                | ⚠️ half buttons does nothing? ( getuserbyid )         |
| html-dnd-api       |       | 🚫                                               | ✅                                                     |
| ignore-mouse       |       | ✅                                                | ✅                                                     |
| keybindings        |       | ✅                                                | ✅                                                     |
| menu               |       | ⚠️ Hide/Unhide issue                             | ✅                                                     |
| notifications      |       | ✅                                                | ⚠️ nothing happens on button click                    |
| panic-handling     |       | ✅                                                | ✅                                                     |
| plain              |       | ✅                                                | ✅                                                     |
| raw-message        |       | ✅                                                | ✅                                                     |
| screen             |       | ✅                                                | ⚠️ slider bubble drags window                         |
| services           |       | ✅                                                | ✅                                                     |
| show-macos-toolbar |       | ➖                                                | ✅                                                     |
| single-instance    |       | ✅                                                | ✅                                                     |
| systray-basic      |       | ✅                                                | ✅                                                     |
| systray-custom     |       | ✅                                                | ✅                                                     |
| systray-menu       |       | ✅                                                | ✅                                                     |
| video              |       | ✅                                                | ✅                                                     |
| window             |       | ⚠️ SetPos 0,0 is going to 5,0. GetPos is correct | ✅                                                     |
| window-api         |       | ✅                                                | ✅                                                     |
| window-call        |       | ✅                                                | ✅                                                     |
| window-menubar     |       | ✅                                                | ⚠️ not sure what should happen in osx                 |
| wml                |       | ✅                                                | ✅                                                     |

## Open Bugs

- https://github.com/wailsapp/wails/issues/4151
- https://github.com/wailsapp/wails/issues/4131
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
