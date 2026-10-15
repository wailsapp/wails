---
title: "다른 프런트엔드 프레임워크 사용하기"
description: "기본 제공 템플릿이 없는 프레임워크를 사용하기 위해 직접 만든 Vite 프로젝트를 frontend 디렉터리에 배치하는 방법"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails는 의도적으로 소수의 프레임워크에 대해서만 기본 시작 템플릿을 제공합니다:

| 템플릿 | 언어 |
| --- | --- |
| `vanilla` | TypeScript(기본값) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

그렇다고 해서 선택지가 이것뿐인 것은 아닙니다. Wails 앱의 프런트엔드는 <strong>그저 하나의 웹 프로젝트</strong>일 뿐이므로, 정적 HTML/CSS/JS로 빌드할 수 있는 프레임워크라면 무엇이든 사용할 수 있습니다. 원하는 프레임워크(Solid, Preact, Lit, Qwik, SvelteKit, Angular 등)에 템플릿이 없다면 몇 분 만에 직접 스캐폴딩할 수 있습니다.

## `frontend/` 디렉터리의 사용 방식

Wails는 `frontend/`에 어떤 프레임워크가 있는지 신경 쓰지 않습니다. 프레임워크와 무관한 다음의 몇 가지 규약만 사용합니다:

- **`frontend/dist/`가 배포 대상입니다.** `main.go`는 `//go:embed all:frontend/dist`를 사용해 빌드된 프런트엔드를 임베드하고 애셋 서버에서 제공합니다. 빌드 결과로 정적 번들을 `frontend/dist/`(Vite의 기본 출력 디렉터리)에 생성해야 합니다.
- **빌드는 `frontend/package.json`에서 제어합니다.** Wails는 `wails3 build` 중에 프런트엔드의 `build` 스크립트를 실행합니다. `wails3 dev` 중에는 `dev`을 실행하고 핫 리로드를 위해 Vite 개발 서버를 프록시합니다.
- **바인딩은 `frontend/bindings/`에 생성됩니다.** Wails는 등록된 Go 서비스를 검사하고 그 위치에 타입 안전 SDK를 작성합니다. 다른 모듈과 마찬가지로 다음과 같이 가져올 수 있습니다:
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **개발 서버는 고정 포트에서 실행됩니다.** `wails3 dev`는 `WAILS_VITE_PORT`에 지정된 포트(기본값 `9245`)에서 Vite를 프록시하므로, `strictPort: true`을 사용해 `server.port`를 해당 포트로 설정하세요. 이는 일반적인 Vite 설정이며 Wails 플러그인은 관여하지 않습니다.
- **선택 사항 — 타입이 지정된 사용자 정의 이벤트.** 기본 제공 템플릿은 `@wailsio/runtime/plugins/vite` 플러그인도 등록합니다. 이 플러그인은 *타입이 지정된* 사용자 정의 이벤트를 사용할 때만 필요합니다. 생성된 이벤트 타입 정의를 런타임에 주입하며, 바인딩이 생성될 때까지 빌드가 실패하게 합니다. 문자열 기반 `Events.On("time", …)` API만 사용한다면 이 플러그인을 제외해도 됩니다.

그 밖의 구성 요소, 라우팅, 상태 및 스타일링은 전적으로 사용하는 프레임워크에 달려 있습니다.

## Vite로 원하는 프레임워크 스캐폴딩하기

가장 빠른 방법은 기본 제공 템플릿으로 시작한 다음(그러면 `main.go`, `Taskfile`, 빌드 애셋 및 작동하는 Go 서비스를 얻게 됩니다), `frontend/`를 해당 프레임워크용 새 Vite 프로젝트로 교체하는 것입니다.

@steps
### 기본 템플릿으로 프로젝트 만들기
```bash
wails3 init -n myapp
cd myapp
```

### `frontend/`를 해당 프레임워크용 Vite 앱으로 교체하기
Vite에서는 명령 하나로 대부분의 프레임워크를 스캐폴딩할 수 있습니다. 템플릿을 선택하세요:

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

`solid`를 `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla` 등 원하는 Vite 템플릿이나 그 템플릿의 `-ts` 변형(`solid-ts`, `preact-ts` 등)으로 바꾸세요.

### 런타임을 설치하고 Vite가 Wails 개발 서버를 가리키도록 설정하기
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime`는 JS API(`Events`, `Browser`, 대화 상자 등)를 제공합니다. `vite.config`에서 *필수로* 변경해야 하는 것은 개발 서버 포트뿐입니다. `wails3 dev`에서 개발 서버를 찾을 수 있도록 설정하세요:

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

<strong>타입이 지정된 사용자 정의 이벤트</strong>를 사용하려는 경우에만 플러그인도 추가하세요. 이 플러그인은 생성된 이벤트 타입을 주입하며, 빌드 전에 바인딩이 존재해야 합니다:

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Go 서비스 호출하기
바인딩을 한 번 생성한 다음 구성 요소 어디에서나 가져오세요:

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### 실행하기
```bash
wails3 dev
```

@end

@note{type="tip" title="일반 JavaScript 프로젝트 생성하기"}
같은 명령으로 TypeScript를 사용하지 않는 프로젝트도 스캐폴딩할 수 있습니다. `-ts`가 아닌 Vite 템플릿을 사용하기만 하면 됩니다:

```bash
npm create vite@latest frontend -- --template solid
```

@end

## 자체 스캐폴더가 있는 프레임워크

일부 프레임워크는 Vite의 `create` 템플릿으로 생성하지 않고 자체 도구를 사용합니다. 이러한 프레임워크도 사용할 수 있습니다. 각 프레임워크의 기본 명령으로 스캐폴딩한 다음 Wails 플러그인을 추가하세요:

- **SvelteKit:** `npx sv create frontend`. 정적 번들로 빌드되도록 정적 어댑터(`@sveltejs/adapter-static`)를 사용하고 SSR을 비활성화하세요.
- **Qwik:** `npm create qwik@latest`. 정적(SSG) 어댑터를 사용하세요.
- **Angular:** `ng new`로 스캐폴딩하고, `outputPath`을 `dist`로 설정한 다음, 빌드 스크립트가 `ng build`를 가리키도록 설정하세요.

규칙은 항상 같습니다. `frontend/dist/`에 정적 빌드를 생성하고, `@wailsio/runtime` Vite 플러그인을 유지하거나 런타임을 직접 가져오며, `frontend/bindings/`에서 Go 바인딩을 가져오세요.

@note{type="info"}
어떤 프레임워크에 대해 완성도 높은 구성을 만들었다면 다른 사용자도 이를 직접 `wails3 init -t`할 수 있도록 [사용자 정의 템플릿](/guides/advanced/custom-templates/)으로 게시하는 것을 고려해 보세요.

@end
