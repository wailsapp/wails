---
title: "macOS 플랫폼 통합"
description: "macOS의 문서 창, Dock 및 메뉴 추가 기능, 네이티브 패널, 상태 항목, 피드백, 다양한 형식을 지원하는 클립보드, 외부로 드래그하기, 권한, 전원 및 수명 주기"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

관련 플랫폼: macOS

Wails v3는 macOS 사용자가 네이티브 앱에 기대하는 동작을 애플리케이션에 제공합니다. 프록시 아이콘과 계단식 배치를 갖춘 문서 창, 기호와 배지가 있는 메뉴, 진행률을 표시하는 Dock 메뉴, 네이티브 경고와 패널, 제거 가능한 상태 항목, 햅틱과 음성, 다양한 형식을 지원하는 클립보드와 외부로 드래그하기, 권한·전원·로캘에 관한 시스템 정보, Services 메뉴와 Handoff, AppleScript, Quick Look 통합이 포함됩니다. 모든 기능은 `application` 패키지를 통해 Go에서 제어합니다.

동일한 코드는 Windows와 Linux에서도 컴파일됩니다. 설정 메서드는 값을 저장하고, 조회 메서드는 0 값을 반환하며, macOS가 필요한 작업은 `ErrMacOnly`, `ErrDialogNotSupported`, `ErrClipboardNotSupported`와 같은 문서화된 오류를 반환합니다. 아래 [플랫폼 참고 사항](#platform-notes)에 macOS 이외 플랫폼에서 각 영역의 동작이 나와 있습니다.

네이티브 창 구성 요소(도구 막대, 사이드바, 속성 패널, 보조 컨트롤 및 창 탭)는 [네이티브 macOS 창 구성 요소](/guides/macos-native-chrome) 가이드를 참조하세요.

## 문서 창

문서 창은 제목 표시줄에 해당 파일을 표시하고, 저장하지 않은 변경 사항이 있으면 닫기 버튼에 점을 표시하며, 새 창을 계단식으로 엽니다. 모두 `WebviewWindow`의 메서드와 `MacWindow`의 옵션으로 제공됩니다.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Report.md",
    URL:             "/editor",
    InitialPosition: application.WindowCascade,
    Mac: application.MacWindow{
        FrameAutosaveName: "editor.main",
        TitleBar: application.MacTitleBar{
            WindowButtonsOffset: &application.Point{X: 12, Y: 8},
        },
    },
})

window.SetRepresentedFile("/Users/me/Documents/Report.md")
window.SetSubtitle("Documents")
window.SetDocumentEdited(true)
```

- `SetRepresentedFile`은 제목 표시줄에 파일의 프록시 아이콘을 표시합니다. 사용자는 아이콘을 다른 앱으로 드래그하거나 command-클릭하여 경로를 볼 수 있습니다. 제거하려면 `""`을 전달하세요. `RepresentedFile`로 설정값을 읽습니다.
- `SetDocumentEdited`는 닫기 버튼에 저장하지 않은 변경 사항을 나타내는 점을 표시하고 프록시 아이콘을 흐리게 합니다. `IsDocumentEdited`로 상태를 읽습니다.
- `SetSubtitle`은 macOS 11 이상에서 제목 아래에 두 번째 줄을 표시합니다.
- `InitialPosition: application.WindowCascade`는 새 문서를 열 때처럼 마지막 계단식 창의 아래쪽과 오른쪽에 창을 배치합니다. `CascadeFrom(other)`는 기존 창을 기준으로 동일하게 배치하고 이후 창을 위한 계단식 배치 지점을 업데이트합니다.
- `Mac.FrameAutosaveName`은 창이 처음 표시되기 전에 저장된 위치와 크기를 복원하고 창이 이동할 때마다 계속 저장합니다. 복원된 창 프레임은 `X`, `Y`, `Width`, `Height`, `InitialPosition`보다 우선합니다. `SetFrameAutosaveName`은 실행 중인 창의 이름을 변경합니다.
- `MacTitleBar.WindowButtonsOffset`은 닫기, 최소화, 확대 버튼을 지정한 포인트 수만큼 이동합니다. `SetWindowButtonsOffset`과 `ResetWindowButtonsOffset`으로 실행 중에 변경할 수 있습니다.

세 설정 메서드는 모두 네이티브 창이 생성되기 전에 호출할 수 있으며, 값은 창이 생성될 때 적용됩니다.

### 사용자 주의 요청

`RequestAttention`은 애플리케이션이 백그라운드에 있는 동안 Dock 아이콘을 튀어 오르게 합니다. 정보성 요청은 한 번 튀어 오릅니다. 긴급 요청은 사용자가 애플리케이션을 활성화하거나 요청을 취소할 때까지 계속 튀어 오릅니다.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

한 번 주의를 요청하는 크로스 플랫폼 방법으로 `Flash`가 계속 제공됩니다.

### 인쇄 및 내보내기

`PrintWithOptions`는 명시적인 페이지 설정으로 WebView를 인쇄합니다. 0 값은 공유 인쇄 설정이 적용된 인쇄 패널을 표시합니다. `Print`는 기존 동작(가로 방향, 30포인트 여백)을 유지합니다.

`ExportPDF`는 페이지를 PDF 문서로 렌더링하고 `Snapshot`은 PNG로 캡처합니다. 둘 다 WebKit을 기다리므로 goroutine에서 호출하고 애플리케이션 스레드에서는 절대 호출하지 마세요. 해당 스레드에서 호출하면 `ErrMacExportOnMainThread`를 반환합니다.

```go
err := window.PrintWithOptions(application.PrintOptions{
    Orientation: application.PrintOrientationPortrait,
    Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
    Silent:      false,
})
if err != nil {
    log.Println("print:", err)
}

go func() {
    pdf, err := window.ExportPDF(application.PDFExportOptions{})
    if err != nil {
        log.Println("export:", err)
        return
    }
    if err := os.WriteFile("report.pdf", pdf, 0o644); err != nil {
        log.Println("write:", err)
    }

    png, err := window.Snapshot(application.SnapshotOptions{Width: 800})
    if err != nil {
        log.Println("snapshot:", err)
        return
    }
    if err := os.WriteFile("preview.png", png, 0o644); err != nil {
        log.Println("write:", err)
    }
}()
```

`PrintOptions`는 `PrinterName`, `PaperName`(`"iso-a4"` 같은 PostScript 이름), `Scale`도 받습니다. `PDFExportOptions`와 `SnapshotOptions`에는 캡처 영역을 제한하는 선택적 `Rect`와 기본값이 `DefaultMacExportTimeout`(30초)인 `Timeout`을 지정할 수 있습니다.

## 시트

시트는 저장 패널처럼 부모 창 위쪽에 연결되는 두 번째 창입니다. 모든 `WebviewWindow`는 `PresentSheet`를 사용해 다른 창의 시트로 표시할 수 있으며, `EndSheet`와 응답 코드로 종료할 수 있습니다. 응답 코드는 시트의 `OnSheetEnd` 콜백에 전달됩니다. 연결하기 전에 화면에 잠깐 나타나지 않도록 시트 창을 `Hidden`으로 생성하세요. 종료 시 AppKit이 다시 화면에서 치우므로 같은 창을 반복해서 표시할 수 있습니다.

```go
sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Rename",
    URL:    "/rename",
    Width:  420,
    Height: 180,
    Hidden: true,
})
sheet.OnSheetEnd(func(code int) {
    if code == application.MacSheetResponseOK {
        log.Println("renamed")
    }
})

