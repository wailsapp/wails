---
title: "権限"
description: "Web コンテンツからのカメラ、マイク、位置情報などの機能へのアクセス要求を制御します"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

`navigator.mediaDevices.getUserMedia()`、Geolocation API、または Notifications API を呼び出す Web コンテンツでは、ホストアプリケーションがそれらの要求を許可または拒否する必要があります。Wails は、`WebviewWindowOptions` にクロスプラットフォームの `Permissions` マップを提供しており、プラットフォーム固有のコードを使わずに宣言的に制御できます。

## クイックスタート

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

そのウィンドウの Web コンテンツからのカメラおよびマイクへのアクセス要求は、ブラウザーのプロンプトを表示せずに許可されます。

## 権限の種類

`PermissionType`（uint8）は、Web コンテンツが要求できる機能を識別します。

| 定数 | 機能 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 権限の値

`Permission`（uint8）は、指定した種類に適用するポリシーです。

| 定数 | 値 | 意味 |
| --- | --- | --- |
| `PermissionDefault` | 0 | プラットフォームのネイティブ処理を使用（下記参照） |
| `PermissionAllow` | 1 | プロンプトを表示せずに許可 |
| `PermissionDeny` | 2 | プロンプトを表示せずに拒否 |

`PermissionDefault` はゼロ値であるため、マップに設定されていないエントリーはデフォルトとして動作します。

## プラットフォームごとの動作

基盤となる WebView のネイティブ動作が異なるため、`PermissionDefault` の処理方法はプラットフォームごとに異なります。

### Linux（WebKitGTK）

WebKitGTK には**ネイティブの権限プロンプトがありません**。ハンドラーが設定されていない場合、すべての要求を通知せずに拒否します。このため、この機能が追加される前は `getUserMedia` が常に `NotAllowedError` を返していました。

現在、Wails は Linux で**カメラおよびマイク**へのアクセス要求を処理します。位置情報、通知、クリップボードの読み取りはまだ接続されていないため、設定したポリシーにかかわらず拒否されます。

| ポリシー | カメラ／マイク | 位置情報、通知、クリップボード |
| --- | --- | --- |
| `PermissionDefault` | **許可**（getUserMedia を再び利用可能にする） | 常に拒否 |
| `PermissionAllow` | 許可 | 常に拒否（未実装） |
| `PermissionDeny` | 拒否 | 常に拒否 |

### Windows（WebView2）

WebView2 にはネイティブの権限プロンプトと、種類ごとの権限 API があります。5 種類の機能すべてが完全にサポートされています。

| ポリシー | 動作 |
| --- | --- |
| `PermissionDefault` | WebView2 がネイティブの OS／ブラウザー権限プロンプトを表示 |
| `PermissionAllow` | 通知せずに許可 |
| `PermissionDeny` | 通知せずに拒否 |

**重要：**この機能が導入される前、Wails は無条件に `SetGlobalPermission(Allow)` を呼び出し、すべての機能を通知せずに許可していました。現在は、`Permissions` にエントリーが 1 つでも存在する場合、この一括許可は設定**されません**。未設定の機能は自動的に許可されず、WebView2 のネイティブプロンプトにフォールスルーします。

つまり、Windows で `Permissions` を少しでも設定すると、明示的に指定していない機能は通知せずに許可されるのではなく、プロンプトが表示されます。必要な機能を明示的に設定してください。

### macOS（WKWebView + TCC）

macOS では、2 つのレイヤーの両方で許可される必要があります。WKWebView はキャプチャセッションを開始する前にアプリケーションへ問い合わせ、この要求に `Permissions` マップが応答します。その下では、システムのプライバシーフレームワーク（TCC）がデバイス自体へのアクセスを制御します。アプリが実際にカメラまたはマイクへ初めてアクセスしたときに OS のプロンプトが表示され、ユーザーの選択はアプリごとに「システム設定」→「プライバシーとセキュリティ」に記憶されます。

Wails は macOS 12 以降で、**カメラおよびマイク**へのアクセス要求を処理します。位置情報、通知、クリップボードの読み取りには `WKUIDelegate` に相当するものがないため、接続されておらず、TCC のみに委ねられます。Linux と同様に、これらに設定したポリシーは効果がありません。

| ポリシー | カメラ／マイク | 位置情報、通知、クリップボード |
| --- | --- | --- |
| `PermissionDefault` | WebKit が独自の権限プロンプトを表示 | TCC のみ |
| `PermissionAllow` | WebKit のプロンプトを省略 — **TCC は引き続き適用** | TCC のみ |
| `PermissionDeny` | WebKit が拒否。システムの承認は別途処理 | TCC のみ |

`PermissionAllow` が許可するのは WebView の要求であり、デバイスへのアクセスではありません。システムの承認は引き続き独立して適用され、「システム設定」で拒否されたアプリは拒否されたままです。省略されるのは WebKit 独自のプロンプトだけです。

macOS 12 より前のバージョンにはデリゲートメソッドが存在しないため、マップは無視され、すべての要求で WebKit のプロンプトが使用されます。

