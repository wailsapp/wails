---
title: "Windows UAC 구성"
description: "Windows용 Wails 애플리케이션의 사용자 계정 컨트롤(UAC) 구성"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

관련 플랫폼: <span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Windows 사용자 계정 컨트롤(UAC)은 Wails 애플리케이션의 실행 권한을 결정합니다. 기본적으로 Wails v3 애플리케이션은 Windows 매니페스트에 명시적인 UAC 구성을 포함하므로 여러 컴퓨터에서 일관되게 작동합니다.

## UAC 실행 수준

Windows 애플리케이션은 매니페스트 파일을 통해 여러 실행 수준을 요청할 수 있습니다. Wails v3는 기본 실행 수준이 지정된 UAC 구성을 자동으로 포함하며, 애플리케이션의 요구 사항에 맞게 이를 사용자 지정할 수 있습니다.

### 사용 가능한 실행 수준

| 수준 | 설명 | 사용 사례 |
| --- | --- | --- |
| `asInvoker` | 부모 프로세스와 동일한 권한으로 실행됩니다. | 대부분의 애플리케이션에 사용하는 기본값 |
| `highestAvailable` | 사용자에게 허용된 가장 높은 권한으로 실행됩니다. | 상승된 권한이 필요할 수 있는 애플리케이션 |
| `requireAdministrator` | 항상 관리자 권한이 필요합니다. | 시스템 유틸리티, 설치 프로그램 |

### 기본 구성

Wails v3 애플리케이션의 Windows 매니페스트에는 다음과 같은 기본 UAC 구성이 포함됩니다.

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

이 구성은 애플리케이션이 다음과 같이 작동하도록 보장합니다.

- 애플리케이션을 시작한 프로세스와 동일한 권한으로 실행됩니다.
- 기본적으로 권한 상승이 필요하지 않습니다.
- 여러 컴퓨터에서 일관되게 작동합니다.
- 일반 사용자에게 UAC 프롬프트를 표시하지 않습니다.

## UAC 구성 사용자 지정

Wails v3에서는 사용자가 빌드 자산을 사용자 지정하도록 권장하므로 Windows 매니페스트 템플릿을 직접 편집하여 UAC 구성을 변경할 수 있습니다.

### 매니페스트 템플릿 위치 찾기

Windows 매니페스트 템플릿은 다음 위치에 있습니다.

```
build/windows/wails.exe.manifest
```

### 실행 수준 변경

실행 수준을 변경하려면 `requestedExecutionLevel` 요소의 `level` 특성을 편집하세요.

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### 예제

#### 표준 애플리케이션(기본값)

대부분의 애플리케이션에서는 기본 `asInvoker` 수준을 사용하는 것이 좋습니다:

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### 시스템 유틸리티

상승된 권한을 사용할 수 있을 때 해당 권한이 필요한 애플리케이션:

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### 관리 도구

항상 관리자 권한이 필요한 애플리케이션:

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## UI 액세스

`uiAccess` 특성은 애플리케이션이 더 높은 권한의 UI 요소와 상호 작용할 수 있는지를 제어합니다. 대부분의 경우 이 값은 `false`로 유지하는 것이 좋습니다.

애플리케이션에서 다음 작업이 필요한 경우에만 `true`로 설정하세요.

- 다른 애플리케이션에 입력 보내기
- 다른 애플리케이션의 UI 제어하기
- 더 높은 권한으로 실행되는 프로세스의 UI 요소에 액세스하기

@note{type="caution" title="UI 액세스 요구 사항"}
`uiAccess="true"`로 설정하려면 애플리케이션이 다음 요구 사항을 충족해야 합니다.

- 신뢰할 수 있는 인증 기관에서 발급한 인증서로 디지털 서명되어 있어야 합니다.
- 보안 위치(Program Files 또는 Windows\System32)에 설치되어 있어야 합니다.

@end

## 사용자 지정 UAC 설정으로 빌드

매니페스트 템플릿을 수정한 후 평소와 같이 애플리케이션을 빌드하세요.

```bash
wails3 build
```

빌드 프로세스에서 사용자 지정 UAC 구성을 실행 파일에 자동으로 포함합니다.

## UAC 구성 확인

`go-winres` 도구를 사용하여 UAC 설정이 올바르게 포함되었는지 확인할 수 있습니다.

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

그런 다음 추출된 매니페스트 파일을 검사하여 UAC 구성이 포함되어 있는지 확인하세요.

@note{type="tip" title="매니페스트 유지"}
일부 다른 프레임워크와 달리 Wails v3의 UAC 구성은 컴파일 중에 실행 파일에 직접 포함되므로 애플리케이션을 다른 컴퓨터로 복사해도 유지됩니다.

@end

## 문제 해결

### UAC 프롬프트가 표시되지 않음

`requireAdministrator`로 설정했지만 UAC 프롬프트가 표시되지 않는 경우:

- 매니페스트가 실행 파일에 올바르게 포함되었는지 확인하세요.
- 이미 권한이 상승된 프로세스에서 실행하고 있지 않은지 확인하세요.
- 매니페스트 구문이 유효한 XML인지 확인하세요.

### 애플리케이션이 시작되지 않음

UAC 변경 후 애플리케이션이 시작되지 않는 경우:

- 매니페스트 구문에 XML 오류가 있는지 확인하세요
- 실행 수준 값이 유효한지 확인하세요
- 문제를 분리해서 확인하려면 `asInvoker`(으)로 되돌려 보세요

### 컴퓨터마다 동작이 일관되지 않음

컴퓨터마다 UAC 동작이 다른 경우:

- 매니페스트가 외부 파일이 아니라 실행 파일에 포함되어 있는지 확인하세요
- 빌드 후 실행 파일이 수정되지 않았는지 확인하세요
- 대상 컴퓨터에서 Windows UAC 설정이 활성화되어 있는지 확인하세요
