---
title: "测试与持续集成"
description: "Wails v3 如何通过单元测试、集成测试套件、竞态检测和 GitHub Actions CI 确保质量。"
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

健壮的桌面框架需要坚如磐石的测试。 Wails v3 采用<strong>分层策略</strong>：

| 层级 | 目标 | 工具 |
| --- | --- | --- |
| 单元测试 | 为独立函数提供快速反馈 | `go test ./...` |
| 生成器/CLI 测试 | 验证`wails3 generate bindings`和 CLI 基础设施 | `task test:generator`、`task test:cli` |
| 模板测试 | 确保发布的每个模板仍能成功构建 | `task test:templates` |
| 竞态检测 | 捕获运行时和桥接层中的数据竞态 | `go test -race ./...` |
| CI 矩阵 | 确保每个 PR 在不同操作系统上均可靠 | GitHub Actions |

本文档说明<strong>测试位于何处</strong>、**如何运行测试**，以及<strong>Taskfile 编排哪些任务</strong>。

> 旧版草稿中以`pkg/application/RACE.md`引用的竞态指南
>
> 目前位于`v3/TESTING.md`。

---

## 1. 目录约定

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

准则：

- **将单元测试放在代码旁边**（`foo.go` ↔ `foo_test.go`）。
- 对于`pkg/`包，在有助于保持 API 整洁时，请采用<strong>黑盒风格</strong>（`package application_test`）。
- 共享测试夹具放在使用它们的位置（此工作树中没有集中的`internal/testutil/`包，请改为按包接入辅助函数）。

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
- 当多个输入、边界情况或预期结果用于检验同一行为时，优先使用<strong>表驱动</strong>测试。为每个测试用例指定描述性名称。
- 必要时，通过构建标签（`foo_windows_test.go`、`foo_darwin_test.go`等）为平台特有行为提供桩实现。

### 覆盖率要求

新增和变更的逻辑应达到100% 的 Go 语句覆盖率。请测量 所更改包的覆盖率，不要依赖整个仓库的覆盖率百分比：

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

在常规测试环境中，有些路径确实无法合理测试，例如仅在特定平台出现的故障、 依赖硬件的行为，或无法安全触发的防御性回退逻辑。此类例外应严格限制在 必要范围内，并在 PR 描述中说明每条未覆盖的路径。

### 在本地运行

```bash
cd v3
go test ./... -cover
```

也可以通过 Taskfile 运行（请使用实际目标，不存在`task test`快捷目标）：

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. 集成测试

`v3/tests/`包含跨包集成测试工具。Taskfile 的 `test:example:*`和`test:examples:*`目标会在 darwin / windows / linux 上执行 构建和启动检查（包括 Linux 上基于 Docker 的 GTK3 / GTK4 矩阵）。

> 可运行的示例位于`v3/examples/`。测试目标会选择并构建
>
> 适用于宿主平台或 CI 矩阵的示例。

使用以下命令运行宿主平台冒烟测试套件：

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. 竞态检测

数据竞态对 GUI 运行时而言是致命的。

### 竞态指南

有关以下内容，请参阅`v3/TESTING.md`：

- 已知的无害竞态及其抑制理由
- 如何解读跨越 Cgo 边界的堆栈跟踪（Linux GTK + WebKit2GTK）

### 本地竞态测试套件

```
go test -race ./...
```

> `wails3 dev`没有`-race`标志，其 CLI 标志为`--config`、`--port`，
>
> 以及`-s`（启用 HTTPS）。若要在竞态检测器下检验运行时，
>
> 请使用`go build -race`构建测试应用并直接运行。

---

## 5. GitHub Actions 工作流

`.github/workflows/`下的实际工作流文件（已对照工作树验证）：

