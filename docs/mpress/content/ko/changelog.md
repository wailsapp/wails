---
title: "변경 이력"
description: "Wails v3의 버전 기록 및 릴리스 정보"
slug: "changelog"
sourcePath: "changelog.md"
---

범례:

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 이 프로젝트의 모든 주요 변경 사항은 이 파일에 기록됩니다.

형식은 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)를 기반으로 하며, 이 프로젝트는 [유의적 버전 관리](https://semver.org/spec/v2.0.0.html)를 준수합니다.

- 새 기능은 `Added`에 기록합니다.
- 기존 기능의 변경 사항은 `Changed`에 기록합니다.
- 곧 제거될 기능은 `Deprecated`에 기록합니다.
- 제거된 기능은 `Removed`에 기록합니다.
- 모든 버그 수정은 `Fixed`에 기록합니다.
- 취약점이 있는 경우 `Security`에 기록합니다.

_/

/_   * 자동 업데이트됨 *   릴리스 정보는 병합된 PR에서 생성됩니다. 일반 기여 PR에서는 이 파일을 편집하거나 제거된 `v3/UNRELEASED_CHANGELOG.md`을 다시 만들지 마세요. 변경 이력 유지 관리 방법은 `AGENTS.md`을 참조하세요. _/

## [미출시]

## v3.2.0-nightly.20261103 - 2026-11-03

## 변경됨
- [v3] Menu.Update 시 시스템 트레이 메뉴를 새로 고침함 — @johnmaguire의 [PR](https://github.com/wailsapp/wails/pull/6314)
- [v3/Windows] 창에 SetIcon을 구현함 — @hypherionmc의 [PR](https://github.com/wailsapp/wails/pull/6502)

## v3.2.0-nightly.20261102 - 2026-11-02

## 변경됨
- DisableBounce 대신 DisableOverscroll을 사용함 — @lethusang의 [PR](https://github.com/wailsapp/wails/pull/6077)

## v3.2.0-nightly.20261101 - 2026-11-01

## 추가됨
- 실험적인 CGO 비의존 macOS 및 Linux 백엔드를 추가함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6521)
- 네이티브 YAML 빌드, 이식 가능한 패키징 및 서명을 추가함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6531)

## 수정됨
- 여러 v3 나이틀리 빌드 문제를 수정하고 종속성을 정리함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6519)
- v3 릴리스 검증 및 픽스처 분류를 수정함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6533)
- 교차 플랫폼 릴리스 검증 실패를 해결함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6538)

## v3.1.0 - 2026-11-01

## 추가됨
- 실험적인 CGO 비의존 macOS 및 Linux 백엔드를 추가함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6521)
- 네이티브 YAML 빌드, 이식 가능한 패키징 및 서명을 추가함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6531)

## 변경됨
- 폐기된 데스크톱 메타데이터를 복사하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6391)
- fmt 없이 콘텐츠 길이를 형식화함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6393)
- 후보를 수집하지 않고 인덱스 루트를 선택함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6404)
- 탐색할 수 없는 파일의 복사 버퍼를 재사용함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6408)
- 접두사를 탐지할 때 리더를 할당하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6410)
- 알려진 확장자의 MIME 탐지 읽기를 건너뜀 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6414)
- 큐 증가와 축소가 반복되는 현상을 방지함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6415)
- 폐기될 요청 컨텍스트를 할당하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6419)
- 비활성화된 디버그 로그 준비를 건너뜀 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6420)
- MIME 탐지 버퍼를 지연 할당함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6428)
- 하이퍼링크 형식 지정 시 할당을 줄임 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6443)
- 글머리 기호 형식 지정 시 임시 할당을 제거함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6449)
- UUID 문자열 변환의 할당을 한 번으로 줄임 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6453)
- 명령문 핸들 형식 지정 시 할당을 줄임 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6454)
- 이름 기반 UUID 다이제스트의 임시 할당을 제거함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6455)
- 완성되지 않은 명령 출력 줄을 복사하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6460)
- 패키지 스테이징 중 실행 파일을 스트리밍 방식으로 복사함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6472)
- 완성되지 않은 출력 줄을 다시 스캔하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6473)
- 컴파일된 명령 라우팅 패턴을 재사용함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6475)
- 런타임 검색 중 메타데이터 조회를 피함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6482)
- 현지화된 대체 페이지의 할당을 줄임 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6485)
- 각 프로세서에서 패턴을 다시 컴파일하지 않도록 함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6490)

## 수정됨
- 시스템 트레이에 SetTooltip을 구현함 — @fan711의 [PR](https://github.com/wailsapp/wails/pull/6017)
- 설정 포털에서 다크 모드를 읽음 — @TechHutTV의 [PR](https://github.com/wailsapp/wails/pull/6072)
- Window.SetZoom을 1.0 미만으로 설정할 수 있게 함 — @APshenkin의 [PR](https://github.com/wailsapp/wails/pull/6166)
- 각 Go 리스너에 emit 순서대로 전달함 — @DevLumuz의 [PR](https://github.com/wailsapp/wails/pull/6275)
- Emit에서 hookLock을 유지한 상태로 이벤트 훅을 읽음 — @DevLumuz의 [PR](https://github.com/wailsapp/wails/pull/6276)
- Run 이후에 생성된 트레이가 누락되지 않도록 트레이의 Menu 속성을 항상 알림 — @AlbinoGeek의 [PR](https://github.com/wailsapp/wails/pull/6284)
- Run 전에 전송된 메인 스레드 작업을 Run이 시작될 때까지 보류함 — @AlbinoGeek의 [PR](https://github.com/wailsapp/wails/pull/6287)
- github 업데이터 프로바이더에서 올바른 애셋 URL을 사용함 — @gavidroselj의 [PR](https://github.com/wailsapp/wails/pull/6304)
- WebView2가 Content-Length를 생략해도 요청 본문을 유지함 — @strraaa의 [PR](https://github.com/wailsapp/wails/pull/6326)
- macOS 및 GTK 화면을 DIP 공간에서 보고함(혼합 DPI 경계) — @julianstorer의 [PR](https://github.com/wailsapp/wails/pull/6346)
- 해제된 알림을 클릭으로 보고하지 않음 — @bercedev의 [PR](https://github.com/wailsapp/wails/pull/6367)
- 터미널이 닫히면 개발 프로세스를 중지함 — @xmasdev의 [PR](https://github.com/wailsapp/wails/pull/6371)
- 증거 검사에서 변경 이력 실행 범위를 squash-merge 커밋에 매핑함 — @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/6373)
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6386)에서 작업 표시줄 재시작 디바운스를 1초 동안 유지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6387)에서 전체 화면으로 시작한 경우에도 크기 제약 조건을 유지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6388)에서 전체 화면 또는 최대화된 창을 새 제약 조건에 맞추기 위해 크기 조정하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6389)에서 liquid-glass 예제가 Run 전에 대화 상자를 표시하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6390)에서 Run 전에 요청된 대화 상자가 무기한 대기하는 대신 아무 작업도 하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6392)에서 데스크톱 항목에 명령이 없을 때 대체 경로를 사용하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6394)에서 파일 이름을 보관하지 않고 현재 콘텐츠를 감지하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6395)에서 대소문자가 혼합된 애셋 확장자를 인식하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6396)에서 사용이 끝난 콘텐츠 감지기 버퍼를 해제
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6397)에서 제거된 ringqueue 참조를 해제
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6398)에서 대체 항목을 선택하기 전에 헤더 삭제를 반영하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6399)에서 플러시 전에 대체 응답 헤더를 확정하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6400)에서 이름이 정확히 일치하는 인덱스 파일만 검색하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6401)에서 탐색할 수 없는 HEAD 응답의 본문과 전체 읽기를 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6402)에서 짧은 읽기가 발생해도 MIME 감지 접두부를 채우도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6403)에서 성공한 빈 비탐색 응답을 확정하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6405)에서 확장된 URI 대상을 중복 생성하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6406)에서 필드 확장 중 경로를 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6407)에서 XDG 데스크톱 파일 검색 디렉터리를 준수하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6409)에서 빈 플러시 시 감지기 헤더를 확정하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6411)에서 데스크톱 항목을 가리는 디렉터리를 건너뛰도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6412)에서 실행 파일 경로 내부에서 발견된 대상을 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6413)에서 타이머가 실행된 후 콜백 캡처를 해제
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6416)에서 잘못된 ldconfig 라이브러리 일치를 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6417)에서 파일 시스템 라이브러리 대체 검색의 일치 범위를 제한
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6418)에서 캐시 무효화 데이터 경합을 제거
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6421)에서 선택적 로거 구성을 준수하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6422)에서 공유 시작 URL 변경과 오래된 포트 사용을 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6423)에서 요청이 취소되면 개발 프록시 재시도를 중단하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6424)에서 번들 애셋이 없을 때 404을 반환하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6425)에서 번들 런타임 HEAD 요청을 처리하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6426)에서 서비스 라우팅 시 이스케이프된 경로를 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6427)에서 지나치게 큰 숫자형 시험판 식별자의 순서를 올바르게 지정하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6429)에서 도구 다운로드를 저장하기 전에 HTTP 오류를 거부하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6430)에서 데스크톱 검색 중 실행 제한 시간을 준수하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6432)에서 개발 프록시 재시도 간에 요청 본문을 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6434)에서 모든 반환 경로에서 HTTP 응답 본문을 닫도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6435)에서 병렬 라이브러리 검색 중 취소 오류를 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6436)에서 디코딩 전에 실패한 GitHub 태그 응답을 거부하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6437)에서 대괄호가 포함된 리터럴 템플릿 경로를 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6438)에서 wails3를 한 번만 설치하고 설치 프로그램 오류를 보존하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6439)에서 주석 블록의 숨겨진 항목과 헤더를 무시하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6440)에서 안정 릴리스 조회 오류가 드러나도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6441)에서 형식을 해석하지 않고 일반 메시지를 출력하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6442)에서 JSON null을 불러온 후 쓰기 패닉이 발생하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6444)에서 전체 저장소 읽기에 스냅샷을 반환하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6446)에서 일부 행을 순회한 후 발생한 오류를 전파하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6447)에서 압축 해제 대상 경로를 리터럴 그대로 유지하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6448)에서 닫기 및 다시 열기 시 준비된 문 핸들을 해제하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6450)에서 문 준비 작업이 닫기 또는 다시 열기와 경합하는 상황을 처리하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6451)에서 대체 인덱스 페이지에 대한 HEAD 요청을 준수하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6452)에서 완전한 16진수 문 핸들을 요구하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6456)에서 동시 실행 중에도 기본 인메모리 데이터베이스를 유지하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6457)에서 파싱된 HEAD 출력과 stderr를 분리
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6458)에서 임베드된 번들로 런타임 애셋을 생성하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6459)에서 TODO, FIXME 및 TBD 자리표시자를 감지하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6461)에서 릴리스 노트 API 응답을 검증하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6462)에서 상태 확인 시 작업 디렉터리를 준수하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6463)에서 동일 파일 복사로 인한 데이터 손실을 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6464)에서 캐시 식별 정보에 작업 디렉터리를 포함하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6465)에서 확인할 때마다 포트를 한 번만 탐색하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6466)에서 캐시 항목이 null로 디코딩될 때 패닉이 발생하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6467)에서 정수 범위 제한 없이 숫자형 시험판 버전을 증가시키도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6468)에서 레거시 검사로 인해 캐시 미스가 감춰지지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6469)에서 기본 작업을 실행기 기본 디렉터리에서 실행하도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6470)에서 탐색기 명령의 리터럴 퍼센트 기호가 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6471)에서 사전 조건과 작업 명령의 디렉터리를 일치시킴
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6474)에서 심볼릭 링크를 생성된 일반 파일로 대체
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6476)에서 빌드 실패 시 기존 DMG가 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6477)에서 따옴표로 묶인 정규식과 경로의 백슬래시가 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6478)에서 이미지를 할당하기 전에 생성할 아이콘 크기를 검증하도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6479)에서 명령 체인과 리디렉션을 셸을 통해 처리하도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6480)에서 빌드 실패 시 이전 Linux 패키지가 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6481)에서 플랫폼별 셸 대체 실행을 사용하도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6483)에서 경로 앞에 슬래시가 반복되어도 에셋을 제공하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6484)에서 빌드 메타데이터가 포함된 안정 버전을 인식하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6486)에서 리소스 생성 실패를 정확히 보고하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6487)에서 활성 리포터를 교체할 때 경합이 발생하지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6488)에서 동시에 수행되는 캐시 파일 저장을 직렬화
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6489)에서 여러 템플릿 동등성 후보를 지원하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6491)에서 저장된 null 카테고리 맵을 정규화
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6492)에서 실패한 변수 셸 명령의 진단 정보를 표시하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6493)에서 셸 명령 체인을 암시적 Go 캐시에서 제외
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6494)에서 Linux 카테고리 저장을 직렬화
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6495)에서 따옴표로 묶거나 이스케이프한 공백이 포함된 Go 출력을 캐시하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6496)에서 인코딩 실패 시 Mac 아이콘이 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6497)에서 유효하지 않은 Go 테스트 횟수 플래그를 더 이상 버리지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6498)에서 새 파일로 복사할 때 원본 권한이 유지되도록 수정
- @taliesin-ai가 여러 v3 nightly 빌드 문제를 수정하고 [PR](https://github.com/wailsapp/wails/pull/6519)에서 종속성을 정리
- @taliesin-ai가 v3 릴리스 검증과 픽스처 분류를 [PR](https://github.com/wailsapp/wails/pull/6533)에서 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6538)에서 크로스 플랫폼 릴리스 검증 실패를 해결

## v3.1.0-nightly.20261028 - 2026-10-28

## 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6490)에서 컴파일된 정규식 패턴을 재사용하여 변경 로그 생성 속도를 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6485)에서 현지화된 대체 페이지 제공 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6482)에서 불필요한 메타데이터 조회를 피하여 AppImage 런타임 검색 속도를 높임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6475)에서 컴파일된 명령 라우팅 패턴을 재사용하여 wake 성능을 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6473)에서 완료되지 않은 출력 줄을 다시 스캔하지 않아 성능을 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6472)에서 MSIX 패키지 스테이징 중 실행 파일 복사본을 스트리밍하여 메모리 사용량을 줄이고 빌드 속도를 높임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6460)에서 완료되지 않은 명령 출력 줄을 복사하지 않아 성능을 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6455)에서 이름 기반 UUID 생성 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6454)에서 SQLite 문 핸들 형식 지정 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6453)에서 UUID 문자열 변환에 필요한 할당을 한 번으로 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6449)에서 변경 로그 글머리표 항목 형식 지정 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6443)에서 터미널 하이퍼링크 형식 지정 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6428)에서 MIME 버퍼를 지연 감지하여 에셋 제공 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6420)에서 디버그 로깅이 비활성화된 경우 에셋 서버의 디버그 로그 준비를 건너뛰도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6419)에서 버려지는 요청 컨텍스트를 생성하지 않아 webview 요청 처리 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6414)에서 알려진 파일 확장자의 MIME 감지를 건너뛰어 에셋 제공 속도를 높임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6410)에서 콘텐츠 감지용 리더를 사용하지 않아 에셋 제공 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6408)에서 탐색할 수 없는 에셋 파일 제공 시 할당을 줄임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6404)에서 후보를 수집하지 않아 에셋 서버의 인덱스 루트 선택 속도를 높임
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6393)에서 fmt 없이 콘텐츠 길이의 형식을 지정하여 에셋 서버 성능을 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6391)에서 데스크톱 파일 탐색기 메타데이터를 읽을 때 메모리 복사를 줄임

