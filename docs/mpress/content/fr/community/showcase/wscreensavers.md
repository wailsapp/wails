---
title: "WScreenSavers"
description: "Suite d’économiseur d’écran Wayland, fond d’écran animé et verrouillage PAM pour Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![Fenêtre studio de WScreenSavers affichant les aperçus WebGL 3D et XScreenSaver](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) est une suite d’économiseur d’écran Wayland, de fond d’écran animé et de verrouillage d’écran PAM pour les environnements de bureau Linux (`wlroots`, Hyprland, Sway, KDE, GNOME), développée avec Go, Wails v3, Vue 3, Three.js et `gtk-layer-shell`.

Elle affiche les modules classiques XScreenSaver (`-window-id` via XWayland/Xvfb embarqué) ainsi que des scènes WebGL Three.js sur des surfaces Wayland layer-shell (`BACKGROUND` et `OVERLAY`), avec détection d’inactivité, prise en charge multi-écrans, une interface en ligne de commande IPC (`wscreensavers ctl`) et un greffon Noctalia Shell.

[Voir sur GitHub](https://github.com/chrisfair/wscreensavers) | [Versions](https://github.com/chrisfair/wscreensavers/releases/latest)
