---
title: "テンプレートシステム"
description: "Wails v3 が新規プロジェクトをスキャフォールディングする仕組み、テンプレートの構成、独自テンプレートの作成方法について説明します。"
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails には、`wails3 init`でそのまま実行できるプロジェクトを生成できる<strong>テンプレートシステム</strong>が付属しています。組み込みテンプレートが用意されているフレームワークは、意図的に少数（Vanilla、React、Vue、Svelte）に絞られています。その他のフレームワークも、[独自のフロントエンドを持ち込む](/guides/dev/frontend-frameworks/)か、[カスタムテンプレート](/guides/advanced/custom-templates/)を公開することで使用できます。

このページでは、以下について説明します。

1. テンプレートディレクトリの構成
2. CLI がテンプレートを選択してレンダリングする仕組み
3. 新しいテンプレートを作成する手順
4. 既存テンプレートの更新または上書き
5. トラブルシューティングとベストプラクティス

---

## 1. テンプレートの格納場所

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — すべてのプロジェクトにマージされる共通のボイラープレート（Taskfile、`build/`ディレクトリ、共有インフラストラクチャ）。
- **`base/`** — すべてのテンプレートの基盤となる Go 側のコード。注意：`base/`自体には`template.json`が<strong>含まれていません</strong>。このファイルは、フレームワーク固有の各テンプレート内にあります。
- **フレームワークフォルダー** — フロントエンド（`frontend/`）、フレームワークの設定、およびテンプレートのメタデータを記述する`template.json`が含まれます。
- フォルダー名は、CLI に渡す<strong>テンプレート ID</strong>（`wails3 init -t react`）と一致します。
- <strong>言語の規則：</strong>TypeScript がデフォルトで、接尾辞のない名前（`react`）を使用します。JavaScript 版が存在する場合は、`-js`という接尾辞を付けます（`react-js`）。組み込みテンプレートでは、`template.yaml`内の`typescript: true|false`で言語を明示的に宣言します。コミュニティテンプレートでは、従来の`-ts`接尾辞が引き続き使用されている場合があり、フォールバックとして認識されます。

> `internal/templates/`ディレクトリ全体が CLI バイナリにコンパイルされます
>
> その際に`//go:embed *`を使用するため、ユーザーはオフラインでもプロジェクトをスキャフォールディングできます。

---

## 2. `wails3 init`がテンプレートを使用する仕組み

呼び出しチェーン（`cmd/wails3/init.go`はありません。CLI は`cmd/wails3/main.go`内で直接構成されています）：

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

`Template.Load()` / `Template.CopyTo()` / `Template.Validate()` API はありません。展開は、埋め込まれた`fs.FS`に対して`gosod`（`github.com/leaanthony/gosod`）が処理します。

### `wails3 init`のフラグ

`internal/flags/init.go`で定義されています。

| フラグ | 用途 | デフォルト |
| --- | --- | --- |
| `-p` | パッケージ名 | `main` |
| `-t` | 組み込みテンプレート名、ローカルパス、または URL | `vanilla` |
| `-n` | プロジェクト名 | （空） |
| `-d` | プロジェクトディレクトリ | `.` |
| `-q` | コンソール出力を抑制 | false |
| `-l` | テンプレートを一覧表示 | false |
| `-skipgomodtidy` | 展開後の`go mod tidy`の実行をスキップ | false |
| `-git` | 初期化する Git リポジトリの URL | （空） |
| `-mod` | Go モジュールパス（未設定の場合は`-git`から導出） | （空） |
| `-s` | リモートテンプレート使用時の警告をスキップ | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | 生成されるビルドアセットに組み込まれるメタデータ | 適切なデフォルト値 |

`-list`の長いエイリアスは<strong>なく</strong>（`-l`のみ）、テンプレート単位の`--help`もありません。

### 置換

プレースホルダーは標準の Go テンプレートディレクティブです。先頭の`.`はフィールドアクセサーの一部です。

| プレースホルダー | 例 | 値の取得元 |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n` フラグ／ディレクトリ名 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod` フラグ、または `-git` から導出 |
| `{{.WailsVersion}}` | `v3.0.0-…` | `internal/version` から取得したコンパイル時組み込み定数 |
| `{{.ProductName}}`、`{{.ProductDescription}}`、`{{.ProductVersion}}`、`{{.ProductCompany}}`、`{{.ProductCopyright}}`、`{{.ProductComments}}`、`{{.ProductIdentifier}}` | ビルド時のメタデータ | 対応する `-product*` フラグ |

新しいプレースホルダーが必要な場合は、`internal/templates/templates.go` のテンプレートデータにフィールドを追加し、`internal/flags/init.go` に対応するフィールドとフラグを追加します（または `internal/commands/init.go` から値を設定します）。

### コピー後フック

`gosod` によるテンプレートの展開が完了すると、CLI は次の処理を実行します。

```
go mod tidy
```

ただし、`-skipgomodtidy` を渡した場合は実行しません。`task deps` ステップはありません。

---

## 3. 新しいテンプレートの作成

> 例：**Solid** テンプレートを追加する

