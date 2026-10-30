---
title: "WScreenSavers"
description: "Linux 向け Wayland スクリーンセーバー、ライブ壁紙、PAM 画面ロックスイート"
slug: "community/showcase/wscreensavers"
sourcePath: "community/showcase/wscreensavers.md"
---

![3D WebGL と XScreenSaver のプレビューを表示する WScreenSavers スタジオウィンドウ](/assets/showcase-images/wscreensavers.png)

[WScreenSavers](https://github.com/chrisfair/wscreensavers) は、Go、Wails v3、Vue 3、Three.js、`gtk-layer-shell` で構築された Linux デスクトップ（`wlroots`、Hyprland、Sway、KDE、GNOME）向けの Wayland スクリーンセーバー、ライブ壁紙、PAM 画面ロックスイートです。

埋め込み XWayland/Xvfb 経由の従来型 XScreenSaver モジュール（`-window-id`）とカスタム Three.js WebGL シーンを Wayland の `BACKGROUND` および `OVERLAY` レイヤーシェルサーフェスに描画し、アイドル検出、マルチモニター対応、IPC CLI（`wscreensavers ctl`）、Noctalia Shell プラグインを備えています。

[GitHubで見る](https://github.com/chrisfair/wscreensavers) | [リリース](https://github.com/chrisfair/wscreensavers/releases/latest)
