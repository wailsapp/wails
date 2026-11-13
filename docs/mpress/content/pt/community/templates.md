---
title: "Templates"
description: "Templates de projeto e kits iniciais para Wails"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
Esta página pode estar desatualizada para o Wails v3.

@end

<!-- TODO: Update this link -->

Esta página apresenta uma lista de templates mantidos pela comunidade. Para criar seu próprio template, consulte o guia de [Templates](https://v2.wails.io/docs/guides/templates).

@note{type="tip" title="Como enviar um template"}
Clique em `Edit this page` na parte inferior para incluir seus templates.

@end

Para usar esses templates, execute `wails init -n "Your Project Name" -t [the link below[@version]]`

Se não houver um sufixo de versão, o template de código da branch principal será usado por padrão. Se houver um sufixo de versão, será usado o template de código correspondente à tag dessa versão.

Exemplo: `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="Atenção"}
**O projeto Wails não mantém nem se responsabiliza por templates de terceiros!3**

Se tiver dúvidas sobre um template, verifique `package.json` e `wails.json` para saber quais scripts são executados e quais pacotes são instalados.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) — Template do Wails baseado no ecossistema Vue (TypeScript integrado, tema escuro, internacionalização, roteamento de página única e TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) — Um template que usa JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint e Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) — Um template que usa TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier e Composition API com &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) — Template do Wails baseado em Naive UI (uma biblioteca de componentes Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) — Template do Wails que usa Nuxt3 puro e TypeScript, com importações automáticas para o runtime JavaScript do Wails
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) — Template do Wails que usa Vue+TypeScript+Vite+Element-plus (仿网易云)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) — Angular 15 repleto de recursos e pronto para produção.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) — Angular com TypeScript, Sass, recarregamento automático, divisão de código e i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) — Um template que usa ReactJS
- [wails-react-template](https://github.com/flin7/wails-react-template) — Um template mínimo para React com suporte ao desenvolvimento ao vivo
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) — Um template que usa Next.js e TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) — Um template que usa Next.js e TypeScript com App Router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) — Um template que usa Next.js e TypeScript com App Router, diretório src e exemplo
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) — Um template para React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) — Um template com Vite, React, TypeScript, TailwindCSS e shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) — Um template que usa Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) — Um template que usa Svelte e Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) — Um template que usa Svelte e Vite com TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) — Um template atualizado que usa Svelte v4.2.0 e Vite com TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) — Um template que usa SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) — Um template que usa SvelteKit e Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) — Um template que usa Solid + TypeScript + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) — Um template que usa Solid + JavaScript + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) — Desenvolva seu aplicativo com interface gráfica usando programação funcional e uma configuração de recarregamento automático **ágil** :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) — Combine o poder :muscle: do Elm + Tailwind CSS + Wails! Compatível com recarregamento automático.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) — Use uma combinação exclusiva de htmx puro para interatividade com templ para criar componentes e formulários

## JavaScript puro (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) — Um template que contém apenas JavaScript, HTML e CSS básicos

## Lit (componentes Web)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) — Template do Wails que fornece um frontend com Lit, a biblioteca de componentes Shoelace, além de Prettier e TypeScript pré-configurados.
