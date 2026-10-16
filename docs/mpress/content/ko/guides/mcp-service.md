---
title: "LLM 제어(MCP)"
description: "Model Context Protocol을 통해 LLM 에이전트가 앱을 테스트하고 제어할 수 있도록 합니다"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="실험적 기능"}
내장 MCP 서버는 실험적 기능이며 향후 릴리스에서 API가 변경될 수 있습니다.

@end

Wails v3에는 내장 [Model Context Protocol](https://modelcontextprotocol.io)(MCP) 서버가 있어 LLM 에이전트(Claude Code, IDE 어시스턴트 또는 모든 MCP 클라이언트)가 **실행 중인** Wails 애플리케이션을 검사하고 테스트하며 제어할 수 있습니다.

프로젝트 수명 주기를 자동화하려면 별도의 `wails3 mcp` CLI 서버를 사용하세요. 이 서버를 사용하면 에이전트가 프로젝트를 검사하고 초기화하며, 진단을 실행하고, 빌드와 개발 작업을 시작하고, 바인딩을 생성하고, 이름이 지정된 Taskfile 태스크를 실행하고, 제한된 범위의 작업 출력을 가져올 수 있습니다. CLI 서버의 범위는 기본적으로 현재 디렉터리로 제한되며 임의의 셸 명령 실행 기능을 노출하지 않습니다. 전송 방식, 인증 및 도구에 대한 자세한 내용은 [CLI MCP 문서](/guides/cli/#mcp)를 참조하세요.

이 기능을 활성화하면 앱에 연결된 에이전트가 다음 작업을 수행할 수 있습니다.

- **창 나열 및 제어** — 크기, 위치, 포커스, 전체 화면, 개발자 도구, 새로고침 등
- **DOM 검사** — 요소 쿼리, HTML 가져오기, 구조 스냅샷 생성
- **JavaScript 평가** — 어떤 창에서든 임의의 코드를 실행하고 결과 가져오기
- **사용자 입력 시뮬레이션** — 에이전트의 작업을 지켜볼 수 있도록 <strong>화면에 표시되는 애니메이션 커서</strong>로 마우스 이동, 클릭, 드래그 및 스크롤 렌더링
- **문자 입력 및 키 누르기** — React 제어 입력에서 작동하는 사실적인 문자별 이벤트
- **바인딩된 Go 메서드 호출** 및 애플리케이션 이벤트 발생/대기

## 작동 방식

MCP 서버는 <strong>`mcp` 빌드 태그</strong>가 있을 때만 애플리케이션에 컴파일됩니다. 이 태그가 없으면 서버 코드가 바이너리에 전혀 포함되지 않으므로 런타임 오버헤드도, 열린 포트도, 공격 표면도 없습니다.

태그가 있으면 서버가 `App.Run()` 내부에서 자동으로 시작되고, 기본적으로 `127.0.0.1:9099`에 바인딩되며 엔드포인트를 로그에 기록합니다. 사용자 코드는 필요하지 않습니다.

## 튜토리얼

### 1단계 — 일반적인 Wails 애플리케이션 작성

MCP에는 import나 등록이 필요하지 않습니다. 평소와 같은 방식으로 앱을 만드세요.

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 2단계 — `mcp` 태그로 빌드 또는 실행

@tabs
[Wails CLI(권장)]
`WAILS_MCP=1`을 설정하면 Wails CLI가 `mcp` 태그를 자동으로 추가합니다.

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[Go 직접 사용]
태그를 `go run` 또는 `go build`에 직접 전달하세요.

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows(PowerShell)]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

애플리케이션은 시작할 때 MCP 엔드포인트를 로그에 기록합니다.

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### 3단계 — 클라이언트 연결

서버는 <strong>MCP 스트리밍 가능 HTTP 전송 방식</strong>을 사용합니다. MCP 호환 클라이언트라면 무엇이든 연결할 수 있습니다.

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

그런 다음 Claude에게 앱과 상호작용하도록 요청하세요.

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code(GitHub Copilot)]
`.vscode/settings.json`에 추가하세요.

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[기타 클라이언트]
스트리밍 가능 HTTP 전송 방식을 지원하는 MCP 클라이언트가 다음 주소를 사용하도록 설정하세요.

