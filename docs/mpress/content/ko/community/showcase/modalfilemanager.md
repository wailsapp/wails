---
title: "Modal File Manager"
description: "Wails로 빌드한 데스크톱 애플리케이션"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager)는 웹 기술을 사용하는 이중 창 파일 관리자입니다. 제가 처음 설계한 버전은 NW.js를 기반으로 했으며 [여기](https://github.com/raguay/ModalFileManager-NWjs)에서 확인할 수 있습니다. 이 버전은 동일한 Svelte 기반 프런트엔드 코드를 사용하지만(NW.js에서 벗어난 후 크게 수정되었습니다), 백엔드는 [Wails 2](https://wails.io/) 구현입니다. 이 구현을 사용하면서 저는 더 이상 명령줄 `rm`, `cp` 등의 명령을 사용하지 않지만, 테마와 확장 기능을 다운로드하려면 시스템에 git이 설치되어 있어야 합니다. 전체가 Go로 작성되었으며 이전 버전보다 훨씬 빠르게 실행됩니다.

이 파일 관리자는 Vim과 같은 원칙, 즉 상태로 제어되는 키보드 동작을 중심으로 설계되었습니다. 상태 수는 고정되어 있지 않으며 매우 자유롭게 프로그래밍할 수 있습니다. 따라서 무한히 많은 키보드 구성을 만들어 사용할 수 있습니다. 이것이 다른 파일 관리자와의 주요 차이점입니다. GitHub에서 테마와 확장 기능을 다운로드할 수 있습니다.
