---
title: "前端路由"
description: "在 Wails 應用程式中使用前端路由"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

前端路由是在單頁應用程式中切換檢視畫面的常見方式。 本指南說明搭配 Wails 使用不同前端框架時的建議做法。

@tabs{sync-key="framework"}
[Vue]
在 Vue 中，建議使用 [Hash 模式](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode)進行路由：

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Hash 模式使用 URL 雜湊來呈現不同的檢視畫面；採用以雜湊為基礎的 URL 格式，可避免 Wails 執行階段 干擾路由所造成的問題。

[Angular]
在 Angular 中，建議使用 [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy)進行路由：

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

使用以雜湊為基礎的 URL，可確保路由在所有平台上都能與 Wails 的視窗處理機制正確搭配運作。

[React]
在 React 中，建議使用 [HashRouter](https://reactrouter.com/en/main/router-components/hash-router)進行路由：

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

HashRouter 使用 URL 雜湊而非路徑，因此在所有平台上都能與 Wails 穩定搭配運作。

[Svelte]
在 Svelte 中，建議使用 [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router)進行路由：

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

svelte-spa-router 支援以雜湊為基礎的路由，因此可與 Wails 應用程式相容。

@end

## 為什麼要使用 Hash 路由？

Wails 會將前端嵌入原生 WebView 視窗。使用以雜湊為基礎的路由（#/page，而非 /page） 可避免與下列項目發生衝突：

- Wails 執行階段的內部路由
- 不同平台上的原生視窗 URL 處理機制
- 從非根目錄路徑提供的正式環境資源

## 疑難排解

如果正式版本中的前端路由發生問題：

- 請確認前端建置工具已設定為輸出供 <strong>Hash 模式</strong>路由使用的內容
- 若是以 Vite 為基礎的專案，請將 `base: "./"`新增至 `vite.config.js`
- 請確認 `index.html`能正確處理 SPA 路由的備援
