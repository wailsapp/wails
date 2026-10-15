---
title: "코드 서명"
description: "모든 플랫폼에서 Wails 애플리케이션에 서명하는 방법을 설명하는 가이드"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## 애플리케이션 코드 서명

이 가이드에서는 macOS, Windows 및 Linux용 Wails 애플리케이션에 서명하는 방법을 설명합니다. Wails v3는 코드 서명, 공증 및 PGP 키 관리를 위한 CLI 도구를 기본 제공합니다.

- **macOS** - macOS 애플리케이션 서명 및 공증
- **Windows** - Windows 실행 파일 및 패키지 서명
- **Linux** - PGP 키로 DEB 및 RPM 패키지 서명

## 크로스 플랫폼 서명 매트릭스

이 매트릭스는 각 소스 플랫폼에서 서명할 수 있는 대상을 보여 줍니다.

| 대상 형식 | Windows에서 | macOS에서 | Linux에서 |
| --- | :---: | :---: | :---: |
| Windows EXE/MSI | ✅ | ✅ | ✅ |
| macOS .app 번들 | ❌ | ✅ | ❌ |
| macOS 공증 | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Windows 및 Linux 패키지는 <strong>모든 플랫폼</strong>에서 서명할 수 있습니다. macOS 서명에는 Apple 도구가 필요하므로 Mac이 있어야 합니다.

@end

### 서명 백엔드

Wails는 사용 가능한 최적의 서명 백엔드를 자동으로 선택합니다.

| 플랫폼 | 네이티브 백엔드 | 크로스 플랫폼 백엔드 |
| --- | --- | --- |
| Windows | `signtool.exe`(Windows SDK) | 기본 제공 |
| macOS | `codesign`(Xcode) | 사용할 수 없음 |
| Linux | 해당 없음 | 기본 제공 |

네이티브 플랫폼에서 실행할 때 Wails는 호환성을 극대화하기 위해 네이티브 도구를 사용합니다. 크로스 컴파일할 때는 기본 제공 서명 기능을 사용합니다.

## 빠른 시작

서명을 가장 빠르게 구성하려면 설정 마법사를 사용합니다.

```bash
wails3 setup
```

