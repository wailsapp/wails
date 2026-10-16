---
title: "创建自定义模板"
description: "如何生成、自定义和托管自己的 Wails v3 项目模板"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails 附带了一组内置模板，但你也可以创建自己的模板并与社区共享。自定义模板就是一个 Git 仓库——公开托管后，任何人都可以用一条命令基于它搭建项目框架。

## 生成模板骨架

`wails3 generate template`命令会生成一个可供自定义的模板目录：

```bash
wails3 generate template -name MyTemplate
```

所有标志：

| 标志 | 说明 | 默认值 |
| --- | --- | --- |
| `-name` | 模板名称（必填） | — |
| `-author` | 作者名称 | — |
| `-description` | CLI 中显示的简短说明 | — |
| `-helpurl` | 此模板的文档 URL | — |
| `-version` | 初始版本 | `v0.0.1` |
| `-frontend` | 将现有前端目录复制到模板中 | — |
| `-dir` | 模板目录的写入位置 | 当前目录 |

使用所有标志的示例：

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

生成的目录结构如下：

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="阅读 NEXTSTEPS.md"}
生成的`NEXTSTEPS.md`包含模板各个部分的详细指南。请在自定义前阅读它，并在发布前将其删除——通过你的模板创建的项目中不得包含此文件。

@end

## 配置模板元数据

打开`template.yaml`，设置模板的元数据：

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

`wailsVersion`字段为<strong>必填项</strong>，且必须为`3`。顶部的`# yaml-language-server`注释可在 VS Code（安装[YAML 扩展](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)后）和 JetBrains IDE 中启用自动补全和内联验证——你可以保留或删除它；它不会影响运行时行为。

## 自定义模板

### 前端

`frontend/`目录会原样复制到通过你的模板创建的每个项目中。请将占位内容替换为实际的前端：

@tabs
[从头开始]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

按照提示操作，然后安装依赖项：

```bash
npm install
```

[使用现有项目]
生成模板时传入`-frontend`，即可一次性复制现有前端：

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

也可以稍后手动将其复制到`frontend/`目录中。

@end

### 构建任务

`Taskfile.tmpl.yml`定义了构建工作流。请更新`install:frontend:deps`和`build:frontend`任务，使其与你的前端工具链相匹配：

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go 应用程序

`main.go.tmpl`文件是应用程序的入口点。创建项目时，Wails 的模板引擎会处理此文件——`{{.ProductName}}`之类的模板变量会替换为用户提供的值。

若要将其作为真正的 Go 文件进行编辑（并获得 IDE 支持），请暂时将其重命名为`main.go`，完成更改后，在提交前将其改回`main.go.tmpl`。

#### 模板变量

任何`.tmpl`文件中都可以使用以下变量：

| 变量 | 说明 | 示例 |
| --- | --- | --- |
| `{{.ProjectName}}` | 用户提供的项目名称 | `"MyApp"` |
| `{{.BinaryName}}` | 二进制文件名 | `"myapp"` |
| `{{.ProductName}}` | 产品显示名称 | `"My Application"` |
| `{{.ProductDescription}}` | 产品说明 | `"An awesome application"` |
| `{{.ProductVersion}}` | 产品版本 | `"1.0.0"` |
| `{{.ProductCompany}}` | 公司/作者名称 | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | 版权字符串 | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | 其他产品备注 | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | 反向 DNS 产品标识符 | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Go 模块路径 | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | 创建项目时使用的 Wails 版本 | `"3.0.0"` |
| `{{.Typescript}}` | 如果模板名称以`-ts`结尾，则为`true` | `true` |
| `{{.Opn}}` | 字面量`{{`——在模板内进行转义 | `{{` |
| `{{.Cls}}` | 字面量`}}`——在模板内进行转义 | `}}` |

@note{type="tip"}
模板中的任何文件都可以是`.tmpl`文件，包括 HTML、JSON 和 YAML 文件。不带`.tmpl`后缀的文件将按原样复制。

@end

## 在本地测试模板

发布前，请从本地路径创建项目来测试模板：

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

然后验证项目能否正常运行：

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

检查以下各项：

- 前端热重载正常工作
- 更改 Go 代码后，应用会重新构建并重新启动
- `bin/`中的生产环境二进制文件能正常运行

## 发布到 GitHub

@steps
### **为模板创建一个公开的 GitHub 仓库**。仓库根目录必须包含 `template.yaml`。
### **删除 `NEXTSTEPS.md`**——此文件用于指导模板作者，不得出现在用户通过你的模板创建的项目中。
### **提交并推送**模板目录中的内容，将其作为仓库根目录的内容：
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **使用语义化版本控制为发布版本添加标签**：
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

用户现在可以通过你的模板创建项目：

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="第三方模板警告"}
用户安装远程模板时，Wails 会显示警告，说明该模板是第三方代码，并且 Wails 项目不对其内容承担任何责任。用户必须明确确认后，系统才会创建项目。

作为模板作者，你需要对模板中所有代码的安全性和正确性负责。

@end

## 最佳实践

- **编写清晰的`README.md`**——用户创建项目后会看到此文件。请说明如何运行、构建和自定义项目。
- **填写`helpurl`**——链接到你的仓库或专门的文档。用户会在 Wails CLI 模板列表中看到它。
- **固定前端依赖项的版本**，在`package.json`中固定版本，以免上游更新导致安装失败。
- **添加标签前进行测试**——向社区公布之前，先从已添加标签的发布版本创建一个全新项目。
- **保留`wailsVersion: 3`**——此字段用于告知 Wails 模板面向哪个主版本。请勿更改。
- **定期更新**——及时更新依赖项，并针对 Wails 的新版本进行测试。
