---
title: "첫 번째 애플리케이션"
description: "첫 번째 Wails 데스크톱 애플리케이션을 단계별로 만듭니다"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

이 가이드에서는 프로젝트 설정, 빌드, 개발 워크플로를 포함하여 첫 번째 Wails v3 애플리케이션을 만드는 방법을 설명합니다.

<br/>

<br/>

@steps
### 새 프로젝트 만들기
터미널을 열고 다음 명령을 실행하여 새 Wails 프로젝트를 만드세요.

```bash
wails3 init -n myfirstapp
```

이 명령은 필요한 모든 파일이 포함된 `myfirstapp` 디렉터리를 새로 만듭니다.

   <video src="/assets/wails_init.mp4" controls></video>

### 프로젝트 구조 살펴보기
`myfirstapp` 디렉터리로 이동하세요. 다음과 같은 여러 파일과 폴더가 있습니다.

@filetree
- build/           빌드 과정에서 사용하는 파일 포함
  - appicon.png  애플리케이션 아이콘
  - config.yml   빌드 구성
  - Taskfile.yml Build tasks
  - darwin/      macOS 전용 빌드 파일
    - Info.dev.plist Development configuration
    - Info.plist    프로덕션 구성
    - Taskfile.yml  macOS 빌드 작업
    - icons.icns    macOS 애플리케이션 아이콘
  - linux/       Linux 전용 빌드 파일
    - Taskfile.yml  Linux 빌드 작업
    - appimage/     AppImage 패키징
      - build.sh  AppImage 빌드 스크립트
    - nfpm/        NFPM 패키징
      - nfpm.yaml Package configuration
      - scripts/  빌드 스크립트
  - windows/     Windows 전용 빌드 파일
    - Taskfile.yml        Windows 빌드 작업
    - icon.ico           Windows 애플리케이션 아이콘
    - info.json          애플리케이션 메타데이터
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS 설치 프로그램 파일
      - project.nsi                    NSIS 프로젝트 파일
      - wails_tools.nsh               NSIS 도우미 스크립트
- frontend/        프런트엔드 애플리케이션 파일
  - index.html   기본 HTML 파일
  - main.js      기본 JavaScript 파일
  - package.json NPM package configuration
  - public/      정적 자산
  - Inter Font License.txt Font license
- .gitignore      Git 무시 파일
- README.md       프로젝트 문서
- Taskfile.yml    프로젝트 작업
- go.mod          Go 모듈 파일
- go.sum          Go 모듈 체크섬
- greetservice.go Greeting service
- main.go         애플리케이션 기본 코드
@end

잠시 시간을 내어 이 파일들을 살펴보고 프로젝트 구조를 익혀 보세요.

