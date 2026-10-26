---
title: "M-Press でドキュメントを作成する"
description: "M-Press をダウンロードしてインストールし、Wails ドキュメントをプレビューして、ドキュメントへのコントリビューションを提出します。"
sourcePath: "contributing/documentation.md"
---

Wails v3 では、Markdown からこのドキュメントサイトを生成するために **M-Press** を使用しています。誤字の修正、説明の改善、サンプルやガイドの追加、翻訳の修正が可能です。ドキュメントのみの修正では、事前の Issue や失敗するコードテストは必要ありません。

コントリビューションの準備には Git、プルリクエストの提出には GitHub アカウントが必要です。ダウンロード可能な M-Press 実行ファイルには、ドキュメントジェネレーターと図のレンダラーが含まれています。ドキュメントのプレビューに、Go、Node.js、Wails CLI、デスクトップビルド用の依存関係は必要ありません。リポジトリの追加検証スクリプトには Python 3 が必要です。

## ドキュメントページから開始する

このガイドを開くには、サイドバーで **コントリビューション → ドキュメント** を選択します。特定のページで作業するには、そのページのヘッダーにある **編集アイコン**を選択し、続いて **このページをローカルで編集** を選択します。お使いのオペレーティングシステム用のコマンドをコピーし、ターミナルで実行してください。

M-Press がインストールされている場合、コマンドはそれを使用します。インストールされていない場合は、インストーラーが公式リリースをダウンロードし、チェックサムを検証します。このコマンドはコントリビューション用のローカル作業コピーを準備し、選択したページを M-Press で開きます。自動インストーラーはより新しいリリースをダウンロードする可能性があるため、リポジトリのチェックを実行する前に、`mpress version` を以下の固定バージョンと照合してください。

ローカルプレビューを再ビルドさせながら、お使いのエディターでソースを編集できます。ウィザードでは、**変更をレビュー**、**変更をコミット**、**コントリビューションブランチをプッシュ**、**ドラフトプルリクエストを開く**も、それぞれ個別の操作として利用できます。提出する前に、以下のチェックとプルリクエストに関するガイダンスに従ってください。翻訳の生成は任意であり、設定済みのプロバイダーが必要になる場合があります。Markdown の編集とプレビューには必要ありません。

再現可能な手動セットアップを行うには、以下の手順に従ってください。

## M-Press をダウンロードしてインストールする

リポジトリのビルドスクリプトは、**最新の安定版 M-Press リリース**をダウンロードします。ローカルでも同じリリースを使用してください。以下からインストール方法を選び、`mpress version`を実行して、リリースページに表示されているバージョンと比較します。

### 実行ファイルをダウンロードする

