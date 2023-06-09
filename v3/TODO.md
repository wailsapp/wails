# TODO

Informal and incomplete list of things needed in v3.

## General

- [x] Generate Bindings
- [x] Generate TS Models
- [ ] Dev Mode
- [ ] Generate Info.Plist from `info.json`

- [ ] Windows Port
- [ ] Linux Port

## Runtime

- [x] Pass window ID with window calls in JS
- [x] Implement alias for `window` in JS
- [x] Implement runtime dispatcher
  - [x] Log
  - [x] Same Window
  - [ ] Other Window
  - [x] Dialogs
    - [x] Info
    - [x] Warning
    - [x] Error
    - [x] Question
    - [x] OpenFile
    - [x] SaveFile
  - [x] Events
  - [x] Screens
  - [x] Clipboard
  - [x] Application
- [ ] Create `.d.ts` file

## Templates

- [ ] Create plain template
- [ ] Improve default template

## Runtime

- [ ] To log or not to log?
- [ ] Unify cross-platform events, eg. `onClose`

## Plugins

- [ ] Move logins to `v3/plugins`
- [ ] Expose application logger to plugins
