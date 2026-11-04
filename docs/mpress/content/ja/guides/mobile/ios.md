---
title: "iOS"
description: "iOS で Wails アプリケーションをビルドして実行する方法 — セットアップ、シミュレータ、実機向けビルド、設定、ネイティブ機能"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="実験的機能"}
iOS サポートは実験的な機能であり、今後のリリースで変更される可能性があります。

@end

@note{type="tip"}
Wails を使ったモバイル開発が初めての場合は、まず [初めてのモバイルアプリ →](/guides/mobile/first-mobile-app/) の手順に沿って進めてから、完全なリファレンスとしてこのページに戻ってください。

@end

Wails v3 アプリは、完全なネイティブアプリとして iOS 上で動作します。さらに、デスクトップ版と*まったく同じ*ように動作します。同じ Go バックエンド、同じフロントエンド、同じ `@wailsio/runtime` を使用し、サービスバインディング、イベント、ダイアログ、クリップボードもすべて同じように動作するため、モバイル固有の再配線は**一切**不要です。モバイル専用のコードベースも、移植レイヤーも、習得すべき特別な API もありません。既存の Wails アプリがそのまま iOS 上で動作します。移植は文字どおりシームレスです。アプリをそのまま移してリリースできます。

同じ `main.go` でデスクトップ向けと iOS 向けの両方をビルドします。iOS 固有の調整は `application.Options.IOS` で設定します。

## 要件

- **完全版の Xcode** がインストールされた macOS（コマンドラインツールだけでは不十分です）— `wails3 doctor` で検出可能な iOS SDK が表示されます
- Go 1.26 以降および npm

## シミュレータ

プロジェクトディレクトリで次を実行します。

```bash
wails3 task ios:run
```

これによりアプリがビルドされ、シミュレータがまだ起動していなければ起動され、その上でアプリが立ち上がります。

併せて役立つコマンド：

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

デバッグビルドでは、Safari の「開発」メニューから WebView を検証できます。

## パッケージ化

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

これらは最適化され、不要なシンボルを削除した本番用ビルドです。

## 実機向けビルド

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` は実機向けにビルドします。エンタイトルメントは `build/ios/entitlements.plist` から取得され、実機向けビルドにのみ適用されます。アプリに必要な機能のキーを追加してください。

@note{type="tip"}
署名、プロビジョニング、App Store アーカイブを自動管理するには、`wails3 task ios:xcode` で生成済みの Xcode プロジェクトを開き、代わりに Xcode からビルドしてください。

@end

## 設定

`build/config.yml`：

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

起動オプション（`application.Options.IOS`）には、`DisableScroll`、`DisableBounce`、`DisableScrollIndicators`、`DisableInputAccessoryView`、`EnableBackForwardNavigationGestures`、`DisableLinkPreview`、`EnableInlineMediaPlayback`、`EnableAutoplayWithoutUserAction`、`DisableInspectable`、`UserAgent`、`ApplicationNameForUserAgent`、`BackgroundColour`、および `EnableNativeTabs` と `NativeTabsItems` によるネイティブの下部タブが含まれます。

## ネイティブ機能

iOS 固有の機能は `application.IOS` を通じて利用できます。共有コードをプラットフォーム非依存に保つため、`//go:build ios` ファイル内の Go コードから呼び出します。Android では、`application.Android` を通じて同じ機能セットを利用できます。

一度きりの処理は即座に制御を返します。

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

クエリ用ヘルパーは、`SafeAreaJSON()`、`AppInfoJSON()`、`PowerJSON()`、`NetworkJSON()`、`StorageJSON()`、`GetOrientation()`、`GetBrightness()` の結果を JSON として返します。`StoragePath()` は、アプリの Application Support ディレクトリへの絶対パスを返します。このディレクトリは、データベースなどの永続ファイルの保存先に適しています（Android の `getFilesDir()` に相当します）。ディレクトリは初回アクセス時に作成されます。作成できなかった場合、`StoragePath()` は空文字列を返すため、使用前に `""` かどうかを確認してください。

### イベント

権限プロンプト、センサーストリーム、カメラ撮影など、後から完了する処理は、戻り値ではなく**イベント**として結果を通知します。このイベントは Go またはフロントエンドでリッスンできます。Android と共通の機能には名前の先頭に `common:` が付き、iOS 専用機能には `ios:` が付きます。

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| イベント | トリガー | ペイロード |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

`v3/examples/mobile` にある全機能を網羅したサンプルでは、上記の各機能がエンドツーエンドで接続されています。

## WebView の制御

一部の WebView の動作は、実行時に Go から変更することもできます。

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

同梱の `@wailsio/runtime` は、フロントエンド**向け**の小規模な iOS 名前空間も公開しています。

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

ネイティブの下部タブの選択は、`window` に `nativeTabSelected` イベントとして届きます。

## サポート状況

| 領域 | 状況 |
| --- | --- |
| フロントエンドのレンダリングとアセット | ✅ |
| サービスバインディング、イベント（双方向） | ✅ |
| メッセージダイアログ | ✅ |
| ファイル／複数ファイル／ディレクトリを開くダイアログ | ✅ サンドボックス内のコピーとしてインポート |
| ファイル保存ダイアログ | ❌ 代わりにアプリのサンドボックス内へ書き込む |
| クリップボード | ✅ |
| Screens API | ✅ セーフエリアを考慮した作業領域を含む |
| ライフサイクルイベント | ✅ |
| ウィンドウのジオメトリ、メニュー、システムトレイ | iOS では何も行いません |
| 複数ウィンドウ | 最初のウィンドウだけが表示されます |

## 移植時の注意事項

- デスクトップ向けコードは変更せずにiOS用としてビルドできます。ウィンドウ、メニュー、システムトレイの呼び出しは何も行いません。
- ファイル保存ダイアログを、アプリのサンドボックスへの書き込みと共有に置き換えてください。
- フロントエンドはレスポンシブに設計してください。セーフエリアは自動的に処理されます。