### 3.1 フォルダーと ID

```
internal/templates/solid/
```

フォルダー名がテンプレート ID になります。**kebab-case** を使用してください。

### 3.2 最小限のファイル構成

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

まず `react` をコピーし、不要なファイルを削除します。`template.yaml` の作成を忘れないでください。TypeScript テンプレートでは `typescript: true` を設定します。これがないフォルダーは `base/` だけです。

### 3.3 プレースホルダーの更新

リテラルのサンプル値を検索し、Go テンプレートディレクティブに置換します。例：

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 組み込み

`templates.go` は初期化時に埋め込みファイルシステムを走査するため、通常は `internal/templates/<id>/` の下に新しいフォルダーを追加するだけで十分です。手動の登録呼び出しは必要ありません。追加のロジック（カスタム検証やコピー後の処理）が必要な場合は、`internal/templates/templates.go` 内の `templates.Install` に追加します。

### 3.5 テスト

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

次の点を確認します。

- 開発サーバーが `WAILS_VITE_PORT` で公開されているポートで起動すること
- 生成されたバインディングが `frontend/bindings/...` の下に作成されること
- ホットリロードが動作すること

---

## 4. 既存テンプレートの変更

1. `internal/templates/<id>/` の下にあるファイルを編集します。
2. CLI を再ビルドします（`cd v3 && go build -o ../wails3 ./cmd/wails3`）。`//go:embed *` ディレクティブによって新しい内容が取り込まれます。
3. `frontend/package.json` と `Taskfile.yml` にある <strong>依存関係のバージョン</strong>を上げます。
4. 動作が変わる場合は、テンプレートの `template.json` の説明を更新します。

### 一般的な変更

| 作業 | 変更箇所 |
| --- | --- |
| 開発サーバーのポートを変更する | `frontend/vite.config.ts` — `WAILS_VITE_PORT` を読み取る |
| 環境変数を追加する | `build/Taskfile.yml` または `frontend/.env` |
| JavaScript パッケージマネージャーを置き換える | `build/Taskfile.yml` 内の `npm` を `pnpm`/`bun` に置き換える |

---

## 5. テンプレート作成のヒント

- **フロントエンドを汎用的に保つ** — Wails 固有のグローバルを参照しないでください。`/wails/runtime.js` は実行時にアセットサーバーから配信されます。
- **コンパイル済み成果物を含めない** — `node_modules`、`dist`、`.DS_Store` を埋め込みディレクトリから除外します（または、決してコミットされないように `.gitignore` します）。
- **前提条件を記載する** — Node のバージョンや追加の CLI ツールなどを、`template.json` または `NEXTSTEPS.md` に記載します。
- **破壊的変更を避ける** — 大規模な全面改修の場合は、既存のテンプレートを変更するのではなく、新しいテンプレート ID を作成します。

---

## 6. トラブルシューティング

| 症状 | 原因 | 解決方法 |
| --- | --- | --- |
| `unknown template name` | `-t` の入力ミス、またはテンプレートが埋め込まれていない | `wails3 init -l` を実行して、利用可能なテンプレートを一覧表示する |
| プレースホルダーが置換されない | `{{.ProjectName}}` ではなく `{{ProjectName}}` を使用している | 先頭に `.` を追加する（Go テンプレートのフィールドアクセス） |
| 開発サーバーで空白ページが表示される | Vite の設定で `WAILS_VITE_PORT` が読み込まれていない | `vite.config.ts` を確認する |
| 本番環境向けのフロントエンドのビルドに失敗する | Vite の `base` パスの設定を忘れている | `vite.config.ts` で `base: "./"` を設定する |

---

## 7. 主要ソースファイル一覧

| ファイル | 役割 |
| --- | --- |
| `internal/templates/templates.go` | テンプレートのファイルシステムを埋め込み、`Install(options *flags.Init) error`、`GetDefaultTemplates()`、`ValidTemplateName(name)` を公開する |
| `internal/templates/<id>/**` | 実際のテンプレート内容 |
| `internal/commands/init.go` | CLI の連携処理：テンプレートを選択し、メタデータを設定して、`templates.Install` を呼び出す |
| `internal/commands/generate_template.go` | `wails3 generate template` — 実際に開発・使用しているプロジェクトをテンプレートへ<em>エクスポート</em>し直すためのユーティリティ（更新時に便利） |
| `internal/flags/init.go` | `wails3 init` のフラグ定義 |

---

## 8. まとめ

- テンプレートは **`internal/templates/`** にあり、`//go:embed *` を介して CLI に組み込まれます。
- `wails3 init -t <id>` は `gosod` を介してテンプレートを展開し、`go mod tidy` を実行します（`-skipgomodtidy` でスキップできます）。
- テンプレートの作成は簡単です。<strong>フォルダーを作成</strong>し、ファイルと `template.json` を追加して、`{{.ProjectName}}` 形式のプレースホルダーを使用するだけです。
- このシステムは<strong>拡張可能</strong>かつ<strong>自己完結型</strong>であり、カスタムスタックをチームやコミュニティと共有するのに最適です。

テンプレート作成を楽しんでください！
