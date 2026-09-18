---
title: "初めてのモバイルアプリ"
description: "Wails アプリを数分で iOS Simulator または Android Emulator 上で実行する"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

このガイドでは、標準的な Wails デスクトップアプリを iOS Simulator または Android Emulator 上で実行します。<strong>Go コードを変更する必要はありません。</strong>すべてのターゲットで同じ `main.go` をビルドします。

<strong>所要時間：</strong>15～30 分（大部分は初回実行時のツールチェーンのインストールにかかる時間です）

## デスクトッププロジェクトから始める

まだプロジェクトがない場合は、新規作成します：

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

まず、デスクトップアプリが動作することを確認します：

```bash
wails3 dev
```

アプリが開いたら終了し、次に進みます。デスクトップで動作するものはすべてモバイルでも動作します。このガイドでは、`main.go`にも Go コードにも手を加える必要はありません。

---

## プラットフォームを選択する

@tabs{sync-key="mobile-platform"}
[iOS Simulator]
### 要件

- **macOS**（iOS のビルドは macOS でのみ行えます）
- **完全版の Xcode** — コマンドラインツールだけでは不十分です。App Store からインストールしてから、次を実行します：
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- <strong>Go 1.25 以降</strong>と **npm**（`wails3 init`を実行済みであれば、すでにインストールされています）

確認するには `wails3 doctor` を実行します。検出された iOS SDK の一覧が表示されます。

### Simulator 上で実行する

@steps
### アプリを起動する
```bash
wails3 task ios:run
```

これだけです。アプリがビルドされ、Simulator がまだ起動していなければ起動し、その上でアプリが起動します。

@note{type="tip"}
初回実行には数分かかります（iOS 用の Wails フレームワークをコンパイルしてキャッシュするためです）。2 回目以降は大幅に速くなります。

@end

起動すると、変更を加えていないデスクトップアプリが iOS Simulator 上で動作します。`main.go`もフロントエンドも同じです：

![iOS Simulator 上で動作するデフォルトの Wails アプリ](/assets/ios-simulator-first-app.png)

### ログをリアルタイム表示する
別のターミナルで次を実行します：

```bash
wails3 task ios:logs:dev
```

Simulator のログをアプリに絞り込み、リアルタイムで表示します。`fmt.Println`と `log.Println` の出力がここに表示されます。

### WebView を調査する
Safari で <strong>開発 → Simulator → アプリ</strong>を選択します。コンソール、デバッガー、ネットワークパネルなど、Web インスペクタの全機能を使用できます。

### 変更を加える
任意のフロントエンドファイル（`frontend/src/main.js`、`index.html`など）を編集し、`wails3 task ios:run`を再実行します。Wails がフロントエンドを再ビルドし、アプリを再起動します。

Go コードを変更した場合も、`wails3 task ios:run`を再実行します。Go の再コンパイルは増分方式のため、変更されたパッケージだけが再ビルドされます。

@end

### Xcode で開く（任意）

```bash
wails3 task ios:xcode
```

これにより、`build/ios/`が Xcode で開きます。Xcode は、実機へのデプロイ、高度なプロファイリング、プロビジョニングプロファイルの管理に使用できます。Wails はビルドのたびに Xcode プロジェクトを再生成するため、生成されたファイルを直接変更しないでください。

[Android Emulator]
### 要件

**Android SDK**、**NDK**、<strong>JDK</strong>が必要です。最も簡単な方法は Android Studio を使用することですが、コマンドラインツールも使用できます：

@steps
### Android コマンドラインツールをインストールする
[developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only)からダウンロードし、`~/android-sdk/cmdline-tools/latest/`に展開します。

### SDK コンポーネントをインストールする
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### エミュレーターを作成する
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### 環境変数を設定する
`~/.zshrc`または `~/.bashrc` に追加します：

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

再読み込みします：`source ~/.zshrc`

### JDK をインストールする
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

必要なものがすべて検出されることを確認するため、`wails3 doctor`を実行します。

### Emulator 上で実行する

