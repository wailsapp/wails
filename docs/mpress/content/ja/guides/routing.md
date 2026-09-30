---
title: "フロントエンドルーティング"
description: "Wails アプリケーションでフロントエンドルーティングを使用する"
slug: "guides/routing"
sourcePath: "guides/routing.md"
---

フロントエンドルーティングは、シングルページアプリケーションでビューを切り替える一般的な方法です。 このガイドでは、Wails を使用する際に各種フロントエンドフレームワークで推奨される方法について説明します。

@tabs{sync-key="framework"}
[Vue]
Vue でのルーティングには、[ハッシュモード](https://next.router.vuejs.org/guide/essentials/history-mode.html#hash-mode)を推奨します：

```javascript
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    //...
  ],
});
```

ハッシュモードでは、URL ハッシュを使用して異なるビューをレンダリングします。ハッシュベースの URL 形式を使用することで、Wails ランタイムがルーティングに干渉する問題を回避できます。

[Angular]
Angular でのルーティングには、[HashLocationStrategy](https://codecraft.tv/courses/angular/routing/routing-strategies#%5Fhashlocationstrategy)を推奨します：

```typescript
RouterModule.forRoot(routes, { useHash: true });
```

ハッシュベースの URL を使用すると、すべてのプラットフォームで Wails のウィンドウ処理とルーティングが正しく連携します。

[React]
React でのルーティングには、[HashRouter](https://reactrouter.com/en/main/router-components/hash-router)を推奨します：

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

HashRouter はパスの代わりに URL ハッシュを使用するため、すべてのプラットフォームで Wails と安定して連携します。

[Svelte]
Svelte でのルーティングには、[svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router)を推奨します：

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

svelte-spa-router はハッシュベースのルーティングをサポートしているため、Wails アプリケーションと互換性があります。

@end

## ハッシュルーティングを使用する理由

Wails はフロントエンドをネイティブの WebView ウィンドウに埋め込みます。パスベースのルーティング（/page）ではなくハッシュベースのルーティング（#/page）を使用すると、次の要素との競合を回避できます：

- Wails ランタイムの内部ルーティング
- プラットフォームごとに異なるネイティブウィンドウの URL 処理
- ルート以外のパスから配信される本番環境用アセット

## トラブルシューティング

本番ビルドでフロントエンドルーティングに問題が発生する場合は、次を確認してください：

- フロントエンドのビルドツールが、<strong>ハッシュモード</strong>のルーティング向けに出力するよう設定されていることを確認します
- Vite ベースのプロジェクトでは、`vite.config.js`に`base: "./"`を追加します
- `index.html`が SPA ルーティングのフォールバックを適切に処理することを確認します
