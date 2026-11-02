---
title: "WScreenSavers"
description: "適用於 Linux 的 Wayland 螢幕保護程式、動態桌布與 PAM 螢幕鎖定套件"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![顯示 3D WebGL 與 XScreenSaver 預覽的 WScreenSavers 工作室視窗](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) 是面向 Linux 桌面環境（`wlroots`、Hyprland、Sway、KDE、GNOME）的 Wayland 螢幕保護程式、動態桌布與 PAM 螢幕鎖定套件，基於 Go、Wails v3、Vue 3、Three.js 與 `gtk-layer-shell` 建構。

它透過內嵌 XWayland/Xvfb 將經典 XScreenSaver 模組（`-window-id`）與自訂 Three.js WebGL 場景渲染至 Wayland 的 `BACKGROUND` 與 `OVERLAY` 層級表面，支援閒置偵測、多顯示器、IPC 命令列工具（`wscreensavers ctl`）以及 Noctalia Shell 外掛程式。

[在 GitHub 上查看](https://github.com/chrisfair/wscreensavers) | [發行版](https://github.com/chrisfair/wscreensavers/releases/latest)
