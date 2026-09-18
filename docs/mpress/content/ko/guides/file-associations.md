---
title: "파일 연결"
description: "Wails 애플리케이션의 파일 연결 구성"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

관련 플랫폼: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

파일 연결을 사용하면 사용자가 특정 유형의 파일을 열 때 애플리케이션에서 해당 파일을 처리할 수 있습니다. 이는 텍스트 편집기, 이미지 뷰어 또는 특정 파일 형식을 사용하는 모든 애플리케이션에 특히 유용합니다. 이 가이드에서는 Wails v3 애플리케이션에서 파일 연결을 구현하는 방법을 설명합니다.

## 개요

현재 Wails v3에서는 다음 환경에서 파일 연결을 지원합니다:

- Windows(NSIS 설치 관리자 패키지)
- macOS(애플리케이션 번들)

## 구성

파일 연결은 프로젝트의 `build` 디렉터리에 있는 `config.yml` 파일에서 구성합니다.

### 기본 구성

파일 연결을 설정하려면 다음 단계를 따르세요:

1. `build/config.yml`을 엽니다.
2. `fileAssociations` 섹션 아래에 파일 연결을 추가합니다.
3. `wails3 update build-assets`을 실행하여 빌드 자산을 업데이트합니다.
4. 애플리케이션 옵션에서 `FileAssociations` 필드를 설정합니다.
5. `wails3 package`을 사용하여 애플리케이션을 패키징합니다.

구성 예시는 다음과 같습니다:

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### 구성 속성

| 속성 | 설명 | 플랫폼 |
| --- | --- | --- |
| ext | 맨 앞의 마침표를 제외한 파일 확장자(예: `txt`) | 모든 플랫폼 |
| name | 파일 유형의 표시 이름 | 모든 플랫폼 |
| description | 파일 속성에 표시되는 설명 | Windows |
| iconName | 빌드 폴더에 있는 아이콘 파일의 이름(확장자 제외) | 모든 플랫폼 |
| role | 이 파일 유형에 대한 애플리케이션의 역할(예: `Editor`, `Viewer`) | macOS |
| mimeType | 파일의 MIME 유형(예: `image/jpeg`) | macOS |

## 파일 열기 이벤트 수신

애플리케이션에서 파일 열기 이벤트를 처리하려면 `events.Common.ApplicationOpenedWithFile` 이벤트를 수신할 수 있습니다:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## 단계별 튜토리얼

간단한 텍스트 편집기의 파일 연결을 설정하는 과정을 단계별로 살펴보겠습니다:

@steps
### 아이콘 만들기
- 파일 유형에 사용할 아이콘을 만듭니다(권장 크기: 16x16, 32x32, 48x48, 256x256).
- 아이콘을 프로젝트의 `build` 폴더에 저장합니다.
- `iconName` 구성에 따라 아이콘 이름을 지정합니다(예: `textFileIcon.png`).

@note{type="tip"}
`wails3 generate icons`을 사용하여 필요한 아이콘을 생성할 수 있습니다. 자세한 내용을 확인하려면 `wails3 generate icons --help`을 실행하세요.

@end

- macOS의 경우 `create:app:bundle:` 작업에 `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources`과 같은 복사 구문을 추가하세요.

### 파일 연결 구성
`build/config.yml` 파일을 편집하여 파일 연결을 추가하세요:

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### 빌드 자산 업데이트
다음 명령을 실행하여 빌드 자산을 업데이트하세요:

```bash
wails3 update build-assets
```

### 애플리케이션 옵션에서 파일 연결 설정
`main.go` 파일의 애플리케이션 옵션에서 `FileAssociations` 필드를 설정하세요:

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="애플리케이션 구성과 config.yml 모두에 파일 확장자가 필요한 이유는 무엇인가요?"}
Windows에서 파일 연결을 통해 파일을 열면 파일 이름이 애플리케이션의 첫 번째 인수로 전달되어 애플리케이션이 실행됩니다. 애플리케이션은 첫 번째 인수가 파일인지 명령줄 인수인지 알 수 없으므로, 애플리케이션 옵션의 `FileAssociations` 필드를 사용하여 첫 번째 인수가 연결된 파일인지 확인합니다.

@end

### 애플리케이션 패키징
다음 명령을 사용하여 애플리케이션을 패키징하세요:

```bash
wails3 package
```

패키징된 애플리케이션은 `bin` 디렉터리에 생성됩니다. 그런 다음 애플리케이션을 설치하고 테스트할 수 있습니다.

## 추가 참고 사항

- 아이콘은 빌드 폴더에 PNG 형식으로 제공하는 것이 좋습니다.
- 파일 연결을 테스트하려면 패키징된 애플리케이션을 설치해야 합니다.

@end
