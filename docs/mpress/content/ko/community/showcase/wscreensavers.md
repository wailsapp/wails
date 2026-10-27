---
title: "WScreenSavers"
description: "Linux용 Wayland 화면 보호기, 라이브 배경화면 및 PAM 화면 잠금 도구 모음"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![3D WebGL 및 XScreenSaver 미리보기를 보여주는 WScreenSavers 스튜디오 창](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers)는 Go, Wails v3, Vue 3, Three.js 및 `gtk-layer-shell`로 구축된 Linux 데스크톱(`wlroots`, Hyprland, Sway, KDE, GNOME)용 Wayland 화면 보호기, 라이브 배경화면 및 PAM 화면 잠금 도구 모음입니다.

내장 XWayland/Xvfb를 통한 기존 XScreenSaver 모듈(`-window-id`)과 사용자 지정 Three.js WebGL 장면을 Wayland `BACKGROUND` 및 `OVERLAY` 레이어 셸 표면에 렌더링하며, 유휴 상태 감지, 다중 모니터 지원, IPC CLI(`wscreensavers ctl`), Noctalia Shell 플러그인을 포함합니다.

[GitHub에서 보기](https://github.com/chrisfair/wscreensavers) | [릴리스](https://github.com/chrisfair/wscreensavers/releases/latest)
