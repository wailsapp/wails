---
title: "M-Press로 문서 작성하기"
description: "M-Press를 다운로드하여 설치하고, Wails 문서를 미리 본 후 문서 기여를 제출합니다."
sourcePath: "contributing/documentation.md"
---

Wails v3에서는 **M-Press**를 사용하여 Markdown을 이 문서 사이트로 변환합니다. 오타를 수정하거나, 설명을 개선하거나, 예제 또는 가이드를 추가하거나, 번역을 바로잡을 수 있습니다. 문서만 수정하는 경우에는 사전에 이슈를 만들거나 실패하는 코드 테스트를 준비할 필요가 없습니다.

기여를 준비하려면 Git이 필요하고, 풀 리퀘스트를 제출하려면 GitHub 계정이 필요합니다. 다운로드 가능한 M-Press 실행 파일에는 문서 생성기와 다이어그램 렌더러가 포함되어 있으므로 문서를 미리 보는 데 Go, Node.js, Wails CLI 또는 데스크톱 빌드 종속성은 필요하지 않습니다. 저장소의 추가 검증 스크립트를 실행하려면 Python 3이 필요합니다.

## 문서 페이지에서 시작하기

이 가이드를 열려면 사이드바에서 **기여 → 문서**를 선택합니다. 특정 페이지에서 작업하려면 해당 페이지 헤더의 **편집 아이콘**을 선택한 다음 **이 페이지를 로컬에서 편집**을 선택합니다. 사용 중인 운영 체제에 맞는 명령을 복사하여 터미널에서 실행합니다.

M-Press가 설치되어 있으면 이 명령은 M-Press를 사용합니다. 설치되어 있지 않으면 설치 프로그램이 공식 릴리스를 다운로드하고 체크섬을 검증합니다. 그런 다음 로컬 기여용 작업 트리를 준비하고 선택한 페이지를 M-Press에서 엽니다. 저장소 검사를 실행하기 전에 `mpress version`을 아래에 고정된 버전과 대조해 확인하세요. 자동 설치 프로그램이 더 최신 릴리스를 다운로드할 수 있습니다.

로컬 미리보기가 다시 빌드되는 동안 원하는 편집기에서 소스를 편집할 수 있습니다. 마법사에는 **변경 사항 검토**, **변경 사항 커밋**, **기여 브랜치 푸시**, **초안 풀 리퀘스트 열기**도 각각 별도의 작업으로 제공됩니다. 제출하기 전에 아래의 검사 및 PR 지침을 따르세요. 번역 생성은 선택 사항이며 구성된 제공자가 필요할 수 있지만, Markdown 편집과 미리보기에는 필요하지 않습니다.

재현 가능한 수동 설정을 구성하려면 아래 단계를 따르세요.

## M-Press 다운로드 및 설치

저장소 빌드 스크립트는 **최신 안정 버전의 M-Press 릴리스**를 다운로드합니다. 로컬에서도 같은 릴리스를 사용하세요. 아래에서 설치 방법을 선택하고 `mpress version`을 실행한 다음, 릴리스 페이지에 표시된 버전과 비교하세요.

### 실행 파일 다운로드