// From a menu item or a bound method:
if err := window.PresentSheet(sheet); err != nil {
    log.Println(err)
}

// From the sheet's own page, through a bound method:
sheet.EndSheet(application.MacSheetResponseOK)
```

`PresentCriticalSheet`는 이미 연결된 일반 시트 뒤에 대기시키지 않고 그 앞에 시트를 표시합니다. `PresentNativeSheet`는 `NativeWindow`에 대해 동일하게 작동합니다. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet`, `HasAttachedSheet`는 현재 상태를 알려줍니다. 시트를 종료하는 대신 시트 창을 닫으면 `MacSheetResponseStop`이 전달됩니다.

## 팝오버

`MacPopover`는 `NSPopover`입니다. 창 안의 사각형, 도구 막대 항목 또는 메뉴 막대의 상태 항목을 기준으로 표시되는 임시 패널입니다. 콘텐츠는 네이티브 `MacAccessory` 컨트롤 띠이며, [네이티브 macOS 창 구성 요소](/guides/macos-native-chrome) 가이드에서 제목 표시줄 보조 컨트롤에 사용하는 것과 같은 유형입니다. 처음 표시하기 전에 모든 컨트롤을 추가하세요.

```go
strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
strip.AddLabel("Sort by")
strip.AddSegmented([]string{"Date", "Title"}, 0)
strip.AddFlexibleSpace()
strip.AddButton("Apply").OnClick(func(*application.Context) {
    // apply the sort
})

popover := application.NewMacPopover(application.MacPopoverOptions{
    Width:    320,
    Behavior: application.MacPopoverBehaviorTransient,
    Content:  strip,
})
popover.OnClose(func() {
    log.Println("popover closed")
})

// Anchored to a rectangle in the page, in window content coordinates:
err := popover.ShowRelativeTo(application.Rect{X: 20, Y: 60, Width: 120, Height: 28}, window, application.MacRectEdgeMaxY)
if err != nil {
    log.Println(err)
}
```

`MacToolbarItem.ShowPopover`와 `SystemTray.ShowPopover`는 같은 팝오버를 도구 막대 항목이나 상태 항목에 고정합니다. `MacPopoverBehaviorTransient`는 팝오버 바깥을 클릭하면 닫고, `Semitransient`는 팝오버를 표시한 창을 클릭한 경우에만 닫습니다. 기본 동작은 `Close`를 호출할 때까지 열린 상태를 유지합니다. `MacRectEdge`는 팝오버가 나타날 쪽을 선택합니다. `SetContentSize`와 `SetBehavior`는 표시 중인 팝오버를 조정하며, `Destroy`는 네이티브 팝오버를 해제하고 콘텐츠 띠를 다른 곳에서 사용할 수 있도록 합니다.

## 상태 복원

macOS는 충돌, 강제 종료 또는 재부팅 후 애플리케이션의 창을 다시 엽니다. 시스템 설정에서 "애플리케이션 종료 시 창 닫기"가 꺼져 있으면 정상 종료 후에도 다시 엽니다. 창에 `Mac.RestorationID`를 지정하고, `SetRestorationData`로 다시 만드는 데 필요한 정보를 저장한 다음, 다음 실행 시 창을 다시 만들도록 `app.Window.OnRestore`를 등록하세요.

```go
app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
    if id != "editor" {
        return nil
    }
    path := state.Get("path")
    restored := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: filepath.Base(path),
        URL:   "/editor?path=" + path,
        Mac:   application.MacWindow{RestorationID: id},
    })
    restored.SetRestorationData(state.Data)
    return restored
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Report.md",
    URL:   "/editor?path=/Users/me/Documents/Report.md",
    Mac:   application.MacWindow{RestorationID: "editor"},
})
window.SetRestorationData(map[string]string{"path": "/Users/me/Documents/Report.md"})
```

애플리케이션이 종료될 때 표시 중인 창만 저장됩니다. `RestorationState.Data`는 문자열 맵이므로 식별자, 경로, 위치만 저장하세요. `InteractionState`는 WebView의 뒤로·앞으로 목록과 스크롤 위치를 불투명한 데이터 블롭(macOS 12 이상)으로 반환합니다. `RestoreInteractionState`는 이를 다시 만든 창에 적용합니다. 보통 복원 데이터에 base64로 인코딩해 저장합니다. `SetRestorationID`와 `RestorationID`는 실행 중인 창의 식별자를 변경하고 읽습니다.

## 표시 옵션

`MacPresentationOptions`는 `NSApplication.presentationOptions`에 대응합니다. 애플리케이션이 활성화된 동안 Dock 또는 메뉴 막대를 숨기고 프로세스 전환, 강제 종료, 로그아웃 또는 숨기기 명령을 비활성화하는 비트마스크입니다. 실행 시 적용하려면 애플리케이션 옵션에서 `Mac.PresentationOptions`를 설정하고, 실행 중 변경하려면 `SetPresentationOptions`를 사용하세요. 잘못된 조합은 AppKit에 전달되기 전에 `ErrMacPresentationOptionsInvalid`를 감싼 오류로 거부됩니다.

```go
kiosk := application.MacPresentationHideDock |
    application.MacPresentationHideMenuBar |
    application.MacPresentationDisableProcessSwitching |
    application.MacPresentationDisableForceQuit

if err := app.SetPresentationOptions(kiosk); err != nil {
    log.Println(err)
}
log.Println("presentation:", app.PresentationOptions())

// Restore the standard Dock and menu bar:
_ = app.SetPresentationOptions(application.MacPresentationDefault)
```

메뉴 막대를 숨기려면(`HideMenuBar` 또는 `AutoHideMenuBar`) Dock 옵션 중 하나가 필요하며, `AutoHideToolbar`에는 `FullScreen`과 `AutoHideMenuBar`가 모두 필요합니다. `Validate`는 값이 위반한 첫 번째 규칙을 보고하고 `Has`는 개별 플래그를 검사합니다.

## 메뉴와 Dock

메뉴 항목에 SF Symbols, 배지, 섹션 머리글, 색상 팔레트, 혼합 체크 상태, 대체 항목 및 들여쓰기를 사용할 수 있습니다. 모두 `MenuItem`과 `Menu`의 메서드이므로 애플리케이션 메뉴, 컨텍스트 메뉴, 트레이 메뉴 및 Dock 메뉴에서 작동합니다.

```go
menu := app.NewMenu()
menu.AddRole(application.AppMenu)

fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Open...").SetSymbol("folder").SetAccelerator("CmdOrCtrl+o").
    OnClick(func(*application.Context) {
        // open the file, then note it in the recent list:
        app.Menu.AddRecentDocument("/Users/me/Documents/Report.md")
    })
fileMenu.AddRole(application.OpenRecent)

view := menu.AddSubmenu("View")
view.AddSectionHeader("Mailboxes")
inbox := view.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
inbox.OnClick(func(*application.Context) {
    inbox.ClearBadge()
})
view.Add("Updates").SetBadgeText("New")

view.AddSeparator()
wrap := view.AddCheckbox("Wrap lines", false).SetMixed()
view.Add("Reset").SetIndentationLevel(1).OnClick(func(*application.Context) {
    wrap.SetMixed()
})

view.AddSeparator()
view.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").OnClick(func(*application.Context) {})
view.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).
    OnClick(func(*application.Context) {})

colours := []application.RGBA{
    application.NewRGB(255, 59, 48),
    application.NewRGB(255, 149, 0),
    application.NewRGB(52, 199, 89),
}
view.AddPalette([]string{"tag.fill"}, colours, 0, func(_ *application.Context, index int) {
    log.Println("tag colour", index)
}).SetLabel("Tag colour")

app.Menu.Set(menu)
```

