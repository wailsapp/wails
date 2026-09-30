---
title: "构建系统"
description: "了解 Wails 如何构建和打包你的应用程序"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## 统一构建系统

Wails 提供了一个<strong>统一构建系统</strong>，只需一条命令即可编译 Go 代码、打包前端资源、将所有内容嵌入单个可执行文件，并处理特定于平台的构建。

```bash
wails3 build
```

<strong>输出：</strong>嵌入所有内容的原生可执行文件。

## 构建流程概览

**[构建流程图占位符]**

## 构建阶段

### 1. 分析阶段

Wails 扫描你的 Go 代码以了解其中的服务：

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Wails 提取的内容：**

- 服务名称：`GreetService`
- 方法名称：`Greet`
- 参数类型：`string`
- 返回类型：`string`

<strong>用途：</strong>生成 TypeScript 绑定

### 2. 生成阶段

#### TypeScript 绑定

Wails 生成类型安全的绑定：

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**优势：**

- 完整的类型安全
- IDE 自动补全
- 编译时错误
- JSDoc 注释

#### 前端构建

运行你的前端打包工具（Vite、webpack 等）：

```bash
# Vite example
vite build --outDir dist
```

**执行的操作：**

- 编译 JavaScript/TypeScript
- 处理并压缩 CSS
- 优化资源
- 生成源映射（仅开发环境）
- 输出到`frontend/dist/`

### 3. 编译阶段

#### Go 编译

使用优化选项编译 Go 代码：

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**标志：**

- `-s`：移除符号表
- `-w`：移除 DWARF 调试信息
- 结果：二进制文件更小（缩小约30%）

**特定于平台的处理：**

- Windows：嵌入图标的`.exe`
- macOS：`.app`包结构
- Linux：ELF 二进制文件

#### 资源嵌入

前端资源会嵌入 Go 二进制文件：

```go
//go:embed frontend/dist
var assets embed.FS
```

<strong>结果：</strong>包含所有内容的单个可执行文件。

### 4. 输出

**单个原生二进制文件：**

- Windows：`myapp.exe`（约 15MB）
- macOS：`myapp.app`（约 15MB）
- Linux：`myapp`（约 15MB）

**无依赖项**（系统 WebView 除外）。

## 开发环境与生产环境

@tabs{sync-key="mode"}
[开发环境（wails3 dev）]
**针对速度进行了优化：**

```bash
wails3 dev
```

**执行的操作：**

1. 启动前端开发服务器（Vite 默认使用端口9245）
2. 编译 Go 时不启用优化
3. 启动应用程序并将其指向开发服务器
4. 启用热重载
5. 包含源映射

**特征：**

- **快速重新构建**（前端更改耗时&lt;1 秒）
- **不嵌入资源**（由开发服务器提供）
- 包含<strong>调试符号</strong>
- 启用<strong>源映射</strong>
- **详细日志**

<strong>文件大小：</strong>较大（包含调试符号时约为 50 MB）

[生产构建（wails3 build）]
**针对体积和性能进行了优化：**

```bash
wails3 build
```

**构建过程：**

1. 以生产模式构建前端（压缩）
2. 启用优化编译 Go 代码
3. 移除调试符号
4. 嵌入资源
5. 生成单个二进制文件

**特性：**

- **经过优化的代码**（已压缩并进行 tree-shaking）
- **资源已嵌入**（无外部文件）
- **已移除调试符号**
- **无 source map**
- **最少量日志**

<strong>文件大小：</strong>较小（约 15 MB）

@end

## 构建命令

### 基本构建

```bash
wails3 build
```

**输出：**`bin/<APP_NAME>`（Windows 上为`bin/<APP_NAME>.exe`）。`bin/`目录位于项目根目录。

`wails3 build`是`wails3 task build`的轻量封装。它仅转发`--tags`这一构建时标志，该标志会成为 Taskfile 的`EXTRA_TAGS`变量：

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build`没有`-platform`、`-o`、`-skipbindings`、`-clean`、`-debug`、`-devbuild`、`-icon`、`-ldflags`或`-package`标志。交叉编译、输出路径、图标和打包均通过项目的 Taskfile（`Taskfile.yml` + `build/config.yml`）控制。

### 跨平台构建和平台特定构建

平台构建以 Taskfile 任务的形式提供，分别位于`darwin:` / `windows:` / `linux:`命名空间下（在`build/Taskfile.<platform>.yml`中定义）。例如：

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

要查看当前项目中的所有可用任务，请运行：

```bash
wails3 task --list
```

### 图标和打包

从源 PNG 生成各平台的图标（`build/icons.icns`、`build/icon.ico`等）：

```bash
wails3 generate icons -input appicon.png
```

构建平台特定的安装程序或软件包：

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## 构建配置

### Taskfile.yml

Wails 3项目使用[Taskfile](https://taskfile.dev/)作为构建编排工具。根目录中的`Taskfile.yml`会包含`build/`中各平台的任务文件：

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

使用`wails3 task <name>`或`task <name>`运行任务：

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### 项目配置：`build/config.yml`

项目元数据（名称、标识符、版本、info-plist 值、NSIS 设置、`.desktop`字段、自定义协议等）存储在`build/config.yml`中。Taskfile 在生成图标、清单、安装程序等内容时会读取此文件。Wails 3中<strong>没有</strong>`build/build.json`文件。

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

运行`wails3 generate build-assets`（或`wails3 update build-assets`），根据此配置刷新平台特定的构建资源。

## 资源嵌入

### 工作原理

Wails 使用 Go 的`embed`包：

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**构建时：**

1. 前端构建至`frontend/dist/`
2. `//go:embed`指令包含文件
3. 将文件编译到二进制文件中
4. 二进制文件包含所有内容

