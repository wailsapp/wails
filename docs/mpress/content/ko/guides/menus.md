---
title: "메뉴"
description: "Wails v3에서 메뉴를 만들고 사용자 지정하는 방법을 설명하는 가이드"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3는 애플리케이션 메뉴와 컨텍스트 메뉴를 모두 만들 수 있는 강력한 메뉴 시스템을 제공합니다. 이 가이드에서는 메뉴 시스템의 다양한 기능을 설명합니다.

## 메뉴 만들기

새 메뉴를 만들려면 Menus 관리자의 `New()` 메서드를 사용하세요:

```go
menu := app.Menu.New()
```

### 메뉴 항목 추가하기

Wails는 각각 특정 용도로 사용하는 여러 유형의 메뉴 항목을 지원합니다:

#### 일반 메뉴 항목

일반 메뉴 항목은 메뉴를 구성하는 기본 요소입니다. 텍스트를 표시하며 클릭하면 동작을 실행할 수 있습니다:

```go
menuItem := menu.Add("Click Me")
```

#### 체크박스

체크박스 메뉴 항목은 전환 가능한 상태를 제공하며, 기능이나 설정을 활성화하거나 비활성화할 때 유용합니다:

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### 라디오 그룹

라디오 그룹을 사용하면 상호 배타적인 선택지 중 하나를 선택할 수 있습니다. 라디오 항목을 서로 나란히 배치하면 자동으로 그룹이 만들어집니다:

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### 구분선

구분선은 메뉴 항목을 논리적인 그룹으로 구성하는 데 도움이 되는 가로선입니다:

```go
menu.AddSeparator()
```

#### 하위 메뉴

하위 메뉴는 메뉴 항목에 마우스 포인터를 올리거나 메뉴 항목을 클릭하면 나타나는 중첩 메뉴입니다. 복잡한 메뉴 구조를 구성할 때 유용합니다:

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### 메뉴 결합하기

메뉴를 다른 메뉴의 뒤나 앞에 추가할 수 있습니다.

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
기본적으로 `prepend`과 `append`은 원본 메뉴와 상태를 공유합니다. 자체 상태를 가진 새 메뉴를 만들려면 메뉴에서 `.Clone()`을 호출하세요.

예: `menu.Append(secondaryMenu.Clone())`

@end

#### 메뉴 비우기

메뉴 항목 수가 가변적이라면 경우에 따라 완전히 새로운 메뉴를 구성하는 편이 더 나을 수 있습니다.

이렇게 하면 기존 메뉴의 모든 항목을 비우고 항목을 다시 추가할 수 있습니다.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
메뉴를 비우면 최상위 수준의 메뉴 항목만 제거됩니다. 하위 메뉴는 표시되지 않지만 여전히 메모리를 차지하므로 메뉴를 신중하게 관리하세요.

@end

#### 메뉴 제거하기

메뉴를 비우고 리소스까지 해제하려면 `Destroy()` 메서드를 사용하세요:

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### 메뉴 항목 속성

메뉴 항목에는 구성할 수 있는 여러 속성이 있습니다:

| 속성 | 메서드 | 설명 |
| --- | --- | --- |
| 레이블 | `SetLabel(string)` | 표시할 텍스트를 설정합니다 |
| 활성화 여부 | `SetEnabled(bool)` | 항목을 활성화하거나 비활성화합니다 |
| 선택 여부 | `SetChecked(bool)` | 선택 상태를 설정합니다(체크박스/라디오 항목용) |
| 도구 설명 | `SetTooltip(string)` | 도구 설명 텍스트를 설정합니다 |
| 숨김 여부 | `SetHidden(bool)` | 항목을 표시하거나 숨깁니다 |
| 단축키 | `SetAccelerator(string)` | 키보드 단축키를 설정합니다 |

### 메뉴 항목 상태

메뉴 항목은 표시 여부와 상호작용 가능 여부를 제어하는 여러 상태를 가질 수 있습니다:

#### 표시 여부

`SetHidden()` 메서드를 사용하여 메뉴 항목을 동적으로 표시하거나 숨길 수 있습니다:

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

숨겨진 메뉴 항목은 다시 표시할 때까지 메뉴에서 완전히 제거됩니다. 이는 애플리케이션이 특정 상태일 때만 나타나야 하는 상황별 메뉴 항목에 유용합니다.

#### 활성화 상태

`SetEnabled()` 메서드를 사용하여 메뉴 항목을 활성화하거나 비활성화할 수 있습니다:

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

비활성화된 메뉴 항목은 계속 표시되지만 회색으로 나타나며 클릭할 수 없습니다. 이는 다음과 같이 현재 동작을 사용할 수 없음을 나타내는 데 흔히 사용됩니다:

- 저장할 변경 사항이 없을 때 "저장" 비활성화
- 선택된 항목이 없을 때 "복사" 비활성화
- 실행 취소할 동작이 없을 때 "실행 취소" 비활성화

