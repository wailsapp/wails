---
title: "문서 수정하기"
description: "M-Press를 사용하여 Wails v3 문서의 수정 PR을 제출합니다."
sourcePath: "contributing/documentation.md"
---

수정 PR을 환영합니다. 오타, 깨진 링크, 오래된 예제, 불명확한 설명 또는 번역을 수정해 주세요. 문서만 수정하는 경우에는 이슈나 실패하는 코드 테스트가 필요하지 않습니다.

## 로컬에서 미리 보기

[wailsapp/wails](https://github.com/wailsapp/wails/fork)를 포크하고, 포크한 저장소를 클론한 다음, `master`에서 브랜치를 생성하세요.

고정된 버전의 문서 생성기를 설치하세요:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

[M-Press v1.0.17 릴리스](https://github.com/leaanthony/mpress/releases/tag/v1.0.17)에서도 검증된 바이너리를 다운로드할 수 있습니다.

Wails 저장소 루트에서 다음을 실행하세요:

```sh
mpress version
mpress dev
```

`docs/mpress/content/`에 있는 `.md` 소스 파일을 편집하세요. 영어가 기본 언어이며 해당 디렉터리에 직접 저장됩니다. 기존 번역은 `fr/` 및 `id/` 같은 언어 폴더에 저장됩니다. 파일을 저장하면 미리 보기가 다시 빌드됩니다.

각 페이지 상단의 메타데이터 블록과 쌍을 이루는 `@...` / `@end` 구성 요소를 그대로 유지하세요. 일반 단락, 제목, 목록 및 펜스로 구분된 코드는 텍스트로 편집할 수 있습니다. `docs/mpress/site/`에 있는 생성된 파일은 편집하지 마세요.

## 수정 사항 확인하기

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

변경한 페이지를 브라우저에서 확인하고, 수정한 코드 예제를 모두 실행하세요. 번역을 수정한 경우에는 변경한 구절 전체를 영어 원문과 비교하세요. 명령어, API 이름, 링크, 코드 샘플 및 다이어그램 연결 관계를 그대로 유지하세요. 자연스러운 기술 용어를 사용하고 요구 사항과 주의 사항을 보존하며, 본문뿐 아니라 표시되는 다이어그램 레이블, 탐색 레이블 및 이미지 설명도 번역하세요.

게시되는 모든 언어에는 각 영어 페이지의 완전한 번역이 있어야 합니다. 영어 자리표시자나 대체 페이지를 사용하지 마세요. 특정 번역 하나를 수정할 때는 해당 언어만 변경해도 됩니다. 영어 원문의 의미를 변경하는 경우에는 게시되는 다른 언어의 해당 페이지도 업데이트하세요. 관련 없는 페이지는 다시 생성할 필요가 없습니다.

## 풀 리퀘스트 제출하기

`master`을 대상으로 PR을 생성하세요. 무엇이 잘못되었는지 기술하고, 수정 내용을 설명하며, 실행한 검사 항목을 나열하세요. 눈에 보이는 레이아웃 변경에는 스크린샷을 포함하고, 변경한 코드 예제에는 플랫폼 및 버전 세부 정보를 포함하세요.

Cloudflare 자격 증명과 비공개 서비스는 필요하지 않습니다. 공개 PR 검사는 배포 자격 증명 없이 정적 사이트를 빌드하고 검증합니다.

코드 변경 및 기능 제안에 대해서는 [Wails에 기여하기](/contributing/)를 참조하세요. 내부 구조에 대해서는 [기술 개요](/contributing/overview/)를 참조하세요.
