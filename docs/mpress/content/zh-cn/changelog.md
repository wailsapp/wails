---
title: "更新日志"
description: "Wails v3 的版本历史和发行说明"
slug: "changelog"
sourcePath: "changelog.md"
---

图例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 本项目的所有重要变更都将记录在此文件中。

本文件的格式基于[维护更新日志](https://keepachangelog.com/en/1.0.0/)， 本项目遵循[语义化版本](https://semver.org/spec/v2.0.0.html)规范。

- `Added`表示新增功能。
- `Changed`表示现有功能的变更。
- `Deprecated`表示即将移除的功能。
- `Removed`表示现已移除的功能。
- `Fixed`表示各类错误修复。
- `Security`表示安全漏洞。

_/

/_   * 请勿更新此文件 *   更新应添加到`v3/UNRELEASED_CHANGELOG.md`   谢谢！ _/

## [尚未发布]

## v3.0.0-beta.21 - 2026-09-13

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/6116)中使用 M-Press 提供 Wails v3文档，由 @leaanthony 贡献

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/6118)中解析 MPD 头部元数据中的 JSON slug 值以生成变更日志，由 @leaanthony 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6080)中修复更新程序：备份失败后清除辅助程序环境变量并重新启动原始目标，由 @cnmax 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6098)中于 App.Run 期间启动默认信号处理程序，由 @leaanthony 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6112)中修复 Windows 菜单：处理 nil 菜单、释放被替换的资源并重绘菜单栏，由 @taliesin-ai 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6115)中恢复使用共享 YAML 配置的新建项目的 MSIX 打包功能，由 @leaanthony 贡献
- 修复以下问题：当泛型模型创建器引用后续声明的辅助函数时，生成的 JavaScript 和 TypeScript 绑定无法加载；并防止创建相互依赖的泛型模型时发生堆栈溢出（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 变更

- 由 @01xR4in 在[PR](https://github.com/wailsapp/wails/pull/6082)中将 Clave 展示项目的链接更新为当前网站和代码仓库的链接

## 修复

- 取消已中止的 Windows 资源请求（包括 Worker 请求），同时在导航期间保留 keepalive 处理程序。在 Apple 平台上，通过应用程序包装器转发原生请求上下文。(#5963、#5969)
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送之间保留更新日志条目
- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中移除引用模块中从未随附的二进制文件的嵌入项，修复`go mod vendor`在所有平台上均因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，并修复[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)

## 移除

- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中移除已被纯 Go 加载器取代的原生 WebView2 加载器支持。此项变更会移除内嵌的`WebView2Loader.dll`二进制文件和`github.com/jchv/go-winloader`依赖项。`native_webview2loader`构建标签仍会被接受且不再报错，但对 v3 构建不起作用
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6097)中从 macOS API 指南中移除未使用的构建标签和 FPS 选项

## v3.0.0-beta.19 - 2026-09-09

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6087)中使用构建标签限制对 macOS 私有 API 的访问，使其仅在主动选择启用时可用——请参阅[文档](https://v3.wails.io/features/browser/integration)、[文档](https://v3.wails.io/features/environment/info)、[文档](https://v3.wails.io/features/windows/basics)、[文档](https://v3.wails.io/features/windows/frameless)、[文档](https://v3.wails.io/features/windows/notch-windows)、[文档](https://v3.wails.io/features/windows/options)、[文档](https://v3.wails.io/guides/build/macos)、[文档](https://v3.wails.io/guides/build/private-macos-apis)以及[文档](https://v3.wails.io/reference/overview)

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6091)中拒绝超过64 MiB 的运行时请求，并返回 HTTP 413

## 安全

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6092)中使用令牌身份验证加强 MCP 来源和远程访问的安全性

## v3.0.0-beta.18 - 2026-09-08

## 修复

- 由 @4RH1T3CT0R7 在[PR](https://github.com/wailsapp/wails/pull/6083)中使用指针接收器，修复 Linux 和 Darwin 上的 Calloc 内存泄漏

## v3.0.0-beta.17 - 2026-09-06

## 修复

- Windows：WebResourceRequested 处理程序中失败或为 nil 的`GetRequest`不再导致进程终止（`log.Fatal`/nil 指针解引用 panic）；现在会丢弃该请求并记录日志。此项修复由 @midagedev 在[PR](https://github.com/wailsapp/wails/pull/6006)中完成

## v3.0.0-beta.16 - 2026-08-29

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6029)中改为在新的终端窗口中提示输入公证密码

## 修复

- 由 @ChewbaccaCookie 在[PR](https://github.com/wailsapp/wails/pull/5919)中正确处理 macOS 上的系统托盘点击类型
- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6041)中让 CI 在更新前移除未使用的 Microsoft apt 软件源

## v3.0.0-beta.15 - 2026-08-27

## 修复

- 由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6043)中将 WebView2 嵌入超时时间增加到60秒

## v3.0.0-beta.14 - 2026-08-26

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6032)中正确命名 macOS 上的 Control-字母组合键
- 由 @nik9play 在[PR](https://github.com/wailsapp/wails/pull/6016)中修复 Windows 上的 ICO 托盘图标，并使其跟随任务栏主题

## v3.0.0-beta.13 - 2026-08-25

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6026)中确保 macOS 运行模态循环时仍继续处理主线程工作
- 由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5923)中让移动端安全存储操作可以失败，并在失败时默认拒绝访问
- 由 @archy-rock3t-cloud 在[PR](https://github.com/wailsapp/wails/pull/5999)中确保即使未注册监听器，也会运行应用程序事件钩子
- 由 @haoku123 在[PR](https://github.com/wailsapp/wails/pull/6023)中更正注释和本地化文档中的拼写错误
- 由 @4RH1T3CT0R7 在[PR](https://github.com/wailsapp/wails/pull/6025)中移除提交在`v3/examples`下的预编译 macOS 二进制文件

## v3.0.0-beta.12 - 2026-08-21

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6010)中新增带生命周期和遥测示例的 macOS 刘海通知窗口——请参阅[文档](https://v3.wails.io/features/windows/notch-windows)
- 新增 macOS NSPanel 窗口支持，包括新选项和原生集成——参见 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/6008)中提供的[文档](https://v3.wails.io/features/windows/options)

## 修复

- 防止并发调用时 SQLite Prepare 挂起，见 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 移除

- 从文档中移除过时的实现跟踪器，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 修复

- 修复 GTK4 Linux 宿主丢弃通过自定义协议和文件关联启动时所传参数的问题，见 @midagedev 的[PR](https://github.com/wailsapp/wails/pull/6000)
- 在变更日志验证中正确处理已删除的行和同一来源的修正，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 新增

- 新增安全的 wails3 mcp 服务器，用于智能体辅助的项目管理，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5896)
- 新增绑定中模型的文档——参见 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5988)中提供的[文档](https://v3.wails.io/features/bindings/models)
- 支持在原子化 Linux 系统上通过 rpm-ostree 安装，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5987)
- 新增原生的每周 Star 历史图表生成和发布功能——参见[文档](https://v3.wails.io/credits)、[文档](https://v3.wails.io/de/credits)、[文档](https://v3.wails.io/fr/credits)、[文档](https://v3.wails.io/id/credits)、[文档](https://v3.wails.io/ja/credits)、[文档](https://v3.wails.io/ko/credits)、[文档](https://v3.wails.io/pt/credits)、[文档](https://v3.wails.io/ru/credits)、[文档](https://v3.wails.io/zh-cn/credits)和[文档](https://v3.wails.io/zh-tw/credits)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5986)
- 新增仅适用于 Darwin 的 mac 包，用于解析应用程序包资源——参见 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5965)中提供的[文档](https://v3.wails.io/guides/build/macos)
- 新增 Condui 展示页面和索引条目——参见[文档](https://v3.wails.io/community/showcase/condui)和[文档](https://v3.wails.io/community/showcase)，见 @mgueregath 的[PR](https://github.com/wailsapp/wails/pull/5962)
- 新增 Redis Viewer 展示页面，包括屏幕截图和项目链接——参见[文档](https://v3.wails.io/community/showcase)和[文档](https://v3.wails.io/community/showcase/redisviewer)，见 @redisviewer 的[PR](https://github.com/wailsapp/wails/pull/5984)

## 变更

- 在 Linux 上将 GTK 应用程序标志更新为 G<em>APPLICATION</em>NON_UNIQUE，见 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5971)
- 将缺失窗口事件的日志级别由警告改为调试，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5914)

## 修复

- 允许已注册的 macOS 键盘快捷键优先于 WebView 处理，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5902)
- 修复文档侧边栏中失效的链接，见 @northes 的[PR](https://github.com/wailsapp/wails/pull/5937)
- 当 WebKit 中止对应的自定义 URL 协议任务时，取消 macOS 和 iOS 的资源请求上下文（#5963）
- 处理 WindowSetFullscreenButtonEnabled 消息，见 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5976)
- 在 preact-ts 模板中导入 Fragment，以解决构建失败问题，见 @haoku123 的[PR](https://github.com/wailsapp/wails/pull/5979)
- 防止旧版 GTK3 纯服务应用程序在尚无可用活动窗口或显示器时执行屏幕发现而崩溃（#5966）
- 当未发布的变更日志为空时，允许明确指定版本的发布任务继续运行（#5977）

## 安全

- 将网站 nanoid 锁定文件更新至已修补的3.3.18版本，以解决安全公告所述问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 新增

- 为自动生成的变更日志条目新增文档 URL 生成功能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5957)
- 新增 Streams：使用 WebSocket 编程模型在 Go 与 JavaScript 之间传输双向字节流，且无需监听套接字。在 Go 中使用`app.HandleStream(name, handler)`声明流，并在前端使用`Stream(name)`连接；后者返回一个结构与`WebSocket`相同的对象。Go→JS 通过资源服务器为每个窗口保持一个轮询请求来传输，JS→Go 则通过普通 POST 传输；不会绑定任何 TCP 端口，也不会有任何内容经过`evaluateJavaScript`。在服务器构建（`-tags server`）中，同一处理程序改为通过真正的 WebSocket 提供服务，因此各构建使用相同的应用程序代码。作者：@leaanthony
- 将 mailbox 变更日志条目移至“未发布”，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5935)

## 变更

- 更新文档侧边栏的自动生成方式和博客作者类型的派生方式，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5938)

## 修复

- WebView2 初始化使用截止时间和消息泵，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5952)
- WebView2 Cookie 测试在 CI 中默认跳过，仅在明确启用时运行，并将执行锁定到当前操作系统线程，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5951)
- Windows 菜单构建器恢复子菜单父项的命令 ID，见 @gilad-ch 的[PR](https://github.com/wailsapp/wails/pull/5944)
- 使官方交叉编译镜像与 GTK 4.14+ Linux 支持基线保持一致（#5928）
- 配置 iOS Xcode 项目以保留继承的链接器标志，并添加 -ObjC，见 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5915)
- 修复大型前端中`wails3 dev`资源代理频繁新建和断开 TCP 连接的问题；该问题可能耗尽宿主机的临时端口，并导致无关进程因`EADDRNOTAVAIL`而失败
- 按窗口将事件 JavaScript 加入队列，以保证有序分发并实现背压，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5934)

## 移除

- 移除桌面二进制文件发布流水线：v3 版本仅发布标签，`wails3` CLI 使用`go install`安装。删除`release-v3.yml`以及每夜构建中触发该流程的步骤，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5946)

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 添加 macOS 自动播放偏好设置，以取消媒体播放必须由用户操作触发的要求，由 @Eyalm321 在[PR](https://github.com/wailsapp/wails/pull/5512)中实现
- 将邮箱的变更日志条目移至“未发布”部分，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5935)中完成

## 已更改

- macOS 缩放动画改用 CADisplayLink 或 NSTimer，以获得更流畅的性能，由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5945)中实现

## 已修复

- 配置 iOS Xcode 项目以保留继承的链接器标志，并添加 -ObjC，由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5915)中实现
- 修复大型前端中`wails3 dev`资源代理频繁创建和销毁 TCP 连接的问题；该问题可能耗尽主机的临时端口，并导致无关进程因`EADDRNOTAVAIL`而失败
- 按窗口将事件 JavaScript 加入队列，以实现有序分派和背压，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5934)中实现

### 已添加

- 实现通用异步 FIFO 邮箱，以有序传递事件，由 @savely-krasovsky 和 @DevLumuz 在[PR](https://github.com/wailsapp/wails/pull/5851)中实现

## v3.0.0-beta.6 - 2026-08-09

## 已添加

- 为超大事件实现有界的主机端存储和有序的 JavaScript 传递，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5930)中实现
- 实现通过 macOS Dock 图标弹跳来提示窗口，由 @julianstorer 在[PR](https://github.com/wailsapp/wails/pull/5921)中实现

## 已修复

- 资源服务器在刷新期间会保留内容类型嗅探错误和尚未写出的前缀，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5931)中实现
- 防止从 Wails 回调替换应用程序菜单时导致 macOS 应用程序崩溃
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨认的问题。深色模式 uxtheme 导出项此前仅对内部版本18334及以上启用，因此应用级深色模式选择加入逻辑从未在这些主机上运行：菜单背景被绘制为深色，但 Windows 仍按浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号导出项从17763起便已存在，因此现已相应调整版本门槛。
- 修复`w32.GetStockObject`调用`GetDeviceCaps`而非`GetStockObject`的问题；该问题导致它对所有预定义对象都返回0。
- 改进 WebView2 引导程序的下载错误处理和报告，由 @jannskiee 在[PR](https://github.com/wailsapp/wails/pull/5924)中实现

## v3.0.0-beta.5 - 2026-08-07

## 已修复

- macOS 应用激活现在仅对常规应用遵循激活策略，由 @julianstorer 在[PR](https://github.com/wailsapp/wails/pull/5897)中实现
- 在 Linux 构建中防范未初始化的 GTK 窗口，由 @julianstorer 在[PR](https://github.com/wailsapp/wails/pull/5898)中实现
- 在加载 URL 前，为 Linux WebKit 窗口设置明确的不透明背景色，由 @julianstorer 在[PR](https://github.com/wailsapp/wails/pull/5899)中实现

## v3.0.0-beta.4 - 2026-08-05

## 已更改

- Android 构建任务默认使用 arm64，deploy-emulator 会选择主机架构，由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5890)中实现

## 已修复

- 在拖动期间保留 macOS 窗口的缩放状态并减少动画，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5900)中实现
- 通过向`webview_window_windows_nonclient.go`构建约束添加`!server`，修复 Windows 服务器模式构建

## v3.0.0-beta.3 - 2026-08-03

## 已添加

- 在实现细节中记录第10阶段 Beta 验证的完成情况，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5881)中完成

## 已修复

- 将窗口句柄传递给 Windows 深色模式 API，并验证参数，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5877)中实现
- 集中处理 macOS 无边框窗口的标题栏按钮状态解析，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5870)中实现
- 修复 Windows 应用程序请求深色模式而 Windows 应用主题为浅色时，原生菜单文本无法辨认的问题。在 Windows 能够渲染深色菜单文本之前，菜单现在会使用与之匹配的浅色原生背景。
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨认的问题。深色模式 uxtheme 导出项此前仅对内部版本18334及以上启用，因此应用级深色模式选择加入逻辑从未在这些主机上运行：菜单背景被绘制为深色，但 Windows 仍按浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号导出项从17763起便已存在，因此现已相应调整版本门槛。

## v3.0.0-beta.2 - 2026-08-02

## 已更改

- 将 v3 从 Alpha 阶段提升至 Beta 阶段
- 记录系统托盘的智能默认设置和弹出菜单自动隐藏行为，并为点击处理程序选择添加回归测试覆盖（#5840）。
- GitHub 更新程序默认排除 Windows 安装程序资源，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5861)中实现
- 为 macOS 无边框窗口添加圆角、直角和自定义圆角半径支持，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5866)中实现

## 已修复

