---
title: "テストと継続的インテグレーション"
description: "Wails v3 が単体テスト、統合テストスイート、競合検出、GitHub Actions CI によって品質を確保する仕組み。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

堅牢なデスクトップフレームワークには、確実なテストが不可欠です。 Wails v3 では、次の<strong>多層的な戦略</strong>を採用しています。

| レイヤー | 目的 | ツール |
| --- | --- | --- |
| 単体テスト | 独立した関数に対する迅速なフィードバック | `go test ./...` |
| ジェネレーター／CLI テスト | `wails3 generate bindings`と CLI の連携処理を検証 | `task test:generator`、`task test:cli` |
| テンプレートテスト | 配布されるすべてのテンプレートが引き続きビルドできることを確認 | `task test:templates` |
| 競合検出 | ランタイムとブリッジのデータ競合を検出 | `go test -race ./...` |
| CI マトリックス | すべての PR で複数 OS にわたる信頼性を確保 | GitHub Actions |

このドキュメントでは、**テストの配置場所**、**テストの実行方法**、および<strong>Taskfile がオーケストレーションする内容</strong>について説明します。

> 以前のドラフトで`pkg/application/RACE.md`として参照されていた競合に関するガイドは、
>
> 現在、`v3/TESTING.md`にあります。

---

## 1. ディレクトリ規約

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

ガイドライン：

- **単体テストはコードの隣に配置します**（`foo.go` ↔ `foo_test.go`）。
- API の健全性が向上する場合は、`pkg/`パッケージ（`package application_test`）に<strong>ブラックボックス形式</strong>を使用します。
- 共有フィクスチャは、それを使用する場所に配置します（この作業ツリーには一元管理された`internal/testutil/`パッケージがないため、代わりにパッケージごとにヘルパーを組み込みます）。

---

## 2. 単体テスト

### テストの記述

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

推奨事項：

