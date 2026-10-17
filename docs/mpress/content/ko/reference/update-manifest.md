---
title: "업데이트 매니페스트 프로토콜"
description: "Wails 앱이 자체 업데이트를 검색하고 검증하는 데 사용하는 개방형 JSON 프로토콜로, 모든 정적 파일 호스트 또는 동적 업데이트 서버에서 제공할 수 있습니다."
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Wails Update Manifest 프로토콜은 Wails 애플리케이션과 업데이트 소스 사이의 작고 개방적인 JSON 규약입니다. HTTPS를 통해 JSON 파일을 제공할 수 있는 곳이라면 어디서든 Wails 업데이트를 제공할 수 있습니다. 예를 들어 S3 버킷, GitHub Pages, CDN 또는 라이선스에 따라 릴리스 제공 여부를 제어하는 동적 업데이트 서버를 사용할 수 있습니다.

클라이언트 측 기능은 프레임워크에 `endpoint` 공급자(`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`)로 포함되어 있습니다. 이 페이지는 서버 측을 구현하는 사용자를 위한 통신 형식 참조 문서입니다.

## 설계 목표

1. **정적 호스트에 적합합니다.** 채널마다 모든 플랫폼의 아티팩트를 나열하는 매니페스트 파일 하나만 있으면 완전한 구현이 됩니다. 서버 코드는 필요하지 않습니다.
2. **동적 서버에 적합합니다.** 클라이언트는 확인할 때마다 `platform`, `arch`, `version` 및 `channel`을 전송하므로 서버는 정확히 하나의 아티팩트로 응답하거나, 라이선스 규칙을 적용하거나, 호출자가 최신 상태이면 `204 No Content`을 반환할 수 있습니다.
3. **검증을 우선합니다.** 매니페스트에는 아티팩트별 체크섬과 서명이 포함되며, Wails 업데이터는 빌드 시 애플리케이션 바이너리에 고정된 공개 키를 기준으로 이를 검증합니다. 업데이트 소스가 자체 신뢰 루트를 선택하는 일은 없습니다.

## 요청

클라이언트는 설정된 매니페스트 URL로 `GET`을 보내며, 이때 `Accept: application/json`과 애플리케이션에 설정된 모든 헤더(예: `Authorization: License <key>`)를 함께 전송합니다.

URL에는 자리표시자를 포함할 수 있으며, 클라이언트는 확인할 때마다 이를 다음 값으로 치환합니다.

| 자리표시자 | 치환 값 |
| --- | --- |
| `{{platform}}` | 실행 중인 OS를 나타내는 Go `GOOS` 값(`darwin`, `windows`, `linux`) |
| `{{arch}}` | 실행 중인 아키텍처를 나타내는 Go `GOARCH` 값(`amd64`, `arm64`, ...) |
| `{{version}}` | 현재 설치된 버전 |
| `{{channel}}` | 설정된 경우 해당 릴리스 채널 |

네 값 중 자리표시자에서 사용하지 않은 값은 각각 같은 이름의 쿼리 매개변수로 추가됩니다(`channel`은 설정된 경우에만 추가). 따라서 다음 두 구성은 모두 유효하며 서로 동일합니다.

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

정적 호스트는 수신한 쿼리 매개변수를 그대로 무시합니다.

## 응답

| 상태 | 의미 |
| --- | --- |
| `200 OK` | 매니페스트가 이어집니다. 업그레이드인지 여부는 클라이언트가 판단합니다. |
| `204 No Content` | 서버가 버전을 비교했으며 호출자는 최신 상태입니다. |
| `404 Not Found` | 게시된 항목이 없습니다(최신 상태와 동일하게 처리됨). |
| 그 밖의 모든 상태 | 오류입니다. 업데이터는 다음으로 설정된 공급자를 시도합니다. |

`200` 본문은 매니페스트 문서입니다.

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### 최상위 필드

| 필드 | 형식 | 필수 여부 | 참고 |
| --- | --- | --- | --- |
| `schemaVersion` | int | 아니요 | 프로토콜 버전입니다. 생략하면 `1`을 의미합니다. 클라이언트는 자신이 이해할 수 있는 버전보다 새로운 값을 거부합니다. |
| `version` | string | **예** | 선행 `v`이 있거나 없는 SemVer 2.0.0입니다. |
| `channel` | string | 아니요 | 정보 제공용입니다. 다른 채널로 설정된 클라이언트는 이 매니페스트를 업데이트 없음으로 처리합니다. |
| `name` | string | 아니요 | 업데이트 창에 표시되는 사람이 읽을 수 있는 릴리스 제목입니다. |
| `notes` | string | 아니요 | 업데이트 창에 렌더링되는 Markdown 형식의 릴리스 정보입니다. |
| `publishedAt` | string | 아니요 | RFC 3339 타임스탬프입니다. |
| `artifacts` | array | **예** | 다운로드 가능한 아티팩트마다 항목 하나를 둡니다. 순서는 게시자의 선호도를 나타냅니다. |
| `metadata` | object | 아니요 | 애플리케이션에 그대로 전달되는 자유 형식의 키/값 데이터입니다. |

클라이언트는 알 수 없는 필드를 무시하므로 서버는 호환성을 깨뜨리지 않고 자체 필드를 추가할 수 있습니다. 서버별 추가 항목은 `metadata`에 넣습니다.

### 아티팩트 필드

