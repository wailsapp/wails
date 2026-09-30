---
title: "MQ Studio"
description: "优先在本地运行的桌面客户端，支持 RocketMQ、RabbitMQ、Kafka 等系统"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![MQ Studio 新建连接对话框，显示 RocketMQ、Kafka 和 RabbitMQ 驱动](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** 是一款 **本地优先消息队列桌面客户端**，使用 **Go、Wails 和 React**构建。每个消息代理都有自己的控制台：RocketMQ 有一套，Kafka 有另一套，RabbitMQ 则提供管理插件。它们的界面和术语各不相同，而且每一套都是需要部署和维护的服务。MQ Studio 用一个应用取代它们：通过统一界面背后的驱动连接各个消息代理，无论连接哪种系统，页面和操作流程都保持一致。

## 主要亮点

- **一个界面，连接所有消息代理** — 目前支持 RocketMQ、RabbitMQ 和 Kafka，计划支持 Pulsar、NATS、MQTT 和 SQS。
- **主题、队列和消息** — 查看主题、队列、交换机和绑定；查询与追踪消息、实时查看日志、发送带有键和消息头的消息、重新发送消息以及处理死信。
- **消费者与积压** — 查看消费组、客户端、订阅和各分区的消费积压，支持重置偏移量以及处理重试和死信队列（DLQ）。
- **集群与告警** — 查看消息代理健康状态、运行指标、吞吐量和磁盘使用量，并接收原生桌面通知。
- **如实呈现连接能力** — 每个驱动都声明其端点实际支持的功能，界面只提供消息代理支持的操作。
- **默认保护隐私** — 配置保存在你的设备上，凭据在存储时加密。

无需部署服务器组件，无需维护网页控制台，也没有遥测。Wails 让这一切成为可能：这些消息代理的管理客户端都是 Go 库，因此驱动层可以在进程内直接调用它们，而前端仍是普通的 React 应用。每个平台只需一个二进制文件，不必再有人运维一项服务。

支持 macOS、Windows 和 Linux，提供英语和中文界面。

[网站](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[下载](https://github.com/amigoer/mq-studio/releases/latest)
