---
title: "Go-프런트엔드 브리지"
description: "Wails가 Go와 JavaScript 간의 직접 통신을 지원하는 방식을 자세히 알아봅니다"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Go-JavaScript 직접 통신

Wails는 Go와 JavaScript 사이에 <strong>직접 연결되는 인메모리 브리지</strong>를 제공하여 HTTP 오버헤드, 프로세스 경계 또는 직렬화 병목 현상 없이 원활하게 통신할 수 있도록 합니다.

## 전체 구조

```d2
direction: right

Frontend: 프런트엔드(JavaScript) {
  UI: React/Vue/Vanilla {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: 자동 생성된 바인딩 {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Wails 브리지 {
  Encoder: JSON 인코더 {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: 메서드 라우터 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 디코더 {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: 타입 생성기 {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: 백엔드(Go) {
  Services: 사용자 서비스 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: 서비스 레지스트리 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "Method('arg') 호출"
Bridge.Encoder -> Bridge.Router: JSON으로 인코딩
Bridge.Router -> Backend.Registry: 서비스 찾기
Backend.Registry -> Backend.Services: 메서드 호출
Backend.Services -> Bridge.Decoder: 결과 반환
Bridge.Decoder -> Frontend.Bindings: JS로 디코딩
Frontend.Bindings -> Frontend.UI: Promise 이행
Bridge.TypeGen -> Frontend.Bindings: 타입 생성
```

**핵심:** HTTP도, IPC도, 프로세스 경계도 없습니다. <strong>타입 안전성</strong>이 보장되는 <strong>직접 함수 호출</strong>만 사용합니다.

## 작동 방식: 단계별 설명

### 1. 서비스 등록(시작 시)

애플리케이션이 시작되면 Wails가 서비스를 검사합니다.

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) Add(a, b int) int {
    return a + b
}

// Register service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Wails가 수행하는 작업:**

1. 내보낸 메서드를 찾기 위해 **구조체를 검사합니다**
2. **타입 정보를 추출합니다**(매개변수, 반환 타입)
3. 메서드 이름을 함수에 매핑하는 **레지스트리를 구축합니다**
4. 전체 타입 정의가 포함된 **TypeScript 바인딩을 생성합니다**

### 2. 바인딩 생성(빌드 시)

Wails는 TypeScript 바인딩을 자동으로 생성합니다.

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**타입 매핑:**

| Go 타입 | TypeScript 타입 |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | 예외(throw됨) |

### 3. 프런트엔드 호출(런타임)

개발자가 JavaScript에서 Go 메서드를 호출합니다.

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**처리 과정:**

1. **바인딩 함수 호출** - `Greet("World")`
2. **메시지 생성** - `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **브리지로 전송** - WebView의 JavaScript 브리지를 통해 전송
4. **Promise 반환** - 응답 대기

### 4. 브리지 처리(런타임)

브리지가 메시지를 수신하여 처리합니다:

```d2
direction: down

