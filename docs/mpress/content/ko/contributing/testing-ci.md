---
title: "테스트 및 지속적 통합"
description: "Wails v3가 단위 테스트, 통합 테스트 모음, 경쟁 상태 감지 및 GitHub Actions CI를 통해 품질을 보장하는 방법을 설명합니다."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

견고한 데스크톱 프레임워크에는 철저한 테스트가 필요합니다. Wails v3는 다음과 같은 <strong>계층형 전략</strong>을 사용합니다:

| 계층 | 목표 | 도구 |
| --- | --- | --- |
| 단위 테스트 | 개별 함수에 대한 빠른 피드백 | `go test ./...` |
| 생성기/CLI 테스트 | `wails3 generate bindings` 및 CLI 연결부 검증 | `task test:generator`, `task test:cli` |
| 템플릿 테스트 | 배포되는 모든 템플릿이 계속 빌드되는지 확인 | `task test:templates` |
| 경쟁 상태 감지 | 런타임 및 브리지의 데이터 경쟁 감지 | `go test -race ./...` |
| CI 매트릭스 | 모든 PR에서 운영체제 간 신뢰성 확보 | GitHub Actions |

이 문서에서는 **테스트가 있는 위치**, **테스트를 실행하는 방법**, 그리고 <strong>Taskfile이 조정하는 작업</strong>을 설명합니다.

> 이전 초안에서 `pkg/application/RACE.md`로 참조했던 경쟁 상태 관련 지침은
>
> 현재 `v3/TESTING.md`에 있습니다.

---

## 1. 디렉터리 규칙

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

지침:

- **단위 테스트는 코드 옆에 두세요**(`foo.go` ↔ `foo_test.go`).
- API의 건전성을 개선할 수 있다면 `pkg/` 패키지(`package application_test`)에 <strong>블랙박스 방식</strong>을 사용하세요.
- 공유 픽스처는 사용되는 곳에 두세요(이 작업 트리에는 중앙 `internal/testutil/` 패키지가 없으므로, 대신 패키지별로 헬퍼를 연결하세요).

---

## 2. 단위 테스트

### 테스트 작성

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

권장 사항:

