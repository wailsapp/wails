---
title: Vorlagen
description: Projektvorlagen und Starter-Kits für Wails
---

:::caution

Diese Seite könnte für Wails v3 veraltet sein.

:::

<!-- TODO: Update this link -->

Diese Seite dient als Liste von von der Community unterstützten Vorlagen. Um Ihre eigene
Vorlage zu erstellen, sehen Sie bitte im Leitfaden [Vorlagen](https://wails.io/docs/guides/templates)
nach.

:::tip[Wie man eine Vorlage einreicht]

Sie können unten auf `Edit this page` klicken, um Ihre Vorlagen einzufügen.

:::

Um diese Vorlagen zu verwenden, führen Sie
`wails init -n "Ihr Projektname" -t [der Link unten[@version]]` aus.

Wenn kein Versions-Suffix vorhanden ist, wird standardmäßig die Vorlage des main-Branches verwendet.
Wenn ein Versions-Suffix vorhanden ist, wird die Vorlage verwendet, die dem Tag dieser
Version entspricht.

Beispiel:
`wails init -n "Ihr Projektname" -t https://github.com/misitebao/wails-template-vue`

:::danger[Achtung]

**Das Wails-Projekt pflegt keine, ist nicht verantwortlich und haftet nicht für
Vorlagen von Drittanbietern!**

Wenn Sie sich bei einer Vorlage unsicher sind, prüfen Sie `package.json` und `wails.json`,
welche Skripte ausgeführt werden und welche Pakete installiert sind.

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Wails
  Vorlage basierend auf der Vue-Ökologie (Integriertes TypeScript, Dunkles Thema,
  Internationalisierung, Single-Page-Routing, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  Eine Vorlage, die JavaScript + Quasar V2 verwendet (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  Eine Vorlage, die TypeScript + Quasar V2 verwendet (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, Composition API mit &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Wails Vorlage basierend auf Naive UI (eine Vue 3 Komponentenbibliothek)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Wails
  Vorlage, die sauberes Nuxt3 und TypeScript mit Auto-Imports für die Wails JS
  Laufzeit verwendet
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Wails
  Vorlage, die Vue+TypeScript+Vite+Element-plus (im Stil von NetEase Cloud) verwendet

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ actiongeladen & bereit für die Produktion.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  Angular mit TypeScript, Sass, Hot-Reload, Code-Splitting und i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  Eine Vorlage, die reactjs verwendet
- [wails-react-template](https://github.com/flin7/wails-react-template) - Eine
  minimale Vorlage für React, die Live-Entwicklung unterstützt
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Eine
  Vorlage, die Next.js und TypeScript verwendet
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  Eine Vorlage, die Next.js und TypeScript mit App-Router verwendet
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  Eine Vorlage, die Next.js und TypeScript mit App-Router src + Beispiel verwendet
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  Eine Vorlage für React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Eine Vorlage mit Vite, React, TypeScript, TailwindCSS und shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Eine Vorlage, die Svelte verwendet
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  Eine Vorlage, die Svelte und Vite verwendet
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  Eine Vorlage, die Svelte und Vite mit TailwindCSS v3 verwendet
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Eine aktualisierte Vorlage, die Svelte v4.2.0 und Vite mit TailwindCSS v3.3.3 verwendet
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  Eine Vorlage, die SvelteKit verwendet
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  Eine Vorlage, die Sveltekit und Shadcn-Svelte verwendet

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Eine Vorlage, die Solid + Ts + Vite verwendet
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Eine Vorlage, die Solid + Js + Vite verwendet

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  Entwickeln Sie Ihre GUI-App mit funktionaler Programmierung und einem **schnellen**
  Hot-Reload-Setup :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Kombinieren Sie die Kräfte :muscle: von Elm + Tailwind CSS + Wails! Hot-Reloading
  unterstützt.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  Verwenden Sie eine einzigartige Kombination aus purem htmx für Interaktivität sowie templ zum
  Erstellen von Komponenten und Formularen

## Reines JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - Eine
  Vorlage mit nichts als grundlegendem JavaScript, HTML und CSS

## Lit (Web-Komponenten)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  Wails-Vorlage, die dem Frontend lit, die Shoelace-Komponentenbibliothek +
  vorkonfiguriertes Prettier und TypeScript bietet.