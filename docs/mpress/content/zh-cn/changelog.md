---
title: "更新日志"
description: "Wails v3的版本历史和发布说明"
slug: "changelog"
sourcePath: "changelog.md"
---

图例：

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- 此项目的所有重要变更都将记录在此文件中。

格式基于[维护更新日志](https://keepachangelog.com/en/1.0.0/)，并且此项目遵循[语义化版本](https://semver.org/spec/v2.0.0.html)。

- `Added`表示新增功能。
- `Changed`表示现有功能的变更。
- `Deprecated`表示即将移除的功能。
- `Removed`表示现已移除的功能。
- `Fixed`表示任何错误修复。
- `Security`表示安全漏洞。

_/

/_   * 自动更新 *   发布说明根据已合并的 PR 生成。请勿在普通贡献 PR 中编辑此文件，也不要重新创建已移除的`v3/UNRELEASED_CHANGELOG.md`。有关更新日志维护的信息，请参阅`AGENTS.md`。 _/

## [未发布]

## v3.2.0-nightly.20261102 - 2026-11-02

## 变更
- 使用 DisableOverscroll 替代 DisableBounce，由 @lethusang 在[PR](https://github.com/wailsapp/wails/pull/6077)中完成

## v3.2.0-nightly.20261101 - 2026-11-01

## 新增
- 添加实验性的无 CGO macOS 和 Linux 后端，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6521)中完成
- 添加原生 YAML 构建、可移植打包和签名，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6531)中完成

## 修复
- 修复多项v3夜间构建问题并清理依赖项，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6519)中完成
- 修复v3发布验证和固件分类，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6533)中完成
- 解决跨平台发布验证失败问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6538)中完成

## v3.1.0 - 2026-11-01

## 新增
- 添加实验性的无 CGO macOS 和 Linux 后端，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6521)中完成
- 添加原生 YAML 构建、可移植打包和签名，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6531)中完成

## 变更
- 避免复制已丢弃的桌面元数据，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6391)中完成
- 不使用 fmt 格式化内容长度，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6393)中完成
- 无需收集候选项即可选择索引根目录，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6404)中完成
- 复用不可定位文件的复制缓冲区，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6408)中完成
- 避免为探测前缀分配读取器，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6410)中完成
- 对已知扩展名跳过 MIME 探测读取，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6414)中完成
- 防止队列增长和收缩振荡，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6415)中完成
- 避免分配随后会被丢弃的请求上下文，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6419)中完成
- 跳过已禁用的调试日志准备工作，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6420)中完成
- 延迟分配 MIME 探测缓冲区，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6428)中完成
- 减少超链接格式化产生的内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6443)中完成
- 消除项目符号格式化产生的临时内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6449)中完成
- 将 UUID 字符串转换减少为一次内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6453)中完成
- 减少语句句柄格式化产生的内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6454)中完成
- 消除基于名称的 UUID 摘要产生的临时内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6455)中完成
- 避免复制尚未完成的命令输出行，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6460)中完成
- 在包暂存期间以流式方式复制可执行文件，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6472)中完成
- 避免重新扫描尚未完成的输出行，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6473)中完成
- 复用已编译的命令路由模式，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6475)中完成
- 避免在运行时发现期间查询元数据，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6482)中完成
- 减少本地化后备页面产生的内存分配，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6485)中完成
- 避免为每个处理器重新编译模式，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6490)中完成

