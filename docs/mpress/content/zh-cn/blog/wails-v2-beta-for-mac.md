---
title: "Wails v2 Mac 版 Beta"
description: "Wails 发布说明与公告"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![wails-mac 屏幕截图](/assets/blog-images/wails-mac.webp)

今天，Wails v2 Mac 版迎来了首个 Beta 版本！我们花了相当长的时间才走到这一步，希望今天发布的版本能为你提供一些切实可用的功能。一路走来经历了不少波折，希望在你的帮助下，我们能够解决这些问题，进一步完善 Mac 移植版，为 v2 正式发布做好准备。

你的意思是，这还不能用于生产环境吗？对你的用例而言，它很可能已经可用，但目前仍有一些已知问题，因此请持续关注[此项目看板](https://github.com/wailsapp/wails/projects/7)。如果你愿意贡献力量，我们将非常欢迎！

那么，与 v1 相比，Wails v2 Mac 版有哪些新变化？提示：它与 Windows Beta 版非常相似 :wink:

## 新功能

![wails-menus-mac 屏幕截图](/assets/blog-images/wails-menus-mac.webp)

许多用户都希望获得原生菜单支持。Wails 现在终于满足了这一需求。应用菜单现已可用，并支持大多数原生菜单功能，包括标准菜单项、复选框、单选项组、子菜单和分隔符。

在 v1 中，大量用户希望能够更细致地控制窗口本身。很高兴地宣布，我们为此新增了专门的运行时 API。它功能丰富，并支持多显示器配置。对话框 API 也得到了改进：现在，你可以使用配置丰富的现代原生对话框，满足各种对话框需求。

### Mac 专用选项

除了常规应用选项外，Wails v2 Mac 版还带来了一些 Mac 专属功能：

- 让窗口变得新潮而通透，就像那些精美的 Swift 应用一样！
- 高度可定制的标题栏
- 支持应用的 NSAppearance 选项
- 通过简单配置自动创建“关于”菜单

### 无需打包资源

v1 的一大痛点是必须将整个应用压缩为单个 JS 文件和单个 CSS 文件。很高兴地宣布，在 v2 中完全不再要求以任何方式打包资源。想加载本地图片？使用带有本地 src 路径的`<img>`标签。想使用一款很酷的字体？将它复制进来，然后在 CSS 中添加其路径即可。

> 哇，听起来就像一个 Web 服务器……

没错，它的工作方式就像 Web 服务器，只不过它并不是 Web 服务器。

> 那么，该如何加入我的资源？

只需将一个包含所有资源的`embed.FS`传入应用配置即可。这些资源甚至不必位于顶层目录中——Wails 会自动为你处理。

### 全新的开发体验

由于资源不再需要打包，一种全新的开发体验由此成为可能。新的`wails dev`命令会构建并运行应用，但它不会使用`embed.FS`中的资源，而是直接从磁盘加载资源。

它还提供以下附加功能：

- 热重载——前端资源的任何更改都会触发应用前端自动重新加载
- 自动重新构建——Go 代码的任何更改都会触发应用重新构建并重新启动

此外，还会在端口34115上启动一个 Web 服务器。它会向任何连接到该端口的浏览器提供你的应用。所有已连接的 Web 浏览器都会响应系统事件，例如资源更改时的热重载。

在 Go 中，我们习惯在应用里处理结构体。将结构体发送到前端并用作应用状态通常很有用。在 v1 中，这个过程非常依赖手动操作，给开发者带来了一些负担。很高兴地宣布，在 v2 中，任何以开发模式运行的应用都会为绑定方法的所有结构体输入或输出参数自动生成 TypeScript 模型，从而让数据模型能够在前后端之间无缝交换。

此外，系统还会动态生成另一个 JS 模块，用于封装所有绑定方法。该模块会为你的方法提供 JSDoc，从而在 IDE 中提供代码补全和提示。当你在封装 Go 代码的自动生成模块中按下 Tab 键时，数据模型会被自动导入，效果真的很酷！

### 远程模板

![remote-mac 屏幕截图](/assets/blog-images/remote-mac.webp)

快速启动并运行应用一直是 Wails 项目的关键目标。项目刚发布时，我们曾尝试覆盖当时众多现代框架，包括 react、vue 和 angular。前端开发领域充满各种不同的主张，发展迅速，很难始终紧跟变化！因此，我们发现基础模板很快就会过时，造成了维护难题。这也意味着，我们无法为最新、最优秀的技术栈提供出色的现代模板。

在 v2 中，我希望让社区拥有更大的自主权，使你们能够自行创建和托管模板，而不必依赖 Wails 项目。因此，你现在可以使用由社区支持的模板创建项目！我希望这能激励开发者共同打造一个充满活力的项目模板生态系统。我非常期待开发者社区能够创造出怎样的成果！

### 原生 M1 支持

感谢[Mat Ryer](https://github.com/matryer/)的大力支持，Wails 项目现在支持 M1 原生构建：

![build-darwin-arm 屏幕截图](/assets/blog-images/build-darwin-arm.webp)

你也可以将`darwin/amd64`指定为目标：

![build-darwin-amd 屏幕截图](/assets/blog-images/build-darwin-amd.webp)

哦，我差点忘了……你还可以使用`darwin/universal`…… :wink:

![build-darwin-universal 屏幕截图](/assets/blog-images/build-darwin-universal.webp)

### 交叉编译到 Windows

由于 Wails v2 Windows 版完全使用 Go 编写，因此无需 docker 即可将 Windows 指定为构建目标。

![build-cross-windows 屏幕截图](/assets/blog-images/build-cross-windows.webp)  
bu

### WKWebView 渲染器

V1 依赖一个现已弃用的 WebView 组件。V2 使用最新的 WKWebKit 组件，因此你可以享用 Apple 提供的最新优秀功能。

### 结语

正如我在 Windows 发布说明中所说，Wails v2 为项目奠定了新的基础。此次发布旨在收集大家对新方案的反馈，并在正式发布前解决所有缺陷。非常期待您的意见！请前往 [v2 Beta](https://github.com/wailsapp/wails/discussions/828) 讨论区提交任何反馈。

最后，我要特别感谢所有[项目赞助者](/credits/#sponsors)，其中包括[JetBrains](https://www.jetbrains.com?from=Wails)。他们的支持在幕后以多种方式推动着项目发展。

项目即将迈入令人振奋的新阶段，我期待看到大家使用 Wails 构建出怎样的作品！

Lea。

附言：Linux 用户，接下来就轮到你们了！

再附言：如果您或您的公司认为 Wails 很有用，请考虑[赞助本项目](https://github.com/sponsors/leaanthony)。谢谢！
