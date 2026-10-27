---
title: "WScreenSavers"
description: "Набор экранных заставок Wayland, живых обоев и PAM-блокировки экрана для Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![Окно студии WScreenSavers с предпросмотром 3D WebGL и модулей XScreenSaver](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) — набор экранных заставок Wayland, живых обоев и блокировки экрана через PAM для рабочих столов Linux (`wlroots`, Hyprland, Sway, KDE, GNOME), созданный на базе Go, Wails v3, Vue 3, Three.js и `gtk-layer-shell`.

Приложение выводит классические модули XScreenSaver (`-window-id` через встроенный XWayland/Xvfb) и трёхмерные сцены Three.js WebGL на поверхности Wayland layer-shell (`BACKGROUND` и `OVERLAY`), поддерживает отслеживание простоя, несколько мониторов, управление через IPC CLI (`wscreensavers ctl`) и плагин для Noctalia Shell.

[Посмотреть на GitHub](https://github.com/chrisfair/wscreensavers) | [Релизы](https://github.com/chrisfair/wscreensavers/releases/latest)