## 修复
- 为系统托盘实现 SetTooltip，由 @fan711 在[PR](https://github.com/wailsapp/wails/pull/6017)中完成
- 从设置门户读取深色模式，由 @TechHutTV 在[PR](https://github.com/wailsapp/wails/pull/6072)中完成
- 允许 Window.SetZoom 低于1.0，由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/6166)中完成
- 按发出顺序传递给每个 Go 监听器，由 @DevLumuz 在[PR](https://github.com/wailsapp/wails/pull/6275)中完成
- 在 Emit 中于 hookLock 保护下读取事件钩子，由 @DevLumuz 在[PR](https://github.com/wailsapp/wails/pull/6276)中完成
- 始终发布托盘的 Menu 属性，避免丢弃在 Run 之后创建的托盘，由 @AlbinoGeek 在[PR](https://github.com/wailsapp/wails/pull/6284)中完成
- 将在 Run 之前发送的主线程工作保留到 Run 启动，由 @AlbinoGeek 在[PR](https://github.com/wailsapp/wails/pull/6287)中完成
- 在 GitHub 更新程序提供程序中使用正确的资源 URL，由 @gavidroselj 在[PR](https://github.com/wailsapp/wails/pull/6304)中完成
- 当 WebView2 省略 Content-Length 时保留请求正文，由 @strraaa 在[PR](https://github.com/wailsapp/wails/pull/6326)中完成
- 在 DIP 空间中报告 macOS 和 GTK 屏幕（混合 DPI 边界），由 @julianstorer 在[PR](https://github.com/wailsapp/wails/pull/6346)中完成
- 不要将已关闭的通知报告为点击，由 @bercedev 在[PR](https://github.com/wailsapp/wails/pull/6367)中完成
- 终端关闭时停止开发进程，由 @xmasdev 在[PR](https://github.com/wailsapp/wails/pull/6371)中完成
- 在证据检查中将更新日志运行映射到压缩合并提交，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6373)中完成
- 将任务栏重启防抖保持一秒，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6386)中完成
- 进入全屏时若窗口已处于全屏状态，则保留尺寸约束，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6387)中完成
- 不要为适应新约束而调整全屏或最大化窗口的大小，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6388)中完成
- 在液态玻璃示例中，不要在 Run 之前弹出对话框，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6389)
- 在 Run 之前请求的对话框不执行任何操作，而不是无限等待，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6390)
- 当桌面条目没有命令时使用回退方案，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6392)
- 探测当前内容时不保留文件名，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6394)
- 识别大小写混合的资源扩展名，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6395)
- 释放已使用的内容探测缓冲区，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6396)
- 释放已移除的 ringqueue 引用，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6397)
- 在选择回退方案之前遵循标头删除操作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6398)
- 在刷新之前提交回退响应标头，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6399)
- 仅查找名称完全匹配的索引文件，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6400)
- 对不可查找的 HEAD 请求避免返回正文和执行完整读取，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6401)
- 跨多次短读取填满 MIME 探测前缀，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6402)
- 提交成功但为空的不可查找响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6403)
- 避免重复已展开的 URI 目标，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6405)
- 在字段展开期间保留路径，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6406)
- 遵循 XDG 桌面文件搜索目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6407)
- 在空刷新时提交探测器标头，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6409)
- 跳过遮蔽桌面条目的目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6411)
- 保留在可执行文件路径内找到的目标，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 在计时器触发后释放回调捕获的内容，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 避免 ldconfig 库误匹配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 限制文件系统库回退匹配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6417)
- 消除缓存失效的数据竞争，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 遵循可选的日志记录器配置，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 避免修改共享的启动 URL 以及使用过期端口，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 请求取消时停止开发代理重试，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 缺少捆绑资源时返回404，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 处理捆绑运行时的 HEAD 请求，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 路由服务时保留转义路径，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 正确排序超大数值型预发布标识符，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 在保存工具下载内容之前拒绝 HTTP 错误响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在桌面发现期间遵循启动截止时间，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 在开发代理重试期间保留请求正文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 在每条返回路径上关闭 HTTP 响应正文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 在并行库搜索中保留取消错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 在解码之前拒绝失败的 GitHub 标签响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 保留包含方括号的模板路径原文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 只安装一次 wails3，并保留安装程序错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 忽略注释块中的隐藏条目和标头，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 公开稳定版本查找错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 不进行格式解释，直接输出纯文本消息，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 避免加载 JSON null 后写入时发生 panic，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 读取整个存储时返回快照，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 在部分行迭代后传播错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 保留解压目标路径原文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 在关闭和重新打开时释放预处理语句句柄，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 处理预处理操作与关闭或重新打开操作之间的竞争，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- 正确处理回退索引页的 HEAD 请求，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 要求提供完整的十六进制语句句柄，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6452)
- 在并发情况下保留默认内存数据库，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 将 stderr 与解析后的 HEAD 输出分离，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 从嵌入式捆绑包生成运行时资源，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 检测 TODO、FIXME 和 TBD 占位符，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 验证发行说明 API 响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6461)
- 在状态检查中遵循任务目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 防止复制同一文件导致数据丢失，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 将任务目录纳入缓存标识，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 每次检查仅探测一次端口，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 避免缓存条目解码为 null 时发生 panic，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 不受整数范围限制地递增数值型预发布版本，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 防止旧版检查掩盖缓存未命中，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 在执行器的基准目录中运行默认任务，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 在资源管理器命令中保留百分号原文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 在[PR](https://github.com/wailsapp/wails/pull/6471)中统一前置条件与任务的命令目录，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6474)中用生成的常规文件替换符号链接，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6476)中确保构建失败时保留现有 DMG，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6477)中保留带引号的正则表达式和路径中的反斜杠，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6478)中分配图像前验证生成的图标尺寸，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6479)中通过 shell 处理命令链和重定向，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6480)中确保构建失败时保留以前的 Linux 软件包，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6481)中使用特定于平台的 shell 后备执行方式，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6483)中支持提供路径开头含多个斜杠的资源，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6484)中识别带有构建元数据的稳定版本，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6486)中准确报告资源生成失败，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6487)中防止替换当前报告器时出现竞态，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6488)中串行执行并发的缓存文件保存操作，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6489)中支持多个模板相等候选项，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6491)中规范化已保存的 null 类别映射，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6492)中显示执行失败的变量 shell 命令所产生的诊断信息，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6493)中不将 shell 命令链纳入隐式 Go 缓存，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6494)中串行执行 Linux 类别的持久化操作，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6495)中缓存含带引号或转义空格的 Go 输出，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6496)中确保编码失败时保留 Mac 图标，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6497)中不再丢弃无效的 Go 测试计数标志，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6498)中将文件复制到新文件时保留源文件权限，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6519)中修复各类v3夜间构建问题并清理依赖项，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6533)中修复v3版本验证和固件分类，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6538)中解决跨平台版本验证失败问题，作者：@taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## 变更
- 在[PR](https://github.com/wailsapp/wails/pull/6490)中通过复用已编译的正则表达式模式提高变更日志生成速度，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6485)中减少提供本地化后备页面时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6482)中避免不必要的元数据查询，加快 AppImage 运行时发现，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6475)中复用已编译的命令路由模式以提高 wake 性能，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6473)中避免重复扫描未完成的输出行以提高性能，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6472)中为 MSIX 软件包暂存以流式方式复制可执行文件，减少内存用量并加快构建，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6460)中避免复制未完成的命令输出行以提高性能，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6455)中减少生成基于名称的 UUID 时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6454)中减少格式化 SQLite 语句句柄时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6453)中将 UUID 字符串转换所需的内存分配减少至一次，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6449)中减少格式化变更日志项目符号条目时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6443)中减少格式化终端超链接时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6428)中通过延迟探测 MIME 缓冲区，减少提供资源时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6420)中让资源服务器在禁用调试日志时跳过调试日志准备，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6419)中避免创建会被丢弃的请求上下文，减少处理 webview 请求时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6414)中对已知文件扩展名跳过 MIME 探测，加快资源提供，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6410)中避免使用读取器探测内容，减少提供资源时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6408)中减少提供不可寻址资源文件时的内存分配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6404)中避免收集候选项，加快资源服务器选择索引根目录，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6393)中不使用 fmt 格式化内容长度，提高资源服务器性能，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6391)中减少读取桌面文件资源管理器元数据时的内存复制，作者：@taliesin-ai

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6498)中将源文件复制到新的目标文件时保留其权限，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6388)中防止全屏和最大化窗口在最小或最大尺寸约束发生变化时调整大小，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6500)中发布 wailsapp/wails 仓库的定期星标历史更新，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6497)中让 wake 不再丢弃无效的 Go 测试计数标志，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6496)中确保编码失败时保留 Mac 图标，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6495)中修复含带引号或转义空格的 Go 输出出现构建缓存未命中的问题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6494)中串行执行 Linux 上通知类别的持久化操作，防止并发写入竞态，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6493)中不将 shell 命令链纳入隐式 Go 构建缓存，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6491)中修复保存的数据为 null 时通知类别映射无法加载的问题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6492)中显示执行失败的变量 shell 命令所产生的诊断信息，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6489)中修复存在多个相等候选项时的 wake 模板匹配，作者：@taliesin-ai
- 序列化并发的缓存文件保存操作，防止缓存写入损坏，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6488)
- 防止替换当前报告器时发生竞态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6487)
- 在 CLI 中准确报告资源生成失败，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6486)
- 让 CLI 能够识别包含构建元数据的稳定版本，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6484)
- 提供请求路径含有重复前导斜杠的资源，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6483)
- 修复 wake 在各平台上的 shell 回退执行，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6481)
- 构建失败时保留原有的 Linux 软件包，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6480)
- 通过 shell 执行命令链和重定向，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6479)
- 解析 wake 命令时保留带引号的正则表达式和路径中的反斜杠，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6477)
- 打包构建失败时保留现有的 DMG 文件，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6476)
- 在绑定输出中用生成的常规文件替换符号链接，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6474)
- 在 Linux 文件浏览器命令中保留字面百分号，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 对齐前置条件与任务命令的目录，使其解析行为保持一致，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6471)
- 在执行器的基准目录中运行默认任务，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 防止旧版检查掩盖缓存未命中，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 修复 CLI 预发布版本编号，使其递增不再受整数溢出限制，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6467)
- 避免缓存条目解码为 null 时发生 panic，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6466)
- 在 CLI 中每次检查仅探测一次端口，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6465)
- 将任务目录纳入缓存标识，避免使用过期的 wake 结果，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6464)
- 防止打包过程将文件复制到其自身时丢失数据，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6463)
- 使状态检查遵循任务目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6462)
- 检测变更日志条目中的 TODO、FIXME 和 TBD 占位符，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6459)
- 在 CLI 中从嵌入式捆绑包生成运行时资源，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6458)
- 在 git 命令中将 stderr 与解析后的 HEAD 输出分离，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6457)
- 在并发访问时保持默认的内存 SQLite 数据库不变，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6456)
- 对回退索引页面的 HEAD 请求返回正确响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6451)
- 处理 SQLite 语句预备操作与数据库关闭或重新打开之间的竞态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6450)
- 关闭或重新打开 SQLite 数据库时释放预备语句句柄，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6448)
- 使用 gosod 时保持提取目标路径的字面形式，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6447)
- 在部分迭代行之后传播 SQLite 查询错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6446)
- 在 kvstore 中读取整个存储时返回快照，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6444)
- 防止 kvstore 加载 JSON null 值后写入时发生 panic，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6442)
- 输出纯文本终端消息，不解释格式化动词，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6441)
- 在更新程序中查询最新稳定版本时显示错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6440)
- 忽略变更日志注释块中的隐藏条目和标题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6439)
- 仅安装一次 wails3，并在更新程序中显示安装程序错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6438)
- 在 gosod 中保留带方括号的模板路径字面值，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6437)
- 更新程序在解码 GitHub 标签响应前拒绝失败的响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6436)
- 在并行搜索库时保留取消错误，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6435)
- 在每条返回路径上关闭 HTTP 响应体，防止资源泄漏，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6434)
- 开发代理重试请求时保留请求体，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6432)
- 发现桌面文件浏览器时遵守启动截止时间，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6430)
- 保存 AppImage 工具下载内容前拒绝 HTTP 错误响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6429)
- 在语义版本比较中正确排序超大数值型预发布标识符，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6427)
- 在资源服务器中路由服务时保留转义路径，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6426)
- 在资源服务器中处理捆绑运行时的 HEAD 请求，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6425)
- 缺少捆绑资源时返回404，而非错误响应，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6424)
- 请求取消时停止开发代理重试，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6423)
- 避免资源服务器中共享的起始 URL 被修改以及端口信息过期，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6422)
- 使资源服务器遵循可选的日志记录器配置，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6421)
- 消除库路径缓存失效处理中的数据竞态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6418)
- 搜索 ldconfig 输出时避免错误匹配库，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6416)
- 防止资源服务器队列在负载下增长与收缩振荡，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 防抖计时器触发后释放回调捕获，避免保留过期引用，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6413)
- 保留在可执行文件路径中找到的文件浏览器目标，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6412)
- 在文件浏览器中跳过遮蔽桌面条目的目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6411)
- 修复响应体为空时不发送资源服务器响应头的问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6409)
- 查找应用程序图标时遵循 XDG 桌面文件搜索目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6407)
- 在[PR](https://github.com/wailsapp/wails/pull/6406)中保留文件资源管理器字段展开期间的文件路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6405)中避免文件资源管理器重复显示展开后的 URI 目标，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6403)中提交成功的空且不可查找的资源响应，而不再使其保持待处理状态，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6402)中修复资源服务器跨多次短读取填充 MIME 探测前缀的问题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6401)中避免资源服务器为不可查找的 HEAD 请求读取完整响应正文，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6400)中使资源服务器解析目录请求时仅提供精确匹配的索引文件，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6367)中修复 Linux 通知被关闭后，在不支持操作的服务器上误报为点击的问题，作者：@bercedev
- 在[PR](https://github.com/wailsapp/wails/pull/6371)中使终端关闭时停止开发进程，作者：@xmasdev
- 在[PR](https://github.com/wailsapp/wails/pull/6386)中将 Windows 任务栏重启的防抖标志保持一秒后再清除，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6399)中修复资源服务器丢失回退响应标头的问题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6398)中使资源服务器在选择回退响应前遵循标头删除操作，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6397)中释放资源服务器环形队列中已移除的引用，以避免内存泄漏，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6396)中释放资源服务器已使用的内容探测缓冲区，以减少内存占用，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6395)中使资源服务能够识别大小写混合的资源文件扩展名，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6394)中使资源服务器无需保留文件名即可正确检测内容类型，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6392)中使桌面条目没有命令时回退到默认值，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6389)中阻止 liquid-glass 示例在非 macOS 系统上于应用启动前显示对话框，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6387)中，当 Windows 窗口从受约束尺寸进入全屏时保留窗口尺寸约束，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6390)中修复在 Run 前调用的对话框永久阻塞而不返回的问题，作者：@taliesin-ai

## 安全性
- 在[PR](https://github.com/wailsapp/wails/pull/6478)中，在分配图像内存前验证生成的图标尺寸，以防止内存占用过高，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6461)中使用发布说明 API 响应前对其进行验证，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6452)中要求 SQLite 插件使用完整的十六进制语句句柄，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6417)中限制文件系统库的回退匹配，以防止加载非预期的库，作者：@taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## 变更
- 在[PR](https://github.com/wailsapp/wails/pull/6366)中，将404页面链接到v2.wails.io 上对应的 Wails v2文档页面，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6365)中，将v2文档站点迁移到v2.wails.io，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6364)中，通过 wails.io 提供v3文档，并通过v2.wails.io 提供v2文档，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6363)中更新v3.0.0发布说明并更正安全公告中的说法，作者：@taliesin-ai

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6275)中按照触发顺序串行执行 Go 事件监听器回调，作者：@DevLumuz
- 在[PR](https://github.com/wailsapp/wails/pull/6284)中始终公布托盘的 Menu 属性，以免 Linux 丢弃在 Run 后创建的托盘，作者：@AlbinoGeek
- 在[PR](https://github.com/wailsapp/wails/pull/6287)中将 Run 前提交的主线程工作暂缓到应用启动后执行，作者：@AlbinoGeek
- 在[PR](https://github.com/wailsapp/wails/pull/6304)中使用正确的资源 API URL 下载 GitHub 更新器的发布资源，作者：@gavidroselj
- 在[PR](https://github.com/wailsapp/wails/pull/6346)中修复 macOS 和 Linux 混合 DPI 配置下的屏幕边界和位置，作者：@julianstorer
- 在[PR](https://github.com/wailsapp/wails/pull/6326)中，当 WebView2 省略 Content-Length 标头时，在 Windows 上保留请求正文，作者：@strraaa
- 在[PR](https://github.com/wailsapp/wails/pull/6276)中通过在锁内创建钩子快照，修复触发事件钩子时的数据竞争，作者：@DevLumuz
- 修复`SystemTray.SetTooltip()`在 Linux 上不执行任何操作的问题：现已设置 StatusNotifierItem 工具提示说明（#6017）
- 在[PR](https://github.com/wailsapp/wails/pull/6374)中，在每周发布工作流内运行变更日志证据检查，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6373)中修复发布变更日志证据检查，使其将工作流运行映射到压缩合并提交，作者：@taliesin-ai
- 修复在深色桌面环境下，Linux 深色模式检测在配色方案变化前返回浅色，以及`linux:SystemThemeChanged`在 GTK4 后端从不触发的问题。两个后端现在共用一套实现，并以`org.freedesktop.appearance`门户命名空间为基础（#6072）
- `Window.SetZoom`不再将带边框窗口（Windows、macOS、Linux）中低于1.0的缩放系数提高到1.0：WebView2 `ZoomFactor`、WKWebView `magnification`和 WebKitGTK 的缩放级别都接受任意正值，而`Zoom`窗口选项本就会直接传递这些值。无边框窗口仍保留1.0的下限，因为其拖动和调整大小区域不会根据缩放进行调整（[#4590](https://github.com/wailsapp/wails/issues/4590)）；现在会忽略非正数、NaN 和 Inf 系数，而不再将其提高到1.0——作者：@APshenkin

## v3.0.0 - 2026-10-26

## 变更
- 每日预发布版本以`vX.Y.0-nightly.YYYYMMDD`形式发布，用于预览下一个每周版本；每周版本会递增次版本号（`vX.Y.0`），因此 Go 的`@latest`始终解析为每周版本或热修复版本；`@wailsio/runtime` npm 软件包的版本始终与 Wails 发布版本一致，每日版本使用`nightly` dist-tag
- 项目模板将`@wailsio/runtime`固定为生成它们的 CLI 版本
- Wails v3要求使用 Go 1.26或更高版本（Go 1.26.10或1.27.3及更高版本包含`crypto/x509`安全修复）

## 修复
- 修复`wails3 update cli`找不到v3版本的问题：它只考虑v2标签，并且从错误的路径安装
- 修复 Windows 应用在某些 Windows 版本上绘制应用菜单栏时因`panic: SetTextColor failed`而崩溃的问题：菜单栏现在会在有效的设备上下文上绘制，或回退到默认绘制；`UAHMENUITEM.Position`现在与 C `int`匹配，并且菜单绘制永远不会发生 panic（#6311）
- 在等待重启期间保护暂存更新的安全：关机过程中启动的下载或第二次重启无法删除交换辅助程序正等待安装的载荷，也无法与其争用该载荷，作者：@julianstorer
- 修复 Windows 应用在`applyDarkMode`中崩溃的问题：`pkg/w32`的初始化过程卸载了`uxtheme.dll`，但其序号函数指针仍在使用；现在会在整个进程生命周期内保持该模块已加载（#6356）
- 在交换前验证暂存的有效载荷，并在丢弃备份前验证已安装的应用程序，使损坏的更新能够回滚，而不会替换正常工作的应用程序，作者 @julianstorer

## 已移除
- 最终翻译同步完成后，在 [PR](https://github.com/wailsapp/wails/pull/5736) 中移除 Crowdin 集成，作者 @github-actions[bot]

## 安全性
- 通过向所有帧注入防护脚本，阻止 Linux（GTK3 和 GTK4）上跨源子帧发送保留的 `wails:` 消息，从而缓解 GHSA-2p5f-2gx8-mh8m
- 将 `github.com/klauspost/compress` 更新到 1.18.7（GHSA-259r-337f-4rfw），将 `golang.org/x/image` 更新到 0.45.0，并使用 Go 1.26.10 构建以包含 `crypto/x509` 修复（GO-2026-6809）
- 在运行时包、设置向导、示例和项目模板中，将 `vite` 更新到 8.2.4 或更高版本（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 已知限制
- **Linux**：由于 WebKitGTK 不会向处理程序提供帧信息，子帧防护机制无法阻止*同源*子帧（例如`srcdoc`或`about:blank` iframe）通过其父帧访问原生消息处理程序。跨源帧以及不含`allow-same-origin`的沙箱帧会被阻止。请勿在 Wails 窗口中嵌入不受信任的同源内容。
- **Android（实验性）：**`window.wails` JavaScript 桥接会向 WebView 中的每个帧公开，且不检查主帧或来源。嵌入式 iframe 中的内容（包括跨源内容）可以调用绑定方法并发送保留的 `wails:` 消息。请勿在应用的 WebView 中嵌入不受信任的第三方内容。这与 macOS 和 iOS 不同，后两者会忽略子帧发出的 `wails:` 消息。计划在 3.0.1 中修复此问题。
- <strong>所有平台：</strong>子帧对绑定方法的调用通过 HTTP 传输到 `/wails/runtime`，而不经过原生消息处理程序，因此该防护不会限制这些调用。应将嵌入的任何 iframe 视为能够访问绑定的 Go 方法。

## v3.0.0-beta.34 - 2026-10-20

## 已修复
- 修复 Windows WebView 在窗口激活时无法重新获得焦点，以及 Tab 键被困在 WebView 内而不能将焦点移至原生窗口再移回的问题（#3783）
- 修复 `internal/webview2/pkg/webview2` COM 虚函数表结构（`ICoreWebView2_2`..`_27` 及其他派生接口）遗漏继承的方法链，导致调用落到错误虚函数表槽位的问题（#6129）
- 修复 Windows 询问对话框忽略自定义按钮标签且从不触发按钮回调的问题；现在使用 TaskDialog 显示自定义按钮（comctl32 v6 不可用时，回退到 MessageBox 并按位置映射回调）（#5794）
- 修复面向 Vite 时，`wails3 dev` 会挂起 60 秒的问题：自动前端就绪门控现在使用 TCP 探测，不再使用未发送 `Accept` 标头并收到 404 的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 在 [PR](https://github.com/wailsapp/wails/pull/6301) 中，使系统托盘菜单在打开期间更新后仍能正常工作，作者 @taliesin-ai
- 在 [PR](https://github.com/wailsapp/wails/pull/6302) 中，让运行时操作妥善处理暂时性的 WebView2 COM 错误，作者 @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## 新增
- 绑定生成器现在会生成一个 `bindings/events` 模块，为使用编译期常量名称通过 `application.RegisterEvent` 注册的每个可发现事件导出带类型的名称常量（例如，为 `user:login` 导出 `UserLogin`）；无需 Vite 插件即可与 `Events.On`/`Emit` 配合使用（#5645）

## 已修复
- 修复 Windows 询问对话框忽略自定义按钮标签且从不触发按钮回调的问题；现在使用 TaskDialog 显示自定义按钮（comctl32 v6 不可用时，回退到 MessageBox 并按位置映射回调）（#5794）
- 修复面向 Vite 时，`wails3 dev` 会挂起 60 秒的问题：自动前端就绪门控现在使用 TCP 探测，不再使用未发送 `Accept` 标头并收到 404 的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 在 [PR](https://github.com/wailsapp/wails/pull/6301) 中，使系统托盘菜单在打开期间更新后仍能正常工作，作者 @taliesin-ai
- 在 [PR](https://github.com/wailsapp/wails/pull/6302) 中，让运行时操作妥善处理暂时性的 WebView2 COM 错误，作者 @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## 新增
- 在 [PR](https://github.com/wailsapp/wails/pull/6297) 中为 macOS 14+ 添加原生 macOS 菜单分节标题——请参阅[文档网站](https://v3.wails.io)（已更新 20 个页面），作者 @leaanthony

## 已修复
- 在 [PR](https://github.com/wailsapp/wails/pull/6299) 中，使 SetRelativePosition 在 macOS 上相对于屏幕工作区域生效，作者 @leaanthony
- 在 [PR](https://github.com/wailsapp/wails/pull/6298) 中，使更新程序辅助进程在父进程超时或替换失败后重新启动未更改的目标，作者 @taliesin-ai
- 在 [PR](https://github.com/wailsapp/wails/pull/6293) 中，将版本化文档使用的 pnpm 固定到版本 10，作者 @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## 新增
- 公开 `@wailsio/runtime/cancellable` 子路径导出，使导入 `CancellablePromise` 时无需执行运行时根入口（#6153）

## 已更改
- 要求使用 Go 1.26 或更高版本：v3 模块现在声明 `go 1.26.0`，由 `wails3 init` 和 `wails3 service init` 生成的项目也声明相同的最低版本

## 已修复
- 修复使用 Go 1.27 构建时绑定生成器输出发生变化的问题。当时，`encoding/json` 类型（如 `Marshaler`、`Unmarshaler` 和 `RawMessage`）变成了 `encoding/json/v2` 和 `encoding/json/jsontext` 类型的别名；现在它们会再次呈现为 `any`，且不会为 `encoding/json/v2` 生成模型（#6292）
- 克隆前验证传给 `wails3 init -t` 的远程模板 URL：仅接受 `https://`、`http://`、`ssh://`、`git://` 和 `user@host:path` URL；这些值会在 `--` 之后传给 git，并且克隆期间会将 git 限制为仅使用这些传输协议
- 拒绝 HTTP 运行时传输中超出范围的 `x-wails-window-id` 标头值，而不再静默截断为另一窗口的 ID
- 加强更新程序的归档解压防护以抵御 zip-slip：现在使用 `filepath.IsLocal` 检查条目路径，并对拼接后的路径执行包含关系检查；同时为 zip 和 tar.gz 归档中嵌套的 `..`、反斜杠、绝对路径及通过符号链接写入的条目添加测试
- 修复来自 `wails3 generate appimage` 的 AppImage 在 Debian 和 Ubuntu 之外启动时因“Unable to spawn a new child process: … WebKitNetworkProcess”而崩溃的问题。捆绑的 WebKit 现在会运行 AppImage 随附的辅助进程。WebKit 的 bubblewrap 沙盒保持启用
- 修复 `-tags server` 构建无法为 Android 和 iOS 编译的问题，从而修复 [#6271](https://github.com/wailsapp/wails/issues/6271)
- 修复临时目录位于另一文件系统（例如 tmpfs `/tmp`）时更新程序无法在 Linux 上执行交换的问题：现在会在已安装的二进制文件旁暂存，并跨文件系统复制
- 修复更新程序试图替换 AppImage 的只读挂载点而不是 `.AppImage` 文件的问题
- 修复 macOS 和 iOS 上来自 iframe 的保留`wails:`窗口控制消息会被接受的问题。现在仅接受来自窗口主框架的此类消息；来自 iframe 的其他消息仍会连同其来源信息传递给`RawMessageHandler`
- 停止在 macOS 上发送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`会在每次运行循环迭代时触发，因此页面每秒会收到约120次空的`evaluateJavaScript:`调用；而在 macOS 上，每次27.0调用都会消耗用户点击的瞬时激活状态，导致`navigator.clipboard.writeText()`及其他需要激活状态的 API 对真实点击返回`NotAllowedError`并拒绝执行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 完成
- 将 Windows 的默认 info.json 语言修复为 en-US（0409），以便资源管理器显示版本详细信息（#6210）
- 修复配有 NVIDIA 独立 GPU 和 AMD 或 Intel 集成 GPU 的笔记本电脑会禁用 Linux DMA-BUF 渲染器的问题：现在 Wayland 仅在 WebKit 使用 NVIDIA 渲染时（应用被卸载到独立 GPU 上运行，或合成器在其上运行）才禁用该渲染器，而不是只要加载了 nvidia 模块就禁用；后者会使集成 GPU 上的渲染速度慢得几乎无法使用。X11 会话仍保留原有行为，因为 Xorg 使用的渲染 GPU 由其配置决定
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中，通过分阶段复制回退方案处理 Linux 更新程序跨设备重命名失败的问题，由 @john-okeefe 完成
- 修复 Windows 上的`WebviewWindow.SetMenu`仅在配置了菜单栏主题时才重绘菜单栏的问题，解决[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修复 Windows 上的`WebviewWindow.SetMenu(nil)`会触发 panic，而不是像 macOS 上一样不执行任何操作的问题，解决[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送期间保留变更日志条目，由 @leaanthony 完成
- 移除引用模块中从未发布的二进制文件的嵌入项，修复`go mod vendor`在所有平台上都会因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，解决[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)；见[PR](https://github.com/wailsapp/wails/pull/6031)，由 @Grantmartin2002 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 上重复销毁窗口，由 @charlesbalderson 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中处理系统托盘菜单的 WM_CONTEXTMENU，使其可通过右键单击打开，由 @wzk2025 完成
- 修复全部五种`InvokeSync*`变体在分派的回调触发 panic 时会永久阻塞调用方的问题，解决[#6107](https://github.com/wailsapp/wails/issues/6107)；见[PR](https://github.com/wailsapp/wails/pull/6108)，由 @RALIST 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6236)中安全地重新构建已销毁的 GTK3 子菜单，由 @leaanthony 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6261)中为嵌套 worker 实现 Windows worker 转发垫片，由 @taliesin-ai 完成
- 修复 Windows 上`OpenFile()`在`CanChooseDirectories(true)`时忽略`PromptForMultipleSelection()`的问题：用户现在可以像在 macOS 和 Linux 上一样，一次选择多个文件夹（#6263）
- 修复 windows/arm64 上按值传递给 WebView2 COM 方法的`double`参数问题；该问题影响`ZoomFactor`、`RasterizationScale`、Cookie 的`Expires`，以及打印设置的缩放比例、页面大小和边距（此前打印设置的 setter 会返回`ErrDoubleArgUnsupported`）。受影响的调用现在通过一个可加载`d0`的`runtime.cgocall`跳板执行，而 Go 的 syscall 路径无法加载该寄存器（golang.org/issue/62583），由 @APshenkin 完成

## 安全性
- 服务器模式（`-tags server`）现在仅接受来自服务器自身来源或`ServerOptions.WebSocketOriginPatterns`中所列来源、以`POST`请求发起的`/wails/runtime`调用；其他方法会收到 HTTP 405，其他浏览器来源会收到 HTTP 403。如果部署位于会重写`Host`标头的反向代理之后，应将其公共来源添加到`WebSocketOriginPatterns`
- 绑定到环回地址（默认的`localhost`、`127.0.0.0/8`或`::1`）的服务器模式（`-tags server`）现在会拒绝请求并返回 HTTP 403，除非`Host`标头为`localhost`、IP 字面量、绑定主机，或与`ServerOptions.WebSocketOriginPatterns`匹配的主机，从而阻止针对`/wails/runtime`和 WebSocket 端点的 DNS 重绑定攻击
- Windows 通知在 WinRT Toast API 失败时不再回退到 PowerShell，而是返回错误

## v3.0.0-beta.30 - 2026-10-15

## 修复
- 修复更新程序尝试替换 AppImage 的只读挂载点，而不是`.AppImage`文件的问题
- 停止在 macOS 上发送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`会在每次运行循环迭代时触发，因此页面每秒会收到约120次空的`evaluateJavaScript:`调用；而在 macOS 上，每次27.0调用都会消耗用户点击的瞬时激活状态，导致`navigator.clipboard.writeText()`及其他需要激活状态的 API 对真实点击返回`NotAllowedError`并拒绝执行（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 完成
- 将 Windows 的默认 info.json 语言修复为 en-US（0409），以便资源管理器显示版本详细信息（#6210）
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中，通过分阶段复制回退方案处理 Linux 更新程序跨设备重命名失败的问题，由 @john-okeefe 完成
- 修复 Windows 上的`WebviewWindow.SetMenu`仅在配置了菜单栏主题时才重绘菜单栏的问题，解决[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修复 Windows 上的`WebviewWindow.SetMenu(nil)`会触发 panic，而不是像 macOS 上一样不执行任何操作的问题，解决[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送期间保留变更日志条目，由 @leaanthony 完成
- 移除引用模块中从未发布的二进制文件的嵌入项，修复`go mod vendor`在所有平台上都会因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，解决[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)；见[PR](https://github.com/wailsapp/wails/pull/6031)，由 @Grantmartin2002 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 上重复销毁窗口，由 @charlesbalderson 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中处理系统托盘菜单的 WM_CONTEXTMENU，使其可通过右键单击打开，由 @wzk2025 完成
- 修复全部五种`InvokeSync*`变体在分派的回调触发 panic 时会永久阻塞调用方的问题，解决[#6107](https://github.com/wailsapp/wails/issues/6107)；见[PR](https://github.com/wailsapp/wails/pull/6108)，由 @RALIST 完成

## v3.0.0-beta.29 - 2026-10-13

## 新增
- macOS 上的全局快捷键现在可感知键盘布局：在非 QWERTY 布局（AZERTY、QWERTZ、Dvorak 等）上，字母快捷键通过`UCKeyTranslate`绑定到标有相应字母的按键，而不是固定的 QWERTY 物理键位。

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6233)中，先初始化原生菜单，再在 Linux 上设置该菜单，由 @abhicodes-007 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6237)中解决 macOS 正式版就绪过程中输出处理和设置向导的回归问题，由 @taliesin-ai 完成
- 修复 Windows 上 WebView2 COM 的`AddRef`/`Release`包装器将过期的线程最后错误报告为失败的问题。它们现在会返回引用计数，因此释放操作不再会中止应用程序或导致资源请求失败（#3931）
- 修复`ICoreWebView2Deferral.Release`调用`AddRef`槽位，导致 Windows 上每个延迟资源请求的延迟对象泄漏的问题（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6225)中将 macOS 应用输出保留在文件提供程序之外，由 @taliesin-ai 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6224)中完善设置向导并显示系统详细信息，由 @taliesin-ai 完成
- 在[PR](https://github.com/wailsapp/wails/pull/6222)中改进 Windows 上 WebView2 恢复过程的错误处理和资源清理，由 @leaanthony 完成
- 改进 Windows 上的 WebView2 恢复机制，防止应用崩溃和控制器滞留，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6220)
- 更正 macOS 媒体授权指南，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6221)
- 修复右键单击 Linux 系统托盘图标时，在打开菜单的同时还会运行单击处理程序，导致每次右键单击都会切换所附加窗口的问题（#6018）
- 在 Windows 上保持 WebView2 回调处理程序存活，直至原生层将其释放，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6184)

## v3.0.0-beta.27 - 2026-10-01

## 修复
- 修复 macOS 12及更高版本忽略`Permissions`选项的问题：现已实现`WKUIDelegate`媒体捕获方法，因此`PermissionAllow`和`PermissionDeny`会像在 Linux 和 Windows 上一样应用于摄像头和麦克风请求。仍需配置`NSCameraUsageDescription`/`NSMicrophoneUsageDescription`；使用沙盒时，还需配置对应的设备权利

## v3.0.0-beta.26 - 2026-09-25

## 新增
- 监视器对前端开发服务器执行 HTTP 就绪检查，见 @atterpac 的[PR](https://github.com/wailsapp/wails/pull/6150)

## 修复
- 从 WebView2 进程故障中恢复，而不是留下空白窗口，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6002)

## v3.0.0-beta.25 - 2026-09-22

## 修复
- 主应用退出时停止`wails3 dev`及其后台进程（#6048）

## v3.0.0-beta.24 - 2026-09-20

## 安全
- 限制构建和发布工作流中的 GITHUB_TOKEN 权限，以遵循最小权限原则
- 修复 screen 示例资源中间件中的路径遍历和符号链接逃逸漏洞
- 将设置向导的依赖项安装限制为获准的软件包管理器命令，同时保持其正常工作

## v3.0.0-beta.23 - 2026-09-16

## 新增

- 允许在 Linux 上覆盖 GTK 应用 ID，见 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5972)

## 修复

- Linux 单实例机制以应用自身的 ID 声明总线名称，见 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5973)
- 取消 Linux 上已中止的自定义方案请求，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6126)

## v3.0.0-beta.22 - 2026-09-14

## 新增

- 完成翻译并迁移社区媒体内容——请参阅[文档网站](https://v3.wails.io)（已更新1812个页面），见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6122)
- 记录可嵌入的 Wails 按钮及其各种样式和尺寸，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6120)

## v3.0.0-beta.21 - 2026-09-13

## 新增

- 使用 M-Press 提供 Wails v3文档，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6116)

## 修复

- 解析 MPD 前置元数据中的 JSON slug 值以生成变更日志，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6118)
- 备份失败后，更新程序会清除辅助环境变量并重新启动原始目标，见 @cnmax 的[PR](https://github.com/wailsapp/wails/pull/6080)
- 在 App.Run 期间启动默认信号处理程序，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6098)
- Windows 菜单现在可处理 nil 菜单、释放被替换的资源并重绘菜单栏，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6112)
- 使用共享 YAML 配置恢复新项目的 MSIX 打包，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6115)
- 修复以下问题：当泛型模型创建器引用后续声明的辅助项时，生成的 JavaScript 和 TypeScript 绑定无法加载；创建相互依赖的泛型模型时发生栈溢出（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 变更

- 将 Clave 展示链接更新为当前网站和代码仓库，见 @01xR4in 的[PR](https://github.com/wailsapp/wails/pull/6082)

## 修复

- 取消已中止的 Windows 资源请求（包括 worker 请求），同时在导航期间保留 keepalive 处理程序。在 Apple 平台上，通过应用包装器转发原生请求上下文。（#5963、#5969）
- 通过重试机制，在相互竞争的推送之间保留变更日志条目，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6094)
- 移除引用模块中从未发布的二进制文件的嵌入项，修复`go mod vendor`在所有平台上因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，并修复[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)，见 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6031)

## 移除

- 移除已被纯 Go 加载器取代的原生 WebView2 加载器支持。此项变更会移除嵌入的`WebView2Loader.dll`二进制文件和`github.com/jchv/go-winloader`依赖项。`native_webview2loader`构建标签仍会被接受且不再报错，但对v3构建没有任何作用，见 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6031)
- 从 macOS API 指南中移除未使用的构建标签和 FPS 选项，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6097)

## v3.0.0-beta.19 - 2026-09-09

## 新增

- 使用构建标签控制私有 macOS API，仅在选择启用时使用——请参阅[文档](https://v3.wails.io/features/browser/integration)、[文档](https://v3.wails.io/features/environment/info)、[文档](https://v3.wails.io/features/windows/basics)、[文档](https://v3.wails.io/features/windows/frameless)、[文档](https://v3.wails.io/features/windows/notch-windows)、[文档](https://v3.wails.io/features/windows/options)、[文档](https://v3.wails.io/guides/build/macos)、[文档](https://v3.wails.io/guides/build/private-macos-apis)和[文档](https://v3.wails.io/reference/overview)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6087)

## 修复

- 拒绝超过64 MiB 的运行时请求并返回 HTTP 413，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6091)

## 安全

- 通过令牌身份验证强化 MCP 来源和远程访问安全性，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6092)

## v3.0.0-beta.18 - 2026-09-08

## 修复

- 通过使用指针接收器修复 Linux 和 Darwin 上的 Calloc 内存泄漏，见 @4RH1T3CT0R7 的[PR](https://github.com/wailsapp/wails/pull/6083)

## v3.0.0-beta.17 - 2026-09-06

## 修复

- Windows：WebResourceRequested 处理程序中失败或为 nil 的`GetRequest`不再导致进程终止（`log.Fatal`/nil 解引用 panic），而是丢弃该请求并记录日志，见 @midagedev 的[PR](https://github.com/wailsapp/wails/pull/6006)

## v3.0.0-beta.16 - 2026-08-29

## 变更

- 在新终端窗口中提示输入公证密码，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6029)

## 修复

- 正确处理 macOS 上的系统托盘点击类型，见 @ChewbaccaCookie 的[PR](https://github.com/wailsapp/wails/pull/5919)
- CI 在更新前移除未使用的 Microsoft apt 仓库，见 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6041)

## v3.0.0-beta.15 - 2026-08-27

## 修复

- 将 WebView2 嵌入超时时间增加到60秒，见 @Grantmartin2002 的[PR](https://github.com/wailsapp/wails/pull/6043)

## v3.0.0-beta.14 - 2026-08-26

## 修复

- 在 macOS 上正确命名 Control-字母组合键，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6032)
- 修复 Windows 上的 ICO 托盘图标，并使其跟随任务栏主题，见 @nik9play 的[PR](https://github.com/wailsapp/wails/pull/6016)

## v3.0.0-beta.13 - 2026-08-25

## 修复

- 在模态循环运行期间继续处理 macOS 主线程上的工作，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6026)
- 使移动端安全存储操作可失败，并在失败时保持关闭状态，见 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5923)
- 即使未注册监听器，也运行应用程序事件钩子，见 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5999)
- 更正注释和本地化文档中的拼写错误，见 @haoku123 的[PR](https://github.com/wailsapp/wails/pull/6023)
- 移除提交在`v3/examples`下的预编译 macOS 二进制文件，见 @4RH1T3CT0R7 的[PR](https://github.com/wailsapp/wails/pull/6025)

## v3.0.0-beta.12 - 2026-08-21

## 新增

- 添加带生命周期和遥测示例的 macOS 刘海通知窗口——参阅[文档](https://v3.wails.io/features/windows/notch-windows)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6010)
- 添加 macOS NSPanel 窗口支持，包括新选项和原生集成——参阅[文档](https://v3.wails.io/features/windows/options)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6008)

## 修复

- 防止并发调用导致 SQLite Prepare 挂起，见 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5998)

## v3.0.0-beta.11 - 2026-08-20

## 移除

- 从文档中移除已过时的实现跟踪器，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/6005)

## v3.0.0-beta.10 - 2026-08-19

## 修复

- 修复 GTK4 Linux 宿主丢弃自定义协议和文件关联启动参数的问题，见 @midagedev 的[PR](https://github.com/wailsapp/wails/pull/6000)
- 在变更日志验证中正确处理已删除的行和同源修正，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5993)

## v3.0.0-beta.9 - 2026-08-16

## 新增

- 添加安全的 wails3 mcp 服务器，用于代理辅助的项目管理，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5896)
- 添加绑定中模型的文档——参阅[文档](https://v3.wails.io/features/bindings/models)，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5988)
- 支持在原子化 Linux 系统上使用 rpm-ostree 安装，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5987)
- 添加原生的每周星标历史图表生成和发布功能——参阅[文档](https://v3.wails.io/credits)、[文档](https://v3.wails.io/de/credits)、[文档](https://v3.wails.io/fr/credits)、[文档](https://v3.wails.io/id/credits)、[文档](https://v3.wails.io/ja/credits)、[文档](https://v3.wails.io/ko/credits)、[文档](https://v3.wails.io/pt/credits)、[文档](https://v3.wails.io/ru/credits)、[文档](https://v3.wails.io/zh-cn/credits)和[文档](https://v3.wails.io/zh-tw/credits)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5986)
- 添加仅适用于 Darwin 的 mac 包，用于解析应用程序包资源——参阅[文档](https://v3.wails.io/guides/build/macos)，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5965)
- 添加 Condui 展示页面和索引条目——参阅[文档](https://v3.wails.io/community/showcase/condui)和[文档](https://v3.wails.io/community/showcase)，见 @mgueregath 的[PR](https://github.com/wailsapp/wails/pull/5962)
- 添加 Redis Viewer 展示页面，包括屏幕截图和项目链接——参阅[文档](https://v3.wails.io/community/showcase)和[文档](https://v3.wails.io/community/showcase/redisviewer)，见 @redisviewer 的[PR](https://github.com/wailsapp/wails/pull/5984)

## 变更

- 将 Linux 的 GTK 应用程序标志更新为 G<em>APPLICATION</em>NON_UNIQUE，见 @overlordtm 的[PR](https://github.com/wailsapp/wails/pull/5971)
- 将缺失窗口事件的日志级别从警告改为调试，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5914)

## 修复

- 允许已注册的 macOS 键盘快捷键优先于 WebView 处理，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5902)
- 修复文档侧边栏中失效的链接，见 @northes 的[PR](https://github.com/wailsapp/wails/pull/5937)
- 当 WebKit 中止匹配的自定义方案任务时，取消 macOS 和 iOS 资源请求上下文（#5963）
- 处理 WindowSetFullscreenButtonEnabled 消息，见 @archy-rock3t-cloud 的[PR](https://github.com/wailsapp/wails/pull/5976)
- 在 preact-ts 模板中导入 Fragment 以解决构建失败问题，见 @haoku123 的[PR](https://github.com/wailsapp/wails/pull/5979)
- 防止旧版 GTK3 纯服务应用程序在活动窗口或显示器尚不可用时执行屏幕发现而崩溃（#5966）
- 当未发布变更日志为空时，允许显式指定版本的发布运行继续执行（#5977）

## 安全

- 将网站 nanoid 锁文件更新到已修补的3.3.18版本，以解决安全公告中的问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5985)

## v3.0.0-beta.8 - 2026-08-12

## 新增

- 为自动生成的变更日志条目添加文档 URL 生成功能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5957)
- 添加 Streams：采用 WebSocket 编程模型，在 Go 与 JavaScript 之间提供双向字节流，且无需监听套接字。在 Go 中使用`app.HandleStream(name, handler)`声明流，并在前端使用`Stream(name)`连接；该调用返回一个形如`WebSocket`的对象。Go→JS 通过资源服务器上每个窗口一个保持中的轮询传输，JS→Go 则通过普通 POST 传输；不会绑定任何 TCP 端口，也不会有任何内容经过`evaluateJavaScript`。在服务器构建（`-tags server`）中，同一处理程序改为通过真正的 WebSocket 提供服务，因此应用程序代码在不同构建中完全相同。作者：@leaanthony
- 将 mailbox 变更日志条目移至“未发布”，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5935)

## 变更

- 更新文档侧边栏自动生成和博客作者类型派生逻辑，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5938)

## 修复

- WebView2 初始化使用截止时间和消息泵，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5952)
- WebView2 Cookie 测试在 CI 中默认跳过，除非明确启用；测试执行会锁定到当前操作系统线程，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5951)
- Windows 菜单构建器恢复子菜单父项的命令 ID，见 @gilad-ch 的[PR](https://github.com/wailsapp/wails/pull/5944)
- 使官方交叉编译镜像与 GTK 4.14及 Linux 支持基线保持一致（#5928）
- 配置 iOS Xcode 项目以保留继承的链接器标志并添加 -ObjC，见 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5915)
- 修复大型前端中`wails3 dev`资源代理的 TCP 连接过度频繁更换问题；该问题可能耗尽主机的临时端口，导致无关进程因`EADDRNOTAVAIL`而失败
- 按窗口对事件 JavaScript 进行排队，以实现有序分发和背压，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5934)

## 移除

- 移除桌面二进制文件发布流水线：v3发布仅包含标签，`wails3` CLI 通过`go install`安装。删除`release-v3.yml`以及触发它的夜间构建步骤，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5946)

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 新增 macOS 自动播放首选项，以停用媒体播放必须由用户操作触发的要求，见 @Eyalm321 的[PR](https://github.com/wailsapp/wails/pull/5512)
- 将邮箱的变更日志条目移至“未发布”部分，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5935)

## 变更

- macOS 缩放动画改用 CADisplayLink 或 NSTimer，以获得更流畅的性能，见 @savely-krasovsky 的[PR](https://github.com/wailsapp/wails/pull/5945)

## 修复

- 配置 iOS Xcode 项目以保留继承的链接器标志并添加 -ObjC，见 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5915)
- 修复大型前端中`wails3 dev`资源代理的 TCP 连接过度频繁更换问题；该问题可能耗尽主机的临时端口，导致无关进程因`EADDRNOTAVAIL`而失败
- 按窗口对事件 JavaScript 进行排队，以实现有序分发和背压，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5934)

### 新增

- 实现通用异步 FIFO 邮箱，以便按顺序传递事件，见 @savely-krasovsky 和 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/5851)

## v3.0.0-beta.6 - 2026-08-09

## 新增

- 为超大事件实现有界的主机端存储和有序的 JavaScript 传递，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5930)
- 实现 macOS Dock 弹跳效果，用于窗口闪烁提醒，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5921)

## 修复

- 资源服务器在刷新期间保留内容类型嗅探器错误和未写入的前缀，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5931)
- 防止从 Wails 回调中替换应用程序菜单时 macOS 应用程序崩溃
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨读的问题。深色模式 uxtheme 导出受内部版本18334限制，因此应用级深色模式选择加入逻辑从未在这些主机上运行：菜单背景被绘制成深色，但 Windows 仍按浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号导出自17763起便已存在，因此现已相应调整该限制条件。
- 修复`w32.GetStockObject`调用`GetDeviceCaps`而非`GetStockObject`的问题；该问题导致它对每个预定义对象都返回0。
- 改进 WebView2 引导程序下载错误的处理和报告，见 @jannskiee 的[PR](https://github.com/wailsapp/wails/pull/5924)

## v3.0.0-beta.5 - 2026-08-07

## 修复

- macOS 应用激活现在仅对常规应用遵循激活策略，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5897)
- 在 Linux 构建中防护未初始化的 GTK 窗口，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5898)
- 在加载 URL 前，为 Linux WebKit 窗口显式设置不透明背景色，见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/5899)

