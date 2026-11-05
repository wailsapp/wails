---
title: "원격 데스크톱(RDP)에서 WebView2가 멈추는 문제"
description: "세션 도중 모니터 DPI가 변경되는 RDP 세션에서 Wails 앱을 실행할 때 WebView2 UI가 수초간 멈추는 문제를 해결합니다."
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## 문제

원격 데스크톱(RDP) 세션에서 Wails 애플리케이션을 사용하면 다음과 같은 일반적인 상호 작용 중에 UI가 수초간 멈출 수 있습니다.

- 팝업 창을 열기 위해 클릭한 후 콘텐츠가 표시되기까지 약 4~8초가 걸립니다.
- 창을 닫으면 상위 창이 약 2초 동안 응답하지 않습니다.
- 느려진 상태는 다시 연결해도 지속되며 호스트 머신을 재부팅해야만 해제됩니다.

이 문제는 세션 도중 Retina에 최적화된 가상 모니터를 프로비저닝하는 iOS용 Microsoft Remote Desktop 클라이언트에서 가장 흔히 발생합니다. 세션 도중 DPI 컨텍스트가 다른 모니터를 추가하는 모든 RDP 클라이언트에서 같은 동작이 발생할 수 있습니다.

## 발생 원인

기본적으로 WebView2는 컴포지터 표면이 자식 창에 존재하는 창 기반 호스팅을 사용합니다. RDP 클라이언트가 세션과 DPI 컨텍스트가 다른 모니터를 추가하면 WebView2 컨트롤러를 호출할 때마다(`PutIsVisible`, `MoveFocus`, 첫 화면 그리기 및 표면 해제) 동기식 DirectComposition 리마샬링이 강제로 수행됩니다. 각 리마샬링은 UI 스레드를 약 2초 동안 차단하므로 팝업을 많이 사용하는 앱에서는 멈춤 시간이 누적됩니다.

동일한 머신에서 실행되는 네이티브 Win32 및 WebView2 애플리케이션은 비주얼 호스팅을 사용하므로 영향을 받지 않습니다. 이는 일반적인 WebView2 또는 Windows 컴포지터 문제가 아니라 호스팅 모드가 원인임을 나타냅니다.

## 해결 방법

Windows 옵션에서 `UseVisualHosting`을 설정하여 비주얼 호스팅을 활성화하세요. 비주얼 호스팅에서는 호스트가 소유한 DirectComposition 비주얼을 통해 WebView2의 컴포지터 표면을 소유하므로 DPI 컨텍스트가 변경되어도 더 이상 동기식 리마샬링이 발생하지 않습니다.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

활성화하면 팝업이 일반적인 탐색 시간(약 150~500ms) 내에 열리고 창을 닫아도 더 이상 상위 창이 차단되지 않습니다.

@note{type="caution"}
`UseVisualHosting`은 `app.Run()`보다 먼저 설정해야 합니다. Wails는 애플리케이션 시작 시 이 값을 읽고 WebView2 환경을 초기화하기 전에 `COREWEBVIEW2_FORCED_HOSTING_MODE` 환경 변수를 `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL`으로 설정합니다. 나중에 설정하면 적용되지 않습니다.

@end

이 옵션의 기본값은 `false`이므로 창 기반 호스팅이 계속 기본 방식으로 사용됩니다. 기존 앱은 직접 활성화하지 않는 한 동작이 변경되지 않습니다.

## 활성화해야 하는 경우

앱이 RDP, 특히 iOS용 Microsoft Remote Desktop 클라이언트를 통해 정기적으로 사용되고 창을 열거나 닫을 때 수초간 멈춘다면 `UseVisualHosting: true`을 설정하세요. 앱을 RDP에서 실행하지 않는다면 이 옵션은 필요하지 않으며 기본값을 그대로 사용해도 됩니다.

## 참고 자료

- [WebView2: 창 기반 호스팅과 비주얼 호스팅 비교](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [WebView2Feedback 이슈 #5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [WebView2Feedback 이슈 #4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
