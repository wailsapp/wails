---
title: "你的第一个应用程序"
description: "逐步创建你的第一个 Wails 桌面应用程序"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

本指南介绍如何创建你的第一个 Wails v3应用程序，涵盖项目设置、构建和开发工作流。

<br/>

<br/>

@steps
### 创建新项目
打开终端并运行以下命令，创建一个新的 Wails 项目：

```bash
wails3 init -n myfirstapp -t vanilla-js
```

本教程明确选择 JavaScript Vanilla 模板。基础的`vanilla`模板使用 TypeScript。该命令会创建`myfirstapp`，其中包含所有必要文件。

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
  - src/         前端源代码
    - main.js    主 JavaScript 文件
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

花些时间浏览这些文件并熟悉项目结构。

@note{type="info"}
尽管 Wails v3默认使用[Task](https://taskfile.dev/)作为构建系统，但你仍然可以使用`make`或其他任何替代构建系统。

@end

### 构建应用程序
要构建应用程序，请执行：

```bash
wails3 build
```

使用生成的 Taskfile 时，此命令会在`bin`中创建生产构建。使用`wails3 dev`可进入开发模式，也可以使用`wails3 build DEV=true`明确执行开发构建。

@note{type="info"}
`wails3 build`是`wails3 task build`的简写，会运行`Taskfile.yml`中的`build`任务。

@end

     <video src="/assets/wails_build.mp4" controls></video>

构建完成后，你可以像运行普通应用程序一样运行它：

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

你会看到初始界面。生产构建不保证开发者检查器或控制台的行为；调试时请使用开发模式。

### 开发模式
你也可以在开发模式下运行应用程序。在此模式下，你可以修改前端代码，并立即在运行中的应用程序里看到变化，无需重新构建整个应用程序。

1. 打开一个新的终端窗口。
2. 运行`wails3 dev`。应用程序将以调试模式编译并运行。
3. 使用你选择的编辑器打开`frontend/index.html`。
4. 将姓名输入框的`placeholder="Your name"`更改为`placeholder="Your full name"`，然后保存文件。
5. 保存文件。

此更改会立即反映在应用程序中。

对后端代码的任何更改都会触发重新构建：

1. 打开`greetservice.go`。
2. 将包含`return "Hello " + name + "!"`的那一行更改为`return "Hello there " + name + "!"`。
3. 保存文件。

应用程序将在几秒内更新。

     <video src="/assets/wails_dev.mp4" controls></video>

### 打包应用程序
当应用程序准备好分发时，你可以创建特定于平台的软件包：

@tabs{sync-key="platform"}
[Mac]
要创建`.app`捆绑包：

```bash
wails3 package
```

这会创建生产构建，并将其打包为`.app`捆绑包，存放在`bin`目录中。

[Windows]
要创建 NSIS 安装程序：

```bash
wails3 package
```

这会创建生产构建，并将其打包为 NSIS 安装程序，存放在`bin`目录中。

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
项目创建时使用占位模块名称`changeme`。建议将其更新为与你的仓库 URL 一致：

1. 在 GitHub（或你偏好的 Git 托管平台）上创建一个新仓库
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


这可确保 Go 模块名称符合 Go 的模块命名约定，并让代码更易于共享。

@note{type="tip" title="专业提示"}
创建项目时使用`-git`标志，即可自动完成所有初始化步骤：

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

此功能支持多种 Git URL 格式：

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project`或`ssh://git@github.com/username/project`
- Git 协议：`git://github.com/username/project`
- 文件系统：`file:///path/to/project.git`

@end

@end

## 恭喜！

你刚刚创建、开发并打包了自己的第一个 Wails 应用程序。这只是你使用 Wails v3 所能实现的众多成果的起点。

## 后续步骤

如果你刚开始使用 Wails，建议接下来阅读我们的教程，其中将通过实践指导你了解 Wails 的各项功能。第一篇教程是[创建服务](/tutorials/01-creating-a-service/)。

如果你是较有经验的用户，请参阅[构建与打包指南](/guides/build/building/)，详细了解如何使用 Wails。