## v3.0.0-beta.4 - 2026-08-05

## 变更

- Android 构建任务默认使用 arm64，deploy-emulator 会选择主机架构，见 @mortenolsrud 的[PR](https://github.com/wailsapp/wails/pull/5890)

## 修复

- 在拖动过程中保留 macOS 窗口的缩放状态并减少动画效果，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5900)
- 通过将`!server`添加到`webview_window_windows_nonclient.go`构建约束中，修复 Windows 服务器模式构建

## v3.0.0-beta.3 - 2026-08-03

## 新增

- 在实现细节中记录第10阶段 Beta 验证已完成，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5881)

## 修复

- 将窗口句柄传递给 Windows 深色模式 API 并验证参数，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5877)
- 集中解析 macOS 无边框窗口的标题栏按钮状态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5870)
- 防止 Windows 应用程序在 Windows 应用主题为浅色时请求深色模式而导致原生菜单文本无法辨读。现在，在 Windows 能够渲染深色菜单文本之前，菜单会使用与之匹配的浅色原生背景。
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨读的问题。深色模式 uxtheme 导出受内部版本18334限制，因此应用级深色模式选择加入逻辑从未在这些主机上运行：菜单背景被绘制成深色，但 Windows 仍按浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号导出自17763起便已存在，因此现已相应调整该限制条件。

## v3.0.0-beta.2 - 2026-08-02

## 变更

- 将v3从 Alpha 提升为 Beta
- 记录系统托盘的智能默认值和弹出窗口自动隐藏行为，并为点击处理程序的选择添加回归测试覆盖（#5840）。
- GitHub 更新程序默认排除 Windows 安装程序资源，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5861)
- 为 macOS 无边框窗口新增圆角、直角及自定义圆角半径支持，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5866)

## 修复

- 报告实时 GTK4 窗口尺寸，并从已配置的表面发出调整大小、最大化、最小化及全屏状态事件（#5830）。
- 修复 Linux WebKit 在 fetch 请求中发送 Blob 或 FormData 时崩溃的问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/5854)
- 当 Blob/FormData 标头缺失时，Fetch shim 会传递 undefined，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5865)

## v3.0.0-alpha2.122 - 2026-08-01

## 新增

## 变更

- 为 macOS 无边框窗口新增圆角、直角及自定义圆角半径支持，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5866)

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5865)中使 fetch shim 对缺失的 Blob/FormData 标头传递 undefined

## v3.0.0-alpha2.121 - 2026-07-31

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5857)中新增 macOS DMG 打包支持，包括新选项和构建任务

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5861)中使 GitHub 更新程序默认排除 Windows 安装程序资产

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5854)中修复 Linux WebKit 在 fetch 请求中发送 Blob 或 FormData 时崩溃的问题

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5853)中实现 macOS 标题栏双击时最大化或最小化的操作

## 变更

- 由 @jeongkyu 在[PR](https://github.com/wailsapp/wails/pull/5849)中更新二维码服务教程，改用 NewServiceWithOptions 并添加间距

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5856)中使 WKWebView 在 macOS 缩放期间保持响应
- 修复 GTK4 窗口大小查询，并从配置的`GdkSurface`发出调整大小、最大化、最小化和全屏状态事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5833)中更新多语言文档，加入架构图

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5825)中为图标生成的输入和输出提供默认路径
- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5797)中将源入口模块添加到运行时 package.json 的 sideEffects 中
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5816)中向贡献指南添加许可证和来源信息章节

## 修复

- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5800)中应用限定作用域的 GTK4 无边框 CSS 以移除圆角
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5789)中使 Windows 上的弹出菜单和屏幕枚举能够妥善处理光标位置获取失败
- 由 @phergul 在[PR](https://github.com/wailsapp/wails/pull/5678)中使 macOS 打开文件对话框正确筛选扩展名，并按后缀验证允许的文件
- 由 @roachadam 在[PR](https://github.com/wailsapp/wails/pull/5793)中防止 Windows 深色模式初始化调用 nil API
- 修复更新程序的32位构建失败：在`GOARCH=386`上将`maxArchiveTotalSize`常量（2 GiB）传给`fmt.Errorf`时，该常量超出了平台`int`的范围。现在已将其显式指定为`int64`类型。
- 修复启动时的 nil 指针 panic：在未加载深色模式 uxtheme API 的 Windows 构建（例如 Windows 10 1809 / Windows Server 2019（内部版本17763））中，当窗口使用深色（或跟随系统的深色）标题栏时会发生此问题。现在，窗口主题设置中的`AllowDarkModeForWindow`调用已添加 nil 检查，与`w32.SetMenuTheme`中已有的检查保持一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 由 @savely-krasovsky 在[PR](https://github.com/wailsapp/wails/pull/5462)中为 Windows 上的非客户区实现自定义命中测试逻辑

## 变更

- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5761)中根据 UseVisualHosting 配置 WebView2 的显示器缩放检测

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5763)中更新常见问题文档，使其重点介绍 Wails v3的功能和指南

## v3.0.0-alpha2.115 - 2026-07-06

## 修复

- 修复`Menu.Update()`无法在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 修复显示配置变化时枚举 macOS 屏幕导致的崩溃，方法是复制屏幕 id/name 字符串并保存数量快照（#5565；@x-haose 在 #5584 中独立诊断并修复）
- 修复 Windows 上的崩溃：当`WM_ERASEBKGND`在最小化/恢复过渡期间绘制纯色背景，而`GetClientRect`返回 nil 时会发生此问题（防护措施由 @sinspired 在 #5636 中报告）
- 由 @mbaklor 在 #5690 中修复前端绑定错误始终按文本解析的问题
- 修复使用`server`构建标签时的 Windows 构建失败；其原因是 Windows GUI 文件缺少 macOS 和 Linux 对应文件已有的`!server`构建约束（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5720)中实现更新清单协议和端点提供程序

## 变更

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5711)中将`webview2`绑定作为`v3/internal/webview2`并入v3模块，移除独立模块、其每夜版发布/同步工作流以及 go.mod 版本调整流程（v3是其唯一使用方）

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5750)中将 WebView2 显示器缩放检测和 DPI 变化时的宿主重新同步修复移至“未发布”部分
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5741)中更新 WebView2 对 float64 和 BOOL 参数的 COM 编组
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5703)中防止 Windows 系统托盘图标更新和销毁时发生 panic 和 nil 解引用
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5743)中修复隐藏窗口在 Windows 上无法正确再次隐藏的问题
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5742)中使 WebView2 控制器的可见性与窗口最小化、最大化和恢复操作同步

### 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5734)中重新启用 WebView2 显示器缩放检测，并仅在 DPI 变化时执行宿主重新同步；此变更基于 @randalmurphal 验证的修复，并由 @eleclin 验证根本原因、@qq540491950 进行硬件测试

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中新增警告：构建发布版 AAB 时若未设置`ANDROID_KEYSTORE_FILE`，则发出警告（Google Play 会拒绝使用调试签名的 bundle）；同时记录 App Bundle 打包和签名方法
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5739)中新增多语言版“为何选择 Wails”文档
- 由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5398)中支持在绑定中将 Go time.Time 映射为 JS Date 或 string
- 新增用于提交到 Play Store 的 Android App Bundle（AAB）打包任务（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 任务仍用于本地/模拟器测试。由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5728)中实现（修复[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 新增 Android 物理设备任务目标，并恢复相机/位置权限。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5735)中实现

## 变更

- 将`webview2`升级到v1.0.28（[发行说明](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 将 Android 模板的`compileSdk`/`targetSdk`从34升级到35，这是 Google Play 对新应用提交的要求。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中实现

## 修复

- 修复 sponsorkit 中已烘焙的头像蒙版。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5745)中实现
- 修复 Android AVD 自动创建功能因按字典序对版本排序而选错系统映像或 cmdline-tools 版本的问题。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中实现
- 修复设置向导建议使用过时 Android NDK 版本的问题（现为26.3.11579264，与文档所述要求一致）。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中实现
- 更新有关 SvelteKit 和选项的法语文档。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5744)中实现
- 修复显示配置变化期间枚举 macOS 屏幕时发生 SIGSEGV 的问题。由 @flofreud 在[PR](https://github.com/wailsapp/wails/pull/5516)中实现

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 新增基于 Go 的贡献者 SVG 生成器，并更新 docs/website 的鸣谢页面。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5724)中实现
- 支持在绑定中将 Go time.Time 映射为 JS Date 或字符串。由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5398)中实现

