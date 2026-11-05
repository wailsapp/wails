---
title: "대화 상자 API"
description: "네이티브 대화 상자 API 전체 레퍼런스"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## 개요

Dialogs API는 네이티브 파일 대화 상자와 메시지 대화 상자를 표시하는 메서드를 제공합니다. `app.Dialog` 관리자를 통해 대화 상자에 접근합니다.

**대화 상자 유형:**

- **파일 대화 상자** - 열기 및 저장 대화 상자
- **메시지 대화 상자** - 정보, 오류, 경고 및 질문 대화 상자

모든 대화 상자는 플랫폼의 디자인과 사용 방식을 따르는 <strong>운영 체제 네이티브 대화 상자</strong>입니다.

## 대화 상자 접근

`app.Dialog` 관리자를 통해 대화 상자에 접근합니다:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## 파일 대화 상자

### OpenFile()

파일 열기 대화 상자를 생성합니다.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**예:**

```go
dialog := app.Dialog.OpenFile()
```

### OpenFileDialogStruct 메서드

#### SetTitle()

대화 상자의 제목을 설정합니다.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**예:**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

파일 형식 필터를 추가합니다.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**매개변수:**

- `displayName` - 사용자에게 표시되는 필터 설명(예: "이미지", "문서")
- `pattern` - 세미콜론으로 구분된 확장자 목록(예: "*.png;*.jpg")

**예:**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

초기 디렉터리를 설정합니다.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**예:**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

디렉터리 선택을 활성화하거나 비활성화합니다.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**예(폴더 선택):**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

파일 선택을 활성화하거나 비활성화합니다.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

새 디렉터리 생성을 활성화하거나 비활성화합니다.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

숨김 파일을 표시하거나 숨깁니다.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

대화 상자를 특정 창에 연결합니다.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

대화 상자를 표시하고 선택한 파일을 반환합니다.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**반환값:**

- `string` - 선택한 파일 경로. 빈 문자열은 "선택 항목 없음"으로 처리하십시오(취소할 경우 운영 체제에 따라 플랫폼 구현에서 빈 문자열이나 nil이 아닌 오류를 반환할 수 있습니다).
- `error` - 대화 상자 자체를 표시하지 못한 경우 nil이 아닌 값

**예:**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

대화 상자를 표시하고 선택한 여러 파일을 반환합니다.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**반환값:**

- `[]string` - 선택한 파일 경로의 배열
- `error` - 대화 상자 표시 실패 시 오류

**예:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

파일 저장 대화 상자를 생성합니다.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**예:**

```go
dialog := app.Dialog.SaveFile()
```

### SaveFileDialogStruct 메서드

#### SetTitle()

대화 상자의 제목을 설정합니다.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

기본 파일 이름을 설정합니다.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**예:**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

파일 형식 필터를 추가합니다.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**예:**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

초기 디렉터리를 설정합니다.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

대화 상자를 특정 창에 연결합니다.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

대화 상자를 표시하고 저장 경로를 반환합니다.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**예:**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### 폴더 선택

별도의 `SelectFolderDialog`은 없습니다. 디렉터리 옵션과 함께 `OpenFile()`을 사용하세요.

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## 메시지 대화 상자

모든 메시지 대화 상자는 `*MessageDialog`을 반환하며 동일한 메서드를 공유합니다.

### Info()

정보 대화 상자를 생성합니다.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**예:**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

오류 대화 상자를 생성합니다.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**예:**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

경고 대화 상자를 생성합니다.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**예:**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

사용자 지정 버튼이 있는 질문 대화 상자를 생성합니다.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**예:**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### MessageDialog 메서드

#### SetTitle()

대화 상자의 제목을 설정합니다.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

대화 상자의 메시지를 설정합니다.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

대화 상자의 사용자 지정 아이콘을 설정합니다.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

대화 상자에 버튼을 추가하고 구성할 수 있도록 해당 버튼을 반환합니다.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**반환값:** `*Button` - 추가로 구성할 수 있는 버튼 인스턴스

**예:**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

기본 버튼으로 사용할 버튼을 설정합니다(Enter 키를 누르면 활성화됨).

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**예:**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

취소 버튼으로 사용할 버튼을 설정합니다(Escape 키를 누르면 활성화됨).

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**예:**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

대화 상자를 특정 창에 연결합니다.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

대화 상자를 표시합니다. 버튼 콜백에서 사용자 응답을 처리합니다.

```go
func (d *MessageDialog) Show()
```

**참고:** `Show()`은 값을 반환하지 않습니다. 버튼 콜백을 사용하여 사용자 응답을 처리하세요.

### 버튼 메서드

#### OnClick()

버튼을 클릭했을 때 호출할 콜백 함수를 설정합니다.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

이 버튼을 기본 버튼으로 지정합니다.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

