---
title: "Wails v3.1.0: Go는 더 많이, 빌드 장벽은 더 낮게"
description: "실험적인 CGO 없는 데스크톱 빌드와 Wake: Wails 애플리케이션을 더 간편하게 빌드, 패키징 및 서명하는 방법입니다."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

이번 주 Wails에서는 CGO 없는 빌드, 더 쉬운 프로젝트 설정, 패키징 및 서명 과정의 불편 감소 등 많은 변화가 있었습니다.

Wails v3.1.0부터 매주 기능 릴리스를 제공하여 새로운 기능을 정기적으로 선보입니다. 버그 수정은 준비되는 대로 nightly 패치 업데이트를 통해 계속 제공됩니다.

## 이번 주 수정 및 최적화

이번 주에는 프레임워크, CLI, 빌드 도구 및 CI 전반에 걸쳐 **100건의 수정, 21건의 성능 개선, 보안 강화 변경 네 건**이 적용되었습니다. 주요 변경 사항으로는 혼합 DPI 환경에서 더 안정적인 창 배치, Linux 트레이 및 다크 모드 수정, 더 효율적인 애셋 제공, MSIX 패키징 중 메모리 사용량 감소 등이 있습니다. [아래에서 주요 변경 사항 몇 가지를 확인해 보세요](#------).

## 실험적인 CGO 없는 빌드

purego를 사용하면 Go에서 C 코드로 연결하는 일반적인 브리지인 CGo 없이 네이티브 시스템 라이브러리를 호출할 수 있습니다. purego가 처음 등장했을 때부터 검토하기 시작했지만, 두 백엔드를 유지 관리하기는 쉽지 않았습니다. 이제 purego가 성숙하고 생성형 AI 덕분에 유지 관리 부담도 줄어들어, 실험적인 macOS 및 Linux 백엔드를 공개할 준비가 되었습니다.

이제 C 컴파일러 없이 데스크톱 앱을 빌드하고 크로스 컴파일할 수 있습니다.

```sh
WAILS_PUREGO=1 wails3 build
```

프런트엔드 애셋을 빌드한 후에는 Go로 직접 빌드할 수도 있습니다.

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

macOS 백엔드는 AppKit과 WebKit을 사용하며, Linux에서는 GTK4/WebKitGTK 6.0 및 GTK3를 지원합니다. 사용자에게는 여전히 네이티브 런타임 라이브러리가 필요하며, 패키징과 서명에는 별도의 요구 사항이 있습니다. 피드백을 수집하는 동안에는 CGo가 계속 기본값으로 유지됩니다.

[구현 및 플랫폼 세부 정보](https://github.com/wailsapp/wails/pull/6521)

## 왜 Wake인가요?

Wails v2에서는 빌드 과정이 지나치게 많이 감춰져 있었습니다. v3에서는 Taskfile을 통해 모든 과정을 공개했지만, 대부분의 사용자는 빌드 시스템을 직접 유지 관리하지 않고도 몇 가지 항목만 맞춤 설정하기를 원합니다.

Wake는 애플리케이션, 보조 CLI 도구, 서버, 대상 및 패키징과 같은 일반적인 선택 사항을 `wails.yaml`에 모아 둡니다. Task는 계속 핵심 빌드 시스템으로 사용되며, 유연성이 필요할 때 이용할 수 있습니다.

이제 Wake는 라이브 리로드, 증분 빌드, 출력 캐싱 및 병렬 패키징을 지원합니다. 바이너리가 이미 있나요? Go나 프런트엔드 도구를 설치하지 않아도 `--skip-build`에서 패키징할 수 있습니다.

`WAILS_EXP_USE_WAKE=1`을 사용해 활성화하세요.

### 안내에 따라 시작하기

Wake를 활성화하면 `wails3 init -ui`에서 백엔드, 대상 및 패키지 선택 과정을 단계별로 안내하고, 도구를 검사한 뒤 프로젝트 구성을 저장합니다.

@image{src="/assets/screenshots/wake/application.png" alt="실험적인 Pure Go를 포함하여 데스크톱 애플리케이션과 백엔드 선택 항목을 표시하는 Wake 프로젝트 마법사." expand}

각 대상에 사용할 패키지 형식과 서명 옵션을 선택하세요.

@image{src="/assets/screenshots/wake/packages.png" alt="AppImage, DEB, RPM, Void 및 Arch Linux 형식과 릴리스 출력 서명 옵션을 표시하는 Linux용 Wake 패키지 선택 화면." expand}

아래 도구 검사 화면은 표준 Linux 빌드를 보여 줍니다.

@image{src="/assets/screenshots/wake/tools.png" alt="표준 Linux 빌드에 필요한 C 컴파일러, Go, npm 및 pkg-config를 검사하는 Wake." expand}

### 내 컴퓨터에서 무엇을 빌드할 수 있나요?

컴파일, 패키징 및 서명에는 각각 다른 요구 사항이 있습니다. 각 호스트에서 수행할 수 있는 작업은 다음과 같습니다.

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| 작업 | Windows 호스트 | macOS 호스트 | Linux 호스트 | 필요한 항목 |
| purego로 데스크톱 앱 컴파일 | ✓ | ✓ | ✓ | Go 및 macOS/Linux 대상용 실험적 purego 백엔드 |
| MSIX 패키지 생성 | ✓ | ✓ | ✓ | 빌드된 애플리케이션. Windows SDK는 필요하지 않음 |
| NSIS 설치 프로그램 생성 | ✓ | ✓ | ✓ | 빌드된 애플리케이션 및 `makensis` |
| macOS 앱 및 DMG 생성 | ✓ | ✓ | ✓ | 빌드된 애플리케이션. Go에서 패키지 생성을 처리함 |
| DEB, RPM, Arch 및 Void 패키지 생성 | ✓ | ✓ | ✓ | 빌드된 애플리케이션 |
| AppImage 생성 | ✓ | ✓ | ✓ | 대상 AppImage 런타임. 크로스 패키징에는 Linux GUI 런타임 트리가 필요함 |
| 내보낸 ID로 Windows 패키지 서명 | ✓ | ✓ | ✓ | 내보낸 RSA PFX 인증서 및 비공개 키 |
| 내보낸 ID로 Apple 앱 및 DMG 서명 | ✓ | ✓ | ✓ | 내보낸 RSA PKCS#12 ID. 아래 검증 참고 사항 참조 |
| DEB, RPM 및 Arch 패키지 서명 | ✓ | ✓ | ✓ | 내보낸 OpenPGP 키 |
| Apple 공증 제출 및 티켓 첨부 | ✓ | ✓ | ✓ | Apple 팀 API 키 및 적절하게 서명된 애플리케이션 |
| Apple의 `codesign` 또는 키체인 프로필 사용 | – | ✓ | – | macOS 및 적절한 서명 ID |
@end

CGo 빌드에는 여전히 C 도구 체인이 필요합니다. DMG를 수정하는 사용자 지정 훅에는 macOS와 `hdiutil`이 필요합니다.

### 서명 및 공증

Wake는 임시 Apple 서명을 포함하여 Windows, Apple 및 Linux 패키지 서명을 지원합니다. 크로스 플랫폼 서명에는 내보낸 ID를 사용하며, 인증서 저장소 및 키체인 ID에는 여전히 네이티브 도구가 필요합니다. 사용자 수준 자격 증명을 사용하면 비공개 키를 프로젝트 구성에 포함하지 않아도 됩니다. AppImage 및 Void 서명은 아직 지원되지 않습니다.

Apple 공증은 팀 API 키를 사용해 Go에서 처리됩니다. Linux에서 제출하고 티켓을 첨부한 결과를 Gatekeeper가 승인하는 것까지 테스트했습니다. 해당 앱은 Apple의 [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html)을 사용해 서명했습니다. Wake가 Go로 생성한 Developer ID 서명을 Apple이 승인하는지는 아직 검증되지 않았습니다.

### Wake로 이전하기

마이그레이션은 기본 제공 템플릿 아홉 개를 모두 지원합니다. 이 과정은 되돌릴 수 없습니다. `wails.yaml`을 활성화하면 Taskfile 라우팅이 대체됩니다. 지원되지 않는 사용자 지정 명령과 재정의는 검토할 수 있도록 거부됩니다.

[Wake 기능, 요구 사항 및 마이그레이션 세부 정보](https://github.com/wailsapp/wails/pull/6531)

## 이번 주 주요 변경 사항 몇 가지

- **혼합 DPI 디스플레이에서도 창이 올바른 위치에 배치됩니다.** 이제 디스플레이 배율이 서로 다른 macOS 및 Linux 구성에서도 화면 경계와 위치가 올바르게 처리됩니다. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **Linux 트레이 및 테마 감지가 더 안정적으로 작동합니다.** 앱 시작 후 생성된 트레이가 더 이상 누락되지 않고 툴팁이 채워지며, 다크 모드 감지는 시작 시점부터 데스크톱의 화면 표시 설정을 읽습니다.
- **Go 이벤트 콜백은 리스너별로 순서를 유지합니다.** 각 리스너는 이벤트가 발생한 순서대로 직렬 실행됩니다. 서로 다른 리스너는 동시에 실행될 수 있으며, 동시에 발생한 이벤트의 순서는 정의되지 않습니다. 데이터 경합을 방지하기 위해 잠금 상태에서 이벤트 훅의 스냅샷을 생성합니다. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **빌드 캐싱이 실제 프로젝트의 더 다양한 사례를 처리합니다.** 이제 따옴표로 감싸거나 이스케이프한 공백 때문에 Go 출력 캐시가 누락되지 않으며, 파일 손상을 방지하도록 동시 캐시 저장을 직렬화합니다. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **애셋 제공 과정에서 불필요한 작업이 줄었습니다.** 알려진 파일 확장자는 MIME 유형 감지를 건너뛰며, 디버그 로깅이 비활성화된 경우 더 이상 폐기될 메시지를 미리 준비하지 않습니다. 두 변경 사항 모두 요청 처리 경로의 작업량을 줄입니다. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **MSIX 패키징에서 바이너리 크기만 한 버퍼를 사용하지 않습니다.** 이제 실행 파일 전체를 메모리로 읽는 대신 패키지에 스트리밍합니다. 이전에는 대표적인 16 MiB 파일을 복사할 때 작업당 약 16.8 MB가 할당되었습니다. [#6472](https://github.com/wailsapp/wails/pull/6472)

## 직접 사용해 보세요

두 기능 모두 실험적입니다. 크로스 플랫폼 빌드를 시도하고 앱을 패키징한 뒤, 어느 부분에서 불편을 겪었는지 알려 주세요. 테스트하고 문제를 보고하며 기여해 주시는 모든 분께 감사드립니다.

저희 Discord 커뮤니티의 beesophie에게 특별히 감사드립니다. beesophie는 놀라운 빌드 시스템인 [Cook](https://github.com/SimTXNG/cook)을 만들어 저희 빌드 시스템의 몇 가지 부족한 부분을 깨닫게 해 주었습니다. Cook이 보여 준 가능성 덕분에 Wails 빌드를 더 쉽게 개선할 부분을 파악할 수 있었으며, 이번 작업에도 큰 영감을 받았습니다.

Beesophie, 감사합니다!
