---
title: "Dockとタスクバー"
description: "macOSおよびWindowsでDockアイコンの表示／非表示を管理し、バッジを表示します"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## はじめに

Wailsは、デスクトップアプリケーション向けにクロスプラットフォームのDockサービスを提供します。このサービスでは、次の操作を実行できます。

- macOSのDockでアプリケーションアイコンを非表示または表示する
- アプリケーションタイルまたはDock／タスクバーアイコンにバッジを表示する（macOSおよびWindows）

## 基本的な使い方

### サービスの作成

まず、Dockサービスを初期化します。

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

### カスタムバッジオプションを指定したサービスの作成（Windowsのみ）

Windowsでは、さまざまなオプションを使用してバッジの外観をカスタマイズできます。

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

## Dockの操作

### Dockのアプリアイコンを非表示にする

macOSのDockからアプリアイコンを非表示にします。

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Dockのアプリアイコンを表示する

macOSのDockにアプリアイコンを表示します。

```go
// Show the app icon
dockService.ShowAppIcon()
```

## バッジの操作

### バッジの設定

アプリケーションタイル／Dockアイコンにバッジを設定します。

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### カスタムバッジの設定（Windowsのみ）

今回の呼び出しにのみ適用するオプションを指定して、バッジを設定します。

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

### バッジの削除

アプリケーションアイコンからバッジを削除します。

```go
dockService.RemoveBadge()
```

### 設定済みバッジの取得

```go
dockService.GetBadge()
```

## プラットフォームに関する考慮事項

@tabs
[macOS]
macOSでは、次のように動作します。

- Dockアイコンを<strong>非表示</strong>または<strong>表示</strong>にできます
- バッジはDockアイコン上に直接表示されます
- バッジオプションは<strong>カスタマイズできません</strong>（`NewWithOptions`/`SetCustomBadge`に渡したオプションはすべて無視されます）
- macOS標準のDockバッジスタイルが使用され、外観に合わせて自動的に調整されます
- ラベルのオーバーフローはシステムによって処理されます
- 空のラベルを指定すると、デフォルトの「●」バッジが表示されます

[Windows]
Windowsでは、次のように動作します。

- 現在、このサービスではタスクバーアイコンの非表示／表示はサポートされていません
- バッジはタスクバーにオーバーレイアイコンとして表示されます
- バッジにはテキスト値を使用できます
- `BadgeOptions`を使用してバッジの外観をカスタマイズできます
- バッジを表示するには、アプリケーションにウィンドウが必要です
- 複数文字のラベルには、小さいフォントサイズが自動的に使用されます
- ラベルのオーバーフローは処理されません
- カスタマイズオプション：
  - **TextColour**：テキストの色（デフォルト：白）
  - **BackgroundColour**：バッジの背景色（デフォルト：赤）
  - **FontName**：フォントファイル名（デフォルト：「segoeuib.ttf」）
  - **FontSize**：1文字の場合のフォントサイズ（デフォルト：18）
  - **SmallFontSize**：複数文字の場合のフォントサイズ（デフォルト：14）


[Linux]
Linuxでは、次のように動作します。

- Dockアイコンの表示／非表示およびバッジ機能は利用できません

@end

## ベストプラクティス

1. **Dockアイコンを非表示にする場合（macOS）：**
  - ユーザーが引き続きアプリにアクセスできるようにします（例：[システムトレイ](/features/menus/systray/)経由）
  - 代替UIに「終了」オプションを含めます
  - アプリはCommand+Tabのアプリ切り替え画面に表示されません
  - 開いているウィンドウは引き続き表示され、操作できます
  - すべてのウィンドウを閉じても、アプリが終了しない場合があります（macOSの動作によって異なります）
  - ユーザーは、Dockの右クリックメニューから終了する標準的な方法を利用できなくなります


2. **バッジは控えめに使用する：**
  - バッジを頻繁に更新しすぎると、ユーザーの注意をそらすおそれがあります
  - バッジは重要な通知にのみ使用します


3. **バッジのテキストを短くする：**
  - 数字のバッジが最も効果的です
  - macOSでは、テキストバッジを簡潔にします


4. **Windowsでバッジをカスタマイズする場合：**
  - テキストと背景色のコントラストを十分に確保します
  - 文字数が増えるとフォントサイズが小さくなるため、さまざまな長さのテキストでテストします
  - 確実に利用できるよう、一般的なシステムフォントを使用します


## APIリファレンス

### サービス管理

| メソッド | 説明 |
| --- | --- |
| `New()` | 新しい Dock サービスを作成します |
| `NewWithOptions(options BadgeOptions)` | カスタムバッジオプションを指定して新しい Dock サービスを作成します（Windows のみ。macOS および Linux ではオプションは無視されます） |

### Dock の操作

| メソッド | 説明 |
| --- | --- |
| `HideAppIcon()` | macOS の Dock からアプリアイコンを非表示にします（macOS のみ） |
| `ShowAppIcon()` | macOS の Dock にアプリアイコンを表示します（macOS のみ） |

### バッジの操作

| メソッド | 説明 |
| --- | --- |
| `SetBadge(label string) error` | 指定したラベルのバッジを設定します |
| `SetCustomBadge(label string, options BadgeOptions) error` | 指定したラベルとカスタムスタイルオプションを使用してバッジを設定します（Windows のみ） |
| `RemoveBadge() error` | アプリケーションアイコンからバッジを削除します |
| `GetBadge() *string` | 現在のバッジを取得します |

### 構造体と型

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