이 작업은 모든 플랫폼에서 **서명 시 적용되는** 공유 서명 구성을 `~/.config/wails/defaults.yaml`에 기록합니다([구성 우선순위](#--2) 참조). 서명 단계에서는 다음 작업을 수행합니다.

- **어떤 호스트에서든** 대상 플랫폼용 서명 도구 중 로컬에 설치된 도구를 감지하고, 사용 중인 OS에 맞는 설치 명령을 표시합니다(예: `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- macOS에서는 키체인에 있는 Developer ID 인증서를 나열합니다.
- Linux를 대상으로 하는 경우 GPG 키를 나열하고 새 키를 **생성하고 내보낼** 수 있습니다.
- Windows를 대상으로 하는 경우 OpenSSL을 통해 테스트용 <strong>자체 서명 인증서를 생성</strong>할 수 있습니다.
- 암호를 Taskfile이 아닌 **시스템 키체인에 안전하게 저장합니다**.

@note{type="tip"}
암호는 시스템의 네이티브 자격 증명 저장소(macOS Keychain, Windows Credential Manager 또는 Linux Secret Service)에 저장됩니다. 따라서 서명 구성이 안전하게 보호되며 모든 Wails 프로젝트에서 작동합니다.

@end

## 구성 우선순위

서명 작업을 실행하면 각 서명 옵션은 다음 순서로 결정됩니다(처음 일치하는 값이 적용됨).

1. `wails3 tool sign`에 명시적으로 전달한 플래그(예: `--pgp-key`, `--certificate`, `--identity`).
2. 일치하는 **프로젝트 Taskfile 변수**(`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY` 등).
3. `~/.config/wails/defaults.yaml`에 있는 **전역 구성**(`wails3 setup`에서 기록).

즉, Taskfile 변수는 <strong>선택적 재정의 값</strong>입니다. 변수가 설정되지 않으면 전역으로 구성된 키/인증서/ID를 사용합니다. 세 가지 소스 중 어느 곳에도 값이 없으면 서명 명령에서 구성 방법을 알려 주는 명확한 오류를 보고합니다.

## 프로젝트별 구성

전역이 아닌 단일 프로젝트에 서명을 구성하려면 프로젝트 안에서 프로젝트별 마법사를 실행합니다. 그러면 해당 프로젝트의 `build/<platform>/Taskfile.yml` 파일에 `vars`을 기록합니다.

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### 수동 구성

또는 플랫폼별 Taskfile을 직접 편집할 수 있습니다. 각 파일 위쪽의 `vars` 섹션을 편집합니다.

@tabs
[macOS]
`build/darwin/Taskfile.yml`을 편집합니다.

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

그런 다음 다음을 실행합니다.

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
`build/windows/Taskfile.yml`을 편집합니다.

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

암호는 시스템 키체인에서 가져옵니다(구성하려면 `wails3 setup signing`을 실행하십시오).

그런 다음 다음을 실행하세요:

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
`build/linux/Taskfile.yml`을 편집하세요:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

비밀번호는 시스템 키체인에서 가져옵니다(구성하려면 `wails3 setup signing`을 실행하세요).

그런 다음 다음을 실행하세요:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

시스템에서 직접 서명 상태를 확인할 수도 있습니다:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

모든 플랫폼의 서명 구성을 대화형으로 설정하려면 마법사를 사용하세요:

```bash
wails3 setup signing
```

## macOS 코드 서명

### 사전 요구 사항

- Apple Developer 계정(연간 $99)
- Developer ID Application 인증서
- Xcode Command Line Tools 설치

### 서명 ID

사용 가능한 서명 ID를 확인하세요:

```bash
security find-identity -v -p codesigning
```

출력:

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
App Store 외부에 배포하려면 **Developer ID Application** 인증서가 필요합니다.

@end

### 구성

`build/darwin/Taskfile.yml`을 편집하여 서명 변수를 설정하세요:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| 변수 | 필수 여부 | 설명 |
| --- | --- | --- |
| `SIGN_IDENTITY` | 예 | Developer ID(예: "Developer ID Application: Your Company (TEAMID)") |
| `KEYCHAIN_PROFILE` | 공증용 | 자격 증명이 저장된 키체인 프로필 이름 |
| `ENTITLEMENTS` | 아니요 | 엔타이틀먼트 파일 경로 |

그런 다음 다음을 실행하세요:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### 엔타이틀먼트

엔타이틀먼트는 앱에서 접근할 수 있는 기능을 제어합니다. 일반적으로 Wails 앱에는 개발 환경과 프로덕션 환경에서 서로 다른 엔타이틀먼트가 필요합니다:

- **개발**: JIT, 서명되지 않은 메모리 및 디버깅 엔타이틀먼트 필요
- **프로덕션**: 최소한의 엔타이틀먼트(네트워크 접근만 허용)

두 파일을 모두 생성하려면 대화형 설정 마법사를 사용하세요:

```bash
wails3 setup entitlements
```

그러면 다음 파일이 생성됩니다:

- `build/darwin/entitlements.dev.plist` - 개발 빌드용
- `build/darwin/entitlements.plist` - 프로덕션/서명된 빌드용

**사용 가능한 프리셋:**\

| 프리셋 | 설명 |
| --- | --- |
| 개발 | JIT, 서명되지 않은 메모리, 디버깅, 네트워크 |
| 프로덕션 | 네트워크만 허용(최소 권한, 가장 안전함) |
| 둘 다 | 개발 및 프로덕션 파일을 모두 생성(권장) |
| App Store | 네트워크 및 파일 접근 권한과 함께 샌드박스 활성화 |
| 사용자 지정 | 개별 엔타이틀먼트 선택 |

@note{type="note"}
darwin Taskfile의 `run` 작업은 `entitlements.dev.plist`을 자동으로 사용합니다. `sign` 작업은 프로덕션 빌드에 `entitlements.plist`을 사용합니다.

@end

그런 다음 Taskfile 변수의 `ENTITLEMENTS`이 해당 파일을 가리키도록 설정하세요.

### 공증

Apple은 배포되는 모든 앱을 공증하도록 요구합니다.

@steps
### **자격 증명을 키체인에 저장하세요**(최초 한 번만 설정). `wails3 setup signing`을 실행하거나(필요한 값을 입력하라는 메시지를 표시하고 내부적으로 `notarytool`을 호출함) `notarytool`을 직접 호출하세요:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Taskfile의 KEYCHAIN_PROFILE을 설정하여** 위의 프로필 이름과 일치시키세요.
### **앱에 서명하고 공증하세요**:
```bash
wails3 task darwin:sign:notarize
```

### **공증을 확인하세요**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
공증에는 일반적으로 1-2분이 걸립니다. 티켓은 앱에 자동으로 스테이플됩니다.

@end

## Windows 코드 서명

### 사전 요구 사항

- 코드 서명 인증서(DigiCert, Sectigo 등에서 발급)
- Windows에서 네이티브 서명을 사용하는 경우: Windows SDK 설치(`signtool.exe`용)
- macOS/Linux에서 크로스 플랫폼 서명을 수행하려면 [`osslsigncode`](https://github.com/mtrojnar/osslsigncode)이 필요합니다(서명 단계의 `wails3 setup`에 호스트별 설치 명령이 표시됩니다).

### 자체 서명 인증서 생성(테스트용)

서명 파이프라인만 테스트하려면 `wails3 setup`을 실행하고 서명 단계에서 **Windows** 탭을 연 다음 <strong>자체 서명 인증서 생성</strong>을 선택하세요. 그러면 OpenSSL을 사용하여 코드 서명용 `.pfx`을 생성하고 전역 구성에 해당 경로를 기록합니다.

@note{type="caution"}
자체 서명 인증서는 **테스트 및 내부 배포 용도로만** 사용해야 합니다. 최종 사용자에게 SmartScreen 경고가 표시되기 때문입니다. 공개 릴리스에는 신뢰할 수 있는 CA에서 발급한 인증서가 필요합니다.

@end

### 구성

`build/windows/Taskfile.yml`을 편집하여 서명 변수를 설정하세요.

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| 변수 | 필수 여부 | 설명 |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | 재정의 | .pfx/.p12 인증서 파일 경로(설정하지 않으면 전역 구성 사용) |
| `SIGN_THUMBPRINT` | 재정의 | Windows 인증서 저장소에 있는 인증서 지문(`SIGN_CERTIFICATE`의 대안) |
| `TIMESTAMP_SERVER` | 아니요 | 타임스탬프 서버 URL(기본값: http://timestamp.digicert.com) |

@note{type="note"}
이 변수들은 선택적 <strong>재정의</strong>입니다. 설정하지 않으면 `wails3 setup`을 통해 전역으로 구성한 인증서를 사용합니다. [구성 우선순위](#--2)를 참조하세요.

@end

@note{type="note"}
인증서 비밀번호는 Taskfile이 아니라 시스템 키체인에 저장됩니다. `wails3 setup signing`을 실행하여 구성하거나 CI에서 `WAILS_WINDOWS_CERT_PASSWORD` 환경 변수를 설정하세요.

@end

그런 다음 다음 명령을 실행하세요.

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### 크로스 플랫폼 서명

Windows 실행 파일은 어느 플랫폼에서든 서명할 수 있습니다. 동일한 Taskfile 구성과 명령을 macOS 및 Linux에서도 사용할 수 있습니다.

### 지원되는 Windows 형식

| 형식 | 확장자 | 참고 |
| --- | --- | --- |
| 실행 파일 | .exe | 표준 PE 서명 |
| 설치 프로그램 | .msi | Windows Installer 패키지 |
| 앱 패키지 | .msix, .appx | 최신 Windows 앱 |

## Linux 패키지 서명

Linux 패키지(DEB 및 RPM)는 PGP/GPG 키를 사용하여 서명합니다. Windows 및 macOS 코드 서명과 달리 Linux 패키지 서명은 OS가 코드를 신뢰한다는 것이 아니라 패키지가 신뢰할 수 있는 출처에서 제공되었음을 입증합니다.

### 사전 요구 사항

- PGP 키 쌍(Wails로 생성 가능)

### PGP 키 생성

가장 쉬운 방법은 설정 마법사를 사용하는 것입니다. `wails3 setup`을 실행하고 서명 단계에서 **Linux** 탭을 연 다음 <strong>새 GPG 키 생성</strong>을 선택하세요. 마법사는 다음 작업을 수행합니다.

- GPG 키링에 RSA 4096 키를 생성합니다(무인 실행 및 CI에 적합한 키로 만들려면 암호를 비워 두세요).
- 키를 <strong>`~/.wails/signing/<keyid>.asc`</strong>로 내보냅니다(빌드에서 서명에 사용하는 파일).
- 서명할 때 자동으로 사용되도록 키 ID와 내보낸 파일 경로를 모두 `~/.config/wails/defaults.yaml`에 기록합니다.

@note{type="note"}
키가 **ID만으로** 구성된 경우(예: 마법사가 키를 자동으로 내보내는 기능이 도입되기 전에 생성한 경우), 다음에 Linux 탭을 열면 키를 파일로 내보내고 경로를 입력합니다. 수동 작업은 필요하지 않습니다. 빌드는 키 <em>파일</em>을 사용하여 서명하므로 중요한 것은 경로입니다.

@end

`gpg`을 사용하여 수동으로 수행할 수도 있습니다.

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

권장 설정: RSA 4096비트, 만료 기간 1년, 강력한 비밀번호로 보호.

@note{type="caution"}
개인 키를 안전하게 보관하세요! 암호화하여 저장하고 안전하게 백업하세요.

@end

### 구성

`build/linux/Taskfile.yml`을 편집하여 서명 변수를 설정하세요.

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| 변수 | 필수 여부 | 설명 |
| --- | --- | --- |
| `PGP_KEY` | 재정의 | 내보낸 PGP 비공개 키 파일의 경로(설정하지 않으면 `wails3 setup`을 통해 전역으로 구성한 키 사용) |
| `SIGN_ROLE` | 아니요 | DEB 서명 역할(기본값: builder) |

@note{type="note"}
PGP 키 암호는 Taskfile이 아닌 시스템 키체인에 저장됩니다. 암호를 구성하려면 `wails3 setup signing`을 실행하거나 CI에서 `WAILS_PGP_PASSWORD` 환경 변수를 설정하세요.

@end

그런 다음 다음을 실행하세요:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### DEB 서명 역할

DEB 패키지의 경우 `SIGN_ROLE`을 통해 서명 역할을 지정할 수 있습니다:

- `origin`: 패키지 출처의 서명
- `maint`: 패키지 유지관리자의 서명
- `archive`: 아카이브 유지관리자의 서명
- `builder`: 패키지 빌더의 서명(기본값)

### 크로스 플랫폼 서명

Linux 패키지는 어떤 플랫폼에서든 서명할 수 있습니다. 동일한 Taskfile 구성과 명령을 Windows 및 macOS에서도 사용할 수 있습니다.

### 키 정보 보기

```bash
gpg --show-keys signing-key.asc
```

출력:

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Linux 패키지 검증

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### 공개 키 배포

사용자가 패키지를 검증하려면 공개 키가 필요합니다:

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## GitHub Actions 통합

CI 환경에서는 시스템 키체인 대신 환경 변수를 통해 암호를 제공합니다:

| 환경 변수 | 설명 |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Windows 인증서 암호 |
| `WAILS_PGP_PASSWORD` | Linux 패키지용 PGP 키 암호 |

Taskfile 변수를 직접 전달할 수도 있습니다:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### macOS 워크플로

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Windows 워크플로

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### 크로스 플랫폼 워크플로(Linux 러너)

단일 Linux 러너에서 Windows 및 Linux 패키지에 서명하세요:

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
크로스 플랫폼 서명에 Linux 러너를 사용하면 별도의 Windows 러너가 필요하지 않으므로 CI/CD가 간소화됩니다. macOS 서명에는 Apple의 네이티브 도구가 필요하므로 여전히 macOS 러너를 사용해야 합니다.

@end

## CLI 참조

### wails3 setup signing

프로젝트의 서명을 구성하는 대화형 마법사입니다.

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

마법사에서 다음 작업을 안내합니다:

- **macOS**: Developer ID 인증서를 선택하고 공증 자격 증명을 구성합니다(`xcrun notarytool store-credentials` 호출).
- **Windows**: 인증서 파일과 지문 중 하나를 선택하고 암호와 타임스탬프 서버를 설정합니다.
- **Linux**: 기존 PGP 키를 사용하거나 새 키를 생성하고(`gpg` 호출) 서명 역할을 구성합니다.

### wails3 setup entitlements

macOS 권한을 구성하는 대화형 마법사입니다.

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**프리셋:**

- **개발**: JIT, 디버깅 및 네트워크 권한이 포함된 `entitlements.dev.plist`을 생성합니다.
- **프로덕션**: 최소한의 권한이 포함된 `entitlements.plist`을 생성합니다.
- **둘 다**: 두 파일을 모두 생성합니다(권장).
- **App Store**: Mac App Store용 샌드박스 권한을 생성합니다.
- **사용자 지정**: 개별 권한과 대상 파일을 선택합니다.

### wails3 sign

현재 플랫폼 또는 지정한 플랫폼용 바이너리와 패키지에 서명합니다. 적절한 플랫폼별 서명 작업을 호출하는 래퍼입니다.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Taskfile의 서명 구성을 사용하는 해당 `<platform>:sign` 작업을 실행합니다.

### wails3 tool sign

특정 파일에 직접 서명하는 저수준 명령입니다. Taskfile에서 내부적으로 사용됩니다.

```bash
wails3 tool sign [flags]
```

**공통 플래그:**\

| 플래그 | 설명 |
| --- | --- |
| `--input` | 서명할 파일의 경로 |
| `--output` | 출력 경로(선택 사항, 기본적으로 원본 위치에 저장) |
| `--verbose` | 상세 출력 활성화 |

**Windows/macOS 플래그:**\

| 플래그 | 설명 |
| --- | --- |
| `--certificate` | PKCS#12 인증서(.pfx/.p12)의 경로 |
| `--password` | 인증서 비밀번호 |
| `--timestamp` | 타임스탬프 서버 URL |

**macOS 전용 플래그:**\

| 플래그 | 설명 |
| --- | --- |
| `--identity` | 서명 ID(임시 서명에는 '-' 사용) |
| `--entitlements` | 권한 plist 파일의 경로 |
| `--hardened-runtime` | 강화된 런타임 활성화(기본값: true) |
| `--notarize` | 공증을 위해 제출 |
| `--keychain-profile` | 공증용 키체인 프로필 |

**Windows 전용 플래그:**\

| 플래그 | 설명 |
| --- | --- |
| `--thumbprint` | Windows 인증서 저장소의 인증서 지문 |

**Linux 전용 플래그:**\

| 플래그 | 설명 |
| --- | --- |
| `--pgp-key` | PGP 개인 키의 경로 |
| `--pgp-password` | PGP 키 비밀번호 |
| `--role` | DEB 서명 역할(origin/maint/archive/builder) |

### 서명 상태 확인(네이티브 도구)

v3에는 **`wails3 signing`** 명령이 없습니다. 서명 상태를 확인하려면 네이티브 도구를 직접 사용하세요:

| 작업 | 명령 |
| --- | --- |
| macOS 코드 서명 ID 나열 | `security find-identity -v -p codesigning` |
| 공증 자격 증명 저장 | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| PGP 키 파일 확인 | `gpg --show-keys <key.asc>` |
| PGP 키 쌍 생성 | `gpg --full-generate-key` |
| 공개 키 내보내기 | `gpg --armor --export <email>` |

## 문제 해결

### macOS 문제

**"Developer ID 인증서를 찾을 수 없음"**

- 인증서가 키체인에 설치되어 있는지 확인하세요
- `security find-identity -v -p codesigning` 명령으로 인증서가 만료되지 않았는지 확인하세요
- "Apple Development" 인증서만 있는 것이 아니라 "Developer ID Application" 인증서가 있는지 확인하세요

**"공증 실패"**

- 공증 로그를 확인하세요: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- 강화된 런타임이 활성화되어 있는지 확인하세요
- 앱에 서명되지 않은 바이너리가 포함되어 있지 않은지 확인하세요

**"코드 서명 실패"**

- 키체인이 잠금 해제되어 있는지 확인하세요: `security unlock-keychain`
- 앱 번들의 파일 권한을 확인하세요

### Windows 문제

**"인증서를 찾을 수 없음"**

- 인증서 경로가 올바른지 확인하세요
- 인증서 비밀번호를 확인하세요
- 인증서가 유효한지 확인하세요(만료되거나 해지되지 않았는지 확인)

**"타임스탬프 서버 오류"**

- 다른 타임스탬프 서버를 사용해 보세요:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Linux 문제

**"유효하지 않은 PGP 키"**

- 키 파일이 ASCII-armored 형식인지 확인하세요
- `gpg --show-keys <key.asc>`을 사용하여 키가 만료되지 않았는지 확인하세요
- 비밀번호가 올바른지 확인하세요

**"서명 검증 실패"**

- 공개 키를 올바르게 가져왔는지 확인하세요
- 서명 후 패키지가 수정되지 않았는지 확인하세요

## 추가 자료

### 공식 문서

- [Apple 코드 서명 가이드](https://developer.apple.com/support/code-signing/)
- [Apple 공증 문서](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Microsoft 코드 서명](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Debian 패키지 서명](https://wiki.debian.org/SecureApt)
- [RPM 패키지 서명](https://rpm-software-management.github.io/rpm/manual/signatures.html)
