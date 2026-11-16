---
title: "메뉴 레퍼런스"
description: "메뉴 항목 유형, 속성 및 메서드에 대한 전체 레퍼런스"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## 메뉴 레퍼런스

메뉴 항목 유형, 속성 및 동적 동작에 대한 전체 레퍼런스입니다. 체크박스, 라디오 그룹, 구분선 및 동적 업데이트를 사용하여 전문적이고 반응성이 뛰어난 메뉴를 구축하세요.

## 메뉴 항목 유형

### 일반 메뉴 항목

가장 일반적인 유형으로, 텍스트를 표시하고 동작을 실행합니다.

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**용도:** 명령, 동작, 창 열기

### 체크박스

선택됨/선택되지 않음 상태를 전환할 수 있는 메뉴 항목입니다.

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**용도:** 불리언 설정, 기능 전환, 보기 옵션

**중요:** 클릭하면 선택 상태가 자동으로 전환됩니다.

### 라디오 그룹

상호 배타적인 옵션으로, 하나만 선택할 수 있습니다.

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**용도:** 상호 배타적인 선택 항목(크기, 테마, 모드)

**그룹화 방식:**

- 서로 인접한 라디오 항목은 자동으로 그룹을 구성합니다.
- 하나를 선택하면 그룹의 다른 항목은 선택 해제됩니다.
- 구분선이나 일반 항목으로 그룹을 분리하세요.

**여러 그룹을 사용하는 예:**

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

### 하위 메뉴

항목을 체계적으로 구성하기 위한 중첩 메뉴 구조입니다.

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**용도:** 관련 항목 그룹화, 복잡성 완화

**중첩 제한:** 대부분의 플랫폼은 2-3단계를 지원합니다. 이보다 깊게 중첩하지 마세요.

### 구분선

메뉴 항목 사이를 시각적으로 나누는 요소입니다.

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**용도:** 관련 항목을 시각적으로 그룹화

**권장 사항:** 메뉴의 시작이나 끝에 구분선을 배치하지 마세요.

### 섹션 머리글

<strong>macOS 14+</strong>에서 메뉴 항목 그룹에 사용하는 네이티브 비대화형 레이블입니다.

```go
menu.AddSectionHeader("Build server")
menu.Add("Start build").OnClick(startBuild)
menu.Add("View logs").OnClick(viewLogs)
```

`AddSectionHeader`은 AppKit의 표준 섹션 머리글 모양을 사용합니다. 메뉴 항목을 반환하지 않으며 동작도 연결하지 않습니다. 이전 macOS 버전과 Windows, Linux, 모바일 및 서버 빌드에서는 이 호출이 아무 작업도 하지 않으며 행도 삽입하지 않습니다.