@steps
### アプリを起動する
```bash
wails3 task android:run
```

初回実行時には、次の処理が行われます：

- 実行中のエミュレーターがなければ起動する
- バインディングを生成し、フロントエンドをビルドする
- NDK クロスコンパイラーを使用して Go コードをコンパイルし、`libwails.so` を生成する
- Gradle でデバッグ APK を組み立てる
- エミュレーターにインストールして起動する

@note{type="tip"}
初回ビルドでは Gradle をダウンロードし、NDK ツールチェーンをコンパイルするため、5～10 分かかります。2 回目以降は増分ビルドとなり、1 分以内で完了します。

@end

### ログをリアルタイム表示する
別のターミナルで次を実行します：

```bash
wails3 task android:logs
```

アプリに絞り込んで `adb logcat` を実行します。`fmt.Println` の出力がここに表示されます。

### WebView を調査する
Chrome を開き、`chrome://inspect`に移動します。アプリの WebView が **Remote Target** の下に表示されます。<strong>inspect</strong>をクリックして DevTools を開きます。

### 変更を加える
任意のファイルを編集し、`wails3 task android:run`を再実行します。Gradle の増分ビルドにより、変更されたコードだけが再コンパイルされます。

@end

@end

---

## 実行された処理を理解する

`main.go`には一切変更が加えられていません。すべての処理は Wails が行いました：

- **ビルドシステム** — プロジェクト内の`Taskfile.yml`には、プラットフォーム固有のツールチェーンを実行する`ios:*`タスクと`android:*`タスクが含まれています。
- **Go のクロスコンパイル** — 適切な`GOARCH`と sysroot を使用した`GOOS=ios`または`GOOS=android`。
- **ネイティブホスト** — コンパイル済みの Go コードを組み込み、WebView をホストする、生成済みの Xcode プロジェクト（iOS）または Gradle プロジェクト（Android）。
- **アセットの配信** — `frontend/dist/`は Go バイナリに埋め込まれ、同じプロセス内から配信されます。localhost サーバーは不要です。

---

## アプリをモバイル対応にする

アプリはすでに動作しますが、スマートフォンの画面ではデスクトップアプリのように見えます。少し変更するだけで、見栄えは大きく改善します。

### レスポンシブ CSS

モバイル画面は幅が狭く、入力方法も異なります。`frontend/public/style.css`（または同等のファイル）で次のように設定します。

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Go でプラットフォームを検出する

ビルドタグを使用すると、共有コードを煩雑にせずにプラットフォーム固有の動作を追加できます。

iOS 専用コード用に`mobile_ios.go`を作成します。

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Android 専用コード用に`mobile_android.go`を作成します。

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

デスクトップでも共有コードをコンパイルできるように、スタブとして`mobile_desktop.go`を作成します。

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### JavaScript でプラットフォームを検出し、モバイル専用 UI を制御する

`IOS.*`と`Android.*`のランタイムオブジェクトは、それぞれ対応するプラットフォームにのみ存在します。デスクトップで呼び出すと例外がスローされます。Kitchen Sink で使用されている適切なパターンは、プラットフォームを一度だけ検出し、モバイル専用コントロールを完全に非表示にすることです。

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

次に、HTML では次のように記述します。

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

この方法なら、モバイル専用ボタンはデスクトップでは一切レンダリングされず、個々の呼び出しを毎回`if (isMobile)`チェックで保護する必要もありません。

Go 側では、ビルドタグ付きのスタブと組み合わせ、必要なプラットフォームでのみイベントハンドラーが登録されるようにします。

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

これは[Kitchen Sink](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)で実際に使用されているパターンです。`native_features_stub.go`、`native_features_ios.go`、`native_features_android.go`を参照してください。

@note{type="note" title="ネイティブ機能 API とイベントの命名規則"}
知っておくべき規則が 2 つあります。

