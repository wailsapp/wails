---
title: Шаблоны
description: Шаблоны проектов и стартовые наборы для Wails
---

:::caution

Эта страница может быть устаревшей для Wails v3.

:::

<!-- TODO: Update this link -->

Эта страница служит списком шаблонов, поддерживаемых сообществом. Чтобы создать
собственный шаблон, пожалуйста, ознакомьтесь с руководством [Шаблоны](https://wails.io/docs/guides/templates).

:::tip[Как отправить шаблон]

Вы можете нажать `Edit this page` внизу страницы, чтобы добавить свои шаблоны.

:::

Чтобы использовать эти шаблоны, выполните
`wails init -n "Your Project Name" -t [ссылка ниже[@version]]`

Если суффикс версии отсутствует, по умолчанию используется шаблон кода из основной ветки.
Если суффикс версии присутствует, используется шаблон кода, соответствующий тегу этой
версии.

Пример:
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

:::danger[Внимание]

**Проект Wails не обслуживает, не несет ответственности и не гарантирует работу шаблонов
сторонних разработчиков!**

Если вы не уверены в шаблоне, проверьте `package.json` и `wails.json`, чтобы узнать,
какие скрипты запускаются и какие пакеты устанавливаются.

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Шаблон Wails
  на базе экосистемы Vue (Интегрирован TypeScript, Темная тема,
  Интернационализация, Маршрутизация SPA, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  Шаблон, использующий JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  Шаблон, использующий TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, Composition API с &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Шаблон Wails на базе Naive UI (Библиотека компонентов Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Шаблон Wails
  с использованием чистого Nuxt3 и TypeScript с автоимпортами для wails js
  runtime
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Шаблон Wails
  с использованием Vue+TypeScript+Vite+Element-plus (клон NetEase Cloud Music)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ насыщенный функциями и готовый к развертыванию в продакшен.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  Angular с TypeScript, Sass, Hot-Reload, Code-Splitting и i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  Шаблон, использующий reactjs
- [wails-react-template](https://github.com/flin7/wails-react-template) -
  Минимальный шаблон для React, поддерживающий живую разработку
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) -
  Шаблон, использующий Next.js и TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  Шаблон, использующий Next.js и TypeScript с App router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  Шаблон, использующий Next.js и TypeScript с App router src + пример
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  Шаблон для React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Шаблон с Vite, React, TypeScript, TailwindCSS и shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Шаблон, использующий Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  Шаблон, использующий Svelte и Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  Шаблон, использующий Svelte и Vite с TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Обновленный шаблон, использующий Svelte v4.2.0 и Vite с TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  Шаблон, использующий SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  Шаблон, использующий Sveltekit и Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Шаблон, использующий Solid + Ts + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Шаблон, использующий Solid + Js + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  Разрабатывайте свое GUI-приложение с помощью функционального программирования и **мгновенной**
  настройки hot-reload :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Объедините мощь :muscle: Elm + Tailwind CSS + Wails! Поддерживается горячая перезагрузка.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  Используйте уникальное сочетание чистого htmx для интерактивности и templ для
  создания компонентов и форм

## Чистый JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - Шаблон,
  содержащий только базовый JavaScript, HTML и CSS

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  Шаблон Wails, предоставляющий фронтенд на базе lit, библиотеку компонентов Shoelace +
  предварительно настроенные prettier и typescript.