[最新の安定版 M-Press リリース](https://github.com/leaanthony/mpress/releases/latest)を開くか、以下から対応するアーカイブをダウンロードしてください。ARM64 には Apple Silicon が含まれます。AMD64 は、一般的な Intel/AMD PC のアーキテクチャです。

| コンピューター | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

同じリリースから[checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt)をダウンロードしてください。ダウンロード先のディレクトリでアーカイブの SHA-256 ハッシュを計算し、`checksums.txt`内の該当ファイル名の行と比較します。

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

別のアーキテクチャを選択した場合は、実際にダウンロードしたファイル名を使用してください。ハッシュが一致したら、実行ファイルを展開してインストールします。

@tabs
[macOS]
ダウンロードディレクトリで、Apple Silicon 用アーカイブを使用します。

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Intel Mac の場合は、`mpress-darwin-amd64.tar.gz` に置き換えてください。新しいターミナルでもコマンドを使用できるようにするには、`export PATH` の行を `~/.zshrc` に追加します。

[Linux]
ダウンロードディレクトリで実行します。

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

ARM64 の場合は、`mpress-linux-arm64.tar.gz` に置き換えてください。新しいターミナルでもコマンドを使用できるようにするには、`export PATH` の行を、`~/.bashrc` などのシェル起動ファイルに追加します。

[Windows PowerShell]
ダウンロードディレクトリで実行します。

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

ARM64 の場合は、`mpress-windows-arm64.zip` に置き換えてください。以降のセッションでもコマンドを使用できるようにするには、Windows の **環境変数**で、ユーザーの **Path** にユーザーの `bin` ディレクトリを追加し、ターミナルを開き直します。
@end

### Go でインストールする

すでに Go を使用している場合は、代わりに最新の安定版リリースをビルドできます。必要なツールチェーンについては、そのリリースの`go.mod`を確認してください。Go のツールチェーン自動選択機能によって、そのツールチェーンがダウンロードされる場合があります。このドキュメント更新で使用したリリースには、Go 1.26.10以降が必要です。

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go は、コマンドを `GOBIN` にインストールします。`GOBIN` が設定されていない場合は、`go env GOPATH` 配下の `bin` ディレクトリにインストールします。`mpress` が見つからない場合は、そのディレクトリを `PATH` に追加し、ターミナルを開き直してください。

## ドキュメントのソースを取得する

[wailsapp/wails をフォーク](https://github.com/wailsapp/wails/fork)し、`YOUR_USERNAME` を自分の GitHub ユーザー名に置き換えて、次のコマンドを実行してください。

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

これにより、アップストリームの **master** ブランチを起点として、ドキュメント作業用ブランチが作成されます。`docs/my-change` の代わりに、内容が分かる独自のブランチ名を使用してください。

## ライブプレビューを実行する

`mpress.yaml` がある Wails リポジトリのルートから実行します。

```sh
mpress dev
```

ブラウザーで [http://127.0.0.1:3000](http://127.0.0.1:3000) を開きます。このガイドを表示するには、[http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/) にアクセスしてください。編集中はコマンドを実行したままにします。ソースへの変更を保存すると、プレビューが再ビルドされます。以下の検証コマンドを実行する前に、**Ctrl+C** を押して停止してください。

リポジトリには M-Press の設定がすでに含まれているため、新しいサイトを初期化する必要はありません。別のポートを使用するには、`mpress dev --port 3001` を実行し、対応する localhost アドレスを開いてください。

## ページを編集または追加する

英語のソースページは `docs/mpress/content/` にあります。翻訳は、`fr/` や `id/` などの言語フォルダーにあります。編集アイコンから表示されるコントリビューション手順で、現在のページに対応するソースファイルを確認できます。

既存のページでは、Markdown を編集し、先頭の YAML メタデータを維持してください。見出し、段落、リスト、リンク、フェンス形式のコードブロックを使用します。`@note` / `@end`や`@tabs` / `@end`など、対になっている M-Press コンポーネントは維持してください。サポートされる構文については、[M-Press オーサリングドキュメント](https://m-press.me/authoring/)を参照してください。

新しいガイドでは、`docs/mpress/content/guides/my-guide.md` のようなファイルを作成します。

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

該当するセクションですでに `autogenerate` が使用されている場合を除き、`docs/mpress/content/_nav.yaml` の適切なセクションにページを追加してください。読者がナビゲーションや関連ページからのリンクを通じて、そのページに到達できることを確認します。新しい英語ページでは、カバレッジチェックに合格する前に各ロケールの完全なコピーも必要です。メンテナーまたは翻訳コントリビューターと調整してください。

生成された出力は `docs/mpress/site/` にあります。ソースページと、意図して変更したナビゲーションまたは静的アセットをコミットしてください。生成されたサイトファイルとローカルキャッシュはプルリクエストに含めません。

## コントリビューションを確認する

リポジトリのルートから実行します。

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Windows で Python 3 がその方法でインストールされている場合は、`python` を使用してください。[Task](https://taskfile.dev/) がすでにインストールされている場合、`task docs:build` を実行すると、PATH 上の `mpress` を使用して同じチェックが実行されます。

翻訳チェックでは、翻訳ファイルの欠落、空のファイル、対応する原文がないファイルがあると失敗します。古くなった、または破損した翻訳箇所は、デフォルトでは警告として報告されますが、`--strict`を指定すると、これらの検出結果でも失敗します。既存の英語本文を編集した場合は、翻訳の同期が追いつく前でも提出できます。1つの翻訳を修正する場合は、その言語だけを変更できます。コマンド、API名、リンク、コード例、図の接続関係はそのまま維持し、英語のプレースホルダーを翻訳済みコンテンツとして使用しないでください。

変更したページをブラウザーで確認し、変更したコード例をすべて実行してください。変更していないページでチェックが失敗した場合は、チェックを緩和せず、出力と関連するコンテキストをPRに記載してください。

## プルリクエストを送信する

差分を確認し、ソースの変更をコミットして、コントリビューション用ブランチを自分のフォークにプッシュします。`wailsapp/wails`の**master**を対象にPRを作成してください。変更内容を説明し、影響を受けるページへのリンクと、実行したチェックの一覧を記載します。レイアウトを変更した場合は変更前後のスクリーンショットを、コード例を変更した場合はプラットフォームとバージョンの詳細を含めてください。

本番サイトは静的サイトです。ローカルプレビューと検証に、デプロイ用の認証情報や翻訳プロバイダーは必要ありません。レビュー後の公開はメンテナーが行います。

コードのコントリビューションと機能提案については、[Wailsへのコントリビューション](/contributing/)を参照してください。フレームワークの内部構造については、[技術概要](/contributing/overview/)を参照してください。
