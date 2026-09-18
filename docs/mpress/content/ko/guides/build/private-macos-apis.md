---
title: "비공개 macOS API"
description: "비공개 macOS API에 의존하는 모든 Wails 기능과 옵션 및 명시적 활성화 명령과 공개 빌드 대체 동작입니다."
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails는 기본적으로 공개 macOS API를 사용합니다. 단일 Go 빌드 태그 `private_mac_apis`를 사용하면 이 페이지에 나열된 비공개 WebKit 및 AppKit 호출이 활성화됩니다. 모든 공개 Go 옵션과 메서드는 어느 빌드에서든 계속 사용할 수 있습니다. 이 태그가 없으면 비공개 API 전용 작업은 아무 동작도 하지 않으며, 공개 대안이 있는 기능은 해당 대안을 사용합니다.

@note{type="caution" title="비공개 macOS 동작 활성화"}
창 옵션을 설정하는 것만으로는 비공개 API가 활성화되지 않습니다. 이를 활성화하려면 빌드 명령에 `private_mac_apis`를 추가하세요. 이 태그는 macOS 데스크톱 빌드에만 적용되며 iOS, Android, Windows, Linux 또는 서버 빌드에는 적용되지 않습니다.

@end

## 비공개 API 활성화

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

프로덕션 빌드를 직접 수행하려면 `go build -tags production,private_mac_apis .`를 사용하세요. 프런트엔드 예제는 해당 README의 안내에 따라 바인딩과 애셋을 빌드한 후 실행하세요. 사용자 지정 Taskfile이나 이전 버전의 Taskfile은 `EXTRA_TAGS`를 Go 컴파일러에 전달해야 합니다.

## 기능 목록

| 기능 또는 값 | `private_mac_apis`로 활성화되는 항목 | 태그가 없을 때 |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | 네이티브 창 위의 투명한 WKWebView | 네이티브 창은 구성되지만 웹뷰는 불투명한 상태로 유지됩니다. |
| `Mac.Backdrop: MacBackdropTranslucent` | 네이티브 블러 효과가 비쳐 보이도록 하는 투명한 WKWebView | 불투명한 웹뷰 뒤에 블러 효과가 구성됩니다. |
| `Mac.Backdrop: MacBackdropLiquidGlass` | 비공개 웹뷰 배경 제어를 사용하여 글라스 레이어 위에 표시되는 투명한 WKWebView | 불투명한 웹뷰 뒤에 글라스 레이어가 구성되며 스타일 지정에는 공개 대안을 사용합니다. |
| Liquid Glass 설정 중 웹뷰 배경 지우기 | 비공개 WebKit `backgroundColor` 제어 | macOS 12 이상에서는 공개 `underPageBackgroundColor`를 사용하고, 이전 macOS에서는 레이어 색상을 사용합니다. 웹뷰를 투명하게 만들지는 않습니다. |
| `app.Window.NewNotchWindow(...)` | 형태가 지정된 노치 패널 내부의 투명한 웹뷰 | 배치와 애니메이션을 포함해 패널은 계속 작동하지만 웹뷰는 불투명한 상태로 유지됩니다. |
| `Mac.LiquidGlass.Style` | 문서화되지 않은 다크 스타일 값을 포함한 Wails의 기존 네이티브 스타일 매핑 | 공개 regular/clear 스타일과 light/dark 모양을 사용합니다. 아래 값 표를 참조하세요. |
| `Mac.LiquidGlass.GroupID` | 비어 있지 않은 식별자에 대해 비공개 글라스 그룹화를 요청합니다. | 무시되며 그룹화를 요청하지 않습니다. |
| `Mac.LiquidGlass.GroupSpacing` | 0보다 큰 값에 대해 비공개 그룹 간격을 요청합니다. | 무시됩니다. |
| `window.OpenDevTools()` 및 JavaScript `Window.OpenDevTools()` | macOS 12 이상에서 프로그래밍 방식으로 WebKit 인스펙터를 엽니다. | 아무 동작도 하지 않습니다. |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | macOS 12 이상에서 창이 처음 표시될 때 프로그래밍 방식으로 인스펙터를 열도록 요청합니다. | 아무 동작도 하지 않습니다. |
| macOS 13.3 이전 버전용 레거시 인스펙터 활성화 | 인스펙터 지원이 컴파일에 포함된 경우 WebKit 개발자용 추가 기능을 활성화합니다. | 아무 동작도 하지 않습니다. 공개 Safari 검사를 사용하려면 macOS 13.3 이상이 필요합니다. |

## 웹뷰 투명도와 배경

**비공개 API 필요:** `MacBackdropTransparent`, `MacBackdropTranslucent`, `MacBackdropLiquidGlass` 및 노치 창에서 사용하는 웹뷰 투명도입니다. 내부적으로 Wails는 WebKit의 비공개 `drawsBackground` 키를 설정합니다. HTML 또는 CSS 배경만 투명하게 설정해서는 불투명한 네이티브 WKWebView를 투명하게 만들 수 없습니다.

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

비공개 웹뷰 배경색 작업은 WebKit의 `backgroundColor` 키를 사용합니다. Liquid Glass 설정에서는 이를 사용하여 웹뷰 배경을 지웁니다. 태그가 없으면 이 내부 작업은 macOS 12 이상에서 공개 `underPageBackgroundColor`를 사용하고, 이전 macOS에서는 뷰의 레이어를 사용합니다. 이러한 대안으로는 웹뷰를 투명하게 만들 수 없습니다.

