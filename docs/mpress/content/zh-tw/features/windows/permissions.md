---
title: "權限"
description: "控制網頁內容對相機、麥克風、地理位置及其他功能的請求"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

呼叫`navigator.mediaDevices.getUserMedia()`、Geolocation API 或 Notifications API 的網頁內容，需要由主機應用程式允許或拒絕這些請求。Wails 在`WebviewWindowOptions`上提供跨平台的`Permissions`對應表，讓你以宣告方式控制權限，無須撰寫平台專用程式碼。

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

來自該視窗網頁內容的相機與麥克風請求會直接獲准，不會顯示瀏覽器提示。

## 權限類型

`PermissionType` (uint8) 用於識別網頁內容可請求的功能。

| 常數 | 功能 |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## 權限值

`Permission` (uint8) 是套用至指定類型的原則。

| 常數 | 值 | 含義 |
| --- | --- | --- |
| `PermissionDefault` | 0 | 使用平台的原生處理方式（見下文） |
| `PermissionAllow` | 1 | 不顯示提示並直接允許 |
| `PermissionDeny` | 2 | 不顯示提示並直接拒絕 |

`PermissionDefault`是零值，因此對應表中未設定的項目會採用預設行為。

## 平台行為

由於底層 WebView 的原生行為不同，各平台處理`PermissionDefault`的方式也不同。

### Linux（WebKitGTK）

WebKitGTK**沒有原生權限提示**。如果未附加處理常式，它會在不顯示提示的情況下拒絕所有請求；這就是加入此功能之前，`getUserMedia`一律傳回`NotAllowedError`的原因。

Wails 目前會在 Linux 上處理**相機與麥克風**請求。地理位置、通知及剪貼簿讀取功能尚未接通，因此無論設定何種原則，請求仍會遭到拒絕。

| 原則 | 相機／麥克風 | 地理位置、通知、剪貼簿 |
| --- | --- | --- |
| `PermissionDefault` | **允許**（恢復 getUserMedia） | 一律拒絕 |
| `PermissionAllow` | 允許 | 一律拒絕（尚未實作） |
| `PermissionDeny` | 拒絕 | 一律拒絕 |

### Windows（WebView2）

WebView2 具有原生權限提示及依類型區分的權限 API，完整支援全部五種功能類型。

| 原則 | 行為 |
| --- | --- |
| `PermissionDefault` | WebView2 顯示其原生作業系統／瀏覽器權限提示 |
| `PermissionAllow` | 不顯示提示並直接允許 |
| `PermissionDeny` | 不顯示提示並直接拒絕 |

**重要：**在此功能推出前，Wails 會無條件呼叫`SetGlobalPermission(Allow)`，在不顯示提示的情況下允許所有功能。現在，只要`Permissions`中有任何項目，就**不會**設定這項全面允許。未設定的功能不再自動獲准，而是交由 WebView2 顯示原生提示。

這表示只要在 Windows 上設定了任何`Permissions`，未明確列出的功能就會顯示提示，而不是在不顯示提示的情況下獲准。請明確設定所需功能。

### macOS（WKWebView + TCC）

在 macOS 上，必須取得兩個層級的同意。WKWebView 會在開始擷取工作階段前詢問應用程式，而`Permissions`對應表會回應這項請求。在其下層，系統隱私權架構（TCC）會管控裝置本身：應用程式第一次實際存取相機或麥克風時會顯示作業系統提示，並且系統會在「系統設定」→「隱私權與安全性」中針對每個應用程式記住使用者的選擇。

Wails 會在 macOS 12 及更新版本上處理**相機和麥克風**請求。地理位置、通知和剪貼簿讀取沒有 `WKUIDelegate` 對應機制，因此未接線處理，僅由 TCC 管控——與 Linux 相同，您為這些功能設定的原則不會產生任何作用。

| 原則 | 相機／麥克風 | 地理位置、通知、剪貼簿 |
| --- | --- | --- |
| `PermissionDefault` | WebKit 顯示自己的權限提示 | 僅由 TCC 處理 |
| `PermissionAllow` | 略過 WebKit 的提示，但**TCC 仍然適用** | 僅由 TCC 處理 |
| `PermissionDeny` | 由 WebKit 拒絕；系統授權另行處理 | 僅由 TCC 處理 |

`PermissionAllow`允許的是 WebView 的請求，而不是裝置存取權。系統授權仍會獨立套用，在「系統設定」中遭拒絕的應用程式仍會遭到拒絕。只有 WebKit 自己的提示會被略過。

