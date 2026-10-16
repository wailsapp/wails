---
title: "前端路由"
description: "在 Wails 应用程序中使用前端路由"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

前端路由是在单页应用程序中切换视图的常用方式。 本指南介绍了在使用 Wails 时，针对不同前端框架的推荐方法。

@tabs{sync-key="framework"}
[Vue]
Vue 中推荐使用 [Hash 模式](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode)进行路由：

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Hash 模式使用 URL 哈希渲染不同视图，并通过采用基于哈希的 URL 格式，避免 Wails 运行时 干扰路由。

[Angular]
Angular 中推荐使用 [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy)进行路由：

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

使用基于哈希的 URL，可确保路由在所有平台上都能与 Wails 的窗口处理机制正常协同工作。

[React]
React 中推荐使用 [HashRouter](https://reactrouter.com/en/main/router-components/hash-router)进行路由：

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

HashRouter 使用 URL 哈希而非路径，因此在所有平台上都能可靠地与 Wails 配合使用。

[Svelte]
Svelte 中推荐使用 [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router)进行路由：

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

svelte-spa-router 支持基于哈希的路由，因此与 Wails 应用程序兼容。

@end

## 为什么使用哈希路由？

Wails 将前端嵌入原生 WebView 窗口中。使用基于哈希的路由（#/page，而非 /page） 可以避免与以下内容发生冲突：

- Wails 运行时的内部路由
- 不同平台上的原生窗口 URL 处理机制
- 从非根路径提供的生产环境资源

## 故障排除

如果生产构建中的前端路由出现问题：

- 确保已将前端构建工具配置为输出适用于 <strong>Hash 模式</strong>路由的内容
- 对于基于 Vite 的项目，请将 `base: "./"` 添加到 `vite.config.js`
- 验证 `index.html` 是否正确处理 SPA 路由的回退
