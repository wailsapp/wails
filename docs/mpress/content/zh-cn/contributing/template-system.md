---
title: "模板系统"
description: "Wails v3 如何搭建新项目、如何组织模板，以及如何构建自己的模板。"
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails 随附一套<strong>模板系统</strong>，可让`wails3 init`生成一个可直接运行的项目。内置模板特意只涵盖少数框架（Vanilla、React、Vue、Svelte）；其他任何框架都可以通过[自带前端](/guides/dev/frontend-frameworks/)或发布[自定义模板](/guides/advanced/custom-templates/)来使用。

本页涵盖：

1. 模板目录布局
2. CLI 如何选择和渲染模板
3. 逐步创建新模板
4. 更新或覆盖现有模板
5. 故障排除和最佳实践

---

## 1. 模板所在位置

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`**——通用样板内容（Taskfile、`build/`目录和共享基础设施），会合并到每个项目中。
- **`base/`**——每个模板都以此作为 Go 端的起点。注意：`base/`本身<strong>不</strong>包含`template.json`；该文件位于各框架专用模板中。
- **框架文件夹**——包含前端（`frontend/`）、框架配置，以及用于描述模板元数据的`template.json`。
- 文件夹名称与传给 CLI 的<strong>模板 ID</strong>（`wails3 init -t react`）一致。
- <strong>语言命名约定：</strong>TypeScript 是默认语言，使用不带后缀的名称（`react`）；如果存在 JavaScript 变体，则添加`-js`后缀（`react-js`）。内置模板通过`template.yaml`中的`typescript: true|false`显式声明其语言。社区模板可能仍使用旧版`-ts`后缀，系统会将其作为回退方案予以识别。

> 整个`internal/templates/`目录都会编译进 CLI 二进制文件
>
> 这是通过`//go:embed *`实现的，因此用户可以离线搭建项目。

---

## 2. `wails3 init`如何使用模板

调用链（不存在`cmd/wails3/init.go`文件——CLI 直接在 `cmd/wails3/main.go`中完成连接）：

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

不存在`Template.Load()`/`Template.CopyTo()`/`Template.Validate()` API——提取由`gosod`（`github.com/leaanthony/gosod`）针对嵌入的`fs.FS`执行。

### `wails3 init`标志

定义于`internal/flags/init.go`：

| 标志 | 用途 | 默认值 |
| --- | --- | --- |
| `-p` | 包名称 | `main` |
| `-t` | 内置模板名称、本地路径或 URL | `vanilla` |
| `-n` | 项目名称 | （空） |
| `-d` | 项目目录 | `.` |
| `-q` | 禁止向控制台输出 | false |
| `-l` | 列出模板 | false |
| `-skipgomodtidy` | 提取后跳过运行`go mod tidy` | false |
| `-git` | 要初始化的 Git 仓库 URL | （空） |
| `-mod` | Go 模块路径（如未设置，则从`-git`派生） | （空） |
| `-s` | 使用远程模板时跳过警告 | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | 嵌入所生成构建资产中的元数据 | 合理的默认值 |

**没有**`-list`这一长名称别名（仅有`-l`），也没有各模板专用的`--help`。

### 替换项

占位符是标准 Go 模板指令——开头的`.`是字段访问器的一部分：

| 占位符 | 示例 | 来源 |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n` 标志/目录名 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod` 标志，或从 `-git` 派生 |
| `{{.WailsVersion}}` | `v3.0.0-…` | 来自 `internal/version` 的编译时常量 |
| `{{.ProductName}}`、`{{.ProductDescription}}`、`{{.ProductVersion}}`、`{{.ProductCompany}}`、`{{.ProductCopyright}}`、`{{.ProductComments}}`、`{{.ProductIdentifier}}` | 构建时元数据 | 对应的 `-product*` 标志 |

如果需要新的占位符，请在 `internal/templates/templates.go` 的模板数据中添加字段，并在 `internal/flags/init.go` 中添加匹配的字段/标志（或通过 `internal/commands/init.go` 设置该字段）。

### 复制后钩子

`gosod` 完成模板提取后，CLI 会运行：

```
go mod tidy
```

除非传入 `-skipgomodtidy`。不存在 `task deps` 步骤。

---

## 3. 创建新模板

> 示例：添加 **Solid** 模板

### 3.1 文件夹和 ID

