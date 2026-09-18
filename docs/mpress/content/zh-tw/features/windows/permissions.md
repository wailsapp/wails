---
title: "權限"
description: "控制網頁內容對相機、麥克風、地理位置及其他功能的請求"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

呼叫`navigator.mediaDevices.getUserMedia()`、Geolocation API 或 Notifications API 的網頁內容，需要由主機應用程式允許或拒絕這些請求。Wails 在`WebviewWindowOptions`上提供跨平台的`Permissions`對應表，讓你能以宣告方式進行控制，無須編寫平台特定程式碼。

## 快速開始

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

來自該視窗網頁內容的相機和麥克風請求會直接獲准，不會顯示瀏覽器提示。

## 權限類型

`PermissionType`（uint8）識別網頁內容可請求的功能。

| 常數 | 功能 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 權限值

`Permission`（uint8）是套用至指定類型的原則。

| 常數 | 值 | 含義 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 使用平台的原生處理方式（見下文） |
| `PermissionAllow` | 1 | 不提示，直接允許 |
| `PermissionDeny` | 2 | 不提示，直接拒絕 |

`PermissionDefault`是零值，因此對應表中未設定的項目會採用預設行為。

## 平台行為

由於各平台底層 WebView 的原生行為不同，因此每個平台處理`PermissionDefault`的方式也不同。

### Linux（WebKitGTK）

WebKitGTK<strong>沒有原生權限提示</strong>。若未附加處理常式，它會無聲地拒絕每個請求；這就是為什麼在加入此功能之前，`getUserMedia`一律傳回`NotAllowedError`。

Wails 目前可在 Linux 上處理<strong>相機和麥克風</strong>請求。地理位置、通知和剪貼簿讀取功能尚未接上，因此無論設定何種原則，都仍會遭到拒絕。

| 原則 | 相機／麥克風 | 地理位置、通知、剪貼簿 |
| --- | --- | --- |
| `PermissionDefault` | **允許**（恢復 getUserMedia） | 一律拒絕 |
| `PermissionAllow` | 允許 | 一律拒絕（尚未實作） |
| `PermissionDeny` | 拒絕 | 一律拒絕 |

### Windows（WebView2）

WebView2 具有原生權限提示，以及依權限種類區分的權限 API。五種功能類型全都受到完整支援。

| 原則 | 行為 |
| --- | --- |
| `PermissionDefault` | WebView2 顯示其原生作業系統／瀏覽器權限提示 |
| `PermissionAllow` | 無聲地允許 |
| `PermissionDeny` | 無聲地拒絕 |

<strong>重要：</strong>在此功能推出之前，Wails 會無條件呼叫`SetGlobalPermission(Allow)`，無聲地允許所有功能。現在，只要`Permissions`中有任何項目，就<strong>不會</strong>設定這項全面允許。未設定的功能不會自動獲准，而是交由 WebView2 顯示原生提示。

這表示在 Windows 上，只要你設定了任何`Permissions`項目，未明確列出的功能就會顯示提示，而非無聲地獲准。請明確設定所需的功能。

### macOS (WKWebView + TCC)

macOS 上需要兩層授權。WKWebView 在啟動擷取工作階段前詢問應用程式，由 `Permissions` 對應表回答該請求。在其下層，系統隱私權架構 TCC 控制裝置本身：應用程式首次實際存取相機或麥克風時會顯示系統提示，使用者的選擇依應用程式儲存在「系統設定 → 隱私權與安全性」中。

Wails 在 macOS 12 及更新版本上處理**相機和麥克風**請求。地理位置、通知和剪貼簿讀取沒有對應的 `WKUIDelegate` 介面，因此尚未串接，仍交由 TCC 處理。與 Linux 一樣，為這些功能設定的原則沒有作用。

| 原則 | 相機 / 麥克風 | 地理位置、通知、剪貼簿 |
| --- | --- | --- |
| `PermissionDefault` | WebKit 顯示自己的權限提示 | 僅限 TCC |
| `PermissionAllow` | 略過 WebKit 提示 — **TCC 仍然適用** | 僅限 TCC |
| `PermissionDeny` | 在存取裝置之前拒絕 | 僅限 TCC |

`PermissionAllow` 允許的是 WebView 請求，而不是裝置本身。首次擷取仍會觸發 TCC 提示；使用者在系統設定中拒絕的應用程式仍會被拒絕，應用程式無法自行授予裝置存取權。`PermissionAllow` 省略的只有前置的 WebKit 提示。