[최신 안정 버전의 M-Press 릴리스](https://github.com/leaanthony/mpress/releases/latest)를 열거나 아래에서 해당 아카이브를 다운로드하세요. ARM64에는 Apple Silicon이 포함되며, AMD64는 일반적인 Intel/AMD PC 아키텍처입니다.

| 컴퓨터 | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

같은 릴리스에서 [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt)를 다운로드하세요. 다운로드 디렉터리에서 아카이브의 SHA-256 해시를 계산하고 `checksums.txt`에서 해당 파일 이름이 있는 줄과 비교하세요.

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

다른 아키텍처를 선택했다면 실제로 다운로드한 파일 이름을 사용하세요. 해시가 일치하면 실행 파일을 압축 해제하고 설치하세요.

@tabs
[macOS]
다운로드 디렉터리에서 Apple Silicon 아카이브를 사용하여 실행하세요.

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Intel Mac에서는 `mpress-darwin-amd64.tar.gz`으로 바꾸세요. 새 터미널에서도 명령을 사용할 수 있도록 `export PATH` 줄을 `~/.zshrc`에 추가하세요.

[Linux]
다운로드 디렉터리에서 실행하세요.

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

ARM64에서는 `mpress-linux-arm64.tar.gz`으로 바꾸세요. 새 터미널에서도 명령을 사용할 수 있도록 `export PATH` 줄을 `~/.bashrc` 등의 셸 시작 파일에 추가하세요.

[Windows PowerShell]
다운로드 디렉터리에서 실행하세요.

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

ARM64에서는 `mpress-windows-arm64.zip`으로 바꾸세요. 이후 세션에서도 명령을 사용할 수 있도록 사용자 `bin` 디렉터리를 Windows **환경 변수**의 사용자 **Path**에 추가한 다음 터미널을 다시 여세요.
@end

### Go로 설치하기

이미 Go를 사용하고 있다면 최신 안정 버전 릴리스를 직접 빌드할 수도 있습니다. 필요한 도구 체인은 해당 릴리스의 `go.mod`에서 확인하세요. Go의 자동 도구 체인 선택 기능이 이를 다운로드할 수 있습니다. 이 문서 작업에 사용된 릴리스에는 Go 1.26.10 이상이 필요합니다.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go는 `GOBIN`에 명령을 설치합니다. `GOBIN`이 설정되지 않은 경우에는 `go env GOPATH` 아래의 `bin` 디렉터리에 설치합니다. `mpress`을 찾을 수 없다면 해당 디렉터리를 `PATH`에 추가하고 터미널을 다시 여세요.

## 문서 소스 가져오기

[wailsapp/wails를 포크](https://github.com/wailsapp/wails/fork)한 다음, `YOUR_USERNAME`을 자신의 GitHub 사용자 이름으로 바꾸어 다음 명령을 실행하세요.

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

그러면 업스트림 **master** 브랜치를 기준으로 문서 브랜치가 생성됩니다. `docs/my-change` 대신 내용을 잘 나타내는 고유한 브랜치 이름을 사용하세요.

## 실시간 미리보기 실행하기

`mpress.yaml`이 있는 Wails 저장소 루트에서 실행하세요.

```sh
mpress dev
```

브라우저에서 [http://127.0.0.1:3000](http://127.0.0.1:3000)을 여세요. 이 가이드는 [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/)에서 확인할 수 있습니다. 편집하는 동안 명령을 계속 실행해 두세요. 소스 변경 사항을 저장하면 미리보기가 다시 빌드됩니다. 아래의 검증 명령을 실행하기 전에 **Ctrl+C**를 눌러 중지하세요.

저장소에는 M-Press 구성이 이미 포함되어 있습니다. 새 사이트를 초기화할 필요가 없습니다. 다른 포트를 사용하려면 `mpress dev --port 3001`을 실행하고 해당 localhost 주소를 여세요.

## 페이지 편집 또는 추가하기

영어 소스 페이지는 `docs/mpress/content/`에 있습니다. 번역은 `fr/` 및 `id/` 같은 언어 폴더 아래에 있습니다. 편집 아이콘의 기여 안내에서 현재 페이지의 소스 파일을 확인할 수 있습니다.

기존 페이지를 수정하려면 해당 Markdown을 편집하고 상단의 YAML 메타데이터를 그대로 유지하세요. 제목, 단락, 목록, 링크 및 펜스 코드 블록을 사용하세요. `@note` / `@end` 및 `@tabs` / `@end`처럼 쌍을 이루는 M-Press 구성 요소를 그대로 유지하세요. 지원되는 구문은 [M-Press 작성 문서](https://m-press.me/authoring/)에서 설명합니다.

새 가이드는 `docs/mpress/content/guides/my-guide.md` 같은 파일을 만드세요.

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

해당 섹션에서 이미 `autogenerate`을 사용하고 있지 않다면 `docs/mpress/content/_nav.yaml`의 적절한 섹션에 페이지를 추가하세요. 독자가 탐색 메뉴와 관련 페이지의 링크를 통해 해당 페이지에 접근할 수 있는지 확인하세요. 새 영어 페이지는 모든 로케일의 사본이 완성되어야 범위 검사를 통과할 수 있으므로 관리자 또는 번역 기여자와 조율하세요.

생성된 출력은 `docs/mpress/site/`에 저장됩니다. 소스 페이지와 의도한 탐색 또는 정적 자산 변경 사항은 커밋하되, 생성된 사이트 파일과 로컬 캐시는 PR에서 제외하세요.

## 기여 내용 검사하기

저장소 루트에서 실행하세요.

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Windows에서 Python 3을 그런 방식으로 설치했다면 `python`을 사용하세요. [Task](https://taskfile.dev/)가 이미 설치되어 있다면 `task docs:build`을 실행하여 PATH에 있는 `mpress`으로 동일한 검사를 수행할 수 있습니다.

번역 검사는 번역 파일이 누락되었거나 비어 있거나 대응하는 원문 없이 남아 있으면 실패합니다. 오래되었거나 손상된 번역 구절은 기본적으로 경고로 보고되며, `--strict`을 사용하면 이러한 문제도 검사 실패로 처리됩니다. 기존 영어 문구를 수정한 경우 번역 동기화가 완료되기 전에도 제출할 수 있습니다. 특정 언어의 번역을 바로잡을 때는 해당 언어만 변경할 수 있습니다. 명령어, API 이름, 링크, 코드 예제 및 다이어그램 연결은 그대로 유지하고, 영어 자리표시자를 번역된 콘텐츠로 사용하지 마세요.

변경한 페이지를 브라우저에서 확인하고 수정한 코드 예제를 모두 실행하세요. 수정하지 않은 페이지 때문에 검사가 실패하면 검사를 완화하지 말고 출력과 관련 맥락을 PR에 포함하세요.

## 풀 리퀘스트 제출

diff를 검토하고 소스 변경 사항을 커밋한 다음, 기여 브랜치를 자신의 포크에 푸시하세요. `wailsapp/wails`에서 **master** 브랜치를 대상으로 PR을 여세요. 변경 사항을 설명하고, 영향을 받는 페이지의 링크를 추가하며, 실행한 검사 목록을 작성하세요. 레이아웃을 변경한 경우 변경 전후의 스크린샷을 포함하고, 코드 예제를 변경한 경우 플랫폼 및 버전 세부 정보를 포함하세요.

프로덕션 사이트는 정적 사이트이므로 로컬 미리 보기와 검증에는 배포 자격 증명이나 번역 제공업체가 필요하지 않습니다. 검토가 끝난 뒤 게시는 유지관리자가 처리합니다.

코드 기여 및 기능 제안에 관해서는 [Wails에 기여하기](/contributing/)를 참조하세요. 프레임워크 내부 구조에 관해서는 [기술 개요](/contributing/overview/)를 참조하세요.