- 이미 `go.mod`에 있는 [`stretchr/testify`](https://github.com/stretchr/testify)을 사용하세요.
- 여러 입력, 경계 사례 또는 예상 결과로 동일한 동작을 검사할 때는 가능하면 **테이블 기반** 테스트를 사용하세요. 각 사례에 설명적인 이름을 지정하세요.
- 필요한 경우 플랫폼별 특이 동작을 빌드 태그(`foo_windows_test.go`, `foo_darwin_test.go`, …) 뒤에 스텁으로 구현하세요.

### 커버리지 기준

새로 추가하거나 변경한 로직은 Go 문장 커버리지가 100%여야 합니다. 저장소 전체 비율에 의존하지 말고 변경한 패키지를 측정하세요:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

일반적인 테스트 환경에서는 합리적으로 테스트할 수 없는 경로도 있습니다. 예를 들면 특정 플랫폼에서만 발생하는 오류, 하드웨어에 종속된 동작, 안전하게 유발할 수 없는 방어적 대체 경로가 있습니다. 이러한 예외는 최소한으로 제한하고, 커버되지 않은 각 경로를 PR 설명에 명시하세요.

### 로컬에서 실행

```bash
cd v3
go test ./... -cover
```

또는 Taskfile을 통해 실행하세요(실제 대상이며 `task test` 단축 명령은 없습니다):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. 통합 테스트

`v3/tests/`에는 패키지 간 통합 테스트 하네스가 있습니다. Taskfile의 `test:example:*` 및 `test:examples:*` 대상은 darwin / windows / linux 전반에서 빌드 및 실행 검사를 수행합니다(Linux에서는 Docker 기반 GTK3 / GTK4 매트릭스 포함).

> 실행 가능한 예제는 `v3/examples/`에 있습니다. 테스트 대상은 호스트 플랫폼 또는 CI 매트릭스에
>
> 적합한 예제를 선택하고 빌드합니다.

다음 명령으로 호스트 플랫폼의 스모크 테스트 모음을 실행하세요:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 경쟁 상태 감지

GUI 런타임에서 데이터 경쟁은 치명적입니다.

### 경쟁 상태 지침

다음 내용은 `v3/TESTING.md`에서 확인하세요:

- 알려진 무해한 경쟁 상태와 이를 억제하는 근거
- Cgo 경계를 가로지르는 스택 추적을 해석하는 방법(Linux GTK + WebKit2GTK)

### 로컬 경쟁 상태 테스트 모음

```
go test -race ./...
```

> `wails3 dev`에는 `-race` 플래그가 없습니다. CLI 플래그는 `--config`, `--port`,
>
> 그리고 `-s`(HTTPS 활성화)입니다. 경쟁 상태 감지기를 사용해 런타임을 검사하려면
>
> `go build -race`으로 테스트 앱을 빌드한 후 직접 실행하세요.

---

## 5. GitHub Actions 워크플로

`.github/workflows/` 아래에 있는 실제 워크플로 파일(작업 트리와 대조하여 확인함):

| 파일 | 용도 |
| --- | --- |
| `build-and-test-v3.yml` | 기본 v3 빌드 및 테스트 매트릭스입니다. `go-version: 1.25`과 함께 `actions/setup-go@v5`을 사용합니다. `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples`(GTK4 경로에서는 `BUILD_TAGS=gtk4 task test:examples`도 포함), `task generator:test:check`, `task install`을 실행한 다음, 스모크 검사로 `wails3 build`을 실행합니다. Linux 작업은 `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk`을 설치하고 `dbus-run-session -- xvfb-run`에서 테스트 모음을 실행합니다. |
| `cross-compile-test-v3.yml` | 교차 컴파일 정상 동작 검사 |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | 변경 로그 자동화 |
| `nightly-release-v3.yml` | v3 나이틀리 릴리스 아티팩트 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 종속성/릴리스 관리 |
| `build-cross-image.yml` | 크로스 컴파일러 컨테이너 이미지를 빌드 |
| `publish-npm.yml` | 내장된 `@wailsio/runtime` JS 런타임을 npm에 게시 |
| `pr-master.yml` | `master` 브랜치를 기준으로 하는 PR 측 검사 |
| `semgrep.yml` | Semgrep 정적 분석 |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | 저장소 유지 관리/v2 측 워크플로 |

이 트리에는 `qodana.yaml`도 **없고** `runtime.yml`도 **없습니다**. 이 페이지의 이전 초안에는 둘 다 명시되어 있었지만, 정적 분석은 `semgrep.yml`만 담당하며 런타임 JS 패키지는 `publish-npm.yml`을 통해 게시됩니다.

CI 단계는 위의 Taskfile 대상(`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …)과 대응하므로 로컬에서 CI를 일대일로 재현할 수 있습니다. `build-and-test-v3.yml`의 스모크 `wails3 build` 단계는 **추가 플래그 없이** 호출됩니다. `wails3 build`에는 `-skip-package` 플래그가 없습니다.

---

## 6. 로컬 CI와의 일치

모든 작업을 포괄하는 단일 `task ci` 대상은 없습니다. 실제 대상을 연이어 실행하여 CI를 재현하세요:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 실패한 테스트 문제 해결

| 증상 | 가능성 높은 원인 | 해결 방법 |
| --- | --- | --- |
| <strong>`webview_window_darwin.go`</strong>의 경쟁 상태 | 메인 스레드가 아닌 곳에서 창 상태를 변경함 | 호출이 메인 스레드에서 실행되도록 `application.InvokeAsync`/`Invoke`을 통해 마샬링 |
| **헤드리스 CI에서 Linux 테스트가 멈춤** | GTK에 디스플레이가 필요함 | `xvfb-run`에서 실행(예: `xvfb-run task test:examples:linux`) |
| **템플릿 빌드 실패** | 프런트엔드 잠금 파일이 최신 상태가 아님 | 깨끗한 디렉터리를 대상으로 `wails3 init`을 다시 실행하여 템플릿 새로 고침 |
| **Coverpkg 오류** | 통합 테스트에서 `main`을 가져옴 | 빌드 태그를 `//go:build integration`으로 전환하고 가져오기를 조건부로 제한 |

---

## 8. 새 테스트 추가

1. **단위 테스트** — `*_test.go`을 만들고 `go test ./...`을 실행하세요
2. **생성기/CLI** — `internal/generator/testcases/` 또는 `internal/commands/*_test.go` 아래의 테스트 사례를 확장하고 `task test:generator`/`task test:cli`을 다시 실행하세요
3. **템플릿/예제** — 배포되는 템플릿이 `task test:templates`으로 계속 빌드되는지 확인하세요

---

## 9. 주요 파일 안내

| 항목 | 경로 |
| --- | --- |
| 생성기 왕복 테스트 | `internal/generator/generate_test.go` |
| 빌드 자산 테스트 | `internal/commands/build-assets_test.go` |
| 경쟁 상태/Cgo 가이드 | `v3/TESTING.md` |
| Taskfile 테스트 대상 | `v3/Taskfile.yaml` |
| 이벤트 상수 생성기 | `v3/tasks/events/generate.go` |
| CI 워크플로 | `.github/workflows/build-and-test-v3.yml`(`actions/setup-go@v5`을 통한 Go 1.25) |
| 정적 분석 | `.github/workflows/semgrep.yml` |
| 런타임 npm 게시 | `.github/workflows/publish-npm.yml` |

---

Wails v3에서 품질은 나중에 고려하는 사항이 아닙니다. 단위 테스트, 생성기/템플릿 테스트 모음, 경쟁 상태 감지, 크로스 플랫폼 CI 매트릭스를 갖추고 있으므로 변경 사항이 지원되는 모든 OS에서 문제없이 실행된다는 확신을 가지고 기여할 수 있습니다. 즐겁게 테스트하세요!