macOS 12 之前不存在此委派方法，因此會忽略對應表，所有請求都回到 WebKit 提示。

請確保你的`Info.plist`包含適當的用途說明鍵：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="拒絕 Info.plist 未宣告用途的功能"}

如果請求在缺少對應用途說明鍵的情況下到達 AVFoundation，結果不只是請求失敗：macOS 會終止應用程式。

`PermissionDefault` 是零值，因此只設定 `{PermissionMicrophone: PermissionAllow}` 會讓相機繼續使用 WebKit 提示。如果使用者接受提示，而應用程式只宣告了 `NSMicrophoneUsageDescription`，應用程式就會終止。請為每個沒有用途說明的功能明確設定 `PermissionDeny`：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionDeny,
},
```

@end

## 常見模式

### 媒體擷取應用程式

在所有平台上允許使用相機和麥克風：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在<strong>Linux</strong>上，這會明確允許使用這兩種裝置；其他功能仍會遭到拒絕。 在<strong>Windows</strong>上，這會允許使用這兩種裝置；任何未列出的其他功能都會顯示原生提示。 在 **macOS** 上，這會在 WebKit 層允許兩者，不顯示瀏覽器提示。首次使用時，TCC 仍會詢問裝置本身的存取權，且 `Info.plist` 必須包含兩個用途說明鍵。

### 在 Linux 上拒絕媒體擷取

Linux 預設允許使用攝影機和麥克風。若要停用：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### 在 Windows 上允許所有功能

若要直接授予所有功能而不顯示提示：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### 各視窗的原則

不同視窗可以使用不同的原則：

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

## Windows 專用覆寫

各視窗的`Windows.Permissions`欄位（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）仍然有效，並可在套用跨平台對應表後覆寫個別功能。需要存取沒有跨平台對應項目的 WebView2 權限類型（例如`CoreWebView2PermissionKindOtherSensors`）時，請使用此欄位。

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows 上的評估順序如下：

1. 跨平台`Permissions`對應表（透過`SetPermission`設定各類型的狀態）
2. `Windows.Permissions`對應表（覆寫個別類型）
3. 對於兩個對應表均未涵蓋的任何類型：若已設定原則，則顯示 WebView2 的原生提示；若未設定原則，則自動允許（舊版行為）

## 平台支援矩陣

| 功能 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 麥克風 | ✅ | ✅ | ✅ (macOS 12+) |
| 攝影機 | ✅ | ✅ | ✅ (macOS 12+) |
| 地理位置 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |
| 通知 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |
| 讀取剪貼簿 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |

macOS 欄為 ✅ 時，原則回答 WebKit 請求，TCC 仍額外控制裝置。為 ❌ 時，該功能僅交由 TCC 處理。

## 疑難排解

**升級後，`getUserMedia`在 Linux 上仍然失敗**

請確認未明確設定`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。預設值（未設定）允許在 Linux 上擷取媒體。

**Windows 正在提示我授予未設定的權限**

`Permissions`中只要有任何項目，Wails 就不再設定全面授予權限的`Allow`。未列出的功能會顯示 WebView2 的原生提示。請為應用程式使用的每項功能加入明確的`PermissionAllow`項目。

**macOS 權限無法運作**

`Permissions` 支援 macOS 12 及更新版本的相機和麥克風；地理位置、通知和剪貼簿讀取尚未串接，會忽略原則。TCC 仍控制裝置：`PermissionAllow` 只移除 WebKit 提示，不移除系統提示。請檢查 `Info.plist` 中的 `NSMicrophoneUsageDescription`、`NSCameraUsageDescription`，以及「系統設定 → 隱私權與安全性」中的授權。

**網頁內容請求相機或麥克風時，macOS 應用程式結束**

如果請求在缺少對應用途說明鍵的情況下到達 AVFoundation，結果不只是請求失敗：macOS 會終止應用程式。 請加入對應的鍵，或為該功能設定 `PermissionDeny`，使請求無法到達 AVFoundation。`PermissionDefault` 會保留使用者可以接受的 WebKit 提示。

**地理位置、通知和剪貼簿在 Linux 上不會生效**

目前 Linux 僅處理攝影機和麥克風。其他功能類型尚未實作支援；無論設定何種原則，這些功能仍會遭到拒絕。
