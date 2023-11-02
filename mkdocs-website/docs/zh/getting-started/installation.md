# 安装

要安装Wails CLI，请确保已经安装[Go 1.21+](https://go.dev/dl/)并运行以下命令：

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3-alpha
cd v3/cmd/wails3
go install
```

## 支持的平台

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Ubuntu 22.04 AMD64/ARM64（其他Linux可能也可以工作！）

## 依赖项

在安装之前，Wails有一些常见的依赖项需要安装：

=== "Go 1.21+"

    从[Go下载页面](https://go.dev/dl/)下载Go。

    确保按照官方的[Go安装指南](https://go.dev/doc/install)进行操作。您还需要确保您的`PATH`环境变量中包含`~/go/bin`目录的路径。重新启动终端并进行以下检查：

    - 检查Go是否已正确安装：`go version`
    - 检查`~/go/bin`是否在您的PATH变量中：`echo $PATH | grep go/bin`

=== "npm（可选）"

    虽然Wails不需要安装npm，但如果您想使用捆绑的模板，则需要安装npm。

    从[Node下载页面](https://nodejs.org/en/download/)下载最新的node安装程序。最好使用最新的版本，因为这是我们通常进行测试的版本。

    运行`npm --version`进行验证。

=== "Task（可选）"

    Wails CLI嵌入了一个名为[Task](https://taskfile.dev/#/installation)的任务运行器。这是可选的，但建议安装。如果您不想安装Task，可以使用`wails3 task`命令代替`task`。
    安装Task将给您最大的灵活性。

## 平台特定的依赖项

您还需要安装特定于平台的依赖项：

=== "Mac"

    Wails要求安装xcode命令行工具。可以通过运行以下命令来完成：

    ```
    xcode-select --install
    ```

=== "Windows"

    Wails要求安装[WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)。一些Windows安装可能已经安装了此软件。您可以使用`wails doctor`命令检查。

=== "Linux"

    Linux需要标准的`gcc`构建工具以及`libgtk3`和`libwebkit`。而不是列出不同发行版的大量命令，Wails可以尝试确定您特定发行版的安装命令。安装后运行<code>wails doctor</code>，将显示如何安装依赖项。如果您的发行版/软件包管理器不受支持，请在discord上告诉我们。

## 系统检查

运行`wails3 doctor`将检查您是否安装了正确的依赖项。如果没有安装，它将提供缺失的内容，并帮助您解决任何问题。

## 看起来缺少`wails3`命令？

如果系统报告缺少`wails3`命令，请检查以下内容：

- 确保您已正确按照Go安装指南进行操作。
- 检查`go/bin`目录是否在`PATH`环境变量中。
- 关闭/重新打开当前终端以使用新的`PATH`变量。