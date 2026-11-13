---
title: "モバイル概要"
description: "デスクトップアプリと同じ Go コードベースから iOS および Android アプリをビルドする"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 は、デスクトップ向けに作成済みのものと同じ `main.go` およびフロントエンドを使用して、**iOS と Android** で動作します。別個のモバイルプロジェクトも、コード共有ブリッジも、書き直しも必要ありません。Go バイナリはモバイルターゲット向けにコンパイルされ、ネイティブ WebView が既存のフロントエンドをレンダリングします。

@cards{cols="2"}
iOS
WKWebView + UIKit ホスト。アセットはカスタム `wails://` スキーム経由で配信されるため、ポートは開きません。完全版 Xcode を備えた **macOS** が必要です。

[iOS ガイド →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`。Go は NDK により `libwails.so` としてコンパイルされます。macOS、Linux、Windows で動作します。

[Android ガイド →](/guides/mobile/android/)

@end

## 動作を確認する：Kitchen Sink の例

何が可能かを理解するには、**Kitchen Sink** を見るのが最適です。これは、1 つのコードベースから iOS、Android、デスクトップで同じように動作する単一の Wails アプリです。

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/mobile" description="バインディング · イベント · ダイアログ · ハプティクス · 位置情報 · 生体認証 · 通知 · セキュアストレージ · その他 — すべて 1 つの main.go から利用可能"}
7 個のタブを通じて、主要なモバイル API サーフェスをすべて実演しており、デスクトップでも動作します。デスクトップでは、フロントエンドのプラットフォームチェックにより **Mobile** タブと **Hardware** タブが非表示になります。デスクトップ向けにビルドした場合、Go 側では `common:*` モバイルイベント用のハンドラーを登録しません。これは、単一のコードベースをすべての環境に提供するための推奨パターンです。

| タブ | プラットフォーム | 実演内容 |
| --- | --- | --- |
| **バインディング** | すべて | 値、構造体、エラーを返す JS → Go サービス呼び出し |
| **イベント** | すべて | Go → JS クロック、JS → Go → JS の ping/pong、OS システムイベント（バッテリー、ネットワーク、テーマ） |
| **ダイアログ** | すべて | 各プラットフォームのネイティブメッセージダイアログ |
| **システム** | すべて | クリップボード、画面メトリクス、デバイス情報 |
| **モバイル** | iOS + Android | 共有シート、スリープ防止、ライト、明るさ、生体認証、ローカル通知、セキュアストレージ |
| **ハードウェア** | iOS + Android | ハプティクス、位置情報、加速度センサー、近接センサー、テキスト読み上げ |
| **ネイティブ** | iOS + Android | iOS：ハプティクス + WKWebView の切り替え · Android：バイブレーション + トースト |

自分で実行するには：

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## 仕組み

すべてのプラットフォームで同じアプリケーションモデルが適用されます。

1. **Go バックエンド** — サービス、イベントハンドラー、アプリケーションロジックは、変更せずに `GOOS=ios` および `GOOS=android` 向けにコンパイルされます。
2. **フロントエンド** — まったく同じ HTML/JS/CSS を使用します。`@wailsio/runtime` パッケージも同じように動作し、サービスバインディング、イベント、ダイアログ、クリップボードはすべて同じプロセス内トランスポートを経由します。
3. **WebView ホスト** — iOS では `UIViewController` 内の `WKWebView`、Android では `Activity` 内の `WebView` です。Wails がメッセージブリッジを自動的に接続します。
4. **プロセス内アセット配信** — アセットは localhost サーバーではなく、Go のメモリから直接配信されます。開くポートも、ループバックも、追加のレイテンシーもありません。

プラットフォーム固有の動作は、`//go:build ios` または `//go:build android` で保護されたファイルに配置されるため、共有コードをクリーンに保てます。

## 前提条件の概要

| 要件 | iOS | Android |
| --- | --- | --- |
| オペレーティングシステム | macOS のみ | macOS、Linux、Windows |
| ツールチェーン | 完全版 Xcode（CLI ツールのみでは不可） | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| 確認コマンド | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
ツールチェーンをセットアップした後に `wails3 doctor` を実行してください。各プラットフォームで検出されたものと不足しているものが正確に表示されます。

@end

## サポートされる機能

両プラットフォームで、同じ中核機能セットを利用できます。

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
| ハプティクス | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| デバイス情報 | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| ネイティブタブ（iOS） | ✅ UITabBar | — |
| トーストメッセージ（Android） | — | ✅ `Android.Toast.Show` |
| 複数ウィンドウ | ❌ 最初のウィンドウのみ | ❌ 最初のウィンドウのみ |
| ウィンドウのジオメトリ／メニュー／トレイ | 意図的に何も行わない | 意図的に何も行わない |

## ビルドタグの規則

プラットフォーム条件付きコードを記述する際には、次の2つの重要な規則を理解しておく必要があります。

- **`ios` は `darwin`** を暗黙に含むため、`//go:build darwin` タグを付けたファイルは iOS 向けにもコンパイルされます。macOS のみを対象にするには、`//go:build darwin && !ios` を使用します。
- **`android` は `linux`** を暗黙に含むため、`//go:build linux` タグを付けたファイルは Android 向けにもコンパイルされます。デスクトップ版 Linux のみを対象にするには、`//go:build linux && !android` を使用します。

実行時には、`runtime.GOOS` はそれぞれ `"ios"` と `"android"` を返します。

## 実行時のプラットフォーム検出

ビルドタグは、特定のプラットフォームでのみ*コンパイル*できるコードに使用します。共有コード内で通常の条件分岐を行う場合は、`application.System` を使用します。これはすべてのビルドで利用でき（ビルドタグは不要）、同じファイルがどのプラットフォームでも動作します。

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

利用可能な値は、`IsMobile()`、`IsDesktop()`、`IsServer()`（`server` ビルドタグ）、`IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)` です。

フロントエンドには、`@wailsio/runtime` に対応するヘルパーがあります。

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## 次のステップ

@cards{cols="2"}
🚀 初めてのモバイルアプリ
デスクトップ向け Wails アプリを、わずか数分で iOS Simulator または Android Emulator 上で実行します。

[始める →](/guides/mobile/first-mobile-app/)

---
iOS ガイド
iOS ツールチェーンの完全なセットアップ、シミュレーター、デバイス向けビルド、署名、構成、API リファレンスについて説明します。

[iOS ガイド →](/guides/mobile/ios/)

---
Android ガイド
Android SDK/NDK の完全なセットアップ、エミュレーター、APK の署名、Play Store 向けパッケージング、API リファレンスについて説明します。

[Android ガイド →](/guides/mobile/android/)

@end