#### 동적 상태 관리

이러한 상태를 이벤트 핸들러와 결합하여 동적 메뉴를 만들 수 있습니다:

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

### 이벤트 처리

메뉴 항목은 `OnClick` 메서드를 사용하여 클릭 이벤트에 응답할 수 있습니다.

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

컨텍스트는 클릭한 메뉴 항목에 관한 정보를 제공합니다.

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### 역할 기반 메뉴 항목

Wails는 표준 기능이 있는 메뉴 항목을 자동으로 생성하는 사전 정의된 메뉴 역할을 제공합니다. 지원되는 메뉴 역할은 다음과 같습니다.

#### 전체 메뉴 구조

이러한 역할은 일반적인 기능을 갖춘 전체 메뉴 구조를 생성합니다.

| 역할 | 설명 | 플랫폼 참고 사항 |
| --- | --- | --- |
| `AppMenu` | 정보, 서비스, 숨기기/표시 및 종료 항목이 있는 애플리케이션 메뉴 | macOS 전용 |
| `EditMenu` | 실행 취소, 다시 실행, 잘라내기, 복사, 붙여넣기 등이 포함된 표준 편집 메뉴 | 모든 플랫폼 |
| `ViewMenu` | 새로고침, 확대/축소 및 전체 화면 제어 기능이 있는 보기 메뉴 | 모든 플랫폼 |
| `WindowMenu` | 창 제어 기능(최소화, 확대/축소 등) | 모든 플랫폼 |
| `HelpMenu` | Wails 웹사이트로 연결되는 "자세히 알아보기" 링크가 있는 도움말 메뉴 | 모든 플랫폼 |

#### 개별 메뉴 항목

이러한 역할을 사용하여 개별 메뉴 항목을 추가할 수 있습니다.

| 역할 | 설명 | 플랫폼 참고 사항 |
| --- | --- | --- |
| `About` | 애플리케이션 정보 대화 상자 표시 | 모든 플랫폼 |
| `Hide` | 애플리케이션 숨기기 | macOS 전용 |
| `HideOthers` | 다른 애플리케이션 숨기기 | macOS 전용 |
| `UnHide` | 숨겨진 애플리케이션 표시 | macOS 전용 |
| `CloseWindow` | 현재 창 닫기 | 모든 플랫폼 |
| `Minimise` | 창 최소화 | 모든 플랫폼 |
| `Zoom` | 창 확대/축소 | macOS 전용 |
| `Front` | 창을 맨 앞으로 가져오기 | macOS 전용 |
| `Quit` | 애플리케이션 종료 | 모든 플랫폼 |
| `Undo` | 마지막 작업 실행 취소 | 모든 플랫폼 |
| `Redo` | 마지막 작업 다시 실행 | 모든 플랫폼 |
| `Cut` | 선택 영역 잘라내기 | 모든 플랫폼 |
| `Copy` | 선택 영역 복사 | 모든 플랫폼 |
| `Paste` | 클립보드에서 붙여넣기 | 모든 플랫폼 |
| `PasteAndMatchStyle` | 스타일에 맞춰 붙여넣기 | macOS 전용 |
| `SelectAll` | 모두 선택 | 모든 플랫폼 |
| `Delete` | 선택 영역 삭제 | 모든 플랫폼 |
| `Reload` | 현재 페이지 새로고침 | 모든 플랫폼 |
| `ForceReload` | 현재 페이지 강제로 새로고침 | 모든 플랫폼 |
| `ToggleFullscreen` | 전체 화면 모드 전환 | 모든 플랫폼 |
| `ResetZoom` | 확대/축소 수준 초기화 | 모든 플랫폼 |
| `ZoomIn` | 확대 | 모든 플랫폼 |
| `ZoomOut` | 축소 | 모든 플랫폼 |

다음은 완전한 메뉴와 개별 역할을 모두 사용하는 방법을 보여 주는 예제입니다:

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

## 애플리케이션 메뉴

애플리케이션 메뉴는 애플리케이션 창 상단(Windows/Linux) 또는 화면 상단(macOS)에 표시되는 메뉴입니다.

### 애플리케이션 메뉴 동작

`app.Menu.Set()`을 사용하여 애플리케이션 메뉴를 설정하면 macOS에서는 기본 메뉴가 됩니다. Windows/Linux에서는 창별로 메뉴가 설정됩니다.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

다음은 이러한 메뉴의 서로 다른 동작을 보여 주는 전체 예제입니다:

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

## 컨텍스트 메뉴

컨텍스트 메뉴는 애플리케이션의 요소를 마우스 오른쪽 버튼으로 클릭할 때 표시되는 팝업 메뉴입니다. 클릭한 요소와 관련된 작업에 빠르게 접근할 수 있습니다.

