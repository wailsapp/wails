---
title: "빌드 및 패키징 파이프라인"
description: "`wails3 build`를 실행할 때 내부에서 수행되는 작업, 크로스 플랫폼 바이너리가 생성되는 방식, 각 OS용 설치 프로그램이 만들어지는 방식을 설명합니다."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build`은 의도적으로 **간결하게** 설계되었습니다. 호스트 프로젝트의 `build` 태스크에 추가 빌드 태그를 전달하는 Taskfile 래퍼입니다. 핵심 작업은 프로젝트 자체의 `build/Taskfile.yml`(`wails3 init`에서 생성), 베이크 시점 자산을 관리하는 `internal/commands/build-assets.go`, Linux nfpm 패키징을 담당하는 `internal/packager`, 그리고 플랫폼별 설치 프로그램을 담당하는 `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, 및 `internal/commands/dot_desktop.go`에서 수행됩니다.

이 페이지에서 다루는 내용은 다음과 같습니다.

1. 실제 CLI 진입점
2. Taskfile 기반 빌드 흐름
3. 자산 베이크 및 빌드 정보 삽입
4. 플랫폼별 패키징 백엔드
5. 파이프라인 사용자 지정
6. 문제 해결

---

## 1. 실제 CLI 진입점

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`:

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build`에서 제공하는 플래그는 <strong>단 하나</strong>인 `--tags`뿐이며, 이 플래그는 `EXTRA_TAGS=`로 전달됩니다. `wails3 build`에는 `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` 또는 `-clean` 플래그가 **전혀** 없습니다. 크로스 컴파일, 출력 경로, 아이콘 등은 **`Taskfile.yml`**, **`build/config.yml`** 및 이를 지원하는 `wails3 generate icons` / `wails3 generate build-assets` 명령에서 구성합니다.

`build/build.json`은 v3의 일부가 **아닙니다**. 구성에는 `Taskfile.yml`와 `build/config.yml`를 사용합니다.

---

## 2. Taskfile 기반 빌드 흐름

새로 초기화한 프로젝트에는 대략 다음과 같은 네임스페이스를 포함하는 `build/Taskfile.yml`이 제공됩니다.

| 네임스페이스 | 태스크(일부) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

`wails3 build`은 기본적으로 호스트 OS의 `build` 네임스페이스를 호출하며, 프로젝트의 Taskfile은 다시 호스트별 플래그를 사용하여 `go build`을 실행합니다. 다른 OS용으로 빌드하려면 `wails3 build`에 플래그를 전달하지 말고 해당 태스크를 직접 실행해야 합니다(예: `wails3 task darwin:build:universal`).

기본 출력 디렉터리는 <strong>`bin/<APP_NAME>`</strong>이며, `build/bin/` 접두사는 없습니다.

---

## 3. 베이크 시점 자산 및 빌드 정보

