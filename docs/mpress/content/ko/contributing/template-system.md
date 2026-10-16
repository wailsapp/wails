---
title: "템플릿 시스템"
description: "Wails v3에서 새 프로젝트의 기본 구조를 생성하는 방식, 템플릿 구성 방식, 자체 템플릿을 만드는 방법을 설명합니다."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails에는 `wails3 init`로 바로 실행할 수 있는 프로젝트를 생성할 수 있는 <strong>템플릿 시스템</strong>이 포함되어 있습니다. 의도적으로 소수의 프레임워크(Vanilla, React, Vue, Svelte)에만 기본 제공 템플릿이 있으며, 그 밖의 프레임워크는 [자체 프런트엔드를 가져오거나](/guides/dev/frontend-frameworks/) [사용자 정의 템플릿](/guides/advanced/custom-templates/)을 게시하여 사용할 수 있습니다.

이 페이지에서 다루는 내용은 다음과 같습니다.

1. 템플릿 디렉터리 구조
2. CLI가 템플릿을 선택하고 렌더링하는 방식
3. 새 템플릿을 단계별로 만드는 방법
4. 기존 템플릿을 업데이트하거나 재정의하는 방법
5. 문제 해결 및 모범 사례

---

## 1. 템플릿 위치

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — 모든 프로젝트에 병합되는 공통 상용구(Taskfile, `build/` 디렉터리, 공유 인프라)입니다.
- **`base/`** — 모든 템플릿의 기반이 되는 Go 측 코드입니다. 참고: `base/` 자체에는 `template.json`이 **포함되어 있지 않습니다**. 이 파일은 각 프레임워크별 템플릿 안에 있습니다.
- **프레임워크 폴더** — 프런트엔드(`frontend/`), 프레임워크 설정, 템플릿 메타데이터를 설명하는 `template.json`을 포함합니다.
- 폴더 이름은 CLI에 전달하는 **템플릿 ID**(`wails3 init -t react`)와 일치합니다.
- **언어 규칙:** TypeScript가 기본이며 접미사 없는 이름(`react`)을 사용합니다. JavaScript 변형이 있는 경우 이름에 `-js` 접미사를 붙입니다(`react-js`). 기본 제공 템플릿은 `template.yaml`의 `typescript: true|false`을 사용하여 언어를 명시적으로 선언합니다. 커뮤니티 템플릿은 기존 `-ts` 접미사를 계속 사용할 수 있으며, 이 접미사는 대체 방식으로 인식됩니다.

> 전체 `internal/templates/` 디렉터리는 CLI 바이너리에 컴파일됩니다.
>
> 이는 `//go:embed *`을 통해 이루어지므로 사용자는 오프라인에서도 프로젝트의 기본 구조를 생성할 수 있습니다.

---

## 2. `wails3 init`에서 템플릿을 사용하는 방식

호출 체인은 다음과 같습니다(`cmd/wails3/init.go` 없음 — CLI는 `cmd/wails3/main.go`에서 직접 연결됨).

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

`Template.Load()` / `Template.CopyTo()` / `Template.Validate()` API는 없습니다. 템플릿 파일 추출은 내장된 `fs.FS`을 대상으로 `gosod`(`github.com/leaanthony/gosod`)에서 처리합니다.

### `wails3 init` 플래그

`internal/flags/init.go`에 정의되어 있습니다.

| 플래그 | 용도 | 기본값 |
| --- | --- | --- |
| `-p` | 패키지 이름 | `main` |
| `-t` | 기본 제공 템플릿 이름, 로컬 경로 또는 URL | `vanilla` |
| `-n` | 프로젝트 이름 | (비어 있음) |
| `-d` | 프로젝트 디렉터리 | `.` |
| `-q` | 콘솔 출력 숨기기 | false |
| `-l` | 템플릿 목록 표시 | false |
| `-skipgomodtidy` | 템플릿 파일 추출 후 `go mod tidy` 실행 건너뛰기 | false |
| `-git` | 초기화할 Git 저장소 URL | (비어 있음) |
| `-mod` | Go 모듈 경로(설정하지 않으면 `-git`에서 파생) | (비어 있음) |
| `-s` | 원격 템플릿 사용 시 경고 건너뛰기 | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | 생성된 빌드 자산에 포함되는 메타데이터 | 적절한 기본값 |

긴 별칭인 `-list`은 **없으며**(`-l`만 있음), 템플릿별 `--help`도 없습니다.

### 치환

자리표시자는 표준 Go 템플릿 지시문입니다. 앞의 `.`은 필드 접근자의 일부입니다.

| 플레이스홀더 | 예시 | 출처 |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n` 플래그 / 디렉터리 이름 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod` 플래그 또는 `-git`에서 파생 |
| `{{.WailsVersion}}` | `v3.0.0-…` | `internal/version`의 컴파일 시 포함되는 상수 |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | 빌드 시점 메타데이터 | 해당 `-product*` 플래그 |

새 플레이스홀더가 필요하면 `internal/templates/templates.go`의 템플릿 데이터에 필드를 추가하고, `internal/flags/init.go`에 일치하는 필드와 플래그를 추가하세요(또는 `internal/commands/init.go`에서 값을 설정하세요).

### 복사 후 후크

`gosod`에서 템플릿 추출을 마치면 CLI는 다음을 실행합니다:

```
go mod tidy
```

단, `-skipgomodtidy`을 전달한 경우에는 실행하지 않습니다. `task deps` 단계는 없습니다.

---

## 3. 새 템플릿 만들기

> 예시: **Solid** 템플릿 추가

### 3.1 폴더 및 ID

```
internal/templates/solid/
```