| 文件 | 用途 |
| --- | --- |
| `build-and-test-v3.yml` | 主要的 v3 构建和测试矩阵。通过`go-version: 1.25`使用`actions/setup-go@v5`。依次运行`task runtime:check`、`task runtime:test`、`task runtime:build`、`task test:examples`（GTK4 路径还会运行`BUILD_TAGS=gtk4 task test:examples`）、`task generator:test:check`、`task install`，最后运行`wails3 build`进行冒烟检查。Linux 作业会安装`libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk`，并在`dbus-run-session -- xvfb-run`下运行测试套件。 |
| `cross-compile-test-v3.yml` | 交叉编译健全性检查 |
| `auto-changelog-v3.yml`、`changelog-v3.yml` | 变更日志自动化 |
| `nightly-release-v3.yml` | v3 每夜版发布产物 |
| `bump-webview2-v3.yml`、`release-webview2.yml` | WebView2 依赖项/发布管理 |
| `build-cross-image.yml` | 构建交叉编译器容器镜像 |
| `publish-npm.yml` | 将嵌入式`@wailsio/runtime` JS 运行时发布到 npm |
| `pr-master.yml` | 针对`master`分支的 PR 端检查 |
| `semgrep.yml` | Semgrep 静态分析 |
| `stale-issues.yml`、`issue-labeler.yml`、`file-labeler.yml`、`claude.yml`、`generate-sponsor-image.yml`、`sync-translated-documents.yml`、`upload-source-documents.yml`、`build-and-test.yml`、`weekly-release-v2.yml` | 仓库维护/v2 端流程 |

此工作树中<strong>没有</strong>`qodana.yaml`，也<strong>没有</strong>`runtime.yml`——本页的旧草稿曾提及两者，但只有`semgrep.yml`负责静态分析，而运行时 JS 包通过`publish-npm.yml`发布。

CI 步骤与上述 Taskfile 目标（`task test:cli`、`task test:generator`、`task test:templates`、`task test:examples`……）相对应，因此你可以在本地逐项复现 CI。`build-and-test-v3.yml`中的冒烟`wails3 build`步骤在调用时<strong>不带任何额外标志</strong>——`wails3 build`没有`-skip-package`标志。

---

## 6. 本地 CI 对等验证

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
| <strong>`webview_window_darwin.go`</strong>中的竞态 | 在主线程之外修改窗口状态 | 通过`application.InvokeAsync`/`Invoke`封送调用，使其在主线程上运行 |
| **Linux 测试在无头 CI 中挂起** | GTK 需要显示服务器 | 在`xvfb-run`下运行，例如`xvfb-run task test:examples:linux` |
| **模板构建失败** | 前端锁文件已过期 | 针对一个干净目录重新运行`wails3 init`，以刷新模板 |
| **Coverpkg 错误** | 集成测试导入了`main` | 改用构建标签`//go:build integration`，并以此限制该导入 |

---

## 8. 添加新测试

1. **单元测试**——创建`*_test.go`，运行`go test ./...`
2. **生成器/CLI**——扩展`internal/generator/testcases/`或`internal/commands/*_test.go`下的测试用例，然后重新运行`task test:generator`/`task test:cli`
3. **模板/示例**——确保随附模板仍可使用`task test:templates`构建

---

## 9. 关键文件一览

| 内容 | 路径 |
| --- | --- |
| 生成器往返测试 | `internal/generator/generate_test.go` |
| 资源构建测试 | `internal/commands/build-assets_test.go` |
| 竞态/Cgo 指南 | `v3/TESTING.md` |
| Taskfile 测试目标 | `v3/Taskfile.yaml` |
| 事件常量生成器 | `v3/tasks/events/generate.go` |
| CI 工作流 | `.github/workflows/build-and-test-v3.yml`（通过`actions/setup-go@v5`配置使用 Go 1.25） |
| 静态分析 | `.github/workflows/semgrep.yml` |
| 运行时 npm 发布 | `.github/workflows/publish-npm.yml` |

---

在 Wails v3 中，质量并非事后才考虑的问题。有了单元测试、生成器和模板测试套件、竞态检测以及跨平台 CI 矩阵，你可以放心地贡献代码，因为你的更改能够在我们支持的每个操作系统上顺利通过测试。祝测试愉快！
