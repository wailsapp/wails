---
title: "Modèles"
description: "Modèles de projet et kits de démarrage pour Wails"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
Cette page est peut-être obsolète pour Wails v3.

@end

<!-- TODO: Update this link -->

Cette page répertorie les modèles pris en charge par la communauté. Pour créer votre propre modèle, consultez le guide sur les [modèles](https://wails.io/docs/guides/templates).

@note{type="tip" title="Comment proposer un modèle"}
Vous pouvez cliquer sur `Edit this page` en bas de la page pour ajouter vos modèles.

@end

Pour utiliser ces modèles, exécutez : `wails init -n "Your Project Name" -t [the link below[@version]]`

En l’absence de suffixe de version, le modèle de code de la branche principale est utilisé par défaut. Si un suffixe de version est présent, le modèle de code correspondant au tag de cette version est utilisé.

Exemple : `wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="Attention"}
**Le projet Wails n’assure pas la maintenance des modèles tiers (de 3e partie) et décline toute responsabilité à leur égard !**

Si vous avez des doutes sur un modèle, examinez `package.json` et `wails.json` afin de vérifier les scripts exécutés et les paquets installés.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) — Modèle Wails basé sur l’écosystème Vue (TypeScript intégré, thème sombre, internationalisation, routage monopage, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) — Modèle utilisant JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) — Modèle utilisant TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier, API de composition avec &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) — Modèle Wails basé sur Naive UI (une bibliothèque de composants Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) — Modèle Wails utilisant une installation épurée de Nuxt3 et TypeScript, avec importations automatiques pour le runtime JavaScript de Wails
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) — Modèle Wails utilisant Vue+TypeScript+Vite+Element-plus (imitant NetEase Cloud Music)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) — Angular 15+ riche en fonctionnalités et prêt pour la production.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) — Angular avec TypeScript, Sass, rechargement à chaud, fractionnement du code et i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) — Modèle utilisant ReactJS
- [wails-react-template](https://github.com/flin7/wails-react-template) — Modèle minimal pour React prenant en charge le développement avec rechargement en direct
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) — Modèle utilisant Next.js et TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) — Modèle utilisant Next.js et TypeScript avec App Router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) — Modèle utilisant Next.js et TypeScript avec App Router, un répertoire src et un exemple
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) — Modèle pour React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) — Modèle avec Vite, React, TypeScript, TailwindCSS et shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) — Modèle utilisant Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) — Modèle utilisant Svelte et Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) — Modèle utilisant Svelte et Vite avec TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) — Modèle mis à jour utilisant Svelte v4.2.0 et Vite avec TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) — Modèle utilisant SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) — Modèle utilisant SvelteKit et Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) — Modèle utilisant Solid + TypeScript + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) — Modèle utilisant Solid + JavaScript + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) — Développez votre application avec interface graphique en programmation fonctionnelle et profitez d’une configuration de rechargement à chaud **ultraréactive** :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) — Combinez la puissance :muscle: d’Elm + Tailwind CSS + Wails ! Rechargement à chaud pris en charge.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) — Utilisez une combinaison unique de htmx pur pour l’interactivité et de templ pour créer des composants et des formulaires

## JavaScript pur (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) — Modèle ne contenant que du JavaScript, du HTML et du CSS de base

## Lit (composants web)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) — Modèle Wails fournissant une interface frontend avec Lit, la bibliothèque de composants Shoelace, ainsi que Prettier et TypeScript préconfigurés.
