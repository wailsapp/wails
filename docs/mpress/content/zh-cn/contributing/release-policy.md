---
title: "发布策略"
description: "Wails v3 如何确定每夜版和每周版的版本、维护者如何发布这些版本，以及用户为何应固定版本。"
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## 用一张表说明模型

Wails v3按**时间表而非功能**发布。在`v3.0.0`基线之后：

| 发布类型 | 标签 | 发布频率 | GitHub 发布版 | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| 每夜版 | `v3.<next>.0-nightly.<YYYYMMDD>`（UTC 日期；若在同一天重新运行，则追加`.N`） | 自动、每日，仅当`UNRELEASED_CHANGELOG.md`有内容时发布 | 预发布版，**不**标记为 Latest | 相同版本，但不含`v`，dist-tag 为`nightly` |
| 每周版 | `v3.<next>.0` | 自动，于 Australia/Sydney 时区星期日20:00发布<strong>（日期尚未确认）</strong> | 常规发布版，标记为 Latest | 相同版本，dist-tag 为`latest` |
| 热修复版 | `v3.<current>.<N>`, `N > 0` | 很少发布，手动操作 | 正常发布，Latest | 相同版本，dist-tag 为`latest` |
| 主版本 | `v4.0.0` | 仅用于破坏性变更 | 常规发布版 | `latest` |

`<next>`是已**发布**的最高 v3 次版本号加一。nightly 是**下一个 weekly** 的预览，因此其基础版本取自最高的发布标签（忽略预发布版本）：在发布热修复 `v3.1.1` 后，下一个 nightly 仍为 `v3.2.0-nightly.*`。它必须使用更高的版本号，因为预发布版本的排序低于其对应的正式版本：`v3.0.0-nightly.<date>`会排在 `v3.0.0` 之前。

**提升次版本号并不代表功能里程碑。**`v3.4.0`表示“`v3.3.0`之后的每周发布版”。要了解具体变更，请阅读发布说明，而不要只看版本号。

策略由版本号的**形式**决定：带`-nightly.`后缀的是每夜版；带任何其他后缀（`-alpha`、`-beta`、`-rc`）的是 GA 之前的旧式构建；不带后缀的是发布版。

### 版本同步一致

CLI、Go 模块和 npm 包`@wailsio/runtime`使用取自`v3/internal/version/version.txt`的**相同版本**（对于标签`v3.2.0-nightly.20261102`，npm 上为`3.2.0-nightly.20261102`）。使用`wails3 init`创建的项目会将`@wailsio/runtime`固定为 CLI 自身的版本；CLI 的开发构建则回退到`latest`。

## 完整示例

| 日期 | 事件 | 标签 | `go install ...@latest`解析为 |
| --- | --- | --- | --- |
| 10月25日（星期日） | GA，手动发布 | `v3.0.0` | `v3.0.0` |
| 10月26日（星期一）至10月31日（星期六） | 每夜版 | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| 11月1日（星期日）20:00（悉尼时间） | 每周版 | `v3.1.0` | `v3.1.0` |
| 11月2日（星期一） | 每夜版（预览下一个每周版） | `v3.2.0-nightly.20261102` | `v3.1.0` |
| 11月4日（星期三） | 热修复版，从发布分支手动发布 | `v3.1.1` | `v3.1.1` |
| 11月8日（星期日） | 每周版 | `v3.2.0` | `v3.2.0` |

每周版的排序始终高于同一版本号的所有每夜版，而每夜版绝不会变成`@latest`。

## 用户应运行什么