## 变更

- 使用 Go 生成器替换基于 Node 的赞助者图像处理流水线。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5719)中实现

## 修复

- 修复 Android 构建资产依赖项安装脚本。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5729)中实现
- 在`ValidateAndSanitizeURL`中拒绝 U+0085（NEXT LINE）控制字符，从而补全 URL 验证器对空白字符的覆盖
- 在无边框窗口的 DPI 变化时重新计算 DWM 框架。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4785)中实现
- 修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题。由 @yulesxoxo 在[PR](https://github.com/wailsapp/wails/pull/4632)中实现
- 为 Darwin 对话框、菜单、托盘和通知中的 Cocoa 对象新增显式 Objective-C 内存管理。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5714)中实现
- 修复 Linux CGO 后端错误和系统托盘问题。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5718)中实现

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 将 HappyTools 添加到社区展示中。由 @Aliuyanfeng 在[PR](https://github.com/wailsapp/wails/pull/5061)中实现
- 新增印度尼西亚语区域设置支持和完整文档。由 @triadmoko 在[PR](https://github.com/wailsapp/wails/pull/5643)中实现
- 为 WindowsWindow 新增 DisableMenu 选项。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4813)中实现

## 变更

- 更新 Taskfile 模板和 CLI，以使用 GOOS 和 ARCH 分派构建/打包任务。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5617)中实现

## 修复

- 修复 Mac 窗口标签页功能的问题。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5708)中实现

## 移除

- 从贡献、功能和指南中移除德语翻译的 MDX 文件。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5702)中实现

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 实现 macOS WebView 重新加载和强制重新加载，并新增 WebContent 进程终止后的恢复机制。由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5129)中实现
- 新增有关贡献、功能和指南的完整德语文档。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5396)中实现
- 增强通知功能，新增声音、附件、调度和更新 API。由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/5333)中实现

## 修复

- 在无边框窗口的 DPI 变化时重新计算 DWM 框架。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4785)中实现
- 修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题。由 @yulesxoxo 在[PR](https://github.com/wailsapp/wails/pull/4632)中实现

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 向 EventsEmit 文档添加代码示例。由 @iamhabbeboy 在[PR](https://github.com/wailsapp/wails/pull/5026)中实现
- 新增 Windows WebView2 可视化托管选项。由 @MerIijn 在[PR](https://github.com/wailsapp/wails/pull/5380)中实现
- 将 Klustr 添加到社区展示文档中。由 @SametKUM 在[PR](https://github.com/wailsapp/wails/pull/5536)中实现
- 将 Kira 添加到社区展示中，并新增相关页面和变更日志条目。由 @thiennguyen93 在[PR](https://github.com/wailsapp/wails/pull/5685)中实现
- 向 MCP 服务指南添加反馈部分。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5694)中实现

## 变更

- 服务器模式现在具备一等支持的生产构建，与桌面端构建任务保持一致（#5693）。`task build:server`默认构建生产二进制文件（`-tags server,production`、`-trimpath`、已剥离符号），并接受`DEV=true`（开发服务器）、`OBFUSCATED=true`（garble）和`EXTRA_TAGS`。`task run:server`运行开发服务器。`Dockerfile.server`/`task build:docker`构建生产服务器（`-tags server,production`），并先构建生产前端；该映像默认是在 distroless/static 上构建的纯 Go 静态版本，同时将`CGO_ENABLED`、`GO_IMAGE`和`RUNTIME_IMAGE`公开为可覆盖的构建参数，以供 CGO 应用使用。

## 修复

- 防止在存在待处理异步调用时关闭窗口导致崩溃。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4435)中实现
- 防止在 Windows 上打开隐藏应用时激活窗口。由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5249)中实现
- 确保 WebKit 请求元数据、响应完成处理和正文流处理在 GTK 主线程上运行。由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5668)中实现
- 修复`Menu.Update()`未在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对数量创建快照，修复显示配置变化时枚举 macOS 屏幕导致崩溃的问题（#5565；@x-haose 在 #5584 中独立诊断并修复）
- 通过在`WM_DPICHANGED`处理程序中重新设置控制器边界，与取消最小化时的 DPI 重新同步保持一致，修复在 Windows 上将窗口拖过不同 DPI 的显示器后 WebView2 内容先缩小再消失的问题（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 通过`app.GlobalShortcut`（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）添加全局（系统级）键盘快捷键。即使应用程序未获得焦点，快捷键也会触发。各平台均采用原生实现，不依赖第三方库：macOS 使用 Carbon 热键，Windows 使用`RegisterHotKey`，X11 使用`XGrabKey`，Wayland 使用 XDG Desktop Portal 全局快捷键接口。
- 添加内置 MCP 服务器：使用`mcp`标签构建应用程序时，会自动启动一个模型上下文协议服务器，使 LLM 智能体能够测试和控制正在运行的 Wails 应用程序，包括控制窗口、检查 DOM、执行 JavaScript、调用绑定方法、处理事件，以及模拟鼠标和键盘输入；模拟输入会通过屏幕上的动画光标呈现。无需编写用户代码：设置`WAILS_MCP=1`后，`wails3 build`/`wails3 dev`会自动添加`mcp`标签。完全通过环境变量（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）进行配置。

## 修复

- 修复`Menu.Update()`未在 GTK4 Linux 上重新构建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对数量创建快照，修复显示器发生变化时枚举 macOS 屏幕导致的崩溃（#5565；@x-haose 在 #5584 中独立诊断并修复）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5613)中添加带侧边栏导航的实验性 Wake 文档

## v3.0.0-alpha2.106 - 2026-06-24

## 变更

- 将`webview2`升级到v1.0.27。
  - ci(webview2)：修复发布构建（交叉编译 Windows + 完整的 go.sum）（#5671）\

  **完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5672)中从 webview2 发布工作流的交叉编译步骤移除 go vet
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5670)中将自动生成变更日志所用的 OpenRouter 模型更新为 google/gemini-2.5-flash-lite
- 将`webview2`升级到v1.0.26。

### 修复

- **从瞬态运行时 COM 错误中恢复，而不是退出**（#5658、#5580）。此前，`Chromium.errorCallback`会对*任何* COM 错误调用`os.Exit(1)`，因此启动后一次可恢复的短暂故障就会终止整个应用程序。运行时路径（`Resize`/`GetClientRect`、`Navigate`/`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）现在会记录错误并恢复。特别是，`MessageReceived`中格式错误或不受信任的 Web 消息现在会被丢弃，而不会导致进程终止。这解决了跨越混合 DPI 显示器时发生的一类崩溃（#5544、#5650）。环境/控制器创建路径仍会导致致命错误。\

**完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5671)中修复 release-webview2 工作流，使其能够正确处理 go.sum 文件
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5659)中通过清除并重新构建原生菜单，修复 Linux GTK4 菜单更新

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 添加`application.System`，以便从共享代码中检测运行时平台：`System.IsMobile()`（iOS/Android）、`System.IsDesktop()`（macOS/Windows/Linux）、`System.IsServer()`（`server`构建标签），以及用于直接测试单个目标的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在所有目标上编译，因此无需构建标签即可创建条件分支。`@wailsio/runtime`中还提供了对应的前端辅助函数（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）
- 添加“使用其他前端框架”指南，说明如何将自己的 Vite 项目直接放入`frontend/`（涵盖 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`向导现在会检查移动端（iOS/Android）工具链，包括 Xcode 和 iOS Simulator 运行时、JDK、Android SDK/NDK 及模拟器；适用时还提供一键安装和可复制的 Shell 配置修复命令
- 生成的项目附带一个`frontend/.npmrc`，其中设置了7天的`minimum-release-age`，以降低刚发布的（可能已遭入侵的）软件包带来的风险（pnpm 和 bun 会遵循此设置；npm 会忽略它，但不会产生不良影响）

## 变更

- 采用全新的霓虹山峦主视觉，重新设计所有内置入门模板（Web、iOS 和 Android）
- **TypeScript 现在是入门模板的默认语言，并占用不带语言后缀的模板名称。** `wails3 init`（不含`-t`）会搭建 TypeScript 项目；`-t vanilla`、`-t react`、`-t vue`和`-t svelte`均为 TypeScript 模板，对应的 JavaScript 变体分别为`-t vanilla-js`、`-t react-js`、`-t vue-js`和`-t svelte-js`。内置模板通过`template.yaml`中的`typescript:`声明所用语言；使用`-ts`后缀的社区模板仍可作为回退方案继续工作
- 使用霓虹“digital Wails”主题重新设计`wails3 setup`向导（山峦背景上的磨砂玻璃鲜活效果）

## 修复

- 修复 Windows 上的崩溃：当应用最小化的时间足以让 WebView2 挂起或回收其渲染/GPU 进程时，恢复应用会触发此崩溃。最小化/恢复时的 DPI 重新同步（#5544）现在仅在窗口 DPI 确实发生变化时才操作 WebView2 控制器，从而避免在常见的相同 DPI 恢复场景中对已挂起的控制器执行会导致致命错误的 COM 调用（#5605）
- 修复长期运行且频繁加载资产/媒体的 Linux 应用中反复发生的原生`SIGABRT`/`SIGSEGV`崩溃（通常发生在 GTK 主循环期间的`g_object_unref`内部）。资产服务器从工作 goroutine 完成`WebKitURISchemeRequest`，因而在 GTK 主线程之外调用了线程不安全的 WebKit2GTK 函数；现在，完成操作（`webkit_uri_scheme_request_finish_with_response`/`finish_error`）会在主线程上运行。此项工作完成了 #5566中的部分修复。GTK3 和 GTK4/WebKitGTK 6.0构建均受影响（#5631、#5557）
- 修复 Linux/GTK3 上`setupSignalHandlers`中间歇性出现的`fatal error: invalid pointer found on stack`。作为信号`user_data`传递的窗口 ID 保存在 Go 的`unsafe.Pointer`局部变量中，因此垃圾回收器在复制栈时扫描到这个（非指针）值便会中止。现在 Go 端将该 ID 保持为整数类型（`uintptr_t`），把 #4958对 GTK4 路径所做的同一修复向后移植到旧版 GTK3 路径（该修复将 C 信号函数改为`uintptr_t`，以消除`-race`/checkptr 错误）（#5631）

## 移除

- 移除`react-swc`、`preact`、`lit`、`solid`、`qwik`和`sveltekit`入门模板（及其`-ts`变体）。目前支持的内置模板为`vanilla`、`react`、`vue`和`svelte`——每个模板默认均使用 TypeScript，并提供`-js` JavaScript 变体。仍可通过[自带前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自定义模板使用任何其他框架

## v3.0.0-alpha2.104 - 2026-06-18

## 修复

- 修复绑定的 Go 服务方法返回空字符串时发生的 iOS 崩溃（SIGABRT）。iOS 资产响应写入器使用`buf != nil`而非正文长度来保护正文指针，因此长度为零的正文会导致`&buf[0]`发生 panic；现在改为按长度进行保护，与桌面端写入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 变更

- 将 iOS 和 Android 原生功能迁移到平台管理器：通过`application.IOS.*`和`application.Android.*`调用它们（例如`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），不再使用旧的`application.IOS*`/`application.Android*`自由函数（#5602）
- 重命名移动端桥接事件：跨平台事件现在使用`common:*`前缀（例如`common:haptic`、`common:location`），平台专属事件使用`ios:*`/`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）；不再使用`native:*`前缀（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增实验性的`wails3 setup`向导，用于交互式项目设置和依赖项检查
- 为`wails3 doctor`新增`--json`标志，以提供机器可读的输出
- 在`wails3 doctor`命令中新增签名状态部分

## 修复

- 修复 Linux 上的 npm 检测，使其除检查包管理器外还检查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生消息对话框（UIAlertController）以及打开文件、多个文件和目录的对话框（UIDocumentPickerViewController）；保存对话框会返回明确的错误
- iOS：通过 UIPasteboard 支持剪贴板
- iOS：通过 UIScreen 获取真实屏幕指标（点、像素、缩放比例和安全区域内的工作区）
- iOS：设备构建（`IOS_PLATFORM=device`）、代码签名身份/预置描述文件/权利支持、`.ipa`打包，以及通过 devicectl 实现`deploy-device`
- iOS：可配置最低 iOS 版本（build/config.yml 中的`ios.minIOSVersion`）
- iOS：`wails3 doctor`会在 macOS 上报告 Xcode 和 iOS SDK 的可用性
- iOS：系统事件——电池、网络、主题、屏幕锁定和内存不足事件以`events.IOS.*`及平台中立的`events.Common.*`应用事件形式公开
- iOS：原生移动功能桥接（导出的`application.IOS*`）——共享面板、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（Face ID/Touch ID）、本地通知和 Keychain 安全存储
- iOS：传感器与硬件——触觉反馈、单次地理定位、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距和屏幕捕获检测
- iOS：文档（IOS.md 和文档站点指南）
- Android：原生消息对话框（AlertDialog）以及打开文件和多个文件的对话框（Storage Access Framework，导入为缓存副本）；打开目录和保存对话框会返回明确的错误
- Android：通过 ClipboardManager 支持剪贴板
- Android：通过 WindowMetrics/DisplayMetrics 获取真实屏幕指标（dp、像素、缩放比例和系统栏以内的工作区）
- Android：触觉反馈（`Android.Haptics.Vibrate`）、设备信息（`Android.Device.Info`）和 Toast 提示（`Android.Toast.Show`）运行时方法
- Android：带类型的生命周期事件（`events.Android.*`，由 events.txt 生成），其中`ActivityCreated`映射到`Common.ApplicationStarted`
- Android：构建流水线会生成可安装的调试版和发布版 APK（`android:run`、`android:package`、`android:package:fat`）；发布签名默认使用调试密钥库，也可通过`ANDROID_KEYSTORE_*`环境变量使用真实密钥库
- Android：`wails3 doctor`会报告 Android SDK、NDK 和 JDK
- Android：系统事件——电池、网络、主题、屏幕锁定和内存不足事件以`events.Android.*`及平台中立的`events.Common.*`应用事件形式公开
- Android：原生移动功能桥接（导出的`application.Android*`）——共享、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（BiometricPrompt）、本地通知和 EncryptedSharedPreferences 安全存储
- Android：传感器与硬件——触觉反馈、单次地理定位、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距，以及使用 FLAG_SECURE 阻止屏幕捕获
- Android：文档（ANDROID.md 和文档站点指南）
- 示例：`mobile`综合示例新增 Mobile 和 Hardware 选项卡，用于演示 iOS 和 Android 上的原生功能桥接（胶囊式选项卡可换行显示为多行）
- 移动端：电池——应用进入后台时会暂停加速度计、接近传感器、手电筒和示例中的周期性时钟，并在返回前台时恢复（Android 会让进程继续在后台运行，而手电筒是在 iOS 上持续保持的硬件状态）；此外，Android 系统事件接收器仅在应用位于前台时注册
- iOS：相机拍摄——`application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → 携带 base64 缩略图的`native:capture`事件）
- iOS：后台执行——`application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication 后台任务时间窗口），以及可配置的`ios.backgroundModes`（build/config.yml），用于将`UIBackgroundModes`模板化写入生成的 Info.plist
- Android：相机拍摄——`application.AndroidCapturePhoto`/`AndroidCaptureVideo`（通过 FileProvider 调用系统相机 → `native:capture`事件）
- Android：前台服务——`application.AndroidStartForegroundService`/`AndroidStopForegroundService`（带有持续通知的`WailsForegroundService`可让进程保持运行，以执行长时间运行的后台工作）
- 示例：新增 Camera 选项卡，用于演示照片/视频拍摄和后台执行（Android 使用前台服务，iOS 使用后台任务时间窗口）

## 修复

- 修复 Linux 上`getUserMedia`总是因`NotAllowedError`而失败的问题：WebKitGTK 会拒绝无人处理的权限请求，而`permission-request`信号此前未连接。现在通过新的跨平台`WebviewWindowOptions.Permissions`映射（`map[PermissionType]Permission`）处理摄像头/麦克风权限，Linux（WebKitGTK）和 Windows（WebView2）均会遵循该映射。在没有原生提示框的 Linux 上，摄像头/麦克风默认允许使用（恢复`getUserMedia`），并可通过`PermissionDeny`将其关闭（#5552）
- iOS：`GOOS=ios`可再次编译（导出`events.IOS`并添加移动端方法名存根），带 production 标签的构建也可编译（修复 pkg/application 和多个服务中的构建标签）
- iOS：Go→JS 事件和 ExecJS 现已正常工作——页面启动时不再加载两次，`wails:runtime:ready`握手也不会再丢失
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted`不再与应用启动发生竞态；移除了固定2秒的启动等待
- iOS：修复每次执行 Go→JS JavaScript 时都会发生的 C 字符串泄漏
- iOS：`hasListeners`现在会反映真实的监听器注册状态
- iOS：生产构建不再编译框架调试日志代码
- Android：`GOOS=android`可再次编译——定义了`events.Android`，移除了越界的`events_android.go`监听器数组，添加了移动端方法名存根，并阻止桌面 Linux 文件（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 构建
- Android：JS→Go 绑定现已正常工作——WebView 无法将`fetch()` POST 请求体传递给`shouldInterceptRequest`，因此运行时调用现在通过 JavascriptInterface 传输（`nativeHandleRuntimeCall`）路由，而不会因请求体为 nil 而崩溃
- Android：`Screens.*`运行时调用现在会返回真实数据——ScreenManager 现已在启动时填充（此前从未接入，因此`GetAll`返回 nil）
- Android：生产构建不再编译框架调试日志代码；在调试构建中，日志会以`Wails`标签输出到 logcat
- Android：实现真正的`hasListeners`注册表、JNI 引用/异常处理，以及仅加载一次的页面生命周期（不再重复导航）
- 修复在 Windows 上运行 Vite 开发服务器时`wails3 generate bindings`因“拒绝访问”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5515）
- 修复 macOS 上显示器发生变化后读取屏幕信息时偶发的致命崩溃：所存储的屏幕 ID 和名称指针指向自动释放的`UTF8String`缓冲区，这些缓冲区可能在 Go 复制其内容之前被释放（释放后使用）。现在会对字符串执行`strdup`，并在转换后将其释放；屏幕枚举也会在显式的自动释放池中运行，因此从 Go goroutine 调用时不再发生泄漏（#5556）
- 修复 Linux 上 assetserver 关闭`WebKitURISchemeRequest`时偶发的 SIGSEGV：最后一个`g_object_unref`在 assetserver goroutine 上运行，导致 WebKit GObject 在 GTK 主线程之外被终结。现在会通过`g_main_context_invoke`将 unref 操作编排到 GTK 主上下文中执行（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 新增