```
internal/templates/solid/
```

文件夹名称即模板 ID。请保持为 **kebab-case** 格式。

### 3.2 最小文件集

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

首先复制 `react`，然后删减文件。别忘了编写 `template.yaml`；对于 TypeScript 模板，请设置 `typescript: true`。`base/` 是唯一没有该文件的文件夹。

### 3.3 更新占位符

搜索示例字面值并将其替换为 Go 模板指令，例如：

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 接入

由于 `templates.go` 会在初始化时遍历嵌入式文件系统，因此通常只需在 `internal/templates/<id>/` 下添加新文件夹，无需手动调用注册函数。如果需要额外逻辑（自定义验证、复制后步骤），请将其添加到 `internal/templates/templates.go` 中的 `templates.Install`。

### 3.5 测试

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

请确保：

- 开发服务器在 `WAILS_VITE_PORT` 中公布的端口上启动
- 生成的绑定出现在 `frontend/bindings/...` 下
- 热重载正常工作

---

## 4. 修改现有模板

1. 编辑 `internal/templates/<id>/` 下的文件。
2. 重新构建 CLI（`cd v3 && go build -o ../wails3 ./cmd/wails3`）；`//go:embed *` 指令会包含新内容。
3. 升级`frontend/package.json`和`Taskfile.yml`中的<strong>依赖项版本</strong>。
4. 如果行为发生变化，请更新模板的 `template.json` 描述。

### 常见调整

| 任务 | 位置 |
| --- | --- |
| 更改开发服务器端口 | `frontend/vite.config.ts` — 读取 `WAILS_VITE_PORT` |
| 添加环境变量 | `build/Taskfile.yml` 或 `frontend/.env` |
| 替换 JS 包管理器 | 在 `build/Taskfile.yml` 中将 `npm` 替换为 `pnpm`/`bun` |

---

## 5. 模板编写技巧

- **保持前端通用**——避免引用 Wails 专用全局对象；运行时，`/wails/runtime.js` 由资源服务器提供。
- **不要包含编译产物**——从嵌入的目录中排除`node_modules`、`dist`、`.DS_Store`（或通过`.gitignore`忽略它们，确保它们绝不会被提交）。
- **记录先决条件**——在 `template.json` 或 `NEXTSTEPS.md` 中记录 Node 版本、额外的 CLI 工具等。
- **避免破坏性变更**——如果改动很大，请创建新的模板 ID，而不要修改现有模板。

---

## 6. 故障排除

| 症状 | 原因 | 解决方法 |
| --- | --- | --- |
| `unknown template name` | `-t` 中有拼写错误，或模板未嵌入 | 运行 `wails3 init -l` 列出可用模板 |
| 占位符未被替换 | 使用了 `{{ProjectName}}`，而不是 `{{.ProjectName}}` | 添加开头的 `.`（Go 模板字段访问语法） |
| 开发服务器打开空白页面 | Vite 配置未读取 `WAILS_VITE_PORT` | 检查你的 `vite.config.ts` |
| 前端在生产环境中构建失败 | 遗漏了 Vite 的 `base` 路径 | 在 `vite.config.ts` 中设置 `base: "./"` |

---

## 7. 关键源文件一览

| 文件 | 职责 |
| --- | --- |
| `internal/templates/templates.go` | 嵌入模板文件系统，并公开 `Install(options *flags.Init) error`、`GetDefaultTemplates()` 和 `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | 实际的模板内容 |
| `internal/commands/init.go` | CLI 衔接代码：选择模板、填充元数据并调用 `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — 将现有项目<em>导出</em>回模板的实用工具（便于更新） |
| `internal/flags/init.go` | `wails3 init` 的标志定义 |

---

## 8. 回顾

- 模板位于 **`internal/templates/`** 中，并通过 `//go:embed *` 内置到 CLI 中。
- `wails3 init -t <id>` 通过 `gosod` 提取模板并运行 `go mod tidy`（可通过 `-skipgomodtidy` 跳过）。
- 创建模板很简单：**创建文件夹**，添加所需文件和一个 `template.json`，然后使用 `{{.ProjectName}}` 风格的占位符。
- 该系统既<strong>可扩展</strong>又<strong>自包含</strong>，非常适合与团队或社区分享自定义技术栈。

祝你模板开发愉快！
