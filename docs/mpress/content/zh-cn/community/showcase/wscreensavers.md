---
title: "WScreenSavers"
description: "适用于 Linux 的 Wayland 屏幕保护程序、动态壁纸与 PAM 锁屏套件"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![显示 3D WebGL 与 XScreenSaver 预览的 WScreenSavers 工作室窗口](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) 是面向 Linux 桌面环境（`wlroots`、Hyprland、Sway、KDE、GNOME）的 Wayland 屏幕保护程序、动态壁纸与 PAM 锁屏套件，基于 Go、Wails v3、Vue 3、Three.js 和 `gtk-layer-shell` 构建。

它通过内嵌 XWayland/Xvfb 将经典 XScreenSaver 模块（`-window-id`）与自定义 Three.js WebGL 场景渲染到 Wayland 的 `BACKGROUND` 和 `OVERLAY` 层级表面，支持空闲检测、多显示器、IPC 命令行工具（`wscreensavers ctl`）以及 Noctalia Shell 插件。

[在 GitHub 上查看](https://github.com/chrisfair/wscreensavers) | [发布版本](https://github.com/chrisfair/wscreensavers/releases/latest)
