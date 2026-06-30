---
title: Frontend Routing
description: Using frontend routing in your Wails application
sidebar:
  order: 35
---

import { Tabs, TabItem } from "@astrojs/starlight/components";

Frontend routing is a popular way to switch views in a single-page application.
This guide covers recommended approaches for different frontend frameworks when using Wails.


<Tabs syncKey="framework">
    <TabItem label="Vue" icon="seti:vue">

The recommended approach for routing in Vue is [Hash Mode](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode):

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

Hash mode uses the URL hash to render different views, avoiding issues with the Wails runtime
interfering with routing by using the hash-based URL format.

    </TabItem>
    <TabItem label="Angular" icon="seti:angular">

The recommended approach for routing in Angular is [HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy):

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

Using hash-based URLs ensures routing works correctly with Wails' window handling on all platforms.

    </TabItem>
    <TabItem label="React" icon="seti:react">

The recommended approach for routing in React is [HashRouter](https://reactrouter.com/en/main/router-components/hash-router):

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

HashRouter uses URL hashes instead of paths, which works reliably with Wails across all platforms.

    </TabItem>
    <TabItem label="Svelte" icon="seti:svelte">

The recommended approach for routing in Svelte is [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router):

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

svelte-spa-router supports hash-based routing, making it compatible with Wails applications.

    </TabItem>
</Tabs>

## Why Hash Routing?

Wails embeds your frontend into a native webview window. Using hash-based routing (#/page instead of /page)
avoids conflicts with:

- The Wails runtime's internal routing
- Native window URL handling on different platforms
- Production assets served from non-root paths

## Troubleshooting

If you're having issues with frontend routing in production builds:

- Make sure your frontend build tool is configured to output for **hash mode** routing
- For Vite-based projects, add `base: "./"` to your `vite.config.js`
- Verify your `index.html` handles the fallback properly for SPA routing