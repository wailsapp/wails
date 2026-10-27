---
title: "다음 단계"
description: "첫 번째 앱을 빌드한 후 살펴볼 내용"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

첫 번째 Wails 앱을 빌드했습니다. 아래 학습 경로를 따라가거나, 무엇을 빌드할지 이미 알고 있다면 해당 가이드로 바로 이동하세요.

## 더 알아보기

한 번에 하나씩 완전한 애플리케이션을 만들어 보며 학습을 이어가세요.

1. **[QR 코드 서비스 만들기](/tutorials/01-creating-a-service/).** 서비스를 통해 Go 로직을 프런트엔드에 노출하는 방법을 연습하세요.
2. **[TODO 목록 만들기](/tutorials/02-todo-vanilla/).** 서비스 호출을 통해 애플리케이션 데이터를 추가, 변경 및 삭제하세요.
3. **[메모 앱 만들기](/tutorials/03-notes-vanilla/).** 파일, 네이티브 대화 상자 및 저장된 데이터를 다뤄 보세요.

릴리스를 살펴볼 준비가 되면 [자동 업데이트 앱 튜토리얼](/tutorials/04-self-update-a-wails-app/)을 따라 해 보세요. 각 학습 과정에서 배울 내용을 확인하려면 [모든 튜토리얼](/tutorials/overview/)을 살펴보세요.

## 기능

**이미 직접 앱을 만들고 있나요?** 현재 수행하려는 작업에 맞는 방법 가이드를 이용하세요.

| 하려는 작업 | 시작할 문서 |
| --- | --- |
| 프런트엔드에서 Go 호출하기 | [JavaScript에서 Go 호출하기](/features/bindings/methods/) |
| 창 열기 또는 관리하기 | [창 기본 사항](/features/windows/basics/) |
| 네이티브 메뉴 또는 트레이 아이콘 추가하기 | [애플리케이션 메뉴](/features/menus/application/) 또는 [시스템 트레이 메뉴](/features/menus/systray/) |
| 사용자에게 파일을 열거나 저장하도록 요청하기 | [파일 대화 상자](/features/dialogs/file/) |
| 컴포넌트 간에 이벤트 보내기 | [이벤트 시스템](/features/events/system/) |
| 텍스트 복사 또는 드롭된 파일 받기 | [클립보드](/features/clipboard/basics/) 또는 [파일 드래그 앤 드롭](/features/drag-and-drop/files/) |
| 키보드 단축키 추가하기 | [키보드 단축키](/features/keyboard/shortcuts/) |
| 앱 빌드 및 패키징하기 | [애플리케이션 빌드하기](/guides/build/building/)를 먼저 진행한 다음, [Windows](/guides/build/windows/), [macOS](/guides/build/macos/) 또는 [Linux](/guides/build/linux/)용으로 패키징하기 |

이 페이지들은 사이드바의 **방법 가이드** 아래에 있습니다. 전체 섹션을 순서대로 읽지 않고도 관련 주제를 바로 열 수 있습니다.

## Wails의 작동 방식 이해하기

코드 이면의 설계를 이해하려면 **핵심 개념**을 참고하세요.

- [Wails의 작동 방식](/concepts/architecture/)에서는 Go 백엔드, 프런트엔드 및 네이티브 WebView가 어떻게 연결되는지 설명합니다.
- [Go-프런트엔드 브리지](/concepts/bridge/)에서는 Go와 JavaScript 사이에서 호출이 어떻게 전달되는지 설명합니다.
- [애플리케이션 수명 주기](/concepts/lifecycle/)에서는 시작 및 종료 과정을 설명합니다.
- [빌드 시스템](/concepts/build-system/)에서는 프로젝트 작업과 프런트엔드 도구의 역할을 설명합니다.

## API 찾아보기

학습 과정이 아니라 형식, 메서드, 옵션 또는 명령이 필요할 때는 [API 참조](/reference/overview/)를 이용하세요. 명령줄 작업은 [CLI 명령 개요](/reference/cli/)부터 살펴보고, 개별 명령의 세부 정보는 [CLI 플래그 및 예제](/guides/cli/)를 참고하세요.

## 도움받기

[자주 묻는 질문](/faq/)을 확인하고, [실행 가능한 예제](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples)를 살펴보거나, [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하세요. 플랫폼별 문제는 사이드바의 **문제 해결** 아래에 분류되어 있습니다.