- 为`MacWebviewPreferences`扩展更多 WKWebView 配置选项：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize`和`ApplicationNameForUserAgent`（#5549）

## 修复

- 修复在 Vite 开发服务器运行时，`wails3 generate bindings`在 Windows 上因“拒绝访问”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录中的文件（#5561）
- 修复 Linux 上无边框窗口不触发 JS 调整大小事件的问题；修复无边框窗口的滚动条边缘检测（#5368）
- 修复 Windows 更新程序在临时目录与安装目录位于不同卷时因“无效的跨设备链接”而失败的问题（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修复

- 修复在 Vite 开发服务器运行时，`wails3 generate bindings`在 Windows 上因“拒绝访问”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录中的文件（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修复

- 修复 Linux 上 WebKit 空闲时（例如检查器处于打开状态）界面冻结的问题：不再对`SIGUSR1`强制设置`SA_ONSTACK`，因为这会破坏 JavaScriptCore 的 GC 线程同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 新增

- 新增调试页面以及使用`runtime/trace`的相关内容

## 变更

- 移除一些不必要的`_ "embed"`导入，略微整理代码

## 修复

- 修复 Windows 上窗口取消最大化后未强制执行最小宽度/高度限制的问题（#4593）
- 修复使用 Frameless + Transparent 窗口选项时，全屏模式下鼠标点击会穿透的问题（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 新增

- 新增 Garble 混淆支持（[#4563](https://github.com/wailsapp/wails/issues/4563)）：稳定的绑定方法 ID、构建/Taskfile 配套支持（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），以及所有面向运行时的有效载荷上的 JSON 结构体标签（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），使线格式在 Garble 重命名导出字段后仍能保持不变。

## v3.0.0-alpha.95 - 2026-05-20

## 新增

- 补充缺失的项目结构页面

## 变更

- 文档：将架构页面上的几幅图改为使用时序图，以获得更清晰的显示效果
- 文档：补充说明运行前需安装 D2

## 修复

- 修复 GTK4 默认构建中的`wails3 generate appimage`：打包器现在会先从二进制文件检测 GTK 技术栈，再搜索运行时文件，因此 GTK4 构建会选择`webkitgtk-6.0/`下的`libwebkitgtkinjectedbundle.so`，而`-tags gtk3`构建会选择`webkit2gtk-4.1/`下的`libwebkit2gtkinjectedbundle.so`。`.relr.dyn`探测现在也会检查`libgtk-4.so.1`，因此无论使用哪种技术栈，现代工具链都会正确禁用剥离。（#5475）
- 修复使用相对`-builddir`调用`wails3 generate appimage`时失败的问题：打包器现在会预先将`-binary`、`-icon`、`-desktopfile`、`-builddir`和`-outputdir`解析为绝对路径，因此流程中途执行`s.CD`时，不会再破坏 AppRun 下载 goroutine 或复制后的`ldd`探测。
- 修复桌面文件的`Name=`字段与二进制文件的基本名称不匹配时，`wails3 generate appimage`无法将最终 AppImage 移至`-outputdir`的问题：打包器现在会通过`OUTPUT`环境变量强制 linuxdeploy 的 appimage 插件将 AppImage 写入`<binary>-<arch>.AppImage`，而不是使用从桌面文件派生的名称。
- 修复在 alpha.93中将 GTK4 + WebKitGTK 6.0技术栈提升为默认技术栈后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`在 Linux 上不触发的问题。新的默认`application_linux.go` `run()`没有调用`setupCommonEvents()`（该函数会将`Linux.*`事件转发给对应的`Common.*`事件）或`monitorPowerEvents()`。现在通过`application_linux_dbus.go`在 GTK3 和 GTK4 构建路径之间共享 DBus 电源监视器辅助程序。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修复