- `SetSymbol`은 제목 옆에 SF Symbol을 표시합니다(macOS 11 이상). `SetBitmap`으로 설정한 이미지를 대체합니다.
- `SetBadge`는 개수를, `SetBadgeText`는 제목 뒤에 짧은 문자열을 표시합니다(macOS 14 이상). `ClearBadge`는 배지를 제거하고 `BadgeCount`와 `BadgeText`는 설정값을 읽습니다.
- `AddSectionHeader`는 상호작용할 수 없는 머리글을 추가합니다(macOS 14 이상). 이전 릴리스에서는 제목이 같은 비활성화된 항목으로 표시됩니다.
- `AddPalette`는 `NSMenu`의 팔레트 메뉴를 이용하는 색상 견본 행을 추가합니다(macOS 14 이상). 견본마다 기호 하나, 색상마다 하나씩 전달하거나 채워진 원을 사용하려면 빈 슬라이스를 전달하세요. 레이블이 없으면 팔레트가 부모 메뉴에 인라인으로 표시됩니다. `SetLabel`을 사용하면 제목이 있는 하위 메뉴로 표시됩니다. `PaletteSelected`는 선택된 인덱스를 반환합니다.
- `SetMixed`는 확인란을 대시로 그려지는 혼합 상태로 만듭니다. 클릭하면 AppKit의 동작과 같이 완전히 켜진 상태가 됩니다.
- `SetAlternate(true)`는 서로 다른 수정 키를 누르고 있는 동안 바로 위 항목을 대신해 해당 항목을 표시합니다. 두 항목은 같은 키를 공유하고 수정 키는 달라야 합니다.
- `SetIndentationLevel`은 제목을 최대 15단계까지 들여씁니다.

### 최근 항목 열기

`fileMenu.AddRole(application.OpenRecent)`은 표준 최근 항목 열기 하위 메뉴를 추가합니다. macOS에서는 열릴 때마다 `NSDocumentController`가 메뉴를 채우며 메뉴 지우기 항목도 포함합니다. `app.Menu.AddRecentDocument`로 파일을 추가하고 `RecentDocuments`로 목록을 조회하며 `ClearRecentDocuments`로 목록을 비우세요. 목록은 재시작 후에도 유지됩니다.

최근 파일을 선택하면 Finder에서 파일을 열었을 때와 동일한 이벤트가 전달됩니다:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Dock 메뉴

`app.Menu.SetDockMenu`는 Dock 아이콘을 마우스 오른쪽 버튼으로 클릭했을 때 표시할 정적 메뉴를 설치합니다. `OnDockMenu`는 표시 직전에 요청에 따라 메뉴를 매번 구성하므로 항목이 변화하는 상태를 반영해야 할 때 적합합니다. 구성 함수가 정적 메뉴보다 우선합니다.

```go
dockMenu := application.NewMenu()
dockMenu.Add("New Document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
    // create a document
})
app.Menu.SetDockMenu(dockMenu)

app.Menu.OnDockMenu(func() *application.Menu {
    dynamic := application.NewMenu()
    dynamic.AddSectionHeader("Recent")
    for _, path := range app.Menu.RecentDocuments() {
        path := path
        dynamic.Add(filepath.Base(path)).OnClick(func(*application.Context) {
            app.Menu.OpenRecentDocument(path)
        })
    }
    return dynamic
})
```

### Dock 진행률

Dock 서비스는 기존 배지 지원과 함께 Dock 아이콘 위에 진행률 표시줄을 그립니다. `dock.New()`를 서비스로 등록하고 0에서 1 사이의 비율을 `SetProgress`에 전달하세요.

```go
import "github.com/wailsapp/wails/v3/pkg/services/dock"

dockService := dock.New()

app := application.New(application.Options{
    Name: "Exporter",
    Services: []application.Service{
        application.NewService(dockService),
    },
})

app.Event.On("export:progress", func(event *application.CustomEvent) {
    if fraction, ok := event.Data.(float64); ok {
        _ = dockService.SetProgress(fraction)
    }
})
app.Event.On("export:done", func(*application.CustomEvent) {
    _ = dockService.ClearProgress()
})
```

`GetProgress`는 현재 비율을 반환하며 표시줄이 보이지 않으면 `nil`을 반환합니다.

## 대화상자

메시지, 열기 및 저장 대화상자는 macOS 옵션을 받으며, 대화상자 관리자에는 텍스트 입력 프롬프트와 시스템 색상 및 글꼴 패널이 추가됩니다.

### 경고

`SetSuppression`은 "이 메시지를 다시 표시하지 않음" 확인란을 추가하고 `SetHelp`는 도움말 버튼을 표시합니다. 버튼 콜백에서 `Suppressed`로 확인란 상태를 읽거나 `OnSuppression`을 등록하여 먼저 전달받으세요.

```go
dialog := app.Dialog.Warning().
    SetTitle("Delete 3 items?").
    SetMessage("The items will be moved to the Bin.").
    SetSuppression("Do not warn me again").
    SetHelp(func() {
        log.Println("help requested")
    }).
    AttachToWindow(window)

dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
    if dialog.Suppressed() {
        // remember not to ask again
    }
})
dialog.AddButton("Cancel").SetAsCancel()
dialog.Show()
```

### 텍스트 입력 프롬프트

`Prompt`는 텍스트 필드가 있는 경고를 표시하고 닫힐 때까지 호출을 차단하므로 goroutine 또는 바인딩된 메서드에서 호출하세요. `Secure`는 필드를 비밀번호 필드로 바꾸고 `Window`는 경고를 시트로 표시합니다.

```go
go func() {
    value, ok, err := app.Dialog.Prompt(application.PromptOptions{
        Title:        "Name this document",
        Message:      "The name is used for the exported file.",
        Placeholder:  "Untitled",
        DefaultValue: "Quarterly report",
        OKLabel:      "Rename",
        Window:       window,
    })
    if err != nil || !ok {
        return
    }
    log.Println("renamed to", value)
}()
```

### 파일 패널

`AddContentType`은 열기와 저장 대화상자 모두에서 통일 형식 식별자로 필터링합니다. `AddFilter`와 함께 사용할 수 있으므로 `"public.image"`는 시스템이 아는 모든 이미지 유형과 일치하게 하고, 필터로는 확장자를 기준으로 PDF를 찾을 수 있습니다.

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Choose images").
    AddFilter("PDF", "*.pdf").
    AddContentType("public.image").
    AttachToWindow(window).
    PromptForMultipleSelection()
