---
title: "컨텍스트 메뉴"
description: "애플리케이션용 오른쪽 클릭 컨텍스트 메뉴 만들기"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## 문제

사용자는 상황에 맞는 작업을 제공하는 오른쪽 클릭 메뉴를 기대합니다. 요소마다 필요한 메뉴가 다릅니다.

- **텍스트**: 잘라내기, 복사, 붙여넣기
- **이미지**: 저장, 복사, 열기
- **사용자 지정 요소**: 애플리케이션별 작업

컨텍스트 메뉴를 직접 만들려면 마우스 이벤트, 위치 지정 및 플랫폼 차이를 처리해야 합니다.

## Wails 솔루션

Wails는 CSS 속성을 사용하는 <strong>선언형 컨텍스트 메뉴</strong>를 제공합니다. 메뉴를 HTML 요소와 연결하고 데이터를 전달하며 클릭을 처리할 수 있으며, 이 모든 기능이 플랫폼 네이티브 동작으로 작동합니다.

![macOS에서 WebView 위에 표시된 Wails 사용자 지정 컨텍스트 메뉴](/assets/screenshots/context-menu-macos.png)

메뉴는 플랫폼 네이티브이지만, 메뉴를 연 요소는 WebView의 일부로 유지됩니다. 이 macOS 화면은 아래 예제의 사용자 지정 컨텍스트 메뉴 등록을 사용합니다.

## 빠른 시작

**Go 코드:**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML:**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**이것으로 끝입니다!** textarea를 오른쪽 클릭하면 사용자 지정 메뉴가 표시됩니다.

## 컨텍스트 메뉴 만들기

### 기본 컨텍스트 메뉴

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

**메뉴 ID:** 고유해야 합니다. 메뉴를 HTML 요소와 연결하는 데 사용합니다.

### 하위 메뉴 사용

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### 체크박스 및 라디오 그룹 사용

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

<strong>모든 메뉴 항목 유형</strong>은 [메뉴 레퍼런스](/features/menus/reference/)를 참조하세요.

## HTML 요소와 연결하기

CSS 사용자 지정 속성을 사용하여 컨텍스트 메뉴를 연결하세요.

### 기본 연결

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**CSS 속성:** `--custom-contextmenu: <menu-id>`

### 컨텍스트 데이터 사용

HTML에서 Go로 데이터를 전달하세요.

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Go 핸들러:**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**CSS 속성:**

- `--custom-contextmenu: <menu-id>` - 표시할 메뉴
- `--custom-contextmenu-data: <data>` - 핸들러에 전달할 데이터

### 동적 데이터

JavaScript에서 데이터를 동적으로 생성하세요.

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### 여러 요소에서 같은 메뉴 사용

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**하나의 메뉴를 사용하되 요소마다 다른 데이터를 전달합니다.**

## 컨텍스트 데이터

### 컨텍스트 데이터에 접근하기

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**데이터 유형:** 항상 `string`입니다. 필요에 따라 파싱하세요.

### 복합 데이터 전달하기

복합 데이터에는 JSON을 사용하세요.

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Go 핸들러:**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="보안"}
프런트엔드에서 전달된 **컨텍스트 데이터는 항상 검증하세요**. 사용자가 CSS 속성을 조작할 수 있으므로 이 데이터를 신뢰할 수 없는 입력으로 취급하세요.

@end

### 검증 예제

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## 기본 컨텍스트 메뉴

WebView는 표준 작업(복사, 붙여넣기, 검사)을 위한 내장 컨텍스트 메뉴를 제공합니다. `--default-contextmenu`을 사용하여 제어하세요.

### 기본 메뉴 숨기기

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**사용 사례:** 기본 메뉴가 적합하지 않은 사용자 지정 UI 요소.

### 기본 메뉴 표시하기

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**사용 사례:** 텍스트 영역, 입력 필드, 편집 가능한 콘텐츠.

### 자동(스마트) 모드

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**기본 동작입니다.** 다음 경우에 기본 메뉴를 표시합니다.

- 텍스트가 선택된 경우
- 텍스트 입력 필드인 경우
- 편집 가능한 콘텐츠(`contenteditable`)인 경우

그 외의 경우에는 기본 메뉴를 숨깁니다.

### 사용자 지정 메뉴와 기본 메뉴 함께 사용하기

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**동작:**

1. 사용자 지정 메뉴가 먼저 표시됩니다.
2. 사용자 지정 메뉴가 비어 있거나 없으면 기본 메뉴가 표시됩니다.
3. 두 메뉴가 함께 표시될 수 있습니다(플랫폼에 따라 다름).

## 동적 컨텍스트 메뉴

애플리케이션 상태에 따라 메뉴를 업데이트하세요.

### 항목 활성화/비활성화하기

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="항상 Update() 호출"}
메뉴 상태를 변경한 후에는 <strong>반드시 `contextMenu.Update()`</strong>을 호출하세요. 이는 Windows에서 매우 중요합니다.

