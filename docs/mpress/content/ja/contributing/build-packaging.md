---
title: "ビルドおよびパッケージングパイプライン"
description: "`wails3 build` の実行時に内部で行われる処理、クロスプラットフォームバイナリの生成方法、各 OS 向けインストーラーの生成方法について説明します。"
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` は意図的に**薄い**作りになっています。これは Taskfile のラッパーであり、追加のビルドタグをホストプロジェクトの `build` タスクへ転送します。主要な処理は、プロジェクト自身の `build/Taskfile.yml`（`wails3 init` によって生成）、ベイク時のアセットを管理する `internal/commands/build-assets.go`、Linux の nfpm パッケージングを担う `internal/packager`、およびプラットフォーム別インストーラーを担う `internal/commands/appimage.go`、`internal/commands/msix.go`、`internal/commands/dmg/dmg.go`、`internal/commands/dot_desktop.go` に実装されています。

このページでは、以下について説明します。

1. 実際の CLI エントリーポイント
2. Taskfile 主導のビルドフロー
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

`flags.Build` が公開するフラグは**1 つだけ**で、`--tags`（`EXTRA_TAGS=` として転送）です。`wails3 build` には、`-platform`、`-o`、`-skipbindings`、`-skip-package`、`-package`、`-ldflags`、`-verbose`、`-debug`、`-devbuild`、`-icon`、`-clean` の各フラグは**ありません**。クロスコンパイル、出力パス、アイコンなどは、**`Taskfile.yml`**、**`build/config.yml`**、および補助的な `wails3 generate icons` / `wails3 generate build-assets` コマンドで設定します。

`build/build.json` は v3 の一部では**ありません**。設定は `Taskfile.yml` と `build/config.yml` です。

---

## 2. Taskfile 主導のビルドフロー

新しく初期化したプロジェクトには、おおむね次の名前空間を持つ `build/Taskfile.yml` が含まれます。

| 名前空間 | タスク（抜粋） |
| --- | --- |
| `darwin:` | `build`、`build:universal`、`package`、`run`、`dev` |
| `windows:` | `build`、`package`、`run`、`dev` |
| `linux:` | `build`、`package`、`run`、`dev` |
| `common:` | `update:build-assets`、`generate:icons`、`generate:syso` |

`wails3 build` は、デフォルトでホスト OS の `build` 名前空間を呼び出します。続いてプロジェクトの Taskfile が、ホスト固有のフラグを指定して `go build` をシェルから実行します。別の OS 向けにビルドするには、`wails3 build` にフラグを渡すのではなく、その OS のタスクを直接実行します（例：`wails3 task darwin:build:universal`）。

デフォルトの出力ディレクトリは **`bin/<APP_NAME>`** です（`build/bin/` プレフィックスは付きません）。

---

## 3. ベイク時のアセットとビルド情報

| 対象 | ファイル |
| --- | --- |
| ビルドアセットの生成／更新 | `internal/commands/build-assets.go` |
| ビルド情報の出力（CLI：`wails3 tool buildinfo`） | `internal/commands/tool_buildinfo.go` — 情報を出力します。`ldflags` を注入するものでは**ありません** |
| 本番用スタブ | `internal/assetserver/build_production.go` — `//go:build production` |
| フロントエンドバンドル | アプリケーション自身のパッケージ内で `//go:embed` により埋め込み（例：`main.go` の隣） |
| Windows リソース（`.syso`） | `internal/commands/syso.go` — `rsrc_windows_<arch>.syso` を生成 |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG の入力 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI は、アプリケーション用の `bundled_assetserver.go` を自動的にはベイクしません。`internal/assetserver/bundled_assetserver.go` は**手書き**で、`bundledassets/` 配下に埋め込まれた JS ランタイムをラップします。

---

## 4. パッケージングバックエンド

### Linux

Linux のパッケージングは、`fpm` ではなく **nfpm** によって実行されます。

- `internal/packager/packager.go` は `github.com/goreleaser/nfpm/v2` をラップし、`CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)` を公開します。
- 生成されたプロジェクトには、`internal/commands/` 配下に nfpm 形式の設定 `myapp.DEB`、`myapp.RPM`、`myapp.ARCHLINUX` が含まれます（`wails3 tool package` が使用します）。
- AppImage の生成処理は `internal/commands/appimage.go` にあり、そこから `linuxdeploy` + `linuxdeploy-plugin-gtk` を呼び出します（プラグインは `internal/commands/linuxdeploy-plugin-gtk.sh` に同梱されています）。

`wails3 build` には `-package deb`/`rpm` フラグは**ありません**。`wails3 tool package` またはプラットフォーム固有の Taskfile ターゲットを使用してください。

