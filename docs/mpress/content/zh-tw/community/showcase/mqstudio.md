---
title: "MQ Studio"
description: "優先在本機執行的桌面用戶端，支援 RocketMQ、RabbitMQ、Kafka 等系統"
slug: "community/showcase/mqstudio"
sourcePath: "community/showcase/mqstudio.md"
---

![MQ Studio 新增連線對話方塊，顯示 RocketMQ、Kafka 和 RabbitMQ 驅動程式](/assets/showcase-images/mqstudio.webp)

**[MQ Studio](https://mq-studio.amigoer.com)** 是一款 **本機優先訊息佇列桌面用戶端**，以 **Go、Wails 和 React**打造。每個訊息代理都有自己的主控台：RocketMQ 有一套，Kafka 有另一套，RabbitMQ 則提供管理外掛。它們的介面與術語各不相同，而且每一套都是需要部署與維護的服務。MQ Studio 以單一應用程式取代它們：透過共用介面背後的驅動程式連接各個訊息代理，無論連線至哪種系統，頁面與操作流程都保持一致。

## 主要特色

- **一個介面，連接所有訊息代理** — 目前支援 RocketMQ、RabbitMQ 和 Kafka，並規劃支援 Pulsar、NATS、MQTT 和 SQS。
- **主題、佇列與訊息** — 檢視主題、佇列、交換器與繫結；查詢及追蹤訊息、即時查看日誌、傳送帶有鍵與標頭的訊息、重新傳送訊息，以及處理無法投遞的訊息。
- **消費者與積壓** — 檢視消費者群組、用戶端、訂閱及各分割區的消費積壓，支援重設位移量、重試及死信佇列（DLQ）處理。
- **叢集與警示** — 檢視訊息代理健康狀態、執行階段指標、輸送量與磁碟使用量，並接收原生桌面通知。
- **如實呈現連線能力** — 每個驅動程式都宣告其端點實際支援的功能，介面僅提供訊息代理支援的操作。
- **預設保護隱私** — 設定保留在你的裝置上，憑證以加密方式儲存。

無須部署伺服器元件，無須維護網頁主控台，也沒有遙測。Wails 讓這一切成為可能：這些訊息代理的管理用戶端都是 Go 程式庫，因此驅動層可以在同一處理程序內直接呼叫它們，而前端仍是一般的 React 應用程式。每個平台只需一個二進位檔案，不必另外維運服務。

支援 macOS、Windows 和 Linux，提供英文與中文介面。

[網站](https://mq-studio.amigoer.com) |
[GitHub](https://github.com/amigoer/mq-studio) |
[下載](https://github.com/amigoer/mq-studio/releases/latest)
