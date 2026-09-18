---
title: "Snippet Expander"
description: "Wails로 빌드한 데스크톱 애플리케이션"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Snippet Expander 스크린샷](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Snippet Expander의 스니펫 선택 창 스크린샷

![Snippet Expander 스크린샷](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Snippet Expander의 스니펫 추가 화면 스크린샷

![Snippet Expander 스크린샷](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Snippet Expander의 검색 및 붙여넣기 창 스크린샷

[Snippet Expander](https://snippetexpander.org)는 Linux용 "간편한 확장형 텍스트 스니펫 도우미"입니다.

Snippet Expander는 스니펫과 설정을 관리하기 위해 Wails로 빌드한 GUI 애플리케이션으로 구성되며, 스니펫을 빠르게 선택하고 붙여넣을 수 있는 검색 및 붙여넣기 창 모드를 제공합니다.

Wails 기반 GUI, go-lang CLI 및 vala-lang 자동 확장기 데몬은 모두 D-Bus를 통해 go-lang 데몬과 통신합니다. 이 데몬은 대부분의 작업을 수행하며 스니펫 데이터베이스와 공통 설정을 관리하고, 스니펫 확장 및 붙여넣기 등의 서비스를 제공합니다.

[소스 코드](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38)를 살펴보면 Wails 앱이 UI에서 백엔드로 메시지를 보내고, 백엔드가 다시 데몬으로 메시지를 전달하는 방식을 확인할 수 있습니다. 또한 앱의 다른 인스턴스나 CLI를 통한 스니펫 변경 사항을 모니터링하도록 D-Bus 이벤트를 구독하고, 이를 Wails 이벤트를 통해 UI에 즉시 표시하는 방식도 확인할 수 있습니다.
