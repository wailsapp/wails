---
title: "Wails v3.1.0：更多 Go，更少构建障碍"
description: "实验性的无 CGO 桌面构建与 Wake：一种更简单的 Wails 应用构建、打包和签名方式。"
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

Wails 忙碌的一周：无 CGO 构建、更简便的项目设置，以及更顺畅的打包和签名流程。

Wails v3.1.0开启了每周功能发布，让你能够定期获得新功能。错误修复仍会在准备就绪后通过每夜补丁更新持续发布。

## 本周的修复与优化

本周在框架、CLI、构建工具和 CI 中完成了**100项修复、21项性能改进和四项安全加固变更**。重点包括：在混合 DPI 环境中更可靠地定位窗口、修复 Linux 托盘和深色模式问题、精简资源服务，以及降低 MSIX 打包期间的内存占用。[请参阅下面列出的部分变更](#heading-4)。

## 实验性的无 CGO 构建

借助 purego，Go 无需使用 CGo（Go 通常用于对接 C 代码的桥梁）即可调用原生系统库。purego 刚出现时，我们就开始探索这项技术，但维护两个后端颇具挑战。如今 purego 日趋成熟，生成式 AI 也让维护工作更易于管理，因此我们准备分享实验性的 macOS 和 Linux 后端。

现在，无需 C 编译器即可构建和交叉编译桌面应用：

```sh
WAILS_PUREGO=1 wails3 build
```

构建前端资源后，还可以直接使用 Go 构建：

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

macOS 后端使用 AppKit 和 WebKit；Linux 支持 GTK4/WebKitGTK 6.0和 GTK3。用户仍需安装原生运行时库，打包和签名也各有额外要求。在我们收集反馈期间，CGo 仍是默认选项。

[实现与平台详情](https://github.com/wailsapp/wails/pull/6521)

## 为什么选择 Wake？

Wails v2隐藏了过多构建过程。在v3中，Taskfile 公开了整个过程，但大多数人只想进行少量自定义，而不想维护一套构建系统。

Wake 将常见选项集中到`wails.yaml`中，包括应用、配套 CLI 工具、服务器、目标和打包配置。Task 仍是我们的核心构建系统；需要其灵活性时，依然可以使用。

Wake 现已支持实时重载、增量构建、输出缓存和并行打包。已经有二进制文件？即使未安装 Go 或前端工具，`--skip-build`也能对其进行打包。

使用`WAILS_EXP_USE_WAKE=1`启用它。

### 引导式入门

启用 Wake 后，`wails3 init -ui`会引导你选择后端、目标和软件包，检查所需工具，并保存项目配置。

@image{src="/assets/screenshots/wake/application.png" alt="Wake 项目向导显示桌面应用和后端选项，其中包括实验性的 Pure Go。" expand}

为每个目标选择软件包格式和签名选项：

@image{src="/assets/screenshots/wake/packages.png" alt="Wake 的 Linux 软件包选择界面，显示 AppImage、DEB、RPM、Void 和 Arch Linux 格式，以及为发布输出签名的选项。" expand}

下面的工具检查展示了标准 Linux 构建：

@image{src="/assets/screenshots/wake/tools.png" alt="Wake 正在为标准 Linux 构建检查 C 编译器、Go、npm 和 pkg-config。" expand}

### 我的计算机可以构建什么？

编译、打包和签名的要求各不相同。以下是你可以在各类主机上执行的操作：

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| 操作 | Windows 主机 | macOS 主机 | Linux 主机 | 所需条件 |
| 使用 purego 编译桌面应用 | ✓ | ✓ | ✓ | Go，以及用于 macOS/Linux 目标的实验性 purego 后端 |
| 创建 MSIX 软件包 | ✓ | ✓ | ✓ | 已构建的应用；无需 Windows SDK |
| 创建 NSIS 安装程序 | ✓ | ✓ | ✓ | 已构建的应用和`makensis` |
| 创建 macOS 应用和 DMG | ✓ | ✓ | ✓ | 已构建的应用；由 Go 创建软件包 |
| 创建 DEB、RPM、Arch 和 Void 软件包 | ✓ | ✓ | ✓ | 已构建的应用 |
| 创建 AppImage | ✓ | ✓ | ✓ | 目标 AppImage 运行时；交叉打包还需要 Linux GUI 运行时目录树 |
| 使用导出的身份签署 Windows 软件包 | ✓ | ✓ | ✓ | 导出的 RSA PFX 证书和私钥 |
| 使用导出的身份签署 Apple 应用和 DMG | ✓ | ✓ | ✓ | 导出的 RSA PKCS#12身份；请参阅下方的验证说明 |
| 签署 DEB、RPM 和 Arch 软件包 | ✓ | ✓ | ✓ | 导出的 OpenPGP 密钥 |
| 提交 Apple 公证并附加票据 | ✓ | ✓ | ✓ | Apple 团队 API 密钥和经过适当签名的应用 |
| 使用 Apple 的`codesign`或钥匙串配置文件 | – | ✓ | – | macOS 和相应的签名身份 |
@end

CGo 构建仍需要 C 工具链。修改 DMG 的自定义钩子需要 macOS 和`hdiutil`。

### 签名与公证

Wake 支持 Windows、Apple 和 Linux 软件包签名，包括 Apple 临时签名。跨平台签名使用导出的身份；证书存储和钥匙串中的身份仍需使用原生工具。用户级凭据可避免将私钥写入项目配置。目前尚不支持 AppImage 和 Void 签名。

使用团队 API 密钥，即可通过 Go 完成 Apple 公证。我们已在 Linux 上测试提交和票据附加，Gatekeeper 接受了处理结果。该应用使用 Apple 的[`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html)进行签名；Wake 通过 Go 生成的 Developer ID 签名能否获得 Apple 接受，仍未经验证。

### 迁移到 Wake

迁移支持全部九个内置模板。此过程不可逆：激活`wails.yaml`会替换 Taskfile 路由。不受支持的自定义命令和覆盖项会被拒绝，需人工审查。

[Wake 功能、要求和迁移详情](https://github.com/wailsapp/wails/pull/6531)

## 本周部分亮点

- <strong>窗口现在能在混合 DPI 显示器上正确定位。</strong>在采用不同显示缩放比例的 macOS 和 Linux 配置中，屏幕边界和位置现在均能得到正确处理。[#6346](https://github.com/wailsapp/wails/pull/6346)
- <strong>Linux 托盘和主题检测更加可靠。</strong>应用启动后创建的托盘不再丢失，工具提示会正确填充，深色模式检测也会从启动时起读取桌面的外观设置。
- **Go 事件回调会为每个监听器保持顺序。** 每个监听器都按事件发出顺序串行运行；不同监听器可以并发运行，而并发发出的事件没有确定的顺序。事件钩子会在锁的保护下生成快照，以避免数据竞争。[#6275](https://github.com/wailsapp/wails/pull/6275)、[#6276](https://github.com/wailsapp/wails/pull/6276)
- <strong>构建缓存可处理更多实际项目。</strong>带引号或转义的空格不再导致 Go 输出缓存未命中，并发缓存保存操作也会串行执行，以防文件损坏。[#6495](https://github.com/wailsapp/wails/pull/6495)、[#6488](https://github.com/wailsapp/wails/pull/6488)
- <strong>资源服务减少了不必要的工作。</strong>对于已知文件扩展名，不再探测 MIME 类型；禁用调试日志后，也不再预先生成最终会被丢弃的消息。两项变更都减少了请求处理路径上的工作量。[#6414](https://github.com/wailsapp/wails/pull/6414)、[#6420](https://github.com/wailsapp/wails/pull/6420)
- <strong>MSIX 打包不再使用与二进制文件等大的缓冲区。</strong>可执行文件会以流式方式写入软件包，而不再整体读入内存。此前，复制一个具有代表性的16 MiB 文件，每次操作约会分配16.8 MB 内存。[#6472](https://github.com/wailsapp/wails/pull/6472)

## 亲自试用

这两项功能都处于实验阶段。请尝试进行跨平台构建、打包你的应用，并告诉我们你在哪里遇到了阻碍。感谢所有参与测试、报告问题和贡献代码的人。

还要特别感谢 Discord 社区的 beesophie，她构建了令人惊叹的构建系统[Cook](https://github.com/SimTXNG/cook)，让我们看到了自身系统中的一些不足。Cook 所展现的可能性帮助我们认识到，可以在哪些方面简化 Wails 的构建流程，也为这项工作带来了真正的灵感。

谢谢你，Beesophie！
