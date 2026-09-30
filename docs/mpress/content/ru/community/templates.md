---
title: "Шаблоны"
description: "Шаблоны проектов и стартовые наборы для Wails"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
Эта страница может содержать устаревшую информацию для Wails v3.

@end

<!-- TODO: Update this link -->

На этой странице перечислены шаблоны, поддерживаемые сообществом. Чтобы создать собственный шаблон, обратитесь к руководству [«Шаблоны»](https://wails.io/docs/guides/templates).

@note{type="tip" title="Как отправить шаблон"}
Чтобы добавить свои шаблоны, нажмите `Edit this page` внизу страницы.

@end

Чтобы использовать эти шаблоны, выполните команду: `wails init -n "Your Project Name" -t [the link below[@version]]`

Если суффикс версии не указан, по умолчанию используется шаблон кода из основной ветки. Если суффикс версии указан, используется шаблон кода, соответствующий тегу этой версии.

Пример: `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="Внимание"}
**Проект Wails не сопровождает шаблоны 3-х лиц, не отвечает за них и не несёт за них юридической ответственности!**

Если вы не уверены в шаблоне, проверьте `package.json` и `wails.json`, чтобы узнать, какие скрипты выполняются и какие пакеты устанавливаются.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) — шаблон Wails на основе экосистемы Vue (интегрированы TypeScript, тёмная тема, интернационализация, одностраничная маршрутизация и TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) — шаблон на основе JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) — шаблон на основе TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier, Composition API с &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) — шаблон Wails на основе Naive UI (библиотеки компонентов для Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) — шаблон Wails на основе чистого Nuxt3 и TypeScript с автоматическим импортом среды выполнения Wails для JavaScript
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) — шаблон Wails на основе Vue+TypeScript+Vite+Element-plus (имитация NetEase Cloud Music)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) — Angular 15+ с широким набором возможностей, готовый к развёртыванию в рабочей среде.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) — Angular с TypeScript, Sass, горячей перезагрузкой, разделением кода и i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) — шаблон на основе ReactJS
- [wails-react-template](https://github.com/flin7/wails-react-template) — минимальный шаблон для React с поддержкой разработки в реальном времени
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) — шаблон на основе Next.js и TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) — шаблон на основе Next.js и TypeScript с App Router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) — шаблон на основе Next.js и TypeScript с App Router, каталогом src и примером
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) — шаблон для React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) — шаблон с Vite, React, TypeScript, TailwindCSS и shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) — шаблон на основе Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) — шаблон на основе Svelte и Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) — шаблон на основе Svelte и Vite с TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) — обновлённый шаблон на основе Svelte v4.2.0 и Vite с TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) — шаблон на основе SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) — шаблон на основе SvelteKit и Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) — шаблон на основе Solid + TypeScript + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) — шаблон на основе Solid + JavaScript + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) — разрабатывайте приложения с графическим интерфейсом, используя функциональное программирование и **быстрое** окружение с горячей перезагрузкой :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) — объедините мощь :muscle: Elm + Tailwind CSS + Wails! Поддерживается горячая перезагрузка.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) — используйте уникальное сочетание чистого htmx для интерактивности и templ для создания компонентов и форм

## Чистый JavaScript (без фреймворков)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) — шаблон, содержащий только базовые JavaScript, HTML и CSS

## Lit (веб-компоненты)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) — шаблон Wails с фронтендом на основе Lit, библиотекой компонентов Shoelace и предварительно настроенными Prettier и TypeScript.
