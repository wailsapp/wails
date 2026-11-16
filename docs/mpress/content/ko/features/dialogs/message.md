---
title: "메시지 대화상자"
description: "정보, 경고, 오류 및 질문 표시"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## 메시지 대화상자

Wails는 플랫폼에 적합한 모양의 **네이티브 메시지 대화상자**를 제공합니다. 정보, 경고, 오류 및 질문 대화상자의 제목, 메시지, 버튼을 사용자 지정할 수 있습니다. API가 간단하고 네이티브 방식으로 동작하며 기본적으로 접근성을 지원합니다.

## 대화상자 만들기

메시지 대화상자는 `app.Dialog` 관리자를 통해 사용합니다.

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

모든 메서드는 메서드 체이닝으로 구성할 수 있는 `*MessageDialog`을 반환합니다.

## 정보 대화상자

정보 메시지를 표시합니다.

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**사용 사례:**

- 성공 확인
- 완료 알림
- 정보 메시지
- 상태 업데이트

**예시 - 저장 확인:**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## 경고 대화상자

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
- 잠재적 문제

**예시 - 디스크 공간 경고:**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## 오류 대화상자

오류를 표시합니다.

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**사용 사례:**

- 오류 메시지
- 실패 알림
- 예외 처리
- 심각한 문제

**예시 - 네트워크 오류:**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## 질문 대화상자

사용자에게 질문하고 버튼 콜백을 통해 응답을 처리합니다.

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**사용 사례:**

- 작업 확인
- 예/아니요 질문
- 다중 선택
- 사용자 결정

**예시 - 저장하지 않은 변경 사항:**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## 대화상자 옵션

### 제목 및 메시지

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**권장 사항:**

- **제목:** 짧고 설명이 명확하게 작성(2-5단어)
- **메시지:** 명확하고 구체적이며 실행 가능한 내용으로 작성
- **전문 용어 지양:** 쉬운 표현 사용

### 버튼

**단일 버튼(정보/경고/오류):**

정보, 경고 및 오류 대화상자에는 기본 "확인" 버튼이 표시됩니다.

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

사용자 지정 버튼도 추가할 수 있습니다.

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**여러 버튼(질문):**

`AddButton()`을 사용하여 버튼을 추가합니다. 이 메서드는 구성 가능한 `*Button`을 반환합니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

dialog.Show()
```

**기본 버튼 및 취소 버튼:**

`SetDefaultButton()`을 사용하여 강조 표시되고 Enter 키로 실행될 버튼을 지정합니다. `SetCancelButton()`을 사용하여 Escape 키로 실행될 버튼을 지정합니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

버튼에서 플루언트 방식의 `SetAsDefault()` 및 `SetAsCancel()` 메서드를 사용할 수도 있습니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**권장 사항:**

- **1-3개 버튼:** 너무 많은 선택지를 제시하지 마세요.
- **명확한 레이블:** "확인" 대신 "저장" 사용
- **안전한 기본값:** 비파괴적 작업
- **중요한 버튼 순서:** 가장 가능성이 높은 작업을 먼저 배치(취소 제외)

### 사용자 지정 아이콘

대화상자에 사용자 지정 아이콘을 설정합니다.

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### 창에 연결

특정 창에 연결합니다.

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**장점:**

- 올바른 창에 대화상자가 표시됨
- 대화상자가 표시되는 동안 부모 창이 비활성화됨
- 향상된 다중 창 사용자 경험

## 전체 예제

### 파괴적 작업 확인

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### 종료 확인

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### 다운로드 옵션이 있는 업데이트 대화상자

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### 사용자 지정 아이콘이 있는 질문

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## 권장 사항

### ✅ 권장

- **구체적으로 작성하세요** - "성공"이 아니라 "파일이 문서 폴더에 저장되었습니다"라고 표시하세요.
- **적절한 유형을 사용하세요** - 오류에는 오류 유형을, 경고에는 경고 유형을 사용하세요.
- **맥락을 제공하세요** - 관련 세부 정보를 포함하세요.
- **명확한 버튼 레이블을 사용하세요** - "확인"이 아니라 "삭제"를 사용하세요.
- **안전한 기본값을 설정하세요** - 비파괴적 작업을 사용하세요.
- **취소를 처리하세요** - 사용자가 대화상자를 닫을 수 있습니다.

### ❌ 비권장

- **과도하게 사용하지 마세요** - 작업 흐름을 방해합니다.
- **빈번한 업데이트에 사용하지 마세요** - 대신 알림을 사용하세요.
- **일반적인 메시지를 사용하지 마세요** - "오류"만으로는 아무 정보도 전달할 수 없습니다.
- **오류를 무시하지 마세요** - dialog.Show() 오류를 처리하세요.
- **불필요하게 차단하지 마세요** - 비동기 대안을 고려하세요.
- **기술 전문 용어를 사용하지 마세요** - 쉬운 표현을 사용하세요.

## 플랫폼별 차이점

### macOS

- 창에 연결하면 시트 스타일로 표시됨
- 표준 키보드 단축키 지원(취소는 ⌘.)
- 시스템 테마를 자동으로 따름
- 접근성 기능 기본 제공

### Windows

- 모달 대화 상자
- TaskDialog 모양
- Esc 키로 취소
- 시스템 테마를 따름

### Linux

- GTK 대화 상자
- 데스크톱 환경에 따라 다름
- 데스크톱 테마를 따름
- 표준 키보드 탐색

## 다음 단계

@cards{cols="2"}
📖 파일 대화 상자
파일 열기와 저장 및 폴더 선택.

[자세히 알아보기 →](/features/dialogs/file/)

---
◆ 사용자 지정 대화 상자
사용자 지정 대화 상자 창을 만듭니다.

[자세히 알아보기 →](/features/dialogs/custom/)

---
● 알림
사용자 작업을 방해하지 않는 알림입니다.

[자세히 알아보기 →](/features/notifications/overview/)

---
★ 이벤트
이벤트를 사용하여 비차단 방식으로 통신합니다.

[자세히 알아보기 →](/features/events/system/)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [대화 상자 예제](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/dialogs)를 확인하세요.
