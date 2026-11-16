---
title: "ビルドシステム"
description: "Wails がアプリケーションをビルドしてパッケージ化する仕組み"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## 統合ビルドシステム

Wails は、Go コードのコンパイル、フロントエンドアセットのバンドル、すべての単一実行ファイルへの埋め込み、プラットフォーム固有のビルドを、1 つのコマンドで処理する**統合ビルドシステム**を提供します。

```bash
wails3 build
```

<strong>出力：</strong>すべてが埋め込まれたネイティブ実行ファイル。

## ビルドプロセスの概要

**[ビルドプロセス図のプレースホルダー]**

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
- コンパイル時エラー
- JSDoc コメント

#### フロントエンドのビルド

フロントエンドのバンドラー（Vite、webpack など）が実行されます：

```bash
# Vite example
vite build --outDir dist
```

**実行される処理：**

- JavaScript／TypeScript のコンパイル
- CSS の処理と圧縮
- アセットの最適化
- ソースマップの生成（開発時のみ）
- `frontend/dist/` に出力

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

**プラットフォーム固有の出力：**

- Windows：アイコンが埋め込まれた `.exe`
- macOS：`.app` バンドル構造
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

- Windows：`myapp.exe`（約15MB）
- macOS：`myapp.app`（約15MB）
- Linux：`myapp`（約15MB）

**依存関係なし**（システムの WebView を除く）。

## 開発環境と本番環境

@tabs{sync-key="mode"}
[開発（wails3 dev）]
**速度を重視して最適化：**

```bash
wails3 dev
```

**実行される処理：**

1. フロントエンド開発サーバーを起動（デフォルトではポート 9245 で Vite を実行）
2. 最適化せずに Go をコンパイル
3. 開発サーバーを参照するアプリを起動
4. ホットリロードを有効化
5. ソースマップを含める

**特性：**

- **高速な再ビルド**（フロントエンド変更時は &lt;1 秒）
- **アセットの埋め込みなし**（開発サーバーから配信）
- **デバッグシンボル**を含む
- **ソースマップ**を有効化
- **詳細ログ**

<strong>ファイルサイズ：</strong>大きい（デバッグシンボルを含めて約50MB）

[本番（wails3 build）]
**サイズとパフォーマンスを重視して最適化：**

```bash
wails3 build
```

**実行される処理：**

1. 本番用にフロントエンドをビルド（圧縮）
2. 最適化を有効にして Go をコンパイル
3. デバッグシンボルを削除
4. アセットを埋め込み
5. 単一のバイナリを作成

**特性：**

- **最適化されたコード**（圧縮、ツリーシェイキング済み）
- **埋め込み済みアセット**（外部ファイルなし）
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

**出力：**`bin/<APP_NAME>`（Windows では `bin/<APP_NAME>.exe`）。`bin/` ディレクトリはプロジェクトルートに配置されます。

`wails3 build` は `wails3 task build` の薄いラッパーです。ビルド時に転送する唯一のフラグは `--tags` で、これは Taskfile の変数 `EXTRA_TAGS` になります：

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build` には、`-platform`、`-o`、`-skipbindings`、`-clean`、`-debug`、`-devbuild`、`-icon`、`-ldflags`、`-package` の各フラグはありません。クロスコンパイル、出力パス、アイコン、パッケージ化は、プロジェクトの Taskfile（`Taskfile.yml` + `build/config.yml`）で制御します。

### クロスプラットフォームおよびプラットフォーム固有のビルド

プラットフォーム向けビルドは、`build/Taskfile.<platform>.yml` で定義された `darwin:` / `windows:` / `linux:` 名前空間内の Taskfile タスクとして提供されます。例：

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

現在のプロジェクトで利用可能なすべてのタスクを表示するには、次を実行します。

```bash
wails3 task --list
```

### アイコンとパッケージング

ソース PNG からプラットフォーム用アイコン（`build/icons.icns`、`build/icon.ico` など）を生成します。

```bash
wails3 generate icons -input appicon.png
```

プラットフォーム固有のインストーラー／パッケージをビルドします。

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## ビルド設定

### Taskfile.yml

Wails 3 プロジェクトは、ビルドのオーケストレーターとして [Taskfile](https://taskfile.dev/) を使用します。ルートの `Taskfile.yml` は、`build/` にあるプラットフォーム別のタスクファイルをインクルードします。

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

`wails3 task <name>` または `task <name>` でタスクを実行します。

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### プロジェクト設定：`build/config.yml`

プロジェクトのメタデータ（名前、識別子、バージョン、info-plist の値、NSIS 設定、`.desktop` フィールド、カスタムプロトコルなど）は `build/config.yml` にあります。Taskfile は、アイコン、マニフェスト、インストーラーなどを生成するときにこのファイルを読み取ります。Wails 3 には **`build/build.json` ファイルはありません**。

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

この設定からプラットフォーム固有のビルドアセットを更新するには、`wails3 generate build-assets`（または `wails3 update build-assets`）を実行します。

## アセットの埋め込み

### 仕組み

Wails は Go の `embed` パッケージを使用します。

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

1. フロントエンドを `frontend/dist/` にビルド
2. `//go:embed` ディレクティブでファイルを取り込む
3. ファイルをバイナリにコンパイル
4. バイナリにすべてを格納