- **Go 側のネイティブ機能ではプラットフォームマネージャーを使用します。**`application.IOS.*`と`application.Android.*`のシングルトン経由で呼び出します。たとえば、`application.IOS.Haptic("medium")`や`application.Android.Share(payload)`です。各マネージャーは対応するプラットフォームにのみ存在するため、その呼び出しは`//go:build ios`ファイルまたは`//go:build android`ファイルに記述します。
- <strong>イベントは到達範囲に応じて名前空間で分けられます。</strong>両方のプラットフォームが認識するイベントには`common:*`プレフィックスを使用します（`common:haptic`、`common:location`など）。一方、一方のプラットフォームだけが生成または処理できるイベントには`ios:*`または`android:*`を使用します（例：`ios:backgroundTask`、`android:foregroundService`）。ほぼすべてのモバイル機能は共通しているため、フロントエンドでは`common:*`配下でイベントごとに 1 つのリスナーを維持します。

@end

### 触覚フィードバックを追加する（iOS）

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### バイブレーションを追加する（Android）

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## 本番向けにビルドする

@tabs{sync-key="mobile-platform"}
[iOS]
**シミュレーター向けビルド**（シミュレーターでのテスト用。署名は不要）：

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**実機向けビルド**（署名 ID とプロビジョニングプロファイルが必要）：

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**配布用 IPA**（App Store または TestFlight 用）：

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
App Store Connect にアップロードする場合は、`wails3 task ios:xcode`を使用し、署名とアーカイブを Xcode に管理させます。証明書、プロファイル、公証に伴う複雑な処理は Xcode が自動的に行います。

@end

[Android]
**デバッグ APK**（Android のデバッグ用キーストアで署名され、直接インストール可能）：

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**リリース APK**（独自のキーストアで署名）：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**ユニバーサル APK**（arm64 と x86_64 を 1 つのファイルに収録）：

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Play Store にアップロードする場合は、APK ではなく`.aab`（Android App Bundle）を生成します。Android Studio で`build/android/`を開き、<strong>Build → Generate Signed Bundle / APK</strong>を使用します。

@end

@end

---

## トラブルシューティング

### `wails3 task ios:run`が「no iOS SDKs found」で失敗する

完全版の Xcode をインストールし、選択しておく必要があります。

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run`が「SDK not found」で失敗する

`ANDROID_HOME`が設定され、エクスポートされていることを確認します。次のコマンドで検証してください。

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### シミュレーターが起動しない

利用可能なシミュレーターを一覧表示し、いずれかを手動で起動します。

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect`にターゲットが表示されない

WebView はデバッグモードである必要があります（`android:run`ではデフォルトです）。本番ビルドではなく、デバッグビルドを実行していることを確認してください。また、`adb devices`でエミュレーターが接続済みと表示されていることも確認します。

#### セーフエリアのインセットが適用されない

HTML に viewport メタタグが含まれていることを確認してください。

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Kitchen Sink を試す

最初のアプリが動作したら、ほかに何ができるかを学ぶには、**Kitchen Sink** のサンプルを試すのが最も近道です。これは、単一のコードベースから iOS、Android、デスクトップで動作する完全な Wails アプリで、触覚フィードバック、位置情報、生体認証、ローカル通知、セキュアストレージなどを網羅しています。

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

[`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile) でソースを参照してください。特に `native_features_ios.go` ファイルと `native_features_android.go` ファイルは、プラットフォーム固有の機能をコピー＆ペーストで実装する際の出発点として役立ちます。

## 次のステップ

@cards{cols="2"}
iOS ガイド
完全なリファレンス：設定オプション、ネイティブタブ、WKWebView の機能・設定の有効／無効を切り替える項目、実機向けビルド、署名。

[iOS ガイド →](/guides/mobile/ios/)

---
Android ガイド
完全なリファレンス：設定、トースト、Play Store 向けパッケージング、NDK の詳細。

[Android ガイド →](/guides/mobile/android/)

---
📖 Kitchen Sink のソース
触覚フィードバック、位置情報、生体認証、通知、セキュアストレージを、実行可能な 1 つのアプリにまとめています。

[GitHub で表示 →](https://github.com/wailsapp/wails/tree/master/v3/examples/mobile)

@end