이 버튼을 취소 버튼으로 지정합니다.

```go
func (b *Button) SetAsCancel() *Button
```

## 전체 예제

### 파일 선택 예제

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### 확인 대화 상자 예제

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### 변경 사항 저장 대화 상자

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### 여러 파일 처리

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### 대화 상자를 사용한 오류 처리

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### 플랫폼별 기본값

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## 모범 사례

### 권장 사항

- **네이티브 대화 상자를 사용하세요** - 플랫폼의 디자인과 사용 방식을 따릅니다.
- **명확한 제목을 제공하세요** - 사용자가 목적을 이해하는 데 도움이 됩니다.
- **적절한 필터를 설정하세요** - 사용자가 올바른 파일 형식을 선택하도록 안내합니다.
- **취소를 처리하세요** - 오류가 있는지 확인하세요(사용자가 취소할 수 있음).
- **데이터를 삭제하거나 손상할 수 있는 작업에는 확인 절차를 표시하세요** - Question 대화 상자를 사용하세요.
- **피드백을 제공하세요** - 성공 메시지에는 Info 대화 상자를 사용하세요.
- **합리적인 기본값을 설정하세요** - 기본 디렉터리, 파일 이름 등을 설정하세요.
- **버튼 동작에 콜백을 사용하세요** - 사용자 응답을 올바르게 처리하세요.

### 금지 사항

- **오류를 무시하지 마세요** - 사용자가 취소하면 오류가 반환됩니다.
- **모호한 버튼 레이블을 사용하지 마세요** - "저장"/"취소"처럼 구체적으로 작성하세요.
- **대화 상자를 지나치게 사용하지 마세요** - 작업 흐름을 방해합니다
- **취소 시 오류를 표시하지 마세요** - 취소는 정상적인 동작입니다
- **파일 필터를 빠뜨리지 마세요** - 사용자가 올바른 파일을 찾는 데 도움이 됩니다
- **경로를 하드코딩하지 마세요** - os.UserHomeDir() 또는 이와 유사한 함수를 사용하세요

## 플랫폼별 대화 상자 유형

### macOS

- 대화 상자가 제목 표시줄에서 아래로 펼쳐집니다
- 상위 창에 연결되는 "시트" 스타일입니다
- macOS 네이티브 디자인을 사용합니다

### Windows

- 표준 Windows 대화 상자를 사용합니다
- Windows 디자인 지침을 따릅니다
- 최신 Windows 10/11 디자인을 사용합니다

### Linux

- GTK 기반 시스템에서는 GTK 대화 상자를 사용합니다
- Qt 기반 시스템에서는 Qt 대화 상자를 사용합니다
- 데스크톱 환경과 일치합니다

#### Linux 대화 상자의 동작

Linux에서 기본 GTK4 빌드는 파일 대화 상자에 <strong>xdg-desktop-portal</strong>을 사용합니다. 이를 통해 데스크톱과 네이티브로 통합되지만 일부 옵션은 적용되지 않습니다. 레거시 GTK3 경로 (`-tags gtk3`)에서는 다음 옵션을 프로그래밍 방식으로 완전히 제어할 수 있습니다:

| 옵션 | GTK3(`-tags gtk3`) | GTK4(기본값) | 참고 |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ 작동함 | ❌ 적용되지 않음 | 사용자가 대화 상자의 UI 토글(Ctrl+H 또는 메뉴)로 제어 |
| `CanCreateDirectories()` | ✅ 작동함 | ❌ 적용되지 않음 | 포털에서 항상 활성화됨 |
| `ResolvesAliases()` | ✅ 작동함 | ❌ 적용되지 않음 | 포털에서 심볼릭 링크 해석을 처리함 |
| `SetButtonText()` | ✅ 작동함 | ✅ 작동함 | 사용자 지정 수락 버튼 텍스트가 작동함 |

**이러한 제한이 있는 이유:** GTK4의 포털 기반 대화 상자는 UI 제어를 데스크톱 환경(GNOME, KDE 등)에 위임합니다. 이는 의도된 설계입니다. 포털은 애플리케이션 전반에 걸쳐 일관된 사용자 경험을 제공하고 사용자 설정을 따릅니다.

@note{type="info"}
기본 GTK4 빌드는 포털 기반 대화 상자를 사용합니다. 애플리케이션에서 위의 대화 상자 옵션을 프로그래밍 방식으로 완전히 제어해야 한다면 레거시 `-tags gtk3` 경로로 빌드하세요 (v3.0.x까지 지원되며 v3.1에서 제거됨). 자세한 내용은 [Linux 패키징 - 레거시 GTK3 지원](/guides/build/linux/#legacy-gtk3-support)을 참조하세요.

@end

## 일반적인 패턴

### "다른 이름으로 저장" 패턴

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### "최근 항목 열기" 패턴

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
