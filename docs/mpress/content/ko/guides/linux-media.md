---
title: "로컬 오디오 및 비디오 재생"
description: "Linux에서 크기가 제한된 blob URL로 앱에 포함된 미디어를 재생하고 사용 후 해제합니다."
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

`Media.SetSource`를 사용하여 Wails v3 애플리케이션에서 짧은 로컬 클립을 재생하세요. Linux에서 WebKitGTK는 미디어 재생을 GStreamer에 위임하지만, GStreamer는 `wails://` URL을 직접 로드할 수 없습니다. 이 도우미는 Wails 스트림으로 클립을 받아 플레이어에 blob URL을 지정합니다.

데스크톱 전송은 수신 소켓을 열지 않고 기존 Wails 자산 전송을 사용합니다. API는 오디오 및 비디오 요소에서 작동합니다. 일반 HTTP/HTTPS 미디어는 플레이어의 `src`를 직접 사용할 수 있습니다.

## 미디어 파일 등록

프런트엔드에서 재생할 수 있는 클립을 담은 파일 시스템을 노출한 다음, 이름이 지정된 스트림에 미디어 핸들러를 등록하세요.

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

앱을 만든 후 `app.Run()`을 호출하기 전에 `registerMedia(app)`을 호출하세요.

디스크 파일에는 `os.OpenRoot(directory)`를 사용하고 `root.FS()`를 `media.NewHandler`에 전달하세요. `app.Run()`이 반환될 때까지 루트를 열어 두었다가 닫으세요. 프런트엔드가 읽어도 되는 파일만 포함된 디렉터리를 선택하세요. `os.Root`는 심볼릭 링크가 해당 디렉터리 밖을 가리켜도 접근을 제한합니다.

## 클립 로드

플레이어를 만드세요.

```html
<video id="player" controls></video>
```

npm 런타임을 사용하는 프런트엔드의 경우:

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

번들 런타임을 사용하는 애플리케이션은 가져오기를 다음과 같이 변경하세요.

```javascript
import { Media } from '/wails/runtime.js';
```

두 번째 인수는 등록된 스트림 이름입니다. 세 번째는 파일 시스템 루트 기준의 슬래시로 구분된 상대 파일 경로로, `welcome.mp4` 또는 `tutorials/intro.mp4`와 같습니다. URL이나 운영 체제 경로가 아닙니다.

소스가 지정되면 프로미스가 완료됩니다. 이후 플레이어가 이를 디코딩합니다. 지원하지 않는 코덱을 감지하려면 플레이어의 `error` 이벤트를 처리하고, 직접 재생을 시작해야 한다면 사용자 상호작용에서 `player.play()`를 호출하세요.

## 클립 교체 또는 해제

클립을 바꾸려면 `Media.SetSource`를 다시 호출하세요. 해당 플레이어의 이전 대기 중인 로드를 취소하므로 느린 응답이 최신 선택을 덮어쓸 수 없습니다. 새 클립이 성공적으로 로드될 때까지 이전 클립을 사용할 수 있으며, 그 후 이전 blob URL이 폐기됩니다.

플레이어를 닫거나 컴포넌트를 언마운트할 때 `Media.ClearSource(player)`를 호출하세요.

```javascript
Media.ClearSource(player);
```

대기 중인 로드를 취소하고 플레이어를 초기화하며 blob URL을 해제합니다. 문서에서 요소를 제거하기 전에 호출하세요. 프레임워크 컴포넌트는 언마운트 또는 폐기 훅에서 호출해야 합니다. 요소만 제거하면 blob URL은 해제되지 않습니다. 플레이어 소스에는 이 도우미들을 일관되게 사용하세요.

