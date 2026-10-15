---
title: "绑定系统"
description: "Wails v3 如何让 Go 和 JavaScript 无需任何样板代码即可相互调用"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> “绑定”是一个<strong>类型安全的契约</strong>，让你可以编写：

```go
msg, err := chatService.Send("Hello")
```

在 Go 中<em>和</em>

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

在 TypeScript 中<strong>，而无需手动编写任何 IPC 粘合代码</strong>。 本文详细说明这一过程是<em>如何</em>实现的：从构建时的<strong>静态分析</strong>，到<strong>代码生成</strong>，再到通过 WebView 传输字节的<strong>运行时桥接层</strong>。

> 有关以下内容，请参阅[`contributing/architecture/bindings`](/contributing/architecture/bindings/)：
>
> 生成器流水线的权威深入解析——本页则是
>
> 面向贡献者的概述。

---

## 1. 30 秒概览

| 阶段 | 组件 | 输出 |
| --- | --- | --- |
| **收集/分析** | `internal/generator/collect/`、`internal/generator/analyse.go` | 已导出的 Go 服务、方法、参数、返回类型和模型的内存模型 |
| **生成** | `internal/generator/render/templates/*.tmpl`（`service.{js,ts}.tmpl`、`models.{js,ts}.tmpl`、`index.tmpl`、`eventcreate.js.tmpl`、`eventdata.d.ts.tmpl`、`newline.tmpl`） | `frontend/bindings/<full Go import path>/...`下每个服务对应的 ES 模块 |
| **运行时** | `pkg/application/messageprocessor*.go`以及`internal/runtime/desktop/@wailsio/runtime/src/`下的嵌入式 JS 运行时（`calls.ts`、`events.ts`……） | 通过 WebView 原生桥接层传输的调用/事件消息 |

此流程由`wails3 generate bindings`命令协调；该命令针对一组 Go 包驱动`generator.Generate`（定义于`internal/generator/generate.go`）。

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. 静态分析

### 入口点

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

收集器遍历每个已加载的包，并记录：

- `collect.ServiceInfo`——每个已导出并绑定的 Go 结构体对应一个。
- `collect.ServiceMethodInfo`/`collect.MethodInfo`——每个方法的签名信息（名称、参数、结果、错误位置、接收者、文档）。
- `collect.ModelInfo`/`collect.StructInfo`——生成为 TS/JS 模型。
- `//wails:inject`、`//wails:include`、`//wails:internal`、`//wails:ignore`、`//wails:id <hex>`等指令注释（参见`internal/generator/collect/directive.go`）。

不受支持的类型会导致生成器报错，使错误在构建时而非运行时暴露。

### 模型标识符

运行时调用封装使用方法完全限定名称（`pkg.Struct.Method`）的<strong>确定性 FNV-1a 哈希值</strong>来标识该方法。在生成的绑定中，它显示为`$Call.ByID(<numeric-id>, …)`；使用`-names`运行生成过程时，则显示为`$Call.ByName("pkg.Struct.Method", …)`。

---

## 3. 代码生成

### 模板

`internal/generator/render/templates/`：

| 模板 | 用途 |
| --- | --- |
| `service.js.tmpl` | 每个已绑定服务对应一个 JS 模块 |
| `service.ts.tmpl` | 配套 TypeScript 文件（使用`-ts`选项时生成） |
| `models.js.tmpl` | 模型类输出（每个包） |
| `models.ts.tmpl` | 模型`.d.ts`输出（每个包） |
| `index.tmpl` | 每个包的`index.{js,ts}`桶式重导出 |
| `eventcreate.js.tmpl`/`eventdata.d.ts.tmpl` | 事件构造函数/载荷类型定义 |
| `newline.tmpl` | 末尾换行符规范化器 |

输出位于`frontend/bindings/<full Go import path>/...`下——例如，在`github.com/you/yourapp/services/chat`中定义的服务会生成到`frontend/bindings/github.com/you/yourapp/services/chat/`下。v3 中没有`frontend/src/wailsjs/`目录。

### JavaScript 输出

生成的绑定是 ES 模块，它们从`/wails/runtime.js`导入运行时辅助函数：

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