if err == nil {
    log.Println(paths)
}
```

저장 패널에는 형식 팝업, 이름 필드의 사용자 지정 레이블 및 Finder 태그가 추가됩니다. 사용자가 팝업의 선택을 변경하면 `SetFormats`가 허용되는 유형과 이름 필드의 확장자를 바꾸고, `SelectedFormat`은 최종 선택을 보고합니다.

```go
formats := []application.DialogFormat{
    {Label: "PNG image", Extension: "png", UTI: "public.png"},
    {Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
}
save := app.Dialog.SaveFile().
    SetFilename("Quarterly report.png").
    SetNameFieldLabel("Export As:").
    SetTags([]string{"Reports", "Draft"}).
    AttachToWindow(window)
save.SetFormats(formats, 0, nil)

path, err := save.PromptForSingleSelection()
if err == nil && path != "" {
    log.Println("export", path, "as", formats[save.SelectedFormat()].Label)
}
```

### 색상 및 글꼴 패널

`PickColor`와 `PickFont`는 공유 시스템 패널을 열고 패널이 닫힐 때까지 호출을 차단합니다. `OnChange`는 패널이 열려 있는 동안 선택이 바뀔 때마다 전달하므로 페이지에서 선택 결과를 실시간으로 미리 볼 수 있습니다. 종류별로 한 번에 패널 하나만 열 수 있으며 두 번째 호출은 `ErrDialogInProgress`를 반환합니다.

```go
go func() {
    colour, changed, err := app.Dialog.PickColor(application.ColorPickerOptions{
        Initial:    application.NewRGB(52, 120, 246),
        ShowsAlpha: true,
        Title:      "Accent colour",
        OnChange: func(colour application.RGBA) {
            app.Event.Emit("accent:preview", colour)
        },
    })
    if err == nil && changed {
        log.Println("accent", colour)
    }

    font, changed, err := app.Dialog.PickFont(application.FontPickerOptions{
        Family: "Helvetica Neue",
        Size:   18,
    })
    if err == nil && changed {
        log.Println(font.Family, font.Face, font.PostScriptName, font.Size)
    }
}()
```

## 상태 항목과 피드백

### 상태 항목

macOS의 시스템 트레이 항목은 `NSStatusItem`입니다. SF Symbol로 그릴 수 있고, 도구 설명을 표시할 수 있으며, 기본 제공 항목처럼 사용자가 제거할 수 있습니다.

```go
tray := app.SystemTray.New()
tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
tray.SetTooltip("Recorder. Command-drag to remove.")
tray.SetRemovable(true, "com.example.recorder.tray")
tray.OnVisibilityChange(func(visible bool) {
    log.Println("status item visible:", visible)
})

trayMenu := app.Menu.New()
trayMenu.Add("Show").OnClick(func(*application.Context) {
    window.Show().Focus()
})
tray.SetMenu(trayMenu)
tray.Run()
```

- `SetSymbol`은 메뉴 막대의 모양을 따르도록 기호를 템플릿 이미지로 렌더링합니다(macOS 11 이상). `SetSymbolConfiguration`은 포인트 크기와 굵기를 설정합니다.
- `SetTooltip`은 마우스를 올렸을 때 표시할 텍스트를 설정합니다. `Tooltip`은 이를 읽습니다.
- `SetRemovable(true, name)`을 사용하면 사용자가 command-드래그로 메뉴 막대에서 항목을 제거할 수 있습니다. macOS가 실행 간 제거 상태를 기억하도록 안정적인 자동 저장 이름을 지정하세요. `Show` 또는 `SetVisible(true)`로 다시 표시할 수 있습니다.
- `IsVisible`은 `NSStatusItem.visible`을 읽으므로 사용자가 항목을 제거한 뒤에는 false입니다. `OnVisibilityChange`는 모든 변경 사항을 보고합니다.

### 햅틱

`app.Haptics.Perform`은 애플리케이션이 활성화된 동안 Force Touch 트랙패드 또는 Magic Trackpad에서 패턴을 재생합니다.

```go
app.Haptics.Perform(application.HapticAlignment)
```

종류는 `HapticGeneric`, `HapticAlignment`(항목이 제자리에 맞춰질 때), `HapticLevelChange`(걸림 또는 클릭 단계)입니다. `IsSupported`는 플랫폼에서 햅틱 피드백을 재생할 수 있는지 보고합니다.

### 소리

`app.Sound`는 경고음, 이름이 지정된 시스템 소리 또는 오디오 파일을 재생합니다.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play`는 `SystemSounds`에 있는 이름 또는 Core Audio가 디코딩할 수 있는 파일의 절대 경로를 받습니다. `PlayData`는 메모리에 있는 전체 오디오 파일을 재생합니다.

### 음성

`app.Speech.Speak`는 시스템 음성으로 읽을 텍스트를 대기열에 넣고 `Utterance`를 반환합니다. 발화는 차례대로 재생됩니다. `Stop`은 발화 하나를 취소하고 `StopAll`은 대기열을 비웁니다. `Voices`는 설치된 음성을 식별자와 언어와 함께 나열합니다.

```go
utterance, err := app.Speech.Speak("Export finished", application.SpeechOptions{
    Voice: "com.apple.voice.compact.en-GB.Daniel",
    Rate:  0.5,
})
if err != nil {
    log.Println(err)
    return
}
utterance.OnFinished(func() {
    log.Println("done, stopped:", utterance.WasStopped())
})
```

`Recognize`는 `SFSpeechRecognizer`를 사용해 기본 마이크의 음성을 텍스트로 변환합니다. 첫 호출에서는 마이크와 음성 인식 권한을 요청하고 사용자가 응답할 때까지 차단하므로 goroutine에서 호출하세요. 중간 변환 결과는 `OnPartial`로 전달됩니다. `Stop`은 캡처를 끝내고 최종 텍스트를 반환합니다.

```go
go func() {
    session, err := app.Speech.Recognize(application.RecognitionOptions{
        Locale: "en-US",
        OnPartial: func(text string) {
            app.Event.Emit("dictation:partial", text)
        },
    })
    if err != nil {
        if errors.Is(err, application.ErrSpeechRecognitionDenied) {
            _ = app.Permissions.OpenSystemSettings(application.PermissionKindMicrophone)
        }
        log.Println(err)
        return
    }
    time.Sleep(5 * time.Second)
    text, err := session.Stop()
    if err != nil {
        log.Println(err)
        return
    }
    app.Event.Emit("dictation:final", text)
}()
```

음성 인식에는 `Info.plist`에 `NSSpeechRecognitionUsageDescription`과 `NSMicrophoneUsageDescription`을 선언한 번들 애플리케이션이 필요합니다. 두 키가 없으면 macOS가 접근을 거부하고 `Recognize`가 `ErrSpeechRecognitionUsageDescription`을 반환합니다.

## 클립보드와 드래그

### 다양한 형식을 지원하는 클립보드

`app.Clipboard`는 일반 텍스트 외에도 이미지, 파일 참조, HTML, RTF 및 임의의 통일 형식 식별자로 지정된 원시 데이터를 읽고 씁니다. `Types`는 클립보드에 있는 유형을 나열하고 `OnChange`는 어느 애플리케이션에서든 발생한 변경을 보고합니다.

```go
if err := app.Clipboard.SetHTML("<p>Rich <b>HTML</b></p>", "Rich HTML"); err != nil {
    log.Println(err)
}
_ = app.Clipboard.SetFiles([]string{"/Users/me/Documents/Report.md"})
_ = app.Clipboard.SetData("com.example.record", []byte(`{"id":42}`))

stop := app.Clipboard.OnChange(func() {
    log.Println("clipboard changed", app.Clipboard.ChangeCount(), app.Clipboard.Types())
    if files, err := app.Clipboard.Files(); err == nil && len(files) > 0 {
        log.Println("files:", files)
    }
})
defer stop()
```

`SetImage`와 `Image`는 PNG 바이트를 사용하며, 다른 앱이 TIFF로 복사한 이미지는 자동으로 변환됩니다. 클립보드 변경에 대한 시스템 알림이 없으므로 `OnChange`는 리스너가 하나 이상 있는 동안 500 ms마다 변경 횟수를 폴링합니다.

### 외부로 드래그하기

`StartDrag`는 사용자가 Finder에서 항목을 집어 든 것처럼 창에서 시스템 드래그를 시작합니다. 기존 파일, 대상이 드롭을 받아들일 때에만 콘텐츠를 생성하는 파일 프로미스 또는 일반 텍스트를 제공합니다. 마우스 동작 중에 시작하세요. Go 메서드를 바인딩하고 페이지의 드래그 가능한 요소에서 `mousedown` 또는 `pointerdown` 처리기가 이를 호출하도록 하세요. WebKit이 자체 드래그를 시작하지 않도록 HTML `draggable` 속성은 `false`로 설정하세요.

```go
// DragService is bound to the page and called from a mousedown handler.
type DragService struct {
    window *application.WebviewWindow
}

func (s *DragService) DragExport() error {
    return s.window.StartDrag(application.DragItems{
        Promises: []application.DragPromise{{
            Filename: "export.csv",
            Data: func() ([]byte, error) {
                return []byte("id,name\n1,Wails\n"), nil
            },
        }},
        Operations: application.DragOperationCopy,
    })
}
```

```go
window.OnDragEnd(func(operation application.DragOperation) {
    log.Println("drag ended:", operation)
})
```

동작 중이 아닐 때 `StartDrag`를 호출하면 `ErrDragOutNoGesture`를 반환합니다. `DragItems.Image`와 `ImageOffset`은 커서 아래에 표시할 이미지를 설정합니다.

### 다른 애플리케이션에서 드롭하기

파일 드롭은 계속 `WindowFilesDropped` 이벤트를 사용합니다. 다른 앱에서 드래그한 텍스트, URL 또는 이미지를 받으려면 `DropTypes`에 유형을 나열하고 `OnDrop`을 등록하세요. 이런 드롭은 페이지 자체의 HTML5 드롭 처리기가 아닌 Go로 전달됩니다.

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Inbox",
    URL:   "/",
    DropTypes: []application.DropType{
        application.DropFiles,
        application.DropText,
        application.DropURLs,
        application.DropImages,
    },
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    log.Println("files:", event.Context().DroppedFiles())
})

