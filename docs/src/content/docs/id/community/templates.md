---
title: Template
description: Template proyek dan starter kit untuk Wails
---

:::caution

Halaman ini mungkin sudah usang untuk Wails v3.

:::

<!-- TODO: Perbarui tautan ini -->

Halaman ini berfungsi sebagai daftar template yang didukung komunitas. Untuk membangun template
Anda sendiri, lihat panduan [Template](https://wails.io/docs/guides/templates).

:::tip[Cara Mengirimkan Template]

Anda dapat mengklik `Edit this page` di bagian bawah untuk memasukkan template Anda.

:::

Untuk menggunakan template ini, jalankan
`wails init -n "Nama Proyek Anda" -t [tautan di bawah[@versi]]`

Jika tidak ada sufiks versi, template kode branch main digunakan secara default.
Jika ada sufiks versi, template kode yang sesuai dengan tag versi
tersebut digunakan.

Contoh:
`wails init -n "Nama Proyek Anda" -t https://github.com/misitebao/wails-template-vue`

:::danger[Perhatian]

**Proyek Wails tidak memelihara, tidak bertanggung jawab, dan tidak berkewajiban atas
template pihak ketiga!**

Jika Anda ragu tentang template, periksa `package.json` dan `wails.json` untuk
skrip apa yang dijalankan dan paket apa yang diinstal.

:::

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Template
  Wails berbasis ekosistem Vue (TypeScript terintegrasi, tema gelap,
  internasionalisasi, routing halaman tunggal, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) -
  Template menggunakan JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) -
  Template menggunakan TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint,
  Prettier, Composition API dengan &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) -
  Template Wails berbasis Naive UI (Perpustakaan Komponen Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Template
  Wails menggunakan Nuxt3 bersih dan TypeScript dengan auto-import untuk wails js
  runtime
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Template
  Wails menggunakan Vue+TypeScript+Vite+Element-plus(仿网易云)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) -
  Angular 15+ lengkap dan siap produksi.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) -
  Angular dengan TypeScript, Sass, Hot-Reload, Code-Splitting, dan i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) -
  Template menggunakan reactjs
- [wails-react-template](https://github.com/flin7/wails-react-template) - Template
  minimal untuk React yang mendukung live development
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Template
  menggunakan Next.js dan TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) -
  Template menggunakan Next.js dan TypeScript dengan App router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) -
  Template menggunakan Next.js dan TypeScript dengan App router src + contoh
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) -
  Template untuk React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) -
  Template dengan Vite, React, TypeScript, TailwindCSS, dan shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) -
  Template menggunakan Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) -
  Template menggunakan Svelte dan Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) -
  Template menggunakan Svelte dan Vite dengan TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) -
  Template diperbarui menggunakan Svelte v4.2.0 dan Vite dengan TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) -
  Template menggunakan SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) -
  Template menggunakan Sveltekit dan Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) -
  Template menggunakan Solid + Ts + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) -
  Template menggunakan Solid + Js + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) -
  Kembangkan aplikasi GUI Anda dengan pemrograman fungsional dan setup hot-reload
  yang **responsif** :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) -
  Gabungkan kekuatan :muscle: Elm + Tailwind CSS + Wails! Hot reloading
  didukung.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) -
  Gunakan kombinasi unik htmx murni untuk interaktivitas plus templ untuk
  membuat komponen dan form

## Pure JavaScript (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) -
  Template dengan JavaScript, HTML, dan CSS dasar saja

## Lit (web components)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) -
  Template Wails dengan frontend lit, perpustakaan komponen Shoelace +
  prettier dan typescript yang sudah dikonfigurasi.
