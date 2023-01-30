# TODO

Informal and incomplete list of things needed in v3.

## General

- [ ] Generate Bindings
- [ ] Generate TS Models
- [ ] Port NSIS creation 
- [ ] Dev Mode
- [ ] Generate Info.Plist from `info.json`

- [ ] Windows Port
- [ ] Linux Port

## Runtime

- [x] Pass window ID with window calls in JS
- [ ] Implement alias for `window` in JS
- [ ] Implement runtime dispatcher
  - [ ] Log
  - [x] Same Window
  - [ ] Other Window
  - [ ] Dialogs
    - [x] Info
    - [x] Warning
    - [x] Error
    - [x] Question
    - [x] OpenFile
    - [x] SaveFile
  - [ ] Events
  - [ ] Screens
  - [x] Clipboard
  - [ ] Application
- [ ] Create `.d.ts` file

## Templates

- [ ] Create plain template
- [ ] Improve default template

## Runtime

- [ ] To log or not to log?
- [ ] Unify cross-platform events, eg. `onClose`