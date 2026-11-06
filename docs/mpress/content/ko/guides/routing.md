---
title: "프런트엔드 라우팅"
description: "Wails 애플리케이션에서 프런트엔드 라우팅 사용하기"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

프런트엔드 라우팅은 단일 페이지 애플리케이션에서 뷰를 전환하는 데 널리 사용되는 방식입니다. 이 가이드에서는 Wails를 사용할 때 프런트엔드 프레임워크별로 권장되는 접근 방식을 설명합니다.

@tabs{sync-key="framework"}
[Vue]
Vue에서 권장되는 라우팅 방식은 [해시 모드](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode)입니다:

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

해시 모드는 URL 해시를 사용하여 서로 다른 뷰를 렌더링합니다. 해시 기반 URL 형식을 사용하므로 Wails 런타임이 라우팅에 간섭하여 발생하는 문제를 방지할 수 있습니다.

[Angular]
Angular에서 권장되는 라우팅 방식은 [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy)입니다:

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

해시 기반 URL을 사용하면 모든 플랫폼에서 Wails의 창 처리 방식과 함께 라우팅이 올바르게 작동합니다.

[React]
React에서 권장되는 라우팅 방식은 [HashRouter](https://reactrouter.com/en/main/router-components/hash-router)입니다:

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

HashRouter는 경로 대신 URL 해시를 사용하므로 모든 플랫폼에서 Wails와 안정적으로 작동합니다.

[Svelte]
Svelte에서 권장되는 라우팅 방식은 [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router)입니다:

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

svelte-spa-router는 해시 기반 라우팅을 지원하므로 Wails 애플리케이션과 호환됩니다.

@end

## 해시 라우팅을 사용해야 하는 이유

Wails는 프런트엔드를 네이티브 웹뷰 창에 임베드합니다. 해시 기반 라우팅(/page 대신 #/page)을 사용하면 다음 항목과의 충돌을 방지할 수 있습니다:

- Wails 런타임의 내부 라우팅
- 플랫폼별 네이티브 창의 URL 처리 방식
- 루트가 아닌 경로에서 제공되는 프로덕션 애셋

## 문제 해결

프로덕션 빌드에서 프런트엔드 라우팅에 문제가 발생하는 경우 다음을 확인하세요:

- 프런트엔드 빌드 도구가 **해시 모드** 라우팅용으로 출력하도록 구성되어 있는지 확인하세요
- Vite 기반 프로젝트에서는 `vite.config.js`에 `base: "./"`을 추가하세요
- `index.html`이 SPA 라우팅을 위한 폴백을 올바르게 처리하는지 확인하세요