| 目标 | 命令 |
| --- | --- |
| 最新发布版（每周版或热修复版） | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`或`wails3 update cli` |
| **可复现（推荐）** | 固定到确切标签：在`go.mod`、CI 和文档中使用`v3.1.0` |
| 某个每夜版 | 其确切标签（`@v3.2.0-nightly.20261102`）、`wails3 update cli -nightly` |
| 任意类型的最新构建 | `wails3 update cli -pre` |
| 运行时 npm 包 | `@wailsio/runtime`（`latest` = 每周版/热修复版，`nightly` = 每夜版） |

其工作原理是：Go 会将`@latest`解析为**不带预发布后缀**的最高标签，并且仅在不存在发布版时才考虑预发布版。GitHub 的“pre-release”复选框对 Go 代理不可见，因此每夜版使用 semver 预发布后缀，而不依赖该复选框。生产环境中务必固定版本。

## 分支

| 分支 | 用途 |
| --- | --- |
| `master` | 用于开发和每夜版。 |
| `releases/**` | 用于一次性发布的短期分支：`v3.0.0` GA 和热修复版。所需 CI 与`master`相同。 |

每周版从`master`发布，但**并非从其分支顶端**发布：工作流会在上一个每周版之后的提交中，选取所有必需检查（`v3 Build Results`、`Run Go Tests`、`Bundled runtime matches source`）均已通过的最新提交，并准确地为该提交打标签。如果没有任何提交通过全部检查，则不发布任何版本。每夜版从`master`发布；HEAD 未通过全部检查时仅发出警告。

## 如何发布版本

发布权限仅授予项目负责人（参见治理文档）。贡献者绝不打标签或发布；他们只向`v3/UNRELEASED_CHANGELOG.md`添加条目。

所有操作都通过**Nightly + Weekly Release v3**（`.github/workflows/nightly-release-v3.yml`）运行：

- <strong>每夜版（自动）</strong>于 UTC 15:00运行。未发布的变更日志为空时跳过。
- 悉尼时间 20:00 执行**每周发布（自动）**。GitHub cron 仅支持 UTC，因此会安排两个 cron（UTC 时间 09:00 和 10:00），而守卫作业仅在悉尼时间为 20:00 时继续执行（可感知夏令时）。每周发布通过发布门禁后，会选择最新的绿色提交，递增次版本号，为该提交添加标签而不改动变更日志，然后仅更新 `master` 上的版本字段（`version.txt`和运行时软件包版本），因此会保留所选提交之后合并的依赖项和锁文件更新。npm 从标签发布，而不是从 `master` 发布。
- <strong>手动。</strong>使用`channel`、`dry_run`（默认启用）以及可选的`release_version`运行工作流。试运行会输出版本、类型、GitHub 标志、npm dist-tag 和说明，而不会改动 GitHub 或 npm。
- <strong>GA 或热修复。</strong>设置 `release_version`（`v3.0.0`、`v3.1.1`）。针对旧发布线的热修复还可以设置 `ref: releases/...`；其发布提交会被 cherry-pick 回 `master`。

本地试运行：`cd v3/tasks/release && go run release.go --dry-run --channel weekly`；`go run release.go --print-policy v3.2.0-nightly.20261102`会显示推导出的标志。

### 安全保障

- 自动生成的版本必须严格高于所有现有标签（按语义化版本顺序，预发布版本低于正式版本），否则工具会拒绝执行。
- 同一时间只能运行一个发布作业（使用共享并发组）。
- npm 发布由发布**标签**触发，并构建该标签所指向的提交。

## 发布门禁

只要存在以下任一情况，**每周发布**（绝不包括每夜发布）就会拒绝创建标签：

1. 带有`P0`或`release-blocker`标签的未关闭议题或拉取请求；
2. 自上次每周发布以来创建且带有`regression`标签的未关闭议题；
3. 严重级别为严重或高的未解决 Dependabot 警报、严重级别为严重或高的未解决代码扫描警报，以及任何未解决的密钥扫描警报，但**仅限已发布的代码**。如果警报仅位于`v3/examples/`、`v2/examples/`、`website/`、`docs/`、`.github/`或任意`examples/`目录下，则会被忽略（`v3/tasks/releasegate/main.go`中的`SHIPPED_ALERT_IGNORE`，**由负责人审核**）；
4. 仍处于`triage`状态的仓库安全公告。

**谁可以设置阻塞项。**只有当某标签最近一次 `labeled` 事件（之后没有 `unlabeled`）由 `wailsapp/maintainers` 团队成员或拥有 `maintain`/`admin` 权限的协作者执行时，该标签才有效。其他任何人应用的同一标签都会被忽略并记录到日志中，因此外部人员无法阻塞发布或解除发布阻塞。自动应用的 `high-priority` 标签特意**不**用于此目的。`taliesin-ai`等机器人账户只要拥有 `maintain`，即可通过允许列表检查；这是有意为之，分类代理可以在提供证据的情况下应用阻塞标签。

**失败时的行为。**门禁采用**故障关闭**策略：如果无法读取警报 API 或允许列表，每周发布就会停止。失败时不会创建任何标签；阻止项列表（含链接）会写入作业摘要；如果已设置仓库变量`WEEKLY_STATUS_ISSUE`，还会更新唯一的“每周发布状态”议题中的评论。

**负责人必须在 GitHub 上完成的设置**（不会自动创建任何内容）：创建`P0`、`release-blocker`和`regression`标签；创建“每周发布状态”议题，并将`WEEKLY_STATUS_ISSUE`设为其编号；授予`WAILS_REPO_TOKEN`对议题、拉取请求、管理（协作者权限）、Dependabot 警报、代码扫描警报、密钥扫描警报、仓库安全公告及组织成员的读取权限，并授予议题写入权限以发布状态评论。默认的`GITHUB_TOKEN`无法读取这些内容。

**绕过。**`override_reason`（工作流手动触发输入）会跳过门禁。绕过原因、执行者和所有被跳过的阻止项都会写入摘要和状态议题。只有发布负责人才能使用此选项。GitHub 无法将某个输入限制为仅一人可用，因此应在设有必需审核人的环境中运行发布作业（推荐），确保绕过操作必须获得明确批准。

## 相关内容

- [测试与持续集成](/contributing/testing-ci/)
- [维护者分诊](/contributing/maintainer-triage/)
