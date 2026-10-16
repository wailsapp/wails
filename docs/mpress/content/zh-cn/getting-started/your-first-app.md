---
title: "你的第一个应用程序"
description: "逐步创建你的第一个 Wails 桌面应用程序"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

本指南将介绍如何创建你的第一个 Wails v3 应用程序，涵盖项目设置、构建和开发工作流。

<br/>

<br/>

@steps
### 创建新项目
打开终端并运行以下命令，以创建新的 Wails 项目：

```bash
wails3 init -n myfirstapp
```

此命令会创建一个名为`myfirstapp`的新目录，其中包含所有必要的文件。

   <video src="/assets/wails_init.mp4" controls></video>

### 了解项目结构
进入`myfirstapp`目录。你会看到以下文件和文件夹：

@filetree
- build/           包含构建过程使用的文件
  - appicon.png  应用程序图标
  - config.yml   构建配置
  - Taskfile.yml Build tasks
  - darwin/      macOS 专用构建文件
    - Info.dev.plist Development configuration
    - Info.plist    生产环境配置
    - Taskfile.yml  macOS 构建任务
    - icons.icns    macOS 应用程序图标
  - linux/       Linux 专用构建文件
    - Taskfile.yml  Linux 构建任务
    - appimage/     AppImage 打包
      - build.sh  AppImage 构建脚本
    - nfpm/        NFPM 打包
      - nfpm.yaml Package configuration
      - scripts/  构建脚本
  - windows/     Windows 专用构建文件
    - Taskfile.yml        Windows 构建任务
    - icon.ico           Windows 应用程序图标
    - info.json          应用程序元数据
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS 安装程序文件
      - project.nsi                    NSIS 项目文件
      - wails_tools.nsh               NSIS 辅助脚本
- frontend/        前端应用程序文件
  - index.html   主 HTML 文件
  - main.js      主 JavaScript 文件
  - package.json NPM package configuration
  - public/      静态资源
  - Inter Font License.txt Font license
- .gitignore      Git 忽略文件
- README.md       项目文档
- Taskfile.yml    项目任务
- go.mod          Go 模块文件
- go.sum          Go 模块校验和
- greetservice.go Greeting service
- main.go         应用程序主代码
@end

花一点时间浏览这些文件并熟悉项目结构。

@note{type="info"}
虽然 Wails v3 默认使用[Task](https://taskfile.dev/)作为构建系统，但你也完全可以使用`make`或任何其他替代构建系统。

@end

### 构建应用程序
要构建应用程序，请执行：

```bash
wails3 build
```

此命令会编译应用程序的调试版本，并将其保存到新建的`bin`目录中。

@note{type="info"}
`wails3 build`是`wails3 task build`的简写形式，它会运行`Taskfile.yml`中的`build`任务。

@end

     <video src="/assets/wails_build.mp4" controls></video>

构建完成后，可以像运行任何普通应用程序一样运行它：

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

你会看到一个简单的用户界面，它是应用程序的起点。由于运行的是调试版本，你还会在控制台窗口中看到日志。这对调试很有帮助。

### 开发模式
还可以在开发模式下运行应用程序。在此模式下，你可以修改前端代码，并立即在运行中的应用程序内看到变化，无需重新构建整个应用程序。

1. 打开一个新的终端窗口。
2. 运行`wails3 dev`。应用程序将以调试模式编译并运行。
3. 使用你选择的编辑器打开`frontend/index.html`。
4. 编辑代码，将`Please enter your name below`改为`Please enter your name below!!!`。
5. 保存文件。

此更改会立即反映在应用程序中。

对后端代码的任何更改都会触发重新构建：

1. 打开`greetservice.go`。
2. 将包含`return "Hello " + name + "!"`的那一行改为`return "Hello there " + name + "!"`。
3. 保存文件。

应用程序将在几秒内完成更新。

     <video src="/assets/wails_dev.mp4" controls></video>

### 打包应用程序
应用程序准备好分发后，可以创建特定于平台的软件包：

@tabs{sync-key="platform"}
[Mac]
要创建`.app`捆绑包：

```bash
wails3 package
```

这会创建生产版本，并将其打包为`bin`目录中的`.app`捆绑包。

[Windows]
要创建 NSIS 安装程序：

```bash
wails3 package
```

这会创建生产版本，并将其打包为`bin`目录中的 NSIS 安装程序。

[Linux]
Wails 支持多种用于 Linux 分发的软件包格式：

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

有关打包选项和配置的更多详细信息，请参阅我们的[构建与打包指南](/guides/build/building/)。

### 设置版本控制和模块名称
创建项目时使用的是占位模块名称`changeme`。建议将其更新为与你的仓库 URL 一致：

1. 在 GitHub（或你首选的 Git 托管平台）上创建一个新仓库
2. 在项目目录中初始化 Git：
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. 设置远程仓库（请替换为你的仓库 URL）：
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. 更新`go.mod`中的模块名称，使其与你的仓库 URL 一致：
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. 推送代码：
  ```bash
  git push -u origin main
  ```


这样可以确保 Go 模块名称符合 Go 的模块命名约定，也更便于共享代码。

@note{type="tip" title="实用技巧"}
创建项目时使用`-git`标志，即可自动完成所有初始化步骤：

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

此标志支持多种 Git URL 格式：

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project`或`ssh://git@github.com/username/project`
- Git 协议：`git://github.com/username/project`
- 文件系统：`file:///path/to/project.git`

@end

@end

## 恭喜！

你刚刚创建、开发并打包了自己的第一个 Wails 应用程序。这只是你使用 Wails v3 实现更多成果的开始。

## 后续步骤

如果你刚开始使用 Wails，建议接下来阅读我们的教程，通过实践指南了解 Wails 的各项功能。第一个教程是[创建服务](/tutorials/01-creating-a-service/)。

如果你是更高级的用户，请参阅[构建和打包指南](/guides/build/building/)，详细了解如何使用 Wails。
