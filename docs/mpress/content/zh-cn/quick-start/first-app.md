---
title: "你的第一个应用"
description: "在 10 分钟内构建一个可运行的 Wails 应用"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

我们将构建一个简单的问候应用，以演示 Wails 的核心概念：

- 由 Go 后端管理逻辑
- 前端调用 Go 函数
- 类型安全的绑定
- 开发期间热重载

<strong>完成所需时间：</strong>10 分钟

@note{type="tip" title="Windows 11 用户性能提示"}
建议使用[开发驱动器](https://learn.microsoft.com/en-us/windows/dev-drive/)存储项目。开发驱动器针对开发者工作负载进行了优化，与普通 NTFS 驱动器相比，可将构建速度和磁盘访问速度显著提升，最高可达 30%。

@end

## 创建项目

@steps
### 生成项目
```bash
wails3 init -n myapp
cd myapp
```

这将使用默认的 Vanilla + Vite 模板（HTML/CSS/TypeScript，使用 Vite 打包器）创建一个新项目。

@note{type="tip" title="其他模板"}
可根据偏好的框架尝试 `-t react`、`-t vue` 或 `-t svelte`。这些模板默认使用 TypeScript；如需纯 JavaScript，请使用 `-t vanilla-js` 或 `-t react-js`。 运行 `wails3 init -l` 可查看所有可用模板，也可以 [使用自己的前端框架](/guides/dev/frontend-frameworks/)。

@end

### 了解项目结构
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### 运行应用
```bash
wails3 dev
```

@note{type="info" title="首次运行"}
首次运行可能比预期耗时更长，因为需要安装前端依赖项、生成绑定等。后续运行会快得多。

@end

应用打开后会显示问候界面。输入你的姓名并单击“问候”——Go 后端会处理输入并返回问候语。

@end

## 工作原理

下面来了解实现这一功能的代码。

### Go 后端

打开 `greetservice.go`：

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**核心概念：**

1. **服务**——包含导出方法的 Go 结构体
2. **导出方法**——`Greet`以大写字母开头，因此前端可以调用它
3. **简单逻辑**——接收姓名并返回问候语
4. **类型安全**——输入和输出类型均有明确定义

@note{type="tip" title="了解服务和绑定"}
<strong>服务</strong>是独立的 Go 模块，用于向前端公开功能。它们就是普通的 Go 结构体，包含导出方法，并注册在应用配置的 `Services` 字段中。

<strong>绑定</strong>是自动生成的 TypeScript/JavaScript SDK，前端通过它调用这些服务。运行 `wails3 dev` 或 `wails3 build` 时，Wails 会分析已注册的服务，并在 `frontend/bindings/` 中生成类型安全的绑定。

可以将服务视为后端 API，将绑定视为与该 API 通信的客户端库。

@end

### 注册服务

打开 `main.go`，找到服务注册代码：

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

这会向 Wails 注册 `GreetService`，使其所有导出方法都可供前端调用。

### 前端

打开 `frontend/src/main.js`：

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**核心概念：**

1. **自动生成的绑定**——`GreetService`从生成的代码中导入
2. **类型安全调用**——方法名称和签名与 Go 代码一致
3. **默认异步**——所有 Go 调用都会返回 Promise
4. **错误处理**——来自 Go 的错误会在 try/catch 中被捕获

@note{type="info" title="绑定在哪里？"}
生成的绑定位于 `frontend/bindings/` 中。运行 `wails3 dev` 或 `wails3 build` 时会自动创建这些绑定。

**切勿手动编辑这些文件**——每次构建时都会重新生成它们。

@end

## 自定义应用

下面添加一项新功能，以了解整个工作流程。

### 添加“批量问候”功能

@steps
### 向 GreetService 添加方法
将以下内容添加到 `greetservice.go`：

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### 应用将自动重新构建
保存文件后，`wails3 dev`会自动重新构建 Go 代码并重启应用。

@note{type="info" title="自动重新构建"}
更改 Go 代码会触发自动重新构建和重启。更改前端代码会进行热重载，无需重启。

@end

### 在前端中使用该功能
将以下内容添加到 `frontend/src/main.js`：

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

打开浏览器控制台并调用 `greetMany()`——你将看到问候语数组。

@end

## 构建生产版本

准备分发应用时：

```bash
wails3 build
```

**此操作会：**

- 启用优化编译 Go 代码
- 为生产环境构建前端（压缩代码）
- 在 `bin/` 中创建原生可执行文件

@tabs{sync-key="os"}
[Windows]
**输出：**`bin/myapp.exe`

双击即可运行。无需安装依赖项（WebView2 是 Windows 的一部分）。

[macOS]
**输出：**`bin/myapp.app`

拖入“应用程序”文件夹，或双击运行。

[Linux]
**输出：**`bin/myapp`

使用`./bin/myapp`运行，或为应用程序启动器创建一个`.desktop`文件。

@end

@note{type="tip" title="跨平台构建"}
想为其他平台构建吗？请参阅[跨平台构建 →](/guides/build/cross-platform/)

@end

## 我们学到了什么

**项目结构**

- `main.go`用于 Go 后端
- `frontend/`用于 UI 代码
- `Taskfile.yml`用于构建任务

**服务**

- 创建包含导出方法的 Go 结构体
- 使用`application.NewService()`注册
- 方法自动可供前端使用

**绑定**

- 自动生成 TypeScript 定义
- 类型安全的函数调用
- 默认异步（Promise）

**开发工作流**

- `wails3 dev`用于热重载
- Go 代码发生更改时自动重新构建并重启
- 前端代码发生更改时立即热重载

---

<strong>有问题？</strong>加入[Discord](https://discord.gg/JDdSxwjhGf)并向社区提问。
