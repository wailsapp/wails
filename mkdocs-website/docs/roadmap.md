# Roadmap

The roadmap is a living document and is subject to change. If you have any suggestions, please open an issue.
Each milestone will have a set of goals that we are aiming to achieve. These are subject to change.

## Alpha milestones

### Alpha 1

#### Goals

Alpha 1 is the initial release. It is intended to get feedback on the new API and to get people experimenting with it.
The main goal is to get most of the examples working on all platforms.

#### Status

- W - Working
- P - Partially working
- N - Not working

| Example       | Mac | Windows | Linux |
|---------------|-----|---------|-------|
| binding       | W   | W       |       |
| build         | W   | W       |       |
| clipboard     | W   | W       |       |
| context menus | W   | W       |       |
| dialogs       | P   | W       |       |
| drag-n-drop   | W   | N       |       |
| events        | W   | W       |       |
| frameless     | W   | W       |       |
| keybindings   | W   | W       |       |
| plain         | W   | W       |       |
| screen        | W   | W       |       |
| systray       | W   | W       |       |
| video         |     | W       |       |
| window        | P   | W       |       |
| wml           |     | W       |       |

- Mac Dialogs work, however the file dialogs issue a warning that needs to be fixed.

#### TODO:

- [ ] Fix `+[CATransaction synchronize] called within transaction` warnings on Mac
- [ ] When hiding window, application terminates

### Alpha 2

- [ ] Most examples working on Linux
- [ ] Project creation via `wails init`

