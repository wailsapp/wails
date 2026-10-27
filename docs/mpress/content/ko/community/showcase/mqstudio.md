---
title: "MQ Studio"
description: "RocketMQ, RabbitMQ, Kafka 등을 위한 로컬 우선 데스크톱 클라이언트"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![RocketMQ, Kafka, RabbitMQ 드라이버가 표시된 MQ Studio의 새 연결 대화상자](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** 앱은 **로컬 우선 메시지 큐 데스크톱 클라이언트**이며, **Go, Wails, React**로 만들었습니다. 각 브로커에는 자체 콘솔이 있습니다. RocketMQ와 Kafka는 서로 다른 콘솔을 사용하고 RabbitMQ는 관리 플러그인을 제공합니다. 인터페이스와 용어가 다르고, 각각 배포하고 유지해야 하는 서비스이기도 합니다. MQ Studio는 이들을 하나의 애플리케이션으로 통합합니다. 공통 인터페이스 뒤의 드라이버를 통해 각 브로커에 연결하므로 어떤 시스템에 연결하든 화면과 작업 흐름이 같습니다.

## 주요 특징

- **모든 브로커를 하나의 인터페이스에서** — 현재 RocketMQ, RabbitMQ, Kafka를 지원하며 Pulsar, NATS, MQTT, SQS도 계획되어 있습니다.
- **토픽, 큐, 메시지** — 토픽, 큐, 익스체인지, 바인딩을 확인하고 메시지를 조회·추적하며 로그를 실시간으로 확인할 수 있습니다. 키와 헤더를 포함한 메시지 생성, 재전송, 배달하지 못한 메시지 처리도 가능합니다.
- **컨슈머와 지연** — 그룹, 클라이언트, 구독, 파티션별 지연을 확인하고 오프셋 재설정, 재시도 및 데드 레터 큐(DLQ) 처리를 수행합니다.
- **클러스터와 알림** — 브로커 상태, 런타임 지표, 처리량, 디스크 사용량 및 운영체제의 데스크톱 알림을 제공합니다.
- **연결 대상의 실제 기능 반영** — 각 드라이버가 엔드포인트의 실제 기능을 선언하며, 인터페이스는 브로커가 지원하는 작업만 제공합니다.
- **기본적인 개인정보 보호** — 설정은 기기에 저장되고 자격 증명은 저장 시 암호화됩니다.

배포할 서버 구성 요소도, 유지할 웹 콘솔도 없으며 텔레메트리도 없습니다. Wails 덕분에 가능합니다. 브로커 관리 클라이언트가 Go 라이브러리이므로 드라이버 계층이 같은 프로세스 안에서 직접 호출하고, 프런트엔드는 일반 React 앱으로 유지됩니다. 별도로 운영할 서비스 대신 플랫폼마다 하나의 바이너리를 제공합니다.

macOS, Windows, Linux에서 사용할 수 있으며 영어와 중국어를 지원합니다.

[웹사이트](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[다운로드](https://github.com/amigoer/mq-studio/releases/latest)
