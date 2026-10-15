---
title: "カスタムテンプレートの作成"
description: "独自の Wails v3 プロジェクトテンプレートを生成、カスタマイズ、ホスティングする方法"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails には組み込みテンプレート一式が付属していますが、独自のテンプレートを作成してコミュニティと共有することもできます。カスタムテンプレートは単なる Git リポジトリです。公開すれば、誰でも 1 つのコマンドでそのテンプレートからプロジェクトのひな形を生成できます。

## テンプレートのひな形を生成する

`wails3 generate template` コマンドを実行すると、すぐにカスタマイズできるテンプレートディレクトリが生成されます。

```bash
wails3 generate template -name MyTemplate
```

すべてのフラグ：

| フラグ | 説明 | デフォルト |
| --- | --- | --- |
| `-name` | テンプレート名（必須） | — |
| `-author` | 作成者名 | — |
| `-description` | CLI に表示される短い説明 | — |
| `-helpurl` | このテンプレートのドキュメントの URL | — |
| `-version` | 初期バージョン | `v0.0.1` |
| `-frontend` | 既存のフロントエンドディレクトリをテンプレートにコピーする | — |
| `-dir` | テンプレートディレクトリの出力先 | 現在のディレクトリ |

すべてのフラグを指定する例：

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

生成されるディレクトリの構成は次のとおりです。

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="NEXTSTEPS.md を読む"}
生成された `NEXTSTEPS.md` には、テンプレートの各部分に関する詳しいガイダンスが記載されています。カスタマイズする前に読んでください。公開前に削除してください。このファイルがテンプレートから作成されたプロジェクトに含まれていてはいけません。

@end

## テンプレートのメタデータを設定する

`template.yaml` を開き、テンプレートのメタデータを設定します。

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

`wailsVersion` フィールドは<strong>必須</strong>であり、`3`でなければなりません。先頭の `# yaml-language-server` コメントにより、VS Code（[YAML 拡張機能](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)を使用）および JetBrains IDE で自動補完とインライン検証が有効になります。このコメントは残しても削除しても構いません。実行時の動作には影響しません。

## テンプレートをカスタマイズする

### フロントエンド

`frontend/` ディレクトリは、テンプレートから作成されるすべてのプロジェクトにそのままコピーされます。プレースホルダーの内容を実際のフロントエンドに置き換えてください。

@tabs
[ゼロから作成する]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

画面の指示に従い、依存関係をインストールします。

```bash
npm install
```

[既存のプロジェクトを使用する]
テンプレートの生成時に `-frontend` を指定すると、既存のフロントエンドを一度にコピーできます。

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

または、後から `frontend/` ディレクトリに手動でコピーします。

@end

### ビルドタスク

`Taskfile.tmpl.yml` ではビルドワークフローを定義します。フロントエンドのツールチェーンに合わせて、`install:frontend:deps` タスクと `build:frontend` タスクを更新してください。

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go アプリケーション

`main.go.tmpl` ファイルはアプリケーションのエントリーポイントです。プロジェクトの作成時に Wails のテンプレートエンジンによって処理され、`{{.ProductName}}` などのテンプレート変数がユーザーの指定値に置き換えられます。

IDE のサポートを利用して通常の Go ファイルとして編集するには、一時的にファイル名を `main.go` に変更して編集し、コミットする前に `main.go.tmpl` に戻してください。

#### テンプレート変数

次の変数は、どの `.tmpl` ファイルでも使用できます。

| 変数 | 説明 | 例 |
| --- | --- | --- |
| `{{.ProjectName}}` | ユーザーが指定したプロジェクト名 | `"MyApp"` |
| `{{.BinaryName}}` | バイナリファイル名 | `"myapp"` |
| `{{.ProductName}}` | 製品の表示名 | `"My Application"` |
| `{{.ProductDescription}}` | 製品の説明 | `"An awesome application"` |
| `{{.ProductVersion}}` | 製品バージョン | `"1.0.0"` |
| `{{.ProductCompany}}` | 会社名／作成者名 | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | 著作権表示文字列 | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | 製品に関する追加コメント | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | ドメイン名を逆順にした形式の製品識別子 | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Goモジュールパス | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | プロジェクトの作成に使用したWailsのバージョン | `"3.0.0"` |
| `{{.Typescript}}` | テンプレート名が`-ts`で終わる場合は`true` | `true` |
| `{{.Opn}}` | リテラルの`{{` — テンプレート内でエスケープ | `{{` |
| `{{.Cls}}` | リテラルの`}}` — テンプレート内でエスケープ | `}}` |

@note{type="tip"}
HTML、JSON、YAMLファイルを含め、テンプレート内のどのファイルも`.tmpl`ファイルにできます。`.tmpl`サフィックスがないファイルは、そのままコピーされます。

@end

## テンプレートをローカルでテストする

公開する前に、ローカルパスからプロジェクトを作成してテンプレートをテストします。

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

次に、プロジェクトが動作することを確認します。

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

次の点を確認してください。

- フロントエンドのホットリロードが動作すること
- Goコードを変更すると、アプリが再ビルドされて再起動すること
- `bin/`内の本番用バイナリが正しく動作すること

## GitHubで公開する

@steps
### **テンプレート用の公開GitHubリポジトリを作成します**。リポジトリのルートには`template.yaml`が必要です。
### **`NEXTSTEPS.md`を削除します** — このファイルはテンプレート作成者向けのガイドであり、ユーザーがテンプレートから作成するプロジェクトに含めてはいけません。
### **コミットしてプッシュします**。テンプレートディレクトリの内容をリポジトリのルートに配置してください。
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **セマンティックバージョニングを使用してリリースにタグを付けます**。
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

これで、ユーザーはこのテンプレートからプロジェクトを作成できます。

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="サードパーティ製テンプレートに関する警告"}
ユーザーがリモートテンプレートをインストールすると、Wailsは、そのテンプレートがサードパーティ製のコードであり、その内容についてWailsプロジェクトは一切の責任を負わないことを説明する警告を表示します。プロジェクトを作成するには、ユーザーが明示的に確認する必要があります。

テンプレート作成者は、テンプレート内のすべてのコードのセキュリティと正確性に責任を負います。

@end

## ベストプラクティス

- <strong>わかりやすい`README.md`</strong>を作成します — これは、ユーザーがプロジェクトを作成した後に表示されます。プロジェクトの実行、ビルド、カスタマイズの方法を説明してください。
- <strong>`helpurl`</strong>を入力します — リポジトリまたは専用ドキュメントへのリンクを指定してください。ユーザーには、Wails CLIのテンプレート一覧で表示されます。
- <strong>フロントエンドの依存関係のバージョンを`package.json`で固定</strong>し、アップストリームの更新によってインストールが失敗するのを防ぎます。
- **タグを付ける前にテストします** — コミュニティに告知する前に、タグ付きリリースから新規プロジェクトを作成してください。
- <strong>`wailsVersion: 3`</strong>は変更せずに維持します — このフィールドは、テンプレートが対象とするWailsのメジャーバージョンをWailsに示します。変更しないでください。
- **定期的に更新します** — 依存関係を最新の状態に保ち、Wailsの新しいリリースに対してテストしてください。
