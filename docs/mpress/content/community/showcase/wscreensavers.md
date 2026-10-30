---
title: "WScreenSavers"
description: "Wayland screensaver, live wallpaper, and PAM screen-lock suite for Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![WScreenSavers studio window showing 3D WebGL and XScreenSaver previews](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) is a Wayland screensaver, live wallpaper, and PAM screen-lock suite for Linux desktops (`wlroots`, Hyprland, Sway, KDE, GNOME), built with Go, Wails v3, Vue 3, Three.js, and `gtk-layer-shell`.

It bridges classic XScreenSaver modules (`-window-id` via embedded XWayland/Xvfb) alongside custom Three.js WebGL scenes onto Wayland `BACKGROUND` and `OVERLAY` layer-shell surfaces, with idle detection, multi-monitor support, an IPC CLI (`wscreensavers ctl`), and a Noctalia Shell plugin.

[View on GitHub](https://github.com/chrisfair/wscreensavers) | [Releases](https://github.com/chrisfair/wscreensavers/releases/latest)