@note{type="info"}
Wails v3에서는 [Task](https://taskfile.dev/)를 기본 빌드 시스템으로 사용하지만, `make` 또는 다른 빌드 시스템을 사용해도 됩니다.

@end

### 애플리케이션 빌드하기
애플리케이션을 빌드하려면 다음을 실행하세요.

```bash
wails3 build
```

이 명령은 애플리케이션의 디버그 버전을 컴파일하고 새로 생성된 `bin` 디렉터리에 저장합니다.

@note{type="info"}
`wails3 build`는 `wails3 task build`의 축약형이며 `Taskfile.yml`에 있는 `build` 작업을 실행합니다.

@end

     <video src="/assets/wails_build.mp4" controls></video>

빌드가 완료되면 일반 애플리케이션과 같은 방식으로 실행할 수 있습니다.

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

애플리케이션의 출발점인 간단한 UI가 표시됩니다. 디버그 버전이므로 콘솔 창에는 로그도 표시됩니다. 이 로그는 디버깅할 때 유용합니다.

### 개발 모드
애플리케이션을 개발 모드로 실행할 수도 있습니다. 이 모드에서는 전체 애플리케이션을 다시 빌드하지 않고도 프런트엔드 코드를 변경하고 실행 중인 애플리케이션에 변경 사항이 반영되는 것을 확인할 수 있습니다.

1. 새 터미널 창을 여세요.
2. `wails3 dev`을 실행하세요. 애플리케이션이 컴파일된 후 디버그 모드로 실행됩니다.
3. 원하는 편집기에서 `frontend/index.html`을 여세요.
4. 코드를 편집하여 `Please enter your name below`을 `Please enter your name below!!!`(으)로 변경하세요.
5. 파일을 저장하세요.

이 변경 사항은 애플리케이션에 즉시 반영됩니다.

백엔드 코드를 변경하면 다시 빌드됩니다.

1. `greetservice.go`을 여세요.
2. `return "Hello " + name + "!"`이 있는 줄을 `return "Hello there " + name + "!"`(으)로 변경하세요.
3. 파일을 저장하세요.

몇 초 안에 애플리케이션이 업데이트됩니다.

     <video src="/assets/wails_dev.mp4" controls></video>

### 애플리케이션 패키징하기
애플리케이션을 배포할 준비가 되면 플랫폼별     패키지를 만들 수 있습니다.

@tabs{sync-key="platform"}
[Mac]
`.app` 번들을 만들려면 다음을 실행하세요.

```bash
wails3 package
```

프로덕션 빌드를 생성하고 `bin` 디렉터리에 `.app` 번들로 패키징합니다.

[Windows]
NSIS 설치 프로그램을 만들려면 다음을 실행하세요.

```bash
wails3 package
```

프로덕션 빌드를 생성하고 `bin` 디렉터리에 NSIS 설치 프로그램으로 패키징합니다.

[Linux]
Wails는 Linux 배포를 위한 여러 패키지 형식을 지원합니다.

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

패키징 옵션과 구성에 관한 자세한 내용은     [빌드 및 패키징 가이드](/guides/build/building/)를 참조하세요.

### 버전 관리 및 모듈 이름 설정하기
프로젝트는 자리표시자 모듈 이름인 `changeme`(으)로 생성됩니다. 저장소 URL과 일치하도록 이 이름을 변경하는 것이 좋습니다.

1. GitHub 또는 원하는 Git 호스트에서 새 저장소를 만드세요.
2. 프로젝트 디렉터리에서 git을 초기화하세요.
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. 원격 저장소를 설정하세요(실제 저장소 URL로 바꾸세요).
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. `go.mod`의 모듈 이름을 저장소 URL과 일치하도록 변경하세요.
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. 코드를 푸시하세요.
  ```bash
  git push -u origin main
  ```


이렇게 하면 Go 모듈 이름이 Go의 모듈 명명 규칙에 맞게 설정되어 코드를 더 쉽게 공유할 수 있습니다.

@note{type="tip" title="유용한 팁"}
프로젝트를 만들 때 `-git` 플래그를 사용하면 모든 초기화 단계를 자동화할 수 있습니다.

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

다양한 Git URL 형식을 지원합니다.

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` 또는 `ssh://git@github.com/username/project`
- Git 프로토콜: `git://github.com/username/project`
- 파일 시스템: `file:///path/to/project.git`

@end

@end

## 축하합니다!

방금 첫 번째 Wails 애플리케이션을 만들고 개발하여 패키징까지 마쳤습니다. 이는 Wails v3로 이룰 수 있는 일의 시작에 불과합니다.

## 다음 단계

Wails를 처음 사용한다면, 다음으로 Wails의 다양한 기능을 실습하며 익힐 수 있는 튜토리얼을 읽어보시기 바랍니다. 첫 번째 튜토리얼은 [서비스 만들기](/tutorials/01-creating-a-service/)입니다.

숙련된 사용자라면 Wails 사용 방법에 관한 더 자세한 내용은 [빌드 및 패키징 가이드](/guides/build/building/)에서 확인하세요.
