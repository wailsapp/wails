---
title: "바인딩 시스템"
description: "Wails v3에서 상용구 코드 없이 Go와 JavaScript가 서로를 호출하는 방식"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> "바인딩"은 다음과 같은 코드를 작성할 수 있게 해 주는 <strong>타입 안전 계약</strong>입니다.

```go
msg, err := chatService.Send("Hello")
```

Go에서는 *그리고*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

TypeScript에서는 **IPC 연결 코드를 직접 작성하지 않고도** 가능합니다. 이 문서에서는 빌드 시의 <strong>정적 분석</strong>부터 <strong>코드 생성</strong>을 거쳐 WebView를 통해 바이트를 전달하는 <strong>런타임 브리지</strong>에 이르기까지, 이 과정이 *어떻게* 이루어지는지 설명합니다.

> 다음에 관한 자세하고 공식적인 설명은 [`contributing/architecture/bindings`](/contributing/architecture/bindings/)에서 확인하세요.
>
> 생성기 파이프라인을 심층적으로 다루며, 이 페이지에서는
>
> 기여자에게 초점을 맞춰 개괄적으로 설명합니다.

---

## 1. 30초 개요

| 단계 | 컴포넌트 | 출력 |
| --- | --- | --- |
| **수집/분석** | `internal/generator/collect/`, `internal/generator/analyse.go` | 내보낸 Go 서비스, 메서드, 매개변수, 반환 타입 및 모델의 인메모리 모델 |
| **생성** | `internal/generator/render/templates/*.tmpl`(`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | `frontend/bindings/<full Go import path>/...` 아래의 서비스별 ES 모듈 |
| **런타임** | `pkg/application/messageprocessor*.go` + `internal/runtime/desktop/@wailsio/runtime/src/` 아래의 내장 JS 런타임(`calls.ts`, `events.ts` 등) | WebView의 네이티브 브리지를 통한 호출/이벤트 메시지 |

이 흐름은 `wails3 generate bindings` 명령이 조정하며, 이 명령은 Go 패키지 집합을 대상으로 `generator.Generate`(`internal/generator/generate.go`에 정의됨)을 실행합니다.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. 정적 분석

### 진입점

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

수집기 패스는 로드된 모든 패키지를 순회하며 다음을 기록합니다.

- `collect.ServiceInfo` — 내보내고 바인딩한 Go 구조체마다 하나씩 기록합니다.
- `collect.ServiceMethodInfo`/`collect.MethodInfo` — 메서드별 시그니처 정보(이름, 매개변수, 결과, 오류 위치, 리시버, 문서)입니다.
- `collect.ModelInfo`/`collect.StructInfo` — TS/JS 모델로 출력합니다.
- `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore`, `//wails:id <hex>` 같은 지시문 주석입니다(`internal/generator/collect/directive.go` 참조).

지원하지 않는 타입은 생성기 오류를 발생시키므로, 실수가 런타임이 아닌 빌드 시점에 드러납니다.

### 모델 식별자

런타임 호출 엔벌로프는 메서드의 정규화된 이름(`pkg.Struct.Method`)에 대한 <strong>결정적 FNV-1a 해시</strong>로 메서드를 식별합니다. 생성된 바인딩에서는 `$Call.ByID(<numeric-id>, …)`로 표시되며, `-names` 옵션으로 생성을 실행하면 `$Call.ByName("pkg.Struct.Method", …)`로 표시됩니다.

---

## 3. 코드 생성

### 템플릿

`internal/generator/render/templates/`:

| 템플릿 | 용도 |
| --- | --- |
| `service.js.tmpl` | 바인딩된 서비스별 JS 모듈 하나 |
| `service.ts.tmpl` | TypeScript 컴패니언(`-ts` 포함) |
| `models.js.tmpl` | 모델 클래스 출력(패키지별) |
| `models.ts.tmpl` | 모델 `.d.ts` 출력(패키지별) |
| `index.tmpl` | 패키지별 `index.{js,ts}` 배럴 재내보내기 |
| `eventcreate.js.tmpl`/`eventdata.d.ts.tmpl` | 이벤트 생성자/페이로드 타입 정의 |
| `newline.tmpl` | 후행 줄 바꿈 정규화기 |

출력은 `frontend/bindings/<full Go import path>/...` 아래에 생성됩니다. 예를 들어 `github.com/you/yourapp/services/chat`에 정의된 서비스는 `frontend/bindings/github.com/you/yourapp/services/chat/` 아래에 생성됩니다. v3에는 `frontend/src/wailsjs/` 디렉터리가 없습니다.

### JavaScript 출력

생성된 바인딩은 `/wails/runtime.js`에서 런타임 헬퍼를 가져오는 ES 모듈입니다.

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

`-names` 옵션으로 생성을 실행하면 대신 `$Call.ByName("pkg.Struct.Method", ...)`이 출력됩니다. 이는 항상 <strong>정규화된 전체 이름</strong>이며, `"Method"`만 출력되는 경우는 없습니다.

생성된 모델 클래스는 필드별 `if (!("X" in $$source))` 기본값 및 따옴표로 묶은 필드 이름과 함께 `$$source` 생성자 패턴을 사용하며, 문자열 입력에 `JSON.parse`을 실행하는 `static createFrom(...)`도 사용합니다.

### 주요 타입 매핑

`internal/generator/render/`을 기준으로 검증했습니다.

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V`(문자열이 아닌 `K`) | `{ [_ in K]?: V }`(`Map<K, V>`도 `Record<K, V>`도 아님) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string`(JSON ISO 8601) |
| `error`(반환 위치) | 거부된 프로미스 |

### 리플렉션 참고 사항

`pkg/application/bindings.go`은 **수작업으로 작성되었으며**, `BoundMethod` 레지스트리에서 메서드 디스패치를 구동하는 데 `reflect`을 사용합니다. 이전에 언급된 "런타임 리플렉션 없음"이라는 설명을 너무 문자 그대로 받아들이지 마세요. 생성기는 리플렉션을 사용하지 않지만 런타임 디스패처는 사용합니다.

---

## 4. 런타임 호출 프로토콜

### JavaScript 측

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

런타임 헬퍼는 `internal/runtime/desktop/@wailsio/runtime/src/calls.ts`(호출 디스패치), `events.ts`(이벤트) 등의 파일에 있습니다. 이 트리에는 `invoke.ts`이나 `errors.ts`이 없습니다. 정확한 전송 엔벌로프는 JS 측의 `calls.ts`에서 인코딩되고 Go 측의 `pkg/application/messageprocessor_call.go`에서 디코딩됩니다. 브리지를 디버깅할 때는 이 두 파일을 함께 살펴보세요.

### Go 측

1. `pkg/application/messageprocessor_call.go`이 호출 메시지를 수신합니다.
2. `pkg/application/bindings.go`에서 ID나 이름으로 바인딩된 메서드를 조회합니다(`reflect`에서 구동).
3. 바인딩된 메서드를 호출하고 `{result, error}`을 직렬화하여 JS로 돌려보냅니다.

### 오류 매핑

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise`이 결과와 함께 이행됨 |
| `error != nil` | `Promise`이 거부되며, 이때 `Error`의 `message`은 Go 오류 문자열임 |

---

## 5. Go에서 JavaScript 호출하기

바인딩 생성기는 단방향입니다(Go 메서드를 JS에 노출). Go → JS 통신에는 이벤트 버스를 사용하거나 창에서 JS를 실행하세요.

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

JS 측에서는 `/wails/runtime.js`의 `Events.On(name, cb)`을 사용하여 구독하세요.

---

## 6. 확장 및 문제 해결

### 지원되지 않는 타입 오류

```
error: field "Client" uses unsupported type: chan struct{}
```

→ 채널을 메서드 API 뒤에 래핑하거나, 생성기가 해당 필드를 건너뛰도록 필드에 `//wails:internal`을 표시하세요.

### 오래된 바인딩

생성된 출력은 `wails3 generate bindings` / `wails3 dev` / `wails3 build`을 실행할 때마다 덮어씁니다. IDE IntelliSense에 오래된 스텁이 표시되면 `frontend/bindings/`을 삭제하고 생성기를 다시 실행하세요. `-clean` 플래그(현재 빌드의 기본값은 `true`)는 실행할 때마다 먼저 바인딩 디렉터리를 비웁니다.

### 성능 팁

- 브리지를 통해 큰 바이트 슬라이스를 스트리밍하지 말고, 대신 애셋 서버를 통해 제공하세요.
- 지연 시간이 중요하다면 빠른 호출 여러 개를 하나의 메서드로 묶으세요.
- 할당을 줄이려면 작은 매개변수 구조체에 값 리시버를 우선 사용하세요.

---

## 7. 주요 파일 맵

| 관심 영역 | 파일 |
| --- | --- |
| 생성기 오케스트레이션 | `internal/generator/generate.go` |
| 의미 체계 검사 | `internal/generator/analyse.go` |
| 수집(서비스, 메서드, 모델) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| 렌더링 템플릿 | `internal/generator/render/templates/*.tmpl` |
| 생성된 바인딩 위치 | `frontend/bindings/<full Go import path>/...` |
| Go 측 디스패처 | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| JS 런타임 | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

브리지 버그를 추적할 때 이 요약표를 가까이 두고 활용하세요.

---

## 8. 요약

1. <strong>수집기</strong>가 Go 코드를 스캔하여 인메모리 의미 모델을 만듭니다.
2. <strong>템플릿</strong>이 서비스별 ES 모듈과 패키지별 모델/인덱스 파일을 생성합니다.
3. <strong>메시지 프로세서</strong>가 바인딩 레지스트리를 통해 Go 측에서 호출을 디스패치합니다.
4. <strong>JS 런타임</strong>이 이 모든 기능을 취소를 지원하는 관용적인 프로미스로 래핑합니다.

IPC 상용구 코드를 단 한 줄도 작성할 필요가 없습니다. 이것이 Wails v3 바인딩 시스템입니다. 이제 마음껏 바인딩하세요!