- 报告 GTK4 窗口的实时尺寸，并从已配置的表面发出调整大小、最大化、最小化和全屏状态事件（#5830）。
- 修复在 fetch 请求中发送 Blob 或 FormData 时导致 Linux WebKit 崩溃的问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5854)中实现
- 对于缺失的 Blob/FormData 标头，fetch shim 会传递 undefined，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5865)中实现

## v3.0.0-alpha2.122 - 2026-08-01

## 已添加

## 已更改

- 为 macOS 无边框窗口添加圆角、直角和自定义圆角半径支持，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5866)中实现

## 已修复

- 对于缺失的 Blob/FormData 标头，fetch shim 会传递 undefined，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5865)中实现

## v3.0.0-alpha2.121 - 2026-07-31

## 已添加

- 添加 macOS DMG 打包支持，包括新选项和构建任务，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5857)

## 变更

- GitHub 更新程序现在默认排除 Windows 安装程序资产，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5861)

## 修复

- 修复在 fetch 请求中发送 Blob 或 FormData 时 Linux WebKit 崩溃的问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5854)

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 实现 macOS 标题栏双击时最大化或最小化窗口的操作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5853)

## 变更

- 更新二维码服务教程以使用 NewServiceWithOptions，并添加间距，见 @jeongkyu 的[PR](https://github.com/wailsapp/wails/pull/5849)

## 修复

- 使 WKWebView 在 macOS 缩放期间保持响应，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5856)
- 修复 GTK4 窗口尺寸查询，并从已配置的`GdkSurface`发出调整大小、最大化、最小化和全屏状态事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修复

- 更新多语言文档以加入架构图，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5833)

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 为图标生成的输入和输出提供默认路径，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5825)
- 将源入口模块添加到运行时 package.json 的 sideEffects 中，见 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5797)
- 在贡献指南中添加许可证和来源说明章节，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5816)

## 修复

- 应用限定作用域的 GTK4 无边框 CSS 以移除圆角，见 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5800)
- 在 Windows 上妥善处理弹出菜单和屏幕枚举过程中的光标位置获取失败，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5789)
- macOS 打开文件对话框现在能正确按扩展名筛选，并根据文件名后缀验证允许的文件，见 @phergul 的[PR](https://github.com/wailsapp/wails/pull/5678)
- 防止 Windows 深色模式初始化时调用 nil API，见 @roachadam 的[PR](https://github.com/wailsapp/wails/pull/5793)
- 修复更新程序的32位构建失败问题：在`GOARCH=386`上将`maxArchiveTotalSize`常量（2 GiB）传递给`fmt.Errorf`时，会溢出平台的`int`。该常量现在已显式指定为`int64`类型。
- 修复启动时的 nil 指针 panic：在未加载深色模式 uxtheme API 的 Windows 构建上，如果窗口使用深色（或跟随系统深色）标题栏，就会出现此问题，例如 Windows 10 1809 / Windows Server 2019（内部版本17763）。窗口主题设置中的`AllowDarkModeForWindow`调用现在加入了 nil 检查，与`w32.SetMenuTheme`中已有的检查保持一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 为 Windows 上的非客户端区域实现自定义命中测试逻辑，见 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5462)

## 变更

- 根据 UseVisualHosting 配置 WebView2 的显示器缩放检测，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5761)

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 更新常见问题文档，使其重点介绍 Wails v3 的功能和指南，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5763)

## v3.0.0-alpha2.115 - 2026-07-06

## 修复

- 修复`Menu.Update()`未在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对数量创建快照，修复显示器变更时枚举 macOS 屏幕导致的崩溃（#5565；@x-haose 在 #5584 中独立诊断并修复）
- 修复 Windows 上的崩溃：在最小化/还原转换期间，当`GetClientRect`返回 nil 时，`WM_ERASEBKGND`绘制纯色背景会导致崩溃（防护措施由 @sinspired 在 #5636 中报告）
- 修复前端绑定错误始终按文本解析的问题，由 @mbaklor 在 #5690 中修复
- 修复使用`server`构建标签时 Windows 构建失败的问题；原因是 Windows GUI 文件缺少 macOS 和 Linux 对应文件中已有的`!server`构建约束（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 实现更新清单协议和端点提供程序，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5720)

## 变更

- 将`webview2`绑定并入 v3 模块并命名为`v3/internal/webview2`，移除独立模块、其 nightly 发布/同步工作流以及 go.mod 版本协调流程（v3 是其唯一使用方），见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5711)

## 修复

- 将 WebView2 显示器缩放检测和 DPI 变更时的宿主重新同步修复移至“未发布”章节，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5750)
- 更新 WebView2 对 float64 和 BOOL 参数的 COM 编组处理，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5741)
- 防止在 Windows 系统托盘图标更新和销毁期间发生 panic 和 nil 解引用，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5703)
- 修复隐藏窗口在 Windows 上无法正确再次隐藏的问题，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5743)
- 在窗口最小化、最大化和还原时同步 WebView2 控制器的可见性，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5742)

### 修复

- 重新启用 WebView2 显示器缩放检测，并仅在 DPI 变更时执行宿主重新同步，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5734)；该更改基于 @randalmurphal 验证的修复，并由 @eleclin 验证根本原因、@qq540491950 完成硬件测试

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 在未设置`ANDROID_KEYSTORE_FILE`的情况下构建发布版 AAB 时添加警告（Google Play 会拒绝使用调试签名的应用包），并在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5730)中记录 App Bundle 的打包和签名方法
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5739)中新增多语言 Why Wails 文档
- 在由 @fbbdev 提交的[PR](https://github.com/wailsapp/wails/pull/5398)中支持在绑定中将 Go time.Time 映射为 JS Date 或字符串
- 在由 @mortenolsrud 提交的[PR](https://github.com/wailsapp/wails/pull/5728)中新增用于提交到 Play Store 的 Android App Bundle (AAB) 打包任务（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 任务仍用于本地/模拟器测试（修复[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5735)中新增 Android 实体设备任务目标，并恢复相机/位置权限

## 变更

- 将`webview2`升级到 v1.0.28（[发行说明](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5730)中，将 Android 模板的`compileSdk`/`targetSdk`从34升级到35；这是 Google Play 对新应用提交的要求

## 修复

- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/5745)中修复 sponsorkit 中已烘焙的头像蒙版
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5730)中修复 Android AVD 自动创建功能因按字典序排序版本而选错系统映像或 cmdline-tools 版本的问题
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5730)中修复设置向导建议使用过时 Android NDK 版本的问题（现为26.3.11579264，与文档要求一致）
- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/5744)中更新有关 SvelteKit 和选项的法语文档
- 在由 @flofreud 提交的[PR](https://github.com/wailsapp/wails/pull/5516)中修复显示器变化期间枚举 macOS 屏幕时发生的 SIGSEGV

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5724)中新增基于 Go 的贡献者 SVG 生成器，并更新文档和网站的致谢页面
- 在由 @fbbdev 提交的[PR](https://github.com/wailsapp/wails/pull/5398)中支持在绑定中将 Go time.Time 映射为 JS Date 或字符串

## 变更

- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5719)中使用 Go 生成器替换基于 Node 的赞助者图片处理流水线

## 修复

- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5729)中修复 Android 构建资产依赖项安装脚本
- 在`ValidateAndSanitizeURL`中拒绝 U+0085 (NEXT LINE) 控制字符，补全 URL 验证器对空白字符的覆盖
- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/4785)中，使无边框窗口在 DPI 变化时重新计算 DWM 窗口边框
- 在由 @yulesxoxo 提交的[PR](https://github.com/wailsapp/wails/pull/4632)中修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5714)中，为 Darwin 的对话框、菜单、托盘和通知所涉及的 Cocoa 对象添加显式 Objective-C 内存管理
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5718)中修复 Linux CGO 后端错误和系统托盘问题

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 在由 @Aliuyanfeng 提交的[PR](https://github.com/wailsapp/wails/pull/5061)中将 HappyTools 添加到社区展示
- 在由 @triadmoko 提交的[PR](https://github.com/wailsapp/wails/pull/5643)中新增印度尼西亚语区域设置支持和完整文档
- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/4813)中为 WindowsWindow 新增 DisableMenu 选项

## 变更

- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/5617)中更新 Taskfile 模板和 CLI，以使用 GOOS 和 ARCH 分派构建/打包任务

## 修复

- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5708)中修复 Mac 窗口标签页功能的问题

## 移除

- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5702)中从贡献、功能和指南文档中移除德语翻译的 MDX 文件

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 在由 @wayneforrest 提交的[PR](https://github.com/wailsapp/wails/pull/5129)中实现 macOS WebView 的重新加载和强制重新加载，并新增 WebContent 进程终止后的恢复机制
- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/5396)中新增涵盖贡献、功能和指南的完整德语文档
- 在由 @popaprozac 提交的[PR](https://github.com/wailsapp/wails/pull/5333)中增强通知功能，支持声音、附件、定时通知及通知更新 API

## 修复

- 在由 @leaanthony 提交的[PR](https://github.com/wailsapp/wails/pull/4785)中，使无边框窗口在 DPI 变化时重新计算 DWM 窗口边框
- 在由 @yulesxoxo 提交的[PR](https://github.com/wailsapp/wails/pull/4632)中修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 在由 @iamhabbeboy 提交的[PR](https://github.com/wailsapp/wails/pull/5026)中为 EventsEmit 文档新增代码示例
- 在由 @MerIijn 提交的[PR](https://github.com/wailsapp/wails/pull/5380)中新增 Windows WebView2 可视化托管选项
- 在由 @SametKUM 提交的[PR](https://github.com/wailsapp/wails/pull/5536)中将 Klustr 添加到社区展示文档
- 在由 @thiennguyen93 提交的[PR](https://github.com/wailsapp/wails/pull/5685)中将 Kira 添加到社区展示，并新增页面和变更日志条目
- 在由 @taliesin-ai 提交的[PR](https://github.com/wailsapp/wails/pull/5694)中为 MCP 服务指南新增反馈章节

## 变更

- 服务器模式现已提供与桌面端构建任务一致的一等生产构建支持（#5693）。`task build:server`默认构建生产二进制文件（`-tags server,production`、`-trimpath`、已剥离符号），并接受`DEV=true`（开发服务器）、`OBFUSCATED=true`（garble）和`EXTRA_TAGS`。`task run:server`运行开发服务器。`Dockerfile.server`/`task build:docker`先构建生产服务器（`-tags server,production`）和生产前端；镜像默认采用基于 distroless/static 的纯 Go 静态构建，并将`CGO_ENABLED`、`GO_IMAGE`和`RUNTIME_IMAGE`公开为可覆盖的构建参数，以支持 CGO 应用。

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/4435)中防止关闭仍有待处理异步调用的窗口时发生崩溃，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5249)中防止在 Windows 上打开隐藏应用时激活窗口，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5668)中确保 WebKit 请求元数据、响应完成处理和正文流处理均在 GTK 主线程上运行，作者：@taliesin-ai
- 修复`Menu.Update()`无法在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并创建数量快照，修复显示配置变化时枚举 macOS 屏幕导致的崩溃（#5565；@x-haose 在 #5584 中独立诊断并修复）
- 通过在`WM_DPICHANGED`处理程序中重新设置控制器边界（与取消最小化时的 DPI 重新同步方式一致），修复在 Windows 上将窗口拖过不同 DPI 的显示器后 WebView2 内容先缩小再消失的问题（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 通过`app.GlobalShortcut`新增全局（系统级）键盘快捷键（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。即使应用未获得焦点，快捷键也会触发。各平台均采用原生实现，不依赖第三方库：macOS 使用 Carbon 热键，Windows 使用`RegisterHotKey`，X11 使用`XGrabKey`，Wayland 使用 XDG Desktop Portal 全局快捷键接口。
- 新增内置 MCP 服务器：这是一个 Model Context Protocol 服务器，当应用使用`mcp`标签构建时会自动启动，使 LLM 智能体能够测试和控制正在运行的 Wails 应用，包括窗口控制、DOM 检查、JavaScript 求值、绑定方法调用、事件，以及使用屏幕动画光标呈现的模拟鼠标/键盘输入。无需编写用户代码：设置`WAILS_MCP=1`后，`wails3 build`/`wails3 dev`会自动添加`mcp`标签。完全通过环境变量（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）配置。

## 修复

- 修复`Menu.Update()`无法在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并创建数量快照，修复显示配置变化时枚举 macOS 屏幕导致的崩溃（#5565；@x-haose 在 #5584 中独立诊断并修复）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5613)中新增实验性 Wake 文档及侧边栏导航，作者：@leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## 变更

- 将`webview2`升级到 v1.0.27。
  - ci(webview2)：修复发布构建（交叉编译 Windows + 补全 go.sum）（#5671）\

  **完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- 在[PR](https://github.com/wailsapp/wails/pull/5672)中从 webview2 发布工作流的交叉编译步骤中移除 go vet，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5670)中将自动生成变更日志所用的 OpenRouter 模型更新为 google/gemini-2.5-flash-lite，作者：@taliesin-ai
- 将`webview2`升级到 v1.0.26。

### 修复

- **从暂时性运行时 COM 错误中恢复，而非退出**（#5658、#5580）。此前，`Chromium.errorCallback`会对<em>任何</em> COM 错误调用`os.Exit(1)`，因此启动后一次可恢复的短暂故障就会终止整个应用。运行时路径（`Resize`/`GetClientRect`、`Navigate`/`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）现在会记录错误并恢复。特别是，`MessageReceived`中格式错误或不可信的 Web 消息现在会被丢弃，而不会导致进程终止。这解决了跨越不同 DPI 显示器时发生的此类崩溃（#5544、#5650）。环境/控制器创建路径中发生的错误仍按致命错误处理。\

**完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5671)中修复 release-webview2 工作流，使其能够正确处理 go.sum 文件，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5659)中通过清除并重新构建原生菜单，修复 Linux GTK4 菜单更新问题，作者：@taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 新增`application.System`，用于从共享代码中检测运行时平台：`System.IsMobile()`（iOS/Android）、`System.IsDesktop()`（macOS/Windows/Linux）、`System.IsServer()`（`server`构建标签），以及用于直接测试单个目标的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在每个目标上编译，因此可以在不使用构建标签的情况下进行条件分支。`@wailsio/runtime`中还提供了对应的前端辅助函数（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）
- 新增“使用其他前端框架”指南，说明如何将自己的 Vite 项目直接放入`frontend/`（涵盖 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`向导现在会检查移动端（iOS/Android）工具链，包括 Xcode 和 iOS Simulator 运行时、JDK、Android SDK/NDK 及模拟器；在适用情况下，还提供一键安装和可复制的 shell 配置修复方案
- 生成的项目附带一个`frontend/.npmrc`，其中设置了7天的`minimum-release-age`，以减少接触刚发布且可能已遭入侵的软件包（pnpm 和 bun 会遵循此设置；npm 会安全地忽略它）

## 变更

- 采用全新的霓虹山景主视觉，重新设计所有内置入门模板（Web、iOS 和 Android）
- **TypeScript 现已成为入门模板的默认语言，并使用不带后缀的模板名称。**`wails3 init`（不含`-t`）会搭建 TypeScript 项目；`-t vanilla`、`-t react`、`-t vue`和`-t svelte`均为 TypeScript 模板，对应的 JavaScript 变体分别为`-t vanilla-js`、`-t react-js`、`-t vue-js`和`-t svelte-js`。内置模板通过`template.yaml`中的`typescript:`声明其语言；使用`-ts`后缀的社区模板仍可作为后备方案继续使用
- 采用霓虹“数字 Wails”主题（以山景为背景的磨砂玻璃鲜活效果）重新设计`wails3 setup`向导

## 修复

- 修复 Windows 上的崩溃：当应用最小化时间过长，导致 WebView2 挂起或其渲染/GPU 进程被回收后，恢复应用时不再崩溃。最小化/恢复时的 DPI 重新同步（#5544）现在仅在窗口 DPI 确实发生变化时才操作 WebView2 控制器，从而避免在常见的相同 DPI 恢复场景中对已挂起的控制器进行致命的 COM 调用（#5605）
- 修复长期运行的 Linux 应用在频繁加载资源/媒体时反复发生的原生`SIGABRT`/`SIGSEGV`崩溃（通常发生在 GTK 主循环期间的`g_object_unref`内部）。资源服务器此前从工作 goroutine 完成`WebKitURISchemeRequest`，导致在 GTK 主线程之外调用线程不安全的 WebKit2GTK 函数；现在完成操作（`webkit_uri_scheme_request_finish_with_response`/`finish_error`）会在主线程上运行。此修复补全了 #5566 中的部分修复。GTK3 和 GTK4/WebKitGTK 6.0构建均受影响（#5631、#5557）
- 修复 Linux/GTK3 上`setupSignalHandlers`中间歇性出现的`fatal error: invalid pointer found on stack`。作为信号`user_data`传递的窗口 ID 保存在 Go 的`unsafe.Pointer`局部变量中，因此垃圾回收器在复制栈期间扫描到该（非指针）值时会中止。现在 Go 端将该 ID 保持为整数类型（`uintptr_t`），把 #4958 针对 GTK4 路径实施的同一修复向后移植到旧版 GTK3 路径（GTK4 路径已将 C 信号函数改为`uintptr_t`，以消除`-race`/checkptr 错误）（#5631）

## 移除

- 移除`react-swc`、`preact`、`lit`、`solid`、`qwik`和`sveltekit`入门模板（以及它们的`-ts`变体）。目前支持的内置模板为`vanilla`、`react`、`vue`和`svelte`——每个模板默认使用 TypeScript，并提供`-js` JavaScript 变体。仍可通过[自带前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自定义模板使用任何其他框架

## v3.0.0-alpha2.104 - 2026-06-18

## 修复

- 修复绑定的 Go 服务方法返回空字符串时发生的 iOS 崩溃（SIGABRT）。iOS 资源响应写入器此前根据`buf != nil`而非正文长度来保护正文指针，因此长度为零的正文会导致`&buf[0]`发生 panic；现在改为根据长度进行保护，与桌面端写入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 变更

- 将 iOS 和 Android 原生功能迁移到平台管理器：通过`application.IOS.*`和`application.Android.*`调用它们（例如`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），不再使用旧的`application.IOS*`/`application.Android*`自由函数（#5602）
- 重命名移动端桥接事件：跨平台事件现在使用`common:*`前缀（例如`common:haptic`、`common:location`），平台专属事件使用`ios:*`/`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）；不再使用`native:*`前缀（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增实验性的`wails3 setup`向导，用于交互式项目设置和依赖项检查
- 为`wails3 doctor`新增`--json`标志，以输出机器可读内容
- 为`wails3 doctor`命令新增签名状态部分

## 修复

- 修复 Linux 上的 npm 检测，使其除检查包管理器外还会检查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生消息对话框（UIAlertController）以及打开文件/多个文件/目录的对话框（UIDocumentPickerViewController）；保存对话框会返回明确错误
- iOS：通过 UIPasteboard 支持剪贴板
- iOS：通过 UIScreen 获取真实屏幕指标（点、像素、缩放比例、安全区域内的工作区）
- iOS：设备构建（`IOS_PLATFORM=device`）、代码签名身份/预置描述文件/权限配置支持、`.ipa`打包，以及通过 devicectl 执行`deploy-device`
- iOS：可配置最低 iOS 版本（build/config.yml 中的`ios.minIOSVersion`）
- iOS：`wails3 doctor`会报告 macOS 上 Xcode 和 iOS SDK 的可用情况
- iOS：系统事件——电池、网络、主题、屏幕锁定和低内存状态以`events.IOS.*`及平台无关的`events.Common.*`应用事件形式提供
- iOS：原生移动功能桥接（导出的`application.IOS*`）——系统分享面板、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（Face ID/Touch ID）、本地通知和 Keychain 安全存储
- iOS：传感器和硬件——触觉反馈、单次地理定位、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距以及屏幕捕获检测
- iOS：文档（IOS.md 和文档网站指南）
- Android：原生消息对话框（AlertDialog）以及打开文件/多个文件的对话框（Storage Access Framework，以缓存副本形式导入）；打开目录和保存对话框会返回明确错误
- Android：通过 ClipboardManager 支持剪贴板
- Android：通过 WindowMetrics/DisplayMetrics 获取真实屏幕指标（dp、像素、缩放比例、扣除系统栏占用后的可用工作区）
- Android：触觉反馈（`Android.Haptics.Vibrate`）、设备信息（`Android.Device.Info`）和 Toast 消息（`Android.Toast.Show`）运行时方法
- Android：类型化生命周期事件（`events.Android.*`，从 events.txt 生成），其中`ActivityCreated`映射到`Common.ApplicationStarted`
- Android：构建流水线会生成可安装的调试版和发布版 APK（`android:run`、`android:package`、`android:package:fat`）；发布签名默认使用调试密钥库，也可通过`ANDROID_KEYSTORE_*`环境变量使用真实密钥库
- Android：`wails3 doctor`会报告 Android SDK、NDK 和 JDK
- Android：系统事件——电池、网络、主题、屏幕锁定和低内存状态以`events.Android.*`及平台无关的`events.Common.*`应用事件形式提供
- Android：原生移动功能桥接（导出的`application.Android*`）——共享、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（BiometricPrompt）、本地通知和 EncryptedSharedPreferences 安全存储
- Android：传感器和硬件——触觉反馈、单次地理定位、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距以及通过 FLAG_SECURE 阻止屏幕捕获
- Android：文档（ANDROID.md 和文档网站指南）
- 示例：`mobile`综合示例新增“移动端”和“硬件”选项卡，用于演示跨 iOS 和 Android 的原生功能桥接（胶囊式选项卡会换行到多行）
- 移动端：电池——应用进入后台时会暂停加速度计、接近传感器、手电筒和示例中的周期性时钟，返回前台时恢复（Android 会让进程继续在后台运行，而手电筒是会在 iOS 上持续保持的硬件状态）；此外，Android 系统事件接收器仅在应用位于前台时注册
- iOS：相机拍摄 — `application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → 携带 base64 缩略图的`native:capture`事件）
- iOS：后台执行 — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication 后台任务时间窗口），以及可配置的`ios.backgroundModes`（build/config.yml），用于将`UIBackgroundModes`模板化写入生成的 Info.plist
- Android：相机拍摄 — `application.AndroidCapturePhoto`/`AndroidCaptureVideo`（通过 FileProvider 调用系统相机 → `native:capture`事件）
- Android：前台服务 — `application.AndroidStartForegroundService`/`AndroidStopForegroundService`（带有常驻通知的`WailsForegroundService`可使进程保持运行，以执行长时间运行的后台工作）
- 示例：新增相机选项卡，用于演示照片/视频拍摄和后台执行（Android 使用前台服务，iOS 使用后台任务时间窗口）

## 修复

- 修复 Linux 上`getUserMedia`始终因`NotAllowedError`而失败的问题：WebKitGTK 会拒绝无人处理的权限请求，而`permission-request`信号此前未连接。现在通过新的跨平台`WebviewWindowOptions.Permissions`映射（`map[PermissionType]Permission`）处理相机和麦克风权限，Linux（WebKitGTK）和 Windows（WebView2）均会遵循该映射。Linux 没有原生权限提示，因此默认允许相机和麦克风（恢复`getUserMedia`），也可通过`PermissionDeny`将其关闭（#5552）
- iOS：`GOOS=ios`现在可再次编译（导出了`events.IOS`并添加了移动端方法名存根），带 production 标签的构建也可正常编译（修复了 pkg/application 和多个服务中的构建标签）
- iOS：Go→JS 事件和 ExecJS 现已可用——页面启动时不再加载两次，`wails:runtime:ready`握手也不会再丢失
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted`不再与应用启动发生竞态；移除了固定的2秒启动等待
- iOS：修复了每次执行 Go→JS JavaScript 时发生的 C 字符串泄漏
- iOS：`hasListeners`现在会反映实际的监听器注册状态
- iOS：生产构建在编译时会移除框架调试日志
- Android：`GOOS=android`现在可再次编译——定义了`events.Android`，移除了越界的`events_android.go`监听器数组，添加了移动端方法名存根，并阻止桌面 Linux 文件（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 构建
- Android：JS→Go 绑定现已可用——WebView 无法将`fetch()`的 POST 请求体传递给`shouldInterceptRequest`，因此运行时调用现在通过 JavascriptInterface 传输（`nativeHandleRuntimeCall`）进行路由，不再因请求体为 nil 而崩溃
- Android：`Screens.*`运行时调用现在会返回真实数据——ScreenManager 现已在启动时填充（此前从未接入，因此`GetAll`返回 nil）
- Android：生产构建在编译时会移除框架调试日志；调试构建中的日志则以`Wails`标签输出到 logcat
- Android：实现了真正的`hasListeners`注册表、JNI 引用/异常处理以及单次加载的页面生命周期（不再重复导航）
- 修复在 Vite 开发服务器运行时，Windows 上`wails3 generate bindings`因“访问被拒绝”而失败的问题：现在将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5515）
- 修复 macOS 上显示器发生变化后读取屏幕信息时偶发的致命崩溃：屏幕 ID 和名称所保存的指针指向自动释放的`UTF8String`缓冲区，这些缓冲区可能在 Go 完成复制前被释放（释放后使用）。现在会对这些字符串调用`strdup`，并在转换后将其释放；屏幕枚举也会在显式的自动释放池中运行，因此从 Go goroutine 调用时不再发生泄漏（#5556）
- 修复 Linux 上 assetserver 关闭`WebKitURISchemeRequest`时偶发的 SIGSEGV：最后一次`g_object_unref`在 assetserver goroutine 上运行，导致 WebKit GObject 在 GTK 主线程之外完成终结。现在通过`g_main_context_invoke`将 unref 调度到 GTK 主上下文中执行（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 新增

- 扩展`MacWebviewPreferences`，新增以下 WKWebView 配置选项：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`和`ApplicationNameForUserAgent`（#5549）

## 修复

- 修复在 Vite 开发服务器运行时，Windows 上`wails3 generate bindings`因“访问被拒绝”而失败的问题：现在将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5561）
- 修复 Linux 上无边框窗口不触发 JS 调整大小事件的问题；修复无边框窗口的滚动条边缘检测（#5368）
- 修复 Windows 更新程序在临时目录与安装目录位于不同卷时因“无效的跨设备链接”而失败的问题（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修复

- 修复在 Vite 开发服务器运行时，Windows 上`wails3 generate bindings`因“访问被拒绝”而失败的问题：现在将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修复

- 修复 Linux 上 WebKit 空闲时（例如打开检查器时）界面冻结的问题：不再对`SIGUSR1`强制设置`SA_ONSTACK`，因为这会破坏 JavaScriptCore 的 GC 线程同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 新增

- 新增调试页面以及使用`runtime/trace`的相关内容

## 变更

- 移除了一些不必要的`_ "embed"`导入，稍微整理了代码

## 修复

- 修复 Windows 上窗口取消最大化后未强制执行最小宽度/高度约束的问题（#4593）
- 修复同时使用 Frameless 和 Transparent 窗口选项时，全屏模式下鼠标点击会穿透的问题（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 新增

- 新增 Garble 混淆支持（[#4563](https://github.com/wailsapp/wails/issues/4563)）：稳定的绑定方法 ID、构建/Taskfile 衔接（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），以及所有面向运行时的载荷上的 JSON 结构体标签（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），确保即使 Garble 重命名导出字段，线协议格式也能保持不变。

## v3.0.0-alpha.95 - 2026-05-20

## 新增

- 补充缺失的项目结构页面

## 变更

- 文档：调整架构页面中的几个图表，改用序列图以获得更清晰的显示效果
- 文档：添加说明，指出运行前需安装 D2

## 修复

- 修复 GTK4 默认配置下的`wails3 generate appimage`：打包器现在会先从二进制文件检测 GTK 技术栈，再搜索运行时文件，因此 GTK4 构建会选取`libwebkitgtkinjectedbundle.so`（位于`webkitgtk-6.0/`下），`-tags gtk3`构建会选取`libwebkit2gtkinjectedbundle.so`（位于`webkit2gtk-4.1/`下）。`.relr.dyn`探测还会检查`libgtk-4.so.1`，因此无论使用哪种技术栈，现代工具链都能正确禁用符号剥离。(#5475)
- 修复使用相对`-builddir`调用`wails3 generate appimage`时失败的问题：打包器现在会预先将`-binary`、`-icon`、`-desktopfile`、`-builddir`和`-outputdir`解析为绝对路径，使流程中途执行的`s.CD`不会再破坏 AppRun 下载 goroutine 或复制后的`ldd`探测。
- 修复桌面文件中的`Name=`字段与二进制文件的基本名称不匹配时，`wails3 generate appimage`无法将最终 AppImage 移至`-outputdir`的问题：打包器现在通过`OUTPUT`环境变量强制 linuxdeploy 的 appimage 插件将 AppImage 写入`<binary>-<arch>.AppImage`，而不再使用根据桌面文件派生的名称。
- 修复在 alpha.93 中将 GTK4 + WebKitGTK 6.0技术栈提升为默认配置后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`在 Linux 上不触发的问题。新的默认`application_linux.go` `run()`没有调用`setupCommonEvents()`（该函数会将`Linux.*`事件转发给对应的`Common.*`事件）或`monitorPowerEvents()`。现在通过`application_linux_dbus.go`在 GTK3 和 GTK4 构建路径之间共享 DBus 电源监视辅助程序。(#5474)

## v3.0.0-alpha.94 - 2026-05-19

## 修复

- 修复在 alpha.93 中将 GTK4 + WebKitGTK 6.0技术栈提升为默认配置后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`在 Linux 上不触发的问题。新的默认`application_linux.go` `run()`没有调用`setupCommonEvents()`（该函数会将`Linux.*`事件转发给对应的`Common.*`事件）或`monitorPowerEvents()`。现在通过`application_linux_dbus.go`在 GTK3 和 GTK4 构建路径之间共享 DBus 电源监视辅助程序。(#5474)

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 上为`wails3 doctor`输出添加`XDG_SESSION_TYPE`

## 修复

- 由 @leaanthony 修复 appmenu-gtk-module 访问尚未创建底层原生窗口资源的窗口所导致的 Wayland 窗口菜单崩溃问题（#4769）
- 由 @leaanthony 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题
- 由 @overlordtm 修复在 Windows 上初始化拖放时出现“内存不足”错误的问题 (#4701)
- 由 @leaanthony 修复主线程回调存储在删除映射项时错误使用 RLock 所导致的竞态条件（Linux、macOS、iOS）(#4424)
- 修复将命令行参数传递给任务时的变量处理。现在会正确初始化以 KEY=VALUE 键值对形式指定的 CLI 变量，并在整个任务执行过程中传递这些变量。
- 修复 macOS 上的 NSWindowZoomButton 冲突：`MaximiseButtonState`和`FullscreenButtonState`现在都会在启动时和运行时采用限制更严格的状态；任一 setter 都无法再悄然覆盖另一个 setter (#5319)
- 修复 CodeRabbit 在 #5463 中发现的旧版 GTK3 构建路径（`-tags gtk3`）中一组原有缺陷：通过文件关联启动时不再跳过启动处理程序；`getTheme`现在具备边界和类型安全性；`appName`不再释放由 GLib 所有的内存；`clipboardGet`不再泄漏 GTK 返回的`gchar*`；`Calloc`现在使用指针接收者（且`NewCalloc`返回`*Calloc`），使池能够真正跟踪分配；`zoomOut`现在使用`zoomInFactor`的倒数，不再使用会被钳制为1.0的负乘数；`execJS`会复用预分配的空 world-name，不再每次调用都泄漏一个`C.CString("")`；已从`menuItem.setAccelerator`中移除开发期间使用的`fmt.Println`。解决 #5465。
- 修复默认 GTK4 构建路径（`linux_cgo.go`）中相同的`Calloc`值接收者泄漏问题：通过指针接收者和`NewCalloc() *Calloc`，现在可以真正跟踪并释放每个窗口的`c.String(...)`分配。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，以便通过`PACKAGE_MANAGER`选项控制所使用的前端包管理器
- 在模板数据中添加`{{.Opn}}`和`{{.Cls}}`，使 Taskfile 模板的编写更可预测

## 变更

- 修改部分现有 Taskfile，使其使用`{{.Opn}} and {{.Cls}}`

## 修复

- 修复面板读取托盘菜单期间更新该菜单时，`linuxSystemTray`中发生`concurrent map read and map write`运行时致命错误的问题。
- WebView2 错误和堆栈跟踪输出改用`log`而非`fmt`，以免应用在 Windows 上未附加控制台的情况下运行时丢失消息。

## v3.0.0-alpha.91 - 2026-05-12

## 变更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5414)中更新赞助者 SVG
- <strong>破坏性变更 (macOS)：</strong>统一 macOS 坐标系，使`GetScreens`、`Position`和`SetPosition`均使用同一坐标空间——采用逻辑点、Y 轴向下，并以主屏幕左上角的`(0,0)`为原点。这与 Windows、GTK 以及 Electron 和 Web 的公共 API 一致。物理位置在主屏幕上方的屏幕现在会报告负`Bounds.Y`（以前为正值），而`Position()`/`SetPosition()`值现在使用逻辑点，不再使用`points × primaryScale`。仍可保持`Position()` → `SetPosition()`往返转换不变；但需要更新早期 alpha 构建所记录的绝对值或手动计算的权宜方案（例如乘以`primaryScale`，或根据屏幕高度翻转 Y 轴）。解决[#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5416)中以防御性方式验证 DBus 信号名称和消息体长度，防止发生 panic
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5363)中修复 Linux 上 GTK 菜单处理的内存安全问题
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中检测 NVIDIA GPU，并在 Linux 上禁用 DMA-BUF 渲染器
- 在[#5117](https://github.com/wailsapp/wails/issues/5117)中修复 macOS 上`SetPosition`跨屏幕 Y 坐标转换：以主屏幕高度作为全局基准，使窗口在相对于主显示器存在垂直偏移的显示器上落于正确位置
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修复 Git PR 模板，使其指向正确的反馈 URL
- 修复了一系列 Windows 系统托盘`SetMenu`崩溃问题。这些问题由损坏的`DestroyMenu`系统调用引起：该调用错误地传入了四个参数，而非一个，导致每次调用都返回 FALSE 且不释放任何内容。此外，在重新构建菜单时释放 HMENU 和 HBITMAP 句柄（包括运行时通过`MenuItem.SetBitmap`分配的句柄），重置`Win32Menu.Update`中过时的复选框/单选项映射，并移除`systemtray.updateMenu`中一个导致分配次数翻倍的冗余`Update()`调用。长时间运行的系统托盘应用在每次重新构建菜单时不再泄漏 GDI/USER 对象。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5261)中新增在 macOS 上为 WKWebView User-Agent 配置应用名称的功能，贡献者：@vinhvoit225
- 在[PR](https://github.com/wailsapp/wails/pull/5400)中为 gin-service 示例新增间接依赖 github.com/coder/websocket，贡献者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5402)中为构建资源测试新增深度相等比较支持，贡献者：@leaanthony

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5401)中将构建输出统一到 assets 目录，贡献者：@taliesin-ai
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5399)中更新赞助者 SVG

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5289)中为 macOS 单实例消息使用通知对象，贡献者：@overlordtm
- 在[PR](https://github.com/wailsapp/wails/pull/5383)中批量处理 Windows 回调，以防止高负载下 Promise 丢失，贡献者：@taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5316)中新增 go<em>test</em>results 作业，用于汇总 Go 测试结果，贡献者：@leaanthony

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5369)中根据条件将大型 RPC 载荷拆分为分块 POST 请求，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5386)中将所有前端模板的 Vite 从5.x.x 升级到8.0.0，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5365)中将 Vite 开发服务器端口配置迁移到环境变量，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5361)中配置所有模板，使 Vite 开发服务器绑定到127.0.0.1，贡献者：@leaanthony
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5384)中更新赞助者 SVG

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5312)中，于 build-assets 更新期间清理 Info.plist 模板存根，贡献者：@leaanthony
- 通过在主线程上同步应用菜单项修改方法（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修复 macOS 菜单中的过时状态；这消除了`dispatch_async`竞态，避免菜单快速重新打开时仍呈现先前状态（#5002）
- 在[PR](https://github.com/wailsapp/wails/pull/5203)中于开发模式下忽略`*_test.go`文件，以防止不必要的重新构建，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止应用未运行时 Menu.Update() 发生段错误，贡献者：@wucm667
- 在[PR](https://github.com/wailsapp/wails/pull/5382)中使用 lastSizeWParam 控制 Windows 上菜单栏的重绘，贡献者：@taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5371)中将 HiddenOnTaskbar 改为使用 WS<em>EX</em>TOOLWINDOW，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5370)中重新排列依赖项，并移除 go.mod 中的 webview2 replace 指令，贡献者：@atterpac
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5358)中更新赞助者 SVG

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5331)中移除泛型间接别名，并统一映射键类型，贡献者：@fbbdev

## 移除

- 在[PR](https://github.com/wailsapp/wails/pull/5377)中删除 PR-master 工作流，移除文档、Go 测试和跳过测试，贡献者：@leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5352)中新增 Wails v3 韩语文档，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5354)中新增安装和快速入门的法语文档，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5355)中新增快速入门、概念和社区的葡萄牙语文档，贡献者：@leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5328)中新增法语文档本地化，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5343)中为 文档站点新增德语区域设置，贡献者：@leaanthony

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5347)中将全部8个已翻译的区域设置注册到 文档配置，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5317)中更新与 WebView2 有关的多个 Windows 文件，贡献者：@leaanthony

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5340)中拆分 Linux 上 GTK3 与 GTK4 的对话框分派逻辑，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5339)中确保对话框回调在 GTK 线程上执行，从而修复段错误，贡献者：@leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5179)中向仓库添加 PR 模板 URL，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5330)中新增 Wails v3 德语文档，贡献者：@leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5307)中新增一个选项，用于禁止在 macOS 上按 Escape 键退出全屏模式，贡献者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/5310)中新增一个选项，用于禁止在 macOS 上按 Escape 键退出全屏模式，贡献者：@leaanthony
- 由 @yuseferi 在[PR](https://github.com/wailsapp/wails/pull/5288)中添加 Pausa 社区展示文档

## 变更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5308)中更新赞助者 SVG
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5309)中更新图标生成命令，以处理不受支持的平台
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5224)中以三态 ButtonState 替换布尔型全屏 API，并实现平台绑定

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5315)中防止 WebView2 在控制器状态为 nil 时执行焦点操作
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5313)中更新 GitHub Actions 工作流，使其正确引用 PR 的基础分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5203)中使开发模式忽略`*_test.go`文件，以防止不必要的重新构建
- 由 @wucm667 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止应用未运行时 Menu.Update() 发生段错误

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5094)中添加 InstallScope 标志以及计算机/用户安装构建选项
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5294)中为 BrowserWindow 添加无操作的 SetScreen 方法，以满足 Window 接口要求

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中检测 NVIDIA GPU，并在 Linux 上禁用 DMA-BUF 渲染器
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修复 Git PR 模板，使其指向正确的反馈 URL
- 修复了一系列 Windows 系统托盘`SetMenu`崩溃问题。这些问题由损坏的`DestroyMenu`系统调用引起：该调用错误地传递了四个参数而非一个，导致每次调用都返回 FALSE，且不释放任何内容。此外，在重新构建菜单时释放 HMENU 和 HBITMAP 句柄（包括运行时通过`MenuItem.SetBitmap`分配的句柄），重置`Win32Menu.Update`中过期的复选框/单选按钮映射，并移除`systemtray.updateMenu`中导致分配量翻倍的冗余`Update()`调用。长期运行的系统托盘应用不再在每次重新构建菜单时泄漏 GDI/USER 对象。

