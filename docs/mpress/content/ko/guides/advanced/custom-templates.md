---
title: "사용자 지정 템플릿 만들기"
description: "Wails v3 프로젝트 템플릿을 직접 생성하고 사용자 지정하여 호스팅하는 방법"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails에는 기본 제공 템플릿 모음이 포함되어 있지만, 직접 템플릿을 만들어 커뮤니티와 공유할 수도 있습니다. 사용자 지정 템플릿은 단순한 Git 저장소입니다. 공개적으로 호스팅하면 누구나 명령 하나로 해당 템플릿에서 프로젝트 뼈대를 생성할 수 있습니다.

## 템플릿 뼈대 생성하기

`wails3 generate template` 명령은 바로 사용자 지정할 수 있는 템플릿 디렉터리를 생성합니다.

```bash
wails3 generate template -name MyTemplate
```

모든 플래그:

| 플래그 | 설명 | 기본값 |
| --- | --- | --- |
| `-name` | 템플릿 이름(필수) | — |
| `-author` | 작성자 이름 | — |
| `-description` | CLI에 표시되는 간단한 설명 | — |
| `-helpurl` | 이 템플릿의 문서 URL | — |
| `-version` | 초기 버전 | `v0.0.1` |
| `-frontend` | 기존 프런트엔드 디렉터리를 템플릿에 복사 | — |
| `-dir` | 템플릿 디렉터리를 생성할 위치 | 현재 디렉터리 |

모든 플래그를 사용하는 예:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

생성된 디렉터리의 구조는 다음과 같습니다.

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="NEXTSTEPS.md 읽기"}
생성된 `NEXTSTEPS.md`에는 템플릿의 각 부분에 관한 자세한 지침이 들어 있습니다. 사용자 지정하기 전에 이 파일을 읽으세요. 게시하기 전에는 삭제해야 합니다. 이 파일은 템플릿으로 만든 프로젝트에 포함되어서는 안 됩니다.

@end

## 템플릿 메타데이터 구성하기

`template.yaml`을 열어 템플릿의 메타데이터를 설정하세요.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

`wailsVersion` 필드는 **필수**이며 `3`이어야 합니다. 파일 맨 위의 `# yaml-language-server` 주석은 VS Code([YAML 확장 프로그램](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) 사용)와 JetBrains IDE에서 자동 완성과 인라인 유효성 검사를 활성화합니다. 이 주석은 그대로 두거나 제거해도 되며 런타임에는 아무런 영향을 주지 않습니다.

## 템플릿 사용자 지정하기

### 프런트엔드

`frontend/` 디렉터리는 템플릿에서 생성되는 모든 프로젝트에 그대로 복사됩니다. 자리 표시자 콘텐츠를 실제 프런트엔드로 교체하세요.

@tabs
[처음부터 시작하기]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

안내에 따라 진행한 다음 종속성을 설치하세요.

```bash
npm install
```

[기존 프로젝트 사용하기]
템플릿을 생성할 때 `-frontend`을 전달하면 기존 프런트엔드를 한 번에 복사할 수 있습니다.

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

또는 나중에 `frontend/` 디렉터리로 직접 복사하세요.

@end

### 빌드 작업

`Taskfile.tmpl.yml`은 빌드 워크플로를 정의합니다. 프런트엔드 도구 체인에 맞게 `install:frontend:deps` 및 `build:frontend` 작업을 수정하세요.

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go 애플리케이션

`main.go.tmpl` 파일은 애플리케이션 진입점입니다. 프로젝트를 생성할 때 Wails의 템플릿 엔진이 이 파일을 처리하며, `{{.ProductName}}` 같은 템플릿 변수는 사용자가 입력한 값으로 대체됩니다.

IDE 지원을 받는 실제 Go 파일로 편집하려면 파일 이름을 임시로 `main.go`으로 변경하고 수정한 다음, 커밋하기 전에 다시 `main.go.tmpl`으로 변경하세요.

#### 템플릿 변수

모든 `.tmpl` 파일에서 다음 변수를 사용할 수 있습니다.