**実行時：**

1. アプリを起動
2. アセットをメモリから配信
3. アセットのディスク I/O は不要
4. 高速な読み込み

### カスタムアセット

追加のファイルを埋め込みます。

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

- JavaScript を圧縮（約 70% 削減）
- CSS を圧縮（約 60% 削減）
- 画像を最適化
- ツリーシェイキングを適用

### Go の最適化

**コンパイラーフラグ：**

```bash
-ldflags="-s -w"
```

- `-s`：シンボルテーブルを削除（約 10% 削減）
- `-w`：DWARF デバッグ情報を削除（約 20% 削減）

**追加の最適化：**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`：ビルド時に変数の値を設定
- バージョン番号やビルド日付に便利

### バイナリ圧縮

**UPX（任意）：**

```bash
# After building
upx --best bin/myapp.exe
```

**結果：**

- サイズを約 50% 削減
- 起動がわずかに遅くなる（約 100ms）
- macOS では非推奨（コード署名の問題）

## プラットフォーム固有のビルド

### Windows

**出力：** `myapp.exe`

**含まれるもの：**

- アプリケーションアイコン
- バージョン情報
- マニフェスト（UAC 設定）

**アイコン：**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

続いて Windows の `tool package` ステップで、生成された `.ico` を実行可能ファイルに埋め込みます。

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

macOS の Taskfile には、両方のアーキテクチャ向けにビルドし、`wails3 tool lipo` で結合する `darwin:build:universal`（および `darwin:package:universal`）タスクが用意されています。

```bash
wails3 task darwin:build:universal
```

### Linux

**出力：** `myapp`（ELF バイナリ）

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

## ビルド性能

### 標準的なビルド時間

| フェーズ | 時間 | 備考 |
| --- | --- | --- |
| 解析 | &lt;1秒 | Go コードのスキャン |
| バインディング生成 | &lt;1秒 | TypeScript の生成 |
| フロントエンドのビルド | 5-30秒 | プロジェクトの規模による |
| Go のコンパイル | 2-10秒 | コードの規模による |
| アセットの埋め込み | &lt;1秒 | フロントエンドの埋め込み |
| **合計** | **10-45秒** | 初回ビルド |
| **差分ビルド** | **5-15秒** | 2回目以降のビルド |

### ビルドの高速化

**1. ビルドキャッシュを使用する:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. 必要なものだけを実行する:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. 並列ビルド（複数マシン／CI）:**

v3 での Linux／Windows／macOS 間のクロスコンパイルは、通常、Docker の `wails-cross` コンテナ内、またはプラットフォームごとの専用ランナー上で行います。`wails3 build` 自体のターゲットはホスト OS です。サポートされるワークフローについては、[クロスプラットフォームビルド](/guides/build/cross-platform/)を参照してください。

**4. より高速なツールを使用する:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## トラブルシューティング

### ビルドが失敗する

**症状:** `wails3 build` がエラーで終了する

**一般的な原因:**

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


3. **依存関係の不足**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### バイナリが大きすぎる

**症状:** バイナリが 50MB を超えている

**解決策:**

1. **デバッグシンボルを削除する**（同梱の Taskfile では、すでに `go build` に `-ldflags="-s -w"` を渡しています）。

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

**症状:** ビルドに 1 分を超える時間がかかる

**解決策:**

1. **ビルドキャッシュを使用する**
  - Go のキャッシュは自動的に使用されます
  - フロントエンドのキャッシュ（Vite）は自動的に使用されます


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

- **開発中は `wails3 dev` を使用する** - 素早く反復できる
- **リリースには `wails3 build` を使用する** - 最適化された出力を生成できる
- **ビルドにバージョンを付ける** - バージョンの埋め込みには `-ldflags` を使用する
- **ターゲットプラットフォーム上でビルドをテストする** - クロスコンパイルは完全ではない
- **フロントエンドのビルドを高速に保つ** - バンドラーの設定を最適化する
- **ビルドキャッシュを使用する** - 2回目以降のビルドを高速化できる

### ❌ 禁止事項

- **`build/` ディレクトリをコミットしない** - `.gitignore` に追加する
- **ビルドのテストを省略しない** - リリース前に必ずテストする
- **不要なアセットを埋め込まない** - バイナリを小さく保つ
- **本番環境でデバッグビルドを使用しない** - 最適化されたビルドを使用する
- **コード署名を忘れない** - 配布には必須

## 次のステップ

**アプリケーションのビルド** - ビルドとパッケージ化の詳細ガイド [詳細を見る →](/guides/build/building/)

**クロスプラットフォームビルド** - 1台のマシンですべてのプラットフォーム向けにビルド [詳細を見る →](/guides/build/cross-platform/)

**インストーラーの作成** - エンドユーザー向けのインストーラーを作成 [詳細を見る →](/guides/installers/)

---

**ビルドについて質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[ビルド例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/build)を確認してください。
