---
title: "ビルドとパッケージングのパイプライン"
description: "`wails3 build` の実行時に内部で行われる処理、クロスプラットフォームバイナリの生成方法、および各 OS 向けインストーラーの生成方法。"
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` は意図的に<strong>薄い</strong>実装になっています。これは、追加のビルドタグをホストプロジェクトの `build` タスクに転送する Taskfile ラッパーです。主要な処理は、プロジェクト独自の `build/Taskfile.yml`（`wails3 init` によって生成）、`internal/commands/build-assets.go`（ベイク時のアセットを管理）、`internal/packager`（Linux の nfpm パッケージング）、および `internal/commands/appimage.go`、`internal/commands/msix.go`、`internal/commands/dmg/dmg.go`、`internal/commands/dot_desktop.go`（プラットフォーム別インストーラー）にあります。

このページでは、次の内容を説明します。

1. 実際の CLI エントリーポイント
2. Taskfile 駆動のビルドフロー
3. アセットのベイクとビルド情報の注入
4. プラットフォーム別のパッケージングバックエンド
5. パイプラインのカスタマイズ
6. トラブルシューティング

---

## 1. 実際の CLI エントリーポイント

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`：

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` が公開するフラグは、<strong>1 つ</strong>だけです。それは `--tags`（`EXTRA_TAGS=` として転送）です。`wails3 build` には、**いずれも** `-platform`、`-o`、`-skipbindings`、`-skip-package`、`-package`、`-ldflags`、`-verbose`、`-debug`、`-devbuild`、`-icon`、`-clean` の各フラグはありません。クロスコンパイル、出力パス、アイコンなどは、**`Taskfile.yml`**、**`build/config.yml`**、および補助的な `wails3 generate icons` / `wails3 generate build-assets` コマンドで構成します。

`build/build.json` は v3 の一部では<strong>ありません</strong>。構成には `Taskfile.yml` と `build/config.yml` を使用します。

---

## 2. Taskfile 駆動のビルドフロー

新しく初期化したプロジェクトには、概ね次の名前空間を持つ `build/Taskfile.yml` が含まれています。

| 名前空間 | タスク（抜粋） |
| --- | --- |
| `darwin:` | `build`、`build:universal`、`package`、`run`、`dev` |
| `windows:` | `build`、`package`、`run`、`dev` |
| `linux:` | `build`、`package`、`run`、`dev` |
| `common:` | `update:build-assets`、`generate:icons`、`generate:syso` |

`wails3 build` は、デフォルトでホスト OS の `build` 名前空間を呼び出します。続いて、プロジェクトの Taskfile がホスト固有のフラグを指定して `go build` を実行します。別の OS 向けにビルドするには、`wails3 build` にフラグを渡すのではなく、その OS のタスクを直接実行します（例：`wails3 task darwin:build:universal`）。

デフォルトの出力ディレクトリは **`bin/<APP_NAME>`** です（`build/bin/` プレフィックスは付きません）。

---

## 3. ベイク時のアセットとビルド情報

| 対象 | ファイル |
| --- | --- |
| ビルドアセットの生成／更新 | `internal/commands/build-assets.go` |
| ビルド情報の表示（CLI：`wails3 tool buildinfo`） | `internal/commands/tool_buildinfo.go` — 情報を表示します。`ldflags` インジェクターでは<strong>ありません</strong> |
| 本番用スタブ | `internal/assetserver/build_production.go` — `//go:build production` |
| フロントエンドバンドル | アプリケーション独自のパッケージ内で `//go:embed` を使用して埋め込みます（例：`main.go` の隣） |
| Windows リソース（`.syso`） | `internal/commands/syso.go` — `rsrc_windows_<arch>.syso` を生成します |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG の入力 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI は、アプリケーション用の `bundled_assetserver.go` を自動的にベイクしません。`internal/assetserver/bundled_assetserver.go` は<strong>手書き</strong>で、`bundledassets/` 配下に埋め込まれた JS ランタイムをラップします。

---

## 4. パッケージングバックエンド

### Linux

Linux のパッケージングは、`fpm` ではなく **nfpm** によって実行されます。

- `internal/packager/packager.go` は `github.com/goreleaser/nfpm/v2` をラップし、`CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)` を公開します。
- 生成されたプロジェクトには、nfpm 形式の構成 `myapp.DEB`、`myapp.RPM`、`myapp.ARCHLINUX` が `internal/commands/` 配下に含まれています（`wails3 tool package` が使用）。
- AppImage の生成処理は `internal/commands/appimage.go` にあり、`linuxdeploy` + `linuxdeploy-plugin-gtk` を呼び出します（プラグインは `internal/commands/linuxdeploy-plugin-gtk.sh` に同梱されています）。

