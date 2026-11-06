---
title: "Foundation"
description: "生产可用的 Wails v3 桌面应用脚手架"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Foundation 布局设置页，可切换九种布局预设](/assets/showcase-images/foundation-layouts.webp)

![Foundation 深色主题首页](/assets/showcase-images/foundation-home-dark.webp)

![Foundation 黑曜主题下的 X-Pro 图表](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** 是一个生产可用的 **Wails v3 桌面应用脚手架**，使用 **Go、React 19、TypeScript 和 Mantine 9** 构建。它为新项目准备好了通常要花几周才能搭好的部分——分层的 Go 后端、分平台的窗口外观、持久化、主题、国际化和端到端测试——让你可以直接进入业务开发。

## 主要亮点

- **各平台原生窗口** — Windows 和 Linux 使用无边框、由前端自绘的标题栏；macOS 使用沉浸式标题栏，并保留系统红绿灯按钮。
- **九种可切换布局** — 图标栏、展开侧栏、顶部导航、融合标题栏、右侧栏、底部 Dock、可拖拽的悬浮 Dock、浮动卡片和专注模式，均基于 CSS Grid 外壳；一行代码即可注册自定义预设。
- **主题与国际化** — 内置浅色、深色和黑曜预设，以及自定义调色板编辑器；开箱支持中文和英文；主窗口、子窗口和托盘面板的主题与语言实时同步。
- **SQLite 持久化** — GORM AutoMigrate，无需迁移文件；数据库位置可由用户选择，并提供逐表占用图表。
- **桌面集成** — 系统托盘面板、子窗口、原生文件与消息对话框，以及通过原生「另存为」对话框导出 CSV / XLSX。
- **可靠的工具层** — HTTP 客户端、绑定 Windows Job Object 的子进程、AES-GCM 加密、滚动日志和原子写入。
- **端到端测试** — Playwright 驱动 Wails 服务器模式构建（真实 Go 后端 + 真实前端），同一套测试还会生成 README 截图。

Wails 让这一切都打包在每个平台的单个二进制文件中：Go 服务通过生成的绑定直接调用，前端则遵循严格的 MVVM 结构，功能增加后依然易于维护。

支持 Windows、macOS 和 Linux。

[GitHub](https://github.com/FxRayHughes/foundation) |
[发布版本](https://github.com/FxRayHughes/foundation/releases/latest)
