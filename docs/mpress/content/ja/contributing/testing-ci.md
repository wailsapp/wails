---
title: "テストと継続的インテグレーション"
description: "Wails v3 が単体テスト、統合テストスイート、競合検出、GitHub Actions CI によって品質を確保する仕組み。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

堅牢なデスクトップフレームワークには、信頼性の高いテストが不可欠です。Wails v3 では、次の**多層戦略**を採用しています。

| レイヤー | 目的 | ツール |
| --- | --- | --- |
| 単体テスト | 独立した関数に対する迅速なフィードバック | `go test ./...` |
| ジェネレーター／CLI テスト | `wails3 generate bindings` と CLI の連携処理を検証 | `task test:generator`, `task test:cli` |
| テンプレートテスト | 提供されるすべてのテンプレートが引き続きビルドできることを確認 | `task test:templates` |
| 競合検出 | ランタイムとブリッジのデータ競合を検出 | `go test -race ./...` |
| CI マトリックス | すべての PR で各 OS にわたる信頼性を確保 | GitHub Actions |

このドキュメントでは、<strong>テストの配置場所</strong>、<strong>テストの実行方法</strong>、および<strong> Taskfile が統括する処理</strong>について説明します。

`v3/TESTING.md` はクロスプラットフォームのサンプルビルドについて説明するものであり、競合を抑制する方法の一覧ではありません。

---

## 1. ディレクトリ規約

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

ガイドライン：

- **単体テストはコードの隣に配置します**（`foo.go` ↔ `foo_test.go`）。
- API の健全性が向上する場合は、`pkg/` パッケージ（`package application_test`）に**ブラックボックス形式**を使用します。
- 共有フィクスチャは使用箇所に配置します（この作業ツリーには共通の `internal/testutil/` パッケージがないため、代わりにパッケージごとにヘルパーを組み込みます）。

---

## 2. 単体テスト

### テストの記述

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

推奨事項：

