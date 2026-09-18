---
title: "選單"
description: "在 Wails v3 中建立及自訂選單的指南"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 提供功能強大的選單系統，可用來建立應用程式選單和快顯選單。本指南將逐一介紹選單系統的各項功能。

## 建立選單

若要建立新選單，請使用 Menus 管理器的 `New()` 方法：

```go
menu := app.Menu.New()
```

### 新增選單項目

Wails 支援數種選單項目，各有特定用途：

#### 一般選單項目

一般選單項目是選單的基本組成元素。它們會顯示文字，並可在點選時觸發動作：

```go
menuItem := menu.Add("Click Me")
```

#### 核取方塊

核取方塊選單項目提供可切換的狀態，適合用來啟用或停用功能或設定：

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### 選項按鈕群組

選項按鈕群組可讓使用者從一組互斥選項中選擇一項。相鄰放置選項按鈕項目時，系統會自動建立群組：

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### 分隔線

分隔線是水平線，可協助將選單項目整理成邏輯群組：

```go
menu.AddSeparator()
```

#### 子選單

子選單是巢狀選單，將滑鼠游標停留在選單項目上或點選該項目時便會顯示。它們很適合用來整理複雜的選單結構：

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### 合併選單

您可以將一個選單附加或前置到另一個選單中。

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
依預設，`prepend` 和 `append` 會與原始選單共用狀態。如果要建立具有獨立狀態的新選單， 可以對該選單呼叫 `.Clone()`。

例如：`menu.Append(secondaryMenu.Clone())`

@end

#### 清除選單

如果選單項目的數量不固定，在某些情況下，重新建構整個選單會更合適。

這會清除現有選單中的所有項目，之後您可以重新新增項目。

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
清除選單只會清除最上層的選單項目。子選單雖然不再顯示，但仍會占用記憶體， 因此請務必妥善管理選單。

@end

#### 銷毀選單

若要清除並釋放選單，請使用 `Destroy()` 方法：

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### 選單項目屬性

選單項目有數個可設定的屬性：

| 屬性 | 方法 | 說明 |
| --- | --- | --- |
| 標籤 | `SetLabel(string)` | 設定顯示文字 |
| 已啟用 | `SetEnabled(bool)` | 啟用或停用項目 |
| 已勾選 | `SetChecked(bool)` | 設定勾選狀態（適用於核取方塊或選項按鈕項目） |
| 工具提示 | `SetTooltip(string)` | 設定工具提示文字 |
| 已隱藏 | `SetHidden(bool)` | 顯示或隱藏項目 |
| 快速鍵 | `SetAccelerator(string)` | 設定鍵盤快速鍵 |

### 選單項目狀態

選單項目可處於不同狀態，以控制其可見性及互動能力：

#### 可見性

您可以使用 `SetHidden()` 方法，動態顯示或隱藏選單項目：

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

隱藏的選單項目會從選單中完全移除，直到再次顯示為止。這適合用於只應在特定應用程式狀態下出現的情境式選單項目。

#### 啟用狀態

您可以使用 `SetEnabled()` 方法啟用或停用選單項目：

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

停用的選單項目仍然可見，但會呈現灰色且無法點選。這通常用來表示某項動作目前無法使用，例如：

- 沒有任何變更可儲存時，停用「儲存」
- 未選取任何內容時，停用「複製」
- 沒有可復原的動作時，停用「復原」

#### 動態狀態管理

您可以將這些狀態與事件處理常式結合，以建立動態選單：

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### 事件處理

選單項目可使用`OnClick`方法回應點擊事件：

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

事件內容提供所點擊選單項目的相關資訊：

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### 以角色為基礎的選單項目

Wails 提供一組預先定義的選單角色，可自動建立具備標準功能的選單項目。以下是支援的選單角色：

#### 完整選單結構

這些角色會建立具備常用功能的完整選單結構：

| 角色 | 說明 | 平台備註 |
| --- | --- | --- |
| `AppMenu` | 包含「關於」、「服務」、「隱藏／顯示」及「結束」的應用程式選單 | 僅限 macOS |
| `EditMenu` | 包含「復原」、「重做」、「剪下」、「複製」、「貼上」等項目的標準「編輯」選單 | 所有平台 |
| `ViewMenu` | 包含「重新載入」、「縮放」及「全螢幕」控制項的「檢視」選單 | 所有平台 |
| `WindowMenu` | 視窗控制項（最小化、縮放等） | 所有平台 |
| `HelpMenu` | 包含連至 Wails 網站之「深入瞭解」連結的「說明」選單 | 所有平台 |

#### 個別選單項目

這些角色可用來新增個別選單項目：

