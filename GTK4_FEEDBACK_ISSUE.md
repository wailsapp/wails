# GitHub Issue Template: GTK4 Feedback

**Copy the content below to create a new GitHub issue after merging the PR.**

---

## Title
`[Feedback Wanted] Experimental GTK4 + WebKitGTK 6.0 Support`

## Labels
`Enhancement`, `Linux`, `v3`, `feedback-wanted`

## Body

Wails v3 now includes **experimental** support for GTK4 and WebKitGTK 6.0 on Linux via the `-tags gtk4` build flag.

## Why GTK4?

- **Better Wayland support** - Native Wayland rendering without X11 compatibility layer
- **Improved performance** - GPU-accelerated rendering pipeline
- **Modern features** - Fractional scaling, HDR support, better HiDPI handling
- **Future-proof** - GTK3 is in maintenance mode, GTK4 is the future of Linux desktop

## How to Try It

### Prerequisites

Install GTK4 and WebKitGTK 6.0 development packages:

| Distro | Command |
|--------|---------|
| **Ubuntu/Debian** | `sudo apt install libgtk-4-dev libwebkitgtk-6.0-dev` |
| **Fedora** | `sudo dnf install gtk4-devel webkit2gtk6.0-devel` |
| **Arch Linux** | `sudo pacman -S gtk4 webkitgtk-6.0` |
| **openSUSE** | `sudo zypper install gtk4-devel libwebkit2gtk-6_0-devel` |

### Build Your App

```bash
# Development mode
wails3 dev -tags gtk4

# Production build
wails3 build -tags gtk4
```

### Verify GTK Version

Your app's `wails doctor` output should show:
```
GTK Version: 4.x.x
WebKit Version: 6.0.x
```

## What We Need Feedback On

Please test and report on:

- [ ] **Build success** - Does it compile on your distro?
- [ ] **Basic functionality** - Windows, menus, dialogs work?
- [ ] **Performance** - Better/worse than GTK3? (especially on Wayland)
- [ ] **Visual issues** - Any rendering glitches or UI problems?
- [ ] **Missing features** - Anything that works in GTK3 but not GTK4?
- [ ] **Crashes** - Any segfaults or panics?

## How to Provide Feedback

### Comment Format

Please include:
```
**Distro**: Ubuntu 24.04 / Fedora 40 / Arch / etc.
**Desktop**: GNOME / KDE / Sway / etc.
**Display Server**: Wayland / X11
**GTK4 Version**: (from pkg-config --modversion gtk4)
**WebKitGTK Version**: (from pkg-config --modversion webkitgtk-6.0)

**What I tested**: [description]
**Result**: [works / broken / partial]
**Details**: [any relevant info, error messages, screenshots]
```

### For Bugs

Open a separate issue with `[GTK4]` prefix in the title and link it here.

Example: `[GTK4] Menu accelerators not working on KDE Wayland`

### For Fixes

PRs are welcome! Please:
1. Target the `v3-alpha` branch
2. Include `gtk4` in your PR title
3. Test on both GTK3 (default) and GTK4 builds

## Known Limitations

| Limitation | Reason |
|------------|--------|
| Window positioning on Wayland | Wayland protocol doesn't allow arbitrary positioning |
| `ShowHiddenFiles()` in file dialogs | GTK4 uses portal-based dialogs controlled by DE |
| `CanCreateDirectories()` | Same as above - portal controls this |
| X11 may have varied performance | GTK4 is optimized for Wayland |

## Related

- Original feature request: #3193
- Implementation tracker: See `IMPLEMENTATION.md` in repo

---

**Thank you for helping test GTK4 support!** Your feedback directly shapes whether this becomes the default in a future release.
