---
title: "Frontend-Routing"
description: "Frontend-Routing in Ihrer Wails-Anwendung verwenden"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

Frontend-Routing ist eine verbreitete Methode, um in einer Single-Page-Anwendung zwischen Ansichten zu wechseln. Dieser Leitfaden beschreibt die empfohlenen Ansätze für verschiedene Frontend-Frameworks bei der Verwendung von Wails.

@tabs{sync-key="framework"}
[Vue]
Für das Routing in Vue wird der [Hash-Modus](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode) empfohlen:

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Der Hash-Modus verwendet den URL-Hash, um verschiedene Ansichten darzustellen. Das hashbasierte URL-Format verhindert dabei Probleme durch Eingriffe der Wails-Laufzeit in das Routing.

[Angular]
Für das Routing in Angular wird [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy) empfohlen:

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

Hashbasierte URLs stellen sicher, dass das Routing auf allen Plattformen korrekt mit der Fensterverwaltung von Wails funktioniert.

[React]
Für das Routing in React wird [HashRouter](https://reactrouter.com/en/main/router-components/hash-router) empfohlen:

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

HashRouter verwendet URL-Hashes anstelle von Pfaden und funktioniert daher mit Wails auf allen Plattformen zuverlässig.

[Svelte]
Für das Routing in Svelte wird [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router) empfohlen:

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

svelte-spa-router unterstützt hashbasiertes Routing und ist daher mit Wails-Anwendungen kompatibel.

@end

## Warum Hash-Routing?

Wails bettet Ihr Frontend in ein natives Webview-Fenster ein. Hashbasiertes Routing (#/page anstelle von /page) vermeidet Konflikte mit:

- dem internen Routing der Wails-Laufzeit
- der nativen Verarbeitung von Fenster-URLs auf verschiedenen Plattformen
- Produktions-Assets, die unter anderen URL-Pfaden als dem Root-Pfad bereitgestellt werden

## Fehlerbehebung

Wenn bei Produktions-Builds Probleme mit dem Frontend-Routing auftreten:

- Stellen Sie sicher, dass Ihr Frontend-Build-Tool für die Ausgabe von Routing im **Hash-Modus** konfiguriert ist
- Fügen Sie bei Vite-basierten Projekten `base: "./"` zu Ihrer `vite.config.js` hinzu
- Überprüfen Sie, ob Ihr `index.html` den Fallback für SPA-Routing korrekt verarbeitet
