---
title: "Templat"
description: "Templat proyek dan kit pemula untuk Wails"
slug: "community/templates"
sourcePath: "community/templates.md"
---

@note{type="caution"}
Halaman ini mungkin sudah tidak mutakhir untuk Wails v3.

@end

<!-- TODO: Update this link -->

Halaman ini berisi daftar templat yang didukung komunitas. Untuk membuat templat Anda sendiri, lihat panduan [Templat](https://wails.io/docs/guides/templates).

@note{type="tip" title="Cara Mengirimkan Templat"}
Anda dapat mengeklik `Edit this page` di bagian bawah untuk menyertakan templat Anda.

@end

Untuk menggunakan templat ini, jalankan  
`wails init -n "Your Project Name" -t [the link below[@version]]`

Jika tidak ada akhiran versi, templat kode dari cabang utama akan digunakan secara default.  
Jika ada akhiran versi, templat kode yang sesuai dengan tag versi tersebut akan digunakan.

Contoh:  
`wails init -n "Your Project Name" -t https://github.com/misitebao/wails-template-vue`

@note{type="danger" title="Perhatian"}
**Proyek Wails tidak memelihara, tidak bertanggung jawab, dan tidak menanggung tanggung jawab hukum atas templat pihak ke-3!**

Jika Anda ragu mengenai suatu templat, periksa `package.json` dan `wails.json` untuk mengetahui skrip yang dijalankan dan paket yang diinstal.

@end

## Vue

- [wails-template-vue](https://github.com/misitebao/wails-template-vue) - Templat Wails berbasis ekosistem Vue (TypeScript terintegrasi, tema gelap, internasionalisasi, perutean satu halaman, TailwindCSS)
- [wails-template-quasar-js](https://github.com/sgosiaco/wails-template-quasar-js) - Templat yang menggunakan JavaScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier)
- [wails-template-quasar-ts](https://github.com/sgosiaco/wails-template-quasar-ts) - Templat yang menggunakan TypeScript + Quasar V2 (Vue 3, Vite, Sass, Pinia, ESLint, Prettier, Composition API dengan &lt;script setup&gt;)
- [wails-template-naive](https://github.com/tk103331/wails-template-naive) - Templat Wails berbasis Naive UI (pustaka komponen Vue 3)
- [wails-template-nuxt](https://github.com/gornius/wails-template-nuxt) - Templat Wails yang menggunakan Nuxt3 murni dan TypeScript dengan impor otomatis untuk runtime JS Wails
- [Wails-Tool-Template](https://github.com/xisuo67/Wails-Tool-Template) - Templat Wails yang menggunakan Vue+TypeScript+Vite+Element-plus (meniru NetEase Cloud Music)

## Angular

- [wails-template-angular](https://github.com/mateothegreat/wails-template-angular) - Angular 15+ kaya fitur dan siap digunakan dalam produksi.
- [wails-angular-template](https://github.com/TAINCER/wails-angular-template) - Angular dengan TypeScript, Sass, hot reload, pemisahan kode, dan i18n

## React

- [wails-react-template](https://github.com/AlienRecall/wails-react-template) - Templat yang menggunakan ReactJS
- [wails-react-template](https://github.com/flin7/wails-react-template) - Templat minimal untuk React yang mendukung pengembangan langsung
- [wails-template-nextjs](https://github.com/LGiki/wails-template-nextjs) - Templat yang menggunakan Next.js dan TypeScript
- [wails-template-nextjs-app-router](https://github.com/thisisvk-in/wails-template-nextjs-app-router) - Templat yang menggunakan Next.js dan TypeScript dengan App Router
- [wails-template-nextjs-app-router-src](https://github.com/edai-git/wails-template-nextjs-app-router) - Templat yang menggunakan Next.js dan TypeScript dengan App Router, direktori src, dan contoh
- [wails-vite-react-ts-tailwind-template](https://github.com/hotafrika/wails-vite-react-ts-tailwind-template) - Templat untuk React + TypeScript + Vite + TailwindCSS
- [wails-vite-react-ts-tailwind-shadcnui-template](https://github.com/Mahcks/wails-vite-react-tailwind-shadcnui-ts) - Templat dengan Vite, React, TypeScript, TailwindCSS, dan shadcn/ui

## Svelte

- [wails-svelte-template](https://github.com/raitonoberu/wails-svelte-template) - Templat yang menggunakan Svelte
- [wails-vite-svelte-template](https://github.com/BillBuilt/wails-vite-svelte-template) - Templat yang menggunakan Svelte dan Vite
- [wails-vite-svelte-tailwind-template](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template) - Templat yang menggunakan Svelte dan Vite dengan TailwindCSS v3
- [wails-svelte-tailwind-vite-template](https://github.com/PylotLight/wails-vite-svelte-tailwind-template/tree/master) - Templat yang diperbarui menggunakan Svelte v4.2.0 dan Vite dengan TailwindCSS v3.3.3
- [wails-sveltekit-template](https://github.com/h8gi/wails-sveltekit-template) - Templat yang menggunakan SvelteKit
- [wails-template-shadcn-svelte](https://github.com/xijaja/wails-template-shadcn-svelte) - Templat yang menggunakan SvelteKit dan Shadcn-Svelte

## Solid

- [wails-template-vite-solid-ts](https://github.com/xijaja/wails-template-solid-ts) - Templat yang menggunakan Solid + TS + Vite
- [wails-template-vite-solid-js](https://github.com/xijaja/wails-template-solid-js) - Templat yang menggunakan Solid + JS + Vite

## Elm

- [wails-elm-template](https://github.com/benjamin-thomas/wails-elm-template) - Kembangkan aplikasi GUI Anda dengan pemrograman fungsional dan konfigurasi hot reload yang **responsif** :tada: :rocket:
- [wails-template-elm-tailwind](https://github.com/rnice01/wails-template-elm-tailwind) - Gabungkan keunggulan :muscle: Elm + Tailwind CSS + Wails! Mendukung hot reload.

## HTMX

- [wails-htmx-templ-chi-tailwind](https://github.com/PylotLight/wails-hmtx-templ-template) - Gunakan kombinasi unik htmx murni untuk interaktivitas dan templ untuk membuat komponen serta formulir

## JavaScript Murni (Vanilla)

- [wails-pure-js-template](https://github.com/KiddoV/wails-pure-js-template) - Templat yang hanya berisi JavaScript, HTML, dan CSS dasar

## Lit (komponen web)

- [wails-lit-shoelace-esbuild-template](https://github.com/Braincompiler/wails-lit-shoelace-esbuild-template) - Templat Wails yang menyediakan frontend dengan Lit, pustaka komponen Shoelace, serta Prettier dan TypeScript yang telah dikonfigurasi sebelumnya.