트레이, 컨텍스트 또는 애플리케이션 메뉴에서 관련 명령에 레이블을 지정할 때 사용하세요. 생략된 머리글은 라디오 그룹을 분리하지 않으므로 모든 플랫폼에서 라디오 그룹을 분리하려면 `AddSeparator()`을 사용하세요. 메서드 시그니처는 [메뉴 API](/reference/menu#addsectionheader)를, 호스트 이름 머리글은 `systray-menu` 예제를 참조하세요.

## 메뉴 항목 속성

### 레이블

메뉴 항목에 표시되는 텍스트입니다.

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**동적 레이블:**

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

### 활성화 상태

메뉴 항목과 상호 작용할 수 있는지를 제어합니다.

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Windows 메뉴 동작"}
Windows에서는 메뉴 상태가 변경될 때 메뉴를 다시 구성해야 합니다. 메뉴 항목을 활성화하거나 비활성화한 후에는 **항상 `menu.Update()`을 호출하세요**. 특히 항목이 비활성화된 상태로 생성된 경우에는 반드시 호출해야 합니다.

**이유:** Windows 메뉴는 업데이트할 때 처음부터 다시 빌드됩니다. `Update()`을 호출하지 않으면 클릭 핸들러가 제대로 실행되지 않습니다.

@end

**예: 동적 활성화/비활성화**

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

**일반적인 패턴: 조건에 따라 활성화**

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

### 선택 상태

체크박스 및 라디오 항목의 선택 상태를 제어하거나 조회합니다.

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**자동 전환:** 체크박스는 클릭하면 자동으로 전환됩니다. 클릭 핸들러에서 `SetChecked()`을 호출할 필요가 없습니다.

**수동 제어:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### 액셀러레이터(키보드 단축키)

메뉴 항목에 키보드 단축키를 추가합니다.

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**액셀러레이터 형식:**

- `CmdOrCtrl` - macOS에서는 Cmd, Windows/Linux에서는 Ctrl
- `Shift`, `Alt`, `Option` - 보조 키
- `A-Z`, `0-9` - 문자/숫자 키
- `F1-F12` - 기능 키
- `Enter`, `Space`, `Backspace` 등 - 특수 키

**예:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**플랫폼별 액셀러레이터:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### 도구 설명

메뉴 항목에 마우스를 올렸을 때 표시할 텍스트를 추가합니다(플랫폼별 지원 여부는 다름).

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**플랫폼 지원:**

- **Windows:** ✅ 지원됨
- **macOS:** ❌ 지원되지 않음(도구 설명은 메뉴의 표준 기능이 아님)
- **Linux:** ⚠️ 데스크톱 환경에 따라 다름

### 숨김 상태

메뉴 항목을 제거하지 않고 숨깁니다.

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**용도:** 디버그 옵션, 기능 플래그, 조건부 기능

## 이벤트 처리

### OnClick 핸들러

메뉴 항목을 클릭하면 코드를 실행합니다.

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**컨텍스트에서 제공하는 항목:**

- `ctx.ClickedMenuItem()` - 클릭한 메뉴 항목
- 창 컨텍스트(창 메뉴에서 호출된 경우)
- 애플리케이션 컨텍스트

**예: 핸들러에서 메뉴 항목에 접근**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### 여러 핸들러

여러 핸들러를 설정할 수 있으며 마지막으로 설정한 핸들러가 적용됩니다.

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

**권장 사항:** 핸들러는 한 번만 설정하고, 필요한 경우 그 안에서 조건부 로직을 사용하세요.

## 동적 메뉴

### 메뉴 항목 업데이트

**핵심 원칙:** 메뉴 상태를 변경한 후에는 항상 `menu.Update()`을 호출하세요.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**이것이 중요한 이유:**

- **Windows:** 업데이트할 때 메뉴가 다시 구성됩니다.
- **macOS/Linux:** 중요도는 낮지만 여전히 권장됩니다
- **클릭 핸들러:** Update()를 호출하지 않으면 제대로 실행되지 않습니다

### 메뉴 다시 빌드하기

변경 사항이 크다면 전체 메뉴를 다시 빌드하세요.

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

**다시 빌드해야 하는 경우:**

- 최근 파일 목록이 변경된 경우
- 플러그인 메뉴가 변경된 경우
- 상태가 크게 전환되는 경우

**업데이트해야 하는 경우:**

- 항목 활성화/비활성화
- 레이블 변경
- 체크박스 전환

### 상황별 메뉴

애플리케이션 상태에 따라 메뉴를 조정하세요.

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

## 플랫폼별 차이점

### 메뉴 모음 위치

| 플랫폼 | 위치 | 참고 |
| --- | --- | --- |
| **macOS** | 화면 상단 | 전역 메뉴 모음 |
| **Windows** | 창 상단 | 창별 메뉴 |
| **Linux** | 창 상단 | 창별 메뉴(일반적) |

### 표준 메뉴

**macOS:**

- 앱 이름이 표시되는 "애플리케이션" 메뉴가 있습니다
- 애플리케이션 메뉴에 "환경설정"이 있습니다
- 애플리케이션 메뉴에 "종료"가 있습니다

**Windows/Linux:**

- 애플리케이션 메뉴가 없습니다
- 편집 또는 도구 메뉴에 "환경설정"이 있습니다
- 파일 메뉴에 "종료"가 있습니다

**예: 플랫폼에 적합한 구조**

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

### 단축키 규칙

**macOS:**

- 대부분의 단축키에는 `Cmd+`
- 환경설정에는 `Cmd+,`
- 종료에는 `Cmd+Q`

**Windows:**

- 대부분의 단축키에는 `Ctrl+`
- 환경설정에는 `Ctrl+P` 또는 `Ctrl+,`
- 종료에는 `Alt+F4`(또는 `Ctrl+Q`)

**Linux:**

- 일반적으로 Windows 규칙을 따릅니다
- 데스크톱 환경에서 재정의할 수 있습니다

## 권장 사항

### ✅ 권장

- 메뉴 상태를 변경한 후에는(특히 Windows에서) **menu.Update()를 호출하세요**
- 상호 배타적인 옵션에는 **라디오 그룹을 사용하세요**
- 켜고 끌 수 있는 기능에는 **체크박스를 사용하세요**
- 자주 사용하는 작업에는 **단축키를 추가하세요**
- 구분선을 사용하여 **관련 항목을 그룹화하세요**
- 플랫폼마다 동작이 다르므로 **모든 플랫폼에서 테스트하세요**

### ❌ 금지

- **menu.Update() 호출을 잊지 마세요**. 호출하지 않으면 클릭 핸들러가 제대로 작동하지 않습니다
- **너무 깊게 중첩하지 마세요**. 최대 2-3단계까지 사용하세요
- **구분선으로 시작하거나 끝내지 마세요**. 전문적이지 않아 보입니다
- **macOS에서 툴팁을 사용하지 마세요**. 지원되지 않습니다
- **플랫폼별 단축키를 하드코딩하지 마세요**. `CmdOrCtrl`을 사용하세요

## 문제 해결

### 메뉴 항목이 응답하지 않음

**증상:** 클릭 핸들러가 실행되지 않습니다

**원인:** 항목을 활성화한 후 `menu.Update()` 호출을 잊었습니다

**해결 방법:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### 메뉴 항목이 회색으로 표시됨

**증상:** 메뉴 항목을 클릭할 수 없습니다

**원인:** 항목이 비활성화되어 있습니다

**해결 방법:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### 단축키가 작동하지 않음

**증상:** 키보드 단축키로 메뉴 항목이 실행되지 않습니다

**원인:**

1. 단축키 형식이 잘못되었습니다
2. 시스템 단축키와 충돌합니다
3. 창에 포커스가 없습니다

**해결 방법:**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## 다음 단계

- [애플리케이션 메뉴](/features/menus/application/) - 애플리케이션 메뉴 모음 만들기
- [컨텍스트 메뉴](/features/menus/context/) - 마우스 오른쪽 버튼 클릭 컨텍스트 메뉴
- [시스템 트레이 메뉴](/features/menus/systray/) - 시스템 트레이/메뉴 모음 메뉴
- [메뉴 패턴](/guides/menus/) - 일반적인 메뉴 패턴과 권장 사항

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [메뉴 예제](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/menu)를 확인하세요.
