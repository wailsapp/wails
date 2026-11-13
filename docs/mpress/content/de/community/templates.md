---
title: "Vorlagen"
description: "Projektvorlagen und Starter-Kits für Wails"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
Diese Seite ist für Wails v3 möglicherweise veraltet.

@end

<!-- TODO: Update this link -->

Diese Seite enthält eine Liste von Vorlagen, die von der Community unterstützt werden. Informationen zum Erstellen einer eigenen Vorlage finden Sie im Leitfaden [Vorlagen](https://v2.wails.io/docs/guides/templates).

@note{type="tip" title="Eine Vorlage einreichen"}
Klicken Sie unten auf `Edit this page`, um Ihre Vorlagen hinzuzufügen.

@end

Führen Sie zur Verwendung dieser Vorlagen `wails init -n "Your Project Name" -t [the link below[@version]]` aus.

Ist kein Versionssuffix angegeben, wird standardmäßig die Codevorlage des Hauptbranches verwendet. Ist ein Versionssuffix angegeben, wird die Codevorlage verwendet, die dem Tag dieser Version entspricht.

Beispiel: `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="Achtung"}
**Das Wails-Projekt pflegt keine Vorlagen von Drittanbietern und übernimmt dafür weder Verantwortung noch Haftung!3**

Wenn Sie sich bei einer Vorlage unsicher sind, prüfen Sie anhand von `package.json` und `wails.json`, welche Skripte ausgeführt und welche Pakete installiert werden.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) – Wails-Vorlage auf Basis des Vue-Ökosystems (integriertes TypeScript, dunkles Theme, Internationalisierung, Single-Page-Routing, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) – Eine Vorlage mit JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) – Eine Vorlage mit TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier, Composition API mit &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) – Wails-Vorlage auf Basis von Naive UI (einer Komponentenbibliothek für Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) – Wails-Vorlage mit unverändertem Nuxt3 und TypeScript sowie automatischen Imports für die Wails-JavaScript-Runtime
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) – Wails-Vorlage mit Vue+TypeScript+Vite+Element-plus (im Stil von NetEase Cloud Music)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) – Angular 15 mit umfangreicher Funktionalität und bereit für den Produktionseinsatz.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) – Angular mit TypeScript, Sass, Hot Reloading, Code-Splitting und i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) – Eine Vorlage mit ReactJS
- [wails-react-template](https://github.com/flin7/wails-react-template) – Eine minimale React-Vorlage mit Unterstützung für die Live-Entwicklung
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) – Eine Vorlage mit Next.js und TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) – Eine Vorlage mit Next.js, TypeScript und App Router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) – Eine Vorlage mit Next.js, TypeScript, App Router, src-Verzeichnis und Beispiel
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) – Eine Vorlage für React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) – Eine Vorlage mit Vite, React, TypeScript, TailwindCSS und shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) – Eine Vorlage mit Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) – Eine Vorlage mit Svelte und Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) – Eine Vorlage mit Svelte, Vite und TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) – Eine aktualisierte Vorlage mit Svelte v4.2.0, Vite und TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) – Eine Vorlage mit SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) – Eine Vorlage mit SvelteKit und Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) – Eine Vorlage mit Solid + TypeScript + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) – Eine Vorlage mit Solid + JavaScript + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) – Entwickeln Sie Ihre GUI-App mit funktionaler Programmierung und einer **flotten** Hot-Reload-Konfiguration :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) – Kombinieren Sie die Stärken :muscle: von Elm + Tailwind CSS + Wails! Hot Reloading wird unterstützt.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) – Verwenden Sie eine einzigartige Kombination aus reinem htmx für Interaktivität und templ zum Erstellen von Komponenten und Formularen

## Reines JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) – Eine Vorlage, die ausschließlich grundlegendes JavaScript, HTML und CSS enthält

## Lit (Webkomponenten)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) – Wails-Vorlage mit einem Frontend auf Basis von Lit, der Komponentenbibliothek Shoelace sowie vorkonfiguriertem Prettier und TypeScript.
