---
title: "テンプレートシステム"
description: "Wails v3 が新規プロジェクトをスキャフォールディングする仕組み、テンプレートの構成、および独自テンプレートの作成方法。"
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails には、`wails3 init` で実行可能なプロジェクトを生成できる**テンプレートシステム**が付属しています。意図的に絞り込まれた少数のフレームワーク（Vanilla、React、Vue、Svelte）には組み込みテンプレートがあります。それ以外のフレームワークも、[独自のフロントエンドを持ち込む](/guides/dev/frontend-frameworks/)か、[カスタムテンプレート](/guides/advanced/custom-templates/)を公開することで使用できます。

このページでは、以下について説明します。

1. テンプレートディレクトリの構成
2. CLI がテンプレートを選択してレンダリングする仕組み
3. 新しいテンプレートを段階的に作成する方法
4. 既存のテンプレートを更新または上書きする方法
5. トラブルシューティングとベストプラクティス

---

## 1. テンプレートの配置場所

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — すべてのプロジェクトにマージされる共通ボイラープレート（Taskfile、`build/` ディレクトリ、共有基盤）。
- **`base/`** — カスタムテンプレート用の最小限のフロントエンドスケルトンであり、自動的に追加されるバックエンド層ではありません。通常のプロジェクトの作業コピー作成では、`_common/` に選択したフレームワークテンプレートが重ね合わされます。
- **フレームワークフォルダー** — フロントエンド（`frontend/`）、フレームワーク設定、およびテンプレートのメタデータを記述する `template.yaml` が含まれます。
- フォルダー名は、CLI に渡す**テンプレート ID**（`wails3 init -t react`）と一致します。
- **言語の命名規則：** TypeScript がデフォルトで、サフィックスのない名前（`react`）を使用します。JavaScript バリアントが存在する場合は、`-js` サフィックスを付けます（`react-js`）。組み込みテンプレートは、`template.yaml` 内の `typescript: true|false` で言語を明示的に宣言します。コミュニティテンプレートでは、従来の `-ts` サフィックスを引き続き使用でき、フォールバックとして認識されます。

> `internal/templates/` ディレクトリ全体が CLI バイナリにコンパイルされます
>
> これは `//go:embed *` を介して行われるため、ユーザーはオフラインでもプロジェクトをスキャフォールディングできます。

---

ローダーは従来の `template.json` メタデータにも対応しています。新しい組み込みテンプレートでは YAML を使用し、`typescript` 設定も明示します。

## 2. `wails3 init` がテンプレートを使用する仕組み

呼び出しチェーン（`cmd/wails3/init.go` はありません。CLI は `cmd/wails3/main.go` で直接接続されています）：

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

`Template.Load()` / `Template.CopyTo()` / `Template.Validate()` API はありません。作業コピーの作成は、埋め込まれた `fs.FS` に対して `gosod`（`github.com/leaanthony/gosod`）が処理します。

### `wails3 init` のフラグ

`internal/flags/init.go` で定義されています：

| フラグ | 用途 | デフォルト |
| --- | --- | --- |
| `-p` | パッケージ名 | `main` |
| `-t` | 組み込みテンプレート名、ローカルパス、または URL | `vanilla` |
| `-n` | プロジェクト名 | （空） |
| `-d` | プロジェクトディレクトリ | `.` |
| `-q` | コンソール出力を抑制 | false |
| `-l` | テンプレートを一覧表示 | false |
| `-skipgomodtidy` | 作業コピーの作成後に `go mod tidy` を実行しない | false |
| `-git` | 初期化する Git リポジトリの URL | （空） |
| `-mod` | Go モジュールパス（未設定の場合は `-git` から導出） | （空） |
| `-s` | リモートテンプレート使用時の警告を省略 | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | 生成されるビルドアセットに埋め込まれるメタデータ | 妥当なデフォルト値 |

`-list` には長い形式のエイリアスはなく（`-l` のみ）、テンプレートごとの `--help` もありません。

### 置換

プレースホルダーは標準の Go テンプレートディレクティブです。先頭の `.` はフィールドアクセサーの一部です：

| プレースホルダー | 例 | ソース |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n` フラグ / ディレクトリ名 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod` フラグ、または `-git` から導出 |
| `{{.WailsVersion}}` | `v3.0.0` | `internal/version` から組み込まれるコンパイル時定数 |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | ビルド時のメタデータ | 対応する `-product*` フラグ |

新しいプレースホルダーが必要な場合は、`internal/templates/templates.go` のテンプレートデータにフィールドを追加し、`internal/flags/init.go` に対応するフィールドとフラグを追加します（または `internal/commands/init.go` から設定します）。

### コピー後のフック

`gosod` によるテンプレートからの作業コピー作成が完了すると、CLI は次を実行します：

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

フォルダー名がテンプレート ID になります。**kebab-case** にしてください。

