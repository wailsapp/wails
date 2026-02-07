# Dialog Manual Tests

Comprehensive test suite for the GTK4 dialog implementation.

## Building

```bash
cd v3/test/manual/dialog
task build:all
```

Binaries are output to `bin/` directory with GTK3/GTK4 variants.

## Test Categories

### Message Dialogs

#### 1. message-info

Tests info/information dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Info | Dialog with title and message |
| Title Only | Dialog with only title |
| Message Only | Dialog with only message |
| Custom Icon | Dialog displays custom Wails icon |
| Long Message | Text wraps properly |
| Attached to Window | Dialog is modal to main window |

#### 2. message-question

Tests question dialogs with buttons.

| Test | Expected Behavior |
|------|-------------------|
| Two Buttons | Yes/No buttons, callbacks work |
| Three Buttons | Save/Don't Save/Cancel buttons |
| With Default Button | Default button highlighted, Enter selects it |
| With Cancel Button | Escape key triggers cancel button |
| Custom Icon | Dialog displays custom icon |
| Attached to Window | Dialog is modal to main window |
| Button Callbacks | Each button triggers correct callback |

#### 3. message-warning

Tests warning dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Warning | Warning dialog with title and message |
| Title Only | Warning with only title |
| Message Only | Warning with only message |
| Custom Icon | Warning with custom icon |
| Long Warning | Text wraps properly |
| Attached to Window | Dialog is modal to main window |

#### 4. message-error

Tests error dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Error | Error dialog with title and message |
| Title Only | Error with only title |
| Message Only | Error with only message |
| Custom Icon | Error with custom icon |
| Technical Error | Long error message wraps properly |
| Attached to Window | Dialog is modal to main window |

### File Dialogs

#### 5. file-open

Tests single file open dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Open | File picker opens, selection returned |
| With Title | Dialog has custom title |
| Show Hidden Files | Hidden files (.*) visible |
| Start in Home | Dialog opens in home directory |
| Start in /tmp | Dialog opens in /tmp |
| Filter: Text Files | Only .txt, .md, .log files shown |
| Filter: Images | Only image files shown |
| Multiple Filters | Filter dropdown with multiple options |
| Custom Button Text | Open button has custom text |
| Attached to Window | Dialog is modal to main window |

#### 6. file-open-multi

Tests multiple file selection.

| Test | Expected Behavior |
|------|-------------------|
| Select Multiple Files | Can select multiple files with Ctrl+click |
| With Hidden Files | Hidden files visible in selection |
| Filter: Source Code | Only source files shown |
| Filter: Documents | Only document files shown |
| Attached to Window | Dialog is modal to main window |

#### 7. file-save

Tests save file dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Save | Save dialog opens |
| With Message | Dialog has custom message |
| With Default Filename | Filename field pre-populated |
| Start in Home | Dialog opens in home directory |
| Start in /tmp | Dialog opens in /tmp |
| Show Hidden Files | Hidden files visible |
| Can Create Directories | New folder button works |
| Cannot Create Directories | New folder button hidden/disabled |
| Custom Button Text | Save button has custom text |
| Attached to Window | Dialog is modal to main window |

#### 8. file-directory

Tests directory selection dialogs.

| Test | Expected Behavior |
|------|-------------------|
| Basic Directory | Can only select directories |
| Start in Home | Dialog opens in home directory |
| Start in / | Dialog opens at root |
| Can Create Directories | New folder button works |
| Show Hidden | Hidden directories visible |
| Resolve Aliases/Symlinks | Symlinks resolved to real paths |
| Custom Button Text | Open button has custom text |
| Multiple Directories | Can select multiple directories |
| Attached to Window | Dialog is modal to main window |

## GTK Version Matrix

| Test | GTK4 | GTK3 |
|------|------|------|
| message-info | | |
| message-question | | |
| message-warning | | |
| message-error | | |
| file-open | | |
| file-open-multi | | |
| file-save | | |
| file-directory | | |

## Running Individual Tests

```bash
# GTK4 (default)
./bin/message-info-gtk4
./bin/message-question-gtk4
./bin/message-warning-gtk4
./bin/message-error-gtk4
./bin/file-open-gtk4
./bin/file-open-multi-gtk4
./bin/file-save-gtk4
./bin/file-directory-gtk4

# GTK3
./bin/message-info-gtk3
./bin/message-question-gtk3
./bin/message-warning-gtk3
./bin/message-error-gtk3
./bin/file-open-gtk3
./bin/file-open-multi-gtk3
./bin/file-save-gtk3
./bin/file-directory-gtk3
```

## Checklist for Full Verification

### Message Dialogs

- [ ] Dialog appears centered or attached correctly
- [ ] Title displays correctly
- [ ] Message displays correctly
- [ ] Custom icons display correctly
- [ ] Long text wraps properly
- [ ] OK/Close button dismisses dialog
- [ ] Escape key closes dialog (where applicable)

### Question Dialogs

- [ ] All buttons display correctly
- [ ] Button callbacks fire correctly
- [ ] Default button is highlighted
- [ ] Enter key activates default button
- [ ] Escape key activates cancel button
- [ ] Multiple buttons layout correctly

### File Dialogs

- [ ] Dialog opens in correct directory
- [ ] Filters work correctly
- [ ] Hidden files toggle works
- [ ] Create directory works (where enabled)
- [ ] Cancel returns empty string
- [ ] Selection returns correct path(s)
- [ ] Multiple selection works (multi tests)
- [ ] Custom button text displays

### Known Issues

Document any issues found during testing:

```
[GTK Version] [Test] - Issue description
Example: GTK4 file-open - Filter dropdown not visible
```