| 角色 | 說明 | 平台備註 |
| --- | --- | --- |
| `About` | 顯示應用程式的「關於」對話方塊 | 所有平台 |
| `Hide` | 隱藏應用程式 | 僅限 macOS |
| `HideOthers` | 隱藏其他應用程式 | 僅限 macOS |
| `UnHide` | 顯示已隱藏的應用程式 | 僅限 macOS |
| `CloseWindow` | 關閉目前視窗 | 所有平台 |
| `Minimise` | 將視窗最小化 | 所有平台 |
| `Zoom` | 縮放視窗 | 僅限 macOS |
| `Front` | 將視窗移至最前方 | 僅限 macOS |
| `Quit` | 結束應用程式 | 所有平台 |
| `Undo` | 復原上一個動作 | 所有平台 |
| `Redo` | 重做上一個動作 | 所有平台 |
| `Cut` | 剪下選取內容 | 所有平台 |
| `Copy` | 複製選取內容 | 所有平台 |
| `Paste` | 從剪貼簿貼上 | 所有平台 |
| `PasteAndMatchStyle` | 貼上並符合樣式 | 僅限 macOS |
| `SelectAll` | 全選 | 所有平台 |
| `Delete` | 刪除選取內容 | 所有平台 |
| `Reload` | 重新載入目前頁面 | 所有平台 |
| `ForceReload` | 強制重新載入目前頁面 | 所有平台 |
| `ToggleFullscreen` | 切換全螢幕模式 | 所有平台 |
| `ResetZoom` | 重設縮放比例 | 所有平台 |
| `ZoomIn` | 放大 | 所有平台 |
| `ZoomOut` | 縮小 | 所有平台 |

以下範例示範如何同時使用完整選單和個別角色：

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## 應用程式選單

應用程式選單會顯示在應用程式視窗頂端（Windows/Linux），或螢幕頂端（macOS）。

### 應用程式選單行為

使用`app.Menu.Set()`設定應用程式選單後，該選單會成為 macOS 上的主選單。 在 Windows/Linux 上，選單則按視窗個別設定。

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

以下完整範例示範這些不同的選單行為：

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## 內容選單

內容選單是在應用程式中以滑鼠右鍵按一下元素時出現的彈出式選單，可讓使用者快速存取與所按元素相關的動作。

### 預設內容選單

預設內容選單是 WebView 的內建內容選單，提供下列系統層級操作：

- 用於文字處理的複製、剪下和貼上
- 文字選取控制項
- 拼字檢查選項

#### 控制預設內容選單

您可以使用`--default-contextmenu` CSS 屬性控制預設內容選單何時出現：

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
此功能只有在[前端執行階段準備就緒](/reference/frontend-runtime/)後，才會如預期運作。

@end

#### 巢狀內容選單的行為

在巢狀元素上使用`--default-contextmenu`屬性時，適用下列規則：

1. 除非明確覆寫，否則子元素會繼承其父元素的內容選單設定
2. 最具體（距離最近）的設定優先
3. 可使用`auto`值重設為預設行為

巢狀內容選單行為範例：

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### 自訂內容選單

自訂內容選單可提供與所按元素相關的應用程式專屬動作，特別適合用於：

- 文件管理器中的檔案操作
- 影像處理工具
- 資料格中的自訂動作
- 元件特定操作

#### 建立自訂內容功能表

建立自訂內容功能表時，請提供一個唯一識別碼（名稱），以便將功能表連結至 HTML 元素：

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

name 參數（本例為「imageMenu」）是唯一識別碼，用於：

1. 將 HTML 元素連結至這個特定的內容功能表
2. 識別按一下滑鼠右鍵時應顯示的功能表
3. 允許更新及清理功能表

#### 內容資料

處理內容功能表事件時，您可以存取被按下的功能表項目及其關聯的內容資料：

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

內容資料會從 HTML 元素的`--custom-contextmenu-data`屬性傳入，並可在點擊處理常式中透過`ctx.ContextMenuData()`取得。這在下列情況中特別實用：

- 處理每個項目都需要唯一識別的清單或資料格
- 對特定元件或元素執行操作
- 將狀態或中繼資料從前端傳遞至後端

#### 內容功能表管理

變更內容功能表後，請呼叫`Update()`方法以套用變更：

```go
contextMenu.Update()
```

不再需要內容功能表時，可以將其銷毀：

```go
contextMenu.Destroy()
```

@note{type="danger" title="警告"}
呼叫`Destroy()`後，再次使用該內容功能表的參照將導致 panic。

@end

### 實務範例：圖片庫

以下是在圖片庫中實作自訂內容功能表的完整範例：

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

在此範例中：

1. 內容功能表使用識別碼「imageMenu」建立
2. 每個圖片容器都使用`--custom-contextmenu: imageMenu`連結至功能表
3. 每個容器都使用`--custom-contextmenu-data`將其圖片 ID 作為內容資料提供
4. 後端會在點擊處理常式中接收圖片 ID，並可執行特定操作
5. 所有圖片都重複使用同一個功能表，但內容資料會指出要操作哪張圖片

此模式特別適合：

- 需要對各資料列執行特定操作的資料格
- 需要對檔案執行內容相關動作的檔案管理員
- 需要對不同元素執行不同操作的設計工具
- 將相同操作套用至多個執行個體的任何元件