使用`-names`运行生成过程时，改为生成`$Call.ByName("pkg.Struct.Method", ...)`——始终是<strong>完全限定名称</strong>，绝不会只有`"Method"`。

生成的模型类采用`$$source`构造函数模式，其中包含每个字段的`if (!("X" in $$source))`默认值、带引号的字段名，以及一个对字符串输入运行`JSON.parse`的`static createFrom(...)`。

### 类型映射要点

已根据`internal/generator/render/`验证：

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V`（非字符串`K`） | `{ [_ in K]?: V }`（不是`Map<K, V>`，也不是`Record<K, V>`） |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string`（JSON ISO 8601） |
| `error`（返回值位置） | 被拒绝的 Promise |

### 反射说明

`pkg/application/bindings.go`是<strong>手写的</strong>，并使用`reflect`根据`BoundMethod`注册表进行方法分派。不要过于字面地理解早期所称的“运行时零反射”——生成器避免使用反射，但运行时分派器会使用反射。

---

## 4. 运行时调用协议

### JavaScript 端

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

运行时辅助程序位于`internal/runtime/desktop/@wailsio/runtime/src/calls.ts`（调用分派）、`events.ts`（事件）及相关文件中——此工作树中没有`invoke.ts`或`errors.ts`。传输消息的具体封装格式由 JS 端的`calls.ts`编码，并由 Go 端的`pkg/application/messageprocessor_call.go`解码；调试桥接时，请结合查看这两个文件。

### Go 端

1. `pkg/application/messageprocessor_call.go`接收调用消息。
2. 在`pkg/application/bindings.go`中按 ID 或名称查找已绑定的方法（由`reflect`驱动）。
3. 调用已绑定的方法，并将`{result, error}`序列化后返回给 JS。

### 错误映射

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise`以结果兑现 |
| `error != nil` | `Promise`以`Error`拒绝，其`message`为 Go 错误字符串 |

---

## 5. 从 Go 调用 JavaScript

绑定生成器是单向的（将 Go 方法公开给 JS）。若要从 Go 与 JS 通信，请使用事件总线，或在窗口中运行 JS：

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

在 JS 端，使用`/wails/runtime.js`中的`Events.On(name, cb)`订阅。

---

## 6. 扩展与故障排除

### 不支持的类型错误

```
error: field "Client" uses unsupported type: chan struct{}
```

→ 将通道封装在方法 API 后面，或使用`//wails:internal`标记该字段，使生成器跳过它。

### 过期的绑定

每次运行`wails3 generate bindings`、`wails3 dev`或`wails3 build`时，都会覆盖生成的输出。如果 IDE 的 IntelliSense 显示过期的存根，请删除`frontend/bindings/`并重新运行生成器。`-clean`标志（当前构建中的默认值为`true`）会在每次运行前清空绑定目录。

### 性能提示

- 避免通过桥接流式传输大型字节切片，改由资源服务器提供这些数据。
- 对延迟敏感时，将多个快速调用批量合并为一次方法调用。
- 对于较小的参数结构体，优先使用值接收者以减少内存分配。

---

## 7. 关键文件一览

| 关注点 | 文件 |
| --- | --- |
| 生成器编排 | `internal/generator/generate.go` |
| 语义检查 | `internal/generator/analyse.go` |
| 收集（服务、方法、模型） | `internal/generator/collect/{service,method,model,struct,package}.go` |
| 渲染模板 | `internal/generator/render/templates/*.tmpl` |
| 生成的绑定所在位置 | `frontend/bindings/<full Go import path>/...` |
| Go 端分派器 | `pkg/application/bindings.go`、`messageprocessor_call.go` |
| JS 运行时 | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

跟踪桥接错误时，请将这份速查表放在手边。

---

## 8. 总结

1. <strong>收集器</strong>扫描 Go 代码 → 内存中的语义模型。
2. <strong>模板</strong>为每个服务生成 ES 模块，并为每个包生成模型和索引文件。
3. <strong>消息处理器</strong>通过绑定注册表在 Go 端分派调用。
4. <strong>JS 运行时</strong>将这一切封装为支持取消的惯用 Promise。

整个过程无需你编写一行 IPC 样板代码。这就是 Wails v3 绑定系统。开始绑定吧！
