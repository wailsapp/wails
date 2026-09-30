---
title: "ビルドシステム"
description: "Wails がアプリケーションをビルドしてパッケージ化する仕組み"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## 統合ビルドシステム

Wails は、Go コードのコンパイル、フロントエンドアセットのバンドル、すべての要素の単一実行ファイルへの埋め込み、プラットフォーム固有のビルドを、1 つのコマンドで処理する<strong>統合ビルドシステム</strong>を提供します。

```bash
wails3 build
```

<strong>出力：</strong>すべてが埋め込まれたネイティブ実行ファイル。

## ビルドプロセスの概要

**［ビルドプロセス図のプレースホルダー］**

## ビルドフェーズ

### 1. 解析フェーズ

Wails は Go コードをスキャンしてサービスを把握します：

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Wails が抽出する情報：**

- サービス名：`GreetService`
- メソッド名：`Greet`
- パラメーターの型：`string`
- 戻り値の型：`string`

<strong>用途：</strong>TypeScript バインディングの生成

### 2. 生成フェーズ

#### TypeScript バインディング

Wails は型安全なバインディングを生成します：

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**利点：**

- 完全な型安全性
- IDE の自動補完
- コンパイル時のエラー検出
- JSDoc コメント

#### フロントエンドのビルド

フロントエンドバンドラー（Vite、webpack など）が実行されます：

```bash
# Vite example
vite build --outDir dist
```

**実行される処理：**

- JavaScript／TypeScript のコンパイル
- CSS の処理と圧縮
- アセットの最適化
- ソースマップの生成（開発時のみ）
- `frontend/dist/`への出力

### 3. コンパイルフェーズ

#### Go のコンパイル

Go コードは最適化を有効にしてコンパイルされます：

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**フラグ：**

- `-s`：シンボルテーブルを削除
- `-w`：DWARF デバッグ情報を削除
- 結果：バイナリを小型化（約30% 削減）

**プラットフォーム固有の形式：**

- Windows：アイコンが埋め込まれた `.exe`
- macOS：`.app`バンドル構造
- Linux：ELF バイナリ

#### アセットの埋め込み

フロントエンドアセットは Go バイナリに埋め込まれます：

```go
//go:embed frontend/dist
var assets embed.FS
```

<strong>結果：</strong>すべてを内包した単一の実行ファイル。

### 4. 出力

**単一のネイティブバイナリ：**

- Windows：`myapp.exe`（約 15MB）
- macOS：`myapp.app`（約 15MB）
- Linux：`myapp`（約 15MB）

**依存関係なし**（システムの WebView を除く）。

## 開発環境と本番環境

@tabs{sync-key="mode"}
[開発（wails3 dev）]
**速度を重視して最適化：**

```bash
wails3 dev
```

**実行される処理：**

1. フロントエンド開発サーバーを起動（デフォルトでは Vite をポート 9245 で起動）
2. 最適化せずに Go をコンパイル
3. 開発サーバーを参照するようにアプリを起動
4. ホットリロードを有効化
5. ソースマップを含める

**特徴：**

- **高速な再ビルド**（フロントエンドの変更は &lt;1 秒）
- **アセットを埋め込まない**（開発サーバーから配信）
- <strong>デバッグシンボル</strong>を含む
- <strong>ソースマップ</strong>を有効化
- **詳細なログ出力**

<strong>ファイルサイズ：</strong>大きい（デバッグシンボルを含めて約50MB）

[本番環境（wails3 build）]
**サイズとパフォーマンスを最適化：**

```bash
wails3 build
```

**実行される処理：**

1. 本番環境向けにフロントエンドをビルド（縮小）
2. 最適化を有効にしてGoをコンパイル
3. デバッグシンボルを削除
4. アセットを埋め込み
5. 単一のバイナリを作成

**特性：**

- **最適化されたコード**（縮小、ツリーシェイキング済み）
- **アセットを埋め込み済み**（外部ファイルなし）
- **デバッグシンボルを削除済み**
- **ソースマップなし**
- **最小限のログ出力**

<strong>ファイルサイズ：</strong>小さい（約15MB）

@end

## ビルドコマンド

### 基本的なビルド

