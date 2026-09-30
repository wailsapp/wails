---
title: "モバイル概要"
description: "デスクトップアプリと同じ Go コードベースから iOS／Android アプリをビルド"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 は、デスクトップ向けにすでに作成しているものと同じ `main.go` とフロントエンドを使用して、**iOS と Android** で動作します。モバイル専用の別プロジェクトも、コード共有用のブリッジも、書き直しも必要ありません。Go バイナリがモバイルターゲット向けにコンパイルされ、ネイティブ WebView が既存のフロントエンドをレンダリングします。

@cards{cols="2"}
iOS
WKWebView と UIKit ホストを使用します。アセットはカスタム `wails://` スキーム経由で配信されるため、ポートは開きません。 完全版の Xcode がインストールされた **macOS** が必要です。

[iOS ガイド →](/guides/mobile/ios/)

---
Android
Android WebView と `WebViewAssetLoader` を使用します。Go は NDK を介して `libwails.so` としてコンパイルされます。 macOS、Linux、Windows で動作します。

[Android ガイド →](/guides/mobile/android/)

@end

## 実際の動作を見る：Kitchen Sink の例

何ができるかを理解するには、**Kitchen Sink** を見るのが最適です。これは、1 つのコードベースから iOS、Android、デスクトップで同じように動作する単一の Wails アプリです。

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="バインディング · イベント · ダイアログ · 触覚フィードバック · 位置情報 · 生体認証 · 通知 · セキュアストレージ · その他 — すべて 1 つの main.go から利用可能"}
7 個のタブで主要なモバイル API サーフェスをすべて実演し、デスクトップでも動作します。**Mobile** タブと **Hardware** タブは、フロントエンドでプラットフォームを確認してデスクトップでは非表示にします。デスクトップ向けにビルドした場合、Go 側では `common:*` モバイルイベント用のハンドラーを登録しません。1 つのコードベースをあらゆる環境に提供するには、このパターンを推奨します。

| タブ | プラットフォーム | 実演内容 |
| --- | --- | --- |
| **バインディング** | すべて | 値、構造体、エラーを返す JS → Go サービス呼び出し |
| **イベント** | すべて | Go → JS の時計、JS → Go → JS の ping/pong、OS のシステムイベント（バッテリー、ネットワーク、テーマ） |
| **ダイアログ** | すべて | 各プラットフォームのネイティブメッセージダイアログ |
| **システム** | すべて | クリップボード、画面メトリクス、デバイス情報 |
| **モバイル** | iOS + Android | 共有シート、スリープ抑止、ライト、明るさ、生体認証、ローカル通知、セキュアストレージ |
| **ハードウェア** | iOS + Android | 触覚フィードバック、位置情報、加速度センサー、近接センサー、テキスト読み上げ |
| **ネイティブ** | iOS + Android | iOS：触覚フィードバック + WKWebView の切り替え · Android：バイブレーション + トースト |

自分で実行するには、次の手順に従います。

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 仕組み

すべてのプラットフォームで同じアプリケーションモデルが適用されます。

1. **Go バックエンド** — サービス、イベントハンドラー、アプリケーションロジックは、変更せずに `GOOS=ios` と `GOOS=android` 向けにコンパイルできます。
2. **フロントエンド** — まったく同じ HTML／JS／CSS を使用します。`@wailsio/runtime` パッケージも同じように動作し、サービスバインディング、イベント、ダイアログ、クリップボードはすべて同じプロセス内トランスポートを経由します。
3. **WebView ホスト** — iOS では `UIViewController` 内の `WKWebView`、Android では `Activity` 内の `WebView` を使用します。Wails がメッセージブリッジを自動的に接続します。
4. **プロセス内アセット配信** — アセットは localhost サーバーではなく、Go のメモリから直接配信されます。ポートの開放も、ループバックも、余分なレイテンシーもありません。

プラットフォーム固有の動作は、`//go:build ios` または `//go:build android` で保護されたファイルに配置されるため、共有コードを簡潔に保てます。

## 前提条件の概要

| 要件 | iOS | Android |
| --- | --- | --- |
| オペレーティングシステム | macOS のみ | macOS、Linux、Windows |
| ツールチェーン | 完全版の Xcode（CLI ツールのみでは不可） | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 確認コマンド | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
ツールチェーンのセットアップ後に `wails3 doctor` を実行してください。各プラットフォームについて、検出されたものと不足しているものが正確に表示されます。

@end

## サポートされる機能

両方のプラットフォームで、同じ中核機能セットを利用できます。

| 機能 | iOS | Android |
| --- | --- | --- |
| サービスバインディング（JS → Go） | ✅ | ✅ |
| イベント（双方向） | ✅ | ✅ |
| メッセージダイアログ | ✅ UIAlertController | ✅ AlertDialog |
| ファイルを開くダイアログ | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| ファイル保存ダイアログ | ❌ 代わりにサンドボックスへ書き込み | ❌ 代わりにサンドボックスへ書き込み |
| クリップボード | ✅ UIPasteboard | ✅ ClipboardManager |
| 画面／セーフエリアのメトリクス | ✅ | ✅ |
| ライフサイクルイベント | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| 触覚フィードバック | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| デバイス情報 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| ネイティブタブ（iOS） | ✅ UITabBar | — |
| トーストメッセージ（Android） | — | ✅ `Android.Toast.Show` |
| 複数ウィンドウ | ❌ 最初のウィンドウのみ | ❌ 最初のウィンドウのみ |
| ウィンドウの位置とサイズ／メニュー／トレイ | 意図的に何も処理しない | 意図的に何も処理しない |

## ビルドタグのルール

プラットフォーム条件付きコードを記述する際に知っておくべき重要なルールが2つあります。

- **`ios` は `darwin`** を暗黙的に含みます。そのため、`//go:build darwin` タグを付けたファイルはiOS向けにもコンパイルされます。macOSのみを対象にするには、`//go:build darwin && !ios` を使用します。
- **`android` は `linux`** を暗黙的に含みます。そのため、`//go:build linux` タグを付けたファイルはAndroid向けにもコンパイルされます。デスクトップ版Linuxのみを対象にするには、`//go:build linux && !android` を使用します。

実行時には、`runtime.GOOS` はそれぞれ `"ios"` と `"android"` を返します。

## 実行時のプラットフォーム判定

ビルドタグは、特定のプラットフォームでしか<em>コンパイル</em>できないコードに使用します。共有コード内で通常の条件分岐を行う場合は、`application.System` を使用します。これはすべてのビルドで利用できるため（ビルドタグは不要）、同じファイルがどのプラットフォームでも動作します。

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

利用可能なメソッドは、`IsMobile()`、`IsDesktop()`、`IsServer()`（`server` ビルドタグ）、および `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`です。

フロントエンドでは、`@wailsio/runtime` に対応するヘルパーが用意されています。

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 次のステップ

@cards{cols="2"}
🚀 初めてのモバイルアプリ
デスクトップ向けWailsアプリを、わずか数分でiOS SimulatorまたはAndroid Emulator上で実行できます。

[始める →](/guides/mobile/first-mobile-app/)

---
iOSガイド
iOSツールチェーンの完全なセットアップ、シミュレーター、デバイス向けビルド、署名、構成、およびAPIリファレンス。

[iOSガイド →](/guides/mobile/ios/)

---
Androidガイド
Android SDK/NDKの完全なセットアップ、エミュレーター、APK署名、Play Store向けパッケージング、およびAPIリファレンス。

[Androidガイド →](/guides/mobile/android/)

@end