在 macOS 12 以下版本中，該委派方法不存在，因此對應表會被忽略，而且每個請求都會改由 WebKit 顯示提示。

請確認`Info.plist`包含適當的用途說明鍵：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="必須提供用途說明"}

WebKit 可能會在呼叫 Wails 權限委派之前請求系統授權。`PermissionDeny`會封鎖 WebKit 請求，但無法保證不會出現 TCC 提示或 AVFoundation 授權請求。

請針對應用程式所請求的裝置加入適當的用途說明；若應用程式在沙箱中執行，也要加入相符的裝置權利。請勿以`PermissionDeny`取代這些宣告。

@end

## 常見模式

### 媒體擷取應用程式

在所有平台上允許相機與麥克風：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在**Linux**上，這會明確允許使用這兩種裝置；其他功能仍會遭到拒絕。在**Windows**上，這會允許使用這兩種裝置；任何未列出的其他功能都會顯示原生提示。在**macOS**上，這會在 WebKit 層允許使用這兩種裝置，因此不會顯示瀏覽器提示；首次使用時，TCC 仍會針對裝置本身要求授權，而且`Info.plist`中必須包含這兩個用途說明鍵。

### 在 Linux 上拒絕媒體擷取

Linux 預設允許使用攝影機和麥克風。若要停用：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### 在 Windows 上允許所有功能

若要允許所有功能而不顯示提示：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### 個別視窗原則

不同視窗可以採用不同的原則：

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

## Windows 專用覆寫設定

套用跨平台對應表後，個別視窗的`Windows.Permissions`欄位（`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`）仍然有效，且可覆寫個別功能。需要存取沒有跨平台對應項目的 WebView2 權限種類（例如`CoreWebView2PermissionKindOtherSensors`）時，請使用此欄位。

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Windows 上的評估順序如下：

1. 跨平台`Permissions`對應表（透過`SetPermission`設定各種類的狀態）
2. `Windows.Permissions`對應表（覆寫個別種類）
3. 對於兩個對應表均未涵蓋的任何種類：已設定原則時顯示 WebView2 的原生提示；未設定原則時自動允許（舊有行為）

## 平台支援矩陣

| 功能 | Linux | Windows | macOS |
| --- | --- | --- | --- |
| 麥克風 | ✅ | ✅ | ✅（macOS 12+） |
| 攝影機 | ✅ | ✅ | ✅（macOS 12+） |
| 地理位置 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |
| 通知 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |
| 讀取剪貼簿 | ❌ 尚未支援 | ✅ | ❌ 尚未支援 |

macOS 欄標示為 ✅ 時，原則會回應 WebKit 的要求，但 TCC 仍會在此之外限制裝置存取。標示為 ❌ 時，該功能完全交由 TCC 處理。

## 疑難排解

**`getUserMedia`升級後在 Linux 上仍然失敗**

請確認您未明確設定`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。Linux 上的預設值（未設定）允許媒體擷取。

**Windows 針對我未設定的權限顯示提示**

只要`Permissions`中存在任何項目，Wails 就不再設定一律允許的`Allow`授權。未列出的功能會顯示 WebView2 的原生提示。請為應用程式使用的每項功能加入明確的`PermissionAllow`項目。

**macOS 權限無法運作**

`Permissions`涵蓋 macOS 12及更新版本上的攝影機和麥克風；地理位置、通知和讀取剪貼簿尚未接通，因此會忽略此原則。在原則生效的情況下，TCC 仍會在此之外限制裝置存取：`PermissionAllow`會移除 WebKit 的提示，但不會移除系統提示。請確認`Info.plist`包含正確的用途說明鍵（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription`），而且使用者已在「系統設定」→「隱私權與安全性」中授予存取權。

**網頁內容要求使用攝影機或麥克風時，我的 macOS 應用程式會結束**

請確認應用程式要求使用的每個裝置，在`Info.plist`中都有對應的`NSCameraUsageDescription`或`NSMicrophoneUsageDescription`。缺少說明可能導致要求遭到拒絕；在缺少必要說明的情況下存取 AVFoundation，可能會終止應用程式。`PermissionDeny`控制 WebKit 的決策，不能取代這些鍵。

**地理位置／通知／剪貼簿在 Linux 上沒有作用**

Linux 目前只處理攝影機和麥克風。其他功能類型尚未實作支援，因此無論您設定何種原則，它們仍會遭到拒絕。