### 기본 컨텍스트 메뉴

기본 컨텍스트 메뉴는 다음과 같은 시스템 수준 작업을 제공하는 웹뷰의 내장 컨텍스트 메뉴입니다:

- 텍스트 편집을 위한 복사, 잘라내기 및 붙여넣기
- 텍스트 선택 컨트롤
- 맞춤법 검사 옵션

#### 기본 컨텍스트 메뉴 제어

`--default-contextmenu` CSS 속성을 사용하여 기본 컨텍스트 메뉴가 표시되는 시점을 제어할 수 있습니다:

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
이 기능은 [프런트엔드 런타임이 준비된 후에만](/reference/frontend-runtime/) 예상대로 작동합니다.

@end

#### 중첩된 컨텍스트 메뉴 동작

중첩된 요소에 `--default-contextmenu` 속성을 사용할 때는 다음 규칙이 적용됩니다:

1. 명시적으로 재정의하지 않는 한 자식 요소는 부모 요소의 컨텍스트 메뉴 설정을 상속합니다
2. 가장 구체적인(가장 가까운) 설정이 우선합니다
3. `auto` 값을 사용하여 기본 동작으로 재설정할 수 있습니다

중첩된 컨텍스트 메뉴 동작의 예제:

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

### 사용자 정의 컨텍스트 메뉴

사용자 정의 컨텍스트 메뉴를 사용하면 클릭한 요소와 관련된 애플리케이션별 작업을 제공할 수 있습니다. 특히 다음 용도에 유용합니다:

- 문서 관리자의 파일 작업
- 이미지 조작 도구
- 데이터 그리드의 사용자 정의 작업
- 컴포넌트별 작업

#### 사용자 지정 컨텍스트 메뉴 만들기

사용자 지정 컨텍스트 메뉴를 만들 때는 메뉴를 HTML 요소와 연결하는 고유 식별자(이름)를 지정합니다:

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

name 매개변수(이 예제에서는 "imageMenu")는 다음 용도로 사용되는 고유 식별자입니다:

1. HTML 요소를 이 특정 컨텍스트 메뉴에 연결
2. 마우스 오른쪽 버튼을 클릭할 때 표시할 메뉴 식별
3. 메뉴 업데이트 및 정리 지원

#### 컨텍스트 데이터

컨텍스트 메뉴 이벤트를 처리할 때는 클릭한 메뉴 항목과 연결된 컨텍스트 데이터에 모두 접근할 수 있습니다:

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

컨텍스트 데이터는 HTML 요소의 `--custom-contextmenu-data` 속성에서 전달되며 클릭 핸들러에서 `ctx.ContextMenuData()`을 통해 사용할 수 있습니다. 이는 특히 다음과 같은 경우에 유용합니다:

- 각 항목을 고유하게 식별해야 하는 목록이나 그리드로 작업하는 경우
- 특정 컴포넌트나 요소에 대한 작업을 처리하는 경우
- 프런트엔드에서 백엔드로 상태나 메타데이터를 전달하는 경우

#### 컨텍스트 메뉴 관리

컨텍스트 메뉴를 변경한 후에는 `Update()` 메서드를 호출하여 변경 사항을 적용합니다:

```go
contextMenu.Update()
```

컨텍스트 메뉴가 더 이상 필요하지 않으면 제거할 수 있습니다:

```go
contextMenu.Destroy()
```

@note{type="danger" title="경고"}
`Destroy()`을 호출한 후 컨텍스트 메뉴 참조를 다시 사용하면 패닉이 발생합니다.

@end

### 실제 활용 예제: 이미지 갤러리

다음은 이미지 갤러리에 사용자 지정 컨텍스트 메뉴를 구현하는 전체 예제입니다:

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

이 예제에서는 다음과 같이 동작합니다:

1. "imageMenu" 식별자로 컨텍스트 메뉴를 만듭니다.
2. `--custom-contextmenu: imageMenu`을 사용하여 각 이미지 컨테이너를 메뉴에 연결합니다.
3. 각 컨테이너는 `--custom-contextmenu-data`을 사용하여 해당 이미지 ID를 컨텍스트 데이터로 제공합니다.
4. 백엔드는 클릭 핸들러에서 이미지 ID를 수신하고 해당 이미지에 맞는 작업을 수행할 수 있습니다.
5. 모든 이미지에 동일한 메뉴를 재사용하지만, 컨텍스트 데이터를 통해 작업할 이미지를 식별할 수 있습니다.

이 패턴은 특히 다음과 같은 경우에 유용합니다:

- 각 행에 맞는 작업이 필요한 데이터 그리드
- 각 파일의 컨텍스트에 맞는 작업이 필요한 파일 관리자
- 요소마다 서로 다른 작업이 필요한 디자인 도구
- 동일한 작업을 여러 인스턴스에 적용하는 모든 컴포넌트
