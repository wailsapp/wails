---
title: "Android"
description: "Android で Wails アプリケーションをビルドして実行する方法 — ツールチェーンのセットアップ、エミュレーター、APK の署名、Play ストア向けパッケージング、API リファレンス"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="実験的機能"}
Android サポートは実験的な機能であり、今後のリリースで変更される可能性があります。

@end

@note{type="tip"}
Wails で初めてモバイル開発を行う場合は、まず[初めてのモバイルアプリ →](/guides/mobile/first-mobile-app/)の手順に沿って進めてから、完全なリファレンスとしてこのページに戻ってください。

@end

Wails v3 アプリケーションは、Android 上でネイティブアプリとして動作します。Android の`WebView`がフロントエンドをレンダリングし、アセットは Go アセットサーバーを基盤とする`WebViewAssetLoader`を通じて<strong>インプロセス</strong>で配信されます（localhost サーバーも開放ポートも使用しません）。また、標準の`@wailsio/runtime`は変更なしで動作し、サービスバインディング、イベント、ダイアログ、クリップボードは Go メッセージプロセッサーを経由します。

同じ`main.go`からデスクトップ向けと Android 向けの両方をビルドできます。Go コードは C 共有ライブラリ（`libwails.so`、`GOOS=android` + NDK ツールチェーン）としてコンパイルされ、小さな Java ホストによって読み込まれます。Android 固有の動作は、`//go:build android`で保護されたプラットフォーム別の Go ファイルに実装されています。

## 要件

- platform-tools、SDK プラットフォーム（API 35）、build-tools、**NDK**（26.3.x）を含む<strong>Android SDK</strong> — 検出された項目は`wails3 doctor`で確認できます
- Gradle 用の<strong>JDK</strong>（例：OpenJDK 21）。`java`が`PATH`に含まれていない場合は、`JAVA_HOME`を設定してください
- Go 1.25以降と npm
- SDK を指す`ANDROID_HOME`（または`ANDROID_SDK_ROOT`）

コマンドラインツールを使用して、必要な SDK コンポーネントをインストールします：

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## エミュレーターでの実行

プロジェクトディレクトリで次を実行します：

```bash
wails3 task android:run
```

エミュレーターが実行されていない場合は起動し、バインディングの生成、フロントエンドのビルド、エミュレーターの ABI 向けの Go コードの`libwails.so`へのコンパイル、Gradle によるデバッグ APK の作成を順に行ってから、その APK をインストールして起動します。

併用すると便利なコマンド：

```bash
wails3 task android:logs    # stream the app's logcat output
```

デバッグビルドでは、Chrome の`chrome://inspect`から WebView を検査できます。

## パッケージング

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

本番ビルドでは`-tags production,android`を使用し、シンボルを除去するとともに、フレームワークの内部診断機能をコンパイル対象から除外します。`wails3 task android:package:fat`は`arm64-v8a`と`x86_64`の両方をビルドし、1 つの APK にまとめます。

Google Play では、新しいアプリの申請に Android App Bundle（`.aab`）形式が必要です。また、新規申請では Android 15（API 35）以降を対象にする必要があります。プロジェクトテンプレートでは、`build/android/app/build.gradle`内の`compileSdk`と`targetSdk`が35に設定されています。`wails3 task android:bundle:fat`を実行すると、両方の ABI を含む`bin/<AppName>.aab`が生成されます。Google Play はそこからデバイスごとに最適化された APK を生成するため、このファットバンドルがストアへのアップロードに適した成果物です。`.aab`は`adb`で直接インストールできないため、ローカルテストやエミュレーターでのテストには引き続き APK が最も手軽です。

`android:run`と`android:deploy-emulator`はエミュレーター向けのタスクです。物理 Android デバイスでは、デバッグ APK に`android:run:device`、リリース APK に`android:deploy-device`を使用してください。どちらも`arm64`向けにビルドし、`adb devices`に表示される接続済みの非エミュレーターエントリのうち最初のものを選択してインストールし、`com.wails.app.MainActivity`を起動します。特定のデバイスを対象にするには、`DEVICE_ID=<serial>`を渡してください。

