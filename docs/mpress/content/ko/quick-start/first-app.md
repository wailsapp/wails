---
title: "첫 번째 앱"
description: "10분 만에 작동하는 Wails 애플리케이션 만들기"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Wails의 핵심 개념을 보여 주는 간단한 인사말 애플리케이션을 만들어 보겠습니다.

- 로직을 관리하는 Go 백엔드
- Go 함수를 호출하는 프런트엔드
- 타입 안전 바인딩
- 개발 중 핫 리로드

**완료 예상 시간:** 10분

@note{type="tip" title="Windows 11 사용자를 위한 성능 팁"}
프로젝트 저장소로 [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/)를 사용하는 것을 고려해 보세요. Dev Drive는 개발자 워크로드에 최적화되어 있어 일반 NTFS 드라이브보다 빌드 시간과 디스크 액세스 속도를 최대 30%까지 크게 개선할 수 있습니다.

@end

## 프로젝트 만들기

@steps
### 프로젝트 생성
```bash
wails3 init -n myapp
cd myapp
```

기본 Vanilla + Vite 템플릿(Vite 번들러를 사용하는 HTML/CSS/TypeScript)으로 새 프로젝트를 만듭니다.

@note{type="tip" title="다른 템플릿"}
선호하는 프레임워크에 따라 `-t react`, `-t vue` 또는 `-t svelte`을 사용해 보세요. 기본적으로 TypeScript를 사용하며, 일반 JavaScript를 사용하려면 `-t vanilla-js` 또는 `-t react-js`을 사용하세요. 사용 가능한 모든 템플릿을 확인하려면 `wails3 init -l`을 실행하거나, [원하는 프런트엔드 프레임워크를 직접 사용하세요](/guides/dev/frontend-frameworks/).

@end

### 프로젝트 구조 이해하기
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### 앱 실행
```bash
wails3 dev
```

@note{type="info" title="첫 실행"}
처음 실행할 때는 프런트엔드 종속성을 설치하고 바인딩을 생성하는 등의 작업으로 예상보다 오래 걸릴 수 있습니다. 이후 실행부터는 훨씬 빨라집니다.

@end

앱이 열리고 인사말 인터페이스가 표시됩니다. 이름을 입력하고 "인사하기"를 클릭하면 Go 백엔드가 입력을 처리하여 인사말을 반환합니다.

@end

## 작동 방식

이 기능을 구현하는 코드를 살펴보겠습니다.

### Go 백엔드

`greetservice.go`을 여세요.

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**핵심 개념:**

1. **서비스** - 내보낸 메서드가 있는 Go 구조체
2. **내보낸 메서드** - `Greet`은 대문자로 시작하므로 프런트엔드에서 사용할 수 있습니다.
3. **간단한 로직** - 이름을 받아 인사말을 반환합니다.
4. **타입 안전성** - 입력 및 출력 타입이 정의되어 있습니다.

@note{type="tip" title="서비스와 바인딩 이해하기"}
<strong>서비스</strong>는 프런트엔드에 기능을 제공하는 독립적인 Go 모듈입니다. 내보낸 메서드가 있는 일반적인 Go 구조체이며, 애플리케이션 설정의 `Services` 필드에 등록합니다.

<strong>바인딩</strong>은 프런트엔드에서 이러한 서비스를 호출할 수 있도록 자동 생성되는 TypeScript/JavaScript SDK입니다. `wails3 dev` 또는 `wails3 build`을 실행하면 Wails가 등록된 서비스를 분석하고 `frontend/bindings/`에 타입 안전 바인딩을 생성합니다.

서비스는 백엔드 API이고, 바인딩은 이 API와 통신하는 클라이언트 라이브러리라고 생각하면 됩니다.

@end

### 서비스 등록하기

`main.go`을 열고 서비스 등록 부분을 찾으세요.

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

