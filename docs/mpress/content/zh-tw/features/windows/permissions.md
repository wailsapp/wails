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

### macOS（TCC）

macOS 透過其系統隱私權架構管理相機、麥克風、地理位置和通知的存取權。網頁內容首次請求某項功能時，作業系統會自動顯示提示，並在「系統設定」→「隱私權與安全性」中依應用程式記住使用者的選擇。

即使完全不設定`Permissions`，此機制也能正常運作。目前 macOS 會<strong>忽略此對應表</strong>；無論如何設定，所有請求都會交由 TCC 處理。實際的功能缺口是`PermissionDeny`在 macOS 上不起作用：若 TCC 已在系統層級授予某項功能的權限，你便無法阻止 WebView 使用該功能。

請確保你的`Info.plist`包含適當的用途說明鍵：

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## 常見模式

### 媒體擷取應用程式

在所有平台上允許使用相機和麥克風：

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

在<strong>Linux</strong>上，這會明確允許使用這兩種裝置；其他功能仍會遭到拒絕。 在<strong>Windows</strong>上，這會允許使用這兩種裝置；任何未列出的其他功能都會顯示原生提示。 在<strong>macOS</strong>上，這項設定不會生效；所有權限均由 TCC 處理。

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
| 麥克風 | ✅ | ✅ | 僅限 TCC |
| 攝影機 | ✅ | ✅ | 僅限 TCC |
| 地理位置 | ❌ 尚未支援 | ✅ | 僅限 TCC |
| 通知 | ❌ 尚未支援 | ✅ | 僅限 TCC |
| 讀取剪貼簿 | ❌ 尚未支援 | ✅ | 僅限 TCC |

## 疑難排解

**升級後，`getUserMedia`在 Linux 上仍然失敗**

請確認未明確設定`PermissionMicrophone: PermissionDeny`或`PermissionCamera: PermissionDeny`。預設值（未設定）允許在 Linux 上擷取媒體。

**Windows 正在提示我授予未設定的權限**

`Permissions`中只要有任何項目，Wails 就不再設定全面授予權限的`Allow`。未列出的功能會顯示 WebView2 的原生提示。請為應用程式使用的每項功能加入明確的`PermissionAllow`項目。

**macOS 權限無法運作**

`Permissions`對應表在 macOS 上不會生效。請確認`Info.plist`包含正確的用途說明鍵（`NSMicrophoneUsageDescription`、`NSCameraUsageDescription`等），而且使用者已在「系統設定」→「隱私權與安全性」中授予存取權。

**地理位置、通知和剪貼簿在 Linux 上不會生效**

目前 Linux 僅處理攝影機和麥克風。其他功能類型尚未實作支援；無論設定何種原則，這些功能仍會遭到拒絕。
