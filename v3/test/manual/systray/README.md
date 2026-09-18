# Systray Manual Tests

Comprehensive test suite for the Systray API v2 implementation.

## Building

```bash
cd v3/test/manual/systray
task build:all
```

Binaries are output to `bin/` directory with GTK3/GTK4 variants.

## Test Scenarios

### 1. window-only

Tests systray with only an attached window (no menu).

| Action | Expected Behavior |
|--------|-------------------|
| Left-click | Toggle window visibility |
| Right-click | Nothing |
| Double-click | Nothing |

### 2. menu-only

Tests systray with only a menu (no attached window).

| Action | Expected Behavior |
|--------|-------------------|
| Left-click | Nothing |
| Right-click | Show menu |
| Double-click | Nothing |

### 3. window-menu

Tests systray with both attached window and menu.

| Action | Expected Behavior |
|--------|-------------------|
| Left-click | Toggle window visibility |
| Right-click | Show menu |
| Double-click | Nothing |

### 4. custom-handlers

Tests custom click handlers overriding smart defaults. Demonstrates that custom handlers
can add logging/analytics while still using the standard systray methods.

| Action | Expected Behavior |
|--------|-------------------|
| Left-click | Custom handler logs + toggles window |
| Right-click | Custom handler logs + opens menu |
| Double-click | Logs to console |

### 5. hide-options

Tests `HideOnEscape` and `HideOnFocusLost` window options.

| Action | Expected Behavior |
|--------|-------------------|
| Left-click systray | Toggle window |
| Press Escape | Hide window |
| Click outside | Hide window (standard WMs only) |

## Platform Test Matrix

### Desktop Environments

| Test | GNOME | KDE | XFCE | Hyprland | Sway | i3 |
|------|-------|-----|------|----------|------|-----|
| window-only | | | | | | |
| menu-only | | | | | | |
| window-menu | | | | | | |
| custom-handlers | | | | | | |
| hide-options | | | | | | |

### GTK Version Matrix

| Test | GTK4 | GTK3 |
|------|------|------|
| window-only | | |
| menu-only | | |
| window-menu | | |
| custom-handlers | | |
| hide-options | | |

## Focus-Follows-Mouse Behavior

On tiling WMs with focus-follows-mouse (Hyprland, Sway, i3):

- `HideOnFocusLost` is **automatically disabled**
- This prevents the window from hiding immediately when the mouse moves away
- `HideOnEscape` still works normally

### Testing HideOnFocusLost

1. **Standard WM (GNOME, KDE)**:
   - Run `hide-options` test
   - Click systray to show window
   - Click anywhere outside the window
   - Window should hide

2. **Tiling WM (Hyprland, Sway, i3)**:
   - Run `hide-options` test
   - Click systray to show window
   - Move mouse outside window
   - Window should NOT hide (feature disabled)
   - Press Escape to hide (still works)

## Checklist for Full Verification

### Per-Environment Checklist

- [ ] Systray icon appears in system tray
- [ ] Left-click behavior matches expected
- [ ] Right-click behavior matches expected
- [ ] Window appears near systray icon (not random position)
- [ ] Window stays on top when shown
- [ ] Multiple show/hide cycles work correctly
- [ ] Menu items are clickable and work
- [ ] Application quits cleanly

### Hide Options Specific

- [ ] Escape key hides window
- [ ] Focus lost hides window (standard WMs)
- [ ] Focus lost does NOT hide window (tiling WMs)
- [ ] Re-clicking systray after hide works

### Known Issues

Document any issues found during testing:

```
[Environment] [GTK Version] [Test] - Issue description
Example: Hyprland GTK4 window-menu - Menu appears at wrong position
```

## Running Individual Tests

```bash
# GTK4 (default)
./bin/window-only-gtk4
./bin/menu-only-gtk4
./bin/window-menu-gtk4
./bin/custom-handlers-gtk4
./bin/hide-options-gtk4

# GTK3
./bin/window-only-gtk3
./bin/menu-only-gtk3
./bin/window-menu-gtk3
./bin/custom-handlers-gtk3
./bin/hide-options-gtk3
```
