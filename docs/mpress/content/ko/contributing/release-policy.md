---
title: "릴리스 정책"
description: "Wails v3의 나이틀리 및 위클리 릴리스 버전 관리 방식, 관리자가 릴리스하는 방법, 사용자가 버전을 고정해야 하는 이유를 설명합니다."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## 표 하나로 보는 모델

Wails v3 릴리스는 **기능이 아닌 일정에 따라** 진행됩니다. `v3.0.0` 기준 버전 이후에는 다음과 같습니다.

| 릴리스 | 태그 | 주기 | GitHub 릴리스 | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| 나이틀리 | `v3.<next>.0-nightly.<YYYYMMDD>`(UTC 날짜, 같은 날 다시 실행하면 `.N` 추가) | 자동으로 매일 실행하며, `UNRELEASED_CHANGELOG.md`에 내용이 있을 때만 릴리스 | 시험판 릴리스, Latest가 **아님** | `v`을 제외하면 같은 버전, dist-tag `nightly` |
| 위클리 | `v3.<next>.0` | 자동으로 일요일 20:00 Australia/Sydney에 실행 **(요일 미확정)** | 정식 릴리스, Latest | 같은 버전, dist-tag `latest` |
| 핫픽스 | `v3.<current>.<N>`, `N > 0` | 드물게 수동으로 실행 | 정식 릴리스, Latest | 같은 버전, dist-tag `latest` |
| 메이저 | `v4.0.0` | 호환성을 깨는 변경이 있을 때만 실행 | 정식 릴리스 | `latest` |

`<next>`은 릴리스된 **최고** v3 마이너 버전에 1을 더한 버전입니다. 나이틀리는 **다음 주간 릴리스**의 프리뷰이므로, 가장 높은 릴리스 태그를 기준으로 합니다(프리릴리스는 무시함). 따라서 핫픽스 `v3.1.1` 이후에도 다음 나이틀리는 여전히 `v3.2.0-nightly.*`입니다. 프리릴리스는 해당 정식 릴리스보다 낮게 정렬되므로 더 높은 버전으로 지정해야 합니다. `v3.0.0-nightly.<date>`은 `v3.0.0`보다 앞에 정렬됩니다.

**마이너 버전 증가는 기능 마일스톤이 아닙니다.** `v3.4.0`은 "`v3.3.0` 다음의 위클리 릴리스"를 뜻합니다. 변경 사항을 확인하려면 버전 번호가 아니라 릴리스 노트를 읽으십시오.

정책은 버전의 **형식**에 따라 결정됩니다. `-nightly.` 접미사가 있으면 나이틀리이고, 그 밖의 접미사(`-alpha`, `-beta`, `-rc`)가 있으면 GA 이전의 레거시 빌드이며, 접미사가 없으면 릴리스입니다.

### 동일하게 맞춘 버전

CLI, Go 모듈 및 npm 패키지 `@wailsio/runtime`에는 `v3/internal/version/version.txt`에서 가져온 **동일한 버전**이 지정됩니다(태그가 `v3.2.0-nightly.20261102`이면 npm에서는 `3.2.0-nightly.20261102`). `wails3 init`으로 생성한 프로젝트는 `@wailsio/runtime`을 CLI 자체 버전에 고정하며, CLI 개발 빌드는 `latest`을 대체 값으로 사용합니다.

## 실제 예시

| 날짜 | 이벤트 | 태그 | `go install ...@latest`의 해석 결과 |
| --- | --- | --- | --- |
| 25 10월 일요일 | GA, 수동 릴리스 | `v3.0.0` | `v3.0.0` |
| 10월 26 월요일~31 토요일 | 나이틀리 | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| Sydney 기준 11월 1 일요일 20:00 | 위클리 | `v3.1.0` | `v3.1.0` |
| 11월 2 월요일 | 나이틀리(다음 위클리의 미리보기) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| 11월 4 수요일 | 릴리스 브랜치에서 수동으로 진행하는 핫픽스 | `v3.1.1` | `v3.1.1` |
| 11월 8 일요일 | 위클리 | `v3.2.0` | `v3.2.0` |

위클리는 항상 같은 버전 번호의 모든 나이틀리보다 높은 순서로 정렬되며, 나이틀리는 절대로 `@latest`이 되지 않습니다.