- [`stretchr/testify`](https://github.com/stretchr/testify) を使用します。これはすでに `go.mod` に含まれています。
- 複数の入力、エッジケース、または期待結果で同じ動作を検証する場合は、可能な限り**テーブル駆動**テストを使用します。各ケースには内容が分かる名前を付けてください。
- 必要に応じて、プラットフォーム固有の挙動をビルドタグ（`foo_windows_test.go`、`foo_darwin_test.go`、…）の背後でスタブ化します。

### カバレッジ要件

新規および変更されたロジックには、Go ステートメントカバレッジ 100% が求められます。リポジトリ全体の割合に依存せず、変更したパッケージを測定してください。

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

通常のテスト環境では現実的にテストできないパスもあります。たとえば、特定のプラットフォームでのみ発生する障害、ハードウェア依存の動作、安全には誘発できない防御的フォールバックなどです。このような例外は最小限にとどめ、カバーされていない各パスについて PR の説明に記載してください。

### ローカルでの実行

```bash
cd v3
go test ./... -cover
```

Taskfile 経由でも実行できます（実在するターゲットを使用してください。`task test` というショートカットはありません）。

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. 統合テスト

`v3/tests/` には、パッケージ横断の統合テストハーネスがあります。Taskfile の `test:example:*` および `test:examples:*` ターゲットは、darwin／windows／linux 全体のコンパイルチェックを実行します（Linux では Docker ベースの GTK3／GTK4 マトリックスも含みます）。

> 実行可能なサンプルは `v3/examples/` にあります。テストターゲットは、
>
> ホストプラットフォームまたは CI マトリックスに適したサンプルを選択してビルドします。

ホストプラットフォームのスモークテストスイートは、次のコマンドで実行します。

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 競合検出

GUI ランタイムにおいて、データ競合は致命的です。

### 競合の調査

競合している両方のアクセスを調査し、影響を受けるプラットフォームで再現してください。根拠なしにレポートを無害とみなしたり、抑制したりしないでください。`v3/TESTING.md` はサンプルのコンパイルについて説明するものです。これらのビルドに成功しても、競合が存在しないことやネイティブランタイムの動作は確認できません。

### ローカル競合テストスイート

```
go test -race ./...
```

> `wails3 dev` に `-race` フラグはありません。CLI フラグは `--config`、`--port`、
>
> および `-s`（HTTPS を有効化）です。競合検出器を有効にしてランタイムを実行するには、
>
> `go build -race` を使用してテストアプリをビルドし、直接実行します。

---

## 5. GitHub Actions ワークフロー

`.github/workflows/` 配下に実在するワークフローファイル（作業ツリーと照合済み）：

| ファイル | 目的 |
| --- | --- |
| `build-and-test-v3.yml` | Go 1.26/1.27 マトリックス、ランタイムとジェネレーターのチェック、テンプレート／サンプルのコンパイル、およびアプリケーションのビルドチェック。Linux ではデフォルトで GTK4/WebKitGTK 6.0 を使用し、レガシーバリアントでは `BUILD_TAGS=gtk3` を設定します。ネイティブでの起動と操作は別途確認する必要があります。 |
| `cross-compile-test-v3.yml` | クロスコンパイルの健全性チェック |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | 変更履歴の自動化 |
| `nightly-release-v3.yml` | ナイトリー v3 リリース成果物 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 の依存関係／リリース管理 |
| `build-cross-image.yml` | クロスコンパイラーのコンテナイメージをビルド |
| `publish-npm.yml` | 組み込みの `@wailsio/runtime` JS ランタイムを npm に公開 |
| `pr-master.yml` | `master` ブランチに対する PR 側のチェック |
| `semgrep.yml` | Semgrep による静的解析 |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | リポジトリの保守作業と v2 ワークフロー。旧 Crowdin 同期／アップロードワークフローは廃止されています。 |

CI の各ステップは、上記の Taskfile ターゲット（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`、…）に対応しているため、ローカルで CI を一対一に再現できます。`build-and-test-v3.yml` のスモークテスト用 `wails3 build` ステップは、**追加フラグなし**で呼び出されます。`wails3 build` に `-skip-package` フラグはありません。

---

## 6. ローカル環境での CI 再現

すべてをまとめて実行する単一の `task ci` ターゲットはありません。実在するターゲットを連結して CI を再現してください。

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 失敗したテストのトラブルシューティング

| 症状 | 考えられる原因 | 修正方法 |
| --- | --- | --- |
| **`webview_window_darwin.go`での競合状態** | メインスレッド以外からウィンドウの状態を変更している | 呼び出しがメインスレッドで実行されるよう、`application.InvokeAsync` / `Invoke`を介してマーシャリングする |
| **ヘッドレスCIでLinuxテストがハングする** | GTKにはディスプレイが必要 | `xvfb-run`環境下で実行する（例：`xvfb-run task test:examples:linux`） |
| **テンプレートのビルドが失敗する** | フロントエンドのロックファイルが古い | テンプレートを更新するため、空のディレクトリを対象に`wails3 init`を再実行する |
| **Coverpkgエラー** | 統合テストが`main`をインポートしている | ビルドタグ`//go:build integration`に切り替え、インポートを条件付きにする |

---

## 8. 新しいテストの追加

1. **ユニットテスト** — `*_test.go`を作成し、`go test ./...`を実行する
2. **ジェネレーター / CLI** — `internal/generator/testcases/`または`internal/commands/*_test.go`配下のテストケースを拡充し、`task test:generator` / `task test:cli`を再実行する
3. **テンプレート / サンプル** — 配布されるテンプレートが`task test:templates`で引き続きビルドできることを確認する

---

## 9. 主要ファイル一覧

| 項目 | パス |
| --- | --- |
| ジェネレーターのラウンドトリップテスト | `internal/generator/generate_test.go` |
| ビルドアセットのテスト | `internal/commands/build-assets_test.go` |
| クロスプラットフォームのサンプルビルド | `v3/TESTING.md` |
| Taskfileのテストターゲット | `v3/Taskfile.yaml` |
| イベント定数ジェネレーター | `v3/tasks/events/generate.go` |
| CIワークフロー | `.github/workflows/build-and-test-v3.yml`（`actions/setup-go@v5` による Go 1.26/1.27） |
| 静的解析 | `.github/workflows/semgrep.yml` |
| ランタイムnpmパッケージの公開 | `.github/workflows/publish-npm.yml` |

---

Wails v3では、品質を後回しにしていません。ユニットテスト、ジェネレーターおよびテンプレートのテストスイート、競合状態の検出、クロスプラットフォームCIマトリックスにより、変更がサポート対象のすべてのOSで正常に動作することを確認しながら、安心してコントリビュートできます。テストを楽しみましょう！
