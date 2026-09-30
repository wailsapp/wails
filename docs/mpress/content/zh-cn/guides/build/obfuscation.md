---
title: "混淆构建"
description: "使用 Garble 构建 Wails 应用，保护源代码免遭逆向工程"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) 是一款 Go 构建工具，它会替代 `go build`，对符号重命名、混淆常量，并从生成的二进制文件中移除调试信息。Wails v3 通过两个新命令提供对 Garble 的原生支持。

## 前提条件

- **Go 1.26.2 或更高版本**——Garble v0.16.0 要求使用该版本
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Garble 要求的最低 Go 版本会随发行版本而变化。如果使用较旧的 Go 工具链，请在安装前查看 [Garble 发行版本页面](https://github.com/burrowers/garble/releases)，找到与所用工具链匹配的版本。

@end

## 必需：为服务类型添加 JSON 标签

绑定服务方法返回或接受的任何结构体，其每个导出字段都必须具有显式 JSON 标签：

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble 会重命名结构体的导出字段，而 Wails 会通过 Garble 无法静态追踪的 `interface{}` 参数将这些结构体传递给 `json.Marshal`。缺少 JSON 标签的混淆构建可以成功编译，但在运行时，前端会收到被混淆的字段名或空字段名。请在运行混淆构建前添加标签。

@end

Wails 自身的类型——`Screen`、`Rect`、`Point`、`Size`、`EnvironmentInfo`、`OSInfo`、`Capabilities`——已经添加了标签。你只需为自己的类型添加标签。

## 使用混淆进行构建

@steps
### 生成稳定 ID 文件
每当添加、重命名或移除绑定服务方法时，都应运行此命令：

```bash
wails3 generate bindings -obfuscated
```

此操作会在主包目录中创建 `wails_obfuscated.gen.go`——请提交此文件。

### 使用 Garble 构建
```bash
wails3 build --obfuscated
```

使用经过混淆的绑定构建应用。

@end

## 向 Garble 传递额外标志

使用 `--garbleargs` 将选项直接转发给 `garble`：

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

有关受支持标志的完整列表，请参阅 [Garble 文档](https://github.com/burrowers/garble#flags)。

## 高级：将 ID 文件写入其他包

默认情况下，`wails_obfuscated.gen.go` 会写入 `main` 包所在的目录。如果项目将服务放在由 `main` 导入的子包中，可以使用 `-obfuscated-output` 将文件写入该子包：

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
目标包必须由 `main` 包直接或间接导入，以便其 `init()` 在启动时运行。如果该包不可达，稳定 ID 将永远不会注册，绑定调用也会失败（例如，运行时出现 `binding not found` 错误）。

@end

## 故障排除

### `garble: command not found`

Garble 未安装，或者 `$(go env GOPATH)/bin` 不在 `PATH` 中。

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 前端收到错误字段值或空字段值

服务返回类型缺少 `json:"..."` 标签。检查绑定方法返回的每个结构体，并为每个导出字段添加显式标签。

### 浏览器控制台中出现 `binding not found` 错误

稳定 ID 文件缺失或未被编译。请检查：

- 主包目录（或传递给 `-obfuscated-output` 的目录）中存在 `wails_obfuscated.gen.go`
- 你运行了 `wails3 build --obfuscated`，该命令会添加 `wails_obfuscated` 构建标签
- 如果使用了 `-obfuscated-output`，目标包已由 `main` 导入

### Windows Defender 将构建标记为病毒

在构建期间，经 Garble 混淆的 Go 二进制文件会被 Windows Defender 根据启发式规则标记，因为这些文件缺少调试符号，并且类似于加壳的可执行文件。构建会失败并显示：

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

将临时目录（Go 写入中间构建产物的位置）和项目目录添加到 Defender 的排除列表：

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

这些排除项仅适用于指定路径，不会在全局范围内禁用 Defender。

### 构建失败并显示 `unsupported Go version`

Garble v0.16.0 要求使用 Go 1.26.2 或更高版本。请升级 Go，或查看 [Garble 发行版本页面](https://github.com/burrowers/garble/releases)，找到与所用工具链兼容的版本。