**运行时：**

1. 应用启动
2. 从内存提供资源
3. 资源无需磁盘 I/O
4. 加载迅速

### 自定义资源

嵌入其他文件：

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## 构建优化

### 前端优化

**Vite（默认）：**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**结果：**

- JavaScript 已压缩（体积减少约70%）
- CSS 已压缩（体积减少约60%）
- 图像已优化
- 已应用 tree-shaking

### Go 优化

**编译器标志：**

```bash
-ldflags="-s -w"
```

- `-s`：移除符号表（大小减少约10%）
- `-w`：移除 DWARF 调试信息（大小减少约20%）

**其他优化：**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`：在构建时设置变量值
- 适用于版本号、构建日期

### 二进制文件压缩

**UPX（可选）：**

```bash
# After building
upx --best bin/myapp.exe
```

**结果：**

- 大小减少约50%
- 启动速度略慢（约 100ms）
- 不建议用于 macOS（存在代码签名问题）

## 特定于平台的构建

### Windows

**输出：**`myapp.exe`

**包含：**

- 应用程序图标
- 版本信息
- 清单（UAC 设置）

**图标：**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

随后，Windows 的`tool package`步骤会将生成的`.ico`嵌入可执行文件。

**清单：**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**输出：**`myapp.app`（应用程序包）

**结构：**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist：**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**通用二进制文件：**

macOS Taskfile 提供`darwin:build:universal`（以及`darwin:package:universal`）任务，用于构建两种架构并通过`wails3 tool lipo`将其合并：

```bash
wails3 task darwin:build:universal
```

### Linux

**输出：**`myapp`（ELF 二进制文件）

**依赖项：**

- GTK3
- WebKitGTK

**桌面文件：**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**安装：**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## 构建性能

### 典型构建时间

| 阶段 | 耗时 | 备注 |
| --- | --- | --- |
| 分析 | &lt;1s | 扫描 Go 代码 |
| 生成绑定 | &lt;1s | 生成 TypeScript 代码 |
| 前端构建 | 5-30s | 取决于项目大小 |
| Go 编译 | 2-10s | 取决于代码规模 |
| 资源嵌入 | &lt;1s | 嵌入前端资源 |
| **总计** | **10-45s** | 首次构建 |
| **增量构建** | **5-15s** | 后续构建 |

### 加快构建速度

**1. 使用构建缓存：**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. 仅运行所需任务：**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. 并行构建（多台计算机/CI）：**

在 v3 中，Linux、Windows 和 macOS 之间的交叉编译通常在 Docker `wails-cross`容器中进行，或分别在各平台的专用运行器上进行——`wails3 build`本身以宿主操作系统为目标。有关受支持的工作流，请参阅[跨平台构建](/guides/build/cross-platform/)。

**4. 使用更快的工具：**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## 故障排除

### 构建失败

**症状：**`wails3 build`退出并报错

**常见原因：**

1. **Go 编译错误**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **前端构建错误**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **缺少依赖项**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### 二进制文件过大

<strong>症状：</strong>二进制文件大于 50 MB

**解决方案：**

1. **移除调试符号**（随附的 Taskfile 已将`-ldflags="-s -w"`传递给`go build`）。

2. **检查嵌入的资源**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **使用 UPX 压缩**
  ```bash
  upx --best bin/myapp.exe
  ```


### 构建缓慢

<strong>症状：</strong>构建耗时超过1分钟

**解决方案：**

1. **使用构建缓存**
  - Go 会自动使用缓存
  - 前端缓存（Vite）会自动启用


2. **只运行所需的任务**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **优化前端构建**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## 最佳实践

### ✅ 应该做

- **开发期间使用`wails3 dev`**——快速迭代
- **发布时使用`wails3 build`**——生成优化后的产物
- **为构建设置版本**——使用`-ldflags`嵌入版本信息
- **在目标平台上测试构建产物**——交叉编译并非万无一失
- **保持前端构建快速**——优化打包工具配置
- **使用构建缓存**——加快后续构建

### ❌ 不应该做

- **不要提交`build/`目录**——将其添加到`.gitignore`
- **不要跳过构建测试**——发布前务必测试
- **不要嵌入不必要的资源**——减小二进制文件大小
- **不要在生产环境中使用调试构建**——请使用优化构建
- **不要忘记代码签名**——分发时必须进行代码签名

## 后续步骤

**构建应用程序**——构建和打包的详细指南 [了解更多 →](/guides/build/building/)

**跨平台构建**——在一台计算机上为所有平台构建 [了解更多 →](/guides/build/cross-platform/)

**创建安装程序**——为最终用户创建安装程序 [了解更多 →](/guides/installers/)

---

<strong>对构建有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[构建示例](https://github.com/wailsapp/wails/tree/master/v3/examples/build)。