`<video><source ...></video>` 마크업에서는 선택한 파일 이름을 `Media.SetSource(video, "media", name)`에 전달하세요. 도우미가 부모 플레이어의 `src`를 설정하며, 이는 자식 `<source>`보다 우선합니다. 이 속성을 지우면 브라우저가 자식 요소를 다시 고려할 수 있으므로, 이 API로 모든 소스를 관리할 때는 빈 플레이어를 사용하세요.

문서에 연결되지 않은 오디오 객체도 작동합니다.

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## 다운로드 제한 및 로드 취소

기본 제한은 **소스당 32 MiB**입니다. Go에 설정한 제한 내에서 바이트 단위의 더 작거나 큰 양의 정수 제한을 선택할 수 있습니다.

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

너무 큰 파일은 `RangeError`로 거부됩니다. Go 핸들러는 내용을 읽기 전에 파일 크기를 확인하고, 설정된 제한과 프런트엔드 제한 중 작은 값까지만 전송합니다. 프런트엔드도 수신 크기를 확인하고 불완전한 전송을 거부합니다. 취소 시에는 `AbortError` 또는 `AbortController.abort(reason)`에 전달한 이유로 거부됩니다.

파일은 64 KiB 프레임으로 전송됩니다. Wails 데스크톱 스트림 전송은 폴링 응답을 1 MiB로 제한하므로, WebView2가 전체 응답을 버퍼링하더라도 하나의 응답에 미디어 파일 전체를 버퍼링하지 않습니다. 스트림 큐에는 별도의 제한된 역압 제어가 있습니다. 이 전송 제한은 플레이어 도우미가 완전한 blob을 사용하는 동작을 바꾸지 않습니다.

**재생 전에 파일 전체가 다운로드됩니다.** 바이트 제한은 파일별 제한이며 앱 전체 메모리 제한이 아닙니다. 여러 플레이어, 교체 중인 이전 클립, blob 생성 및 디코딩된 미디어가 추가 메모리를 사용할 수 있습니다. 도우미는 플레이어의 `preload` 설정과 관계없이 호출 시 전송합니다. 사용자가 클립 로드를 선택할 때 호출하세요.

큰 로컬 파일의 경우 이 도우미는 스트리밍 솔루션이 아닙니다. 제한을 늘리면 메모리 사용량도 늘어납니다. 이미 HTTP/HTTPS로 호스팅되는 미디어는 브라우저가 범위 요청으로 스트리밍하고 탐색할 수 있도록 기본 미디어 로드를 사용해야 합니다.

## Linux 재생 문제 해결

- 직접 로컬 재생 시 **No URI handler implemented for "wails"** 오류가 표시되면 `Media.SetSource`로 클립을 로드하세요. 기본 GTK4 스택과 레거시 `-tags gtk3` 스택 모두에 해당합니다.
- 로드는 성공하지만 디코딩에 실패하면 대상 시스템에 설치된 GStreamer 코덱을 확인하세요. MP4에는 일반적으로 H.264 비디오 및 AAC 오디오 지원이 필요하고, MP3에는 MP3 디코더가 필요합니다. 지원하는 배포판에서 제공할 형식을 테스트하세요.
- 전송에 실패하면 등록된 스트림 이름, 상대 파일 이름 및 파일 시스템 권한을 확인하세요. 전송 중 파일을 변경하면 불완전한 전송이 발생할 수 있으므로 쓰기가 끝난 후 다시 시도하세요.
- 앱에 Content Security Policy가 설정되어 있다면 `connect-src`에서 Wails 자산 출처를, `media-src`에서 `blob:`을 허용하세요. 로컬 전용 정책에서는 `connect-src 'self'; media-src 'self' blob:`으로 지정할 수 있습니다. 다른 지시문은 유지하세요.
- 파일이 제한을 초과하면 더 짧거나 작은 클립을 선택하거나 앱의 메모리 예산에 맞는 명시적 제한을 설정하세요.

[audio-video 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video)를 `go run .`로 실행하여 컴퓨터에서 포함된 MP3 및 MP4 샘플을 확인하세요.
