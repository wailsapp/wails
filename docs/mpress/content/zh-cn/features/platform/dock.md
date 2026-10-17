---
title: "程序坞与任务栏"
description: "管理 macOS 和 Windows 上程序坞图标的可见性并显示徽标"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## 简介

Wails 为桌面应用提供跨平台的程序坞服务。通过此服务，你可以：

- 在 macOS 程序坞中隐藏和显示应用图标
- 在应用磁贴或程序坞/任务栏图标上显示徽标（macOS 和 Windows）

## 基本用法

### 创建服务

首先，初始化程序坞服务：

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

### 使用自定义徽标选项创建服务（仅限 Windows）

在 Windows 上，可以使用各种选项自定义徽标外观：

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

## 程序坞操作

### 隐藏程序坞中的应用图标

从 macOS 程序坞中隐藏应用图标：

```go
// Hide the app icon
dockService.HideAppIcon()
```

### 显示程序坞中的应用图标

在 macOS 程序坞中显示应用图标：

```go
// Show the app icon
dockService.ShowAppIcon()
```

## 徽标操作

### 设置徽标

在应用磁贴/程序坞图标上设置徽标：

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### 设置自定义徽标（仅限 Windows）

设置徽标并应用仅用于本次设置的选项：

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

### 移除徽标

从应用图标中移除徽标：

```go
dockService.RemoveBadge()
```

### 获取已设置的徽标

```go
dockService.GetBadge()
```

## 平台注意事项

@tabs
[macOS]
在 macOS 上：

- 可以<strong>隐藏</strong>和<strong>显示</strong>程序坞图标
- 徽标直接显示在程序坞图标上
- 徽标选项<strong>不可自定义</strong>（传递给`NewWithOptions`/`SetCustomBadge`的所有选项都会被忽略）
- 使用标准的 macOS 程序坞徽标样式，并自动适应外观设置
- 标签溢出由系统处理
- 提供空标签时，将显示默认徽标“●”

[Windows]
在 Windows 上：

- 此服务目前不支持隐藏或显示任务栏图标
- 徽标以叠加图标的形式显示在任务栏中
- 徽标支持文本值
- 可以通过`BadgeOptions`自定义徽标外观
- 应用必须拥有窗口才能显示徽标
- 对于多字符标签，会自动使用较小的字号
- 不处理标签溢出
- 自定义选项：
  - **TextColour**：文本颜色（默认：白色）
  - **BackgroundColour**：徽标背景颜色（默认：红色）
  - **FontName**：字体文件名（默认：“segoeuib.ttf”）
  - **FontSize**：单字符字号（默认：18）
  - **SmallFontSize**：多字符字号（默认：14）


[Linux]
在 Linux 上：

- 程序坞图标可见性控制和徽标功能不可用

@end

## 最佳实践

1. **隐藏程序坞图标时（macOS）：**
  - 确保用户仍可访问你的应用（例如通过[系统托盘](/features/menus/systray/)）
  - 在替代用户界面中提供“退出”选项
  - 应用不会出现在 Command+Tab 切换器中
  - 已打开的窗口仍然可见且可正常使用
  - 关闭所有窗口不一定会退出应用（macOS 的行为有所不同）
  - 用户将无法通过右键单击程序坞图标这一标准方式退出应用


2. **谨慎使用徽标：**
  - 过于频繁地更新徽标可能会分散用户的注意力
  - 仅将徽标用于重要通知


3. **保持徽标文本简短：**
  - 数字徽标最有效
  - 在 macOS 上，文本徽标应保持简短


4. **自定义 Windows 徽标时：**
  - 确保文本与背景颜色之间具有较高的对比度
  - 使用不同长度的文本进行测试，因为文本越长，字号越小
  - 使用常见的系统字体以确保字体可用


## API 参考

### 服务管理

| 方法 | 说明 |
| --- | --- |
| `New()` | 创建新的程序坞服务 |
| `NewWithOptions(options BadgeOptions)` | 创建带有自定义徽章选项的新程序坞服务（仅限 Windows；在 macOS 和 Linux 上会忽略这些选项） |

### 程序坞操作

| 方法 | 说明 |
| --- | --- |
| `HideAppIcon()` | 在 macOS 程序坞中隐藏应用图标（仅限 macOS） |
| `ShowAppIcon()` | 在 macOS 程序坞中显示应用图标（仅限 macOS） |

### 徽章操作

| 方法 | 说明 |
| --- | --- |
| `SetBadge(label string) error` | 使用指定标签设置徽章 |
| `SetCustomBadge(label string, options BadgeOptions) error` | 使用指定标签和自定义样式选项设置徽章（仅限 Windows） |
| `RemoveBadge() error` | 从应用图标中移除徽章 |
| `GetBadge() *string` | 获取当前徽章 |

### 结构体和类型

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
