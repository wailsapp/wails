---
title: "Perutean Frontend"
description: "Menggunakan perutean frontend dalam aplikasi Wails Anda"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

Perutean frontend merupakan cara populer untuk beralih tampilan dalam aplikasi satu halaman. Panduan ini membahas pendekatan yang direkomendasikan untuk berbagai framework frontend saat menggunakan Wails.

@tabs{sync-key="framework"}
[Vue]
Pendekatan yang direkomendasikan untuk perutean di Vue adalah [Mode Hash](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode):

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Mode hash menggunakan hash URL untuk merender tampilan yang berbeda sehingga menghindari masalah akibat runtime Wails yang mengganggu perutean, dengan menggunakan format URL berbasis hash.

[Angular]
Pendekatan yang direkomendasikan untuk perutean di Angular adalah [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy):

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

Penggunaan URL berbasis hash memastikan perutean berfungsi dengan benar bersama penanganan jendela Wails di semua platform.

[React]
Pendekatan yang direkomendasikan untuk perutean di React adalah [HashRouter](https://reactrouter.com/en/main/router-components/hash-router):

```jsx
import ReactDOM from "react-dom/client";
import { HashRouter, Routes, Route } from "react-router-dom";

ReactDOM.createRoot(root).render(
  <HashRouter basename={"/"}>
    {/* The rest of your app goes here */}
    <Routes>
      <Route path="/" element={<Page0 />} />
      <Route path="/page1" element={<Page1 />} />
      <Route path="/page2" element={<Page2 />} />
      {/* more... */}
    </Routes>
  </HashRouter>
);
```

HashRouter menggunakan hash URL, bukan jalur, sehingga berfungsi secara andal dengan Wails di semua platform.

[Svelte]
Pendekatan yang direkomendasikan untuk perutean di Svelte adalah [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router):

```svelte
<script>
    import Router from "svelte-spa-router";
</script>

<Router
    routes={{
        "/": Home,
        "/products": wrap({
            asyncComponent: () => import("./routes/Products.svelte"),
        }),
        "/settings": Settings,
        "*": NotFound,
    }}
/>
```

svelte-spa-router mendukung perutean berbasis hash sehingga kompatibel dengan aplikasi Wails.

@end

## Mengapa Menggunakan Perutean Hash?

Wails menyematkan frontend Anda ke dalam jendela webview native. Penggunaan perutean berbasis hash (#/page, bukan /page) menghindari konflik dengan:

- Perutean internal runtime Wails
- Penanganan URL jendela native pada berbagai platform
- Aset produksi yang disajikan dari jalur selain root

## Pemecahan Masalah

Jika Anda mengalami masalah dengan perutean frontend dalam build produksi:

- Pastikan alat build frontend Anda dikonfigurasi untuk menghasilkan output bagi perutean **mode hash**
- Untuk proyek berbasis Vite, tambahkan `base: "./"` ke `vite.config.js`
- Pastikan `index.html` Anda menangani fallback dengan benar untuk perutean SPA