- `go.mod`にすでに含まれている[`stretchr/testify`](https://github.com/stretchr/testify)を使用します。
- 複数の入力、エッジケース、または期待される結果で同じ動作を検証する場合は、可能な限り<strong>テーブル駆動</strong>テストを使用します。各ケースには、その内容が分かる名前を付けてください。
- 必要に応じて、プラットフォーム固有の差異をビルドタグ（`foo_windows_test.go`、`foo_darwin_test.go`、…）の背後でスタブ化します。

### カバレッジ要件

新規または変更されたロジックには、Go ステートメントカバレッジ100% が求められます。リポジトリ全体の割合に依存せず、変更したパッケージを測定してください。

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

通常のテスト環境では、合理的にテストできないパスもあります。たとえば、プラットフォーム固有の障害、ハードウェア依存の動作、安全に発生させることができない防御的なフォールバックなどです。そのような例外は最小限に抑え、カバーされていない各パスについて PR の説明に記載してください。

### ローカルでの実行

```bash
cd v3
go test ./... -cover
```

Taskfile を使用して実行することもできます（実在するターゲットを使用してください。`task test`というショートカットはありません）。

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. 統合テスト

`v3/tests/`には、パッケージ横断の統合テストハーネスがあります。Taskfile の`test:example:*`ターゲットと`test:examples:*`ターゲットは、darwin / windows / linux の各環境でビルドおよび起動チェックを実行します（Linux では Docker ベースの GTK3 / GTK4 マトリックスを含みます）。

> 実行可能なサンプルは`v3/examples/`にあります。テストターゲットは、
>
> ホストプラットフォームまたは CI マトリックスに適したサンプルを選択してビルドします。

ホストプラットフォーム用のスモークテストスイートは、次のコマンドで実行します。

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 競合検出

GUI ランタイムにおいて、データ競合は致命的です。

### 競合ガイド

次の内容については、`v3/TESTING.md`を参照してください。

- 既知の無害な競合と、それを抑制する根拠
- Cgo 境界をまたぐスタックトレースの解釈方法（Linux GTK + WebKit2GTK）

### ローカル競合テストスイート

```
go test -race ./...
```

> `wails3 dev`には`-race`フラグがありません。使用できる CLI フラグは`--config`、`--port`、
>
> および`-s`（HTTPS を有効化）です。競合検出器を有効にしてランタイムをテストするには、
>
> `go build -race`を指定してテストアプリをビルドし、直接実行します。

---

## 5. GitHub Actions ワークフロー

`.github/workflows/`配下にある実際のワークフローファイルは次のとおりです（作業ツリーと照合済み）。

| ファイル | 目的 |
| --- | --- |
| `build-and-test-v3.yml` | v3 のメインのビルド＋テストマトリックス。`go-version: 1.25`とともに`actions/setup-go@v5`を使用します。`task runtime:check`、`task runtime:test`、`task runtime:build`、`task test:examples`（GTK4 パスでは`BUILD_TAGS=gtk4 task test:examples`も使用）、`task generator:test:check`、`task install`を実行した後、スモークチェックとして`wails3 build`を実行します。Linux ジョブでは`libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk`をインストールし、`dbus-run-session -- xvfb-run`上でテストスイートを実行します。 |
| `cross-compile-test-v3.yml` | クロスコンパイルの健全性チェック |
| `auto-changelog-v3.yml`、`changelog-v3.yml` | 変更履歴の自動化 |
| `nightly-release-v3.yml` | v3 のナイトリーリリース成果物 |
| `bump-webview2-v3.yml`、`release-webview2.yml` | WebView2 の依存関係／リリース管理 |
| `build-cross-image.yml` | クロスコンパイラーのコンテナーイメージをビルド |
| `publish-npm.yml` | 組み込みの `@wailsio/runtime` JS ランタイムを npm に公開 |
| `pr-master.yml` | `master` ブランチに対する PR 側のチェック |
| `semgrep.yml` | Semgrep による静的解析 |
| `stale-issues.yml`、`issue-labeler.yml`、`file-labeler.yml`、`claude.yml`、`generate-sponsor-image.yml`、`sync-translated-documents.yml`、`upload-source-documents.yml`、`build-and-test.yml`、`weekly-release-v2.yml` | リポジトリの保守／v2 側のフロー |

このツリーには `qodana.yaml` は<strong>存在せず</strong>、`runtime.yml` も<strong>存在しません</strong>。以前のこのページの草稿では両方に言及していましたが、静的解析を扱うのは `semgrep.yml` だけで、ランタイム JS パッケージは `publish-npm.yml` を通じて公開されます。

CI の各ステップは、上記の Taskfile ターゲット（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`、…）に対応しているため、CI をローカルで一対一に再現できます。`build-and-test-v3.yml` のスモーク `wails3 build` ステップは、<strong>追加のフラグなしで</strong>呼び出されます。`wails3 build` に `-skip-package` フラグはありません。

---

## 6. ローカル環境での CI の再現

単一の包括的な `task ci` ターゲットはありません。実際のターゲットを連続して実行し、CI を再現します：

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 失敗するテストのトラブルシューティング

| 症状 | 考えられる原因 | 修正方法 |
| --- | --- | --- |
| **`webview_window_darwin.go`** での競合 | メインスレッド以外からウィンドウの状態を変更している | 呼び出しがメインスレッドで実行されるよう、`application.InvokeAsync`／`Invoke` を介してマーシャリングする |
| **ヘッドレス CI で Linux テストがハングする** | GTK にはディスプレイが必要 | `xvfb-run` 環境下で実行する（例：`xvfb-run task test:examples:linux`） |
| **テンプレートのビルドが失敗する** | フロントエンドのロックファイルが古い | クリーンなディレクトリに対して `wails3 init` を再実行し、テンプレートを更新する |
| **Coverpkg エラー** | 統合テストが `main` をインポートしている | ビルドタグ `//go:build integration` に切り替え、インポートを条件付きにする |

---

## 8. 新しいテストの追加

1. **ユニットテスト** — `*_test.go` を作成し、`go test ./...` を実行します
2. **ジェネレーター／CLI** — `internal/generator/testcases/` または `internal/commands/*_test.go` 配下のケースを拡充し、`task test:generator`／`task test:cli` を再実行します
3. **テンプレート／サンプル** — 同梱のテンプレートが引き続き `task test:templates` でビルドできることを確認します

---

## 9. 主要ファイル一覧

| 内容 | パス |
| --- | --- |
| ジェネレーターのラウンドトリップテスト | `internal/generator/generate_test.go` |
| アセットのビルドテスト | `internal/commands/build-assets_test.go` |
| 競合／Cgo ガイド | `v3/TESTING.md` |
| Taskfile のテストターゲット | `v3/Taskfile.yaml` |
| イベント定数ジェネレーター | `v3/tasks/events/generate.go` |
| CI ワークフロー | `.github/workflows/build-and-test-v3.yml`（`actions/setup-go@v5` を介した Go 1.25） |
| 静的解析 | `.github/workflows/semgrep.yml` |
| ランタイムの npm 公開 | `.github/workflows/publish-npm.yml` |

---

Wails v3 では、品質は後付けではありません。ユニットテスト、ジェネレーター／テンプレートのテストスイート、競合検出、クロスプラットフォームの CI マトリックスがあるため、変更がサポート対象のすべての OS で正常に動作することを確認しながら、安心してコントリビュートできます。テストを楽しんでください！