`Info.plist` に適切な使用目的の説明キーが含まれていることを確認してください。

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="使用目的の説明が必要です"}

WebKit は、Wails の権限デリゲートを呼び出す前にシステムの承認を要求する場合があります。`PermissionDeny` は WebKit の要求をブロックしますが、TCC のプロンプトや AVFoundation の承認要求が発生しないことを保証するものではありません。

アプリが要求するデバイスに適した使用目的の説明を含め、サンドボックス化する場合は対応するデバイスエンタイトルメントも含めてください。`PermissionDeny` をこれらの宣言の代わりに使用しないでください。

@end

## 一般的なパターン

### メディアキャプチャアプリ

すべてのプラットフォームでカメラとマイクを許可します。

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Linux** では、これにより両方のデバイスが明示的に許可され、その他の機能は引き続き拒否されます。**Windows** では両方が許可され、一覧に含めていないその他の機能についてはネイティブプロンプトが表示されます。**macOS** では WebKit レイヤーで両方が許可されるため、ブラウザのプロンプトは表示されません。ただし、初回使用時には TCC がデバイス自体へのアクセス許可を引き続き要求し、`Info.plist`には両方の使用目的説明キーが必要です。

### Linux でメディアキャプチャを拒否する

Linux では、デフォルトでカメラとマイクが許可されます。無効にするには、次のようにします。

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Windows ですべての機能を許可する

プロンプトを表示せずにすべての機能を許可するには、次のようにします。

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### ウィンドウごとのポリシー

ウィンドウごとに異なるポリシーを設定できます。

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Windows 固有のオーバーライド

ウィンドウごとの `Windows.Permissions` フィールド（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）も引き続き機能し、クロスプラットフォームのマップが適用された後に個別の機能をオーバーライドできます。クロスプラットフォームに相当するものがない WebView2 の権限種別（例: `CoreWebView2PermissionKindOtherSensors`）にアクセスする必要がある場合に使用してください。

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows での評価順序は次のとおりです。

1. クロスプラットフォームの `Permissions` マップ（`SetPermission` を介して種別ごとの状態を設定）
2. `Windows.Permissions` マップ（個別の種別をオーバーライド）
3. どちらのマップにも含まれない種別: ポリシーが設定されている場合は WebView2 のネイティブプロンプト、ポリシーが設定されていない場合は自動許可（従来の動作）

## プラットフォーム対応表

| 機能 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| マイク | ✅ | ✅ | ✅（macOS 12 以降） |
| カメラ | ✅ | ✅ | ✅（macOS 12 以降） |
| 位置情報 | ❌ 未対応 | ✅ | ❌ 未対応 |
| 通知 | ❌ 未対応 | ✅ | ❌ 未対応 |
| クリップボードの読み取り | ❌ 未対応 | ✅ | ❌ 未対応 |

macOS の欄が ✅ の場合、ポリシーが WebKit の要求に応答しますが、それに加えて TCC が引き続きデバイスへのアクセスを制御します。❌ の場合、その機能は TCC のみに委ねられます。

## トラブルシューティング

**`getUserMedia`アップグレード後も Linux で失敗する**

`PermissionMicrophone: PermissionDeny` または `PermissionCamera: PermissionDeny` を明示的に設定していないことを確認してください。デフォルト（未設定）では、Linux でメディアキャプチャが許可されます。

**Windows で、設定していない権限のプロンプトが表示される**

`Permissions` にエントリが1つでも存在すると、Wails はすべてを対象とする `Allow` 許可を設定しなくなります。一覧に含めていない機能については、WebView2 のネイティブプロンプトが表示されます。アプリが使用するすべての機能について、明示的な `PermissionAllow` エントリを追加してください。

**macOS の権限が機能しない**

`Permissions` は macOS 12 以降のカメラとマイクに対応しています。位置情報、通知、クリップボードの読み取りはまだ接続されておらず、ポリシーを無視します。ポリシーが適用される場合でも、それに加えて TCC が引き続きデバイスへのアクセスを制御します。`PermissionAllow` によって WebKit のプロンプトは表示されなくなりますが、システムのプロンプトは表示されます。`Info.plist` に正しい使用目的説明キー（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription`）が含まれていること、およびユーザーが「システム設定」→「プライバシーとセキュリティ」でアクセスを許可していることを確認してください。

**Web コンテンツがカメラまたはマイクへのアクセスを要求すると macOS アプリが終了する**

アプリが要求する各デバイスについて、`Info.plist` に `NSCameraUsageDescription` または `NSMicrophoneUsageDescription` が含まれていることを確認してください。説明がない場合、要求が拒否されることがあります。必要な説明なしで AVFoundation にアクセスすると、アプリが終了する可能性があります。`PermissionDeny` は WebKit の判断を制御するものであり、これらのキーの代わりにはなりません。

**Linux で位置情報、通知、クリップボードの設定が機能しない**

現在 Linux で処理されるのは、カメラとマイクのみです。その他の機能種別はまだ実装されていないため、設定したポリシーにかかわらず拒否されたままになります。
