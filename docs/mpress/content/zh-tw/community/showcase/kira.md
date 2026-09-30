---
title: "Kira"
description: "原生 macOS 桌面用戶端，將 AWS 的 ECS、RDS、S3、DynamoDB 等服務整合到單一鍵盤操作介面中"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** 是<strong>適用於 AWS 的原生 macOS 桌面用戶端</strong>，主打<em>「流暢無阻的 AWS 體驗」</em>。它使用<strong>Go、Wails 和 React + TypeScript</strong>建置，將 AWS 操作整合到單一、以鍵盤操作的應用程式中。只需透過 AWS SSO 驗證一次，即可管理多個帳戶的基礎架構，完全不必在主控台、CLI 和一堆資料庫工具之間來回切換。

它源自一個切身的不便：只是為了發布一項變更，卻得在瀏覽器分頁、`aws`命令和獨立的 SQL 用戶端之間來回切換，實在太慢。Kira 將這些工作流程整合到單一快速的原生視窗中：登入、選擇帳戶，所需的一切只需按下快捷鍵即可使用。

![Kira 的多帳戶概覽：在同一處查看正式、預備和開發帳戶](/assets/showcase-images/kira_screenshot_1.png)

## 主要特色

- **多帳戶 AWS SSO**：只需登入一次，即可在同一處切換帳戶與區域
- **ECS**：瀏覽叢集、服務和任務；重新部署、擴縮及回復服務；查看任務定義；監控服務指標；並透過 ECS Exec 開啟互動式 Shell
- **資料庫**：在 RDS 上執行 SQL、查詢及掃描 DynamoDB，並連線至 PostgreSQL、MySQL 和 Redshift；支援安全的 SSH 通道，且憑證儲存在 macOS 鑰匙圈中
- **S3**：瀏覽儲存貯體和前置詞；預覽、上傳、下載、複製、重新命名及刪除物件；並建立資料夾
- **Secrets Manager**：列出祕密並擷取作用中帳戶的祕密值
- **CloudWatch Logs**：即時追蹤及搜尋日誌串流
- **智慧查詢**：選用的 AI 輔助 SQL 生成功能，由`claude` CLI 提供支援
- **擴充功能**：安裝自訂`.kext`套件，加入由小型 Go 指令碼驅動的操作按鈕
- **快速導覽**：全域喚出快捷鍵、`Cmd+K`命令選擇區，以及`kira://`深層連結

## 深入瞭解

即時監看 ECS 服務的任務健康狀態、CPU、記憶體及部署狀態，並可直接在清單中重新部署、擴縮或回復。

![Kira 瀏覽 ECS 服務並顯示即時指標](/assets/showcase-images/kira_screenshot_17.png)

像使用檔案管理員一樣瀏覽 S3。預覽物件、檢視中繼資料和版本，並直接上傳、下載、重新命名或刪除。

![Kira 的 S3 物件瀏覽器，顯示物件預覽和中繼資料](/assets/showcase-images/kira_screenshot_5.png)

以適用於 macOS、經過簽署與公證的`.dmg`形式發布。

[造訪 Kira](https://kira.thiennguyen.dev) | [閱讀文件](https://docs.kira.thiennguyen.dev)
