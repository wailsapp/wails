# Roadmap

The roadmap is a living document and is subject to change. If you have any
suggestions, please open an issue. Each milestone will have a set of goals that
we are aiming to achieve. These are subject to change.

## Alpha milestones

### Alpha 1

#### Goals

Alpha 1 is the initial release. It is intended to get feedback on the new API
and to get people experimenting with it. The main goal is to get most of the
examples working on all platforms.

#### Status

- :material-check-bold: - Working
- :material-minus: - Partially working
- :material-close: - Not working

{{ read_csv("status.csv") }}

- Mac Dialogs work, however the file dialogs issue a warning that needs to be
  fixed.

#### TODO:

- [ ] Fix `+[CATransaction synchronize] called within transaction` warnings on
      Mac
- [ ] When hiding window, application terminates

### Alpha 2

- [ ] Most examples working on Linux
- [ ] Project creation via `wails init`