```bash
wails3 build
```

**出力：**`bin/<APP_NAME>`（Windowsでは`bin/<APP_NAME>.exe`）。`bin/`ディレクトリはプロジェクトルートにあります。

`wails3 build`は`wails3 task build`の薄いラッパーです。転送するビルド時フラグは`--tags`のみで、これはTaskfile変数`EXTRA_TAGS`になります：

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build`には、`-platform`、`-o`、`-skipbindings`、`-clean`、`-debug`、`-devbuild`、`-icon`、`-ldflags`、`-package`の各フラグはありません。クロスコンパイル、出力パス、アイコン、パッケージ化は、プロジェクトのTaskfile（`Taskfile.yml` + `build/config.yml`）で制御します。

### クロスプラットフォームおよびプラットフォーム固有のビルド

プラットフォーム向けビルドは、`build/Taskfile.<platform>.yml`で定義された`darwin:` / `windows:` / `linux:`名前空間配下のTaskfileタスクとして公開されています。例：

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

現在のプロジェクトで利用可能なすべてのタスクを表示するには、次を実行します：

```bash
wails3 task --list
```

### アイコンとパッケージ化

ソースPNGからプラットフォーム用アイコン（`build/icons.icns`、`build/icon.ico`など）を生成します：

```bash
wails3 generate icons -input appicon.png
```

プラットフォーム固有のインストーラーまたはパッケージをビルドします：

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## ビルド設定

### Taskfile.yml

Wails 3プロジェクトでは、ビルドオーケストレーターとして[Taskfile](https://taskfile.dev/)を使用します。ルートの`Taskfile.yml`は、`build/`にあるプラットフォーム別のタスクファイルをインクルードします：

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

`wails3 task <name>`または`task <name>`でタスクを実行します：

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### プロジェクト設定：`build/config.yml`

プロジェクトのメタデータ（名前、識別子、バージョン、Info.plistの値、NSIS設定、`.desktop`フィールド、カスタムプロトコルなど）は`build/config.yml`にあります。Taskfileは、アイコン、マニフェスト、インストーラーなどを生成するときにこのファイルを読み取ります。Wails 3には`build/build.json`ファイルは<strong>存在しません</strong>。

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

この設定からプラットフォーム固有のビルドアセットを更新するには、`wails3 generate build-assets`（または`wails3 update build-assets`）を実行します。

## アセットの埋め込み

### 仕組み

WailsはGoの`embed`パッケージを使用します：

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**ビルド時：**

1. フロントエンドを`frontend/dist/`にビルド
2. `//go:embed`ディレクティブでファイルを取り込み
3. ファイルをバイナリにコンパイル
4. バイナリにすべてを格納

**実行時：**

1. アプリを起動
2. メモリからアセットを提供
3. アセットのディスクI/Oなし
4. 高速な読み込み

### カスタムアセット

追加ファイルを埋め込みます：

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## ビルドの最適化

### フロントエンドの最適化

**Vite（デフォルト）：**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**結果：**

- JavaScriptを縮小（約70%削減）
- CSSを縮小（約60%削減）
- 画像を最適化
- ツリーシェイキングを適用

### Goの最適化

**コンパイラフラグ：**

```bash
-ldflags="-s -w"
```

- `-s`：シンボルテーブルを削除（約10%削減）
- `-w`：DWARFデバッグ情報を削除（約20%削減）

