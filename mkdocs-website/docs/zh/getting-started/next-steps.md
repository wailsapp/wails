# 下一步

现在您已经安装了Wails，可以开始探索alpha版本。

最好的起点是Wails存储库中的`examples`目录。
其中包含了一些可以运行和玩耍的示例。

## 运行示例

要运行示例，您可以简单地使用：

```shell
go run .
```

在示例目录中。

示例的状态在[路线图](../roadmap.md)中说明。

## 创建新项目

要创建新项目，可以使用`wails3 init`命令。这将在当前目录中创建一个新项目。

Wails3默认使用[Task](https://taskfile.dev)作为其构建系统，尽管您可以使用自己的构建系统，或直接使用`go build`。Wails内置了任务构建系统，可以使用`wails3 task`运行。

如果查看`Taskfile.yaml`文件，您会看到有一些任务被定义。最重要的任务是`build`任务。这是在使用`wails3 build`时运行的任务。

任务文件可能不完整，并且可能会随时间变化而改变。

## 构建项目

要构建项目，可以使用`wails3 build`命令。这是`wails3 task build`的快捷方式。