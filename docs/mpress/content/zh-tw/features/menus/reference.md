---
title: "選單參考"
description: "選單項目類型、屬性與方法的完整參考"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## 選單參考

選單項目類型、屬性與動態行為的完整參考。使用核取方塊、選項群組、分隔線與動態更新，建立專業且反應靈敏的選單。

## 選單項目類型

### 一般選單項目

最常見的類型——顯示文字並觸發動作：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

<strong>適用於：</strong>命令、動作、開啟視窗

### 核取方塊

可切換勾選／未勾選狀態的選單項目：

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

<strong>適用於：</strong>布林設定、功能開關、檢視選項

<strong>重要：</strong>按一下後，勾選狀態會自動切換。

### 選項群組

互斥選項——一次只能選取一個：

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

<strong>適用於：</strong>互斥選擇（大小、佈景主題、模式）

**群組運作方式：**

- 相鄰的選項項目會自動形成群組
- 選取其中一個項目時，會取消選取群組中的其他項目
- 使用分隔線或一般項目分隔不同群組

**多個群組的範例：**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### 子選單

用於組織內容的巢狀選單結構：

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

<strong>適用於：</strong>將相關項目分組、減少雜亂

<strong>巢狀層級限制：</strong>大多數平台支援2-3層。請避免更深的巢狀結構。

### 分隔線

選單項目之間的視覺分隔元素：

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

<strong>適用於：</strong>在視覺上將相關項目分組

<strong>最佳做法：</strong>不要在選單開頭或結尾放置分隔線。

## 選單項目屬性

### 標籤

選單項目顯示的文字：

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**動態標籤：**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### 啟用狀態

控制選單項目是否可供互動：

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Windows 選單行為"}
在 Windows 上，選單狀態變更時必須重新建構。啟用或停用選單項目後<strong>一律呼叫`menu.Update()`</strong>，尤其是該項目建立時處於停用狀態的情況。

<strong>原因：</strong>Windows 選單更新時會從頭重新建構。如果未呼叫`Update()`，點擊處理常式將無法正常觸發。

@end

**範例：動態啟用／停用**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**常見模式：符合條件時啟用**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### 勾選狀態

對於核取方塊和選項項目，可控制或查詢其勾選狀態：

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

<strong>自動切換：</strong>按一下核取方塊後，其狀態會自動切換。您不需要在點擊處理常式中呼叫`SetChecked()`。

**手動控制：**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### 快速鍵（鍵盤快捷鍵）

為選單項目新增鍵盤快捷鍵：

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**快速鍵格式：**

- `CmdOrCtrl`——在 macOS 上為 Cmd，在 Windows/Linux 上為 Ctrl
- `Shift`、`Alt`、`Option`——輔助按鍵
- `A-Z`、`0-9`——字母／數字鍵
- `F1-F12`——功能鍵
- `Enter`、`Space`、`Backspace`等——特殊按鍵

**範例：**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**平台特定的快速鍵：**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### 工具提示

為選單項目新增滑鼠懸停文字（各平台的支援情況不一）：

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**平台支援情況：**

- **Windows：**✅ 支援
- **macOS：**❌ 不支援（工具提示並非選單的標準功能）
- **Linux：**⚠️ 視桌面環境而異

### 隱藏狀態

隱藏選單項目而不移除：

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

<strong>適用於：</strong>偵錯選項、功能旗標、條件式功能

## 事件處理

### OnClick 處理常式

按一下選單項目時執行程式碼：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Context 提供：**

- `ctx.ClickedMenuItem()`－被按下的選單項目
- 視窗上下文（若來自視窗選單）
- 應用程式上下文

**範例：在處理常式中存取選單項目**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### 多個處理常式

您可以設定多個處理常式（以最後一個為準）：

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

<strong>最佳實務：</strong>只設定一次處理常式，並視需要在其中使用條件邏輯。

## 動態選單

### 更新選單項目