window.OnDrop(func(_ *application.Context, data application.DropData) {
    log.Println("text:", data.Text, "urls:", data.URLs, "images:", len(data.Images), "at", data.X, data.Y)
})
```

## 시스템

### 권한

`app.Permissions`는 시스템 개인정보 보호 권한(카메라, 마이크, 화면 녹화, 손쉬운 사용, 위치, 알림, 입력 모니터링 및 전체 디스크 접근)을 보고하고 요청합니다. `Status`는 권한을 요청하는 메시지를 표시하지 않습니다. `Request`는 아직 결정되지 않은 권한을 요청하고 사용자가 응답할 때까지 차단하므로 goroutine에서 호출하세요. `OpenSystemSettings`는 해당 개인정보 보호 및 보안 설정 패널을 엽니다.

```go
go func() {
    status := app.Permissions.Status(application.PermissionKindCamera)
    if status == application.PermissionStatusNotDetermined {
        status, _ = app.Permissions.Request(application.PermissionKindCamera)
    }
    if status == application.PermissionStatusDenied {
        _ = app.Permissions.OpenSystemSettings(application.PermissionKindCamera)
    }
}()
```

전체 디스크 접근 권한은 요청할 수 없으며 `ErrPermissionNotRequestable`을 반환합니다. 사용자를 설정 패널로 안내하세요. 응답을 받지 못한 요청은 `ErrPermissionRequestTimeout`을 반환합니다. macOS에서는 대개 해당 권한의 사용 목적 설명 키가 `Info.plist`에 없다는 뜻입니다.

이제 macOS에서는 `Permissions` 창 옵션이 적용됩니다. 이 옵션은 페이지의 `getUserMedia` 요청을 처리하는 방식을 결정합니다. `PermissionAllow`는 WebView 자체의 권한 요청을 건너뛰고, `PermissionDeny`는 묻지 않고 거부하며, `PermissionDefault`는 권한 요청을 표시합니다. 카메라나 마이크를 처음 사용할 때는 시스템 수준의 TCC 권한 요청이 여전히 나타납니다.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Meeting",
    URL:   "/",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionDefault,
    },
})
```

### 전원

`app.Power.PreventSleep`은 반환된 해제 함수를 호출할 때까지 시스템이 잠자기 상태로 들어가지 않도록 합니다. `Display`를 사용하면 화면도 켜 둡니다. 유지 요청은 개별적으로 계산되므로 앱의 여러 부분에서 동시에 사용할 수 있습니다. 이유는 Activity Monitor에 표시됩니다.

```go
release, err := app.Power.PreventSleep("Exporting video", application.PreventSleepOptions{Display: true})
if err != nil {
    log.Println(err)
}
defer release()

state := app.Power.State()
if state.LowPowerMode || state.ThermalState >= application.ThermalStateSerious {
    // trim background work
}
log.Println("battery", state.BatteryLevel, "charging", state.Charging, "on battery", state.OnBattery)
```

변경 사항은 `events.Mac.ApplicationDidChangePowerState`(저전력 모드 전환)와 `events.Mac.ApplicationDidChangeThermalState`로 전달됩니다.

### 수명 주기

앱이 `NSSupportsSuddenTermination`으로 허용하면 macOS는 로그아웃이나 종료 시 유휴 앱을 즉시 종료할 수 있습니다. `NSSupportsAutomaticTermination`으로 허용하면 창이 없는 유휴 앱을 종료할 수 있습니다. `app.Lifecycle.HoldTermination`은 파일 저장과 같은 중요한 구간에서 두 종료 방식을 모두 일시 중지합니다.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled`는 실행 중 갑작스러운 종료를 켜거나 끕니다. `SuddenTerminationEnabled`는 `Info.plist` 키의 값으로 시작하는 현재 상태를 보고합니다.

### 환경

`app.Env`에는 사용자 설정에 관한 세 가지 조회 기능이 추가됩니다.

```go
a11y := app.Env.Accessibility()
if a11y.ReduceMotion || a11y.ReduceTransparency {
    app.Event.Emit("theme:calm", true)
}

layout := app.Env.KeyboardLayout()
log.Println(layout.ID, layout.Name, layout.Languages)

