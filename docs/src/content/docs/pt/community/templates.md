---
title: Modelos
description: Modelos de projeto e kits iniciais para Wails
---

:::caution

Esta página pode estar desatualizada para o Wails v3.

:::

<!-- TODO: Update this link -->

Esta página serve como uma lista de modelos suportados pela comunidade. Para criar seu próprio
modelo, consulte o guia [Modelos](https://wails.io/docs/guides/templates).

:::tip[Como Submeter um Modelo]

Você pode clicar em `Editar esta página` na parte inferior para incluir seus modelos.

:::

Para usar estes modelos, execute
`wails init -n "Nome do Seu Projeto" -t [o link abaixo[@version]]`

Se não houver sufixo de versão, o modelo de código do branch principal é usado por padrão.
Se houver um sufixo de versão, o modelo de código correspondente à tag desta
versão é usado.

Exemplo:
`wails init -n "Nome do Seu Projeto" -t https://github.com/misitebao/wails-template-vue`

:::danger[Atenção]

**O projeto Wails não mantém, não é responsável nem tem responsabilidade por modelos de terceiros!**

Se você tiver dúvidas sobre um modelo, inspecione `package.json` e `wails.json` para
ver quais scripts são executados e quais pacotes estão instalados.

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Modelo Wails
  baseado na ecologia Vue (TypeScript integrado, tema escuro,
  internacionalização, roteamento de página única, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  Um modelo usando JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  Um modelo usando TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, Composition API com &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Modelo Wails baseado no Naive UI (Uma biblioteca de componentes Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Modelo Wails
  usando Nuxt3 limpo e TypeScript com importação automática para o
  runtime JS do Wails
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Modelo Wails
  usando Vue+TypeScript+Vite+Element-plus (inspirado no NetEase Cloud)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ repleto de recursos e pronto para produção.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  Angular com TypeScript, Sass, Hot-Reload, Code-Splitting e i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  Um modelo usando reactjs
- [wails-react-template](https://github.com/flin7/wails-react-template) - Um
  modelo minimalista para React que suporta desenvolvimento em tempo real
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Um
  modelo usando Next.js e TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  Um modelo usando Next.js e TypeScript com roteador de aplicativo
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  Um modelo usando Next.js e TypeScript com roteador de aplicativo src + exemplo
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  Um modelo para React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Um modelo com Vite, React, TypeScript, TailwindCSS e shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Um modelo usando Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  Um modelo usando Svelte e Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  Um modelo usando Svelte e Vite com TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Um modelo atualizado usando Svelte v4.2.0 e Vite com TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  Um modelo usando SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  Um modelo usando Sveltekit e Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Um modelo usando Solid + Ts + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Um modelo usando Solid + Js + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  Desenvolva seu aplicativo GUI com programação funcional e uma configuração de
  hot-reload **rápida** :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Combine os poderes :muscle: de Elm + Tailwind CSS + Wails! Hot reloading
  suportado.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  Use uma combinação única de htmx puro para interatividade mais templ para
  criação de componentes e formulários

## JavaScript Puro (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - Um
  modelo com nada além de JavaScript básico, HTML e CSS

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  Modelo Wails fornecendo ao frontend lit, biblioteca de componentes Shoelace +
  prettier e typescript pré-configurados.