- 修复在 alpha.93中将 GTK4 + WebKitGTK 6.0技术栈提升为默认技术栈后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep`和`Common.SystemDidWake`在 Linux 上不触发的问题。新的默认`application_linux.go` `run()`没有调用`setupCommonEvents()`（该函数会将`Linux.*`事件转发给对应的`Common.*`事件）或`monitorPowerEvents()`。现在通过`application_linux_dbus.go`在 GTK3 和 GTK4 构建路径之间共享 DBus 电源监视器辅助程序。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中新增`XDG_SESSION_TYPE`

## 修复

- 由 @leaanthony 修复 appmenu-gtk-module 访问尚未实现的窗口而导致 Wayland 上窗口菜单崩溃的问题（#4769）
- 由 @leaanthony 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题
- 由 @overlordtm 修复 Windows 上初始化拖放功能时出现“内存不足”错误的问题（#4701）
- 由 @leaanthony 修复主线程回调存储在删除映射条目时错误使用 RLock 所导致的竞态条件（Linux、macOS、iOS）（#4424）
- 修复向任务传递命令行参数时的变量处理。现在，以 KEY=VALUE 对形式指定的 CLI 变量会被正确初始化，并在整个任务执行过程中传播。
- 修复 macOS 上的 NSWindowZoomButton 冲突：`MaximiseButtonState`和`FullscreenButtonState`现在都会在启动时和运行时应用限制更严格的状态；任一 setter 都无法再静默覆盖另一个（#5319）
- 修复 CodeRabbit 在 #5463中发现的旧版 GTK3 构建路径（`-tags gtk3`）中的一组既有错误：通过文件关联启动时不再跳过启动处理程序；`getTheme`现在可确保边界和类型安全；`appName`不再释放 GLib 所有的内存；`clipboardGet`不再泄漏 GTK 返回的`gchar*`；`Calloc`现在使用指针接收者（且`NewCalloc`返回`*Calloc`），使池能够实际跟踪分配；`zoomOut`使用`zoomInFactor`的倒数，而不是会被钳制为1.0的负乘数；`execJS`复用预分配的空 world-name，不再每次调用都泄漏一个`C.CString("")`；已从`menuItem.setAccelerator`中移除开发期间使用的`fmt.Println`。解决 #5465。
- 修复默认 GTK4 构建路径（`linux_cgo.go`）中相同的`Calloc`值接收器泄漏问题：改用指针接收器和`NewCalloc() *Calloc`，从而真正跟踪并释放每个窗口的`c.String(...)`分配。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，以便通过`PACKAGE_MANAGER`选项控制所使用的前端包管理器
- 使用`{{.Opn}}`和`{{.Cls}}`扩充模板数据，使 Taskfile 模板的编写更加可预测

## 变更

- 修改了几个现有 Taskfile，以使用`{{.Opn}} and {{.Cls}}`

## 修复

- 修复面板读取托盘菜单期间更新该菜单时，`linuxSystemTray`中的`concurrent map read and map write`运行时致命错误。
- 使用`log`而非`fmt`输出 WebView2 错误和堆栈跟踪，避免应用在 Windows 上运行但未连接控制台时丢失消息。

## v3.0.0-alpha.91 - 2026-05-12

## 变更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5414)中更新赞助者 SVG
- <strong>破坏性变更（macOS）：</strong>统一 macOS 坐标系，使`GetScreens`、`Position`和`SetPosition`均使用同一空间——逻辑点、Y 轴向下，且`(0,0)`位于主屏幕左上角。这与 Windows、GTK 以及 Electron 和 Web 的公共 API 一致。物理位置在主屏幕上方的屏幕现在会报告负的`Bounds.Y`（此前为正值），`Position()`/`SetPosition()`值现在使用逻辑点，而非`points × primaryScale`。仍可保持`Position()` → `SetPosition()`往返转换不变；需要更新早期 alpha 构建中记录的绝对值或手动计算的变通方案（例如乘以`primaryScale`，或根据屏幕高度翻转 Y 轴）。解决[#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5416)中以防御性方式验证 DBus 信号名称和消息体长度，防止发生 panic
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5363)中修复 Linux 上 GTK 菜单处理中的内存安全问题
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中检测 NVIDIA GPU，并在 Linux 上禁用 DMA-BUF 渲染器
- 修复 macOS 上`SetPosition`的跨屏幕 Y 坐标转换：使用主屏幕高度作为全局基准，使窗口在相对于主显示器有垂直偏移的显示器上落在正确位置，见[#5117](https://github.com/wailsapp/wails/issues/5117)
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修复 Git PR 模板，使其指向正确的反馈 URL
- 修复一系列由有缺陷的`DestroyMenu`系统调用引起的 Windows 系统托盘`SetMenu`崩溃：该调用错误地传入了四个参数而非一个，导致每次调用都返回 FALSE，且未释放任何内容。菜单重建时还会释放 HMENU 和 HBITMAP 句柄（包括运行时通过`MenuItem.SetBitmap`分配的句柄），重置`Win32Menu.Update`中过时的复选框/单选框映射，并移除`systemtray.updateMenu`中会使分配量翻倍的冗余`Update()`调用。长时间运行的系统托盘应用不再在每次菜单重建时泄漏 GDI/USER 对象。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 由 @vinhvoit225 在[PR](https://github.com/wailsapp/wails/pull/5261)中为 macOS 上的 WKWebView User-Agent 添加可配置的应用名称
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5400)中为 gin-service 示例添加间接依赖 github.com/coder/websocket
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5402)中为构建资产测试添加深度相等比较支持

## 变更

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5401)中将构建输出整合到 assets 目录
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5399)中更新赞助者 SVG

## 修复

- 由 @overlordtm 在[PR](https://github.com/wailsapp/wails/pull/5289)中为 macOS 单实例消息使用通知对象
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5383)中批量处理 Windows 回调，防止高负载下丢失 promise

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5316)中添加 go<em>test</em>results 作业，以汇总 Go 测试结果

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5369)中根据条件将大型 RPC 载荷拆分为分块 POST 请求
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5386)中将所有前端模板的 Vite 从5.x.x 升级到8.0.0
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5365)中将 Vite 开发服务器的端口配置迁移到环境变量
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5361)中配置所有模板，使 Vite 开发服务器绑定到127.0.0.1
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5384)中更新赞助者 SVG

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5312)中更新构建资产时清理 Info.plist 模板存根
- 通过在主线程上同步应用菜单项修改器（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修复 macOS 菜单中的过时状态，消除导致快速重新打开菜单时仍呈现先前状态的`dispatch_async`竞态（#5002）
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5203)中使开发模式忽略`*_test.go`文件，防止不必要的重新构建
- 由 @wucm667 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止应用未运行时 Menu.Update() 发生段错误
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5382)中使用 lastSizeWParam 控制 Windows 上菜单栏的重绘

## v3.0.0-alpha.88 - 2026-05-09

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5371)中将 HiddenOnTaskbar 改为使用 WS<em>EX</em>TOOLWINDOW
- 由 @atterpac 在[PR](https://github.com/wailsapp/wails/pull/5370)中重新排列依赖项，并移除 go.mod 中的 webview2 replace 指令
- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5358)中更新赞助者 SVG

## 修复

- 由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5331)中移除通用间接别名，并整合映射键类型

## 移除

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5377)中删除 PR-master 工作流，同时移除文档、Go 测试和跳过测试

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5352)中添加 Wails v3的韩语文档
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5354)中添加安装和快速入门的法语文档
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5355)中添加快速入门、概念和社区的葡萄牙语文档

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5328)中添加法语文档本地化
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5343)中为文档站点添加德语区域设置

## 变更

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5347)中将全部8翻译区域设置注册到文档配置中
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5317)中更新与 WebView2 相关的各个 Windows 文件

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5340)中为 Linux 分离 GTK3 和 GTK4 的对话框分派逻辑
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5339)中确保对话框回调在 GTK 线程上执行，修复段错误

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5179)中将 PR 模板 URL 添加到仓库
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5330)中添加 Wails v3的德语文档

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5307)中添加用于禁止 Escape 键在 macOS 上退出全屏模式的选项
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5310)中添加用于禁止 Escape 键在 macOS 上退出全屏模式的选项
- 由 @yuseferi 在[PR](https://github.com/wailsapp/wails/pull/5288)中添加 Pausa 社区展示文档

## 变更

- 由`@github-actions[bot]`在[PR](https://github.com/wailsapp/wails/pull/5308)中更新赞助商 SVG
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5309)中更新图标生成命令，使其能够处理不受支持的平台
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5224)中用三态 ButtonState 替换布尔型全屏 API，并实现平台绑定

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5315)中防止 WebView2 聚焦操作遇到 nil 控制器状态
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5313)中更新 GitHub Actions 工作流，使其正确引用 PR 的基础分支
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5203)中使开发模式忽略`*_test.go`文件，以避免不必要的重新构建
- 由 @wucm667 在[PR](https://github.com/wailsapp/wails/pull/5291)中防止应用未运行时 Menu.Update() 发生段错误

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5094)中添加 InstallScope 标志以及计算机级/用户级安装的构建选项
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5294)中为 BrowserWindow 添加无操作的 SetScreen 方法，以满足 Window 接口

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5295)中检测 NVIDIA GPU，并在 Linux 上禁用 DMA-BUF 渲染器
- 由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5109)中修复 Git PR 模板，使其指向正确的反馈 URL
- 修复一系列 Windows 系统托盘`SetMenu`崩溃问题。这些问题源于损坏的`DestroyMenu`系统调用传入了四个参数而非一个，导致每次调用都返回 FALSE，且未释放任何资源。还会在重建菜单时释放 HMENU 和 HBITMAP 句柄（包括运行时通过`MenuItem.SetBitmap`分配的句柄），重置`Win32Menu.Update`中过期的复选框/单选按钮映射，并移除`systemtray.updateMenu`中会使分配次数翻倍的冗余`Update()`调用。长时间运行的系统托盘应用不再在每次重建菜单时泄漏 GDI/USER 对象。

## v3.0.0-alpha.82 - 2026-05-01

## 修复

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5232)中修复桌面文件生成，使其正确处理桌面名称

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
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5272)中更新v3的文档和示例

## 修复

- 由 @AkagiYui 在[PR](https://github.com/wailsapp/wails/pull/5265)中为反向代理增加重试逻辑，并在开发环境中强制使用 IPv4
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5281)中重写未发布变更日志的触发工作流

## 移除

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5267)中移除用于各种测试目的的 shell 测试脚本
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5266)中删除v3-alpha 文档部署工作流和 CNAME 记录

### 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5196)中向侧边栏导航添加“前端路由”条目
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5185)中添加包含各框架专属建议的前端路由指南
- 添加对模态表单（macOS）的支持
- 由 @leaanthony 升级 ghw 版本，以更好地支持 Apple 设备 (#4977)
- 向 Dock 服务添加`GetBadge`方法
- 向`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）(#4957)
- 添加绑定生成器自动生成枚举的文档，包括专门的 Enums 页面和侧边栏导航 (#4972)
- 向`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）(#4957)
- 在`v3/examples/web-apis/`中添加 Web API 示例，演示41种浏览器 API，包括存储（localStorage、sessionStorage、IndexedDB、Cache API）、网络（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒体（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、设备（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、性能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 添加 WebView API 兼容性检查器示例（`v3/examples/webview-api-check/`），用于跨平台测试200种以上的浏览器 API
- 添加`internal/libpath`包，用于在 Linux 上查找原生库路径，支持并行搜索、缓存以及 Flatpak/Snap/Nix
- <strong>开发中：</strong>为 Linux 添加实验性 WebKitGTK 6.0 / GTK4 支持，可通过`-tags gtk4`启用（GTK3/WebKit2GTK 4.1仍为默认选项）
- 注意：在平铺式窗口管理器（例如 Hyprland、Sway）上，由于窗口管理器控制窗口几何布局，最小化/最大化操作可能无法按预期工作
- 由 @AbdelhadiSeddar 在**在 JavaScript 中监听事件**文档中添加了如何使用**一次性处理程序**的说明
- 由 @leaanthony 为`WebviewWindowOptions`添加`UseApplicationMenu`选项，使 Windows/Linux 上的窗口可继承通过`app.Menu.Set()`设置的应用程序菜单
- 由 @wimaha 添加对使用`.icon`文件（Apple Icon Composer 格式）生成 Liquid Glass 图标和资源目录（macOS）的支持（#4934）
- 添加用于无头/Web 部署的实验性服务器模式（`-tags server`）。该模式支持在没有原生 GUI 依赖项的情况下，将 Wails 应用作为 HTTP 服务器运行。使用`wails3 task build:server`构建。详情请参阅`examples/server`。
- 添加`internal/libpath`包，用于在 Linux 上查找原生库路径，支持并行搜索、缓存以及 Flatpak/Snap/Nix
- 由 @leaanthony 为`MacWindow`添加`CollectionBehavior`选项，用于控制窗口在 macOS Spaces 和全屏模式下的行为（#4756）
- 由 @leaanthony 为 pkg/application 添加单元测试
- 由 @leaanthony 为 MSIX 打包添加自定义协议支持
- 添加 Linux 桌面环境检测[PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 为 JavaScript 运行时添加`Window.Print()`方法，用于从前端触发打印对话框（#4290）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 由 @leaanthony 为 Linux 添加更多 WebKit2 加载状态变更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 在 Linux 构建期间生成`.desktop`文件，而非仅在打包时生成（#4575）
- 由 @leaanthony 添加 Linux 运行时依赖项文档，其中包含各发行版对应的软件包名称和 nfpm 打包示例（#4339）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加 NVIDIA 驱动程序版本信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中为 macOS 添加通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @chinenual 在[PR](https://github.com/wailsapp/wails/pull/4760)中为 helloworld 模板添加 aria-label 标识符，使 Appium 测试客户端能够轻松测试示例应用
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中为 macOS 添加通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @fbbdev 和 @ianvs 在[#4633](https://github.com/wailsapp/wails/pull/4633)中添加类型化事件
- 添加`systray-clock`示例，展示带有实时工具提示更新的无头托盘（#4653）。
- 由 @Tolfx 在 #4510 中为 Windows 添加 NSIS 协议模板
- 由 @Tolfx 在 #4510 中为 build-assets 添加测试
- macOS：由 @nidib 在[#4588](https://github.com/wailsapp/wails/pull/4588)中使菜单栏显示原生窗口控件
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于在程序坞中隐藏/显示应用图标
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于在程序坞中隐藏/显示应用图标
- 由 @leaanthony 在[#4534](https://github.com/wailsapp/wails/pull/4534)中为 macOS 添加原生 Liquid Glass 效果支持，使用 NSGlassEffectView（macOS 15.0及更高版本）并以 NSVisualEffectView 作为回退方案，同时提供全面的材质自定义选项
- 由 @leaanthony 在[#4500](https://github.dev/wailsapp/wails/pull/4500)中添加浏览器 URL 清理功能。该功能基于 @APShenkin 在[#4484](https://github.com/wailsapp/wails/pull/4484)中的工作。
- 由[@leaanthony](https://github.com/leaanthony)基于[@Taiterbase](https://github.com/Taiterbase)的原始工作，在此[PR](https://github.com/wailsapp/wails/pull/4241)中为 Windows/Mac 添加内容保护
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4488)中添加支持，可通过`wails3 build`和`wails3 package`别名将 CLI 变量传递给 Task 命令（#4422）
- 由[@atterpac](https://github.com/atterpac)在[#4318](https://github.com/wailsapp/wails/pull/4318)中添加对放置区的支持，其中事件源包含被放置元素的数据
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
- `app.Env.GetAccentColor`用于获取用户系统的强调色。支持 macOS。由[@etesam913](https://github.com/etesam913)实现
- 由[@atterpac](https://github.com/atterpac)在[#4137](https://github.com/wailsapp/wails/pull/4137)中添加`window.ToggleFrameless()` API
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4345)中为 Linux 添加各发行版特定的构建依赖项
- 由 @atterpac 在[PR](https://github.com/wailsapp/wails/pull/4404)中添加绑定指南
- **整理测试基础设施**：由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中将 Docker 测试文件移至专用的`test/docker/`目录，并优化镜像、提高构建可靠性
- **改进资源管理模式**：由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中为示例添加正确的事件处理程序清理和上下文感知的 goroutine 管理
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中添加对 aarch64 AppImage 构建的支持
- 由[@leaanthony](https://github.com/leaanthony)向`wails doctor`添加诊断章节
- 由[@leaanthony](https://github.com/leaanthony)在调用服务方法时将窗口添加到上下文中
- 由[@leaanthony](https://github.com/leaanthony)添加`window-call`示例，演示如何确定调用服务的是哪个窗口
- 由[@leaanthony](https://github.com/leaanthony)编写新的菜单指南
- 由[@leaanthony](https://github.com/leaanthony)改进 panic 处理
- 由[@leaanthony](https://github.com/leaanthony)编写新的菜单指南
- 由[@fbbdev](https://github.com/fbbdev)在[#4024](https://github.com/wailsapp/wails/pull/4024)中为 Service API 添加文档注释
- 由[@leaanthony](https://github.com/leaanthony)在[#4024](https://github.com/wailsapp/wails/pull/4024)中添加函数`application.NewServiceWithOptions`，用于通过附加配置初始化服务
- 由[@FalcoG](https://github.com/FalcoG)和[@leaanthony](https://github.com/leaanthony)在[#4031](https://github.com/wailsapp/wails/pull/4031)中改进菜单控制
- 由[@leaanthony](https://github.com/leaanthony)添加更多文档
- 由[@leaanthony](https://github.com/leaanthony)支持在标准事件监听器中取消事件
- 由[@leaanthony](https://github.com/leaanthony)添加对系统托盘`Hide`、`Show`和`Destroy`的支持
- 由[@leaanthony](https://github.com/leaanthony)添加对系统托盘`SetTooltip`的支持。最初的构想来自[@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中，使绑定生成器有关不受支持类型的警告能够报告包路径
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中为绑定生成器添加对泛型别名的支持
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中为绑定生成器添加对`omitzero` JSON 标志的支持
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中添加`//wails:ignore`指令，以阻止为选定的服务方法生成绑定
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中为服务和模型添加`//wails:internal`指令，允许类型在 Go 中导出但不在 JS/TS 中导出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中为绑定生成器添加对别名类型常量的支持，以便使用弱类型枚举
- 由[@fbbdev](https://github.com/fbbdev)在[#4068](https://github.com/wailsapp/wails/pull/4068)中添加针对 Go 1.24功能的绑定生成器测试
- 在[#4065](https://github.com/wailsapp/wails/pull/4065)中向`OSInfo.Branding`添加对 macOS 15“Sequoia”的支持，以改进操作系统版本检测
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中添加`PostShutdown`钩子，用于在关闭流程完成后运行自定义代码
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中添加`FatalError`结构体，以支持在自定义错误处理程序中检测致命错误
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中统一服务启动和关闭顺序并为其编写文档
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中添加应用启动/关闭序列测试框架，以及服务启动/关闭测试
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中添加`RegisterService`方法，用于在应用创建后注册服务
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中向应用和服务选项添加`MarshalError`字段，用于自定义处理绑定调用中的错误
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中添加可取消的 Promise 包装器，以便通过 Promise 链传播取消请求
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中添加将绑定调用的取消操作关联到`AbortSignal`的功能
- 由[@leaanthony](https://github.com/leaanthony)为 WML 添加对`data-wml-*`属性的支持，同时继续支持常规的`wml-*`属性
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为所有服务添加`Configure`方法，用于后期配置和动态重新配置
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中使`fileserver`服务在未配置时发送503“服务不可用”响应
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中使`kvstore`服务在未配置时默认提供内存键值存储
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`kvstore`服务添加`Load`方法，用于在配置更改后从文件重新加载数据
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`kvstore`服务添加`Clear`方法，用于删除所有键
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中向`log`服务添加`Level`类型，以提供 JS 端的日志级别常量
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`log`服务添加`Log`方法，用于动态指定日志级别
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中使`sqlite`服务在未配置时默认提供内存数据库
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`sqlite`服务添加`Close`方法，用于手动关闭数据库
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`sqlite`服务的查询方法添加取消支持
- 由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中为`sqlite`服务添加带 JS 绑定的预处理语句支持
- 基于[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始工作，由[Lea Anthony](https://github.com/leaanthony)在[PR](https://github.com/wailsapp/wails/pull/3537)中添加 Gin 支持
- 由[@oSethoum](https://github.com/osethoum)在[#4134](https://github.com/wailsapp/wails/pull/4134)中修复自动保存和密码自动保存始终启用的问题
- 由[@leaanthony](https://github.com/leaanthony)为窗口添加`SetMenu()`，以便为窗口设置菜单
- 由[@popaprozac](https://github.com/popaprozac)在[#4098](https://github.com/wailsapp/wails/pull/4098)中添加通知支持
-  由[@wimaha](https://github.com/wimaha)在[#4177](https://github.com/wailsapp/wails/pull/4177)中添加 macOS 文件关联支持
- 由[@leaanthony](https://github.com/leaanthony)添加用于提升语义化版本号的`wails3 tool version`
- 由[@popaprozac](https://github.com/popaprozac)在[#](https://github.com/wailsapp/wails/pull/4234)中添加 macOS 和 Windows 徽标支持
- 由[@fbbdev](https://github.com/fbbdev)和[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中添加已注册/强类型事件支持
- 由[@fbbdev](https://github.com/fbbdev)和[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中添加为自定义事件注册钩子的功能
- 由[@Krzysztofz01](https://github.com/Krzysztofz01)和[@rcalixte](https://github.com/rcalixte)添加`app.OpenFileManager(path string, selectFile bool)`，用于在系统文件管理器中打开路径`path`，并可选择通过`selectFile`高亮显示
- 由[@leaanthony](https://github.com/leaanthony)为`wails3 init`命令添加新的`-git`标志
- 由[@leaanthony](https://github.com/leaanthony)添加新的`wails3 generate webview2bootstrapper`命令
- 由[@leaanthony](https://github.com/leaanthony)在运行时中添加`init()`方法，以便手动初始化运行时
- 由[@leaanthony](https://github.com/leaanthony)向窗口的 WindowOptions 添加`WindowDidMoveDebounceMS`选项
- 由[@leaanthony](https://github.com/leaanthony)添加单实例功能。基于 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)。
- 由[@leaanthony](https://github.com/leaanthony)添加`wails3 generate template`命令
- 由[@leaanthony](https://github.com/leaanthony)添加`wails3 releasenotes`命令
- 由[@leaanthony](https://github.com/leaanthony)添加`wails3 update cli`命令
- 由[@leaanthony](https://github.com/leaanthony)为`wails3 generate bindings`命令添加`-clean`选项
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中支持构建 aarch64 (arm64) AppImage Linux 应用
- 由 @ansxuman 在[#3958](https://github.com/wailsapp/wails/pull/3958)中添加赞助者超链接
- 支持构建 Linux deb、rpm 和 Arch Linux 软件包，由
- 添加 Darwin 通用构建和软件包支持，由
- 向网站添加事件文档，由
- 为非 SSR 开发配置的 sveltekit 和 sveltekit-ts 模板
- 使用新的`wails3 update build-assets`命令更新构建资源，由
- 用于测试 HTML 拖放 API 的示例，由
- 文件关联支持，由[leaanthony](https://github.com/leaanthony)在
- 新的`wails3 generate runtime`命令，由
- 用于指定窗口是否应居中的新`InitialPosition`选项，由
- 向`application`包添加`Path`和`Paths`方法，由
- 添加`GeneralAutofillEnabled`和`PasswordAutosaveEnabled` Windows 选项
- 添加获取调用服务方法的窗口的功能，由
- 为 WebView2 添加`EnabledFeatures`和`DisabledFeatures`选项，由
- ⊞ 用于增强高 DPI 显示器支持的新 DIP 系统，由
- ⊞ 窗口类名选项，由[windom](https://github.com/windom/)在
- 扩展服务以提供插件功能。由
- 🐧 WindowDidMove / WindowDidResize 事件，位于
- ⊞ WindowDidResize 事件，位于
-  添加 Event ApplicationShouldHandleReopen，以便处理程序坞
-  由 @tmclane 在实现中添加 getPrimaryScreen/getScreens，位于
-  添加在 macOS 全屏模式下显示工具栏的选项，由
- 🐧 添加 onKeyPress 逻辑，将 Linux 按键转换为快捷键
- 🐧 添加任务`run:linux`，由
- 由[@almas-x](https://github.com/almas-x)在以下位置导出`SetIcon`方法
- 由[@almas-x](https://github.com/almas-x)在以下位置改进`OnShutdown`
- 在`Window`接口中恢复`ToggleMaximise`方法，由
- 向`Environment()`添加更多信息。由 @leaanthony 在
- 在`Window`接口上公开`WebviewWindow.IsFocused`方法，由
- 支持在 WML 系统中使用多个以空格分隔的触发事件，由
- 从捆绑的 JS 运行时脚本添加 ESM 导出，由
- 添加绑定生成器标志，以使用捆绑的 JS 运行时脚本，而不是
- 由[@abichinger](https://github.com/abichinger)在 Linux 上实现`setIcon`
- 向 dev 命令添加`-port`标志并支持环境变量
- 添加绑定方法调用的测试，由
- ⊞ 为已创建的窗口添加`SetIgnoreMouseEvents`，由
-  添加设置窗口堆叠层级（顺序）的功能，由

### 修复

- 在`Physical*`字段中将 NSScreen 点值转换为设备像素，修复 Retina Mac 上减半的`Screen.Bounds`、`WorkArea`和`Size`；同时填充顶层`Screen.X`/`Y`，使多显示器接触检测和工作区放置正确。由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5168)中完成
- 修复 ScreenManager 中的数据竞争；该问题会在显示配置发生变化时（例如睡眠/唤醒期间热插拔外接显示器）导致 WebKit DisplayLink 死锁
- 当 Assets.car 存在时，直接将 CFBundleIconName 设置为 appicon，由 @symball 在[PR](https://github.com/wailsapp/wails/pull/5154)中完成
- 修复`wails3 doctor`在 Fedora、openSUSE、Arch 和 NixOS 上报告错误 WebKitGTK 软件包的问题——由于v3在编译时需要4.1 API，已移除4.0回退条目（#5071）
- 修复 openSUSE webkit2gtk doctor 软件包名称（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，后者是正确的 openSUSE 软件包名称）（#5071）
- 修复桌面开发模式下缺少`/wails/custom.js`时出现的`Unexpected token '<'`错误。为`/wails/custom.js`添加了显式404处理程序，并在`loadOptionalScript`中添加了不区分大小写的`Content-Type`验证，以防止将 HTML SPA 回退内容作为 JavaScript 注入。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修复 macOS 上系统托盘菜单的高亮状态——菜单打开时，图标现在会显示选中状态（#4910）
- 修复 macOS 上附加到系统托盘的窗口显示在其他窗口后方的问题——现在使用正确的弹出窗口层级（#4910）
- 修复文档中错误的`@wailsio/runtime`导入示例（#4989）
- 修复 Darwin 上无边框窗口无法最小化的问题（#4294）
- 通过从 go-task 的最新状态检查中排除`node_modules/`，修复`wails3 build`和`wails3 dev`期间长达20-30分钟的卡顿。此前，`sources: "**/*"`通配模式会导致 go-task 枚举`node_modules/`中的每个文件并计算校验和（超过50000-100000个文件，其中包含 MUI 等大型依赖项），在 Windows/NTFS 上尤其缓慢（#4939）
- 修复 C `Screen` typedef 与 X11 Xlib.h 冲突所导致的 GTK4 构建失败（#4957）
- 修复 macOS 上程序坞徽标方法的一致性问题
- 修复`InvisibleTitleBarHeight`被应用于所有 macOS 窗口而非仅应用于无边框或标题栏透明窗口的问题（#4960）
- 启用`InvisibleTitleBarHeight`时，通过跳过窗口边缘附近的拖动启动，修复从顶部角落调整窗口大小时窗口晃动或抖动的问题（#4960）
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 贡献
- 修复 Windows 上显示缩放比例不为100%时文件拖放无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上使用悬停效果时文件拖放无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过使用`gtk_window_present()`，修复 Linux/GTK4 上显示或隐藏窗口时有时会恢复为最小化状态的问题（#4957）
- 通过`XTranslateCoordinates`/`XMoveWindow`添加以 X11 为条件的支持，修复 Linux/GTK4 上获取或设置窗口位置始终返回0,0的问题（#4957）
- 添加基于信号的尺寸限制来替代已移除的`gtk_window_set_geometry_hints`，修复 Linux/GTK4 上未强制实施窗口最大尺寸的问题（#4957）
- 通过实现正确的 PhysicalBounds 计算以及借助`gdk_monitor_get_scale`提供的分数缩放支持，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 贡献
- 修复 Windows 上显示缩放比例不为100%时文件拖放无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上使用悬停效果时文件拖放无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过实现正确的 PhysicalBounds 计算以及借助`gdk_monitor_get_scale`提供的分数缩放支持，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 贡献
- 修复 macOS 上因 App.Window.Current() 未从主线程访问 AppKit API 而导致的“幽灵窗口”问题（#4947），由 @wimaha 贡献
- 通过实现 WKUIDelegate runOpenPanelWithParameters，修复 HTML `<input type="file">`在 macOS 上无法工作的问题（#4862）
- 修复在 macOS/Linux 上使用`@wailsio/runtime` npm 模块时原生文件拖放无法工作的问题（#4953），由 @leaanthony 贡献
- 修复跨包类型别名的绑定生成问题（#4578），由 @fbbdev 贡献
- 修复因违反 GTK 线程安全要求而导致 OpenFileDialog 在 Linux 上崩溃的问题（#3683），由 @ddmoney420 贡献
- 修复对隐藏或已销毁的窗口调用`Focus()`时发生 SIGSEGV 崩溃的问题（#4890），由 @ddmoney420 贡献
- 修复在 Linux 上设置空图标或位图时可能发生 panic 的问题（#4923），由 @ddmoney420 贡献
- 修复从服务绑定调用 ErrorDialog 时在 macOS 上崩溃的问题（#3631），由 @leaanthony 贡献
- 使菜单能够在 Windows 操作系统上的`v3\examples\dialogs`中显示，由 @ndianabasi 贡献
- 修复页面重新加载期间因竞态条件而导致 TypeError 的问题（#4872），由 @ddmoney420 贡献
- 通过移除`Collector.IsVoidAlias()`方法中的全局状态，修复绑定生成器测试输出错误的问题（#4941），由 @fbbdev 贡献
- 修复`<input type="file">`文件选择器在 macOS 上无法工作的问题（#4862），由 @leaanthony 贡献
- 修复`Position()`和`SetPosition()`在 macOS 上使用不一致的坐标系，导致保存或恢复状态时窗口位置偏移的问题（#4816），由 @leaanthony 贡献
- 修复已通过应用程序清单设置 DPI 感知时出现 SetProcessDpiAwarenessContext“访问被拒绝”错误的问题（#4803）
- 更新键盘快捷键文档页面，并更正`KeyBinding.Add`回调参数的类型，由 @ndianabasi 贡献
- 修复有关生成自定义绑定的文档，必须使用`-d String`而非`-o String`
- 修复执行`menu.Update()`时菜单未清除子项的问题
- 修复文档中过时的 Manager API 引用（更新31文件，改用`app.Window.New()`、`app.Event.Emit()`等新模式），由 @leaanthony 贡献
- 修复因 WebKit 覆盖信号处理程序，导致绑定到 JS 的 Go 方法发生 panic 时 Linux 崩溃的问题（#3965），由 @leaanthony 贡献
- 修复 SaveFileDialog.SetFilename() 在 Linux 上不起作用的问题（#4841），由 @samstanier 贡献
- 修复拖放示例中放置坐标显示为 undefined 的问题
- 修复 APP_NAME 包含空格时 macOS 应用程序包创建失败的问题（大括号展开问题）
- 修复在 Windows 上调用服务方法时发生索引越界 panic 的问题（回退 goccy/go-json）
- 修复 Windows 上显示缩放比例不为100%时文件拖放无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上使用悬停效果时文件拖放无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 更新所有操作系统的 Taskfile.yml 文件中的全部命令，以支持`APP_NAME`等变量包含空格的情况，由 @ndianabasi 贡献
- 修复在 Linux 上执行 'build:universal:lipo:go' 任务时的命令参数错误，由 @wux1an 贡献
- 修复在 Linux 上运行 'wails3 build GOOS=darwin GOARCH=arm64' 时出现的 Docker 错误“undefined symbol: **<em>ubsan</em>handle_xxxxxxx”，由 @wux1an 完成
- 整合自定义协议文档并添加 Universal Links 章节，由 @leaanthony 完成
- 通过添加保护机制来防止并发调用 TrackPopupMenuEx，修复在 Windows 上反复单击系统托盘图标时菜单崩溃的问题 (#4151)，由 @leaanthony 完成
- 防止在 app.Run() 之前调用 systray.Run() 时应用崩溃，由 @leaanthony 完成
- 修复启用 ApplicationShouldTerminateAfterLastWindowClosed 后，通过 Hide()/Show() 切换窗口可见性时 macOS 崩溃的问题 (#4389)，由 @leaanthony 完成
- 修复在 macOS 和 Windows 上反复打开上下文菜单时发生的内存泄漏 (#4012)，由 @leaanthony 完成
- 修复 macOS 上上下文菜单的原生资源未被复用、导致每次显示时都新建菜单的问题 (#4012)，由 @leaanthony 完成
- 修复应用以`Hidden: true`启动后，单击 macOS 程序坞图标无法显示隐藏窗口的问题 (#4583)，由 @leaanthony 完成
- 修复因 CGO 调用中窗口指针类型不正确而导致 macOS 打印对话框无法打开的问题 (#4290)，由 @leaanthony 完成
- 修复 appmenu-gtk-module 访问尚未实现的窗口而导致 Wayland 上窗口菜单崩溃的问题 (#4769)，由 @leaanthony 完成
- 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题，由 @leaanthony 完成
- 修复在 Windows 上初始化拖放时出现“内存不足”错误的问题 (#4701)，由 @overlordtm 完成
- 修复因 URI 转义不正确而导致 Linux 文件资源管理器打开错误目录的问题 (#4397)，由 @leaanthony 完成
- 通过自动检测`.relr.dyn` ELF 节并禁用剥离，修复现代 Linux 发行版（Arch、Fedora 39+、Ubuntu 24.04+）上的 AppImage 构建失败问题 (#4642)，由 @leaanthony 完成
- 修复`wails doctor`在 Fedora/DNF 系统上错误报告已安装 webkit 软件包的问题 (#4457)，由 @leaanthony 完成
- 修复默认`config.yml`会使用生产构建运行`wails3 dev`的问题，由 @mbaklor 完成
- 修复 iOS 服务桩因导入不存在的软件包而导致构建失败的问题，由 @leaanthony 完成
- 修复 debug/info 方法中的结构化日志记录导致“no formatting directives”错误的问题，由 @leaanthony 完成
- 移除合并移动平台代码时意外包含的临时调试打印语句，由 @leaanthony 完成
- 通过自动禁用 DMA-BUF 渲染器，修复使用 NVIDIA GPU 时 WebKitGTK 在 Wayland 上崩溃的问题（错误 71 Protocol error），由 @leaanthony 完成
- 解决 Linux 上`application.WebviewWindowOptions.BackgroundColour`中的 alpha 值被忽略的问题（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修复未提供自定义图标时，Windows 系统托盘图标未默认使用应用图标的问题 (#4704)
- 跟踪`HICON`的所有权，仅销毁用户创建的句柄，防止 Explorer 回收时崩溃 (#4653)。
- 销毁时释放 Windows 系统主题监听器和保留的托盘图标，以避免 goroutine 和设备上下文泄漏 (#4653)。
- 将托盘工具提示截断至127个 UTF-16单元，以免损坏代理项对和多字节字形 (#4653)。
- 修复 Windows 打包任务失败的问题 (#4667)
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修复 go-webview2 v1.0.22签名变更导致的 Windows 构建错误 (#4513, #4645)
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 通过移除`<.Info.Protocol>`到`<.Protocol>`，修复 Linux desktop.tmpl 的协议范围，由 @Tolfx 在 #4510 中完成
- 修复 liquid glass 演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中完成
- 修复 Linux 上的系统托盘菜单更新问题，由[@JackDoan](https://github.com/JackDoan)在[#4604](https://github.com/wailsapp/wails/issues/4604)中完成
- 修复在 Windows 上创建隐藏窗口时出现白色窗口的问题，由 @leaanthony 在[#4612](https://github.com/wailsapp/wails/pull/4612)中完成
- 修复文档中的 notifications 软件包导入路径，由 @rxliuli 在[#4617](https://github.com/wailsapp/wails/pull/4617)中完成
- 修复使用 npm 软件包 @wailsio/runtime 时拖放功能无法工作的问题 (#4489)，由 @leaanthony 在 #4616 中完成
- Windows：修复启动时窗口闪烁以及隐藏窗口被错误显示的问题，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4600)中完成。
- 修复 Wayland 窗口最大化时的尺寸问题（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)完成
- 修复 Wayland 窗口最大化时的尺寸问题（https://github.com/wailsapp/wails/issues/4429），由[@samstanier](https://github.com/samstanier)完成
- 修复 liquid glass 演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中完成
- 修复 AssetServer 在 macOS 上可能崩溃的问题，由 @jghiloni 在[#4576](https://github.com/wailsapp/wails/pull/4576)中完成
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修复
- 修复夜间版本发布流水线，由 @riadafridishibly 在[#4597](https://github.com/wailsapp/wails/pull/4597)中完成
- 修复 liquid glass 演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中完成
- 修复 AssetServer 在 macOS 上可能崩溃的问题，由 @jghiloni 在[#4576](https://github.com/wailsapp/wails/pull/4576)中完成
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在[#4585](https://github.com/wailsapp/wails/pull/4585)中修复
- 修复夜间版本发布流水线，由 @riadafridishibly 在[#4597](https://github.com/wailsapp/wails/pull/4597)中完成
- 修复 liquid glass 演示中的重复定义错误，由 @Etesam913 在[#4542](https://github.com/wailsapp/wails/pull/4542)中完成
- 修复 Windows 上的 SetBackgroundColour，由 @PPTGamer 在[PR](https://github.com/wailsapp/wails/pull/4492)中完成
- 更新文档以反映 Manager API 重构带来的变更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 修复 Linux taskfile 中 Linux .desktop 文件的 appicon 变量 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 更新文档以反映 Manager API 重构带来的变更，由 @yulesxoxo 在[PR #4476](https://github.com/wailsapp/wails/pull/4476)中完成
- 修复[#4456](https://github.com/wailsapp/wails/issues/4456)中报告的 Windows 空指针解引用错误，由 @leaanthony 在[#4460](https://github.com/wailsapp/wails/pull/4460)中完成
- 在 macOS WKWebView 中添加对`allowsBackForwardNavigationGestures`的支持，以启用双指轻扫导航手势 (#1857)
- 修复初始设置为禁用的菜单项无法触发 onClick 的问题，由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中完成。感谢 @IanVS 的初步调查。
- 修复构建失败时未清理 Vite 服务器的问题 (#4403)
- 修复在 Windows 上关闭或取消`SaveFileDialog`时发生 panic 的问题。由 @hkhere 在[PR](https://github.com/wailsapp/wails/pull/4284)中修复
- 修复了 Windows 上 HTML 层级的拖放问题，由[@mbaklor](https://github.com/mbaklor)在[#4259](https://github.com/wailsapp/wails/pull/4259)中完成
- 在 macOS WKWebView 中添加对`allowsBackForwardNavigationGestures`的支持，以启用双指轻扫导航手势（#1857）
- 修复了最初设为禁用的菜单项无法使用 onClick 的问题，由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中完成。感谢 @IanVS 的前期调查。
- 修复构建失败时未清理 Vite 服务器的问题（#4403）
- 修复了 Windows 上的通知解析问题，由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4450)中完成
- 修复了 doctor 命令对 Windows SDK 依赖项的检查，由[@kodumulo](https://github.com/kodumulo)在[#4390](https://github.com/wailsapp/wails/issues/4390)中完成
- 修复了 Mac 上 processURLRequest 中的空指针解引用问题，由[@etesam913](https://github.com/etesam913)在[#4366](https://github.com/wailsapp/wails/pull/4366)中完成
- 修复了 Linux 上导致带筛选条件的对话框无法使用的问题，由[@bh90210](https://github.com/bh90210)在[#4287](https://github.com/wailsapp/wails/pull/4287)中完成
- 修复了 Windows 和 Linux 上的编辑菜单问题，由[@leaanthony](https://github.com/leaanthony)在[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)中完成
- 将 macOS .plist 文件中的最低系统版本从10.13.0更新为10.15.0，由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中完成
- 修复了窗口 ID 跳号问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复调用 RegisterContextMenu 时菜单为 nil 的问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了绑定生成器输出中的依赖循环，由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中完成
- 修复了绑定生成器输出中的先使用后定义错误，由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中完成
- 将构建标志传递给绑定生成器，由[@fbbdev](https://github.com/fbbdev)在[#4023](https://github.com/wailsapp/wails/pull/4023)中完成
- 将 Windows Taskfile 中的路径改为使用正斜杠，以确保其可在非 Windows 平台上运行，由[@leaanthony](https://github.com/leaanthony)完成
- 现已修复 Mac + Mac JS 事件问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了 macOS 上的事件死锁问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了 Windows 上初始化 Window 时，在提供 HTML 但未提供 JS 的情况下出现的`Parameter incorrect`错误，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了资源服务器中用于探测内容类型的响应前缀大小，由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中完成
- 修复了资源服务器对根索引路径上非404响应的处理，由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中完成
- 修复了绑定生成器在测试泛型类型的属性时出现的未定义行为，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了底层类型与具名包装类型属性不同时，绑定生成器为模型生成的输出，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器对映射键类型生成的输出及预处理，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器为实现编组器接口的结构体生成的输出，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器对涉及泛型类型的类型循环的检测，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器输出中对未导出模型的无效引用，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 将注入的代码移至服务文件末尾，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了绑定生成器对文件关闭操作所返回错误的处理，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 不再对仅定义生命周期或 http 方法而未定义其他绑定方法的服务发出警告，由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中完成
- 修复了非 React 模板在使用浅色系统配色方案时无法显示 Hello World 页脚的问题，由[@marcus-crane](https://github.com/marcus-crane)在[#4056](https://github.com/wailsapp/wails/pull/4056)中完成
- 修复了 macOS 上隐藏菜单项的问题，由[@leaanthony](https://github.com/leaanthony)完成
- 修复了消息处理器中的错误处理和格式化，由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
-  修复退出应用程序时跳过服务关闭的问题，由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中完成
-  确保菜单更新在主线程上进行，由[@leaanthony](https://github.com/leaanthony)完成
- 拖动和调整大小机制现在更加稳健，也更符合平台的预期行为，由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 修复了[#4097](https://github.com/wailsapp/wails/issues/4097)中 Webpack/angular 丢弃运行时初始化代码的问题，由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中完成
- 修复了初始状态为隐藏的菜单项，由[@IanVS](https://github.com/IanVS)在[#4116](https://github.com/wailsapp/wails/pull/4116)中完成
- 修复了 assetFileServer 在无扩展名请求中，当`[request]`不存在但`[request].html`存在时不提供`.html`文件的问题
- 修复了图标生成路径，由[@robin-samuel](https://github.com/robin-samuel)在[#4125](https://github.com/wailsapp/wails/pull/4125)中完成
- 修复了未发出`fullscreen`、`unfullscreen`、`unminimise`和`unmaximise`事件的问题，由[@oSethoum](https://github.com/osethoum)在[#4130](https://github.com/wailsapp/wails/pull/4130)中完成
- 修复了配置中默认版本的前缀不正确所导致的 NSIS 错误，由[@robin-samuel](https://github.com/robin-samuel)在[#4126](https://github.com/wailsapp/wails/pull/4126)中完成
- 修复了 Windows 上 Dialogs 运行时函数返回经过转义的路径的问题，由[TheGB0077](https://github.com/TheGB0077)在[#4188](https://github.com/wailsapp/wails/pull/4188)中完成
- 修复了 HKCU 中的 Webview2 检测路径，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 macOS 上的输入问题，由[@leaanthony](https://github.com/leaanthony)完成。
- 修复了 Windows 图标生成任务的文件名，由[@yulesxoxo](https://github.com/yulesxoxo)在[#4219](https://github.com/wailsapp/wails/pull/4219)中完成。
- 修复了无边框窗口的透明度问题，由[@leaanthony](https://github.com/leaanthony)基于 @kron 的工作完成。
- 修复了窗口被禁用或最小化时的聚焦调用问题，由[@leaanthony](https://github.com/leaanthony)基于 @kron 的工作完成。
- 修复了任务栏重启后系统托盘不显示的问题，由[@leaanthony](https://github.com/leaanthony)基于 @kron 的工作完成。
- 修复了 fallbackResponseWriter 未在 [#4245](https://github.com/wailsapp/wails/pull/4245) 中实现 Flush() 的问题
- 修复了 fallbackResponseWriter 未实现 Flush() 的问题，由 [@superDingda] 在 [#4236](https://github.com/wailsapp/wails/issues/4236) 中完成
- 修复了 macOS 上存在待处理的异步 Go 绑定函数调用时关闭窗口会崩溃的问题，由 [@joshhardy](https://github.com/joshhardy) 在 [#4354](https://github.com/wailsapp/wails/pull/4354) 中完成
- 修复了 Windows 效率模式启动时的竞态条件，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了 Windows 图标句柄的清理问题，由 [@leaanthony](https://github.com/leaanthony) 完成。
- 修复了 Windows 上的 `OpenFileManager`，由 [@PPTGamer](https://github.com/PPTGamer) 在 [#4375](https://github.com/wailsapp/wails/pull/4375) 中完成。
- 修复了 Linux 的最小/最大宽度选项，由 @atterpac 在 [#3979](https://github.com/wailsapp/wails/pull/3979) 中完成
- 通过提升 npm 版本为 TypeScript 模板提供类型定义，由 @atterpac 在 [#3966](https://github.com/wailsapp/wails/pull/3966) 中完成
- 修复了 SvelteKit 模板中的 CSS 引用，由 @atterpac 在 [#3945](https://github.com/wailsapp/wails/pull/3945) 中完成
- 确保 window run() 中的关键回调在主线程上调用，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了对话框目录选择器示例，由 [@leaanthony](https://github.com/leaanthony) 完成
- 为缺少 index.html 的情况创建了新的中文错误页面，由 [@leaanthony](https://github.com/leaanthony) 完成
-  确保 `windowDidBecomeKey` 回调在主线程上运行，由 [@leaanthony](https://github.com/leaanthony) 完成
-  支持无边框窗口全屏显示，由 [@leaanthony](https://github.com/leaanthony) 完成
-  改进了窗口销毁逻辑，由 [@leaanthony](https://github.com/leaanthony) 完成
-  修复了窗口附加到系统托盘时的位置逻辑，由 [@leaanthony](https://github.com/leaanthony) 完成
-  支持无边框窗口全屏显示，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了事件处理，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了窗口关闭逻辑，由 [@leaanthony](https://github.com/leaanthony) 完成
- 通用 Taskfile 现在默认会为 TypeScript 模板生成 TypeScript 绑定，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了没有打开任何窗口或仅有系统托盘时，收到 WM_CLOSE 消息会关闭应用程序的问题，由 [@mmalcek](https://github.com/mmalcek) 在 [#3990](https://github.com/wailsapp/wails/pull/3990) 中完成
- 修复了 garble 构建，由 @5aaee9 在 [#3192](https://github.com/wailsapp/wails/pull/3192) 中完成
- 修复了 Windows NSIS 构建，由 [@leaanthony](https://github.com/leaanthony) 完成
- 修复了 Linux 多选对话框中因未关闭而导致的死锁
- 修复了 Windows 构建期间对 .syso 文件的跨平台清理，由
- 修复了 amd64 AppImage 编译，由 @atterpac 在
- 修复了构建资源更新，由 @ansxuman 在
- 修复了 Linux 系统托盘的 `OnClick` 和 `OnRightClick` 实现，由 @atterpac 完成
- 修复了 `AlwaysOnTop` 在 Mac 上无法工作的问题，由
-  修复了 `application.NewEditMenu` 包含重复项的问题
- 🐧 修复了 aarch64 编译
- ⊞ 修复了单选组菜单项，由
- 修复了在 MacOS 上构建可运行的 .app 时，'name' 和 'outputfilename' 导致的错误
- 修复了拖放示例中使用 customEventProcessor 时的错误，由
- 🐧 修复了添加 IgnoreMouseEvents 后引入的 Linux 编译错误，由
- ⊞ 修复了 syso 图标文件生成错误，由
- 🐧 合并了在 Wayland 上原生运行的修复，来源于
- 不要绑定内部服务方法，位于
- ⊞ 修复了系统托盘启动时的 panic，位于
- 不要绑定内部服务方法，位于
- ⊞ 修复了系统托盘启动时的 panic，位于
- 大幅重构了菜单项和事件处理。目前主要改进了 macOS。由
- 修复了插件和事件重构后的测试，位于
- ⊞ 修复了 `Failed to unregister class Chrome_WidgetWin_0` 警告。由
- 模块问题
- 修复了调整大小事件的消息传递，由 [atterpac](https://github.com/atterpac) 在
- 🐧修复了 NixOS 上的主题处理错误，由
- 修复了 Windows 上跨卷安装项目的问题，由
- 修复了 React 模板的 CSS，使其显示页脚，由
- 通过更新到最新版本的 refresh，修复了开发模式下的僵尸进程
- 修复了 AppImage 的 WebKit 文件获取问题，由 [Atterpac](https://github.com/atterpac) 完成
- 修复了 Doctor 对 apt 软件包的验证，由 [Atterpac](https://github.com/Atterpac) 在
- 修复了应用程序退出时冻结的问题（Darwin），由 @5aaee9 在
- 修复了 Windows 上示例的背景颜色，由
- 修复了默认上下文菜单，由 [mmghv](https://github.com/mmghv) 在
- 修复了 Darwin 上方向键的十六进制值，由
- 使 Windows 上的拖放功能恢复正常。添加者：
- 修复了用户没有正确驱动程序时 Doctor 在 Linux 上的错误
- 修复了启动时的 DPI 缩放问题（Windows）。修改者：[@almas-x](https://github.com/almas-x)，位于
- 修复了 `go.mod` 中的替换行，使其使用相对路径。修复包含以下内容的 Windows 路径：
- 修复了 macOS 系统托盘未附加窗口时的点击处理，由
- 修复了未知选项导致 Windows 构建失败的问题，由
- 修复了没有附加窗口时在 Windows 上左键单击系统托盘图标导致的崩溃
- 修复了第二次打开窗口时 baseURL 错误的问题，由 @5aaee9 在 PR
- 修复了 `WebviewWindow.Restore` 方法中 if 分支的顺序，由
- 在多次调用 `GetStartURL` 时正确计算 `startURL`，当
- 修复了 `Screen` 结构体的 JS 类型，使其与对应的 Go 类型一致，由
- 修复了 `WML.Reload` 方法，以确保正确清理已注册的事件
- 修复了自定义上下文菜单在 Linux 上立即关闭的问题，由
- 修复了绑定生成的模型文件的输出路径和扩展名
- 修复了绑定生成的 JS 代码中模型文件的导入路径
- 修复了部分 Linux 发行版上的拖放功能，由
- 修复了使用 `wails3 task dev` 时缺少 macOS 任务的问题，由
- 修复了注册事件导致向 nil map 赋值的问题，由
- 修复了绑定方法参数的反序列化，由
- 修复绑定方法返回多个值时的处理问题，作者：
- 修复 doctor 无法检测并非通过系统包管理器安装的 npm 的问题
- 修复缺少 MicrosoftEdgeWebview2Setup.exe 的问题。感谢
- 修复因窗口 ID 处理而导致 Linux 上随机崩溃的问题，作者：@leaanthony。基于
- 修复 systemTray.setIcon 在 Linux 上导致崩溃的问题，作者：
- 修复首次调用 `setFrameless` 函数时未应用窗口边框的问题，平台：

### 变更

- **破坏性变更**：生成的 JS/TS 绑定中的映射键现已标记为可选，以准确反映 Go 映射的语义。现在，在 Typescript 中访问映射值会返回 `T | undefined` 而非 `T`，因此需要进行空值检查或断言（#4943），作者：`@fbbdev`
- 根据 `@wailsio/runtime` 中的变更，将 `Event` 改为使用 `Events`，并相应修改 `Features/Events/Event System` 文档中的函数调用，作者：@AbdelhadiSeddar
- 将 `EnabledFeatures`、`DisabledFeatures` 和 `AdditionalBrowserArgs` 从各窗口选项移至应用程序级 `Options.Windows`（#4559），作者：@leaanthony
- 更新 `Drag N Drop` 示例的 README，并重点说明该示例演示了 `Internal Drag and Drop`，作者：@ndianabasi
- 将多处调试日志的级别从 Info 改为 Debug（作者：@mbaklor）
- <strong>破坏性变更：</strong>在窗口选项中将 `EnableDragAndDrop` 重命名为 `EnableFileDrop`
- <strong>破坏性变更：</strong>在事件上下文中将 `DropZoneDetails` 重命名为 `DropTargetDetails`
- <strong>破坏性变更：</strong>将 `WindowEventContext` 上的 `DropZoneDetails()` 方法重命名为 `DropTargetDetails()`
- <strong>破坏性变更：</strong>移除 `WindowDropZoneFilesDropped` 事件，改用 `WindowFilesDropped`
- <strong>破坏性变更：</strong>将 HTML 属性从 `data-wails-dropzone` 改为 `data-file-drop-target`
- <strong>破坏性变更：</strong>将 CSS 悬停类从 `wails-dropzone-hover` 改为 `file-drop-target-active`
- <strong>破坏性变更：</strong>从 Windows 中移除 `DragEffect`、`OnEnterEffect` 和 `OnOverEffect` 选项（它们原属于已移除的 IDropTarget）
- 所有运行时 JSON 处理（方法绑定、事件、WebView 请求、通知、kvstore）均改用 goccy/go-json，性能提升 21-63%，内存分配减少 40-60%
- 优化 BoundMethod 结构体布局并缓存 isVariadic 标志，以降低每次调用的开销
- 对具有 `<=8` 个参数的方法使用栈上分配的参数缓冲区，以避免堆分配
- 优化方法调用中的结果收集，避免为单个返回值分配切片
- 使用 sync.Map 作为 MIME 类型缓存，以提升并发性能
- 使用缓冲池读取 HTTP 传输请求正文
- 在内容类型嗅探器中延迟分配 CloseNotify 通道，以减少每次请求的内存分配
- 移除资源服务器的 CSS 调试日志
- 扩展 MIME 类型扩展名映射，涵盖 50+ 种常见 Web 格式（字体、音频、视频等）
- 更新窗口 `X/Y` 选项的文档，作者：@ruhuang2001
- 更新 `Frontend Runtime` 文档，添加更多用于生成前端绑定的选项，作者：@ndianabasi
- 更新 Wails v3 资源服务器的文档页面，作者：@ndianabasi
- **破坏性变更**：移除包级对话框函数（`application.InfoDialog()`、`application.QuestionDialog()` 等）。请改用 `app.Dialog` 管理器：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新对话框文档以匹配实际 API：使用带回调的 `app.Dialog.*`、`AddButton()`（而非 `SetButtons()`）、`SetDefaultButton(*Button)`（而非字符串）、`AddFilter()`（而非 `SetFilters()`）、`SetFilename()`（而非 `SetDefaultFilename()`），并使用 `app.Dialog.OpenFile().CanChooseDirectories(true)` 选择文件夹
- **破坏性变更**：现在默认创建生产构建。要创建开发构建，请在 Taskfile 中设置 `DEV=true`。生成一个新项目作为示例，作者：@leaanthony
- 发出包含零个或一个数据参数的自定义事件时，数据值将直接赋给 Data 字段，而不会封装在切片中，作者：[@fbbdev](https://github.com/fbbdev)，见 [#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows 托盘现在会通过切换 `NIS_HIDDEN` 来遵循 `SystemTray.Show()`/`Hide()`，使应用可以真正消失并重新显示（#4653）。
- 托盘注册会复用已解析的图标，仅设置一次 `NOTIFYICON_VERSION_4`，并启用 `NIF_SHOWTIP`，以便工具提示在资源管理器重启后恢复（#4653）。
- macOS：窗口居中时使用 `visibleFrame` 而非 `frame`，以排除菜单栏和程序坞区域
- macOS：窗口居中时使用 `visibleFrame` 而非 `frame`，以排除菜单栏和程序坞区域
- 使用 `-config` 参数运行 `wails3 update build-assets` 时，通过 `-product*` 参数设置的值为
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，作者：@leaanthony，见 [#4471](https://github.com/wailsapp/wails/pull/4471)
- 重构内部窗口处理，作者：@leaanthony，见 [#4471](https://github.com/wailsapp/wails/pull/4471)
- 移除 `application.WindowIDKey` 和 `application.WindowNameKey`（由 `application.WindowKey` 取代），作者：[@leaanthony](https://github.com/leaanthony)
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，固定长度数组类型的类字段现在会按预期长度初始化，而不再初始化为空数组，作者：[@fbbdev](https://github.com/fbbdev)，见 [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- `application.NewService` 不再接受 options 作为可选参数（请改用 `application.NewServiceWithOptions`），作者：[@leaanthony](https://github.com/leaanthony)，见 [#4024](https://github.com/wailsapp/wails/pull/4024)
- 移除 `nanoid` 依赖，作者：[@leaanthony](https://github.com/leaanthony)
- 更新窗口示例，以展示 mica/acrylic/tabbed 窗口样式，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，已移除 `internal.js/ts` 模型文件；现在可在 `models.js/ts` 中找到所有模型，作者：[@fbbdev](https://github.com/fbbdev)，见 [#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 绑定中，具名类型绝不会渲染为其他具名类型的别名；旧行为现在仅限于别名，作者：[@fbbdev](https://github.com/fbbdev)，见 [#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 绑定的类模式下，类型为类型参数的结构体字段会标记为可选，且绝不会自动初始化，作者：[@fbbdev](https://github.com/fbbdev)，见 [#4045](https://github.com/wailsapp/wails/pull/4045)
- 从模板中移除 ESLint，作者：[@IanVS](https://github.com/IanVS)，见 [#4059](https://github.com/wailsapp/wails/pull/4059)
- 将版权年份更新为 2025，作者：[@IanVS](https://github.com/IanVS)，见 [#4037](https://github.com/wailsapp/wails/pull/4037)
- 由 [@IanVS](https://github.com/IanVS) 在 [#4075](https://github.com/wailsapp/wails/pull/4075) 中添加 event.Sender 文档
- 由 [@leaanthony](https://github.com/leaanthony) 添加 Go 1.24 支持
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中更改：现在调用 `App.Run` 时会调用 `ServiceStartup` 钩子，而不是在 `application.New` 中调用
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中更改：现在从 `App.Run` 返回 `ServiceStartup` 错误，而不是终止进程
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中更改：现在从 JS 发起的绑定和对话框调用会以错误对象而不是字符串拒绝
- 由 [@leaanthony](https://github.com/leaanthony) 改进 Windows 上的系统托盘菜单定位
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中将 JS 运行时移植到 TypeScript
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中更改：运行时在导入后立即初始化，无需等待窗口加载
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中更改：运行时不再导出 init 方法，可以通过仅用于副作用的导入来初始化运行时
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中更改：绑定方法现在返回一个 `CancellablePromise`，取消时会以 `CancelError` 拒绝；调用的实际结果将被丢弃
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中统一将内置服务类型命名为 `Service`
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中统一将带选项的内置服务创建函数命名为 `NewWithConfig`
- 由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中将 `sqlite` 服务的 `Select` 方法重命名为 `Query`，以便与 Go API 保持一致
- 模板：由 [@IanVS](https://github.com/IanVS) 在 [#4133](https://github.com/wailsapp/wails/pull/4133) 中将运行时移至 "dependencies"，并整理 package.json 文件
- 由 [@popaprozac](https://github.com/popaprozac) 在 [#4171](https://github.com/wailsapp/wails/pull/4171) 中更改：在开发环境中创建应用程序包并对其进行临时签名，以启用某些 macOS API
- 由 [@leaanthony](https://github.com/leaanthony) 将构建资源移至平台专用目录
- 由 [@leaanthony](https://github.com/leaanthony) 将 Taskfile 移至平台专用目录并重命名
- 由 [@leaanthony](https://github.com/leaanthony) 大幅改善缺少 `index.html` 时的体验
- [Windows] 由 [@leaanthony](https://github.com/leaanthony) 改进最小化和还原的性能。基于 [562589540](https://github.com/562589540) 提交的原始 [PR](https://github.com/wailsapp/wails/pull/3955)
- 由 [@leaanthony](https://github.com/leaanthony) 移除 `ShouldClose` 选项（请改为注册 events.Common.WindowClosing 事件钩子）
- [Windows] 由 [@leaanthony](https://github.com/leaanthony) 减少打开窗口时的闪烁
- 由 [@leaanthony](https://github.com/leaanthony) 移除 `Window.Destroy`，因为它原本应为内部函数
- 由 [@leaanthony](https://github.com/leaanthony) 将 `WindowClose` 事件重命名为 `WindowClosing`
- 由 [@leaanthony](https://github.com/leaanthony) 更改：前端构建现在根据构建类型使用 vite 环境 "development" 或 "production"
- 由 [@leaanthony](https://github.com/leaanthony) 更新到 go-webview2 v1.19
- 确保使用 taskfile 的 fork，由 @leaanthony 完成
- 更新 Taskfile 的 fork，以修复使用以下方式安装时的版本问题
- 使用 Taskfile 的 fork，以修复使用以下方式安装时的版本问题
- `service.OnStartup` 现在会在出错时关闭应用程序并运行
- 重构系统托盘点击消息传递，以更好地契合用户交互，由
- 资源嵌入中包含 `all:frontend/dist`，以支持会生成以下内容的框架
- Taskfile 重构，由 [leaanthony](https://github.com/leaanthony) 在以下内容中完成
- 由以下人员升级到 `go-webview2` v1.0.16
- 修复 `Screen` 类型，使其包含 `ID` 而不是 `Id`，由
- 更新 `go.mod.tmpl` 的 Wails 版本以支持 `application.ServiceOptions`，由
- 修复服务名称的确定逻辑，由 [windom](https://github.com/windom/) 在以下内容中完成
- mkdocs serve 现在使用 docker，由 [leaanthony](https://github.com/leaanthony) 完成
- 将开发配置整合到 `config.yml` 中，由
- 系统托盘对话框现在默认使用应用程序图标（如有）（Windows），由
- 改进 macOS 上的 GPU 和内存报告，由
- 移除 `WebviewGpuIsDisabled` 和 `EnableFraudulentWebsiteWarnings`
- Events API 变更：`On`/`Emit` -> 用户事件，`OnApplicationEvent` ->
- 修复 Linux 上的 Events API，由 [TheGB0077](https://github.com/TheGB0077) 在以下内容中完成
- [CI] 改进 actions，并使 actions 也能在 fork 和以下环境中运行
- 将 `AbsolutePosition()` 重命名为 `Position()`，由
- 将 Linux WebKit 依赖从 webkitgtk2-4.0 更新为 webkit2gtk-4.1，以
- 捆绑的 JS 运行时脚本现在是 ESM 模块：导入它的 script 标签
- `@wailsio/runtime` 包不会在 `window.wails` 上发布其 API
- 窗口 API 模块 `@wailsio/runtime/src/window` 现在会公开其所属的
- JS 窗口 API 已更新，以匹配当前的 Go `WebviewWindow`
- 绑定生成器现在默认按 ID 调用。`-id` CLI 选项
- 新的绑定代码布局：输出文件以前按文件夹组织
- 结构体字段 `application.Options.Bind` 已重命名为
- 绑定服务的新语法：现在必须将服务实例包装在
- 在非终端或 CI 环境中禁用旋转指示器，由

### 移除

- **破坏性变更**：从每个窗口的 `WindowsWindow` 选项中移除 `EnabledFeatures`、`DisabledFeatures` 和 `AdditionalLaunchArgs`。请改用应用程序级的 `Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures` 和 `Options.Windows.AdditionalBrowserArgs`。这些标志全局应用于共享的 WebView2 环境（#4559），由 @leaanthony 完成
- 移除 Windows 上原生的 `IDropTarget` 实现，改用基于 JavaScript 的方案（与 v2 的行为一致）
- 移除 github.com/wailsapp/mimetype 依赖，改用扩展后的扩展名映射和标准库 http.DetectContentType，将二进制文件大小减少约 1.2MB
- 通过为 Linux 文件资源管理器实现精简的 .desktop 文件解析器，移除 gopkg.in/ini.v1 依赖，节省约 45KB
- 通过使用 Go 1.21+ 标准库的 slices 包和精简的内部辅助函数，从运行时代码中移除 samber/lo，节省约 310KB
- 从 Darwin URL scheme 处理程序中移除调试 printf 语句（#4834）
- **破坏性变更**：移除 `linux:WindowLoadChanged` 事件；请改用 `linux:WindowLoadFinished` 检测 WebView 何时加载完毕（#3896），由 @leaanthony 完成

### 破坏性变更

- **Manager API 重构**：将应用程序 API 从扁平结构重组为分门别类的 manager，以改善代码组织并使 API 更易于发现；由 [@leaanthony](https://github.com/leaanthony) 在 [#4359](https://github.com/wailsapp/wails/pull/4359) 中完成
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
- 重命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`，由 [@leaanthony](https://github.com/leaanthony) 完成
- 将 `Path` 和 `Paths` 方法移至 `application` 包，由 [@leaanthony](https://github.com/leaanthony) 完成
- 应用程序菜单现在仅适用于 macOS，由 [@leaanthony](https://github.com/leaanthony) 完成

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

## 已移除

## v3.0.0-alpha.54 - 2025-12-29

## 新增

## 已修复

## 已移除

## v3.0.0-alpha.53 - 2025-12-27

## 新增

## 已修复

## v3.0.0-alpha.52 - 2025-12-26

## 已修复

## v3.0.0-alpha.51 - 2025-12-23

## 已修复

## v3.0.0-alpha.50 - 2025-12-21

## 已更改

## v3.0.0-alpha.49 - 2025-12-18

## 已更改

## v3.0.0-alpha.48 - 2025-12-16

## 新增

## 已更改

## 已修复

## v3.0.0-alpha.47 - 2025-12-15

## 新增

## 已修复

## v3.0.0-alpha.46 - 2025-12-14

## 新增

## 已移除

## v3.0.0-alpha.45 - 2025-12-13

## 新增

## 已修复

## v3.0.0-alpha.44 - 2025-12-12

## 新增

## 已更改

## 已修复

## v3.0.0-alpha.43 - 2025-12-11

## 新增

## v3.0.0-alpha.42 - 2025-12-10

## 新增

## v3.0.0-alpha.41 - 2025-11-23

## 已修复

## v3.0.0-alpha.40 - 2025-11-13

## 已修复

## v3.0.0-alpha.39 - 2025-11-12

## 新增

## 已更改

## v3.0.0-alpha.38 - 2025-11-04

## 新增

## 已更改

## 已修复

## v3.0.0-alpha.37 - 2025-11-02

## 已修复

## v3.0.0-alpha.36 - 2025-10-15

## 已修复

## v3.0.0-alpha.35 - 2025-10-14

## 已修复

## v3.0.0-alpha.34 - 2025-10-06

## 新增

## 已修复

## v3.0.0-alpha.33 - 2025-10-04

## 已修复

## v3.0.0-alpha.32 - 2025-10-02

## 已修复

## v3.0.0-alpha.31 - 2025-09-27

## 已修复

## v3.0.0-alpha.30 - 2025-09-26

## 已修复

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 已更改

## 已修复

## v3.0.0-alpha.29 - 2025-09-25

## 新增

## 已更改

## 已修复

## v3.0.0-alpha.27 - 2025-09-07

## 已修复

## v3.0.0-alpha.26 - 2025-08-24

## 新增

## v3.0.0-alpha.25 - 2025-08-16

## 已更改

不再被忽略，并会覆盖配置值。

## v3.0.0-alpha.24 - 2025-08-13

## 新增

## v3.0.0-alpha.23 - 2025-08-11

## 已修复

## v3.0.0-alpha.22 - 2025-08-10

## 新增

## 已更改

+ 修复范围过宽的 Linux 软件包依赖项，并修复过时的 RPM 依赖项。

## v3.0.0-alpha.21 - 2025-08-07

## 已修复

## v3.0.0-alpha.20 - 2025-08-06

## 已修复

## v3.0.0-alpha.19 - 2025-08-05

## 新增

## 已修复

## v3.0.0-alpha.18 - 2025-08-03

## 新增

## 已修复

## v3.0.0-alpha.17 - 2025-07-31

## 已修复

## v3.0.0-alpha.16 - 2025-07-25

## 新增

## v3.0.0-alpha.15 - 2025-07-25

## 新增

## v3.0.0-alpha.14 - 2025-07-25

## 新增

## v3.0.0-alpha.12 - 2025-07-15

### 新增

### 已修复

## v3.0.0-alpha.11 - 2025-07-12

## 新增

## v3.0.0-alpha.10 - 2025-07-06

### 破坏性变更

### 新增

### 已修复

### 已更改

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

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 变更

`service.OnShutdown`适用于此前已启动的所有服务，由 @atterpac 在[#3920](https://github.com/wailsapp/wails/pull/3920)中完成；@atterpac 在[#3907](https://github.com/wailsapp/wails/pull/3907)中；子文件夹，由 @atterpac 在[#3887](https://github.com/wailsapp/wails/pull/3887)中完成；[#3748](https://github.com/wailsapp/wails/pull/3748)；[leaanthony](https://github.com/leaanthony)；[etesam913](https://github.com/etesam913)在[#3778](https://github.com/wailsapp/wails/pull/3778)中；[northes](https://github.com/northes)在[#3836](https://github.com/wailsapp/wails/pull/3836)中；[#3827](https://github.com/wailsapp/wails/pull/3827)；[leaanthony](https://github.com/leaanthony)；[@leaanthony](https://github.com/leaanthony)；[@leaanthony](https://github.com/leaanthony)；（已由`EnabledFeatures`和`DisabledFeatures`选项取代），由[leaanthony](https://github.com/leaanthony)完成

### 修复

channel 变量，由 @michael-freling 在[#3925](https://github.com/wailsapp/wails/pull/3925)中完成；[ansxuman](https://github.com/ansxuman)在[#3924](https://github.com/wailsapp/wails/pull/3924)中；[#3898](https://github.com/wailsapp/wails/pull/3898)；[#3901](https://github.com/wailsapp/wails/pull/3901)；在[#3886](https://github.com/wailsapp/wails/pull/3886)中；[leaanthony](https://github.com/leaanthony)在[#3841](https://github.com/wailsapp/wails/pull/3841)中；Darwin 编辑菜单中的`PasteAndMatchStyle`角色，由[johnmccabe](https://github.com/johnmccabe)在[#3839](https://github.com/wailsapp/wails/pull/3839)中完成；[#3840](https://github.com/wailsapp/wails/issues/3840)在[#3854](https://github.com/wailsapp/wails/pull/3854)中，由[kodflow](https://github.com/kodflow)完成；[@leaanthony](https://github.com/leaanthony)；二者不同。由 @nickisworking 在[#3789](https://github.com/wailsapp/wails/pull/3789)中完成

## v3.0.0-alpha.7 - 2024-09-18

### 新增

[mmghv](https://github.com/mmghv)在[#3665](https://github.com/wailsapp/wails/pull/3665)中；[#3682](https://github.com/wailsapp/wails/pull/3682)；[atterpac](https://github.com/atterpac)和[leaanthony](https://github.com/leaanthony)在[#3570](https://github.com/wailsapp/wails/pull/3570)中

### 变更

应用程序事件`OnWindowEvent` -> 窗口事件，由[leaanthony](https://github.com/leaanthony)完成；[#3734](https://github.com/wailsapp/wails/pull/3734)；以`v3/`或`v3-`为前缀的分支，由[stendler](https://github.com/stendler)在[#3747](https://github.com/wailsapp/wails/pull/3747)中完成

### 修复

[etesam913](https://github.com/etesam913)在[#3742](https://github.com/wailsapp/wails/pull/3742)中；[atterpac](https://github.com/atterpac)在[#3721](https://github.com/wailsapp/wails/pull/3721)中；[atterpac](https://github.com/atterpac)在[#3675](https://github.com/wailsapp/wails/pull/3675)中；[#1811](https://github.com/wailsapp/wails/pull/1811)在[#3614](https://github.com/wailsapp/wails/pull/3614)中，由[@stendler](https://github.com/stendler)完成；[#3720](https://github.com/wailsapp/wails/pull/3720)，由[leaanthony](https://github.com/leaanthony)完成；[#3693](https://github.com/wailsapp/wails/issues/3693)，由[@DeltaLaboratory](https://github.com/DeltaLaboratory)完成；[#3720](https://github.com/wailsapp/wails/pull/3720)，由[leaanthony](https://github.com/leaanthony)完成；[#3693](https://github.com/wailsapp/wails/issues/3693)，由[@DeltaLaboratory](https://github.com/DeltaLaboratory)完成；[leaanthony](https://github.com/leaanthony)；[#3746](https://github.com/wailsapp/wails/pull/3746)，由[@stendler](https://github.com/stendler)完成；[leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修复

## v3.0.0-alpha.5 - 2024-07-30

### 新增

[#3580](https://github.com/wailsapp/wails/pull/3580)；[#3580](https://github.com/wailsapp/wails/pull/3580)；点击图标，由 @5aaee9 在[#2991](https://github.com/wailsapp/wails/pull/2991)中完成；[#2618](https://github.com/wailsapp/wails/pull/2618)；[@fbbdev](https://github.com/fbbdev)在[#3282](https://github.com/wailsapp/wails/pull/3282)中；@[Atterpac](https://github.com/Atterpac)在[#3022](https://github.com/wailsapp/wails/pull/3022])中；[@marcus-crane](https://github.com/marcus-crane)在[#3146](https://github.com/wailsapp/wails/pull/3146)中；[PR](https://github.com/wailsapp/wails/pull/3147)；[PR](https://github.com/wailsapp/wails/pull/3189)；[@fbbdev](https://github.com/fbbdev)在[#3281](https://github.com/wailsapp/wails/pull/3281)中；[aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)；基于 @Mai-Lapyst 的[PR](https://github.com/wailsapp/wails/pull/2044)；[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中；[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中；[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中；npm 软件包，由[@fbbdev](https://github.com/fbbdev)在[#3334](https://github.com/wailsapp/wails/pull/3334)中完成；在[#3354](https://github.com/wailsapp/wails/pull/3354)中；`WAILS_VITE_PORT`，由[@abichinger](https://github.com/abichinger)在[#3429](https://github.com/wailsapp/wails/pull/3429)中完成；[@abichinger](https://github.com/abichinger)在[#3431](https://github.com/wailsapp/wails/pull/3431)中；[@bruxaodev](https://github.com/bruxaodev)在[#3667](https://github.com/wailsapp/wails/pull/3667)中；[@OlegGulevskyy](https://github.com/OlegGulevskyy)在[#3674](https://github.com/wailsapp/wails/pull/3674)中

### 修复

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) 于   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) 于   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) 于   [#3477](https://github.com/wailsapp/wails/pull/3477)   由[Atterpac](https://github.com/atterpac)于   [#3320](https://github.com/wailsapp/wails/pull/3320)。   于[#3306](https://github.com/wailsapp/wails/pull/3306)。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) 于   [#2750](https://github.com/wailsapp/wails/pull/2750)。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   [jaybeecave](https://github.com/jaybeecave) 于   [#3052](https://github.com/wailsapp/wails/pull/3052)。   [@pylotlight](https://github.com/pylotlight) 于   [PR](https://github.com/wailsapp/wails/pull/3039)   已安装。由[@pylotlight](https://github.com/pylotlight)在   [PR](https://github.com/wailsapp/wails/pull/3032)中添加   [PR](https://github.com/wailsapp/wails/pull/3145)   空格 - @leaanthony。   [thomas-senechal](https://github.com/thomas-senechal) 于 PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) 于 PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   附加窗口[tw1nk](https://github.com/tw1nk)于 PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) 于   [#3279](https://github.com/wailsapp/wails/pull/3279)   存在`FRONTEND_DEVSERVER_URL`。   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) 于   [#3295](https://github.com/wailsapp/wails/pull/3295)   监听器，由[@fbbdev](https://github.com/fbbdev)于   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) 于   [#3330](https://github.com/wailsapp/wails/pull/3330)   生成器，由[@fbbdev](https://github.com/fbbdev)于   [#3334](https://github.com/wailsapp/wails/pull/3334)   生成器，由[@fbbdev](https://github.com/fbbdev)于   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) 于   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) 于   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) 于   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) 于   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) 于   [#3431](https://github.com/wailsapp/wails/pull/3431)   由[@pekim](https://github.com/pekim)于   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)，作者为   [@5aaee9](https://github.com/5aaee9)。   [@windom](https://github.com/windom/) 于   [#3636](https://github.com/wailsapp/wails/pull/3636)。   Windows，由[@bruxaodev](https://github.com/bruxaodev/)于   [#3691](https://github.com/wailsapp/wails/pull/3691)。

### 变更

[mmghv](https://github.com/mmghv) 于   [#3611](https://github.com/wailsapp/wails/pull/3611)   支持 Ubuntu 24.04 LTS，由[atterpac](https://github.com/atterpac)于   [#3461](https://github.com/wailsapp/wails/pull/3461)   必须具有`type="module"`属性。由   [@fbbdev](https://github.com/fbbdev)于   [#3295](https://github.com/wailsapp/wails/pull/3295)   对象，并且不启动 WML 系统。此项变更旨在改进封装。如有需要，可以调用新的`WML.Enable`方法手动启动 WML 系统。捆绑的 JS 运行时脚本仍会自动执行这两项操作。由[@fbbdev](https://github.com/fbbdev)于   [#3295](https://github.com/wailsapp/wails/pull/3295)   窗口对象作为默认导出。现在无法再通过 ESM 命名导入或命名空间导入语法导入单个方法。   API。部分方法的名称或原型已更改，具体如下：`Screen`变为`GetScreen`；`GetZoomLevel`/`SetZoomLevel`变为`GetZoom`/`SetZoom`；   `GetZoom`、`Width`和`Height`现在直接返回值，不再将其包装在对象中。由[@fbbdev](https://github.com/fbbdev)于   [#3295](https://github.com/wailsapp/wails/pull/3295)   已被移除。使用`-names` CLI 选项可切换回按名称调用。   由[@fbbdev](https://github.com/fbbdev)于   [#3468](https://github.com/wailsapp/wails/pull/3468)   以前以其所属包命名；现在使用完整的 Go 导入路径，包括模块路径。由[@fbbdev](https://github.com/fbbdev)于   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`。由[@fbbdev](https://github.com/fbbdev)于   [#3468](https://github.com/wailsapp/wails/pull/3468)   调用`application.NewService`。由[@fbbdev](https://github.com/fbbdev)于   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) 于   [#3574](https://github.com/wailsapp/wails/pull/3574)
