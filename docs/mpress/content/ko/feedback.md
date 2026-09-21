---
title: "피드백"
description: "Wails v3에 피드백을 제공하고 문제를 보고하는 방법"
slug: "feedback"
sourcePath: "feedback.md"
---

여러분의 피드백을 환영하며 적극 권장합니다! 새 이슈나 토론을 만들기 전에 기존 이슈 또는 토론을 검색해 주세요. 다음과 같은 방법으로 기여할 수 있습니다:

@tabs
[버그]
버그를 발견한 경우, 버그 보고 템플릿을 사용하여 GitHub에서 [이슈를 등록해](https://github.com/wailsapp/wails/issues/new/choose) 주세요.

- 간단하고 재현 가능한 예제와 함께 버그를 명확하게 설명해 주세요. 문서에서 어떤 동작이 발생해야 *하는지* 명확하지 않다면 보고서에 그 내용도 포함해 주세요.
- 보고서에 `wails3 doctor`의 출력을 포함해 주세요.
- 버그가 현재 문서에 설명된 동작과 일치하지 않는 경우에는 다음 작업도 수행해 주세요:
  - `v3/examples` 디렉터리에 있는 기존 예제를 업데이트하거나, 문제를 명확하게 보여 주는 새 예제를 만들어 주세요.
  - 해당 이슈를 참조하는 [PR](https://github.com/wailsapp/wails/pulls)을 열어 주세요.


@note{type="caution"}
*기억해 주세요*. 예상하지 못한 동작이라고 해서 반드시 버그인 것은 아닙니다. 단지 기대한 대로 동작하지 않는 것일 수 있습니다. 그런 경우에는 `Suggestions`을 사용해 주세요.

@end

Discord의 [#v3](https://discord.gg/bdj28QNHmT) 채널에서도 버그에 관해 자유롭게 논의할 수 있습니다.

[수정]
버그 수정 사항이나 문서 개선 사항이 있다면 다음을 수행해 주세요:

- [기여 가이드라인](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md)에 따라 [Wails 저장소](https://github.com/wailsapp/wails)에 풀 리퀘스트를 열어 주세요.
- 관련 이슈가 있다면 PR 설명에서 참조해 주세요.

[개선 사항]
새 기능과 공개 동작에 대한 변경은 기능 요청 이슈가 아니라 **WEP(Wails Enhancement Proposal)** 초안 풀 리퀘스트를 통해 제안해야 합니다.

- [WEP 절차](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)를 읽어 주세요.
- 템플릿을 복사한 다음, WEP와 지원 자료만 포함하고 제목이 `[WEP] <title>`인 초안 PR을 열어 주세요.
- 먼저 [GitHub Discussions](https://github.com/wailsapp/wails/discussions) 또는 Discord의 [#v3](https://discord.gg/bdj28QNHmT) 채널에서 아이디어를 비공식적으로 논의해도 되지만, 메인테이너의 결정을 받으려면 WEP PR이 필요합니다.

[추천]
- GitHub에서 :thumbsup: 반응을 사용하여 버그, WEP, 토론에 대한 지지를 표시해 주세요.
- "+1" 또는 "저도 그렇습니다"와 같은 댓글만 추가하지는 *말아* 주세요.
- "이 버그는 ARM 빌드에도 영향을 줍니다" 또는 "다른 접근 방식으로는 ...이 있습니다"처럼 의미 있게 기여할 내용이 있다면 댓글을 추가해 주세요.

@end

알려진 문제와 진행 중인 작업은 [여기](https://github.com/orgs/wailsapp/projects/6)에서 확인할 수 있습니다.