locale := app.Env.Locale()
log.Println(locale.Identifier, locale.Language, locale.Region, locale.Preferred)
```

- `Accessibility`는 동작 줄이기, 투명도 줄이기, 대비 증가, 색상 없이 구별하기, 색상 반전, VoiceOver 및 스위치 제어 설정을 반영합니다.
- `KeyboardLayout`은 식별자, 현지화된 이름 및 언어와 함께 현재 입력 소스를 반환합니다.
- `Locale`은 AppKit이 앱에 대해 선택한 로캘과 사용자의 전체 `Preferred` 목록을 우선순위대로 반환합니다. `Identifier`는 번들이 `CFBundleLocalizations`에 선언한 언어만 반영하므로 직접 언어를 선택하려면 `Preferred`를 사용하세요.

### 이벤트

다음 애플리케이션 이벤트가 새로 추가되었습니다. 각 이벤트는 `app.Event.OnApplicationEvent`를 통해 전달됩니다. 최신 값은 해당 관리자를 조회하세요.

| 이벤트 | 발생 시점 | 조회 방법 |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | 저전력 모드가 전환될 때 | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | 열 부하 상태가 변할 때 | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | 손쉬운 사용 디스플레이 설정이 변할 때 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | 동일한 변경이 macOS 이벤트로 전달될 때 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | 입력 소스가 변할 때 | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | 로캘이 변할 때 | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## 통합

### Services 메뉴

`app.ServicesProvider.Register`는 모든 macOS 애플리케이션이 선택한 텍스트나 파일에 대해 표시하는 Services 하위 메뉴에 항목을 추가합니다. 처리기는 클립보드를 `ServiceRequest`로 받고, 다시 쓸 내용을 `ServiceResponse`로 반환합니다. 빈 응답은 선택 항목을 변경하지 않습니다. AppKit이 메인 스레드에서 처리기를 기다리므로 빠르게 완료되도록 하세요.

```go
err := app.ServicesProvider.Register(application.ServiceDefinition{
    Name:          "summarise",
    MenuTitle:     "Summarise with Notes",
    SendTypes:     []string{"public.utf8-plain-text"},
    ReturnTypes:   []string{"public.utf8-plain-text"},
    KeyEquivalent: "S",
    Handler: func(_ *application.Context, request application.ServiceRequest) (application.ServiceResponse, error) {
        return application.ServiceResponse{Text: "Summary: " + request.Text}, nil
    },
})
if err != nil {
    log.Println(err)
}

// Paste this inside the top-level <dict> of Info.plist.
log.Println(app.ServicesProvider.InfoPlistXML())
```

`Name`은 AppKit이 보내는 메시지이며 일반 식별자여야 합니다. `SendTypes`와 `ReturnTypes`는 클립보드 유형이며 서비스에는 둘 중 하나 이상이 필요합니다. 등록만으로는 서비스가 표시되지 않습니다. 번들의 `Info.plist`에서 `NSServices` 아래에 선언해야 합니다. `InfoPlistXML`은 바로 붙여 넣을 수 있는 블록을 반환하고, `InfoPlistEntries`는 plist 직렬화기에 사용할 수 있도록 동일한 데이터를 맵으로 반환합니다. 해당 항목의 `NSPortName`은 애플리케이션의 `Name`이며 `CFBundleName`과 일치해야 합니다.

Wails CLI로 빌드하는 프로젝트는 `build/config.yml`에 동일한 서비스를 한 번 선언하고 패키징 과정에서 `NSServices` 블록을 생성하도록 할 수 있습니다:

```yaml
services:
  - name: SummariseText
    menuTitle: Summarise with My Product
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S
```

각 항목은 Go에 등록된 `ServiceDefinition`과 `Name`이 같아야 합니다. 새 빌드를 설치한 뒤 `pbs -update`를 실행하면 로그아웃하지 않고도 Services 메뉴에 변경 사항이 반영됩니다.

### Handoff와 사용자 활동

`app.Activity.Publish`는 `NSUserActivity`를 현재 활동으로 설정하여 사용자가 다른 기기에서 이어서 작업하거나 Spotlight에서 찾거나 Siri의 제안을 받을 수 있게 합니다. 반환된 `PublishedActivity`는 상태가 변할 때 업데이트하고 문서를 닫을 때 무효화할 수 있습니다. 새 활동을 게시하면 이전 활동을 대체합니다.

```go
activity, err := app.Activity.Publish(application.UserActivity{
    Type:               "com.example.notes.editing",
    Title:              "Editing Quarterly report",
    UserInfo:           map[string]any{"note": "quarterly-report"},
    WebpageURL:         "https://example.com/notes/quarterly-report",
    EligibleForHandoff: true,
    EligibleForSearch:  true,
    Keywords:           []string{"report", "quarterly"},
})
if err != nil {
    log.Println(err)
    return
}
if err := activity.Update(map[string]any{"note": "quarterly-report", "cursor": 120}); err != nil {
    log.Println(err)
}
// When the document closes:
activity.Invalidate()
```

들어오는 활동은 `OnContinue`로 전달됩니다. 유니버설 링크는 `UserActivityTypeBrowsingWeb` 유형을 사용하고 `WebpageURL`에 페이지를 담습니다. 앱에서 URL 처리를 하나의 코드 경로로 유지할 수 있도록 `events.Common.ApplicationLaunchedWithUrl`로도 전달됩니다. `OnWillContinue`, `OnFailed`, `OnUpdated`는 나머지 델리게이트 동작을 처리합니다.

```go
app.Activity.OnContinue(func(_ *application.Context, incoming application.UserActivity) bool {
    if incoming.Type == application.UserActivityTypeBrowsingWeb {
        log.Println("universal link", incoming.WebpageURL)
        return true
    }
    note, _ := incoming.UserInfo["note"].(string)
    log.Println("continue editing", note)
    return true
})
```

모든 활동 유형은 `Info.plist`의 `NSUserActivityTypes`에 나열해야 합니다. 유니버설 링크에는 `applinks:example.com` 항목이 포함된 `com.apple.developer.associated-domains` 권한과 해당 도메인의 일치하는 `apple-app-site-association` 파일도 필요합니다.

### Apple Events

`app.AppleEvents.Handle`은 이벤트 클래스와 ID에 대한 처리기를 등록하여 AppleScript, Shortcuts 및 다른 애플리케이션에서 앱을 제어할 수 있게 합니다. 코드는 네 글자 문자열입니다. 직접 매개변수는 Go 값(`string`, 파일 경로의 `[]string`, `int64`, `float64`, `bool`, `[]any` 또는 `AppleEventRawData`)으로 디코딩되며 응답의 `Result`에도 같은 종류를 사용할 수 있습니다. 이벤트가 일시 중지된 동안 처리기는 자체 goroutine에서 실행됩니다.

```go
err := app.AppleEvents.Handle("WAIL", "note", func(_ *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
    text, _ := event.DirectObject.(string)
    return application.AppleEventReply{Result: "noted: " + text}, nil
})
if err != nil {
    log.Println(err)
}

go func() {
    result, err := app.AppleEvents.Send("com.apple.finder", "misc", "actv", nil)
    log.Println(result, err)
}()

sdef := app.AppleEvents.ScriptingDefinition()
if err := os.WriteFile("Notes.sdef", []byte(sdef), 0o644); err != nil {
    log.Println(err)
}
```

스크립트는 원시 이벤트 구문으로 처리기를 바로 호출할 수 있습니다:

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition`은 각 처리기에 명령 이름을 부여하는 최소한의 `.sdef`를 생성합니다. 이를 `Contents/Resources`에 포함하고 `Info.plist`의 `NSAppleScriptEnabled`와 `OSAScriptingDefinition`으로 지정하세요. 그러면 Script Editor에서 File > Open Dictionary 아래에 표시됩니다. Wails는 사용자 지정 URL 스킴의 Get URL 이벤트를 이미 처리합니다. `"GURL"`/`"GURL"` 처리기는 기존 처리에 연결되고, `"aevt"`/`"odoc"` 처리기는 기본 제공 Open Documents 전달을 대체합니다. `Send`는 번들 식별자로 실행 중인 애플리케이션을 대상으로 합니다. 번들 앱에는 `NSAppleEventsUsageDescription`이 필요하며 응답이 도착할 때까지 호출한 goroutine을 차단합니다.