| 필드 | 유형 | 필수 여부 | 참고 |
| --- | --- | --- | --- |
| `url` | string | **예** | 절대 URL 또는 매니페스트 URL 기준의 상대 URL입니다. `http(s)`만 지원합니다. |
| `platform` | string | 아니요 | Go `GOOS` 값입니다. 일반적인 별칭(`macos`, `win`, ...)도 허용됩니다. 비어 있으면 모든 플랫폼과 일치합니다. |
| `arch` | string | 아니요 | Go `GOARCH` 값입니다. 일반적인 별칭(`x86_64`, `aarch64`, ...)도 허용됩니다. 비어 있으면 모든 아키텍처와 일치합니다. |
| `filename` | string | 아니요 | 기본값은 `url` 경로의 마지막 세그먼트입니다. |
| `filetype` | string | 아니요 | 기본값은 파일 이름 확장자입니다. |
| `size` | int | 아니요 | 바이트 단위이며 다운로드 진행률을 표시하는 데 사용됩니다. |
| `digestAlgo` / `digest` | string / base64 | 아니요 | `sha256` 또는 `sha512`입니다. |
| `signatureAlgo` / `signature` | string / base64 | 아니요 | `ed25519`, `ed25519ph` 또는 `ecdsa-p256`입니다. `signature`가 있으면 항상 `signatureAlgo`가 필요합니다. 각 알고리즘이 무엇에 서명하는지는 [업데이터 가이드](/guides/updater/#cryptographic-verification)를 참조하세요. |

클라이언트는 `platform` 및 `arch`가 실행 중인 시스템과 일치하는 **첫 번째** 아티팩트를 선택합니다. Base64 값은 패딩 유무와 관계없이 허용됩니다.

### 버전 비교

매니페스트가 업그레이드인지 여부는 항상 클라이언트 측에서 SemVer 2.0.0 우선순위 규칙에 따라 결정됩니다. 즉, 매니페스트의 `version`는 설치된 버전보다 반드시 더 최신이어야 합니다. 따라서 정적 호스팅에서도 간단히 올바르게 동작합니다(매니페스트는 항상 최신 릴리스를 기술하고, 최신 상태인 클라이언트는 아무 작업도 하지 않습니다). 한편 동적 서버에서는 대역폭 절감을 위해 `204`를 계속 사용할 수 있습니다.

## 검증 및 신뢰

체크섬과 서명은 매니페스트에 포함되지만 신뢰 루트는 포함되지 않습니다. 서명은 빌드 시 애플리케이션이 `updater.Config.PublicKey`를 통해 고정한 공개 키로 검증됩니다. 업데이트 소스가 침해되거나 바뀌더라도 자체 키를 제공할 수 없습니다. 애플리케이션에 고정된 키가 없는데 아티팩트에 서명이 포함된 경우 검증은 안전하게 실패합니다. 선언된 `signatureAlgo`가 없는 서명이나 디코딩할 수 없는 서명도 마찬가지입니다. 클라이언트는 절대로 암묵적으로 다이제스트만 사용하는 검증으로 대체하지 않습니다.

다이제스트만 있는 아티팩트는 다이제스트 검사 후 설치됩니다. 이 검사는 손상을 방지하지만 변조 방지를 위해 TLS와 호스트 자체의 무결성에 의존합니다. 보안에 민감한 항목에는 반드시 서명을 제공하세요.

프레임워크의 `ed25519ph` 방식으로 아티팩트에 서명하려면 Go 코드 몇 줄이면 됩니다.

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

실제로는 이 코드를 작성할 일이 거의 없습니다. CLI가 대신 처리합니다.

## wails3 CLI로 게시하기

`wails3 updater` 명령 그룹은 전체 게시 파이프라인을 지원합니다. 릴리스에는 다음 세 명령이 필요합니다.

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest`는 파일이나 디렉터리를 받습니다(키 자료, `.json`, 체크섬 및 릴리스 노트 사이드카 파일은 자동으로 건너뜁니다). 각 아티팩트를 SHA-512로 스트리밍하고, `-key`가 지정되면 Ed25519ph로 다이제스트에 서명하며, `MyApp-2.1.0-darwin-arm64.zip` 같은 일반적인 파일 이름에서 `platform`와 `arch`를 추론합니다(`macOS`, `win64`, `x86_64`, `aarch64` 같은 일반적인 별칭도 인식합니다. 추론할 수 없는 항목에는 경고를 출력하며, 해당 항목은 이후 모든 플랫폼과 일치합니다). 상대 URL을 출력하려면 `-url-prefix`를 생략하세요. 매니페스트는 아티팩트와 같은 위치에 업로드하세요.

`verify`는 불일치가 하나라도 있으면 0이 아닌 코드로 종료되므로 빌드와 게시 사이의 자연스러운 CI 게이트로 사용할 수 있습니다. 매니페스트를 직접 구성하는 서버를 위해 `wails3 updater sign -key updater.key <files...>`는 각 파일의 `digest`/`signature` 필드를 자체 문서에 바로 병합할 수 있는 JSON으로 출력합니다.

## 인증

인증은 서버가 처리하며, 프로토콜은 헤더만 전달합니다. 클라이언트는 매니페스트를 요청할 때마다 구성된 헤더를 다시 전송합니다. 아티팩트를 다운로드할 때는 아티팩트 URL의 호스트가 매니페스트와 같고 `https`에서 `http`로 다운그레이드되지 않는 경우에만 `Authorization` 헤더를 전송합니다. 또한 교차 출처 리디렉션이나 다운그레이드 리디렉션이 발생하면 이 헤더를 제거하므로 자격 증명이 CDN이나 오브젝트 스토리지로 유출되거나 평문으로 전송되는 일이 없습니다.

호스팅형 라이선스 서비스와 자연스럽게 연동되는 라이선스 제한 예제:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## 클라이언트 구성

전체 `endpoint.Config` 참조와 GitHub, keygen.sh 및 AppCast 공급자와 함께 폴백 체인에 공급자를 연결하는 방법은 [업데이터 가이드](/guides/updater/#providers)를 참조하세요.