## v3.0.0-alpha.82 - 2026-05-01

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5232)中修复桌面文件生成，使其能够正确处理桌面名称

## v3.0.0-alpha.81 - 2026-04-30

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5286)中将每夜版发布时间调整为 UTC 15:00

## 修复

- 修复 Retina Mac 上 Screen Bounds、WorkArea 和 Size 值减半的问题 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5285)中更新文档依赖项和内容集合加载器

## v3.0.0-alpha.79 - 2026-04-29

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5270)中向 trigger-release 作业授予 actions: write 权限

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5283)中将发布任务的默认分支设为 master，并更新变更日志措辞
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5282)中更新自动生成变更日志的工作流，使其使用最新版本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5280)中添加路径过滤器并移除失效工作流，以提高工作流效率
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5274)中更新文档，使示例链接引用 master 分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5272)中更新 v3 的文档和示例

## 修复

- 由 @AkagiYui 在[PR](https://github.com/wailsapp/wails/pull/5265)中为反向代理添加重试逻辑，并在开发环境中强制使用 IPv4
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5281)中重写未发布变更日志的触发工作流

## 移除

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5267)中移除用于各种测试目的的 shell 测试脚本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5266)中删除 v3-alpha 文档部署工作流和 CNAME 记录

### 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5196)中向侧边栏导航添加“前端路由”条目
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5185)中添加前端路由指南，其中包含针对不同框架的建议
- 添加对窗口附属模态对话框（macOS）的支持
- 由 @leaanthony 升级 ghw 版本，以更好地支持 Apple 设备 (#4977)
- 向程序坞服务添加`GetBadge`方法
- 向`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）(#4957)
- 添加绑定生成器自动生成枚举的文档，包括专门的“枚举”页面和侧边栏导航 (#4972)
- 向`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）(#4957)
- 在`v3/examples/web-apis/`中添加 Web API 示例，演示41个浏览器 API，包括存储（localStorage、sessionStorage、IndexedDB、Cache API）、网络（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒体（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、设备（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、性能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 添加 WebView API 兼容性检查器示例（`v3/examples/webview-api-check/`），用于跨平台测试 200+ 个浏览器 API
- 添加`internal/libpath`包，用于在 Linux 上查找原生库路径，支持并行搜索、缓存以及 Flatpak/Snap/Nix
- <strong>开发中：</strong>为 Linux 添加实验性 WebKitGTK 6.0 / GTK4 支持，可通过`-tags gtk4`使用（GTK3/WebKit2GTK 4.1仍为默认选项）
- 注意：在平铺式窗口管理器（例如 Hyprland、Sway）上，由于窗口管理器会控制窗口几何属性，最小化/最大化操作可能无法按预期工作
- 由 @AbdelhadiSeddar 在<strong>在 JavaScript 中监听事件</strong>文档中添加了如何使用<strong>一次性处理函数</strong>的说明
- 由 @leaanthony 为`WebviewWindowOptions`添加`UseApplicationMenu`选项，使 Windows/Linux 上的窗口可以继承通过`app.Menu.Set()`设置的应用程序菜单
- 由 @wimaha 添加使用`.icon`文件（Apple Icon Composer 格式）生成 Liquid Glass 图标和资源目录的支持（macOS）（#4934）
- 为无头/Web 部署添加实验性服务器模式（`-tags server`）。该模式允许 Wails 应用在不依赖原生 GUI 的情况下作为 HTTP 服务器运行。使用`wails3 task build:server`构建。详情请参阅`examples/server`。
- 添加`internal/libpath`包，用于在 Linux 上查找原生库路径，支持并行搜索、缓存以及 Flatpak/Snap/Nix
- 由 @leaanthony 为`MacWindow`添加`CollectionBehavior`选项，用于控制窗口在 macOS 空间和全屏模式下的行为（#4756）
- 由 @leaanthony 为 pkg/application 添加单元测试
- 由 @leaanthony 为 MSIX 打包添加自定义协议支持
- 在 Linux 上添加桌面环境检测[PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 为 JavaScript 运行时添加`Window.Print()`方法，以便从前端触发打印对话框（#4290）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 由 @leaanthony 为 Linux 添加更多 WebKit2 加载状态变更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 在 Linux 构建期间生成`.desktop`文件，而不再仅在打包时生成（#4575）
- 由 @leaanthony 添加 Linux 运行时依赖项文档，其中包含各发行版的软件包名称和 nfpm 打包示例（#4339）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加 NVIDIA 驱动程序版本信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中为 macOS 添加通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @chinenual 在[PR](https://github.com/wailsapp/wails/pull/4760)中为 helloworld 模板添加 aria-label 标识符，使 Appium 测试客户端可以轻松测试示例应用
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中为 macOS 添加通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @fbbdev 和 @ianvs 在[#4633](https://github.com/wailsapp/wails/pull/4633)中添加类型化事件
- 添加`systray-clock`示例，展示实时更新工具提示的无头托盘（#4653）。
- 由 @Tolfx 在 #4510 中添加适用于 Windows 的 NSIS Protocol 模板
- 由 @Tolfx 在 #4510 中为 build-assets 添加测试
- macOS：由 @nidib 在[#4588](https://github.com/wailsapp/wails/pull/4588)中使原生窗口控件显示在菜单栏中
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于在程序坞中隐藏/显示应用图标
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于在程序坞中隐藏/显示应用图标
- 由 @leaanthony 在[#4534](https://github.com/wailsapp/wails/pull/4534)中使用 NSGlassEffectView（macOS 15.0+）为 macOS 添加原生 Liquid Glass 效果支持，并以 NSVisualEffectView 作为回退方案，同时提供全面的材质自定义选项
- 由 @leaanthony 在[#4500](https://github.dev/wailsapp/wails/pull/4500)中添加浏览器 URL 清理功能。该功能基于 @APShenkin 在[#4484](https://github.com/wailsapp/wails/pull/4484)中的工作。
- 由[@leaanthony](https://github.com/leaanthony)基于[@Taiterbase](https://github.com/Taiterbase)在此[PR](https://github.com/wailsapp/wails/pull/4241)中的原始工作，为 Windows/Mac 添加内容保护功能
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4488)中添加支持，可通过`wails3 build`和`wails3 package`别名将 CLI 变量传递给 Task 命令（#4422）
- 由[@atterpac](https://github.com/atterpac)在[#4318](https://github.com/wailsapp/wails/pull/4318)中添加对放置区域的支持，通过事件源提供被放置元素的数据
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中为`WindowsWindow`选项添加`AdditionalLaunchArgs`，以便向 WebView2 浏览器传递额外的命令行参数。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中添加在 wails init 后自动运行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中添加 Windows Snap Assist 功能
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中为`WindowsWindow`选项添加`AdditionalLaunchArgs`，以便向 WebView2 浏览器传递额外的命令行参数。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中添加在 wails init 后自动运行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中添加 Windows Snap Assist 功能
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中添加 Windows `getAccentColor`实现
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中添加 Windows `getAccentColor`实现
- Windows 深色主题菜单和菜单栏。由 @leaanthony 在[a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)中实现
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4405)中重命名内置服务，使 JS/TS 绑定更加清晰
- `app.Env.GetAccentColor`用于获取用户系统的强调色。适用于 MacOS。由[@etesam913](https://github.com/etesam913)实现
- 由[@atterpac](https://github.com/atterpac)在[#4137](https://github.com/wailsapp/wails/pull/4137)中添加`window.ToggleFrameless()` API
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4345)中添加 Linux 各发行版专用的构建依赖项
- 由 @atterpac 在[PR](https://github.com/wailsapp/wails/pull/4404)中添加绑定指南
- **整理测试基础设施**：将 Docker 测试文件移至专用的`test/docker/`目录，优化镜像并提高构建可靠性；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
- **改进资源管理模式**：在示例中添加正确的事件处理程序清理和上下文感知的 goroutine 管理；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
- 支持 aarch64 AppImage 构建；由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中完成
- 由[@leaanthony](https://github.com/leaanthony)为`wails doctor`添加诊断章节
- 调用服务方法时将窗口添加到上下文中；由[@leaanthony](https://github.com/leaanthony)完成
- 添加`window-call`示例，演示如何确定哪个窗口正在调用服务；由[@leaanthony](https://github.com/leaanthony)完成
- 新增菜单指南；由[@leaanthony](https://github.com/leaanthony)完成
- 改进 panic 处理；由[@leaanthony](https://github.com/leaanthony)完成
- 新增菜单指南；由[@leaanthony](https://github.com/leaanthony)完成
- 为 Service API 添加文档注释；由[@fbbdev](https://github.com/fbbdev)在[#4024](https://github.com/wailsapp/wails/pull/4024)中完成
- 添加`application.NewServiceWithOptions`函数，以使用额外配置初始化服务；由[@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中完成
- 改进菜单控制；由[@FalcoG](https://github.com/FalcoG)和[@leaanthony](https://github.com/leaanthony)在[#4031](https://github.com/wailsapp/wails/pull/4031)中完成
- 添加更多文档；由[@leaanthony](https://github.com/leaanthony)完成
- 支持在标准事件监听器中取消事件；由[@leaanthony](https://github.com/leaanthony)完成
- 支持系统托盘的`Hide`、`Show`和`Destroy`；由[@leaanthony](https://github.com/leaanthony)完成
- 支持系统托盘的`SetTooltip`；由[@leaanthony](https://github.com/leaanthony)完成。最初构想来自[@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- 在绑定生成器针对不支持类型发出的警告中报告包路径；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 添加绑定生成器对泛型别名的支持；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 添加绑定生成器对`omitzero` JSON 标志的支持；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 添加`//wails:ignore`指令，以阻止为选定的服务方法生成绑定；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 为服务和模型添加`//wails:internal`指令，以允许类型在 Go 中导出但不在 JS/TS 中导出；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 添加绑定生成器对别名类型常量的支持，以便使用弱类型枚举；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 添加针对 Go 1.24功能的绑定生成器测试；由[@fbbdev](https://github.com/fbbdev)在[#4068](https://github.com/wailsapp/wails/pull/4068)中完成
- 为`OSInfo.Branding`添加对 macOS 15“Sequoia”的支持，以改进操作系统版本检测；在[#4065](https://github.com/wailsapp/wails/pull/4065)中完成
- 添加`PostShutdown`钩子，用于在关闭流程完成后运行自定义代码；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 添加`FatalError`结构体，以支持在自定义错误处理程序中检测致命错误；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 统一服务的启动和关闭顺序并为其编写文档；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 为应用程序启动/关闭序列以及服务启动/关闭测试添加测试框架；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 添加`RegisterService`方法，用于在应用程序创建后注册服务；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 在应用程序和服务选项中添加`MarshalError`字段，用于自定义处理绑定调用中的错误；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 添加可取消的 Promise 包装器，可沿 Promise 链传播取消请求；由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 支持将绑定调用的取消操作关联到`AbortSignal`；由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- WML 除了支持常规的`wml-*`属性外，还支持`data-wml-*`属性；由[@leaanthony](https://github.com/leaanthony)完成
- 为所有服务添加`Configure`方法，用于后期配置和动态重新配置；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- `fileserver`服务在未配置时发送503 Service Unavailable 响应；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- `kvstore`服务在未配置时默认提供内存键值存储；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- 为`kvstore`服务添加`Load`方法，用于在配置更改后从文件重新加载数据；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- 为`kvstore`服务添加`Clear`方法，用于删除所有键；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- 在`log`服务中添加`Level`类型，以提供 JS 端的日志级别常量；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- 为`log`服务添加`Log`方法，用于动态指定日志级别；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- `sqlite`服务在未配置时默认提供内存数据库；由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中完成
- 在`sqlite`服务上添加`Close`方法，以便手动关闭数据库，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 为`sqlite`服务的查询方法添加取消支持，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 为`sqlite`服务添加带 JS 绑定的预处理语句支持，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 添加 Gin 支持，由[Lea Anthony](https://github.com/leaanthony)在[PR](https://github.com/wailsapp/wails/pull/3537)中实现，基于[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始工作
- 修复自动保存和密码自动保存始终启用的问题，由[@oSethoum](https://github.com/osethoum)在[#4134](https://github.com/wailsapp/wails/pull/4134)中实现
- 为窗口添加`SetMenu()`，以便为窗口设置菜单，由[@leaanthony](https://github.com/leaanthony)实现
- 添加通知支持，由[@popaprozac](https://github.com/popaprozac)在[#4098](https://github.com/wailsapp/wails/pull/4098)中实现
-  添加 macOS 文件关联支持，由[@wimaha](https://github.com/wimaha)在[#4177](https://github.com/wailsapp/wails/pull/4177)中实现
- 添加用于递增语义版本号的`wails3 tool version`，由[@leaanthony](https://github.com/leaanthony)实现
- 添加 macOS 和 Windows 徽标支持，由[@popaprozac](https://github.com/popaprozac)在[#](https://github.com/wailsapp/wails/pull/4234)中实现
- 添加对已注册/强类型事件的支持，由[@fbbdev](https://github.com/fbbdev)和[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中实现
- 添加为自定义事件注册钩子的功能，由[@fbbdev](https://github.com/fbbdev)和[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中实现
- 添加`app.OpenFileManager(path string, selectFile bool)`，用于打开系统文件管理器并定位到路径`path`，还可选择通过`selectFile`突出显示，由[@Krzysztofz01](https://github.com/Krzysztofz01)和[@rcalixte](https://github.com/rcalixte)实现
- 为`wails3 init`命令添加新的`-git`标志，由[@leaanthony](https://github.com/leaanthony)实现
- 添加新的`wails3 generate webview2bootstrapper`命令，由[@leaanthony](https://github.com/leaanthony)实现
- 在运行时中添加`init()`方法，以支持手动初始化运行时，由[@leaanthony](https://github.com/leaanthony)实现
- 在窗口的 WindowOptions 中添加`WindowDidMoveDebounceMS`选项，由[@leaanthony](https://github.com/leaanthony)实现
- 添加单实例功能，由[@leaanthony](https://github.com/leaanthony)实现。该功能基于 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)。
- 添加`wails3 generate template`命令，由[@leaanthony](https://github.com/leaanthony)实现
- 添加`wails3 releasenotes`命令，由[@leaanthony](https://github.com/leaanthony)实现
- 添加`wails3 update cli`命令，由[@leaanthony](https://github.com/leaanthony)实现
- 为`wails3 generate bindings`命令添加`-clean`选项，由[@leaanthony](https://github.com/leaanthony)实现
- 支持为 aarch64 (arm64) 构建 Linux AppImage，由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中实现
- 添加赞助者超链接，由 @ansxuman 在[#3958](https://github.com/wailsapp/wails/pull/3958)中实现
- 支持构建 deb、rpm 和 Arch Linux 安装包，由
- 添加对 Darwin 通用构建和安装包的支持，由
- 为网站添加事件文档，由
- 添加配置为非 SSR 开发模式的 sveltekit 和 sveltekit-ts 模板
- 使用新的`wails3 update build-assets`命令更新构建资产，由
- 添加用于测试 HTML 拖放 API 的示例，由
- 添加文件关联支持，由[leaanthony](https://github.com/leaanthony)在
- 添加新的`wails3 generate runtime`命令，由
- 添加新的`InitialPosition`选项，用于指定窗口应居中还是
- 向`application`包添加`Path`和`Paths`方法，由
- 添加`GeneralAutofillEnabled`和`PasswordAutosaveEnabled` Windows 选项
- 添加获取调用服务方法的窗口的功能，由
- 为 WebView2 添加`EnabledFeatures`和`DisabledFeatures`选项，由
- ⊞ 添加用于增强高 DPI 显示器支持的新 DIP 系统，由
- ⊞ 添加窗口类名选项，由[windom](https://github.com/windom/)在
- 扩展服务以提供插件功能。由
- 🐧 添加 WindowDidMove / WindowDidResize 事件，位于
- ⊞ 添加 WindowDidResize 事件，位于
-  添加 ApplicationShouldHandleReopen 事件，以便处理程序坞
-  向实现中添加 getPrimaryScreen/getScreens，由 @tmclane 在
-  添加在 macOS 全屏模式下显示工具栏的选项，由
- 🐧 添加 onKeyPress 逻辑，将 Linux 按键转换为快捷键
- 🐧 添加任务`run:linux`，由
- 导出`SetIcon`方法，由[@almas-x](https://github.com/almas-x)在
- 改进`OnShutdown`，由[@almas-x](https://github.com/almas-x)在
- 在`Window`接口中恢复`ToggleMaximise`方法，由
- 为`Environment()`添加更多信息。由 @leaanthony 在
- 在`Window`接口上公开`WebviewWindow.IsFocused`方法，由
- 支持在 WML 系统中使用多个以空格分隔的触发事件，由
- 从捆绑的 JS 运行时脚本中添加 ESM 导出，由
- 添加绑定生成器标志，以使用捆绑的 JS 运行时脚本，而不是
- 在 Linux 上实现`setIcon`，由[@abichinger](https://github.com/abichinger)完成
- 向 dev 命令添加`-port`标志并支持环境变量
- 添加绑定方法调用测试，由
- ⊞ 为已创建的窗口添加`SetIgnoreMouseEvents`，由
-  添加设置窗口堆叠级别（顺序）的功能，由

### 修复

- 将 NSScreen 的点值转换为`Physical*`字段中的设备像素，修复 Retina Mac 上减半的`Screen.Bounds`、`WorkArea`和`Size`；同时填充顶层`Screen.X`/`Y`，使[PR](https://github.com/wailsapp/wails/pull/5168)中的多显示器相接检测和工作区定位正确，由 @wayneforrest 完成
- 修复 ScreenManager 中的数据竞争，该竞争会在显示配置发生变化时（例如休眠/唤醒期间热插拔外接显示器）导致 WebKit DisplayLink 死锁
- 当 Assets.car 存在时，直接将 CFBundleIconName 设置为 appicon，见[PR](https://github.com/wailsapp/wails/pull/5154)，由 @symball 完成
- 修复`wails3 doctor`在 Fedora、openSUSE、Arch 和 NixOS 上报告错误 WebKitGTK 软件包的问题——由于 v3 在编译时需要4.1 API，已移除4.0回退条目（#5071）
- 修复 openSUSE 的 webkit2gtk doctor 软件包名称（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，后者是正确的 openSUSE 软件包名称）（#5071）
- 修复桌面开发模式下缺少`/wails/custom.js`时出现的`Unexpected token '<'`错误。为`/wails/custom.js`添加了显式404处理程序，并在`loadOptionalScript`中添加了不区分大小写的`Content-Type`验证，以防止将 HTML SPA 回退内容作为 JavaScript 注入。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修复 macOS 上系统托盘菜单的高亮状态——菜单打开时，图标现在会显示选中状态（#4910）
- 修复 macOS 上附加到系统托盘的窗口显示在其他窗口后方的问题——现在使用正确的弹出窗口级别（#4910）
- 修复文档中错误的`@wailsio/runtime`导入示例（#4989）
- 修复 darwin 上无边框窗口无法最小化的问题（#4294）
- 通过从 go-task 的最新状态检查中排除`node_modules/`，修复`wails3 build`和`wails3 dev`期间持续20-30分钟的卡顿。此前，`sources: "**/*"`通配模式会导致 go-task 枚举`node_modules/`中的每个文件并计算校验和（使用 MUI 等大型依赖项时，文件数可达50000-100000以上），在 Windows/NTFS 上尤其缓慢（#4939）
- 修复 C 的`Screen` typedef 与 X11 Xlib.h 冲突导致的 GTK4 构建失败（#4957）
- 修复 macOS 上程序坞徽标方法的一致性问题
- 修复`InvisibleTitleBarHeight`应用于所有 macOS 窗口，而非仅应用于无边框窗口或标题栏透明窗口的问题（#4960）
- 启用`InvisibleTitleBarHeight`时，通过跳过窗口边缘附近的拖动启动，修复从顶部两角调整窗口大小时窗口晃动/抖动的问题（#4960）
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 完成
- 修复 Windows 上文件拖放在显示缩放比例不是100% 时无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上文件拖放与悬停效果配合时无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过使用`gtk_window_present()`，修复 Linux/GTK4 上显示/隐藏窗口时有时会恢复为最小化状态的问题（#4957）
- 通过`XTranslateCoordinates`/`XMoveWindow`添加以 X11 为条件的支持，修复 Linux/GTK4 上获取/设置窗口位置始终返回0,0的问题（#4957）
- 添加基于信号的大小限制来替代已移除的`gtk_window_set_geometry_hints`，修复 Linux/GTK4 上未强制执行窗口最大尺寸的问题（#4957）
- 通过实现正确的 PhysicalBounds 计算并借助`gdk_monitor_get_scale`支持分数缩放，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 完成
- 修复 Windows 上文件拖放在显示缩放比例不是100% 时无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上文件拖放与悬停效果配合时无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过实现正确的 PhysicalBounds 计算并借助`gdk_monitor_get_scale`支持分数缩放，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 完成
- 修复因 App.Window.Current() 未从主线程访问 AppKit API 而在 macOS 上引发的“幽灵窗口”问题（#4947），由 @wimaha 完成
- 通过实现 WKUIDelegate runOpenPanelWithParameters，修复 HTML `<input type="file">`在 macOS 上无法工作的问题（#4862）
- 修复在 macOS/Linux 上使用`@wailsio/runtime` npm 模块时原生文件拖放无法工作的问题（#4953），由 @leaanthony 完成
- 修复跨软件包类型别名的绑定生成问题（#4578），由 @fbbdev 完成
- 修复 Linux 上因违反 GTK 线程安全要求而导致 OpenFileDialog 崩溃的问题（#3683），由 @ddmoney420 完成
- 修复对隐藏或已销毁的窗口调用`Focus()`时发生 SIGSEGV 崩溃的问题（#4890），由 @ddmoney420 完成
- 修复 Linux 上设置空图标或位图时可能发生 panic 的问题（#4923），由 @ddmoney420 完成
- 修复 macOS 上从服务绑定调用 ErrorDialog 时发生崩溃的问题（#3631），由 @leaanthony 完成
- 使菜单可在 Windows 操作系统的`v3\examples\dialogs`中显示，由 @ndianabasi 完成
- 修复页面重新加载期间导致 TypeError 的竞争条件（#4872），由 @ddmoney420 完成
- 通过移除`Collector.IsVoidAlias()`方法中的全局状态，修复绑定生成器测试输出错误的问题（#4941），由 @fbbdev 完成
- 修复`<input type="file">`文件选择器在 macOS 上无法工作的问题（#4862），由 @leaanthony 完成
- 修复 macOS 上 `Position()` 和 `SetPosition()` 使用不一致的坐标系，导致保存/恢复状态时窗口位置漂移的问题（#4816），由 @leaanthony 修复
- 修复已通过应用程序清单设置 DPI 感知时 SetProcessDpiAwarenessContext 报“访问被拒绝”错误的问题（#4803）
- 更新键盘快捷键文档页面，并更正 `KeyBinding.Add` 回调参数的类型，由 @ndianabasi 完成
- 修复有关生成自定义绑定的文档：必须使用 `-d String`，而不是 `-o String`
- 修复调用 `menu.Update()` 时菜单未清除子项的问题
- 修复文档中过时的 Manager API 引用（更新 31 文件，改用 `app.Window.New()`、`app.Event.Emit()` 等新模式），由 @leaanthony 完成
- 修复因 WebKit 覆盖信号处理程序，导致绑定到 JS 的 Go 方法发生 panic 时 Linux 崩溃的问题（#3965），由 @leaanthony 修复
- 修复 SaveFileDialog.SetFilename() 在 Linux 上不生效的问题（#4841），由 @samstanier 修复
- 修复拖放示例中的放置坐标显示为 undefined 的问题
- 修复 APP_NAME 包含空格时创建 macOS 应用程序包失败的问题（花括号展开问题）
- 修复在 Windows 上调用服务方法时发生索引越界 panic 的问题（回退 goccy/go-json）
- 修复 Windows 在显示缩放比例不是 100% 时文件拖放不起作用的问题
- 修复在 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上配合悬停效果使用时文件拖放无法可靠工作的问题
- 修复在 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 更新所有操作系统的 Taskfile.yml 文件中的全部命令，以支持 `APP_NAME` 等变量包含空格的情况，由 @ndianabasi 完成
- 修复在 Linux 上执行 'build:universal:lipo:go' 任务时的命令参数错误，由 @wux1an 修复
- 修复在 Linux 上运行 'wails3 build GOOS=darwin GOARCH=arm64' 时出现的 Docker 错误“undefined symbol: **<em>ubsan</em>handle_xxxxxxx”，由 @wux1an 修复
- 整合自定义协议文档并添加通用链接（Universal Links）章节，由 @leaanthony 完成
- 通过添加防止并发调用 TrackPopupMenuEx 的保护措施，修复在 Windows 上反复单击系统托盘图标时菜单崩溃的问题（#4151），由 @leaanthony 修复
- 防止在 app.Run() 之前调用 systray.Run() 时应用程序崩溃，由 @leaanthony 完成
- 修复启用 ApplicationShouldTerminateAfterLastWindowClosed 后，通过 Hide()/Show() 切换窗口可见性时应用在 macOS 上崩溃的问题（#4389），由 @leaanthony 修复
- 修复在 macOS 和 Windows 上反复打开上下文菜单时的内存泄漏问题（#4012），由 @leaanthony 修复
- 修复 macOS 上未复用上下文菜单原生资源、导致每次显示时都重新创建菜单的问题（#4012），由 @leaanthony 修复
- 修复应用程序以 `Hidden: true` 启动时，单击 macOS 程序坞图标不会显示隐藏窗口的问题（#4583），由 @leaanthony 修复
- 修复因 CGO 调用中的窗口指针类型不正确而导致 macOS 打印对话框无法打开的问题（#4290），贡献者：@leaanthony
- 修复 Wayland 上因 appmenu-gtk-module 访问尚未创建底层原生窗口资源的窗口而导致窗口菜单崩溃的问题（#4769），贡献者：@leaanthony
- 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题，贡献者：@leaanthony
- 修复在 Windows 上初始化拖放功能时出现“内存不足”错误的问题（#4701），贡献者：@overlordtm
- 修复 Linux 上因 URI 转义不正确而导致文件资源管理器打开错误目录的问题（#4397），贡献者：@leaanthony
- 通过自动检测`.relr.dyn` ELF 节并禁用剥离，修复 AppImage 在现代 Linux 发行版（Arch、Fedora 39+、Ubuntu 24.04+）上的构建失败问题（#4642），贡献者：@leaanthony
- 修复`wails doctor`在 Fedora/DNF 系统上错误报告已安装 webkit 软件包的问题（#4457），贡献者：@leaanthony
- 修复默认`config.yml`会使用生产构建运行`wails3 dev`的问题，贡献者：@mbaklor
- 修复 iOS 服务存根因导入不存在的软件包而导致构建失败的问题，贡献者：@leaanthony
- 修复 debug/info 方法中的结构化日志记录导致“没有格式化指令”错误的问题，贡献者：@leaanthony
- 移除从移动平台合并中意外带入的临时调试打印语句，贡献者：@leaanthony
- 通过自动禁用 DMA-BUF 渲染器，修复 WebKitGTK 在配备 NVIDIA GPU 的 Wayland 上崩溃的问题（错误 71：协议错误），贡献者：@leaanthony
- 解决 Linux 上`application.WebviewWindowOptions.BackgroundColour`中的 alpha 值被忽略的问题（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修复未提供自定义图标时 Windows 系统托盘图标不默认使用应用图标的问题（#4704）
- 跟踪`HICON`的所有权，从而仅销毁用户创建的句柄，防止 Explorer 回收时崩溃（#4653）。
- 销毁时释放 Windows 系统主题监听器和保留的托盘图标，以免泄漏 goroutine 和设备上下文（#4653）。
- 将托盘工具提示截断为127个 UTF-16单元，以免破坏代理项对和多字节字形（#4653）。
- 修复 Windows 打包任务失败的问题（#4667）
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量，[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修复 go-webview2 v1.0.22签名变更导致的 Windows 构建错误（#4513、#4645）
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量，[PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 通过将`<.Info.Protocol>`改为`<.Protocol>`，修复 Linux desktop.tmpl 中的协议遍历问题，贡献者：@Tolfx，见 #4510
- 修复 Liquid Glass 演示中的重复定义错误，见[#4542](https://github.com/wailsapp/wails/pull/4542)，贡献者：@Etesam913
- 修复 Linux 上的系统托盘菜单更新问题，见[#4604](https://github.com/wailsapp/wails/issues/4604)，贡献者：[@JackDoan](https://github.com/JackDoan)
- 修复在 Windows 上创建隐藏窗口时出现白色窗口的问题，贡献者：@leaanthony，见[#4612](https://github.com/wailsapp/wails/pull/4612)
- 修复文档中的 notifications 软件包导入路径，贡献者：@rxliuli，见[#4617](https://github.com/wailsapp/wails/pull/4617)
- 修复使用 npm 包 @wailsio/runtime 时拖放功能不工作的问题（#4489），由 @leaanthony 在 #4616 中修复
- Windows：修复启动时窗口闪烁以及隐藏窗口错误显示的问题，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4600)中修复。
- 修复 Wayland 窗口最大化时的尺寸问题（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)修复
- 修复 Wayland 窗口最大化时的尺寸问题（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)修复
- 修复液态玻璃演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中修复
- 修复 AssetServer 在 macOS 上可能崩溃的问题，由 @jghiloni 在[#4576](https://github.com/wailsapp/wails/pull/4576)中修复
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修复
- 修复夜间版本发布流水线，由 @riadafridishibly 在[#4597](https://github.com/wailsapp/wails/pull/4597)中修复
- 修复液态玻璃演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中修复
- 修复 AssetServer 在 macOS 上可能崩溃的问题，由 @jghiloni 在[#4576](https://github.com/wailsapp/wails/pull/4576)中修复
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修复
- 修复夜间版本发布流水线，由 @riadafridishibly 在[#4597](https://github.com/wailsapp/wails/pull/4597)中修复
- 修复液态玻璃演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中修复
- 修复 Windows 上的 SetBackgroundColour，由 @PPTGamer 在[PR](https://github.com/wailsapp/wails/pull/4492)中修复
- 更新文档以反映 Manager API 重构带来的变更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 修复 Linux taskfile 中 .desktop 文件的 appicon 变量，见[PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 更新文档以反映 Manager API 重构带来的变更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 由 @leaanthony 在[#4460](https://github.com/wailsapp/wails/pull/4460)中修复[#4456](https://github.com/wailsapp/wails/issues/4456)中报告的 Windows 空指针解引用错误
- 在 macOS WKWebView 中添加对`allowsBackForwardNavigationGestures`的支持，以启用双指轻扫导航手势（#1857）
- 修复初始设置为禁用的菜单项无法使用 onClick 的问题，由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中修复。感谢 @IanVS 开展初步调查。
- 修复构建失败时未清理 Vite 服务器的问题（#4403）
- 修复在 Windows 上关闭或取消`SaveFileDialog`时发生 panic 的问题。由 @hkhere 在[PR](https://github.com/wailsapp/wails/pull/4284)中修复
- 修复 Windows 上 HTML 层级的拖放功能，由[@mbaklor](https://github.com/mbaklor)在[#4259](https://github.com/wailsapp/wails/pull/4259)中修复
- 在 macOS WKWebView 中添加对`allowsBackForwardNavigationGestures`的支持，以启用双指轻扫导航手势（#1857）
- 修复初始设置为禁用的菜单项无法使用 onClick 的问题，由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中修复。感谢 @IanVS 开展初步调查。
- 修复构建失败时未清理 Vite 服务器的问题（#4403）
- 修复 Windows 上的通知解析问题，由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4450)中修复
- 修复 doctor 命令，使其检查 Windows SDK 依赖项，由[@kodumulo](https://github.com/kodumulo)在[#4390](https://github.com/wailsapp/wails/issues/4390)中修复
- 修复 Mac 上 processURLRequest 中的空指针解引用问题，由[@etesam913](https://github.com/etesam913)在[#4366](https://github.com/wailsapp/wails/pull/4366)中修复
- 修复 Linux 上导致筛选对话框无法使用的错误，由[@bh90210](https://github.com/bh90210)在[#4287](https://github.com/wailsapp/wails/pull/4287)中修复
- 修复 Windows 和 Linux 上的“编辑”菜单问题，由[@leaanthony](https://github.com/leaanthony)在[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)中修复
- 将 macOS .plist 文件中的最低系统版本从10.13.0更新为10.15.0，由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中完成
- 修复窗口 ID 跳号问题，由[@leaanthony](https://github.com/leaanthony)修复
- 修复调用 RegisterContextMenu 时菜单为 nil 的问题，由[@leaanthony](https://github.com/leaanthony)修复
- 修复绑定生成器输出中的依赖循环，由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修复
- 修复绑定生成器输出中的先使用后定义错误，由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修复
- 将构建标志传递给绑定生成器，由[@fbbdev](https://github.com/fbbdev)在[#4023](https://github.com/wailsapp/wails/pull/4023)中完成
- 将 Windows Taskfile 中的路径改为使用正斜杠，以确保其可在非 Windows 平台上运行，由[@leaanthony](https://github.com/leaanthony)完成
- 现已修复 Mac 和 Mac JS 事件问题，由[@leaanthony](https://github.com/leaanthony)修复
- 修复 macOS 上的事件死锁问题，由[@leaanthony](https://github.com/leaanthony)修复
- 修复 Windows 上初始化 Window 时，提供了 HTML 但未提供 JS 所导致的`Parameter incorrect`错误，由[@leaanthony](https://github.com/leaanthony)修复
- 修复资源服务器中用于探测内容类型的响应前缀大小，由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修复
- 修复资源服务器对根索引路径上非404响应的处理，由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修复
- 修复绑定生成器测试泛型类型的属性时出现的未定义行为，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复
- 修复底层类型与命名包装器的属性不同时，绑定生成器针对模型生成的输出，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复
- 修复绑定生成器针对映射键类型和预处理生成的输出，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复
- 修复绑定生成器针对实现编组器接口的结构体生成的输出，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复
- 修复绑定生成器对涉及泛型类型的类型循环的检测，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复
- 修复了绑定生成器输出中对未导出模型的无效引用，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 将注入的代码移至服务文件末尾，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器对文件关闭操作所产生错误的处理，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 对于定义了生命周期方法或 HTTP 方法、但未定义其他绑定方法的服务，不再显示警告，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了非 React 模板在使用浅色系统配色方案时无法显示 Hello World 页脚的问题，由[@marcus-crane](https://github.com/marcus-crane)在[#4056](https://github.com/wailsapp/wails/pull/4056)中完成
- 修复了 macOS 上菜单项被隐藏的问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了消息处理器中的错误处理和格式化，由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
-  修复了退出应用程序时跳过服务关闭的问题，由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
-  确保菜单更新在主线程上执行，由[@leaanthony](https://github.com/leaanthony)完成
- 拖动和调整大小机制现已更加稳健，也更符合平台的预期行为，由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 修复了[#4097](https://github.com/wailsapp/wails/issues/4097) 中 Webpack/angular 丢弃运行时初始化代码的问题，由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 修复了初始状态为隐藏的菜单项，由[@IanVS](https://github.com/IanVS)在[#4116](https://github.com/wailsapp/wails/pull/4116)中完成
- 修复了以下问题：当无扩展名请求中的`[request]`不存在但`[request].html`存在时，assetFileServer 不提供`.html`文件
- 修复了图标生成路径，由[@robin-samuel](https://github.com/robin-samuel)在[#4125](https://github.com/wailsapp/wails/pull/4125)中完成
- 修复了未发出`fullscreen`、`unfullscreen`、`unminimise`和`unmaximise`事件的问题，由[@oSethoum](https://github.com/osethoum)在[#4130](https://github.com/wailsapp/wails/pull/4130)中完成
- 修复了配置中默认版本的前缀不正确所导致的 NSIS 错误，由[@robin-samuel](https://github.com/robin-samuel)在[#4126](https://github.com/wailsapp/wails/pull/4126)中完成
- 修复了 Windows 上 Dialogs 运行时函数返回转义后路径的问题，由[TheGB0077](https://github.com/TheGB0077)在[#4188](https://github.com/wailsapp/wails/pull/4188)中完成
- 修复了 HKCU 中的 Webview2 检测路径，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 macOS 上的输入问题，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 Windows 图标生成任务的文件名，由[@yulesxoxo](https://github.com/yulesxoxo)在[#4219](https://github.com/wailsapp/wails/pull/4219)中完成。
- 基于 @kron 的工作，修复了无边框窗口的透明度问题，由[@leaanthony](https://github.com/leaanthony)完成。
- 基于 @kron 的工作，修复了窗口被禁用或最小化时的聚焦调用，由[@leaanthony](https://github.com/leaanthony)完成。
- 基于 @kron 的工作，修复了任务栏重启后系统托盘不显示的问题，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 fallbackResponseWriter 未实现 Flush() 的问题，见[#4245](https://github.com/wailsapp/wails/pull/4245)
- 修复了 fallbackResponseWriter 未实现 Flush() 的问题，由 [@superDingda] 在[#4236](https://github.com/wailsapp/wails/issues/4236)中完成
- 修复了 macOS 窗口在异步 Go 绑定函数调用仍待处理时关闭会导致崩溃的问题，由[@joshhardy](https://github.com/joshhardy)在[#4354](https://github.com/wailsapp/wails/pull/4354)中完成
- 修复了 Windows 效率模式启动时的竞态条件，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了 Windows 图标句柄的清理，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 Windows 上的`OpenFileManager`，由[@PPTGamer](https://github.com/PPTGamer)在[#4375](https://github.com/wailsapp/wails/pull/4375)中完成。
- 修复了 Linux 的最小/最大宽度选项，由 @atterpac 在[#3979](https://github.com/wailsapp/wails/pull/3979)中完成
- 通过提升 npm 版本修复了 TypeScript 模板的类型定义，由 @atterpac 在[#3966](https://github.com/wailsapp/wails/pull/3966)中完成
- 修复了 SvelteKit 模板中的 CSS 引用，由 @atterpac 在[#3945](https://github.com/wailsapp/wails/pull/3945)中完成
- 确保 window run() 中的关键回调在主线程上调用，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了对话框目录选择器示例，由[@leaanthony](https://github.com/leaanthony)完成
- 新增了 index.html 缺失时显示的中文错误页面，由[@leaanthony](https://github.com/leaanthony)完成
-  确保`windowDidBecomeKey`回调在主线程上运行，由[@leaanthony](https://github.com/leaanthony)完成
-  支持无边框窗口全屏显示，由[@leaanthony](https://github.com/leaanthony)完成
-  改进了窗口销毁逻辑，由[@leaanthony](https://github.com/leaanthony)完成
-  修复了窗口附加到系统托盘时的位置逻辑，由[@leaanthony](https://github.com/leaanthony)完成
-  支持无边框窗口全屏显示，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了事件处理，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了窗口关闭逻辑，由[@leaanthony](https://github.com/leaanthony)完成
- 通用 taskfile 现在默认会为 TypeScript 模板生成 TypeScript 绑定，由[@leaanthony](https://github.com/leaanthony)完成
- 修复没有打开任何窗口或仅有系统托盘时，收到 WM_CLOSE 消息后的应用程序关闭行为，由[@mmalcek](https://github.com/mmalcek)在[#3990](https://github.com/wailsapp/wails/pull/3990)中完成
- 修复了 garble 构建，由 @5aaee9 在[#3192](https://github.com/wailsapp/wails/pull/3192)中完成
- 修复了 Windows NSIS 构建，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了 Linux 多选对话框中因未关闭而导致的死锁
- 修复了 Windows 构建期间对 .syso 文件的跨平台清理，由
- 修复了 amd64 AppImage 编译，由 @atterpac 在
- 修复了构建资源更新，由 @ansxuman 在
- 修复了 Linux 系统托盘的`OnClick`和`OnRightClick`实现，作者：@atterpac
- 修复了`AlwaysOnTop`在 Mac 上无法工作的问题，作者：
-  修复了`application.NewEditMenu`包含重复项的问题
- 🐧 修复了 aarch64 编译问题
- ⊞ 修复了单选组菜单项，作者：
- 修复在 MacOS 上构建可运行的 .app 时出现的错误（触发条件：当 'name' 和 'outputfilename'……）
- 修复了拖放示例中使用 customEventProcessor 时的错误，作者：
- 🐧 修复了添加 IgnoreMouseEvents 后引入的 Linux 编译错误，作者：
- ⊞ 修复了 syso 图标文件生成错误，作者：
- 🐧 合入了来自以下来源的修复，使其可在 Wayland 中原生运行：
- 不要在以下位置绑定内部服务方法：
- ⊞ 修复了以下位置的系统托盘启动时 panic：
- 不要在以下位置绑定内部服务方法：
- ⊞ 修复了以下位置的系统托盘启动时 panic：
- 对菜单项和事件处理进行了重大重构。目前主要改进了 macOS。作者：
- 修复了以下位置在插件和事件重构后的测试：
- ⊞ 修复了`Failed to unregister class Chrome_WidgetWin_0`警告。作者：
- 模块问题
- 由[atterpac](https://github.com/atterpac)修复了以下位置的调整大小事件消息传递：
- 🐧修复了 NixOS 上的主题处理错误，作者：
- 修复了 Windows 上跨卷安装项目的问题，作者：
- 修复了 React 模板的 CSS，使其能够显示页脚，作者：
- 通过更新至最新的 refresh，修复了在开发模式下工作时出现的僵尸进程
- 修复了 AppImage 的 WebKit 文件来源问题，作者：[Atterpac](https://github.com/atterpac)
- 由[Atterpac](https://github.com/Atterpac)修复了以下位置的 Doctor apt 软件包验证：
- 修复了应用程序退出时冻结的问题（Darwin），由 @5aaee9 在以下位置完成：
- 修复了 Windows 上示例的背景颜色，作者：
- 由[mmghv](https://github.com/mmghv)修复了以下位置的默认上下文菜单：
- 修复了 Darwin 上方向键的十六进制值，作者：
- 使 Windows 上的拖放功能恢复正常。添加者：
- 修复了用户没有正确驱动程序时 Doctor 在 Linux 上出现的错误
- 修复了启动时的 DPI 缩放问题（Windows）。由[@almas-x](https://github.com/almas-x)在以下位置更改：
- 修复`go.mod`中的替换行以使用相对路径。修复包含以下内容的 Windows 路径：
- 修复了 MacOS 系统托盘未关联窗口时的点击处理，作者：
- 修复了未知选项导致 Windows 构建失败的问题，作者：
- 修复了 Windows 上未配置以下内容时左键单击系统托盘图标导致的崩溃：
- 修复了两次打开窗口时 baseURL 错误的问题，由 @5aaee9 在 PR 中完成：
- 修复了`WebviewWindow.Restore`方法中 if 分支的顺序，作者：
- 当出现以下情况时，跨多次`GetStartURL`调用正确计算`startURL`：
- 修复`Screen`结构体的 JS 类型，使其与对应的 Go 类型一致，作者：
- 修复`WML.Reload`方法，以确保正确清理已注册的事件
- 修复了 Linux 上自定义上下文菜单立即关闭的问题，作者：
- 修复绑定所生成模型文件的输出路径和扩展名
- 修复绑定所生成 JS 代码中模型文件的导入路径
- 修复了某些 Linux 发行版上的拖放功能，作者：
- 修复了使用`wails3 task dev`时缺少 macOS 任务的问题，作者：
- 修复了注册事件导致向 nil map 赋值的问题，作者：
- 修复了绑定方法参数的反序列化，作者：
- 修复了对绑定方法多个返回值的处理，作者：
- 修复了 Doctor 对并非通过系统软件包管理器安装的 npm 的检测
- 修复了缺少 MicrosoftEdgeWebview2Setup.exe 的问题。感谢：
- 修复了窗口 ID 处理导致 Linux 上随机崩溃的问题，作者：@leaanthony。基于：
- 修复了 systemTray.setIcon 在 Linux 上导致崩溃的问题，作者：
- 修复了以下平台上首次调用`setFrameless`函数时未能应用窗口边框的问题：

### 变更

- **破坏性变更**：生成的 JS/TS 绑定中的 map 键现在标记为可选，以准确反映 Go map 的语义。TypeScript 中访问 map 值现在返回`T | undefined`而非`T`，因此需要执行 null 检查或使用类型断言（#4943），作者：`@fbbdev`
- 根据`@wailsio/runtime`中的变更，将`Event`改为`Events`，并在`Features/Events/Event System`的文档中采用适当的函数调用，作者：@AbdelhadiSeddar
- 将`EnabledFeatures`、`DisabledFeatures`和`AdditionalBrowserArgs`从每窗口选项移至应用程序级别的`Options.Windows`（#4559），作者：@leaanthony
- 更新`Drag N Drop`示例的 README，并强调该示例演示了`Internal Drag and Drop`，作者：@ndianabasi
- 将多处调试日志的级别从 Info 改为 Debug（作者：@mbaklor）
- <strong>破坏性变更：</strong>在窗口选项中将`EnableDragAndDrop`重命名为`EnableFileDrop`
- <strong>破坏性变更：</strong>在事件上下文中将`DropZoneDetails`重命名为`DropTargetDetails`
- <strong>破坏性变更：</strong>将`WindowEventContext`上的`DropZoneDetails()`方法重命名为`DropTargetDetails()`
- <strong>破坏性变更：</strong>移除`WindowDropZoneFilesDropped`事件，改用`WindowFilesDropped`
- <strong>破坏性变更：</strong>将 HTML 属性从`data-wails-dropzone`更改为`data-file-drop-target`
- <strong>破坏性变更：</strong>将 CSS 悬停类从`wails-dropzone-hover`更改为`file-drop-target-active`
- <strong>破坏性变更：</strong>从 Windows 中移除`DragEffect`、`OnEnterEffect`和`OnOverEffect`选项（它们曾是已移除的 IDropTarget 的一部分）
- 所有运行时 JSON 处理（方法绑定、事件、WebView 请求、通知和 kvstore）改用 goccy/go-json，性能提升21-63%，内存分配减少40-60%
- 优化 BoundMethod 结构体布局并缓存 isVariadic 标志，以减少每次调用的开销
- 对于具有`<=8`个参数的方法，使用在栈上分配的参数缓冲区，以避免堆分配
- 优化方法调用中的结果收集，避免为单个返回值分配切片
- MIME 类型缓存改用 sync.Map，以提升并发性能
- 读取 HTTP 传输请求正文时使用缓冲池
- 在内容类型探测器中延迟分配 CloseNotify 通道，以减少每个请求的内存分配
- 从资源服务器中移除 CSS 调试日志
- 扩展 MIME 类型扩展名映射，以涵盖50多种常见 Web 格式（字体、音频、视频等）
- 更新 Window `X/Y`选项的文档 @ruhuang2001
- 更新`Frontend Runtime`文档，添加更多用于生成前端绑定的选项，作者：@ndianabasi
- 更新 Wails v3资源服务器文档页面，作者：@ndianabasi
- **破坏性变更**：移除包级对话框函数（`application.InfoDialog()`、`application.QuestionDialog()`等）。请改用`app.Dialog`管理器：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新对话框文档以匹配实际 API：使用`app.Dialog.*`、带回调的`AddButton()`（而非`SetButtons()`）、`SetDefaultButton(*Button)`（而非字符串）、`AddFilter()`（而非`SetFilters()`）、`SetFilename()`（而非`SetDefaultFilename()`），并使用`app.Dialog.OpenFile().CanChooseDirectories(true)`选择文件夹
- **破坏性变更**：现在默认使用生产构建。要创建开发构建，请在 Taskfile 中设置`DEV=true`。请生成一个新项目以查看配置示例，作者：@leaanthony
- 发出包含零个或一个数据参数的自定义事件时，数据值将直接赋给 Data 字段，而不会封装到切片中；由[@fbbdev](https://github.com/fbbdev)在[#4633](https://github.com/wailsapp/wails/pull/4633)中完成
- Windows 托盘现在会通过切换`NIS_HIDDEN`来遵循`SystemTray.Show()`/`Hide()`，使应用能够真正消失并重新出现（#4653）。
- 注册托盘时会复用已解析的图标，仅设置一次`NOTIFYICON_VERSION_4`，并启用`NIF_SHOWTIP`，以便工具提示在 Explorer 重启后恢复（#4653）。
- macOS：窗口居中时使用`visibleFrame`而非`frame`，以排除菜单栏和程序坞区域
- macOS：窗口居中时使用`visibleFrame`而非`frame`，以排除菜单栏和程序坞区域
- 使用`-config`参数运行`wails3 update build-assets`时，通过`-product*`参数设置的值会
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，作者：@leaanthony，见[#4471](https://github.com/wailsapp/wails/pull/4471)
- 重构内部窗口处理，作者：@leaanthony，见[#4471](https://github.com/wailsapp/wails/pull/4471)
- 移除`application.WindowIDKey`和`application.WindowNameKey`（由`application.WindowKey`取代），作者：[@leaanthony](https://github.com/leaanthony)
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，固定长度数组类型的类字段现在会按预期长度初始化，而不再初始化为空；由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中完成
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- `application.NewService`不再接受 options 作为可选参数（请改用`application.NewServiceWithOptions`）；由[@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中完成
- 移除`nanoid`依赖，作者：[@leaanthony](https://github.com/leaanthony)
- 更新 Window 示例，以展示 mica/acrylic/tabbed 窗口样式，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，已移除`internal.js/ts`模型文件；现在可在`models.js/ts`中找到所有模型；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 在 JS/TS 绑定中，命名类型绝不会渲染为其他命名类型的别名；旧行为现在仅适用于别名；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 在 JS/TS 绑定的类模式下，类型为类型参数的结构体字段会被标记为可选，并且绝不会自动初始化；由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 从模板中移除 ESLint，作者：[@IanVS](https://github.com/IanVS)，见[#4059](https://github.com/wailsapp/wails/pull/4059)
- 将版权日期更新为2025，作者：[@IanVS](https://github.com/IanVS)，见[#4037](https://github.com/wailsapp/wails/pull/4037)
- 添加 event.Sender 文档，作者：[@IanVS](https://github.com/IanVS)，见[#4075](https://github.com/wailsapp/wails/pull/4075)
- 支持 Go 1.24，作者：[@leaanthony](https://github.com/leaanthony)
- `ServiceStartup`钩子现在会在调用`App.Run`时触发，而非在`application.New`中触发；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- `ServiceStartup`错误现在会从`App.Run`返回，而不会终止进程；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 来自 JS 的绑定调用和对话框调用现在会以错误对象而非字符串拒绝；由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
- 改进 Windows 上系统托盘菜单的定位，作者：[@leaanthony](https://github.com/leaanthony)
- JS 运行时已移植到 TypeScript；由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 运行时一经导入便会初始化，无需等待窗口加载，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4100](https://github.com/wailsapp/wails/pull/4100)
- 运行时不再导出 init 方法。可以通过仅为副作用而进行的导入来初始化运行时，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4100](https://github.com/wailsapp/wails/pull/4100)
- 绑定方法现在返回一个`CancellablePromise`；如果调用被取消，它将以`CancelError`拒绝。调用的实际结果会被丢弃，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4100](https://github.com/wailsapp/wails/pull/4100)
- 内置服务类型现在统一命名为`Service`，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4067](https://github.com/wailsapp/wails/pull/4067)
- 带选项的内置服务创建函数现在统一命名为`NewWithConfig`，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4067](https://github.com/wailsapp/wails/pull/4067)
- 为与 Go API 保持一致，`sqlite`服务的`Select`方法现已更名为`Query`，贡献者：[@fbbdev](https://github.com/fbbdev)，见[#4067](https://github.com/wailsapp/wails/pull/4067)
- 模板：将运行时移至 "dependencies"，并整理 package.json 文件，贡献者：[@IanVS](https://github.com/IanVS)，见[#4133](https://github.com/wailsapp/wails/pull/4133)
- 在开发环境中创建应用包并进行临时签名，以启用某些 macOS API，贡献者：[@popaprozac](https://github.com/popaprozac)，见[#4171](https://github.com/wailsapp/wails/pull/4171)
- 将构建资源移至各平台专用目录，贡献者：[@leaanthony](https://github.com/leaanthony)
- 将 Taskfile 移至各平台专用目录并重命名，贡献者：[@leaanthony](https://github.com/leaanthony)
- 显著改善缺少`index.html`时的使用体验，贡献者：[@leaanthony](https://github.com/leaanthony)
- [Windows] 提升最小化和还原操作的性能，贡献者：[@leaanthony](https://github.com/leaanthony)。基于[562589540](https://github.com/562589540)最初提交的[PR](https://github.com/wailsapp/wails/pull/3955)
- 移除`ShouldClose`选项（请改为给 events.Common.WindowClosing 注册钩子），贡献者：[@leaanthony](https://github.com/leaanthony)
- [Windows] 减少打开窗口时的闪烁，贡献者：[@leaanthony](https://github.com/leaanthony)
- 移除`Window.Destroy`，因为它原本应是内部函数，贡献者：[@leaanthony](https://github.com/leaanthony)
- 将`WindowClose`事件重命名为`WindowClosing`，贡献者：[@leaanthony](https://github.com/leaanthony)
- 前端构建现在会根据构建类型使用 Vite 环境 "development" 或 "production"，贡献者：[@leaanthony](https://github.com/leaanthony)
- 更新至 go-webview2 v1.19，贡献者：[@leaanthony](https://github.com/leaanthony)
- 确保使用 Taskfile 的分支版本，贡献者：@leaanthony
- 更新 Taskfile 的分支版本，以修复使用以下方式安装时的版本问题
- 使用 Taskfile 的分支版本，以修复使用以下方式安装时的版本问题
- `service.OnStartup`现在会在发生错误时关闭应用程序并运行
- 重构系统托盘点击消息传递，使其更符合用户交互，贡献者：
- 嵌入资源时包含`all:frontend/dist`，以支持会生成以下内容的框架
- Taskfile 重构，贡献者：[leaanthony](https://github.com/leaanthony)，见
- 升级至`go-webview2` v1.0.16，贡献者：
- 修正`Screen`类型，使其包含`ID`而非`Id`，贡献者：
- 更新`go.mod.tmpl`中的 Wails 版本以支持`application.ServiceOptions`，贡献者：
- 修复服务名称的确定逻辑，贡献者：[windom](https://github.com/windom/)，见
- mkdocs serve 现在使用 Docker，贡献者：[leaanthony](https://github.com/leaanthony)
- 将开发配置整合到`config.yml`中，贡献者：
- 系统托盘对话框现在默认使用应用程序图标（如果可用，仅限 Windows），贡献者：
- 改进 macOS 上的 GPU 和内存信息报告，贡献者：
- 移除`WebviewGpuIsDisabled`和`EnableFraudulentWebsiteWarnings`
- Events API 变更：`On`/`Emit` -> 用户事件，`OnApplicationEvent` ->
- 修复 Linux 上的 Events API，贡献者：[TheGB0077](https://github.com/TheGB0077)，见
- [CI] 改进 Actions，并允许在分支仓库中也运行 Actions，以及
- 将`AbsolutePosition()`重命名为`Position()`，贡献者：
- 将 Linux WebKit 依赖从 webkitgtk2-4.0更新为 webkit2gtk-4.1，以
- 内置的 JS 运行时脚本现在是 ESM 模块：导入该脚本的 script 标签
- `@wailsio/runtime`包不再将其 API 发布到`window.wails`上
- Window API 模块`@wailsio/runtime/src/window`现在会公开其所属的
- JS 窗口 API 已更新，以与当前的 Go`WebviewWindow`保持一致
- 绑定生成器现在默认按 ID 调用。`-id` CLI 选项
- 新的绑定代码布局：输出文件此前按文件夹组织
- 结构体字段`application.Options.Bind`已重命名为
- 绑定服务的新语法：现在必须将服务实例包装在一个
- 在非终端或 CI 环境中禁用加载指示器，贡献者：

### 已移除

- **破坏性变更**：从每个窗口的`WindowsWindow`选项中移除`EnabledFeatures`、`DisabledFeatures`和`AdditionalLaunchArgs`。请改用应用程序级别的`Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`和`Options.Windows.AdditionalBrowserArgs`。这些标志会全局应用于共享的 WebView2 环境（#4559），贡献者：@leaanthony
- 移除 Windows 上`IDropTarget`的原生实现，改用基于 JavaScript 的方案（与 v2 的行为一致）
- 移除 github.com/wailsapp/mimetype 依赖，改用扩展后的扩展名映射和标准库 http.DetectContentType，将二进制文件大小减少约1.2MB
- 通过为 Linux 文件资源管理器实现精简的 .desktop 文件解析器，移除 gopkg.in/ini.v1 依赖，节省约 45KB
- 使用 Go 1.21+ 标准库的 slices 包和少量内部辅助函数，从运行时代码中移除 samber/lo，节省约 310KB
- 从 Darwin URL scheme 处理程序中移除调试 printf 语句（#4834）
- **破坏性变更**：移除 `linux:WindowLoadChanged` 事件；请改用 `linux:WindowLoadFinished` 检测 WebView 何时完成加载（#3896），由 @leaanthony 完成

### 破坏性变更

- **管理器 API 重构**：将应用程序 API 从扁平结构重组为按管理器组织的结构，以改善代码组织并提高 API 的可发现性；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- 重命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`；由 [@leaanthony](https://github.com/leaanthony) 完成
- 将 `Path` 和 `Paths` 方法移至 `application` 包；由 [@leaanthony](https://github.com/leaanthony) 完成
- 应用程序菜单现在仅适用于 macOS；由 [@leaanthony](https://github.com/leaanthony) 完成

## v3.0.0-alpha.78 - 2026-04-21

## 新增

## 修复

## v3.0.0-alpha.77 - 2026-04-18

## 修复

## v3.0.0-alpha.76 - 2026-04-17

## 修复

## v3.0.0-alpha.75 - 2026-04-16

## 修复

## v3.0.0-alpha.74 - 2026-03-01

## 新增

## 修复

## v3.0.0-alpha.73 - 2026-02-27

## 修复

## v3.0.0-alpha.72 - 2026-02-16

## 修复

## v3.0.0-alpha.71 - 2026-02-10

## 新增

## 修复

## v3.0.0-alpha.70 - 2026-02-09

## 新增

## 修复

## v3.0.0-alpha.69 - 2026-02-08

## 新增

## 修复

## v3.0.0-alpha.68 - 2026-02-07

## 新增

## 变更

## 修复

## v3.0.0-alpha.67 - 2026-02-04

## 新增

## 变更

## 修复

## v3.0.0-alpha.66 - 2026-02-03

## 新增

## 变更

## 修复

## 移除

## v3.0.0-alpha.65 - 2026-02-01

## 新增

## v3.0.0-alpha.64 - 2026-01-26

## 新增

## v3.0.0-alpha.63 - 2026-01-25

## 修复

## v3.0.0-alpha.62 - 2026-01-22

## 修复

## v3.0.0-alpha.61 - 2026-01-20

## 修复

## v3.0.0-alpha.60 - 2026-01-14

## 修复

## v3.0.0-alpha.59 - 2026-01-11

## 变更

## v3.0.0-alpha.58 - 2026-01-09

## 修复

## v3.0.0-alpha.57 - 2026-01-05

## 变更

## 修复

## v3.0.0-alpha.56 - 2026-01-04

## 新增

## 变更

## 修复

## 移除

## v3.0.0-alpha.55 - 2026-01-02

## 变更

## 修复

## 移除

## v3.0.0-alpha.54 - 2025-12-29

## 新增

## 修复

## 移除

## v3.0.0-alpha.53 - 2025-12-27

## 新增

## 修复

## v3.0.0-alpha.52 - 2025-12-26

## 修复

## v3.0.0-alpha.51 - 2025-12-23

## 修复

## v3.0.0-alpha.50 - 2025-12-21

## 变更

## v3.0.0-alpha.49 - 2025-12-18

## 变更

## v3.0.0-alpha.48 - 2025-12-16

## 新增

## 变更

## 修复

## v3.0.0-alpha.47 - 2025-12-15

## 新增

## 修复

## v3.0.0-alpha.46 - 2025-12-14

## 新增

## 移除

## v3.0.0-alpha.45 - 2025-12-13

## 新增

## 修复

## v3.0.0-alpha.44 - 2025-12-12

## 新增

## 变更

## 修复

## v3.0.0-alpha.43 - 2025-12-11

## 新增

## v3.0.0-alpha.42 - 2025-12-10

## 新增

## v3.0.0-alpha.41 - 2025-11-23

## 修复

## v3.0.0-alpha.40 - 2025-11-13

## 修复

## v3.0.0-alpha.39 - 2025-11-12

## 新增

## 变更

## v3.0.0-alpha.38 - 2025-11-04

## 新增

## 变更

## 修复

## v3.0.0-alpha.37 - 2025-11-02

## 修复

## v3.0.0-alpha.36 - 2025-10-15

## 修复

## v3.0.0-alpha.35 - 2025-10-14

## 修复

## v3.0.0-alpha.34 - 2025-10-06

## 新增

## 修复

## v3.0.0-alpha.33 - 2025-10-04

## 修复

## v3.0.0-alpha.32 - 2025-10-02

## 修复

## v3.0.0-alpha.31 - 2025-09-27

## 修复

## v3.0.0-alpha.30 - 2025-09-26

## 修复

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 变更

## 修复

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 变更

## 修复

## v3.0.0-alpha.27 - 2025-09-07

## 修复

## v3.0.0-alpha.26 - 2025-08-24

## 新增

## v3.0.0-alpha.25 - 2025-08-16

## 变更

不再被忽略，并会覆盖配置值。

## v3.0.0-alpha.24 - 2025-08-13

## 新增

## v3.0.0-alpha.23 - 2025-08-11

## 修复

## v3.0.0-alpha.22 - 2025-08-10

## 新增

## 变更

+ 修复范围过宽的 Linux 软件包依赖项，并修复已过时的 RPM 依赖项。

## v3.0.0-alpha.21 - 2025-08-07

## 修复

## v3.0.0-alpha.20 - 2025-08-06

## 修复

## v3.0.0-alpha.19 - 2025-08-05

## 新增

## 修复

## v3.0.0-alpha.18 - 2025-08-03

## 新增

## 修复

## v3.0.0-alpha.17 - 2025-07-31

## 修复

## v3.0.0-alpha.16 - 2025-07-25

## 新增

## v3.0.0-alpha.15 - 2025-07-25

## 新增

## v3.0.0-alpha.14 - 2025-07-25

## 新增

## v3.0.0-alpha.12 - 2025-07-15

### 新增

### 修复

## v3.0.0-alpha.11 - 2025-07-12

## 新增

## v3.0.0-alpha.10 - 2025-07-06

### 破坏性变更

### 新增

### 修复

### 变更

## v3.0.0-alpha.9 - 2025-01-13

### 新增

### 修复

### 变更

## v3.0.0-alpha.8.3 - 2024-12-07

### 变更

## v3.0.0-alpha.8.2 - 2024-12-07

### 变更

`go install`，由 @leaanthony 完成

## v3.0.0-alpha.8.1 - 2024-12-07

### 变更

`go install`，由 @leaanthony 完成

## v3.0.0-alpha.8 - 2024-12-06

### 新增

@atterpac 于 [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) 于   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) 于   [#3867](https://github.com/wailsapp/wails/pull/3867)   由 [atterpac](https://github.com/atterpac) 于   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) 于   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   由 [leaanthony](https://github.com/leaanthony) 于   [#3885](https://github.com/wailsapp/wails/pull/3885)将其定位在给定的 X/Y 位置   [ansxuman](https://github.com/ansxuman) 和   [leaanthony](https://github.com/leaanthony) 于   [#3823](https://github.com/wailsapp/wails/pull/3823)   由 [leaanthony](https://github.com/leaanthony) 于   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) 于   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony)。 -

### 变更

`service.OnShutdown`适用于此前已启动的任何服务，由 @atterpac 在以下变更中完成：   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac 于[#3907](https://github.com/wailsapp/wails/pull/3907)   子文件夹，由 @atterpac 于   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)于   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes)于   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （已由`EnabledFeatures`和`DisabledFeatures`选项取代），作者：   [leaanthony](https://github.com/leaanthony)

### 修复

通道变量，由 @michael-freling 于   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) 于   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   于 [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) 于   [#3841](https://github.com/wailsapp/wails/pull/3841)   Darwin 编辑菜单中的 `PasteAndMatchStyle` 角色，由   [johnmccabe](https://github.com/johnmccabe) 于   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) 于   [#3854](https://github.com/wailsapp/wails/pull/3854)，作者：   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   不同。由 @nickisworking 于   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### 新增

[mmghv](https://github.com/mmghv) 于   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) 和   [leaanthony](https://github.com/leaanthony) 于   [#3570](https://github.com/wailsapp/wails/pull/3570)

### 变更

应用程序事件 `OnWindowEvent` -> 窗口事件，作者：   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   以 `v3/` 或 `v3-` 为前缀的分支，由   [stendler](https://github.com/stendler) 于   [#3747](https://github.com/wailsapp/wails/pull/3747)

### 修复

[etesam913](https://github.com/etesam913) 于   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) 于   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) 于   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) 于   [#3614](https://github.com/wailsapp/wails/pull/3614)，作者：   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720)，作者：   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，作者：   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720)，作者：   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693)，作者：   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)，作者：   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修复

## v3.0.0-alpha.5 - 2024-07-30

### 新增

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   图标点击，由 @5aaee9 于 [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) 于   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   于[#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) 于   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) 于   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   基于 @Mai-Lapyst 的 [PR](https://github.com/wailsapp/wails/pull/2044)   [@fbbdev](https://github.com/fbbdev) 于   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) 于   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) 于   [#3295](https://github.com/wailsapp/wails/pull/3295)   npm 软件包，由 [@fbbdev](https://github.com/fbbdev) 于   [#3334](https://github.com/wailsapp/wails/pull/3334)   于 [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT`，由 [@abichinger](https://github.com/abichinger) 于   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) 于   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) 于   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) 于   [#3674](https://github.com/wailsapp/wails/pull/3674)

### 修复

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) 在   [#3515](https://github.com/wailsapp/wails/pull/3515) 中   [atterpac](https://github.com/atterac) 在   [#3512](https://github.com/wailsapp/wails/pull/3512) 中   [atterpac](https://github.com/atterpac) 在   [#3477](https://github.com/wailsapp/wails/pull/3477) 中   由[Atterpac](https://github.com/atterpac)在   [#3320](https://github.com/wailsapp/wails/pull/3320) 中完成。   位于[#3306](https://github.com/wailsapp/wails/pull/3306)中。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) 在   [#2750](https://github.com/wailsapp/wails/pull/2750) 中。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave) 在   [#3052](https://github.com/wailsapp/wails/pull/3052) 中。   [@pylotlight](https://github.com/pylotlight) 在   [PR](https://github.com/wailsapp/wails/pull/3039) 中   安装。由[@pylotlight](https://github.com/pylotlight)在   [PR](https://github.com/wailsapp/wails/pull/3032) 中添加   [PR](https://github.com/wailsapp/wails/pull/3145)   空格——@leaanthony。   [thomas-senechal](https://github.com/thomas-senechal) 在 PR   [#3207](https://github.com/wailsapp/wails/pull/3207) 中   [thomas-senechal](https://github.com/thomas-senechal) 在 PR   [#3208](https://github.com/wailsapp/wails/pull/3208) 中   附加窗口[tw1nk](https://github.com/tw1nk)在 PR   [#3271](https://github.com/wailsapp/wails/pull/3271) 中   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) 在   [#3279](https://github.com/wailsapp/wails/pull/3279) 中   存在`FRONTEND_DEVSERVER_URL`。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) 在   [#3295](https://github.com/wailsapp/wails/pull/3295) 中   监听器由[@fbbdev](https://github.com/fbbdev)在   [#3295](https://github.com/wailsapp/wails/pull/3295) 中完成   [@abichinger](https://github.com/abichinger) 在   [#3330](https://github.com/wailsapp/wails/pull/3330) 中   生成器由[@fbbdev](https://github.com/fbbdev)在   [#3334](https://github.com/wailsapp/wails/pull/3334) 中完成   生成器由[@fbbdev](https://github.com/fbbdev)在   [#3334](https://github.com/wailsapp/wails/pull/3334) 中完成   [@abichinger](https://github.com/abichinger) 在   [#3346](https://github.com/wailsapp/wails/pull/3346) 中   [@hfoxy](https://github.com/hfoxy) 在   [#3417](https://github.com/wailsapp/wails/pull/3417) 中   [@hfoxy](https://github.com/hfoxy) 在   [#3426](https://github.com/wailsapp/wails/pull/3426) 中   [@fbbdev](https://github.com/fbbdev) 在   [#3431](https://github.com/wailsapp/wails/pull/3431) 中   [@fbbdev](https://github.com/fbbdev) 在   [#3431](https://github.com/wailsapp/wails/pull/3431) 中   由[@pekim](https://github.com/pekim)在   [#3458](https://github.com/wailsapp/wails/pull/3458) 中完成   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)由   [@5aaee9](https://github.com/5aaee9)提交。   [@windom](https://github.com/windom/) 在   [#3636](https://github.com/wailsapp/wails/pull/3636) 中。   Windows 相关更改由[@bruxaodev](https://github.com/bruxaodev/)在   [#3691](https://github.com/wailsapp/wails/pull/3691) 中完成。

### 变更

[mmghv](https://github.com/mmghv) 在   [#3611](https://github.com/wailsapp/wails/pull/3611) 中   支持 Ubuntu 24.04 LTS，由[atterpac](https://github.com/atterpac)在   [#3461](https://github.com/wailsapp/wails/pull/3461) 中完成   必须具有`type="module"`属性。由   [@fbbdev](https://github.com/fbbdev)在   [#3295](https://github.com/wailsapp/wails/pull/3295) 中完成   对象，并且不启动 WML 系统。这样做是为了改进   封装。如果需要，可以调用新的`WML.Enable`方法手动启动   WML 系统。随附的 JS 运行时脚本仍会自动执行这两项   操作。由[@fbbdev](https://github.com/fbbdev)在   [#3295](https://github.com/wailsapp/wails/pull/3295) 中完成   窗口对象作为默认导出。现在无法再通过 ESM 命名导入或命名空间导入语法   导入单独的方法。   API。部分方法的名称或原型已更改，具体如下：`Screen`   变为`GetScreen`；`GetZoomLevel`/`SetZoomLevel`变为`GetZoom`/`SetZoom`；   `GetZoom`、`Width`和`Height`现在直接返回值，而不再将其   包装在对象中。由[@fbbdev](https://github.com/fbbdev)在   [#3295](https://github.com/wailsapp/wails/pull/3295) 中完成   已移除。使用`-names` CLI 选项可切换回按名称调用。   由[@fbbdev](https://github.com/fbbdev)在   [#3468](https://github.com/wailsapp/wails/pull/3468) 中完成   以前按照其所在包命名；现在使用完整的 Go 导入路径，   包括模块路径。由[@fbbdev](https://github.com/fbbdev)在   [#3468](https://github.com/wailsapp/wails/pull/3468) 中完成   `application.Options.Services`。由[@fbbdev](https://github.com/fbbdev)在   [#3468](https://github.com/wailsapp/wails/pull/3468) 中完成   调用`application.NewService`。由[@fbbdev](https://github.com/fbbdev)在   [#3468](https://github.com/wailsapp/wails/pull/3468) 中完成   [@DeltaLaboratory](https://github.com/DeltaLaboratory) 在   [#3574](https://github.com/wailsapp/wails/pull/3574) 中