```
http://127.0.0.1:9099/mcp
```

@end

### 4단계 — 테스트 세션 실행

에이전트에게 애플리케이션의 기능을 실행해 보도록 요청하세요. 다음은 몇 가지 프롬프트 예시입니다.

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## 구성

모든 구성은 환경 변수를 통해 이루어지므로 코드를 변경할 필요가 없습니다.

| 환경 변수 | 기본값 | 설명 |
| --- | --- | --- |
| `WAILS_MCP` | (설정되지 않음) | Wails CLI 사용 시 `1`, `true`, `on` 또는 `yes`으로 설정하면 `mcp` 빌드 태그가 자동으로 추가됩니다. |
| `WAILS_MCP_HOST` | `127.0.0.1` | 바인딩할 인터페이스입니다. 루프백이 아닌 주소에 바인딩하려면 `WAILS_MCP_TOKEN`이 필요합니다. |
| `WAILS_MCP_TOKEN` | 설정되지 않음 | 루프백에서는 선택 사항인 Bearer 토큰이며, 다른 바인딩 주소에서는 필수입니다. 클라이언트는 `Authorization: Bearer <token>`을 전송합니다. |
| `WAILS_MCP_PORT` | `9099` | 수신 대기할 포트입니다. 임의로 할당되는 사용 가능한 포트를 사용하려면 `0`으로 설정하세요. 할당된 포트는 로그에 출력됩니다. |
| `WAILS_MCP_TIMEOUT` | `30000` | 기본 JS 평가 제한 시간이며 단위는 <strong>밀리초</strong>입니다. |
| `WAILS_MCP_HIDE_CURSOR` | (설정되지 않음) | 애니메이션 커서 오버레이를 비활성화하려면 `1` 또는 `true`로 설정하세요. |

예시 — 사용자 지정 포트와 60초 제한 시간:

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## 사용 가능한 도구

| 도구 | 용도 |
| --- | --- |
| `app_info` | 애플리케이션 정보: 플랫폼, 아키텍처, 모든 창, MCP 엔드포인트 |
| `windows_list` | 모든 창의 위치, 크기 및 상태 나열 |
| `window_control` | 포커스, 크기 조정, 이동, 전체 화면, 개발자 도구, 새로고침, URL 설정 등(22개 작업) |
| `js_eval` | 창에서 JavaScript 평가(비동기 본문, 값에는 `return` 사용) |
| `dom_html` | 페이지 또는 특정 요소의 HTML 가져오기 |
| `dom_query` | CSS 선택자로 요소 찾기 — 태그, 텍스트, 경계, 표시 여부 |
| `screenshot_dom` | 표시된 페이지의 구조적 스냅샷(DOM 기반, 픽셀 정보 없음) |
| `mouse_move` | 커서를 특정 지점 또는 CSS 선택자로 애니메이션 이동 |
| `mouse_click` | 애니메이션 커서로 클릭(왼쪽/오른쪽/가운데 버튼, 두 번 클릭, 보조 키) |
| `mouse_drag` | 애니메이션 커서로 드래그(HTML5 드래그 앤 드롭 요소 지원) |
| `mouse_scroll` | 특정 지점 또는 요소에서 스크롤 |
| `keyboard_type` | 실제와 같은 이벤트를 발생시키며 텍스트를 한 글자씩 입력 |
| `keyboard_press` | 보조 키를 선택적으로 사용하여 단일 키(Enter, Tab, Escape, ArrowDown, …) 누르기 |
| `call_bound_method` | 바인딩된 Go 서비스 메서드 호출(예: `main.GreetService.Greet`) |
| `emit_event` | Wails 애플리케이션 이벤트 발생시키기 |
| `wait_for_event` | Wails 애플리케이션 이벤트를 기다린 후 해당 데이터 반환하기 |

### 다중 창 지원

창에 작용하는 모든 도구는 창의 **이름**(`WebviewWindowOptions.Name`을 통해 설정)을 포함하는 선택적 `window` 인수를 받습니다. 이 인수를 생략하면 현재 포커스된 창을 대상으로 하며, 포커스된 창이 없으면 첫 번째 창을 대상으로 합니다.