| 항목 | 파일 |
| --- | --- |
| 빌드 자산 생성/업데이트 | `internal/commands/build-assets.go` |
| 빌드 정보 출력기(CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` — 정보를 출력하며 `ldflags` 삽입기가 **아닙니다** |
| 프로덕션 스텁 | `internal/assetserver/build_production.go` — `//go:build production` |
| 프런트엔드 번들 | 애플리케이션 자체 패키지에서 `//go:embed`을 통해 임베드됩니다(예: `main.go` 옆). |
| Windows 리소스(`.syso`) | `internal/commands/syso.go` — `rsrc_windows_<arch>.syso` 생성 |
| Windows MSIX | `internal/commands/msix.go` + `internal/commands/webview2/` |
| macOS DMG 입력 | `internal/commands/dmg/` |
| Linux `.desktop` | `internal/commands/dot_desktop.go` |

CLI는 애플리케이션용 `bundled_assetserver.go`을 자동으로 베이크하지 않습니다. `internal/assetserver/bundled_assetserver.go`은 **직접 작성되며**, `bundledassets/` 아래에 임베드된 JS 런타임을 래핑합니다.

---

## 4. 패키징 백엔드

### Linux

Linux 패키징은 `fpm`이 아니라 <strong>nfpm</strong>에서 처리합니다.

- `internal/packager/packager.go`은 `github.com/goreleaser/nfpm/v2`을 래핑하고 `CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)`을 제공합니다.
- 생성된 프로젝트에는 `internal/commands/` 아래에 nfpm 형식의 `myapp.DEB`, `myapp.RPM` 및 `myapp.ARCHLINUX` 구성 파일이 포함되며, 이 파일들은 `wails3 tool package`에서 사용됩니다.
- AppImage 생성 기능은 `internal/commands/appimage.go`에 있으며, 여기서 `linuxdeploy` + `linuxdeploy-plugin-gtk`을 호출합니다. 플러그인은 `internal/commands/linuxdeploy-plugin-gtk.sh`에 포함되어 있습니다.

`wails3 build`에는 `-package deb`/`rpm` 플래그가 **없습니다**. `wails3 tool package` 또는 플랫폼별 Taskfile 대상을 사용하세요.

### macOS

- 프로젝트 Taskfile의 `darwin:package`은 `.app` 번들을 생성합니다.
- DMG 자산은 `internal/commands/dmg/` 아래에 있습니다. 프로젝트에서는 `darwin:package` 작업이 끝난 후 `hdiutil`을 사용하여 번들을 DMG로 패키징할 수 있습니다(최신 템플릿의 Taskfile에는 `dmg` 헬퍼가 포함되어 있습니다).
- CFBundle 식별자, 버전 및 저작권 정보는 `wails3 init` 실행 시 사용하는 `-product*` 플래그와 `build/config.yml`에서 가져옵니다.

### Windows

- Windows 패키징은 WiX/MSI가 아닌 <strong>MSIX</strong>를 대상으로 합니다. 전체 워크플로는 `internal/commands/msix.go` 및 `internal/commands/webview2/`에서 확인하세요.
- **`internal/commands/packager.go`는 없으며**, **`internal/commands/windows_resources/` 디렉터리도 없습니다**.
- 실행 파일의 선택적 코드 서명은 `wails3 tool sign`(Authenticode)을 통해 실행됩니다. 자세한 내용은 `internal/commands/sign.go`을 참조하세요.

---

## 5. 파이프라인 사용자 지정

| 요구 사항 | 방법 |
| --- | --- |
| 추가 빌드 태그 | `wails3 build --tags myFeature,otherTag` |
| 린터/빌드 전 단계 | `build/Taskfile.yml`에 태스크를 추가하고 OS별 `build` 태스크가 이 태스크에 종속되도록 설정하세요. |
| 크로스 컴파일 | 관련 OS 태스크(예: `wails3 task linux:build`)를 실행하세요. `-platform` 플래그는 없습니다. |
| 패키징 건너뛰기 | `build` 태스크만 실행하세요. `package`은 별도 태스크입니다. |
| 사용자 지정 패키저 | `internal/commands/myapp.*` 아래에 설정을 추가하고 `-config <file>`을 지정하여 `wails3 tool package`을 호출하세요. |
| 심벌 제거 | `darwin:/windows:/linux:`의 `build` 태스크를 편집하여 `-ldflags "-s -w"`을 `go build`에 직접 전달하세요. `wails3 build` 자체에는 `-ldflags` 플래그가 없습니다. |

모든 Taskfile 대상은 Wails가 제공하는 환경 변수(`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL` 등)를 따르므로 사용자 지정 태스크에서도 이러한 변수를 사용할 수 있습니다.

---

## 6. 문제 해결

| 증상 | 가능성이 높은 원인 | 해결 방법 |
| --- | --- | --- |
| **`ld: framework not found WebKit`(mac)** | Xcode CLI 도구가 설치되어 있지 않음 | `xcode-select --install` |
| **프로덕션 빌드에서 빈 창이 표시됨** | 프런트엔드 빌드 실패 또는 SPA 라우팅 문제 | `frontend/dist/index.html`이 존재하며 자산 핸들러가 이를 폴백으로 사용하는지 확인하세요. |
| **MSIX 패키징 도구가 설치되어 있지 않음** | `WebView2` SDK/MSIX 도구가 설치되어 있지 않음 | `wails3 task install:msix:tools`을 실행하세요. |
| **`linuxdeploy`을 찾을 수 없음** | PATH에 플러그인이 없음 | `linuxdeploy`을 설치하고 CLI의 자동 설치 단계를 통해 `internal/commands/linuxdeploy-plugin-gtk.sh`을 실행하세요. |

`wails3 build`에는 `-verbose` 플래그가 없습니다. `TASK_X_VERBOSE=1`(Taskfile)을 설정하거나 태스크 대상을 직접 살펴보고 실행되는 명령을 확인하세요.

---

## 7. 주요 소스 맵

| 관련 기능 | 파일 |
| --- | --- |
| 빌드 래퍼 | `internal/commands/task_wrapper.go`(`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| 빌드 자산 생성 | `internal/commands/build-assets.go`(`GenerateBuildAssets`, `UpdateBuildAssets`) |
| 빌드 정보 출력기 | `internal/commands/tool_buildinfo.go` |
| AppImage 빌더 | `internal/commands/appimage.go` |
| Linux 패키징(nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| Windows MSIX | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Windows 리소스 생성기 | `internal/commands/syso.go` |
| macOS DMG 자산 | `internal/commands/dmg/` |
| `.desktop` 생성기 | `internal/commands/dot_desktop.go` |
| 버전 상수 | `internal/version/version.go` |

빌드 실패를 추적할 때 이 표를 가까이 두고 활용하세요.

---

이제 <strong>소스 코드</strong>부터 <strong>설치 프로그램</strong>까지 전체 흐름을 파악했습니다. 요약하면 `wails3 build` 자체는 얇은 래퍼에 불과합니다. 거의 모든 사용자 지정은 프로젝트의 `Taskfile.yml`/`build/config.yml`에서 이루어지거나 명시적인 `wails3 generate …`/`wails3 tool …` 하위 명령을 통해 수행됩니다. 성공적으로 배포하시기 바랍니다!
