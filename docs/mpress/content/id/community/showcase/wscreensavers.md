---
title: "WScreenSavers"
description: "Rangkaian screensaver Wayland, wallpaper hidup, dan pengunci layar PAM untuk Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![Jendela studio WScreenSavers menampilkan pratinjau WebGL 3D dan XScreenSaver](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) adalah rangkaian screensaver Wayland, wallpaper hidup, dan pengunci layar PAM untuk desktop Linux (`wlroots`, Hyprland, Sway, KDE, GNOME), dibangun menggunakan Go, Wails v3, Vue 3, Three.js, dan `gtk-layer-shell`.

Aplikasi ini menjembatani modul XScreenSaver klasik (`-window-id` melalui XWayland/Xvfb tertanam) bersama adegan WebGL Three.js ke permukaan layer-shell Wayland (`BACKGROUND` dan `OVERLAY`), dilengkapi deteksi diam, dukungan multi-monitor, CLI IPC (`wscreensavers ctl`), dan plugin Noctalia Shell.

[Lihat di GitHub](https://github.com/chrisfair/wscreensavers) | [Rilis](https://github.com/chrisfair/wscreensavers/releases/latest)
