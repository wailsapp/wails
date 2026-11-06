---
title: "노치 창"
description: "카메라 하우징에 부착되는 네이티브 macOS 창을 만듭니다."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow`은 카메라 하우징에 부착되는, 활성화되지 않는 형태가 지정된 macOS 패널을 만듭니다. Wails는 네이티브 배치, 검은색 연결 날개, 투명한 외부 캔버스, 창 레벨, Spaces 동작, 선택적 표시 및 숨기기 애니메이션을 관리합니다. 웹 콘텐츠는 요청한 내부 사각형 영역만 차지합니다. 포인터를 창 안으로 이동하면 애플리케이션을 활성화하지 않고도 해당 웹뷰가 즉시 키 입력을 받는 상태가 됩니다.

@note{type="caution" title="웹뷰 투명도는 비공개 API를 사용합니다"}
투명한 웹 콘텐츠를 통해 네이티브 노치 형태가 드러나게 하려면 `NewNotchWindow`에 `-tags private_mac_apis`이 필요합니다. 이 태그가 없어도 패널, 배치, 애니메이션은 계속 작동하지만 웹뷰는 불투명하게 유지됩니다. [비공개 macOS API](/guides/build/private-macos-apis/#webview-transparency-and-background)를 참조하세요.

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="MacBook 카메라 하우징에서 아래로 미끄러져 내려와 실시간 시스템 지표를 표시한 뒤 다시 노치 아래로 숨는 Wails 노치 알림"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    영구 웹뷰와 네이티브 표시 및 숨기기 전환 효과를 사용하는
    애니메이션 노치 알림입니다.
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## 옵션

| 필드 | 타입 | 기본값 | 설명 |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | 네이티브 형태의 가장자리 안쪽에서 사용할 수 있는 웹뷰 너비입니다. |
| `Height` | `int` | `92` | 네이티브 형태의 가장자리 안쪽에서 사용할 수 있는 웹뷰 높이입니다. |
| `Animated` | `bool` | `false` | `Show` 시 창을 아래로, `Hide` 시 위로 밀어 이동시킵니다. |
| `AnimationSpeed` | `time.Duration` | `420ms` | 표시 지속 시간입니다. 숨기기에는 이 값의 3분의 2를 사용하며, 기본값은 `280ms`입니다. |
| `Screen` | `*Screen` | 주 디스플레이 | 특정 디스플레이를 대상으로 지정합니다. 지정하지 않으면 Wails는 주 디스플레이를 사용합니다. |
| `WindowOptions` | `WebviewWindowOptions` | 기본값 | 이름, URL 또는 HTML, CSS, JavaScript, 키 바인딩 및 기타 웹뷰 동작을 제공합니다. |

`NewNotchWindow`은 외부 크기, 위치, 프레임, 투명도, 크기 조정 정책, 네이티브 패널 클래스, 창 레벨 및 컬렉션 동작을 관리합니다. `WindowOptions`에서 해당 필드에 지정한 값은 의도적으로 대체됩니다. 다른 필드는 유지됩니다. 창이 카메라 하우징에 계속 부착되도록 네이티브 배경 드래그와 CSS 드래그 영역은 비활성화됩니다.

반환되는 `NotchWindow`은 의도적으로 `Show`, `Hide`, `Visibility` 및 `Close`만 노출합니다. 상위 수준 핸들을 통해서는 네이티브 지오메트리를 변경할 수 없습니다.

@note{type="note"}
노치 창에는 macOS가 필요합니다. 카메라 하우징이 없는 Mac에서는 Wails가 창을 메뉴 막대 아래의 상단 중앙에 배치합니다. 지원되지 않는 플랫폼에서 `NewNotchWindow`은 수명 주기 메서드가 안전하게 아무 작업도 하지 않고 `Visibility`이 항상 false인 비활성 핸들을 반환합니다.

@end

## 수명 주기

- `Show`은 기존 네이티브 창을 표시합니다. 애니메이션을 활성화하면 디스플레이 위쪽에서 아래로 미끄러져 내려옵니다.
- `Hide`은 창을 재사용할 수 있도록 유지합니다. 애니메이션을 활성화하면 창을 표시 순서에서 제외하기 전에 디스플레이 위쪽으로 다시 밀어 올립니다. 웹뷰, JavaScript 상태, 바인딩 및 이벤트 리스너는 계속 로드된 상태로 유지되지만 숨겨진 창에는 호버 대상이 없습니다. 다시 표시하려면 애플리케이션에서 `Show`을 호출해야 합니다.
- `Visibility`은 현재 네이티브 표시 상태를 보고합니다.
- `Close`은 네이티브 창을 영구적으로 제거합니다. 해당 알림을 다시 표시하려면 먼저 새 창을 만드세요.
- 포인터가 들어오면 해당 노치 창을 맨 앞으로 가져오고 웹뷰에 포커스를 주어 즉시 키보드로 조작할 수 있게 하면서, 패널이 애플리케이션을 활성화하지 않는 동작은 그대로 유지합니다.

`NewNotchWindow`을 호출할 때마다 자체 콘텐츠, 표시 상태 및 애니메이션 상태를 가진 독립적인 창이 생성됩니다. 창들은 동일한 네이티브 레벨을 사용하므로 가장 최근에 표시되었거나 포인터가 들어온 인스턴스가 맨 앞에 나타나 이전 인스턴스를 덮을 수 있습니다. 창들이 서로 다른 화면을 대상으로 하면 각 창은 해당 화면의 카메라 하우징 중앙에 배치됩니다. macOS는 서로 다른 애플리케이션에 속한 노치 창들을 조정하지 않습니다. 별도의 앱들이 동일한 위치와 레벨에 창을 표시하면 가장 최근에 표시 순서가 지정된 창이 맨 앞에 나타납니다.

알림 용도에서는 일반적으로 애플리케이션이 숨겨진 창을 재사용하거나 하나의 웹뷰에서 자체 대기열을 관리합니다. Wails는 대기열 처리, 교체, 자동 닫기 또는 단일 창 정책을 강제하지 않습니다.

@note{type="note"}
호버 시 다시 열리는 영구적인 축소 표면은 의도적으로 알림 숨기기와 별도로 취급되며 [#6009](https://github.com/wailsapp/wails/issues/6009)에서 추적됩니다.

@end

알림을 반복해서 표시하고 숨겨도 실시간 JavaScript 상태가 유지되는 소형 시스템 모니터는 [notch-notification 예제](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notch-notification)를 참조하세요.