<strong>黃金準則：</strong>變更選單狀態後，一律呼叫`menu.Update()`。

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**這一點很重要的原因：**

- <strong>Windows：</strong>選單會在更新時重新建構
- <strong>macOS/Linux：</strong>重要性較低，但仍建議這麼做
- <strong>按一下處理常式：</strong>若未呼叫 Update()，便無法正常觸發

### 重建選單

若有重大變更，請重建整個選單：

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**何時應重建：**

- 最近使用的檔案清單變更
- 外掛程式選單變更
- 重大狀態轉換

**何時應更新：**

- 啟用或停用項目
- 變更標籤
- 切換核取方塊

### 內容相關選單

根據應用程式狀態調整選單：

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## 平台差異

### 選單列位置

| 平台 | 位置 | 備註 |
| --- | --- | --- |
| **macOS** | 螢幕頂端 | 全域選單列 |
| **Windows** | 視窗頂端 | 各視窗各自具備選單 |
| **Linux** | 視窗頂端 | 通常各視窗各自具備選單 |

### 標準選單

**macOS：**

- 具備「應用程式」選單（顯示應用程式名稱）
- 「偏好設定」位於「應用程式」選單中
- 「結束」位於「應用程式」選單中

**Windows/Linux：**

- 沒有「應用程式」選單
- 「偏好設定」位於「編輯」或「工具」選單中
- 「結束」位於「檔案」選單中

**範例：符合平台慣例的結構**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### 快速鍵慣例

**macOS：**

- 大多數快速鍵使用`Cmd+`
- 「偏好設定」使用`Cmd+,`
- 「結束」使用`Cmd+Q`

**Windows：**

- 大多數快速鍵使用`Ctrl+`
- 「偏好設定」使用`Ctrl+P`或`Ctrl+,`
- 「結束」使用`Alt+F4`（或`Ctrl+Q`）

**Linux：**

- 通常遵循 Windows 慣例
- 桌面環境可能會覆寫這些設定

## 最佳實務

### ✅ 建議做法

- 變更選單狀態後，**呼叫 menu.Update()**（尤其是在 Windows 上）
- 對互斥選項<strong>使用選項群組</strong>
- 對可切換的功能<strong>使用核取方塊</strong>
- **為常用動作新增快速鍵**
- **使用分隔線將相關項目分組**
- **在所有平台上測試**——行為因平台而異

### ❌ 請勿

- **別忘了呼叫 menu.Update()**——否則點擊處理函式將無法正常運作
- **請勿巢狀排列過深**——最多 2-3 層
- **請勿以分隔線開頭或結尾**——這樣看起來不專業
- **請勿在 macOS 上使用工具提示**——macOS 不支援此功能
- **請勿寫死特定平台的快速鍵**——請使用 `CmdOrCtrl`

## 疑難排解

### 選單項目沒有回應

<strong>症狀：</strong>點擊處理函式未觸發

<strong>原因：</strong>啟用項目後忘了呼叫 `menu.Update()`

**解決方法：**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### 選單項目呈灰色

<strong>症狀：</strong>無法點擊選單項目

<strong>原因：</strong>項目已停用

**解決方法：**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### 快速鍵無法使用

<strong>症狀：</strong>鍵盤快速鍵未觸發選單項目

**原因：**

1. 快速鍵格式不正確
2. 與系統快速鍵衝突
3. 視窗未取得焦點

**解決方法：**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## 後續步驟

- [應用程式選單](/features/menus/application/)——建立應用程式選單列
- [內容選單](/features/menus/context/)——按一下滑鼠右鍵開啟的內容選單
- [系統匣選單](/features/menus/systray/)——系統匣／選單列選單
- [選單模式](/guides/menus/)——常見的選單模式與最佳實務

---

<strong>有問題嗎？</strong>請在 [Discord](https://discord.gg/JDdSxwjhGf) 中提問，或查看[選單範例](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)。
