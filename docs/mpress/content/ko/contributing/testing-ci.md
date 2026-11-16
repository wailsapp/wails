---
title: "테스트 및 지속적 통합"
description: "Wails v3에서 단위 테스트, 통합 테스트 모음, 경쟁 상태 감지 및 GitHub Actions CI를 통해 품질을 보장하는 방법을 설명합니다."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

견고한 데스크톱 프레임워크에는 빈틈없는 테스트가 필요합니다. Wails v3에서는 다음과 같은 **계층형 전략**을 사용합니다.

| 계층 | 목표 | 도구 |
| --- | --- | --- |
| 단위 테스트 | 격리된 함수에 대한 빠른 피드백 | `go test ./...` |
| 생성기/CLI 테스트 | `wails3 generate bindings` 및 CLI 연결부 검증 | `task test:generator`, `task test:cli` |
| 템플릿 테스트 | 배포되는 모든 템플릿이 계속 빌드되는지 확인 | `task test:templates` |
| 경쟁 상태 감지 | 런타임 및 브리지의 데이터 경쟁 감지 | `go test -race ./...` |
| CI 매트릭스 | 모든 PR에서 운영체제 간 안정성 확인 | GitHub Actions |

이 문서에서는 **테스트가 있는 위치**, **테스트 실행 방법**, **Taskfile이 조정하는 작업**을 설명합니다.

`v3/TESTING.md`에서는 크로스 플랫폼 예제 빌드를 설명하며, 경쟁 상태 억제 목록은 제공하지 않습니다.

---

## 1. 디렉터리 규칙

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

지침:

- **단위 테스트는 코드 옆에 두세요**(`foo.go` ↔ `foo_test.go`).
- API 위생을 개선할 수 있다면 `pkg/` 패키지(`package application_test`)에 **블랙박스 방식**을 사용하세요.
- 공유 픽스처는 사용하는 위치에 둡니다. 이 작업 트리에는 중앙 `internal/testutil/` 패키지가 없으므로 패키지별로 헬퍼를 연결하세요.

---

## 2. 단위 테스트

### 테스트 작성

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

권장 사항:

