---
title: "Android"
description: "Android で Wails アプリケーションをビルドして実行する方法 — ツールチェーンのセットアップ、エミュレーター、APK の署名、Play ストア向けパッケージング、API リファレンス"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="試験的機能"}
Android サポートは試験的な機能であり、今後のリリースで変更される可能性があります。

@end

@note{type="tip"}
Wails で初めてモバイルアプリを開発する場合は、まず[初めてのモバイルアプリ →](/guides/mobile/first-mobile-app/)の手順に沿って進めてから、完全なリファレンスとしてこのページに戻ってください。

@end

Wails v3アプリケーションは、Android 上でネイティブアプリとして動作します。Android の`WebView`がフロントエンドをレンダリングし、アセットは Go アセットサーバーを基盤とする`WebViewAssetLoader`を通じて**プロセス内**で提供されます（localhost サーバーも開放ポートも使用しません）。標準の`@wailsio/runtime`は変更なしで動作し、サービスバインディング、イベント、ダイアログ、クリップボードは Go メッセージプロセッサーを経由します。

同じ`main.go`からデスクトップ向けと Android 向けの両方をビルドできます。Go コードは C 共有ライブラリとしてコンパイルされ（`libwails.so`、`GOOS=android`と NDK ツールチェーンを使用）、小さな Java ホストによって読み込まれます。Android 固有の動作は、`//go:build android`で保護されたプラットフォーム別の Go ファイルに実装されています。

## 要件

- platform-tools、SDK プラットフォーム（API 35）、build-tools、**NDK**（26.3.x）を含む**Android SDK** — 検出された内容は`wails3 doctor`で確認できます
- Gradle 用の**JDK**（例：OpenJDK 21）。`java`が`PATH`に含まれていない場合は、`JAVA_HOME`を設定してください
- Go 1.26以降、および npm
- SDK を指す`ANDROID_HOME`（または`ANDROID_SDK_ROOT`）

コマンドラインツールを使用して SDK の各コンポーネントをインストールします。

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## エミュレーターでの実行

プロジェクトディレクトリで次を実行します。

```bash
wails3 task android:run
```

実行中のエミュレーターがなければ起動し、バインディングを生成してフロントエンドをビルドし、エミュレーターの ABI 向けに Go コードを`libwails.so`へコンパイルします。その後、Gradle でデバッグ APK を作成し、インストールして起動します。

併せて役立つコマンド：

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

本番ビルドでは`-tags production,android`を使用し、シンボルを削除して、フレームワーク内部の診断機能をコンパイル時に除外します。`wails3 task android:package:fat`は、`arm64-v8a`と`x86_64`の両方を単一の APK にビルドします。

Google Play では、新しいアプリの申請に Android App Bundle（`.aab`）形式が必須です。また、新規申請では Android 15（API 35）以降をターゲットにする必要があります。プロジェクトテンプレートでは、`build/android/app/build.gradle`内の`compileSdk`と`targetSdk`が35に設定されています。`wails3 task android:bundle:fat`は、両方の ABI を含む`bin/<AppName>.aab`を生成します。Google Play はそこからデバイスごとに最適化された APK を生成するため、ストアへのアップロードには両 ABI を含むバンドルが適切な成果物です。`.aab`は`adb`で直接インストールできないため、ローカルおよびエミュレーターでのテストには引き続き APK が最も迅速です。

`android:run`と`android:deploy-emulator`は、エミュレーター向けのタスクです。物理 Android デバイスでは、デバッグ APK には`android:run:device`、リリース APK には`android:deploy-device`を使用してください。どちらも`arm64`向けにビルドし、`adb devices`に表示される接続済みの非エミュレーターのうち最初のエントリを選択してインストールし、`com.wails.app.MainActivity`を起動します。特定のデバイスを対象にするには、`DEVICE_ID=<serial>`を渡してください。

## 署名とリリースビルド

キーストアがない場合、リリースビルドはテスト用にインストールできるよう、Android の**デバッグ**キーストアで署名されます。独自のキーストアで署名するには、次を設定します。

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

同じ変数で App Bundle にも署名できます。これらを設定して`wails3 task android:bundle:fat`を実行すると、Google Play に提出可能な`.aab`が生成されます。設定しない場合、バンドルはデバッグキーストアで署名され、Google Play に拒否されるため、タスクによって警告が表示されます。

@note{type="tip"}
[Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756)を使用する場合、ローカルで署名に使うキーストアは**アップロード鍵**です。Google はこの鍵でアップロードを検証した後、管理しているアプリ署名鍵でアプリに再署名します。また、Google Play ではアップロードのたびに、より大きい`versionCode`が必要です。`build/android/app/build.gradle`で値を増やしてください。

@end

## 設定

フロントエンドは実行時に`Android`ランタイムオブジェクト（`Android.Haptics.Vibrate(durationMs)`、`Android.Device.Info()`、`Android.Toast.Show(message)`）を通じて Android の機能を操作します。パッケージ名は、ビルドタスク内の`APP_ID`で制御します。

## 動作する機能と動作しない機能

| 領域 | 状態 |
| --- | --- |
| WebView とプロセス内アセット（`WebViewAssetLoader`） | ✅ |
| サービスバインディング、イベント（双方向） | ✅ |
| メッセージダイアログ | ✅ ボタンコールバック付き AlertDialog |
| ファイルを開く／複数ファイルを開くダイアログ | ✅ Storage Access Framework（ファイルはキャッシュコピーとしてインポート） |
| ディレクトリを開く／ファイルを保存するダイアログ | ❌ エラーを返します — 代わりにアプリのサンドボックス内へ書き込んでください |
| クリップボード | ✅ ClipboardManager |
| Screens API | ✅ システムバーを除いた作業領域を含む WindowMetrics |
| ライフサイクルイベント（`events.Android.*`） | ✅ |
| ハプティクス、デバイス情報、トースト | ✅ `Android.*`ランタイム API |
| エミュレーターおよび物理デバイス向けビルド | ✅ `android:run`、`android:run:device`、`android:deploy-emulator`、`android:deploy-device` |
| ウィンドウジオメトリ、メニュー、システムトレイ | 意図的に何もしない操作 |
| 複数ウィンドウ | 最初のウィンドウのみ表示されます |

## 移植時の注意事項

- デスクトップコードは`GOOS=android`でも変更なしでコンパイルできます。Android アプリは全画面表示されるため、ジオメトリ、メニュー、トレイに対する呼び出しは何もしない操作になります。
- `android`では、`linux`ビルドタグ**である**も暗黙的に有効になります（Android は Linux カーネルを使用します）。デスクトップ Linux 専用ファイルには`//go:build linux && !android`が必要で、実行時の`runtime.GOOS`は`"android"`です。
- ファイル保存ダイアログとディレクトリ選択ダイアログは、アプリのサンドボックスへの書き込みと、インテントによる共有フローに置き換えてください。ファイルを開くダイアログは動作し、選択したドキュメントをキャッシュディレクトリのコピーとしてインポートするため、実際のファイルシステムパスを取得できます。
- 実際のアプリは必ず`CGO_ENABLED=1`と NDK を使用してビルドされます。非 cgo パスは、`wails3 generate bindings`などのツールがパッケージを読み込めるようにするためだけに存在します。
- フロントエンドはレスポンシブに設計してください。`Screens`の作業領域には、ステータスバーとナビゲーションバーは含まれません。