자세한 내용은 [메뉴 레퍼런스](/features/menus/reference/#enabled-state)를 참조하세요.

@end

### 레이블 변경

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### 메뉴 다시 빌드

변경 사항이 큰 경우 전체 메뉴를 다시 빌드하세요:

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## 플랫폼별 동작

컨텍스트 메뉴는 **플랫폼 네이티브** 방식으로 작동합니다:

@tabs{sync-key="platform"}
[macOS]
**macOS 네이티브 컨텍스트 메뉴:**

- 시스템 애니메이션 및 전환 효과
- 마우스 오른쪽 클릭 = Control+Click(자동)
- 시스템 화면 모드(라이트/다크)에 맞게 조정
- 기본 메뉴에서 표준 텍스트 작업 제공
- 긴 메뉴에 네이티브 스크롤 적용

**macOS 규칙:**

- 메뉴 항목에는 문장형 대소문자 표기 사용
- 대화 상자를 여는 항목에는 줄임표(...) 사용
- 일반적인 단축키: ⌘C(복사), ⌘V(붙여넣기)

[Windows]
**Windows 네이티브 컨텍스트 메뉴:**

- Windows 네이티브 스타일
- Windows 테마를 따름
- 기본 메뉴에서 표준 Windows 작업 제공
- 터치 및 펜 입력 지원

**Windows 규칙:**

- 메뉴 항목에는 제목형 대소문자 표기 사용
- 대화 상자를 여는 항목에는 줄임표(...) 사용
- 일반적인 단축키: Ctrl+C(복사), Ctrl+V(붙여넣기)

[Linux]
**데스크톱 환경 통합:**

- 데스크톱 테마(GTK, Qt 등)에 맞게 조정
- 마우스 오른쪽 클릭 동작은 시스템 설정을 따름
- 기본 메뉴 내용은 환경에 따라 다름
- 위치는 데스크톱 환경(DE) 규칙을 따름

**Linux 고려 사항:**

- 대상 데스크톱 환경에서 테스트
- GTK와 Qt는 서로 다르게 동작함
- 일부 데스크톱 환경(DE)은 컨텍스트 메뉴를 사용자 지정함

@end

## 전체 예제

**Go 코드:**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## 모범 사례

### ✅ 권장 사항

- **메뉴의 초점을 명확하게 유지하세요** - 해당 요소와 관련된 작업만 포함하세요
- **컨텍스트 데이터를 검증하세요** - 신뢰할 수 없는 입력으로 취급하세요
- **명확한 레이블을 사용하세요** - "삭제"가 아닌 "파일 삭제"를 사용하세요
- **menu.Update()를 호출하세요** - 메뉴 상태를 변경한 후 호출하세요
- **모든 플랫폼에서 테스트하세요** - 동작이 플랫폼마다 다릅니다
- **키보드 단축키를 제공하세요** - 자주 사용하는 작업에 제공하세요
- **관련 항목을 그룹화하세요** - 구분선을 사용하세요

### ❌ 금지 사항

- **컨텍스트 데이터를 신뢰하지 마세요** - 항상 검증하세요
- **메뉴를 너무 길게 만들지 마세요** - 최대 7-10개 항목
- **menu.Update() 호출을 잊지 마세요** - 그렇지 않으면 메뉴가 제대로 작동하지 않습니다
- **메뉴를 너무 깊게 중첩하지 마세요** - 최대 2단계
- **전문 용어를 사용하지 마세요** - 사용자가 이해하기 쉬운 레이블을 사용하세요
- **핸들러를 블로킹하지 마세요** - 빠르게 실행되도록 유지하세요

## 문제 해결

### 컨텍스트 메뉴가 표시되지 않음

**가능한 원인:**

1. 메뉴 ID 불일치
2. CSS 속성 오타
3. 런타임이 초기화되지 않음

**해결 방법:**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### 컨텍스트 데이터가 수신되지 않음

**가능한 원인:**

1. CSS 속성이 설정되지 않음
2. 데이터에 특수 문자가 포함됨

**해결 방법:**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

또는 JavaScript를 사용하세요:

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### 메뉴 항목이 반응하지 않음

**원인:** 활성화한 후 `menu.Update()` 호출을 누락함

**해결 방법:**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## 다음 단계

@cards{cols="2"}
📖 메뉴 참조
메뉴 항목 유형과 속성에 대한 전체 참조 문서입니다.

[자세히 알아보기 →](/features/menus/reference/)

---
☰ 애플리케이션 메뉴
애플리케이션 메뉴 모음을 만듭니다.

[자세히 알아보기 →](/features/menus/application/)

---
★ 시스템 트레이 메뉴
시스템 트레이/메뉴 막대 연동 기능을 추가합니다.

[자세히 알아보기 →](/features/menus/systray/)

---
📖 메뉴 패턴
일반적인 메뉴 패턴과 모범 사례를 소개합니다.

[자세히 알아보기 →](/guides/menus/)

@end

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [컨텍스트 메뉴 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus)를 확인하세요.
