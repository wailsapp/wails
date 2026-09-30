---
title: "iOS"
description: "iOS で Wails アプリケーションをビルドして実行する方法 — セットアップ、シミュレーター、実機ビルド、設定、ネイティブ機能"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="試験的機能"}
iOS サポートは試験段階であり、今後のリリースで変更される可能性があります。

@end

@note{type="tip"}
Wails を使ったモバイル開発が初めての場合は、まず [初めてのモバイルアプリ →](/guides/mobile/first-mobile-app/) の手順に沿って進めてから、完全なリファレンスとしてこのページに戻ってください。

@end

Wails v3 アプリは、完全なネイティブアプリとして iOS 上で動作します。さらに、デスクトップ版と<em>まったく同じ</em>ように動作します。同じ Go バックエンドと同じフロントエンドを使用し、`@wailsio/runtime`も同じです。サービスバインディング、イベント、ダイアログ、クリップボードはすべて同一に動作し、モバイル固有の再配線は<strong>一切</strong>必要ありません。モバイル専用のコードベースも、移植レイヤーも、習得すべき特別な API もありません。既存の Wails アプリがそのまま iOS 上で動作します。移植は文字どおりシームレスです。アプリを変更せずに持ち込み、そのままリリースできます。

同じ `main.go` からデスクトップ向けと iOS 向けの両方をビルドできます。iOS 固有の調整は `application.Options.IOS` で設定します。

## 要件

- **完全版の Xcode** がインストールされた macOS（コマンドラインツールだけでは不十分です）— `wails3 doctor` には検出された iOS SDK が表示されます
- Go 1.25 以降と npm

## シミュレーター

プロジェクトディレクトリで次を実行します。

```bash
wails3 task ios:run
```

アプリがビルドされ、シミュレーターがまだ実行されていない場合は起動したうえで、アプリが起動します。

併せて役立つコマンド：

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

デバッグビルドでは、Safari の「開発」メニューから WebView を検査できます。

## パッケージング

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

これらは最適化およびストリップ済みの本番ビルドです。

## 実機ビルド

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` は実機向けにビルドします。エンタイトルメントは `build/ios/entitlements.plist` から取得され、実機ビルドにのみ適用されます。アプリに必要なケイパビリティキーを追加してください。

@note{type="tip"}
署名、プロビジョニング、App Store アーカイブを自動管理するには、`wails3 task ios:xcode` を使って生成済みの Xcode プロジェクトを開き、代わりに Xcode からビルドしてください。

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

起動オプション（`application.Options.IOS`）には、`DisableScroll`、`DisableBounce`、`DisableScrollIndicators`、`DisableInputAccessoryView`、`EnableBackForwardNavigationGestures`、`DisableLinkPreview`、`EnableInlineMediaPlayback`、`EnableAutoplayWithoutUserAction`、`DisableInspectable`、`UserAgent`、`ApplicationNameForUserAgent`、`BackgroundColour`、および `EnableNativeTabs` + `NativeTabsItems` によるネイティブの下部タブが含まれます。

## ネイティブ機能

iOS 固有の機能は `application.IOS` を通じて利用できます。共有コードをプラットフォーム非依存に保つため、`//go:build ios` ファイル内の Go から呼び出します。Android では、`application.Android` を通じて同じ機能一式を利用できます。

単発のアクションは即座に戻ります。

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

クエリヘルパー（`SafeAreaJSON()`、`AppInfoJSON()`、`PowerJSON()`、`NetworkJSON()`、`StorageJSON()`、`GetOrientation()`、`GetBrightness()`）は、結果を JSON として返します。`StoragePath()` はアプリの Application Support ディレクトリへの絶対パスを返します。このディレクトリは、データベースやその他の永続ファイルの保存先に適しています（Android の `getFilesDir()` に相当します）。ディレクトリは初回アクセス時に作成されます。作成できない場合、`StoragePath()` は空文字列を返すため、使用前に `""` を確認してください。

### イベント

権限プロンプト、センサーストリーム、カメラ撮影など、後から完了する処理は、戻り値ではなく<strong>イベント</strong>として結果を通知します。このイベントは Go またはフロントエンドでリッスンできます。Android と共有する機能の名前には `common:`、iOS 専用機能の名前には `ios:` がプレフィックスとして付きます。

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

`v3/examples/mobile` にある kitchen-sink サンプルでは、上記のすべての機能がエンドツーエンドで接続されています。

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

同梱の `@wailsio/runtime` は、フロントエンド向けの小規模な <strong>iOS 名前空間</strong>も公開します。

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

ネイティブの下部タブで選択が行われると、`window` に `nativeTabSelected` イベントが送られます。

## サポート状況

| 領域 | 状況 |
| --- | --- |
| フロントエンドのレンダリングとアセット | ✅ |
| サービスバインディング、イベント（双方向） | ✅ |
| メッセージダイアログ | ✅ |
| ファイル（単一／複数）／ディレクトリを開くダイアログ | ✅ サンドボックス内のコピーとしてインポート |
| ファイル保存ダイアログ | ❌ 代わりにアプリのサンドボックス内へ書き込み |
| クリップボード | ✅ |
| 画面 API | ✅ セーフエリア内の作業領域を含む |
| ライフサイクルイベント | ✅ |
| ウィンドウの位置とサイズ、メニュー、システムトレイ | iOS では何も行いません |
| 複数ウィンドウ | 最初のウィンドウのみ表示 |

## 移植時の注意事項

- デスクトップ向けコードは変更せずに iOS 用としてビルドできます。ウィンドウ、メニュー、システムトレイの呼び出しは単に何も行いません。
- ファイル保存ダイアログの代わりに、アプリのサンドボックスへ書き込んでから共有してください。
- フロントエンドはレスポンシブに設計してください。セーフエリアは自動的に処理されます。