## 사용자가 실행할 명령

| 목표 | 명령 |
| --- | --- |
| 최신 릴리스(위클리 또는 핫픽스) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` 또는 `wails3 update cli` |
| **재현 가능(권장)** | `go.mod`, CI 및 문서에서 정확한 태그 `v3.1.0`으로 고정 |
| 나이틀리 | 정확한 태그(`@v3.2.0-nightly.20261102`), `wails3 update cli -nightly` |
| 종류와 관계없이 최신 빌드 | `wails3 update cli -pre` |
| 런타임 npm 패키지 | `@wailsio/runtime`(`latest` = 위클리/핫픽스, `nightly` = 나이틀리) |

이 방식이 작동하는 이유는 Go가 `@latest`을 **시험판 접미사가 없는** 가장 높은 태그로 해석하고, 릴리스가 없을 때만 시험판 릴리스를 고려하기 때문입니다. GitHub의 "pre-release" 체크박스는 Go 프록시에서 확인할 수 없으므로, 나이틀리는 이 체크박스에 의존하지 않고 semver 시험판 접미사를 사용합니다. 프로덕션에서는 항상 버전을 고정하십시오.

## 브랜치

| 브랜치 | 용도 |
| --- | --- |
| `master` | 개발 및 나이틀리용입니다. |
| `releases/**` | `v3.0.0` GA 및 핫픽스 같은 일회성 릴리스를 위한 단기 브랜치입니다. `master`과 동일한 CI가 필수입니다. |

위클리는 `master`에서 만들지만 **그 브랜치의 최신 커밋에서 만들지는 않습니다**. 워크플로는 이전 위클리 이후의 커밋 중 필수 검사(`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`)를 모두 통과한 가장 최신 커밋을 선택하고 정확히 그 커밋에 태그를 지정합니다. 통과한 커밋이 하나도 없으면 아무것도 릴리스하지 않습니다. 나이틀리는 `master`에서 만들며, HEAD가 검사를 통과하지 않았을 때도 경고만 표시합니다.

## 릴리스 방법

릴리스 권한은 프로젝트 Lead에게만 있습니다(거버넌스 문서 참조). 기여자는 태그를 지정하거나 게시하지 않으며, `v3/UNRELEASED_CHANGELOG.md`에 항목을 추가합니다.

모든 작업은 **Nightly + Weekly Release v3**(`.github/workflows/nightly-release-v3.yml`)에서 실행됩니다.

- <strong>나이틀리(자동)</strong>는 UTC 15:00에 실행됩니다. 미릴리스 변경 로그가 비어 있으면 건너뜁니다.
- <strong>주간 릴리스(자동)</strong>는 시드니 시간 20:00에 실행됩니다. GitHub cron은 UTC만 지원하므로 두 개의 cron(UTC 09:00 및 10:00)을 예약하며, 가드 작업은 시드니에서 20:00일 때만 진행합니다(DST 반영). 주간 릴리스는 릴리스 게이트를 통과한 뒤 성공 상태인 최신 커밋을 선택하고, 마이너 버전을 올린 다음, 변경 로그를 건드리지 않고 해당 커밋에 태그를 지정합니다. 이후 `master`에서 버전 필드(`version.txt` 및 런타임 패키지 버전)만 업데이트하므로 선택된 커밋 이후에 병합된 종속성과 잠금 파일 업데이트가 유지됩니다. npm은 `master`가 아니라 태그에서 게시합니다.
- **수동.** `channel`, `dry_run`(기본적으로 켜짐), 선택적으로 `release_version`을 지정해 워크플로를 실행합니다. 시험 실행은 GitHub 또는 npm을 변경하지 않고 버전, 종류, GitHub 플래그, npm dist-tag 및 노트를 출력합니다.
- **GA 또는 핫픽스.** `release_version`(`v3.0.0`, `v3.1.1`)을 설정합니다. 이전 릴리스 계열의 핫픽스에는 `ref: releases/...`도 설정할 수 있으며, 해당 릴리스 커밋은 `master`에 다시 체리픽됩니다.

로컬 드라이런: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`. 파생된 플래그는 `go run release.go --print-policy v3.2.0-nightly.20261102`에서 확인할 수 있습니다.

### 안전장치

- 자동으로 정한 버전은 기존의 모든 태그보다 반드시 커야 하며(유의적 버전 순서에서 시험판 < 정식 릴리스), 그렇지 않으면 도구가 실행을 거부합니다.
- 한 번에 하나의 릴리스 작업만 실행됩니다(공유 동시 실행 그룹).
- npm 게시는 릴리스의 **태그**에 의해 트리거되며, 태그가 지정된 커밋을 빌드합니다.

## 릴리스 게이트

다음 항목 중 하나라도 있으면 **주간** 릴리스는 태그 생성을 거부합니다. 나이틀리 릴리스에는 절대 적용되지 않습니다.

1. `P0` 또는 `release-blocker` 레이블이 지정된 열린 이슈나 풀 리퀘스트;
2. 지난 주간 릴리스 이후 생성되었으며 `regression` 레이블이 지정된 열린 이슈;
3. Dependabot의 치명적 또는 높음 심각도 열린 알림, 코드 스캔의 치명적 또는 높음 심각도 열린 알림, 모든 열린 비밀 스캔 알림. 이러한 알림은 <strong>출시되는 코드에 한해서만</strong> 적용됩니다. `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/` 또는 `examples/` 디렉터리에만 있는 알림은 무시합니다(`v3/tasks/releasegate/main.go`의 `SHIPPED_ALERT_IGNORE`, <strong>소유자가 검토</strong>);
4. 아직 `triage` 상태인 저장소 보안 권고.

**차단자를 설정할 수 있는 사람.** 레이블은 가장 최근의 `labeled` 이벤트(그 이후에 `unlabeled` 이벤트가 없는 경우)를 `wailsapp/maintainers` 팀의 구성원이나 `maintain`/`admin` 권한을 가진 협업자가 수행한 경우에만 차단자로 인정됩니다. 다른 사람이 같은 레이블을 적용하면 무시되고 로그에 기록되므로 외부인은 릴리스를 차단하거나 차단을 해제할 수 없습니다. 자동으로 적용되는 `high-priority` 레이블은 의도적으로 사용하지 **않습니다**. `taliesin-ai` 같은 봇 계정도 `maintain` 권한이 있으면 허용 목록 검사를 통과합니다. 이는 의도된 동작이며, 트리아지 에이전트는 근거와 함께 차단 레이블을 적용할 수 있습니다.

**실패 시 동작.** 게이트는 **실패 시 차단**합니다. 알림 API나 허용 목록을 읽을 수 없으면 주간 릴리스가 중단됩니다. 실패하면 태그를 전혀 생성하지 않고 목록을 링크와 함께 작업 요약에 기록하며, 저장소 변수 `WEEKLY_STATUS_ISSUE`가 설정되어 있으면 단일 "주간 릴리스 상태" 이슈의 댓글을 업데이트합니다.

**GitHub에서 소유자가 수행해야 하는 설정**(자동으로 생성되는 항목은 없음): `P0`, `release-blocker`, `regression` 레이블을 생성합니다. "주간 릴리스 상태" 이슈를 생성하고 `WEEKLY_STATUS_ISSUE`를 해당 이슈 번호로 설정합니다. `WAILS_REPO_TOKEN`에 Issues, Pull requests, Administration(협업자 권한), Dependabot alerts, Code scanning alerts, Secret scanning alerts, Repository security advisories 및 조직 Members에 대한 읽기 권한과 상태 댓글 작성을 위한 Issues 쓰기 권한을 부여합니다. 기본 `GITHUB_TOKEN`는 이러한 항목을 읽을 수 없습니다.

**재정의.** `override_reason`(workflow-dispatch 입력)은 게이트를 건너뜁니다. 사유, 실행자 및 건너뛴 모든 차단 항목은 요약과 상태 이슈에 기록됩니다. 릴리스 책임자만 이를 사용해야 합니다. GitHub에서는 입력 사용자를 한 사람으로 제한할 수 없으므로, 재정의에 명시적 승인이 필요하도록 필수 검토자가 지정된 환경에서 릴리스 작업을 실행하는 것이 좋습니다.

## 관련 항목

- [테스트 및 지속적 통합](/contributing/testing-ci/)
- [유지관리자 트리아지](/contributing/maintainer-triage/)
