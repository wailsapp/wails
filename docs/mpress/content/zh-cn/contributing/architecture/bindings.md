---
title: "绑定系统"
description: "绑定系统如何收集和处理信息并生成 JavaScript/TypeScript 代码"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

本指南介绍 Wails 绑定系统的内部工作原理，帮助希望了解自动代码生成机制的开发者深入理解其实现。

## 架构概述

Wails 绑定系统由三个主要组件构成：

1. **收集**：分析 Go 代码，提取服务、模型及其他声明的相关信息
2. **配置**：管理绑定生成过程的设置和选项
3. **渲染**：根据收集的信息生成 JavaScript/TypeScript 代码

@filetree
- internal/generator/
  - collect/     # 包分析和信息提取
  - config/      # 配置结构和接口
  - render/      # 生成 JS/TS 代码
@end

## 收集过程

收集过程负责分析 Go 包，并提取服务、模型及其他声明的相关信息。此过程由`collect`包处理。

### 关键组件

- **收集器**：管理包信息并缓存已收集的数据
- **包**：表示正在分析的 Go 包，并存储收集到的服务、模型和指令
- **服务**：收集服务类型及其方法的相关信息
- **模型**：收集模型类型的详细信息，包括字段、值和类型参数
- **指令**：解析并解释 Go 源代码中的`//wails:`指令

### 收集流程

1. 收集器扫描项目中指定的 Go 包
2. 识别服务类型（其方法将向前端公开的结构体）
3. 收集每个服务的方法信息
4. 识别模型类型（用作服务方法参数或返回值的结构体）
5. 收集每个模型的字段和类型参数信息
6. 处理代码中发现的所有`//wails:`指令

## 渲染过程

渲染过程负责根据收集的信息生成 JavaScript/TypeScript 代码。此过程由`render`包处理。

### 关键组件

- **渲染器**：协调服务、模型和索引文件的渲染
- **模块**：表示单个生成的 JavaScript/TypeScript 模块
- **模板**：用于生成代码的文本模板

### 渲染流程

1. 渲染器为每个服务生成一个 JavaScript/TypeScript 文件，其中包含与服务方法一一对应的函数
2. 渲染器为每个模型生成一个与模型结构体对应的 JavaScript/TypeScript 类
3. 渲染器生成索引文件，重新导出所有服务和模型
4. 渲染器应用`//wails:inject`指令指定的所有自定义代码注入

## 类型映射

绑定系统最重要的方面之一，是如何将 Go 类型映射到 JavaScript/TypeScript 类型。以下是映射关系概要：

| Go 类型 | JavaScript 类型 | TypeScript 类型 |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V`（非字符串`K`） | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | 自定义类 |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | 不支持 | 不支持 |
| `chan` | 不支持 | 不支持 |

## 指令系统

绑定系统支持多种用于自定义生成代码的指令。这些指令以注释的形式添加到 Go 代码中。

### 可用指令

- `//wails:inject`：将自定义 JavaScript/TypeScript 代码注入生成的绑定中
- `//wails:include`：在生成绑定时包含其他文件
- `//wails:internal`：将类型或方法标记为内部使用，防止其导出到前端
- `//wails:ignore`：在生成绑定时完全忽略某个方法
- `//wails:id`：为方法指定自定义 ID，覆盖默认的基于哈希的 ID

### 指令处理

1. 在收集阶段，收集器会识别并解析 Go 代码中的指令
2. 指令与相应的声明（服务、方法、模型等）一起存储
3. 在渲染阶段，渲染器会应用指令来自定义生成的代码

## 高级功能

### 条件代码生成

绑定系统支持使用两字符条件前缀，对`include`和`inject`指令进行条件代码生成：

```
<language><style>:<content>
```

其中：

- `<language>`可以是：
  - `*` - JavaScript 和 TypeScript
  - `j` - 仅 JavaScript
  - `t` - 仅 TypeScript


- `<style>`可以是：
  - `*` - 类和接口
  - `c` - 仅类
  - `i` - 仅接口


例如：

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### 自定义方法 ID

默认情况下，方法通过基于哈希的 ID 来标识。不过，你可以使用`//wails:id`指令指定自定义 ID：

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

这有助于在重构代码时保持兼容性。

## 性能注意事项

绑定生成器以高效为设计目标，但仍需注意以下几点：

1. 首次运行时需要建立待扫描软件包的缓存，因此速度较慢
2. 后续运行会使用缓存的信息，因此速度更快
3. 生成器会处理项目中的所有软件包，对于大型项目，这可能比较耗时
4. 可以使用`-clean`标志在生成前清理输出目录

## 调试

如果在生成绑定时遇到问题，可以使用`-v`标志启用调试输出：

```bash
wails3 generate bindings -v
```

这将提供有关收集和渲染过程的详细信息，有助于确定问题的根源。
