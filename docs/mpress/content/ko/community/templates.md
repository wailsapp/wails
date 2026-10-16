---
title: "템플릿"
description: "Wails용 프로젝트 템플릿 및 스타터 키트"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
이 페이지의 내용은 Wails v3에 맞게 업데이트되지 않았을 수 있습니다.

@end

<!-- TODO: Update this link -->

이 페이지에는 커뮤니티에서 지원하는 템플릿이 나열되어 있습니다. 자체 템플릿을 만들려면 [템플릿](https://wails.io/docs/guides/templates) 가이드를 참조하세요.

@note{type="tip" title="템플릿 제출 방법"}
하단의 `Edit this page`을 클릭하여 템플릿을 추가할 수 있습니다.

@end

이 템플릿을 사용하려면 다음 명령을 실행하세요.  
`wails init -n "Your Project Name" -t [the link below[@version]]`

버전 접미사가 없으면 기본적으로 main 브랜치의 코드 템플릿을 사용합니다.  
버전 접미사가 있으면 해당 버전의 태그에 대응하는 코드 템플릿을 사용합니다.

예:  
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="주의"}
**Wails 프로젝트는 제3자 템플릿을 유지 관리하지 않으며, 이에 대한 책임이나 법적 책임을 지지 않습니다!**

템플릿이 확실하지 않다면 어떤 스크립트가 실행되고 어떤 패키지가 설치되는지 `package.json`과 `wails.json`에서 확인하세요.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Vue 생태계를 기반으로 하는 Wails 템플릿(TypeScript, 다크 테마, 국제화, 단일 페이지 라우팅, TailwindCSS 통합)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) - JavaScript + Quasar V2를 사용하는 템플릿(Vue 3, Vite, Sass, Pinia, ESLint, Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) - TypeScript + Quasar V2를 사용하는 템플릿(Vue 3, Vite, Sass, Pinia, ESLint, Prettier, &lt;script setup&gt;을 사용하는 Composition API)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) - Vue 3 컴포넌트 라이브러리인 Naive UI를 기반으로 하는 Wails 템플릿
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - 깔끔한 Nuxt3 및 TypeScript를 사용하고 Wails JS 런타임 자동 가져오기를 지원하는 Wails 템플릿
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Vue+TypeScript+Vite+Element-plus를 사용하는 Wails 템플릿(NetEase Cloud Music을 모방)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) - 다양한 기능을 갖추고 프로덕션에 바로 배포할 수 있는 Angular 15+ 템플릿
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) - TypeScript, Sass, 핫 리로드, 코드 분할 및 i18n을 지원하는 Angular 템플릿

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) - ReactJS를 사용하는 템플릿
- [wails-react-template](https://github.com/flin7/wails-react-template) - 라이브 개발을 지원하는 최소 구성의 React 템플릿
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Next.js와 TypeScript를 사용하는 템플릿
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) - App Router와 함께 Next.js 및 TypeScript를 사용하는 템플릿
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) - App Router src 및 예제와 함께 Next.js 및 TypeScript를 사용하는 템플릿
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) - React + TypeScript + Vite + TailwindCSS용 템플릿
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) - Vite, React, TypeScript, TailwindCSS 및 shadcn/ui를 사용하는 템플릿

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) - Svelte를 사용하는 템플릿
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) - Svelte와 Vite를 사용하는 템플릿
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) - TailwindCSS v3와 함께 Svelte 및 Vite를 사용하는 템플릿
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) - Svelte v4.2.0, Vite 및 TailwindCSS v3.3.3을 사용하는 업데이트된 템플릿
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) - SvelteKit을 사용하는 템플릿
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) - SvelteKit과 Shadcn-Svelte를 사용하는 템플릿

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) - Solid + TypeScript + Vite를 사용하는 템플릿
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) - Solid + JavaScript + Vite를 사용하는 템플릿

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) - 함수형 프로그래밍과 **빠른** 핫 리로드 설정으로 GUI 앱을 개발하세요 :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) - Elm + Tailwind CSS + Wails의 강력한 기능을 결합하세요 :muscle:! 핫 리로드를 지원합니다.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) - 상호작용을 위한 순수 htmx와 컴포넌트 및 폼 생성을 위한 templ의 독특한 조합을 사용합니다.

## 순수 JavaScript(바닐라)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - 기본 JavaScript, HTML 및 CSS만 포함하는 템플릿

## Lit(웹 컴포넌트)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) - Lit과 Shoelace 컴포넌트 라이브러리로 프런트엔드를 제공하며 Prettier와 TypeScript가 미리 구성된 Wails 템플릿
