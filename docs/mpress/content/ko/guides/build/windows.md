---
title: "Windows 패키징"
description: "Windows 배포용 Wails 애플리케이션 패키징"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## NSIS 설치 프로그램

기본 패키징 형식은 NSIS 설치 프로그램을 생성합니다:

```bash
wails3 package GOOS=windows
```

이 명령은 `wails3 task windows:package`을 실행하여 다음 작업을 수행합니다:

1. 애플리케이션 빌드
2. WebView2 부트스트래퍼 생성
3. NSIS 설치 프로그램 생성

출력: `build/windows/nsis/<AppName>-installer.exe`

### MSIX 패키지

Microsoft Store에 배포하거나 최신 Windows 배포 방식을 사용하려면 다음을 실행합니다:

```bash
wails3 package GOOS=windows FORMAT=msix
```

출력: `bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX를 사용하려면 `makeappx.exe`(Windows SDK) 또는 독립 실행형 MSIX 도구가 필요합니다. Windows Taskfile에서는 설치 작업을 `wails3 task install:msix:tools`로 제공합니다.

@end

## 설치 프로그램 사용자 지정

NSIS 구성은 `build/windows/nsis/project.nsi`에 있습니다. 다음 항목을 사용자 지정하려면 이 파일을 편집하세요:

- 설치 프로그램 UI 및 브랜딩
- 설치 디렉터리
- 시작 메뉴 및 바탕 화면 바로 가기
- 파일 연결
- 사용권 계약

애플리케이션 메타데이터는 `build/windows/info.json`에서 가져옵니다:

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## 코드 서명

SmartScreen 경고가 표시되지 않도록 실행 파일과 설치 프로그램에 서명하세요:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

`build/windows/Taskfile.yml`에서 서명을 구성하세요:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

인증서 암호를 안전하게 보관하세요:

```bash
wails3 setup signing
```

자세한 내용은 [애플리케이션 서명](/guides/build/signing/)을 참조하세요.

## ARM용 빌드

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## 문제 해결

### makensis를 찾을 수 없음

NSIS를 설치하세요:

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### SmartScreen 경고

실행 파일이 서명되지 않았습니다. 위의 [코드 서명](#-)을 참조하세요.

### WebView2 누락

설치 프로그램에는 필요할 때 런타임을 다운로드하는 WebView2 부트스트래퍼가 포함되어 있습니다. 오프라인 설치가 필요하면 Microsoft에서 Evergreen Standalone Installer를 다운로드하세요.
