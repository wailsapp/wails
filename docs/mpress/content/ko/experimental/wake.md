---
title: "Wake"
description: "기존 Taskfile을 실행하면서 더 빠른 증분 빌드, 구조화된 출력, 기본 병렬 실행을 제공하는 실험적인 Wails 인식 빌드 러너입니다."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="실험적 기능"}
Wake는 `WAILS_USE_WAKE=true`을 통해 명시적으로 활성화해야 하며, 기본 러너가 **아닙니다**. 이 변수를 설정하지 않으면 `wails3 build / package / sign / task`은 이전과 완전히 동일하게 작동합니다. 지원 기능과 동작은 릴리스마다 변경될 수 있습니다.

@end

Wake는 `wails3`용 <strong>실험적 대체 빌드 러너</strong>입니다. 프로젝트에 이미 있는 것과 동일한 `Taskfile.yml`을 읽으며, task, dep, var, template, include 및 플랫폼 네임스페이스 구문도 그대로 사용합니다. 그런 다음 범용 [Task](https://taskfile.dev) 런타임 대신 Wails를 인식하는 실행기를 통해 이를 실행합니다.

목표는 Task를 대체하는 것이 아닙니다. 실제 Wails 프로젝트의 빌드 방식에 특화되고, 나머지 `wails3` CLI와 일치하는 의미 체계, 출력 및 기본값을 갖춘 러너를 제공하는 것이 목표입니다. **Wake만 사용하더라도 Taskfile을 변경할 필요는 없습니다.**

## Wake가 필요한 이유

Wake와 Task 런타임은 모두 `wails3`에 컴파일되어 있으므로 별도의 바이너리를 설치할 필요가 없습니다. 차이점은 Wake가 <strong>도메인을 이해한다는 것</strong>입니다. 범용 러너는 Taskfile에 나열된 단계를 지정된 순서대로 실행합니다. Wake는 Wails 빌드가 실제로 무엇<em>인지</em> 압니다. 즉, 프런트엔드 번들은 바이너리에 포함되고, 바이너리는 플랫폼별 아티팩트로 패키징되며, 아이콘과 바인딩도 함께 생성된다는 점을 이해합니다. Wake는 이 지식을 활용해 범용 러너로는 불가능한 방식으로 빌드를 최적화합니다.

- **실제 빌드에 필요한 작업만 수행합니다.** Wake는 각 단계의 실제 입력과 출력을 직접 추적합니다. Go 빌드에서는 모듈 그래프와 종속 단계의 출력이 이에 해당하므로, 관련 항목이 변경되지 않았다면 컴파일러와 링커를 다시 실행하지 않고 완전히 건너뜁니다. 범용 러너는 Taskfile에 감시할 파일이 정확히 미리 명시된 경우에만 단계를 건너뛸 수 있지만, Wake는 빌드에 관해 이미 알고 있는 정보를 바탕으로 이를 판단합니다. 아무 작업도 필요 없는 재빌드에는 대략 <strong>~20 ms(Wake)가 걸리는 반면 Task는 ~316 ms</strong>가 걸립니다. 콜드 빌드의 실제 소요 시간은 동일합니다. 콜드 빌드 시간은 `npm install`, Vite 및 Go 컴파일러가 대부분을 차지하기 때문입니다.

- **동시에 실행할 수 있는 작업을 압니다.** Wake는 서로 독립적인 단계를 이해하므로 기본적으로 병렬 실행하며, 최종 결과 행에는 그로 인해 얻은 속도 향상이 표시됩니다. 형제 단계의 출력이 뒤섞여 조사에 방해가 될 때는 `WAKE_SERIAL=true`을 사용해 병렬 실행을 비활성화하세요.

- **wails3에서 제어하는 구조화된 출력입니다.** Wake는 wails3 자체 리포터를 통해 출력을 표시합니다. 계획된 단계마다 한 행을 표시하고, 상태를 실시간으로 갱신하며, 마지막에 단계별 내역을 색상으로 구분해 보여 주고, 실패 패널 안에 클릭 가능한 `file:line` 링크를 제공합니다. `NO_COLOR` 및 비 TTY 환경(CI 로그)에서도 출력이 깔끔하게 단순화됩니다.

- **내장되어 있어 Wails와 함께 발전할 수 있습니다.** Wake는 타사 도구가 아니라 `wails3`의 일부이므로, 별도 프로젝트에서 구현하기를 기다리지 않고 새로운 빌드 기능을 직접 추가할 수 있습니다. 또한 현재 Taskfile이 `wails3` 바이너리를 셸에서 호출해 실행하는 크로스 플랫폼 스크립트와 도구를 네이티브로 실행할 가능성도 열립니다. 현재 방식은 호출할 때마다 프로세스를 생성합니다. 이 작업을 프로세스 내부로 통합하면 오버헤드가 줄어들며, 앞으로 더 많은 속도 향상을 기대할 수 있습니다.

## Wake 활성화

Wake는 전적으로 `WAILS_USE_WAKE=true` 환경 변수로 제어됩니다. 이 변수를 설정하지 않거나 `true` 이외의 값으로 설정하면 모든 `wails3` 명령은 이전과 완전히 동일하게 내장 Task 런타임을 사용합니다.

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

이 플래그는 `wails3 build`, `wails3 package`, `wails3 sign` 및 `wails3 task <name>`에 적용됩니다. `wails3 dev`은 아직 영향을 **받지 않습니다**. 개발용 감시기는 여전히 자체 파이프라인을 사용합니다.

@note{type="tip" title="Wake를 안전하게 활성화할 수 있습니다"}
Wake가 구현하지 않은 Taskfile 기능을 발견하면 전체 실행을 프로세스 내부의 내장 Task 런타임으로 넘깁니다. 외부 `task` 바이너리를 설치할 필요가 없습니다. 최악의 경우에도 이 플래그를 사용하지 않았을 때와 완전히 동일하게 작동합니다.

@end

## 계층형 로컬 재정의

Wake는 <strong>기본 Taskfile과 로컬 재정의</strong>를 지원합니다. `Taskfile.yml` 옆에 파일을 배치하면 그 파일의 정의가 우선 적용됩니다.

| 파일 | 용도 | 우선순위 |
| --- | --- | --- |
| `Taskfile.yml` | 커밋된 기본 파일 | 가장 낮음 |
| `Taskfile.override.yml` / `.yaml` | 팀 전체에 적용되는 커밋된 재정의 | 중간 |
| `Taskfile.local.yml` / `.yaml` | 개인용이며 일반적으로 Git에서 무시됨 | 가장 높음 |

**병합 의미 체계(로컬 정의 우선):**

- 이름이 **같은** task는 기본 task를 재정의합니다. 재정의에 목록 필드(`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`, `aliases`)가 있으면 기본값을 <strong>대체</strong>합니다. 재정의에서 생략한 필드는 기본 정의의 값을 유지합니다.
- `env`과 `vars`은 <strong>키별로 병합</strong>되며, 키가 충돌하면 재정의 값이 우선합니다.
- 재정의 파일에만 **존재하는** task는 <strong>추가</strong>됩니다.

예를 들어 커밋된 `Taskfile.yml`은 개발용 플래그로 빌드하지만 사용자의 컴퓨터에서는 항상 프로덕션 빌드를 해야 한다면 다음과 같이 설정합니다.

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

이제 `build`은 사용자의 프로덕션 명령을 실행하고 `smoke`도 사용할 수 있으며, 커밋된 Taskfile은 변경되지 않습니다.

@note{type="note" title="신뢰 모델"}
재정의 파일은 프롬프트 없이 자동으로 검색되어 적용됩니다. 그렇다고 새로운 기능이 허용되는 것은 아닙니다. Taskfile은 이미 임의의 셸 명령을 실행하므로, 재정의 파일은 `Taskfile.yml`을 편집해서 할 수 없는 작업을 수행할 수 없습니다. 커밋된 `Taskfile.override.*`은 PR diff에 표시되고, `Taskfile.local.*`는 사용자 컴퓨터에서 생성됩니다. 형식이 잘못된 재정의 파일은 조용히 건너뛰지 않고 실행을 중단시킵니다. 결정론적인 CI 빌드를 위해 재정의 검색을 완전히 건너뛰려면 `WAILS_NO_OVERRIDES=true`을 설정하세요.

@end

## 자동 대체 실행

Wake가 구현하지 않은 Taskfile 기능을 발견하면 전체 실행을 내장 Task 런타임으로 넘깁니다. 현재 다음 항목이 이 대체 실행을 트리거합니다.

- taskfile 수준의 `dotenv`
- `interleaved` 이외의 `output` 모드
- `requires` 블록
- `interval`(taskfile 또는 task 수준)
- `always` 이외의 `run` 모드
- task의 `short`
- 작업의 `defer`

## 환경 변수

| 변수 | 효과 |
| --- | --- |
| `WAILS_USE_WAKE` | 값이 `true`이면 라우팅 가능한 `wails3` 명령에 Wake를 활성화하며, 그 외의 값이면 Task 런타임을 사용합니다 |
| `WAILS_NO_OVERRIDES` | `true`는 `Taskfile.local.*` / `.override.*` 검색을 건너뜁니다(결정적 빌드) |
| `WAKE_VERBOSE` | 하위 프로세스의 stdout/stderr를 캡처해 실패할 때만 표시하는 대신 실시간으로 스트리밍합니다 |
| `WAKE_SILENT` | 작업 출력을 완전히 숨깁니다 |
| `WAKE_SERIAL` | `true`는 병렬 `deps:` 팬아웃을 비활성화합니다(기본값은 병렬 실행) |
| `WAKE_FORCE` | `true`는 완전히 새로 다시 빌드할 수 있도록 모든 캐시를 우회합니다 |
| `WAKE_DEBUG` | 리졸버 내부 동작(DAG, 종속성, 변수 참조, 실행 라우팅)을 기록합니다 |
| `WAKE_NOTICE` | 실행할 때마다 표시되는 "wake (experimental)" 알림을 숨기려면 `off`로 설정합니다 |

빌드 캐시는 `.wake/cache.json`에 저장됩니다(Task는 `.task/`을 사용합니다).

## 피드백

Wake는 실험적인 기능이며, 여러분의 피드백에 따라 향후 방향이 결정됩니다. Wake를 사용해 보셨다면 더 빨랐는지, 더 명확했는지, 그리고 문제가 발생하지는 않았는지 알려 주세요. 가장 유용한 보고에는 실행한 작업, 예상한 결과, 실제로 발생한 일이 포함됩니다. [Wake 피드백 토론](https://github.com/wailsapp/wails/discussions/5679)에서 의견을 남겨 주세요.
