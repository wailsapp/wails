# 更新日志

此项目的所有重要变更都将记录在此文件中。

格式基于 [维护更新日志](https://keepachangelog.com/en/1.0.0/)， 并且该项目遵循 [语义化版本](https://semver.org/spec/v2.0.0.html)。

## [Unreleased][]

## [v2.0.0-beta.39.2] - 2022-07-20

## Added

* 由 @acheong08 在 https://github.com/wailsapp/wails/pull/1600 中更新macOS 菜单示例

## Fixed

* Reinstate Go 1.17 compatibility by @leaanthony in https://github.com/wailsapp/wails/pull/1605

## [v2.0.0-beta.39] - 2022-07-19

## Added

* 新的屏幕尺寸运行时间API由 @skamensky 在 https://github.com/wailsapp/wails/pull/1519
* 通过 @leaanthony 在https://github.com/wailsapp/wails/pull/1547 自动发现vite devserver 端口
* 添加 nixpkgs 支持医生命令。 由 @ianmjones 在 https://github.com/wailsapp/wails/pull/1551
* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1578 开发的新预构建钩子
* 新的生产日志级别选项由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1555

## Fixed

* 修复Windows中使用 ICoreWebView2HttpheadersCollectionIterator 的 @stffabi 在 https://github.com/wailsapp/wails/pull/1589
* 将 WindowGet * 移动到主线程由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1464
* 允许 -appargs 标志传递标记到二进制文件。 由 @ianmjones 在 https://github.com/wailsapp/wails/pull/1534
* 以无英文会话修复已安装的 apt 软件包。 由 @ianmjones 在 https://github.com/wailsapp/wails/pull/1548
* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1558 修复Mac OnBEforeClos代码
* 支持由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1435 中转换TS 的地图
* Check for line length when scanning for local devserver url by @leaanthony in https://github.com/wailsapp/wails/pull/1566
* 在 https://github.com/wailsapp/wails/pull/1556 中删除winc 中由 @stffabi 和 @leaanthony 在 winc 中使用 unsafe.Pointer

## Changed

* 重命名WindowSetRGBA -> WindowSetBackgroundColour by @leaanthony in https://github.com/wailsapp/wails/pull/1506
* 由 @stffabi 在 https://github.com/wailsapp/wails/pull/1510 中对dev 命令的改进
* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1398 更新vscode 模板
* 在 /v2/internal/frontend/runtime/dev 由 @dependabot 从 3.42.2到 3.49.0 堆放在https://github.com/wailsapp/wails/pull/1572
* Bump svelte from 3.42.5 to 3.49.0 in /v2/internal/frontend/runtime by @dependabot in https://github.com/wailsapp/wails/pull/1573
* 在 https://github.com/wailsapp/wails/pull/1586 中添加 `在 @acheong08 中发现` 错误
* https://github.com/wailsapp/wails/pull/1591

## 新建贡献者

* @skamensky在https://github.com/wailsapp/wails/pull/1519 中做了他们的首次贡献
* @acheong08 在https://github.com/wailsapp/wails/pull/1586 中做了他们的首次贡献

**完整更新日志**: https://github.com/wailsapp/wails/compare/v2.0.0-beta.38...v2.0.0-beta.39

## [v2.0.0-beta.38] - 2022-06-27

### Added

* 在 https://github.com/wailsapp/wails/pull/1426 通过@Lyimmi 添加种族检测器 & dev
* [linux] 支持 `linux/arm` 架构由 @Lyimmi 在 https://github.com/wailsapp/wails/pull/1427
* 使用 `-g` 选项在 https://github.com/wailsapp/wails/pull/1430 由 @jaesung9507 创建 gitnore
* [windows] 在 https://github.com/wailsapp/wails/pull/1474 由 @leaanthony 添加 Suspend/Resume 回调支持
* 添加运行时函数 `WindowSetAlwaysOnTop` 由 @chenxiao1990 在 https://github.com/wailsapp/wails/pull/1442
* [windows] 允许用@NanoNik设置浏览器路径在 https://github.com/wailsapp/wails/pull/1448

### Fixed

* [linux] 改进由 @stffabi 在 https://github.com/wailsapp/wails/pull/1392 中转换回调的主线程
* [windows] 修复 WebView2 最小运行时间检查通过 @stffabi 在 https://github.com/wailsapp/wails/pull/1456
* [linux] 由 @abtin 在 https://github.com/wailsapp/wails/pull/1461 修复apt 命令语法 (#1458)
* [windows] 如果在 https://github.com/wailsapp/wails/pull/1466 通过 @leaanthony 设置窗口背景颜色
* 由 @LukenSkyne 在 https://github.com/wailsapp/wails/pull/1449 修复文档中的小类型
* 修复网址由 @andywenk 在 https://github.com/wailsapp/wails/pull/1460
* 修复运行时由 https://github.com/wailsapp/wails/pull/1473 @leaanthony 更改主题
* 修复：如果无法移除由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1465 构建的临时绑定，不要停止
* [windows] 将正确的安装状态传递到webview安装策略，由 @stffabi 在 https://github.com/wailsapp/wails/pull/1483
* [windows] 使 `设置背景颜色` 兼容 `windows/386` 由 @stffabi 在 https://github.com/wailsapp/wails/pull/1493
* 由 @Orijhins 修复lit-ts 模板在 https://github.com/wailsapp/wails/pull/1494

### Changed

* [windows] 只能通过@stffabi 在 https://github.com/wailsapp/wails/pull/1432 从嵌入的 WebView2 加载程序
* 添加展示条目，10月份由 @marcus-crane 更新主页旋转至10月份。https://github.com/wailsapp/wails/pull/1436
* 在https://github.com/wailsapp/wails/pull/1410 中使用@leaanthony 包装的退货方式
* [windows] Unlock OSThread after native calls have been finished by @stffabi in https://github.com/wailsapp/wails/pull/1441
* 添加 `背景颜色` 并废弃 `RGBA` 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1475
* AssetsHandler 删除 dev 模式 @stffabi 在 https://github.com/wailsapp/wails/pull/1479 中的重试逻辑。
* 在 https://github.com/wailsapp/wails/pull/1492 通过 @sidwebworks添加Solid JS 模板到文档
* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1488 中更好地处理信号
* 在 https://github.com/wailsapp/wails/pull/1489 通过 @tomanagle 创建 root 的18。

## 新建贡献者

* @jaesung9507在https://github.com/wailsapp/wails/pull/1430中做了他们的首次贡献
* @LukenSkyne 在https://github.com/wailsapp/wails/pull/1449 中做出了他们的首次贡献
* @andywenk 在https://github.com/wailsapp/wails/pull/1460中做了他们的首次贡献
* @abtin 在https://github.com/wailsapp/wails/pull/1461 中做了他们的首次贡献
* @chenxiao1990在https://github.com/wailsapp/wails/pull/1442 中做了他们的首次贡献
* @NanoNik在https://github.com/wailsapp/wails/pull/1448 中做了他们的首次贡献
* @sidwebworks在https://github.com/wailsapp/wails/pull/1492 中做出了他们的首次贡献
* @tomanagle 在https://github.com/wailsapp/wails/pull/1489中做了他们的首次贡献

## [v2.0.0-beta.37] - 2022-05-26

### Added

* 在 https://github.com/wailsapp/wails/pull/1413 中以 @mondy 在wails dev 命令中添加 `nogen` 标志
* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1400在 Windows 预览中对新本地半透明性的初步支持

### Fixed

* 由 @leaanthony 在 https://github.com/wailsapp/wails/pull/1383 编写的 Bugfix/不正确绑定
* 由 @polikow 修复运行时间.js 在 https://github.com/wailsapp/wails/pull/1369
* 修复由 @antimatter96 格式的文档：https://github.com/wailsapp/wails/pull/1372
* 事件 | 修正 #1388 由 @lambdajack 在 https://github.com/wailsapp/wails/pull/1390
* bugfix: correctly typo by @tmclane in https://github.com/wailsapp/wails/pull/1391
* 修复文档中的 @LGiki 在 https://github.com/wailsapp/wails/pull/1393
* 修复由 @rayshoo 在 https://github.com/wailsapp/wails/pull/1466 通过 ipc.js 绑定的
* 请确保在 https://github.com/wailsapp/wails/pull/1403 由 @stffabi 在一个新的goroutine 上执行菜单回调
* 更新运行时间.d.ts & 模板由 @Yz4230 在 https://github.com/wailsapp/wails/pull/1421
* 在 https://github.com/wailsapp/wails/pull/1419 由 @edwardbrowncres 添加缺失的类名称到React 和Preact 模板的输入

### Changed
* 由 @stffabi 在 https://github.com/wailsapp/wails/pulliture 改进多平台版本
* 在 https://github.com/wailsapp/wails/pull/1385 中@stffabi 使用的 AssetsHandler 仅使用重新加载逻辑。
* 在 https://github.com/wailsapp/wails/pull/1387 通过 @Junkher 更新 events.mdx
* 添加 Next.js 模板由 @LGiki 在 https://github.com/wailsapp/wails/pull/1394
* 在 https://github.com/wailsapp/wails/pull/1414 由 @TechplexEngineer 在生成器上添加文档
* 在 https://github.com/wailsapp/wails/pull/1423 通过 @daodao97 添加 macos 自定义菜单编辑器提示

### 新建贡献者
* @polikow在https://github.com/wailsapp/wails/pull/1369 中做出了他们的首次贡献
* @antimatter96 在https://github.com/wailsapp/wails/pull/1372 中做出了他们的首次贡献
* @Junkher在https://github.com/wailsapp/wails/pull/1387中做出了他们的首次贡献
* @lambdajack在https://github.com/wailsapp/wails/pull/1390中做出了他们的首次贡献
* @LGiki 在https://github.com/wailsapp/wails/pull/1393 中做了他们的首次贡献
* @rayshoo在https://github.com/wailsapp/wails/pull/1406中做出了他们的首次贡献
* @TechplexEngineer 在https://github.com/wailsapp/wails/pull/1414中做了他们的首次贡献
* @mondy 在https://github.com/wailsapp/wails/pull/1413中做了他们的首次贡献
* @Yz4230 在https://github.com/wailsapp/wails/pull/1421 中做了他们的首次贡献
* @daodao97 在https://github.com/wailsapp/wails/pull/1423中做了他们的首次贡献
* @edwardbrowncross在https://github.com/wailsapp/wails/pull/1419中做出了他们的第一个贡献


## [v2.0.0-beta.36] - 2022-04-27

### Fixed
- [v2] 验证 devServer 属性为正确的表单，由 [@stffabi](https://github.com/stffabi) 在 https://github.com/wailsapp/wails/pull/1359
- [v2, darwin] 初始化堆栈上的本地变量，以防止由 [@stffabi](https://github.com/stffabi) 在https://github.com/wailsapp/wails/pull/1362 中产生的断层故障
- Vue-TS 模板修复

### Changed
- 将 `Onstartup` 方法添加到默认模板

## [v2.0.0-beta.35] - 2022-04-27

### 打破更改

- 当数据被发送到 `EventsOn` 回调时 它是作为数值的分割发送的， 而不是方法的可选参数。 `事件` 现在可以正常工作，但如果您 目前使用此功能，您需要更新您的代码 [更多信息](https://github.com/wailsapp/wails/issues/1324)
- 已损坏的 `bindings.js` 和 `bindings.d.ts` 文件已被新的 JS/TS 代码生成系统所取代。 更多 详情 [在这里](https://wails.io/docs/howdoesitwork#calling-bound-go-methods)

### Added

- **新模板**: Swelte, React, Vue, Preact, Lit and Vanilla 模板, 既有JS 版本又有TS 版本。 `等待 -l` 获取更多信息 信息。
- 默认模板现在由 [Vite](https://vitejs.dev) 供电。 当你 使用 `wails dev` 时，启用闪电快速重新加载！
- 为外部前端开发服务器添加支持。 请参阅 `frontend:dev:serverUrl` [工程配置](https://wails.io/docs/reference/project-config) - [@stffabi](https://github.com/stffabi)
- [Windows完全配置的暗色模式](https://wails.io/docs/reference/options#theme)。
- 极大改进了 [WailsJS 生成](https://wails.io/docs/howdoesitwork#calling-bound-go-methods) (都是 Javascript 和类型)
- Wails医生现在报告有关安装wail的信息 - [@stffabi](https://github.com/stffabi)
- [代码签名](https://wails.io/docs/guides/signing) 和 [NSIS 安装程序](https://wails.io/docs/guides/windows-installer) - [@gardc](https://github.com/gardc)
- 添加对 `-trimpath` [构建标志](https://wails.io/docs/reference/cli#build)
- 添加对默认资产处理程序的支持 - [@stffabi](https://github.com/stffabi)

### Fixed

- 改进了 BOM 标记和注释的 mimetype 检测 - [@napalu](https://github.com/napalu)
- 删除重复的 mimetype 条目 - [@napalu](https://github.com/napalu)
- 删除生成的定义文件中重复的类型导入- [@adalessa](https://github.com/adalessa)
- 添加缺少的方法声明 - [@adalessa](https://github.com/adalessa)
- 启动时修复 Linux sigabrt - [@napalu](https://github.com/napalu)
- 双击事件现在可以使用 `数据wails-drag` 属性 - [@jicg](https://github.com/jicg)
- 最小化帧率窗口时禁止调整大小 - [@stffabi](https://github.com/stffabi)
- 固定的 TS/JS 生成用于不返回的去方法
- 修复工程目录中生成的WailsJS

### Changed

- 网站文档现已版本
- 改进 `runtime.Environment` 调用
- 改进Mac 的关闭操作
- 一堆依赖物的安全更新
- 改进网站内容 - [@misitebao](https://github.com/misitebao)
- 升级问题模板 - [@misitebao](https://github.com/misitebao)
- 将不需要版本管理的文档转换为单个页面
  - [@misitebao](https://github.com/misitebao)
- 正在使用Algolia搜索的网站

## [v2.0.0-beta.34] - 2022-03-26

### Added

- 在 [@napalu](https://github.com/napalu) 在 Linux 上添加对 'DomReady' 回调的支持 #1249
- MacOS - 默认显示扩展 [@leaanthony](https://github.com/leaanthony)  在 #1228

### Fixed

- [v2, nsis] 看起来像/作为路径分隔器只适用于一些指令。在#1227 中使用 由 [@stffabi](https://github.com/stffabi)
- 导入绑定定义模型由 [@adalessa](https://github.com/adalessa) 在 #1231
- 使用 [@leaanthony](https://github.com/leaanthony)  在网站上进行本地搜索 #1234
- 确保二进制资源可以由 [@napalu](https://github.com/napalu) 在 #1240
- 在#1241中从 [@leaanthony](https://github.com/leaanthony)  从磁盘中加载时仅重试加载素材
- [v2, window] 修复最大起始状态由 [@stffabi](https://github.com/stffabi) 在 #1243
- 确保Linux IsFullscreen使用GDK_WINDOW_STATE_FULLSCREEN bit掩码。 由 [@ianmjones](https://github.com/ianmjones) 在 #1245
- 修复Mac 的ExecJS内存泄漏由 [@leaanthony](https://github.com/leaanthony)  在 #1230
- 由 [@BillBuilt](https://github.com/BillBuilt) 在 #1247 中修复或至少一个工作区 (#1232)
- [v2] 在 #1258 中使用 os.Args[0] 来启动自己的外挂 [@stffabi](https://github.com/stffabi)
- [v2, windows] Windows开关方案: https -> http by @stefpap, in #1255
- 在 [@leaanthony](https://github.com/leaanthony)  在 #1257 中在 Web 视图2 中恢复聚焦。
- 在 Show() 被调用时尝试聚焦窗口。 由 [@leaanthony](https://github.com/leaanthony)  在 #1212
- 检查用户安装的 Linux 依赖关系的系统，由 [@leaanthony](https://github.com/leaanthony)  在 #1180 中

### Changed

- 功能 (网站)：在 [@misitebao](https://github.com/misitebao) 在 #1215 中同步文档并添加内容
- refactory(片段)：优化默认模板由 [@misitebao](https://github.com/misitebao) 在 #1214 中
- 由 [@leaanthony](https://github.com/leaanthony)  在 #1216 中初始构建后运行监视器
- 由 [@leaanthony](https://github.com/leaanthony)  在 #1218 中更新特色/文档
- 功能 (网站)：优化网站并同步文档由 [@misitebao](https://github.com/misitebao) 在 #1219
- 文档：由 [@misitebao](https://github.com/misitebao) 在 #1224 中同步文档
- 默认索引页面由 [@leaanthony](https://github.com/leaanthony)  在 #1229 中
- Build added win32 compatibility by [@fengweiqiang](https://github.com/fengweiqiang) in #1238
- 文档：由 [@misitebao](https://github.com/misitebao) 在 #1260 中同步文档

## [v2.0.0-beta.33][] - 2022-03-05

### Added

- NSIS Installer support for creating installers for Windows applications - Thanks [@stffabi](https://github.com/stffabi) 🎉
- New frontend:dev:watcher command to spin out 3rd party watchers when using wails dev - Thanks [@stffabi](https://github.com/stffabi)🎉
- Remote templates now support version tags - Thanks [@misitebao](https://github.com/misitebao) 🎉

### Fixed

- A number of fixes for ARM Linux providing a huge improvement - Thanks [@ianmjones](https://github.com/ianmjones) 🎉
- Fixed potential Nil reference when discovering the path to `index.html`
- Fixed crash when using `runtime.Log` methods in a production build
- Improvements to internal file handling meaning webworkers will now work on Windows - Thanks [@stffabi](https://github.com/stffabi)🎉

### Changed

- The Webview2 bootstrapper is now run as a normal user and doesn't require admin rights
- The docs have been improved and updated
- Added troubleshooting guide

[Unreleased]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.33...HEAD
[v2.0.0-beta.33]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.32...v2.0.0-beta.33
