---
title: "Routage frontend"
description: "Utiliser le routage frontend dans votre application Wails"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

Le routage frontend est une méthode courante pour passer d’une vue à une autre dans une application monopage. Ce guide présente les approches recommandées pour différents frameworks frontend avec Wails.

@tabs{sync-key="framework"}
[Vue]
L’approche recommandée pour le routage dans Vue est le [mode hash](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode) :

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Le mode hash utilise le fragment d’URL pour afficher différentes vues et évite ainsi que le runtime Wails interfère avec le routage, grâce au format d’URL basé sur un fragment.

[Angular]
L’approche recommandée pour le routage dans Angular est [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy) :

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

L’utilisation d’URL basées sur des fragments garantit le bon fonctionnement du routage avec la gestion des fenêtres de Wails sur toutes les plateformes.

[React]
L’approche recommandée pour le routage dans React est [HashRouter](https://reactrouter.com/en/main/router-components/hash-router) :

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

HashRouter utilise des fragments d’URL plutôt que des chemins, ce qui assure un fonctionnement fiable avec Wails sur toutes les plateformes.

[Svelte]
L’approche recommandée pour le routage dans Svelte est [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router) :

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

svelte-spa-router prend en charge le routage basé sur des fragments, ce qui le rend compatible avec les applications Wails.

@end

## Pourquoi utiliser le routage par fragment ?

Wails intègre votre frontend dans une fenêtre de vue web native. L’utilisation du routage basé sur des fragments (#/page au lieu de /page) évite les conflits avec :

- Le routage interne du runtime Wails
- La gestion native des URL des fenêtres sur les différentes plateformes
- Les ressources de production servies depuis des chemins autres que la racine

## Résolution des problèmes

Si vous rencontrez des problèmes de routage frontend dans les builds de production :

- Assurez-vous que votre outil de build frontend est configuré pour produire une sortie destinée au routage en **mode hash**
- Pour les projets basés sur Vite, ajoutez `base: "./"` à votre `vite.config.js`
- Vérifiez que votre `index.html` gère correctement le mécanisme de repli pour le routage des applications monopages
