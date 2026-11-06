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

/_-- 本文件将记录此项目的所有重要变更。

格式以[维护更新日志](https://keepachangelog.com/en/1.0.0/)为基础，并且本项目遵循[语义化版本](https://semver.org/spec/v2.0.0.html)规范。

- `Added`表示新增功能。
- `Changed`表示现有功能的变更。
- `Deprecated`表示即将移除的功能。
- `Removed`表示现已移除的功能。
- `Fixed`表示错误修复。
- `Security`表示安全漏洞修复。

_/

/_   * 自动更新 *   发布说明根据已合并的 PR 生成。请勿在常规贡献 PR 中编辑此文件，也不要重新创建已移除的`v3/UNRELEASED_CHANGELOG.md`。有关更新日志维护的信息，请参阅`AGENTS.md`。 _/

## [未发布]

## v3.2.0-nightly.20261103 - 2026-11-03

## 变更
- [v3] 在 Menu.Update 时刷新系统托盘菜单，见 @johnmaguire 的[PR](https://github.com/wailsapp/wails/pull/6314)
- [v3/Windows] 为窗口实现 SetIcon，见 @hypherionmc 的[PR](https://github.com/wailsapp/wails/pull/6502)

## v3.2.0-nightly.20261102 - 2026-11-02

## 变更
- 使用 DisableOverscroll 取代 DisableBounce，见 @lethusang 的[PR](https://github.com/wailsapp/wails/pull/6077)

## v3.2.0-nightly.20261101 - 2026-11-01

## 新增
- 添加实验性的无 CGO macOS 和 Linux 后端，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 添加原生 YAML 构建、可移植打包和签名，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 修复
- 修复多项v3夜间构建问题并清理依赖项，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6519)
- 修复v3发布验证和固件分类，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6533)
- 解决跨平台发布验证失败问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6538)

## v3.1.0 - 2026-11-01

## 新增
- 添加实验性的无 CGO macOS 和 Linux 后端，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6521)
- 添加原生 YAML 构建、可移植打包和签名，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6531)

## 变更
- 避免复制已丢弃的桌面元数据，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6391)
- 不使用 fmt 格式化内容长度，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6393)
- 选择索引根目录时不收集候选项，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6404)
- 复用不可定位文件的复制缓冲区，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6408)
- 避免为探测前缀分配读取器，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6410)
- 对已知扩展名跳过 MIME 探测读取，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6414)
- 防止队列增长与收缩振荡，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6415)
- 避免分配随后会被丢弃的请求上下文，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6419)
- 跳过已禁用调试日志的准备工作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6420)
- 延迟分配 MIME 探测缓冲区，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6428)
- 减少超链接格式化产生的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6443)
- 消除项目符号格式化产生的临时内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6449)
- 将 UUID 字符串转换减少到一次内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6453)
- 减少语句句柄格式化产生的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6454)
- 消除基于名称的 UUID 摘要所产生的临时内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6455)
- 避免复制尚未完成的命令输出行，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6460)
- 在包暂存期间以流式方式复制可执行文件，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6472)
- 避免重新扫描尚未完成的输出行，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6473)
- 复用已编译的命令路由模式，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6475)
- 避免在运行时发现期间查询元数据，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6482)
- 减少本地化后备页面产生的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6485)
- 避免为每个处理器重新编译模式，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6490)

## 修复
- 为系统托盘实现 SetTooltip，见 @fan711 的[PR](https://github.com/wailsapp/wails/pull/6017)
- 从设置门户读取深色模式，见 @TechHutTV 的[PR](https://github.com/wailsapp/wails/pull/6072)
- 允许 Window.SetZoom 低于1.0，见 @APshenkin 的[PR](https://github.com/wailsapp/wails/pull/6166)
- 按照触发顺序向每个 Go 监听器传递事件，见 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6275)
- 在 Emit 中持有 hookLock 时读取事件钩子，见 @DevLumuz 的[PR](https://github.com/wailsapp/wails/pull/6276)
- 始终公布托盘的 Menu 属性，以免在 Run 之后创建的托盘被丢弃，见 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6284)
- 将在 Run 之前发送的主线程工作保留到 Run 启动后再执行，见 @AlbinoGeek 的[PR](https://github.com/wailsapp/wails/pull/6287)
- 在 GitHub 更新程序提供程序中使用正确的资源 URL，见 @gavidroselj 的[PR](https://github.com/wailsapp/wails/pull/6304)
- 当 WebView2 省略 Content-Length 时保留请求正文，见 @strraaa 的[PR](https://github.com/wailsapp/wails/pull/6326)
- 在 DIP 空间中报告 macOS 和 GTK 屏幕（混合 DPI 边界），见 @julianstorer 的[PR](https://github.com/wailsapp/wails/pull/6346)
- 不要将已关闭的通知报告为点击，见 @bercedev 的[PR](https://github.com/wailsapp/wails/pull/6367)
- 终端关闭时停止开发进程，见 @xmasdev 的[PR](https://github.com/wailsapp/wails/pull/6371)
- 在证据检查中将更新日志运行映射到压缩合并提交，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6373)
- 在[PR](https://github.com/wailsapp/wails/pull/6386)中将任务栏重启的防抖延迟保持一秒，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6387)中保留从全屏状态开始时的尺寸约束，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6388)中不再为适应新约束而调整全屏或最大化窗口的大小，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6389)中不再于 liquid-glass 示例的 Run 之前弹出对话框，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6390)中让 Run 之前请求的对话框不执行任何操作，而不是永远等待，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6392)中为桌面条目没有命令的情况提供回退处理，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6394)中嗅探当前内容且不保留文件名，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6395)中识别大小写混合的资源扩展名，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6396)中释放已使用的内容嗅探器缓冲区，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6397)中释放已移除的 ringqueue 引用，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6398)中先遵循响应头删除操作，再选择回退方案，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6399)中于刷新前提交回退响应头，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6400)中仅发现名称完全匹配的索引文件，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6401)中避免为不可查找的 HEAD 响应生成正文和执行完整读取，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6402)中通过多次短读取填满 MIME 嗅探前缀，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6403)中提交成功但内容为空的不可查找响应，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6405)中避免重复展开后的 URI 目标，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6406)中进行字段展开时保留路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6407)中遵循 XDG 桌面文件搜索目录，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6409)中执行空刷新时提交嗅探器响应头，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6411)中跳过遮蔽桌面条目的目录，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6412)中保留在可执行文件路径内找到的目标，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6413)中于定时器触发后释放回调捕获项，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6416)中避免 ldconfig 库误匹配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6417)中限制文件系统库的回退匹配，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6418)中消除缓存失效的数据竞争，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6421)中遵循可选的日志记录器配置，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6422)中避免修改共享的启动 URL 及使用过期端口，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6423)中于请求取消时停止开发代理重试，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6424)中针对缺失的捆绑资源返回404，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6425)中处理捆绑运行时的 HEAD 请求，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6426)中路由服务时保留转义路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6427)中正确排序超大数值型预发布标识符，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6429)中保存工具下载内容前拒绝 HTTP 错误响应，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6430)中进行桌面发现时遵循启动截止时间，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6432)中在开发代理重试期间保留请求正文，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6434)中的每条返回路径上关闭 HTTP 响应正文，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6435)中保留并行库搜索中的取消错误，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6436)中于解码前拒绝失败的 GitHub 标签响应，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6437)中按字面保留带方括号的模板路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6438)中仅安装一次 wails3 并保留安装程序错误，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6439)中忽略注释块内的隐藏条目和标题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6440)中呈现稳定版本查找错误，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6441)中打印纯文本消息且不作格式解析，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6442)中避免加载 JSON null 后写入时发生 panic，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6444)中为整个存储区的读取返回快照，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6446)中对部分行完成迭代后传播错误，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6447)中按字面保留解压目标路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6448)中于关闭和重新打开时释放预处理语句句柄，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6450)中处理预处理操作与关闭或重新打开操作之间的竞争，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6451)中正确处理回退索引页面的 HEAD 请求，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6452)中要求完整的十六进制语句句柄，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6456)中在并发情况下保留默认内存数据库，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6457)中将 stderr 与已解析的 HEAD 输出分离，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6458)中从嵌入式捆绑包生成运行时资源，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6459)中检测 TODO、FIXME 和 TBD 占位符，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6461)中验证发行说明 API 响应，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6462)中进行状态检查时遵循任务目录，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6463)中防止复制到同一文件导致数据丢失，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6464)中将任务目录纳入缓存标识，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6465)中每次检查仅探测一次端口，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6466)中避免缓存条目解码为 null 时发生 panic，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6467)中递增数值型预发布版本号且不受整数范围限制，作者：@taliesin-ai
- 防止旧版检查掩盖缓存未命中问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6468)
- 在执行器的基础目录中运行默认任务，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6469)
- 保留资源管理器命令中的字面百分号，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6470)
- 统一前置条件与任务命令的目录，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6471)
- 将符号链接替换为生成的普通文件，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6474)
- 构建失败时保留现有 DMG，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6476)
- 保留带引号的正则表达式及路径中的反斜杠，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6477)
- 分配图像前验证生成的图标尺寸，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6478)
- 通过 shell 执行命令链和重定向，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6479)
- 构建失败时保留之前的 Linux 软件包，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6480)
- 使用特定于平台的 shell 回退执行方式，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6481)
- 支持提供路径开头含重复斜杠的资源，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6483)
- 识别带有构建元数据的稳定版本，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6484)
- 准确报告资源生成失败，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6486)
- 防止替换当前报告器时发生竞态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6487)
- 串行执行并发的缓存文件保存操作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6488)
- 支持多个模板相等条件候选项，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6489)
- 规范化保存值为 null 的类别映射，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6491)
- 显示变量 shell 命令失败时的诊断信息，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6492)
- 不将 shell 命令链放入隐式 Go 缓存，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6493)
- 串行执行 Linux 类别持久化操作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6494)
- 缓存路径含带引号或转义空格的 Go 输出，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6495)
- 编码失败时保留 Mac 图标，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6496)
- 不再丢弃无效的 Go 测试计数标志，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6497)
- 复制到新文件时保留源文件权限，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6498)
- 修复多项v3夜间构建问题并清理依赖项，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6519)
- 修复v3版本发布验证和测试夹具分类，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6533)
- 解决跨平台版本发布验证失败问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6538)

## v3.1.0-nightly.20261028 - 2026-10-28

## 变更
- 通过复用已编译的正则表达式模式提升变更日志生成速度，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6490)
- 减少提供本地化回退页面时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6485)
- 通过避免不必要的元数据查询加快 AppImage 运行时发现，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6482)
- 复用已编译的命令路由模式以提升 wake 性能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6475)
- 避免重复扫描尚未完成的输出行以提升性能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6473)
- 在 MSIX 软件包暂存期间以流式方式复制可执行文件，从而减少内存使用并加快构建，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6472)
- 避免复制尚未完成的命令输出行以提升性能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6460)
- 减少生成基于名称的 UUID 时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6455)
- 减少格式化 SQLite 语句句柄时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6454)
- 将 UUID 字符串转换所需的内存分配减少到一次，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6453)
- 减少格式化变更日志项目符号条目时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6449)
- 减少格式化终端超链接时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6443)
- 通过延迟探测 MIME 缓冲区，减少提供资源时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6428)
- 禁用调试日志记录时，跳过资源服务器中的调试日志准备工作，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6420)
- 避免创建随后丢弃的请求上下文，减少 webview 请求处理中的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6419)
- 对扩展名已知的文件跳过 MIME 探测，从而加快资源提供速度，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6414)
- 避免使用读取器探测内容，减少提供资源时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6410)
- 减少提供不可定位的资源文件时的内存分配，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6408)
- 避免收集候选项，加快资源服务器选择索引根目录的速度，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6404)
- 不使用 fmt 格式化内容长度，从而提升资源服务器性能，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6393)
- 减少读取桌面文件资源管理器元数据时的内存复制，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6391)

