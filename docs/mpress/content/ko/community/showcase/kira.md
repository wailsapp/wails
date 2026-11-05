---
title: "Kira"
description: "AWS의 ECS, RDS, S3, DynamoDB 등을 하나의 키보드 중심 인터페이스로 통합한 네이티브 macOS 데스크톱 클라이언트"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

<strong>[Kira](https://kira.thiennguyen.dev)</strong>는 <strong>AWS용 네이티브 macOS 데스크톱 클라이언트</strong>로, <em>"번거로움 없는 AWS"</em>를 지향합니다. <strong>Go, Wails 및 React + TypeScript</strong>로 구축되었으며, AWS 작업을 하나의 키보드 중심 앱으로 통합합니다. AWS SSO로 한 번 인증하면 콘솔과 CLI, 수많은 데이터베이스 도구 사이를 전환하지 않고도 여러 계정의 인프라를 관리할 수 있습니다.

처음에는 개인적인 불편함에서 시작되었습니다. 변경 사항 하나를 배포하기 위해 브라우저 탭과 `aws` 명령어, 별도의 SQL 클라이언트를 오가는 과정이 느리게 느껴졌습니다. Kira는 이러한 워크플로를 빠른 네이티브 창 하나로 통합합니다. 로그인하고 계정을 선택하면 필요한 모든 기능을 키 입력만으로 사용할 수 있습니다.

![프로덕션, 스테이징 및 개발 계정을 한곳에서 보여 주는 Kira의 다중 계정 개요](/assets/showcase-images/kira_screenshot_1.png)

## 주요 기능

- **다중 계정 AWS SSO** - 한 번 로그인한 후 한곳에서 계정과 리전을 전환할 수 있습니다
- **ECS** - 클러스터, 서비스 및 태스크를 탐색하고, 서비스를 재배포하거나 확장·축소하고 롤백하며, 태스크 정의를 확인하고, 서비스 메트릭을 모니터링하고, ECS Exec를 통해 대화형 셸을 열 수 있습니다
- **데이터베이스** - 안전한 SSH 터널링을 사용하고 자격 증명을 macOS 키체인에 저장하면서 RDS에서 SQL을 실행하고, DynamoDB를 쿼리 및 스캔하고, PostgreSQL, MySQL 및 Redshift에 연결할 수 있습니다
- **S3** - 버킷과 접두사를 탐색하고, 객체를 미리 보거나 업로드, 다운로드, 복사, 이름 변경 및 삭제하고, 폴더를 만들 수 있습니다
- **Secrets Manager** - 활성 계정의 보안 암호 목록을 확인하고 값을 가져올 수 있습니다
- **CloudWatch Logs** - 로그 스트림을 실시간으로 확인하고 검색할 수 있습니다
- **Smart Query** - `claude` CLI 기반의 선택적 AI 지원 SQL 생성 기능입니다
- **확장 기능** - 소규모 Go 스크립트로 작동하는 작업 버튼을 추가하는 사용자 지정 `.kext` 번들을 설치할 수 있습니다
- **빠른 탐색** - 전역 호출 단축키, `Cmd+K` 명령 팔레트 및 `kira://` 딥 링크를 제공합니다

## 자세히 살펴보기

목록을 벗어나지 않고 ECS 서비스의 태스크 상태, CPU와 메모리, 배포 상태를 실시간으로 확인하고 재배포, 확장·축소 또는 롤백할 수 있습니다.

![실시간 메트릭과 함께 ECS 서비스를 탐색하는 Kira](/assets/showcase-images/kira_screenshot_17.png)

파일 관리자처럼 S3를 탐색할 수 있습니다. 객체를 미리 보고 메타데이터와 버전을 확인하며, 현재 화면에서 바로 업로드, 다운로드, 이름 변경 또는 삭제할 수 있습니다.

![객체 미리 보기와 메타데이터를 제공하는 Kira의 S3 객체 브라우저](/assets/showcase-images/kira_screenshot_5.png)

macOS용으로 서명 및 공증된 `.dmg` 형태로 배포됩니다.

[Kira 방문하기](https://kira.thiennguyen.dev) | [문서 읽기](https://docs.kira.thiennguyen.dev)