### Quick Look

`app.QuickLook.Preview`는 하나 이상의 파일에 대해 공유 Quick Look 패널을 엽니다. 경로가 여러 개이면 패널에 파일 간 이동 화살표가 표시됩니다. `Thumbnail`은 시스템의 썸네일 제공자를 통해 파일을 렌더링하고 PNG를 반환하므로 문서, 이미지, PDF 및 동영상에 사용할 수 있습니다.

```go
if err := app.QuickLook.Preview([]string{"/Users/me/Documents/Report.pdf"}); err != nil {
    log.Println(err)
}

go func() {
    png, err := app.QuickLook.Thumbnail("/Users/me/Documents/Report.pdf", application.ThumbnailOptions{
        Width: 256,
        Scale: 2,
    })
    if err != nil {
        log.Println(err)
        return
    }
    if err := os.WriteFile("thumbnail.png", png, 0o644); err != nil {
        log.Println(err)
    }
}()
```

경로는 절대 경로여야 하며 파일이 존재해야 합니다. `ClosePreview`와 `IsPreviewOpen`으로 패널을 관리합니다. `ThumbnailOptions.IconMode`는 Finder 스타일의 문서 테두리를 그리고 `Scale: 2`는 Retina 이미지를 생성합니다. `Thumbnail`은 호출한 goroutine을 차단하므로 goroutine 또는 바인딩된 메서드에서 호출하세요.

### 작업 공간 도우미

`app.Browser`에는 `NSWorkspace`를 활용하는 도우미 세 가지가 추가됩니다. `OpenWith`는 번들 식별자 또는 번들 경로로 지정한 애플리케이션으로 파일을 엽니다. `ApplicationsForFile`은 파일을 열 수 있는 설치된 애플리케이션을 기본 처리기부터 나열합니다. `ActivateApplication`은 실행 중인 애플리케이션을 맨 앞으로 가져옵니다.

```go
apps := app.Browser.ApplicationsForFile("/Users/me/Documents/Report.md")
for _, info := range apps {
    log.Println(info.Name, info.BundleID, info.Path)
}

if err := app.Browser.OpenWith("/Users/me/Documents/Report.md", "com.apple.TextEdit"); err != nil {
    log.Println(err)
}

if err := app.Browser.ActivateApplication("com.apple.TextEdit"); err != nil {
    log.Println(err)
}
```

### Spotlight

`app.Spotlight.Index`는 Core Spotlight를 통해 애플리케이션 콘텐츠를 시스템 검색 색인에 추가합니다. 각 `SearchableItem`에는 `ID`, `Title`이 있으며, 선택적으로 일괄 제거를 위한 `Domain`, `Description`, `Keywords`, `ContentType`, PNG 썸네일, 딥 링크 `URL`, 만료 시간을 지정할 수 있습니다. 사용자가 Spotlight에서 항목을 선택하거나 검색어로 "앱에서 검색"을 선택하면 `OnOpen`이 호출됩니다.

```go
err := app.Spotlight.Index([]application.SearchableItem{{
    ID:          "note:quarterly-report",
    Domain:      "notes",
    Title:       "Quarterly report",
    Description: "Draft for the board meeting",
    Keywords:    []string{"finance", "q3"},
    ContentType: "public.plain-text",
    URL:         "notes://open/quarterly-report",
}})
if err != nil {
    log.Println(err)
}

stop := app.Spotlight.OnOpen(func(_ *application.Context, id string, query string) {
    if query != "" {
        log.Println("search in app:", query)
        return
    }
    log.Println("open item", id)
})
defer stop()

// Later:
_ = app.Spotlight.Delete([]string{"note:quarterly-report"})
_ = app.Spotlight.DeleteDomain("notes")
```

`IsAvailable`은 색인이 항목을 받아들이는지 보고합니다. 색인 생성에는 번들 애플리케이션이 필요합니다. 번들로 묶이지 않은 `go run` 바이너리에서 색인한 항목은 Spotlight에 표시되지 않습니다. `DeleteAll`은 앱이 색인한 모든 항목을 제거합니다.

## 향후 추가 예정

시트, 팝오버, 표시 옵션 및 상태 복원을 다루는 `mac-windows-extra` 예제를 추가 중이며, 반영되면 이곳에 링크합니다.

## 버전 요구 사항

이 페이지의 모든 기능은 macOS 전용이며 Go API는 모든 플랫폼에서 동일합니다. Wails는 macOS 10.13 이상을 대상으로 합니다. 더 새로운 릴리스가 필요한 기능은 아래와 같이 제한됩니다.

