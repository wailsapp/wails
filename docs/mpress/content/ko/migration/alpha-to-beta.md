---
title: "v3 알파에서 업그레이드"
description: "기존 Wails v3 알파 프로젝트를 고정된 베타 버전으로 업데이트"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

이 안내서는 기존 v3 알파 프로젝트를 위한 것입니다. Wails v2에서는 [v2에서 v3로 이전하는 안내서](/migration/v2-to-v3/)를 사용하세요.

## 업그레이드 전

프로젝트를 커밋하거나 백업하세요. 현재 알파 버전부터 선택한 베타까지의 [변경 로그](/changelog/)를 읽으세요. 소스 코드, API 또는 빌드 설정을 변경해야 할 수 있습니다. [데스크톱 호환성 정책](/status/)과 플랫폼 요구 사항을 확인하세요.

아래 명령은 게시된 `v3.0.0-beta.23` 릴리스를 정확한 버전 지정의 예로 사용하며, 항상 최신 릴리스를 추적하라는 권장이 아닙니다. 다른 릴리스를 선택한다면 해당 CLI, Go 모듈, npm 런타임 버전을 확인하고 명령을 함께 수정하세요. 이 예에서 npm 버전은 Go 버전에서 `v` 접두사를 뺀 것과 같습니다.

## 1. CLI 업데이트

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

`wails3 version`이 설치한 버전을 표시하는지 확인하세요. `PATH`에서 앞에 있는 이전 실행 파일 때문에 새 CLI가 실행되지 않을 수 있습니다.

## 2. Go 모듈 업데이트

프로젝트 루트에서 실행하세요. 의존성 변경을 검토하고 관련 없는 모듈을 일괄 업데이트하지 마세요.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. 프런트엔드 런타임 업데이트

npm과 `frontend` 디렉터리를 사용하는 프로젝트의 경우:

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

잠금 파일을 유지하고 변경 사항을 검토하세요. 프런트엔드에서 다른 패키지 관리자나 디렉터리를 사용한다면 정확한 런타임 버전을 유지하면서 이 단계를 조정하세요.

## 4. 재생성, 빌드 및 테스트

프로젝트 루트에서 Go 서비스로부터 바인딩을 다시 생성하고 빌드하세요.

```sh
wails3 generate bindings
wails3 build
```

빌드한 애플리케이션을 실행하고 배포하는 각 지원 플랫폼에서 작업 흐름을 테스트하세요. 소스, 생성된 바인딩, 모듈 파일 및 프런트엔드 잠금 파일의 변경 사항을 검토하고 함께 커밋하세요.

## 업그레이드에 실패한 경우

`PATH`의 CLI를 확인하고, `go list -m github.com/wailsapp/wails/v3`로 모듈 버전을, `npm --prefix frontend ls @wailsio/runtime`으로 설치된 런타임을 확인하세요. 버전 불일치를 해결한 뒤 바인딩을 다시 생성하세요. 모든 알파 버전이 코드 변경 없이 업그레이드될 수 있다고 가정하지 마세요.

문제가 계속되면 이전 및 새 버전, 정확한 오류, `wails3 doctor` 출력을 포함하여 [재현 가능한 이슈를 보고](https://github.com/wailsapp/wails/issues/new/choose)하세요. 취약점은 [보안 정책](https://github.com/wailsapp/wails/blob/master/SECURITY.md)을 따르세요.