## 署名とリリースビルド

キーストアがない場合、リリースビルドはテスト用にインストールできるよう、Android の<strong>デバッグ</strong>キーストアで署名されます。独自のキーストアで署名するには、次を設定します：

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

同じ変数を使用して App Bundle にも署名できます。これらを設定して`wails3 task android:bundle:fat`を実行すると、Play に提出可能な`.aab`が生成されます。設定しない場合、バンドルはデバッグキーストアで署名され、Google Play に拒否されるため、タスクは警告を表示します。

@note{type="tip"}
[Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756)を使用する場合、ローカルでの署名に使うキーストアが<strong>アップロード鍵</strong>になります。Google はこの鍵でアップロードを検証してから、管理しているアプリ署名鍵でアプリに再署名します。また、Google Play ではアップロードのたびに、より大きな`versionCode`が必要です。`build/android/app/build.gradle`で値を増やしてください。

@end

## 設定

フロントエンドは実行時に`Android`ランタイムオブジェクトを介して、`Android.Haptics.Vibrate(durationMs)`、`Android.Device.Info()`、`Android.Toast.Show(message)`などの Android 機能を操作します。パッケージ名は、ビルドタスク内の`APP_ID`で指定します。

## 動作する機能と動作しない機能

| 領域 | 状態 |
| --- | --- |
| WebView + インプロセスアセット（`WebViewAssetLoader`） | ✅ |
| サービスバインディング、イベント（双方向） | ✅ |
| メッセージダイアログ | ✅ ボタンコールバックに対応した AlertDialog |
| ファイルを開く／複数ファイルを開くダイアログ | ✅ Storage Access Framework（ファイルはキャッシュのコピーとしてインポート） |
| ディレクトリを開く／ファイルを保存するダイアログ | ❌ エラーを返す — 代わりにアプリのサンドボックス内へ書き込んでください |
| クリップボード | ✅ ClipboardManager |
| 画面 API | ✅ システムバーを除いた作業領域を含む WindowMetrics |
| ライフサイクルイベント（`events.Android.*`） | ✅ |
| ハプティクス、デバイス情報、トースト | ✅ `Android.*`ランタイム API |
| エミュレーターおよび物理デバイス向けビルド | ✅ `android:run`、`android:run:device`、`android:deploy-emulator`、`android:deploy-device` |
| ウィンドウジオメトリ、メニュー、システムトレイ | 意図的に何もしない処理 |
| 複数ウィンドウ | 最初のウィンドウのみ表示 |

## 移植時の注意事項

- デスクトップ用コードは`GOOS=android`でも変更せずにコンパイルできます。Android アプリはフルスクリーンで動作するため、ジオメトリ、メニュー、トレイに対する呼び出しは何もしない処理になります。
- `android`は`linux`ビルドタグ<strong>を</strong>暗黙的に有効にします（Android は Linux カーネルを使用しています）。デスクトップ Linux 専用のファイルには`//go:build linux && !android`が必要であり、実行時の`runtime.GOOS`は`"android"`になります。
- ファイル保存ダイアログとディレクトリ選択ダイアログは、アプリのサンドボックスへの書き込みと、インテントを使った共有フローに置き換えてください。ファイルを開くダイアログは動作し、選択したドキュメントをキャッシュディレクトリ内のコピーとしてインポートするため、実際のファイルシステムパスを取得できます。
- 実際のアプリは常に `CGO_ENABLED=1` と NDK を使用してビルドします。cgo を使用しないパスは、`wails3 generate bindings` などのツールがパッケージを読み込めるようにするためだけに存在します。
- フロントエンドはレスポンシブに設計してください。`Screens` の作業領域には、ステータスバーとナビゲーションバーは含まれません。