Receive: 메시지 수신 {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: JSON 파싱 {
  shape: rectangle
}

Validate: 검증 {
  Check: 서비스가 존재하는가? {
    shape: diamond
  }

  CheckMethod: 메서드가 존재하는가? {
    shape: diamond
  }

  CheckTypes: 타입이 올바른가? {
    shape: diamond
  }
}

Invoke: Go 메서드 호출 {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: 결과 인코딩 {
  shape: rectangle
}

Send: 응답 전송 {
  shape: rectangle
  style.fill: "#10B981"
}

Error: 오류 전송 {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: 예
Validate.Check -> Error: 아니요
Validate.CheckMethod -> Validate.CheckTypes: 예
Validate.CheckMethod -> Error: 아니요
Validate.CheckTypes -> Invoke: 예
Validate.CheckTypes -> Error: 아니요
Invoke -> Encode: 성공
Invoke -> Error: 오류
Encode -> Send
```

**보안:** 등록된 서비스와 외부에 공개된 메서드만 호출할 수 있습니다.

### 5. Go 실행(런타임)

Go 메서드가 실행됩니다:

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**실행 컨텍스트:**

- <strong>goroutine</strong>에서 실행됩니다(비차단 방식).
- **모든 Go 기능**(파일 시스템, 네트워크, 데이터베이스)을 사용할 수 있습니다.
- <strong>다른 Go 코드</strong>를 자유롭게 호출할 수 있습니다.
- 결과 또는 오류를 반환합니다.

### 6. 응답(런타임)

결과가 JavaScript로 다시 전송됩니다:

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**오류 처리:**

```go
func (g *GreetService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```javascript
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Go error:", error)  // "division by zero"
}
```

## 성능 특성

### 속도

**일반적인 호출 오버헤드:** &lt;1ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**다른 방식과의 비교:**

- **HTTP/REST:** 5-50ms(네트워크 스택, 직렬화)
- **IPC:** 1-10ms(프로세스 경계, 마샬링)
- **Wails Bridge:** &lt;1ms(메모리 내 직접 호출)

### 메모리

**호출당 오버헤드:** ~1KB(메시지 버퍼)

**제로 카피 최적화:** 가능한 경우 대용량 데이터(>1MB)에 공유 메모리를 사용합니다.

### 동시성

**호출은 동시에 실행됩니다:**

- 각 호출은 자체 goroutine에서 실행됩니다.
- 여러 호출을 동시에 실행할 수 있습니다.
- 호출 간에 블로킹이 발생하지 않습니다.

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## 타입 시스템

### 지원되는 타입

#### 기본 타입

```go
// Go
func Example(
    s string,
    i int,
    f float64,
    b bool,
) (string, int, float64, bool) {
    return s, i, f, b
}
```

```typescript
// TypeScript (auto-generated)
function Example(
    s: string,
    i: number,
    f: number,
    b: boolean,
): Promise<[string, number, number, boolean]>
```

#### 슬라이스와 배열

```go
// Go
func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```typescript
// TypeScript
function Sum(numbers: number[]): Promise<number>

// Usage
const total = await Sum([1, 2, 3, 4, 5])  // 15
```

#### 맵

```go
// Go
func GetConfig() map[string]interface{} {
    return map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
        "enabled": true,
    }
}
```

```typescript
// TypeScript
function GetConfig(): Promise<Record<string, any>>

// Usage
const config = await GetConfig()
console.log(config.theme)  // "dark"
```

#### 구조체

```go
// Go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(id int) (*User, error) {
    return &User{
        ID:    id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}
```

```typescript
// TypeScript (auto-generated)
interface User {
    id: number
    name: string
    email: string
}

function GetUser(id: number): Promise<User>

// Usage
const user = await GetUser(1)
console.log(user.name)  // "Alice"
```

**JSON 태그:** TypeScript의 필드 이름을 제어하려면 `json:` 태그를 사용하세요.

#### 시간

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<Date>

// Usage
const timestamp = await GetTimestamp()
console.log(timestamp.toISOString())
```

#### 오류

```go
// Go
func Validate(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}
```

```typescript
// TypeScript
function Validate(input: string): Promise<void>

// Usage
try {
    await Validate("")
} catch (error) {
    console.error(error)  // "input cannot be empty"
}
```

### 지원되지 않는 타입

다음 타입은 브리지를 통해 전달할 수 **없습니다**:

- **채널**(`chan T`)
- **함수**(`func()`)
- **인터페이스**(`interface{}` / `any` 제외)
- **포인터**(구조체를 가리키는 포인터 제외)
- **내보내지 않은 필드**(소문자로 시작)

**해결 방법:** ID 또는 핸들을 사용하세요.

```go
// ❌ Can't pass file handle
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

// ✅ Return file ID instead
var files = make(map[string]*os.File)

func OpenFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    id := generateID()
    files[id] = file
    return id, nil
}

func ReadFile(id string) ([]byte, error) {
    file := files[id]
    return io.ReadAll(file)
}

func CloseFile(id string) error {
    file := files[id]
    delete(files, id)
    return file.Close()
}
```

## 고급 패턴

### 컨텍스트 전달

서비스에서 호출 컨텍스트에 접근할 수 있습니다.

```go
type UserService struct{}

func (s *UserService) GetCurrentUser(ctx context.Context) (*User, error) {
    // Access the calling window via the context value
    window, _ := ctx.Value(application.WindowKey).(application.Window)
    _ = window

    // Access the application
    app := application.Get()
    _ = app

    // Your logic
    return getCurrentUser(), nil
}
```

**컨텍스트에서 제공하는 정보:**

- 호출한 창
- 애플리케이션 인스턴스
- 요청 메타데이터

### 데이터 스트리밍

대용량 데이터에는 반환 값 대신 이벤트를 사용하세요.

```go
func ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        // Emit progress events
        app.Event.Emit("file-progress", map[string]interface{}{
            "line": lineNum,
            "text": scanner.Text(),
        })
    }

    return scanner.Err()
}
```

```javascript
import { Events } from '@wailsio/runtime'
import { ProcessLargeFile } from './bindings/FileService'

// Listen for progress
Events.On('file-progress', (data) => {
    console.log(`Line ${data.line}: ${data.text}`)
})

// Start processing
await ProcessLargeFile('/path/to/large/file.txt')
```

### 취소

취소 가능한 작업에는 컨텍스트를 사용하세요.

```go
func LongRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check if cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

**참고:** 프런트엔드 연결이 끊어지면 컨텍스트가 자동으로 취소됩니다.

### 일괄 작업

작업을 일괄 처리하여 브리지 오버헤드를 줄이세요.

```go
// ❌ Inefficient: N bridge calls
for _, item := range items {
    await ProcessItem(item)
}

// ✅ Efficient: 1 bridge call
await ProcessItems(items)
```

```go
func ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = processItem(item)
    }
    return results, nil
}
```

## 브리지 디버깅

### 디버그 로깅 활성화

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**출력에 표시되는 정보:**

- 메서드 호출
- 매개변수
- 반환 값
- 오류
- 실행 시간 정보

### 생성된 바인딩 확인

생성된 TypeScript를 확인하려면 `frontend/bindings/`을 살펴보세요.

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### 서비스 직접 테스트

프런트엔드 없이 Go 서비스를 테스트하세요.

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## 성능 최적화 팁

### ✅ 권장 사항

- **작업 일괄 처리** - 브리지 호출 횟수를 줄이세요.
- **스트리밍에 이벤트 사용** - 대용량 배열을 반환하지 마세요.
- **메서드를 빠르게 유지** - &lt;100ms가 이상적입니다.
- **goroutine 사용** - 오래 걸리는 작업에 사용하세요.
- **Go 측에서 캐싱** - 반복 계산을 피하세요.

### ❌ 금지 사항

- **지나치게 많이 호출하지 않기** - 가능하면 일괄 처리하세요.
- **방대한 데이터를 반환하지 않기** - 페이지네이션 또는 스트리밍을 사용하세요.
- **블로킹하지 않기** - 오래 걸리는 작업에는 goroutine을 사용하세요.
- **복잡한 타입을 전달하지 않기** - 단순하게 유지하세요.
- **오류를 무시하지 않기** - 항상 처리하세요.

## 보안

브리지는 기본적으로 안전합니다.

1. **허용 목록만 사용** - 등록된 서비스만 호출할 수 있습니다.
2. **타입 검증** - 인수를 Go 타입과 대조하여 검사합니다.
3. **eval() 미사용** - 프런트엔드에서 임의의 Go 코드를 실행할 수 없습니다.
4. **리플렉션 악용 방지** - 내보낸 메서드에만 접근할 수 있습니다.

**권장 사례:**

- Go에서 **입력을 검증하세요**(프런트엔드를 신뢰하지 마세요).
- 인증/권한 부여에 **컨텍스트를 사용하세요**.
- 비용이 큰 작업은 **호출 빈도를 제한하세요**.
- 파일 경로와 사용자 입력을 **정제하세요**.

## 다음 단계

**빌드 시스템** - Wails가 애플리케이션을 빌드하고 번들링하는 방법을 알아보세요. [자세히 알아보기 →](/concepts/build-system/)

**서비스** - 서비스 시스템을 자세히 살펴보세요. [자세히 알아보기 →](/features/bindings/services/)

**이벤트** - 게시/구독 통신에 이벤트 사용 [자세히 알아보기 →](/features/events/system/)

---

**브리지에 대해 궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [바인딩 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)를 확인하세요.