`WebviewWindowOptions.BackgroundColour` 및 `window.SetBackgroundColour()`는 macOS에서 <strong>네이티브 창</strong>의 색상을 설정하며, 그 자체로는 비공개 API가 필요하지 않습니다. 마찬가지로 `Frameless` 및 `Mac.TitleBar.AppearsTransparent`는 공개 AppKit API를 사용합니다. 비공개 API에 의존하는 것은 제목 표시줄 투명도가 아니라 웹뷰 투명도입니다. macOS 배경 효과를 적용하려면 `BackgroundType`에만 의존하지 말고 `Mac.Backdrop`를 구성하세요.

[창 옵션](/features/windows/options/#mac-options), [프레임 없는 창](/features/windows/frameless/#with-transparent-background) 및 [노치 창](/features/windows/notch-windows/)을 참조하세요.

## Liquid Glass 값

네이티브 `NSGlassEffectView`를 사용할 수 있는 경우(macOS 26 이상) 다음 매핑이 적용됩니다. 네이티브 스타일 값 중 `0`(regular)과 `1`(clear)만 문서화되어 있습니다. Go 상수는 두 빌드 모두에서 기존 값을 유지합니다.

| `MacLiquidGlassStyle` 값 | `private_mac_apis` 사용 시 | 태그가 없을 때 |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic` (`0`) | 네이티브 regular 스타일(`0`) | 네이티브 regular 스타일(`0`) |
| `LiquidGlassStyleLight` (`1`) | 기존 네이티브 스타일 매핑(`1`, clear) | Aqua 모양을 적용한 네이티브 regular 스타일(`0`) |
| `LiquidGlassStyleDark` (`2`) | **문서화되지 않은 네이티브 스타일 값 `2`** | Dark Aqua 모양이 적용된 네이티브 일반 스타일(`0`) |
| `LiquidGlassStyleVibrant`(`3`) | 기존 밝은/네이티브 투명 스타일(`1`)에 매핑됨 | 네이티브 투명 스타일(`1`) |

Automatic과 Vibrant는 문서화된 네이티브 스타일 값을 사용하지만, 전체 창 Liquid Glass 배경을 표시하려면 여전히 <strong>웹뷰 투명도</strong>를 위한 태그가 필요합니다. Light의 모양은 두 빌드에서 서로 다릅니다. 비공개 스타일 매핑으로 향후 macOS 버전에서도 동일한 효과가 렌더링된다고 보장할 수는 없습니다.

**항상 비공개:** `GroupID` 및 `GroupSpacing`. Wails는 그룹화를 요청하기 전에 비공개 `setGroupIdentifier:`, `setGroupName:` 및 `setGroupSpacing:` 셀렉터가 있는지 확인합니다. 태그가 없으면 이러한 작업은 아무 효과도 없습니다. 태그를 활성화해도 실행 중인 macOS 버전에서 이러한 셀렉터를 지원한다고 보장되지는 않습니다.

`MacLiquidGlass.Material`, `CornerRadius` 및 `TintColor` 자체에는 비공개 API가 필요하지 않습니다. 네이티브 Liquid Glass를 지원하지 않는 macOS 버전에서는 Wails가 반투명 대체 효과를 구성하지만, 이를 웹뷰를 통해 표시하려면 여전히 태그가 필요합니다.

## 웹 검사기

**비공개 API 필요:** 애플리케이션에서 WebKit 검사기를 열기 위해 `OpenDevTools()`을 호출하거나 `OpenInspectorOnStartup: true`을 설정하는 경우입니다. Wails는 비공개 `_inspector` 셀렉터를 사용합니다. `private_mac_apis`이 없으면 이러한 작업은 아무 효과 없이 조용히 무시됩니다.

검사기 지원도 빌드에 포함되어야 합니다. 기존 `production` 및 `devtools` 태그의 의미는 그대로 유지됩니다.

| 빌드 태그 | 프로그래밍 방식으로 검사기 열기 | macOS 13.3 이상에서 공개 API를 사용하는 Safari 검사 |
| --- | --- | --- |
| 없음 | 효과 없음 | 활성화됨 |
| `private_mac_apis` | macOS 12 이상에서 활성화됨 | 활성화됨 |
| `production` | 효과 없음 | 비활성화됨 |
| `production,private_mac_apis` | 효과 없음 | 비활성화됨 |
| `production,devtools` | 효과 없음 | 활성화됨 |
| `production,devtools,private_mac_apis` | macOS 12 이상에서 활성화됨 | 활성화됨 |

macOS 13.3 이상에서 Wails는 공개 `WKWebView.inspectable`을 사용해 Safari 검사를 활성화하며, 여기에는 비공개 API가 필요하지 않습니다. macOS 13.3 이전 버전에서는 대체 구현이 비공개 `developerExtrasEnabled` 환경설정을 사용하므로, 빌드에 포함된 검사기 지원뿐 아니라 `private_mac_apis`도 필요합니다.

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## 이 목록 유지 관리

모든 비공개 네이티브 호출은 `v3/pkg/application/mac_private_api_darwin.go`에 격리되어 있으며, 기본 빌드에서는 `mac_public_api_darwin.go`을 선택합니다. 이 목록은 투명도, 웹뷰 배경색, 글래스 스타일, 글래스 그룹화, 검사기 열기 및 레거시 검사기 활성화를 다룹니다. 이러한 구현을 변경할 때는 이 페이지와 영향을 받는 옵션 문서를 함께 업데이트해야 합니다.