## 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6498)에서 새 대상 파일로 복사할 때 원본 파일 권한이 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6388)에서 최소 또는 최대 크기 제약 조건이 변경되어도 전체 화면 및 최대화된 창의 크기가 바뀌지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6500)에서 wailsapp/wails 저장소의 예약된 스타 기록 업데이트를 게시하도록 변경
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6497)에서 wake가 유효하지 않은 Go 테스트 횟수 플래그를 더 이상 버리지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6496)에서 인코딩 실패 시 Mac 아이콘이 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6495)에서 따옴표로 묶거나 이스케이프한 공백이 포함된 Go 출력의 빌드 캐시 미스를 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6494)에서 동시 쓰기 경합을 방지하도록 Linux의 알림 카테고리 저장을 직렬화
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6493)에서 셸 명령 체인을 암시적 Go 빌드 캐시에서 제외
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6491)에서 저장된 데이터가 null일 때 알림 카테고리 맵을 불러오지 못하는 문제 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6492)에서 실패한 변수 셸 명령의 진단 정보 표시
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6489)에서 동등한 후보가 여러 개일 때 wake 템플릿 매칭 문제 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6488)에서 캐시 파일의 동시 저장을 직렬화하여 캐시 쓰기 손상 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6487)에서 활성 리포터를 교체할 때 발생하는 경합 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6486)에서 리소스 생성 실패를 CLI에 정확히 보고
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6484)에서 빌드 메타데이터가 포함된 안정 버전을 CLI가 인식하도록 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6483)에서 요청 경로가 여러 개의 선행 슬래시를 포함하는 애셋도 제공하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6481)에서 각 플랫폼의 wake 셸 대체 실행 문제 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6480)에서 빌드 실패 시 기존 Linux 패키지 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6479)에서 명령 체인과 리디렉션을 셸을 통해 처리
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6477)에서 wake 명령을 파싱할 때 따옴표로 묶인 정규식과 경로의 백슬래시 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6476)에서 빌드 패키징 실패 시 기존 DMG 파일 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6474)에서 바인딩 출력의 심볼릭 링크를 생성된 일반 파일로 교체
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6470)에서 Linux 파일 탐색기 명령의 리터럴 퍼센트 기호 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6471)에서 사전 조건과 작업 명령의 디렉터리를 일치시켜 일관되게 해석되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6469)에서 기본 작업을 실행기 기본 디렉터리에서 실행
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6468)에서 레거시 검사가 캐시 미스를 숨기지 않도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6467)에서 정수 오버플로 제한 없이 증가하도록 CLI 시험판 버전 번호 처리 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6466)에서 캐시 항목이 null로 디코딩될 때 패닉 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6465)에서 CLI가 검사당 한 번만 포트를 탐색하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6464)에서 캐시 식별 정보에 작업 디렉터리를 포함하여 오래된 wake 결과 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6463)에서 패키징 중 파일을 자기 자신에게 복사할 때 발생하는 데이터 손실 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6462)에서 상태 검사가 작업 디렉터리를 따르도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6459)에서 변경 로그 항목의 TODO, FIXME 및 TBD 자리표시자 감지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6458)에서 내장 번들로부터 런타임 애셋을 CLI에서 생성
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6457)에서 git 명령의 파싱된 HEAD 출력과 stderr 분리
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6456)에서 동시 접근 시에도 기본 인메모리 SQLite 데이터베이스 유지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6451)에서 대체 인덱스 페이지에 대한 HEAD 요청에 올바른 응답 반환
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6450)에서 SQLite 문 준비와 데이터베이스 닫기 또는 다시 열기 사이의 경합 처리
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6448)에서 SQLite 데이터베이스를 닫거나 다시 열 때 준비된 문 핸들 해제
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6447)에서 gosod 사용 시 추출 대상 경로를 리터럴로 유지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6446)에서 일부 행을 순회한 후 발생한 SQLite 쿼리 오류 전파
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6444)에서 kvstore 전체 저장소 읽기에 스냅샷 반환
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6442)에서 JSON null 값을 불러온 후 kvstore 쓰기 시 발생하는 패닉 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6441)에서 서식 지정 동사를 해석하지 않고 일반 터미널 메시지 출력
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6440)에서 업데이터가 최신 안정 릴리스를 조회할 때 발생한 오류 표시
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6439)에서 변경 로그 주석 블록의 숨겨진 항목과 헤더 무시
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6438)에서 wails3를 한 번만 설치하고 업데이터에 설치 프로그램 오류 표시
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6437)에서 gosod의 대괄호로 묶인 템플릿 경로를 리터럴로 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6436)에서 업데이터가 실패한 GitHub 태그 응답을 디코딩 전에 거부하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6435)에서 병렬 라이브러리 검색 중 취소 오류 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6434)에서 모든 반환 경로에서 HTTP 응답 본문을 닫아 리소스 누수 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6432)에서 개발 프록시가 요청을 재시도할 때 요청 본문 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6430)에서 데스크톱 파일 탐색기 검색 중 실행 기한 준수
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6429)에서 AppImage 도구 다운로드를 저장하기 전에 HTTP 오류 거부
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6427)에서 semver 비교 시 지나치게 큰 숫자형 시험판 식별자를 올바르게 정렬
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6426)에서 애셋 서버가 서비스를 라우팅할 때 이스케이프된 경로 보존
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6425)에서 애셋 서버가 번들 런타임에 대한 HEAD 요청을 처리하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6424)에서 번들 애셋이 없을 때 잘못된 응답 대신 404 반환
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6423)에서 요청이 취소되면 개발 프록시의 재시도 중지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6422)에서 애셋 서버의 공유 시작 URL 변경과 오래된 포트 사용 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6421)에서 애셋 서버가 선택적 로거 설정을 따르도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6418)에서 라이브러리 경로 캐시 무효화의 데이터 경합 제거
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6416)에서 ldconfig 출력을 검색할 때 잘못된 라이브러리 일치 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6415)에서 부하 상황의 애셋 서버 큐 증가 및 축소 진동 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6413)에서 디바운스 타이머 실행 후 콜백이 캡처한 항목을 해제하여 오래된 참조 유지 방지
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6412)에서 실행 파일 경로 내부에서 찾은 파일 탐색기 대상 보존
- 파일 탐색기에서 데스크톱 항목을 가리는 디렉터리를 건너뛰도록 수정 [PR](https://github.com/wailsapp/wails/pull/6411) 제공: @taliesin-ai
- 응답 본문이 비어 있을 때 애셋 서버 헤더가 전송되지 않는 문제 수정 [PR](https://github.com/wailsapp/wails/pull/6409) 제공: @taliesin-ai
- 애플리케이션 아이콘을 찾을 때 XDG 데스크톱 파일 검색 디렉터리를 따르도록 수정 [PR](https://github.com/wailsapp/wails/pull/6407) 제공: @taliesin-ai
- 파일 탐색기의 필드 확장 과정에서 파일 경로를 보존하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6406) 제공: @taliesin-ai
- 파일 탐색기에서 확장된 URI 대상이 중복되지 않도록 수정 [PR](https://github.com/wailsapp/wails/pull/6405) 제공: @taliesin-ai
- 성공한 빈 비탐색형 애셋 응답을 보류 상태로 두지 않고 커밋하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6403) 제공: @taliesin-ai
- 애셋 서버에서 짧은 읽기가 여러 번 발생해도 MIME 감지용 접두부를 올바르게 채우도록 수정 [PR](https://github.com/wailsapp/wails/pull/6402) 제공: @taliesin-ai
- 애셋 서버에서 비탐색형 HEAD 요청의 전체 응답 본문을 읽지 않도록 수정 [PR](https://github.com/wailsapp/wails/pull/6401) 제공: @taliesin-ai
- 애셋 서버에서 디렉터리 요청을 해석할 때 정확히 일치하는 인덱스 파일만 제공하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6400) 제공: @taliesin-ai
- 액션을 지원하지 않는 서버에서 닫은 Linux 알림이 클릭으로 보고되는 문제 수정 [PR](https://github.com/wailsapp/wails/pull/6367) 제공: @bercedev
- 터미널이 닫히면 개발 프로세스를 중지하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6371) 제공: @xmasdev
- Windows 작업 표시줄 재시작 디바운스 플래그를 지우기 전에 1초 동안 유지하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6386) 제공: @taliesin-ai
- 애셋 서버에서 폴백 응답 헤더가 손실되는 문제 수정 [PR](https://github.com/wailsapp/wails/pull/6399) 제공: @taliesin-ai
- 애셋 서버에서 폴백을 선택하기 전에 헤더 삭제를 반영하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6398) 제공: @taliesin-ai
- 메모리 누수를 방지하도록 애셋 서버의 링 큐에서 제거된 참조를 해제 [PR](https://github.com/wailsapp/wails/pull/6397) 제공: @taliesin-ai
- 메모리 사용량을 줄이도록 애셋 서버에서 사용이 끝난 콘텐츠 감지 버퍼를 해제 [PR](https://github.com/wailsapp/wails/pull/6396) 제공: @taliesin-ai
- 애셋 제공 시 대소문자가 혼합된 애셋 파일 확장자를 인식하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6395) 제공: @taliesin-ai
- 애셋 서버에서 파일 이름을 유지하지 않고 콘텐츠 유형을 올바르게 감지하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6394) 제공: @taliesin-ai
- 데스크톱 항목에 명령이 없으면 기본값을 사용하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6392) 제공: @taliesin-ai
- macOS 이외의 플랫폼에서 liquid-glass 예제가 앱 시작 전에 대화 상자를 표시하지 않도록 수정 [PR](https://github.com/wailsapp/wails/pull/6389) 제공: @taliesin-ai
- Windows에서 제한된 크기인 상태로 전체 화면을 시작해도 창 크기 제약 조건을 유지하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6387) 제공: @taliesin-ai
- Run 전에 호출된 대화 상자가 반환하지 않고 영원히 차단되는 문제 수정 [PR](https://github.com/wailsapp/wails/pull/6390) 제공: @taliesin-ai

## 보안
- 과도한 메모리 사용을 방지하도록 이미지를 할당하기 전에 생성된 아이콘 크기를 검증 [PR](https://github.com/wailsapp/wails/pull/6478) 제공: @taliesin-ai
- 릴리스 정보 API 응답을 사용하기 전에 검증하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6461) 제공: @taliesin-ai
- SQLite 플러그인에서 완전한 16진수 문 핸들을 요구하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6452) 제공: @taliesin-ai
- 의도하지 않은 라이브러리가 로드되지 않도록 파일 시스템 라이브러리의 폴백 일치 범위를 제한 [PR](https://github.com/wailsapp/wails/pull/6417) 제공: @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## 변경됨
- 404 페이지를 v2.wails.io에서 일치하는 Wails v2 문서 페이지에 연결 [PR](https://github.com/wailsapp/wails/pull/6366) 제공: @taliesin-ai
- v2 문서 사이트를 v2.wails.io로 이전 [PR](https://github.com/wailsapp/wails/pull/6365) 제공: @taliesin-ai
- v3 문서는 wails.io에서, v2 문서는 v2.wails.io에서 제공하도록 변경 [PR](https://github.com/wailsapp/wails/pull/6364) 제공: @taliesin-ai
- v3.0.0 릴리스 정보를 업데이트하고 보안 권고의 주장을 수정 [PR](https://github.com/wailsapp/wails/pull/6363) 제공: @taliesin-ai

## 수정됨
- Go 이벤트 리스너 콜백을 emit 순서대로 직렬 전달하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6275) 제공: @DevLumuz
- Run 이후에 생성된 트레이가 Linux에서 누락되지 않도록 트레이의 Menu 속성을 항상 알리도록 수정 [PR](https://github.com/wailsapp/wails/pull/6284) 제공: @AlbinoGeek
- Run 전에 제출된 메인 스레드 작업을 애플리케이션이 시작될 때까지 보류하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6287) 제공: @AlbinoGeek
- GitHub 업데이터의 릴리스 애셋을 다운로드할 때 올바른 애셋 API URL을 사용하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6304) 제공: @gavidroselj
- DPI가 혼합된 macOS 및 Linux 구성에서 화면 경계와 위치를 수정 [PR](https://github.com/wailsapp/wails/pull/6346) 제공: @julianstorer
- Windows에서 WebView2가 Content-Length 헤더를 생략해도 요청 본문을 보존하도록 수정 [PR](https://github.com/wailsapp/wails/pull/6326) 제공: @strraaa
- 잠금 상태에서 훅의 스냅샷을 만들어 이벤트 훅 emit의 데이터 경합을 수정 [PR](https://github.com/wailsapp/wails/pull/6276) 제공: @DevLumuz
- Linux에서 `SystemTray.SetTooltip()`이 아무 작업도 하지 않던 문제 수정: 이제 StatusNotifierItem 툴팁 설명이 설정됨(#6017)
- 주간 릴리스 워크플로에서 변경 로그 근거 검사를 실행하도록 변경 [PR](https://github.com/wailsapp/wails/pull/6374) 제공: @taliesin-ai
- 실행을 스쿼시 병합 커밋에 매핑하도록 릴리스 변경 로그 근거 검사를 수정 [PR](https://github.com/wailsapp/wails/pull/6373) 제공: @taliesin-ai
- 어두운 데스크톱에서도 색 구성표가 변경되기 전까지 Linux 다크 모드 감지가 라이트 모드를 반환하고, GTK4 백엔드에서 `linux:SystemThemeChanged`이 전혀 실행되지 않던 문제를 수정했습니다. 이제 두 백엔드는 `org.freedesktop.appearance` 포털 네임스페이스를 기반으로 하는 하나의 구현을 공유합니다(#6072).
- 프레임이 있는 창(Windows, macOS, Linux)에서 `Window.SetZoom`은 더 이상 1.0 미만의 배율을 1.0로 올리지 않습니다. WebView2 `ZoomFactor`, WKWebView `magnification` 및 WebKitGTK의 확대/축소 수준은 모든 양수를 허용하며, `Zoom` 창 옵션도 이미 이러한 값을 그대로 전달했습니다. 프레임 없는 창에서는 드래그 및 크기 조절 영역이 확대/축소에 맞춰 조정되지 않으므로 1.0 하한을 유지합니다([#4590](https://github.com/wailsapp/wails/issues/4590)). 이제 양수가 아닌 배율, NaN 및 Inf 배율은 1.0로 올리는 대신 무시합니다 — @APshenkin

## v3.0.0 - 2026-10-26

## 변경됨
- 나이틀리 프리릴리스는 다음 주간 릴리스를 미리 보여 주는 `vX.Y.0-nightly.YYYYMMDD` 형식으로 게시되고, 주간 릴리스에서는 마이너 버전이 올라가므로(`vX.Y.0`) Go의 `@latest`는 항상 주간 릴리스 또는 핫픽스 릴리스로 해석됩니다. `@wailsio/runtime` npm 패키지 버전은 항상 Wails 릴리스 버전과 일치하며, 나이틀리는 `nightly` dist-tag로 게시됩니다.
- 프로젝트 템플릿은 `@wailsio/runtime`을 해당 템플릿을 생성한 CLI 버전으로 고정합니다.
- Wails v3에는 Go 1.26 이상이 필요합니다(Go 1.26.10 또는 1.27.3 이상에는 `crypto/x509` 보안 수정이 포함되어 있습니다).

## 수정됨
- `wails3 update cli`이 v3 릴리스를 찾지 못하던 문제 수정: v2 태그만 고려하고 잘못된 경로에서 설치했습니다.
- 일부 Windows 빌드에서 애플리케이션 메뉴 모음을 그릴 때 Windows 앱이 `panic: SetTextColor failed`와 함께 충돌하는 문제를 수정했습니다. 이제 메뉴 모음은 유효한 디바이스 컨텍스트에 그려지거나 기본 그리기로 대체되고, `UAHMENUITEM.Position`가 C `int`와 일치하며, 메뉴를 그릴 때 패닉이 발생하지 않습니다(#6311).
- 재시작 대기 중에도 스테이징된 업데이트를 안전하게 유지합니다. 종료 중에 시작된 다운로드나 두 번째 재시작이 스왑 도우미가 설치하려고 대기 중인 페이로드를 삭제하거나 두고 경합하지 못하도록 했습니다. 제공: @julianstorer
- `pkg/w32` 초기화 과정에서 `uxtheme.dll`의 서수 함수 포인터를 아직 사용 중인데도 해당 모듈을 언로드하여 Windows 앱이 `applyDarkMode`에서 충돌하는 문제를 수정했습니다. 이제 프로세스가 종료될 때까지 모듈을 로드된 상태로 유지합니다(#6356).
- 스왑 전에 스테이징된 페이로드를 검증하고 백업을 폐기하기 전에 설치된 애플리케이션을 검증하여, 손상된 업데이트가 정상 작동하는 애플리케이션을 대체하지 않고 롤백되도록 했습니다. 제공: @julianstorer

## 제거됨
- [PR](https://github.com/wailsapp/wails/pull/5736)에서 최종 번역 동기화 후 Crowdin 통합을 제거했습니다. 제공: @github-actions[bot]

## 보안
- 모든 프레임에 가드 스크립트를 삽입하여 Linux(GTK3 및 GTK4)의 교차 출처 자식 프레임에서 오는 예약된 `wails:` 메시지를 차단하고 GHSA-2p5f-2gx8-mh8m을 완화했습니다.
- `github.com/klauspost/compress`을 1.18.7(GHSA-259r-337f-4rfw)로, `golang.org/x/image`을 0.45.0로 업데이트하고, `crypto/x509` 수정 사항(GO-2026-6809)을 포함하도록 Go 1.26.10로 빌드했습니다.
- 런타임 패키지, 설정 마법사, 예제 및 프로젝트 템플릿의 `vite`을 8.2.4 이상으로 업데이트했습니다(GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv).

### 알려진 제한 사항
- **Linux:** WebKitGTK가 핸들러에 프레임 정보를 제공하지 않으므로, 자식 프레임(예: `srcdoc` 또는 `about:blank` iframe)이 *동일 출처*인 경우 가드가 부모를 통해 네이티브 메시지 핸들러에 접근하는 것을 막을 수 없습니다. 교차 출처 프레임과 `allow-same-origin`가 없는 샌드박스 프레임은 차단됩니다. 신뢰할 수 없는 동일 출처 콘텐츠를 Wails 창에 임베드하지 마세요.
- **Android(실험적):** `window.wails` JavaScript 브리지가 메인 프레임 또는 출처 검사 없이 WebView의 모든 프레임에 노출됩니다. 교차 출처 콘텐츠를 포함하여 임베드된 iframe의 콘텐츠가 바인딩된 메서드를 호출하고 예약된 `wails:` 메시지를 전송할 수 있습니다. 신뢰할 수 없는 서드 파티 콘텐츠를 앱의 WebView에 임베드하지 마세요. 이는 자식 프레임에서 오는 `wails:` 메시지를 무시하는 macOS 및 iOS와 다릅니다. 3.0.1에서 수정할 예정입니다.
- **모든 플랫폼:** 자식 프레임의 바인딩된 메서드 호출은 네이티브 메시지 핸들러가 아니라 HTTP를 통해 `/wails/runtime`로 전달되므로 가드가 이를 제한하지 않습니다. 임베드하는 모든 iframe이 바인딩된 Go 메서드에 접근할 수 있다고 간주하세요.

## v3.0.0-beta.34 - 2026-10-20

## 수정됨
- 창을 활성화할 때 Windows webview가 포커스를 되찾지 못하고, Tab 키가 네이티브 창으로 이동했다가 다시 돌아오는 대신 webview 내부에 갇히는 문제를 수정했습니다(#3783).
- `internal/webview2/pkg/webview2` COM vtable 구조체(`ICoreWebView2_2`..`_27` 및 기타 파생 인터페이스)에서 상속된 메서드 체인이 누락되어 호출이 잘못된 vtable 슬롯으로 전달되던 문제를 수정했습니다(#6129).
- Windows 질문 대화 상자가 사용자 지정 버튼 레이블을 무시하고 버튼 콜백을 전혀 실행하지 않는 문제를 수정했습니다. 이제 사용자 지정 버튼은 TaskDialog로 표시되며, comctl32 v6을 사용할 수 없으면 위치 기반 콜백 매핑을 사용하는 MessageBox로 대체됩니다(#5794).
- Vite를 대상으로 60할 때 `wails3 dev`가 멈추는 문제를 수정했습니다. 이제 자동 프런트엔드 준비 게이트는 `Accept` 헤더 없이 전송되어 404 응답을 받던 HTTP GET 대신 TCP 프로브를 사용합니다([#6209](https://github.com/wailsapp/wails/issues/6209)).
- 열려 있는 systray 메뉴를 업데이트해도 계속 작동하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6301), 제공: @taliesin-ai
- 런타임 작업에서 일시적인 WebView2 COM 오류를 정상적으로 처리하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6302), 제공: @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## 추가됨
- 이제 바인딩 생성기가 컴파일 시간 상수 이름을 사용하여 `application.RegisterEvent`에 등록된 검색 가능한 모든 이벤트의 형식화된 이름 상수를 내보내는 `bindings/events` 모듈을 생성합니다(예: `user:login`에 대한 `UserLogin`). Vite 플러그인 없이 `Events.On`/`Emit`와 함께 사용할 수 있습니다(#5645).

## 수정됨
- Windows 질문 대화 상자가 사용자 지정 버튼 레이블을 무시하고 버튼 콜백을 전혀 실행하지 않는 문제를 수정했습니다. 이제 사용자 지정 버튼은 TaskDialog로 표시되며, comctl32 v6을 사용할 수 없으면 위치 기반 콜백 매핑을 사용하는 MessageBox로 대체됩니다(#5794).
- Vite를 대상으로 60할 때 `wails3 dev`가 멈추는 문제를 수정했습니다. 이제 자동 프런트엔드 준비 게이트는 `Accept` 헤더 없이 전송되어 404 응답을 받던 HTTP GET 대신 TCP 프로브를 사용합니다([#6209](https://github.com/wailsapp/wails/issues/6209)).
- 열려 있는 systray 메뉴를 업데이트해도 계속 작동하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6301), 제공: @taliesin-ai
- 런타임 작업에서 일시적인 WebView2 COM 오류를 정상적으로 처리하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6302), 제공: @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## 추가됨
- macOS 14 이상에 네이티브 macOS 메뉴 섹션 헤더를 추가했습니다. [문서 사이트](https://v3.wails.io)를 참조하세요(20개 페이지 업데이트). [PR](https://github.com/wailsapp/wails/pull/6297), 제공: @leaanthony

## 수정됨
- macOS에서 SetRelativePosition이 화면 작업 영역을 기준으로 동작하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6299), 제공: @leaanthony
- 부모 프로세스 시간 초과 또는 교체 실패 후 업데이터 도우미가 변경되지 않은 대상을 다시 실행하도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6298), 제공: @taliesin-ai
- 버전별 문서에서 pnpm 버전을 10으로 고정했습니다. [PR](https://github.com/wailsapp/wails/pull/6293), 제공: @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## 추가됨
- 루트 런타임 진입점을 평가하지 않고 `CancellablePromise`을 가져올 수 있도록 `@wailsio/runtime/cancellable` 하위 경로 내보내기를 공개했습니다(#6153).

## 변경됨
- Go 1.26 이상을 요구합니다. 이제 v3 모듈이 `go 1.26.0`을 선언하며, `wails3 init` 및 `wails3 service init`에서 생성한 프로젝트도 동일한 최소 버전을 선언합니다.

## 수정됨
- Go 1.27로 빌드할 때 `Marshaler`, `Unmarshaler`, `RawMessage` 같은 `encoding/json` 형식이 `encoding/json/v2` 및 `encoding/json/jsontext` 형식의 별칭이 되어 바인딩 생성기 출력이 달라지는 문제를 수정했습니다. 이제 다시 `any`로 렌더링되며 `encoding/json/v2`에 대한 모델은 생성되지 않습니다(#6292).
- `wails3 init -t`에 전달된 원격 템플릿 URL을 복제 전에 검증합니다. `https://`, `http://`, `ssh://`, `git://` 및 `user@host:path` URL만 허용하고, 값은 `--` 뒤에 git으로 전달하며, 복제 중에는 git이 해당 전송 방식만 사용하도록 제한합니다.
- HTTP 런타임 전송에서 범위를 벗어난 `x-wails-window-id` 헤더 값을 다른 창의 ID로 조용히 잘라내는 대신 거부하도록 했습니다.
- zip-slip을 방지하도록 업데이터의 아카이브 추출을 강화했습니다. 이제 항목 경로를 `filepath.IsLocal` 및 결합된 경로에 대한 포함 여부 검사로 확인하며, zip 및 tar.gz 아카이브에서 중첩된 `..`, 백슬래시, 절대 경로 및 심볼릭 링크를 통한 쓰기 항목을 테스트합니다.
- Debian 및 Ubuntu 이외의 환경에서 `wails3 generate appimage`의 AppImage가 시작할 때 "Unable to spawn a new child process: … WebKitNetworkProcess" 오류로 충돌하던 문제를 수정했습니다. 이제 번들 WebKit이 AppImage에 포함된 헬퍼 프로세스를 실행합니다. WebKit의 bubblewrap 샌드박스는 계속 활성화됩니다.
- Android 및 iOS용 `-tags server` 빌드가 컴파일되지 않던 문제를 수정하여 [#6271](https://github.com/wailsapp/wails/issues/6271) 문제를 해결했습니다.
- 임시 디렉터리가 다른 파일 시스템(예: tmpfs `/tmp`)에 있을 때 Linux에서 업데이터가 교체에 실패하던 문제를 수정했습니다. 이제 설치된 바이너리 옆에 파일을 스테이징한 후 파일 시스템 간에 복사합니다.
- 업데이터가 `.AppImage` 파일 대신 AppImage의 읽기 전용 마운트를 교체하려던 문제를 수정했습니다.
- macOS 및 iOS에서 iframe이 예약된 `wails:` 창 제어 메시지를 보낼 수 있던 문제를 수정했습니다. 이제 이러한 메시지는 창의 메인 프레임에서만 허용되며, iframe의 다른 메시지는 출처 정보와 함께 계속 `RawMessageHandler`에 전달됩니다.
- macOS에서 `WindowDidUpdate` 및 `WindowWillUpdate` 게시를 중단했습니다. `NSWindowDidUpdateNotification`가 모든 실행 루프 패스에서 발생하여 페이지가 초당 약 120회의 빈 `evaluateJavaScript:` 호출을 받았으며, macOS에서는 각 27.0 호출이 사용자 클릭의 일시적 활성화를 소모했습니다. 이로 인해 실제 클릭에서도 `navigator.clipboard.writeText()` 및 활성화가 필요한 다른 API가 `NotAllowedError` 오류로 거부되었습니다([#6170](https://github.com/wailsapp/wails/issues/6170)). 기여: @APshenkin
- Explorer가 버전 세부 정보를 표시하도록 기본 Windows info.json 언어를 en-US(0409)로 수정했습니다(#6210).
- AMD 또는 Intel iGPU와 NVIDIA dGPU가 함께 있는 노트북에서 Linux DMA-BUF 렌더러가 비활성화되던 문제를 수정했습니다. 이제 Wayland에서는 WebKit이 NVIDIA에서 렌더링할 때만, 즉 앱이 dGPU로 오프로딩되었거나 컴포지터가 해당 GPU에서 실행될 때만 비활성화됩니다. 더 이상 nvidia 모듈이 로드되었다는 이유만으로 비활성화되어 iGPU의 렌더링 속도가 극도로 느려지지 않습니다. X11 세션에서는 Xorg의 렌더링 GPU가 구성으로 지정되므로 기존 동작을 유지합니다.
- @john-okeefe가 [PR](https://github.com/wailsapp/wails/pull/6135)에서 Linux 업데이터의 장치 간 이름 변경 실패를 스테이징 복사 폴백으로 처리했습니다.
- 메뉴 모음 테마를 구성하지 않으면 Windows에서 `WebviewWindow.SetMenu`이 메뉴 모음을 다시 그리지 않던 문제를 수정하여 [#6103](https://github.com/wailsapp/wails/issues/6103) 문제를 해결했습니다.
- `WebviewWindow.SetMenu(nil)`이 macOS에서처럼 아무 작업도 하지 않는 대신 Windows에서 패닉을 일으키던 문제를 수정하여 [#6104](https://github.com/wailsapp/wails/issues/6104) 문제를 해결했습니다.
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6094)에서 재시도를 추가하여 푸시가 경합할 때도 변경 로그 항목을 보존했습니다.
- 모듈에 배포된 적 없는 바이너리를 참조하는 임베드를 제거하여 모든 플랫폼에서 `go mod vendor`이 `pattern arm64/WebView2Loader.dll: no matching files found` 오류로 실패하던 문제를 수정했습니다. @Grantmartin2002가 [PR](https://github.com/wailsapp/wails/pull/6031)에서 [#5782](https://github.com/wailsapp/wails/issues/5782) 및 [#5376](https://github.com/wailsapp/wails/issues/5376) 문제를 해결했습니다.
- @charlesbalderson이 [PR](https://github.com/wailsapp/wails/pull/6159)에서 macOS 창이 중복으로 해제되지 않도록 했습니다.
- @wzk2025가 [PR](https://github.com/wailsapp/wails/pull/6162)에서 WM_CONTEXTMENU를 처리하여 오른쪽 클릭으로 시스템 트레이 메뉴가 열리도록 했습니다.
- 디스패치된 콜백에서 패닉이 발생하면 `InvokeSync*`의 다섯 가지 변형 모두 호출자를 영원히 차단하던 문제를 수정했습니다. @RALIST가 [PR](https://github.com/wailsapp/wails/pull/6108)에서 [#6107](https://github.com/wailsapp/wails/issues/6107) 문제를 해결했습니다.
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6236)에서 제거된 GTK3 하위 메뉴를 안전하게 다시 빌드하도록 했습니다.
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6261)에서 중첩 워커를 위한 Windows 워커 전달 심을 구현했습니다.
- Windows에서 `CanChooseDirectories(true)`일 때 `OpenFile()`이 `PromptForMultipleSelection()`을 무시하던 문제를 수정했습니다. 이제 macOS 및 Linux에서처럼 여러 폴더를 한 번에 선택할 수 있습니다(#6263).
- windows/arm64에서 WebView2 COM 메서드에 값으로 전달되는 `double` 인수 문제를 수정했습니다. 이 문제는 `ZoomFactor`, `RasterizationScale`, 쿠키 `Expires`, 인쇄 설정의 배율·페이지 크기·여백에 영향을 미쳤으며, 이전에는 인쇄 설정 setter가 `ErrDoubleArgUnsupported`을 반환했습니다. 이제 영향을 받는 호출은 `d0`을 로드하는 `runtime.cgocall` 트램펄린을 거칩니다. Go의 syscall 경로에서는 이를 로드할 수 없습니다(golang.org/issue/62583). 기여: @APshenkin

## 보안
- 이제 서버 모드(`-tags server`)는 서버 자체 출처 또는 `ServerOptions.WebSocketOriginPatterns`에 나열된 출처에서 온 `POST` 요청으로만 `/wails/runtime` 호출을 허용합니다. 다른 메서드에는 HTTP 405를, 다른 브라우저 출처에는 HTTP 403를 반환합니다. `Host` 헤더를 다시 작성하는 리버스 프록시 뒤에 배포한 경우 공개 출처를 `WebSocketOriginPatterns`에 추가해야 합니다.
- 이제 루프백 주소(기본값 `localhost`, `127.0.0.0/8` 또는 `::1`)에 바인딩된 서버 모드(`-tags server`)는 `Host` 헤더가 `localhost`, IP 리터럴, 바인딩 호스트 또는 `ServerOptions.WebSocketOriginPatterns`과 일치하는 호스트가 아니면 HTTP 403로 요청을 거부합니다. 이를 통해 `/wails/runtime` 및 WebSocket 엔드포인트를 대상으로 한 DNS 리바인딩 공격을 차단합니다.
- WinRT 토스트 API가 실패할 때 Windows 알림이 더 이상 PowerShell로 폴백하지 않으며, 대신 오류를 반환합니다.

## v3.0.0-beta.30 - 2026-10-15

## 수정됨
- 업데이터가 `.AppImage` 파일 대신 AppImage의 읽기 전용 마운트를 교체하려던 문제를 수정했습니다.
- macOS에서 `WindowDidUpdate` 및 `WindowWillUpdate` 게시를 중단했습니다. `NSWindowDidUpdateNotification`가 모든 실행 루프 패스에서 발생하여 페이지가 초당 약 120회의 빈 `evaluateJavaScript:` 호출을 받았으며, macOS에서는 각 27.0 호출이 사용자 클릭의 일시적 활성화를 소모했습니다. 이로 인해 실제 클릭에서도 `navigator.clipboard.writeText()` 및 활성화가 필요한 다른 API가 `NotAllowedError` 오류로 거부되었습니다([#6170](https://github.com/wailsapp/wails/issues/6170)). 기여: @APshenkin
- Explorer가 버전 세부 정보를 표시하도록 기본 Windows info.json 언어를 en-US(0409)로 수정했습니다(#6210).
- @john-okeefe가 [PR](https://github.com/wailsapp/wails/pull/6135)에서 Linux 업데이터의 장치 간 이름 변경 실패를 스테이징 복사 폴백으로 처리했습니다.
- 메뉴 모음 테마를 구성하지 않으면 Windows에서 `WebviewWindow.SetMenu`이 메뉴 모음을 다시 그리지 않던 문제를 수정하여 [#6103](https://github.com/wailsapp/wails/issues/6103) 문제를 해결했습니다.
- `WebviewWindow.SetMenu(nil)`이 macOS에서처럼 아무 작업도 하지 않는 대신 Windows에서 패닉을 일으키던 문제를 수정하여 [#6104](https://github.com/wailsapp/wails/issues/6104) 문제를 해결했습니다.
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6094)에서 재시도를 추가하여 푸시가 경합할 때도 변경 로그 항목을 보존했습니다.
- 모듈에 배포된 적 없는 바이너리를 참조하는 임베드를 제거하여 모든 플랫폼에서 `go mod vendor`이 `pattern arm64/WebView2Loader.dll: no matching files found` 오류로 실패하던 문제를 수정했습니다. @Grantmartin2002가 [PR](https://github.com/wailsapp/wails/pull/6031)에서 [#5782](https://github.com/wailsapp/wails/issues/5782) 및 [#5376](https://github.com/wailsapp/wails/issues/5376) 문제를 해결했습니다.
- @charlesbalderson이 [PR](https://github.com/wailsapp/wails/pull/6159)에서 macOS 창이 중복으로 해제되지 않도록 했습니다.
- @wzk2025가 [PR](https://github.com/wailsapp/wails/pull/6162)에서 WM_CONTEXTMENU를 처리하여 오른쪽 클릭으로 시스템 트레이 메뉴가 열리도록 했습니다.
- 디스패치된 콜백에서 패닉이 발생하면 `InvokeSync*`의 다섯 가지 변형 모두 호출자를 영원히 차단하던 문제를 수정했습니다. @RALIST가 [PR](https://github.com/wailsapp/wails/pull/6108)에서 [#6107](https://github.com/wailsapp/wails/issues/6107) 문제를 해결했습니다.

## v3.0.0-beta.29 - 2026-10-13

## 추가됨
- 이제 macOS의 전역 단축키가 키보드 레이아웃을 인식합니다. QWERTY의 고정된 물리적 키 위치를 사용하는 대신 `UCKeyTranslate`을 통해 QWERTY가 아닌 레이아웃(AZERTY, QWERTZ, Dvorak 등)에서 문자 가속기를 해당 문자가 표시된 키에 바인딩합니다.

## 수정됨
- @abhicodes-007가 [PR](https://github.com/wailsapp/wails/pull/6233)에서 Linux의 네이티브 메뉴를 설정하기 전에 초기화하도록 했습니다.
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6237)에서 출력 처리 및 설정 마법사의 macOS GA 준비 상태 회귀 문제를 해결했습니다.
- Windows의 WebView2 COM `AddRef`/`Release` 래퍼가 오래된 스레드 마지막 오류를 실패로 보고하던 문제를 수정했습니다. 이제 참조 횟수를 반환하므로 release가 더 이상 애플리케이션을 중단하거나 자산 요청을 실패하게 만들 수 없습니다(#3931).
- Windows에서 지연된 모든 에셋 요청의 연기 처리가 누수되게 한 `ICoreWebView2Deferral.Release`의 `AddRef` 슬롯 호출 문제 수정(#3931)

## v3.0.0-beta.28 - 2026-10-05

## 수정됨
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6225)에서 macOS 앱 출력이 파일 제공자 외부에 유지되도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6224)에서 설정 마법사를 완성하고 시스템 세부 정보를 표시하도록 개선
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6222)에서 Windows의 WebView2 복구 오류 처리 및 리소스 정리를 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6220)에서 앱 충돌과 컨트롤러 고립을 방지하도록 Windows의 WebView2 복구 기능을 개선
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6221)에서 macOS 미디어 권한 부여 안내를 수정
- Linux 시스템 트레이 아이콘을 마우스 오른쪽 버튼으로 클릭할 때 메뉴가 열리는 동시에 클릭 핸들러도 실행되어, 연결된 창이 오른쪽 클릭할 때마다 전환되던 문제 수정(#6018)
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6184)에서 Windows의 WebView2 콜백 핸들러가 네이티브 해제 시점까지 유지되도록 수정

## v3.0.0-beta.27 - 2026-10-01

## 수정됨
- macOS 12 이상에서 `Permissions` 옵션이 무시되던 문제 수정: 이제 `WKUIDelegate` 미디어 캡처 메서드가 구현되어 `PermissionAllow` 및 `PermissionDeny`가 Linux 및 Windows와 마찬가지로 카메라와 마이크 요청에 적용됩니다. 여전히 `NSCameraUsageDescription` / `NSMicrophoneUsageDescription`가 필요하며, 샌드박스 환경에서는 이에 맞는 기기 엔타이틀먼트도 필요합니다.

## v3.0.0-beta.26 - 2026-09-25

## 추가됨
- @atterpac이 [PR](https://github.com/wailsapp/wails/pull/6150)에서 감시자가 프런트엔드 개발 서버에 HTTP 준비 상태 검사를 적용하도록 개선

## 수정됨
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6002)에서 빈 창을 그대로 두는 대신 WebView2 프로세스 장애로부터 복구하도록 수정

## v3.0.0-beta.25 - 2026-09-22

## 수정됨
- 주 애플리케이션이 종료될 때 `wails3 dev` 및 해당 백그라운드 프로세스도 중지하도록 수정(#6048)

## v3.0.0-beta.24 - 2026-09-20

## 보안
- 최소 권한 원칙을 따르도록 빌드 및 게시 워크플로의 GITHUB_TOKEN 권한을 제한
- 화면 예제의 에셋 미들웨어에서 경로 순회 및 심볼릭 링크 이탈 문제 수정
- 설정 마법사의 종속성 설치를 승인된 패키지 관리자 명령으로 제한하면서도 정상적으로 작동하도록 유지

## v3.0.0-beta.23 - 2026-09-16

## 추가됨

- @overlordtm이 [PR](https://github.com/wailsapp/wails/pull/5972)에서 Linux의 GTK 애플리케이션 ID를 재정의할 수 있도록 지원

## 수정됨

- @overlordtm이 [PR](https://github.com/wailsapp/wails/pull/5973)에서 Linux 단일 인스턴스가 앱 자체 ID로 버스 이름을 점유하도록 수정
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6126)에서 Linux의 중단된 사용자 지정 스킴 요청을 취소하도록 수정

## v3.0.0-beta.22 - 2026-09-14

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6122)에서 번역을 완료하고 커뮤니티 미디어를 이전 — [문서 사이트](https://v3.wails.io) 참조(1812개 페이지 업데이트)
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6120)에서 다양한 스타일과 크기의 삽입 가능한 Wails 버튼을 문서화

## v3.0.0-beta.21 - 2026-09-13

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6116)에서 Wails v3 문서를 M-Press로 제공하도록 추가

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6118)에서 변경 로그 생성을 위해 MPD 프런트매터의 JSON 슬러그 값을 파싱하도록 수정
- @cnmax가 [PR](https://github.com/wailsapp/wails/pull/6080)에서 백업 실패 후 업데이터가 도우미 환경 변수를 지우고 원래 대상을 다시 실행하도록 수정
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6098)에서 App.Run 중 기본 시그널 핸들러를 시작하도록 수정
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/6112)에서 Windows 메뉴가 nil 메뉴를 처리하고, 교체된 리소스를 해제하며, 메뉴 모음을 다시 그리도록 수정
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6115)에서 공유 YAML 구성을 사용하는 새 프로젝트의 MSIX 패키징을 복원
- 제네릭 모델 생성자가 뒤에 선언된 도우미를 참조할 때 생성된 JavaScript 및 TypeScript 바인딩을 불러오지 못하던 문제를 수정하고, 상호 의존적인 제네릭 모델을 생성할 때 스택 오버플로가 발생하지 않도록 방지(#6062)

## v3.0.0-beta.20 - 2026-09-10

## 변경됨

- @01xR4in이 [PR](https://github.com/wailsapp/wails/pull/6082)에서 Clave 쇼케이스 링크를 현재 웹사이트와 저장소로 업데이트

## 수정됨

- 탐색 전반에서 keepalive 핸들러를 유지하면서 워커 요청을 포함한 중단된 Windows 에셋 요청을 취소합니다. Apple 플랫폼에서는 네이티브 요청 컨텍스트를 애플리케이션 래퍼를 통해 전달합니다. (#5963, #5969)
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6094)에서 푸시가 경합할 때 재시도를 통해 변경 로그 항목이 보존되도록 수정
- 모듈에 배포된 적이 없는 바이너리를 참조하던 임베드를 제거하여 모든 플랫폼에서 `go mod vendor`가 `pattern arm64/WebView2Loader.dll: no matching files found` 오류로 실패하던 문제를 수정하고, [#5782](https://github.com/wailsapp/wails/issues/5782) 및 [#5376](https://github.com/wailsapp/wails/issues/5376)을 해결함. @Grantmartin2002의 [PR](https://github.com/wailsapp/wails/pull/6031)

## 제거됨

- 순수 Go 로더로 대체된 네이티브 WebView2 로더 지원을 제거합니다. 이에 따라 내장된 `WebView2Loader.dll` 바이너리와 `github.com/jchv/go-winloader` 종속성이 제거됩니다. `native_webview2loader` 빌드 태그는 계속 허용되어 더 이상 오류가 발생하지 않지만 v3 빌드에는 아무런 영향도 주지 않습니다. @Grantmartin2002의 [PR](https://github.com/wailsapp/wails/pull/6031)
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6097)에서 macOS API 가이드의 사용하지 않는 빌드 태그와 FPS 옵션을 제거

## v3.0.0-beta.19 - 2026-09-09

## 추가됨

- 선택적으로 사용하도록 비공개 macOS API를 빌드 태그로 제한 — [문서](https://v3.wails.io/features/browser/integration), [문서](https://v3.wails.io/features/environment/info), [문서](https://v3.wails.io/features/windows/basics), [문서](https://v3.wails.io/features/windows/frameless), [문서](https://v3.wails.io/features/windows/notch-windows), [문서](https://v3.wails.io/features/windows/options), [문서](https://v3.wails.io/guides/build/macos), [문서](https://v3.wails.io/guides/build/private-macos-apis) 및 [문서](https://v3.wails.io/reference/overview) 참조. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/6087)

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6091)에서 64 MiB를 초과하는 런타임 요청을 HTTP 413로 거부하도록 수정

## 보안

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/6092)에서 토큰 인증으로 MCP 출처 및 원격 액세스 보안을 강화

## v3.0.0-beta.18 - 2026-09-08

## 수정됨

- 포인터 리시버를 사용하여 Linux 및 Darwin의 Calloc 메모리 누수를 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6083), @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## 수정됨

- Windows: WebResourceRequested 핸들러에서 `GetRequest`이 실패하거나 nil이어도 더 이상 프로세스가 종료되지 않습니다(`log.Fatal` / nil 역참조 패닉). 대신 요청을 폐기하고 로그에 기록합니다. [PR](https://github.com/wailsapp/wails/pull/6006), @midagedev

## v3.0.0-beta.16 - 2026-08-29

## 변경됨

- 새 터미널 창에서 공증 암호를 입력하도록 변경했습니다. [PR](https://github.com/wailsapp/wails/pull/6029), @leaanthony

## 수정됨

- macOS에서 시스템 트레이 클릭 유형을 올바르게 처리하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5919), @ChewbaccaCookie
- CI가 업데이트 전에 사용하지 않는 Microsoft apt 저장소를 제거하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6041), @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## 수정됨

- WebView2 임베드 제한 시간을 60초로 늘렸습니다. [PR](https://github.com/wailsapp/wails/pull/6043), @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## 수정됨

- macOS에서 Control-문자 키 입력의 이름을 올바르게 지정하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6032), @taliesin-ai
- Windows에서 ICO 트레이 아이콘을 수정하고 작업 표시줄 테마를 따르도록 했습니다. [PR](https://github.com/wailsapp/wails/pull/6016), @nik9play

## v3.0.0-beta.13 - 2026-08-25

## 수정됨

- macOS에서 모달 루프가 실행되는 동안에도 메인 스레드 작업을 계속 처리하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6026), @leaanthony
- 모바일 보안 저장소 작업이 실패할 수 있도록 하고, 실패 시 닫힌 상태를 유지하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5923), @mortenolsrud
- 등록된 리스너가 없어도 애플리케이션 이벤트 훅을 실행하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5999), @archy-rock3t-cloud
- 주석과 현지화된 문서의 오타를 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6023), @haoku123
- `v3/examples` 아래에 커밋된 사전 컴파일 macOS 바이너리를 제거했습니다. [PR](https://github.com/wailsapp/wails/pull/6025), @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## 추가됨

- 수명 주기 및 텔레메트리 예제가 포함된 macOS 노치 알림 창을 추가했습니다. [문서](https://v3.wails.io/features/windows/notch-windows)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/6010), @leaanthony
- 새 옵션과 네이티브 통합을 제공하는 macOS NSPanel 창 지원을 추가했습니다. [문서](https://v3.wails.io/features/windows/options)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/6008), @leaanthony

## 수정됨

- 동시 호출 시 SQLite Prepare가 멈추는 문제를 방지했습니다. [PR](https://github.com/wailsapp/wails/pull/5998), @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## 제거됨

- 문서에서 더 이상 사용하지 않는 구현 추적기를 제거했습니다. [PR](https://github.com/wailsapp/wails/pull/6005), @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## 수정됨

- GTK4 Linux 호스트가 사용자 지정 프로토콜 및 파일 연결 실행 인수를 누락하는 문제를 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/6000), @midagedev
- 변경 로그 검증에서 삭제된 줄과 동일 소스의 수정 사항을 올바르게 처리하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5993), @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## 추가됨

- 에이전트 지원 프로젝트 관리를 위한 안전한 wails3 mcp 서버를 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5896), @leaanthony
- 바인딩의 모델에 관한 문서를 추가했습니다. [문서](https://v3.wails.io/features/bindings/models)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/5988), @taliesin-ai
- Atomic Linux 시스템에서 rpm-ostree 설치를 지원합니다. [PR](https://github.com/wailsapp/wails/pull/5987), @leaanthony
- 네이티브 주간 스타 기록 차트 생성 및 게시 기능을 추가했습니다. [문서](https://v3.wails.io/credits), [문서](https://v3.wails.io/de/credits), [문서](https://v3.wails.io/fr/credits), [문서](https://v3.wails.io/id/credits), [문서](https://v3.wails.io/ja/credits), [문서](https://v3.wails.io/ko/credits), [문서](https://v3.wails.io/pt/credits), [문서](https://v3.wails.io/ru/credits), [문서](https://v3.wails.io/zh-cn/credits) 및 [문서](https://v3.wails.io/zh-tw/credits)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/5986), @leaanthony
- 애플리케이션 번들 리소스를 확인하기 위한 Darwin 전용 mac 패키지를 추가했습니다. [문서](https://v3.wails.io/guides/build/macos)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/5965), @leaanthony
- Condui 쇼케이스 페이지와 색인 항목을 추가했습니다. [문서](https://v3.wails.io/community/showcase/condui) 및 [문서](https://v3.wails.io/community/showcase)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/5962), @mgueregath
- 스크린샷과 프로젝트 링크가 포함된 Redis Viewer 쇼케이스 페이지를 추가했습니다. [문서](https://v3.wails.io/community/showcase) 및 [문서](https://v3.wails.io/community/showcase/redisviewer)를 참조하세요. [PR](https://github.com/wailsapp/wails/pull/5984), @redisviewer

## 변경됨

- Linux의 GTK 애플리케이션 플래그를 G<em>APPLICATION</em>NON_UNIQUE로 업데이트했습니다. [PR](https://github.com/wailsapp/wails/pull/5971), @overlordtm
- 누락된 창 이벤트를 경고 수준 대신 디버그 수준으로 기록하도록 변경했습니다. [PR](https://github.com/wailsapp/wails/pull/5914), @julianstorer

## 수정됨

- 등록된 macOS 키보드 단축키가 webview보다 우선하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5902), @julianstorer
- 끊어진 문서 사이드바 링크를 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5937), @northes
- WebKit이 일치하는 사용자 지정 스킴 작업을 중단하면 macOS 및 iOS 애셋 요청 컨텍스트를 취소하도록 수정했습니다(#5963).
- WindowSetFullscreenButtonEnabled 메시지를 처리하도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5976), @archy-rock3t-cloud
- 빌드 실패를 해결하기 위해 preact-ts 템플릿에서 Fragment를 가져오도록 수정했습니다. [PR](https://github.com/wailsapp/wails/pull/5979), @haoku123
- 활성 창이나 디스플레이를 사용할 수 있기 전에 화면 검색이 실행될 때 레거시 GTK3 서비스 전용 애플리케이션이 충돌하지 않도록 수정했습니다(#5966).
- 미출시 변경 로그가 비어 있어도 명시적 버전을 지정한 릴리스 실행을 계속할 수 있도록 수정했습니다(#5977).

## 보안

- 보안 권고를 해결하기 위해 웹사이트 nanoid 잠금 파일을 패치된 3.3.18 버전으로 업데이트했습니다. [PR](https://github.com/wailsapp/wails/pull/5985), @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## 추가됨

- 자동 변경 로그 항목에 문서 URL 생성 기능을 추가했습니다. [PR](https://github.com/wailsapp/wails/pull/5957), @taliesin-ai
- Streams 추가: 수신 대기 소켓 없이 WebSocket 프로그래밍 모델을 사용하여 Go와 JavaScript 간에 양방향 바이트 스트림을 제공합니다. Go에서 `app.HandleStream(name, handler)`을 사용해 스트림을 선언하고 프런트엔드에서 `Stream(name)`을 사용해 연결하면 `WebSocket` 형태의 객체가 반환됩니다. Go→JS 전송은 애셋 서버를 통해 창마다 하나의 대기 폴로 처리되고, JS→Go 전송은 일반 POST로 처리됩니다. TCP 포트에 바인딩되는 것은 없으며 `evaluateJavaScript`을 통과하는 것도 없습니다. 서버 빌드(`-tags server`)에서는 동일한 핸들러가 실제 WebSocket을 통해 대신 제공되므로 빌드가 달라도 애플리케이션 코드는 동일합니다. 작성자: @leaanthony
- 메일박스 변경 로그 항목을 Unreleased로 이동([PR](https://github.com/wailsapp/wails/pull/5935)). 작성자: @leaanthony

## 변경됨

- 문서 사이드바 자동 생성 및 블로그 작성자 유형 파생 방식 업데이트([PR](https://github.com/wailsapp/wails/pull/5938)). 작성자: @leaanthony

## 수정됨

- WebView2 초기화에서 기한과 메시지 펌프 사용([PR](https://github.com/wailsapp/wails/pull/5952)). 작성자: @leaanthony
- WebView2 쿠키 테스트는 명시적으로 사용 설정하지 않으면 CI에서 건너뛰고, 실행을 현재 OS 스레드에 고정([PR](https://github.com/wailsapp/wails/pull/5951)). 작성자: @leaanthony
- Windows 메뉴 빌더에서 하위 메뉴의 상위 항목에 대한 명령 ID 복원([PR](https://github.com/wailsapp/wails/pull/5944)). 작성자: @gilad-ch
- 공식 크로스 컴파일 이미지가 GTK 4.14+ Linux 지원 기준과 일치하도록 조정(#5928)
- 상속된 링커 플래그를 유지하고 -ObjC를 추가하도록 iOS Xcode 프로젝트 구성([PR](https://github.com/wailsapp/wails/pull/5915)). 작성자: @mortenolsrud
- 대규모 프런트엔드에서 `wails3 dev` 애셋 프록시가 TCP 연결을 과도하게 반복 생성하던 문제를 수정했습니다. 이 문제로 호스트의 임시 포트가 고갈되어 관련 없는 프로세스가 `EADDRNOTAVAIL` 오류로 실패할 수 있었습니다.
- 순차 디스패치와 백프레셔를 위해 창별 이벤트 JavaScript를 큐에 저장([PR](https://github.com/wailsapp/wails/pull/5934)). 작성자: @leaanthony

## 제거됨

- 데스크톱 바이너리 릴리스 파이프라인 제거: v3 릴리스는 태그로만 제공되며 `wails3` CLI는 `go install`으로 설치합니다. `release-v3.yml` 및 이를 디스패치하던 나이틀리 단계를 삭제했습니다([PR](https://github.com/wailsapp/wails/pull/5946)). 작성자: @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## 추가됨

- 미디어 재생 시 사용자 동작 요구 사항을 비활성화하는 macOS 자동 재생 환경설정 추가([PR](https://github.com/wailsapp/wails/pull/5512)). 작성자: @Eyalm321
- 메일박스 변경 로그 항목을 Unreleased로 이동([PR](https://github.com/wailsapp/wails/pull/5935)). 작성자: @leaanthony

## 변경됨

- 더 부드럽게 동작하도록 macOS 확대/축소 애니메이션에서 CADisplayLink 또는 NSTimer 사용([PR](https://github.com/wailsapp/wails/pull/5945)). 작성자: @savely-krasovsky

## 수정됨

- 상속된 링커 플래그를 유지하고 -ObjC를 추가하도록 iOS Xcode 프로젝트 구성([PR](https://github.com/wailsapp/wails/pull/5915)). 작성자: @mortenolsrud
- 대규모 프런트엔드에서 `wails3 dev` 애셋 프록시가 TCP 연결을 과도하게 반복 생성하던 문제를 수정했습니다. 이 문제로 호스트의 임시 포트가 고갈되어 관련 없는 프로세스가 `EADDRNOTAVAIL` 오류로 실패할 수 있었습니다.
- 순차 디스패치와 백프레셔를 위해 창별 이벤트 JavaScript를 큐에 저장([PR](https://github.com/wailsapp/wails/pull/5934)). 작성자: @leaanthony

### 추가됨

- 이벤트를 순서대로 전달하기 위한 범용 비동기 FIFO 메일박스 구현([PR](https://github.com/wailsapp/wails/pull/5851)). 작성자: @savely-krasovsky 및 @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## 추가됨

- 크기가 지나치게 큰 이벤트를 위한 용량 제한 호스트 측 저장소와 순차적 JavaScript 전달 구현([PR](https://github.com/wailsapp/wails/pull/5930)). 작성자: @leaanthony
- 창 깜박임을 위한 macOS Dock 바운스 구현([PR](https://github.com/wailsapp/wails/pull/5921)). 작성자: @julianstorer

## 수정됨

- 애셋 서버가 플러시 중 콘텐츠 유형 스니퍼 오류와 기록되지 않은 접두사를 보존하도록 수정([PR](https://github.com/wailsapp/wails/pull/5931)). 작성자: @leaanthony
- Wails 콜백에서 애플리케이션 메뉴를 교체할 때 macOS 애플리케이션이 충돌하지 않도록 수정
- Windows 10 1809 / Windows Server 2019(빌드 17763)에서 네이티브 메뉴를 읽을 수 없던 문제를 수정했습니다. 다크 모드 uxtheme 내보내기는 빌드 18334을 기준으로 제한되어 있었으므로 해당 호스트에서는 앱 수준의 다크 모드 사용 설정이 실행되지 않았습니다. 메뉴 배경은 어둡게 그려졌지만 Windows는 계속 라이트 테마로 메뉴 텍스트를 그려 어두운 배경에 어두운 텍스트가 표시되었습니다. 해당 서수는 17763부터 존재하므로 이제 조건이 이에 맞게 조정되었습니다.
- `w32.GetStockObject`이 `GetStockObject` 대신 `GetDeviceCaps`을 호출하여 모든 스톡 객체에 대해 0을 반환하던 문제를 수정했습니다.
- WebView2 부트스트래퍼 다운로드 오류 처리 및 보고 개선([PR](https://github.com/wailsapp/wails/pull/5924)). 작성자: @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## 수정됨

- 이제 macOS 앱 활성화에서 일반 앱에만 활성화 정책을 적용([PR](https://github.com/wailsapp/wails/pull/5897)). 작성자: @julianstorer
- Linux 빌드에서 초기화되지 않은 GTK 창 보호([PR](https://github.com/wailsapp/wails/pull/5898)). 작성자: @julianstorer
- URL을 로드하기 전에 Linux WebKit 창에 명시적인 불투명 배경색 설정([PR](https://github.com/wailsapp/wails/pull/5899)). 작성자: @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## 변경됨

- Android 빌드 작업의 기본 아키텍처를 arm64로 설정하고 deploy-emulator가 호스트 아키텍처를 선택하도록 변경([PR](https://github.com/wailsapp/wails/pull/5890)). 작성자: @mortenolsrud

## 수정됨

- 드래그 중 macOS 창의 확대/축소 상태를 유지하고 모션을 줄이도록 수정([PR](https://github.com/wailsapp/wails/pull/5900)). 작성자: @leaanthony
- `webview_window_windows_nonclient.go` 빌드 제약 조건에 `!server`을 추가하여 Windows 서버 모드 빌드 수정

## v3.0.0-beta.3 - 2026-08-03

## 추가됨

- 구현 세부 정보에 10단계 베타 검증 완료를 문서화([PR](https://github.com/wailsapp/wails/pull/5881)). 작성자: @leaanthony

## 수정됨

- Windows 다크 모드 API에 창 핸들을 전달하고 인수를 검증하도록 수정([PR](https://github.com/wailsapp/wails/pull/5877)). 작성자: @leaanthony
- 프레임 없는 창의 macOS 제목 표시줄 버튼 상태 결정 로직을 중앙화([PR](https://github.com/wailsapp/wails/pull/5870)). 작성자: @taliesin-ai
- Windows 앱 테마가 밝은 상태에서 Windows 애플리케이션이 다크 모드를 요청할 때 네이티브 메뉴 텍스트를 읽을 수 없던 문제를 방지했습니다. 이제 Windows가 어두운 메뉴 텍스트를 렌더링할 수 있을 때까지 메뉴에서 이에 맞는 밝은 네이티브 배경을 사용합니다.
- Windows 10 1809 / Windows Server 2019(빌드 17763)에서 네이티브 메뉴를 읽을 수 없던 문제를 수정했습니다. 다크 모드 uxtheme 내보내기는 빌드 18334을 기준으로 제한되어 있었으므로 해당 호스트에서는 앱 수준의 다크 모드 사용 설정이 실행되지 않았습니다. 메뉴 배경은 어둡게 그려졌지만 Windows는 계속 라이트 테마로 메뉴 텍스트를 그려 어두운 배경에 어두운 텍스트가 표시되었습니다. 해당 서수는 17763부터 존재하므로 이제 조건이 이에 맞게 조정되었습니다.

## v3.0.0-beta.2 - 2026-08-02

## 변경됨

- v3을 알파에서 베타로 승격
- 시스템 트레이의 스마트 기본값과 팝업 자동 숨김 동작을 문서화하고, 클릭 핸들러 선택에 대한 회귀 테스트를 추가했습니다(#5840).
- GitHub 업데이터가 기본적으로 Windows 설치 프로그램 애셋을 제외하도록 변경([PR](https://github.com/wailsapp/wails/pull/5861)). 작성자: @leaanthony
- [PR](https://github.com/wailsapp/wails/pull/5866)에서 macOS 프레임리스 창에 둥근 모서리, 직각 모서리 및 사용자 지정 반경 모서리 지원 추가, 작성자: @leaanthony

## 수정됨

- 실시간 GTK4 창 크기를 보고하고, 구성된 서피스에서 크기 조정, 최대화, 최소화 및 전체 화면 상태 이벤트를 발생시킵니다(#5830).
- [PR](https://github.com/wailsapp/wails/pull/5854)에서 fetch 요청으로 Blob 또는 FormData를 전송할 때 발생하는 Linux WebKit 충돌 수정, 작성자: @taliesin-ai
- [PR](https://github.com/wailsapp/wails/pull/5865)에서 Blob/FormData 헤더가 없을 때 fetch shim이 undefined를 전달하도록 변경, 작성자: @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## 추가됨

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5866)에서 macOS 프레임리스 창에 둥근 모서리, 직각 모서리 및 사용자 지정 반경 모서리 지원 추가, 작성자: @leaanthony

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5865)에서 Blob/FormData 헤더가 없을 때 fetch shim이 undefined를 전달하도록 변경, 작성자: @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5857)에서 새로운 옵션과 빌드 작업을 포함한 macOS DMG 패키징 지원 추가, 작성자: @leaanthony

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5861)에서 GitHub 업데이터가 기본적으로 Windows 설치 프로그램 자산을 제외하도록 변경, 작성자: @leaanthony

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5854)에서 fetch 요청으로 Blob 또는 FormData를 전송할 때 발생하는 Linux WebKit 충돌 수정, 작성자: @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5853)에서 최대화 또는 최소화를 수행하는 macOS 제목 표시줄 더블 클릭 동작 구현, 작성자: @taliesin-ai

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5849)에서 QR 서비스 튜토리얼이 NewServiceWithOptions를 사용하도록 업데이트하고 간격 추가, 작성자: @jeongkyu

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5856)에서 macOS 확대/축소 중에도 WKWebView가 응답성을 유지하도록 수정, 작성자: @leaanthony
- GTK4 창 크기 조회를 수정하고, 구성된 `GdkSurface`에서 크기 조정, 최대화, 최소화 및 전체 화면 상태 이벤트를 발생시킵니다.

## v3.0.0-alpha2.119 - 2026-07-27

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5833)에서 여러 언어의 문서에 아키텍처 다이어그램이 포함되도록 업데이트, 작성자: @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5825)에서 아이콘 생성 입력 및 출력에 기본 경로 제공, 작성자: @taliesin-ai
- [PR](https://github.com/wailsapp/wails/pull/5797)에서 소스 진입점 모듈을 런타임 package.json의 sideEffects에 추가, 작성자: @savely-krasovsky
- [PR](https://github.com/wailsapp/wails/pull/5816)에서 기여 가이드에 라이선스 및 출처 섹션 추가, 작성자: @taliesin-ai

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5800)에서 범위가 지정된 GTK4 프레임리스 CSS를 적용하여 테두리 반경 제거, 작성자: @savely-krasovsky
- [PR](https://github.com/wailsapp/wails/pull/5789)에서 Windows의 팝업 메뉴 및 화면 열거 시 커서 위치 조회 실패를 안전하게 처리, 작성자: @wayneforrest
- [PR](https://github.com/wailsapp/wails/pull/5678)에서 macOS 파일 열기 대화 상자가 확장자를 올바르게 필터링하고 접미사를 기준으로 허용된 파일을 검증하도록 수정, 작성자: @phergul
- [PR](https://github.com/wailsapp/wails/pull/5793)에서 Windows 다크 모드 초기화 시 nil API 호출 방지, 작성자: @roachadam
- 업데이터의 32비트 빌드 실패를 수정했습니다. `GOARCH=386`에서 `fmt.Errorf`에 전달할 때 `maxArchiveTotalSize` 상수(2 GiB)가 플랫폼 `int`의 범위를 초과했습니다. 이제 명시적으로 `int64` 형식이 지정됩니다.
- Windows 10 1809 / Windows Server 2019(빌드 17763)처럼 다크 모드 uxtheme API를 로드하지 않는 Windows 빌드에서 창이 Dark(또는 시스템 다크) 제목 표시줄을 사용할 때 시작 시 발생하는 nil 포인터 패닉을 수정했습니다. 이제 창 테마 설정의 `AllowDarkModeForWindow` 호출은 `w32.SetMenuTheme`에서 이미 사용하는 보호 처리와 동일하게 nil 여부를 확인합니다.

## v3.0.0-alpha2.117 - 2026-07-08

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5462)에서 Windows의 비클라이언트 영역을 위한 사용자 지정 히트 테스트 로직 구현, 작성자: @savely-krasovsky

## 변경됨

- [PR](https://github.com/wailsapp/wails/pull/5761)에서 UseVisualHosting을 기준으로 WebView2 모니터 배율 감지를 구성, 작성자: @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5763)에서 Wails v3 기능과 지침에 중점을 두도록 FAQ 문서 업데이트, 작성자: @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## 수정됨

- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제 수정(#5659, @puneetdixit200도 #5539에서 독립적으로 진단하고 수정)
- 화면 ID/이름 문자열을 복사하고 개수를 스냅샷으로 저장하여 디스플레이 변경 시 macOS 화면을 열거하다가 발생하는 충돌 수정(#5565, @x-haose도 #5584에서 독립적으로 진단하고 수정)
- 최소화/복원 전환 중 `GetClientRect`이 nil을 반환할 때 `WM_ERASEBKGND`이 단색 배경을 그리면서 Windows에서 발생하는 충돌 수정(@sinspired가 #5636에서 보호 처리 보고)
- 프런트엔드 바인딩 오류가 항상 텍스트로 파싱되는 문제 수정, 작성자: @mbaklor, #5690
- macOS 및 Linux의 해당 파일에는 이미 있는 `!server` 빌드 제약 조건이 Windows GUI 파일에 누락되어 `server` 빌드 태그 사용 시 발생하는 Windows 빌드 실패 수정(#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## 추가됨

- [PR](https://github.com/wailsapp/wails/pull/5720)에서 Update Manifest 프로토콜 및 엔드포인트 공급자 구현, 작성자: @taliesin-ai

## 변경됨

- `webview2` 바인딩을 v3 모듈에 `v3/internal/webview2`로 통합하고, 독립형 모듈과 해당 모듈의 nightly 릴리스/동기화 워크플로 및 go.mod 버전 조정 절차 제거(유일한 소비자는 v3), [PR](https://github.com/wailsapp/wails/pull/5711), 작성자: @taliesin-ai

## 수정됨

- [PR](https://github.com/wailsapp/wails/pull/5750)에서 WebView2 모니터 배율 감지 및 DPI 변경 시 호스트 재동기화 수정 사항을 Unreleased 섹션으로 이동, 작성자: @taliesin-ai
- [PR](https://github.com/wailsapp/wails/pull/5741)에서 float64 및 BOOL 매개변수의 WebView2 COM 마샬링 업데이트, 작성자: @wayneforrest
- [PR](https://github.com/wailsapp/wails/pull/5703)에서 Windows 시스템 트레이 아이콘 업데이트 및 제거 시 패닉과 nil 역참조 방지, 작성자: @wayneforrest
- [PR](https://github.com/wailsapp/wails/pull/5743)에서 Windows의 숨겨진 창이 올바르게 다시 숨겨지지 않는 문제 수정, 작성자: @wayneforrest
- [PR](https://github.com/wailsapp/wails/pull/5742)에서 WebView2 컨트롤러 가시성을 창 최소화/최대화/복원 상태와 동기화, 작성자: @wayneforrest

### 수정됨

- @randalmurphal이 검증한 수정 사항을 바탕으로 WebView2 모니터 배율 감지를 다시 활성화하고, DPI가 변경될 때만 호스트를 재동기화하도록 제한했습니다. @eleclin이 근본 원인을 검증하고 @qq540491950이 하드웨어에서 테스트했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5734)

## v3.0.0-alpha2.113 - 2026-07-04

## 추가됨

- `ANDROID_KEYSTORE_FILE`이 설정되지 않은 상태에서 릴리스 AAB를 빌드하면 경고하도록 하고(Google Play는 디버그 서명된 번들을 거부함), App Bundle 패키징 및 서명 방법을 문서화했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5730)
- 여러 언어로 Why Wails 문서를 추가했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5739)
- 바인딩에서 Go time.Time을 JS Date 또는 string으로 매핑할 수 있도록 지원합니다. @fbbdev의 [PR](https://github.com/wailsapp/wails/pull/5398)
- Play Store 제출용 Android App Bundle(AAB) 패키징 작업(`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`)을 추가했습니다. 로컬 및 에뮬레이터 테스트용 APK 작업은 그대로 유지됩니다. @mortenolsrud의 [PR](https://github.com/wailsapp/wails/pull/5728)([#5726](https://github.com/wailsapp/wails/issues/5726) 수정)
- Android 실제 기기용 작업 대상을 추가하고 카메라/위치 권한을 재개했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5735)

## 변경됨

- `webview2`을 v1.0.28로 업데이트했습니다([릴리스 정보](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Android 템플릿의 `compileSdk`/`targetSdk`을 34에서 신규 앱 제출 시 Google Play가 요구하는 35으로 업데이트했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5730)

## 수정됨

- sponsorkit에서 미리 적용된 아바타 마스크를 수정했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5745)
- 버전을 사전식으로 정렬하여 Android AVD 자동 생성 시 잘못된 시스템 이미지 또는 cmdline-tools 버전을 선택하던 문제를 수정했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5730)
- 설정 마법사가 오래된 Android NDK 버전을 제안하던 문제를 수정했습니다. 이제 문서에 명시된 요구 사항과 일치하는 26.3.11579264을 제안합니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5730)
- SvelteKit 및 옵션에 관한 프랑스어 문서를 업데이트했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5744)
- 디스플레이 변경 중 macOS 화면 열거에서 발생하던 SIGSEGV를 수정했습니다. @flofreud의 [PR](https://github.com/wailsapp/wails/pull/5516)

## v3.0.0-alpha2.112 - 2026-07-03

## 추가됨

- Go 기반 기여자 SVG 생성기를 추가하고 docs/website의 크레딧 페이지를 업데이트했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5724)
- 바인딩에서 Go time.Time을 JS Date 또는 string으로 매핑할 수 있도록 지원합니다. @fbbdev의 [PR](https://github.com/wailsapp/wails/pull/5398)

## 변경됨

- Node 기반 후원자 이미지 파이프라인을 Go 생성기로 교체했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5719)

## 수정됨

- Android 빌드 자산 종속성 설치 스크립트를 수정했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5729)
- `ValidateAndSanitizeURL`에서 U+0085(NEXT LINE) 제어 문자를 거부하도록 하여 URL 유효성 검사기의 공백 문자 검사를 완성했습니다.
- 프레임 없는 창에서 DPI 변경 시 DWM 프레임을 다시 계산하도록 했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/4785)
- Windows에서 배율이 100%가 아닐 때 DnD 드롭 영역 감지가 실패하던 문제를 수정했습니다. @yulesxoxo의 [PR](https://github.com/wailsapp/wails/pull/4632)
- Darwin의 대화 상자, 메뉴, 트레이 및 알림 전반에서 Cocoa 객체에 명시적인 Objective-C 메모리 관리를 추가했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5714)
- Linux CGO 백엔드 버그와 시스템 트레이 문제를 수정했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5718)

## v3.0.0-alpha2.111 - 2026-07-01

## 추가됨

- 커뮤니티 쇼케이스에 HappyTools를 추가했습니다. @Aliuyanfeng의 [PR](https://github.com/wailsapp/wails/pull/5061)
- 인도네시아어 로캘 지원과 포괄적인 문서를 추가했습니다. @triadmoko의 [PR](https://github.com/wailsapp/wails/pull/5643)
- WindowsWindow에 DisableMenu 옵션을 추가했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/4813)

## 변경됨

- GOOS 및 ARCH에 따라 build/package 작업을 디스패치하도록 Taskfile 템플릿과 CLI를 업데이트했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5617)

## 수정됨

- macOS 창 탭 기능 문제를 수정했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5708)

## 제거됨

- 기여, 기능 및 가이드에서 독일어로 번역된 MDX 파일을 제거했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5702)

## v3.0.0-alpha2.110 - 2026-06-30

## 추가됨

- macOS WebView의 다시 로드 및 강제 다시 로드를 구현하고, WebContent 프로세스 종료 시 복구 기능을 추가했습니다. @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5129)
- 기여, 기능 및 가이드에 관한 포괄적인 독일어 문서를 추가했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5396)
- 소리, 첨부 파일, 예약 및 업데이트 API를 추가하여 알림 기능을 개선했습니다. @popaprozac의 [PR](https://github.com/wailsapp/wails/pull/5333)

## 수정됨

- 프레임 없는 창에서 DPI 변경 시 DWM 프레임을 다시 계산하도록 했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/4785)
- Windows에서 배율이 100%가 아닐 때 DnD 드롭 영역 감지가 실패하던 문제를 수정했습니다. @yulesxoxo의 [PR](https://github.com/wailsapp/wails/pull/4632)

## v3.0.0-alpha2.109 - 2026-06-29

## 추가됨

- EventsEmit 문서에 코드 예제를 추가했습니다. @iamhabbeboy의 [PR](https://github.com/wailsapp/wails/pull/5026)
- Windows WebView2 비주얼 호스팅 옵션을 추가했습니다. @MerIijn의 [PR](https://github.com/wailsapp/wails/pull/5380)
- 커뮤니티 쇼케이스 문서에 Klustr를 추가했습니다. @SametKUM의 [PR](https://github.com/wailsapp/wails/pull/5536)
- 새 페이지 및 변경 로그 항목과 함께 Kira를 커뮤니티 쇼케이스에 추가했습니다. @thiennguyen93의 [PR](https://github.com/wailsapp/wails/pull/5685)
- MCP 서비스 가이드에 피드백 섹션을 추가했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5694)

## 변경됨

- 이제 서버 모드에는 데스크톱 빌드 작업과 일관된 정식 프로덕션 빌드가 있습니다(#5693). `task build:server`은 기본적으로 프로덕션 바이너리(`-tags server,production`, `-trimpath`, 심벌 제거됨)를 빌드하며, `DEV=true`(개발 서버), `OBFUSCATED=true`(garble) 및 `EXTRA_TAGS`을 허용합니다. `task run:server`은 개발 서버를 실행합니다. `Dockerfile.server` / `task build:docker`은 프로덕션 프런트엔드를 먼저 빌드한 다음 프로덕션 서버(`-tags server,production`)를 빌드합니다. 이미지는 기본적으로 distroless/static 기반의 순수 Go 정적 빌드를 사용하며, CGO 앱을 위해 `CGO_ENABLED`, `GO_IMAGE` 및 `RUNTIME_IMAGE`을 재정의할 수 있는 빌드 인수로 노출합니다.

## 수정됨

- 대기 중인 비동기 호출이 있는 창을 닫을 때 충돌하지 않도록 했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/4435)
- Windows에서 숨겨진 앱을 열 때 창이 활성화되지 않도록 했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5249)
- WebKit 요청 메타데이터, 응답 완료 및 본문 스트림 처리가 GTK 메인 스레드에서 실행되도록 보장했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5668)
- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제를 수정했습니다(#5659, @puneetdixit200도 #5539에서 독립적으로 진단하고 수정함).
- 디스플레이 변경 시 macOS 화면을 열거하는 과정에서 발생하는 충돌을 수정했습니다. 화면 ID/이름 문자열을 복사하고 개수를 스냅샷으로 저장합니다(#5565, @x-haose도 #5584에서 독립적으로 진단하고 수정함).
- Windows에서 서로 다른 DPI의 모니터 사이로 창을 드래그한 후 WebView2 콘텐츠가 축소되었다가 사라지는 문제를 수정했습니다. 최소화 해제 시 수행하는 DPI 재동기화와 동일하게 `WM_DPICHANGED` 핸들러에서 컨트롤러 경계를 다시 설정합니다(#5677).

## v3.0.0-alpha2.108 - 2026-06-28

## 추가됨

- `app.GlobalShortcut`을 통해 전역(시스템 전체) 키보드 단축키를 추가했습니다(`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). 애플리케이션에 포커스가 없어도 단축키가 작동합니다. 타사 종속성 없이 플랫폼별 네이티브 방식으로 구현했습니다. macOS에서는 Carbon 단축키, Windows에서는 `RegisterHotKey`, X11에서는 `XGrabKey`, Wayland에서는 XDG Desktop Portal 전역 단축키 인터페이스를 사용합니다.
- 내장 MCP 서버를 추가했습니다. 애플리케이션을 `mcp` 태그로 빌드하면 자동으로 시작되는 Model Context Protocol 서버로, LLM 에이전트가 실행 중인 Wails 애플리케이션을 테스트하고 제어할 수 있습니다. 창 제어, DOM 검사, JavaScript 평가, 바인딩된 메서드 호출, 이벤트 및 화면에 애니메이션 커서로 표시되는 마우스/키보드 입력 시뮬레이션을 지원합니다. 사용자 코드는 필요하지 않습니다. `WAILS_MCP=1`이 설정되면 `wails3 build`/`wails3 dev`이 `mcp` 태그를 자동으로 추가합니다. 모든 설정은 환경 변수(`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`)로 구성합니다.

## 수정됨

- GTK4 Linux에서 `Menu.Update()`이 네이티브 메뉴를 다시 빌드하지 않는 문제를 수정했습니다(#5659, @puneetdixit200도 #5539에서 독립적으로 진단하고 수정함).
- 디스플레이 변경 시 macOS 화면을 열거하는 과정에서 발생하는 충돌을 수정했습니다. 화면 ID/이름 문자열을 복사하고 개수를 스냅샷으로 저장합니다(#5565, @x-haose도 #5584에서 독립적으로 진단하고 수정함).

## v3.0.0-alpha2.107 - 2026-06-27

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5613)에서 사이드바 탐색 기능이 포함된 실험적 Wake 문서를 추가했습니다.

## v3.0.0-alpha2.106 - 2026-06-24

## 변경됨

- `webview2`을 v1.0.27(으)로 올렸습니다.
  - ci(webview2): 릴리스 빌드 수정(Windows 교차 컴파일 + 완전한 go.sum)(#5671)\

  **전체 변경 사항:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/5672)에서 webview2 릴리스 워크플로의 교차 컴파일 과정에서 go vet을 제거했습니다.
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/5670)에서 자동 변경 로그의 OpenRouter 모델을 google/gemini-2.5-flash-lite로 업데이트했습니다.
- `webview2`을 v1.0.26(으)로 올렸습니다.

### 수정 사항

- **일시적인 런타임 COM 오류가 발생해도 종료하지 않고 복구하도록 수정했습니다**(#5658, #5580). 이전에는 `Chromium.errorCallback`이 *모든* COM 오류에 대해 `os.Exit(1)`을 호출했으므로, 시작 후 복구 가능한 일시적 오류가 발생해도 애플리케이션 전체가 종료되었습니다. 이제 런타임 경로(`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`)에서는 오류를 기록하고 복구합니다. 특히 `MessageReceived`에서 형식이 잘못되었거나 신뢰할 수 없는 웹 메시지가 들어오면 프로세스를 종료하는 대신 해당 메시지를 버립니다. 이 변경으로 서로 다른 DPI의 모니터 사이를 이동할 때 발생하는 충돌 유형을 해결했습니다(#5544, #5650). 환경/컨트롤러 생성 경로의 오류는 계속 치명적 오류로 처리합니다.\

**전체 변경 사항:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 수정됨

- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/5671)에서 go.sum 파일을 올바르게 처리하도록 release-webview2 워크플로를 수정했습니다.
- @taliesin-ai가 [PR](https://github.com/wailsapp/wails/pull/5659)에서 네이티브 메뉴를 지우고 다시 빌드하여 Linux GTK4 메뉴 업데이트 문제를 수정했습니다.

## v3.0.0-alpha2.105 - 2026-06-21

## 추가됨

- 공유 코드에서 런타임 플랫폼을 감지할 수 있도록 `application.System`을 추가했습니다. `System.IsMobile()`(iOS/Android), `System.IsDesktop()`(macOS/Windows/Linux), `System.IsServer()`(`server` 빌드 태그) 및 단일 대상을 직접 테스트하는 `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`을 제공합니다. 모든 대상에서 컴파일되므로 빌드 태그 없이 분기할 수 있습니다. 이에 대응하는 프런트엔드 도우미(`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`)는 `@wailsio/runtime`에서 사용할 수 있습니다.
- 자체 Vite 프로젝트를 `frontend/`에 배치하는 방법을 설명하는 "다른 프런트엔드 프레임워크 사용" 가이드를 추가했습니다(Solid, Preact, Lit, SvelteKit, Qwik, Angular 등 포함).
- 이제 `wails3 setup` 마법사가 모바일(iOS/Android) 도구 체인을 검사합니다. Xcode와 iOS Simulator 런타임, JDK, Android SDK/NDK 및 에뮬레이터를 검사하며, 해당되는 경우 한 번의 클릭으로 설치하고 복사 가능한 셸 구성 수정 사항을 제공합니다.
- 생성된 프로젝트에는 7일의 `minimum-release-age`을 설정하는 `frontend/.npmrc`이 포함됩니다. 새로 게시된 패키지(침해되었을 가능성이 있는 패키지)에 노출될 위험을 줄이기 위한 설정입니다(pnpm과 bun은 이 설정을 준수하며 npm은 문제없이 무시함).

## 변경됨

- 모든 내장 시작 템플릿을 새로운 네온 산맥 히어로 디자인으로 개편했습니다(웹, iOS 및 Android).
- **이제 TypeScript가 시작 템플릿의 기본값이며 접미사 없는 템플릿 이름을 사용합니다.** `wails3 init`(`-t` 없음)은 TypeScript 프로젝트를 스캐폴딩합니다. `-t vanilla`, `-t react`, `-t vue` 및 `-t svelte`은 TypeScript 템플릿이며, JavaScript 변형은 `-t vanilla-js`, `-t react-js`, `-t vue-js` 및 `-t svelte-js`입니다. 내장 템플릿은 `template.yaml`의 `typescript:`으로 언어를 선언합니다. `-ts` 접미사를 사용하는 커뮤니티 템플릿도 대체 방식으로 계속 작동합니다.
- 네온 "digital Wails" 테마(산맥 배경 위의 반투명 유리 효과)로 `wails3 setup` 마법사를 개편했습니다.

## 수정됨

- WebView2가 일시 중단되거나 렌더링/GPU 프로세스가 재활용될 만큼 오랫동안 최소화된 앱을 Windows에서 복원할 때 발생하는 충돌을 수정했습니다. 이제 최소화/복원 시 DPI 재동기화(#5544)는 창의 DPI가 실제로 변경된 경우에만 WebView2 컨트롤러에 접근합니다. 따라서 일반적인 동일 DPI 복원 시 일시 중단된 컨트롤러에 치명적인 COM 호출을 하지 않습니다(#5605).
- 자산/미디어를 자주 로드하며 장시간 실행되는 Linux 앱에서 네이티브 `SIGABRT`/`SIGSEGV` 충돌이 반복되는 문제를 수정했습니다. 일반적으로 GTK 메인 루프 중 `g_object_unref` 내부에서 발생했습니다. 자산 서버가 워커 goroutine에서 `WebKitURISchemeRequest`을 완료하면서 스레드에 안전하지 않은 WebKit2GTK 함수를 GTK 메인 스레드 밖에서 호출했습니다. 이제 완료 처리(`webkit_uri_scheme_request_finish_with_response`/`finish_error`)는 메인 스레드에서 실행됩니다. #5566의 부분 수정 사항을 완성합니다. GTK3 및 GTK4/WebKitGTK 6.0 빌드 모두에 영향을 줍니다(#5631, #5557).
- Linux/GTK3의 `setupSignalHandlers`에서 간헐적으로 발생하는 `fatal error: invalid pointer found on stack`을 수정했습니다. 신호 `user_data`로 전달된 창 ID가 Go의 로컬 `unsafe.Pointer`에 보관되어, 스택 복사 중 가비지 컬렉터가 이 값(포인터가 아닌 값)을 스캔할 때 중단되었습니다. 이제 Go 측에서는 ID를 정수 형식(`uintptr_t`)으로 유지합니다. 이는 #4958에서 GTK4 경로에 적용한 것과 같은 수정 사항을 기존 GTK3 경로로 백포트한 것입니다. 당시 GTK4 경로에서는 `-race`/checkptr 오류를 없애기 위해 C 신호 함수를 `uintptr_t`으로 변경했습니다(#5631).

## 제거됨

- `react-swc`, `preact`, `lit`, `solid`, `qwik` 및 `sveltekit` 시작 템플릿과 각각의 `-ts` 변형을 제거했습니다. 이제 지원되는 내장 템플릿은 `vanilla`, `react`, `vue` 및 `svelte`입니다. 각 템플릿은 기본적으로 TypeScript를 사용하며 `-js` JavaScript 변형도 제공합니다. 다른 프레임워크도 [자체 프런트엔드 가져오기](https://v3.wails.io/guides/dev/frontend-frameworks) 또는 사용자 지정 템플릿을 통해 계속 사용할 수 있습니다.

## v3.0.0-alpha2.104 - 2026-06-18

## 수정됨

- 바인딩된 Go 서비스 메서드가 빈 문자열을 반환할 때 발생하던 iOS 충돌(SIGABRT)을 수정했습니다. iOS 애셋 응답 작성기가 본문 길이 대신 `buf != nil`으로 본문 포인터를 검사했기 때문에 길이가 0인 본문에서 `&buf[0]`이 패닉을 일으켰습니다. 이제 데스크톱 작성기와 동일하게 길이를 기준으로 검사합니다.

## v3.0.0-alpha2.103 - 2026-06-15

## 변경됨

- iOS 및 Android 네이티브 기능을 플랫폼 관리자로 이전했습니다. 이제 기존 `application.IOS*`/`application.Android*` 자유 함수 대신 `application.IOS.*` 및 `application.Android.*`을 통해 호출합니다(예: `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`)(#5602).
- 모바일 브리지 이벤트의 이름을 변경했습니다. 이제 크로스 플랫폼 이벤트에는 `common:*` 접두사를 사용하고(예: `common:haptic`, `common:location`), 플랫폼 전용 이벤트에는 `ios:*` / `android:*`을 사용합니다(예: `ios:backgroundTask`, `android:foregroundService`). `native:*` 접두사는 더 이상 사용하지 않습니다(#5602).

## v3.0.0-alpha.102 - 2026-06-14

## 추가됨

- 대화형 프로젝트 설정 및 종속성 검사를 위한 실험적 `wails3 setup` 마법사를 추가했습니다.
- 기계 판독 가능한 출력을 위해 `wails3 doctor`에 `--json` 플래그를 추가했습니다.
- `wails3 doctor` 명령에 서명 상태 섹션을 추가했습니다.

## 수정됨

- Linux에서 npm을 감지할 때 패키지 관리자뿐만 아니라 PATH도 확인하도록 수정했습니다.

## v3.0.0-alpha.101 - 2026-06-13

## 추가됨

- iOS: 네이티브 메시지 대화 상자(UIAlertController)와 파일/여러 파일/디렉터리 열기 대화 상자(UIDocumentPickerViewController)를 추가했습니다. 저장 대화 상자는 명시적인 오류를 반환합니다.
- iOS: UIPasteboard를 통한 클립보드 지원을 추가했습니다.
- iOS: UIScreen을 통한 실제 화면 측정값(포인트, 픽셀, 배율, 안전 영역 작업 공간)을 추가했습니다.
- iOS: 기기 빌드(`IOS_PLATFORM=device`), 코드 서명 ID / 프로비저닝 프로파일 / 엔타이틀먼트 지원, `.ipa` 패키징 및 devicectl을 통한 `deploy-device`을 추가했습니다.
- iOS: 구성 가능한 최소 iOS 버전(build/config.yml의 `ios.minIOSVersion`)을 추가했습니다.
- iOS: `wails3 doctor`이 macOS에서 Xcode 및 iOS SDK의 사용 가능 여부를 보고합니다.
- iOS: 배터리, 네트워크, 테마, 화면 잠금 및 메모리 부족 시스템 이벤트가 `events.IOS.*` 및 플랫폼 중립적인 `events.Common.*` 애플리케이션 이벤트로 노출됩니다.
- iOS: 네이티브 모바일 기능 브리지(내보낸 `application.IOS*`) — 공유 시트, URL 열기, 절전 방지, 손전등, 안전 영역 인셋, 밝기, 앱 정보, 방향 잠금, 상태 표시줄, 생체 인증(Face ID/Touch ID), 로컬 알림 및 Keychain 보안 저장소를 추가했습니다.
- iOS: 센서 및 하드웨어 — 햅틱, 일회성 위치 정보, 가속도계, 근접 센서, 텍스트 음성 변환, 저장 공간 정보, 전원/배터리 상태, 네트워크 상태, 키보드 인셋 및 화면 캡처 감지를 추가했습니다.
- iOS: 문서(IOS.md 및 문서 사이트 가이드)를 추가했습니다.
- Android: 네이티브 메시지 대화 상자(AlertDialog)와 파일/여러 파일 열기 대화 상자(Storage Access Framework, 캐시 복사본으로 가져옴)를 추가했습니다. 디렉터리 열기 및 저장 대화 상자는 명시적인 오류를 반환합니다.
- Android: ClipboardManager를 통한 클립보드 지원을 추가했습니다.
- Android: WindowMetrics/DisplayMetrics를 통한 실제 화면 측정값(dp, 픽셀, 배율, 시스템 표시줄 작업 영역)을 추가했습니다.
- Android: 햅틱(`Android.Haptics.Vibrate`), 기기 정보(`Android.Device.Info`) 및 토스트(`Android.Toast.Show`) 런타임 메서드를 추가했습니다.
- Android: `ActivityCreated`이 `Common.ApplicationStarted`에 매핑되는 형식 지정 수명 주기 이벤트(`events.Android.*`, events.txt에서 생성)를 추가했습니다.
- Android: 빌드 파이프라인이 설치 가능한 디버그 및 릴리스 APK(`android:run`, `android:package`, `android:package:fat`)를 생성합니다. 릴리스 서명에는 기본적으로 디버그 키 저장소를 사용하며, `ANDROID_KEYSTORE_*` 환경 변수를 통해 실제 키 저장소를 사용할 수도 있습니다.
- Android: `wails3 doctor`이 Android SDK, NDK 및 JDK를 보고합니다.
- Android: 배터리, 네트워크, 테마, 화면 잠금 및 메모리 부족 시스템 이벤트가 `events.Android.*` 및 플랫폼 중립적인 `events.Common.*` 애플리케이션 이벤트로 노출됩니다.
- Android: 네이티브 모바일 기능 브리지(내보낸 `application.Android*`) — 공유, URL 열기, 절전 방지, 손전등, 안전 영역 인셋, 밝기, 앱 정보, 방향 잠금, 상태 표시줄, 생체 인증(BiometricPrompt), 로컬 알림 및 EncryptedSharedPreferences 보안 저장소를 추가했습니다.
- Android: 센서 및 하드웨어 — 햅틱, 일회성 위치 정보, 가속도계, 근접 센서, 텍스트 음성 변환, 저장 공간 정보, 전원/배터리 상태, 네트워크 상태, 키보드 인셋 및 FLAG_SECURE 화면 캡처 차단을 추가했습니다.
- Android: 문서(ANDROID.md 및 문서 사이트 가이드)를 추가했습니다.
- 예제: `mobile` 키친 싱크에 iOS와 Android의 네이티브 기능 브리지를 시연하는 Mobile 및 Hardware 탭을 추가했습니다(필 모양 탭은 여러 행으로 줄바꿈됩니다).
- 모바일: 배터리 — 앱이 백그라운드로 전환되면 가속도계, 근접 센서, 손전등 및 예제의 주기적 시계를 일시 중지하고, 앱으로 돌아오면 복원합니다(Android에서는 프로세스가 백그라운드에서 계속 실행되며, iOS에서 손전등은 지속되는 하드웨어 상태입니다). 또한 Android 시스템 이벤트 수신기는 앱이 포그라운드에 있을 때만 등록됩니다.
- iOS: 카메라 캡처 — `application.IOSCapturePhoto`/`IOSCaptureVideo`(UIImagePickerController → base64 썸네일이 포함된 `native:capture` 이벤트)을 추가했습니다.
- iOS: 백그라운드 실행 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`(UIApplication 백그라운드 작업 시간 창)과 생성된 Info.plist에 `UIBackgroundModes`을 템플릿으로 삽입하는, 구성 가능한 `ios.backgroundModes`(build/config.yml)을 추가했습니다.
- Android: 카메라 캡처 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`(FileProvider를 통한 시스템 카메라 → `native:capture` 이벤트)을 추가했습니다.
- Android: 포그라운드 서비스 — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`(지속 알림이 있는 `WailsForegroundService`이 장시간 실행되는 백그라운드 작업을 위해 프로세스를 계속 실행함)을 추가했습니다.
- 예제: 사진/동영상 캡처와 백그라운드 실행(Android에서는 포그라운드 서비스, iOS에서는 백그라운드 작업 시간 창)을 시연하는 Camera 탭을 추가했습니다.

## 수정됨

- Linux에서 `getUserMedia`이 항상 `NotAllowedError` 오류로 실패하던 문제를 수정했습니다. WebKitGTK는 아무도 처리하지 않는 권한 요청을 거부하는데 `permission-request` 신호가 연결되지 않았습니다. 이제 새로운 크로스 플랫폼 `WebviewWindowOptions.Permissions` 맵(`map[PermissionType]Permission`)에 따라 카메라/마이크를 처리하며, Linux(WebKitGTK)와 Windows(WebView2) 모두에서 이를 적용합니다. 네이티브 프롬프트가 없는 Linux에서는 카메라/마이크를 기본적으로 허용하여 `getUserMedia`을 복원하며, `PermissionDeny`으로 비활성화할 수 있습니다(#5552).
- iOS: `GOOS=ios`이 다시 컴파일됩니다(내보낸 `events.IOS`, 모바일 메서드 이름 스텁). 프로덕션 태그가 지정된 빌드도 컴파일됩니다(pkg/application 및 여러 서비스의 빌드 태그 수정).
- iOS: 이제 Go→JS 이벤트와 ExecJS가 작동합니다. 시작 시 페이지가 더 이상 두 번 로드되지 않으며 `wails:runtime:ready` 핸드셰이크가 더 이상 유실되지 않습니다.
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted`이 더 이상 앱 시작과 경합하지 않습니다. 고정된 2초 시작 대기 시간을 제거했습니다.
- iOS: Go→JS JavaScript를 실행할 때마다 발생하던 C 문자열 누수를 수정했습니다.
- iOS: 이제 `hasListeners`이 실제 리스너 등록 상태를 반영합니다.
- iOS: 프로덕션 빌드에서 프레임워크 디버그 로깅이 컴파일 단계에서 제외됩니다.
- Android: `GOOS=android`이 다시 컴파일됩니다. `events.Android`을 정의하고, 범위를 벗어나던 `events_android.go` 리스너 배열을 제거하고, 모바일 메서드 이름 스텁을 추가했으며, 데스크톱 Linux 파일(`linux_cgo.*`, `events_linux.*`, `environment_linux.go`)이 Android 빌드에 포함되지 않도록 했습니다.
- Android: 이제 JS→Go 바인딩이 작동합니다. WebView는 `fetch()` POST 본문을 `shouldInterceptRequest`에 전달할 수 없으므로, 런타임 호출은 nil 요청 본문에서 충돌하는 대신 JavascriptInterface 전송 계층(`nativeHandleRuntimeCall`)을 통해 라우팅됩니다.
- Android: `Screens.*` 런타임 호출이 실제 데이터를 반환합니다. 이제 시작 시 ScreenManager가 채워집니다(이전에는 연결된 적이 없어 `GetAll`이 nil을 반환했습니다).
- Android: 프로덕션 빌드에서는 프레임워크 디버그 로깅이 컴파일에서 제외되며, 디버그 빌드에서는 `Wails` 태그로 logcat에 기록됩니다.
- Android: 실제 `hasListeners` 레지스트리, JNI 참조/예외 처리 및 한 번만 로드되는 페이지 수명 주기(이중 탐색 없음)
- Vite 개발 서버가 실행 중일 때 Windows에서 `wails3 generate bindings`이(가) "Access is denied" 오류로 실패하는 문제를 수정했습니다. 생성된 파일을 출력 디렉터리 위로 이름을 바꾸는 대신 해당 디렉터리에 동기화합니다(#5515).
- 디스플레이 변경 후 화면 정보를 읽을 때 macOS에서 간헐적으로 치명적인 충돌이 발생하는 문제를 수정했습니다. 화면 ID와 이름이 자동 해제되는 `UTF8String` 버퍼를 가리키는 포인터를 저장하여 Go가 복사하기 전에 버퍼가 해제될 수 있었습니다(use-after-free). 이제 문자열에 `strdup`을(를) 적용하고 변환 후 해제하며, 화면 열거를 명시적인 autorelease pool에서 실행하므로 Go goroutine에서 호출해도 더 이상 누수가 발생하지 않습니다(#5556).
- assetserver가 `WebKitURISchemeRequest`을(를) 닫을 때 Linux에서 간헐적으로 SIGSEGV가 발생하는 문제를 수정했습니다. 마지막 `g_object_unref`이(가) assetserver goroutine에서 실행되면서 GTK 메인 스레드 외부에서 WebKit GObject를 마무리했습니다. 이제 unref는 `g_main_context_invoke`을(를) 통해 GTK 메인 컨텍스트로 전달됩니다(#5557).

## v3.0.0-alpha.100 - 2026-06-13

## 추가됨

- `MacWebviewPreferences`에 WKWebView 구성 옵션 `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize` 및 `ApplicationNameForUserAgent`을(를) 추가했습니다(#5549).

## 수정됨

- Vite 개발 서버가 실행 중일 때 Windows에서 `wails3 generate bindings`이(가) "Access is denied" 오류로 실패하는 문제를 수정했습니다. 생성된 파일을 출력 디렉터리 위로 이름을 바꾸는 대신 해당 디렉터리에 동기화합니다(#5561).
- Linux의 프레임 없는 창에서 JS 크기 조정 이벤트가 발생하지 않는 문제와 프레임 없는 창의 스크롤바 가장자리 감지 문제를 수정했습니다(#5368).
- 임시 디렉터리와 설치 디렉터리가 서로 다른 볼륨에 있을 때 Windows의 업데이터가 "invalid cross-device link" 오류로 실패하는 문제를 수정했습니다(#5560).

## v3.0.0-alpha.99 - 2026-06-10

## 수정됨

- Vite 개발 서버가 실행 중일 때 Windows에서 `wails3 generate bindings`이(가) "Access is denied" 오류로 실패하는 문제를 수정했습니다. 생성된 파일을 출력 디렉터리 위로 이름을 바꾸는 대신 해당 디렉터리에 동기화합니다(#5515).

## v3.0.0-alpha.98 - 2026-06-03

## 수정됨

- `SIGUSR1`에 `SA_ONSTACK`을(를) 더 이상 강제하지 않도록 하여 유휴 상태(예: 검사기가 열려 있을 때)에서 Linux의 WebKit UI가 멈추는 문제를 수정했습니다. 이 설정은 JavaScriptCore의 GC 스레드 동기화를 손상시켰습니다(#5527).

## v3.0.0-alpha.97 - 2026-05-31

## 추가됨

- 디버깅 페이지와 `runtime/trace` 작업 방법을 추가했습니다.

## 변경됨

- 불필요한 `_ "embed"` import 일부를 제거하여 코드를 약간 정리했습니다.

## 수정됨

- Windows에서 창의 최대화를 해제한 후 최소 너비/높이 제약 조건이 적용되지 않는 문제를 수정했습니다(#4593).
- Frameless + Transparent 창 옵션을 사용한 전체 화면 모드에서 마우스 클릭이 창을 통과하는 문제를 수정했습니다(#4408).

## v3.0.0-alpha.96 - 2026-05-25

## 추가됨

- Garble 난독화 지원을 추가했습니다([#4563](https://github.com/wailsapp/wails/issues/4563)). 안정적인 바인딩 메서드 ID, 빌드/Taskfile 연동(`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), 그리고 런타임에 노출되는 모든 페이로드의 JSON 구조체 태그(`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`)를 추가하여 Garble이 내보낸 필드의 이름을 변경해도 전송 형식이 유지되도록 했습니다.

## v3.0.0-alpha.95 - 2026-05-20

## 추가됨

- 누락된 프로젝트 구조 페이지를 추가했습니다.

## 변경됨

- 문서: 더 깔끔하게 표시되도록 아키텍처 페이지의 다이어그램 두 개를 시퀀스 다이어그램으로 변경했습니다.
- 문서: 실행을 위한 필수 조건으로 D2를 설치해야 한다는 참고 사항을 추가했습니다.

## 수정됨

- GTK4 기본 구성에서 `wails3 generate appimage`을(를) 수정했습니다. 이제 번들러가 런타임 파일을 검색하기 전에 바이너리에서 GTK 스택을 감지하므로 GTK4 빌드에는 `webkitgtk-6.0/` 아래의 `libwebkitgtkinjectedbundle.so`을(를), `-tags gtk3` 빌드에는 `webkit2gtk-4.1/` 아래의 `libwebkit2gtkinjectedbundle.so`을(를) 선택합니다. `.relr.dyn` 프로브는 `libgtk-4.so.1`도 확인하므로 스택과 관계없이 최신 툴체인에서 스트리핑이 올바르게 비활성화됩니다. (#5475)
- 상대 `-builddir`을(를) 사용해 호출할 때 `wails3 generate appimage`이(가) 실패하는 문제를 수정했습니다. 이제 번들러가 `-binary`, `-icon`, `-desktopfile`, `-builddir` 및 `-outputdir`을(를) 처음부터 절대 경로로 해석하므로, 처리 도중의 `s.CD`이(가) AppRun 다운로드 goroutine이나 복사 후 `ldd` 프로브를 손상시키지 않습니다.
- 데스크톱 `Name=` 필드가 바이너리 기본 이름과 일치하지 않을 때 `wails3 generate appimage`이(가) 최종 AppImage를 `-outputdir`(으)로 이동하지 못하는 문제를 수정했습니다. 이제 번들러는 `OUTPUT` 환경 변수를 통해 linuxdeploy의 appimage 플러그인이 데스크톱 파일에서 파생된 이름 대신 `<binary>-<arch>.AppImage`에 AppImage를 쓰도록 강제합니다.
- alpha.93에서 GTK4 + WebKitGTK 6.0 스택이 기본값으로 승격된 후 Linux에서 `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` 및 `Common.SystemDidWake`이(가) 발생하지 않는 문제를 수정했습니다. 새 기본 `application_linux.go` `run()`은(는) `setupCommonEvents()`(`Linux.*` 이벤트를 해당 `Common.*` 이벤트로 전달함)이나 `monitorPowerEvents()`을(를) 호출하지 않았습니다. 이제 `application_linux_dbus.go`을(를) 통해 GTK3 및 GTK4 빌드 경로가 DBus 전원 모니터 도우미를 공유합니다. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## 수정됨

- alpha.93에서 GTK4 + WebKitGTK 6.0 스택이 기본값으로 승격된 후 Linux에서 `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep` 및 `Common.SystemDidWake`이(가) 발생하지 않는 문제를 수정했습니다. 새 기본 `application_linux.go` `run()`은(는) `setupCommonEvents()`(`Linux.*` 이벤트를 해당 `Common.*` 이벤트로 전달함)이나 `monitorPowerEvents()`을(를) 호출하지 않았습니다. 이제 `application_linux_dbus.go`을(를) 통해 GTK3 및 GTK4 빌드 경로가 DBus 전원 모니터 도우미를 공유합니다. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## 추가됨

- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE`을(를) 추가했습니다.

## 수정됨

- appmenu-gtk-module이 실체화되지 않은 창에 접근하여 Wayland에서 창 메뉴가 충돌하는 문제를 @leaanthony가 수정했습니다(#4769).
- 앱 이름에 잘못된 문자(공백, 괄호 등)가 포함될 때 GTK 애플리케이션이 충돌하는 문제를 @leaanthony가 수정했습니다.
- Windows에서 드래그 앤 드롭을 초기화할 때 "not enough memory" 오류가 발생하는 문제를 @overlordtm이 수정했습니다(#4701).
- map 항목을 삭제할 때 잘못된 RLock을 사용하여 mainthread 콜백 저장소에서 발생하던 경쟁 조건을 @leaanthony가 수정했습니다(Linux, macOS, iOS)(#4424).
- 작업에 명령줄 인수를 전달할 때의 변수 처리를 수정했습니다. 이제 KEY=VALUE 쌍으로 지정한 CLI 변수가 올바르게 초기화되고 작업 실행 전반에 전달됩니다.
- macOS의 NSWindowZoomButton 충돌을 수정했습니다. 이제 `MaximiseButtonState` 및 `FullscreenButtonState`은(는) 시작 시점과 런타임 모두에서 더 제한적인 상태를 적용하며, 어느 setter도 다른 setter의 설정을 암묵적으로 재정의할 수 없습니다(#5319).
- CodeRabbit이 #5463에서 발견한 기존 GTK3 빌드 경로(`-tags gtk3`)의 여러 기존 버그를 수정했습니다. 이제 파일 연결을 통한 실행 시 시작 핸들러를 건너뛰지 않으며, `getTheme`은 범위 및 형식 면에서 안전합니다. `appName`은 더 이상 GLib 소유 메모리를 해제하지 않고, `clipboardGet`는 더 이상 GTK가 반환한 `gchar*`을 누수하지 않습니다. 이제 `Calloc`은 포인터 리시버를 사용하고 `NewCalloc`은 `*Calloc`을 반환하므로 풀이 실제로 할당을 추적합니다. `zoomOut`은 1.0로 제한되던 음수 승수 대신 `zoomInFactor`의 역수를 사용합니다. `execJS`는 호출할 때마다 `C.CString("")`을 누수하지 않고 미리 할당된 빈 월드 이름을 재사용합니다. `menuItem.setAccelerator`에서 개발용 `fmt.Println`도 제거했습니다. #5465을 해결합니다.
- 기본 GTK4 빌드 경로(`linux_cgo.go`)에서도 동일한 `Calloc` 값 리시버 누수를 수정했습니다. 포인터 리시버와 `NewCalloc() *Calloc`을 사용해 창별 `c.String(...)` 할당을 실제로 추적하고 해제합니다.

## v3.0.0-alpha.92 - 2026-05-15

## 추가됨

- `PACKAGE_MANAGER` 옵션으로 사용할 프런트엔드 패키지 관리자를 제어할 수 있도록 Taskfile을 수정했습니다.
- Taskfile 템플릿을 더 예측 가능하게 작성할 수 있도록 템플릿 데이터에 `{{.Opn}}` 및 `{{.Cls}}`을 추가했습니다.

## 변경됨

- 기존 Taskfile 중 일부가 `{{.Opn}} and {{.Cls}}`을 사용하도록 수정했습니다.

## 수정됨

- 패널이 트레이 메뉴를 읽는 중에 메뉴가 업데이트될 때 `linuxSystemTray`에서 발생하는 `concurrent map read and map write` 런타임 치명적 오류를 수정했습니다.
- Windows에서 앱을 연결된 콘솔 없이 실행해도 메시지가 손실되지 않도록 WebView2 오류와 스택 추적 출력에 `fmt` 대신 `log`을 사용합니다.

## v3.0.0-alpha.91 - 2026-05-12

## 변경됨

- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5414)에서 스폰서 SVG를 업데이트했습니다.
- **호환성 중단 변경(macOS):** `GetScreens`, `Position` 및 `SetPosition`가 모두 동일한 공간, 즉 논리적 포인트와 아래쪽으로 증가하는 Y축을 사용하고 기본 화면의 왼쪽 위를 `(0,0)`로 삼도록 macOS 좌표계를 정규화했습니다. 이는 Windows, GTK 및 Electron과 웹의 공개 API와 일치합니다. 이제 기본 화면보다 물리적으로 위에 있는 화면은 음수 `Bounds.Y`을 보고하며(이전에는 양수), `Position()`/`SetPosition()` 값은 `points × primaryScale` 대신 논리적 포인트 단위를 사용합니다. `Position()` → `SetPosition()` 왕복 변환은 그대로 유지됩니다. 이전 알파 빌드에서 기록한 절댓값이나 직접 계산한 우회 처리(예: `primaryScale`를 곱하거나 화면 높이를 기준으로 Y축을 뒤집는 처리)는 업데이트해야 합니다. [#5117](https://github.com/wailsapp/wails/issues/5117)을 해결합니다.

## 수정됨

- 패닉을 방지하도록 DBus 신호 이름과 본문 길이를 방어적으로 검증했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5416)
- Linux의 GTK 메뉴 처리에서 발생하는 메모리 안전성 문제를 수정했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5363)
- NVIDIA GPU를 감지해 Linux에서 DMA-BUF 렌더러를 비활성화했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5295)
- macOS의 `SetPosition` 화면 간 Y 좌표 변환을 수정했습니다. 기본 화면 높이를 전역 기준으로 사용하므로 기본 디스플레이와 수직으로 어긋난 모니터에서도 창이 올바른 위치에 배치됩니다. [#5117](https://github.com/wailsapp/wails/issues/5117)
- Git PR 템플릿이 올바른 피드백 URL을 가리키도록 수정했습니다. @wayneforrest의 [PR](https://github.com/wailsapp/wails/pull/5109)
- 인수 하나가 아닌 네 개를 전달하던 잘못된 `DestroyMenu` 시스템 호출로 인해 발생하는 Windows 시스템 트레이의 `SetMenu` 크래시 문제를 수정했습니다. 이 문제로 모든 호출이 FALSE를 반환하고 아무것도 해제하지 못했습니다. 또한 메뉴를 다시 빌드할 때 HMENU 및 HBITMAP 핸들(`MenuItem.SetBitmap`을 통해 런타임에 할당된 핸들 포함)을 해제하고, `Win32Menu.Update`에서 오래된 체크박스/라디오 맵을 초기화하며, 할당을 두 배로 늘리던 `systemtray.updateMenu`의 불필요한 `Update()` 호출을 제거했습니다. 이제 장시간 실행되는 시스템 트레이 앱이 메뉴를 다시 빌드할 때마다 GDI/USER 객체를 누수하지 않습니다.

## v3.0.0-alpha.90 - 2026-05-11

## 추가됨

- macOS의 WKWebView User-Agent에 설정 가능한 애플리케이션 이름을 추가했습니다. @vinhvoit225의 [PR](https://github.com/wailsapp/wails/pull/5261)
- gin-service 예제에 간접 의존성 github.com/coder/websocket을 추가했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5400)
- 빌드 애셋 테스트에 심층 동등성 비교 지원을 추가했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5402)

## 변경됨

- 빌드 출력을 assets 디렉터리로 통합했습니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5401)
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5399)에서 스폰서 SVG를 업데이트했습니다.

## 수정됨

- macOS 단일 인스턴스 메시지에 알림 객체를 사용합니다. @overlordtm의 [PR](https://github.com/wailsapp/wails/pull/5289)
- 부하가 클 때 프로미스가 손실되지 않도록 Windows 콜백을 일괄 처리합니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5383)

## v3.0.0-alpha.89 - 2026-05-10

## 추가됨

- Go 테스트 결과를 집계하는 go<em>test</em>results 작업을 추가했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5316)

## 변경됨

- 큰 RPC 페이로드를 조건에 따라 청크 단위 POST 요청으로 분할합니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5369)
- 모든 프런트엔드 템플릿에서 Vite를 5.x.x에서 8.0.0로 업그레이드했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5386)
- Vite 개발 서버 포트 구성을 환경 변수로 이전했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5365)
- 모든 템플릿에서 Vite 개발 서버가 127.0.0.1에 바인딩되도록 구성했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5361)
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5384)에서 스폰서 SVG를 업데이트했습니다.

## 수정됨

- build-assets 업데이트 중 Info.plist 템플릿 스텁을 정리합니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5312)
- 메뉴 항목 변경 메서드(`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`)를 메인 스레드에서 동기적으로 적용해 macOS 메뉴의 오래된 상태 문제를 수정했습니다. 메뉴를 빠르게 다시 열었을 때 이전 상태가 렌더링되던 `dispatch_async` 경쟁 상태를 제거했습니다(#5002).
- 불필요한 재빌드를 방지하도록 개발 모드에서 `*_test.go` 파일을 무시합니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5203)
- 앱이 실행 중이 아닐 때 Menu.Update()에서 세그멘테이션 오류가 발생하지 않도록 했습니다. @wucm667의 [PR](https://github.com/wailsapp/wails/pull/5291)
- Windows에서 메뉴 모음 다시 그리기를 제한하는 데 lastSizeWParam을 사용합니다. @taliesin-ai의 [PR](https://github.com/wailsapp/wails/pull/5382)

## v3.0.0-alpha.88 - 2026-05-09

## 변경됨

- HiddenOnTaskbar가 WS<em>EX</em>TOOLWINDOW를 사용하도록 변경했습니다. @leaanthony의 [PR](https://github.com/wailsapp/wails/pull/5371)
- go.mod에서 의존성 순서를 변경하고 webview2 replace 지시문을 제거했습니다. @atterpac의 [PR](https://github.com/wailsapp/wails/pull/5370)
- `@github-actions[bot]`이 [PR](https://github.com/wailsapp/wails/pull/5358)에서 스폰서 SVG를 업데이트했습니다.

## 수정됨

- 제네릭 간접 참조 별칭을 제거하고 맵 키 형식을 통합했습니다. @fbbdev의 [PR](https://github.com/wailsapp/wails/pull/5331)

## 제거됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5377)에서 문서화, Go 테스트 및 테스트 건너뛰기를 제거하고 PR-master 워크플로 삭제

## v3.0.0-alpha.87 - 2026-05-07

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5352)에서 Wails v3 한국어 문서 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5354)에서 설치 및 빠른 시작 프랑스어 문서 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5355)에서 빠른 시작, 개념 및 커뮤니티 포르투갈어 문서 추가

## v3.0.0-alpha.86 - 2026-05-06

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5328)에서 프랑스어 문서 현지화 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5343)에서 문서 사이트에 독일어 로케일 추가

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5347)에서 번역된 8개 로케일을 문서 구성에 모두 등록
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5317)에서 WebView2와 관련된 여러 Windows 파일 업데이트

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5340)에서 Linux의 대화 상자 디스패치를 GTK3와 GTK4로 분리
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5339)에서 대화 상자 콜백이 GTK 스레드에서 실행되도록 보장하여 세그멘테이션 오류 수정

## v3.0.0-alpha.85 - 2026-05-05

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5179)에서 저장소에 PR 템플릿 URL 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5330)에서 Wails v3 독일어 문서 추가

## v3.0.0-alpha.84 - 2026-05-03

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5307)에서 macOS에서 Escape 키로 전체 화면을 종료하지 못하게 하는 옵션 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5310)에서 macOS에서 Escape 키로 전체 화면을 종료하지 못하게 하는 옵션 추가
- @yuseferi가 [PR](https://github.com/wailsapp/wails/pull/5288)에서 Pausa 커뮤니티 쇼케이스 문서 추가

## 변경됨

- `@github-actions[bot]`가 [PR](https://github.com/wailsapp/wails/pull/5308)에서 스폰서 SVG 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5309)에서 지원되지 않는 플랫폼을 처리하도록 아이콘 생성 명령 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5224)에서 불리언 전체 화면 API를 3상태 ButtonState로 교체하고 플랫폼 바인딩 구현

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5315)에서 컨트롤러 상태가 nil일 때 WebView2 포커스 작업을 방지하도록 보호 로직 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5313)에서 PR의 베이스 브랜치를 올바르게 참조하도록 GitHub Actions 워크플로 업데이트
- @leaanthony가 불필요한 재빌드를 방지하기 위해 [PR](https://github.com/wailsapp/wails/pull/5203)에서 개발 모드가 `*_test.go` 파일을 무시하도록 수정
- @wucm667이 [PR](https://github.com/wailsapp/wails/pull/5291)에서 앱이 실행 중이 아닐 때 발생하는 Menu.Update() 세그멘테이션 오류 방지

## v3.0.0-alpha.83 - 2026-05-02

## 추가됨

- @symball이 [PR](https://github.com/wailsapp/wails/pull/5094)에서 시스템/사용자 설치를 위한 InstallScope 플래그 및 빌드 옵션 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5294)에서 Window 인터페이스를 충족하도록 BrowserWindow에 아무 작업도 하지 않는 SetScreen 메서드 추가

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5295)에서 NVIDIA GPU를 감지하고 Linux에서 DMA-BUF 렌더러 비활성화
- @wayneforrest가 [PR](https://github.com/wailsapp/wails/pull/5109)에서 올바른 피드백 URL을 가리키도록 git PR 템플릿 수정
- 하나가 아닌 네 개의 인수를 전달해 모든 호출이 FALSE를 반환하고 아무것도 해제하지 못하게 하던 손상된 `DestroyMenu` 시스템 호출 때문에 발생한 Windows 시스템 트레이 `SetMenu` 크래시 계열을 수정했습니다. 메뉴를 다시 빌드할 때 HMENU 및 HBITMAP 핸들(`MenuItem.SetBitmap`를 통해 런타임에 할당된 핸들 포함)도 해제하고, `Win32Menu.Update`에서 오래된 체크박스/라디오 맵을 초기화하며, 할당을 두 배로 늘리던 `systemtray.updateMenu`의 불필요한 `Update()` 호출을 제거했습니다. 이제 장시간 실행되는 시스템 트레이 앱은 메뉴를 다시 빌드할 때마다 GDI/USER 객체를 누수하지 않습니다.

## v3.0.0-alpha.82 - 2026-05-01

## 수정됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5232)에서 데스크톱 이름을 올바르게 처리하도록 데스크톱 파일 생성 수정

## v3.0.0-alpha.81 - 2026-04-30

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5286)에서 나이틀리 릴리스 일정을 UTC 15:00로 조정

## 수정됨

- Retina Mac에서 Screen Bounds, WorkArea 및 Size 값이 절반으로 줄어드는 문제 수정 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5285)에서 문서 종속성 및 콘텐츠 컬렉션 로더 업데이트

## v3.0.0-alpha.79 - 2026-04-29

## 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5270)에서 trigger-release 작업에 actions: write 권한 부여

## 변경됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5283)에서 릴리스 작업의 기본 브랜치를 master로 설정하고 변경 로그 문구 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5282)에서 최신 버전을 사용하도록 auto-changelog 워크플로 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5280)에서 경로 필터를 추가하고 사용되지 않는 워크플로를 제거하여 워크플로 효율 개선
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5274)에서 예제 링크가 master 브랜치를 참조하도록 문서 업데이트
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5272)에서 v3 관련 문서 및 예제 업데이트

## 수정됨

- @AkagiYui가 [PR](https://github.com/wailsapp/wails/pull/5265)에서 개발용 재시도 로직과 IPv4 강제 사용을 추가하여 리버스 프록시 개선
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5281)에서 미릴리스 변경 로그 트리거 워크플로 재작성

## 제거됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5267)에서 다양한 테스트 용도의 셸 테스트 스크립트 제거
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5266)에서 v3-alpha 문서 배포 워크플로 및 CNAME 레코드 삭제

### 추가됨

- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5196)에서 사이드바 탐색에 프런트엔드 라우팅 항목 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/5185)에서 프레임워크별 권장 사항을 포함한 프런트엔드 라우팅 가이드 추가
- 모달 시트 지원 추가(macOS)
- Apple 기기 지원 개선을 위해 @leaanthony가 ghw 버전 업데이트(#4977)
- Dock 서비스에 `GetBadge` 메서드 추가
- 사용자 지정 Go 빌드 태그(예: `wails3 build -tags gtk4`)를 전달할 수 있도록 `wails3 build` 명령에 `-tags` 플래그 추가(#4957)
- 전용 Enums 페이지와 사이드바 탐색 항목을 포함해 바인딩 생성기의 자동 열거형 생성 문서 추가(#4972)
- 사용자 지정 Go 빌드 태그(예: `wails3 build -tags gtk4`)를 전달할 수 있도록 `wails3 build` 명령에 `-tags` 플래그 추가(#4957)
- `v3/examples/web-apis/`에 Web API 예제를 추가하여 Storage(localStorage, sessionStorage, IndexedDB, Cache API), Network(Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media(Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Device(Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performance(Performance API, Mutation Observer, Intersection/Resize Observer), UI(Web Components, Pointer Events, Selection, Dialog, Drag and Drop) 등을 포함한 41 브라우저 API 시연
- 플랫폼 전반에서 200개 이상의 브라우저 API를 테스트하는 WebView API 호환성 검사기 예제(`v3/examples/webview-api-check/`) 추가
- 병렬 검색, 캐싱 및 Flatpak/Snap/Nix 지원을 통해 Linux에서 네이티브 라이브러리 경로를 찾는 `internal/libpath` 패키지 추가
- **작업 진행 중:** Linux용 실험적 WebKitGTK 6.0 / GTK4 지원 추가. `-tags gtk4`을 통해 사용할 수 있음(GTK3/WebKit2GTK 4.1은 계속 기본값)
- 참고: 타일링 창 관리자(예: Hyprland, Sway)에서는 WM이 창의 위치와 크기를 제어하므로 최소화/최대화 작업이 예상대로 작동하지 않을 수 있습니다.
- @AbdelhadiSeddar가 **JavaScript에서 이벤트 수신** 문서에 **일회성 핸들러** 사용 방법 추가
- Windows/Linux의 창이 `app.Menu.Set()`을 통해 설정된 애플리케이션 메뉴를 상속할 수 있도록 @leaanthony가 `WebviewWindowOptions`에 `UseApplicationMenu` 옵션 추가
- @wimaha가 Liquid Glass 아이콘과 애셋 카탈로그를 생성할 때 `.icon` 파일(Apple Icon Composer 형식)을 사용할 수 있도록 지원 추가(macOS)(#4934)
- 헤드리스/웹 배포를 위한 실험적 서버 모드(`-tags server`) 추가. 네이티브 GUI 종속성 없이 Wails 앱을 HTTP 서버로 실행할 수 있습니다. `wails3 task build:server`로 빌드하세요. 자세한 내용은 `examples/server`을 참조하세요.
- 병렬 검색, 캐싱 및 Flatpak/Snap/Nix 지원을 통해 Linux에서 네이티브 라이브러리 경로를 찾는 `internal/libpath` 패키지 추가
- macOS Spaces 및 전체 화면에서 창 동작을 제어할 수 있도록 @leaanthony가 `MacWindow`에 `CollectionBehavior` 옵션 추가(#4756)
- @leaanthony가 pkg/application 단위 테스트 추가
- @leaanthony가 MSIX 패키징에 사용자 지정 프로토콜 지원 추가
- Linux 데스크톱 환경 감지 기능 추가 [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 프런트엔드에서 인쇄 대화 상자를 열 수 있도록 @leaanthony가 JavaScript 런타임에 `Window.Print()` 메서드 추가(#4290)
- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE` 추가
- @leaanthony가 Linux용 WebKit2 로드 변경 이벤트 `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` 추가(#3896)
- @leaanthony가 Linux의 `wails3 doctor` 출력에 `XDG_SESSION_TYPE` 추가
- 패키징할 때뿐 아니라 Linux 빌드 중에도 `.desktop` 파일 생성(#4575)
- 배포판별 패키지 이름과 nfpm 패키징 예제를 포함한 Linux 런타임 종속성 문서를 @leaanthony가 추가(#4339)
- @leaanthony가 Linux의 `wails3 doctor` 출력에 NVIDIA 드라이버 버전 정보 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4710)에서 원시 메시지 핸들러에 출처(origin) 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4712)에서 macOS용 유니버설 링크 지원 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4702)에서 바인딩 전송 계층 리팩터링
- 예제 앱을 Appium 테스트 클라이언트로 쉽게 테스트할 수 있도록 @chinenual이 [PR](https://github.com/wailsapp/wails/pull/4760)에서 helloworld 템플릿에 aria-label 식별자 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4710)에서 원시 메시지 핸들러에 출처(origin) 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4712)에서 macOS용 유니버설 링크 지원 추가
- @APshenkin이 [PR](https://github.com/wailsapp/wails/pull/4702)에서 바인딩 전송 계층 리팩터링
- @fbbdev와 @ianvs가 [#4633](https://github.com/wailsapp/wails/pull/4633)에서 형식화된 이벤트 추가
- 툴팁이 실시간으로 업데이트되는 헤드리스 트레이를 보여 주는 `systray-clock` 예제 추가(#4653).
- @Tolfx가 #4510에서 Windows용 NSIS Protocol 템플릿 추가
- @Tolfx가 #4510에서 build-assets 테스트 추가
- macOS: @nidib가 [#4588](https://github.com/wailsapp/wails/pull/4588)에서 메뉴 막대에 네이티브 창 컨트롤 표시
- @popaprozac이 [PR](https://github.com/wailsapp/wails/pull/4451)에서 Dock의 앱 아이콘을 숨기거나 표시하는 macOS Dock 서비스 추가
- @popaprozac이 [PR](https://github.com/wailsapp/wails/pull/4451)에서 Dock의 앱 아이콘을 숨기거나 표시하는 macOS Dock 서비스 추가
- @leaanthony가 [#4534](https://github.com/wailsapp/wails/pull/4534)에서 NSGlassEffectView(macOS 15.0 이상)와 NSVisualEffectView 대체 구현을 사용하는 macOS용 네이티브 Liquid Glass 효과 지원 및 광범위한 재질 사용자 지정 옵션 추가
- @leaanthony가 [#4500](https://github.dev/wailsapp/wails/pull/4500)에서 브라우저 URL 정제 기능 추가. @APShenkin의 [#4484](https://github.com/wailsapp/wails/pull/4484)을 기반으로 함.
- [@Taiterbase](https://github.com/Taiterbase)의 원래 작업을 기반으로 [@leaanthony](https://github.com/leaanthony)가 이 [PR](https://github.com/wailsapp/wails/pull/4241)에서 Windows/Mac 콘텐츠 보호 기능 추가
- @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/4488)에서 `wails3 build` 및 `wails3 package` 별칭을 통해 CLI 변수를 Task 명령에 전달하는 기능 추가(#4422)
- [@atterpac](https://github.com/atterpac)이 [#4318](https://github.com/wailsapp/wails/pull/4318)에서 이벤트 소싱으로 드롭된 요소 데이터를 제공하는 드롭존 지원 추가
- WebView2 브라우저에 추가 명령줄 인수를 전달할 수 있도록 `WindowsWindow` 옵션에 `AdditionalLaunchArgs` 추가. [PR](https://github.com/wailsapp/wails/pull/4467)에서 구현
- [@triadmoko](https://github.com/triadmoko)가 [PR](https://github.com/wailsapp/wails/pull/4286)에서 wails init 후 go mod tidy를 자동으로 실행하도록 추가
- @leaanthony가 [PR](https://github.dev/wailsapp/wails/pull/4463)에서 Windows Snapassist 기능 추가
- WebView2 브라우저에 추가 명령줄 인수를 전달할 수 있도록 `WindowsWindow` 옵션에 `AdditionalLaunchArgs` 추가. [PR](https://github.com/wailsapp/wails/pull/4467)에서 구현
- [@triadmoko](https://github.com/triadmoko)가 [PR](https://github.com/wailsapp/wails/pull/4286)에서 wails init 후 go mod tidy를 자동으로 실행하도록 추가
- @leaanthony가 [PR](https://github.dev/wailsapp/wails/pull/4463)에서 Windows Snapassist 기능 추가
- [@almas-x](https://github.com/almas-x)가 [PR](https://github.com/wailsapp/wails/pull/4427)에서 Windows용 `getAccentColor` 구현 추가
- [@almas-x](https://github.com/almas-x)가 [PR](https://github.com/wailsapp/wails/pull/4427)에서 Windows용 `getAccentColor` 구현 추가
- Windows 다크 테마 메뉴 및 메뉴 모음. @leaanthony가 [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)에서 기여
- JS/TS 바인딩을 더 명확하게 하기 위해 기본 제공 서비스의 이름을 변경. @popaprozac이 [PR](https://github.com/wailsapp/wails/pull/4405)에서 기여
- 사용자 시스템의 강조 색상을 가져오는 `app.Env.GetAccentColor`. MacOS에서 작동. [@etesam913](https://github.com/etesam913)이 기여
- `window.ToggleFrameless()` API 추가. [@atterpac](https://github.com/atterpac)이 [#4137](https://github.com/wailsapp/wails/pull/4137)에서 기여
- Linux 배포판별 빌드 종속성 추가. @leaanthony가 [PR](https://github.com/wailsapp/wails/pull/4345)에서 기여
- 바인딩 가이드 추가. @atterpac이 [PR](https://github.com/wailsapp/wails/pull/4404)에서 기여
- **테스트 인프라 정비**: Docker 테스트 파일을 전용 `test/docker/` 디렉터리로 이동하고 이미지를 최적화하여 빌드 안정성을 개선. [@leaanthony](https://github.com/leaanthony)가 [#4359](https://github.com/wailsapp/wails/pull/4359)에서 기여
- **리소스 관리 패턴 개선**: 예제에 적절한 이벤트 핸들러 정리와 컨텍스트 인식 goroutine 관리를 추가. [@leaanthony](https://github.com/leaanthony)가 [#4359](https://github.com/wailsapp/wails/pull/4359)에서 기여
- aarch64 AppImage 빌드 지원. [@AkshayKalose](https://github.com/AkshayKalose)가 [#3981](https://github.com/wailsapp/wails/pull/3981)에서 기여
- `wails doctor`에 진단 섹션 추가. [@leaanthony](https://github.com/leaanthony)가 기여
- 서비스 메서드를 호출할 때 컨텍스트에 창 추가. [@leaanthony](https://github.com/leaanthony)가 기여
- 어느 창에서 서비스를 호출하는지 확인하는 방법을 보여 주는 `window-call` 예제 추가. [@leaanthony](https://github.com/leaanthony)가 기여
- 새로운 메뉴 가이드. [@leaanthony](https://github.com/leaanthony)가 기여
- panic 처리 개선. [@leaanthony](https://github.com/leaanthony)가 기여
- 새로운 메뉴 가이드. [@leaanthony](https://github.com/leaanthony)가 기여
- Service API에 문서 주석 추가. [@fbbdev](https://github.com/fbbdev)가 [#4024](https://github.com/wailsapp/wails/pull/4024)에서 기여
- 추가 구성으로 서비스를 초기화하는 `application.NewServiceWithOptions` 함수 추가. [@leaanthony](https://github.com/leaanthony)가 [#4024](https://github.com/wailsapp/wails/pull/4024)에서 기여
- 메뉴 제어 개선. [@FalcoG](https://github.com/FalcoG)와 [@leaanthony](https://github.com/leaanthony)가 [#4031](https://github.com/wailsapp/wails/pull/4031)에서 기여
- 문서 추가. [@leaanthony](https://github.com/leaanthony)가 기여
- 표준 이벤트 리스너에서 이벤트 취소 지원. [@leaanthony](https://github.com/leaanthony)가 기여
- Systray `Hide`, `Show` 및 `Destroy` 지원. [@leaanthony](https://github.com/leaanthony)가 기여
- Systray `SetTooltip` 지원. [@leaanthony](https://github.com/leaanthony)가 기여. 최초 아이디어는 [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)이 제안
- 지원되지 않는 형식에 대한 바인딩 생성기 경고에 패키지 경로 표시. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- 제네릭 별칭에 대한 바인딩 생성기 지원 추가. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- `omitzero` JSON 플래그에 대한 바인딩 생성기 지원 추가. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- 선택한 서비스 메서드의 바인딩 생성을 방지하는 `//wails:ignore` 지시문 추가. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- Go에서는 내보내지만 JS/TS에서는 내보내지 않는 형식을 허용하도록 서비스와 모델에 `//wails:internal` 지시문 추가. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- 느슨하게 형식이 지정된 열거형을 사용할 수 있도록 별칭 형식 상수에 대한 바인딩 생성기 지원 추가. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 기여
- Go 1.24 기능에 대한 바인딩 생성기 테스트 추가. [@fbbdev](https://github.com/fbbdev)가 [#4068](https://github.com/wailsapp/wails/pull/4068)에서 기여
- OS 버전 감지를 개선하기 위해 `OSInfo.Branding`에 macOS 15 "Sequoia" 지원 추가. [#4065](https://github.com/wailsapp/wails/pull/4065)에서 기여
- 종료 절차가 완료된 후 사용자 지정 코드를 실행하는 `PostShutdown` 훅 추가. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- 사용자 지정 오류 핸들러에서 치명적 오류를 감지할 수 있도록 `FatalError` 구조체 추가. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- 서비스 시작 및 종료 순서를 표준화하고 문서화. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- 애플리케이션 시작/종료 순서와 서비스 시작/종료 테스트를 위한 테스트 하네스 추가. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- 애플리케이션 생성 후 서비스를 등록하는 `RegisterService` 메서드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- 바인딩 호출의 사용자 지정 오류 처리를 위해 애플리케이션 및 서비스 옵션에 `MarshalError` 필드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 기여
- promise 체인을 통해 취소 요청을 전파하는 취소 가능한 promise 래퍼 추가. [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 기여
- 바인딩 호출 취소를 `AbortSignal`에 연결하는 기능 추가. [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 기여
- WML에서 일반적인 `wml-*` 특성과 함께 `data-wml-*` 특성 지원. [@leaanthony](https://github.com/leaanthony)가 기여
- 지연 구성 및 동적 재구성을 위해 모든 서비스에 `Configure` 메서드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- `fileserver` 서비스가 구성되지 않은 경우 503 Service Unavailable 응답을 전송. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- `kvstore` 서비스가 구성되지 않은 경우 기본적으로 인메모리 키-값 저장소를 제공. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- 구성 변경 후 파일에서 데이터를 다시 불러오도록 `kvstore` 서비스에 `Load` 메서드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- 모든 키를 삭제하도록 `kvstore` 서비스에 `Clear` 메서드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- JS 측 로그 수준 상수를 제공하도록 `log` 서비스에 `Level` 형식 추가. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- 로그 수준을 동적으로 지정하도록 `log` 서비스에 `Log` 메서드 추가. [@fbbdev](https://github.com/fbbdev)가 [#4067](https://github.com/wailsapp/wails/pull/4067)에서 기여
- `sqlite` 서비스는 구성되지 않은 경우 기본적으로 인메모리 DB를 제공합니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- DB를 수동으로 닫을 수 있도록 `sqlite` 서비스에 `Close` 메서드를 추가했습니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- `sqlite` 서비스의 쿼리 메서드에 취소 지원을 추가했습니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- `sqlite` 서비스에 JS 바인딩을 포함한 준비된 문 지원을 추가했습니다. [@fbbdev](https://github.com/fbbdev), [#4067](https://github.com/wailsapp/wails/pull/4067)
- [@AnalogJ](https://github.com/AnalogJ)의 이 [PR](https://github.com/wailsapp/wails/pull/3537) 원본 작업을 바탕으로 Gin 지원을 추가했습니다. [Lea Anthony](https://github.com/leaanthony), [PR](https://github.com/wailsapp/wails/pull/3537)
- 자동 저장과 비밀번호 자동 저장이 항상 활성화되는 문제를 수정했습니다. [@oSethoum](https://github.com/osethoum), [#4134](https://github.com/wailsapp/wails/pull/4134)
- 창에 메뉴를 설정할 수 있도록 창에 `SetMenu()`을 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- 알림 지원을 추가했습니다. [@popaprozac](https://github.com/popaprozac), [#4098](https://github.com/wailsapp/wails/pull/4098)
-  mac용 파일 연결 지원을 추가했습니다. [@wimaha](https://github.com/wimaha), [#4177](https://github.com/wailsapp/wails/pull/4177)
- 시맨틱 버전 증가를 위한 `wails3 tool version`을 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- macOS 및 Windows용 배지 지원을 추가했습니다. [@popaprozac](https://github.com/popaprozac), [#](https://github.com/wailsapp/wails/pull/4234)
- 등록된/엄격하게 형식이 지정된 이벤트 지원을 추가했습니다. [@fbbdev](https://github.com/fbbdev) 및 [@IanVS](https://github.com/IanVS), [#4161](https://github.com/wailsapp/wails/pull/4161)
- 사용자 정의 이벤트용 훅을 등록하는 기능을 추가했습니다. [@fbbdev](https://github.com/fbbdev) 및 [@IanVS](https://github.com/IanVS), [#4161](https://github.com/wailsapp/wails/pull/4161)
- `path` 경로에서 시스템 파일 관리자를 여는 `app.OpenFileManager(path string, selectFile bool)`을 추가했으며, `selectFile`을 통해 선택적으로 강조 표시할 수 있습니다. [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- `wails3 init` 명령에 새로운 `-git` 플래그를 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- 새로운 `wails3 generate webview2bootstrapper` 명령을 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- 런타임을 수동으로 초기화할 수 있도록 런타임에 `init()` 메서드를 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- Window의 WindowOptions에 `WindowDidMoveDebounceMS` 옵션을 추가했습니다. [@leaanthony](https://github.com/leaanthony)
- 단일 인스턴스 기능을 추가했습니다. [@leaanthony](https://github.com/leaanthony). @APshenkin의 [v2 PR](https://github.com/wailsapp/wails/pull/2951)을 기반으로 합니다.
- `wails3 generate template` 명령. [@leaanthony](https://github.com/leaanthony)
- `wails3 releasenotes` 명령. [@leaanthony](https://github.com/leaanthony)
- `wails3 update cli` 명령. [@leaanthony](https://github.com/leaanthony)
- `wails3 generate bindings` 명령용 `-clean` 옵션. [@leaanthony](https://github.com/leaanthony)
- aarch64(arm64) AppImage Linux 빌드를 지원합니다. [@AkshayKalose](https://github.com/AkshayKalose), [#3981](https://github.com/wailsapp/wails/pull/3981)
- 후원자용 하이퍼링크를 추가했습니다. @ansxuman, [#3958](https://github.com/wailsapp/wails/pull/3958)
- deb, rpm 및 Arch Linux 패키저 빌드를 포함한 Linux 패키징 지원:
- Darwin 유니버설 빌드 및 패키지 지원 추가:
- 웹사이트에 이벤트 문서 추가:
- 비 SSR 개발용으로 설정된 sveltekit 및 sveltekit-ts 템플릿
- 새로운 `wails3 update build-assets` 명령을 사용하여 빌드 자산 업데이트:
- HTML Drag and Drop API를 테스트하는 예제:
- 파일 연결 지원. [leaanthony](https://github.com/leaanthony), 다음에서:
- 새로운 `wails3 generate runtime` 명령:
- 창을 가운데에 배치할지 여부를 지정하는 새로운 `InitialPosition` 옵션:
- `application` 패키지에 `Path` 및 `Paths` 메서드 추가:
- `GeneralAutofillEnabled` 및 `PasswordAutosaveEnabled` Windows 옵션 추가
- 서비스 메서드를 호출한 창을 가져오는 기능 추가:
- Webview2용 `EnabledFeatures` 및 `DisabledFeatures` 옵션 추가:
- ⊞ 향상된 High DPI 모니터 지원을 위한 새로운 DIP 시스템:
- ⊞ 창 클래스 이름 옵션. [windom](https://github.com/windom/), 다음에서:
- 플러그인 기능을 제공하도록 서비스가 확장되었습니다. 작성자:
- 🐧 WindowDidMove / WindowDidResize 이벤트, 다음에서:
- ⊞ WindowDidResize 이벤트, 다음에서:
-  Dock을 처리할 수 있도록 Event ApplicationShouldHandleReopen 추가:
-  구현에 getPrimaryScreen/getScreens 추가. @tmclane, 다음에서:
-  macOS 전체 화면 모드에서 도구 모음을 표시하는 옵션 추가:
- 🐧 Linux 키 입력을 액셀러레이터로 변환하는 onKeyPress 로직 추가:
- 🐧 작업 `run:linux` 추가:
- `SetIcon` 메서드 내보내기. [@almas-x](https://github.com/almas-x), 다음에서:
- `OnShutdown` 개선. [@almas-x](https://github.com/almas-x), 다음에서:
- `Window` 인터페이스에 `ToggleMaximise` 메서드 복원:
- `Environment()`에 추가 정보를 더했습니다. @leaanthony, 다음에서:
- `Window` 인터페이스에 `WebviewWindow.IsFocused` 메서드 공개:
- WML 시스템에서 공백으로 구분된 여러 트리거 이벤트 지원:
- 번들 JS 런타임 스크립트에서 ESM 내보내기 추가:
- 다음 대신 번들 JS 런타임 스크립트를 사용하기 위한 바인딩 생성기 플래그 추가:
- Linux에 `setIcon` 구현. [@abichinger](https://github.com/abichinger)
- dev 명령에 `-port` 플래그를 추가하고 환경 변수를 지원:
- 바인딩된 메서드 호출 테스트 추가:
- ⊞ 이미 생성된 창에 `SetIgnoreMouseEvents` 추가:
-  창의 쌓임 수준(순서)을 설정하는 기능 추가:

### 수정됨

- `Physical*` 필드에서 NSScreen 포인트 값을 장치 픽셀로 변환하여 Retina Mac에서 `Screen.Bounds`, `WorkArea` 및 `Size` 값이 절반으로 줄어드는 문제를 수정하고, 최상위 `Screen.X`/`Y`을 채워 다중 모니터 접촉 감지와 작업 영역 배치가 올바르게 동작하도록 했습니다. @wayneforrest, [PR](https://github.com/wailsapp/wails/pull/5168)
- 디스플레이 구성 변경 시(예: 절전 모드 진입/해제 중 외부 모니터 핫플러그) WebKit DisplayLink 교착 상태를 일으키는 ScreenManager의 데이터 경합을 수정했습니다.
- Assets.car가 있으면 CFBundleIconName을 appicon으로 직접 설정하도록 변경 [PR](https://github.com/wailsapp/wails/pull/5154), @symball
- Fedora, openSUSE, Arch 및 NixOS에서 `wails3 doctor`이 잘못된 WebKitGTK 패키지를 보고하던 문제 수정 — v3이 컴파일 시 4.1 API를 요구하므로 4.0 대체 항목 제거 (#5071)
- openSUSE의 webkit2gtk doctor 패키지 이름 수정(`webkit2gtk4_1-devel` → 올바른 openSUSE 패키지 이름인 `webkit2gtk3-devel`) (#5071)
- 데스크톱 개발 모드에서 `/wails/custom.js`이 없을 때 발생하는 `Unexpected token '<'` 오류 수정. HTML SPA 대체 콘텐츠가 JavaScript로 삽입되지 않도록 `/wails/custom.js`에 명시적인 404 처리기를 추가하고, `loadOptionalScript`에서 대소문자를 구분하지 않는 `Content-Type` 검증 추가. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- macOS의 시스템 트레이 메뉴 강조 상태 수정 - 이제 메뉴가 열려 있으면 아이콘이 선택된 상태로 표시됨 (#4910)
- macOS에서 시스템 트레이에 연결된 창이 다른 창 뒤에 나타나던 문제 수정 - 이제 적절한 팝업 창 레벨 사용 (#4910)
- 문서 전반의 잘못된 `@wailsio/runtime` import 예제 수정 (#4989)
- darwin에서 프레임 없는 창을 최소화할 수 없던 문제 수정 (#4294)
- `node_modules/`을 go-task의 최신 상태 검사에서 제외하여 `wails3 build` 및 `wails3 dev` 실행 중 20-30분 동안 멈추던 문제 수정. 이전에는 `sources: "**/*"` glob 때문에 go-task가 `node_modules/`의 모든 파일(MUI처럼 무거운 종속성이 포함된 50000-100000개 이상의 파일)을 열거하고 체크섬을 계산했으며, 특히 Windows/NTFS에서 느렸음 (#4939)
- C의 `Screen` typedef가 X11 Xlib.h와 충돌하여 발생하던 GTK4 빌드 실패 수정 (#4957)
- macOS의 Dock 배지 메서드 일관성 수정
- `InvisibleTitleBarHeight`이 프레임 없는 창이나 투명 제목 표시줄 창에만 적용되지 않고 모든 macOS 창에 적용되던 문제 수정 (#4960)
- `InvisibleTitleBarHeight`이 활성화된 상태에서 위쪽 모서리로 창 크기를 조절할 때 창이 흔들리거나 떨리던 문제를 창 가장자리 근처에서는 드래그를 시작하지 않도록 하여 수정 (#4960)
- JS/TS 바인딩에서 enum 키가 있는 매핑된 타입 생성 문제 수정 (#4437), @fbbdev
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간에 있던 문제 수정(물리적 픽셀과 CSS 픽셀의 차이)
- Linux에서 hover 효과를 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않던 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- `gtk_window_present()`을 사용하여 Linux/GTK4에서 창을 표시하거나 숨길 때 간혹 최소화된 상태로 복원되던 문제 수정 (#4957)
- `XTranslateCoordinates`/`XMoveWindow`을 통해 X11 조건부 지원을 추가하여 Linux/GTK4에서 창 위치를 가져오거나 설정할 때 항상 0,0이 반환되던 문제 수정 (#4957)
- 제거된 `gtk_window_set_geometry_hints`을 대신할 신호 기반 크기 제한을 추가하여 Linux/GTK4에서 최대 창 크기가 적용되지 않던 문제 수정 (#4957)
- 올바른 PhysicalBounds 계산을 구현하고 `gdk_monitor_get_scale`을 통해 소수 배율 조정을 지원하여 Linux/GTK4의 DPI 배율 조정 문제 수정(GTK 4.14+)
- Linux/GTK4에서 새 창을 만들 때 메뉴 항목이 중복되던 문제 수정
- JS/TS 바인딩에서 enum 키가 있는 매핑된 타입 생성 문제 수정 (#4437), @fbbdev
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간에 있던 문제 수정(물리적 픽셀과 CSS 픽셀의 차이)
- Linux에서 hover 효과를 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않던 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- 올바른 PhysicalBounds 계산을 구현하고 `gdk_monitor_get_scale`을 통해 소수 배율 조정을 지원하여 Linux/GTK4의 DPI 배율 조정 문제 수정(GTK 4.14+)
- Linux/GTK4에서 새 창을 만들 때 메뉴 항목이 중복되던 문제 수정
- JS/TS 바인딩에서 enum 키가 있는 매핑된 타입 생성 문제 수정 (#4437), @fbbdev
- App.Window.Current()에서 Main Thread를 통해 AppKit API에 접근하지 않아 macOS에서 발생하던 "유령 창" 문제 수정 (#4947), @wimaha
- WKUIDelegate runOpenPanelWithParameters를 구현하여 macOS에서 HTML `<input type="file">`이 작동하지 않던 문제 수정 (#4862)
- macOS/Linux에서 `@wailsio/runtime` npm 모듈을 사용할 때 네이티브 파일 드래그 앤 드롭이 작동하지 않던 문제 수정 (#4953), @leaanthony
- 패키지 간 타입 별칭의 바인딩 생성 문제 수정 (#4578), @fbbdev
- GTK 스레드 안전성 위반으로 인해 Linux에서 발생하던 OpenFileDialog 충돌 수정 (#3683), @ddmoney420
- 숨겨지거나 제거된 창에서 `Focus()`을 호출할 때 발생하던 SIGSEGV 충돌 수정 (#4890), @ddmoney420
- Linux에서 빈 아이콘이나 비트맵을 설정할 때 발생할 수 있던 panic 수정 (#4923), @ddmoney420
- macOS에서 서비스 바인딩을 통해 호출할 때 발생하던 ErrorDialog 충돌 수정 (#3631), @leaanthony
- `v3\examples\dialogs`에서 Windows OS에 메뉴가 표시되도록 변경, @ndianabasi
- 페이지를 다시 불러오는 동안 TypeError를 일으키던 경쟁 조건 수정 (#4872), @ddmoney420
- `Collector.IsVoidAlias()` 메서드에서 전역 상태를 제거하여 바인딩 생성기 테스트의 잘못된 출력 수정 (#4941), @fbbdev
- macOS에서 `<input type="file">` 파일 선택기가 작동하지 않던 문제 수정 (#4862), @leaanthony
- macOS에서 `Position()`과 `SetPosition()`이 서로 다른 좌표계를 사용하여 상태를 저장하고 복원할 때 창 위치가 달라지던 문제 수정 (#4816), @leaanthony
- 애플리케이션 매니페스트를 통해 DPI 인식이 이미 설정된 경우 발생하던 SetProcessDpiAwarenessContext "Access is denied" 오류 수정 (#4803)
- 키보드 단축키 문서 페이지를 업데이트하고 `KeyBinding.Add`의 콜백 매개변수 타입 수정, @ndianabasi
- 사용자 지정 바인딩 생성 관련 문서 수정: `-o String` 대신 반드시 `-d String`을 사용해야 함
- `menu.Update()`에서 메뉴의 하위 항목이 지워지지 않던 문제 수정
- 문서의 오래된 Manager API 참조 수정(31 파일을 `app.Window.New()`, `app.Event.Emit()` 등의 새 패턴을 사용하도록 업데이트), @leaanthony
- WebKit이 신호 처리기를 재정의하여 JS에 바인딩된 Go 메서드에서 panic이 발생할 때 Linux에서 충돌하던 문제 수정 (#3965), @leaanthony
- Linux에서 SaveFileDialog.SetFilename()이 아무 효과가 없던 문제 수정 (#4841), @samstanier
- 드래그 앤 드롭 예제에서 드롭 좌표가 undefined로 표시되던 문제 수정
- APP_NAME에 공백이 포함된 경우 macOS 앱 번들 생성에 실패하던 문제 수정(중괄호 확장 문제)
- Windows에서 서비스 메서드를 호출할 때 발생하던 인덱스 범위 초과 panic 수정(goccy/go-json 되돌림)
- Windows에서 디스플레이 배율이 100%가 아닐 때 파일 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않던 문제 수정
- Windows에서 파일 드롭 좌표가 잘못된 픽셀 공간에 있던 문제 수정(물리적 픽셀과 CSS 픽셀의 차이)
- Linux에서 호버 효과를 사용할 때 파일 드래그 앤 드롭이 안정적으로 작동하지 않는 문제 수정
- Linux에서 파일 드롭을 활성화하면 HTML5 내부 드래그 앤 드롭이 작동하지 않는 문제 수정
- @ndianabasi가 `APP_NAME` 같은 변수에 포함된 공백을 처리할 수 있도록 모든 운영 체제용 Taskfile.yml 파일의 모든 명령 업데이트
- @wux1an이 Linux에서 'build:universal:lipo:go' 태스크를 실행할 때 발생하는 명령 인수 오류 수정
- @wux1an이 Linux에서 'wails3 build GOOS=darwin GOARCH=arm64'를 실행할 때 발생하는 Docker 오류 "undefined symbol: **<em>ubsan</em>handle_xxxxxxx" 수정
- @leaanthony가 사용자 지정 프로토콜 문서를 통합하고 Universal Links 섹션 추가
- @leaanthony가 동시 TrackPopupMenuEx 호출을 방지하는 가드를 추가하여 아이콘을 반복해서 클릭할 때 Windows 시스템 트레이 메뉴가 충돌하는 문제 수정(#4151)
- @leaanthony가 app.Run()보다 먼저 systray.Run()을 호출할 때 앱이 충돌하지 않도록 수정
- @leaanthony가 ApplicationShouldTerminateAfterLastWindowClosed를 활성화한 상태에서 Hide()/Show()로 창 표시 여부를 전환할 때 macOS에서 발생하는 충돌 수정(#4389)
- @leaanthony가 메뉴를 반복해서 열 때 macOS와 Windows의 컨텍스트 메뉴에서 발생하는 메모리 누수 수정(#4012)
- @leaanthony가 macOS에서 컨텍스트 메뉴의 네이티브 리소스를 재사용하지 않아 메뉴를 표시할 때마다 새로 생성되는 문제 수정(#4012)
- @leaanthony가 `Hidden: true`로 앱을 시작했을 때 macOS Dock 아이콘을 클릭해도 숨겨진 창이 표시되지 않는 문제 수정(#4583)
- @leaanthony가 CGO 호출에서 잘못된 창 포인터 형식을 사용하여 macOS 인쇄 대화 상자가 열리지 않는 문제 수정(#4290)
- @leaanthony가 appmenu-gtk-module이 아직 실체화되지 않은 창에 접근하여 Wayland에서 창 메뉴가 충돌하는 문제 수정(#4769)
- @leaanthony가 앱 이름에 잘못된 문자(공백, 괄호 등)가 포함되었을 때 GTK 애플리케이션이 충돌하는 문제 수정
- @overlordtm이 Windows에서 드래그 앤 드롭을 초기화할 때 발생하는 "not enough memory" 오류 수정(#4701)
- @leaanthony가 잘못된 URI 이스케이프로 인해 Linux 파일 탐색기가 엉뚱한 디렉터리를 여는 문제 수정(#4397)
- @leaanthony가 `.relr.dyn` ELF 섹션을 자동으로 감지하고 스트리핑을 비활성화하여 최신 Linux 배포판(Arch, Fedora 39+, Ubuntu 24.04+)에서 AppImage 빌드가 실패하는 문제 수정(#4642)
- @leaanthony가 Fedora/DNF 기반 시스템에서 `wails doctor`이 webkit 패키지가 설치되었다고 잘못 보고하는 문제 수정(#4457)
- @mbaklor가 기본 `config.yml`이 프로덕션 빌드에서 `wails3 dev`을 실행하던 문제 수정
- @leaanthony가 존재하지 않는 패키지를 가져와 iOS 서비스 스텁에서 빌드가 실패하는 문제 수정
- @leaanthony가 debug/info 메서드의 구조화된 로깅에서 "no formatting directives" 오류가 발생하는 문제 수정
- @leaanthony가 모바일 플랫폼 병합 과정에서 실수로 포함된 임시 디버그 출력문 제거
- @leaanthony가 DMA-BUF 렌더러를 자동으로 비활성화하여 NVIDIA GPU를 사용하는 Wayland에서 WebKitGTK가 충돌하는 문제(Error 71 Protocol error) 수정
- Linux의 `application.WebviewWindowOptions.BackgroundColour`에서 알파 값이 무시되는 문제 해결([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- 사용자 지정 아이콘을 제공하지 않았을 때 Windows 시스템 트레이 아이콘이 기본적으로 애플리케이션 아이콘을 사용하지 않는 문제 수정(#4704)
- 사용자가 생성한 핸들만 폐기하도록 `HICON` 소유권을 추적하여 Explorer 재활용 시 발생하는 충돌 방지(#4653).
- 폐기 시 Windows 시스템 테마 리스너와 유지 중인 트레이 아이콘을 해제하여 goroutine과 디바이스 컨텍스트 누수 방지(#4653).
- 서로게이트 쌍과 멀티바이트 글리프의 손상을 방지하도록 트레이 도구 설명을 127 UTF-16 단위에서 잘라냄(#4653).
- Windows 패키지 태스크 실패 문제 수정(#4667)
- Linux taskfile의 Linux AppImage appicon 변수 수정 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- go-webview2 v1.0.22 시그니처 변경으로 발생한 Windows 빌드 오류 수정(#4513, #4645)
- Linux taskfile의 Linux AppImage appicon 변수 수정 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- @Tolfx가 linux desktop.tmpl의 프로토콜 범위에서 `<.Info.Protocol>`~`<.Protocol>`을 제거하여 수정(#4510)
- @Etesam913가 liquid glass 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542))
- [@JackDoan](https://github.com/JackDoan)이 Linux의 시스템 트레이 메뉴 업데이트 문제 수정([#4604](https://github.com/wailsapp/wails/issues/4604))
- @leaanthony가 숨겨진 창을 생성할 때 Windows에 흰색 창이 나타나는 문제 수정([#4612](https://github.com/wailsapp/wails/pull/4612))
- @rxliuli가 문서의 notifications 패키지 가져오기 경로 수정([#4617](https://github.com/wailsapp/wails/pull/4617))
- @leaanthony가 npm 패키지 @wailsio/runtime을 사용할 때 드래그 앤 드롭이 작동하지 않는 문제 수정(#4489, #4616)
- Windows: 시작 시 창이 깜박이고 숨겨진 창이 잘못 표시되는 문제를 @leaanthony가 수정([PR](https://github.com/wailsapp/wails/pull/4600)).
- [@samstanier](https://github.com/samstanier)가 Wayland 창 크기 최대화 문제 수정(https://github.com/wailsapp/wails/issues/4429)
- [@samstanier](https://github.com/samstanier)가 Wayland 창 크기 최대화 문제 수정(https://github.com/wailsapp/wails/issues/4429)
- @Etesam913가 liquid glass 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542))
- @jghiloni가 MacOS에서 AssetServer가 충돌할 수 있는 문제 수정([#4576](https://github.com/wailsapp/wails/pull/4576))
- @rev42가 NextJs로 빌드할 때 발생하는 컴파일 문제 수정([#4585](https://github.com/wailsapp/wails/pull/4585))
- @riadafridishibly가 나이틀리 릴리스 파이프라인 수정([#4597](https://github.com/wailsapp/wails/pull/4597))
- @Etesam913가 liquid glass 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542))
- @jghiloni가 MacOS에서 AssetServer가 충돌할 수 있는 문제 수정([#4576](https://github.com/wailsapp/wails/pull/4576))
- @rev42가 NextJs로 빌드할 때 발생하는 컴파일 문제 수정([#4585](https://github.com/wailsapp/wails/pull/4585))
- @riadafridishibly가 나이틀리 릴리스 파이프라인 수정([#4597](https://github.com/wailsapp/wails/pull/4597))
- @Etesam913가 liquid glass 데모의 재정의 오류 수정([#4542](https://github.com/wailsapp/wails/pull/4542))
- @PPTGamer가 Windows의 SetBackgroundColour 수정([PR](https://github.com/wailsapp/wails/pull/4492))
- @yulesxoxo가 Manager API 리팩터링의 변경 사항을 반영하도록 문서 업데이트([PR #4476](https://github.com/wailsapp/wails/pull/4476))
- Linux taskfile의 Linux .desktop 파일 appicon 변수 수정 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- @yulesxoxo가 Manager API 리팩터링의 변경 사항을 반영하도록 문서 업데이트([PR #4476](https://github.com/wailsapp/wails/pull/4476))
- @leaanthony가 [#4456](https://github.com/wailsapp/wails/issues/4456)에서 보고된 Windows nil 포인터 역참조 버그 수정([#4460](https://github.com/wailsapp/wails/pull/4460))
- 두 손가락 스와이프 탐색 제스처를 활성화하도록 macOS WKWebView에 `allowsBackForwardNavigationGestures` 지원 추가(#1857)
- 처음에 비활성화된 상태로 설정된 메뉴 항목에서 onClick이 작동하지 않던 문제를 @leaanthony가 [PR #4469](https://github.com/wailsapp/wails/pull/4469)에서 수정했습니다. 초기 조사에 도움을 준 @IanVS에게 감사드립니다.
- 빌드 실패 시 Vite 서버가 정리되지 않던 문제를 수정했습니다(#4403).
- Windows에서 `SaveFileDialog`을 닫거나 취소할 때 발생하던 패닉을 수정했습니다. @hkhere가 [PR](https://github.com/wailsapp/wails/pull/4284)에서 수정했습니다.
- Windows의 HTML 수준 드래그 앤 드롭을 [@mbaklor](https://github.com/mbaklor)가 [#4259](https://github.com/wailsapp/wails/pull/4259)에서 수정했습니다.
- 두 손가락 쓸어넘기기 탐색 제스처를 활성화하기 위해 macOS WKWebView에서 `allowsBackForwardNavigationGestures`을 지원하도록 추가했습니다(#1857).
- 처음에 비활성화된 상태로 설정된 메뉴 항목에서 onClick이 작동하지 않던 문제를 @leaanthony가 [PR #4469](https://github.com/wailsapp/wails/pull/4469)에서 수정했습니다. 초기 조사에 도움을 준 @IanVS에게 감사드립니다.
- 빌드 실패 시 Vite 서버가 정리되지 않던 문제를 수정했습니다(#4403).
- Windows의 알림 구문 분석 문제를 @popaprozac이 [PR](https://github.com/wailsapp/wails/pull/4450)에서 수정했습니다.
- Windows SDK 종속성을 확인하도록 doctor 명령을 [@kodumulo](https://github.com/kodumulo)가 [#4390](https://github.com/wailsapp/wails/issues/4390)에서 수정했습니다.
- Mac용 processURLRequest의 nil 포인터 역참조를 [@etesam913](https://github.com/etesam913)가 [#4366](https://github.com/wailsapp/wails/pull/4366)에서 수정했습니다.
- 필터가 적용된 대화 상자를 사용할 수 없게 하던 Linux 버그를 [@bh90210](https://github.com/bh90210)이 [#4287](https://github.com/wailsapp/wails/pull/4287)에서 수정했습니다.
- Windows와 Linux의 Edit 메뉴 문제를 [@leaanthony](https://github.com/leaanthony)가 [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)에서 수정했습니다.
- macOS .plist 파일의 최소 시스템 버전을 10.13.0에서 10.15.0로 변경했습니다. [@AkshayKalose](https://github.com/AkshayKalose)가 [#3981](https://github.com/wailsapp/wails/pull/3981)에서 작업했습니다.
- Window ID가 건너뛰어지는 문제를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- RegisterContextMenu 호출 시 발생하는 nil 메뉴 문제를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- 바인딩 생성기 출력의 종속성 순환을 [@fbbdev](https://github.com/fbbdev)가 [#4001](https://github.com/wailsapp/wails/pull/4001)에서 수정했습니다.
- 바인딩 생성기 출력의 정의 전 사용 오류를 [@fbbdev](https://github.com/fbbdev)가 [#4001](https://github.com/wailsapp/wails/pull/4001)에서 수정했습니다.
- 빌드 플래그가 바인딩 생성기에 전달되도록 [@fbbdev](https://github.com/fbbdev)가 [#4023](https://github.com/wailsapp/wails/pull/4023)에서 변경했습니다.
- Windows가 아닌 플랫폼에서도 작동하도록 Windows Taskfile의 경로를 슬래시 형식으로 [@leaanthony](https://github.com/leaanthony)가 변경했습니다.
- Mac + Mac JS 이벤트를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- macOS의 이벤트 교착 상태를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- Windows에서 HTML은 제공되었지만 JS는 제공되지 않은 경우 Window 초기화 중 발생하던 `Parameter incorrect` 오류를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- 애셋 서버에서 콘텐츠 유형 감지에 사용하는 응답 접두사의 크기를 [@fbbdev](https://github.com/fbbdev)가 [#4049](https://github.com/wailsapp/wails/pull/4049)에서 수정했습니다.
- 애셋 서버의 루트 인덱스 경로에서 404이 아닌 응답을 처리하는 방식을 [@fbbdev](https://github.com/fbbdev)가 [#4049](https://github.com/wailsapp/wails/pull/4049)에서 수정했습니다.
- 제네릭 형식의 속성을 테스트할 때 바인딩 생성기에서 발생하던 정의되지 않은 동작을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 기저 형식의 속성이 명명된 래퍼와 동일하지 않은 경우의 모델용 바인딩 생성기 출력을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 맵 키 형식 및 전처리에 대한 바인딩 생성기 출력을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 마샬러 인터페이스를 구현하는 구조체의 바인딩 생성기 출력을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 바인딩 생성기에서 제네릭 형식이 관련된 형식 순환을 감지하는 방식을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 바인딩 생성기 출력에서 내보내지 않은 모델을 잘못 참조하던 문제를 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 삽입되는 코드를 서비스 파일 끝으로 옮겼습니다. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 작업했습니다.
- 바인딩 생성기에서 파일 닫기 작업의 오류를 처리하는 방식을 [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 수정했습니다.
- 수명 주기 또는 http 메서드를 정의하지만 그 밖의 바인딩된 메서드는 없는 서비스에 대한 경고를 숨겼습니다. [@fbbdev](https://github.com/fbbdev)가 [#4045](https://github.com/wailsapp/wails/pull/4045)에서 작업했습니다.
- 밝은 시스템 색 구성표를 사용할 때 React 이외의 템플릿에서 Hello World 바닥글이 표시되지 않던 문제를 [@marcus-crane](https://github.com/marcus-crane)이 [#4056](https://github.com/wailsapp/wails/pull/4056)에서 수정했습니다.
- macOS에서 메뉴 항목이 숨겨지던 문제를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- 메시지 프로세서에서 오류를 처리하고 형식을 지정하는 방식을 [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 수정했습니다.
-  애플리케이션을 종료할 때 서비스 종료 절차를 건너뛰던 문제를 [@fbbdev](https://github.com/fbbdev)가 [#4066](https://github.com/wailsapp/wails/pull/4066)에서 수정했습니다.
-  메뉴 업데이트가 메인 스레드에서 실행되도록 [@leaanthony](https://github.com/leaanthony)가 변경했습니다.
- 드래그 및 크기 조절 메커니즘의 안정성을 높이고 예상되는 플랫폼 동작에 더 가깝게 맞췄습니다. [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 작업했습니다.
- Webpack/angular가 런타임 초기화 코드를 버리는 [#4097](https://github.com/wailsapp/wails/issues/4097) 문제를 [@fbbdev](https://github.com/fbbdev)가 [#4100](https://github.com/wailsapp/wails/pull/4100)에서 수정했습니다.
- 처음에 숨겨진 메뉴 항목과 관련된 문제를 [@IanVS](https://github.com/IanVS)가 [#4116](https://github.com/wailsapp/wails/pull/4116)에서 수정했습니다.
- 확장자가 없는 요청에서 `[request]`은 존재하지 않지만 `[request].html`은 존재하는 경우 assetFileServer가 `.html` 파일을 제공하지 않던 문제를 수정했습니다.
- 아이콘 생성 경로를 [@robin-samuel](https://github.com/robin-samuel)이 [#4125](https://github.com/wailsapp/wails/pull/4125)에서 수정했습니다.
- `fullscreen`, `unfullscreen`, `unminimise` 및 `unmaximise` 이벤트가 발생하지 않던 문제를 [@oSethoum](https://github.com/osethoum)이 [#4130](https://github.com/wailsapp/wails/pull/4130)에서 수정했습니다.
- 구성의 기본 버전에 잘못된 접두사가 붙어 발생하던 NSIS Error를 [@robin-samuel](https://github.com/robin-samuel)이 [#4126](https://github.com/wailsapp/wails/pull/4126)에서 수정했습니다.
- Windows에서 Dialogs 런타임 함수가 이스케이프된 경로를 반환하던 문제를 [TheGB0077](https://github.com/TheGB0077)이 [#4188](https://github.com/wailsapp/wails/pull/4188)에서 수정했습니다.
- HKCU의 Webview2 감지 경로를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- macOS의 입력 문제를 [@leaanthony](https://github.com/leaanthony)가 수정했습니다.
- [@yulesxoxo](https://github.com/yulesxoxo)가 [#4219](https://github.com/wailsapp/wails/pull/4219)에서 Windows 아이콘 생성 작업의 파일 이름을 수정했습니다.
- @kron의 작업을 바탕으로 [@leaanthony](https://github.com/leaanthony)가 프레임 없는 창의 투명도 문제를 수정했습니다.
- @kron의 작업을 바탕으로 [@leaanthony](https://github.com/leaanthony)가 창이 비활성화되거나 최소화되었을 때의 포커스 호출을 수정했습니다.
- @kron의 작업을 바탕으로 [@leaanthony](https://github.com/leaanthony)가 작업 표시줄을 다시 시작한 후 시스템 트레이가 표시되지 않는 문제를 수정했습니다.
- [#4245](https://github.com/wailsapp/wails/pull/4245)에서 fallbackResponseWriter가 Flush()를 구현하지 않던 문제를 수정했습니다.
- [@superDingda]가 [#4236](https://github.com/wailsapp/wails/issues/4236)에서 fallbackResponseWriter가 Flush()를 구현하지 않던 문제를 수정했습니다.
- [@joshhardy](https://github.com/joshhardy)가 [#4354](https://github.com/wailsapp/wails/pull/4354)에서 대기 중인 비동기 Go 바인딩 함수 호출이 있을 때 macOS 창을 닫으면 충돌하는 문제를 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 Windows 효율성 모드 시작 시의 경쟁 조건을 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 Windows 아이콘 핸들 정리 문제를 수정했습니다.
- [@PPTGamer](https://github.com/PPTGamer)가 [#4375](https://github.com/wailsapp/wails/pull/4375)에서 Windows의 `OpenFileManager`을 수정했습니다.
- @atterpac이 [#3979](https://github.com/wailsapp/wails/pull/3979)에서 Linux의 최소/최대 너비 옵션을 수정했습니다.
- @atterpac이 [#3966](https://github.com/wailsapp/wails/pull/3966)에서 npm 버전 상향을 통해 Typescript 템플릿의 타입 정의를 수정했습니다.
- @atterpac이 [#3945](https://github.com/wailsapp/wails/pull/3945)에서 Sveltekit 템플릿의 CSS 참조를 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 window run()의 주요 콜백이 메인 스레드에서 호출되도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)가 대화 상자 디렉터리 선택기 예제를 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 index.html이 없을 때 표시할 새로운 중국어 오류 페이지를 만들었습니다.
-  [@leaanthony](https://github.com/leaanthony)가 `windowDidBecomeKey` 콜백이 메인 스레드에서 실행되도록 했습니다.
-  [@leaanthony](https://github.com/leaanthony)가 프레임 없는 창의 전체 화면을 지원하도록 했습니다.
-  [@leaanthony](https://github.com/leaanthony)가 창 제거 로직을 개선했습니다.
-  [@leaanthony](https://github.com/leaanthony)가 시스템 트레이에 연결되었을 때의 창 위치 로직을 수정했습니다.
-  [@leaanthony](https://github.com/leaanthony)가 프레임 없는 창의 전체 화면을 지원하도록 했습니다.
- [@leaanthony](https://github.com/leaanthony)가 이벤트 처리를 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 창 종료 로직을 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 공통 Taskfile에서 Typescript 템플릿의 Typescript 바인딩을 기본적으로 생성하도록 변경했습니다.
- [@mmalcek](https://github.com/mmalcek)가 [#3990](https://github.com/wailsapp/wails/pull/3990)에서 열린 창이 없고 시스템 트레이만 사용 중일 때 WM_CLOSE 메시지로 애플리케이션이 종료되는 문제를 수정했습니다.
- @5aaee9가 [#3192](https://github.com/wailsapp/wails/pull/3192)에서 garble 빌드를 수정했습니다.
- [@leaanthony](https://github.com/leaanthony)가 Windows NSIS 빌드를 수정했습니다.
- 닫히지 않은 항목으로 인해 발생한 Linux 다중 선택 대화 상자의 교착 상태를 수정했습니다.
- Windows 빌드 중 .syso 파일을 여러 플랫폼에서 정리하는 문제를 수정한 사람:
- @atterpac이 amd64 AppImage 컴파일을 수정한 위치:
- @ansxuman이 빌드 자산 업데이트를 수정한 위치:
- @atterpac이 Linux 시스템 트레이의 `OnClick` 및 `OnRightClick` 구현을 수정했습니다.
- Mac에서 `AlwaysOnTop`이 작동하지 않는 문제를 수정한 사람:
-  중복 항목을 포함하던 `application.NewEditMenu`을 수정했습니다.
- 🐧 aarch64 컴파일을 수정했습니다.
- ⊞ 라디오 그룹 메뉴 항목을 수정한 사람:
- MacOS에서 'name'과 'outputfilename'이 다를 때 실행 가능한 .app을 빌드하면 발생하는 오류를 수정했습니다.
- 드래그 앤 드롭 예제에서 customEventProcessor를 사용할 때 발생하는 버그를 수정한 사람:
- 🐧 IgnoreMouseEvents 추가로 발생한 Linux 컴파일 오류를 수정한 사람:
- ⊞ syso 아이콘 파일 생성 버그를 수정한 사람:
- 🐧 Wayland에서 네이티브로 실행되도록 하는 수정 사항을 다음에서 반영했습니다.
- 다음에서 내부 서비스 메서드를 바인딩하지 않도록 변경했습니다.
- ⊞ 다음에서 시스템 트레이 시작 시 발생하는 패닉을 수정했습니다.
- 다음에서 내부 서비스 메서드를 바인딩하지 않도록 변경했습니다.
- ⊞ 다음에서 시스템 트레이 시작 시 발생하는 패닉을 수정했습니다.
- 메뉴 항목과 이벤트 처리를 대대적으로 리팩터링했습니다. 현재는 주로 macOS가 개선됩니다. 작업자:
- 다음에서 플러그인 및 이벤트 리팩터링 후 테스트를 수정했습니다.
- ⊞ `Failed to unregister class Chrome_WidgetWin_0` 경고를 수정했습니다. 작업자:
- 모듈 문제
- [atterpac](https://github.com/atterpac)이 다음에서 크기 조정 이벤트 메시징을 수정했습니다.
- 🐧 NixOS의 테마 처리 오류를 수정한 사람:
- Windows에서 볼륨 간 프로젝트 설치 문제를 수정한 사람:
- 바닥글이 표시되도록 React 템플릿 CSS를 수정한 사람:
- refresh를 최신 버전으로 업데이트하여 개발 모드에서 작업할 때 발생하는 좀비 프로세스를 수정했습니다.
- [Atterpac](https://github.com/atterpac)이 AppImage WebKit 파일 소싱을 수정했습니다.
- [Atterpac](https://github.com/Atterpac)이 다음에서 Doctor의 apt 패키지 검증을 수정했습니다.
- @5aaee9가 다음에서 애플리케이션 종료 시 멈추는 문제(Darwin)를 수정했습니다.
- Windows에서 예제의 배경색을 수정한 사람:
- [mmghv](https://github.com/mmghv)가 다음에서 기본 컨텍스트 메뉴를 수정했습니다.
- Darwin의 화살표 키 16진수 값을 수정한 사람:
- Windows에서 드래그 앤 드롭이 작동하도록 설정했습니다. 추가한 사람:
- 사용자에게 올바른 드라이버가 없을 때 Linux의 doctor에서 발생하는 버그를 수정했습니다.
- 시작 시 DPI 배율 조정을 수정했습니다(Windows). [@almas-x](https://github.com/almas-x)가 다음에서 변경했습니다.
- 상대 경로를 사용하도록 `go.mod`의 치환 줄을 수정했습니다. 다음이 포함된 Windows 경로 문제를 해결합니다.
- 연결된 창이 없을 때 MacOS 시스템 트레이 클릭 처리를 수정한 사람:
- 알 수 없는 옵션으로 인해 Windows 빌드가 실패하는 문제를 수정한 사람:
- 다음이 없을 때 Windows에서 시스템 트레이 아이콘을 왼쪽 클릭하면 발생하는 충돌을 수정했습니다.
- @5aaee9가 PR에서 창을 두 번 열 때 baseURL이 잘못되는 문제를 수정했습니다.
- `WebviewWindow.Restore` 메서드에서 if 분기 순서를 수정한 사람:
- 다음 상황에서 여러 `GetStartURL` 호출에 걸쳐 `startURL`을 올바르게 계산하도록 했습니다.
- `Screen` 구조체의 JS 타입이 해당 Go 타입과 일치하도록 수정한 사람:
- 등록된 이벤트를 올바르게 정리하도록 `WML.Reload` 메서드를 수정했습니다.
- Linux에서 사용자 지정 컨텍스트 메뉴가 즉시 닫히는 문제를 수정한 사람:
- 바인딩에서 생성한 모델 파일의 출력 경로와 확장자를 수정
- 바인딩에서 생성한 JS 코드에서 모델 파일의 import 경로를 수정
- 일부 Linux 배포판에서 드래그 앤 드롭이 작동하지 않는 문제를 수정한 사람:
- `wails3 task dev` 사용 시 macOS 작업이 누락되는 문제를 수정한 사람:
- 이벤트 등록 시 nil 맵 할당이 발생하는 문제를 수정한 사람:
- 바인딩된 메서드 매개변수의 역마샬링 문제를 수정한 사람:
- 바인딩된 메서드의 여러 반환값 처리 문제를 수정한 사람:
- 시스템 패키지 관리자로 설치하지 않은 npm을 doctor가 감지하지 못하는 문제를 수정
- MicrosoftEdgeWebview2Setup.exe가 누락되는 문제를 수정. 도움을 주신 분:
- 창 ID 처리로 인해 Linux에서 무작위로 충돌하는 문제를 @leaanthony가 수정. 기반 작업:
- Linux에서 systemTray.setIcon 호출 시 충돌하는 문제를 수정한 사람:
- 다음 플랫폼의 `setFrameless` 함수에서 첫 호출 시 창 프레임이 적용되도록 수정:

### 변경

- **호환성 중단**: Go 맵 의미 체계를 정확히 반영하도록 생성된 JS/TS 바인딩의 맵 키를 이제 선택 사항으로 표시합니다. 이제 Typescript에서 맵 값에 접근하면 `T` 대신 `T | undefined`가 반환되므로 null 검사나 단언이 필요합니다(#4943). 기여자: `@fbbdev`
- `@wailsio/runtime`의 변경 사항에 따라 `Event` 사용을 `Events` 사용으로 변경하고, `Features/Events/Event System` 문서의 함수 호출을 적절하게 수정. 기여자: @AbdelhadiSeddar
- `EnabledFeatures`, `DisabledFeatures`, `AdditionalBrowserArgs`를 창별 옵션에서 애플리케이션 수준의 `Options.Windows`로 이동(#4559). 기여자: @leaanthony
- `Drag N Drop` 예제의 README를 업데이트하고, 해당 예제에서 `Internal Drag and Drop`를 시연한다는 점을 강조. 기여자: @ndianabasi
- 여러 디버그 로그의 수준을 Info에서 Debug로 변경(@mbaklor)
- **호환성 중단:** 창 옵션에서 `EnableDragAndDrop`의 이름을 `EnableFileDrop`로 변경
- **호환성 중단:** 이벤트 컨텍스트에서 `DropZoneDetails`의 이름을 `DropTargetDetails`로 변경
- **호환성 중단:** `WindowEventContext`의 `DropZoneDetails()` 메서드 이름을 `DropTargetDetails()`로 변경
- **호환성 중단:** `WindowDropZoneFilesDropped` 이벤트를 제거하고 대신 `WindowFilesDropped` 사용
- **호환성 중단:** HTML 속성을 `data-wails-dropzone`에서 `data-file-drop-target`로 변경
- **호환성 중단:** CSS hover 클래스를 `wails-dropzone-hover`에서 `file-drop-target-active`로 변경
- **호환성 중단:** Windows에서 `DragEffect`, `OnEnterEffect`, `OnOverEffect` 옵션 제거(제거된 IDropTarget의 일부였음)
- 모든 런타임 JSON 처리(메서드 바인딩, 이벤트, webview 요청, 알림, kvstore)에 goccy/go-json을 사용하도록 전환하여 성능을 21-63% 개선하고 메모리 할당을 40-60% 감소
- 호출당 오버헤드를 줄이도록 BoundMethod 구조체 레이아웃을 최적화하고 isVariadic 플래그를 캐시
- 인수가 `<=8`개인 메서드에서 힙 할당을 방지하도록 스택 할당 인수 버퍼 사용
- 반환값이 하나일 때 슬라이스 할당을 방지하도록 메서드 호출의 결과 수집을 최적화
- 동시성 성능을 개선하도록 MIME 유형 캐시에 sync.Map 사용
- HTTP 전송 요청 본문을 읽을 때 버퍼 풀 사용
- 요청당 할당을 줄이도록 콘텐츠 유형 스니퍼에서 CloseNotify 채널을 지연 할당
- 애셋 서버에서 디버그 CSS 로깅 제거
- MIME 유형 확장자 맵을 확장하여 50개 이상의 일반적인 웹 형식(글꼴, 오디오, 비디오 등)을 지원
- Window `X/Y` 옵션 문서 업데이트. 기여자: @ruhuang2001
- 프런트엔드 바인딩 생성 옵션을 더 추가하여 `Frontend Runtime` 문서 업데이트. 기여자: @ndianabasi
- Wails v3 Asset Server 문서 페이지 업데이트. 기여자: @ndianabasi
- **호환성 중단**: 패키지 수준의 대화 상자 함수(`application.InfoDialog()`, `application.QuestionDialog()` 등)를 제거합니다. 대신 `app.Dialog` 관리자인 `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`을 사용하세요.
- 실제 API에 맞게 대화 상자 문서를 업데이트: `SetButtons()`가 아닌 콜백을 사용하는 `app.Dialog.*` 및 `AddButton()`, string이 아닌 `SetDefaultButton(*Button)`, `SetFilters()`가 아닌 `AddFilter()`, `SetDefaultFilename()`이 아닌 `SetFilename()`을 사용하고, 폴더 선택에는 `app.Dialog.OpenFile().CanChooseDirectories(true)` 사용
- **호환성 중단**: 이제 프로덕션 빌드가 기본값입니다. 개발 빌드를 만들려면 Taskfile에 `DEV=true`를 설정하세요. 예제용 새 프로젝트 생성. 기여자: @leaanthony
- 데이터 인수가 0개 또는 1개인 사용자 지정 이벤트를 내보낼 때 데이터 값을 슬라이스로 감싸지 않고 Data 필드에 직접 할당합니다. 기여자: [@fbbdev](https://github.com/fbbdev), 위치: [#4633](https://github.com/wailsapp/wails/pull/4633)
- 이제 Windows 트레이는 `NIS_HIDDEN`를 전환하여 `SystemTray.Show()`/`Hide()`을 준수하므로 앱이 완전히 사라졌다가 다시 나타날 수 있습니다(#4653).
- 트레이 등록 시 확인된 아이콘을 재사용하고, `NOTIFYICON_VERSION_4`을 한 번 설정하며, Explorer 재시작 후 툴팁이 복구되도록 `NIF_SHOWTIP`을 활성화합니다(#4653).
- macOS: 메뉴 막대와 Dock 영역을 제외하고 창을 가운데 정렬하도록 `frame` 대신 `visibleFrame` 사용
- macOS: 메뉴 막대와 Dock 영역을 제외하고 창을 가운데 정렬하도록 `frame` 대신 `visibleFrame` 사용
- `-config` 매개변수와 함께 `wails3 update build-assets`을 실행할 때 `-product*` 매개변수를 통해 설정한 값은
- `window.NativeWindowHandle()` -> `window.NativeWindow()`. 기여자: @leaanthony, 위치: [#4471](https://github.com/wailsapp/wails/pull/4471)
- 내부 창 처리를 리팩터링. 기여자: @leaanthony, 위치: [#4471](https://github.com/wailsapp/wails/pull/4471)
- `application.WindowIDKey` 및 `application.WindowNameKey` 제거(`application.WindowKey`로 대체). 기여자: [@leaanthony](https://github.com/leaanthony)
- 이제 ContextMenuData가 any 대신 string을 반환. 기여자: [@leaanthony](https://github.com/leaanthony)
- 이제 JS/TS 바인딩에서 고정 길이 배열 유형의 클래스 필드는 빈 상태가 아니라 예상 길이로 초기화됩니다. 기여자: [@fbbdev](https://github.com/fbbdev), 위치: [#4001](https://github.com/wailsapp/wails/pull/4001)
- 이제 ContextMenuData가 any 대신 string을 반환. 기여자: [@leaanthony](https://github.com/leaanthony)
- 이제 `application.NewService`은 options를 선택적 매개변수로 받지 않습니다(대신 `application.NewServiceWithOptions` 사용). 기여자: [@leaanthony](https://github.com/leaanthony), 위치: [#4024](https://github.com/wailsapp/wails/pull/4024)
- `nanoid` 종속성 제거. 기여자: [@leaanthony](https://github.com/leaanthony)
- mica/acrylic/tabbed 창 스타일에 맞게 Window 예제 업데이트. 기여자: [@leaanthony](https://github.com/leaanthony)
- JS/TS 바인딩에서 `internal.js/ts` 모델 파일을 제거했으며, 이제 모든 모델은 `models.js/ts`에서 찾을 수 있습니다. 기여자: [@fbbdev](https://github.com/fbbdev), 위치: [#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS 바인딩에서 명명된 유형은 다른 명명된 유형의 별칭으로 렌더링되지 않으며, 기존 동작은 이제 별칭에만 적용됩니다. 기여자: [@fbbdev](https://github.com/fbbdev), 위치: [#4045](https://github.com/wailsapp/wails/pull/4045)
- JS/TS 바인딩의 클래스 모드에서 형식 매개변수 형식의 구조체 필드는 선택 사항으로 표시되며 자동으로 초기화되지 않도록 변경됨 — [#4045](https://github.com/wailsapp/wails/pull/4045)에서 [@fbbdev](https://github.com/fbbdev)
- 템플릿에서 ESLint 제거 — [#4059](https://github.com/wailsapp/wails/pull/4059)에서 [@IanVS](https://github.com/IanVS)
- 저작권 연도를 2025(으)로 업데이트 — [#4037](https://github.com/wailsapp/wails/pull/4037)에서 [@IanVS](https://github.com/IanVS)
- event.Sender 문서 추가 — [#4075](https://github.com/wailsapp/wails/pull/4075)에서 [@IanVS](https://github.com/IanVS)
- Go 1.24 지원 — [@leaanthony](https://github.com/leaanthony)
- 이제 `ServiceStartup` 훅은 `application.New`에서가 아니라 `App.Run` 호출 시 실행됨 — [#4066](https://github.com/wailsapp/wails/pull/4066)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 `ServiceStartup` 오류가 프로세스를 종료하는 대신 `App.Run`에서 반환됨 — [#4066](https://github.com/wailsapp/wails/pull/4066)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 JS의 바인딩 및 대화 상자 호출이 문자열 대신 오류 객체로 거부됨 — [#4066](https://github.com/wailsapp/wails/pull/4066)에서 [@fbbdev](https://github.com/fbbdev)
- Windows에서 시스템 트레이 메뉴 위치 지정 개선 — [@leaanthony](https://github.com/leaanthony)
- JS 런타임을 TypeScript로 포팅 — [#4100](https://github.com/wailsapp/wails/pull/4100)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 런타임을 가져오는 즉시 초기화하므로 창이 로드될 때까지 기다릴 필요가 없음 — [#4100](https://github.com/wailsapp/wails/pull/4100)에서 [@fbbdev](https://github.com/fbbdev)
- 런타임에서 더 이상 init 메서드를 내보내지 않음. 부수 효과 가져오기를 사용해 초기화할 수 있음 — [#4100](https://github.com/wailsapp/wails/pull/4100)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 바인딩된 메서드는 취소될 경우 `CancelError`로 거부되는 `CancellablePromise`을 반환함. 실제 호출 결과는 폐기됨 — [#4100](https://github.com/wailsapp/wails/pull/4100)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 기본 제공 서비스 형식의 이름을 `Service`(으)로 통일 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- 이제 옵션을 받는 기본 제공 서비스 생성 함수의 이름을 `NewWithConfig`(으)로 통일 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- Go API와의 일관성을 위해 `sqlite` 서비스의 `Select` 메서드 이름을 `Query`(으)로 변경 — [#4067](https://github.com/wailsapp/wails/pull/4067)에서 [@fbbdev](https://github.com/fbbdev)
- 템플릿: 런타임을 "dependencies"로 이동하고 package.json 파일 정리 — [#4133](https://github.com/wailsapp/wails/pull/4133)에서 [@IanVS](https://github.com/IanVS)
- 특정 macOS API를 사용할 수 있도록 개발 환경에서 앱 번들을 생성하고 임시 서명함 — [#4171](https://github.com/wailsapp/wails/pull/4171)에서 [@popaprozac](https://github.com/popaprozac)
- 빌드 자산을 플랫폼별 디렉터리로 이동 — [@leaanthony](https://github.com/leaanthony)
- Taskfile을 플랫폼별 디렉터리로 이동하고 이름 변경 — [@leaanthony](https://github.com/leaanthony)
- `index.html`이 없을 때의 사용자 경험을 크게 개선 — [@leaanthony](https://github.com/leaanthony)
- [Windows] 최소화 및 복원 성능 개선 — [@leaanthony](https://github.com/leaanthony). [562589540](https://github.com/562589540)의 원본 [PR](https://github.com/wailsapp/wails/pull/3955)을 기반으로 함
- `ShouldClose` 옵션 제거(대신 events.Common.WindowClosing 이벤트에 훅 등록) — [@leaanthony](https://github.com/leaanthony)
- [Windows] 창을 열 때 발생하는 깜박임 감소 — [@leaanthony](https://github.com/leaanthony)
- 내부 함수로 사용하도록 의도된 `Window.Destroy` 제거 — [@leaanthony](https://github.com/leaanthony)
- `WindowClose` 이벤트의 이름을 `WindowClosing`(으)로 변경 — [@leaanthony](https://github.com/leaanthony)
- 이제 프런트엔드 빌드에서 빌드 형식에 따라 vite 환경 "development" 또는 "production"을 사용 — [@leaanthony](https://github.com/leaanthony)
- go-webview2 v1.19(으)로 업데이트 — [@leaanthony](https://github.com/leaanthony)
- Taskfile 포크가 사용되도록 보장 — @leaanthony
- 다음을 사용해 설치할 때 발생하는 버전 문제를 해결하도록 Taskfile 포크 업데이트
- 다음을 사용해 설치할 때 발생하는 버전 문제를 해결하기 위해 Taskfile 포크 사용
- 이제 `service.OnStartup`이 오류 발생 시 애플리케이션을 종료하고 다음을 실행함
- 사용자 상호 작용에 더 잘 부합하도록 시스템 트레이 클릭 메시징 리팩터링 —
- 다음을 생성하는 프레임워크를 지원하도록 자산 임베드에 `all:frontend/dist` 포함
- Taskfile 리팩터링 — [leaanthony](https://github.com/leaanthony), 다음에서
- `go-webview2` v1.0.16(으)로 업그레이드 —
- `Screen` 형식이 `Id`이 아닌 `ID`을 포함하도록 수정 —
- `application.ServiceOptions`을 지원하도록 `go.mod.tmpl` Wails 버전 업데이트 —
- 서비스 이름 결정 수정 — [windom](https://github.com/windom/), 다음에서
- 이제 mkdocs serve에서 docker 사용 — [leaanthony](https://github.com/leaanthony)
- 개발 구성을 `config.yml`에 통합 —
- 이제 시스템 트레이 대화 상자에서 애플리케이션 아이콘을 사용할 수 있으면 기본값으로 사용(Windows) —
- macOS의 GPU 및 메모리 보고 개선 —
- `WebviewGpuIsDisabled` 및 `EnableFraudulentWebsiteWarnings` 제거
- 이벤트 API 변경: `On`/`Emit` -> 사용자 이벤트, `OnApplicationEvent` ->
- Linux의 이벤트 API 수정 — [TheGB0077](https://github.com/TheGB0077), 다음에서
- [CI] 작업 개선 및 포크에서도 작업을 실행할 수 있도록 설정하고
- `AbsolutePosition()`의 이름을 `Position()`(으)로 변경 —
- Linux WebKit 종속성을 webkitgtk2-4.0 대신 webkit2gtk-4.1(으)로 업데이트하여
- 이제 번들된 JS 런타임 스크립트가 ESM 모듈임: 이를 가져오는 script 태그는
- `@wailsio/runtime` 패키지는 `window.wails`에 API를 게시하지 않음
- 이제 Window API 모듈 `@wailsio/runtime/src/window`에서 이를 포함하는 항목을 노출함
- 현재 Go `WebviewWindow`과 일치하도록 JS 창 API 업데이트
- 이제 바인딩 생성기에서 기본적으로 ID 기반 호출을 사용함. `-id` CLI 옵션은
- 새로운 바인딩 코드 레이아웃: 이전에는 출력 파일이 폴더에 정리되어 있었음
- 구조체 필드 `application.Options.Bind`의 이름을 다음으로 변경함
- 바인딩 서비스의 새로운 구문: 이제 서비스 인스턴스를 다음으로 래핑해야 함
- 터미널이 아닌 환경 또는 CI 환경에서 스피너 비활성화 —

### 제거됨

- **주요 변경**: 창별 `WindowsWindow` 옵션에서 `EnabledFeatures`, `DisabledFeatures`, `AdditionalLaunchArgs`를 제거했습니다. 대신 애플리케이션 수준의 `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures`, `Options.Windows.AdditionalBrowserArgs`를 사용하세요. 이러한 플래그는 공유 WebView2 환경 전체에 적용됩니다(#4559). 작성자: @leaanthony
- Windows의 네이티브 `IDropTarget` 구현을 제거하고 JavaScript 기반 방식으로 대체했습니다(v2 동작과 일치).
- github.com/wailsapp/mimetype 종속성을 제거하고 확장된 확장자 맵과 표준 라이브러리의 http.DetectContentType을 사용하여 바이너리 크기를 약 1.2MB 줄였습니다.
- Linux 파일 탐색기를 위한 최소한의 .desktop 파일 파서를 구현하여 gopkg.in/ini.v1 종속성을 제거하고 약 45KB를 절감했습니다.
- Go 1.21+ 표준 라이브러리의 slices 패키지와 최소한의 내부 도우미를 사용하여 런타임 코드에서 samber/lo를 제거하고 약 310KB를 절감했습니다.
- Darwin URL 스킴 핸들러에서 디버그 printf 문을 제거했습니다(#4834).
- **주요 변경**: `linux:WindowLoadChanged` 이벤트를 제거했습니다. WebView의 로딩 완료를 감지하려면 대신 `linux:WindowLoadFinished`을 사용하세요(#3896). 작성자: @leaanthony

### 주요 변경 사항

- **Manager API 리팩터링**: 코드 구성과 검색 편의성을 개선하기 위해 애플리케이션 API를 평면 구조에서 체계적인 manager 구조로 재구성했습니다. [#4359](https://github.com/wailsapp/wails/pull/4359)에서 [@leaanthony](https://github.com/leaanthony)가 작업했습니다.
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Service 메서드의 이름을 변경했습니다: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`. 작성자: [@leaanthony](https://github.com/leaanthony)
- `Path` 및 `Paths` 메서드를 `application` 패키지로 이동했습니다. 작성자: [@leaanthony](https://github.com/leaanthony)
- 이제 애플리케이션 메뉴는 macOS에서만 지원됩니다. 작성자: [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## 추가

## 수정

## v3.0.0-alpha.77 - 2026-04-18

## 수정

## v3.0.0-alpha.76 - 2026-04-17

## 수정

## v3.0.0-alpha.75 - 2026-04-16

## 수정

## v3.0.0-alpha.74 - 2026-03-01

## 추가

## 수정

## v3.0.0-alpha.73 - 2026-02-27

## 수정

## v3.0.0-alpha.72 - 2026-02-16

## 수정

## v3.0.0-alpha.71 - 2026-02-10

## 추가

## 수정

## v3.0.0-alpha.70 - 2026-02-09

## 추가

## 수정

## v3.0.0-alpha.69 - 2026-02-08

## 추가

## 수정

## v3.0.0-alpha.68 - 2026-02-07

## 추가

## 변경

## 수정

## v3.0.0-alpha.67 - 2026-02-04

## 추가

## 변경

## 수정

## v3.0.0-alpha.66 - 2026-02-03

## 추가

## 변경

## 수정

## 제거

## v3.0.0-alpha.65 - 2026-02-01

## 추가

## v3.0.0-alpha.64 - 2026-01-26

## 추가

## v3.0.0-alpha.63 - 2026-01-25

## 수정

## v3.0.0-alpha.62 - 2026-01-22

## 수정

## v3.0.0-alpha.61 - 2026-01-20

## 수정

## v3.0.0-alpha.60 - 2026-01-14

## 수정

## v3.0.0-alpha.59 - 2026-01-11

## 변경

## v3.0.0-alpha.58 - 2026-01-09

## 수정

## v3.0.0-alpha.57 - 2026-01-05

## 변경

## 수정됨

## v3.0.0-alpha.56 - 2026-01-04

## 추가됨

## 변경됨

## 수정됨

## 제거됨

## v3.0.0-alpha.55 - 2026-01-02

## 변경됨

## 수정됨

## 제거됨

## v3.0.0-alpha.54 - 2025-12-29

## 추가됨

## 수정됨

## 제거됨

## v3.0.0-alpha.53 - 2025-12-27

## 추가됨

## 수정됨

## v3.0.0-alpha.52 - 2025-12-26

## 수정됨

## v3.0.0-alpha.51 - 2025-12-23

## 수정됨

## v3.0.0-alpha.50 - 2025-12-21

## 변경됨

## v3.0.0-alpha.49 - 2025-12-18

## 변경됨

## v3.0.0-alpha.48 - 2025-12-16

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.47 - 2025-12-15

## 추가됨

## 수정됨

## v3.0.0-alpha.46 - 2025-12-14

## 추가됨

## 제거됨

## v3.0.0-alpha.45 - 2025-12-13

## 추가됨

## 수정됨

## v3.0.0-alpha.44 - 2025-12-12

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.43 - 2025-12-11

## 추가됨

## v3.0.0-alpha.42 - 2025-12-10

## 추가됨

## v3.0.0-alpha.41 - 2025-11-23

## 수정됨

## v3.0.0-alpha.40 - 2025-11-13

## 수정됨

## v3.0.0-alpha.39 - 2025-11-12

## 추가됨

## 변경됨

## v3.0.0-alpha.38 - 2025-11-04

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.37 - 2025-11-02

## 수정됨

## v3.0.0-alpha.36 - 2025-10-15

## 수정됨

## v3.0.0-alpha.35 - 2025-10-14

## 수정됨

## v3.0.0-alpha.34 - 2025-10-06

## 추가됨

## 수정됨

## v3.0.0-alpha.33 - 2025-10-04

## 수정됨

## v3.0.0-alpha.32 - 2025-10-02

## 수정됨

## v3.0.0-alpha.31 - 2025-09-27

## 수정됨

## v3.0.0-alpha.30 - 2025-09-26

## 수정됨

## v3.0.0-alpha.29 - 2025-09-25

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.29 - 2025-09-25

## 추가됨

## 변경됨

## 수정됨

## v3.0.0-alpha.27 - 2025-09-07

## 수정됨

## v3.0.0-alpha.26 - 2025-08-24

## 추가됨

## v3.0.0-alpha.25 - 2025-08-16

## 변경됨

더 이상 무시되지 않으며 구성 값을 재정의합니다.

## v3.0.0-alpha.24 - 2025-08-13

## 추가됨

## v3.0.0-alpha.23 - 2025-08-11

## 수정됨

## v3.0.0-alpha.22 - 2025-08-10

## 추가됨

## 변경됨

+ 지나치게 광범위한 Linux 패키지 종속성을 수정하고 오래된 RPM 종속성을 수정했습니다.

## v3.0.0-alpha.21 - 2025-08-07

## 수정됨

## v3.0.0-alpha.20 - 2025-08-06

## 수정됨

## v3.0.0-alpha.19 - 2025-08-05

## 추가됨

## 수정됨

## v3.0.0-alpha.18 - 2025-08-03

## 추가됨

## 수정됨

## v3.0.0-alpha.17 - 2025-07-31

## 수정됨

## v3.0.0-alpha.16 - 2025-07-25

## 추가됨

## v3.0.0-alpha.15 - 2025-07-25

## 추가됨

## v3.0.0-alpha.14 - 2025-07-25

## 추가됨

## v3.0.0-alpha.12 - 2025-07-15

### 추가됨

### 수정됨

## v3.0.0-alpha.11 - 2025-07-12

## 추가됨

## v3.0.0-alpha.10 - 2025-07-06

### 주요 변경 사항

### 추가됨

### 수정됨

### 변경됨

## v3.0.0-alpha.9 - 2025-01-13

### 추가됨

### 수정됨

### 변경됨

## v3.0.0-alpha.8.3 - 2024-12-07

### 변경됨

## v3.0.0-alpha.8.2 - 2024-12-07

### 변경됨

`go install` — @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### 변경됨

`go install` — @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### 추가됨

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 변경됨

`service.OnShutdown`이전에 시작된 모든 서비스에 대해 — @atterpac,   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac, [#3907](https://github.com/wailsapp/wails/pull/3907)   하위 폴더 — @atterpac,   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913),   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes),   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (`EnabledFeatures` 및 `DisabledFeatures` 옵션으로 대체됨) —   [leaanthony](https://github.com/leaanthony)

### 수정됨

채널 변수 — @michael-freling,   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman),   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   [#3886](https://github.com/wailsapp/wails/pull/3886)에서   [leaanthony](https://github.com/leaanthony),   [#3841](https://github.com/wailsapp/wails/pull/3841)   Darwin의 편집 메뉴에서 `PasteAndMatchStyle` 역할 —   [johnmccabe](https://github.com/johnmccabe),   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840),   [#3854](https://github.com/wailsapp/wails/pull/3854) —   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   서로 다름 — @nickisworking,   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### 추가됨

[mmghv](https://github.com/mmghv),   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) 및   [leaanthony](https://github.com/leaanthony),   [#3570](https://github.com/wailsapp/wails/pull/3570)

### 변경됨

Application Events `OnWindowEvent` -> Window Events —   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   `v3/` 또는 `v3-` 접두사가 붙은 브랜치 —   [stendler](https://github.com/stendler),   [#3747](https://github.com/wailsapp/wails/pull/3747)

### 수정됨

[etesam913](https://github.com/etesam913),   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac),   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac),   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811),   [#3614](https://github.com/wailsapp/wails/pull/3614) —   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) —   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) —   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) —   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) —   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) —   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 수정됨

## v3.0.0-alpha.5 - 2024-07-30

### 추가됨

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   아이콘 클릭 — @5aaee9, [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev),   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac),   [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane),   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev),   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   @Mai-Lapyst의 [PR](https://github.com/wailsapp/wails/pull/2044)을 기반으로 함   [@fbbdev](https://github.com/fbbdev),   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev),   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev),   [#3295](https://github.com/wailsapp/wails/pull/3295)   npm 패키지 — [@fbbdev](https://github.com/fbbdev),   [#3334](https://github.com/wailsapp/wails/pull/3334)   [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` — [@abichinger](https://github.com/abichinger),   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger),   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev),   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy),   [#3674](https://github.com/wailsapp/wails/pull/3674)

### 수정됨

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane)이   [#3515](https://github.com/wailsapp/wails/pull/3515)에서 작업함   [atterpac](https://github.com/atterac)이   [#3512](https://github.com/wailsapp/wails/pull/3512)에서 작업함   [atterpac](https://github.com/atterpac)이   [#3477](https://github.com/wailsapp/wails/pull/3477)에서 작업함   [Atterpac](https://github.com/atterpac)이   [#3320](https://github.com/wailsapp/wails/pull/3320)에서 작업함.   [#3306](https://github.com/wailsapp/wails/pull/3306)에서 작업함.   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv)가   [#2750](https://github.com/wailsapp/wails/pull/2750)에서 작업함.   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave)가   [#3052](https://github.com/wailsapp/wails/pull/3052)에서 작업함.   [@pylotlight](https://github.com/pylotlight)가   [PR](https://github.com/wailsapp/wails/pull/3039)에서 작업함   설치됨. [@pylotlight](https://github.com/pylotlight)가   [PR](https://github.com/wailsapp/wails/pull/3032)에서 추가함   [PR](https://github.com/wailsapp/wails/pull/3145)   공백 - @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal)이 PR   [#3207](https://github.com/wailsapp/wails/pull/3207)에서 작업함   [thomas-senechal](https://github.com/thomas-senechal)이 PR   [#3208](https://github.com/wailsapp/wails/pull/3208)에서 작업함   연결된 창 [tw1nk](https://github.com/tw1nk)가 PR   [#3271](https://github.com/wailsapp/wails/pull/3271)에서 작업함   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev)가   [#3279](https://github.com/wailsapp/wails/pull/3279)에서 작업함   `FRONTEND_DEVSERVER_URL`이 있음.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 작업함   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 리스너를 작업함   [@abichinger](https://github.com/abichinger)가   [#3330](https://github.com/wailsapp/wails/pull/3330)에서 작업함   [@fbbdev](https://github.com/fbbdev)가   [#3334](https://github.com/wailsapp/wails/pull/3334)에서 생성기를 작업함   [@fbbdev](https://github.com/fbbdev)가   [#3334](https://github.com/wailsapp/wails/pull/3334)에서 생성기를 작업함   [@abichinger](https://github.com/abichinger)가   [#3346](https://github.com/wailsapp/wails/pull/3346)에서 작업함   [@hfoxy](https://github.com/hfoxy)가   [#3417](https://github.com/wailsapp/wails/pull/3417)에서 작업함   [@hfoxy](https://github.com/hfoxy)가   [#3426](https://github.com/wailsapp/wails/pull/3426)에서 작업함   [@fbbdev](https://github.com/fbbdev)가   [#3431](https://github.com/wailsapp/wails/pull/3431)에서 작업함   [@fbbdev](https://github.com/fbbdev)가   [#3431](https://github.com/wailsapp/wails/pull/3431)에서 작업함   [@pekim](https://github.com/pekim)이   [#3458](https://github.com/wailsapp/wails/pull/3458)에서 작업함   [@robin-samuel](https://github.com/robin-samuel).   [@5aaee9](https://github.com/5aaee9)의 PR [#3466](https://github.com/wailsapp/wails/pull/3622).   [@windom](https://github.com/windom/)이   [#3636](https://github.com/wailsapp/wails/pull/3636)에서 작업함.   [@bruxaodev](https://github.com/bruxaodev/)가   [#3691](https://github.com/wailsapp/wails/pull/3691)에서 Windows를 작업함.

### 변경됨

[mmghv](https://github.com/mmghv)가   [#3611](https://github.com/wailsapp/wails/pull/3611)에서 작업함   [atterpac](https://github.com/atterpac)이   [#3461](https://github.com/wailsapp/wails/pull/3461)에서 Ubuntu 24.04 LTS 지원을 작업함   `type="module"` 속성이 있어야 함.   [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 작업함   객체이며 WML 시스템을 시작하지 않음. 캡슐화를 개선하기 위해 이렇게 변경함. 필요한 경우 새로운 `WML.Enable` 메서드를 호출하여 WML 시스템을 수동으로 시작할 수 있음. 번들로 제공되는 JS 런타임 스크립트는 여전히 두 작업을 모두 자동으로 수행함. [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 작업함   window 객체를 기본 내보내기로 제공함. 이제 ESM 명명 가져오기 또는 네임스페이스 가져오기 구문을 통해 개별 메서드를 가져올 수 없음.   API. 일부 메서드의 이름이나 프로토타입이 변경됨. 구체적으로 `Screen`은 `GetScreen`이 되고, `GetZoomLevel`/`SetZoomLevel`은 `GetZoom`/`SetZoom`이 되며,   `GetZoom`, `Width` 및 `Height`은 이제 값을 객체로 감싸지 않고 직접 반환함. [@fbbdev](https://github.com/fbbdev)가   [#3295](https://github.com/wailsapp/wails/pull/3295)에서 작업함   제거됨. 이름으로 호출하는 방식으로 되돌리려면 `-names` CLI 옵션을 사용함.   [@fbbdev](https://github.com/fbbdev)가   [#3468](https://github.com/wailsapp/wails/pull/3468)에서 작업함   이전에는 포함된 패키지의 이름을 따랐으나, 이제 모듈 경로를 포함한 전체 Go 가져오기 경로를 사용함. [@fbbdev](https://github.com/fbbdev)가   [#3468](https://github.com/wailsapp/wails/pull/3468)에서 작업함   `application.Options.Services`. [@fbbdev](https://github.com/fbbdev)가   [#3468](https://github.com/wailsapp/wails/pull/3468)에서 작업함   `application.NewService` 호출. [@fbbdev](https://github.com/fbbdev)가   [#3468](https://github.com/wailsapp/wails/pull/3468)에서 작업함   [@DeltaLaboratory](https://github.com/DeltaLaboratory)가   [#3574](https://github.com/wailsapp/wails/pull/3574)에서 작업함
