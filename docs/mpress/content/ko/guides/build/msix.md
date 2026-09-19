---
title: "MSIX 패키징"
description: "Wails v3 애플리케이션을 MSIX 패키지로 패키징하기"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX는 최신 Windows 애플리케이션 패키징 형식입니다. Wails는 Windows 빌드 과정에서 MSIX 패키지를 생성할 수 있습니다.

MSIX 패키징 방법은 [Windows 패키징](/guides/build/windows/#msix-package) 가이드에 설명되어 있습니다.

## CLI로 직접 패키지 만들기

CI를 포함하여 Windows에서 MSIX 도구를 실행하세요. 기본 백엔드는 `MakeAppx.exe`를 사용하며 서명에는 `signtool.exe`도 필요합니다. 둘 다 Windows SDK에 포함되어 있습니다. 설치 도우미는 Microsoft Store를 열고 필요한 경우 SDK 다운로드 페이지도 엽니다. 패키지를 만들기 전에 설치를 완료하세요.

식별 정보를 `build/config.yml`에 설정한 다음 실행 파일을 빌드하고 패키지로 만드세요.

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

직접 실행하는 명령은 현재 디렉터리에 `MyApp.msix` 를 생성합니다. 반면 [Windows 패키징 작업](/guides/build/windows/#msix-package) 은 자체 출력 경로를 지정합니다.

## CLI 옵션

| 옵션 | 의미 |
| --- | --- |
| `--config` | 설정 파일입니다. 기본값: `build/config.yml`. |
| `--executable`, `--name` | 기존 실행 파일과 패키지 내부에서 사용할 파일 이름입니다. 둘 다 필수입니다. |
| `--out` | 출력 파일입니다. 기본값: `<ProductName>.msix`. |
| `--arch` | 패키지 아키텍처는 `x64` (기본값), `x86`, `arm`, `arm64`, `x86a64`또는 `neutral`입니다. Go 별칭인 `amd64` 와 `386` 도 사용할 수 있습니다. 실행 파일의 아키텍처와 일치시키세요. |
| `--publisher` | 게시자 식별 정보입니다. 기본값: `CN=<companyName>`. |
| `--cert`, `--cert-password` | 서명에 사용할 PFX 인증서 경로와 암호입니다. |
| `--use-makeappx` | Windows SDK의 기본 패키징 도구를 사용합니다. |
| `--use-msix-tool` | 명시적으로 선택할 도구는 `MsixPackagingTool.exe`이며, 다음 위치에 있어야 합니다: `PATH`. |

## 서명 및 CI

Store 외부에 배포하려면 대상 컴퓨터에서 신뢰하는 인증서로 서명하세요. 인증서의 Subject는 `--publisher` 와 정확히 일치해야 합니다. MakeAppx 백엔드는 `--cert` 가 제공되면 SHA256을 사용하는 SignTool을 호출합니다. 다음 문서를 참고하세요: [Microsoft 서명 가이드](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

이 Windows 워크플로 단계는 Wails와 SDK가 설치되어 있고 앞선 단계에서 PFX 파일을 `CERT_PATH`에 안전하게 배치했다고 가정합니다. 경로가 담긴 시크릿만으로 인증서가 업로드되지는 않습니다.

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## 파일 연결 및 리소스

확장자를 `build/config.yml`에 앞의 점을 제외하고 추가하세요. 생성된 매니페스트가 점을 추가합니다.

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

실행 중 파일 열기는 [파일 연결](/guides/file-associations/)에 설명된 대로 처리하세요. 현재 MakeAppx 백엔드는 실행 파일만 복사하고 투명한 자리 표시자 이미지를 생성합니다. 프로젝트의 `Assets/` 파일을 가져오거나 `iconName` 아이콘을 변환하지 않습니다. 브랜드 이미지나 추가 DLL이 필요하면 사용자 정의 패키징 워크플로를 사용하세요.

| 생성되는 리소스 | 크기(픽셀) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` 는 파일 연결이 설정된 경우에만 생성됩니다. 이 파일들은 패키지의 `Assets/` 디렉터리에 있습니다.

## Store 제출 및 문제 해결

애플리케이션을 [Partner Center 포털](https://partner.microsoft.com/dashboard) 에서 예약하고 제출을 준비할 때 해당 패키지 식별 정보와 게시자를 사용하세요. Store는 제출 과정에서 MSIX 패키지에 서명하므로 이 배포 방식에는 서명 인증서를 구매할 필요가 없습니다. 다음 문서를 참고하세요: [Microsoft 패키지 요구 사항](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

도구 `MakeAppx.exe` 또는 `signtool.exe` 를 찾을 수 없다면 Windows SDK를 설치하거나 복구하세요. Wails는 `PATH` 와 표준 SDK 위치를 검색합니다. 서명 실패 시 인증서의 Subject, 유효 기간, 대상 컴퓨터의 신뢰 여부를 확인하세요. 다음 문서를 참고하세요: [MSIX 문제 해결](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
