---
title: "대화상자 개요"
description: "애플리케이션에 네이티브 시스템 대화상자 표시"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## 네이티브 대화상자

Wails는 모든 플랫폼에서 작동하는 <strong>네이티브 시스템 대화상자</strong>를 제공합니다. 여기에는 메시지 대화상자(정보, 경고, 오류, 질문), 파일 대화상자(열기, 저장, 폴더), 플랫폼 고유의 모양과 동작을 갖춘 사용자 지정 대화상자 창이 포함됩니다.

![취소 및 버리기 버튼이 있는 macOS용 Wails 질문 대화상자](/assets/screenshots/dialog-question-macos.png)

동일한 API가 지원되는 각 플랫폼의 규칙에 맞게 렌더링됩니다. 이 macOS 예제에서는 기본 버튼과 취소 버튼을 포함하여 Wails 창에 연결된 질문 대화상자를 보여 줍니다.

## 빠른 시작

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

**이것으로 끝입니다!** 최소한의 코드로 네이티브 대화상자를 사용할 수 있습니다.

## 대화상자에 접근하기

대화상자는 `app.Dialog` 관리자를 통해 접근합니다.

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## 대화상자 유형

### 정보 대화상자

간단한 메시지를 표시합니다.

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**사용 사례:**

- 성공 메시지
- 정보 알림
- 완료 확인

### 경고 대화상자

경고를 표시합니다.

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**사용 사례:**

- 심각하지 않은 경고
- 지원 중단 예정 알림
- 주의 메시지

### 오류 대화상자

오류를 표시합니다.

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**사용 사례:**

- 오류 메시지
- 실패 알림
- 예외 처리

### 질문 대화상자

사용자에게 질문하고 버튼 콜백을 통해 응답을 처리합니다.

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**사용 사례:**

- 작업 확인
- 예/아니요 질문
- 여러 응답 선택지

## 파일 대화상자

### 파일 열기 대화상자

열 파일을 선택합니다.

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**다중 선택:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### 파일 저장 대화상자

저장할 위치를 선택합니다.

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### 폴더 선택 대화상자

디렉터리 선택을 활성화한 파일 열기 대화상자를 사용하여 디렉터리를 선택합니다.

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## 대화상자 옵션

### 제목 및 메시지

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### 버튼

**간단한 대화상자의 기본 버튼:**

정보, 경고 및 오류 대화상자에는 기본 "확인" 버튼이 표시됩니다.

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**질문 대화상자의 사용자 지정 버튼:**

`AddButton()`을 사용하여 버튼을 추가합니다. 이 메서드는 콜백을 구성할 수 있는 `*Button`을 반환합니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**기본 버튼 및 취소 버튼:**

`SetDefaultButton()`을 사용하여 강조 표시되고 Enter 키로 실행되는 버튼을 지정합니다. `SetCancelButton()`을 사용하여 Escape 키로 실행되는 버튼을 지정합니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### 창에 연결

대화상자를 특정 창에 연결합니다.

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**동작:**

- 대화상자가 부모 창의 중앙에 표시됨
- 대화상자가 표시되는 동안 부모 창이 비활성화됨
- 대화상자가 부모 창과 함께 이동함(macOS)

## 플랫폼별 동작

@tabs{sync-key="platform"}
[macOS]
**macOS 대화상자:**

- 네이티브 NSAlert 모양
- 시스템 테마(라이트/다크) 적용
- 키보드 탐색 지원
- 표준 단축키(취소: ⌘.)
- 접근성 기능 기본 제공
- 창에 연결하면 시트 형태로 표시

**예제:**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Windows 대화상자:**

- 네이티브 TaskDialog 모양
- 시스템 테마 적용
- 키보드 탐색 지원
- 표준 단축키(취소: Esc)
- 접근성 기능 기본 제공
- 부모 창에 모달로 표시

**예:**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Linux 대화 상자:**

- GTK 대화 상자 모양
- 데스크톱 테마를 따름
- 키보드 탐색 지원
- 데스크톱 환경과 통합
- DE(GNOME, KDE 등)에 따라 다름

**예:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## 일반적인 패턴

### 파괴적 작업 전 확인

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### 대화 상자를 사용한 오류 처리

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### 유효성 검사를 포함한 파일 선택

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### 다단계 대화 상자 흐름

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## 권장 사례

### ✅ 권장 사항

- **네이티브 대화 상자를 사용하세요** - 사용자 지정 대화 상자보다 사용자 경험이 좋습니다
- **명확한 메시지를 제공하세요** - 구체적으로 작성하세요
- **적절한 제목을 설정하세요** - 맥락이 중요합니다
- **기본 버튼을 신중하게 사용하세요** - 안전한 옵션을 기본값으로 지정하세요
- **취소를 처리하세요** - 사용자가 취소할 수 있습니다
- **선택한 파일의 유효성을 검사하세요** - 파일 형식을 확인하세요

### ❌ 금지 사항

- **대화 상자를 과도하게 사용하지 마세요** - 작업 흐름을 방해합니다
- **빈번한 메시지에 대화 상자를 사용하지 마세요** - 알림을 사용하세요
- **오류 처리를 빠뜨리지 마세요** - 사용자가 취소할 수 있습니다
- **불필요하게 작업을 차단하지 마세요** - 대안을 고려하세요
- **일반적인 메시지를 사용하지 마세요** - 구체적으로 작성하세요
- **플랫폼 간 차이를 무시하지 마세요** - 모든 플랫폼에서 테스트하세요

## 다음 단계

@cards{cols="2"}
ℹ 메시지 대화 상자
정보, 경고 및 오류 대화 상자입니다.

[자세히 알아보기 →](/features/dialogs/message/)

---
📖 파일 대화 상자
파일 열기, 저장 및 폴더 선택입니다.

[자세히 알아보기 →](/features/dialogs/file/)

---
◆ 사용자 지정 대화 상자
사용자 지정 대화 상자 창을 만듭니다.

[자세히 알아보기 →](/features/dialogs/custom/)

---
▣ 창
창 관리에 관해 알아봅니다.

[자세히 알아보기 →](/features/windows/basics/)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [대화 상자 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs)를 확인하세요.