```
List all windows, then click the "New" button in the window named "editor".
```

### 요소 선택

마우스 및 키보드 도구에는 다음 중 하나를 지정할 수 있습니다:

- **CSS 선택자** — `selector: "#submit-btn"`(요소가 보이도록 자동으로 스크롤됨)
- **좌표** — `x: 400, y: 300`(뷰포트 기준 CSS 픽셀)

드래그 작업에서는 `from_` 및 `to_` 접두사를 사용합니다:

```
Drag from selector: ".card" to selector: ".dropzone"
```

## 보안

@note{type="caution"}
MCP 서버는 애플리케이션을 프로그래밍 방식으로 완전히 제어할 수 있게 합니다. 서버의 도구에 접근할 수 있는 사람은 누구나 DOM을 읽고, JavaScript를 평가하고, 버튼을 클릭하고, Go 메서드를 호출할 수 있습니다.

@end

- 서버는 기본적으로 `127.0.0.1`에 바인딩됩니다. 브라우저 오리진은 HTTP(S) 루프백 오리진이어야 하며, 불투명한(`null`), 잘못된 형식의 오리진과 외부 오리진은 거부됩니다.
- 헤더가 없는 네이티브 MCP 클라이언트도 계속 지원됩니다. `WAILS_MCP_TOKEN`이 없으면 로컬 프로세스와 허용된 로컬 브라우저 오리진을 신뢰합니다. 오리진 검사는 인증이 아닙니다. 모든 `/mcp` 호출에 Bearer 인증을 요구하려면 엔트로피가 높은 토큰을 설정하세요. 클라이언트의 `Authorization: Bearer <token>` 헤더에도 동일한 토큰을 설정하세요. 프리플라이트 요청에는 토큰이 필요하지 않습니다.
- `/eval-result` 콜백은 클라이언트의 Bearer 토큰 대신 평가마다 예측할 수 없는 ID를 사용하므로 WebView 결과 전달과의 호환성이 유지됩니다.
- 프로덕션 빌드에는 `mcp` 태그를 포함하지 **않는** 것이 좋습니다. Wails CLI는 `WAILS_MCP=1`이 명시적으로 설정된 경우에만 이 태그를 추가하며, 기본 `wails3 build`에는 서버 코드가 전혀 포함되지 않습니다.
- 루프백이 아닌 인터페이스(예: LAN 테스트)에 서버를 노출해야 한다면 `WAILS_MCP_HOST=0.0.0.0`과 엔트로피가 높은 `WAILS_MCP_TOKEN`을 설정하세요. 토큰이 없으면 시작에 실패합니다. 기본 제공 리스너는 HTTP를 사용하므로 신뢰할 수 없는 네트워크에서는 암호화된 터널이나 TLS 종료 프록시를 사용하세요.

## 예제 애플리케이션

모든 도구를 시연하는 완전한 플레이그라운드 애플리케이션은 [`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp)에서 사용할 수 있습니다. 다음 항목이 포함되어 있습니다:

- 증가/초기화 버튼이 있는 카운터
- Greet, Add, Shout 바인딩 메서드를 사용하는 이름 입력란
- HTML5 드래그 앤 드롭 소스와 대상
- 스크롤 가능한 목록(항목 50개)
- 이벤트 로그

다음 명령으로 실행하세요:

```shell
cd v3/examples/mcp
go run -tags mcp .
```

그런 다음 Claude Code 또는 임의의 MCP 클라이언트를 `http://127.0.0.1:9099/mcp`에 연결하고 UI를 조작해 보도록 요청하세요.

## 피드백

기본 제공 MCP 서버는 실험적 기능이며, 여러분의 피드백에 따라 향후 방향이 결정됩니다. 사용해 보셨다면 어떤 클라이언트와 도구를 사용했는지, 무엇을 기대했고 실제로는 어떤 결과가 발생했는지, 에이전트가 애플리케이션을 조작하도록 하는 것이 유용했는지 알려주세요. 실행한 내용을 정확히 설명한 보고서가 가장 유용합니다. [MCP 서버 피드백 토론](https://github.com/wailsapp/wails/discussions/5692)에서 의견을 들려주세요.
