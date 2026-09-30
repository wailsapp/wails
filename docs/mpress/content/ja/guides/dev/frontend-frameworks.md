---
title: "その他のフロントエンドフレームワークを使用する"
description: "独自の Vite プロジェクトを frontend ディレクトリに配置して、組み込みテンプレートがないフレームワークを使用する方法"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails には、意図的に絞り込まれた次のフレームワーク向けスターターテンプレートが組み込まれています。

| テンプレート | 言語 |
| --- | --- |
| `vanilla` | TypeScript（デフォルト） |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

ただし、選択肢がこれらに限られるわけではありません。Wails アプリのフロントエンドは<strong>単なる Web プロジェクト</strong>です。静的な HTML/CSS/JS にビルドできるものであれば、何でも動作します。使用したいフレームワーク（Solid、Preact、Lit、Qwik、SvelteKit、Angular など）にテンプレートがなくても、数分で自分で雛形を生成できます。

## `frontend/` ディレクトリの用途

Wails は、`frontend/` 内でどのフレームワークが使われているかを問いません。依存するのは、フレームワークに依存しない次の簡単な規約だけです。

- **`frontend/dist/` が配布物になります。** `main.go` は、ビルド済みのフロントエンドを `//go:embed all:frontend/dist` で埋め込み、アセットサーバーから配信します。ビルドでは、静的バンドルを `frontend/dist/`（Vite のデフォルト出力ディレクトリ）に出力する必要があります。
- **ビルドは `frontend/package.json` によって実行されます。** `wails3 build` の実行中、Wails はフロントエンドの `build` スクリプトを実行します。`wails3 dev` の実行中は `dev` を実行し、ホットリロードのために Vite 開発サーバーをプロキシします。
- **バインディングは `frontend/bindings/` に生成されます。** Wails は登録済みの Go サービスを調べ、そこに型安全な SDK を書き込みます。ほかのモジュールと同じようにインポートできます。
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **開発サーバーは固定ポートで動作します。** `wails3 dev` は `WAILS_VITE_PORT` で指定されたポート（デフォルトは `9245`）上の Vite をプロキシするため、`strictPort: true` を使用して `server.port` をそのポートに設定してください。これは通常の Vite 設定であり、Wails プラグインは関与しません。
- **省略可能 — 型付きカスタムイベント。** 組み込みテンプレートでは、`@wailsio/runtime/plugins/vite` プラグインも登録されます。このプラグインが必要なのは、<em>型付き</em>カスタムイベントを使用する場合だけです。生成されたイベント型定義をランタイムに注入し、バインディングが生成されるまではビルドを失敗させます。文字列ベースの `Events.On("time", …)` API しか使用しない場合は、省略できます。

コンポーネント、ルーティング、状態、スタイルなど、その他のすべては使用するフレームワークに完全に委ねられます。

## Vite で任意のフレームワークの雛形を生成する

最も手早い方法は、組み込みテンプレートから始め（これにより `main.go`、`Taskfile`、ビルドアセット、動作する Go サービスが用意されます）、その後 `frontend/` を使用するフレームワーク向けの新しい Vite プロジェクトに置き換えることです。

@steps
### デフォルトテンプレートからプロジェクトを作成する
```bash
wails3 init -n myapp
cd myapp
```

### `frontend/` を使用するフレームワーク向けの Vite アプリに置き換える
Vite では、ほとんどのフレームワークの雛形を 1 つのコマンドで生成できます。テンプレートを選択してください。

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

`solid` は、任意の Vite テンプレートに置き換えられます。`preact`、`lit`、`svelte`、`vue`、`react`、`vanilla`、またはそれらの `-ts` バリアント（`solid-ts`、`preact-ts` など）を使用できます。

### ランタイムをインストールし、Vite に Wails 開発サーバーを指定する
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` は、JS API（`Events`、`Browser`、ダイアログなど）を提供します。`vite.config` で<em>必須</em>の変更は開発サーバーのポートだけです。`wails3 dev` が検出できるように設定してください。

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

<strong>型付きカスタムイベント</strong>を使用する予定がある場合に限り、プラグインも追加してください。このプラグインは生成されたイベント型を注入し、ビルド前にバインディングが存在することを必須とします。

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Go サービスを呼び出す
バインディングを一度生成した後、コンポーネント内の任意の場所でインポートしてください。

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### 実行する
```bash
wails3 dev
```

@end

@note{type="tip" title="プレーン JavaScript プロジェクトを生成する"}
同じコマンドで TypeScript を使用しないプロジェクトも生成できます。`-ts` ではない Vite テンプレートを使用するだけです。

```bash
npm create vite@latest frontend -- --template solid
```

@end

## 独自の雛形生成ツールを持つフレームワーク

一部のフレームワークは Vite の `create` テンプレートでは作成せず、独自のツールを使用します。そのようなフレームワークも動作します。各フレームワークの標準コマンドで雛形を生成し、その後 Wails プラグインを追加してください。

- **SvelteKit：** `npx sv create frontend`。静的バンドルにビルドされるよう、静的アダプター（`@sveltejs/adapter-static`）を使用し、SSR を無効にしてください。
- **Qwik：** `npm create qwik@latest`。静的（SSG）アダプターを使用してください。
- **Angular：** `ng new` で雛形を生成し、`outputPath` を `dist` に設定して、ビルドスクリプトに `ng build` を指定してください。

規則は常に同じです。`frontend/dist/` に静的ビルドを生成し、`@wailsio/runtime` Vite プラグインを残す（またはランタイムを直接インポートする）とともに、Go バインディングを `frontend/bindings/` からインポートしてください。

@note{type="info"}
フレームワーク向けに完成度の高いセットアップを作成した場合は、ほかの利用者が直接 `wails3 init -t` できるよう、[カスタムテンプレート](/guides/advanced/custom-templates/)として公開することを検討してください。

@end
