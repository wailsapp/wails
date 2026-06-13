---
title: 템플릿
description: Wails를 위한 프로젝트 템플릿 및 시작 키트
---

:::caution

이 페이지는 Wails v3 기준으로 오래되었을 수 있습니다.

:::

<!-- TODO: Update this link -->

이 페이지는 커뮤니티에서 지원하는 템플릿 목록을 제공합니다. 직접 템플릿을 빌드하려면
[템플릿](https://wails.io/docs/guides/templates) 가이드를 참조하세요.

:::tip[템플릿 제출 방법]

하단의 `Edit this page`를 클릭하여 자신의 템플릿을 추가할 수 있습니다.

:::

이 템플릿을 사용하려면 다음 명령어를 실행하세요.
`wails init -n "Your Project Name" -t [아래 링크[@version]]`

버전 접미사가 없으면 기본적으로 main 브랜치의 코드 템플릿이 사용됩니다.
버전 접미사가 있으면 해당 버전의 태그에 해당하는 코드 템플릿이 사용됩니다.

예시:
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

:::danger[주의]

**Wails 프로젝트는 제3자 템플릿을 유지 관리하지 않으며, 이에 대한 책임이나 의무를 지지 않습니다!**

템플릿에 대해 확신이 서지 않는다면 `package.json`과 `wails.json`을 확인하여
어떤 스크립트가 실행되고 어떤 패키지가 설치되는지 확인하세요.

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Vue 생태계를 기반으로 한 Wails
  템플릿 (TypeScript 통합, 다크 테마,
  국제화, 단일 페이지 라우팅, TailwindCSS 포함)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  JavaScript + Quasar V2를 사용하는 템플릿 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  TypeScript + Quasar V2를 사용하는 템플릿 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, &lt;script setup&gt;를 사용한 Composition API)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Naive UI(Vue 3 컴포넌트 라이브러리)를 기반으로 한 Wails 템플릿
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Wails JS
  런타임에 대한 자동 임포트를 갖춘 깔끔한 Nuxt3와 TypeScript를 사용하는 Wails
  템플릿
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Vue+TypeScript+Vite+Element-plus(NetEase Cloud Music 스타일)를 사용하는 Wails
  템플릿

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ 기능 가득 차 있으며 프로덕션 준비 완료.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  TypeScript, Sass, 핫 리로드, 코드 스플리팅 및 i18n이 포함된 Angular

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  reactjs를 사용하는 템플릿
- [wails-react-template](https://github.com/flin7/wails-react-template) - 라이브 개발을 지원하는
  React를 위한 최소 템플릿
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Next.js와 TypeScript를 사용하는
  템플릿
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  App 라우터가 포함된 Next.js와 TypeScript를 사용하는 템플릿
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  App 라우터 src + 예제가 포함된 Next.js와 TypeScript를 사용하는 템플릿
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  React + TypeScript + Vite + TailwindCSS를 위한 템플릿
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Vite, React, TypeScript, TailwindCSS 및 shadcn/ui가 포함된 템플릿

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Svelte를 사용하는 템플릿
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  Svelte와 Vite를 사용하는 템플릿
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  TailwindCSS v3가 포함된 Svelte와 Vite를 사용하는 템플릿
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Svelte v4.2.0과 Vite, TailwindCSS v3.3.3을 사용하는 업데이트된 템플릿
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  SvelteKit을 사용하는 템플릿
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  Sveltekit과 Shadcn-Svelte를 사용하는 템플릿

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Solid + Ts + Vite를 사용하는 템플릿
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Solid + Js + Vite를 사용하는 템플릿

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  함수형 프로그래밍과 **빠른** 핫 리로드 설정으로 GUI 앱을 개발하세요
  :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Elm + Tailwind CSS + Wails의 힘 :muscle:을 결합하세요! 핫 리로딩이
  지원됩니다.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  상호 작용을 위한 순수 htmx와 컴포넌트 및 폼 생성을 위한 templ의 독특한 조합을 사용하세요

## 순수 JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 기본 JavaScript, HTML 및 CSS만 포함된
  템플릿

## Lit (웹 컴포넌트)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  프론트엔드에 lit, Shoelace 컴포넌트 라이브러리 + 사전 구성된 prettier와 typescript를 제공하는
  Wails 템플릿