폴더 이름은 템플릿 ID입니다. <strong>kebab-case</strong>를 유지하세요.

### 3.2 최소 파일 구성

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

먼저 `react`을 복사한 다음 불필요한 파일을 제거하세요. `template.yaml`을 작성하는 것을 잊지 마세요. TypeScript 템플릿이라면 `typescript: true`을 설정하세요. 이 파일이 없는 폴더는 `base/`뿐입니다.

### 3.3 플레이스홀더 업데이트

리터럴 예시 값을 검색하여 Go 템플릿 지시문으로 바꾸세요. 예:

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 연결

`templates.go`은 초기화할 때 임베드된 파일 시스템을 순회하므로, 일반적으로 `internal/templates/<id>/` 아래에 새 폴더를 추가하는 것만으로 충분합니다. 수동 등록 호출은 필요하지 않습니다. 사용자 지정 유효성 검사나 복사 후 단계 같은 추가 로직이 필요하면 `internal/templates/templates.go`의 `templates.Install`에 추가하세요.

### 3.5 테스트

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

다음을 확인하세요:

- 개발 서버가 `WAILS_VITE_PORT`에 지정된 포트에서 시작됨
- 생성된 바인딩이 `frontend/bindings/...` 아래에 나타남
- 핫 리로드가 작동함

---

## 4. 기존 템플릿 수정하기

1. `internal/templates/<id>/` 아래의 파일을 수정하세요.
2. CLI를 다시 빌드하세요(`cd v3 && go build -o ../wails3 ./cmd/wails3`). `//go:embed *` 지시문이 새 콘텐츠를 포함합니다.
3. `frontend/package.json` 및 `Taskfile.yml`에서 필요한 <strong>종속성 버전</strong>을 올리세요.
4. 동작이 변경되면 템플릿의 `template.json` 설명을 업데이트하세요.

### 일반적인 수정

| 작업 | 위치 |
| --- | --- |
| 개발 서버 포트 변경 | `frontend/vite.config.ts` — `WAILS_VITE_PORT` 읽기 |
| 환경 변수 추가 | `build/Taskfile.yml` 또는 `frontend/.env` |
| JS 패키지 관리자 교체 | `build/Taskfile.yml`에서 `npm`을 `pnpm`/`bun`로 교체 |

---

## 5. 템플릿 작성 팁

- **프런트엔드를 범용적으로 유지하세요** — Wails 전용 전역 객체를 참조하지 마세요. `/wails/runtime.js`은 런타임에 에셋 서버에서 제공됩니다.
- **컴파일된 산출물을 포함하지 마세요** — `node_modules`, `dist`, `.DS_Store`을 임베드되는 디렉터리에서 제외하세요(또는 절대 커밋되지 않도록 `.gitignore` 처리하세요).
- **필수 조건을 문서화하세요** — Node 버전, 추가 CLI 도구 등을 `template.json` 또는 `NEXTSTEPS.md`에 기록하세요.
- **호환성을 깨는 변경을 피하세요** — 전면 개편의 규모가 크다면 기존 템플릿을 변경하는 대신 새 템플릿 ID를 만드세요.

---

## 6. 문제 해결

| 증상 | 원인 | 해결 방법 |
| --- | --- | --- |
| `unknown template name` | `-t`의 오타 또는 템플릿이 임베드되지 않음 | 사용 가능한 템플릿을 나열하려면 `wails3 init -l`을 실행하세요 |
| 플레이스홀더가 치환되지 않음 | `{{.ProjectName}}` 대신 `{{ProjectName}}`을 사용함 | 앞에 `.`을 추가하세요(Go 템플릿 필드 접근) |
| 개발 서버에서 빈 페이지가 열림 | Vite 설정이 `WAILS_VITE_PORT`을 읽지 않음 | `vite.config.ts`을 확인하세요 |
| 프로덕션 환경에서 프런트엔드 빌드 실패 | Vite의 `base` 경로 설정을 누락함 | `vite.config.ts`에서 `base: "./"`을 설정하세요 |

---

## 7. 주요 소스 파일 맵

| 파일 | 역할 |
| --- | --- |
| `internal/templates/templates.go` | 템플릿 파일 시스템을 임베드하고 `Install(options *flags.Init) error`, `GetDefaultTemplates()`, `ValidTemplateName(name)`을 노출 |
| `internal/templates/<id>/**` | 실제 템플릿 콘텐츠 |
| `internal/commands/init.go` | CLI 연결 코드: 템플릿을 선택하고 메타데이터를 채운 뒤 `templates.Install`을 호출 |
| `internal/commands/generate_template.go` | `wails3 generate template` — 실제 프로젝트를 다시 템플릿으로 *내보내는* 유틸리티(업데이트할 때 유용) |
| `internal/flags/init.go` | `wails3 init`의 플래그 정의 |

---

## 8. 요약

- 템플릿은 <strong>`internal/templates/`</strong>에 있으며 `//go:embed *`을 통해 CLI에 포함됩니다.
- `wails3 init -t <id>`은 `gosod`을 통해 템플릿을 추출하고 `go mod tidy`을 실행합니다(`-skipgomodtidy`을 사용하면 건너뛸 수 있습니다).
- 템플릿을 만들려면 **폴더를 만들고**, 파일과 `template.json`을 추가한 다음, `{{.ProjectName}}` 형식의 플레이스홀더를 사용하면 됩니다.
- 이 시스템은 **확장 가능하고** **독립적으로 완결되어 있어** 팀이나 커뮤니티와 사용자 지정 기술 스택을 공유하기에 적합합니다.

즐겁게 템플릿을 만들어 보세요!
