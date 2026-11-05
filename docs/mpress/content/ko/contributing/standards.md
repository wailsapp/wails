---
title: "코딩 표준"
description: "Wails v3의 코드 스타일, 규칙 및 모범 사례"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## 코드 스타일 및 규칙

일관된 코딩 표준을 따르면 코드베이스를 더 쉽게 읽고 유지보수하며 기여할 수 있습니다.

## Go 코드 표준

### 코드 서식

표준 Go 서식 도구를 사용하세요.

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

**필수:** 모든 Go 코드는 커밋하기 전에 `gofmt` 및 `goimports` 검사를 통과해야 합니다.

### 명명 규칙

**패키지:**

- 소문자를 사용하고 가능하면 한 단어로 작성하세요.
- `package application`, `package events`
- 밑줄이나 대소문자 혼용을 피하세요.

**내보내는 이름:**

- 타입, 함수, 상수에는 PascalCase를 사용하세요.
- `type WebviewWindow struct`, `func NewApplication()`

**내보내지 않는 이름:**

- 내부 타입, 함수, 변수에는 camelCase를 사용하세요.
- `type windowImpl struct`, `func createWindow()`

**인터페이스:**

- 동작에 따라 이름을 지정하세요: `Reader`, `Writer`, `Handler`
- 메서드가 하나뿐인 인터페이스에는 이름에 `-er` 접미사를 사용하세요.

```go
// Good
type Closer interface {
    Close() error
}

// Avoid
type CloseInterface interface {
    Close() error
}
```

### 오류 처리

**항상 오류를 확인하세요:**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**오류 래핑을 사용하세요:**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**필요한 경우 사용자 정의 오류 타입을 만드세요:**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### 주석 및 문서화

**패키지 주석:**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**내보내는 선언:**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**구현 주석:**

```go
// processEvent handles incoming events from the runtime.
// It dispatches to registered handlers and manages event lifecycle.
func (a *Application) processEvent(event *Event) {
    // Validate event before processing
    if event == nil {
        return
    }

    // Find and invoke handlers
    // ...
}
```

### 함수 및 메서드 구조

**함수가 한 가지 작업에 집중하도록 하세요:**

```go
// Good - single responsibility
func (w *Window) setTitle(title string) {
    w.title = title
    w.updateNativeTitle()
}

// Bad - doing too much
func (w *Window) updateEverything() {
    w.setTitle(w.title)
    w.setSize(w.width, w.height)
    w.setPosition(w.x, w.y)
    // ... 20 more operations
}
```

**조기 반환을 사용하세요:**

```go
// Good
func validate(input string) error {
    if input == "" {
        return errors.New("empty input")
    }

    if len(input) > 100 {
        return errors.New("input too long")
    }

    return nil
}

// Avoid deep nesting
```

### 동시성

**취소 처리에는 context를 사용하세요:**

```go
func (a *Application) RunWithContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-a.done:
        return nil
    }
}
```

**공유 상태는 뮤텍스로 보호하세요:**

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**고루틴 누수를 방지하세요:**

```go
// Good - goroutine has exit condition
func (a *Application) startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            case work := <-a.workChan:
                a.process(work)
            }
        }
    }()
}
```

### 테스트

**테스트 파일 명명 규칙:**

```go
// Implementation: window.go
// Tests: window_test.go
```

**테이블 기반 테스트:**

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty input", "", true},
        {"valid input", "hello", false},
        {"too long", strings.Repeat("a", 101), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## JavaScript/TypeScript 표준

### 코드 서식

일관된 서식을 위해 Prettier를 사용하세요.

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### 명명 규칙

**변수 및 함수:**

- camelCase: `const userName = "John"`

**클래스 및 타입:**

- PascalCase: `class WindowManager`

**상수:**

- UPPER<em>SNAKE</em>CASE: `const MAX_RETRIES = 3`

### TypeScript

**명시적 타입을 사용하세요:**

```typescript
// Good
function greet(name: string): string {
    return `Hello, ${name}`
}

// Avoid implicit any
function process(data) {  // Bad
    return data
}
```

**인터페이스를 정의하세요:**

```typescript
interface WindowOptions {
    title: string
    width: number
    height: number
}

function createWindow(options: WindowOptions): void {
    // ...
}
```

## 커밋 메시지 형식

[Conventional Commits](https://www.conventionalcommits.org/) 형식을 사용하세요.

```
<type>(<scope>): <subject>

<body>

<footer>
```

**유형:**

- `feat`: 새 기능
- `fix`: 버그 수정
- `docs`: 문서 변경
- `refactor`: 코드 리팩터링
- `test`: 테스트 추가 또는 업데이트
- `chore`: 유지보수 작업

**예시:**

```
feat(window): add SetAlwaysOnTop method

Implement SetAlwaysOnTop for keeping windows above others.
Adds platform implementations for macOS, Windows, and Linux.

Closes #123
```

```
fix(events): prevent event handler memory leak

Event listeners were not being properly cleaned up when
windows were closed. This adds explicit cleanup in the
window destructor.
```

## 풀 리퀘스트 지침

### 제출 전 확인 사항

- [ ] 코드가 `gofmt` 및 `goimports` 검사를 통과함
- [ ] 모든 테스트를 통과함(`go test ./...`)
- [ ] 새 코드에 대한 테스트가 있음
- [ ] 필요한 경우 문서를 업데이트함
- [ ] 커밋 메시지가 규칙을 따름
- [ ] `master`과의 병합 충돌이 없음

### PR 설명 템플릿

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## 코드 리뷰 절차

### 리뷰어 지침

- 건설적이고 존중하는 태도를 유지하세요
- 개인적인 선호가 아니라 코드 품질에 집중하세요
- 변경을 제안하는 이유를 설명하세요
- 변경 사항이 만족스러우면 승인하세요

### 작성자 지침

- 모든 의견에 답변하세요
- 필요하면 명확한 설명을 요청하세요
- 요청된 내용을 변경하거나 변경하지 않는 이유를 설명하세요
- 피드백을 열린 자세로 받아들이세요

## 모범 사례

### 성능

- 성급한 최적화를 피하세요
- 최적화 전에 프로파일링하세요
- 성능이 중요한 코드에는 벤치마크를 사용하세요

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### 보안

- 모든 사용자 입력을 검증하세요
- 데이터를 표시하기 전에 정제하세요
- 무작위 데이터에는 `crypto/rand`을 사용하세요
- 민감한 정보를 절대 로그에 기록하지 마세요

### 문서화

- 외부에 공개된 API를 문서화하세요
- 문서에 예제를 포함하세요
- API를 변경할 때 문서도 업데이트하세요
- README 파일을 최신 상태로 유지하세요

## 플랫폼별 코드

### 파일 이름 지정

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### 빌드 태그

```go
//go:build darwin

package application

// macOS-specific code
```

## 린팅

커밋하기 전에 린터를 실행하세요:

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## 질문이 있나요?

표준이 확실하지 않은 경우:

- 기존 코드에서 예제를 확인하세요
- [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하세요
- GitHub에서 토론을 시작하세요
