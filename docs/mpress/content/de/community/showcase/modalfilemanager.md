---
title: "Modal File Manager"
description: "Eine mit Wails entwickelte Desktopanwendung"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager) ist ein mit Webtechnologien entwickelter Dateimanager mit zwei Bereichen. Mein ursprünglicher Entwurf basierte auf NW.js und ist [hier](https://github.com/raguay/ModalFileManager-NWjs) zu finden. Diese Version verwendet denselben Svelte-basierten Frontendcode, der jedoch seit der Abkehr von NW.js umfassend geändert wurde. Das Backend ist dagegen eine [Wails-2](https://wails.io/)Implementierung. Mit dieser Implementierung verwende ich keine Befehlszeilenbefehle wie `rm`, `cp` usw. mehr. Zum Herunterladen von Themes und Erweiterungen muss Git jedoch auf dem System installiert sein. Die Anwendung ist vollständig in Go programmiert und läuft wesentlich schneller als die vorherigen Versionen.

Dieser Dateimanager folgt demselben Prinzip wie Vim: Tastaturaktionen werden durch Zustände gesteuert. Die Anzahl der Zustände ist nicht festgelegt, sondern umfassend programmierbar. Daher lassen sich unendlich viele Tastaturkonfigurationen erstellen und verwenden. Dies ist der Hauptunterschied zu anderen Dateimanagern. Themes und Erweiterungen können von GitHub heruntergeladen werden.