### macOS

- プロジェクトの Taskfile にある `darwin:package` が `.app` バンドルを生成します。
- DMG アセットは `internal/commands/dmg/` 配下にあります。プロジェクトでは `darwin:package` の完了後、`hdiutil` を使用してバンドルを DMG にまとめられます（新しいテンプレートの Taskfile には `dmg` ヘルパーが含まれています）。
- CFBundle の識別子、バージョン、著作権情報は、`wails3 init` の実行時に指定する `-product*` フラグと `build/config.yml` から取得されます。

### Windows

- Windows のパッケージング対象は、WiX/MSI ではなく **MSIX** です。ワークフロー全体については、`internal/commands/msix.go` と `internal/commands/webview2/` を参照してください。
- `internal/commands/packager.go` は**なく**、`internal/commands/windows_resources/` ディレクトリも**ありません**。
- 実行可能ファイルに対する任意のコード署名は、`wails3 tool sign`（Authenticode）を介して実行されます。`internal/commands/sign.go` を参照してください。

---

## 5. パイプラインのカスタマイズ

| 要件 | 方法 |
| --- | --- |
| 追加のビルドタグ | `wails3 build --tags myFeature,otherTag` |
| リンター／ビルド前のステップ | `build/Taskfile.yml` にタスクを追加し、OS 固有の `build` タスクがそのタスクに依存するように設定 |
| クロスコンパイル | 該当する OS タスク（例：`wails3 task linux:build`）を実行します。`-platform` フラグはありません |
| パッケージングを省略 | `build` タスクだけを実行します。`package` は別のタスクです |
| カスタムパッケージャー | `internal/commands/myapp.*` 配下に設定を配置し、`-config <file>` を指定して `wails3 tool package` を呼び出す |
| シンボルを削除する | `darwin:/windows:/linux:` の `build` タスクを編集し、`-ldflags "-s -w"` を `go build` に直接渡す — `wails3 build` 自体には `-ldflags` フラグがない |

Taskfile のすべてのターゲットは、Wails が公開する環境変数（`APP_NAME`、`WAILS_VITE_PORT`、`FRONTEND_DEVSERVER_URL`、…）を使用するため、カスタムタスクでもこれらを利用できます。

---

## 6. トラブルシューティング

| 症状 | 考えられる原因 | 解決方法 |
| --- | --- | --- |
| **`ld: framework not found WebKit`（mac）** | Xcode CLI ツールがない | `xcode-select --install` |
| **本番ビルドでウィンドウが空白になる** | フロントエンドのビルド失敗、または SPA ルーティング | `frontend/dist/index.html` が存在し、アセットハンドラーがこれにフォールバックすることを確認する |
| **MSIX パッケージングツールがない** | `WebView2` SDK／MSIX ツールがインストールされていない | `wails3 task install:msix:tools` を実行する |
| **`linuxdeploy` が見つからない** | プラグインが PATH にない | `linuxdeploy` をインストールし、CLI の自動インストール手順から `internal/commands/linuxdeploy-plugin-gtk.sh` を実行する |

`wails3 build` には `-verbose` フラグがありません。`TASK_X_VERBOSE=1`（Taskfile）を設定するか、タスクターゲットを直接調べて、実行されるコマンドを確認してください。

---

## 7. 主要ソースマップ

| 対象 | ファイル |
| --- | --- |
| ビルドラッパー | `internal/commands/task_wrapper.go`（`Build`、`Package`、`SignWrapper`、`wrapTask`） |
| ビルドアセットの生成 | `internal/commands/build-assets.go`（`GenerateBuildAssets`、`UpdateBuildAssets`） |
| ビルド情報の出力処理 | `internal/commands/tool_buildinfo.go` |
| AppImage ビルダー | `internal/commands/appimage.go` |
| Linux パッケージング（nfpm） | `internal/packager/packager.go`、`internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows MSIX | `internal/commands/msix.go`、`internal/commands/webview2/` |
| Windows リソースジェネレーター | `internal/commands/syso.go` |
| macOS DMG アセット | `internal/commands/dmg/` |
| `.desktop` ジェネレーター | `internal/commands/dot_desktop.go` |
| バージョン定数 | `internal/version/version.go` |

ビルド失敗の原因を追跡するときは、この表を手元に置いてください。

---

これで、**ソースコード**から**インストーラー**までの全体像を把握できました。要するに、`wails3 build` 自体は薄いラッパーにすぎず、ほぼすべてのカスタマイズはプロジェクトの `Taskfile.yml`／`build/config.yml`、または明示的な `wails3 generate …`／`wails3 tool …` サブコマンドで行います。それでは、リリースを成功させましょう！
