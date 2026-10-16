---
title: "Dock 與工作列"
description: "管理 Dock 圖示的顯示狀態，並在 macOS 和 Windows 上顯示徽章"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## 簡介

Wails 為桌面應用程式提供跨平台的 Dock 服務。此服務可讓您：

- 在 macOS Dock 中隱藏及顯示應用程式圖示
- 在應用程式磚或 Dock／工作列圖示上顯示徽章（macOS 和 Windows）

## 基本用法

### 建立服務

首先，初始化 Dock 服務：

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### 使用自訂徽章選項建立服務（僅限 Windows）

在 Windows 上，您可以使用各種選項自訂徽章外觀：

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Dock 操作

### 隱藏 Dock 中的應用程式圖示

從 macOS Dock 隱藏應用程式圖示：

```go
// Hide the app icon
dockService.HideAppIcon()
```

### 顯示 Dock 中的應用程式圖示

在 macOS Dock 中顯示應用程式圖示：

```go
// Show the app icon
dockService.ShowAppIcon()
```

## 徽章操作

### 設定徽章

在應用程式磚／Dock 圖示上設定徽章：

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### 設定自訂徽章（僅限 Windows）

設定套用單次選項的徽章：

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### 移除徽章

從應用程式圖示移除徽章：

```go
dockService.RemoveBadge()
```

### 取得已設定的徽章

```go
dockService.GetBadge()
```

## 平台注意事項

@tabs
[macOS]
在 macOS 上：

- Dock 圖示可以<strong>隱藏</strong>及<strong>顯示</strong>
- 徽章會直接顯示在 Dock 圖示上
- 徽章選項<strong>無法自訂</strong>（傳遞給`NewWithOptions`/`SetCustomBadge`的任何選項都會被忽略）
- 使用標準的 macOS Dock 徽章樣式，並會自動配合外觀調整
- 標籤溢位由系統處理
- 提供空白標籤時，會顯示預設徽章「●」

[Windows]
在 Windows 上：

- 此服務目前不支援隱藏或顯示工作列圖示
- 徽章會以重疊圖示的形式顯示在工作列上
- 徽章支援文字值
- 可透過`BadgeOptions`自訂徽章外觀
- 應用程式必須有視窗，才能顯示徽章
- 多字元標籤會自動使用較小的字型大小
- 不會處理標籤溢位
- 自訂選項：
  - **TextColour**：文字色彩（預設：白色）
  - **BackgroundColour**：徽章背景色彩（預設：紅色）
  - **FontName**：字型檔案名稱（預設：「segoeuib.ttf」）
  - **FontSize**：單一字元的字型大小（預設：18）
  - **SmallFontSize**：多個字元的字型大小（預設：14）


[Linux]
在 Linux 上：

- 無法使用 Dock 圖示顯示狀態及徽章功能

@end

## 最佳實務

1. **隱藏 Dock 圖示時（macOS）：**
  - 確保使用者仍可存取您的應用程式（例如透過[系統匣](/features/menus/systray/)）
  - 在替代使用者介面中加入「結束」選項
  - 應用程式不會出現在 Command+Tab 切換器中
  - 已開啟的視窗仍然可見且可正常使用
  - 關閉所有視窗不一定會結束應用程式（macOS 的行為會有所不同）
  - 使用者將無法使用在 Dock 上按一下滑鼠右鍵的標準方式結束應用程式


2. **謹慎使用徽章：**
  - 過於頻繁地更新徽章可能會使使用者分心
  - 僅將徽章用於重要通知


3. **保持徽章文字簡短：**
  - 數字徽章最有效
  - 在 macOS 上，文字徽章應保持簡短


4. **自訂 Windows 徽章時：**
  - 確保文字與背景色彩之間有高對比度
  - 使用不同的文字長度進行測試，因為字型大小會隨文字變長而縮小
  - 使用常見的系統字型以確保可用


## API 參考

### 服務管理

| 方法 | 說明 |
| --- | --- |
| `New()` | 建立新的 Dock 服務 |
| `NewWithOptions(options BadgeOptions)` | 建立具有自訂徽章選項的新 Dock 服務（僅限 Windows；在 macOS 和 Linux 上會忽略這些選項） |

### Dock 操作

| 方法 | 說明 |
| --- | --- |
| `HideAppIcon()` | 從 macOS Dock 隱藏應用程式圖示（僅限 macOS） |
| `ShowAppIcon()` | 在 macOS Dock 顯示應用程式圖示（僅限 macOS） |

### 徽章操作

| 方法 | 說明 |
| --- | --- |
| `SetBadge(label string) error` | 使用指定的標籤設定徽章 |
| `SetCustomBadge(label string, options BadgeOptions) error` | 使用指定的標籤和自訂樣式選項設定徽章（僅限 Windows） |
| `RemoveBadge() error` | 從應用程式圖示移除徽章 |
| `GetBadge() *string` | 取得目前的徽章 |

### 結構與型別

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