**その他の最適化：**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`：ビルド時に変数の値を設定
- バージョン番号やビルド日時に便利です

### バイナリ圧縮

**UPX（任意）：**

```bash
# After building
upx --best bin/myapp.exe
```

**結果：**

- サイズを約50%削減
- 起動がわずかに遅くなります（約100ms）
- macOSでは推奨されません（コード署名の問題があるため）

## プラットフォーム固有のビルド

### Windows

**出力：** `myapp.exe`

**含まれるもの：**

- アプリケーションアイコン
- バージョン情報
- マニフェスト（UAC設定）

**アイコン：**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

その後、Windowsの`tool package`ステップで、生成された`.ico`を実行可能ファイルに埋め込みます。

**マニフェスト：**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**出力：** `myapp.app`（アプリケーションバンドル）

**構造：**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist：**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**ユニバーサルバイナリ：**

macOSのTaskfileには、両方のアーキテクチャをビルドして`wails3 tool lipo`で結合する`darwin:build:universal`タスク（および`darwin:package:universal`タスク）が用意されています：

```bash
wails3 task darwin:build:universal
```

### Linux

**出力：** `myapp`（ELFバイナリ）

**依存関係：**

- GTK3
- WebKitGTK

**デスクトップファイル：**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**インストール：**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## ビルドパフォーマンス

### 一般的なビルド時間

| フェーズ | 所要時間 | 備考 |
| --- | --- | --- |
| 解析 | &lt;1s | Goコードのスキャン |
| バインディング生成 | &lt;1s | TypeScriptの生成 |
| フロントエンドのビルド | 5-30s | プロジェクトの規模によって異なります |
| Goのコンパイル | 2-10s | コードの規模によって異なります |
| アセットの埋め込み | &lt;1s | フロントエンドの埋め込み |
| **合計** | **10-45s** | 初回ビルド |
| **インクリメンタル** | **5-15s** | 2回目以降のビルド |

### ビルドの高速化

**1. ビルドキャッシュを使用する：**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. 必要な処理だけを実行する：**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. 並列ビルド（複数マシン／CI）：**

v3でのLinux、Windows、macOS間のクロスコンパイルは、通常、Dockerの`wails-cross`コンテナ内、またはプラットフォームごとの専用ランナー上で行います。`wails3 build`自体のターゲットはホストOSです。サポートされているワークフローについては、[クロスプラットフォームビルド](/guides/build/cross-platform/)を参照してください。

**4. 高速なツールを使用する：**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## トラブルシューティング

### ビルドが失敗する

**症状：** `wails3 build` がエラーで終了する

**一般的な原因：**

1. **Go のコンパイルエラー**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **フロントエンドのビルドエラー**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **依存関係が不足している**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### バイナリが大きすぎる

**症状：** バイナリが 50 MB を超えている

**解決策：**

1. **デバッグシンボルを削除する**（同梱の Taskfile では、すでに `-ldflags="-s -w"` が `go build` に渡されています）。

2. **埋め込まれたアセットを確認する**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **UPX 圧縮を使用する**
  ```bash
  upx --best bin/myapp.exe
  ```


### ビルドが遅い

**症状：** ビルドに 1 分を超える時間がかかる

**解決策：**

1. **ビルドキャッシュを使用する**
  - Go のキャッシュは自動的に使用される
  - フロントエンドのキャッシュ（Vite）は自動的に使用される


2. **必要なタスクだけを実行する**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **フロントエンドのビルドを最適化する**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## ベストプラクティス

### ✅ 推奨事項

- **開発中は `wails3 dev` を使用する** — 高速に反復できる
- **リリースには `wails3 build` を使用する** — 出力が最適化される
- **ビルドにバージョンを付ける** — バージョンの埋め込みには `-ldflags` を使用する
- **対象プラットフォームでビルドをテストする** — クロスコンパイルは完全ではない
- **フロントエンドのビルドを高速に保つ** — バンドラーの設定を最適化する
- **ビルドキャッシュを使用する** — 2 回目以降のビルドが高速になる

### ❌ 禁止事項

- **`build/` ディレクトリをコミットしない** — `.gitignore` に追加する
- **ビルドのテストを省略しない** — リリース前に必ずテストする
- **不要なアセットを埋め込まない** — バイナリを小さく保つ
- **本番環境ではデバッグビルドを使用しない** — 最適化されたビルドを使用する
- **コード署名を忘れない** — 配布には必須

## 次のステップ

**アプリケーションのビルド** — ビルドとパッケージ化の詳細ガイド [詳細を見る →](/guides/build/building/)

**クロスプラットフォームビルド** — 1 台のマシンですべてのプラットフォーム向けにビルドする [詳細を見る →](/guides/build/cross-platform/)

**インストーラーの作成** — エンドユーザー向けのインストーラーを作成する [詳細を見る →](/guides/installers/)

---

**ビルドについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[ビルド例](https://github.com/wailsapp/wails/tree/master/v3/examples/build)を確認してください。
