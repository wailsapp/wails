---
title: "Kira"
description: "一款原生 macOS 桌面客户端，通过一个键盘驱动的界面统一管理 AWS 的 ECS、RDS、S3、DynamoDB 等服务"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** 是一款<strong>适用于 AWS 的原生 macOS 桌面客户端</strong>——*“轻松使用 AWS”*。它使用<strong>Go、Wails 和 React + TypeScript</strong>构建，将 AWS 操作整合到一个由键盘驱动的应用中。通过 AWS SSO 完成一次身份验证，即可管理多个账户中的基础设施，全程无需在控制台、CLI 和一堆数据库工具之间来回切换。

它源于一个切身痛点：仅仅为了发布一项更改，却要在浏览器标签页、`aws`命令和单独的 SQL 客户端之间来回切换，实在太慢。Kira 将这些工作流整合到一个快速的原生窗口中——登录、选择账户，此后需要的一切操作都只需按下快捷键。

![Kira 的多账户概览——在一处查看生产、预发布和开发账户](/assets/showcase-images/kira_screenshot_1.png)

## 主要亮点

- **多账户 AWS SSO** - 登录一次，即可在一处切换账户和区域
- **ECS** - 浏览集群、服务和任务；重新部署、扩缩容和回滚服务；查看任务定义；监控服务指标；并通过 ECS Exec 打开交互式 shell
- **数据库** - 对 RDS 运行 SQL，查询和扫描 DynamoDB，并连接 PostgreSQL、MySQL 和 Redshift——支持安全的 SSH 隧道，凭据存储在 macOS 钥匙串中
- **S3** - 浏览存储桶和前缀；预览、上传、下载、复制、重命名和删除对象；以及创建文件夹
- **Secrets Manager** - 列出当前账户中的机密并获取其值
- **CloudWatch Logs** - 实时跟踪和搜索日志流
- **智能查询** - 由`claude` CLI 提供支持的可选 AI 辅助 SQL 生成功能
- **扩展** - 安装自定义`.kext`包，以添加由小型 Go 脚本提供支持的操作按钮
- **快速导航** - 全局唤出快捷键、`Cmd+K`命令面板和`kira://`深层链接

## 深入了解

实时监控 ECS 服务——包括任务运行状况、CPU、内存和部署状态——并且无需离开列表即可重新部署、扩缩容或回滚。

![Kira 中显示实时指标的 ECS 服务浏览界面](/assets/showcase-images/kira_screenshot_17.png)

像使用文件管理器一样浏览 S3。预览对象，查看元数据和版本，并直接上传、下载、重命名或删除对象。

![Kira 的 S3 对象浏览器，支持对象预览和元数据查看](/assets/showcase-images/kira_screenshot_5.png)

以经过签名和公证的 macOS `.dmg`形式分发。

[访问 Kira](https://kira.thiennguyen.dev) | [阅读文档](https://docs.kira.thiennguyen.dev)
