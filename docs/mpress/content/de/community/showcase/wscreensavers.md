---
title: "WScreenSavers"
description: "Wayland-Bildschirmschoner, Live-Hintergrund und PAM-Bildschirmsperre für Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![Studio-Fenster von WScreenSavers mit 3D-WebGL- und XScreenSaver-Vorschauen](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) ist eine Suite für Wayland-Bildschirmschoner, Live-Hintergründe und PAM-Bildschirmsperren unter Linux (`wlroots`, Hyprland, Sway, KDE, GNOME), entwickelt mit Go, Wails v3, Vue 3, Three.js und `gtk-layer-shell`.

Sie bindet klassische XScreenSaver-Module (`-window-id` über eingebettetes XWayland/Xvfb) sowie eigene Three.js-WebGL-Szenen auf Wayland-Layer-Shell-Oberflächen (`BACKGROUND` und `OVERLAY`) ein und unterstützt Leerlauferkennung, mehrere Monitore, eine IPC-CLI (`wscreensavers ctl`) sowie ein Noctalia-Shell-Plugin.

[Auf GitHub ansehen](https://github.com/chrisfair/wscreensavers) | [Releases](https://github.com/chrisfair/wscreensavers/releases/latest)