| 기능 | 최소 macOS 버전 | 이전 릴리스에서의 동작 |
|---------|---------------|-------------------------------|
| 카메라 및 마이크 권한 상태 | 10.14 | 허용된 것으로 보고됨(이전 릴리스는 캡처 장치 접근을 제한하지 않음) |
| `Speech.Speak`와 `Voices` | 10.14 | `ErrSpeechNotSupported` |
| 화면 녹화 및 입력 모니터링 권한 | 10.15 | 허용된 것으로 보고됨 |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | 디버그 로그 항목을 남기고 무시됨 |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| 메뉴 항목 및 상태 항목의 `SetSymbol` | 11 | 이미지가 표시되지 않음 |
| UTI를 사용하는 `AddContentType`, `SetFormats` | 11 | 동일한 식별자가 기존의 허용 파일 유형 API를 통해 적용됨 |
| `StartDrag`의 파일 프로미스 아이콘 | 11 | 일반 문서 아이콘 |
| `PowerState.LowPowerMode` 및 관련 이벤트 | 12 | 항상 false이며 이벤트가 발생하지 않음 |
| `InteractionState`와 `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| `OpenSystemSettings`의 알림 패널 | 13 | 이전 알림 환경설정 패널이 열림 |
| 메뉴 배지, 섹션 머리글, 팔레트 | 14 | 배지는 표시되지 않고 머리글은 비활성화된 항목이며 팔레트는 숨겨짐 |
| 사용자 지정 뷰가 없는 항목의 `MacToolbarItem.ShowPopover` | 14 | `ErrMacPopoverAnchorUnavailable` |

그 밖의 모든 기능에는 Wails 최소 버전 외의 추가 요구 사항이 없습니다.

## Info.plist 키

여러 기능은 애플리케이션의 `Info.plist`에 있는 키에 의존합니다. 사용 목적 설명은 권한 요청에 표시됩니다. 키가 없으면 macOS는 권한 요청을 표시하지 않으며 요청 시간이 초과됩니다.

| 키 | 필요한 기능 |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, 페이지의 카메라 접근 |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, 페이지의 마이크 접근, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | `Lifecycle.SuddenTerminationEnabled`의 초기 상태. `HoldTermination`은 이를 일시 중지함 |
| `NSSupportsAutomaticTermination` | macOS가 유휴 앱을 종료하도록 허용함. `HoldTermination`은 이를 일시 중지함 |
| `CFBundleLocalizations` | `Env.Locale().Identifier`가 보고할 수 있는 언어 |
| `NSServices` | `app.ServicesProvider`에 등록한 서비스마다 항목 하나. `InfoPlistXML`로 블록 생성 |
| `NSUserActivityTypes` | `app.Activity`를 통해 게시하거나 이어받는 모든 `UserActivity.Type` |
| `NSAppleScriptEnabled`와 `OSAScriptingDefinition` | 앱을 스크립트로 제어할 수 있음을 표시하고 `app.AppleEvents.ScriptingDefinition`에서 작성한 `.sdef`를 지정함 |
| `NSAppleEventsUsageDescription` | 다른 애플리케이션을 대상으로 하는 `app.AppleEvents.Send` |

알림 권한과 음성 인식은 앱이 번들 식별자가 있는 번들로 실행되어야 합니다. 번들로 묶이지 않은 `go run` 바이너리는 알림에 대해 `PermissionStatusUnsupported`를 보고합니다.

<a id="platform-notes"></a>

## 플랫폼 참고 사항

이 페이지의 모든 API는 Windows와 Linux에서도 컴파일됩니다. macOS 이외 플랫폼에서는 다음과 같이 작동합니다:

- 창 추가 기능: `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName`, `SetWindowButtonsOffset`은 아무 작업도 하지 않습니다. `WindowCascade`는 `WindowCentered`처럼 작동합니다. `RequestAttention`은 `Cancel`이 아무 작업도 하지 않는 핸들을 반환합니다. `PrintWithOptions`는 `Print`를 호출합니다. `ExportPDF`와 `Snapshot`은 `ErrMacOnly`를 반환합니다.
- 메뉴: 기호, 배지, 혼합 상태, 대체 항목 및 들여쓰기는 저장되지만 그려지지 않습니다. 섹션 머리글은 비활성화된 항목입니다. 팔레트는 숨겨집니다. 최근 항목 열기 하위 메뉴는 메뉴가 생성될 때 Go의 최근 항목 목록에서 채워집니다. Dock 메뉴는 표시되지 않습니다.
- 대화상자: `SetSuppression`, `SetHelp`, `SetNameFieldLabel`, `SetTags`는 무시됩니다. `AddContentType`과 `SetFormats`는 변환 방법이 알려져 있으면 확장자 필터가 됩니다. `Prompt`, `PickColor`, `PickFont`는 `ErrDialogNotSupported`를 반환합니다.
- 상태 항목: `SetSymbol`, `SetRemovable`, `OnVisibilityChange`는 효과가 없습니다. `IsVisible`은 마지막 `Show` 또는 `Hide` 호출을 반영합니다.
- 피드백: 햅틱은 iOS와 Android에서 작동하며 Windows와 Linux에서는 아무 작업도 하지 않습니다. `Sound.Beep`과 `Sound.Play`는 Windows에서 WAV 파일 및 레지스트리 별칭과 함께 작동합니다. 그 밖의 플랫폼에서 `Play`는 `ErrSoundNotSupported`를 반환합니다. 음성 기능은 `ErrSpeechNotSupported`와 `ErrSpeechRecognitionNotSupported`를 반환합니다.
- 클립보드: 다양한 형식을 지원하는 메서드는 `ErrClipboardNotSupported`를 반환하고, `Types`는 비어 있으며, `ChangeCount`는 0이고, `OnChange`는 실행되지 않습니다.
- 드래그: `StartDrag`는 `ErrDragOutUnsupported`를 반환합니다. 파일 이외의 드롭 유형은 전달되지 않으며, 파일 드롭은 계속 `WindowFilesDropped`를 통해 작동합니다.
- 시스템: `Permissions.Status`는 `PermissionStatusUnsupported`를 보고하고 `Request`는 `ErrPermissionsUnsupported`를 반환합니다. `PreventSleep`은 `ErrPreventSleepUnsupported`와 아무 작업도 하지 않는 해제 함수를 반환합니다. `HoldTermination`은 아무 작업도 하지 않는 해제 함수를 반환합니다. `Accessibility`의 모든 값은 false이고 `KeyboardLayout`은 0 값이며 `Locale`은 `LC_ALL`, `LC_MESSAGES`, `LANG`에서 파생됩니다. `Permissions` 창 옵션은 크로스 플랫폼입니다.
- 통합: `ServicesProvider.Register`는 `ErrServicesUnsupported`를, `Activity.Publish`는 `ErrActivityUnsupported`를 반환합니다. `InfoPlistXML`, `InfoPlistEntries` 및 활동 처리기는 계속 작동합니다. `AppleEvents`의 `Handle`과 `Send`는 `ErrAppleEventsNotSupported`를 반환합니다. `ScriptingDefinition`은 모든 플랫폼에서 생성됩니다. `QuickLook.Preview`와 `Thumbnail`은 `ErrQuickLookNotSupported`를 반환합니다. `Spotlight` 색인 메서드는 `ErrSpotlightNotSupported`를 반환하고 `OnOpen`은 실행되지 않습니다. `Browser.OpenWith`는 지정된 실행 파일에 경로를 인수로 전달해 시작하며, `ApplicationsForFile`은 빈 목록이고 `ActivateApplication`은 `ErrApplicationNotRunning`을 반환합니다.
- 표시 옵션: `SetPresentationOptions`는 `ErrMacOnly`를 반환하고 `PresentationOptions`는 `MacPresentationDefault`입니다.
- 시트: `PresentSheet`, `PresentCriticalSheet`, `PresentNativeSheet`는 `ErrMacSheetUnsupported`를 반환합니다. `EndSheet`는 아무 작업도 하지 않고 조회 메서드는 시트가 없다고 보고합니다.
- 팝오버: `NewMacPopover`는 작동하지만 표시 메서드는 `ErrMacPopoverUnsupported`를 반환하고 `IsShown`은 false입니다.
- 상태 복원: `SetRestorationID`와 `SetRestorationData`는 아무 작업도 하지 않고 `OnRestore`는 호출되지 않으며 `InteractionState`와 `RestoreInteractionState`는 `ErrMacOnly`를 반환합니다.

## 예제

각 예제는 완전하게 실행할 수 있는 애플리케이션입니다:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go`가 수정 표시 점, 부제목, PDF 내보내기, 인쇄 옵션 및 계단식 창을 메모 편집기에 연결합니다.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): 기호, 배지, 섹션 머리글, 혼합 상태, 대체 항목, 팔레트, 최근 항목 열기, 동적 Dock 메뉴 및 Dock 진행률을 다룹니다.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): 다시 표시하지 않기와 도움말 버튼, 텍스트 입력 프롬프트, 콘텐츠 유형, 형식 팝업, Finder 태그, 색상 및 글꼴 패널을 다룹니다.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): 제거 가능한 SF Symbol 상태 항목, 햅틱, 시스템 소리, 텍스트 음성 변환 및 음성 인식을 다룹니다.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): 변경 추적이 가능한 다양한 형식의 클립보드, 파일 프로미스를 사용한 외부로 드래그하기, 텍스트·URL·이미지 드롭을 다룹니다.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): 권한, 잠자기 방지, 종료 일시 중지, 전원 상태, 손쉬운 사용, 키보드 레이아웃 및 로캘의 실시간 업데이트를 다룹니다.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): Services 메뉴 항목, 이어받기 처리기가 있는 Handoff 활동, `app.Browser`의 작업 공간 도우미를 다룹니다.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): `OnOpen`을 사용하는 Spotlight 색인, Quick Look 미리보기와 썸네일, 스크립팅 정의가 있는 사용자 지정 Apple Event를 다룹니다.
