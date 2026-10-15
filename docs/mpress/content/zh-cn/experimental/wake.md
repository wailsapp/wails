---
title: "Wake"
description: "一个了解 Wails 的实验性构建运行器，可运行现有 Taskfile，并默认提供更快的增量构建、结构化输出和并行执行。"
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

@note{type="caution" title="实验性功能"}
Wake 需通过 `WAILS_USE_WAKE=true` 主动启用，<strong>并非</strong>默认运行器。 未设置该变量时，`wails3 build / package / sign / task` 的行为与以前完全相同。 功能覆盖范围和行为可能会随版本而变化。

@end

Wake 是用于 `wails3` 的<strong>实验性替代构建运行器</strong>。它读取项目中已有的 同一个 `Taskfile.yml`，支持相同的 task、dep、var、 template、include 和平台命名空间语法，并通过一个 了解 Wails 的执行器运行它，而不是使用通用的 [Task](https://taskfile.dev) 运行时。

目标不是取代 Task，而是提供一个专为 Wails 项目的实际构建方式而设计的运行器，其语义、输出和默认设置 与 `wails3` CLI 的其余部分保持一致。**如果你始终只使用 Wake， 则无需更改 Taskfile。**

## 存在的原因

Wake 和 Task 运行时都已编译到 `wails3` 中，两者都不需要 单独安装二进制文件。区别在于 Wake **理解该领域**。 通用运行器只会按照指定顺序执行 Taskfile 中列出的步骤。 Wake 知道 Wails 构建的实际<em>构成</em>：前端产物包 会嵌入二进制文件，二进制文件会打包成平台特定的 制品，同时还会生成图标和绑定；Wake 会利用这些知识 以通用运行器无法做到的方式优化构建。

- **它只执行构建实际需要的工作。** Wake 会自行跟踪每个步骤的实际输入和输出。对于 Go 构建，输入包括模块图以及依赖步骤的输出，因此当相关内容均未更改时，它会完全跳过编译器和链接器，而不是重新运行它们。只有当 Taskfile 预先明确列出要监视的具体文件时，通用运行器才能跳过某个步骤；Wake 则会根据自己已有的构建信息推导出来。对于无操作的重新构建，耗时大约为 **~20 ms（Wake），而 Task 约为 ~316 ms**。冷构建的实际耗时相同，因为主要时间都花在 `npm install`、Vite 和 Go 编译器上。

- **它知道哪些步骤可以同时运行。** Wake 了解哪些步骤相互独立，因此默认并行运行这些步骤，并在结果行中报告由此获得的加速幅度。如果同级步骤交错的输出会干扰问题调查，可使用 `WAKE_SERIAL=true` 退出并行执行。

- **由 wails3 控制的结构化输出。** Wake 通过 wails3 自带的报告器呈现输出：每个计划步骤占一行、实时显示状态、结束时按颜色分类列出各阶段耗时，并在故障面板中提供可点击的 `file:line` 链接。`NO_COLOR` 和非 TTY 环境（CI 日志）会妥善降级。

- **内置其中，因此可以随 Wails 一同发展。** Wake 是 `wails3` 的一部分，而不是第三方工具，因此可以直接在其上叠加新的构建能力，无需等待另一个项目实现这些功能。这也为原生运行跨平台脚本和工具创造了条件；目前，Taskfile 会调用 `wails3` 二进制文件来运行它们（每次调用都会生成一个进程）。将这项工作移至进程内可减少开销，也意味着未来还会有更多提速空间。

## 启用 Wake

Wake 完全由 `WAILS_USE_WAKE=true` 环境变量控制。 未设置该变量（或将其设置为 `true` 以外的任何值）时，每条 `wails3` 命令 都会与以前一样使用内嵌的 Task 运行时。

```bash
# Default: Task runtime, no Wake involvement
wails3 build

# Opt in: Wake drives build / package / sign / task <name>
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

此标志适用于 `wails3 build`、`wails3 package`、`wails3 sign` 和 `wails3 task <name>`。`wails3 dev` 目前<strong>不</strong>受影响，开发模式的监视器 仍使用自己的流水线。

@note{type="tip" title="启用 Wake 是安全的"}
如果 Wake 遇到尚未实现的 Taskfile 功能，它会在进程内将整个 运行过程交给内嵌的 Task 运行时，无需安装外部 `task` 二进制文件。最坏的结果也只是得到与未启用该标志时 完全相同的行为。

@end

## 分层本地覆盖

Wake 支持<strong>基础 Taskfile 加本地覆盖</strong>。在 `Taskfile.yml` 旁放置一个文件，其中的定义将优先采用：

| 文件 | 用途 | 优先级 |
| --- | --- | --- |
| `Taskfile.yml` | 基础文件，已提交 | 最低 |
| `Taskfile.override.yml` / `.yaml` | 已提交、适用于整个团队的覆盖文件 | 中等 |
| `Taskfile.local.yml` / `.yaml` | 个人文件，通常被 Git 忽略 | 最高 |

**合并语义（以本地定义为准）：**

- 具有<strong>相同名称</strong>的任务会覆盖基础任务。当覆盖定义提供列表字段（`cmds`、`deps`、`sources`、`generates`、`platforms`、`status`、`preconditions`、`aliases`）时，这些字段会<strong>替换</strong>基础定义中的对应字段；覆盖定义中省略的字段则保留基础定义中的值。
- `env` 和 `vars` 会<strong>按键合并</strong>，发生冲突时以覆盖定义为准。
- 如果某个任务<strong>仅</strong>存在于覆盖文件中，则会<strong>添加</strong>该任务。

例如，如果已提交的 `Taskfile.yml` 使用开发标志进行构建，而你的 计算机应始终执行生产构建：

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

现在，`build` 会运行你的生产命令，并且 `smoke` 可供使用，而无需 更改已提交的 Taskfile。

@note{type="note" title="信任模型"}
覆盖文件会被自动发现并应用，不会显示提示。这不会授予任何 新能力：Taskfile 本就可以运行任意 shell 命令，因此覆盖文件 无法执行任何通过编辑 `Taskfile.yml` 所不能执行的操作。已提交的 `Taskfile.override.*` 会出现在 PR 差异中；`Taskfile.local.*` 则在 你自己的计算机上创建。格式错误的覆盖文件会中止运行，而不会被 静默跳过。要获得确定性的 CI 构建，请设置 `WAILS_NO_OVERRIDES=true`， 以完全跳过覆盖文件发现。

@end

## 自动回退

如果 Wake 遇到尚未实现的 Taskfile 功能，它会将整个运行过程 交给内嵌的 Task 运行时。目前，以下情况会触发回退：

- Taskfile 级别的 `dotenv`
- 除 `interleaved` 以外的 `output` 模式
- 一个 `requires` 块
- `interval`（Taskfile 或任务级别）
- 除 `always` 以外的 `run` 模式
- 任务中的 `short`
- 任务中的`defer`

## 环境变量

| 变量 | 作用 |
| --- | --- |
| `WAILS_USE_WAKE` | 设为`true`时，为可路由的`wails3`动词启用 Wake；设为其他任何值时，使用 Task 运行时 |
| `WAILS_NO_OVERRIDES` | `true`跳过`Taskfile.local.*`/`.override.*`发现过程（实现确定性构建） |
| `WAKE_VERBOSE` | 实时流式输出子进程的 stdout/stderr，而不是捕获后仅在失败时显示 |
| `WAKE_SILENT` | 完全禁止任务输出 |
| `WAKE_SERIAL` | `true`禁用并行`deps:`扇出（默认启用并行） |
| `WAKE_FORCE` | `true`绕过所有缓存，以进行真正的全新重建 |
| `WAKE_DEBUG` | 记录解析器内部信息（DAG、依赖项、变量引用、执行路由） |
| `WAKE_NOTICE` | 设为`off`可隐藏每次运行时显示的“wake (experimental)”通知 |

构建缓存位于`.wake/cache.json`中（Task 使用`.task/`）。

## 反馈

Wake 是一项实验性功能，它未来的发展方向取决于您的反馈。如果您进行了尝试，我们很想知道它是否更快、更清晰，以及是否有任何功能出现故障——最有用的报告会说明您运行了什么、预期结果是什么，以及实际发生了什么。请在[Wake 反馈讨论](https://github.com/wailsapp/wails/discussions/5679)中告诉我们。
