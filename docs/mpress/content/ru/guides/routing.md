---
title: "Маршрутизация во фронтенде"
description: "Использование маршрутизации во фронтенде приложения Wails"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

Маршрутизация во фронтенде — популярный способ переключения представлений в одностраничном приложении. В этом руководстве описаны рекомендуемые подходы для разных фронтенд-фреймворков при использовании Wails.

@tabs{sync-key="framework"}
[Vue]
Для маршрутизации во Vue рекомендуется использовать [режим хеш-маршрутизации](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode):

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

В режиме хеш-маршрутизации для отображения разных представлений используется хеш URL. Формат URL на основе хеша позволяет избежать конфликтов между маршрутизацией и средой выполнения Wails.

[Angular]
Для маршрутизации в Angular рекомендуется использовать [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy):

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

URL на основе хеша обеспечивают корректную работу маршрутизации с механизмом управления окнами Wails на всех платформах.

[React]
Для маршрутизации в React рекомендуется использовать [HashRouter](https://reactrouter.com/en/main/router-components/hash-router):

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

HashRouter использует хеши URL вместо путей и поэтому надёжно работает с Wails на всех платформах.

[Svelte]
Для маршрутизации в Svelte рекомендуется использовать [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router):

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

svelte-spa-router поддерживает маршрутизацию на основе хеша, поэтому он совместим с приложениями Wails.

@end

## Зачем нужна хеш-маршрутизация?

Wails встраивает фронтенд в нативное окно webview. Маршрутизация на основе хеша (#/page вместо /page) позволяет избежать конфликтов со следующими компонентами и сценариями:

- Внутренняя маршрутизация среды выполнения Wails
- Обработка URL нативных окон на разных платформах
- Ресурсы рабочей сборки, предоставляемые не из корневых путей

## Устранение неполадок

Если в рабочих сборках возникают проблемы с маршрутизацией во фронтенде:

- Убедитесь, что инструмент сборки фронтенда настроен на создание выходных файлов для маршрутизации в **режиме хеш-маршрутизации**
- В проектах на основе Vite добавьте `base: "./"` в `vite.config.js`
- Убедитесь, что `index.html` правильно обрабатывает резервный маршрут для SPA