### 3.2 最小限のファイルセット

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

まず `react` をコピーし、不要なファイルを削除します。`template.yaml` の作成も忘れないでください。TypeScript テンプレートの場合は `typescript: true` を設定します。これがないフォルダーは `base/` だけです。

### 3.3 プレースホルダーを更新する

次のように、リテラルのサンプル値を検索し、Go テンプレートディレクティブに置換します。

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 組み込む

`templates.go` は初期化時に埋め込み FS を走査するため、通常は `internal/templates/<id>/` の下に新しいフォルダーを追加するだけで十分です。手動で登録関数を呼び出す必要はありません。カスタム検証やコピー後の処理など、追加のロジックが必要な場合は、`internal/templates/templates.go` の `templates.Install` に追加します。

### 3.5 テストする

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

次の点を確認します。

- 開発サーバーが `WAILS_VITE_PORT` で公開されているポートで起動する
- 生成されたバインディングが `frontend/bindings/...` の下に作成される
- ホットリロードが動作する

---

## 4. 既存のテンプレートを変更する

1. `internal/templates/<id>/` の下にあるファイルを編集します。
2. CLI を再ビルドします（`cd v3 && go build -o ../wails3 ./cmd/wails3`）。`//go:embed *` ディレクティブによって新しい内容が取り込まれます。
3. `frontend/package.json` と `Taskfile.yml` にある **依存関係のバージョン**を更新します。
4. 動作が変わる場合は、テンプレートの `template.yaml` の説明を更新します。

### 一般的な調整

| 作業 | 変更箇所 |
| --- | --- |
| 開発サーバーのポートを変更する | `frontend/vite.config.ts` — `WAILS_VITE_PORT` を読み取る |
| 環境変数を追加する | `build/Taskfile.yml` または `frontend/.env` |
| JS パッケージマネージャーを置き換える | `build/Taskfile.yml` で `npm` を `pnpm`/`bun` に置き換える |

---

## 5. テンプレート作成のヒント

- **フロントエンドを汎用的に保つ** — Wails 固有のグローバルを参照しないでください。実行時には `/wails/runtime.js` がアセットサーバーから配信されます。
- **コンパイル済み成果物を含めない** — `node_modules`、`dist`、`.DS_Store` を埋め込みディレクトリから除外します（または、決してコミットされないように `.gitignore` します）。
- **前提条件を文書化する** — Node のバージョンや追加の CLI ツールなどを `template.yaml` または `NEXTSTEPS.md` に記載します。
- **破壊的変更を避ける** — 大幅に刷新する場合は、既存のテンプレートを変更するのではなく、新しいテンプレート ID を作成します。

---

## 6. トラブルシューティング

| 症状 | 原因 | 解決方法 |
| --- | --- | --- |
| `unknown template name` | `-t` の入力ミス、またはテンプレートが埋め込まれていない | `wails3 init -l` を実行し、利用可能なテンプレートを一覧表示する |
| プレースホルダーが置換されない | `{{.ProjectName}}` ではなく `{{ProjectName}}` を使用している | 先頭に `.` を追加する（Go テンプレートのフィールドアクセス） |
| 開発サーバーで空白ページが開く | Vite の設定で `WAILS_VITE_PORT` が読み取られていない | `vite.config.ts` を確認する |
| 本番環境向けのフロントエンドのビルドに失敗する | Vite の `base` パスを設定し忘れている | `vite.config.ts` で `base: "./"` を設定する |

---

## 7. 主要なソースファイルの対応表

| ファイル | 役割 |
| --- | --- |
| `internal/templates/templates.go` | テンプレート FS を埋め込み、`Install(options *flags.Init) error`、`GetDefaultTemplates()`、`ValidTemplateName(name)` を公開する |
| `internal/templates/<id>/**` | 実際のテンプレートの内容 |
| `internal/commands/init.go` | CLI の橋渡し処理：テンプレートを選択し、メタデータを設定して、`templates.Install` を呼び出す |
| `internal/commands/generate_template.go` | `wails3 generate template` — 稼働中のプロジェクトをテンプレートとして*エクスポート*し直すためのユーティリティ（更新時に便利） |
| `internal/flags/init.go` | `wails3 init` のフラグ定義 |

---

## 8. まとめ

- テンプレートは **`internal/templates/`** にあり、`//go:embed *` によって CLI に組み込まれます。
- `wails3 init -t <id>` は `gosod` を介してテンプレートを展開し、`go mod tidy` を実行します（`-skipgomodtidy` を指定すると省略できます）。
- テンプレートの作成は、**フォルダーを作成**し、ファイルと `template.yaml` を追加して、`{{.ProjectName}}` 形式のプレースホルダーを使用するだけです。
- このシステムは**拡張可能**であり、**自己完結型**です。カスタムスタックをチームやコミュニティと共有するのに最適です。

テンプレート作成を楽しんでください！