`wails3 build` には、`-package deb`/`rpm` フラグは<strong>ありません</strong>。`wails3 tool package` またはプラットフォーム固有の Taskfile ターゲットを使用してください。

### macOS

- プロジェクトの Taskfile にある `darwin:package` は、`.app` バンドルを生成します。
- DMG アセットは `internal/commands/dmg/` にあります。プロジェクトでは、`darwin:package` の完了後に `hdiutil` を使用してバンドルを DMG にまとめられます（新しいテンプレートの Taskfile には `dmg` ヘルパーが含まれています）。
- CFBundle の識別子、バージョン、著作権情報は、`wails3 init` の実行時に指定する `-product*` フラグと、`build/config.yml` から取得されます。

### Windows

- Windows のパッケージ化では **MSIX** を対象とします（WiX/MSI ではありません）。ワークフロー全体については、`internal/commands/msix.go` と `internal/commands/webview2/` を参照してください。
- **`internal/commands/packager.go` はなく**、**`internal/commands/windows_resources/` ディレクトリもありません**。
- 実行ファイルへの任意のコード署名は、`wails3 tool sign`（Authenticode）を介して実行されます。`internal/commands/sign.go` を参照してください。

---

## 5. パイプラインのカスタマイズ

| 要件 | 方法 |
| --- | --- |
| 追加のビルドタグ | `wails3 build --tags myFeature,otherTag` |
| リンター／ビルド前処理 | `build/Taskfile.yml` にタスクを追加し、OS 固有の `build` タスクがそのタスクに依存するようにします |
| クロスコンパイル | 該当する OS のタスク（例：`wails3 task linux:build`）を実行します。`-platform` フラグはありません |
| パッケージ化を省略 | `build` タスクだけを実行します。`package` は別のタスクです |
| カスタムパッケージャー | `internal/commands/myapp.*` 配下に設定を配置し、`-config <file>` を指定して `wails3 tool package` を呼び出します |
| シンボルの除去 | `darwin:/windows:/linux:` の `build` タスクを編集し、`-ldflags "-s -w"` を `go build` に直接渡します。`wails3 build` 自体には `-ldflags` フラグがありません |

すべての Taskfile ターゲットは、Wails が公開する環境変数（`APP_NAME`、`WAILS_VITE_PORT`、`FRONTEND_DEVSERVER_URL`、…）に従うため、カスタムタスクでもそれらを利用できます。

---

## 6. トラブルシューティング

| 症状 | 考えられる原因 | 解決方法 |
| --- | --- | --- |
| **`ld: framework not found WebKit`（mac）** | Xcode CLI ツールがありません | `xcode-select --install` |
| **本番ビルドでウィンドウが空白になる** | フロントエンドのビルド失敗、または SPA ルーティング | `frontend/dist/index.html` が存在し、アセットハンドラーがそれにフォールバックすることを確認します |
| **MSIX パッケージ化ツールがありません** | `WebView2` SDK／MSIX ツールがインストールされていません | `wails3 task install:msix:tools` を実行します |
| **`linuxdeploy` が見つかりません** | プラグインが PATH にありません | `linuxdeploy` をインストールし、CLI の自動インストール処理を介して `internal/commands/linuxdeploy-plugin-gtk.sh` を実行します |

`wails3 build` には `-verbose` フラグがありません。`TASK_X_VERBOSE=1`（Taskfile）を設定するか、タスクのターゲットを直接調べて、実行されるコマンドを確認してください。

---

## 7. 主要ソースマップ

| 対象 | ファイル |
| --- | --- |
| ビルドラッパー | `internal/commands/task_wrapper.go`（`Build`、`Package`、`SignWrapper`、`wrapTask`） |
| ビルドアセットの生成 | `internal/commands/build-assets.go`（`GenerateBuildAssets`、`UpdateBuildAssets`） |
| ビルド情報の出力処理 | `internal/commands/tool_buildinfo.go` |
| AppImage ビルダー | `internal/commands/appimage.go` |
| Linux のパッケージ化（nfpm） | `internal/packager/packager.go`、`internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows MSIX | `internal/commands/msix.go`、`internal/commands/webview2/` |
| Windows リソースジェネレーター | `internal/commands/syso.go` |
| macOS DMG アセット | `internal/commands/dmg/` |
| `.desktop` ジェネレーター | `internal/commands/dot_desktop.go` |
| バージョン定数 | `internal/version/version.go` |

ビルド失敗の原因を追跡するときに、この表を手元に置いてください。

---

これで、<strong>ソースコード</strong>から<strong>インストーラー</strong>までの全体像を把握できました。要点は、`wails3 build` 自体は薄いラッパーにすぎず、ほぼすべてのカスタマイズはプロジェクトの `Taskfile.yml`／`build/config.yml`、または明示的な `wails3 generate …`／`wails3 tool …` サブコマンドを介して行うということです。それでは、リリースを成功させましょう！
