---
title: "测试与持续集成"
description: "Wails v3 如何通过单元测试、集成测试套件、数据竞争检测和 GitHub Actions CI 来确保质量。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

可靠的桌面框架需要坚如磐石的测试。Wails v3采用**分层策略**：

| 层级 | 目标 | 工具 |
| --- | --- | --- |
| 单元测试 | 为独立函数提供快速反馈 | `go test ./...` |
| 生成器/CLI 测试 | 验证`wails3 generate bindings`和 CLI 的衔接机制 | `task test:generator`, `task test:cli` |
| 模板测试 | 确保发布的每个模板仍可构建 | `task test:templates` |
| 数据竞争检测 | 发现运行时和桥接层中的数据竞争 | `go test -race ./...` |
| CI 矩阵 | 确保每个 PR 在不同操作系统上均可靠 | GitHub Actions |

本文档说明**测试位于何处**、**如何运行测试**，以及**Taskfile 会编排哪些任务**。

`v3/TESTING.md`记录了跨平台示例的构建方式；它并不是数据竞争抑制项的目录。

---

## 1. 目录约定

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

准则：

- **将单元测试放在代码旁边**（`foo.go` ↔ `foo_test.go`）。
- 对于`pkg/`包，在有助于保持 API 整洁时，请采用**黑盒风格**（`package application_test`）。
- 共享固件应放在其使用位置（此工作树中没有集中的`internal/testutil/`包，请改为在各包中接入辅助程序）。

---

## 2. 单元测试

### 编写测试

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

建议：

- 使用[`stretchr/testify`](https://github.com/stretchr/testify)——它已包含在`go.mod`中。
- 当多组输入、边界情况或预期结果用于检验同一行为时，应优先采用**表驱动**测试。为每个测试用例指定一个描述性名称。
- 必要时，在构建标签（`foo_windows_test.go`、`foo_darwin_test.go`等）之后为平台差异提供桩实现。

### 覆盖率要求

新增和变更的逻辑应达到100% 的 Go 语句覆盖率。请测量所更改包的覆盖率，而不要依赖整个仓库的百分比：

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

有些路径在常规测试环境中确实无法合理测试，例如仅特定平台会出现的故障、依赖硬件的行为，或者无法安全触发的防御性回退逻辑。此类例外应严格限定范围，并在 PR 描述中说明每条未覆盖的路径。

### 在本地运行

```bash
cd v3
go test ./... -cover
```

也可以通过 Taskfile 运行（请使用实际目标——不存在`task test`快捷目标）：

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. 集成测试

`v3/tests/`包含跨包集成测试框架。Taskfile 的`test:example:*`和`test:examples:*`目标会在 darwin/windows/linux 上执行编译检查（包括 Linux 上基于 Docker 的 GTK3/GTK4 矩阵）。

> 可运行的示例位于`v3/examples/`。测试目标会选择并构建
>
> 适用于主机平台或 CI 矩阵的示例。

使用以下命令运行主机平台冒烟测试套件：

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 数据竞争检测

数据竞争对 GUI 运行时而言是致命问题。

### 调查数据竞争

检查发生冲突的两处访问，并在受影响的平台上复现。不要在没有证据的情况下将报告视为无害问题或予以抑制。`v3/TESTING.md`说明的是示例编译；这些构建成功并不能证明不存在数据竞争，也不能证明原生运行时行为正确。

### 本地数据竞争测试套件

```
go test -race ./...
```

> `wails3 dev`没有`-race`标志——它的 CLI 标志是`--config`、`--port`，
>
> 以及`-s`（启用 HTTPS）。要在数据竞争检测器下运行运行时，
>
> 请使用`go build -race`构建测试应用，然后直接运行它。

---

## 5. GitHub Actions 工作流

`.github/workflows/`下的实际工作流文件（已与工作树核对）：

| 文件 | 用途 |
| --- | --- |
| `build-and-test-v3.yml` | Go 1.26/1.27矩阵、运行时和生成器检查、模板/示例编译以及应用构建检查。Linux 默认使用 GTK4/WebKitGTK 6.0；旧版变体会设置`BUILD_TAGS=gtk3`。原生启动和交互必须另行检查。 |
| `cross-compile-test-v3.yml` | 交叉编译健全性检查 |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | 变更日志自动化 |
| `nightly-release-v3.yml` | 每夜构建的v3发布制品 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 依赖项/发布管理 |
| `build-cross-image.yml` | 构建交叉编译器容器镜像 |
| `publish-npm.yml` | 将嵌入式`@wailsio/runtime` JS 运行时发布到 npm |
| `pr-master.yml` | 针对`master`分支执行 PR 侧检查 |
| `semgrep.yml` | Semgrep 静态分析 |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | 仓库维护和v2工作流；旧的 Crowdin 同步/上传工作流已停用。 |

CI 步骤与上述 Taskfile 目标（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`等）一一对应，因此可以在本地逐项复现 CI。`build-and-test-v3.yml`中的冒烟`wails3 build`步骤在调用时**不附加任何额外标志**——`wails3 build`没有`-skip-package`标志。

---

## 6. 本地复现 CI

没有统一的`task ci`总括目标。请串联实际目标来复现 CI：

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. 排查失败的测试

| 症状 | 可能的原因 | 修复方法 |
| --- | --- | --- |
| **`webview_window_darwin.go`中的竞态** | 在主线程之外修改窗口状态 | 通过`application.InvokeAsync` / `Invoke`进行封送，使调用在主线程上运行 |
| **Linux 测试在无头 CI 中挂起** | GTK 需要显示服务器 | 在`xvfb-run`下运行，例如`xvfb-run task test:examples:linux` |
| **模板构建失败** | 前端锁文件已过期 | 对一个干净目录重新运行`wails3 init`以刷新模板 |
| **Coverpkg 错误** | 集成测试导入了`main` | 改用构建标签`//go:build integration`，并对该导入进行条件控制 |

---

## 8. 添加新测试

1. **单元测试** — 创建`*_test.go`，运行`go test ./...`
2. **生成器 / CLI** — 扩展`internal/generator/testcases/`或`internal/commands/*_test.go`下的测试用例，然后重新运行`task test:generator` / `task test:cli`
3. **模板 / 示例** — 确保随附模板仍可使用`task test:templates`构建

---

## 9. 关键文件索引

| 内容 | 路径 |
| --- | --- |
| 生成器往返测试 | `internal/generator/generate_test.go` |
| 资源构建测试 | `internal/commands/build-assets_test.go` |
| 跨平台示例构建 | `v3/TESTING.md` |
| Taskfile 测试目标 | `v3/Taskfile.yaml` |
| 事件常量生成器 | `v3/tasks/events/generate.go` |
| CI 工作流 | `.github/workflows/build-and-test-v3.yml`（通过 `actions/setup-go@v5` 运行 Go 1.26/1.27） |
| 静态分析 | `.github/workflows/semgrep.yml` |
| 运行时 npm 发布 | `.github/workflows/publish-npm.yml` |

---

在 Wails v3中，质量绝非事后考虑。借助单元测试、生成器/模板测试套件、竞态检测和跨平台 CI 矩阵，你可以放心地贡献代码，因为你知道自己的更改能在我们支持的每个操作系统上顺利通过测试。祝测试顺利！
