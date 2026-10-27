---
title: "WScreenSavers"
description: "Suíte de protetor de tela Wayland, papel de parede ao vivo e bloqueio PAM para Linux"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![Janela do estúdio WScreenSavers exibindo prévias WebGL 3D e XScreenSaver](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) é uma suíte de protetor de tela Wayland, papel de parede ao vivo e bloqueio de tela PAM para desktops Linux (`wlroots`, Hyprland, Sway, KDE, GNOME), desenvolvida com Go, Wails v3, Vue 3, Three.js e `gtk-layer-shell`.

O aplicativo integra módulos clássicos do XScreenSaver (`-window-id` via XWayland/Xvfb embutido) e cenas WebGL em Three.js em superfícies layer-shell do Wayland (`BACKGROUND` e `OVERLAY`), com detecção de ociosidade, suporte a múltiplos monitores, CLI IPC (`wscreensavers ctl`) e plugin para Noctalia Shell.

[Ver no GitHub](https://github.com/chrisfair/wscreensavers) | [Versões](https://github.com/chrisfair/wscreensavers/releases/latest)
