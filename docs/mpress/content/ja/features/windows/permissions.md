---
title: "権限"
description: "Web コンテンツからのカメラ、マイク、位置情報、その他の機能へのアクセス要求を制御します"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

`navigator.mediaDevices.getUserMedia()`、Geolocation API、または Notifications API を呼び出す Web コンテンツでは、ホストアプリケーションがそれらの要求を許可または拒否する必要があります。Wails は、`WebviewWindowOptions` にクロスプラットフォームの `Permissions` マップを公開しており、プラットフォーム固有のコードを使わずに宣言的に制御できます。

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

そのウィンドウの Web コンテンツからのカメラとマイクへのアクセス要求は、ブラウザーのプロンプトを表示せずに許可されます。

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

`Permission`（uint8）は、指定した種類に適用されるポリシーです。

| 定数 | 値 | 意味 |
| --- | --- | --- |
| `PermissionDefault` | 0 | プラットフォームのネイティブ処理を使用する（以下を参照） |
| `PermissionAllow` | 1 | 確認せずに許可する |
| `PermissionDeny` | 2 | 確認せずに拒否する |

`PermissionDefault` はゼロ値であるため、マップに設定されていないエントリはデフォルトとして動作します。

## プラットフォームごとの動作

基盤となる WebView のネイティブ動作が異なるため、各プラットフォームでは `PermissionDefault` の処理方法も異なります。

### Linux（WebKitGTK）

WebKitGTK には<strong>ネイティブの権限プロンプトがありません</strong>。ハンドラーが設定されていない場合、すべての要求が通知なしに拒否されます。このため、この機能が追加される前は `getUserMedia` が常に `NotAllowedError` を返していました。

現在、Wails は Linux 上で<strong>カメラとマイク</strong>へのアクセス要求を処理します。位置情報、通知、クリップボードの読み取りはまだ接続されていないため、設定したポリシーにかかわらず拒否されます。

| ポリシー | カメラ／マイク | 位置情報、通知、クリップボード |
| --- | --- | --- |
| `PermissionDefault` | **許可**（getUserMedia を復元） | 常に拒否 |
| `PermissionAllow` | 許可 | 常に拒否（未実装） |
| `PermissionDeny` | 拒否 | 常に拒否 |

### Windows（WebView2）

WebView2 には、ネイティブの権限プロンプトと、種類ごとの権限 API があります。5 種類の機能すべてが完全にサポートされています。

| ポリシー | 動作 |
| --- | --- |
| `PermissionDefault` | WebView2 が OS／ブラウザーのネイティブ権限プロンプトを表示する |
| `PermissionAllow` | 通知なしに許可 |
| `PermissionDeny` | 通知なしに拒否 |

<strong>重要：</strong>この機能が導入される前は、Wails は無条件に `SetGlobalPermission(Allow)` を呼び出し、すべての機能を通知なしに許可していました。現在は、`Permissions` にエントリが 1 つでも存在すると、その一括許可は設定され<strong>ません</strong>。設定されていない機能は、自動的に許可されるのではなく、WebView2 のネイティブプロンプトに委ねられます。

つまり、Windows で `Permissions` を少しでも構成すると、明示的に指定していない機能は通知なしに許可されず、プロンプトが表示されます。必要な機能を明示的に設定してください。

### macOS（TCC）

macOS は、システムのプライバシーフレームワークを通じて、カメラ、マイク、位置情報、通知へのアクセスを管理します。Web コンテンツが初めて機能へのアクセスを要求すると OS のプロンプトが自動的に表示され、ユーザーの選択は「システム設定」→「プライバシーとセキュリティ」にアプリ単位で記憶されます。

これは、`Permissions` を構成しなくても正しく動作します。現在、macOS では<strong>このマップは無視されます</strong>。設定内容にかかわらず、すべての要求が TCC を経由します。実際上の制約として、macOS では `PermissionDeny` が機能しません。TCC がシステムレベルですでに許可している機能を WebView が使用しないようにブロックすることはできません。

`Info.plist` に適切な使用目的の説明キーが含まれていることを確認してください：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## 一般的なパターン

### メディアキャプチャアプリ

すべてのプラットフォームでカメラとマイクを許可します：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Linux** では、これにより両方のデバイスが明示的に許可されます。その他の機能は引き続き拒否されます。 **Windows** では、これにより両方が許可されます。リストに含めていないその他の機能については、ネイティブプロンプトが表示されます。 **macOS** では効果がありません。すべて TCC によって処理されます。

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

### ウィンドウ単位のポリシー

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

ウィンドウ単位の `Windows.Permissions` フィールド（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）も引き続き機能し、クロスプラットフォームのマップが適用された後に、個々の機能をオーバーライドできます。クロスプラットフォームに相当するものがない WebView2 のアクセス許可の種類（`CoreWebView2PermissionKindOtherSensors` など）にアクセスする必要がある場合に使用してください。

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows での評価順序は次のとおりです。

1. クロスプラットフォームの `Permissions` マップ（`SetPermission` を介して種類ごとの状態を設定）
2. `Windows.Permissions` マップ（個々の種類をオーバーライド）
3. どちらのマップにも含まれない種類：ポリシーが設定されている場合は WebView2 のネイティブプロンプト、ポリシーが設定されていない場合は自動許可（従来の動作）

## プラットフォーム対応表

| 機能 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| マイク | ✅ | ✅ | TCC のみ |
| カメラ | ✅ | ✅ | TCC のみ |
| 位置情報 | ❌ 未対応 | ✅ | TCC のみ |
| 通知 | ❌ 未対応 | ✅ | TCC のみ |
| クリップボードの読み取り | ❌ 未対応 | ✅ | TCC のみ |

## トラブルシューティング

**`getUserMedia`がアップグレード後も Linux で失敗する**

`PermissionMicrophone: PermissionDeny` または `PermissionCamera: PermissionDeny` を明示的に設定していないことを確認してください。Linux では、デフォルト（未設定）でメディアキャプチャが許可されます。

**Windows で、設定していないアクセス許可を求めるプロンプトが表示される**

`Permissions` にエントリが1つでも存在すると、Wails は一括許可の `Allow` を設定しなくなります。リストに含めていない機能については、WebView2 のネイティブプロンプトが表示されます。アプリで使用するすべての機能について、明示的な `PermissionAllow` エントリを追加してください。

**macOS のアクセス許可が機能しない**

`Permissions` マップは macOS では効果がありません。`Info.plist` に正しい用途説明キー（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription` など）が含まれていること、およびユーザーが「システム設定」→「プライバシーとセキュリティ」でアクセスを許可していることを確認してください。

**Linux で位置情報／通知／クリップボードが機能しない**

現在、Linux で処理されるのはカメラとマイクのみです。その他の機能タイプへの対応はまだ実装されていません。設定したポリシーにかかわらず、これらは拒否されたままになります。