## 修复
- 将源文件复制到新的目标文件时保留其权限，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6498)
- 最小或最大尺寸约束发生变化时，避免全屏和最大化窗口调整大小，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6388)
- 发布 wailsapp/wails 仓库的定期 Star 历史记录更新，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6500)
- 不再让 wake 丢弃无效的 Go 测试计数标志，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6497)
- 编码失败时保留 Mac 图标，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6496)
- 修复路径含带引号或转义空格的 Go 输出发生构建缓存未命中的问题，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6495)
- 在 Linux 上串行执行通知类别持久化操作，以防止并发写入竞态，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6494)
- 不将 shell 命令链放入隐式 Go 构建缓存，见 @taliesin-ai 的[PR](https://github.com/wailsapp/wails/pull/6493)
- 修复保存的数据为 null 时通知类别映射加载失败的问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6491)中完成
- 显示变量 shell 命令执行失败时的诊断信息，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6492)中完成
- 修复存在多个相等候选项时 wake 模板匹配的问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6489)中完成
- 串行处理并发的缓存文件保存操作，防止缓存写入损坏，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6488)中完成
- 防止替换活动 reporter 时发生竞态，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6487)中完成
- 在 CLI 中准确报告资源生成失败，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6486)中完成
- 在 CLI 中识别包含构建元数据的稳定版本，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6484)中完成
- 提供请求路径包含重复前导斜杠的资源，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6483)中完成
- 修复 wake 在各平台上的 shell 回退执行，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6481)中完成
- 构建失败时保留已有的 Linux 软件包，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6480)中完成
- 通过 shell 执行命令链和重定向，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6479)中完成
- 解析 wake 命令时保留带引号的正则表达式和路径中的反斜杠，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6477)中完成
- 打包构建失败时保留现有的 DMG 文件，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6476)中完成
- 在绑定输出中用生成的常规文件替换符号链接，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6474)中完成
- 在 Linux 文件浏览器命令中保留字面百分号，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6470)中完成
- 对齐前置条件与任务命令的目录，使其解析结果保持一致，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6471)中完成
- 在执行器的基准目录中运行默认任务，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6469)中完成
- 防止旧版检查掩盖缓存未命中，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6468)中完成
- 修复 CLI 预发布版本编号，使其递增不受整数溢出限制，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6467)中完成
- 避免缓存条目解码为 null 时发生 panic，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6466)中完成
- 在 CLI 中每次检查仅探测端口一次，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6465)中完成
- 将任务目录纳入缓存标识，以避免过时的 wake 结果，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6464)中完成
- 防止打包时将文件复制到自身而导致数据丢失，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6463)中完成
- 使状态检查遵循任务目录，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6462)中完成
- 检测变更日志条目中的 TODO、FIXME 和 TBD 占位符，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6459)中完成
- 在 CLI 中从嵌入式捆绑包生成运行时资源，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6458)中完成
- 在 git 命令中将 stderr 与解析后的 HEAD 输出分离，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6457)中完成
- 在并发访问时保留默认的内存中 SQLite 数据库，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6456)中完成
- 为回退索引页面的 HEAD 请求返回正确响应，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6451)中完成
- 处理 SQLite 语句预处理与数据库关闭或重新打开之间的竞态，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6450)中完成
- 关闭或重新打开 SQLite 数据库时释放预处理语句句柄，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6448)中完成
- 使用 gosod 时保持解压目标路径的字面形式，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6447)中完成
- 在部分行迭代后传播 SQLite 查询错误，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6446)中完成
- 在 kvstore 中读取整个存储时返回快照，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6444)中完成
- 防止 kvstore 加载 JSON null 值后写入时发生 panic，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6442)中完成
- 输出纯文本终端消息时不解析格式化动词，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6441)中完成
- 在更新程序中查找最新稳定版本时呈现错误，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6440)中完成
- 忽略变更日志注释块中的隐藏条目和标题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6439)中完成
- 仅安装一次 wails3，并在更新程序中呈现安装程序错误，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6438)中完成
- 在 gosod 中保留带方括号的模板路径字面值，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6437)中完成
- 在更新程序中解码 GitHub 标签响应前拒绝失败的响应，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6436)中完成
- 在并行搜索库时保留取消错误，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6435)中完成
- 在每条返回路径上关闭 HTTP 响应体，以防止资源泄漏，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6434)中完成
- 开发代理重试请求时保留请求体，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6432)中完成
- 发现桌面文件浏览器时遵守启动截止时间，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6430)中完成
- 保存 AppImage 工具下载内容前拒绝 HTTP 错误响应，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6429)中完成
- 在语义化版本比较中正确排序超大数值型预发布标识符，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6427)中完成
- 在资源服务器中路由服务时保留转义后的路径，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6426)中完成
- 在资源服务器中处理捆绑运行时的 HEAD 请求，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6425)中完成
- 缺少捆绑资源时返回404，而不是错误响应，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6424)中完成
- 请求被取消时停止开发代理重试，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6423)中完成
- 避免资源服务器共享的起始 URL 被修改以及使用过时端口，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6422)中完成
- 在资源服务器中遵循可选的日志记录器配置，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6421)中完成
- 消除库路径缓存失效处理中的数据竞态，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6418)中完成
- 搜索 ldconfig 输出时避免错误匹配库，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6416)中完成
- 防止资源服务器队列在负载下增长并反复收缩，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6415)中完成
- 防抖计时器触发后释放捕获的回调，以避免保留过时的引用，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6413)中完成
- 保留在可执行文件路径内找到的文件浏览器目标，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/6412)中完成
- 在文件资源管理器中跳过遮蔽桌面条目的目录，见[PR](https://github.com/wailsapp/wails/pull/6411)，作者：@taliesin-ai
- 修复响应正文为空时资产服务器不发送响应头的问题，见[PR](https://github.com/wailsapp/wails/pull/6409)，作者：@taliesin-ai
- 查找应用程序图标时遵循 XDG 桌面文件搜索目录，见[PR](https://github.com/wailsapp/wails/pull/6407)，作者：@taliesin-ai
- 在文件资源管理器中展开字段时保留文件路径，见[PR](https://github.com/wailsapp/wails/pull/6406)，作者：@taliesin-ai
- 避免在文件资源管理器中重复展开后的 URI 目标，见[PR](https://github.com/wailsapp/wails/pull/6405)，作者：@taliesin-ai
- 提交成功但为空的不可寻址资产响应，不再让其保持待处理状态，见[PR](https://github.com/wailsapp/wails/pull/6403)，作者：@taliesin-ai
- 资产服务器在多次短读取中正确填充 MIME 嗅探前缀，见[PR](https://github.com/wailsapp/wails/pull/6402)，作者：@taliesin-ai
- 避免资产服务器为不可寻址的 HEAD 请求读取完整响应正文，见[PR](https://github.com/wailsapp/wails/pull/6401)，作者：@taliesin-ai
- 资产服务器解析目录请求时仅提供精确匹配的索引文件，见[PR](https://github.com/wailsapp/wails/pull/6400)，作者：@taliesin-ai
- 修复在不支持操作的服务器上，已关闭的 Linux 通知被报告为点击的问题，见[PR](https://github.com/wailsapp/wails/pull/6367)，作者：@bercedev
- 终端关闭时停止开发进程，见[PR](https://github.com/wailsapp/wails/pull/6371)，作者：@xmasdev
- Windows 任务栏重启的防抖标志延迟一秒后再清除，见[PR](https://github.com/wailsapp/wails/pull/6386)，作者：@taliesin-ai
- 修复资产服务器丢失回退响应头的问题，见[PR](https://github.com/wailsapp/wails/pull/6399)，作者：@taliesin-ai
- 资产服务器在选择回退响应前遵循响应头删除操作，见[PR](https://github.com/wailsapp/wails/pull/6398)，作者：@taliesin-ai
- 释放资产服务器中已移除的环形队列引用以避免内存泄漏，见[PR](https://github.com/wailsapp/wails/pull/6397)，作者：@taliesin-ai
- 释放资产服务器中已用完的内容嗅探缓冲区以减少内存使用，见[PR](https://github.com/wailsapp/wails/pull/6396)，作者：@taliesin-ai
- 提供资产时识别大小写混合的资产文件扩展名，见[PR](https://github.com/wailsapp/wails/pull/6395)，作者：@taliesin-ai
- 资产服务器无需保留文件名即可正确检测内容类型，见[PR](https://github.com/wailsapp/wails/pull/6394)，作者：@taliesin-ai
- 桌面条目没有命令时回退到默认值，见[PR](https://github.com/wailsapp/wails/pull/6392)，作者：@taliesin-ai
- 防止 liquid-glass 示例在非 macOS 系统上于应用启动前显示对话框，见[PR](https://github.com/wailsapp/wails/pull/6389)，作者：@taliesin-ai
- 在 Windows 上从受约束的尺寸进入全屏时保留窗口尺寸约束，见[PR](https://github.com/wailsapp/wails/pull/6387)，作者：@taliesin-ai
- 修复在 Run 前调用对话框时永久阻塞而不返回的问题，见[PR](https://github.com/wailsapp/wails/pull/6390)，作者：@taliesin-ai

## 安全性
- 在分配图像内存前验证生成的图标尺寸，防止内存使用过多，见[PR](https://github.com/wailsapp/wails/pull/6478)，作者：@taliesin-ai
- 使用发行说明 API 响应前先进行验证，见[PR](https://github.com/wailsapp/wails/pull/6461)，作者：@taliesin-ai
- SQLite 插件要求提供完整的十六进制语句句柄，见[PR](https://github.com/wailsapp/wails/pull/6452)，作者：@taliesin-ai
- 限制文件系统库的回退匹配，防止加载非预期库，见[PR](https://github.com/wailsapp/wails/pull/6417)，作者：@taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## 变更
- 将404页面链接到v2.wails.io 上相匹配的 Wails v2文档页面，见[PR](https://github.com/wailsapp/wails/pull/6366)，作者：@taliesin-ai
- 将v2文档站点迁移到v2.wails.io，见[PR](https://github.com/wailsapp/wails/pull/6365)，作者：@taliesin-ai
- 通过 wails.io 提供v3文档，并通过v2.wails.io 提供v2文档，见[PR](https://github.com/wailsapp/wails/pull/6364)，作者：@taliesin-ai
- 更新v3.0.0发行说明并更正安全公告中的说法，见[PR](https://github.com/wailsapp/wails/pull/6363)，作者：@taliesin-ai

## 修复
- 按触发顺序串行传递 Go 事件监听器回调，见[PR](https://github.com/wailsapp/wails/pull/6275)，作者：@DevLumuz
- 始终公布托盘的 Menu 属性，避免在 Linux 上丢弃 Run 后创建的托盘，见[PR](https://github.com/wailsapp/wails/pull/6284)，作者：@AlbinoGeek
- 将 Run 前提交的主线程工作暂存至应用程序启动后再执行，见[PR](https://github.com/wailsapp/wails/pull/6287)，作者：@AlbinoGeek
- 下载 GitHub 更新程序的发行版资产时使用正确的资产 API URL，见[PR](https://github.com/wailsapp/wails/pull/6304)，作者：@gavidroselj
- 修复 macOS 和 Linux 混合 DPI 配置下的屏幕边界和位置，见[PR](https://github.com/wailsapp/wails/pull/6346)，作者：@julianstorer
- 在 Windows 上，当 WebView2 省略 Content-Length 响应头时保留请求正文，见[PR](https://github.com/wailsapp/wails/pull/6326)，作者：@strraaa
- 在锁内创建钩子快照，修复事件钩子触发时的数据竞争，见[PR](https://github.com/wailsapp/wails/pull/6276)，作者：@DevLumuz
- 修复`SystemTray.SetTooltip()`在 Linux 上不起作用的问题：现在会设置 StatusNotifierItem 的工具提示描述（#6017）
- 在每周发布工作流中运行变更日志证据检查，见[PR](https://github.com/wailsapp/wails/pull/6374)，作者：@taliesin-ai
- 修复发行版变更日志证据检查，将运行记录映射到压缩合并提交，见[PR](https://github.com/wailsapp/wails/pull/6373)，作者：@taliesin-ai
- 修复 Linux 深色模式检测在深色桌面上返回浅色，且要等配色方案发生变化才恢复正常的问题，以及`linux:SystemThemeChanged`在 GTK4 后端上从不触发的问题。两个后端现在共用一套实现，由`org.freedesktop.appearance`门户命名空间提供支持（#6072）
- 在带边框的窗口（Windows、macOS、Linux）上，`Window.SetZoom`不再将低于1.0的系数提高到1.0：WebView2 `ZoomFactor`、WKWebView `magnification`和 WebKitGTK 的缩放级别均接受任意正值，而`Zoom`窗口选项原本就会直接传递此类值。无边框窗口仍保留1.0下限，因为其拖动和调整大小区域不会根据缩放进行调整（[#4590](https://github.com/wailsapp/wails/issues/4590)）；现在会忽略非正数、NaN 和 Inf 系数，而不再将其提高到1.0 — 作者：@APshenkin

## v3.0.0 - 2026-10-26

## 变更
- 夜间预发行版以`vX.Y.0-nightly.YYYYMMDD`形式发布，用于预览下一个每周发行版；每周发行版会递增次版本号（`vX.Y.0`），因此 Go 的`@latest`始终解析到每周发行版或热修复发行版；`@wailsio/runtime` npm 软件包版本始终与 Wails 发行版本一致，夜间版本使用`nightly` dist-tag
- 项目模板将`@wailsio/runtime`固定为生成这些模板的 CLI 版本
- Wails v3要求 Go 1.26或更高版本（Go 1.26.10或1.27.3及更高版本包含`crypto/x509`安全修复）

## 修复
- 修复`wails3 update cli`找不到v3发行版的问题：它只考虑v2标签，并且从错误的路径安装
- 修复部分 Windows 版本在绘制应用程序菜单栏时因`panic: SetTextColor failed`而崩溃的问题：菜单栏现在会在有效的设备上下文上绘制，或回退到默认绘制；`UAHMENUITEM.Position`现在与 C 的`int`一致；菜单绘制不再引发 panic（#6311）
- 确保等待重启期间暂存更新的安全：关机过程中启动的下载或第二次重启无法删除交换辅助程序正等待安装的有效负载，也不会与其争用，作者 @julianstorer
- 修复 Windows 应用在`applyDarkMode`中崩溃的问题：`pkg/w32`初始化时卸载了`uxtheme.dll`，但其序号函数指针仍在使用。现在该模块会在进程的整个生命周期内保持加载（#6356）
- 在交换前验证暂存的有效负载，并在丢弃备份前验证已安装的应用程序，使损坏的更新能够回滚，而不会替换正常工作的应用程序，作者 @julianstorer

## 已移除
- 在[PR](https://github.com/wailsapp/wails/pull/5736)中完成最后一次翻译同步后移除 Crowdin 集成，作者 @github-actions[bot]

## 安全性
- 通过向所有 frame 注入防护脚本，阻止 Linux（GTK3 和 GTK4）上跨源子 frame 发送保留的`wails:`消息，从而缓解 GHSA-2p5f-2gx8-mh8m
- 将`github.com/klauspost/compress`更新至1.18.7（GHSA-259r-337f-4rfw），将`golang.org/x/image`更新至0.45.0，并使用 Go 1.26.10构建以纳入`crypto/x509`修复（GO-2026-6809）
- 在运行时包、设置向导、示例和项目模板中，将`vite`更新至8.2.4或更高版本（GHSA-9jrq-w75r-8gcw、GHSA-rq7h-c2jc-7f22、GHSA-vfpm-58rq-9qcg、GHSA-r825-w8mm-5pxw、GHSA-587h-7wmj-h5pm、GHSA-698v-4r3x-3vcv）

### 已知限制
- **Linux**：由于 WebKitGTK 不向原生消息处理程序提供任何帧信息，子帧防护机制无法阻止*同源* 子帧（例如`srcdoc`或`about:blank` iframe）通过其父帧访问原生消息处理程序。跨源帧以及不含`allow-same-origin`的沙盒帧会被阻止。请勿在 Wails 窗口中嵌入不受信任的同源内容。
- **Android（实验性）：**`window.wails` JavaScript 桥接会暴露给 WebView 中的每个 frame，且不会检查主 frame 或来源。嵌入式 iframe 中的内容（包括跨源内容）可以调用绑定方法并发送保留的`wails:`消息。请勿在应用的 WebView 中嵌入不受信任的第三方内容。这与 macOS 和 iOS 不同，后两者会忽略来自子 frame 的`wails:`消息。计划在3.0.1中修复此问题。
- <strong>所有平台：</strong>来自子 frame 的绑定方法调用通过 HTTP 传输至`/wails/runtime`，而非原生消息处理程序，因此防护不会限制这些调用。应将嵌入的任何 iframe 视为能够访问绑定的 Go 方法。

## v3.0.0-beta.34 - 2026-10-20

## 已修复
- 修复激活窗口时 Windows WebView 无法重新获得焦点，以及 Tab 键被困在 WebView 内而无法将焦点移至原生窗口再移回的问题（#3783）
- 修复`internal/webview2/pkg/webview2` COM vtable 结构体（`ICoreWebView2_2`..`_27`及其他派生接口）遗漏继承的方法链，导致调用落到错误 vtable 槽位的问题（#6129）
- 修复 Windows 询问对话框忽略自定义按钮标签且从不触发按钮回调的问题；自定义按钮现在通过 TaskDialog 显示（comctl32 v6不可用时，回退到使用按位置映射回调的 MessageBox）（#5794）
- 修复针对 Vite 的60秒`wails3 dev`挂起的问题：自动前端就绪门控现在使用 TCP 探测，而不再使用未发送`Accept`标头并收到404的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 在[PR](https://github.com/wailsapp/wails/pull/6301)中确保系统托盘菜单在打开期间更新后仍能正常工作，作者 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6302)中妥善处理运行时操作期间的瞬态 WebView2 COM 错误，作者 @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## 新增
- 绑定生成器现在会生成一个`bindings/events`模块，为使用编译期常量名称通过`application.RegisterEvent`注册的每个可发现事件导出带类型的名称常量（例如，为`user:login`导出`UserLogin`）；无需 Vite 插件即可与`Events.On`/`Emit`配合使用（#5645）

## 已修复
- 修复 Windows 询问对话框忽略自定义按钮标签且从不触发按钮回调的问题；自定义按钮现在通过 TaskDialog 显示（comctl32 v6不可用时，回退到使用按位置映射回调的 MessageBox）（#5794）
- 修复针对 Vite 的60秒`wails3 dev`挂起的问题：自动前端就绪门控现在使用 TCP 探测，而不再使用未发送`Accept`标头并收到404的 HTTP GET（[#6209](https://github.com/wailsapp/wails/issues/6209)）
- 在[PR](https://github.com/wailsapp/wails/pull/6301)中确保系统托盘菜单在打开期间更新后仍能正常工作，作者 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6302)中妥善处理运行时操作期间的瞬态 WebView2 COM 错误，作者 @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## 新增
- 在[PR](https://github.com/wailsapp/wails/pull/6297)中为 macOS 14+ 添加原生 macOS 菜单分区标题——请参阅[文档网站](https://v3.wails.io)（已更新20个页面），作者 @leaanthony

## 已修复
- 在[PR](https://github.com/wailsapp/wails/pull/6299)中使 SetRelativePosition 在 macOS 上相对于屏幕工作区定位，作者 @leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6298)中，使更新辅助程序在父进程超时或替换失败后重新启动未更改的目标程序，作者 @taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6293)中将版本化文档使用的 pnpm 固定为10版本，作者 @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## 新增
- 公开`@wailsio/runtime/cancellable`子路径导出，使`CancellablePromise`可以在不执行运行时根入口的情况下导入（#6153）

## 已更改
- 要求使用 Go 1.26或更高版本：v3模块现在声明`go 1.26.0`，由`wails3 init`和`wails3 service init`生成的项目也声明相同的最低版本

## 已修复
- 修复使用 Go 1.27构建时绑定生成器输出发生变化的问题。当时，`Marshaler`、`Unmarshaler`和`RawMessage`等`encoding/json`类型变成了`encoding/json/v2`和`encoding/json/jsontext`类型的别名；现在它们会再次渲染为`any`，并且不会为`encoding/json/v2`生成模型（#6292）
- 克隆前验证传给`wails3 init -t`的远程模板 URL：只接受`https://`、`http://`、`ssh://`、`git://`和`user@host:path` URL；值会在`--`之后传给 git；克隆期间 git 仅允许使用这些传输协议
- 拒绝 HTTP 运行时传输中超出范围的`x-wails-window-id`标头值，而不是将其静默截断为另一个窗口的 ID
- 增强更新程序的归档解压安全性以抵御 zip-slip：现在使用`filepath.IsLocal`检查条目路径，并对拼接后的路径执行包含范围检查；同时为 zip 和 tar.gz 归档中嵌套的`..`、反斜杠、绝对路径及通过符号链接写入的条目添加了测试
- 修复通过`wails3 generate appimage`生成的 AppImage 在 Debian 和 Ubuntu 之外的系统上启动时因“Unable to spawn a new child process: … WebKitNetworkProcess”而崩溃的问题。捆绑的 WebKit 现在会运行 AppImage 随附的辅助进程。WebKit 的 bubblewrap 沙箱仍保持启用
- 修复`-tags server`构建无法为 Android 和 iOS 编译的问题，并修复[#6271](https://github.com/wailsapp/wails/issues/6271)
- 修复更新程序在临时目录位于另一文件系统（例如 tmpfs `/tmp`）时无法在 Linux 上完成替换的问题：现在会先在已安装的二进制文件旁暂存，再跨文件系统复制
- 修复更新程序尝试替换 AppImage 的只读挂载点而非`.AppImage`文件的问题
- 修复 macOS 和 iOS 接受来自 iframe 的保留`wails:`窗口控制消息的问题。现在仅接受来自窗口主框架的此类消息；来自 iframe 的其他消息仍会连同其来源信息传递给`RawMessageHandler`
- 停止在 macOS 上发送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`会在运行循环的每次迭代中触发，因此页面每秒会收到约120次空的`evaluateJavaScript:`调用；在 macOS 27.0上，每次调用都会消耗用户点击产生的瞬时激活状态，导致`navigator.clipboard.writeText()`及其他受激活状态限制的 API 拒绝真实点击并返回`NotAllowedError`（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 贡献
- 将 Windows info.json 的默认语言修复为 en-US（0409），以便资源管理器显示版本详细信息（#6210）
- 修复在配有 NVIDIA 独立 GPU 以及 AMD 或 Intel 集成 GPU 的笔记本电脑上，Linux DMA-BUF 渲染器被禁用的问题：在 Wayland 上，现在仅当 WebKit 使用 NVIDIA GPU 渲染时（应用被卸载到独立 GPU 上运行，或合成器在该 GPU 上运行）才禁用该渲染器，而不是只要加载了 nvidia 模块就禁用；后者会使集成 GPU 上的渲染慢如蜗牛。X11 会话仍保留原有行为，因为 Xorg 使用的渲染 GPU 由其配置决定
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中通过暂存复制回退方案处理 Linux 更新程序的跨设备重命名失败，由 @john-okeefe 贡献
- 修复 Windows 上的`WebviewWindow.SetMenu`仅在配置了菜单栏主题时才重绘菜单栏的问题，并修复[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修复`WebviewWindow.SetMenu(nil)`在 Windows 上触发 panic，而不是像在 macOS 上一样不执行任何操作的问题，并修复[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送之间保留变更日志条目，由 @leaanthony 贡献
- 移除引用了模块中从未随附的二进制文件的嵌入项，修复`go mod vendor`在所有平台上均因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，并修复[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)；此项修复由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 上重复拆除窗口，由 @charlesbalderson 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中处理系统托盘菜单的 WM_CONTEXTMENU，使其可通过右键单击打开，由 @wzk2025 贡献
- 修复所有五种`InvokeSync*`变体在分派的回调触发 panic 时永久阻塞调用方的问题，并修复[#6107](https://github.com/wailsapp/wails/issues/6107)；此项修复由 @RALIST 在[PR](https://github.com/wailsapp/wails/pull/6108)中贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6236)中安全地重建已销毁的 GTK3 子菜单，由 @leaanthony 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6261)中为嵌套 worker 实现 Windows worker 转发垫片，由 @taliesin-ai 贡献
- 修复 Windows 上的`OpenFile()`在`CanChooseDirectories(true)`时忽略`PromptForMultipleSelection()`的问题：用户现在可以像在 macOS 和 Linux 上一样一次选择多个文件夹（#6263）
- 修复 windows/arm64 上按值传递给 WebView2 COM 方法的`double`参数：`ZoomFactor`、`RasterizationScale`、Cookie 的`Expires`，以及打印设置的缩放比例、页面大小和边距均受影响（打印设置的 setter 此前会返回`ErrDoubleArgUnsupported`）。受影响的调用现在会经过一个可加载`d0`的`runtime.cgocall`跳板，而 Go 的 syscall 路径无法做到这一点（golang.org/issue/62583），由 @APshenkin 贡献

## 安全性
- 服务器模式（`-tags server`）现在仅接受以`POST`请求形式发出的`/wails/runtime`调用，且请求必须来自服务器自身的源或`ServerOptions.WebSocketOriginPatterns`中列出的源；使用其他方法的请求会收到 HTTP 405，来自其他浏览器源的请求会收到 HTTP 403。如果部署位于会重写`Host`标头的反向代理之后，应将其公共源添加到`WebSocketOriginPatterns`
- 绑定到环回地址（默认的`localhost`、`127.0.0.0/8`或`::1`）的服务器模式（`-tags server`）现在会拒绝请求并返回 HTTP 403，除非`Host`标头为`localhost`、IP 地址字面量、绑定主机或与`ServerOptions.WebSocketOriginPatterns`匹配的主机，从而阻止针对`/wails/runtime`和 WebSocket 端点的 DNS 重绑定攻击
- Windows 通知在 WinRT Toast API 失败时不再回退到 PowerShell，而是返回错误

## v3.0.0-beta.30 - 2026-10-15

## 修复
- 修复更新程序尝试替换 AppImage 的只读挂载点而非`.AppImage`文件的问题
- 停止在 macOS 上发送`WindowDidUpdate`和`WindowWillUpdate`：`NSWindowDidUpdateNotification`会在运行循环的每次迭代中触发，因此页面每秒会收到约120次空的`evaluateJavaScript:`调用；在 macOS 27.0上，每次调用都会消耗用户点击产生的瞬时激活状态，导致`navigator.clipboard.writeText()`及其他受激活状态限制的 API 拒绝真实点击并返回`NotAllowedError`（[#6170](https://github.com/wailsapp/wails/issues/6170)），由 @APshenkin 贡献
- 将 Windows info.json 的默认语言修复为 en-US（0409），以便资源管理器显示版本详细信息（#6210）
- 在[PR](https://github.com/wailsapp/wails/pull/6135)中通过暂存复制回退方案处理 Linux 更新程序的跨设备重命名失败，由 @john-okeefe 贡献
- 修复 Windows 上的`WebviewWindow.SetMenu`仅在配置了菜单栏主题时才重绘菜单栏的问题，并修复[#6103](https://github.com/wailsapp/wails/issues/6103)
- 修复`WebviewWindow.SetMenu(nil)`在 Windows 上触发 panic，而不是像在 macOS 上一样不执行任何操作的问题，并修复[#6104](https://github.com/wailsapp/wails/issues/6104)
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送之间保留变更日志条目，由 @leaanthony 贡献
- 移除引用了模块中从未随附的二进制文件的嵌入项，修复`go mod vendor`在所有平台上均因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，并修复[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)；此项修复由 @Grantmartin2002 在[PR](https://github.com/wailsapp/wails/pull/6031)中贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6159)中防止 macOS 上重复拆除窗口，由 @charlesbalderson 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6162)中处理系统托盘菜单的 WM_CONTEXTMENU，使其可通过右键单击打开，由 @wzk2025 贡献
- 修复所有五种`InvokeSync*`变体在分派的回调触发 panic 时永久阻塞调用方的问题，并修复[#6107](https://github.com/wailsapp/wails/issues/6107)；此项修复由 @RALIST 在[PR](https://github.com/wailsapp/wails/pull/6108)中贡献

## v3.0.0-beta.29 - 2026-10-13

## 新增
- macOS 上的全局快捷键现在可感知键盘布局：在非 QWERTY 布局（AZERTY、QWERTZ、Dvorak 等）中，字母加速键会通过`UCKeyTranslate`绑定到标有相应字母的按键，而不是固定的 QWERTY 物理键位。

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6233)中，先在 Linux 上初始化原生菜单，再进行设置，由 @abhicodes-007 贡献
- 在[PR](https://github.com/wailsapp/wails/pull/6237)中解决 macOS 正式发布准备阶段在输出处理和设置向导方面的回归问题，由 @taliesin-ai 贡献
- 修复 Windows 上的 WebView2 COM `AddRef`/`Release`包装器将过时的线程最终错误报告为失败的问题。它们现在会返回引用计数，因此释放操作不会再中止应用程序或导致资源请求失败（#3931）
- 修复`ICoreWebView2Deferral.Release`调用`AddRef`槽的问题；该问题会导致 Windows 上每个延迟处理的资源请求都泄漏其延迟操作（#3931）

## v3.0.0-beta.28 - 2026-10-05

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6225)中将 macOS 应用输出保留在文件提供程序之外，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6224)中完善设置向导并显示系统详细信息，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6222)中改进 Windows 上的 WebView2 恢复错误处理和资源清理，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6220)中改进 Windows 上的 WebView2 恢复机制，防止应用崩溃和控制器滞留，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6221)中更正 macOS 媒体授权指南，作者：@taliesin-ai
- 修复右键单击 Linux 系统托盘图标时，既运行单击处理程序又打开菜单的问题；该问题会导致关联窗口在每次右键单击时切换状态（#6018）
- 在[PR](https://github.com/wailsapp/wails/pull/6184)中使 Windows 上的 WebView2 回调处理程序保持存活，直至原生层将其释放，作者：@taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## 修复
- 修复 macOS 12及更高版本忽略`Permissions`选项的问题：现已实现`WKUIDelegate`媒体捕获方法，因此`PermissionAllow`和`PermissionDeny`会像在 Linux 和 Windows 上一样应用于摄像头和麦克风请求。仍需提供`NSCameraUsageDescription`/`NSMicrophoneUsageDescription`；使用沙盒时，还需提供相应的设备权限声明

## v3.0.0-beta.26 - 2026-09-25

## 新增
- 在[PR](https://github.com/wailsapp/wails/pull/6150)中使监视器对前端开发服务器执行 HTTP 就绪检查，作者：@atterpac

## 修复
- 在[PR](https://github.com/wailsapp/wails/pull/6002)中使应用从 WebView2 进程故障中恢复，而不是留下空白窗口，作者：@taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## 修复
- 主应用退出时停止`wails3 dev`及其后台进程（#6048）

## v3.0.0-beta.24 - 2026-09-20

## 安全
- 限制构建和发布工作流中的 GITHUB_TOKEN 权限，以遵循最小权限原则
- 修复屏幕示例资源中间件中的路径遍历和符号链接逃逸问题
- 在将设置向导的依赖项安装限制为获准的软件包管理器命令的同时，确保其仍可正常工作

## v3.0.0-beta.23 - 2026-09-16

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5972)中允许覆盖 Linux 上的 GTK 应用 ID，作者：@overlordtm

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5973)中使 Linux 单实例机制以应用自身的 ID 声明总线名称，作者：@overlordtm
- 在[PR](https://github.com/wailsapp/wails/pull/6126)中取消 Linux 上已中止的自定义方案请求，作者：@leaanthony

## v3.0.0-beta.22 - 2026-09-14

## 新增

- 完成翻译并迁移社区媒体内容——参见[文档站点](https://v3.wails.io)（更新了1812个页面）；相关更改见[PR](https://github.com/wailsapp/wails/pull/6122)，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6120)中记录可嵌入的 Wails 按钮及其各种样式和尺寸，作者：@leaanthony

## v3.0.0-beta.21 - 2026-09-13

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/6116)中使用 M-Press 提供 Wails v3文档，作者：@leaanthony

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/6118)中解析 MPD frontmatter 中的 JSON slug 值，以生成变更日志，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6080)中使更新程序在备份失败后清除辅助环境变量并重新启动原始目标，作者：@cnmax
- 在[PR](https://github.com/wailsapp/wails/pull/6098)中于 App.Run 期间启动默认信号处理程序，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6112)中使 Windows 菜单能够处理 nil 菜单、释放被替换的资源并重绘菜单栏，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/6115)中使用共享 YAML 配置恢复新建项目的 MSIX 打包功能，作者：@leaanthony
- 修复以下问题：当泛型模型创建器引用后面声明的辅助项时，生成的 JavaScript 和 TypeScript 绑定无法加载；以及创建相互依赖的泛型模型时发生堆栈溢出（#6062）

## v3.0.0-beta.20 - 2026-09-10

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/6082)中将 Clave 展示链接更新为当前网站和仓库，作者：@01xR4in

## 修复

- 取消 Windows 上已中止的资源请求（包括 worker 请求），同时在导航过程中保留 keepalive 处理程序。在 Apple 平台上，通过应用包装器转发原生请求上下文。（#5963、#5969）
- 在[PR](https://github.com/wailsapp/wails/pull/6094)中通过重试机制，在相互竞争的推送期间保留变更日志条目，作者：@leaanthony
- 在[PR](https://github.com/wailsapp/wails/pull/6031)中移除引用模块中从未发布的二进制文件的嵌入项，修复`go mod vendor`在所有平台上均因`pattern arm64/WebView2Loader.dll: no matching files found`而失败的问题，并修复[#5782](https://github.com/wailsapp/wails/issues/5782)和[#5376](https://github.com/wailsapp/wails/issues/5376)，作者：@Grantmartin2002

## 移除

- 移除已由纯 Go 加载器取代的原生 WebView2 加载器支持。此项更改移除了嵌入的`WebView2Loader.dll`二进制文件和`github.com/jchv/go-winloader`依赖项。`native_webview2loader`构建标签仍会被接受且不再报错，但对v3构建没有任何影响。相关更改见[PR](https://github.com/wailsapp/wails/pull/6031)，作者：@Grantmartin2002
- 在[PR](https://github.com/wailsapp/wails/pull/6097)中从 macOS API 指南移除未使用的构建标签和 FPS 选项，作者：@leaanthony

## v3.0.0-beta.19 - 2026-09-09

## 新增

- 通过构建标签限制私有 macOS API，仅允许选择启用——参见[文档](https://v3.wails.io/features/browser/integration)、[文档](https://v3.wails.io/features/environment/info)、[文档](https://v3.wails.io/features/windows/basics)、[文档](https://v3.wails.io/features/windows/frameless)、[文档](https://v3.wails.io/features/windows/notch-windows)、[文档](https://v3.wails.io/features/windows/options)、[文档](https://v3.wails.io/guides/build/macos)、[文档](https://v3.wails.io/guides/build/private-macos-apis)和[文档](https://v3.wails.io/reference/overview)；相关更改见[PR](https://github.com/wailsapp/wails/pull/6087)，作者：@leaanthony

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/6091)中拒绝超过64 MiB 的运行时请求，并返回 HTTP 413，作者：@leaanthony

## 安全

- 在[PR](https://github.com/wailsapp/wails/pull/6092)中通过令牌身份验证强化 MCP 来源和远程访问的安全性，作者：@leaanthony

## v3.0.0-beta.18 - 2026-09-08

## 修复

- 通过使用指针接收者修复 Linux 和 Darwin 上的 Calloc 内存泄漏，见[PR](https://github.com/wailsapp/wails/pull/6083)，由 @4RH1T3CT0R7 提交

## v3.0.0-beta.17 - 2026-09-06

## 修复

- Windows：WebResourceRequested 处理程序中失败或为 nil 的`GetRequest`不再导致进程终止（`log.Fatal`/nil 解引用 panic）；现在会丢弃该请求并记录日志，见[PR](https://github.com/wailsapp/wails/pull/6006)，由 @midagedev 提交

## v3.0.0-beta.16 - 2026-08-29

## 变更

- 在新的终端窗口中提示输入公证密码，见[PR](https://github.com/wailsapp/wails/pull/6029)，由 @leaanthony 提交

## 修复

- 正确处理 macOS 上的系统托盘点击类型，见[PR](https://github.com/wailsapp/wails/pull/5919)，由 @ChewbaccaCookie 提交
- CI 在更新前移除未使用的 Microsoft apt 软件源，见[PR](https://github.com/wailsapp/wails/pull/6041)，由 @Grantmartin2002 提交

## v3.0.0-beta.15 - 2026-08-27

## 修复

- 将 WebView2 嵌入超时时间增加至60秒，见[PR](https://github.com/wailsapp/wails/pull/6043)，由 @Grantmartin2002 提交

## v3.0.0-beta.14 - 2026-08-26

## 修复

- 在 macOS 上正确命名 Control-字母组合键的按键操作，见[PR](https://github.com/wailsapp/wails/pull/6032)，由 @taliesin-ai 提交
- 修复 Windows 上的 ICO 托盘图标并使其跟随任务栏主题，见[PR](https://github.com/wailsapp/wails/pull/6016)，由 @nik9play 提交

## v3.0.0-beta.13 - 2026-08-25

## 修复

- 在 macOS 上运行模态循环时继续处理主线程任务，见[PR](https://github.com/wailsapp/wails/pull/6026)，由 @leaanthony 提交
- 使移动端安全存储操作可以失败，并在失败时默认拒绝访问，见[PR](https://github.com/wailsapp/wails/pull/5923)，由 @mortenolsrud 提交
- 即使未注册监听器，也运行应用程序事件钩子，见[PR](https://github.com/wailsapp/wails/pull/5999)，由 @archy-rock3t-cloud 提交
- 更正注释和本地化文档中的拼写错误，见[PR](https://github.com/wailsapp/wails/pull/6023)，由 @haoku123 提交
- 移除提交在`v3/examples`下的预编译 macOS 二进制文件，见[PR](https://github.com/wailsapp/wails/pull/6025)，由 @4RH1T3CT0R7 提交

## v3.0.0-beta.12 - 2026-08-21

## 新增

- 新增带生命周期和遥测示例的 macOS 刘海通知窗口——参阅[文档](https://v3.wails.io/features/windows/notch-windows)，见[PR](https://github.com/wailsapp/wails/pull/6010)，由 @leaanthony 提交
- 新增 macOS NSPanel 窗口支持，包括新选项和原生集成——参阅[文档](https://v3.wails.io/features/windows/options)，见[PR](https://github.com/wailsapp/wails/pull/6008)，由 @leaanthony 提交

## 修复

- 防止并发调用时 SQLite Prepare 挂起，见[PR](https://github.com/wailsapp/wails/pull/5998)，由 @archy-rock3t-cloud 提交

## v3.0.0-beta.11 - 2026-08-20

## 移除

- 从文档中移除已过时的实现跟踪器，见[PR](https://github.com/wailsapp/wails/pull/6005)，由 @leaanthony 提交

## v3.0.0-beta.10 - 2026-08-19

## 修复

- 修复 GTK4 Linux 宿主丢弃自定义协议和文件关联启动参数的问题，见[PR](https://github.com/wailsapp/wails/pull/6000)，由 @midagedev 提交
- 在变更日志验证中正确处理已删除的行和同源更正，见[PR](https://github.com/wailsapp/wails/pull/5993)，由 @taliesin-ai 提交

## v3.0.0-beta.9 - 2026-08-16

## 新增

- 新增安全的 wails3 mcp 服务器，用于智能体辅助的项目管理，见[PR](https://github.com/wailsapp/wails/pull/5896)，由 @leaanthony 提交
- 新增绑定中模型的文档——参阅[文档](https://v3.wails.io/features/bindings/models)，见[PR](https://github.com/wailsapp/wails/pull/5988)，由 @taliesin-ai 提交
- 支持在原子化 Linux 系统上通过 rpm-ostree 安装，见[PR](https://github.com/wailsapp/wails/pull/5987)，由 @leaanthony 提交
- 新增原生的每周 star 历史图表生成与发布功能——参阅[文档](https://v3.wails.io/credits)、[文档](https://v3.wails.io/de/credits)、[文档](https://v3.wails.io/fr/credits)、[文档](https://v3.wails.io/id/credits)、[文档](https://v3.wails.io/ja/credits)、[文档](https://v3.wails.io/ko/credits)、[文档](https://v3.wails.io/pt/credits)、[文档](https://v3.wails.io/ru/credits)、[文档](https://v3.wails.io/zh-cn/credits)及[文档](https://v3.wails.io/zh-tw/credits)，见[PR](https://github.com/wailsapp/wails/pull/5986)，由 @leaanthony 提交
- 新增仅限 Darwin 的 mac 包，用于解析应用程序包资源——参阅[文档](https://v3.wails.io/guides/build/macos)，见[PR](https://github.com/wailsapp/wails/pull/5965)，由 @leaanthony 提交
- 新增 Condui 展示页面和索引条目——参阅[文档](https://v3.wails.io/community/showcase/condui)及[文档](https://v3.wails.io/community/showcase)，见[PR](https://github.com/wailsapp/wails/pull/5962)，由 @mgueregath 提交
- 新增带屏幕截图和项目链接的 Redis Viewer 展示页面——参阅[文档](https://v3.wails.io/community/showcase)及[文档](https://v3.wails.io/community/showcase/redisviewer)，见[PR](https://github.com/wailsapp/wails/pull/5984)，由 @redisviewer 提交

## 变更

- 在 Linux 上将 GTK 应用程序标志更新为 G<em>APPLICATION</em>NON_UNIQUE，见[PR](https://github.com/wailsapp/wails/pull/5971)，由 @overlordtm 提交
- 将缺失窗口事件的日志级别从警告改为调试，见[PR](https://github.com/wailsapp/wails/pull/5914)，由 @julianstorer 提交

## 修复

- 允许已注册的 macOS 键盘快捷键优先于 webview 处理，见[PR](https://github.com/wailsapp/wails/pull/5902)，由 @julianstorer 提交
- 修复文档侧边栏中失效的链接，见[PR](https://github.com/wailsapp/wails/pull/5937)，由 @northes 提交
- 当 WebKit 中止匹配的自定义方案任务时，取消 macOS 和 iOS 资源请求的上下文（#5963）
- 处理 WindowSetFullscreenButtonEnabled 消息，见[PR](https://github.com/wailsapp/wails/pull/5976)，由 @archy-rock3t-cloud 提交
- 在 preact-ts 模板中导入 Fragment 以解决构建失败，见[PR](https://github.com/wailsapp/wails/pull/5979)，由 @haoku123 提交
- 防止旧版仅服务 GTK3 应用程序在尚无活动窗口或可用显示器时执行屏幕发现而崩溃（#5966）
- 允许显式指定版本的发布流程在未发布变更日志为空时继续执行（#5977）

## 安全

- 将网站的 nanoid 锁文件更新至已修补的3.3.18版本，以解决安全公告中的问题，见[PR](https://github.com/wailsapp/wails/pull/5985)，由 @taliesin-ai 提交

## v3.0.0-beta.8 - 2026-08-12

## 新增

- 为自动变更日志条目新增文档 URL 生成功能，见[PR](https://github.com/wailsapp/wails/pull/5957)，由 @taliesin-ai 提交
- 新增 Streams：采用 WebSocket 编程模型，在 Go 与 JavaScript 之间提供双向字节流，且无需监听套接字。在 Go 中使用`app.HandleStream(name, handler)`声明流，并在前端使用`Stream(name)`连接；后者返回一个形如`WebSocket`的对象。Go→JS 通过资源服务器上每个窗口一个保持的轮询传输，JS→Go 则通过普通 POST 传输；不会绑定任何 TCP 端口，也不会有任何内容经过`evaluateJavaScript`。在服务器构建（`-tags server`）中，同一处理程序改为通过真正的 WebSocket 提供服务，因此应用程序代码在不同构建中完全相同。由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5935)中将 mailbox 的变更日志条目移至“未发布”，由 @leaanthony 提交

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5938)中更新文档侧边栏自动生成和博客作者类型推导，由 @leaanthony 提交

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5952)中为 WebView2 初始化使用截止时间和消息泵，由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5951)中，使 WebView2 Cookie 测试在 CI 中默认跳过（除非明确启用），并将执行锁定到当前操作系统线程，由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5944)中，使 Windows 菜单构建器恢复子菜单父项的命令 ID，由 @gilad-ch 提交
- 使官方交叉编译镜像与 GTK 4.14+ Linux 支持基线保持一致（#5928）
- 在[PR](https://github.com/wailsapp/wails/pull/5915)中配置 iOS Xcode 项目以保留继承的链接器标志并添加 -ObjC，由 @mortenolsrud 提交
- 修复大型前端中`wails3 dev`资源代理频繁建立和断开 TCP 连接的问题；该问题可能耗尽主机的临时端口，并导致无关进程因`EADDRNOTAVAIL`而失败
- 在[PR](https://github.com/wailsapp/wails/pull/5934)中按窗口将事件 JavaScript 排入队列，以实现有序分派和背压，由 @leaanthony 提交

## 移除

- 移除桌面二进制发布流水线：v3发布仅包含标签，`wails3` CLI 使用`go install`安装。在[PR](https://github.com/wailsapp/wails/pull/5946)中删除`release-v3.yml`以及分派它的 nightly 步骤，由 @leaanthony 提交

## v3.0.0-beta.7 - 2026-08-11

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5512)中新增 macOS 自动播放偏好设置，以禁用媒体播放必须由用户操作触发的要求，由 @Eyalm321 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5935)中将 mailbox 的变更日志条目移至“未发布”，由 @leaanthony 提交

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5945)中让 macOS 缩放动画使用 CADisplayLink 或 NSTimer，以获得更流畅的性能，由 @savely-krasovsky 提交

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5915)中配置 iOS Xcode 项目以保留继承的链接器标志并添加 -ObjC，由 @mortenolsrud 提交
- 修复大型前端中`wails3 dev`资源代理频繁建立和断开 TCP 连接的问题；该问题可能耗尽主机的临时端口，并导致无关进程因`EADDRNOTAVAIL`而失败
- 在[PR](https://github.com/wailsapp/wails/pull/5934)中按窗口将事件 JavaScript 排入队列，以实现有序分派和背压，由 @leaanthony 提交

### 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5851)中实现通用异步 FIFO mailbox，以有序传递事件，由 @savely-krasovsky 和 @DevLumuz 提交

## v3.0.0-beta.6 - 2026-08-09

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5930)中为超大事件实现有界的主机端存储和有序的 JavaScript 传递，由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5921)中实现 macOS Dock 图标弹跳，用于窗口闪烁提醒，由 @julianstorer 提交

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5931)中，使资源服务器在刷新期间保留内容类型嗅探错误和未写入的前缀，由 @leaanthony 提交
- 防止从 Wails 回调替换应用程序菜单时导致 macOS 应用程序崩溃
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨认的问题。深色模式 uxtheme 导出受内部版本18334限制，因此应用级深色模式选择加入从未在这些主机上运行：菜单背景被绘制为深色，但 Windows 仍使用浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号从17763起便已存在，因此现在的限制条件与之匹配。
- 修复`w32.GetStockObject`调用`GetDeviceCaps`而非`GetStockObject`的问题；该问题导致它对每个系统预定义对象都返回0。
- 在[PR](https://github.com/wailsapp/wails/pull/5924)中改进 WebView2 引导程序下载的错误处理和报告，由 @jannskiee 提交

## v3.0.0-beta.5 - 2026-08-07

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5897)中，使 macOS 应用激活仅对常规应用遵循激活策略，由 @julianstorer 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5898)中保护 Linux 构建中尚未初始化的 GTK 窗口，由 @julianstorer 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5899)中，于加载 URL 前为 Linux WebKit 窗口设置明确的不透明背景色，由 @julianstorer 提交

## v3.0.0-beta.4 - 2026-08-05

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5890)中，使 Android 构建任务默认使用 arm64，并让 deploy-emulator 选择主机架构，由 @mortenolsrud 提交

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5900)中，拖动 macOS 窗口时保留其缩放状态并减少动画效果，由 @leaanthony 提交
- 通过向`webview_window_windows_nonclient.go`构建约束添加`!server`，修复 Windows 服务器模式构建

## v3.0.0-beta.3 - 2026-08-03

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5881)中，于实现细节中记录第10阶段 Beta 验证已完成，由 @leaanthony 提交

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5877)中向 Windows 深色模式 API 传递窗口句柄并验证参数，由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5870)中集中处理 macOS 无边框窗口标题栏按钮状态的解析，由 @taliesin-ai 提交
- 防止 Windows 应用程序在 Windows 应用主题为浅色时请求深色模式而导致原生菜单文本无法辨认。在 Windows 能够渲染深色菜单文本前，菜单现在会使用与之匹配的浅色原生背景。
- 修复 Windows 10 1809 / Windows Server 2019（内部版本17763）上原生菜单无法辨认的问题。深色模式 uxtheme 导出受内部版本18334限制，因此应用级深色模式选择加入从未在这些主机上运行：菜单背景被绘制为深色，但 Windows 仍使用浅色主题绘制菜单文本，导致深色背景上显示深色文本。这些序号从17763起便已存在，因此现在的限制条件与之匹配。

## v3.0.0-beta.2 - 2026-08-02

## 变更

- 将v3从 Alpha 提升为 Beta
- 记录系统托盘的智能默认设置和弹出菜单自动隐藏行为，并为点击处理程序选择添加回归测试覆盖（#5840）。
- 在[PR](https://github.com/wailsapp/wails/pull/5861)中，使 GitHub 更新程序默认排除 Windows 安装程序资源，由 @leaanthony 提交
- 在[PR](https://github.com/wailsapp/wails/pull/5866)中添加对 macOS 无边框窗口使用圆角、直角和自定义半径圆角的支持，作者：@leaanthony

## 修复

- 报告实时 GTK4 窗口大小，并从已配置的表面发出调整大小、最大化、最小化和全屏状态事件（#5830）。
- 在[PR](https://github.com/wailsapp/wails/pull/5854)中修复 Linux WebKit 在 fetch 请求中发送 Blob 或 FormData 时崩溃的问题，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5865)中让 fetch shim 对缺失的 Blob/FormData 标头传递 undefined，作者：@leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## 新增

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5866)中添加对 macOS 无边框窗口使用圆角、直角和自定义半径圆角的支持，作者：@leaanthony

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5865)中让 fetch shim 对缺失的 Blob/FormData 标头传递 undefined，作者：@leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5857)中添加 macOS DMG 打包支持及新的选项和构建任务，作者：@leaanthony

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5861)中让 GitHub 更新程序默认排除 Windows 安装程序资产，作者：@leaanthony

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5854)中修复 Linux WebKit 在 fetch 请求中发送 Blob 或 FormData 时崩溃的问题，作者：@taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5853)中实现 macOS 标题栏双击最大化或最小化操作，作者：@taliesin-ai

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5849)中更新二维码服务教程以使用 NewServiceWithOptions，并添加间距，作者：@jeongkyu

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5856)中确保 WKWebView 在 macOS 缩放期间保持响应，作者：@leaanthony
- 修复 GTK4 窗口大小查询，并从已配置的`GdkSurface`发出调整大小、最大化、最小化和全屏状态事件。

## v3.0.0-alpha2.119 - 2026-07-27

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5833)中更新多种语言的文档，加入架构图，作者：@taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5825)中为图标生成的输入和输出提供默认路径，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5797)中将源入口模块添加到运行时 package.json 的 sideEffects 中，作者：@savely-krasovsky
- 在[PR](https://github.com/wailsapp/wails/pull/5816)中为贡献指南添加许可证和来源说明章节，作者：@taliesin-ai

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5800)中应用限定作用域的 GTK4 无边框 CSS 以移除圆角，作者：@savely-krasovsky
- 在[PR](https://github.com/wailsapp/wails/pull/5789)中妥善处理 Windows 上用于弹出菜单和屏幕枚举的光标位置获取失败，作者：@wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5678)中使 macOS 打开文件对话框正确筛选扩展名，并通过后缀验证允许的文件，作者：@phergul
- 在[PR](https://github.com/wailsapp/wails/pull/5793)中防止 Windows 深色模式初始化调用 nil API，作者：@roachadam
- 修复更新程序在32位平台上的构建失败：在`GOARCH=386`上将`maxArchiveTotalSize`常量（2 GiB）传递给`fmt.Errorf`时，该常量会溢出平台的`int`。现在已将其显式指定为`int64`类型。
- 修复启动时的 nil 指针 panic：当窗口使用深色（或系统深色）标题栏，而 Windows 构建未加载深色模式 uxtheme API 时（例如 Windows 10 1809 / Windows Server 2019，构建号17763），会发生此问题。窗口主题设置中的`AllowDarkModeForWindow`调用现在具有 nil 保护，与`w32.SetMenuTheme`中已有的保护一致。

## v3.0.0-alpha2.117 - 2026-07-08

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5462)中为 Windows 上的非客户区实现自定义命中测试逻辑，作者：@savely-krasovsky

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5761)中根据 UseVisualHosting 配置 WebView2 显示器缩放检测，作者：@wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5763)中更新常见问题文档，使其重点介绍 Wails v3的功能和指导，作者：@taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## 修复

- 修复`Menu.Update()`无法在 GTK4 Linux 上重建原生菜单的问题（#5659；@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对数量创建快照，修复显示器变化时枚举 macOS 屏幕所导致的崩溃（#5565；@x-haose 在 #5584 中独立诊断并修复）
- 修复 Windows 上的崩溃：在最小化/还原转换期间，`GetClientRect`返回 nil，而`WM_ERASEBKGND`绘制纯色背景时会发生此问题（保护措施由 @sinspired 在 #5636 中报告）
- 修复前端绑定错误始终被解析为文本的问题，作者：@mbaklor，见 #5690
- 修复使用`server`构建标签时 Windows 构建失败的问题；原因是 Windows GUI 文件缺少 macOS 和 Linux 对应文件已有的`!server`构建约束（#5680）

## v3.0.0-alpha2.114 - 2026-07-05

## 新增

- 在[PR](https://github.com/wailsapp/wails/pull/5720)中实现更新清单协议和端点提供程序，作者：@taliesin-ai

## 变更

- 在[PR](https://github.com/wailsapp/wails/pull/5711)中将`webview2`绑定作为`v3/internal/webview2`并入v3模块，移除独立模块、其 nightly 发布/同步工作流以及 go.mod 版本调整流程（v3是其唯一使用者），作者：@taliesin-ai

## 修复

- 在[PR](https://github.com/wailsapp/wails/pull/5750)中将 WebView2 显示器缩放检测和 DPI 变化时的宿主重新同步修复移至“未发布”章节，作者：@taliesin-ai
- 在[PR](https://github.com/wailsapp/wails/pull/5741)中更新 WebView2 对 float64 和 BOOL 参数的 COM 封送处理，作者：@wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5703)中防止 Windows 系统托盘图标更新和销毁时发生 panic 和 nil 解引用，作者：@wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5743)中修复隐藏窗口在 Windows 上无法正确再次隐藏的问题，作者：@wayneforrest
- 在[PR](https://github.com/wailsapp/wails/pull/5742)中使 WebView2 控制器的可见性与窗口的最小化、最大化和还原保持同步，作者：@wayneforrest

### 修复

- 重新启用 WebView2 显示器缩放检测，并仅在 DPI 变化时执行宿主重新同步，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5734)中完成；该更改基于 @randalmurphal 验证的修复方案，由 @eleclin 验证根本原因，并由 @qq540491950 完成硬件测试

## v3.0.0-alpha2.113 - 2026-07-04

## 新增

- 构建未设置`ANDROID_KEYSTORE_FILE`的发布版 AAB 时添加警告（Google Play 会拒绝使用调试签名的 Bundle），并补充 App Bundle 打包和签名文档，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中完成
- 添加多语言版“为何选择 Wails”文档，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5739)中完成
- 支持在绑定中将 Go time.Time 映射为 JS Date 或字符串，由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5398)中完成
- 添加用于提交到 Play Store 的 Android App Bundle (AAB) 打包任务（`bundle`、`bundle:fat`、`assemble:aab`、`assemble:aab:release`）；APK 任务仍用于本地/模拟器测试，由 @mortenolsrud 在[PR](https://github.com/wailsapp/wails/pull/5728)中完成（修复[#5726](https://github.com/wailsapp/wails/issues/5726)）
- 添加 Android 实体设备任务目标，并恢复相机/位置权限，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5735)中完成

## 变更

- 将`webview2`升级到v1.0.28（[发行说明](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)）。
- 将 Android 模板的`compileSdk`/`targetSdk`从34升级到35，这是 Google Play 对新应用提交的要求，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中完成

## 修复

- 修复 sponsorkit 中已烘焙的头像蒙版，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5745)中完成
- 修复因按字典序排序版本而导致 Android AVD 自动创建时选择错误系统镜像或 cmdline-tools 版本的问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中完成
- 修复设置向导推荐过时 Android NDK 版本的问题（现为26.3.11579264，与文档要求一致），由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5730)中完成
- 更新有关 SvelteKit 和选项的法语文档，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5744)中完成
- 修复显示配置变化期间枚举 macOS 屏幕时发生的 SIGSEGV，由 @flofreud 在[PR](https://github.com/wailsapp/wails/pull/5516)中完成

## v3.0.0-alpha2.112 - 2026-07-03

## 新增

- 添加基于 Go 的贡献者 SVG 生成器，并更新 docs/website 致谢页面，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5724)中完成
- 支持在绑定中将 Go time.Time 映射为 JS Date 或字符串，由 @fbbdev 在[PR](https://github.com/wailsapp/wails/pull/5398)中完成

## 变更

- 使用 Go 生成器替换基于 Node 的赞助者图片处理流水线，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5719)中完成

## 修复

- 修复 Android 构建资源依赖项安装脚本，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5729)中完成
- 在`ValidateAndSanitizeURL`中拒绝 U+0085（NEXT LINE）控制字符，从而补全 URL 验证器对空白字符的覆盖
- 在无边框窗口的 DPI 变化时重新计算 DWM 边框，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4785)中完成
- 修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题，由 @yulesxoxo 在[PR](https://github.com/wailsapp/wails/pull/4632)中完成
- 为 Darwin 对话框、菜单、托盘和通知中的 Cocoa 对象添加显式 Objective-C 内存管理，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5714)中完成
- 修复 Linux CGO 后端错误和系统托盘问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5718)中完成

## v3.0.0-alpha2.111 - 2026-07-01

## 新增

- 将 HappyTools 添加到社区展示，由 @Aliuyanfeng 在[PR](https://github.com/wailsapp/wails/pull/5061)中完成
- 添加印度尼西亚语区域设置支持和全面文档，由 @triadmoko 在[PR](https://github.com/wailsapp/wails/pull/5643)中完成
- 为 WindowsWindow 添加 DisableMenu 选项，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4813)中完成

## 变更

- 更新 Taskfile 模板和 CLI，以使用 GOOS 和 ARCH 分派构建/打包任务，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5617)中完成

## 修复

- 修复 Mac 窗口标签页问题，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5708)中完成

## 移除

- 从贡献、功能和指南文档中移除德语翻译的 MDX 文件，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5702)中完成

## v3.0.0-alpha2.110 - 2026-06-30

## 新增

- 实现 macOS WebView 重新加载和强制重新加载，并添加 WebContent 进程终止后的恢复机制，由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5129)中完成
- 为贡献、功能和指南添加全面的德语文档，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5396)中完成
- 增强通知功能，新增声音、附件和调度支持，并更新 API，由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/5333)中完成

## 修复

- 在无边框窗口的 DPI 变化时重新计算 DWM 边框，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4785)中完成
- 修复 Windows 上缩放比例不为100% 时 DnD 放置区检测失败的问题，由 @yulesxoxo 在[PR](https://github.com/wailsapp/wails/pull/4632)中完成

## v3.0.0-alpha2.109 - 2026-06-29

## 新增

- 为 EventsEmit 文档添加代码示例，由 @iamhabbeboy 在[PR](https://github.com/wailsapp/wails/pull/5026)中完成
- 添加 Windows WebView2 可视化托管选项，由 @MerIijn 在[PR](https://github.com/wailsapp/wails/pull/5380)中完成
- 将 Klustr 添加到社区展示文档，由 @SametKUM 在[PR](https://github.com/wailsapp/wails/pull/5536)中完成
- 将 Kira 添加到社区展示，并新增相关页面和变更日志条目，由 @thiennguyen93 在[PR](https://github.com/wailsapp/wails/pull/5685)中完成
- 在 MCP 服务指南中添加反馈部分，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5694)中完成

## 变更

- 服务器模式现在拥有与桌面构建任务一致的一等生产构建支持（#5693）。`task build:server`默认构建生产二进制文件（`-tags server,production`、`-trimpath`、已剥离符号），并接受`DEV=true`（开发服务器）、`OBFUSCATED=true`（garble）和`EXTRA_TAGS`。`task run:server`运行开发服务器。`Dockerfile.server` / `task build:docker`先构建生产服务器（`-tags server,production`）和生产前端；该镜像默认使用基于 distroless/static 的纯 Go 静态构建，并将`CGO_ENABLED`、`GO_IMAGE`和`RUNTIME_IMAGE`公开为可覆盖的构建参数，以支持 CGO 应用。

## 修复

- 防止关闭仍有待处理异步调用的窗口时崩溃，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4435)中完成
- 防止在 Windows 上打开隐藏应用时激活窗口，由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5249)中完成
- 确保 WebKit 请求元数据、响应完成处理和正文流处理在 GTK 主线程上运行，由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5668)中完成
- 修复`Menu.Update()`未在 GTK4 Linux 上重新构建原生菜单的问题（#5659，@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对屏幕数量进行快照，修复显示配置变更时枚举 macOS 屏幕发生崩溃的问题（#5565，@x-haose 在 #5584 中独立诊断并修复）
- 通过在`WM_DPICHANGED`处理程序中重新设置控制器边界（与取消最小化时的 DPI 重新同步方式一致），修复在 Windows 上将窗口拖过不同 DPI 的显示器后 WebView2 内容先缩小再消失的问题（#5677）

## v3.0.0-alpha2.108 - 2026-06-28

## 新增

- 通过`app.GlobalShortcut`新增全局（系统级）键盘快捷键（`Register`、`Unregister`、`UnregisterAll`、`IsRegistered`、`GetAll`）。即使应用程序未获得焦点，也会触发快捷键。各平台均采用原生实现且不依赖第三方组件：macOS 使用 Carbon 热键，Windows 使用`RegisterHotKey`，X11 使用`XGrabKey`，Wayland 使用 XDG Desktop Portal 全局快捷键接口。
- 新增内置 MCP 服务器：使用`mcp`标签构建应用程序时会自动启动一个 Model Context Protocol 服务器，使 LLM 智能体能够测试并控制正在运行的 Wails 应用程序，包括控制窗口、检查 DOM、执行 JavaScript、调用绑定方法、处理事件，以及通过屏幕上的动画光标呈现模拟的鼠标和键盘输入。无需编写用户代码：设置`WAILS_MCP=1`后，`wails3 build`/`wails3 dev`会自动添加`mcp`标签。完全通过环境变量（`WAILS_MCP_HOST`、`WAILS_MCP_PORT`、`WAILS_MCP_TIMEOUT`、`WAILS_MCP_HIDE_CURSOR`）进行配置。

## 修复

- 修复`Menu.Update()`未在 GTK4 Linux 上重新构建原生菜单的问题（#5659，@puneetdixit200 在 #5539 中独立诊断并修复）
- 通过复制屏幕 ID/名称字符串并对屏幕数量进行快照，修复显示配置变更时枚举 macOS 屏幕发生崩溃的问题（#5565，@x-haose 在 #5584 中独立诊断并修复）

## v3.0.0-alpha2.107 - 2026-06-27

## 新增

- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/5613)中新增带侧边栏导航的实验性 Wake 文档

## v3.0.0-alpha2.106 - 2026-06-24

## 变更

- 将`webview2`升级到v1.0.27。
  - ci(webview2)：修复发布构建（交叉编译 Windows + 完整的 go.sum）（#5671）\

  **完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5672)中从 WebView2 发布工作流的交叉编译流程中移除 go vet
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5670)中将自动生成变更日志所用的 OpenRouter 模型更新为 google/gemini-2.5-flash-lite
- 将`webview2`升级到v1.0.26。

### 修复

- **从暂时性运行时 COM 错误中恢复，而不是退出**（#5658、#5580）。此前，`Chromium.errorCallback`会针对*任何* COM 错误调用`os.Exit(1)`，因此启动后一次可恢复的短暂故障就会终止整个应用程序。运行时路径（`Resize`/`GetClientRect`、`Navigate`/`NavigateToString`、`Init`、`MessageReceived`、`PutZoomFactor`、`OpenDevToolsWindow`）现在会记录错误并恢复。特别是，`MessageReceived`中格式错误或不受信任的 Web 消息现在会被丢弃，而不会导致进程终止。这解决了跨越不同 DPI 显示器时发生的一类崩溃（#5544、#5650）。环境/控制器创建路径仍会导致致命错误。\

**完整差异：** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## 修复

- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5671)中修复 release-webview2 工作流，使其能正确处理 go.sum 文件
- 由 @taliesin-ai 在[PR](https://github.com/wailsapp/wails/pull/5659)中通过清除并重新构建原生菜单，修复 Linux GTK4 菜单更新问题

## v3.0.0-alpha2.105 - 2026-06-21

## 新增

- 新增`application.System`，用于从共享代码中检测运行时平台：`System.IsMobile()`（iOS/Android）、`System.IsDesktop()`（macOS/Windows/Linux）、`System.IsServer()`（`server`构建标签），以及用于直接测试单个目标的`System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)`。它可在所有目标上编译，因此无需构建标签即可进行条件分支。`@wailsio/runtime`中还提供了对应的前端辅助函数（`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`）
- 新增“使用其他前端框架”指南，介绍如何将自己的 Vite 项目放入`frontend/`（涵盖 Solid、Preact、Lit、SvelteKit、Qwik、Angular 等）
- `wails3 setup`向导现在会检查移动端（iOS/Android）工具链，包括 Xcode 和 iOS Simulator 运行时、JDK、Android SDK/NDK 及模拟器；在适用情况下，还提供一键安装和可复制的 shell 配置修复方案
- 生成的项目会附带一个`frontend/.npmrc`，其中设置了7天的`minimum-release-age`，以降低刚发布的软件包（可能已遭入侵）带来的风险（pnpm 和 bun 会遵循此设置；npm 会安全地忽略它）

## 变更

- 使用全新的霓虹山景主视觉重新设计所有内置初始模板（Web、iOS 和 Android）
- **TypeScript 现在是初始模板的默认语言，并占用不带语言限定词的模板名称。** `wails3 init`（不含`-t`）会搭建 TypeScript 项目；`-t vanilla`、`-t react`、`-t vue`和`-t svelte`均为 TypeScript 模板，对应的 JavaScript 变体为`-t vanilla-js`、`-t react-js`、`-t vue-js`和`-t svelte-js`。内置模板通过`template.yaml`中的`typescript:`声明其语言；使用`-ts`后缀的社区模板仍可作为回退方案继续工作
- 使用霓虹“digital Wails”主题（山景背景上的磨砂玻璃鲜明效果）重新设计`wails3 setup`向导

## 修复

- 修复 Windows 上恢复应用程序时发生的崩溃：该应用程序因最小化时间过长，导致 WebView2 已挂起或其渲染/GPU 进程已被回收。最小化/恢复时的 DPI 重新同步（#5544）现在仅在窗口 DPI 实际发生变化时才会操作 WebView2 控制器，从而避免在常见的相同 DPI 恢复场景中对已挂起的控制器进行致命的 COM 调用（#5605）
- 修复长期运行且频繁加载资源/媒体的 Linux 应用程序中反复发生的原生`SIGABRT`/`SIGSEGV`崩溃（通常发生在 GTK 主循环期间的`g_object_unref`内部）。资源服务器从工作 goroutine 中完成`WebKitURISchemeRequest`，因而在 GTK 主线程之外调用了线程不安全的 WebKit2GTK 函数；现在，完成操作（`webkit_uri_scheme_request_finish_with_response`/`finish_error`）会在主线程上运行。此项工作补全了 #5566中的部分修复。GTK3 和 GTK4/WebKitGTK 的6.0构建均受影响（#5631、#5557）
- 修复 Linux/GTK3 上`setupSignalHandlers`中间歇出现的`fatal error: invalid pointer found on stack`。作为信号`user_data`传递的窗口 ID 保存在 Go 的`unsafe.Pointer`局部变量中，因此垃圾回收器在堆栈复制期间扫描到这个（非指针）值时会中止。现在 Go 侧会将该 ID 保持为整数类型（`uintptr_t`），把 #4958对 GTK4 路径应用的同一修复向后移植到旧版 GTK3 路径（该修复将 C 信号函数切换到`uintptr_t`，以消除`-race`/checkptr 错误）（#5631）

## 移除

- 移除`react-swc`、`preact`、`lit`、`solid`、`qwik`和`sveltekit`初始模板（及其`-ts`变体）。现在支持的内置模板集为`vanilla`、`react`、`vue`和`svelte`——每个模板默认均使用 TypeScript，并提供`-js` JavaScript 变体。仍可通过[自带前端](https://v3.wails.io/guides/dev/frontend-frameworks)或自定义模板使用任何其他框架

## v3.0.0-alpha2.104 - 2026-06-18

## 修复

- 修复绑定的 Go 服务方法返回空字符串时发生的 iOS 崩溃（SIGABRT）。iOS 资源响应写入器错误地使用`buf != nil`而非正文长度来保护正文指针，因此长度为零的正文会导致`&buf[0]`发生 panic；现在改为按长度进行保护，与桌面端写入器一致

## v3.0.0-alpha2.103 - 2026-06-15

## 变更

- 将 iOS 和 Android 原生功能迁移到平台管理器：通过`application.IOS.*`和`application.Android.*`调用它们（例如`application.IOS.Haptic("medium")`、`application.Android.Share(payload)`），不再使用旧的`application.IOS*`/`application.Android*`自由函数（#5602）
- 重命名移动端桥接事件：跨平台事件现在使用`common:*`前缀（例如`common:haptic`、`common:location`），平台专属事件使用`ios:*`/`android:*`（例如`ios:backgroundTask`、`android:foregroundService`）；不再使用`native:*`前缀（#5602）

## v3.0.0-alpha.102 - 2026-06-14

## 新增

- 新增实验性`wails3 setup`向导，用于交互式项目设置和依赖项检查
- 为`wails3 doctor`新增`--json`标志，以输出机器可读的结果
- 在`wails3 doctor`命令中新增签名状态部分

## 修复

- 修复 Linux 上的 npm 检测，使其除检查包管理器外还会检查 PATH

## v3.0.0-alpha.101 - 2026-06-13

## 新增

- iOS：原生消息对话框（UIAlertController）以及打开文件/多个文件/目录对话框（UIDocumentPickerViewController）；保存对话框会返回明确的错误
- iOS：通过 UIPasteboard 支持剪贴板
- iOS：通过 UIScreen 获取真实的屏幕指标（点、像素、缩放比例和安全区域工作区）
- iOS：支持设备构建（`IOS_PLATFORM=device`）、代码签名身份/预置描述文件/权利配置、`.ipa`打包，以及通过 devicectl 执行`deploy-device`
- iOS：可配置最低 iOS 版本（build/config.yml 中的`ios.minIOSVersion`）
- iOS：`wails3 doctor`会报告 macOS 上 Xcode 和 iOS SDK 的可用性
- iOS：系统事件——电池、网络、主题、屏幕锁定和内存不足事件会作为`events.IOS.*`及平台中立的`events.Common.*`应用程序事件公开
- iOS：原生移动功能桥接（导出的`application.IOS*`）——共享面板、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（Face ID/Touch ID）、本地通知和 Keychain 安全存储
- iOS：传感器和硬件——触觉反馈、单次地理位置获取、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距和屏幕捕获检测
- iOS：文档（IOS.md 和文档网站指南）
- Android：原生消息对话框（AlertDialog）以及打开文件/多个文件对话框（Storage Access Framework，以缓存副本形式导入）；打开目录和保存对话框会返回明确的错误
- Android：通过 ClipboardManager 支持剪贴板
- Android：通过 WindowMetrics/DisplayMetrics 获取真实的屏幕指标（dp、像素、缩放比例和系统栏工作区）
- Android：触觉反馈（`Android.Haptics.Vibrate`）、设备信息（`Android.Device.Info`）和 Toast（`Android.Toast.Show`）运行时方法
- Android：类型化生命周期事件（`events.Android.*`，由 events.txt 生成），并将`ActivityCreated`映射到`Common.ApplicationStarted`
- Android：构建流水线会生成可安装的调试版和发布版 APK（`android:run`、`android:package`、`android:package:fat`）；发布签名默认使用调试密钥库，也可通过`ANDROID_KEYSTORE_*`环境变量使用真实密钥库
- Android：`wails3 doctor`会报告 Android SDK、NDK 和 JDK
- Android：系统事件——电池、网络、主题、屏幕锁定和内存不足事件会作为`events.Android.*`及平台中立的`events.Common.*`应用程序事件公开
- Android：原生移动功能桥接（导出的`application.Android*`）——共享、打开 URL、保持唤醒、手电筒、安全区域边距、亮度、应用信息、方向锁定、状态栏、生物识别（BiometricPrompt）、本地通知和 EncryptedSharedPreferences 安全存储
- Android：传感器和硬件——触觉反馈、单次地理位置获取、加速度计、接近传感器、文本转语音、存储信息、电源/电池状态、网络状态、键盘边距和使用 FLAG_SECURE 阻止屏幕捕获
- Android：文档（ANDROID.md 和文档网站指南）
- 示例：`mobile`综合示例新增 Mobile 和 Hardware 选项卡，用于演示 iOS 和 Android 上的原生功能桥接（胶囊式选项卡可换行显示为多行）
- 移动端：电池——应用进入后台时，加速度计、接近传感器、手电筒和示例中的周期性时钟会暂停，并在应用返回前台时恢复（Android 会让进程继续在后台运行，而手电筒属于在 iOS 上会持续保持的硬件状态）；此外，Android 系统事件接收器仅在应用位于前台时注册
- iOS：相机捕获——`application.IOSCapturePhoto`/`IOSCaptureVideo`（UIImagePickerController → 携带 base64 缩略图的`native:capture`事件）
- iOS：后台执行——`application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask`（UIApplication 后台任务时间窗口），以及可配置的`ios.backgroundModes`（build/config.yml），它会将`UIBackgroundModes`模板化写入生成的 Info.plist
- Android：相机捕获——`application.AndroidCapturePhoto`/`AndroidCaptureVideo`（通过 FileProvider 调用系统相机 → `native:capture`事件）
- Android：前台服务——`application.AndroidStartForegroundService`/`AndroidStopForegroundService`（带有持续通知的`WailsForegroundService`会让进程保持运行，以执行长时间运行的后台工作）
- 示例：新增 Camera 选项卡，用于演示照片/视频捕获和后台执行（Android 上使用前台服务，iOS 上使用后台任务时间窗口）

## 修复

- 修复`getUserMedia`在 Linux 上总是因`NotAllowedError`而失败的问题：WebKitGTK 会拒绝无人处理的权限请求，而`permission-request`信号此前未连接。现在根据新的跨平台`WebviewWindowOptions.Permissions`映射（`map[PermissionType]Permission`）处理摄像头/麦克风权限，Linux（WebKitGTK）和 Windows（WebView2）都会遵循该映射。在没有原生提示框的 Linux 上，摄像头/麦克风默认允许访问（恢复`getUserMedia`），并可通过`PermissionDeny`关闭（#5552）
- iOS：`GOOS=ios`现已能再次编译（导出`events.IOS`并添加移动端方法名存根），带 production 标签的构建也可编译（修复 pkg/application 和多个服务中的构建标签）
- iOS：Go→JS 事件和 ExecJS 现已正常工作——页面启动时不再加载两次，`wails:runtime:ready`握手也不再可能丢失
- iOS：`ApplicationDidFinishLaunching`/`ApplicationStarted`不再与应用启动发生竞态；移除了固定2秒的启动等待
- iOS：修复每次执行 Go→JS JavaScript 时发生的 C 字符串泄漏
- iOS：`hasListeners`现在会反映真实的监听器注册状态
- iOS：生产构建现已在编译时移除框架调试日志
- Android：`GOOS=android`现已能再次编译——定义了`events.Android`，移除了越界的`events_android.go`监听器数组，新增移动端方法名存根，并阻止桌面 Linux 文件（`linux_cgo.*`、`events_linux.*`、`environment_linux.go`）混入 Android 构建
- Android：JS→Go 绑定现已正常工作——WebView 无法将`fetch()` POST 正文传递给`shouldInterceptRequest`，因此运行时调用现在通过 JavascriptInterface 传输（`nativeHandleRuntimeCall`）进行路由，不再因请求正文为 nil 而崩溃
- Android：`Screens.*`运行时调用现在会返回真实数据——ScreenManager 现已在启动时填充（此前从未接入，因此`GetAll`返回 nil）
- Android：生产构建会编译掉框架调试日志；调试构建则通过 logcat 输出，并使用 `Wails` 标签
- Android：实现真正的 `hasListeners` 注册表、JNI 引用和异常处理，以及仅加载一次的页面生命周期（不会重复导航）
- 修复 Vite 开发服务器运行时，`wails3 generate bindings` 在 Windows 上因“访问被拒绝”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5515）
- 修复 macOS 上显示器变更后读取屏幕信息时偶发的致命崩溃：屏幕 ID 和名称存储的是指向自动释放的 `UTF8String` 缓冲区的指针，这些缓冲区可能在 Go 复制之前被释放（释放后使用）。现在会对字符串执行 `strdup`，并在转换后释放；屏幕枚举也会在显式的自动释放池中运行，因此从 Go goroutine 调用时不再泄漏（#5556）
- 修复 assetserver 关闭 `WebKitURISchemeRequest` 时 Linux 上偶发的 SIGSEGV：最终的 `g_object_unref` 在 assetserver goroutine 上运行，导致 WebKit GObject 在 GTK 主线程之外被终结。现在通过 `g_main_context_invoke` 将 unref 调度到 GTK 主上下文中执行（#5557）

## v3.0.0-alpha.100 - 2026-06-13

## 新增

- 为 `MacWebviewPreferences` 增加更多 WKWebView 配置选项：`EnableAutoplayWithoutUserAction`、`AllowsAirPlayForMediaPlayback`、`AllowsMagnification`、`JavaScriptCanOpenWindowsAutomatically`、`MinimumFontSize` 和 `ApplicationNameForUserAgent`（#5549）

## 修复

- 修复 Vite 开发服务器运行时，`wails3 generate bindings` 在 Windows 上因“访问被拒绝”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5561）
- 修复 Linux 上无边框窗口不触发 JS 调整大小事件的问题；修复无边框窗口的滚动条边缘检测（#5368）
- 修复 Windows 上临时目录与安装目录位于不同卷时，更新程序因“无效的跨设备链接”而失败的问题（#5560）

## v3.0.0-alpha.99 - 2026-06-10

## 修复

- 修复 Vite 开发服务器运行时，`wails3 generate bindings` 在 Windows 上因“访问被拒绝”而失败的问题：现在会将生成的文件同步到输出目录，而不是通过重命名覆盖该目录（#5515）

## v3.0.0-alpha.98 - 2026-06-03

## 修复

- 修复 Linux 上 WebKit 空闲时（例如检查器处于打开状态）界面冻结的问题：不再对 `SIGUSR1` 强制设置 `SA_ONSTACK`，因为这会破坏 JavaScriptCore 的 GC 线程同步（#5527）

## v3.0.0-alpha.97 - 2026-05-31

## 新增

- 添加调试页面以及使用 `runtime/trace` 的说明

## 变更

- 移除一些不必要的 `_ "embed"` 导入，稍微整理代码

## 修复

- 修复 Windows 上窗口取消最大化后未强制执行最小宽度和高度约束的问题（#4593）
- 修复同时启用 Frameless 和 Transparent 窗口选项时，全屏模式下鼠标点击会穿透的问题（#4408）

## v3.0.0-alpha.96 - 2026-05-25

## 新增

- 添加 Garble 混淆支持（[#4563](https://github.com/wailsapp/wails/issues/4563)）：稳定的绑定方法 ID、构建/Taskfile 配套支持（`build --obfuscated --garbleargs`、`generate bindings -obfuscated`），并为所有面向运行时的载荷添加 JSON 结构体标签（`EnvironmentInfo`、`OSInfo`、`Screen`、`Rect`、`Point`、`Size`、`Capabilities`），使线路格式在 Garble 重命名导出字段后仍保持不变。

## v3.0.0-alpha.95 - 2026-05-20

## 新增

- 补充缺失的项目结构页面

## 变更

- 文档：将架构页面上的几个图表改为使用时序图，使显示更清晰
- 文档：添加说明，将安装 D2 列为运行前提条件

## 修复

- 修复 GTK4 默认配置下的 `wails3 generate appimage`：打包器现在会先从二进制文件检测 GTK 技术栈，再搜索运行时文件，因此 GTK4 构建会选择 `libwebkitgtkinjectedbundle.so`（位于 `webkitgtk-6.0/` 下），`-tags gtk3` 构建则会选择 `libwebkit2gtkinjectedbundle.so`（位于 `webkit2gtk-4.1/` 下）。`.relr.dyn` 探测现在也会检查 `libgtk-4.so.1`，因此无论使用哪种技术栈，现代工具链都会正确禁用剥离。（#5475）
- 修复使用相对 `-builddir` 调用 `wails3 generate appimage` 时失败的问题：打包器现在会预先将 `-binary`、`-icon`、`-desktopfile`、`-builddir` 和 `-outputdir` 解析为绝对路径，避免流程中途的 `s.CD` 破坏 AppRun 下载 goroutine 或复制后的 `ldd` 探测。
- 修复桌面文件的 `Name=` 字段与二进制文件基本名称不匹配时，`wails3 generate appimage` 无法将最终 AppImage 移动到 `-outputdir` 的问题：打包器现在通过 `OUTPUT` 环境变量强制 linuxdeploy 的 appimage 插件将 AppImage 写入 `<binary>-<arch>.AppImage`，而不是使用从桌面文件派生的名称。
- 修复在 alpha.93 中将 GTK4 + WebKitGTK 6.0 技术栈提升为默认配置后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep` 和 `Common.SystemDidWake` 在 Linux 上不触发的问题。新的默认 `application_linux.go` `run()` 未调用 `setupCommonEvents()`（该函数会将 `Linux.*` 事件转发给对应的 `Common.*` 事件）或 `monitorPowerEvents()`。现在 GTK3 和 GTK4 构建路径通过 `application_linux_dbus.go` 共享 DBus 电源监控辅助程序。（#5474）

## v3.0.0-alpha.94 - 2026-05-19

## 修复

- 修复在 alpha.93 中将 GTK4 + WebKitGTK 6.0 技术栈提升为默认配置后，`events.Common.ApplicationStarted`、`Common.ThemeChanged`、`Common.SystemWillSleep` 和 `Common.SystemDidWake` 在 Linux 上不触发的问题。新的默认 `application_linux.go` `run()` 未调用 `setupCommonEvents()`（该函数会将 `Linux.*` 事件转发给对应的 `Common.*` 事件）或 `monitorPowerEvents()`。现在 GTK3 和 GTK4 构建路径通过 `application_linux_dbus.go` 共享 DBus 电源监控辅助程序。（#5474）

## v3.0.0-alpha.93 - 2026-05-17

## 新增

- 由 @leaanthony 在 Linux 的 `wails3 doctor` 输出中添加 `XDG_SESSION_TYPE`

## 修复

- 由 @leaanthony 修复 Wayland 上因 appmenu-gtk-module 访问尚未实现的窗口而导致窗口菜单崩溃的问题（#4769）
- 由 @leaanthony 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题
- 由 @overlordtm 修复 Windows 上初始化拖放时出现“内存不足”错误的问题（#4701）
- 由 @leaanthony 修复主线程回调存储在删除映射项时错误使用 RLock 而导致的竞态条件（Linux、macOS、iOS）（#4424）
- 修复将命令行参数传递给任务时的变量处理。现在，以 KEY=VALUE 对形式指定的 CLI 变量会正确初始化，并在整个任务执行过程中传播。
- 修复 macOS 上的 NSWindowZoomButton 冲突：`MaximiseButtonState` 和 `FullscreenButtonState` 现在会在启动时和运行时同时采用限制更严格的状态；任一设置器都无法再静默覆盖另一个（#5319）
- 修复旧版 GTK3 构建路径（`-tags gtk3`）中一组原有缺陷，这些缺陷由 CodeRabbit 在 #5463 中发现：通过文件关联启动时不再跳过启动处理程序；`getTheme` 现在具备边界和类型安全性；`appName` 不再释放由 GLib 所有的内存；`clipboardGet` 不再泄漏 GTK 返回的 `gchar*`；`Calloc` 现在使用指针接收者（且 `NewCalloc` 返回 `*Calloc`），因此池能够真正跟踪内存分配；`zoomOut` 现在使用 `zoomInFactor` 的倒数，不再使用会被钳制为 1.0 的负乘数；`execJS` 现在复用预分配的空世界名称，不再每次调用都泄漏一个 `C.CString("")`；已从 `menuItem.setAccelerator` 中移除用于开发的 `fmt.Println`。解决 #5465。
- 修复默认 GTK4 构建路径（`linux_cgo.go`）中相同的 `Calloc` 值接收者泄漏问题：使用指针接收者和 `NewCalloc() *Calloc`，从而真正跟踪并释放每个窗口的 `c.String(...)` 分配。

## v3.0.0-alpha.92 - 2026-05-15

## 新增

- 修改 Taskfile，以便通过 `PACKAGE_MANAGER` 选项控制所用的前端包管理器
- 使用 `{{.Opn}}` 和 `{{.Cls}}` 扩充模板数据，使 Taskfile 模板的编写更加可预测

## 变更

- 修改了若干现有 Taskfile，以使用 `{{.Opn}} and {{.Cls}}`

## 修复

- 修复面板读取托盘菜单时更新该菜单会在 `linuxSystemTray` 中触发的 `concurrent map read and map write` 运行时致命错误。
- WebView2 的错误和堆栈跟踪输出改用 `log`，不再使用 `fmt`，以免应用在 Windows 上没有附加控制台时丢失消息。

## v3.0.0-alpha.91 - 2026-05-12

## 变更

- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5414) 中更新赞助者 SVG
- <strong>破坏性变更（macOS）：</strong>统一 macOS 坐标系，使 `GetScreens`、`Position` 和 `SetPosition` 全部使用同一空间——逻辑点、Y 轴向下，并以主屏幕左上角为 `(0,0)`。这与 Windows、GTK 以及 Electron 和 Web 的公共 API 一致。物理位置位于主屏幕上方的屏幕现在会报告负的 `Bounds.Y`（此前为正），并且 `Position()`/`SetPosition()` 值现在使用逻辑点，而不再使用 `points × primaryScale`。仍可保持 `Position()` → `SetPosition()` 往返转换不变；对于早期 alpha 构建所记录的绝对值或手动计算的规避方案（例如乘以 `primaryScale`，或依据屏幕高度翻转 Y 轴），需要进行更新。解决 [#5117](https://github.com/wailsapp/wails/issues/5117)。

## 修复

- 防御性验证 DBus 信号名称和消息体长度，以防止 panic；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5416) 中完成
- 修复 Linux 上 GTK 菜单处理中的内存安全问题；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5363) 中完成
- 检测 NVIDIA GPU，并在 Linux 上禁用 DMA-BUF 渲染器；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5295) 中完成
- 修复 macOS 上 `SetPosition` 的跨屏幕 Y 坐标转换：使用主屏幕高度作为全局基准，使窗口能正确定位在相对于主显示器存在垂直偏移的显示器上；见 [#5117](https://github.com/wailsapp/wails/issues/5117)
- 修复 Git PR 模板，使其指向正确的反馈 URL；由 @wayneforrest 在 [PR](https://github.com/wailsapp/wails/pull/5109) 中完成
- 修复一系列 Windows 系统托盘 `SetMenu` 崩溃问题。其原因是损坏的 `DestroyMenu` 系统调用传入了四个参数而非一个，导致每次调用都返回 FALSE，且没有释放任何资源。现在还会在重建菜单时释放 HMENU 和 HBITMAP 句柄（包括运行时通过 `MenuItem.SetBitmap` 分配的句柄），重置 `Win32Menu.Update` 中过时的复选框/单选按钮映射，并移除 `systemtray.updateMenu` 中重复的 `Update()` 调用，该调用曾导致分配量翻倍。长期运行的系统托盘应用不再在每次重建菜单时泄漏 GDI/USER 对象。

## v3.0.0-alpha.90 - 2026-05-11

## 新增

- 在 macOS 上为 WKWebView User-Agent 添加可配置的应用名称；由 @vinhvoit225 在 [PR](https://github.com/wailsapp/wails/pull/5261) 中完成
- 向 gin-service 示例添加间接依赖 github.com/coder/websocket；由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5400) 中完成
- 为构建资产测试添加深度相等比较支持；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5402) 中完成

## 变更

- 将构建输出整合到 assets 目录；由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5401) 中完成
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5399) 中更新赞助者 SVG

## 修复

- 使用通知对象传递 macOS 单实例消息；由 @overlordtm 在 [PR](https://github.com/wailsapp/wails/pull/5289) 中完成
- 批量处理 Windows 回调，以防止高负载下丢失 promise；由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5383) 中完成

## v3.0.0-alpha.89 - 2026-05-10

## 新增

- 添加 go<em>test</em>results 作业以汇总 Go 测试结果；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5316) 中完成

## 变更

- 在满足条件时，将大型 RPC 载荷拆分为分块 POST 请求；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5369) 中完成
- 在所有前端模板中将 Vite 从 5.x.x 升级至 8.0.0；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5386) 中完成
- 将 Vite 开发服务器端口配置迁移到环境变量；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5365) 中完成
- 将所有模板中的 Vite 开发服务器配置为绑定到 127.0.0.1；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5361) 中完成
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5384) 中更新赞助者 SVG

## 修复

- 在更新 build-assets 期间清理 Info.plist 模板存根；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5312) 中完成
- 通过在主线程上同步应用菜单项修改方法（`setMenuItemChecked()`、`setMenuItemLabel()`、`setMenuItemDisabled()`、`setMenuItemHidden()`、`setMenuItemTooltip()`），修复 macOS 菜单中的过时状态，并消除 `dispatch_async` 竞态；该竞态会导致快速重新打开菜单时仍渲染先前状态（#5002）
- 在开发模式下忽略 `*_test.go` 文件，以避免不必要的重新构建；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5203) 中完成
- 防止应用未运行时 Menu.Update() 发生段错误；由 @wucm667 在 [PR](https://github.com/wailsapp/wails/pull/5291) 中完成
- 在 Windows 上使用 lastSizeWParam 控制菜单栏重绘；由 @taliesin-ai 在 [PR](https://github.com/wailsapp/wails/pull/5382) 中完成

## v3.0.0-alpha.88 - 2026-05-09

## 变更

- 将 HiddenOnTaskbar 改为使用 WS<em>EX</em>TOOLWINDOW；由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/5371) 中完成
- 重新排列依赖项，并移除 go.mod 中的 webview2 replace 指令；由 @atterpac 在 [PR](https://github.com/wailsapp/wails/pull/5370) 中完成
- 由 `@github-actions[bot]` 在 [PR](https://github.com/wailsapp/wails/pull/5358) 中更新赞助者 SVG

## 修复

- 移除泛型间接别名并统一映射键类型；由 @fbbdev 在 [PR](https://github.com/wailsapp/wails/pull/5331) 中完成

## 移除

- 删除 PR-master 工作流，并在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5377)中移除文档、Go 测试和跳过测试

## v3.0.0-alpha.87 - 2026-05-07

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5352)中为 Wails v3添加韩语文档
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5354)中添加安装和快速入门的法语文档
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5355)中添加快速入门、概念和社区的葡萄牙语文档

## v3.0.0-alpha.86 - 2026-05-06

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5328)中添加法语文档本地化
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5343)中为文档站点添加德语区域设置

## 变更

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5347)中，将全部8个翻译区域设置注册到文档配置中
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5317)中更新与 WebView2 相关的多个 Windows 文件

## 修复

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5340)中，将 Linux 的对话框分派拆分为 GTK3 和 GTK4 两套实现
- 确保对话框回调在 GTK 线程上执行，修复段错误，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5339)

## v3.0.0-alpha.85 - 2026-05-05

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5179)中向仓库添加 PR 模板 URL
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5330)中为 Wails v3添加德语文档

## v3.0.0-alpha.84 - 2026-05-03

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5307)中添加用于禁止 Escape 键在 macOS 上退出全屏模式的选项
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5310)中添加用于禁止 Escape 键在 macOS 上退出全屏模式的选项
- 在 @yuseferi 的[PR](https://github.com/wailsapp/wails/pull/5288)中添加 Pausa 社区展示文档

## 变更

- 在`@github-actions[bot]`的[PR](https://github.com/wailsapp/wails/pull/5308)中更新赞助者 SVG
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5309)中更新图标生成命令，以处理不受支持的平台
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5224)中用三态 ButtonState 替换布尔型全屏 API，并实现平台绑定

## 修复

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5315)中添加防护，避免 WebView2 在控制器状态为 nil 时执行焦点操作
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5313)中更新 GitHub Actions 工作流，使其正确引用 PR 的基础分支
- 在开发模式下忽略`*_test.go`文件，以避免不必要的重新构建，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5203)
- 防止应用未运行时 Menu.Update() 发生段错误，见 @wucm667 的[PR](https://github.com/wailsapp/wails/pull/5291)

## v3.0.0-alpha.83 - 2026-05-02

## 新增

- 在 @symball 的[PR](https://github.com/wailsapp/wails/pull/5094)中添加 InstallScope 标志以及按计算机/用户安装的构建选项
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5294)中为 BrowserWindow 添加无操作的 SetScreen 方法，以满足 Window 接口

## 修复

- 在 Linux 上检测 NVIDIA GPU 并禁用 DMA-BUF 渲染器，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5295)
- 修复 Git PR 模板，使其指向正确的反馈 URL，见 @wayneforrest 的[PR](https://github.com/wailsapp/wails/pull/5109)
- 修复一系列 Windows 系统托盘`SetMenu`崩溃问题。这些问题由损坏的`DestroyMenu`系统调用引起：该调用传入了四个参数而非一个，导致每次调用都返回 FALSE，且未释放任何资源。此外，在重新构建菜单时释放 HMENU 和 HBITMAP 句柄（包括运行时通过`MenuItem.SetBitmap`分配的句柄），重置`Win32Menu.Update`中过期的复选框/单选框映射，并移除`systemtray.updateMenu`中导致分配量翻倍的冗余`Update()`调用。长期运行的系统托盘应用在每次重新构建菜单时不再泄漏 GDI/USER 对象。

## v3.0.0-alpha.82 - 2026-05-01

## 修复

- 修复桌面文件生成，使其能够正确处理桌面名称，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5232)

## v3.0.0-alpha.81 - 2026-04-30

## 变更

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5286)中将每夜版发布时间调整为 UTC 15:00

## 修复

- 修复 Retina Mac 上 Screen Bounds、WorkArea 和 Size 值减半的问题 -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## 变更

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5285)中更新文档依赖项和内容集合加载器

## v3.0.0-alpha.79 - 2026-04-29

## 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5270)中向 trigger-release 作业授予 actions: write 权限

## 变更

- 发布任务默认使用 master 分支，并更新变更日志措辞，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5283)
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5282)中更新 auto-changelog 工作流，使其使用最新版本
- 通过添加路径筛选器并移除失效的工作流来提高工作流效率，见 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5280)
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5274)中更新文档，使示例链接引用 master 分支
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5272)中更新v3的文档和示例

## 修复

- 在 @AkagiYui 的[PR](https://github.com/wailsapp/wails/pull/5265)中增强反向代理，为开发环境添加重试逻辑并强制使用 IPv4
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5281)中重写未发布变更日志的触发工作流

## 移除

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5267)中移除用于各种测试目的的 shell 测试脚本
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5266)中删除v3-alpha 文档部署工作流和 CNAME 记录

### 新增

- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5196)中向侧边栏导航添加“前端路由”条目
- 在 @leaanthony 的[PR](https://github.com/wailsapp/wails/pull/5185)中添加包含特定框架建议的前端路由指南
- 添加对模态表单（macOS）的支持
- 由 @leaanthony 升级 ghw 版本，以更好地支持 Apple 设备（#4977）
- 为 Dock 服务添加`GetBadge`方法
- 为`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）（#4957）
- 添加有关绑定生成器自动生成枚举的文档，包括专门的 Enums 页面和侧边栏导航（#4972）
- 为`wails3 build`命令添加`-tags`标志，用于传递自定义 Go 构建标签（例如`wails3 build -tags gtk4`）（#4957）
- 在`v3/examples/web-apis/`中添加 Web API 示例，演示41个以上的浏览器 API，包括存储（localStorage、sessionStorage、IndexedDB、Cache API）、网络（Fetch、WebSocket、XMLHttpRequest、EventSource、Beacon）、媒体（Canvas、WebGL、Web Audio、MediaDevices、MediaRecorder、Speech Synthesis）、设备（Geolocation、Clipboard、Fullscreen、Device Orientation、Vibration、Gamepad）、性能（Performance API、Mutation Observer、Intersection/Resize Observer）、UI（Web Components、Pointer Events、Selection、Dialog、Drag and Drop）等
- 添加 WebView API 兼容性检查器示例（`v3/examples/webview-api-check/`），用于跨平台测试200个以上的浏览器 API
- 添加`internal/libpath`包，用于在 Linux 上通过并行搜索和缓存查找原生库路径，并支持 Flatpak/Snap/Nix
- <strong>开发中：</strong>为 Linux 添加实验性 WebKitGTK 6.0 / GTK4 支持，可通过`-tags gtk4`启用（GTK3/WebKit2GTK 4.1仍为默认选项）
- 注意：在平铺式窗口管理器（例如 Hyprland、Sway）上，由于窗口管理器控制窗口的几何布局，最小化/最大化操作可能无法按预期工作
- 由 @AbdelhadiSeddar 在**在 JavaScript 中监听事件**文档中添加了如何使用**一次性处理程序**的说明
- 由 @leaanthony 为`WebviewWindowOptions`添加`UseApplicationMenu`选项，使 Windows/Linux 上的窗口能够继承通过`app.Menu.Set()`设置的应用程序菜单
- 由 @wimaha 添加对使用`.icon`文件（Apple Icon Composer 格式）生成 Liquid Glass 图标和资源目录（macOS）的支持（#4934）
- 添加用于无头/Web 部署的实验性服务器模式（`-tags server`）。该模式允许将 Wails 应用作为 HTTP 服务器运行，无需原生 GUI 依赖项。请使用`wails3 task build:server`构建。详情请参阅`examples/server`。
- 添加`internal/libpath`包，用于在 Linux 上通过并行搜索和缓存查找原生库路径，并支持 Flatpak/Snap/Nix
- 由 @leaanthony 为`MacWindow`添加`CollectionBehavior`选项，用于控制窗口在 macOS 空间和全屏模式下的行为（#4756）
- 由 @leaanthony 为 pkg/application 添加单元测试
- 由 @leaanthony 为 MSIX 打包添加自定义协议支持
- 添加 Linux 桌面环境检测[PR #4797](https://github.com/wailsapp/wails/pull/4797)
- 由 @leaanthony 为 JavaScript 运行时添加`Window.Print()`方法，用于从前端触发打印对话框（#4290）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 由 @leaanthony 为 Linux 添加更多 WebKit2 加载状态变更事件：`WindowLoadStarted`、`WindowLoadRedirected`、`WindowLoadCommitted`、`WindowLoadFinished`（#3896）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加`XDG_SESSION_TYPE`
- 在 Linux 构建期间生成`.desktop`文件，而不再仅在打包时生成（#4575）
- 由 @leaanthony 添加 Linux 运行时依赖项文档，其中包含特定于发行版的软件包名称和 nfpm 打包示例（#4339）
- 由 @leaanthony 在 Linux 上的`wails3 doctor`输出中添加 NVIDIA 驱动程序版本信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中添加 macOS 通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @chinenual 在[PR](https://github.com/wailsapp/wails/pull/4760)中为 helloworld 模板添加 aria-label 标识符，以便 Appium 测试客户端轻松测试示例应用
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4710)中为原始消息处理程序添加来源信息
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4712)中添加 macOS 通用链接支持
- 由 @APshenkin 在[PR](https://github.com/wailsapp/wails/pull/4702)中重构绑定传输层
- 由 @fbbdev 和 @ianvs 在[#4633](https://github.com/wailsapp/wails/pull/4633)中添加类型化事件
- 添加`systray-clock`示例，演示可实时更新工具提示的无头托盘（#4653）。
- 由 @Tolfx 在 #4510 中添加适用于 Windows 的 NSIS 协议模板
- 由 @Tolfx 在 #4510 中为 build-assets 添加测试
- macOS：由 @nidib 在[#4588](https://github.com/wailsapp/wails/pull/4588)中实现将原生窗口控件显示在菜单栏中
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于隐藏/显示 Dock 中的应用图标
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4451)中添加 macOS Dock 服务，用于隐藏/显示 Dock 中的应用图标
- 由 @leaanthony 在[#4534](https://github.com/wailsapp/wails/pull/4534)中添加 macOS 原生 Liquid Glass 效果支持，使用 NSGlassEffectView（macOS 15.0 及更高版本），并回退到 NSVisualEffectView，同时提供全面的材质自定义选项
- 由 @leaanthony 在[#4500](https://github.dev/wailsapp/wails/pull/4500)中实现浏览器 URL 净化。该功能基于 @APShenkin 在[#4484](https://github.com/wailsapp/wails/pull/4484)中的工作。
- 由[@leaanthony](https://github.com/leaanthony)基于[@Taiterbase](https://github.com/Taiterbase)在此[PR](https://github.com/wailsapp/wails/pull/4241)中的原始工作，为 Windows/Mac 添加内容保护
- 由 @leaanthony 在[PR](https://github.com/wailsapp/wails/pull/4488)中添加支持，可通过`wails3 build`和`wails3 package`别名将 CLI 变量传递给 Task 命令（#4422）
- 由[@atterpac](https://github.com/atterpac)在[#4318](https://github.com/wailsapp/wails/pull/4318)中添加放置区支持，事件会提供被放置元素的数据
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中为`WindowsWindow`选项添加`AdditionalLaunchArgs`，以便向 WebView2 浏览器传递额外的命令行参数。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中添加在 wails init 后自动运行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中添加 Windows Snap Assist 功能
- 在[PR](https://github.com/wailsapp/wails/pull/4467)中为`WindowsWindow`选项添加`AdditionalLaunchArgs`，以便向 WebView2 浏览器传递额外的命令行参数。
- 由[@triadmoko](https://github.com/triadmoko)在[PR](https://github.com/wailsapp/wails/pull/4286)中添加在 wails init 后自动运行 go mod tidy 的功能
- 由 @leaanthony 在[PR](https://github.dev/wailsapp/wails/pull/4463)中添加 Windows Snap Assist 功能
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中添加 Windows `getAccentColor`实现
- 由[@almas-x](https://github.com/almas-x)在[PR](https://github.com/wailsapp/wails/pull/4427)中添加 Windows `getAccentColor`实现
- Windows 深色主题菜单和菜单栏。由 @leaanthony 在 [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5) 中实现
- 重命名内置服务，使 JS/TS 绑定更清晰，由 @popaprozac 在 [PR](https://github.com/wailsapp/wails/pull/4405) 中实现
- 使用 `app.Env.GetAccentColor` 获取用户系统的强调色。适用于 macOS。由 [@etesam913](https://github.com/etesam913) 实现
- 添加 `window.ToggleFrameless()` API，由 [@atterpac](https://github.com/atterpac) 在 [#4137](https://github.com/wailsapp/wails/pull/4137) 中实现
- 为 Linux 添加发行版特定的构建依赖项，由 @leaanthony 在 [PR](https://github.com/wailsapp/wails/pull/4345) 中实现
- 添加绑定指南，由 @atterpac 在 [PR](https://github.com/wailsapp/wails/pull/4404) 中实现
- **整理测试基础设施**：将 Docker 测试文件移至专用的 `test/docker/` 目录，优化镜像并提高构建可靠性，由 [@leaanthony](https://github.com/leaanthony) 在 [#4359](https://github.com/wailsapp/wails/pull/4359) 中实现
- **改进资源管理模式**：在示例中添加正确的事件处理程序清理和可感知上下文的 goroutine 管理，由 [@leaanthony](https://github.com/leaanthony) 在 [#4359](https://github.com/wailsapp/wails/pull/4359) 中实现
- 支持构建 aarch64 AppImage，由 [@AkshayKalose](https://github.com/AkshayKalose) 在 [#3981](https://github.com/wailsapp/wails/pull/3981) 中实现
- 为 `wails doctor` 添加诊断部分，由 [@leaanthony](https://github.com/leaanthony) 实现
- 调用服务方法时将窗口添加到上下文中，由 [@leaanthony](https://github.com/leaanthony) 实现
- 添加 `window-call` 示例，演示如何确定哪个窗口正在调用服务，由 [@leaanthony](https://github.com/leaanthony) 实现
- 新增菜单指南，由 [@leaanthony](https://github.com/leaanthony) 编写
- 改进 panic 处理，由 [@leaanthony](https://github.com/leaanthony) 实现
- 新增菜单指南，由 [@leaanthony](https://github.com/leaanthony) 编写
- 为 Service API 添加文档注释，由 [@fbbdev](https://github.com/fbbdev) 在 [#4024](https://github.com/wailsapp/wails/pull/4024) 中实现
- 添加函数 `application.NewServiceWithOptions`，用于通过额外配置初始化服务，由 [@leaanthony](https://github.com/leaanthony) 在 [#4024](https://github.com/wailsapp/wails/pull/4024) 中实现
- 改进菜单控制，由 [@FalcoG](https://github.com/FalcoG) 和 [@leaanthony](https://github.com/leaanthony) 在 [#4031](https://github.com/wailsapp/wails/pull/4031) 中实现
- 添加更多文档，由 [@leaanthony](https://github.com/leaanthony) 编写
- 支持在标准事件监听器中取消事件，由 [@leaanthony](https://github.com/leaanthony) 实现
- 支持系统托盘的 `Hide`、`Show` 和 `Destroy`，由 [@leaanthony](https://github.com/leaanthony) 实现
- 支持系统托盘的 `SetTooltip`，由 [@leaanthony](https://github.com/leaanthony) 实现。最初构想由 [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304) 提出
- 在绑定生成器针对不支持类型发出的警告中报告包路径，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 为绑定生成器添加泛型别名支持，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 为绑定生成器添加对 `omitzero` JSON 标志的支持，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 添加 `//wails:ignore` 指令，以阻止为选定的服务方法生成绑定，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 为服务和模型添加 `//wails:internal` 指令，以允许类型在 Go 中导出但不在 JS/TS 中导出，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 为绑定生成器添加对别名类型常量的支持，以便使用弱类型枚举，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 为 Go 1.24 功能添加绑定生成器测试，由 [@fbbdev](https://github.com/fbbdev) 在 [#4068](https://github.com/wailsapp/wails/pull/4068) 中实现
- 为 `OSInfo.Branding` 添加对 macOS 15“Sequoia”的支持，以改进操作系统版本检测，见 [#4065](https://github.com/wailsapp/wails/pull/4065)
- 添加 `PostShutdown` 钩子，用于在关闭流程完成后运行自定义代码，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 添加 `FatalError` 结构体，以支持在自定义错误处理程序中检测致命错误，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 统一服务启动和关闭顺序并编写相关文档，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 为应用程序启动/关闭序列以及服务启动/关闭测试添加测试框架，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 添加 `RegisterService` 方法，用于在创建应用程序后注册服务，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 在应用程序和服务选项中添加 `MarshalError` 字段，用于自定义处理绑定调用中的错误，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 添加可取消的 Promise 包装器，使取消请求能够沿 Promise 链传播，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- 支持将绑定调用的取消操作关联到 `AbortSignal`，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- WML 除了支持常规的 `wml-*` 属性外，还支持 `data-wml-*` 属性，由 [@leaanthony](https://github.com/leaanthony) 实现
- 为所有服务添加 `Configure` 方法，用于后期配置和动态重新配置，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- `fileserver` 服务在未配置时发送 503 Service Unavailable 响应，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- `kvstore` 服务在未配置时默认提供内存键值存储，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 为 `kvstore` 服务添加 `Load` 方法，以便在配置更改后从文件重新加载数据，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 为 `kvstore` 服务添加 `Clear` 方法，用于删除所有键，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 在 `log` 服务中添加 `Level` 类型，以提供 JS 端日志级别常量，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 为 `log` 服务添加 `Log` 方法，用于动态指定日志级别，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- `sqlite`服务在未配置时默认提供内存数据库，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 为`sqlite`服务添加`Close`方法，以便手动关闭数据库，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 为`sqlite`服务的查询方法添加取消支持，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 为`sqlite`服务添加带 JS 绑定的预处理语句支持，由[@fbbdev](https://github.com/fbbdev)在[#4067](https://github.com/wailsapp/wails/pull/4067)中实现
- 添加 Gin 支持，由[Lea Anthony](https://github.com/leaanthony)在[PR](https://github.com/wailsapp/wails/pull/3537)中实现，基于[@AnalogJ](https://github.com/AnalogJ)在此[PR](https://github.com/wailsapp/wails/pull/3537)中的原始工作
- 修复自动保存和密码自动保存始终启用的问题，由[@oSethoum](https://github.com/osethoum)在[#4134](https://github.com/wailsapp/wails/pull/4134)中实现
- 为窗口添加`SetMenu()`，以便为窗口设置菜单，由[@leaanthony](https://github.com/leaanthony)实现
- 添加通知支持，由[@popaprozac](https://github.com/popaprozac)在[#4098](https://github.com/wailsapp/wails/pull/4098)中实现
-  添加 macOS 文件关联支持，由[@wimaha](https://github.com/wimaha)在[#4177](https://github.com/wailsapp/wails/pull/4177)中实现
- 添加`wails3 tool version`以递增语义化版本号，由[@leaanthony](https://github.com/leaanthony)实现
- 添加 macOS 和 Windows 徽标支持，由[@popaprozac](https://github.com/popaprozac)在[#](https://github.com/wailsapp/wails/pull/4234)中实现
- 添加已注册事件和严格类型事件支持，由[@fbbdev](https://github.com/fbbdev)与[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中实现
- 添加为自定义事件注册钩子的功能，由[@fbbdev](https://github.com/fbbdev)与[@IanVS](https://github.com/IanVS)在[#4161](https://github.com/wailsapp/wails/pull/4161)中实现
- 添加`app.OpenFileManager(path string, selectFile bool)`，用于在路径`path`处打开系统文件管理器，并可选择通过`selectFile`突出显示，由[@Krzysztofz01](https://github.com/Krzysztofz01)和[@rcalixte](https://github.com/rcalixte)实现
- 为`wails3 init`命令添加新的`-git`标志，由[@leaanthony](https://github.com/leaanthony)实现
- 添加新的`wails3 generate webview2bootstrapper`命令，由[@leaanthony](https://github.com/leaanthony)实现
- 在运行时中添加`init()`方法，以便手动初始化运行时，由[@leaanthony](https://github.com/leaanthony)实现
- 向窗口的 WindowOptions 添加`WindowDidMoveDebounceMS`选项，由[@leaanthony](https://github.com/leaanthony)实现
- 添加单实例功能，由[@leaanthony](https://github.com/leaanthony)实现。基于 @APshenkin 的[v2 PR](https://github.com/wailsapp/wails/pull/2951)。
- `wails3 generate template`命令，由[@leaanthony](https://github.com/leaanthony)实现
- `wails3 releasenotes`命令，由[@leaanthony](https://github.com/leaanthony)实现
- `wails3 update cli`命令，由[@leaanthony](https://github.com/leaanthony)实现
- `wails3 generate bindings`命令的`-clean`选项，由[@leaanthony](https://github.com/leaanthony)实现
- 允许为 aarch64 (arm64) 构建 AppImage Linux 应用，由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中实现
- 为赞助者添加超链接，由 @ansxuman 在[#3958](https://github.com/wailsapp/wails/pull/3958)中实现
- 支持 Linux 的 deb、rpm 和 Arch Linux 打包器构建，由
- 添加对 Darwin 通用构建和软件包的支持，由
- 向网站添加事件文档，由
- 添加为非 SSR 开发配置的 sveltekit 和 sveltekit-ts 模板
- 使用新的`wails3 update build-assets`命令更新构建资源，由
- 添加用于测试 HTML 拖放 API 的示例，由
- 添加文件关联支持，由[leaanthony](https://github.com/leaanthony)在
- 添加新的`wails3 generate runtime`命令，由
- 添加新的`InitialPosition`选项，用于指定窗口是否应居中，由
- 向`application`包添加`Path`和`Paths`方法，由
- 添加`GeneralAutofillEnabled`和`PasswordAutosaveEnabled` Windows 选项
- 添加获取调用服务方法的窗口的功能，由
- 为 WebView2 添加`EnabledFeatures`和`DisabledFeatures`选项，由
- ⊞ 添加新的 DIP 系统，以增强高 DPI 显示器支持，由
- ⊞ 添加窗口类名选项，由[windom](https://github.com/windom/)在
- 扩展服务以提供插件功能。由
- 🐧 添加 WindowDidMove / WindowDidResize 事件，见
- ⊞ 添加 WindowDidResize 事件，见
-  添加 ApplicationShouldHandleReopen 事件，以便处理 Dock
-  向实现添加 getPrimaryScreen/getScreens，由 @tmclane 在
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
- 在 Linux 上实现`setIcon`，由[@abichinger](https://github.com/abichinger)实现
- 向 dev 命令添加`-port`标志，并支持环境变量
- 为绑定方法调用添加测试，由
- ⊞ 为已创建的窗口添加`SetIgnoreMouseEvents`，由
-  添加设置窗口堆叠级别（顺序）的功能，由

### 修复

- 修复 Retina Mac 上`Screen.Bounds`、`WorkArea`和`Size`减半的问题：在`Physical*`字段中将 NSScreen 点值转换为设备像素，并填充顶层`Screen.X`/`Y`，使多显示器接触检测和工作区定位正确；由 @wayneforrest 在[PR](https://github.com/wailsapp/wails/pull/5168)中实现
- 修复 ScreenManager 中的数据竞争问题，该问题会在显示配置发生变化时（例如休眠/唤醒期间热插拔外接显示器）导致 WebKit DisplayLink 死锁
- 当 Assets.car 存在时，直接将 CFBundleIconName 设置为 appicon，见 @symball 的[PR](https://github.com/wailsapp/wails/pull/5154)
- 修复`wails3 doctor`在 Fedora、openSUSE、Arch 和 NixOS 上报告错误 WebKitGTK 软件包的问题——由于v3在编译时需要4.1 API，已移除4.0的回退条目（#5071）
- 修复 openSUSE 的 webkit2gtk doctor 软件包名称（`webkit2gtk4_1-devel` → `webkit2gtk3-devel`，后者才是正确的 openSUSE 软件包名称）（#5071）
- 修复桌面开发模式下缺少`/wails/custom.js`时出现的`Unexpected token '<'`错误。为`/wails/custom.js`添加了显式的404处理程序，并在`loadOptionalScript`中添加了不区分大小写的`Content-Type`验证，以防止将 HTML SPA 回退内容作为 JavaScript 注入。（[#5068](https://github.com/wailsapp/wails/issues/5068)）
- 修复 macOS 上系统托盘菜单的高亮状态——现在菜单打开时，图标会显示选中状态（#4910）
- 修复 macOS 上系统托盘附加窗口出现在其他窗口后方的问题——现在使用正确的弹出窗口层级（#4910）
- 修复文档中错误的`@wailsio/runtime`导入示例（#4989）
- 修复 darwin 上无边框窗口无法最小化的问题（#4294）
- 通过在 go-task 的最新状态检查中排除`node_modules/`，修复执行`wails3 build`和`wails3 dev`时卡住20-30分钟的问题。此前，`sources: "**/*"`通配模式会导致 go-task 枚举`node_modules/`中的每个文件并计算其校验和（50000-100000+ 个文件，其中包含 MUI 等大型依赖项），在 Windows/NTFS 上尤其缓慢（#4939）
- 修复 C 的`Screen` typedef 与 X11 Xlib.h 冲突所导致的 GTK4 构建失败（#4957）
- 修复 macOS 上程序坞徽章方法的一致性问题
- 修复`InvisibleTitleBarHeight`被应用于所有 macOS 窗口，而不是仅应用于无边框窗口或标题栏透明窗口的问题（#4960）
- 启用`InvisibleTitleBarHeight`时，通过跳过窗口边缘附近的拖动启动，修复从顶部两角调整窗口大小时窗口晃动或抖动的问题（#4960）
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 提交
- 修复 Windows 上文件拖放在显示缩放比例不是100% 时无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上文件拖放在使用悬停效果时无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过使用`gtk_window_present()`，修复 Linux/GTK4 上显示或隐藏窗口时有时会恢复为最小化状态的问题（#4957）
- 通过`XTranslateCoordinates`/`XMoveWindow`添加仅限 X11 的条件支持，修复 Linux/GTK4 上获取或设置窗口位置时始终返回0,0的问题（#4957）
- 添加基于信号的大小限制来替代已移除的`gtk_window_set_geometry_hints`，修复 Linux/GTK4 上未强制执行窗口最大尺寸的问题（#4957）
- 通过实现正确的 PhysicalBounds 计算，并借助`gdk_monitor_get_scale`支持小数缩放，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 提交
- 修复 Windows 上文件拖放在显示缩放比例不是100% 时无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上文件拖放在使用悬停效果时无法可靠工作的问题
- 修复 Linux 上启用文件放置后 HTML5 内部拖放失效的问题
- 通过实现正确的 PhysicalBounds 计算，并借助`gdk_monitor_get_scale`支持小数缩放，修复 Linux/GTK4 上的 DPI 缩放问题（GTK 4.14+）
- 修复 Linux/GTK4 上创建新窗口时菜单项重复的问题
- 修复 JS/TS 绑定中以枚举为键的映射类型生成问题（#4437），由 @fbbdev 提交
- 修复因 App.Window.Current() 未从主线程访问 AppKit API 而导致 macOS 上出现“幽灵窗口”的问题（#4947），由 @wimaha 提交
- 通过实现 WKUIDelegate runOpenPanelWithParameters，修复 HTML `<input type="file">`在 macOS 上无法工作的问题（#4862）
- 修复在 macOS/Linux 上使用`@wailsio/runtime` npm 模块时原生文件拖放无法工作的问题（#4953），由 @leaanthony 提交
- 修复跨软件包类型别名的绑定生成问题（#4578），由 @fbbdev 提交
- 修复 GTK 线程安全违规导致 OpenFileDialog 在 Linux 上崩溃的问题（#3683），由 @ddmoney420 提交
- 修复对隐藏或已销毁的窗口调用`Focus()`时发生 SIGSEGV 崩溃的问题（#4890），由 @ddmoney420 提交
- 修复在 Linux 上设置空图标或位图时可能发生 panic 的问题（#4923），由 @ddmoney420 提交
- 修复在 macOS 上从服务绑定调用 ErrorDialog 时崩溃的问题（#3631），由 @leaanthony 提交
- 使菜单能够在 Windows 操作系统上的`v3\examples\dialogs`中显示，由 @ndianabasi 提交
- 修复页面重新加载期间因竞态条件导致 TypeError 的问题（#4872），由 @ddmoney420 提交
- 通过移除`Collector.IsVoidAlias()`方法中的全局状态，修复绑定生成器测试输出不正确的问题（#4941），由 @fbbdev 提交
- 修复`<input type="file">`文件选择器在 macOS 上无法工作的问题（#4862），由 @leaanthony 提交
- 修复`Position()`和`SetPosition()`在 macOS 上使用不一致坐标系，导致保存和恢复状态时窗口位置发生偏移的问题（#4816），由 @leaanthony 提交
- 修复已通过应用程序清单设置 DPI 感知时 SetProcessDpiAwarenessContext 报出“Access is denied”错误的问题（#4803）
- 更新键盘快捷键文档页面，并更正`KeyBinding.Add`回调参数的类型，由 @ndianabasi 提交
- 修复有关生成自定义绑定的文档：必须使用`-d String`，而不是`-o String`
- 修复菜单在`menu.Update()`时未清除子项的问题
- 修复文档中过时的 Manager API 引用（更新31文件，改用`app.Window.New()`、`app.Event.Emit()`等新模式），由 @leaanthony 提交
- 修复 WebKit 覆盖信号处理程序，导致绑定到 JS 的 Go 方法发生 panic 时 Linux 崩溃的问题（#3965），由 @leaanthony 提交
- 修复 SaveFileDialog.SetFilename() 在 Linux 上不起作用的问题（#4841），由 @samstanier 提交
- 修复拖放示例中放置坐标显示为 undefined 的问题
- 修复 APP_NAME 包含空格时 macOS 应用程序包创建失败的问题（大括号展开问题）
- 修复在 Windows 上调用服务方法时发生索引越界 panic 的问题（回退 goccy/go-json）
- 修复 Windows 上文件拖放在显示缩放比例不是100% 时无法工作的问题
- 修复 Windows 上启用文件放置后 HTML5 内部拖放失效的问题
- 修复 Windows 上文件放置坐标使用了错误像素空间的问题（物理像素与 CSS 像素）
- 修复 Linux 上文件拖放在使用悬停效果时无法可靠工作的问题
- 修复 Linux 上启用文件拖放后 HTML5 内部拖放失效的问题
- 更新所有操作系统的 Taskfile.yml 文件中的全部命令，以支持 `APP_NAME` 等变量中包含空格，贡献者：@ndianabasi
- 修复 Linux 上执行 'build:universal:lipo:go' 任务时的命令参数错误，贡献者：@wux1an
- 修复在 Linux 上运行 'wails3 build GOOS=darwin GOARCH=arm64' 时出现的 Docker 错误“undefined symbol: **<em>ubsan</em>handle_xxxxxxx”，贡献者：@wux1an
- 整合自定义协议文档并添加 Universal Links 章节，贡献者：@leaanthony
- 通过添加防护以避免并发调用 TrackPopupMenuEx，修复在 Windows 系统托盘菜单中反复单击图标时发生的崩溃（#4151），贡献者：@leaanthony
- 防止在 app.Run() 之前调用 systray.Run() 时应用崩溃，贡献者：@leaanthony
- 修复启用 ApplicationShouldTerminateAfterLastWindowClosed 后，通过 Hide()/Show() 切换窗口可见性时 macOS 上发生崩溃的问题（#4389），贡献者：@leaanthony
- 修复在 macOS 和 Windows 上反复打开上下文菜单时的内存泄漏（#4012），贡献者：@leaanthony
- 修复 macOS 上未复用上下文菜单原生资源，导致每次显示时都重新创建菜单的问题（#4012），贡献者：@leaanthony
- 修复应用以 `Hidden: true` 启动时，单击 macOS 程序坞图标不会显示隐藏窗口的问题（#4583），贡献者：@leaanthony
- 修复因 CGO 调用中的窗口指针类型不正确而无法打开 macOS 打印对话框的问题（#4290），贡献者：@leaanthony
- 修复 appmenu-gtk-module 访问尚未实现的窗口而导致 Wayland 上窗口菜单崩溃的问题（#4769），贡献者：@leaanthony
- 修复应用名称包含无效字符（空格、括号等）时 GTK 应用崩溃的问题，贡献者：@leaanthony
- 修复在 Windows 上初始化拖放时出现“not enough memory”错误的问题（#4701），贡献者：@overlordtm
- 修复因 URI 转义不正确而导致 Linux 文件浏览器打开错误目录的问题（#4397），贡献者：@leaanthony
- 通过自动检测 `.relr.dyn` ELF 节并禁用剥离，修复现代 Linux 发行版（Arch、Fedora 39+、Ubuntu 24.04+）上的 AppImage 构建失败（#4642），贡献者：@leaanthony
- 修复 `wails doctor` 在 Fedora/DNF 系统上错误报告已安装 webkit 软件包的问题（#4457），贡献者：@leaanthony
- 修复默认 `config.yml` 会通过生产构建运行 `wails3 dev` 的问题，贡献者：@mbaklor
- 修复 iOS 服务桩因导入不存在的软件包而导致构建失败的问题，贡献者：@leaanthony
- 修复 debug/info 方法中的结构化日志记录导致“no formatting directives”错误的问题，贡献者：@leaanthony
- 移除从移动平台合并中意外包含的临时调试打印语句，贡献者：@leaanthony
- 通过自动禁用 DMA-BUF 渲染器，修复配备 NVIDIA GPU 的 Wayland 环境中 WebKitGTK 崩溃的问题（Error 71 Protocol error），贡献者：@leaanthony
- 解决 Linux 上 `application.WebviewWindowOptions.BackgroundColour` 中的 alpha 值被忽略的问题（[#4722](https://github.com/wailsapp/wails/pull/4722)，@BradHacker）
- 修复未提供自定义图标时，Windows 系统托盘图标未默认使用应用图标的问题（#4704）
- 跟踪 `HICON` 的所有权，仅销毁用户创建的句柄，防止资源管理器回收时崩溃（#4653）。
- 在销毁期间释放 Windows 系统主题监听器和保留的托盘图标，以避免 goroutine 和设备上下文泄漏（#4653）。
- 将托盘工具提示截断为 127 个 UTF-16 单元，以免损坏代理对和多字节字形（#4653）。
- 修复 Windows 打包任务失败的问题（#4667）
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 修复 go-webview2 v1.0.22 签名变更导致的 Windows 构建错误（#4513、#4645）
- 修复 Linux taskfile 中的 Linux AppImage appicon 变量 [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- 通过移除 `<.Info.Protocol>` 到 `<.Protocol>`，修复 Linux desktop.tmpl 的协议范围，贡献者：@Tolfx（#4510）
- 修复液态玻璃演示中的重复定义错误，贡献者：@Etesam913（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- 修复 Linux 上的系统托盘菜单更新，贡献者：[@JackDoan](https://github.com/JackDoan)（[#4604](https://github.com/wailsapp/wails/issues/4604)）
- 修复在 Windows 上创建隐藏窗口时出现白色窗口的问题，贡献者：@leaanthony（[#4612](https://github.com/wailsapp/wails/pull/4612)）
- 修复文档中的 notifications 软件包导入路径，贡献者：@rxliuli（[#4617](https://github.com/wailsapp/wails/pull/4617)）
- 修复使用 npm 软件包 @wailsio/runtime 时拖放无法工作的问题（#4489），贡献者：@leaanthony（#4616）
- Windows：修复启动时窗口闪烁以及隐藏窗口被错误显示的问题，贡献者：@leaanthony（[PR](https://github.com/wailsapp/wails/pull/4600)）。
- 修复 Wayland 窗口尺寸最大化问题（https://github.com/wailsapp/wails/issues/4429），贡献者：[@samstanier](https://github.com/samstanier)
- 修复 Wayland 窗口尺寸最大化问题（https://github.com/wailsapp/wails/issues/4429），贡献者：[@samstanier](https://github.com/samstanier)
- 修复液态玻璃演示中的重复定义错误，贡献者：@Etesam913（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- 修复 AssetServer 在 MacOS 上可能崩溃的问题，贡献者：@jghiloni（[#4576](https://github.com/wailsapp/wails/pull/4576)）
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修复
- 修复每夜版发布流水线，贡献者：@riadafridishibly（[#4597](https://github.com/wailsapp/wails/pull/4597)）
- 修复液态玻璃演示中的重复定义错误，贡献者：@Etesam913（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- 修复 AssetServer 在 MacOS 上可能崩溃的问题，贡献者：@jghiloni（[#4576](https://github.com/wailsapp/wails/pull/4576)）
- 修复使用 NextJs 构建时的编译问题。由 @rev42 在 [#4585](https://github.com/wailsapp/wails/pull/4585) 中修复
- 修复每夜版发布流水线，贡献者：@riadafridishibly（[#4597](https://github.com/wailsapp/wails/pull/4597)）
- 修复液态玻璃演示中的重复定义错误，贡献者：@Etesam913（[#4542](https://github.com/wailsapp/wails/pull/4542)）
- 修复 Windows 上的 SetBackgroundColour，贡献者：@PPTGamer（[PR](https://github.com/wailsapp/wails/pull/4492)）
- 更新文档以反映 Manager API 重构带来的变更，贡献者：@yulesxoxo（[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- 修复 Linux taskfile 中 Linux .desktop 文件的 appicon 变量 [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- 更新文档以反映 Manager API 重构带来的变更，贡献者：@yulesxoxo（[PR #4476](https://github.com/wailsapp/wails/pull/4476)）
- 修复 [#4456](https://github.com/wailsapp/wails/issues/4456) 中报告的 Windows 空指针解引用错误，由 @leaanthony 在 [#4460](https://github.com/wailsapp/wails/pull/4460) 中修复
- 在 macOS WKWebView 中添加对 `allowsBackForwardNavigationGestures` 的支持，以启用双指轻扫导航手势（#1857）
- 修复了由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中解决的一个问题：初始设置为禁用的菜单项无法使用 onClick。感谢 @IanVS 的前期调查。
- 修复构建失败时未清理 Vite 服务器的问题（#4403）
- 修复了在 Windows 上关闭或取消`SaveFileDialog`时发生 panic 的问题。由 @hkhere 在[PR](https://github.com/wailsapp/wails/pull/4284)中修复
- 由[@mbaklor](https://github.com/mbaklor)在[#4259](https://github.com/wailsapp/wails/pull/4259)中修复了 Windows 上 HTML 层级的拖放功能
- 在 macOS WKWebView 中添加对`allowsBackForwardNavigationGestures`的支持，以启用双指轻扫导航手势（#1857）
- 修复了由 @leaanthony 在[PR #4469](https://github.com/wailsapp/wails/pull/4469)中解决的一个问题：初始设置为禁用的菜单项无法使用 onClick。感谢 @IanVS 的前期调查。
- 修复构建失败时未清理 Vite 服务器的问题（#4403）
- 由 @popaprozac 在[PR](https://github.com/wailsapp/wails/pull/4450)中修复了 Windows 上的通知解析问题
- 由[@kodumulo](https://github.com/kodumulo)在[#4390](https://github.com/wailsapp/wails/issues/4390)中修复了 doctor 命令对 Windows SDK 依赖项的检查
- 由[@etesam913](https://github.com/etesam913)在[#4366](https://github.com/wailsapp/wails/pull/4366)中修复了 Mac 上 processURLRequest 中的 nil 指针解引用问题
- 由[@bh90210](https://github.com/bh90210)在[#4287](https://github.com/wailsapp/wails/pull/4287)中修复了导致 Linux 上无法使用带筛选条件对话框的错误
- 由[@leaanthony](https://github.com/leaanthony)在[#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)中修复了 Windows 和 Linux 上的编辑菜单问题
- 由[@AkshayKalose](https://github.com/AkshayKalose)在[#3981](https://github.com/wailsapp/wails/pull/3981)中将 macOS .plist 文件中的最低系统版本从10.13.0更新为10.15.0
- 由[@leaanthony](https://github.com/leaanthony)修复窗口 ID 跳号问题
- 由[@leaanthony](https://github.com/leaanthony)修复调用 RegisterContextMenu 时菜单为 nil 的问题
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修复了绑定生成器输出中的依赖循环
- 由[@fbbdev](https://github.com/fbbdev)在[#4001](https://github.com/wailsapp/wails/pull/4001)中修复了绑定生成器输出中的先使用后定义错误
- 由[@fbbdev](https://github.com/fbbdev)在[#4023](https://github.com/wailsapp/wails/pull/4023)中实现将构建标志传递给绑定生成器
- 由[@leaanthony](https://github.com/leaanthony)将 Windows Taskfile 中的路径改为使用正斜杠，以确保其可在非 Windows 平台上运行
- 由[@leaanthony](https://github.com/leaanthony)修复了 Mac + Mac JS 事件
- 由[@leaanthony](https://github.com/leaanthony)修复了 macOS 上的事件死锁
- 由[@leaanthony](https://github.com/leaanthony)修复了在 Windows 上初始化窗口时，提供 HTML 但未提供 JS 所导致的`Parameter incorrect`错误
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修复了资源服务器进行内容类型探测时所用响应前缀的大小
- 由[@fbbdev](https://github.com/fbbdev)在[#4049](https://github.com/wailsapp/wails/pull/4049)中修复了资源服务器在根索引路径上对非404响应的处理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器测试泛型类型属性时的未定义行为
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了底层类型与具名包装器属性不同时，绑定生成器为模型生成的输出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器针对映射键类型生成的输出及其预处理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器为实现 marshaler 接口的结构体生成的输出
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器对涉及泛型类型的类型循环的检测
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器输出中对未导出模型的无效引用
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中将注入的代码移至服务文件末尾
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中修复了绑定生成器对文件关闭操作所返回错误的处理
- 由[@fbbdev](https://github.com/fbbdev)在[#4045](https://github.com/wailsapp/wails/pull/4045)中禁止对定义了生命周期或 HTTP 方法但没有其他绑定方法的服务发出警告
- 由[@marcus-crane](https://github.com/marcus-crane)在[#4056](https://github.com/wailsapp/wails/pull/4056)中修复了非 React 模板在使用浅色系统配色方案时无法显示 Hello World 页脚的问题
- 由[@leaanthony](https://github.com/leaanthony)修复了 macOS 上隐藏的菜单项
- 由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修复了消息处理器对错误的处理和格式化
-  由[@fbbdev](https://github.com/fbbdev)在[#4066](https://github.com/wailsapp/wails/pull/4066)中修复了退出应用程序时跳过服务关闭的问题
-  由[@leaanthony](https://github.com/leaanthony)确保菜单更新在主线程上进行
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中增强了拖动和调整大小机制的稳健性，使其更符合预期的平台行为
- 由[@fbbdev](https://github.com/fbbdev)在[#4100](https://github.com/wailsapp/wails/pull/4100)中修复了[#4097](https://github.com/wailsapp/wails/issues/4097)：Webpack/Angular 会丢弃运行时初始化代码
- 由[@IanVS](https://github.com/IanVS)在[#4116](https://github.com/wailsapp/wails/pull/4116)中修复了初始隐藏的菜单项
- 修复了当无扩展名请求中的`[request]`不存在但`[request].html`存在时，assetFileServer 不提供`.html`文件的问题
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4125](https://github.com/wailsapp/wails/pull/4125)中修复了图标生成路径
- 由[@oSethoum](https://github.com/osethoum)在[#4130](https://github.com/wailsapp/wails/pull/4130)中修复了未触发`fullscreen`、`unfullscreen`、`unminimise`和`unmaximise`事件的问题
- 由[@robin-samuel](https://github.com/robin-samuel)在[#4126](https://github.com/wailsapp/wails/pull/4126)中修复了配置中默认版本前缀不正确所导致的 NSIS 错误
- 由[TheGB0077](https://github.com/TheGB0077)在[#4188](https://github.com/wailsapp/wails/pull/4188)中修复了 Dialogs 运行时函数在 Windows 上返回转义路径的问题
- 由[@leaanthony](https://github.com/leaanthony)修复了 HKCU 中的 Webview2 检测路径。
- 由[@leaanthony](https://github.com/leaanthony)修复了 macOS 上的输入问题。
- 由[@yulesxoxo](https://github.com/yulesxoxo)在[#4219](https://github.com/wailsapp/wails/pull/4219)中修复了 Windows 图标生成任务的文件名。
- [@leaanthony](https://github.com/leaanthony)基于 @kron 的工作修复了无边框窗口的透明度问题。
- [@leaanthony](https://github.com/leaanthony)基于 @kron 的工作修复了窗口被禁用或最小化时的聚焦调用。
- [@leaanthony](https://github.com/leaanthony)基于 @kron 的工作修复了任务栏重启后系统托盘不显示的问题。
- 在[#4245](https://github.com/wailsapp/wails/pull/4245)中修复了 fallbackResponseWriter 未实现 Flush() 的问题
- 由 [@superDingda] 在[#4236](https://github.com/wailsapp/wails/issues/4236)中修复了 fallbackResponseWriter 未实现 Flush() 的问题
- 由[@joshhardy](https://github.com/joshhardy)在[#4354](https://github.com/wailsapp/wails/pull/4354)中修复了 macOS 窗口在异步 Go 绑定函数调用尚未完成时关闭会导致崩溃的问题
- 由[@leaanthony](https://github.com/leaanthony)修复了 Windows 效率模式启动时的竞态条件
- 由[@leaanthony](https://github.com/leaanthony)修复了 Windows 图标句柄清理问题。
- 由[@PPTGamer](https://github.com/PPTGamer)在[#4375](https://github.com/wailsapp/wails/pull/4375)中修复了 Windows 上的`OpenFileManager`。
- 由 @atterpac 在[#3979](https://github.com/wailsapp/wails/pull/3979)中修复了 Linux 的最小/最大宽度选项
- 由 @atterpac 在[#3966](https://github.com/wailsapp/wails/pull/3966)中通过提升 npm 版本更新了 TypeScript 模板的类型定义
- 由 @atterpac 在[#3945](https://github.com/wailsapp/wails/pull/3945)中修复了 SvelteKit 模板的 CSS 引用
- 由[@leaanthony](https://github.com/leaanthony)确保 window run() 中的关键回调在主线程上调用
- 由[@leaanthony](https://github.com/leaanthony)修复了对话框目录选择器示例
- 由[@leaanthony](https://github.com/leaanthony)为缺少 index.html 的情况创建了新的中文错误页面
-  由[@leaanthony](https://github.com/leaanthony)确保`windowDidBecomeKey`回调在主线程上运行
-  由[@leaanthony](https://github.com/leaanthony)为无边框窗口添加全屏支持
-  由[@leaanthony](https://github.com/leaanthony)改进了窗口销毁逻辑
-  由[@leaanthony](https://github.com/leaanthony)修复了窗口附加到系统托盘时的位置逻辑
-  由[@leaanthony](https://github.com/leaanthony)为无边框窗口添加全屏支持
- 由[@leaanthony](https://github.com/leaanthony)修复了事件处理
- 由[@leaanthony](https://github.com/leaanthony)修复了窗口关闭逻辑
- 由[@leaanthony](https://github.com/leaanthony)将通用 Taskfile 调整为默认给 TypeScript 模板生成 TypeScript 绑定
- 由[@mmalcek](https://github.com/mmalcek)在[#3990](https://github.com/wailsapp/wails/pull/3990)中修复了没有打开任何窗口或仅有系统托盘时，收到 WM_CLOSE 消息后关闭应用程序的问题
- 由 @5aaee9 在[#3192](https://github.com/wailsapp/wails/pull/3192)中修复了 garble 构建
- 由[@leaanthony](https://github.com/leaanthony)修复了 Windows NSIS 构建
- 修复了未关闭内容导致 Linux 多选对话框死锁的问题
- 修复了 Windows 构建期间对 .syso 文件的跨平台清理，由
- 修复了 amd64 AppImage 编译，由 @atterpac 在
- 修复了构建资源更新，由 @ansxuman 在
- 修复了 Linux 系统托盘的`OnClick`和`OnRightClick`实现，由 @atterpac 完成
- 修复了`AlwaysOnTop`在 Mac 上无法工作的问题，由
-  修复了`application.NewEditMenu`包含重复项的问题
- 🐧 修复了 aarch64 编译
- ⊞ 修复了单选组菜单项，由
- 修复了在 MacOS 上构建可运行的 .app 时，'name' 和 'outputfilename' 导致的错误
- 修复了拖放示例中使用 customEventProcessor 时的错误，由
- 🐧 修复了添加 IgnoreMouseEvents 后引入的 Linux 编译错误，由
- ⊞ 修复了 syso 图标文件生成错误，由
- 🐧 合入了以 Wayland 原生模式运行的修复，来自
- 不在以下位置绑定内部服务方法
- ⊞ 在以下位置修复了系统托盘启动时的 panic
- 不在以下位置绑定内部服务方法
- ⊞ 在以下位置修复了系统托盘启动时的 panic
- 大幅重构菜单项和事件处理。目前主要改进了 macOS。作者：
- 修复插件和事件重构后的测试，见
- ⊞ 修复了`Failed to unregister class Chrome_WidgetWin_0`警告。作者：
- 模块问题
- 由[atterpac](https://github.com/atterpac)在以下位置修复了调整大小事件消息传递
- 🐧修复了 NixOS 上的主题处理错误，由
- 修复了 Windows 上跨卷安装项目的问题，由
- 修复了 React 模板的 CSS，以显示页脚，由
- 通过更新到最新版 refresh，修复了开发模式下的僵尸进程
- 由[Atterpac](https://github.com/atterpac)修复了 AppImage 的 WebKit 文件获取
- 由[Atterpac](https://github.com/Atterpac)在以下位置修复了 Doctor 的 apt 软件包验证
- 由 @5aaee9 在以下位置修复了退出时应用程序冻结的问题（Darwin）
- 修复了 Windows 上示例的背景颜色，由
- 由[mmghv](https://github.com/mmghv)在以下位置修复了默认上下文菜单
- 修复了 Darwin 上方向键的十六进制值，由
- 使 Windows 上的拖放功能可以正常工作。添加者：
- 修复了用户没有正确驱动程序时 Linux 上 Doctor 的错误
- 修复了启动时的 DPI 缩放问题（Windows）。由[@almas-x](https://github.com/almas-x)在以下位置更改
- 修复`go.mod`中的替换行以使用相对路径。修复包含以下内容的 Windows 路径
- 修复了没有附加窗口时 MacOS 系统托盘的点击处理，由
- 修复了未知选项导致的 Windows 构建失败，由
- 修复了 Windows 上没有以下内容时左键单击系统托盘图标导致的崩溃
- 由 @5aaee9 在 PR 中修复了两次打开窗口时 baseURL 错误的问题
- 修复`WebviewWindow.Restore`方法中 if 分支的顺序，由
- 在以下情况下，跨多次`GetStartURL`调用正确计算`startURL`
- 修复`Screen`结构体的 JS 类型，使其与对应的 Go 类型一致，由
- 修复`WML.Reload`方法，确保正确清理已注册的事件
- 修复了 Linux 上自定义上下文菜单立即关闭的问题，由
- 修复绑定生成的模型文件的输出路径和扩展名
- 修复绑定生成的 JS 代码中模型文件的导入路径
- 修复部分 Linux 发行版上的拖放问题，作者：
- 修复使用`wails3 task dev`时缺少 macOS 任务的问题，作者：
- 修复注册事件导致向 nil map 赋值的问题，作者：
- 修复绑定方法参数的反序列化问题，作者：
- 修复绑定方法返回多个值时的处理问题，作者：
- 修复 doctor 无法检测并非通过系统包管理器安装的 npm 的问题
- 修复缺少 MicrosoftEdgeWebview2Setup.exe 的问题。感谢
- 修复因窗口 ID 处理而导致 Linux 上随机崩溃的问题，作者：@leaanthony。基于
- 修复 systemTray.setIcon 在 Linux 上崩溃的问题，作者：
- 修复在以下平台上首次调用`setFrameless`函数时确保应用窗口边框的问题

### 变更

- **破坏性变更**：生成的 JS/TS 绑定中的 map 键现标记为可选，以准确反映 Go map 的语义。TypeScript 中访问 map 值现在返回`T | undefined`而非`T`，因此需要进行空值检查或断言（#4943），作者：`@fbbdev`
- 根据`@wailsio/runtime`中的变更，将`Event`改为使用`Events`，并相应调整`Features/Events/Event System`文档中的函数调用，作者：@AbdelhadiSeddar
- 将`EnabledFeatures`、`DisabledFeatures`和`AdditionalBrowserArgs`从各窗口选项移至应用级`Options.Windows`（#4559），作者：@leaanthony
- 更新`Drag N Drop`示例的 README，并强调该示例演示了`Internal Drag and Drop`，作者：@ndianabasi
- 将多处调试日志的级别从 Info 改为 Debug（作者：@mbaklor）
- <strong>破坏性变更：</strong>在窗口选项中将`EnableDragAndDrop`重命名为`EnableFileDrop`
- <strong>破坏性变更：</strong>在事件上下文中将`DropZoneDetails`重命名为`DropTargetDetails`
- <strong>破坏性变更：</strong>将`WindowEventContext`上的`DropZoneDetails()`方法重命名为`DropTargetDetails()`
- <strong>破坏性变更：</strong>移除`WindowDropZoneFilesDropped`事件，改用`WindowFilesDropped`
- <strong>破坏性变更：</strong>将 HTML 属性从`data-wails-dropzone`改为`data-file-drop-target`
- <strong>破坏性变更：</strong>将 CSS 悬停类从`wails-dropzone-hover`改为`file-drop-target-active`
- <strong>破坏性变更：</strong>从 Windows 中移除`DragEffect`、`OnEnterEffect`和`OnOverEffect`选项（这些选项原属于已移除的 IDropTarget）
- 所有运行时 JSON 处理（方法绑定、事件、WebView 请求、通知、kvstore）均改用 goccy/go-json，性能提升21-63%，内存分配减少40-60%
- 优化 BoundMethod 结构体布局并缓存 isVariadic 标志，以降低每次调用的开销
- 对具有`<=8`个参数的方法使用栈上分配的参数缓冲区，以避免堆分配
- 优化方法调用中的结果收集，避免为单个返回值分配切片
- MIME 类型缓存改用 sync.Map，以提升并发性能
- 读取 HTTP 传输请求体时使用缓冲池
- 在内容类型嗅探器中延迟分配 CloseNotify 通道，以减少每次请求的内存分配
- 移除资源服务器中的 CSS 调试日志
- 扩展 MIME 类型扩展名映射，以涵盖50多种常见 Web 格式（字体、音频、视频等）
- 更新窗口`X/Y`选项的文档，作者：@ruhuang2001
- 更新`Frontend Runtime`文档，添加更多生成前端绑定的选项，作者：@ndianabasi
- 更新 Wails v3资源服务器的文档页面，作者：@ndianabasi
- **破坏性变更**：移除包级对话框函数（`application.InfoDialog()`、`application.QuestionDialog()`等）。请改用`app.Dialog`管理器：`app.Dialog.Info()`、`app.Dialog.Question()`、`app.Dialog.Warning()`、`app.Dialog.Error()`、`app.Dialog.OpenFile()`、`app.Dialog.SaveFile()`
- 更新对话框文档以匹配实际 API：使用带回调的`app.Dialog.*`、`AddButton()`（而非`SetButtons()`）、`SetDefaultButton(*Button)`（而非字符串）、`AddFilter()`（而非`SetFilters()`）、`SetFilename()`（而非`SetDefaultFilename()`），并使用`app.Dialog.OpenFile().CanChooseDirectories(true)`选择文件夹
- **破坏性变更**：现在默认进行生产构建。要创建开发构建，请在 Taskfile 中设置`DEV=true`。生成新项目作为示例，作者：@leaanthony
- 发出包含零个或一个数据参数的自定义事件时，数据值将直接赋给 Data 字段，而不会封装在切片中，作者：[@fbbdev](https://github.com/fbbdev)，见[#4633](https://github.com/wailsapp/wails/pull/4633)
- Windows 托盘现在通过切换`NIS_HIDDEN`来遵循`SystemTray.Show()`/`Hide()`，因此应用可以真正消失并重新显示（#4653）。
- 托盘注册会复用已解析的图标，仅设置一次`NOTIFYICON_VERSION_4`，并启用`NIF_SHOWTIP`，以便工具提示在 Explorer 重启后恢复（#4653）。
- macOS：窗口居中时使用`visibleFrame`而非`frame`，以排除菜单栏和程序坞区域
- macOS：窗口居中时使用`visibleFrame`而非`frame`，以排除菜单栏和程序坞区域
- 使用`-config`参数运行`wails3 update build-assets`时，通过`-product*`参数设置的值会
- `window.NativeWindowHandle()` -> `window.NativeWindow()`，作者：@leaanthony，见[#4471](https://github.com/wailsapp/wails/pull/4471)
- 重构内部窗口处理，作者：@leaanthony，见[#4471](https://github.com/wailsapp/wails/pull/4471)
- 移除`application.WindowIDKey`和`application.WindowNameKey`（由`application.WindowKey`取代），作者：[@leaanthony](https://github.com/leaanthony)
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，固定长度数组类型的类字段现在会按预期长度初始化，而不再初始化为空数组，作者：[@fbbdev](https://github.com/fbbdev)，见[#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData 现在返回 string 而非 any，作者：[@leaanthony](https://github.com/leaanthony)
- `application.NewService`不再接受 options 作为可选参数（请改用`application.NewServiceWithOptions`），作者：[@leaanthony](https://github.com/leaanthony)，见[#4024](https://github.com/wailsapp/wails/pull/4024)
- 移除`nanoid`依赖，作者：[@leaanthony](https://github.com/leaanthony)
- 更新窗口示例，展示 mica/acrylic/tabbed 窗口样式，作者：[@leaanthony](https://github.com/leaanthony)
- 在 JS/TS 绑定中，已移除`internal.js/ts`模型文件；现在所有模型均可在`models.js/ts`中找到，作者：[@fbbdev](https://github.com/fbbdev)，见[#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 绑定中，命名类型绝不会再呈现为其他命名类型的别名；原有行为现在仅适用于别名，作者：[@fbbdev](https://github.com/fbbdev)，见[#4045](https://github.com/wailsapp/wails/pull/4045)
- 在 JS/TS 绑定的类模式下，类型为类型参数的结构体字段会标记为可选，并且绝不会自动初始化，由 [@fbbdev](https://github.com/fbbdev) 在 [#4045](https://github.com/wailsapp/wails/pull/4045) 中实现
- 从模板中移除 ESLint，由 [@IanVS](https://github.com/IanVS) 在 [#4059](https://github.com/wailsapp/wails/pull/4059) 中实现
- 将版权日期更新为 2025，由 [@IanVS](https://github.com/IanVS) 在 [#4037](https://github.com/wailsapp/wails/pull/4037) 中实现
- 添加 event.Sender 文档，由 [@IanVS](https://github.com/IanVS) 在 [#4075](https://github.com/wailsapp/wails/pull/4075) 中实现
- 支持 Go 1.24，由 [@leaanthony](https://github.com/leaanthony) 实现
- `ServiceStartup`钩子现在会在调用 `App.Run` 时触发，而不是在 `application.New` 中触发，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- `ServiceStartup`错误现在会从 `App.Run` 返回，而不会终止进程，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 来自 JS 的绑定和对话框调用现在会以错误对象而非字符串拒绝，由 [@fbbdev](https://github.com/fbbdev) 在 [#4066](https://github.com/wailsapp/wails/pull/4066) 中实现
- 改进 Windows 上系统托盘菜单的定位，由 [@leaanthony](https://github.com/leaanthony) 实现
- JS 运行时已移植到 TypeScript，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- 运行时现在导入后立即初始化，无需等待窗口加载，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- 运行时不再导出 init 方法。可以通过仅为副作用而导入来初始化运行时，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- 绑定方法现在返回一个 `CancellablePromise`；如果调用被取消，它会以 `CancelError` 拒绝。调用的实际结果会被丢弃，由 [@fbbdev](https://github.com/fbbdev) 在 [#4100](https://github.com/wailsapp/wails/pull/4100) 中实现
- 内置服务类型现在统一命名为 `Service`，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 带选项的内置服务创建函数现在统一命名为 `NewWithConfig`，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 为与 Go API 保持一致，`sqlite`服务上的 `Select` 方法现在命名为 `Query`，由 [@fbbdev](https://github.com/fbbdev) 在 [#4067](https://github.com/wailsapp/wails/pull/4067) 中实现
- 模板：将运行时移至 "dependencies"，并整理 package.json 文件，由 [@IanVS](https://github.com/IanVS) 在 [#4133](https://github.com/wailsapp/wails/pull/4133) 中实现
- 在开发期间创建应用程序包并进行临时签名，以启用某些 macOS API，由 [@popaprozac](https://github.com/popaprozac) 在 [#4171](https://github.com/wailsapp/wails/pull/4171) 中实现
- 将构建资源移至特定于平台的目录，由 [@leaanthony](https://github.com/leaanthony) 实现
- 将 Taskfile 移至特定于平台的目录并重命名，由 [@leaanthony](https://github.com/leaanthony) 实现
- 大幅改善缺少 `index.html` 时的体验，由 [@leaanthony](https://github.com/leaanthony) 实现
- [Windows] 改进最小化和还原的性能，由 [@leaanthony](https://github.com/leaanthony) 实现。基于 [562589540](https://github.com/562589540) 最初提交的 [PR](https://github.com/wailsapp/wails/pull/3955)
- 移除 `ShouldClose` 选项（请改为为 events.Common.WindowClosing 注册钩子），由 [@leaanthony](https://github.com/leaanthony) 实现
- [Windows] 减少打开窗口时的闪烁，由 [@leaanthony](https://github.com/leaanthony) 实现
- 移除 `Window.Destroy`，因为它原本应是内部函数，由 [@leaanthony](https://github.com/leaanthony) 实现
- 将 `WindowClose` 事件重命名为 `WindowClosing`，由 [@leaanthony](https://github.com/leaanthony) 实现
- 前端构建现在会根据构建类型使用 vite 环境 "development" 或 "production"，由 [@leaanthony](https://github.com/leaanthony) 实现
- 更新至 go-webview2 v1.19，由 [@leaanthony](https://github.com/leaanthony) 实现
- 确保使用 taskfile 的 fork，由 @leaanthony 实现
- 更新 Taskfile 的 fork，以修复使用以下方式安装时的版本问题
- 使用 Taskfile 的 fork，以修复使用以下方式安装时的版本问题
- `service.OnStartup`现在会在出错时关闭应用程序并运行
- 重构系统托盘点击消息传递，使其更符合用户交互，由
- 资源嵌入包含 `all:frontend/dist`，以支持会生成以下内容的框架
- 重构 Taskfile，由 [leaanthony](https://github.com/leaanthony) 在以下位置实现
- 升级至 `go-webview2` v1.0.16，由
- 修复 `Screen` 类型，使其包含 `ID` 而不是 `Id`，由
- 更新 `go.mod.tmpl` 的 Wails 版本，以支持 `application.ServiceOptions`，由
- 修复服务名称的确定方式，由 [windom](https://github.com/windom/) 在以下位置实现
- mkdocs serve 现在使用 Docker，由 [leaanthony](https://github.com/leaanthony) 实现
- 将开发配置整合到 `config.yml` 中，由
- 系统托盘对话框现在默认使用应用程序图标（如有）（Windows），由
- 改进 macOS 上 GPU 和内存信息的报告，由
- 移除 `WebviewGpuIsDisabled` 和 `EnableFraudulentWebsiteWarnings`
- Events API 变更：`On`/`Emit` -> 用户事件，`OnApplicationEvent` ->
- 修复 Linux 上的 Events API，由 [TheGB0077](https://github.com/TheGB0077) 在以下位置实现
- [CI] 改进 actions，并允许在 fork 和以下环境中也运行 actions
- 将 `AbsolutePosition()` 重命名为 `Position()`，由
- 将 Linux WebKit 依赖项从 webkitgtk2-4.0 更新为 webkit2gtk-4.1，以
- 捆绑的 JS 运行时脚本现在是 ESM 模块：导入它的 script 标签
- `@wailsio/runtime`包不再将其 API 发布到 `window.wails`
- Window API 模块 `@wailsio/runtime/src/window` 现在会公开其所在的
- JS 窗口 API 已更新，以匹配当前的 Go `WebviewWindow`
- 绑定生成器现在默认使用按 ID 调用。`-id` CLI 选项
- 新的绑定代码布局：输出文件以前按文件夹组织
- 结构体字段 `application.Options.Bind` 已重命名为
- 绑定服务的新语法：现在必须将服务实例包装在一个
- 在非终端或 CI 环境中禁用旋转指示器，由

### 已移除

- **破坏性变更**：从每个窗口的`WindowsWindow`选项中移除`EnabledFeatures`、`DisabledFeatures`和`AdditionalLaunchArgs`。请改用应用程序级的`Options.Windows.EnabledFeatures`、`Options.Windows.DisabledFeatures`和`Options.Windows.AdditionalBrowserArgs`。这些标志全局应用于共享的 WebView2 环境（#4559），由 @leaanthony 完成
- 移除 Windows 上原生的`IDropTarget`实现，改用基于 JavaScript 的方式（与v2的行为一致）
- 移除 github.com/wailsapp/mimetype 依赖，改用扩展后的扩展名映射和标准库 http.DetectContentType，使二进制文件大小减少约1.2MB
- 为 Linux 文件资源管理器实现精简的 .desktop 文件解析器，以移除 gopkg.in/ini.v1依赖，节省约45KB
- 使用 Go 1.21+ 标准库的 slices 包和精简的内部辅助函数，从运行时代码中移除 samber/lo，节省约310KB
- 从 Darwin URL 方案处理程序中移除调试 printf 语句（#4834）
- **破坏性变更**：移除`linux:WindowLoadChanged`事件；请改用`linux:WindowLoadFinished`检测 WebView 何时加载完成（#3896），由 @leaanthony 完成

### 破坏性变更

- **Manager API 重构**：将应用程序 API 从扁平结构重组为分类明确的管理器，以改善代码组织并提高可发现性；由[@leaanthony](https://github.com/leaanthony)在[#4359](https://github.com/wailsapp/wails/pull/4359)中完成
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
- 重命名 Service 方法：`Name` -> `ServiceName`、`OnStartup` -> `ServiceStartup`、`OnShutdown` -> `ServiceShutdown`，由[@leaanthony](https://github.com/leaanthony)完成
- 将`Path`和`Paths`方法移至`application`包，由[@leaanthony](https://github.com/leaanthony)完成
- 应用程序菜单现在仅适用于 macOS，由[@leaanthony](https://github.com/leaanthony)完成

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

## 已修复

## v3.0.0-alpha.56 - 2026-01-04

## 新增

## 已更改

## 已修复

## 已移除

## v3.0.0-alpha.55 - 2026-01-02

## 已更改

## 已修复

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

+ 修复范围过宽的 Linux 软件包依赖项，并修复已过时的 RPM 依赖项。

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

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### 变更

`service.OnShutdown`适用于此前已启动的所有服务，由 @atterpac 在[#3920](https://github.com/wailsapp/wails/pull/3920)中完成   @atterpac 在[#3907](https://github.com/wailsapp/wails/pull/3907)中完成   子文件夹，由 @atterpac 在[#3887](https://github.com/wailsapp/wails/pull/3887)中完成   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913)在[#3778](https://github.com/wailsapp/wails/pull/3778)中完成   [northes](https://github.com/northes)在[#3836](https://github.com/wailsapp/wails/pull/3836)中完成   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   （已由`EnabledFeatures`和`DisabledFeatures`选项取代），由[leaanthony](https://github.com/leaanthony)完成

### 修复

通道变量，由 @michael-freling 在[#3925](https://github.com/wailsapp/wails/pull/3925)中完成   [ansxuman](https://github.com/ansxuman)在[#3924](https://github.com/wailsapp/wails/pull/3924)中完成   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   在[#3886](https://github.com/wailsapp/wails/pull/3886)中   [leaanthony](https://github.com/leaanthony)在[#3841](https://github.com/wailsapp/wails/pull/3841)中完成   Darwin 编辑菜单中的`PasteAndMatchStyle`角色，由[johnmccabe](https://github.com/johnmccabe)在[#3839](https://github.com/wailsapp/wails/pull/3839)中完成   [#3840](https://github.com/wailsapp/wails/issues/3840)在[#3854](https://github.com/wailsapp/wails/pull/3854)中，由[kodflow](https://github.com/kodflow)完成   [@leaanthony](https://github.com/leaanthony)   各不相同。由 @nickisworking 在[#3789](https://github.com/wailsapp/wails/pull/3789)中完成

## v3.0.0-alpha.7 - 2024-09-18

### 新增

[mmghv](https://github.com/mmghv)在[#3665](https://github.com/wailsapp/wails/pull/3665)中完成   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac)和[leaanthony](https://github.com/leaanthony)在[#3570](https://github.com/wailsapp/wails/pull/3570)中完成

### 变更

应用程序事件`OnWindowEvent` -> 窗口事件，由[leaanthony](https://github.com/leaanthony)完成   [#3734](https://github.com/wailsapp/wails/pull/3734)   以`v3/`或`v3-`为前缀的分支，由[stendler](https://github.com/stendler)在[#3747](https://github.com/wailsapp/wails/pull/3747)中完成

### 修复

[etesam913](https://github.com/etesam913)在[#3742](https://github.com/wailsapp/wails/pull/3742)中完成   [atterpac](https://github.com/atterpac)在[#3721](https://github.com/wailsapp/wails/pull/3721)中完成   [atterpac](https://github.com/atterpac)在[#3675](https://github.com/wailsapp/wails/pull/3675)中完成   [#1811](https://github.com/wailsapp/wails/pull/1811)在[#3614](https://github.com/wailsapp/wails/pull/3614)中，由[@stendler](https://github.com/stendler)完成   [#3720](https://github.com/wailsapp/wails/pull/3720)，由[leaanthony](https://github.com/leaanthony)完成   [#3693](https://github.com/wailsapp/wails/issues/3693)，由[@DeltaLaboratory](https://github.com/DeltaLaboratory)完成   [#3720](https://github.com/wailsapp/wails/pull/3720)，由[leaanthony](https://github.com/leaanthony)完成   [#3693](https://github.com/wailsapp/wails/issues/3693)，由[@DeltaLaboratory](https://github.com/DeltaLaboratory)完成   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746)，由[@stendler](https://github.com/stendler)完成   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### 修复

## v3.0.0-alpha.5 - 2024-07-30

### 新增

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   点击图标，由 @5aaee9 在[#2991](https://github.com/wailsapp/wails/pull/2991)中完成   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev)在[#3282](https://github.com/wailsapp/wails/pull/3282)中完成   @[Atterpac](https://github.com/Atterpac)在[#3022](https://github.com/wailsapp/wails/pull/3022])中完成   [@marcus-crane](https://github.com/marcus-crane)在[#3146](https://github.com/wailsapp/wails/pull/3146)中完成   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev)在[#3281](https://github.com/wailsapp/wails/pull/3281)中完成   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   基于 @Mai-Lapyst 的[PR](https://github.com/wailsapp/wails/pull/2044)   [@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   [@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   [@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   npm 软件包，由[@fbbdev](https://github.com/fbbdev)在[#3334](https://github.com/wailsapp/wails/pull/3334)中完成   在[#3354](https://github.com/wailsapp/wails/pull/3354)中   `WAILS_VITE_PORT`，由[@abichinger](https://github.com/abichinger)在[#3429](https://github.com/wailsapp/wails/pull/3429)中完成   [@abichinger](https://github.com/abichinger)在[#3431](https://github.com/wailsapp/wails/pull/3431)中完成   [@bruxaodev](https://github.com/bruxaodev)在[#3667](https://github.com/wailsapp/wails/pull/3667)中完成   [@OlegGulevskyy](https://github.com/OlegGulevskyy)在[#3674](https://github.com/wailsapp/wails/pull/3674)中完成

### 修复

[#3606](https://github.com/wailsapp/wails/pull/3606)   由[tmclane](https://github.com/tmclane)在[#3515](https://github.com/wailsapp/wails/pull/3515)中完成   由[atterpac](https://github.com/atterac)在[#3512](https://github.com/wailsapp/wails/pull/3512)中完成   由[atterpac](https://github.com/atterpac)在[#3477](https://github.com/wailsapp/wails/pull/3477)中完成   由[Atterpac](https://github.com/atterpac)在[#3320](https://github.com/wailsapp/wails/pull/3320)中完成。   位于[#3306](https://github.com/wailsapp/wails/pull/3306)中。   [#2972](https://github.com/wailsapp/wails/pull/2972)。   [#2982](https://github.com/wailsapp/wails/pull/2982)   由[mmghv](https://github.com/mmghv)在[#2750](https://github.com/wailsapp/wails/pull/2750)中完成。   [#2753](https://github.com/wailsapp/wails/pull/2753)。   由[jaybeecave](https://github.com/jaybeecave)在[#3052](https://github.com/wailsapp/wails/pull/3052)中完成。   由[@pylotlight](https://github.com/pylotlight)在[PR](https://github.com/wailsapp/wails/pull/3039)中安装。由[@pylotlight](https://github.com/pylotlight)在[PR](https://github.com/wailsapp/wails/pull/3032)中添加   [PR](https://github.com/wailsapp/wails/pull/3145)   空格 - @leaanthony。   由[thomas-senechal](https://github.com/thomas-senechal)在 PR [#3207](https://github.com/wailsapp/wails/pull/3207)中完成   由[thomas-senechal](https://github.com/thomas-senechal)在 PR [#3208](https://github.com/wailsapp/wails/pull/3208)中完成   附加窗口，由[tw1nk](https://github.com/tw1nk)在 PR [#3271](https://github.com/wailsapp/wails/pull/3271)中完成   [#3273](https://github.com/wailsapp/wails/pull/3273)   由[@fbbdev](https://github.com/fbbdev)在[#3279](https://github.com/wailsapp/wails/pull/3279)中完成   存在`FRONTEND_DEVSERVER_URL`。   [#3299](https://github.com/wailsapp/wails/pull/3299)   由[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   监听器，由[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   由[@abichinger](https://github.com/abichinger)在[#3330](https://github.com/wailsapp/wails/pull/3330)中完成   生成器，由[@fbbdev](https://github.com/fbbdev)在[#3334](https://github.com/wailsapp/wails/pull/3334)中完成   生成器，由[@fbbdev](https://github.com/fbbdev)在[#3334](https://github.com/wailsapp/wails/pull/3334)中完成   由[@abichinger](https://github.com/abichinger)在[#3346](https://github.com/wailsapp/wails/pull/3346)中完成   由[@hfoxy](https://github.com/hfoxy)在[#3417](https://github.com/wailsapp/wails/pull/3417)中完成   由[@hfoxy](https://github.com/hfoxy)在[#3426](https://github.com/wailsapp/wails/pull/3426)中完成   由[@fbbdev](https://github.com/fbbdev)在[#3431](https://github.com/wailsapp/wails/pull/3431)中完成   由[@fbbdev](https://github.com/fbbdev)在[#3431](https://github.com/wailsapp/wails/pull/3431)中完成   由[@pekim](https://github.com/pekim)在[#3458](https://github.com/wailsapp/wails/pull/3458)中完成   [@robin-samuel](https://github.com/robin-samuel)。   PR [#3466](https://github.com/wailsapp/wails/pull/3622)，由[@5aaee9](https://github.com/5aaee9)完成。   由[@windom](https://github.com/windom/)在[#3636](https://github.com/wailsapp/wails/pull/3636)中完成。   Windows，由[@bruxaodev](https://github.com/bruxaodev/)在[#3691](https://github.com/wailsapp/wails/pull/3691)中完成。

### 更改

由[mmghv](https://github.com/mmghv)在[#3611](https://github.com/wailsapp/wails/pull/3611)中完成   支持 Ubuntu 24.04 LTS，由[atterpac](https://github.com/atterpac)在[#3461](https://github.com/wailsapp/wails/pull/3461)中完成   必须具有`type="module"`属性。由[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   对象，并且不会启动 WML 系统。这样做是为了改进封装。如有需要，可以调用新的`WML.Enable`方法手动启动 WML 系统。捆绑的 JS 运行时脚本仍会自动执行这两项操作。由[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   将 window 对象作为默认导出。现在无法再通过 ESM 命名导入或命名空间导入语法导入单个方法。   API。部分方法的名称或原型已更改，具体如下：`Screen`变为`GetScreen`；`GetZoomLevel`/`SetZoomLevel`变为`GetZoom`/`SetZoom`；`GetZoom`、`Width`和`Height`现在会直接返回值，而不再将其包装在对象中。由[@fbbdev](https://github.com/fbbdev)在[#3295](https://github.com/wailsapp/wails/pull/3295)中完成   已移除。使用`-names` CLI 选项可切换回按名称调用。由[@fbbdev](https://github.com/fbbdev)在[#3468](https://github.com/wailsapp/wails/pull/3468)中完成   以前按其所在包命名；现在使用完整的 Go 导入路径，包括模块路径。由[@fbbdev](https://github.com/fbbdev)在[#3468](https://github.com/wailsapp/wails/pull/3468)中完成   `application.Options.Services`。由[@fbbdev](https://github.com/fbbdev)在[#3468](https://github.com/wailsapp/wails/pull/3468)中完成   对`application.NewService`的调用。由[@fbbdev](https://github.com/fbbdev)在[#3468](https://github.com/wailsapp/wails/pull/3468)中完成   由[@DeltaLaboratory](https://github.com/DeltaLaboratory)在[#3574](https://github.com/wailsapp/wails/pull/3574)中完成