- 이미 `go.mod`에 포함된 [`stretchr/testify`](https://github.com/stretchr/testify)을 사용하세요.
- 여러 입력, 경계 사례 또는 예상 결과로 동일한 동작을 검증할 때는 가능하면 **테이블 기반** 테스트를 사용하세요. 각 사례에 설명이 분명한 이름을 지정하세요.
- 필요하면 빌드 태그(`foo_windows_test.go`, `foo_darwin_test.go`, …) 뒤에 플랫폼별 특이 동작을 스텁으로 구현하세요.

### 커버리지 요구 수준

새로 추가하거나 변경한 로직에는 Go 구문 커버리지 100%가 요구됩니다. 저장소 전체 백분율에 의존하지 말고 변경한 패키지를 측정하세요.

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

일반적인 테스트 환경에서는 합리적으로 테스트할 수 없는 경로도 있습니다. 예를 들어 특정 플랫폼에서만 발생하는 실패, 하드웨어에 종속된 동작 또는 안전하게 유발할 수 없는 방어적 대체 경로가 이에 해당합니다. 이러한 예외는 최소 범위로 제한하고 커버되지 않은 각 경로를 PR 설명에 명시하세요.

### 로컬에서 실행

```bash
cd v3
go test ./... -cover
```

또는 Taskfile을 통해 실행하세요. 실제 대상만 사용할 수 있으며 `task test` 단축 명령은 없습니다.

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. 통합 테스트

`v3/tests/`에는 패키지 간 통합 테스트 하네스가 있습니다. Taskfile의 `test:example:*` 및 `test:examples:*` 대상은 darwin/windows/linux 전반의 컴파일 검사를 실행하며, Linux에서는 Docker 기반 GTK3/GTK4 매트릭스도 포함합니다.

> 실행 가능한 예제는 `v3/examples/`에 있습니다. 테스트 대상은 다음 항목을 선택하여 빌드합니다.
>
> 호스트 플랫폼 또는 CI 매트릭스에 적합한 예제입니다.

다음 명령으로 호스트 플랫폼 스모크 테스트 모음을 실행하세요.

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 경쟁 상태 감지

GUI 런타임에서 데이터 경쟁은 치명적입니다.

### 경쟁 상태 조사

충돌하는 두 접근을 모두 조사하고 영향을 받는 플랫폼에서 재현하세요. 증거 없이 보고를 무해한 것으로 간주하거나 억제하지 마세요. `v3/TESTING.md`에서는 예제 컴파일을 설명합니다. 해당 빌드가 통과해도 데이터 경쟁이 없거나 네이티브 런타임이 올바르게 동작한다고 입증되는 것은 아닙니다.

### 로컬 경쟁 상태 테스트 모음

```
go test -race ./...
```

> `wails3 dev`에는 `-race` 플래그가 없습니다. CLI 플래그는 `--config`, `--port`,
>
> 그리고 `-s`(HTTPS 활성화)입니다. 경쟁 상태 감지기를 사용해 런타임을 검사하려면,
>
> `go build -race`을 사용하여 테스트 앱을 빌드한 후 직접 실행하세요.

---

## 5. GitHub Actions 워크플로

`.github/workflows/` 아래에 있는 실제 워크플로 파일은 다음과 같습니다(작업 트리와 대조하여 확인함).

| 파일 | 용도 |
| --- | --- |
| `build-and-test-v3.yml` | Go 1.26/1.27 매트릭스, 런타임 및 생성기 검사, 템플릿/예제 컴파일, 애플리케이션 빌드 검사입니다. Linux에서는 기본적으로 GTK4/WebKitGTK 6.0을 사용하며, 레거시 변형은 `BUILD_TAGS=gtk3`을 설정합니다. 네이티브 실행 및 상호작용은 별도로 검사해야 합니다. |
| `cross-compile-test-v3.yml` | 크로스 컴파일 기본 검사 |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | 변경 로그 자동화 |
| `nightly-release-v3.yml` | 야간 v3 릴리스 산출물 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 종속성/릴리스 관리 |
| `build-cross-image.yml` | 크로스 컴파일러 컨테이너 이미지 빌드 |
| `publish-npm.yml` | 내장 `@wailsio/runtime` JS 런타임을 npm에 게시 |
| `pr-master.yml` | `master` 브랜치를 기준으로 하는 PR 검사 |
| `semgrep.yml` | Semgrep 정적 분석 |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | 저장소 유지 관리 및 v2 워크플로입니다. 기존 Crowdin 동기화/업로드 워크플로는 폐기되었습니다. |

CI 단계는 위의 Taskfile 대상(`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …)을 그대로 따르므로 CI를 로컬에서 일대일로 재현할 수 있습니다. `build-and-test-v3.yml`의 스모크 `wails3 build` 단계는 **추가 플래그 없이** 호출됩니다. `wails3 build`에는 `-skip-package` 플래그가 없습니다.

---

## 6. 로컬 CI 동등성

모든 작업을 포괄하는 단일 `task ci` 대상은 없습니다. 실제 대상을 연이어 실행하여 CI를 재현하세요.

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 실패한 테스트 문제 해결

| 증상 | 가능성이 높은 원인 | 해결 방법 |
| --- | --- | --- |
| **`webview_window_darwin.go`에서 발생하는 경쟁 상태** | 메인 스레드 밖에서 창 상태 변경 | 호출이 메인 스레드에서 실행되도록 `application.InvokeAsync` / `Invoke`을 통해 전달합니다 |
| **헤드리스 CI에서 Linux 테스트가 멈춤** | GTK에 디스플레이가 필요함 | `xvfb-run`에서 실행합니다. 예: `xvfb-run task test:examples:linux` |
| **템플릿 빌드 실패** | 프런트엔드 잠금 파일이 최신 상태가 아님 | 템플릿을 새로 고치려면 빈 디렉터리를 대상으로 `wails3 init`을 다시 실행합니다 |
| **Coverpkg 오류** | `main`을 가져오는 통합 테스트 | 빌드 태그를 `//go:build integration`으로 변경하고 해당 가져오기를 조건부로 제한합니다 |

---

## 8. 새 테스트 추가

1. **단위 테스트** — `*_test.go`을 만들고 `go test ./...`을 실행합니다
2. **생성기 / CLI** — `internal/generator/testcases/` 또는 `internal/commands/*_test.go` 아래의 테스트 사례를 확장하고 `task test:generator` / `task test:cli`을 다시 실행합니다
3. **템플릿 / 예제** — 배포되는 템플릿이 `task test:templates`으로 계속 빌드되는지 확인합니다

---

## 9. 주요 파일 맵

| 항목 | 경로 |
| --- | --- |
| 생성기 왕복 테스트 | `internal/generator/generate_test.go` |
| 빌드 자산 테스트 | `internal/commands/build-assets_test.go` |
| 크로스 플랫폼 예제 빌드 | `v3/TESTING.md` |
| Taskfile 테스트 대상 | `v3/Taskfile.yaml` |
| 이벤트 상수 생성기 | `v3/tasks/events/generate.go` |
| CI 워크플로 | `.github/workflows/build-and-test-v3.yml`(`actions/setup-go@v5`를 통한 Go 1.26/1.27) |
| 정적 분석 | `.github/workflows/semgrep.yml` |
| 런타임 npm 게시 | `.github/workflows/publish-npm.yml` |

---

Wails v3에서 품질은 나중에 고려하는 사항이 아닙니다. 단위 테스트, 생성기/템플릿 테스트 모음, 경쟁 상태 감지, 크로스 플랫폼 CI 매트릭스를 통해 변경 사항이 지원하는 모든 OS에서 테스트를 통과한다는 확신을 갖고 기여할 수 있습니다. 즐겁게 테스트하세요!