| 변수 | 설명 | 예 |
| --- | --- | --- |
| `{{.ProjectName}}` | 사용자가 입력한 프로젝트 이름 | `"MyApp"` |
| `{{.BinaryName}}` | 바이너리 파일 이름 | `"myapp"` |
| `{{.ProductName}}` | 제품 표시 이름 | `"My Application"` |
| `{{.ProductDescription}}` | 제품 설명 | `"An awesome application"` |
| `{{.ProductVersion}}` | 제품 버전 | `"1.0.0"` |
| `{{.ProductCompany}}` | 회사/작성자 이름 | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | 저작권 문자열 | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | 추가 제품 설명 | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | 역방향 DNS 형식의 제품 식별자 | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Go 모듈 경로 | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | 프로젝트 생성에 사용된 Wails 버전 | `"3.0.0"` |
| `{{.Typescript}}` | 템플릿 이름이 `-ts`으로 끝나는 경우 `true` | `true` |
| `{{.Opn}}` | 리터럴 `{{` — 템플릿 안에서 이스케이프 | `{{` |
| `{{.Cls}}` | 리터럴 `}}` — 템플릿 안에서 이스케이프 | `}}` |

@note{type="tip"}
HTML, JSON, YAML 파일을 비롯해 템플릿의 모든 파일을 `.tmpl` 파일로 사용할 수 있습니다. `.tmpl` 접미사가 없는 파일은 그대로 복사됩니다.

@end

## 로컬에서 템플릿 테스트하기

게시하기 전에 로컬 경로의 템플릿으로 프로젝트를 생성하여 테스트하세요.

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

그런 다음 프로젝트가 제대로 작동하는지 확인하세요.

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

다음 사항을 확인하세요.

- 프런트엔드 핫 리로드가 작동하는지
- Go 코드 변경 시 앱이 다시 빌드되고 재실행되는지
- `bin/`의 프로덕션 바이너리가 올바르게 실행되는지

## GitHub에 게시하기

@steps
### 템플릿용 **공개 GitHub 리포지토리를 만드세요**. 리포지토리 루트에는 `template.yaml`이(가) 있어야 합니다.
### **`NEXTSTEPS.md`을(를) 삭제하세요**. 이 파일은 템플릿 작성자를 위한 안내 자료이므로 사용자가 템플릿으로 생성하는 프로젝트에 포함되어서는 안 됩니다.
### 템플릿 디렉터리의 내용을 리포지토리 루트로 **커밋하고 푸시하세요**:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### 시맨틱 버전 관리를 사용하여 **릴리스에 태그를 지정하세요**:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

이제 사용자는 템플릿으로 프로젝트를 생성할 수 있습니다:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="서드 파티 템플릿 경고"}
사용자가 원격 템플릿을 설치하면 Wails는 해당 템플릿이 서드 파티 코드이며 Wails 프로젝트는 그 내용에 어떠한 책임도 지지 않는다는 경고를 표시합니다. 프로젝트를 생성하려면 사용자가 명시적으로 확인해야 합니다.

템플릿 작성자는 템플릿에 포함된 모든 코드의 보안과 정확성에 책임이 있습니다.

템플릿은 코드입니다. `Taskfile.yml` 작업은 임의의 셸 명령이며, `wails3 dev`, `wails3 build` 및 `wails3 package`은(는) 이러한 명령과 명령이 호출하는 모든 `package.json` 스크립트를 사용자 권한으로 실행합니다. 이러한 명령을 실행하기 전에 원격 템플릿의 `Taskfile.yml`, `build/` 디렉터리 및 `frontend/package.json`을(를) 검토하세요.

@end

## 권장 사례

- 명확한 **`README.md`을(를) 작성하세요**. 이 내용은 사용자가 프로젝트를 생성한 후 표시됩니다. 프로젝트를 실행하고 빌드하며 사용자 지정하는 방법을 설명하세요.
- **`helpurl`을(를) 입력하세요**. 리포지토리 또는 전용 문서로 연결되는 링크를 지정하세요. 사용자는 Wails CLI 템플릿 목록에서 이 링크를 볼 수 있습니다.
- 업스트림 업데이트로 인해 설치가 실패하지 않도록 `package.json`에서 **프런트엔드 종속성 버전을 고정하세요**.
- **태그를 지정하기 전에 테스트하세요**. 커뮤니티에 공개하기 전에 태그가 지정된 릴리스로 새 프로젝트를 생성하세요.
- **`wailsVersion: 3`을(를) 그대로 유지하세요**. 이 필드는 템플릿이 대상으로 하는 Wails 주 버전을 지정합니다. 변경하지 마세요.
- **정기적으로 업데이트하세요**. 종속성을 최신 상태로 유지하고 새로운 Wails 릴리스에서 테스트하세요.
