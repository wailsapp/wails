---
title: "モバイル API"
description: "クロスプラットフォームの application.Mobile マネージャー — iOS と Android で共通するネイティブモバイル機能のための、ビルドガード付きの単一エントリポイント"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="試験的機能"}
モバイル対応は試験段階であり、今後のリリースで変更される可能性があります。

@end

ネイティブモバイル機能は、次の 2 つの方法で公開されます。

- **プラットフォーム別マネージャー** — `//go:build ios` ファイル内の `application.IOS` と、`//go:build android` ファイル内の `application.Android`。プラットフォーム固有の処理には、これらを使用してください。プラットフォーム別 API の全体については、[iOS](/guides/mobile/ios/) および [Android](/guides/mobile/android/) のリファレンスを参照してください。
- **`application.Mobile`** — 両方のプラットフォームで同じように動作する機能のサブセットを網羅した、ビルドガード付きの単一マネージャー。どの環境でもコンパイルおよび実行できる単一のコードパスが必要な場合に使用してください。

## `application.Mobile`

`application.Mobile` は、iOS では `IOS`、Android では `Android`、デスクトップでは何もしないスタブに処理を振り分けます。ビルド制約がないため、通常のプラットフォーム非依存 Go コードから呼び出せます。独自の `//go:build` ファイルは不要です。

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()` は、アプリのプライベートファイルディレクトリへの絶対パスを返します。Android では `getFilesDir()`、iOS では Application Support ディレクトリであり、データベースやその他の永続ファイルの保存先として推奨されます。デスクトップでは空文字列を返します。デバイス上でもディレクトリを利用できない場合（iOS では作成できない場合）は空文字列を返すため、使用前に `""` かどうかを確認してください。

@note{type="note"}
デバイス外（デスクトップビルド）では、すべての `Mobile` メソッドは何もせず、すべてのクエリはゼロ値（文字列の場合は `""`）を返します。このため、クロスプラットフォームコードから `application.Mobile.*` を無条件に呼び出せます。デスクトップでも実際のパスが必要な場合は、プラットフォームで分岐し、`os.UserConfigDir()` などをフォールバックとして使用してください。

@end

## 機能

`Mobile` マネージャーは、iOS と Android でシグネチャが同一の機能を公開します。

| 機能 | API | 備考 |
| --- | --- | --- |
| 共有シート | `Mobile.Share(json)` | `{text, url}` |
| URL を外部で開く | `Mobile.OpenURL(url)` | システムブラウザー |
| 画面の点灯を維持 | `Mobile.SetKeepAwake(bool)` |  |
| トーチ／フラッシュライト | `Mobile.SetTorch(bool)` | → `common:torch` |
| セーフエリアのインセット | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| アプリ情報 | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| 画面の向きを固定 | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| ステータスバー | `Mobile.SetStatusBar(json)` | スタイル＋表示／非表示 |
| ストレージ情報 | `Mobile.StorageJSON()` | `{free,total}` バイト |
| ストレージパス | `Mobile.StoragePath()` | アプリ専用ファイルディレクトリ |
| 電源／バッテリー | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| ネットワーク状態 | `Mobile.NetworkJSON()` | `{connected,type}` |
| 生体認証 | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| セキュアストレージ | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| 位置情報 | `Mobile.GetLocation()` | 単発取得 → `common:location` |
| 触覚フィードバック | `Mobile.Haptic(type)` | 衝撃／通知／選択 |
| 加速度計 | `Mobile.SetMotion(bool)` | → `common:motion` |
| 近接センサー | `Mobile.SetProximity(bool)` | → `common:proximity` |
| テキスト読み上げ | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| キーボードのインセット | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| 画面キャプチャ | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| カメラ | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

非同期の結果は、プラットフォーム別のマネージャーとまったく同様に、`common:*` イベントとして届きます。ペイロードについては、[イベント](/guides/mobile/ios/#events)を参照してください。

## プラットフォーム固有のままの機能

iOS と Android で形式が異なる機能は、`Mobile`には<strong>含まれません</strong>。ビルドタグ付きファイルから `application.IOS` / `application.Android` を介して呼び出してください。

| 項目 | iOS | Android |
| --- | --- | --- |
| 明るさ（設定） | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| 明るさ／画面の向き（取得） | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| ローカル通知 | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| セキュアストア（書き込み） | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| バックグラウンド実行 | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
`MobileManager` インターフェースは、`application.Mobile` の基盤となる契約です。両方のプラットフォームマネージャーがこのインターフェースを満たす必要があるため、上記の各メソッドは iOS と Android で同一のシグネチャを維持することが保証されます。両者に差異が生じた場合、プラットフォームビルドはコンパイルに失敗します。

@end
