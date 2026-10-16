---
title: "튜토리얼"
description: "애플리케이션을 만들며 Wails 배우기"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

완성된 애플리케이션을 단계별로 만들면서 Wails의 개념을 배우는 튜토리얼입니다. 각 튜토리얼에는 실제로 작동하는 코드, 설명, 실용적인 패턴이 포함되어 있습니다.

@note{type="tip" title="Go가 처음이신가요?"}
튜토리얼을 시작하기 전에 [Go 둘러보기](https://go.dev/tour/)를 완료하세요.

@end

## QR 코드 서비스

![QR 코드 예제](/assets/qr1.png)

QR 코드 생성기를 만들면서 Wails 서비스의 기초를 배우세요. 이 튜토리얼에서는 애플리케이션 로직을 재사용 가능한 서비스로 구성하는 핵심 개념을 소개합니다.

**학습 내용:**

- Wails 서비스를 만들고 구조화하는 방법
- 외부 Go 종속성 관리
- Go 메서드를 프런트엔드에 바인딩하기
- Go와 JavaScript 간에 데이터 전달하기
- 유지보수하기 쉽도록 코드 구성하기

**추천 대상:** 서비스 아키텍처를 이해하려는 Wails 입문자

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>시작하기</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO 목록

![TODO 목록 애플리케이션](/assets/todo-app.png)

아름답고 현대적인 인터페이스를 갖춘 완전한 TODO 목록 애플리케이션을 만들어 보세요. 이 실습형 튜토리얼에서는 바닐라 JavaScript를 사용하는 실용적인 실제 애플리케이션을 통해 Wails의 핵심 패턴을 배웁니다.

**학습 내용:**

- 스레드 안전 상태 관리를 사용하는 서비스 기반 아키텍처
- CRUD 작업(생성, 읽기, 업데이트, 삭제)
- Go와 JavaScript 간의 타입 안전 바인딩
- 프레임워크의 복잡성 없이 현대적인 UI 만들기
- 올바른 오류 처리 및 유효성 검사 패턴

**완료 예상 시간:** 약 20분

**추천 대상:** 첫 번째 완전한 Wails 애플리케이션을 만들려는 분 — 프레임워크의 복잡성을 더하기 전에 기초를 이해하기에 적합합니다.

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>시작하기</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 메모

![메모 애플리케이션](/assets/notes-app.png)

네이티브 파일 대화 상자와 자동 저장 기능을 갖춘 Apple Notes 스타일의 애플리케이션을 만들어 보세요. 이 튜토리얼에서는 파일 작업, 네이티브 대화 상자, 전문적인 UI 패턴과 같은 데스크톱 전용 기능을 보여 줍니다.

**학습 내용:**

- 네이티브 파일 대화 상자(저장, 열기, 정보)
- JSON 기반 데이터 영속성
- 디바운스를 적용한 자동 저장 패턴
- 전문적인 2열 데스크톱 레이아웃
- Go에서 파일 시스템 작업 다루기

**완료 예상 시간:** 약 30분

**추천 대상:** 파일 작업 및 네이티브 OS 대화 상자와 같은 데스크톱 전용 기능을 배우려는 분

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>시작하기</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 자동 업데이트 Wails 앱

새 `wails3 init`에서 시작하여 서명된 릴리스 검증과 헬퍼 모드 바이너리 교체를 거쳐 Wails 애플리케이션에 인앱 자동 업데이트 기능을 추가하세요. GitHub Releases를 업데이트 소스로 사용합니다.

**학습 내용:**

- `app.Updater`를 Wails 앱에 통합하는 방법
- GitHub Releases 공급자 구성하기
- 다이제스트 검증을 위해 `SHA256SUMS`를 포함하여 릴리스 게시하기
- 변조 방지를 위한 Ed25519 서명 추가하기
- CSS, 사용자 지정 HTML 또는 BYO를 통해 기본 창 사용자 지정하기
- `CheckInterval`를 사용한 주기적인 백그라운드 확인

**완료 예상 시간:** 약 25분

**추천 대상:** 업데이트 가능한 데스크톱 앱을 배포하려는 분 — API뿐 아니라 전체 릴리스 파이프라인을 다룹니다.

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>시작하기</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
