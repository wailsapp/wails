---
title: "Roteamento no frontend"
description: "Como usar o roteamento no frontend da sua aplicação Wails"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

O roteamento no frontend é uma forma popular de alternar entre visualizações em uma aplicação de página única. Este guia apresenta as abordagens recomendadas para diferentes frameworks de frontend ao usar o Wails.

@tabs{sync-key="framework"}
[Vue]
A abordagem recomendada para roteamento no Vue é o [modo hash](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode):

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

O modo hash usa o hash da URL para renderizar diferentes visualizações, evitando que a interferência do runtime do Wails no roteamento cause problemas, graças ao uso do formato de URL baseado em hash.

[Angular]
A abordagem recomendada para roteamento no Angular é a [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy):

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

O uso de URLs baseadas em hash garante que o roteamento funcione corretamente com o gerenciamento de janelas do Wails em todas as plataformas.

[React]
A abordagem recomendada para roteamento no React é o [HashRouter](https://reactrouter.com/en/main/router-components/hash-router):

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

O HashRouter usa hashes de URL em vez de caminhos, o que funciona de maneira confiável com o Wails em todas as plataformas.

[Svelte]
A abordagem recomendada para roteamento no Svelte é o [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router):

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

O svelte-spa-router é compatível com roteamento baseado em hash, o que o torna compatível com aplicações Wails.

@end

## Por que usar roteamento baseado em hash?

O Wails incorpora seu frontend a uma janela de webview nativa. O uso de roteamento baseado em hash (#/page em vez de /page) evita conflitos com:

- O roteamento interno do runtime do Wails
- O processamento de URLs de janelas nativas em diferentes plataformas
- Os ativos de produção servidos a partir de caminhos que não estão na raiz

## Solução de problemas

Se você estiver enfrentando problemas com o roteamento no frontend em builds de produção:

- Verifique se a ferramenta de build do frontend está configurada para gerar a saída para roteamento no **modo hash**
- Em projetos baseados no Vite, adicione `base: "./"` ao seu `vite.config.js`
- Verifique se o seu `index.html` processa corretamente o fallback do roteamento de SPA