이 코드는 `GreetService`을 Wails에 등록하여 내보낸 모든 메서드를 프런트엔드에서 사용할 수 있게 합니다.

### 프런트엔드

`frontend/src/main.js`을 여세요.

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**핵심 개념:**

1. **자동 생성된 바인딩** - 생성된 코드에서 `GreetService`을 가져옵니다.
2. **타입 안전 호출** - 메서드 이름과 시그니처가 Go 코드와 일치합니다.
3. **기본 비동기 처리** - 모든 Go 호출은 Promise를 반환합니다.
4. **오류 처리** - Go에서 발생한 오류는 try/catch에서 포착됩니다.

@note{type="info" title="바인딩은 어디에 있나요?"}
생성된 바인딩은 `frontend/bindings/`에 있습니다. `wails3 dev` 또는 `wails3 build`을 실행하면 자동으로 생성됩니다.

**이 파일을 절대로 직접 편집하지 마세요**. 빌드할 때마다 다시 생성됩니다.

@end

## 앱 사용자 지정하기

워크플로를 이해할 수 있도록 새 기능을 추가해 보겠습니다.

### "여러 명에게 인사하기" 기능 추가

@steps
### GreetService에 메서드 추가
`greetservice.go`에 다음 코드를 추가하세요.

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### 앱이 자동으로 다시 빌드됩니다
파일을 저장하면 `wails3 dev`이 Go 코드를 자동으로 다시 빌드하고 앱을 재시작합니다.

@note{type="info" title="자동 재빌드"}
Go 코드가 변경되면 자동으로 다시 빌드하고 재시작합니다. 프런트엔드 변경 사항은 재시작 없이 핫 리로드됩니다.

@end

### 프런트엔드에서 사용하기
`frontend/src/main.js`에 다음 코드를 추가하세요.

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

브라우저 콘솔을 열고 `greetMany()`을 호출하면 인사말 배열이 표시됩니다.

@end

## 프로덕션용 빌드

앱을 배포할 준비가 되면 다음과 같이 진행하세요.

```bash
wails3 build
```

**수행되는 작업:**

- 최적화를 적용하여 Go 코드를 컴파일합니다.
- 프로덕션용 프런트엔드를 빌드합니다(최소화 적용).
- `bin/`에 네이티브 실행 파일을 생성합니다.

@tabs{sync-key="os"}
[Windows]
**출력:** `bin/myapp.exe`

두 번 클릭하여 실행합니다. 별도의 종속성은 필요하지 않습니다(WebView2는 Windows에 포함되어 있습니다).

[macOS]
**출력:** `bin/myapp.app`

Applications 폴더로 드래그하거나 두 번 클릭하여 실행합니다.

[Linux]
**출력:** `bin/myapp`

`./bin/myapp` 명령으로 실행하거나 런처용 `.desktop` 파일을 만듭니다.

@end

@note{type="tip" title="크로스 플랫폼 빌드"}
다른 플랫폼용으로 빌드하려면 [크로스 플랫폼 빌드 →](/guides/build/cross-platform/)를 참조하세요.

@end

## 배운 내용

**프로젝트 구조**

- Go 백엔드용 `main.go`
- UI 코드용 `frontend/`
- 빌드 작업용 `Taskfile.yml`

**서비스**

- 내보낸 메서드가 있는 Go 구조체를 만듭니다.
- `application.NewService()`을 사용하여 등록합니다.
- 메서드를 프런트엔드에서 자동으로 사용할 수 있습니다.

**바인딩**

- 자동 생성된 TypeScript 정의
- 타입 안전 함수 호출
- 기본적으로 비동기 방식(Promise)

**개발 워크플로**

- 핫 리로드용 `wails3 dev`
- Go 변경 시 자동으로 다시 빌드하고 재시작합니다.
- 프런트엔드 변경 사항은 즉시 핫 리로드됩니다.

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에 참여하여 커뮤니티에 질문하세요.
